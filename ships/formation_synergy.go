package ships

// Formation-enhanced ability and gem synergies
// This file extends the abilities and gems systems with formation-aware bonuses.

// AbilityFormationMod defines how abilities are modified when used from specific formation positions.
type AbilityFormationMod struct {
	AbilityID     AbilityID          `bson:"abilityId" json:"abilityId"`
	Position      FormationPosition  `bson:"position" json:"position"`
	Modifications map[string]float64 `bson:"modifications" json:"modifications"` // stat -> multiplier
}

// AbilityFormationModsCatalog contains all ability-position synergies.
var AbilityFormationModsCatalog = []AbilityFormationMod{
	{
		AbilityID: AbilityFocusFire,
		Position:  PositionFront,
		Modifications: map[string]float64{
			"CooldownSeconds":  -0.5,  // 50% faster cooldown
			"DamageMultiplier": 1.2,   // +20% damage
		},
	},
	{
		AbilityID: AbilityCloakWhileAnchored,
		Position:  PositionFlank,
		Modifications: map[string]float64{
			"DurationSeconds":      2.0, // +2x duration
			"DetectionResistance":  0.3, // +30% harder to detect
		},
	},
	{
		AbilityID: AbilityAlphaStrike,
		Position:  PositionFront,
		Modifications: map[string]float64{
			"DamageMultiplier": 1.3,  // +30% first strike damage
			"CritBonus":        0.15, // +15% crit chance
		},
	},
	{
		AbilityID: AbilityEvasiveManeuvers,
		Position:  PositionFlank,
		Modifications: map[string]float64{
			"EvasionBonus":   0.25, // +25% evasion
			"DurationSeconds": 1.5,  // +50% duration multiplier
		},
	},
	{
		AbilityID: AbilityStandoffPattern,
		Position:  PositionBack,
		Modifications: map[string]float64{
			"RangeBonus":       0.30, // +30% range
			"DamageMultiplier": 1.15, // +15% damage
		},
	},
	{
		AbilityID: AbilityPing,
		Position:  PositionBack,
		Modifications: map[string]float64{
			"RangeMultiplier": 1.5, // +50% ping range
			"DurationSeconds": 1.5, // +50% mark duration multiplier
		},
	},
	{
		AbilityID: AbilityTargetingUplink,
		Position:  PositionSupport,
		Modifications: map[string]float64{
			"AccuracyBonus":    0.15, // +15% accuracy
			"DurationSeconds":  1.3,  // +30% duration multiplier
		},
	},
	{
		AbilityID: AbilityPointDefenseScreen,
		Position:  PositionSupport,
		Modifications: map[string]float64{
			"ProtectionRadius": 1.5, // +50% radius multiplier
			"MitigationBonus":  0.20, // +20% damage mitigation
		},
	},
	{
		AbilityID: AbilityOverload,
		Position:  PositionFront,
		Modifications: map[string]float64{
			"DamageMultiplier": 1.25, // +25% damage
			"ShieldPenalty":    -0.10, // -10% less shield penalty
		},
	},
	{
		AbilityID: AbilityInterdictorPulse,
		Position:  PositionFlank,
		Modifications: map[string]float64{
			"RadiusMultiplier": 1.3, // +30% interdiction radius
			"DurationSeconds":  1.2, // +20% duration multiplier
		},
	},
	{
		AbilityID: AbilitySelfRepair,
		Position:  PositionSupport,
		Modifications: map[string]float64{
			"RegenRateBonus": 0.50, // +50% repair rate
		},
	},
	{
		AbilityID: AbilityLongRangeSensors,
		Position:  PositionBack,
		Modifications: map[string]float64{
			"VisibilityBonus": 2.0, // +2 visibility
		},
	},
	{
		AbilityID: AbilityAdaptiveTargeting,
		Position:  PositionFront,
		Modifications: map[string]float64{
			"SwitchSpeed":      -0.30, // 30% faster switch
			"DamageMultiplier": 1.10,  // +10% damage
		},
	},
	{
		AbilityID: AbilitySiegePayload,
		Position:  PositionBack,
		Modifications: map[string]float64{
			"StructureDamageBonus": 0.20, // +20% structure damage
			"SplashRadiusBonus":    1.0,  // +1 splash radius
		},
	},
}

