package ships

import (
    "math"
    "math/rand"
    "sort"
    "strings"
    "time"
)

/*
Cosmic Synthesis Gem System
==========================

Overview
--------
A sophisticated gem system combining Diablo-style upgrading with a rich synthesis system. Gems represent crystallized cosmic energies that can be socketed into ships to enhance their capabilities.

Core Concepts
-------------
1. Gem Families: Thematic groups (Laser, Nuclear, Antimatter, Kinetic, Sensor, Warp, Engineering, Logistics)
2. Tiers: Power levels from 1 (weakest) to 5 (strongest)
3. Kinds: Pure (base family), Hybrid (2 families), or Relic (3+ families)
4. Origins: Astrophysical source (Pulsar, Supernova, etc.) influencing behavior
5. Instability: Risk factor from 0 (stable) to 1 (volatile)

Core Systems
------------
1. Pure Upgrading:
   - Combine 3 identical gems → 1 gem of next tier (up to tier 5)
   - Preserves family and origin
   - 100% success rate

2. Hybrid Synthesis:
   - Combine 2 different families → New hybrid gem
   - Success depends on affinity between families
   - Output tier based on input energy and efficiency
   - May accumulate instability

3. Relic Creation:
   - Combine 3 specific high-tier hybrids → Powerful relic
   - High risk/reward with special effects
   - Example: 3x T5 hybrids → "Singularity Core"

4. Energy Model:
   - Each gem has an energy value based on tier and origin
   - Output energy = sum(input energy) * efficiency
   - Excess energy increases success chance

5. Affinity System:
   - Families have natural affinities (synergy/conflict)
   - Affects success rate and energy efficiency
   - Influences output origin and stats

6. GemWords:
   - Special sequences that unlock bonus effects
   - Order and tier requirements matter
   - Multiple GemWords can activate simultaneously

Usage
-----
- Use `SynthesizeGems()` to merge gems
- Call `EvaluateGemSockets()` to calculate total effects
- Check `GemCatalog` for all possible gems
- See `GemWordsCatalog` for special sequences

Example Flow:
1. Player collects basic gems (Laser I, Nuclear I, etc.)
2. Upgrades to higher tiers (3x Laser I → Laser II)
3. Experiments with hybrid combinations (Laser + Nuclear)
4. Crafts advanced relics (3x T5 hybrids)
5. Discovers powerful GemWord combinations

Balance Notes:
- Pure upgrades are safe but limited
- Hybrids offer unique combinations but with risk
- Relics are powerful but require investment
- Instability adds risk/reward tension
*/

// GemFamily identifies the thematic group of a gem.
type GemFamily string

const (
	GemLaser       GemFamily = "laser"
	GemNuclear     GemFamily = "nuclear"
	GemAntimatter  GemFamily = "antimatter"
	GemKinetic     GemFamily = "kinetic"
	GemSensor      GemFamily = "sensor"
	GemWarp        GemFamily = "warp"
	GemEngineering GemFamily = "engineering"
	GemLogistics   GemFamily = "logistics"
)

// GemID is a stable identifier for a specific family/tier gem.
type GemID string

// Gem is a socketable item granting StatMods and sometimes gating abilities.
// Mods scale with Family and Tier. Use EvaluateGemSockets to aggregate mods and
// detect GemWords.
type Gem struct {
    ID          GemID
    Name        string
    Family      GemFamily
    Tier        int // 1..MaxTier
    Mods        StatMods
    Description string
    // Cosmic synthesis fields
    Origin      GemOrigin   // astrophysical origin flavor (Pulsar, Supernova, ...)
    Kind        GemKind     // pure, hybrid, relic
    Instability float64     // 0..1 accumulated instability from mixed merges
    Parents     []GemID     // provenance for hybrids/relics
    Grants      []AbilityID // per-gem abilities (in addition to GemWords)
}
// GemOrigin is the astrophysical origin flavor of a gem.
type GemOrigin string

const (
    OriginPulsar         GemOrigin = "pulsar"
    OriginSupernova      GemOrigin = "supernova"
    OriginWormhole       GemOrigin = "wormhole"
    OriginSingularity    GemOrigin = "singularity"
    OriginNebula         GemOrigin = "nebula"
    OriginQuasar         GemOrigin = "quasar"
    OriginBigBangRemnant GemOrigin = "bigbang_remnant"
)

