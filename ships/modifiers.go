package ships

import (
	"encoding/json"
	bson "go.mongodb.org/mongo-driver/v2/bson"
)

// Deterministic Combat System
//
// This combat system uses deterministic mechanics instead of RNG for predictable,
// strategic gameplay in hourly turn-based battles:
//
// 1. CRIT: Counter-based, not random chance
//    - CritPct = 0.33 → crit every 3rd attack (1/0.33)
//    - CritPct = 0.50 → crit every 2nd attack (1/0.50)
//    - Crit damage = base damage * 1.5 (+50%)
//
// 2. EVASION: Flat damage reduction, not dodge chance
//    - EvasionPct = 0.35 → 35% damage reduction on all incoming damage
//    - Capped at 75% reduction (EvasionPct = 0.75)
//    - Stacks additively from multiple sources (bio traits, formations, etc.)
//
// 3. FIRST STRIKE: Bonus on attack counter == 1
//    - FirstVolleyPct = 0.30 → +30% damage on first attack only
//    - Resets when battle ends or stack enters cooldown
//
// 4. SHIELDS: Asymptotic mitigation by attack type
//    - Each attack type (Laser/Nuclear/Antimatter) mitigated by corresponding shield
//    - Formula: damage / (1 + shieldValue * 0.15)
//    - Never reaches 100% mitigation (diminishing returns)
//    - Bio debuffs can reduce shields (even below 0, capped at 0 for calculations)
//
// 5. BIO DEBUFFS: Applied post-combat, affect next round
//    - Stack over multiple combat rounds
//    - Can reduce shields, add damage over time, etc.

const statEps = 1e-9

func fz(x float64) bool { return x > -statEps && x < statEps }

// DamageMods expresses percentage damage multipliers by damage type.
// Use values like +0.12 for +12% damage.
// These are additive across sources and applied multiplicatively at resolve time.
// Example: base 100 * (1 + 0.12 + 0.08) = 120 for Laser with +12% and +8%.
// Note: We expose Laser/Nuclear/Antimatter because your combat system has these three channels.
// If you add more types later, extend this struct in a backward-compatible manner.
type DamageMods struct {
	LaserPct      float64
	NuclearPct    float64
	AntimatterPct float64
}

// StatMods are soft modifiers that get applied to a ship type due to
// role modes, runes, runewords, or temporary abilities.
// Positive percentages are buffs, negative are debuffs unless noted.
// Deltas are integer offsets to base stats.
type StatMods struct {
	Damage            DamageMods // % to damage per type
	AttackIntervalPct float64    // % change to attack interval (lower is better)
	SpeedDelta        int        // +/- to base Speed
	VisibilityDelta   int        // +/- to base VisibilityRange
	AttackRangeDelta  int        // +/- to base AttackRange

	LaserShieldDelta      int // +/- to LaserShield
	NuclearShieldDelta    int // +/- to NuclearShield
	AntimatterShieldDelta int // +/- to AntimatterShield

	BucketHPPct         float64 // % change to per-bucket HP (affects survivability)
	OutOfCombatRegenPct float64 // % change to out-of-combat HP regen
	AtCombatRegenPct    float64 // % change to at-combat HP regen
	AbilityCooldownPct  float64 // % change to ability cooldowns (negative reduces CD)

	TransportCapacityPct float64 // % change to TransportCapacity

	// Travel-related
	WarpChargePct         float64 // % change to warp charge time (negative is faster)
	WarpScatterPct        float64 // % change to warp scatter (negative reduces scatter)
	InterdictionResistPct float64 // % chance-based resistance to interdiction effects

	// Combat quality of life
	StructureDamagePct float64 // % bonus vs structures/infrastructure
	SplashRadiusDelta  int     // + radius cells for splash
	AccuracyPct        float64 // % flat accuracy improvement
	CritPct            float64 // DETERMINISTIC: crit interval = 1/CritPct (e.g., 0.33 = every 3rd attack, 0.50 = every 2nd)
	FirstVolleyPct     float64 // % bonus damage on the first attack in combat (attack counter == 1)
	ShieldPiercePct    float64 // % of shields ignored (applied carefully)

	// Economy/logistics
	UpkeepPct           float64 // % change to upkeep
	ConstructionCostPct float64 // % change to build costs

	// Recon/detection (boolean capabilities are OR-composed)
	CloakDetect  bool    // can detect cloaked or mode-switch signals
	PingRangePct float64 // % change to Ping ability range

	// Formation-specific modifiers (applied in formation combat contexts)
	EvasionPct          float64 // DETERMINISTIC: flat % damage reduction (not dodge chance). Capped at 75% reduction.
	FormationSyncBonus  float64 // % bonus when position requirements are met
	PositionFlexibility float64 // % reduced penalty for suboptimal positions

	// Generic mods
	GlobalDefensePct float64 // % global damage reduction
	HPPct            float64 // % change to current HP
}

