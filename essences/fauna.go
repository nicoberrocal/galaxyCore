package essences

import "github.com/nicoberrocal/galaxyCore/ships"

// BuildFauna builds the Fauna tree and uses Path-based essence mutation (applies to each node in chosen path).
func BuildFauna() *BioTree {
	tree := &BioTree{
		Name:        string(Fauna),
		Description: "Fauna biology tree focusing on aggressive hunting, pack tactics, and scavenging",
		Tiers:       make([][]*BioNode, 3),
	}

	// Tier 1: Apex (Solo hunter, high-risk/high-reward)
	apexNodes := []*BioNode{
		// 1. Kill Focus: Crit chance and visibility on new enemy detection
		{
			ID:          "apex_kill_focus",
			Title:       "Kill Focus",
			Description: "When detecting a new enemy, gain +10% crit chance for 5 ticks and +15% increased visibility range. Max: 1 instances.",
			Path:        string(ships.Apex),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionEnemyCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect:  &ships.StatMods{CritPct: 0.1, VisibilityDelta: 15},
					Duration:       5, // ticks
					MaxActivations: 1,
				},
			},
		},
		// 2. Predators Pounce: Bonus damage on first strike
		{
			ID:          "apex_predators_pounce",
			Title:       "Predator's Pounce",
			Description: "First strike attacks deal 30% bonus damage",
			Path:        string(ships.Apex),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCombatState, CompareOp: CompareEqual, Value: "engaging"},
					},
					PrimaryEffect: &ships.StatMods{Damage: ships.DamageMods{LaserPct: 0.3, NuclearPct: 0.3, AntimatterPct: 0.3}},
					Duration:      1, // First strike only
				},
			},
		},
		// 3. Alpha Dominance: Steal ally's damage
		{
			ID:          "apex_alpha_dominance",
			Title:       "Alpha Dominance",
			Description: "When defensively formed ally gets attacked within 200u, you steal 25% of your ally's damage and apply it to target ignoring shields for 1 tick. Cooldown: 10 ticks.",
			Path:        string(ships.Apex),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionTargetIsAttackingAlly, CompareOp: CompareEqual, Value: true},
					},
					Cooldown: 10, // ticks
				},
			},
		},
		// 4. Panicked Herd: Fear on kill
		{
			ID:          "apex_panicked_herd",
			Title:       "Panicked Herd",
			Description: "Killing a target spreads fear among nearby enemies (-25% damage for 2 ticks)",
			Path:        string(ships.Apex),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					Spawn: &SpawnEffect{SpawnType: SpawnFearEffect, Duration: 2},
				},
			},
		},
		// 5. Apex Threat: Heal on crit
		{
			ID:          "apex_apex_threat",
			Title:       "Apex Threat",
			Description: "Crit attacks heal you for 25% of attack value and reduces 5% cooldown of your abilities. 3 Ticks cooldown",
			Path:        string(ships.Apex),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{HPPct: 0.25, AbilityCooldownPct: -0.05},
					Cooldown:      3, // ticks
				},
			},
		},
	}

	// Tier 2: Pack Hunter (Coordinated attacks and buffs)
	packHunterNodes := []*BioNode{
		// 1. Pack Mentality: Damage/evasion bonus per nearby ally
		{
			ID:          "pack_hunter_pack_mentality",
			Title:       "Pack Mentality",
			Description: "For each other allied ship within 150u gain +2% damage and +1% evasion",
			Path:        string(ships.PackHunter),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAllyCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{Damage: ships.DamageMods{LaserPct: 0.02, NuclearPct: 0.02, AntimatterPct: 0.02}, EvasionPct: 0.01},
				},
			},
		},
		// 2. Flanking Instincts: Damage bonus to flanked enemies
		{
			ID:          "pack_hunter_flanking_instincts",
			Title:       "Flanking Instincts",
			Description: "You deal +15% damage to enemies that are being attacked by another ally at rear or flank",
			Path:        string(ships.PackHunter),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionAttackFromBehind, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{Damage: ships.DamageMods{LaserPct: 0.15, NuclearPct: 0.15, AntimatterPct: 0.15}},
				},
			},
		},
		// 3. Communal Defense: Redirect lethal damage
		{
			ID:          "pack_hunter_communal_defense",
			Title:       "Communal Defense",
			Description: "When an enemy attack will kill an allied stack, the highest allied HP stack in 300u receives the damage instead. The stack that received the attack is cleansed of all debuffs. Cooldown 15 ticks, shared by all allied ships in the radius.",
			Path:        string(ships.PackHunter),
		},
		// 4. Hunting Cry: AoE buff on combat start
		{
			ID:          "pack_hunter_hunting_cry",
			Title:       "Hunting Cry",
			Description: "Upon entering combat, grant all allies a +10% speed globally and 20% chance to crit damage in 200u and grants +5% ability cooldown recovery to allies under Pack Mentality",
			Path:        string(ships.PackHunter),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCombatState, CompareOp: CompareEqual, Value: "engaging"},
					},
					PrimaryEffect: &ships.StatMods{SpeedDelta: 10, CritPct: 0.2, AbilityCooldownPct: -0.05},
					AoE:           &AoETraitTarget{Radius: 200, TargetType: AoEAllies},
				},
			},
		},
		// 5. Feeding Frenzy: Stacking damage on kill
		{
			ID:          "pack_hunter_feeding_frenzy",
			Title:       "Feeding Frenzy",
			Description: "Killing an enemy ship applies a +5% damage stack. Max 5 stacks.",
			Path:        string(ships.PackHunter),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{Damage: ships.DamageMods{LaserPct: 0.05, NuclearPct: 0.05, AntimatterPct: 0.05}},
				},
			},
		},
	}

	// Tier 3: Scavengers (Opportunistic and resilient)
	scavengersNodes := []*BioNode{
		// 1. Cannibal Regrowth: Heal on nearby deaths
		{
			ID:          "scavengers_cannibal_regrowth",
			Title:       "Cannibal Regrowth",
			Description: "Each nearby death (enemy or ally within 200u) restores 5% of your max HP",
			Path:        string(ships.Scavengers),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{HPPct: 0.05},
				},
			},
		},
		// 2. Necrotic Bite: Shield/regen debuff on hit
		{
			ID:          "scavengers_necrotic_bite",
			Title:       "Necrotic Bite",
			Description: "Your attacks inflict a stacking debuff that reduces the targets shield by 2%, HP regen and healing received by 15% per stack up to 5 stacks",
			Path:        string(ships.Scavengers),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					PrimaryEffect: &ships.StatMods{LaserShieldDelta: -2, NuclearShieldDelta: -2, AntimatterShieldDelta: -2, AtCombatRegenPct: -0.15},
				},
			},
		},
		// 3. Spoil Extractor: Regen/cooldown bonus on nearby death
		{
			ID:          "scavengers_spoil_extractor",
			Title:       "Spoil Extractor",
			Description: "Upon ship death within 500u, gain a 2% stacking bonus to regen and cooldown recovery (max 5 stacks, lasts 10 ticks).",
			Path:        string(ships.Scavengers),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0},
					},
					PrimaryEffect: &ships.StatMods{AtCombatRegenPct: 0.02, AbilityCooldownPct: -0.02},
					Duration:      10, // ticks
				},
			},
		},
		// 4. Parasitic Latches: Siphon energy on hit
		{
			ID:          "scavengers_parasitic_latches",
			Title:       "Parasitic Latches",
			Description: "When you damage an enemy ship, you siphon a small amount of energy from them, reducing ability cooldown recovery and granting it to you. Drains diminish by 25% per tick",
			Path:        string(ships.Scavengers),
		},
		// 5. Winged Brother: Increased visibility and detection
		{
			ID:          "scavengers_winged_brother",
			Title:       "Winged Brother",
			Description: "You have +10× visibility range and detect destroyed structures and derelict ships globally",
			Path:        string(ships.Scavengers),
			Effect:      ships.StatMods{VisibilityDelta: 1000},
		},
	}

	tree.Tiers[0] = apexNodes
	tree.Tiers[1] = packHunterNodes
	tree.Tiers[2] = scavengersNodes

	return tree
}
