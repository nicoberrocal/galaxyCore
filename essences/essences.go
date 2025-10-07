package essences

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// LORE (top-level summary)
//
// The multiverse is gardened by the Echelons. Their interventions produce
// "Essences" — different physics substrates (Vital, Causality, Null, Entropic).
// Players pick an Essence (the universe their civ hails from), a Biology
// (how their civilization evolved: Mycelia, Aquatica, Flora, Fauna), and a
// Matter expression (what their ships are made of: Plasma, Liquid, Gas, Energy).
//
// - Essence: global, civ-wide modifiers (passive, eco, tactical, scouting).
// - Biology: fleet-level doctrine tree (4 tiers, choices, mutually-exclusive nodes).
//            Some bio trees have node-based essence mutations, some path-based.
// - Matter: ship-level element tree (auras, active abilities, damage interactions).
// - Formation: tactical and positional, intentionally decoupled from Essence/Bio/Matter.
//
// This file models the data structures and provides starter instances for each
// of the 4 essences, 4 biology trees, and 4 matter trees. It shows a simple
// application flow: create a Civ with Essence+Biology+Matter and apply modifiers
// to a Ship instance.
//
// The goal: preserve tactical roster (11 ships) while producing deep strategic
// identity via Essence × Biology × Matter composition.
////////////////////////////////////////////////////////////////////////////////

/////////////////////////
// Basic types & helpers
/////////////////////////

// StatDelta represents a bundle of stat modifications.
// Keep this small and extend as your game needs (HP, Regen, Damage, Speed, Shields, etc.)
type StatDelta struct {
	HPPercent           float64 // additive percents (0.02 = +2%)
	RegenPerTick        float64 // absolute per-tick regen
	SpeedPercent        float64
	DamagePercent       float64
	AccuracyPercent     float64
	ArmorPercent        float64
	ShieldEfficiencyPct float64
	AoEResistPct        float64
	KineticArmorPct     float64
	EvasionPct          float64
	// Add other fields as you need...
}

// Merge applies other into s (simple additive merge).
func (s *StatDelta) Merge(other StatDelta) {
	s.HPPercent += other.HPPercent
	s.RegenPerTick += other.RegenPerTick
	s.SpeedPercent += other.SpeedPercent
	s.DamagePercent += other.DamagePercent
	s.AccuracyPercent += other.AccuracyPercent
	s.ArmorPercent += other.ArmorPercent
	s.ShieldEfficiencyPct += other.ShieldEfficiencyPct
	s.AoEResistPct += other.AoEResistPct
	s.KineticArmorPct += other.KineticArmorPct
	s.EvasionPct += other.EvasionPct
}

// Modifier describes an effect that can be applied (simple wrapper).
type Modifier struct {
	Name        string
	StatDelta   StatDelta
	Description string
	// Condition or trigger fields could be added (OnCrit, OnDeath, OnAbilityUse...)
}

////////////////////////////
// Core domain structures
////////////////////////////

// EssenceType enumerates essences.
type EssenceType string

const (
	EssenceVital     EssenceType = "VitalContinuum"
	EssenceCausality EssenceType = "CausalityWeb"
	EssenceNull      EssenceType = "NullHorizon"
	EssenceEntropic  EssenceType = "EntropicVerge"
)

// Essence is a civ-level package of modifiers and hooks.
type Essence struct {
	Type          EssenceType
	DisplayName   string
	Description   string
	Passive       Modifier // always-on passive effect
	EcoModifier   Modifier // economic modifier (non-combat)
	TacticalMod   Modifier // tactical / combat-layer
	ScoutingMod   Modifier // scouting & intel modifiers
	Notes         string
	ExtraMetadata map[string]interface{}
}

/////////////////////
// Ship & Civ models
/////////////////////

type ShipClass string

const (
	ClassDrone     ShipClass = "Drone"
	ClassScout     ShipClass = "Scout"
	ClassFighter   ShipClass = "Fighter"
	ClassBomber    ShipClass = "Bomber"
	ClassCarrier   ShipClass = "Carrier"
	ClassDestroyer ShipClass = "Destroyer"
	ClassBallista  ShipClass = "Ballista"
	ClassGhost     ShipClass = "Ghost"
	ClassCruiser   ShipClass = "Cruiser"
	ClassCorvette  ShipClass = "Corvette"
	ClassFrigate   ShipClass = "Frigate" // example - use 11 in your real roster
)

// ShipStats is the numeric baseline for a ship type.
type ShipStats struct {
	MaxHP        float64
	Damage       float64
	Speed        float64
	Accuracy     float64
	Armor        float64
	ShieldValue  float64
	RegenPerTick float64
	// etc...
}