func ZeroMods() StatMods { return StatMods{} }

func (d DamageMods) IsZero() bool { return fz(d.LaserPct) && fz(d.NuclearPct) && fz(d.AntimatterPct) }

func (m StatMods) IsZero() bool {
	if !m.Damage.IsZero() { return false }
	if !fz(m.AttackIntervalPct) || m.SpeedDelta != 0 || m.VisibilityDelta != 0 || m.AttackRangeDelta != 0 { return false }
	if m.LaserShieldDelta != 0 || m.NuclearShieldDelta != 0 || m.AntimatterShieldDelta != 0 { return false }
	if !fz(m.BucketHPPct) || !fz(m.OutOfCombatRegenPct) || !fz(m.AtCombatRegenPct) || !fz(m.AbilityCooldownPct) { return false }
	if !fz(m.TransportCapacityPct) || !fz(m.WarpChargePct) || !fz(m.WarpScatterPct) || !fz(m.InterdictionResistPct) { return false }
	if !fz(m.StructureDamagePct) || m.SplashRadiusDelta != 0 || !fz(m.AccuracyPct) || !fz(m.CritPct) { return false }
	if !fz(m.FirstVolleyPct) || !fz(m.ShieldPiercePct) || !fz(m.UpkeepPct) || !fz(m.ConstructionCostPct) { return false }
	if m.CloakDetect || !fz(m.PingRangePct) || !fz(m.EvasionPct) || !fz(m.FormationSyncBonus) || !fz(m.PositionFlexibility) { return false }
	if !fz(m.GlobalDefensePct) || !fz(m.HPPct) { return false }
	return true
}

func CombineMods(a, b StatMods) StatMods {
	a.Damage.LaserPct += b.Damage.LaserPct
	a.Damage.NuclearPct += b.Damage.NuclearPct
	a.Damage.AntimatterPct += b.Damage.AntimatterPct

	a.AttackIntervalPct += b.AttackIntervalPct
	a.SpeedDelta += b.SpeedDelta
	a.VisibilityDelta += b.VisibilityDelta
	a.AttackRangeDelta += b.AttackRangeDelta

	a.LaserShieldDelta += b.LaserShieldDelta
	a.NuclearShieldDelta += b.NuclearShieldDelta
	a.AntimatterShieldDelta += b.AntimatterShieldDelta

	a.BucketHPPct += b.BucketHPPct
	a.OutOfCombatRegenPct += b.OutOfCombatRegenPct
	a.AtCombatRegenPct += b.AtCombatRegenPct
	a.AbilityCooldownPct += b.AbilityCooldownPct

	a.TransportCapacityPct += b.TransportCapacityPct

	a.WarpChargePct += b.WarpChargePct
	a.WarpScatterPct += b.WarpScatterPct
	a.InterdictionResistPct += b.InterdictionResistPct

	a.StructureDamagePct += b.StructureDamagePct
	a.SplashRadiusDelta += b.SplashRadiusDelta
	a.AccuracyPct += b.AccuracyPct
	a.CritPct += b.CritPct
	a.FirstVolleyPct += b.FirstVolleyPct
	a.ShieldPiercePct += b.ShieldPiercePct

	a.UpkeepPct += b.UpkeepPct
	a.ConstructionCostPct += b.ConstructionCostPct

	a.CloakDetect = a.CloakDetect || b.CloakDetect
	a.PingRangePct += b.PingRangePct

	a.EvasionPct += b.EvasionPct
	a.FormationSyncBonus += b.FormationSyncBonus
	a.PositionFlexibility += b.PositionFlexibility

	a.GlobalDefensePct += b.GlobalDefensePct
	a.HPPct += b.HPPct
	return a
}

