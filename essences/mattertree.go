package essences

/////////////////////
// Matter (Element)
/////////////////////

// Aura defines a persistent area effect.
type Aura struct {
	Name        string
	Radius      float64
	Effect      Modifier
	Description string
}

// Matter describes the element tree applied per ship (tiered, ends with Active).
type Matter struct {
	Name        string
	Description string
	Tiers       [][]Modifier // each tier: list of modifiers to pick from (simplified)
	Auras       []Aura
	Active      Modifier // representative active ability modifier (you'll model actual ability code separately)
	Notes       string
}

///////////////////////////////
// Matter trees (one per type)
///////////////////////////////

func BuildPlasma() *Matter {
	return &Matter{
		Name:        "Plasma",
		Description: "Volatile ionized matter: burst damage, chain reactions, high energy.",
		Tiers: [][]Modifier{
			{{Name: "IonBurn", StatDelta: StatDelta{DamagePercent: 0.10, RegenPerTick: -0.01}}},
			{{Name: "OverloadCaps", Description: "On crit chance to chain damage"}},
			{{Name: "RadiantField", StatDelta: StatDelta{DamagePercent: 0.05, ArmorPercent: -0.05}}},
		},
		Auras: []Aura{
			{Name: "Radiant Field", Radius: 150, Effect: Modifier{Name: "PlasmaRadiant", StatDelta: StatDelta{DamagePercent: 0.05}}, Description: "Increases damage, reduces armor"}},
		Active: Modifier{Name: "SolarFlare", Description: "AoE nuke / shield disrupt"},
	}
}

func BuildLiquid() *Matter {
	return &Matter{
		Name:        "Liquid",
		Description: "Flowing matter: sustain, shields and movement synergy.",
		Tiers: [][]Modifier{
			{{Name: "HydroShield", StatDelta: StatDelta{RegenPerTick: 0.02, ArmorPercent: 0.05}}},
			{{Name: "PressureSurge", StatDelta: StatDelta{SpeedPercent: 0.10}}},
			{{Name: "TidalRelay", StatDelta: StatDelta{RegenPerTick: 0.02, SpeedPercent: 0.02}}},
		},
		Auras: []Aura{
			{Name: "Hydrostatic Field", Radius: 200, Effect: Modifier{Name: "HydroField", StatDelta: StatDelta{RegenPerTick: 0.05, ShieldEfficiencyPct: 1}}, Description: "Grants regen & +1 laser shield"}},
		Active: Modifier{Name: "Maelstrom", Description: "Pull / immobilize AoE (flow control)"},
	}
}

func BuildGas() *Matter {
	return &Matter{
		Name:        "Gas",
		Description: "Diffuse matter: concealment, DoT and mobility tactics.",
		Tiers: [][]Modifier{
			{{Name: "VaporShroud", StatDelta: StatDelta{EvasionPct: 0.10, ArmorPercent: -0.05}}},
			{{Name: "Chemcloud", Description: "On hit chance to apply DoT"}},
			{{Name: "Smokescreen", Description: "Aura reduces enemy accuracy"}},
		},
		Auras: []Aura{
			{Name: "Corrosive Miasma", Radius: 200, Effect: Modifier{Name: "CorrosiveMiasma", StatDelta: StatDelta{ArmorPercent: -0.01}}, Description: "Steady armor erosion"}},
		Active: Modifier{Name: "VaporSurge", Description: "Phase through formation & leave toxic trail"},
	}
}

func BuildEnergy() *Matter {
	return &Matter{
		Name:        "Energy",
		Description: "Light & shields: defensive mastery and aura power.",
		Tiers: [][]Modifier{
			{{Name: "ShieldConduction", StatDelta: StatDelta{ShieldEfficiencyPct: 0.10, ArmorPercent: -0.10}}},
			{{Name: "ResonanceCore", Description: "Abilities have higher efficiency when shields are high"}},
			{{Name: "PhotonicLattice", Description: "Aura: ability power amplifier"}},
		},
		Auras: []Aura{
			{Name: "Resonant Field", Radius: 250, Effect: Modifier{Name: "ResonantField", StatDelta: StatDelta{DamagePercent: 0.05}}, Description: "Converts damage into energy"}},
		Active: Modifier{Name: "StormPulse", Description: "Chain lightning that reduces shields"},
	}
}
