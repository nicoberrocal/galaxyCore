package ships

// ShipBlueprints enumerates baseline ship TYPE definitions with recommended
// abilities and default RoleMode. These are data-only; runtime state lives
// elsewhere (see stack.go for HP buckets, movement, etc.).
//
// Notes
// - Numbers are conservative baselines based on our prior design. Tweak freely.
// - Abilities are limited to 3 per type. Some types have additional synergies
//   available via runes (see runes.go) and RoleModes (see roles.go).
// - AttackType is one of: "Laser", "Nuclear", "Antimatter".
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
        AttackRange:      1,
        HP:               100,
        AttackDamage:     5,
        AttackInterval:   3.0,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityResourceHarvester],
            AbilitiesCatalog[AbilitySelfRepair],
            AbilitiesCatalog[AbilityCloakWhileAnchored],
        },
        RoleMode:         RoleEconomic, // default posture
        Sockets:          nil,
        MetalCost:        20,
        CrystalCost:      10,
        PlasmaCost:       0,
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
        AttackRange:      2,
        HP:               100,
        AttackDamage:     10,
        AttackInterval:   1.5,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityLongRangeSensors],
            AbilitiesCatalog[AbilityPing],
            AbilitiesCatalog[AbilityDecoyBeacon],
        },
        RoleMode:    RoleRecon,
        Sockets:     nil,
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
        AttackRange:      2,
        HP:               200,
        AttackDamage:     20,
        AttackInterval:   1.0,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityAdaptiveTargeting],
            AbilitiesCatalog[AbilityFocusFire],
            AbilitiesCatalog[AbilityEvasiveManeuvers],
        },
        RoleMode:    RoleTactical,
        Sockets:     nil,
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
        AttackRange:      6,
        HP:               500,
        AttackDamage:     70,
        AttackInterval:   3.0,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityLightSpeed],
            AbilitiesCatalog[AbilitySiegePayload],
            AbilitiesCatalog[AbilityStandoffPattern],
        },
        RoleMode:    RoleTactical,
        Sockets:     nil,
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
        AttackRange:      3,
        HP:               900,
        AttackDamage:     60,
        AttackInterval:   2.0,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityLightSpeed],
            AbilitiesCatalog[AbilityHangarLaunch],
            AbilitiesCatalog[AbilityPointDefenseScreen],
        },
        RoleMode:          RoleTactical,
        Sockets:           nil,
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
        AttackRange:      3,
        HP:               600,
        AttackDamage:     80,
        AttackInterval:   2.5,
        Abilities: []Ability{
            AbilitiesCatalog[AbilityLightSpeed],
            AbilitiesCatalog[AbilityAlphaStrike],
            AbilitiesCatalog[AbilityInterdictorPulse],
        },
        RoleMode:    RoleTactical,
        Sockets:     nil,
        MetalCost:   500,
        CrystalCost: 700,
        PlasmaCost:  400,
    },
}
