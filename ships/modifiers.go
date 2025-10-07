package ships

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
	CritPct            float64 // % flat critical chance improvement
	FirstVolleyPct     float64 // % bonus to the first volley
	ShieldPiercePct    float64 // % of shields ignored (applied carefully)

	// Economy/logistics
	UpkeepPct           float64 // % change to upkeep
	ConstructionCostPct float64 // % change to build costs

	// Recon/detection (boolean capabilities are OR-composed)
	CloakDetect  bool    // can detect cloaked or mode-switch signals
	PingRangePct float64 // % change to Ping ability range

	// Formation-specific modifiers (applied in formation combat contexts)
	EvasionPct          float64 // % flat evasion chance
	FormationSyncBonus  float64 // % bonus when position requirements are met
	PositionFlexibility float64 // % reduced penalty for suboptimal positions

	// Generic mods
	GlobalDefensePct float64 // % global damage reduction
	HPPct            float64 // % change to current HP
}

// ZeroMods returns a zero-initialized StatMods for convenience.
func ZeroMods() StatMods { return StatMods{} }

// CombineMods adds b into a and returns the result. Simple linear composition.
// Clamping should be enforced at the application layer if needed.
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
