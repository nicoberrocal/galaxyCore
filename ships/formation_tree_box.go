package ships

// initBoxTree initializes the Box Formation mastery tree.
func initBoxTree() {
	FormationTreeCatalog[FormationBox] = FormationTree{
		Formation:   FormationBox,
		Name:        "Box Formation Mastery",
		Description: "Master defensive all-around protection and siege resistance",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			// Tier 1
			{
				ID:           "box_all_around_defense",
				Name:         "All-Around Defense",
				Description:  "All positions gain +1 to all shields.",
				Formation:    FormationBox,
				Tier:         1,
				Row:          0,
				Column:       2,
				Cost:         NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					FormationMods: StatMods{
						LaserShieldDelta:      1,
						NuclearShieldDelta:    1,
						AntimatterShieldDelta: 1,
					},
				},
				Tags: []string{"defense", "starter"},
			},
			// Tier 2
			{
				ID:          "box_siege_resistant",
				Name:        "Siege Resistant",
				Description: "-30% damage from Nuclear attacks, +50% resistance to structure damage bonuses.",
				Formation:   FormationBox,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"box_all_around_defense"},
				},
				Effects: NodeEffects{
					CustomEffect: "siege_resistance",
					FormationMods: StatMods{
						NuclearShieldDelta: 2,
					},
				},
				Tags: []string{"defense", "siege"},
			},
			{
				ID:          "box_even_distribution",
				Name:        "Even Distribution",
				Description: "Damage is distributed perfectly evenly across all positions. +10% HP to all ships.",
				Formation:   FormationBox,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"box_all_around_defense"},
				},
				Effects: NodeEffects{
					CustomEffect: "perfect_distribution",
					FormationMods: StatMods{
						BucketHPPct: 0.10,
					},
				},
				Tags: []string{"defense", "hp"},
			},
			// Tier 3
			{
				ID:          "box_fortress",
				Name:        "Mobile Fortress",
				Description: "All positions gain +2 shields, +15% HP. Box counters Vanguard with 1.4x damage.",
				Formation:   FormationBox,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"box_siege_resistant", "box_even_distribution"},
				},
				Effects: NodeEffects{
					FormationMods: StatMods{
						LaserShieldDelta:      2,
						NuclearShieldDelta:    2,
						AntimatterShieldDelta: 2,
						BucketHPPct:           0.15,
					},
					CounterBonusMultiplier: 0.10,
				},
				Tags: []string{"ultimate", "defense"},
			},
			// Tier 4
			{
				ID:          "box_impregnable",
				Name:        "Impregnable Defense",
				Description: "Box formation takes -40% damage from all sources. Cannot be one-shot.",
				Formation:   FormationBox,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes:  []string{"box_fortress"},
				},
				Effects: NodeEffects{
					CustomEffect: "damage_reduction_global",
					CustomParams: map[string]interface{}{
						"reduction":       0.40,
						"prevent_oneshot": true,
					},
				},
				Tags: []string{"ultimate", "defense"},
			},
		},
	}
}