// GemKind differentiates baseline families from synthesized composites.
type GemKind string

const (
    KindPure   GemKind = "pure"
    KindHybrid GemKind = "hybrid"
    KindRelic  GemKind = "relic"
)
func defaultOrigin(f GemFamily, tier int) GemOrigin {
    switch f {
    case GemLaser:
        return OriginPulsar
    case GemNuclear, GemKinetic:
        return OriginSupernova
    case GemWarp:
        return OriginWormhole
    case GemAntimatter:
        return OriginSingularity
    case GemEngineering, GemLogistics:
        return OriginNebula
    case GemSensor:
        if tier >= 4 {
            return OriginQuasar
        }
        return OriginPulsar
    default:
        return OriginNebula
    }
}

// Affinity encodes how easily two families synthesize and how efficient they are.
type Affinity struct {
    Weight      float64   // -1..+1, negative means clash
    BaseSuccess float64   // baseline success probability for mixed merges
    Efficiency  float64   // energy efficiency multiplier for mixed merges
    OriginHint  GemOrigin // suggested origin for the hybrid
}

var AffinityGraph map[GemFamily]map[GemFamily]Affinity

func initAffinity() {
    AffinityGraph = make(map[GemFamily]map[GemFamily]Affinity)
    set := func(a, b GemFamily, aff Affinity) {
        if AffinityGraph[a] == nil { AffinityGraph[a] = map[GemFamily]Affinity{} }
        if AffinityGraph[b] == nil { AffinityGraph[b] = map[GemFamily]Affinity{} }
        AffinityGraph[a][b] = aff
        AffinityGraph[b][a] = aff
    }
    // Synergies
    set(GemLaser, GemNuclear,  Affinity{Weight: 0.7, BaseSuccess: 0.85, Efficiency: 0.80, OriginHint: OriginSupernova})
    set(GemNuclear, GemKinetic, Affinity{Weight: 0.6, BaseSuccess: 0.85, Efficiency: 0.80, OriginHint: OriginSupernova})
    set(GemLaser, GemSensor,   Affinity{Weight: 0.55, BaseSuccess: 0.80, Efficiency: 0.75, OriginHint: OriginQuasar})
    set(GemSensor, GemWarp,    Affinity{Weight: 0.50, BaseSuccess: 0.75, Efficiency: 0.70, OriginHint: OriginWormhole})
    set(GemEngineering, GemLogistics, Affinity{Weight: 0.65, BaseSuccess: 0.85, Efficiency: 0.80, OriginHint: OriginNebula})
    // Risky/Clash
    set(GemAntimatter, GemKinetic, Affinity{Weight: -0.5, BaseSuccess: 0.50, Efficiency: 0.45, OriginHint: OriginSingularity})
    // Mixed/risky but rewarding
    set(GemAntimatter, GemLaser, Affinity{Weight: 0.30, BaseSuccess: 0.70, Efficiency: 0.65, OriginHint: OriginSingularity})
    set(GemSensor, GemAntimatter, Affinity{Weight: -0.20, BaseSuccess: 0.60, Efficiency: 0.60, OriginHint: OriginQuasar})
    set(GemWarp, GemEngineering,  Affinity{Weight: 0.40, BaseSuccess: 0.75, Efficiency: 0.70, OriginHint: OriginWormhole})
    set(GemWarp, GemLogistics,    Affinity{Weight: 0.35, BaseSuccess: 0.75, Efficiency: 0.70, OriginHint: OriginNebula})
}

func getAffinity(a, b GemFamily) Affinity {
    if a == b { return Affinity{Weight: 1, BaseSuccess: 1.0, Efficiency: 0.90, OriginHint: defaultOrigin(a, 1)} }
    if m, ok := AffinityGraph[a]; ok {
        if aff, ok2 := m[b]; ok2 { return aff }
    }
    return Affinity{Weight: 0.0, BaseSuccess: 0.70, Efficiency: 0.65, OriginHint: OriginNebula}
}

// Merge recipes for known hybrids/relics.
type MergeRecipe struct {
    Key        string
    NameBase   string
    Families   []GemFamily // 2 or 3, order-insensitive
    MinTier    int
    Dominant   GemFamily
    Origin     GemOrigin
    Kind       GemKind
    Grants     []AbilityID
}

