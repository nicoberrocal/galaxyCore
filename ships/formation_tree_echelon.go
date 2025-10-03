package ships

// initEchelonTree initializes the Echelon Formation mastery tree.
func initEchelonTree() {
	FormationTreeCatalog[FormationEchelon] = FormationTree{
		Formation:   FormationEchelon,
		Name:        "Echelon Formation Mastery",
		Description: "Master diagonal staggered lines and asymmetric warfare",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			{
				ID:          "echelon_staggered_advance",
				Name:        "Staggered Advance",
				Description: "Front and flank positions gain +1 shield, +10% damage.",
				Formation:   FormationEchelon,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							LaserShieldDelta: 1,
							Damage:           DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
						},
						PositionFlank: {
							SpeedDelta: 1,
							CritPct:    0.08,
						},
					},
				},
				Tags: []string{"balanced", "starter"},
			},
			{
				ID:          "echelon_crossfire",
				Name:        "Crossfire",
				Description: "When both flank and back positions attack same target, +25% damage bonus.",
				Formation:   FormationEchelon,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"echelon_staggered_advance"},
				},
				Effects: NodeEffects{
					CustomEffect: "coordinated_fire",
					CustomParams: map[string]interface{}{
						"positions":         []string{"flank", "back"},
						"damage_multiplier": 0.25,
					},
				},
				Tags: []string{"synergy", "damage"},
			},
			{
				ID:          "echelon_asymmetric_defense",
				Name:        "Asymmetric Defense",
				Description: "When attacked from one side, opposite side gains +3 shields and +20% damage for 1 turn.",
				Formation:   FormationEchelon,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"echelon_staggered_advance"},
				},
				Effects: NodeEffects{
					CustomEffect: "reactive_defense",
				},
				Tags: []string{"defense", "reactive"},
			},
			{
				ID:          "echelon_concentrated_fire",
				Name:        "Concentrated Fire",
				Description: "Back position +2 range, +15% accuracy. All positions +12% damage vs same target.",
				Formation:   FormationEchelon,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"echelon_crossfire", "echelon_asymmetric_defense"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionBack: {
							AttackRangeDelta: 2,
							AccuracyPct:      0.15,
						},
					},
					FormationMods: StatMods{
						Damage: DamageMods{LaserPct: 0.12, NuclearPct: 0.12, AntimatterPct: 0.12},
					},
				},
				Tags: []string{"ultimate", "focus_fire"},
			},
			{
				ID:          "echelon_perfect_angles",
				Name:        "Perfect Angles",
				Description: "Echelon formation gains optimal positioning: +30% damage, +20% evasion, immune to flank attacks.",
				Formation:   FormationEchelon,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes: []string{"echelon_concentrated_fire"},
				},
				Effects: NodeEffects{
					CustomEffect: "flank_immunity",
					FormationMods: StatMods{
						Damage:     DamageMods{LaserPct: 0.30, NuclearPct: 0.30, AntimatterPct: 0.30},
						EvasionPct: 0.20,
					},
				},
				Tags: []string{"ultimate", "positioning"},
			},
		},
	}
}
