package ships

// ShipBlueprints enumerates baseline ship TYPE definitions with recommended
// abilities and default RoleMode. These are data-only; runtime state lives
// elsewhere (see stack.go for HP buckets, movement, etc.).
//
// Notes
//   - Numbers are conservative baselines based on our prior design. Tweak freely.
//   - Abilities are limited to 3 per type. Some types have additional synergies
//     available via runes (see runes.go) and RoleModes (see roles.go).
//   - AttackType is one of: "Laser", "Nuclear", "Antimatter".
var ShipBlueprints = map[ShipType]Ship{
	// Economic unit: anchors on asteroids/nebulas to gather.
	Drone: {
		ShipType:         string(Drone),
		AttackType:       "Laser",
		LaserShield:      2,
		NuclearShield:    1,
		AntimatterShield: 0,
		Speed:            4,
		VisibilityRange:  3,
		AttackRange:      100,
		HP:               100,
		AttackDamage:     8,
		AttackInterval:   3.0,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityResourceHarvester],
			AbilitiesCatalog[AbilitySelfRepair],
			AbilitiesCatalog[AbilityCloakWhileAnchored],
		},
		MetalCost:         20,
		CrystalCost:       10,
		PlasmaCost:        0,
		TransportCapacity: 0,
		CanTransport:      nil,
	},

	// Recon/light skirmisher with strong sensors and deception tools.
	Scout: {
		ShipType:         string(Scout),
		AttackType:       "Laser",
		LaserShield:      2,
		NuclearShield:    0,
		AntimatterShield: 1,
		Speed:            9,
		VisibilityRange:  10,
		AttackRange:      200,
		HP:               100,
		AttackDamage:     12,
		AttackInterval:   1.5,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLongRangeSensors],
			AbilitiesCatalog[AbilityPing],
			AbilitiesCatalog[AbilityDecoyBeacon],
		},
		MetalCost:   50,
		CrystalCost: 30,
		PlasmaCost:  0,
	},

	// Versatile backbone fighter. Adapts to shields and focuses targets.
	Fighter: {
		ShipType:         string(Fighter),
		AttackType:       "Laser",
		LaserShield:      3,
		NuclearShield:    1,
		AntimatterShield: 1,
		Speed:            6,
		VisibilityRange:  5,
		AttackRange:      350,
		HP:               200,
		AttackDamage:     28,
		AttackInterval:   1.0,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityAdaptiveTargeting],
			AbilitiesCatalog[AbilityFocusFire],
			AbilitiesCatalog[AbilityEvasiveManeuvers],
		},
		MetalCost:   80,
		CrystalCost: 50,
		PlasmaCost:  20,
	},

	// Siege platform with very long range and structure damage bonuses.
	Bomber: {
		ShipType:         string(Bomber),
		AttackType:       "Nuclear",
		LaserShield:      1,
		NuclearShield:    3,
		AntimatterShield: 2,
		Speed:            5,
		VisibilityRange:  6,
		AttackRange:      700,
		HP:               500,
		AttackDamage:     95,
		AttackInterval:   3.0,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLightSpeed],
			AbilitiesCatalog[AbilitySiegePayload],
			AbilitiesCatalog[AbilityStandoffPattern],
		},
		MetalCost:   600,
		CrystalCost: 500,
		PlasmaCost:  200,
	},

	// Mobile hub with bays for escorts. Tanky and defensive utility.
	Carrier: {
		ShipType:         string(Carrier),
		AttackType:       "Nuclear",
		LaserShield:      2,
		NuclearShield:    4,
		AntimatterShield: 3,
		Speed:            4,
		VisibilityRange:  6,
		AttackRange:      470,
		HP:               900,
		AttackDamage:     75,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLightSpeed],
			AbilitiesCatalog[AbilityHangarLaunch],
			AbilitiesCatalog[AbilityPointDefenseScreen],
		},
		MetalCost:         800,
		CrystalCost:       600,
		PlasmaCost:        400,
		TransportCapacity: 40,
		CanTransport:      []string{"drone", "scout", "fighter"},
	},

	// Burst striker with interdiction and lightspeed.
	Destroyer: {
		ShipType:         string(Destroyer),
		AttackType:       "Antimatter",
		LaserShield:      2,
		NuclearShield:    2,
		AntimatterShield: 4,
		Speed:            6,
		VisibilityRange:  6,
		AttackRange:      500,
		HP:               600,
		AttackDamage:     110,
		AttackInterval:   2.5,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLightSpeed],
			AbilitiesCatalog[AbilityAlphaStrike],
			AbilitiesCatalog[AbilityInterdictorPulse],
		},
		MetalCost:   500,
		CrystalCost: 700,
		PlasmaCost:  400,
	},

	// Medium tank/brawler. Fills frontline gap with sustained presence.
	// Combos: Shield Overcharge + Box formation = extreme tankiness
	//         Ramming Speed + Vanguard = aggressive charge
	Cruiser: {
		ShipType:         string(Cruiser),
		AttackType:       "Nuclear",
		LaserShield:      3,
		NuclearShield:    3,
		AntimatterShield: 2,
		Speed:            5,
		VisibilityRange:  5,
		AttackRange:      375,
		HP:               400,
		AttackDamage:     55,
		AttackInterval:   1.8,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLightSpeed],
			AbilitiesCatalog[AbilityShieldOvercharge],
			AbilitiesCatalog[AbilityRammingSpeed],
		},
		MetalCost:   300,
		CrystalCost: 200,
		PlasmaCost:  100,
	},

	// Fast Antimatter attacker/pursuit specialist. Counters Scout swarms.
	// Combos: Antimatter Burst + Ping = devastating single-target nuke
	//         Pursuit Protocol + Target Lock = scout hunter
	//         Target Lock + Interdictor Pulse = warp denial combo
	Corvette: {
		ShipType:         string(Corvette),
		AttackType:       "Antimatter",
		LaserShield:      1,
		NuclearShield:    1,
		AntimatterShield: 3,
		Speed:            8,
		VisibilityRange:  6,
		AttackRange:      290,
		HP:               150,
		AttackDamage:     32,
		AttackInterval:   1.2,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityPursuitProtocol],
			AbilitiesCatalog[AbilityAntimatterBurst],
			AbilitiesCatalog[AbilityTargetLock],
		},
		MetalCost:   100,
		CrystalCost: 150,
		PlasmaCost:  80,
	},

	// AoE specialist/swarm breaker. Punishes tight formations.
	// Combos: Cluster Munitions + Barrage Mode = massive AoE coverage
	//         Suppressive Fire + Standoff Pattern = area lockdown
	//         Cluster + Swarm formation counter = anti-cluster defense
	Ballista: {
		ShipType:         string(Ballista),
		AttackType:       "Nuclear",
		LaserShield:      2,
		NuclearShield:    4,
		AntimatterShield: 1,
		Speed:            3,
		VisibilityRange:  7,
		AttackRange:      900,
		HP:               350,
		AttackDamage:     85,
		AttackInterval:   3.5,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityLightSpeed],
			AbilitiesCatalog[AbilityClusterMunitions],
			AbilitiesCatalog[AbilityBarrageMode],
		},
		MetalCost:   700,
		CrystalCost: 400,
		PlasmaCost:  300,
	},

	// Stealth assassin. Surgical strikes on high-value targets.
	// Combos: Active Camo + Backstab = devastating backline attacks
	//         Smoke Screen + Active Camo = team stealth
	//         Backstab + Flank position = position-based assassination
	Ghost: {
		ShipType:         string(Ghost),
		AttackType:       "Laser",
		LaserShield:      2,
		NuclearShield:    0,
		AntimatterShield: 2,
		Speed:            7,
		VisibilityRange:  5,
		AttackRange:      350,
		HP:               180,
		AttackDamage:     48,
		AttackInterval:   1.5,
		Abilities: []Ability{
			AbilitiesCatalog[AbilityActiveCamo],
			AbilitiesCatalog[AbilityBackstab],
			AbilitiesCatalog[AbilitySmokeScreen],
		},
		MetalCost:   250,
		CrystalCost: 300,
		PlasmaCost:  150,
	},

	// Electronic warfare specialist. Force multiplier through debuffs.
	// Combos: Sensor Jamming + Evasive Maneuvers = extreme evasion
	//         Ability Disruptor + Energy Drain = crippling debuff stack
	//         Energy Drain (stacks) = multiple frigates amplify effect
	Frigate: {
		ShipType:         string(Frigate),
		AttackType:       "Laser",
		LaserShield:      2,
		NuclearShield:    2,
		AntimatterShield: 2,
		Speed:            5,
		VisibilityRange:  6,
		AttackRange:      675,
		HP:               250,
		AttackDamage:     22,
		AttackInterval:   2.0,
		Abilities: []Ability{
			AbilitiesCatalog[AbilitySensorJamming],
			AbilitiesCatalog[AbilityAbilityDisruptor],
			AbilitiesCatalog[AbilityEnergyDrain],
		},
		MetalCost:   200,
		CrystalCost: 250,
		PlasmaCost:  100,
	},
}
