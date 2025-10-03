package ships

// Formation composition bonuses.
// This file contains the data catalogs for fleet composition bonuses.

// GemPositionEffect is DEPRECATED - removed for clean system separation.
// Gems provide their own StatMods only, without position synergies.
type GemPositionEffect struct {
	GemFamily GemFamily         `bson:"gemFamily" json:"gemFamily"`
	Position  FormationPosition `bson:"position" json:"position"`
	Bonus     StatMods          `bson:"bonus" json:"bonus"`
}

// GemPositionEffectsCatalog is DEPRECATED - removed for clean system separation.
// Kept for backward compatibility but should not be used.
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
			BucketHPPct:        0.15,
			LaserShieldDelta:   1,
			NuclearShieldDelta: 1,
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
			Damage:          DamageMods{AntimatterPct: 0.12},
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

// ApplyGemPositionEffects is DEPRECATED - removed for clean system separation.
// Gems provide their own StatMods only, without position synergies.
// This function is kept for backward compatibility but returns zero mods.
func ApplyGemPositionEffects(gems []Gem, position FormationPosition) StatMods {
	// DEPRECATED: Return zero mods
	return ZeroMods()
}

// CompositionBonus is DEPRECATED - removed for clean system separation.
// Fleet composition bonuses create implicit synergies between ship types.
// Each ship should contribute independently.
type CompositionBonus struct {
	Type        string           `bson:"type" json:"type"`
	Description string           `bson:"description" json:"description"`
	Requirement map[ShipType]int `bson:"requirement" json:"requirement"` // min ships of type
	Bonus       StatMods         `bson:"bonus" json:"bonus"`
}

// CompositionBonusesCatalog is DEPRECATED - removed for clean system separation.
// Kept for backward compatibility but should not be used.
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
	{
		Type:        "Tank Division",
		Description: "Heavy armor formation with Cruisers and Carriers.",
		Requirement: map[ShipType]int{
			Cruiser: 2,
			Carrier: 1,
		},
		Bonus: StatMods{
			LaserShieldDelta:      1,
			NuclearShieldDelta:    1,
			AntimatterShieldDelta: 1,
			BucketHPPct:           0.20,
		},
	},
	{
		Type:        "Scout Hunter Pack",
		Description: "Corvette-focused fleet for countering fast ships.",
		Requirement: map[ShipType]int{
			Corvette: 2,
		},
		Bonus: StatMods{
			SpeedDelta:  2,
			AccuracyPct: 0.20,
			Damage:      DamageMods{AntimatterPct: 0.15},
		},
	},
	{
		Type:        "Artillery Battery",
		Description: "AoE-focused fleet with Artillery and support.",
		Requirement: map[ShipType]int{
			Artillery: 1,
			Bomber:    1,
		},
		Bonus: StatMods{
			AttackRangeDelta:   2,
			SplashRadiusDelta:  1,
			StructureDamagePct: 0.20,
		},
	},
	{
		Type:        "Shadow Ops",
		Description: "Stealth-focused fleet for assassination missions.",
		Requirement: map[ShipType]int{
			StealthFrigate: 2,
			Scout:          1,
		},
		Bonus: StatMods{
			CritPct:         0.20,
			FirstVolleyPct:  0.25,
			VisibilityDelta: 2,
		},
	},
	{
		Type:        "Electronic Warfare Wing",
		Description: "Debuff-stacking fleet with Support Frigates.",
		Requirement: map[ShipType]int{
			Frigate: 2,
		},
		Bonus: StatMods{
			AbilityCooldownPct: -0.15,
			AccuracyPct:        0.15,
			VisibilityDelta:    1,
		},
	},
	{
		Type:        "Antimatter Supremacy",
		Description: "Antimatter-focused fleet for shield penetration.",
		Requirement: map[ShipType]int{
			Destroyer: 1,
			Corvette:  2,
		},
		Bonus: StatMods{
			Damage:          DamageMods{AntimatterPct: 0.20},
			ShieldPiercePct: 0.15,
			FirstVolleyPct:  0.15,
		},
	},
	{
		Type:        "Combined Arms",
		Description: "Diverse fleet with all new ship types.",
		Requirement: map[ShipType]int{
			Cruiser:        1,
			Corvette:       1,
			Artillery:      1,
			StealthFrigate: 1,
			Frigate:        1,
		},
		Bonus: StatMods{
			Damage:                DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
			LaserShieldDelta:      1,
			NuclearShieldDelta:    1,
			AntimatterShieldDelta: 1,
			SpeedDelta:            1,
		},
	},
}