var MergeRecipeCatalog []MergeRecipe
var mergeRecipeByKey map[string]MergeRecipe

func initMergeRecipes() {
    MergeRecipeCatalog = []MergeRecipe{
        {NameBase: "Plasma Nova", Families: []GemFamily{GemLaser, GemNuclear}, MinTier: 1, Dominant: GemLaser, Origin: OriginSupernova, Kind: KindHybrid},
        {NameBase: "Shockfront", Families: []GemFamily{GemNuclear, GemKinetic}, MinTier: 1, Dominant: GemNuclear, Origin: OriginSupernova, Kind: KindHybrid},
        {NameBase: "Lensing Array", Families: []GemFamily{GemLaser, GemSensor}, MinTier: 1, Dominant: GemSensor, Origin: OriginQuasar, Kind: KindHybrid},
        {NameBase: "Interdictor Web", Families: []GemFamily{GemSensor, GemWarp}, MinTier: 1, Dominant: GemWarp, Origin: OriginWormhole, Kind: KindHybrid},
        {NameBase: "Annihilation Lattice", Families: []GemFamily{GemAntimatter, GemLaser}, MinTier: 1, Dominant: GemAntimatter, Origin: OriginSingularity, Kind: KindHybrid},
        {NameBase: "Stellar Forge", Families: []GemFamily{GemEngineering, GemLogistics}, MinTier: 1, Dominant: GemLogistics, Origin: OriginNebula, Kind: KindHybrid},
        {NameBase: "Quasar Jet", Families: []GemFamily{GemSensor, GemAntimatter}, MinTier: 1, Dominant: GemSensor, Origin: OriginQuasar, Kind: KindHybrid},
        {NameBase: "Warp Manifold", Families: []GemFamily{GemWarp, GemEngineering}, MinTier: 1, Dominant: GemWarp, Origin: OriginWormhole, Kind: KindHybrid},
        {NameBase: "Supply Nexus", Families: []GemFamily{GemWarp, GemLogistics}, MinTier: 1, Dominant: GemLogistics, Origin: OriginNebula, Kind: KindHybrid},
    }
    mergeRecipeByKey = map[string]MergeRecipe{}
    for i := range MergeRecipeCatalog {
        r := MergeRecipeCatalog[i]
        r.Key = canonicalKey(r.Families)
        MergeRecipeCatalog[i] = r
        mergeRecipeByKey[r.Key] = r
    }
}

