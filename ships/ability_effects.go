package ships

import "time"

// AbilityEffectsCatalog maps abilities to their stat modifier effects.
// This allows abilities to be properly integrated into the compute_v2 modifier system.
// Passive abilities apply permanently, active/toggle abilities apply when active.
var AbilityEffectsCatalog = map[AbilityID]StatMods{
	// Recon/Intel
	AbilityLongRangeSensors: {
		VisibilityDelta: 2,
		PingRangePct:    0.15,
	},
	
	// Carrier Ops
	AbilityPointDefenseScreen: {
		LaserShieldDelta: 2,
	},
	
	// Siege/Fire Support
	AbilitySiegePayload: {
		StructureDamagePct: 0.20,
		SplashRadiusDelta:  1,
	},
	
	AbilityStandoffPattern: {
		AttackRangeDelta:   2,
		AttackIntervalPct:  0.30, // Slower ROF
	},
	
	// Defense
	AbilityEvasiveManeuvers: {
		EvasionPct:       0.25,
		LaserShieldDelta: 2,
		AttackRangeDelta: -1,
	},
	
	// Economy/Utility
	AbilitySelfRepair: {
		OutOfCombatRegenPct: 0.15,
	},
	
	// Active abilities (when activated)
	AbilityAlphaStrike: {
		FirstVolleyPct: 0.50,
		CritPct:        0.20,
	},
	
	AbilityOverload: {
		Damage: DamageMods{
			LaserPct:      0.40,
			NuclearPct:    0.40,
			AntimatterPct: 0.40,
		},
		LaserShieldDelta:      -2,
		NuclearShieldDelta:    -2,
		AntimatterShieldDelta: -2,
	},
	
	AbilityTargetingUplink: {
		AccuracyPct: 0.25,
		CritPct:     0.15,
	},
	
	AbilityFocusFire: {
		Damage: DamageMods{
			LaserPct:      0.30,
			NuclearPct:    0.30,
			AntimatterPct: 0.30,
		},
	},
	
	// GemWord abilities
	AbilityLaserOvercharge: {
		AttackIntervalPct: -0.35,
		Damage: DamageMods{
			LaserPct: 0.20,
		},
	},
	
	AbilityBunkerBuster: {
		StructureDamagePct: 0.50,
		SplashRadiusDelta:  2,
	},
	
	AbilityPhaseLance: {
		ShieldPiercePct: 0.40,
		FirstVolleyPct:  0.30,
	},
	
	AbilityWideAreaPing: {
		VisibilityDelta: 5,
		PingRangePct:    1.0,
	},
	
	AbilityRapidRedeploy: {
		WarpChargePct:  -0.40,
		WarpScatterPct: -0.50,
	},
	
	// New ship abilities
	AbilityShieldOvercharge: {
		LaserShieldDelta:      3,
		NuclearShieldDelta:    3,
		AntimatterShieldDelta: 3,
	},
	
	AbilityRammingSpeed: {
		SpeedDelta: 3,
		// Contact damage handled in combat logic
	},
	
	AbilityRepairDrones: {
		OutOfCombatRegenPct: 0.10,
		// In-combat regen handled separately
	},
	
	AbilityPursuitProtocol: {
		// Speed bonus conditional on target speed - handled in combat logic
		SpeedDelta: 2,
	},
	
	AbilityAntimatterBurst: {
		// Single shot 3x damage - handled in combat logic
		Damage: DamageMods{
			AntimatterPct: 2.0, // +200% = 3x total
		},
	},
	
	AbilityClusterMunitions: {
		SplashRadiusDelta: 2,
	},
	
	AbilityBarrageMode: {
		AttackRangeDelta:   1,
		SplashRadiusDelta:  1,
		AttackIntervalPct:  0.30, // -30% ROF
	},
	
	AbilitySuppressiveFire: {
		// Area denial - handled in combat logic
		SplashRadiusDelta: 2,
	},
	
	AbilityActiveCamo: {
		SpeedDelta: -3,
		// Invisibility handled in visibility logic
	},
	
	AbilityBackstab: {
		// Position-based damage bonus handled in combat logic
		CritPct: 0.50,
	},
	
	AbilitySmokeScreen: {
		// AoE cloak handled in visibility logic
		EvasionPct: 0.20,
	},
	
	AbilitySensorJamming: {
		// Enemy accuracy reduction handled in combat logic
		AccuracyPct: 0.15, // Bonus to allies
	},
	
	AbilityAbilityDisruptor: {
		// Enemy cooldown increase handled in combat logic
		AbilityCooldownPct: -0.10, // Bonus to allies
	},
	
	AbilityEnergyDrain: {
		// Enemy damage reduction handled in combat logic
		Damage: DamageMods{
			LaserPct:      0.10,
			NuclearPct:    0.10,
			AntimatterPct: 0.10,
		},
	},
}

// GetAbilityMods returns the stat modifiers for an ability.
// Returns zero mods if the ability has no stat effects.
func GetAbilityMods(abilityID AbilityID) StatMods {
	if mods, ok := AbilityEffectsCatalog[abilityID]; ok {
		return mods
	}
	return ZeroMods()
}

// AddActiveAbilities adds modifiers from currently active abilities to the builder.
func (mb *ModifierBuilder) AddActiveAbilities(activeAbilities []AbilityID, durations map[AbilityID]time.Duration) *ModifierBuilder {
	for _, abilityID := range activeAbilities {
		mods := GetAbilityMods(abilityID)
		if !isZeroMods(mods) {
			duration := durations[abilityID]
			if duration > 0 {
				mb.AddAbility(abilityID, mods, duration)
			} else {
				// Passive or toggle ability - add as permanent
				mb.stack.AddPermanent(
					SourceAbility,
					string(abilityID),
					string(abilityID),
					mods,
					PriorityAbility,
					mb.now,
				)
			}
		}
	}
	return mb
}
