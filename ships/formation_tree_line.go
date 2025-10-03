package ships

// initLineTree initializes the Line Formation mastery tree.
func initLineTree() {
	FormationTreeCatalog[FormationLine] = FormationTree{
		Formation:   FormationLine,
		Name:        "Line Formation Mastery",
		Description: "Master the balanced line formation with specialized front-line and back-line tactics",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			// ==================
			// TIER 1: BASIC TRAINING
			// ==================
			{
				ID:           "line_defensive_stance",
				Name:         "Defensive Stance",
				Description:  "Front position ships adopt defensive posture. +2 to all shields in front position.",
				Formation:    FormationLine,
				Tier:         1,
				Row:          0,
				Column:       0,
				Cost:         NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							LaserShieldDelta:      2,
							NuclearShieldDelta:    2,
							AntimatterShieldDelta: 2,
						},
					},
				},
				Tags: []string{"defense", "front", "starter"},
			},
			{
				ID:          "line_offensive_posture",
				Name:        "Offensive Posture",
				Description: "Front position prioritizes firepower. +15% damage but -1 to all shields in front position.",
				Formation:   FormationLine,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{
					MutuallyExclusive: []string{"line_defensive_stance"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							Damage: DamageMods{
								LaserPct:      0.15,
								NuclearPct:    0.15,
								AntimatterPct: 0.15,
							},
							LaserShieldDelta:      -1,
							NuclearShieldDelta:    -1,
							AntimatterShieldDelta: -1,
						},
					},
				},
				Tags: []string{"offense", "front", "starter"},
			},
			{
				ID:          "line_balanced_deployment",
				Name:        "Balanced Deployment",
				Description: "Versatile approach. All positions in Line formation gain +5% HP and +5% damage.",
				Formation:   FormationLine,
				Tier:        1,
				Row:         0,
				Column:      4,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{
					MutuallyExclusive: []string{"line_defensive_stance", "line_offensive_posture"},
				},
				Effects: NodeEffects{
					FormationMods: StatMods{
						BucketHPPct: 0.05,
						Damage: DamageMods{
							LaserPct:      0.05,
							NuclearPct:    0.05,
							AntimatterPct: 0.05,
						},
					},
				},
				Tags: []string{"balanced", "starter"},
			},

			// ==================
			// TIER 2: SPECIALIZATION
			// ==================
			{
				ID:          "line_long_range_barrage",
				Name:        "Long-Range Barrage",
				Description: "Back position optimized for fire support. +2 range, +10% accuracy in back position.",
				Formation:   FormationLine,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{}, // Any tier 1
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionBack: {
							AttackRangeDelta: 2,
							AccuracyPct:      0.10,
						},
					},
				},
				Tags: []string{"back", "range", "accuracy"},
			},
			{
				ID:          "line_shield_wall",
				Name:        "Shield Wall",
				Description: "Front position forms impenetrable wall. -10% damage taken, adjacent ships gain +1 shield.",
				Formation:   FormationLine,
				Tier:        2,
				Row:         2,
				Column:      0,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_defensive_stance"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							LaserShieldDelta:      1,
							NuclearShieldDelta:    1,
							AntimatterShieldDelta: 1,
						},
					},
					CustomEffect: "shield_wall_aura",
					CustomParams: map[string]interface{}{
						"damage_reduction":      0.10,
						"adjacent_shield_bonus": 1,
					},
				},
				Tags: []string{"defense", "front", "aura"},
			},
			{
				ID:          "line_breakthrough_assault",
				Name:        "Breakthrough Assault",
				Description: "Front position gains offensive edge. +15% shield pierce, +10% first volley damage.",
				Formation:   FormationLine,
				Tier:        2,
				Row:         2,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_offensive_posture"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							ShieldPiercePct: 0.15,
							FirstVolleyPct:  0.10,
						},
					},
				},
				Tags: []string{"offense", "front", "pierce"},
			},
			{
				ID:          "line_flexible_reserves",
				Name:        "Flexible Reserves",
				Description: "Can swap front/back positions once per battle without reconfiguration time.",
				Formation:   FormationLine,
				Tier:        2,
				Row:         2,
				Column:      4,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_balanced_deployment"},
				},
				Effects: NodeEffects{
					CustomEffect:           "position_swap_front_back",
					ReconfigTimeMultiplier: -0.15,
				},
				Tags: []string{"flexibility", "tactical"},
			},

			// ==================
			// TIER 3: ADVANCED TACTICS
			// ==================
			{
				ID:          "line_enfilade_fire",
				Name:        "Enfilade Fire",
				Description: "Back position exploits line-of-fire. +20% damage vs targets aligned with back position. Line vs Vanguard counter increased to 1.5x.",
				Formation:   FormationLine,
				Tier:        3,
				Row:         4,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_long_range_barrage"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionBack: {
							Damage: DamageMods{
								LaserPct:      0.20,
								NuclearPct:    0.20,
								AntimatterPct: 0.20,
							},
						},
					},
					CounterBonusMultiplier: 0.15, // 1.3 → 1.5 vs Vanguard
				},
				Tags: []string{"back", "damage", "counter"},
			},
			{
				ID:          "line_unbreakable_line",
				Name:        "Unbreakable Line",
				Description: "Front position becomes immovable. When front ships are above 50% HP, entire formation gains +10% evasion. Front takes -20% damage.",
				Formation:   FormationLine,
				Tier:        3,
				Row:         4,
				Column:      0,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_shield_wall"},
				},
				Effects: NodeEffects{
					CustomEffect: "conditional_formation_buff",
					CustomParams: map[string]interface{}{
						"condition":     "front_hp_above_50",
						"evasion_bonus": 0.10,
					},
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							BucketHPPct: 0.10,
						},
					},
				},
				Tags: []string{"defense", "front", "conditional"},
			},
			{
				ID:          "line_hammer_and_anvil",
				Name:        "Hammer and Anvil",
				Description: "Classic combined arms. Front position applies 'Stunned' debuff (1 turn), back position deals +30% damage to stunned targets.",
				Formation:   FormationLine,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_breakthrough_assault"},
				},
				Effects: NodeEffects{
					CustomEffect: "hammer_anvil_combo",
					CustomParams: map[string]interface{}{
						"stun_duration":    1,
						"bonus_vs_stunned": 0.30,
					},
				},
				Tags: []string{"combo", "front", "back"},
			},
			{
				ID:          "line_adaptive_formation",
				Name:        "Adaptive Line",
				Description: "Can use Line formation bonuses with any formation type (but keep Line's counter vulnerabilities).",
				Formation:   FormationLine,
				Tier:        3,
				Row:         4,
				Column:      4,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"line_flexible_reserves"},
				},
				Effects: NodeEffects{
					CustomEffect: "carry_line_bonuses",
				},
				Tags: []string{"flexibility", "hybrid"},
			},

			// ==================
			// TIER 4: MASTERY (Choose 1)
			// ==================
			{
				ID:          "line_supreme_overlord",
				Name:        "Supreme Overlord of the Line",
				Description: "Ultimate mastery. Line formation cannot be countered. All Line bonuses +25%. Unlock ability: 'Perfect Line' (freeze enemy formation for 3 turns, once per battle).",
				Formation:   FormationLine,
				Tier:        4,
				Row:         6,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 8,
					RequiredNodes:  []string{"line_enfilade_fire", "line_unbreakable_line"},
				},
				Effects: NodeEffects{
					CounterResistMultiplier: 1.0, // Immune to counters
					FormationMods: StatMods{
						FormationSyncBonus: 0.25,
					},
					UnlocksAbility: AbilityID("PerfectLine"),
				},
				Tags: []string{"ultimate", "immunity", "ability"},
			},
			{
				ID:          "line_master_of_combined_arms",
				Name:        "Master of Combined Arms",
				Description: "Composition bonuses doubled in Line formation. Can field 8 ship types. Back position abilities have 2x range.",
				Formation:   FormationLine,
				Tier:        4,
				Row:         6,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 8},
				Requirements: NodeRequirements{
					MinNodesInTree: 8,
					RequiredNodes:  []string{"line_hammer_and_anvil"},
				},
				Effects: NodeEffects{
					CompositionBonusMultiplier: 1.0, // Double composition bonuses
					CustomEffect:               "expanded_fleet_cap",
					PositionMods: map[FormationPosition]StatMods{
						PositionBack: {
							AttackRangeDelta: 2,
						},
					},
				},
				Tags: []string{"ultimate", "composition", "synergy"},
			},
		},
	}
}
