package ships

// initGlobalTree initializes the global formation mastery tree (applies to all formations).
func initGlobalTree() {
	FormationTreeCatalog[""] = FormationTree{
		Formation:   "",
		Name:        "Fleet Command Mastery",
		Description: "Universal skills that enhance all formations and improve strategic flexibility",
		MaxTier:     3,
		Nodes: []FormationTreeNode{
			// ====================
			// TIER 1: FOUNDATIONS
			// ====================
			{
				ID:           "global_tactical_awareness",
				Name:         "Tactical Awareness",
				Description:  "Fundamental understanding of fleet operations. +1 visibility to all ships regardless of position.",
				Formation:    "",
				Tier:         1,
				Row:          0,
				Column:       1,
				Cost:         NodeCost{ExperiencePoints: 1},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					GlobalMods: StatMods{
						VisibilityDelta: 1,
					},
				},
				Tags: []string{"vision", "starter"},
			},
			{
				ID:           "global_veteran_training",
				Name:         "Veteran Training",
				Description:  "Improved crew quality. All ships gain +5% HP and +3% accuracy.",
				Formation:    "",
				Tier:         1,
				Row:          0,
				Column:       5,
				Cost:         NodeCost{ExperiencePoints: 1},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					GlobalMods: StatMods{
						BucketHPPct: 0.05,
						AccuracyPct: 0.03,
					},
				},
				Tags: []string{"survivability", "starter"},
			},
			{
				ID:           "global_rapid_deployment",
				Name:         "Rapid Deployment",
				Description:  "Efficient logistics reduce formation reconfiguration time by 20%.",
				Formation:    "",
				Tier:         1,
				Row:          0,
				Column:       3,
				Cost:         NodeCost{ExperiencePoints: 1},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					ReconfigTimeMultiplier: -0.20,
				},
				Tags: []string{"flexibility", "starter"},
			},

			// ====================
			// TIER 2: SPECIALIZATION
			// ====================
			{
				ID:          "global_enhanced_communications",
				Name:        "Enhanced Communications",
				Description: "Superior command coordination. +10% composition bonus effectiveness, +1 range to abilities.",
				Formation:   "",
				Tier:        2,
				Row:         2,
				Column:      0,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"global_tactical_awareness"},
				},
				Effects: NodeEffects{
					CompositionBonusMultiplier: 0.10,
					GlobalMods: StatMods{
						AttackRangeDelta: 1,
					},
				},
				Tags: []string{"synergy", "range"},
			},
			{
				ID:          "global_strategic_vision",
				Name:        "Strategic Vision",
				Description: "See enemy fleet composition before battle. Reveals enemy formation type and loadout preview.",
				Formation:   "",
				Tier:        2,
				Row:         2,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"global_tactical_awareness"},
				},
				Effects: NodeEffects{
					CustomEffect: "reveal_enemy_composition",
					GlobalMods: StatMods{
						VisibilityDelta: 2,
					},
				},
				Tags: []string{"intel", "vision"},
			},
			{
				ID:          "global_superior_logistics",
				Name:        "Superior Logistics",
				Description: "Optimized supply chains. -5% upkeep for all ships, +10% transport capacity.",
				Formation:   "",
				Tier:        2,
				Row:         2,
				Column:      6,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"global_veteran_training"},
				},
				Effects: NodeEffects{
					GlobalMods: StatMods{
						UpkeepPct:            -0.05,
						TransportCapacityPct: 0.10,
					},
				},
				Tags: []string{"economy", "logistics"},
			},
			{
				ID:          "global_adaptive_tactics",
				Name:        "Adaptive Tactics",
				Description: "Flexible command structure. Can change ship positions within formation 1x per battle without reconfiguration time.",
				Formation:   "",
				Tier:        2,
				Row:         2,
				Column:      4,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"global_rapid_deployment"},
				},
				Effects: NodeEffects{
					CustomEffect:           "allow_position_swap_in_combat",
					ReconfigTimeMultiplier: -0.10, // Additional 10% reduction
				},
				Tags: []string{"flexibility", "combat"},
			},

			// ====================
			// TIER 3: MASTERY
			// ====================
			{
				ID:          "global_supreme_commander",
				Name:        "Supreme Commander",
				Description: "Master of all formations. All formation bonuses increased by 15%, all ships gain +5% damage.",
				Formation:   "",
				Tier:        3,
				Row:         4,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 8},
				Requirements: NodeRequirements{
					RequiredNodes:  []string{"global_enhanced_communications", "global_strategic_vision"},
					MinNodesInTree: 5,
				},
				Effects: NodeEffects{
					GlobalMods: StatMods{
						Damage: DamageMods{
							LaserPct:      0.05,
							NuclearPct:    0.05,
							AntimatterPct: 0.05,
						},
						FormationSyncBonus: 0.15,
					},
				},
				Tags: []string{"ultimate", "damage", "formation"},
			},
			{
				ID:          "global_versatile_genius",
				Name:        "Versatile Genius",
				Description: "Unprecedented tactical flexibility. Can run 2 formation types simultaneously (average their bonuses and counters).",
				Formation:   "",
				Tier:        3,
				Row:         4,
				Column:      5,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					RequiredNodes:     []string{"global_strategic_vision", "global_adaptive_tactics"},
					MinNodesInTree:    6,
					MutuallyExclusive: []string{"global_specialist"},
				},
				Effects: NodeEffects{
					CustomEffect: "dual_formation",
				},
				Tags: []string{"ultimate", "flexibility"},
			},
			{
				ID:          "global_specialist",
				Name:        "Formation Specialist",
				Description: "Choose 1 formation to master: +50% effectiveness to that formation, but -15% to all others. Pick carefully!",
				Formation:   "",
				Tier:        3,
				Row:         4,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 8},
				Requirements: NodeRequirements{
					RequiredNodes:     []string{"global_superior_logistics", "global_adaptive_tactics"},
					MinNodesInTree:    5,
					MutuallyExclusive: []string{"global_versatile_genius"},
				},
				Effects: NodeEffects{
					CustomEffect: "formation_specialist",
					CustomParams: map[string]interface{}{
						"bonus_multiplier":   0.50,
						"penalty_multiplier": -0.15,
					},
				},
				Tags: []string{"ultimate", "specialization"},
			},
		},
	}
}