// EvaluateCompositionBonuses is DEPRECATED - removed for clean system separation.
// Fleet composition bonuses create implicit synergies between ship types.
// This function is kept for backward compatibility but returns zero mods.
func EvaluateCompositionBonuses(ships map[ShipType][]HPBucket) (StatMods, []CompositionBonus) {
	return ZeroMods(), []CompositionBonus{}
}

// ComputeCompositionBonuses is DEPRECATED - removed for clean system separation.
// Fleet composition bonuses create implicit synergies between ship types.
// This function is kept for backward compatibility but returns zero mods.
func ComputeCompositionBonuses(ships map[ShipType][]HPBucket) (StatMods, []CompositionBonus) {
	return ZeroMods(), []CompositionBonus{}
}

// countShipsByType counts the total number of ships per type in the fleet.
func countShipsByType(ships map[ShipType][]HPBucket) map[ShipType]int {
	counts := make(map[ShipType]int)
	for shipType, buckets := range ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		counts[shipType] = total
	}
	return counts
}

// meetsRequirements checks if the ship counts meet the composition requirements.
func meetsRequirements(counts map[ShipType]int, requirements map[ShipType]int) bool {
	for reqType, reqCount := range requirements {
		if counts[reqType] < reqCount {
			return false
		}
	}
	return true
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
	{
		Name:        "Tank Wall",
		Description: "Heavy frontline formation with Cruisers and Carriers.",
		Formation:   FormationBox,
		Assignments: map[ShipType]FormationPosition{
			Cruiser:   PositionFront,
			Carrier:   PositionSupport,
			Frigate:   PositionSupport,
			Artillery: PositionBack,
			Fighter:   PositionFront,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Cruiser: 2, Carrier: 1},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Scout Hunter",
		Description: "Corvette-focused formation for countering fast ships.",
		Formation:   FormationSkirmish,
		Assignments: map[ShipType]FormationPosition{
			Corvette:  PositionFlank,
			Destroyer: PositionFront,
			Frigate:   PositionSupport,
			Fighter:   PositionFlank,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Corvette: 2},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Artillery Barrage",
		Description: "Long-range AoE formation with Artillery and Bombers.",
		Formation:   FormationLine,
		Assignments: map[ShipType]FormationPosition{
			Artillery: PositionBack,
			Bomber:    PositionBack,
			Cruiser:   PositionFront,
			Frigate:   PositionSupport,
			Fighter:   PositionFront,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Artillery: 1, Bomber: 1},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Stealth Strike",
		Description: "Assassination formation with Stealth Frigates.",
		Formation:   FormationSkirmish,
		Assignments: map[ShipType]FormationPosition{
			StealthFrigate: PositionFlank,
			Scout:          PositionFlank,
			Corvette:       PositionFlank,
			Fighter:        PositionFront,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{StealthFrigate: 2},
				RoleMode: RoleTactical,
			},
		},
	},
	{
		Name:        "Debuff Stack",
		Description: "Electronic warfare formation with multiple Support Frigates.",
		Formation:   FormationLine,
		Assignments: map[ShipType]FormationPosition{
			Frigate: PositionSupport,
			Cruiser: PositionFront,
			Carrier: PositionSupport,
			Fighter: PositionFront,
		},
		Conditions: []TemplateCondition{
			{
				MinShips: map[ShipType]int{Frigate: 2},
				RoleMode: RoleTactical,
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
