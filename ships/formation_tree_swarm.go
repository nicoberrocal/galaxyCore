package ships

// initSwarmTree initializes the Swarm Formation mastery tree.
func initSwarmTree() {
	FormationTreeCatalog[FormationSwarm] = FormationTree{
		Formation:   FormationSwarm,
		Name:        "Swarm Formation Mastery",
		Description: "Master dispersed anti-AoE tactics and overwhelming numbers",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			{
				ID:          "swarm_dispersal",
				Name:        "Dispersal Tactics",
				Description: "All positions gain +1 speed. -50% splash damage taken.",
				Formation:   FormationSwarm,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					FormationMods: StatMods{
						SpeedDelta: 1,
					},
					CustomEffect: "splash_reduction",
					CustomParams: map[string]interface{}{
						"reduction": 0.50,
					},
				},
				Tags: []string{"speed", "anti_aoe", "starter"},
			},
			{
				ID:          "swarm_overwhelming_numbers",
				Name:        "Overwhelming Numbers",
				Description: "+25% ship capacity in all positions. Each ship type over 20 units grants +5% damage.",
				Formation:   FormationSwarm,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"swarm_dispersal"},
				},
				Effects: NodeEffects{
					CustomEffect: "scale_with_numbers",
					CustomParams: map[string]interface{}{
						"capacity_bonus": 0.25,
						"threshold":      20,
						"damage_per_threshold": 0.05,
					},
				},
				Tags: []string{"scaling", "numbers"},
			},
			{
				ID:          "swarm_rapid_response",
				Name:        "Rapid Response",
				Description: "Swarm can split into 2 smaller swarms or merge instantly without reconfiguration time.",
				Formation:   FormationSwarm,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"swarm_dispersal"},
				},
				Effects: NodeEffects{
					CustomEffect: "split_merge",
					ReconfigTimeMultiplier: -0.50,
				},
				Tags: []string{"flexibility", "tactical"},
			},
			{
				ID:          "swarm_death_by_thousand_cuts",
				Name:        "Death by Thousand Cuts",
				Description: "Each ship in swarm applies 'Bleed' debuff (-2% HP per turn, stacks). Swarm vs Phalanx counter increased to 1.5x.",
				Formation:   FormationSwarm,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"swarm_overwhelming_numbers", "swarm_rapid_response"},
				},
				Effects: NodeEffects{
					CustomEffect: "bleed_stacks",
					CustomParams: map[string]interface{}{
						"bleed_percent": 0.02,
					},
					CounterBonusMultiplier: 0.30,
				},
				Tags: []string{"ultimate", "dot"},
			},
			{
				ID:          "swarm_locust_cloud",
				Name:        "Locust Cloud",
				Description: "Swarm becomes untargetable by single-target abilities. Each destroyed ship grants remaining ships +3% damage (stacking). Immune to splash damage.",
				Formation:   FormationSwarm,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes: []string{"swarm_death_by_thousand_cuts"},
				},
				Effects: NodeEffects{
					CustomEffect: "swarm_ascension",
					CustomParams: map[string]interface{}{
						"untargetable":     true,
						"splash_immunity":  true,
						"vengeful_damage": 0.03,
					},
				},
				Tags: []string{"ultimate", "unstoppable"},
			},
		},
	}
}
