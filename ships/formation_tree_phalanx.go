package ships

// initPhalanxTree initializes the Phalanx Formation mastery tree.
func initPhalanxTree() {
	FormationTreeCatalog[FormationPhalanx] = FormationTree{
		Formation:   FormationPhalanx,
		Name:        "Phalanx Formation Mastery",
		Description: "Master heavy frontal concentration and overwhelming frontal assault",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			{
				ID:          "phalanx_frontal_fortress",
				Name:        "Frontal Fortress",
				Description: "Front position gains +3 shields, +15% HP, +10% damage.",
				Formation:   FormationPhalanx,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							LaserShieldDelta:      3,
							NuclearShieldDelta:    3,
							AntimatterShieldDelta: 1,
							BucketHPPct:           0.15,
							Damage:                DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
						},
					},
				},
				Tags: []string{"front", "starter"},
			},
			{
				ID:          "phalanx_shield_bash",
				Name:        "Shield Bash",
				Description: "Front position has 20% chance to stun enemies for 1 turn on hit.",
				Formation:   FormationPhalanx,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"phalanx_frontal_fortress"},
				},
				Effects: NodeEffects{
					CustomEffect: "stun_on_hit",
					CustomParams: map[string]interface{}{
						"chance":   0.20,
						"duration": 1,
					},
				},
				Tags: []string{"cc", "front"},
			},
			{
				ID:          "phalanx_extended_line",
				Name:        "Extended Line",
				Description: "Back position +3 range. Front position can field 50% more ships.",
				Formation:   FormationPhalanx,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"phalanx_frontal_fortress"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionBack: {
							AttackRangeDelta: 3,
						},
					},
					CustomEffect: "expanded_front_capacity",
				},
				Tags: []string{"range", "capacity"},
			},
			{
				ID:          "phalanx_iron_wall",
				Name:        "Iron Wall",
				Description: "Front position -30% incoming damage. When front holds, back position +40% damage.",
				Formation:   FormationPhalanx,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"phalanx_shield_bash", "phalanx_extended_line"},
				},
				Effects: NodeEffects{
					CustomEffect: "conditional_back_bonus",
					CustomParams: map[string]interface{}{
						"front_damage_reduction": 0.30,
						"back_damage_bonus":      0.40,
					},
				},
				Tags: []string{"ultimate", "defense"},
			},
			{
				ID:          "phalanx_unbreakable",
				Name:        "Unbreakable Phalanx",
				Description: "Front position cannot be destroyed while any ship remains. When flanked, all positions gain front position bonuses.",
				Formation:   FormationPhalanx,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes: []string{"phalanx_iron_wall"},
				},
				Effects: NodeEffects{
					CustomEffect: "last_stand",
					CustomParams: map[string]interface{}{
						"redistribute_damage": true,
					},
				},
				Tags: []string{"ultimate", "immortal"},
			},
		},
	}
}
