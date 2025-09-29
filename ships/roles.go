package ships

// RoleMode defines a soft posture for a ship type. Modes are not hard role swaps,
// they provide modest, capped modifiers and sometimes gate abilities.
// Design goals:
// - Provide tactical flexibility without erasing specialist hull identities.
// - Add commitment/friction to mode switching to reward planning and intel.
// - Support anti-doomstack strategy by enabling dispersed stacks to pivot.
//
// Modes:
//   - Tactical: default combat posture. Minor DPS/ROF bonuses. No economy/science.
//   - Economic: allows anchoring to gather resources at reduced efficiency vs Drones.
//   - Recon: vision-focused. Improves detection and scouting tools.
//   - Scientific: salvage/anomaly analysis focus. Improves OOC repair and cooldowns.
//
// Switching Guidance (gameplay layer, documented here for implementers):
//   - Reconfigure should take time (e.g., 180s baseline; Scouts faster; Engineering gems faster).
//   - Mode-switching should not bypass the 1h post-attack cooldown.
//   - Economic mode while anchored should disable warp travel.
//   - Some abilities are disabled in specific modes (e.g., Bomber StandoffPattern in Economic).
//   - Optionally emit a detectable "reconfiguration signal" for intel play.
type RoleMode string

const (
    RoleTactical  RoleMode = "tactical"
    RoleEconomic  RoleMode = "economic"
    RoleRecon     RoleMode = "recon"
    RoleScientific RoleMode = "scientific"
)

// RoleModeSpec declaratively describes a role mode's intent and its stat modifiers.
// The runtime system can read this spec to apply caps and UI rules.
type RoleModeSpec struct {
    Mode         RoleMode // identifier
    Name         string   // human-readable
    Description  string   // long-form documentation

    // BaseMods are applied while the unit is in this mode (on top of sockets/runewords).
    BaseMods     StatMods

    // UX/rules hints (enforced at higher layers):
    ReconfigureSeconds int      // suggested base reconfigure time
    WarpAllowed        bool     // whether warp is allowed in this mode
    RequiresAnchoring  bool     // whether the mode's benefits require anchoring

    // Ability gates (identifiers from abilities.go)
    DisabledAbilities []AbilityID // abilities not usable in this mode
    EnabledAbilities  []AbilityID // abilities only usable in this mode (optional)
}

// RoleModesCatalog: tweak numbers freely. Values here are conservative baselines.
var RoleModesCatalog = map[RoleMode]RoleModeSpec{
    RoleTactical: {
        Mode:        RoleTactical,
        Name:        "Tactical",
        Description: "Default combat posture. Modest firepower bonuses. No economy/science perks.",
        BaseMods: StatMods{
            AttackIntervalPct: -0.10, // ~10% faster ROF
            // Mild damage bonus expressed per-type for clarity
            Damage: DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
        },
        ReconfigureSeconds: 180,
        WarpAllowed:        true,
    },
    RoleEconomic: {
        Mode:        RoleEconomic,
        Name:        "Economic",
        Description: "Allows anchoring to gather resources at reduced efficiency; weaker combat while anchored.",
        BaseMods: StatMods{
            Damage: DamageMods{LaserPct: -0.25, NuclearPct: -0.25, AntimatterPct: -0.25}, // weaker offense
            // While anchored, an additional -1 shield is recommended (enforce at runtime near asteroids/nebulas)
        },
        ReconfigureSeconds: 180,
        WarpAllowed:        true, // becomes false if Anchored (runtime rule)
        RequiresAnchoring:  false,
        // Example ability restrictions: siege/tactical tools are disabled while anchored mining
        DisabledAbilities: []AbilityID{AbilityStandoffPattern},
        EnabledAbilities:  []AbilityID{AbilityResourceHarvester},
    },
    RoleRecon: {
        Mode:        RoleRecon,
        Name:        "Recon",
        Description: "Vision and detection focus. Slight combat reduction.",
        BaseMods: StatMods{
            VisibilityDelta:  +3,
            AttackRangeDelta: -1,
            Damage: DamageMods{LaserPct: -0.15, NuclearPct: -0.15, AntimatterPct: -0.15},
            CloakDetect: true, // pairs especially well with Sensor gems
            PingRangePct: 0.25,
        },
        ReconfigureSeconds: 120, // faster to adopt
        WarpAllowed:        true,
        EnabledAbilities:   []AbilityID{AbilityPing, AbilityDecoyBeacon},
    },
    RoleScientific: {
        Mode:        RoleScientific,
        Name:        "Scientific",
        Description: "Salvage/anomaly analysis. Stronger OOC repair and shorter ability cooldowns.",
        BaseMods: StatMods{
            OutOfCombatRegenPct: 0.50,
            AbilityCooldownPct:  -0.10,
            AttackIntervalPct:   +0.10, // slower ROF as tradeoff
            Damage: DamageMods{LaserPct: -0.20, NuclearPct: -0.20, AntimatterPct: -0.20},
        },
        ReconfigureSeconds: 180,
        WarpAllowed:        true,
    },
}

// RoleModeMods is a helper to get the StatMods for a mode.
// shipType can be used by game logic to add hull-specific adjustments (not included here).
func RoleModeMods(mode RoleMode, shipType string) StatMods {
    if spec, ok := RoleModesCatalog[mode]; ok {
        return spec.BaseMods
    }
    return ZeroMods()
}