// GemMergeResult is the outcome of a synthesis attempt.
type GemMergeResult struct {
    Success             bool
    Output              Gem
    SuccessProbability  float64
    Consumed            []Gem
    Returned            []Gem
    FailureReason       string
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SynthesizeGems merges 2-3 gems into a new gem following affinity/energy rules.
func SynthesizeGems(inputs []Gem) GemMergeResult {
    res := GemMergeResult{Consumed: append([]Gem(nil), inputs...)}
    if len(inputs) < 2 || len(inputs) > 3 {
        res.FailureReason = "need 2 or 3 gems"
        return res
    }

    // Pure deterministic upgrade: 3 of the same family & tier
    if len(inputs) == 3 && isPureTriple(inputs) {
        out, ok := UpgradeGem(inputs[0])
        if !ok { res.FailureReason = "cannot upgrade"; return res }
        out.Origin = defaultOrigin(out.Family, out.Tier)
        out.Kind = KindPure
        res.Output = out
        res.Success = true
        res.SuccessProbability = 1.0
        return res
    }

    // Relic: Singularity Core (3 distinct T5 hybrids)
    if len(inputs) == 3 && isRelicSingularityCoreInputs(inputs) {
        avgInst := averageInstability(inputs)
        p := clamp(0.40 - 0.25*avgInst, 0.05, 0.95)
        res.SuccessProbability = p
        if rng.Float64() <= p {
            tier := MaxGemTier
            mods := ZeroMods()
            for _, in := range inputs { mods = CombineMods(mods, scaleMods(in.Mods, 0.40)) }
            mods = CombineMods(mods, originFlavorMods(OriginBigBangRemnant, tier))
            out := Gem{
                ID:          GemID("relic-singularity-core-" + itoa(tier)),
                Name:        "Singularity Core " + roman(tier),
                Family:      GemAntimatter,
                Tier:        tier,
                Mods:        mods,
                Description: "Primordial relic forged from cataclysmic hybrids.",
                Origin:      OriginBigBangRemnant,
                Kind:        KindRelic,
                Instability: clamp(avgInst+0.10, 0, 1),
                Parents:     []GemID{inputs[0].ID, inputs[1].ID, inputs[2].ID},
            }
            GemCatalog[out.ID] = out
            res.Output = out
            res.Success = true
        } else {
            var back []Gem
            for _, in := range inputs {
                if in.Tier > 1 { in.Tier = maxInt(1, in.Tier-1); in.ID = GemID(familyID(in.Family, in.Tier)) }
                back = append(back, in)
            }
            res.Returned = back
            res.FailureReason = "relic synthesis failed"
        }
        return res
    }

    fams := uniqueFamilies(inputs)
    if len(fams) < 2 { res.FailureReason = "not enough distinct families for mixed merge"; return res }
    if len(fams) > 2 { res.FailureReason = "3-family hybrids not supported (except relic)"; return res }

    key := canonicalKey(fams)
    recipe, has := mergeRecipeByKey[key]
    if !has {
        aff := getAffinity(fams[0], fams[1])
        recipe = MergeRecipe{
            Key: key, NameBase: strings.Title(string(fams[0])) + "-" + strings.Title(string(fams[1])) + " Hybrid",
            Families: fams, MinTier: 1, Dominant: fams[0], Origin: aff.OriginHint, Kind: KindHybrid,
        }
    }

    var eIn float64
    avgInst := 0.0
    for _, in := range inputs { eIn += energyOf(in); avgInst += in.Instability }
    avgInst /= float64(len(inputs))

    aff := getAffinity(fams[0], fams[1])
    eAvailable := eIn * aff.Efficiency
    tier := tierFromEnergy(eAvailable, recipe.Origin)
    if tier < recipe.MinTier { tier = recipe.MinTier }
    tier = clampInt(tier, 1, MaxGemTier)

    eOut := energyForTierOrigin(tier, recipe.Origin)
    surplus := (eAvailable - eOut) / maxFloat(eOut, 1)
    p := aff.BaseSuccess + 0.25*surplus + 0.20*aff.Weight - 0.50*avgInst
    p = clamp(p, 0.05, 0.99)
    res.SuccessProbability = p

    if rng.Float64() <= p {
        mods := hybridMods(recipe.Families, tier, aff)
        mods = CombineMods(mods, originFlavorMods(recipe.Origin, tier))
        out := Gem{
            ID:          GemID("hybrid-" + slug(recipe.NameBase) + "-" + itoa(tier)),
            Name:        recipe.NameBase + " " + roman(tier),
            Family:      recipe.Dominant,
            Tier:        tier,
            Mods:        mods,
            Description: "Hybrid synthesized from " + strings.Join(familyNames(recipe.Families), "+") + ".",
            Origin:      recipe.Origin,
            Kind:        KindHybrid,
            Instability: clamp(avgInst+instabilityFromWaste(eAvailable, eOut), 0, 1),
            Parents:     []GemID{inputs[0].ID, inputs[1].ID},
            Grants:      recipe.Grants,
        }
        GemCatalog[out.ID] = out
        res.Output = out
        res.Success = true
    } else {
        idx := indexOfMaxTier(inputs)
        back := append([]Gem(nil), inputs...)
        if back[idx].Tier > 1 { back[idx].Tier--; back[idx].ID = GemID(familyID(back[idx].Family, back[idx].Tier)) }
        res.Returned = back
        res.FailureReason = "synthesis failed"
    }
    return res
}

func isPureTriple(in []Gem) bool {
    if len(in) != 3 { return false }
    f := in[0].Family; t := in[0].Tier
    for i := 1; i < 3; i++ { if in[i].Family != f || in[i].Tier != t { return false } }
    return true
}

func isRelicSingularityCoreInputs(in []Gem) bool {
    if len(in) != 3 { return false }
    names := map[string]bool{}
    for _, g := range in {
        if g.Kind != KindHybrid || g.Tier < MaxGemTier { return false }
        names[g.Name] = true
    }
    return len(names) == 3
}

func uniqueFamilies(in []Gem) []GemFamily {
    m := map[GemFamily]bool{}
    for _, g := range in { m[g.Family] = true }
    out := make([]GemFamily, 0, len(m))
    for k := range m { out = append(out, k) }
    return out
}

func averageInstability(in []Gem) float64 {
    if len(in) == 0 { return 0 }
    s := 0.0
    for _, g := range in { s += g.Instability }
    return s / float64(len(in))
}

func energyOf(g Gem) float64 { return energyForTierOrigin(g.Tier, g.Origin) }

func energyForTierOrigin(tier int, origin GemOrigin) float64 {
    if tier < 1 { tier = 1 }
    base := math.Pow(3, float64(tier-1))
    return base * originMultiplier(origin)
}

func originMultiplier(o GemOrigin) float64 {
    switch o {
    case OriginPulsar:
        return 1.00
    case OriginSupernova:
        return 1.10
    case OriginWormhole:
        return 1.10
    case OriginQuasar:
        return 1.05
    case OriginNebula:
        return 0.90
    case OriginSingularity:
        return 1.15
    case OriginBigBangRemnant:
        return 1.30
    default:
        return 1.00
    }
}

func tierFromEnergy(e float64, origin GemOrigin) int {
    best := 1
    for t := 1; t <= MaxGemTier; t++ { if energyForTierOrigin(t, origin) <= e { best = t } }
    return best
}

func instabilityFromWaste(eAvail, eOut float64) float64 {
    if eOut <= 0 { return 0 }
    waste := eAvail - eOut
    if waste <= 0 { return 0.05 }
    frac := waste / eOut
    if frac > 1 { frac = 1 }
    return clamp(frac*0.25, 0.05, 0.25)
}

func hybridMods(fams []GemFamily, tier int, aff Affinity) StatMods {
    m := ZeroMods()
    base := ZeroMods()
    for _, f := range fams { base = CombineMods(base, gemTierMods(f, tier)) }
    scale := 0.50
    if aff.Weight >= 0.4 { scale = 0.60 } else if aff.Weight <= -0.3 { scale = 0.40 }
    m = CombineMods(m, scaleMods(base, scale))
    return m
}

func originFlavorMods(o GemOrigin, tier int) StatMods {
    m := ZeroMods()
    switch o {
    case OriginPulsar:
        m.AttackIntervalPct += -0.03 * float64(tier)
        m.AccuracyPct += 0.01 * float64(tier)
    case OriginSupernova:
        m.StructureDamagePct += 0.03 * float64(tier)
        if tier >= 2 { m.SplashRadiusDelta += 1 }
        if tier >= 4 { m.SplashRadiusDelta += 1 }
    case OriginWormhole:
        m.WarpChargePct += -0.05 * float64(tier)
        m.WarpScatterPct += -0.05 * float64(tier)
        m.InterdictionResistPct += 0.04 * float64(tier)
    case OriginSingularity:
        m.ShieldPiercePct += 0.03 * float64(tier)
        m.CritPct += 0.02 * float64(tier)
        m.FirstVolleyPct += 0.04 * float64(tier)
    case OriginNebula:
        m.UpkeepPct += -0.03 * float64(tier)
        m.ConstructionCostPct += -0.03 * float64(tier)
        m.TransportCapacityPct += 0.05 * float64(tier)
        m.AbilityCooldownPct += -0.02 * float64(tier)
    case OriginQuasar:
        m.PingRangePct += 0.08 * float64(tier)
        m.AccuracyPct += 0.01 * float64(tier)
        if tier >= 2 { m.VisibilityDelta += 1 }
        if tier >= 4 { m.VisibilityDelta += 1 }
    case OriginBigBangRemnant:
        m.CritPct += 0.03 * float64(tier)
        m.BucketHPPct += 0.05 * float64(tier)
        m.AbilityCooldownPct += -0.02 * float64(tier)
    }
    return m
}

// Utility helpers
func slug(s string) string { s = strings.ToLower(s); s = strings.ReplaceAll(s, " ", "-"); return s }

func familyNames(fams []GemFamily) []string {
    out := make([]string, 0, len(fams))
    for _, f := range fams { out = append(out, string(f)) }
    return out
}

func canonicalKey(fams []GemFamily) string {
    parts := make([]string, len(fams))
    for i, f := range fams { parts[i] = string(f) }
    sort.Strings(parts)
    return strings.Join(parts, "+")
}

func indexOfMaxTier(in []Gem) int {
    idx := 0
    best := in[0].Tier
    for i := 1; i < len(in); i++ { if in[i].Tier > best { best = in[i].Tier; idx = i } }
    return idx
}

func scaleMods(m StatMods, s float64) StatMods {
    out := ZeroMods()
    out.Damage.LaserPct = m.Damage.LaserPct * s
    out.Damage.NuclearPct = m.Damage.NuclearPct * s
    out.Damage.AntimatterPct = m.Damage.AntimatterPct * s
    out.AttackIntervalPct = m.AttackIntervalPct * s
    out.SpeedDelta = int(float64(m.SpeedDelta) * s)
    out.VisibilityDelta = int(float64(m.VisibilityDelta) * s)
    out.AttackRangeDelta = int(float64(m.AttackRangeDelta) * s)
    out.LaserShieldDelta = int(float64(m.LaserShieldDelta) * s)
    out.NuclearShieldDelta = int(float64(m.NuclearShieldDelta) * s)
    out.AntimatterShieldDelta = int(float64(m.AntimatterShieldDelta) * s)
    out.BucketHPPct = m.BucketHPPct * s
    out.OutOfCombatRegenPct = m.OutOfCombatRegenPct * s
    out.AbilityCooldownPct = m.AbilityCooldownPct * s
    out.TransportCapacityPct = m.TransportCapacityPct * s
    out.WarpChargePct = m.WarpChargePct * s
    out.WarpScatterPct = m.WarpScatterPct * s
    out.InterdictionResistPct = m.InterdictionResistPct * s
    out.StructureDamagePct = m.StructureDamagePct * s
    out.SplashRadiusDelta = int(float64(m.SplashRadiusDelta) * s)
    out.AccuracyPct = m.AccuracyPct * s
    out.CritPct = m.CritPct * s
    out.FirstVolleyPct = m.FirstVolleyPct * s
    out.ShieldPiercePct = m.ShieldPiercePct * s
    out.UpkeepPct = m.UpkeepPct * s
    out.ConstructionCostPct = m.ConstructionCostPct * s
    out.CloakDetect = m.CloakDetect
    out.PingRangePct = m.PingRangePct * s
    return out
}

func clamp(x, lo, hi float64) float64 { if x < lo { return lo }; if x > hi { return hi }; return x }
func clampInt(x, lo, hi int) int { if x < lo { return lo }; if x > hi { return hi }; return x }
func maxInt(a, b int) int { if a > b { return a }; return b }
func maxFloat(a, b float64) float64 { if a > b { return a }; return b }

const MaxGemTier = 5

// GemCatalog contains all gem definitions (auto-generated per family and tier).
var GemCatalog = map[GemID]Gem{}

// Scientific-flavored base names for gem families.
var gemFamilyBaseName = map[GemFamily]string{
	GemLaser:       "Photon",
	GemNuclear:     "Isotope",
	GemAntimatter:  "Positron",
	GemKinetic:     "Graphene",
	GemSensor:      "Neutrino",
	GemWarp:        "Axion",
	GemEngineering: "Nanoforge",
	GemLogistics:   "Catalyst",
}

func init() {
	// Populate the GemCatalog programmatically to keep data DRY and tunable.
	for family, base := range gemFamilyBaseName {
		for tier := 1; tier <= MaxGemTier; tier++ {
			id := GemID(familyID(family, tier))
			GemCatalog[id] = Gem{
				ID:          id,
				Name:        base + " " + roman(tier),
				Family:      family,
				Tier:        tier,
				Mods:        gemTierMods(family, tier),
				Description: gemDescription(family, tier),
				Origin:      defaultOrigin(family, tier),
				Kind:        KindPure,
				Instability: 0,
				Parents:     nil,
				Grants:      nil,
			}
		}
	}
	// initialize cosmic synthesis systems
	initAffinity()
	initMergeRecipes()
}

// UpgradeRule: 3x same family+tier -> 1x same family at tier+1 (if tier < MaxGemTier)
func CanUpgrade(g Gem) bool { return g.Tier < MaxGemTier }

func UpgradeGem(g Gem) (Gem, bool) {
	if !CanUpgrade(g) {
		return Gem{}, false
	}
	nextID := GemID(familyID(g.Family, g.Tier+1))
	out, ok := GemCatalog[nextID]
	return out, ok
}

// EvaluateGemSockets aggregates gem mods and resolves matching GemWords.
// Returns cumulative StatMods, granted abilities from GemWords, and matched GemWords.
// All matching GemWords apply; order of GemWordsCatalog controls stacking intent.
func EvaluateGemSockets(sockets []Gem) (StatMods, []AbilityID, []GemWord) {
	mods := ZeroMods()
	// Sum gem mods and per-gem grants
	var grants []AbilityID
	for _, g := range sockets {
		mods = CombineMods(mods, g.Mods)
		if len(g.Grants) > 0 {
			grants = append(grants, g.Grants...)
		}
	}
	// Resolve GemWords in order of declaration (first match wins if overlapping by design)
	var matched []GemWord
	for _, gw := range GemWordsCatalog {
		if gemwordMatches(gw, sockets) {
			mods = CombineMods(mods, gw.Effects)
			if len(gw.Grants) > 0 {
				grants = append(grants, gw.Grants...)
			}
			matched = append(matched, gw)
		}
	}
	return mods, grants, matched
}

// GemWord models an ordered recipe that unlocks extra effects when sockets match.
type GemWord struct {
	Name        string
	Sequence    []GemFamily // ordered, must match sockets[0:len(Sequence)] families
	MinTier     int         // each gem in the sequence must be at least this tier
	Effects     StatMods    // additional stat modifiers
	Grants      []AbilityID // abilities unlocked by this GemWord
	Description string
}

// GemWordsCatalog contains a small, thematic set aligned with your earlier gem sets.
var GemWordsCatalog = []GemWord{
	{
		Name:        "Photon Overcharge",
		Sequence:    []GemFamily{GemLaser, GemLaser, GemLaser},
		MinTier:     2,
		Effects:     StatMods{AttackIntervalPct: -0.15},
		Grants:      []AbilityID{AbilityLaserOvercharge},
		Description: "Triple Laser gems grant Laser Overcharge and increase ROF.",
	},
	{
		Name:        "Isotope Bunker Array",
		Sequence:    []GemFamily{GemNuclear, GemNuclear, GemNuclear},
		MinTier:     2,
		Effects:     StatMods{StructureDamagePct: 0.15, SplashRadiusDelta: 1},
		Grants:      []AbilityID{AbilityBunkerBuster},
		Description: "Triple Nuclear gems unlock Bunker Buster vs fortified structures.",
	},
	{
		Name:        "Positron Phase Matrix",
		Sequence:    []GemFamily{GemAntimatter, GemAntimatter, GemAntimatter},
		MinTier:     2,
		Effects:     StatMods{ShieldPiercePct: 0.15, FirstVolleyPct: 0.20},
		Grants:      []AbilityID{AbilityPhaseLance},
		Description: "Triple Antimatter gems phase through shields on opening strikes.",
	},
	{
		Name:        "Neutrino Wide-Area Ping",
		Sequence:    []GemFamily{GemSensor, GemSensor, GemSensor},
		MinTier:     2,
		Effects:     StatMods{VisibilityDelta: 2, PingRangePct: 0.50, AccuracyPct: 0.05},
		Grants:      []AbilityID{AbilityWideAreaPing},
		Description: "Triple Sensor gems enable wide-area recon ping.",
	},
	{
		Name:        "Axion Rapid Redeploy",
		Sequence:    []GemFamily{GemWarp, GemWarp, GemWarp},
		MinTier:     2,
		Effects:     StatMods{WarpChargePct: -0.25, WarpScatterPct: -0.30, InterdictionResistPct: 0.20},
		Grants:      []AbilityID{AbilityRapidRedeploy},
		Description: "Triple Warp gems reduce warp charge and scatter, enabling rapid redeploy.",
	},
	{
		Name:        "Graphene Bulwark Lattice",
		Sequence:    []GemFamily{GemKinetic, GemKinetic, GemKinetic},
		MinTier:     2,
		Effects:     StatMods{BucketHPPct: 0.20, LaserShieldDelta: 1, NuclearShieldDelta: 1, AntimatterShieldDelta: 1},
		Grants:      nil,
		Description: "Triple Kinetic gems increase hull resilience and shields.",
	},
	{
		Name:        "Catalyst Bay Commander",
		Sequence:    []GemFamily{GemLogistics, GemEngineering, GemWarp},
		MinTier:     1,
		Effects:     StatMods{TransportCapacityPct: 0.20, AbilityCooldownPct: -0.05, WarpChargePct: -0.10},
		Grants:      nil,
		Description: "Carrier-focused GemWord: larger bays, faster cooldowns, quicker warps.",
	},
}

// Helpers

func familyID(f GemFamily, tier int) string { return string(f) + "-" + itoa(tier) }

func gemwordMatches(gw GemWord, sockets []Gem) bool {
	if len(sockets) < len(gw.Sequence) {
		return false
	}
	for i, fam := range gw.Sequence {
		if sockets[i].Family != fam {
			return false
		}
		if sockets[i].Tier < gw.MinTier {
			return false
		}
	}
	return true
}

// gemTierMods defines baseline per-tier scaling per family. Numbers are conservative.
func gemTierMods(f GemFamily, tier int) StatMods {
	m := ZeroMods()
	switch f {
	case GemLaser:
		m.Damage.LaserPct = 0.06 * float64(tier)
		m.AccuracyPct = 0.025 * float64(tier)
		if tier >= 3 {
			m.SpeedDelta = 1
		}
	case GemNuclear:
		m.Damage.NuclearPct = 0.06 * float64(tier)
		m.StructureDamagePct = 0.05 * float64(tier)
		if tier >= 2 {
			m.SplashRadiusDelta = 1
		}
		if tier >= 4 {
			m.SplashRadiusDelta = 2
		}
	case GemAntimatter:
		m.Damage.AntimatterPct = 0.06 * float64(tier)
		m.FirstVolleyPct = 0.10 * float64(tier)
		m.ShieldPiercePct = 0.05 * float64(tier)
	case GemKinetic:
		m.BucketHPPct = 0.10 * float64(tier)
		m.OutOfCombatRegenPct = 0.10 * float64(tier)
	case GemSensor:
		m.VisibilityDelta = 1 * tier
		m.CloakDetect = true
		m.PingRangePct = 0.15 * float64(tier)
	case GemWarp:
		m.WarpChargePct = -0.10 * float64(tier)
		m.WarpScatterPct = -0.10 * float64(tier)
		m.InterdictionResistPct = 0.08 * float64(tier)
	case GemEngineering:
		m.AttackIntervalPct = -0.05 * float64(tier)
		m.AbilityCooldownPct = -0.05 * float64(tier)
		m.OutOfCombatRegenPct = 0.05 * float64(tier)
	case GemLogistics:
		m.UpkeepPct = -0.05 * float64(tier)
		m.ConstructionCostPct = -0.04 * float64(tier)
		m.TransportCapacityPct = 0.10 * float64(tier)
	}
	return m
}

func gemDescription(f GemFamily, tier int) string {
	switch f {
	case GemLaser:
		return "Laser damage, accuracy, slight speed at higher tiers."
	case GemNuclear:
		return "Nuclear damage, structure damage, splash at higher tiers."
	case GemAntimatter:
		return "Antimatter damage, first-volley, partial shield pierce."
	case GemKinetic:
		return "Hull durability and out-of-combat regeneration."
	case GemSensor:
		return "Vision, detection, and Ping range."
	case GemWarp:
		return "Warp charge speed, scatter reduction, interdiction resistance."
	case GemEngineering:
		return "Attack speed, ability cooldowns, and self-repair improvements."
	case GemLogistics:
		return "Upkeep/build cost reductions and transport capacity."
	}
	return ""
}

// Small, dependency-free helpers
func roman(n int) string {
	switch n {
	case 1:
		return "I"
	case 2:
		return "II"
	case 3:
		return "III"
	case 4:
		return "IV"
	case 5:
		return "V"
	default:
		return itoa(n)
	}
}

func itoa(n int) string {
	// minimal int to string without importing strconv to keep package light
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