func (m StatMods) MarshalJSON() ([]byte, error) {
	obj := make(map[string]any)
	if !m.Damage.IsZero() {
		dmg := make(map[string]any)
		if !fz(m.Damage.LaserPct) { dmg["LaserPct"] = m.Damage.LaserPct }
		if !fz(m.Damage.NuclearPct) { dmg["NuclearPct"] = m.Damage.NuclearPct }
		if !fz(m.Damage.AntimatterPct) { dmg["AntimatterPct"] = m.Damage.AntimatterPct }
		if len(dmg) > 0 { obj["Damage"] = dmg }
	}
	if !fz(m.AttackIntervalPct) { obj["AttackIntervalPct"] = m.AttackIntervalPct }
	if m.SpeedDelta != 0 { obj["SpeedDelta"] = m.SpeedDelta }
	if m.VisibilityDelta != 0 { obj["VisibilityDelta"] = m.VisibilityDelta }
	if m.AttackRangeDelta != 0 { obj["AttackRangeDelta"] = m.AttackRangeDelta }
	if m.LaserShieldDelta != 0 { obj["LaserShieldDelta"] = m.LaserShieldDelta }
	if m.NuclearShieldDelta != 0 { obj["NuclearShieldDelta"] = m.NuclearShieldDelta }
	if m.AntimatterShieldDelta != 0 { obj["AntimatterShieldDelta"] = m.AntimatterShieldDelta }
	if !fz(m.BucketHPPct) { obj["BucketHPPct"] = m.BucketHPPct }
	if !fz(m.OutOfCombatRegenPct) { obj["OutOfCombatRegenPct"] = m.OutOfCombatRegenPct }
	if !fz(m.AtCombatRegenPct) { obj["AtCombatRegenPct"] = m.AtCombatRegenPct }
	if !fz(m.AbilityCooldownPct) { obj["AbilityCooldownPct"] = m.AbilityCooldownPct }
	if !fz(m.TransportCapacityPct) { obj["TransportCapacityPct"] = m.TransportCapacityPct }
	if !fz(m.WarpChargePct) { obj["WarpChargePct"] = m.WarpChargePct }
	if !fz(m.WarpScatterPct) { obj["WarpScatterPct"] = m.WarpScatterPct }
	if !fz(m.InterdictionResistPct) { obj["InterdictionResistPct"] = m.InterdictionResistPct }
	if !fz(m.StructureDamagePct) { obj["StructureDamagePct"] = m.StructureDamagePct }
	if m.SplashRadiusDelta != 0 { obj["SplashRadiusDelta"] = m.SplashRadiusDelta }
	if !fz(m.AccuracyPct) { obj["AccuracyPct"] = m.AccuracyPct }
	if !fz(m.CritPct) { obj["CritPct"] = m.CritPct }
	if !fz(m.FirstVolleyPct) { obj["FirstVolleyPct"] = m.FirstVolleyPct }
	if !fz(m.ShieldPiercePct) { obj["ShieldPiercePct"] = m.ShieldPiercePct }
	if !fz(m.UpkeepPct) { obj["UpkeepPct"] = m.UpkeepPct }
	if !fz(m.ConstructionCostPct) { obj["ConstructionCostPct"] = m.ConstructionCostPct }
	if m.CloakDetect { obj["CloakDetect"] = true }
	if !fz(m.PingRangePct) { obj["PingRangePct"] = m.PingRangePct }
	if !fz(m.EvasionPct) { obj["EvasionPct"] = m.EvasionPct }
	if !fz(m.FormationSyncBonus) { obj["FormationSyncBonus"] = m.FormationSyncBonus }
	if !fz(m.PositionFlexibility) { obj["PositionFlexibility"] = m.PositionFlexibility }
	if !fz(m.GlobalDefensePct) { obj["GlobalDefensePct"] = m.GlobalDefensePct }
	if !fz(m.HPPct) { obj["HPPct"] = m.HPPct }
	return json.Marshal(obj)
}