// Ship is an instance of a tactical unit. Derived stats will be computed.
type Ship struct {
	ID        string
	Class     ShipClass
	BaseStats ShipStats
	// DerivedStats is the resulting stats after civ modifiers applied.
	DerivedStats ShipStats
	Abilities    []string
	Auras        []Aura
	// Meta
	Essence EssenceType
	Biology string
	Matter  string
}

type Civ struct {
	Name    string
	Essence EssenceType
	// Single biology for fleet (recommended for first iteration)
	Biology *BioTree
	// Matter expression is per-civ default — you may allow per-class override later.
	Matter *Matter
	// You may also want to store selected nodes and active path choices:
	SelectedNodes map[string]*BioNode // nodeID -> selected node
	SelectedPath  string              // for Fauna path selection
}

//////////////////////////////
// Concrete Essences (starter)
//////////////////////////////

var EssenceCatalog = map[EssenceType]Essence{
	EssenceVital: {
		Type:        EssenceVital,
		DisplayName: "Vital Continuum",
		Description: "Energy is abundant, entropy is slow. Biology thrives; life regenerates.",
		Passive:     Modifier{Name: "OrganicRegeneration", StatDelta: StatDelta{RegenPerTick: 0.02}},
		EcoModifier: Modifier{Name: "PhotosyntheticExpansion", Description: "-15% build cost/time on habitable systems"},
		TacticalMod: Modifier{Name: "VitalSurge", StatDelta: StatDelta{DamagePercent: 0.15, SpeedPercent: 0.15}},
		ScoutingMod: Modifier{Name: "BiologicalResonance", Description: "Reveals Biology type on scan"},
		Notes:       "Regenerative baseline; great for sustained, slow-burn engagements.",
	},
	EssenceCausality: {
		Type:        EssenceCausality,
		DisplayName: "Causality Web",
		Description: "Past and future bleed into the present — echoes and pre-echoes are common.",
		Passive:     Modifier{Name: "TemporalEcho", Description: "First action echoes at 60% power"},
		EcoModifier: Modifier{Name: "RecursiveLogistics", Description: "-10% build time; 5% refund chance"},
		TacticalMod: Modifier{Name: "SlipstreamInitiative", StatDelta: StatDelta{AccuracyPercent: 0.10, SpeedPercent: 0.05}},
		ScoutingMod: Modifier{Name: "TemporalGlimpse", Description: "Scouts can reveal formation 1 tick early (chance)"},
		Notes:       "Requires careful timing. Reward for scouting and planning.",
	},
	EssenceNull: {
		Type:        EssenceNull,
		DisplayName: "Null Horizon",
		Description: "Matter flickers; stasis and reflection are common. Stability favored over throughput.",
		Passive:     Modifier{Name: "EventEquilibrium", StatDelta: StatDelta{ShieldEfficiencyPct: 0.10, RegenPerTick: -0.05}},
		EcoModifier: Modifier{Name: "EntropyHarvest", Description: "+20% yield from derelict/decay"},
		TacticalMod: Modifier{Name: "GraviticInversion", Description: "Once per battle deflects projectiles within radius"},
		ScoutingMod: Modifier{Name: "HorizonEcho", Description: "Detect anomalies but reveal self on scan"},
		Notes:       "Strong defensive identity; good vs burst but weak vs sustained pressure.",
	},
	EssenceEntropic: {
		Type:        EssenceEntropic,
		DisplayName: "Entropic Verge",
		Description: "Physical laws fray; corrosion, decay and inversion are part of warfare.",
		Passive:     Modifier{Name: "DecayConstant", Description: "Enemy shields decay while in combat", StatDelta: StatDelta{}},
		EcoModifier: Modifier{Name: "EntropyDividend", Description: "+15% yield from dismantling own ships"},
		TacticalMod: Modifier{Name: "UnstableReversal", Description: "Chance DoT you inflict inverts to heal you"},
		ScoutingMod: Modifier{Name: "EntropyDrift", Description: "Enemy sensors get false positives for your fleet"},
		Notes:       "Chaos-oriented playstyle: attrition, denial, unpredictability.",
	},
}

//////////////////////////////////////
// Simple modifier application model
//////////////////////////////////////

