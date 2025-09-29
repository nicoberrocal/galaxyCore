package ships

// AbilityKind classifies how an ability behaves mechanically.
// - passive: always on, no activation, no duration
// - active: click-to-activate with cooldown and duration
// - toggle: on/off state; may drain or carry opportunity cost
// - aura: passive area effect around the ship while ability is "on"
// - travel: mobility/warp capability
// - conditional: enabled only under certain conditions (e.g., anchored)
type AbilityKind string

const (
    AbilityPassive     AbilityKind = "passive"
    AbilityActive      AbilityKind = "active"
    AbilityToggle      AbilityKind = "toggle"
    AbilityAura        AbilityKind = "aura"
    AbilityTravel      AbilityKind = "travel"
    AbilityConditional AbilityKind = "conditional"
)

// AbilityID is the stable programmatic identifier for an ability.
type AbilityID string

const (
    // Recon/Intel
    AbilityLongRangeSensors AbilityID = "LongRangeSensors" // Passive: +vision; may detect cloak with Sensor gems
    AbilityPing              AbilityID = "Ping"              // Active: mark target stack; synergy with Focus Fire/Uplink
    AbilityDecoyBeacon       AbilityID = "DecoyBeacon"       // Active: create phantom contact in fog-of-war

    // Travel/Logistics
    AbilityLightSpeed     AbilityID = "LightSpeed"     // Travel: enables warp
    AbilityWarpStabilizer AbilityID = "WarpStabilizer" // Aura: reduce allied scatter; resist interdiction

    // Carrier Ops
    AbilityHangarLaunch       AbilityID = "HangarLaunch"       // Active: deploy/garrison escorts; brief launch buffs
    AbilityPointDefenseScreen AbilityID = "PointDefenseScreen" // Toggle/Aura: AoE laser mitigation for nearby allies

    // Strike/Burst
    AbilityAlphaStrike     AbilityID = "AlphaStrike"     // Active: first-volley damage boost; added self-cooldown cost
    AbilityOverload        AbilityID = "Overload"        // Active: +damage, -shields for a window (risky)
    AbilityInterdictorPulse AbilityID = "InterdictorPulse" // Active: blocks enemy warp in area; reveals self

    // Siege/Fire Support
    AbilitySiegePayload   AbilityID = "SiegePayload"   // Passive: bonus vs structures; small splash
    AbilityStandoffPattern AbilityID = "StandoffPattern" // Toggle: +range, slower ROF
    AbilityTargetingUplink AbilityID = "TargetingUplink" // Active: +accuracy/crit vs marked targets

    // Versatility/Defense
    AbilityAdaptiveTargeting AbilityID = "AdaptiveTargeting" // Active: temporarily override attack type
    AbilityFocusFire         AbilityID = "FocusFire"         // Active: bonus vs marked/low-HP-bucket targets
    AbilityEvasiveManeuvers  AbilityID = "EvasiveManeuvers"  // Active: +evasion/+LaserShield, -range briefly

    // Economy/Utility
    AbilityResourceHarvester  AbilityID = "ResourceHarvester"  // Toggle: anchor to gather resources (asteroids/nebulas)
    AbilitySelfRepair         AbilityID = "SelfRepair"         // Passive: out-of-combat bucket regen
    AbilityCloakWhileAnchored AbilityID = "CloakWhileAnchored" // Conditional: cloaked when anchored gathering

    // GemWord-granted abilities
    AbilityLaserOvercharge AbilityID = "LaserOvercharge" // Active: short ROF burst; slight heat/backlash risk handled at combat layer
    AbilityBunkerBuster    AbilityID = "BunkerBuster"    // Active: bonus vs fortified structures; long cooldown
    AbilityPhaseLance      AbilityID = "PhaseLance"      // Active: opener with partial shield ignore
    AbilityWideAreaPing    AbilityID = "WideAreaPing"    // Active: reveal large area; very long cooldown
    AbilityRapidRedeploy   AbilityID = "RapidRedeploy"   // Active: reduce post-attack cooldown or warp charge for a short window
)