// GemPositionEffect defines bonus effects when a gem type is socketed in a ship at a specific formation position.
type GemPositionEffect struct {
	GemFamily GemFamily          `bson:"gemFamily" json:"gemFamily"`
	Position  FormationPosition  `bson:"position" json:"position"`
	Bonus     StatMods           `bson:"bonus" json:"bonus"`
}

// GemPositionEffectsCatalog contains synergies between gem families and formation positions.
var GemPositionEffectsCatalog = []GemPositionEffect{
	{
		GemFamily: GemLaser,
		Position:  PositionFront,
		Bonus: StatMods{
			LaserShieldDelta: 1,
			Damage:           DamageMods{LaserPct: 0.15},
		},
	},
	{
		GemFamily: GemNuclear,
		Position:  PositionFront,
		Bonus: StatMods{
			NuclearShieldDelta: 1,
			BucketHPPct:        0.10,
		},
	},
	{
		GemFamily: GemAntimatter,
		Position:  PositionFront,
		Bonus: StatMods{
			FirstVolleyPct: 0.15,
			CritPct:        0.05,
		},
	},
	{
		GemFamily: GemKinetic,
		Position:  PositionFront,
		Bonus: StatMods{
			BucketHPPct:           0.15,
			LaserShieldDelta:      1,
			NuclearShieldDelta:    1,
		},
	},
	{
		GemFamily: GemSensor,
		Position:  PositionBack,
		Bonus: StatMods{
			AttackRangeDelta: 1,
			VisibilityDelta:  2,
			AccuracyPct:      0.10,
		},
	},
	{
		GemFamily: GemWarp,
		Position:  PositionFlank,
		Bonus: StatMods{
			SpeedDelta:            1,
			WarpChargePct:         -0.10,
			InterdictionResistPct: 0.15,
		},
	},
	{
		GemFamily: GemEngineering,
		Position:  PositionSupport,
		Bonus: StatMods{
			OutOfCombatRegenPct: 0.20,
			AbilityCooldownPct:  -0.10,
		},
	},
	{
		GemFamily: GemLogistics,
		Position:  PositionSupport,
		Bonus: StatMods{
			TransportCapacityPct: 0.15,
			UpkeepPct:            -0.05,
		},
	},
	{
		GemFamily: GemLaser,
		Position:  PositionFlank,
		Bonus: StatMods{
			AccuracyPct: 0.08,
			SpeedDelta:  1,
		},
	},
	{
		GemFamily: GemSensor,
		Position:  PositionFlank,
		Bonus: StatMods{
			VisibilityDelta: 1,
			CloakDetect:     true,
		},
	},
	{
		GemFamily: GemAntimatter,
		Position:  PositionBack,
		Bonus: StatMods{
			Damage:         DamageMods{AntimatterPct: 0.12},
			ShieldPiercePct: 0.08,
		},
	},
	{
		GemFamily: GemKinetic,
		Position:  PositionSupport,
		Bonus: StatMods{
			BucketHPPct:         0.10,
			OutOfCombatRegenPct: 0.15,
		},
	},
}

// GetAbilityFormationMods returns the modifications for an ability when used from a specific position.
func GetAbilityFormationMods(abilityID AbilityID, position FormationPosition) map[string]float64 {
	for _, mod := range AbilityFormationModsCatalog {
		if mod.AbilityID == abilityID && mod.Position == position {
			return mod.Modifications
		}
	}
	return nil
}

// ApplyGemPositionEffects computes bonus stat mods from gem-position synergies.
func ApplyGemPositionEffects(gems []Gem, position FormationPosition) StatMods {
	mods := ZeroMods()
	
	for _, gem := range gems {
		for _, effect := range GemPositionEffectsCatalog {
			if effect.GemFamily == gem.Family && effect.Position == position {
				mods = CombineMods(mods, effect.Bonus)
			}
		}
	}
	
	return mods
}

// CompositionBonus represents a fleet composition bonus that activates when requirements are met.
type CompositionBonus struct {
	Type        string           `bson:"type" json:"type"`
	Description string           `bson:"description" json:"description"`
	Requirement map[ShipType]int `bson:"requirement" json:"requirement"` // min ships of type
	Bonus       StatMods         `bson:"bonus" json:"bonus"`
}

