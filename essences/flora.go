package essences

import "github.com/nicoberrocal/galaxyCore/ships"

// BuildFlora: node-level mutations per earlier design.
func BuildFlora() *BioTree {
	tree := &BioTree{
		Name:        string(Flora),
		Description: "Flora biology tree focusing on area control, resilience, and support",
		Tiers:       make([][]*BioNode, 3),
	}

	// Tier 1: Carnivora (Aggressive, trapping abilities)
	carnivoraNodes := []*BioNode{
		// 1. Thigmonastic Triggers: First strike roots target
		{
			ID:          "carnivora_thigmonastic_triggers",
			Title:       "Thigmonastic Triggers",
			Description: "First strike attack gains 20% critical damage chance and roots the target in place for 1 tick.",
			Path:        string(ships.Carnivora),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnFirstStrike,
					Conditions: []Condition{
						{ConditionType: ConditionCombatState, CompareOp: CompareEqual, Value: "engaging"},
					},
					PrimaryEffect: &ships.StatMods{CritPct: 0.2},
					StatusEffects: []StatusEffect{{Name: "Root", Duration: 1, EffectType: StatusRoot}},
				},
			},
		},
		// 2. Digestive enzymes: Damage over time on attacks
		{
			ID:          "carnivora_digestive_enzymes",
			Title:       "Digestive Enzymes",
			Description: "Your attacks inflict a damage over time effect, dealing a small amount of stacking acid damage every tick for 5 ticks.",
			Path:        string(ships.Carnivora),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnSuccessfulHit,
					Conditions: []Condition{
						{ConditionType: ConditionCriticalHit, CompareOp: CompareEqual, Value: true},
					},
					Spawn: &SpawnEffect{SpawnType: SpawnAcidEffect, Duration: 5},
				},
			},
		},
		// 3. Irresistible aroma: Blinds nearby enemies
		{
			ID:          "carnivora_irresistible_aroma",
			Title:       "Irresistible Aroma",
			Description: "When a ship enters into your range view, blinds it of every stack except this one for 5 ticks.",
			Path:        string(ships.Carnivora),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnEnemyEnterRange,
					Conditions: []Condition{
						{ConditionType: ConditionEnemyCount, CompareOp: CompareGreater, Value: 0}},
					StatusEffects: []StatusEffect{{Name: "Blind", Duration: 5, EffectType: StatusBlind}},
				},
			},
		},
		// 4. Rapid Strike Tendrils: Increased formation countering bonus
		{
			ID:          "carnivora_rapid_strike_tendrils",
			Title:       "Rapid Strike Tendrils",
			Description: "Formation countering bonuses are increased by 20%.",
			Path:        string(ships.Carnivora),
			Effect:      ships.StatMods{FormationSyncBonus: 0.2},
		},
		// 5. Adaptive Apex: Resistance gain on kill
		{
			ID:          "carnivora_adaptive_apex",
			Title:       "Adaptive Apex",
			Description: "After destroying a stack, gain 3% resistance vs. that stack’s attack type (not biotree). Stacks up to 5×.",
			Path:        string(ships.Carnivora),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnKill,
					Conditions: []Condition{
						{ConditionType: ConditionKillCount, CompareOp: CompareGreater, Value: 0}},
					PrimaryEffect: &ships.StatMods{GlobalDefensePct: 0.03},
				},
			},
		},
	}

	// Tier 2: Arbor (Resilience and terrain-based bonuses)
	arborNodes := []*BioNode{
		// 1. Iron Chads: HP/shield bonus near asteroids
		{
			ID:          "arbor_iron_chads",
			Title:       "Iron Chads",
			Description: "Within 500u of asteroids, gain +1% HP and +1% to all shields stacking bonus for every percentage point of full deposits",
			Path:        string(ships.Arbor),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnNearAsteroid,
					Conditions: []Condition{
						{ConditionType: ConditionTerrainNear, CompareOp: CompareEqual, Value: "asteroid"}},
					PrimaryEffect: &ships.StatMods{HPPct: 0.01, LaserShieldDelta: 1, NuclearShieldDelta: 1, AntimatterShieldDelta: 1},
				},
			},
		},
		// 2. Deep root system: Immunity to displacement effects
		{
			ID:          "arbor_deep_root_system",
			Title:       "Deep Root System",
			Description: "You are immune to pull, push and slow effects from every stack that isn't your current target",
			Path:        string(ships.Arbor),
		},
		// 3. Root Coordination: Double damage near allies
		{
			ID:          "arbor_root_coordination",
			Title:       "Root Coordination",
			Description: "Every 3rd tick fighting within 200u of the same allied ship, you have double damage.",
			Path:        string(ships.Arbor),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnTick,
					Conditions: []Condition{
						{ConditionType: ConditionAllyNearby, CompareOp: CompareEqual, Value: true}},
					PrimaryEffect: &ships.StatMods{Damage: ships.DamageMods{LaserPct: 1, NuclearPct: 1, AntimatterPct: 1}},
					Duration:      1,
				},
			},
		},
		// 4. Spectrum Defiance: HP increase when targeted near a star
		{
			ID:          "arbor_spectrum_defiance",
			Title:       "Spectrum Defiance",
			Description: "Within 500u of a star, you gain a 5% stacking max HP increase for every target attacking you",
			Path:        string(ships.Arbor),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnNearStar,
					Conditions: []Condition{
						{ConditionType: ConditionTerrainNear, CompareOp: CompareEqual, Value: "star"},
						{ConditionType: ConditionIsAttacked, CompareOp: CompareEqual, Value: true}},
					PrimaryEffect: &ships.StatMods{HPPct: 0.05},
				},
			},
		},
		// 5. Seismic Root: AoE stun on formation change
		{
			ID:          "arbor_seismic_root",
			Title:       "Seismic Root",
			Description: "When finished changing to a Box formation, release 200u shockwave that deals significant damage and stuns enemies for 3 ticks",
			Path:        string(ships.Arbor),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnFormationChangeComplete,
					Conditions: []Condition{
						{ConditionType: ConditionFormationType, CompareOp: CompareEqual, Value: "Box"}},
					Spawn:         &SpawnEffect{SpawnType: SpawnShockwave, SpawnRadius: 200},
					StatusEffects: []StatusEffect{{Name: "Stun", Duration: 3, EffectType: StatusStun}},
				},
			},
		},
	}

	// Tier 3: Verdant Bloom (Support and healing)
	verdantBloomNodes := []*BioNode{
		// 1. Deep Rooted vitality: Increased HP regen
		{
			ID:          "verdant_bloom_deep_rooted_vitality",
			Title:       "Deep Rooted Vitality",
			Description: "Your base HP regen is increased by 100% in aggressive formation and 200% in defensive formations",
			Path:        string(ships.VerdantBloom),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnTick,
					Conditions: []Condition{
						{ConditionType: ConditionFormationType, CompareOp: CompareEqual, Value: "aggressive"}},
					PrimaryEffect: &ships.StatMods{AtCombatRegenPct: 1},
				},
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnTick,
					Conditions: []Condition{
						{ConditionType: ConditionFormationType, CompareOp: CompareEqual, Value: "defensive"}},
					PrimaryEffect: &ships.StatMods{AtCombatRegenPct: 2},
				},
			},
		},
		// 2. Photosynthetic Reactor: Energy regen and cooldown reduction near a star
		{
			ID:          "verdant_bloom_photosynthetic_reactor",
			Title:       "Photosynthetic Reactor",
			Description: "Your base energy regeneration is increased by 30% while within 500u of a star, you and all nearby allies gain 1% ability cooldown recovery",
			Path:        string(ships.VerdantBloom),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnNearStar,
					Conditions: []Condition{
						{ConditionType: ConditionTerrainNear, CompareOp: CompareEqual, Value: "star"}},
					PrimaryEffect: &ships.StatMods{AbilityCooldownPct: -0.01},
					AoE:           &AoETraitTarget{Radius: 500, TargetType: AoEAllies},
				},
			},
		},
		// 3. Sap Flow: Overhealing becomes a shield
		{
			ID:          "verdant_bloom_sap_flow",
			Title:       "Sap Flow",
			Description: "When within 200u of an ally below 40% HP, 30% of overhealing you receive becomes a temporary damage shield. Attacks on damage dont trigger crits nor attack-based-healing.",
			Path:        string(ships.VerdantBloom),
		},
		// 4. Pollen Cloud: AoE regen buff on ability use
		{
			ID:          "verdant_bloom_pollen_cloud",
			Title:       "Pollen Cloud",
			Description: "All allied ships around a 250u aura gain a 25% HP and Energy Regen increase of 25% your base regeneration rate after using an ability. 1 application max.",
			Path:        string(ships.VerdantBloom),
			ComplexEffects: []ComplexEffect{
				{
					EffectType: ComplexConditional,
					Trigger:    TriggerOnAbilityCast,
					Conditions: []Condition{
						{ConditionType: ConditionAbilityUsed, CompareOp: CompareEqual, Value: true}},
					PrimaryEffect: &ships.StatMods{AtCombatRegenPct: 0.25},
					AoE:           &AoETraitTarget{Radius: 250, TargetType: AoEAllies},
				},
			},
		},
		// 5. Blossom through extinction: On-death AoE heal
		{
			ID:          "verdant_bloom_blossom_through_extinction",
			Title:       "Blossom Through Extinction",
			Description: "At death, restore 25% of max HP and Energy to allied ships in 250u radius.",
			Path:        string(ships.VerdantBloom),
			ComplexEffects: []ComplexEffect{
				{
					EffectType:    ComplexOnDeath,
					Trigger:       TriggerOnDeath,
					PrimaryEffect: &ships.StatMods{HPPct: 0.25},
					AoE:           &AoETraitTarget{Radius: 250, TargetType: AoEAllies},
				},
			},
		},
	}

	tree.Tiers[0] = carnivoraNodes
	tree.Tiers[1] = arborNodes
	tree.Tiers[2] = verdantBloomNodes

	return tree
}