// AbilitiesCatalog provides default configuration for abilities.
// NOTE: Numerical values are sensible baselines. Balance freely.
var AbilitiesCatalog = map[AbilityID]Ability{
    AbilityLongRangeSensors: {
        ID:              AbilityLongRangeSensors,
        Name:            "Long-Range Sensors",
        Kind:            AbilityPassive,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Increases VisibilityRange. With Sensor runes, detects cloaked or switching-mode signals.",
    },
    AbilityPing: {
        ID:              AbilityPing,
        Name:            "Ping",
        Kind:            AbilityActive,
        CooldownSeconds: 300,
        DurationSeconds: 30,
        Description:     "Reveals a target tile and marks a stack, increasing incoming damage for certain weapons.",
    },
    AbilityDecoyBeacon: {
        ID:              AbilityDecoyBeacon,
        Name:            "Decoy Beacon",
        Kind:            AbilityActive,
        CooldownSeconds: 600,
        DurationSeconds: 90,
        Description:     "Projects a phantom contact in fog-of-war to mislead opponents.",
    },
    AbilityLightSpeed: {
        ID:              AbilityLightSpeed,
        Name:            "Light-Speed Travel",
        Kind:            AbilityTravel,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Enables warp travel. Subject to warp charge, scatter, and interdiction.",
    },
    AbilityWarpStabilizer: {
        ID:              AbilityWarpStabilizer,
        Name:            "Warp Stabilizer",
        Kind:            AbilityAura,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Reduces allied warp-in scatter and increases interdiction resistance within aura.",
    },
    AbilityHangarLaunch: {
        ID:              AbilityHangarLaunch,
        Name:            "Hangar Launch",
        Kind:            AbilityActive,
        CooldownSeconds: 180,
        DurationSeconds: 20,
        Description:     "Deploys or recalls escorts. On launch, escorts gain brief buffs.",
    },
    AbilityPointDefenseScreen: {
        ID:              AbilityPointDefenseScreen,
        Name:            "Point-Defense Screen",
        Kind:            AbilityToggle,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Projects an area shield that mitigates incoming Laser fire to allies in range.",
    },
    AbilityAlphaStrike: {
        ID:              AbilityAlphaStrike,
        Name:            "Alpha Strike",
        Kind:            AbilityActive,
        CooldownSeconds: 240,
        DurationSeconds: 10,
        Description:     "Greatly amplifies the first volley against a target at the cost of increased self-cooldown.",
    },
    AbilityOverload: {
        ID:              AbilityOverload,
        Name:            "Overload",
        Kind:            AbilityActive,
        CooldownSeconds: 180,
        DurationSeconds: 20,
        Description:     "+Damage for a short window while reducing shields; risky spike.",
    },
    AbilityInterdictorPulse: {
        ID:              AbilityInterdictorPulse,
        Name:            "Interdictor Pulse",
        Kind:            AbilityActive,
        CooldownSeconds: 300,
        DurationSeconds: 60,
        Description:     "Blocks enemy warp within an area and reveals the caster. Counters light-speed escapes.",
    },
    AbilitySiegePayload: {
        ID:              AbilitySiegePayload,
        Name:            "Siege Payload",
        Kind:            AbilityPassive,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Bonus damage to structures and small splash damage on orbitables.",
    },
    AbilityStandoffPattern: {
        ID:              AbilityStandoffPattern,
        Name:            "Standoff Pattern",
        Kind:            AbilityToggle,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "+AttackRange but slower rate of fire while toggled.",
    },
    AbilityTargetingUplink: {
        ID:              AbilityTargetingUplink,
        Name:            "Targeting Uplink",
        Kind:            AbilityActive,
        CooldownSeconds: 180,
        DurationSeconds: 30,
        Description:     "+Accuracy and +Crit, especially effective against Ping-marked targets.",
    },
    AbilityAdaptiveTargeting: {
        ID:              AbilityAdaptiveTargeting,
        Name:            "Adaptive Targeting",
        Kind:            AbilityActive,
        CooldownSeconds: 120,
        DurationSeconds: 20,
        Description:     "Temporarily overrides the ship's attack type to counter enemy shields.",
    },
    AbilityFocusFire: {
        ID:              AbilityFocusFire,
        Name:            "Focus Fire",
        Kind:            AbilityActive,
        CooldownSeconds: 30,
        DurationSeconds: 10,
        Description:     "Focuses attacks on marked or lowest-HP-bucket targets for increased damage.",
    },
    AbilityEvasiveManeuvers: {
        ID:              AbilityEvasiveManeuvers,
        Name:            "Evasive Maneuvers",
        Kind:            AbilityActive,
        CooldownSeconds: 90,
        DurationSeconds: 12,
        Description:     "Grants temporary evasion and LaserShield at the cost of reduced AttackRange.",
    },
    AbilityResourceHarvester: {
        ID:              AbilityResourceHarvester,
        Name:            "Resource Harvester",
        Kind:            AbilityToggle,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Anchors the ship to gather resources from asteroids or nebulas.",
    },
    AbilitySelfRepair: {
        ID:              AbilitySelfRepair,
        Name:            "Self-Repair",
        Kind:            AbilityPassive,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Restores HP buckets over time while out of combat.",
    },
    AbilityCloakWhileAnchored: {
        ID:              AbilityCloakWhileAnchored,
        Name:            "Cloak While Anchored",
        Kind:            AbilityConditional,
        CooldownSeconds: 0,
        DurationSeconds: 0,
        Description:     "Remains cloaked when anchored for harvesting; breaks on movement or damage.",
    },
    // GemWord abilities
    AbilityLaserOvercharge: {
        ID:              AbilityLaserOvercharge,
        Name:            "Laser Overcharge",
        Kind:            AbilityActive,
        CooldownSeconds: 180,
        DurationSeconds: 12,
        Description:     "Briefly overclocks Laser weapons, increasing rate of fire.",
    },
    AbilityBunkerBuster: {
        ID:              AbilityBunkerBuster,
        Name:            "Bunker Buster",
        Kind:            AbilityActive,
        CooldownSeconds: 300,
        DurationSeconds: 8,
        Description:     "Specialized munition dealing heavy damage to fortified structures.",
    },
    AbilityPhaseLance: {
        ID:              AbilityPhaseLance,
        Name:            "Phase Lance",
        Kind:            AbilityActive,
        CooldownSeconds: 240,
        DurationSeconds: 5,
        Description:     "Opener that phases through a portion of shields for high burst.",
    },
    AbilityWideAreaPing: {
        ID:              AbilityWideAreaPing,
        Name:            "Wide-Area Ping",
        Kind:            AbilityActive,
        CooldownSeconds: 900,
        DurationSeconds: 30,
        Description:     "Reveals a wide area of the map and highlights mode-switching signals.",
    },
    AbilityRapidRedeploy: {
        ID:              AbilityRapidRedeploy,
        Name:            "Rapid Redeploy",
        Kind:            AbilityActive,
        CooldownSeconds: 420,
        DurationSeconds: 60,
        Description:     "Temporarily reduces warp charge time and/or post-attack cooldown for warp-capable ships.",
    },
}
