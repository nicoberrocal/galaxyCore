package ships

// initSkirmishTree initializes the Skirmish Formation mastery tree.
func initSkirmishTree() {
	FormationTreeCatalog[FormationSkirmish] = FormationTree{
		Formation:   FormationSkirmish,
		Name:        "Skirmish Formation Mastery",
		Description: "Master mobile flanking tactics and hit-and-run warfare",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			{
				ID:          "skirmish_mobile_warfare",
				Name:        "Mobile Warfare",
				Description: "Flank position gains +2 speed, +10% accuracy, +15% damage.",
				Formation:   FormationSkirmish,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFlank: {
							SpeedDelta:  2,
							AccuracyPct: 0.10,
							Damage:      DamageMods{LaserPct: 0.15, NuclearPct: 0.15, AntimatterPct: 0.15},
						},
					},
				},
				Tags: []string{"speed", "starter"},
			},
			{
				ID:          "skirmish_hit_and_run",
				Name:        "Hit and Run",
				Description: "After attacking, Skirmish formation can disengage without counterattack (25% chance).",
				Formation:   FormationSkirmish,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"skirmish_mobile_warfare"},
				},
				Effects: NodeEffects{
					CustomEffect: "disengage_chance",
					CustomParams: map[string]interface{}{
						"chance": 0.25,
					},
				},
				Tags: []string{"tactical", "evasion"},
			},
			{
				ID:          "skirmish_flanking_mastery",
				Name:        "Flanking Mastery",
				Description: "Flank position deals +40% damage when attacking from the side or rear.",
				Formation:   FormationSkirmish,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"skirmish_mobile_warfare"},
				},
				Effects: NodeEffects{
					CustomEffect: "flanking_bonus",
					CustomParams: map[string]interface{}{
						"damage_multiplier": 0.40,
					},
				},
				Tags: []string{"damage", "tactical"},
			},
			{
				ID:          "skirmish_guerrilla_tactics",
				Name:        "Guerrilla Tactics",
				Description: "+30% evasion for entire formation. Skirmish vs Phalanx counter increased to 1.5x.",
				Formation:   FormationSkirmish,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"skirmish_hit_and_run", "skirmish_flanking_mastery"},
				},
				Effects: NodeEffects{
					FormationMods: StatMods{
						EvasionPct: 0.30,
					},
					CounterBonusMultiplier: 0.20,
				},
				Tags: []string{"ultimate", "evasion"},
			},
			{
				ID:          "skirmish_phantom_fleet",
				Name:        "Phantom Fleet",
				Description: "Skirmish formation creates 2 decoy formations that enemies must identify (visual only). +50% evasion, cannot be targeted by abilities.",
				Formation:   FormationSkirmish,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes: []string{"skirmish_guerrilla_tactics"},
				},
				Effects: NodeEffects{
					CustomEffect: "phantom_decoys",
					CustomParams: map[string]interface{}{
						"decoy_count": 2,
					},
					FormationMods: StatMods{
						EvasionPct: 0.50,
					},
				},
				Tags: []string{"ultimate", "stealth"},
			},
		},
	}
}