func (m StatMods) MarshalBSON() ([]byte, error) {
	doc := bson.M{}
	if !m.Damage.IsZero() {
		dmg := bson.M{}
		if !fz(m.Damage.LaserPct) { dmg["LaserPct"] = m.Damage.LaserPct }
		if !fz(m.Damage.NuclearPct) { dmg["NuclearPct"] = m.Damage.NuclearPct }
		if !fz(m.Damage.AntimatterPct) { dmg["AntimatterPct"] = m.Damage.AntimatterPct }
		if len(dmg) > 0 { doc["Damage"] = dmg }
	}
	if !fz(m.AttackIntervalPct) { doc["AttackIntervalPct"] = m.AttackIntervalPct }
	if m.SpeedDelta != 0 { doc["SpeedDelta"] = m.SpeedDelta }
	if m.VisibilityDelta != 0 { doc["VisibilityDelta"] = m.VisibilityDelta }
	if m.AttackRangeDelta != 0 { doc["AttackRangeDelta"] = m.AttackRangeDelta }
	if m.LaserShieldDelta != 0 { doc["LaserShieldDelta"] = m.LaserShieldDelta }
	if m.NuclearShieldDelta != 0 { doc["NuclearShieldDelta"] = m.NuclearShieldDelta }
	if m.AntimatterShieldDelta != 0 { doc["AntimatterShieldDelta"] = m.AntimatterShieldDelta }
	if !fz(m.BucketHPPct) { doc["BucketHPPct"] = m.BucketHPPct }
	if !fz(m.OutOfCombatRegenPct) { doc["OutOfCombatRegenPct"] = m.OutOfCombatRegenPct }
	if !fz(m.AtCombatRegenPct) { doc["AtCombatRegenPct"] = m.AtCombatRegenPct }
	if !fz(m.AbilityCooldownPct) { doc["AbilityCooldownPct"] = m.AbilityCooldownPct }
	if !fz(m.TransportCapacityPct) { doc["TransportCapacityPct"] = m.TransportCapacityPct }
	if !fz(m.WarpChargePct) { doc["WarpChargePct"] = m.WarpChargePct }
	if !fz(m.WarpScatterPct) { doc["WarpScatterPct"] = m.WarpScatterPct }
	if !fz(m.InterdictionResistPct) { doc["InterdictionResistPct"] = m.InterdictionResistPct }
	if !fz(m.StructureDamagePct) { doc["StructureDamagePct"] = m.StructureDamagePct }
	if m.SplashRadiusDelta != 0 { doc["SplashRadiusDelta"] = m.SplashRadiusDelta }
	if !fz(m.AccuracyPct) { doc["AccuracyPct"] = m.AccuracyPct }
	if !fz(m.CritPct) { doc["CritPct"] = m.CritPct }
	if !fz(m.FirstVolleyPct) { doc["FirstVolleyPct"] = m.FirstVolleyPct }
	if !fz(m.ShieldPiercePct) { doc["ShieldPiercePct"] = m.ShieldPiercePct }
	if !fz(m.UpkeepPct) { doc["UpkeepPct"] = m.UpkeepPct }
	if !fz(m.ConstructionCostPct) { doc["ConstructionCostPct"] = m.ConstructionCostPct }
	if m.CloakDetect { doc["CloakDetect"] = true }
	if !fz(m.PingRangePct) { doc["PingRangePct"] = m.PingRangePct }
	if !fz(m.EvasionPct) { doc["EvasionPct"] = m.EvasionPct }
	if !fz(m.FormationSyncBonus) { doc["FormationSyncBonus"] = m.FormationSyncBonus }
	if !fz(m.PositionFlexibility) { doc["PositionFlexibility"] = m.PositionFlexibility }
	if !fz(m.GlobalDefensePct) { doc["GlobalDefensePct"] = m.GlobalDefensePct }
	if !fz(m.HPPct) { doc["HPPct"] = m.HPPct }
	return bson.Marshal(doc)
}
