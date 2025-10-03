package ships

// initVanguardTree initializes the Vanguard Formation mastery tree.
func initVanguardTree() {
	FormationTreeCatalog[FormationVanguard] = FormationTree{
		Formation:   FormationVanguard,
		Name:        "Vanguard Formation Mastery",
		Description: "Master aggressive forward deployment and alpha strike tactics",
		MaxTier:     4,
		Nodes: []FormationTreeNode{
			{
				ID:          "vanguard_aggressive_advance",
				Name:        "Aggressive Advance",
				Description: "Front position gains +20% damage, +10% first volley damage.",
				Formation:   FormationVanguard,
				Tier:        1,
				Row:         0,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 2},
				Requirements: NodeRequirements{},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							Damage: DamageMods{LaserPct: 0.20, NuclearPct: 0.20, AntimatterPct: 0.20},
							FirstVolleyPct: 0.10,
						},
					},
				},
				Tags: []string{"offense", "starter"},
			},
			{
				ID:          "vanguard_lightning_strike",
				Name:        "Lightning Strike",
				Description: "Formation reconfiguration time reduced by 40%. +1 speed to all ships.",
				Formation:   FormationVanguard,
				Tier:        2,
				Row:         2,
				Column:      1,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"vanguard_aggressive_advance"},
				},
				Effects: NodeEffects{
					ReconfigTimeMultiplier: -0.40,
					FormationMods: StatMods{SpeedDelta: 1},
				},
				Tags: []string{"speed", "flexibility"},
			},
			{
				ID:          "vanguard_overwhelming_force",
				Name:        "Overwhelming Force",
				Description: "Front position +30% shield pierce, +15% crit chance. Support position abilities have -30% cooldown.",
				Formation:   FormationVanguard,
				Tier:        2,
				Row:         2,
				Column:      3,
				Cost:        NodeCost{ExperiencePoints: 3},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"vanguard_aggressive_advance"},
				},
				Effects: NodeEffects{
					PositionMods: map[FormationPosition]StatMods{
						PositionFront: {
							ShieldPiercePct: 0.30,
							CritPct:         0.15,
						},
						PositionSupport: {
							AbilityCooldownPct: -0.30,
						},
					},
				},
				Tags: []string{"offense", "pierce"},
			},
			{
				ID:          "vanguard_shock_and_awe",
				Name:        "Shock and Awe",
				Description: "First attack in battle deals +50% damage. Vanguard vs Box counter increased to 1.6x.",
				Formation:   FormationVanguard,
				Tier:        3,
				Row:         4,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 5},
				Requirements: NodeRequirements{
					RequiredNodes: []string{"vanguard_lightning_strike", "vanguard_overwhelming_force"},
				},
				Effects: NodeEffects{
					CustomEffect: "first_strike_bonus",
					CustomParams: map[string]interface{}{
						"damage_multiplier": 0.50,
					},
					CounterBonusMultiplier: 0.30,
				},
				Tags: []string{"ultimate", "alpha"},
			},
			{
				ID:          "vanguard_unstoppable",
				Name:        "Unstoppable Assault",
				Description: "Vanguard cannot be slowed or have speed reduced. All damage +25%, -20% shields (high risk/reward).",
				Formation:   FormationVanguard,
				Tier:        4,
				Row:         6,
				Column:      2,
				Cost:        NodeCost{ExperiencePoints: 10},
				Requirements: NodeRequirements{
					MinNodesInTree: 6,
					RequiredNodes: []string{"vanguard_shock_and_awe"},
				},
				Effects: NodeEffects{
					CustomEffect: "unstoppable",
					FormationMods: StatMods{
						Damage: DamageMods{LaserPct: 0.25, NuclearPct: 0.25, AntimatterPct: 0.25},
						LaserShieldDelta:      -2,
						NuclearShieldDelta:    -2,
						AntimatterShieldDelta: -2,
					},
				},
				Tags: []string{"ultimate", "glass_cannon"},
			},
		},
	}
}