// ApplyCivModifiers resolves the civ-level Essence + Biology + Matter on a ship.
// This is intentionally a simple deterministic merger of modifiers.
// Real game logic should respect triggers (onCrit/onDeath/etc.) and do event-driven resolution.
func ApplyCivModifiers(ship *Ship, civ *Civ) {
	// start by copying base to derived
	ship.DerivedStats = ship.BaseStats

	// Apply Essence passive
	if ess, ok := EssenceCatalog[civ.Essence]; ok {
		mergeModifierToShipStats(&ship.DerivedStats, ess.Passive.StatDelta)
		// Tactical and scouting modifiers are handled elsewhere in game loop
		// For example: ess.TacticalMod.Description contains special rules (like "echo")
		_ = ess // keep for reference
	}

	// Apply Biology selected nodes
	// Here we assume civ.SelectedNodes holds the player's chosen nodes by ID.
	for _, node := range civ.SelectedNodes {
		if node == nil {
			continue
		}

		mergeModifierToShipStats(&ship.DerivedStats, node.StatDelta)
		if node.Tradeoff != nil {
			// Convert StatMods to StatDelta for compatibility
			tradeoffDelta := StatDelta{
				HPPercent:           0, // Would need conversion logic
				RegenPerTick:        0,
				SpeedPercent:        0,
				DamagePercent:       0,
				AccuracyPercent:     0,
				ArmorPercent:        0,
				ShieldEfficiencyPct: 0,
				AoEResistPct:        0,
				KineticArmorPct:     0,
				EvasionPct:          0,
			}
			mergeModifierToShipStats(&ship.DerivedStats, tradeoffDelta)
		}

	}

	// Apply Matter tiers / auras
	if civ.Matter != nil {
		// For simplicity apply all tier modifiers (you'll use selected tier in production)
		for _, tier := range civ.Matter.Tiers {
			for _, mod := range tier {
				mergeModifierToShipStats(&ship.DerivedStats, mod.StatDelta)
			}
		}
		// Apply aura effects as well (if ship is within aura radius you'll apply in fleet context)
		for _, a := range civ.Matter.Auras {
			_ = a // aura application is contextual (fleet positioning)
		}
		// Apply matter active ability baseline stat delta (if passive effects)
		mergeModifierToShipStats(&ship.DerivedStats, civ.Matter.Active.StatDelta)
	}
}

// mergeModifierToShipStats applies a StatDelta into ShipStats (simple additive transforms).
func mergeModifierToShipStats(s *ShipStats, d StatDelta) {
	// Apply percents and absolutes. This is intentionally simple — customize to your stat system.
	s.MaxHP *= (1.0 + d.HPPercent)
	s.RegenPerTick += d.RegenPerTick
	s.Speed *= (1.0 + d.SpeedPercent)
	s.Damage *= (1.0 + d.DamagePercent)
	s.Accuracy *= (1.0 + d.AccuracyPercent)
	s.Armor *= (1.0 + d.ArmorPercent)
	s.ShieldValue *= (1.0 + d.ShieldEfficiencyPct)
	// Note: ensure these multipliers are safe in your production code (no negative speeds, etc).
}

//////////////////////
// Example usage/demo
//////////////////////

// BuildStarterCiv constructs a sample civ
func BuildStarterCiv(name string, essence EssenceType, bio *BioTree, matter *Matter) *Civ {
	c := &Civ{
		Name:          name,
		Essence:       essence,
		Biology:       bio,
		Matter:        matter,
		SelectedNodes: make(map[string]*BioNode),
	}
	// Select default nodes (for demo, pick first node in each tier)
	for _, tier := range bio.Tiers {
		if len(tier) > 0 {
			node := tier[0]
			c.SelectedNodes[node.ID] = node
		}
	}
	// For Fauna allow a default path selection (when implemented)
	// if bio.Paths != nil {
	// 	for k := range bio.Paths {
	// 		c.SelectedPath = k
	// 		break
	// 	}
	// }
	return c
}

// DemoRandomShip applies a random ship and prints derived stats.
func DemoRandomShip() {
	// rand.Seed(time.Now().UnixNano())

	// Build trees & matter
	// mycelia := BuildMycelia()
	aquatica := BuildAquatica()
	// flora := BuildFlora()
	// fauna := BuildFauna()
	// plasma := BuildPlasma()
	liquid := BuildLiquid()
	// gas := BuildGas()
	// energy := BuildEnergy()

	// Example civ: Causality Web + Aquatica + Liquid (per your earlier ask)
	civ := BuildStarterCiv("Temporal Mariners", EssenceCausality, aquatica, liquid)

	// Pick a sample ship
	ship := &Ship{
		ID:    "ship-001",
		Class: ClassCruiser,
		BaseStats: ShipStats{
			MaxHP:        450,
			Damage:       40,
			Speed:        4,
			Accuracy:     1,
			Armor:        0.20,
			ShieldValue:  100,
			RegenPerTick: 0,
		},
		Essence: civ.Essence,
		Biology: civ.Biology.Name,
		Matter:  civ.Matter.Name,
	}

	ApplyCivModifiers(ship, civ)

	fmt.Printf("Ship [%s] DerivedStats: HP=%.2f Damage=%.2f Speed=%.2f Acc=%.2f Armor=%.2f Shield=%.2f Regen=%.3f\n",
		ship.Class, ship.DerivedStats.MaxHP, ship.DerivedStats.Damage, ship.DerivedStats.Speed, ship.DerivedStats.Accuracy, ship.DerivedStats.Armor, ship.DerivedStats.ShieldValue, ship.DerivedStats.RegenPerTick)
	// Output is illustrative; numeric design still required.
}