// CompositionBonusesCatalog defines synergies based on fleet composition.
var CompositionBonusesCatalog = []CompositionBonus{
	{
		Type:        "Balanced Fleet",
		Description: "Well-rounded fleet composition with scouts, fighters, and bombers.",
		Requirement: map[ShipType]int{
			Scout:   1,
			Fighter: 2,
			Bomber:  1,
		},
		Bonus: StatMods{
			SpeedDelta: 1,
			Damage:     DamageMods{LaserPct: 0.05, NuclearPct: 0.05, AntimatterPct: 0.05},
		},
	},
	{
		Type:        "Strike Force",
		Description: "Heavy assault specialization with fighters and destroyers.",
		Requirement: map[ShipType]int{
			Fighter:   3,
			Destroyer: 1,
		},
		Bonus: StatMods{
			Damage:  DamageMods{LaserPct: 0.15, NuclearPct: 0.15, AntimatterPct: 0.15},
			CritPct: 0.10,
		},
	},
	{
		Type:        "Siege Armada",
		Description: "Specialized siege fleet with bombers and carriers.",
		Requirement: map[ShipType]int{
			Bomber:  2,
			Carrier: 1,
		},
		Bonus: StatMods{
			StructureDamagePct: 0.25,
			AttackRangeDelta:   1,
			SplashRadiusDelta:  1,
		},
	},
	{
		Type:        "Recon Squadron",
		Description: "Intelligence-focused fleet with scouts and sensor capabilities.",
		Requirement: map[ShipType]int{
			Scout: 3,
		},
		Bonus: StatMods{
			VisibilityDelta: 3,
			SpeedDelta:      2,
			CloakDetect:     true,
			PingRangePct:    0.30,
		},
	},
	{
		Type:        "Mobile Fortress",
		Description: "Defensive powerhouse with carriers and destroyers.",
		Requirement: map[ShipType]int{
			Carrier:   1,
			Destroyer: 2,
		},
		Bonus: StatMods{
			LaserShieldDelta:      2,
			NuclearShieldDelta:    2,
			AntimatterShieldDelta: 2,
			BucketHPPct:           0.15,
		},
	},
	{
		Type:        "Hit and Run",
		Description: "Fast strike force with scouts and fighters.",
		Requirement: map[ShipType]int{
			Scout:   2,
			Fighter: 2,
		},
		Bonus: StatMods{
			SpeedDelta:  2,
			AccuracyPct: 0.15,
			Damage:      DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
		},
	},
	{
		Type:        "Economic Convoy",
		Description: "Resource-focused fleet with drones and carriers.",
		Requirement: map[ShipType]int{
			Drone:   3,
			Carrier: 1,
		},
		Bonus: StatMods{
			TransportCapacityPct: 0.30,
			UpkeepPct:            -0.10,
			ConstructionCostPct:  -0.05,
		},
	},
	{
		Type:        "Rapid Response",
		Description: "Fast deployment force with warp-capable ships.",
		Requirement: map[ShipType]int{
			Bomber:    1,
			Carrier:   1,
			Destroyer: 1,
		},
		Bonus: StatMods{
			WarpChargePct:         -0.20,
			WarpScatterPct:        -0.25,
			InterdictionResistPct: 0.15,
		},
	},
}

// EvaluateCompositionBonuses checks which composition bonuses are active for the given fleet.
func EvaluateCompositionBonuses(ships map[ShipType][]HPBucket) (StatMods, []CompositionBonus) {
	mods := ZeroMods()
	var activeBonuses []CompositionBonus
	
	// Count total ships per type
	counts := make(map[ShipType]int)
	for shipType, buckets := range ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		counts[shipType] = total
	}
	
	// Check each composition bonus
	for _, bonus := range CompositionBonusesCatalog {
		requirementsMet := true
		for reqType, reqCount := range bonus.Requirement {
			if counts[reqType] < reqCount {
				requirementsMet = false
				break
			}
		}
		
		if requirementsMet {
			mods = CombineMods(mods, bonus.Bonus)
			activeBonuses = append(activeBonuses, bonus)
		}
	}
	
	return mods, activeBonuses
}

// FormationTemplate defines a pre-configured formation setup with conditions.
type FormationTemplate struct {
	Name        string                         `bson:"name" json:"name"`
	Description string                         `bson:"description" json:"description"`
	Formation   FormationType                  `bson:"formation" json:"formation"`
	Assignments map[ShipType]FormationPosition `bson:"assignments" json:"assignments"`
	Conditions  []TemplateCondition            `bson:"conditions" json:"conditions"`
}

