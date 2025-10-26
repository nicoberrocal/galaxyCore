package essences

import "github.com/nicoberrocal/galaxyCore/ships"

// BuildAquatica constructs the Aquatica tree with node-level mutations but applies
// the principle that essences mutate chosen nodes (Tidal Surge, Thermocline Reflex, Leviathan/TotalCirculation).
func BuildAquatica() *BioTree {
	tree := &BioTree{
		Name:        string(Aquatica),
		Description: "Aquatic biology tree focusing on fluid adaptation, camouflage, and pack coordination",
		Tiers:       make([][]*BioNode, 3),
	}

	// Tier 1: Cephalopod (Camouflage and Adaptation)
	cephalopodNodes := []*BioNode{
		// 1. Elastic Burst: Can move while changing formation at 50% movement speed
		{
			ID:          "cephalopod_elastic_burst",
			Title:       "Elastic Burst",
			Description: "Can move while changing formation at 50% movement speed. Others cannot move or move very slowly.",
			Path:        string(ships.Cephalopod),
			Effect:      ships.StatMods{FormationSyncBonus: 0.5}, // +50% formation sync bonus (represents movement during formation change)
			StatDelta:   StatDelta{SpeedPercent: 0.5},            // +50% speed (but only during formation changes)
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnFormationChange,
					Conditions: []Condition{
						{ConditionType: ConditionFormationType, CompareOp: CompareEqual, Value: "changing"},
					},
					PrimaryEffect: &ships.StatMods{SpeedDelta: 50}, // Effective speed bonus during formation change
					Duration:      5,                               // ticks
					Cooldown:      0,                               // No cooldown
				},
			},
		},
		// 2. Ink Sac Propulsion: Gain 20% evasion for 5 ticks after activating an ability
		{
			ID:          "cephalopod_ink_sac",
			Title:       "Ink Sac Propulsion",
			Description: "Gain 20% evasion for 5 ticks after activating an ability. Effect can be active only once.",
			Path:        string(ships.Cephalopod),
			StatDelta:   StatDelta{EvasionPct: 0.2}, // +20% evasion
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAbilityCast,
					Conditions: []Condition{
						{ConditionType: ConditionAbilityUsed, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect:  &ships.StatMods{EvasionPct: 0.2}, // +20% evasion
					Duration:       5,                                // ticks
					Cooldown:       0,                                // Can only be active once per ability use
					MaxActivations: 1,
				},
			},
		},
		// 3. Decoy Drones: When attacked from behind, deploy a decoy drone
		{
			ID:          "cephalopod_decoy_drones",
			Title:       "Decoy Drones",
			Description: "When attacked from behind, deploy a decoy drone that draws fire and slows enemy upon destruction.",
			Path:        string(ships.Cephalopod),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAttackFromBehind,
					Conditions: []Condition{
						{ConditionType: ConditionAttackFromBehind, CompareOp: CompareEqual, Value: true},
					},
					Spawn: &SpawnEffect{
						SpawnType:  SpawnDecoyDrone,
						Duration:   3, // ticks
						SpawnCount: 1,
					},
					SecondaryEffect: &ships.StatMods{ // Slow effect when decoy is destroyed
						AttackIntervalPct: 0.5, // +50% attack interval (slower attacks)
					},
					Cooldown: 10, // ticks
				},
			},
		},
		// 4. Neurotoxin Coating: Critical hits apply toxin stacks, stun at 3 stacks
		{
			ID:          "cephalopod_neurotoxin",
			Title:       "Neurotoxin Coating",
			Description: "Your attacks apply a stack of toxin on critical hits. At 3 stacks the target is stunned for 1 tick.",
			Path:        string(ships.Cephalopod),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnCriticalHit,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					StatusEffects: []StatusEffect{
						{
							Name:        "Toxin",
							Description: "Stacking toxin that stuns at 3 stacks",
							Duration:    10, // ticks
							MaxStacks:   3,
							EffectType:  StatusInfection,
						},
					},
					Duration: 10, // ticks
				},
			},
		},
		// 5. Mimetic Genome: Copy enemy attack type for 3 ticks after engaging
		{
			ID:          "cephalopod_mimetic",
			Title:       "Mimetic Genome",
			Description: "Mimics the signature of the targeted stack, changing attack type to its weakest shield after engaging for 3 ticks.",
			Path:        string(ships.Cephalopod),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnCombatStart,
					Conditions: []Condition{
						{ConditionType: ConditionCombatState, CompareOp: CompareEqual, Value: "engaging"},
					},
					Duration: 3, // ticks
					Cooldown: 0, // Permanent until disengage
				},
			},
		},
	}

	// Tier 2: Chondrichthyan (Sensory and Hunting)
	chondrichthyanNodes := []*BioNode{
		// 1. Electroreceptive sense: +50% ping range, detect cloaked wounded ships
		{
			ID:          "chondrichthyan_electroreceptive",
			Title:       "Electroreceptive Sense",
			Description: "50% increased ping range and can detect cloaked wounded ships (<100% HP). Global detection below 20%.",
			Path:        string(ships.Chondrichthyan),
			Effect:      ships.StatMods{PingRangePct: 0.5, CloakDetect: true}, // +50% ping range, cloak detection
			StatDelta:   StatDelta{AccuracyPercent: 0.1},                      // +10% accuracy (represents better detection)
		},
		// 2. Ram ventilation: Stacking speed/damage bonus while sprinting toward enemy
		{
			ID:          "chondrichthyan_ram_ventilation",
			Title:       "Ram Ventilation",
			Description: "Gain a 2% stackable bonus to Speed and Damage for every tick sprinting towards an enemy with an active ability.",
			Path:        string(ships.Chondrichthyan),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnMovementSprint,
					Conditions: []Condition{
						{ConditionType: ConditionMovementState, CompareOp: CompareEqual, Value: "sprinting_toward_enemy"},
						{ConditionType: ConditionHasStatus, CompareOp: CompareEqual, Value: "active_ability"},
					},
					PrimaryEffect: &ships.StatMods{
						SpeedDelta: 2,                                                                       // +2 speed per tick
						Damage:     ships.DamageMods{LaserPct: 0.02, NuclearPct: 0.02, AntimatterPct: 0.02}, // +2% damage per tick
					},
					Duration: 1, // Per tick
				},
			},
		},
		// 3. Spiral Intestine: +50% speed for 5 ticks after destroying a stack
		{
			ID:          "chondrichthyan_spiral_intestine",
			Title:       "Spiral Intestine",
			Description: "Gain 50% speed for 5 ticks after destroying a stack.",
			Path:        string(ships.Chondrichthyan),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnKill,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{SpeedDelta: 50}, // +50 speed
					Duration:      5,                               // ticks
					Cooldown:      0,                               // No cooldown
				},
			},
		},
		// 4. Looming Presence: 4th consecutive attack has 25% crit chance
		{
			ID:          "chondrichthyan_looming_presence",
			Title:       "Looming Presence",
			Description: "The fourth consecutive attack on a target has 25% crit chance.",
			Path:        string(ships.Chondrichthyan),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnConsecutiveAttacks,
					Conditions: []Condition{
						{ConditionType: ConditionConsecutiveAttacks, CompareOp: CompareEqual, Value: 4},
					},
					PrimaryEffect: &ships.StatMods{CritPct: 0.25}, // +25% crit chance
					Duration:      1,                              // One attack
				},
			},
		},
		// 5. Frenzy Scent: +25% damage globally after losing a system for 10 ticks
		{
			ID:          "chondrichthyan_frenzy_scent",
			Title:       "Frenzy Scent",
			Description: "Your damage is increased globally by 25% after losing a system. Lasts 10 ticks.",
			Path:        string(ships.Chondrichthyan),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnSystemLost,
					Conditions: []Condition{
						{ConditionType: ConditionSystemLost, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{
						Damage: ships.DamageMods{LaserPct: 0.25, NuclearPct: 0.25, AntimatterPct: 0.25}, // +25% damage
					},
					Duration: 10, // ticks
					Cooldown: 0,  // No cooldown
				},
			},
		},
	}

	// Tier 3: Cetacean (Social and Communication)
	cetaceanNodes := []*BioNode{
		// 1. Echolocative Sonar: +5% accuracy to self and allies within 200u, minimum 75% accuracy
		{
			ID:          "cetacean_echolocative",
			Title:       "Echolocative Sonar",
			Description: "You and all allies within 200u gain +5% accuracy and cannot have their accuracy reduced below 75%.",
			Path:        string(ships.Cetacean),
			Effect:      ships.StatMods{AccuracyPct: 0.05}, // +5% accuracy
			StatDelta:   StatDelta{AccuracyPercent: 0.05},  // +5% accuracy
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAllyNearby,
					Conditions: []Condition{
						{ConditionType: ConditionAllyNearby, CompareOp: CompareEqual, Value: 200},
						{ConditionType: ConditionAllyCount, CompareOp: CompareGreater, Value: 0},
					},
					AoE: &AoETraitTarget{
						Radius:     200,
						TargetType: AoEAllies,
						Origin:     AoESelf,
					},
					PrimaryEffect: &ships.StatMods{AccuracyPct: 0.05}, // +5% accuracy to allies
					Duration:      0,                                  // Permanent
				},
			},
		},
		// 2. Pod Synchronization: Each ally attacking same target reduces shields by 2%, stack to 10%
		{
			ID:          "cetacean_pod_synchronization",
			Title:       "Pod Synchronization",
			Description: "Each allied stack attacking the same target reduces its shields by 2%, stacking up to 10%.",
			Path:        string(ships.Cetacean),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAllyNearby,
					Conditions: []Condition{
						{ConditionType: ConditionAllyCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{
						LaserShieldDelta:      -2, // -2 laser shield per ally
						NuclearShieldDelta:    -2, // -2 nuclear shield per ally
						AntimatterShieldDelta: -2, // -2 antimatter shield per ally
					},
					Duration: 0, // Permanent while condition is met
				},
			},
		},
		// 3. Bio Acoustic field: +20% speed to allied moving stacks in same formation within 500u
		{
			ID:          "cetacean_bio_acoustic",
			Title:       "Bio Acoustic Field",
			Description: "All ally moving stacks within 500u running the same formation gain +20% speed bonus.",
			Path:        string(ships.Cetacean),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAllyNearby,
					Conditions: []Condition{
						{ConditionType: ConditionFormationType, CompareOp: CompareEqual, Value: "same"},
						{ConditionType: ConditionMovementState, CompareOp: CompareEqual, Value: "moving"},
						{ConditionType: ConditionAllyNearby, CompareOp: CompareLessEq, Value: 500},
					},
					AoE: &AoETraitTarget{
						Radius:     500,
						TargetType: AoEAllies,
						Origin:     AoESelf,
					},
					PrimaryEffect: &ships.StatMods{SpeedDelta: 20}, // +20 speed
					Duration:      0,                               // While conditions are met
				},
			},
		},
		// 4. Communal Lungs: Share 15% healing and 10% energy regen with lowest HP ally within 250u
		{
			ID:          "cetacean_communal_lungs",
			Title:       "Communal Lungs",
			Description: "Shares 15% of healing and 10% of energy regen to the lowest HP ally within 250u.",
			Path:        string(ships.Cetacean),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAllyNearby,
					Conditions: []Condition{
						{ConditionType: ConditionAllyNearby, CompareOp: CompareLessEq, Value: 250},
						{ConditionType: ConditionHPPercent, CompareOp: CompareLess, Value: 1.0}, // Ally has lower HP
					},
					AoE: &AoETraitTarget{
						Radius:     250,
						TargetType: AoEAllies,
						MaxTargets: 1, // Only the lowest HP ally
					},
					PrimaryEffect: &ships.StatMods{
						AtCombatRegenPct: 0.15, // +15% combat regen (healing)
					},
					Duration: 0, // Permanent while conditions met
				},
			},
		},
		// 5. Large Migration: +25% shields while 800u+ from friendly system with allies nearby
		{
			ID:          "cetacean_large_migration",
			Title:       "Large Migration",
			Description: "While 800u+ from any friendly system, gain +25% shields if at least one ally is nearby.",
			Path:        string(ships.Cetacean),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAllyNearby,
					Conditions: []Condition{
						{ConditionType: ConditionDistance, CompareOp: CompareGreater, Value: 800},
						{ConditionType: ConditionAllyCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{
						LaserShieldDelta:      25, // +25 laser shield
						NuclearShieldDelta:    25, // +25 nuclear shield
						AntimatterShieldDelta: 25, // +25 antimatter shield
					},
					Duration: 0, // While conditions met
				},
			},
		},
	}

	tree.Tiers[0] = cephalopodNodes
	tree.Tiers[1] = chondrichthyanNodes
	tree.Tiers[2] = cetaceanNodes

	return tree
}
