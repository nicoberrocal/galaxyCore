package essences

import "github.com/nicoberrocal/galaxyCore/ships"

func BuildMycelia() *BioTree {
	tree := &BioTree{
		Name:        string(Mycelia),
		Description: "Mycelia biology tree focusing on infection, symbiosis, and area denial",
		Tiers:       make([][]*BioNode, 3),
	}

	// Tier 1: Sporeform (Stationary defense, infection, and persistence)
	sporeformNodes := []*BioNode{
		// 1. Dormant Spores: Stacking defense/regen while stationary
		{
			ID:          "sporeform_dormant_spores",
			Title:       "Dormant Spores",
			Description: "When stationary for 3+ ticks, gain a stacking +2% defense and +2% hull regen per tick (max 20%). Loses stacks instantly when moving.",
			Path:        string(ships.Sporeform),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionStationary, CompareOp: CompareGreaterEq, Value: 3},
					},
					PrimaryEffect: &ships.StatMods{
						GlobalDefensePct: 0.02,
						AtCombatRegenPct: 0.02,
					},
					Duration: 1, // Stacks per tick
				},
			},
		},
		// 2. Hyphal Invasion: Infects targets on hit
		{
			ID:          "sporeform_hyphal_invasion",
			Title:       "Hyphal Invasion",
			Description: "Each successful hit infects the target for 5 ticks. Infected enemies lose -5% accuracy and spread infection to nearby enemies within 100u if they die. Attack orders towards infected targets have 10% speed increase.",
			Path:        string(ships.Sporeform),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					StatusEffects: []StatusEffect{
						{
							Name:       "Infection",
							Duration:   5,
							EffectType: StatusInfection,
						},
					},
				},
			},
		},
		// 3. Digestive Mycelium: Heals from damaging infected targets
		{
			ID:          "sporeform_digestive_mycelium",
			Title:       "Digestive Mycelium",
			Description: "Damaging infected targets restores 2% hull per tick. Killing an infected target grants +10% damage for 3 ticks.",
			Path:        string(ships.Sporeform),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionTargetInfected, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{AtCombatRegenPct: 0.02},
					Duration:      1, // Per tick
				},
			},
		},
		// 4. Necrosporic Bloom: On-death spore cloud
		{
			ID:          "sporeform_necrosporic_bloom",
			Title:       "Necrosporic Bloom",
			Description: "Upon your destruction, releases a 300u spore cloud that infects all ships inside. Enemies take 5% HP damage over 3 ticks; allies gain +5% regen.",
			Path:        string(ships.Sporeform),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexOnDeath,
					Spawn: &SpawnEffect{
						SpawnType:   SpawnSporeCloud,
						SpawnRadius: 300,
					},
				},
			},
		},
		// 5. Mycelial Persistence: On-death respawn
		{
			ID:          "sporeform_mycelial_persistence",
			Title:       "Mycelial Persistence",
			Description: "When your stack is destroyed, 25% of your population reforms as a new micro-stack at your nearest controlled system after 10 ticks.",
			Path:        string(ships.Sporeform),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexOnDeath,
					Spawn: &SpawnEffect{
						SpawnType: SpawnMicroStack,
					},
					Duration: 10, // ticks delay
				},
			},
		},
	}

	// Tier 2: Cordyceps (Mind control and parasitic abilities)
	cordycepsNodes := []*BioNode{
		// 1. Neural Spores: Confusion on critical hits
		{
			ID:          "cordyceps_neural_spores",
			Title:       "Neural Spores",
			Description: "Critical hits apply a 'confusion' stack. At 3 stacks, target ships attack random nearby units for 1 tick.",
			Path:        string(ships.Cordyceps),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					StatusEffects: []StatusEffect{
						{
							Name:       "Confusion",
							MaxStacks:  3,
							EffectType: StatusConfusion,
						},
					},
				},
			},
		},
		// 2. Symbiotic Override: Infect and disable a building
		{
			ID:          "cordyceps_symbiotic_override",
			Title:       "Symbiotic Override",
			Description: "When engaging an enemy system, you can infect one building. For 5 ticks it doesn't operate.",
			Path:        string(ships.Cordyceps),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionBuildingInfected, CompareOp: CompareEqual, Value: true},
					},
					Duration: 5, // ticks
				},
			},
		},
		// 3. Spore Zombie: Spawn a husk from a dead infected stack
		{
			ID:          "cordyceps_spore_zombie",
			Title:       "Spore Zombie",
			Description: "When an infected stack dies within 200u, spawn a temporary 'Spore Husk' (20% of original HP, 50% damage) that lasts 3 ticks and fights for you.",
			Path:        string(ships.Cordyceps),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionTargetInfected, CompareOp: CompareEqual, Value: true},
					},
					Spawn: &SpawnEffect{
						SpawnType: SpawnSporeHusk,
						Duration:  3, // ticks
					},
				},
			},
		},
		// 4. Myco-Resonance: Evasion/cooldown reduction per infected target
		{
			ID:          "cordyceps_myco_resonance",
			Title:       "Myco-Resonance",
			Description: "For every infected target within 300u, gain +2% evasion and +2% cooldown reduction.",
			Path:        string(ships.Cordyceps),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionTargetInfected, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{
						EvasionPct:         0.02,
						AbilityCooldownPct: 0.02,
					},
				},
			},
		},
		// 5. Parasitic Singularity: Active ability, AoE damage and heal
		{
			ID:          "cordyceps_parasitic_singularity",
			Title:       "Parasitic Singularity",
			Description: "Activate to merge with all infected enemies within 400u, dealing massive AoE true damage and healing 50% of your max HP per infection consumed. 20 tick cooldown.",
			Path:        string(ships.Cordyceps),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAbilityUsed, CompareOp: CompareEqual, Value: true},
					},
					AoE: &AoETraitTarget{
						Radius:     400,
						TargetType: AoEEnemies,
					},
					Cooldown: 20, // ticks
				},
			},
		},
	}

	// Tier 3: Mycorrhiza (Networked symbiosis and support)
	mycorrhizaNodes := []*BioNode{
		// 1. Networked Roots: Links allies for shared regen and resistance
		{
			ID:          "mycorrhiza_networked_roots",
			Title:       "Networked Roots",
			Description: "All allied stacks within 300u are 'linked.' Linked ships share 10% of hull regen, gain +5% resistance to status effects and 10% movement speed.",
			Path:        string(ships.Mycorrhiza),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAllyInNetwork, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{
						AtCombatRegenPct: 0.1,
						SpeedDelta:       10,
					},
				},
			},
		},
		// 2. Spore Relay: Propagates ability effects through the network
		{
			ID:          "mycorrhiza_spore_relay",
			Title:       "Spore Relay",
			Description: "Abilities you cast propagate through the fungal network, applying 50% of their secondary effects (but not damage) to linked allies.",
			Path:        string(ships.Mycorrhiza),
		},
		// 3. Nutrient Exchange: Regen on ally kills or resource extraction
		{
			ID:          "mycorrhiza_nutrient_exchange",
			Title:       "Nutrient Exchange",
			Description: "When an ally in the network destroys a stack or extracts a resource, all linked units regain 5% energy and 3% HP.",
			Path:        string(ships.Mycorrhiza),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{
						HPPct: 0.03,
					},
				},
			},
		},
		// 4. Fungal Overgrowth: Creates a defensive field with 3+ allies
		{
			ID:          "mycorrhiza_fungal_overgrowth",
			Title:       "Fungal Overgrowth",
			Description: "When 3+ allies are within 300u, create an Overgrowth field for 5 ticks: +20% defense, -20% speed, and immunity to toxins.",
			Path:        string(ships.Mycorrhiza),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAllyCount, CompareOp: CompareGreaterEq, Value: 3},
					},
					Spawn: &SpawnEffect{
						SpawnType: SpawnOvergrowth,
						Duration:  5, // ticks
					},
				},
			},
		},
		// 5. Synaptic Bloom: Active ability, resets ally cooldowns
		{
			ID:          "mycorrhiza_synaptic_bloom",
			Title:       "Synaptic Bloom",
			Description: "Active ability — pulse through the network to instantly reset all allied ability cooldowns by 20% and cleanse one debuff. Costs 30% of your current HP.",
			Path:        string(ships.Mycorrhiza),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAbilityUsed, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{
						AbilityCooldownPct: -0.2,
					},
				},
			},
		},
	}

	tree.Tiers[0] = sporeformNodes
	tree.Tiers[1] = cordycepsNodes
	tree.Tiers[2] = mycorrhizaNodes

	return tree
}