// TemplateCondition defines requirements for using a formation template.
type TemplateCondition struct {
	MinShips map[ShipType]int `bson:"minShips" json:"minShips"`
	RoleMode RoleMode         `bson:"roleMode,omitempty" json:"roleMode,omitempty"`
	Against  FormationType    `bson:"against,omitempty" json:"against,omitempty"` // Counter-formation
}

// FormationTemplatesCatalog contains pre-made formation configurations.
var FormationTemplatesCatalog = []FormationTemplate{
	{
		Name:        "Standard Battle Line",
		Description: "Fighters front, bombers back, balanced approach.",
		Formation:   FormationLine,
		Assignments: map[ShipType]FormationPosition{
			Fighter:   PositionFront,
			Destroyer: PositionFront,
			Bomber:    PositionBack,
			Carrier:   PositionSupport,
			Scout:     PositionFlank,
			Drone:     PositionSupport,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Fighter: 1},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Defensive Box",
		Description: "All-around defense for sieges and defensive operations.",
		Formation:   FormationBox,
		Assignments: map[ShipType]FormationPosition{
			Fighter:   PositionFront,
			Destroyer: PositionFront,
			Carrier:   PositionSupport,
			Bomber:    PositionBack,
			Scout:     PositionFlank,
			Drone:     PositionSupport,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Fighter: 2, Carrier: 1},
			},
		},
	},
	{
		Name:        "Blitz Vanguard",
		Description: "Aggressive alpha strike formation.",
		Formation:   FormationVanguard,
		Assignments: map[ShipType]FormationPosition{
			Destroyer: PositionFront,
			Fighter:   PositionFront,
			Bomber:    PositionSupport,
			Carrier:   PositionSupport,
			Scout:     PositionFlank,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Destroyer: 1, Fighter: 2},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Hit and Run Skirmish",
		Description: "Mobile strike force for quick engagements.",
		Formation:   FormationSkirmish,
		Assignments: map[ShipType]FormationPosition{
			Scout:     PositionFlank,
			Fighter:   PositionFlank,
			Destroyer: PositionFront,
			Bomber:    PositionBack,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Scout: 2, Fighter: 2},
			},
		},
	},
	{
		Name:        "Mining Operation",
		Description: "Economic formation optimized for resource gathering.",
		Formation:   FormationBox,
		Assignments: map[ShipType]FormationPosition{
			Drone:   PositionSupport,
			Carrier: PositionSupport,
			Fighter: PositionFront,
			Scout:   PositionFlank,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Drone: 2},
				RoleMode: RoleEconomic,
			},
		},
	},
	{
		Name:        "Recon Sweep",
		Description: "Scouting formation for intelligence gathering.",
		Formation:   FormationSwarm,
		Assignments: map[ShipType]FormationPosition{
			Scout:   PositionFlank,
			Fighter: PositionFront,
			Carrier: PositionSupport,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Scout: 2},
				RoleMode: RoleRecon,
			},
		},
	},
}

// FindBestTemplate selects the most appropriate formation template for the given conditions.
func FindBestTemplate(ships map[ShipType][]HPBucket, role RoleMode, enemyFormation FormationType) *FormationTemplate {
	// Count total ships per type
	counts := make(map[ShipType]int)
	for shipType, buckets := range ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		counts[shipType] = total
	}
	
	// Score each template
	var bestTemplate *FormationTemplate
	bestScore := 0
	
	for i := range FormationTemplatesCatalog {
		template := &FormationTemplatesCatalog[i]
		score := 0
		
		for _, condition := range template.Conditions {
			// Check ship requirements
			requirementsMet := true
			for reqType, reqCount := range condition.MinShips {
				if counts[reqType] < reqCount {
					requirementsMet = false
					break
				}
			}
			if !requirementsMet {
				continue
			}
			
			score += 10
			
			// Bonus for matching role
			if condition.RoleMode != "" && condition.RoleMode == role {
				score += 5
			}
			
			// Bonus for countering enemy formation
			if condition.Against != "" && condition.Against == enemyFormation {
				score += 8
			}
		}
		
		// Check formation counter advantage
		if enemyFormation != "" {
			counterMult := GetFormationCounterMultiplier(template.Formation, enemyFormation)
			if counterMult > 1.0 {
				score += int((counterMult - 1.0) * 20) // +20 score per 1.0 advantage
			}
		}
		
		if score > bestScore {
			bestScore = score
			bestTemplate = template
		}
	}
	
	return bestTemplate
}
