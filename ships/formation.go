package ships

import (
	"math"
	"time"
)

// FormationType defines the geometric arrangement of ships in a stack.
// Each formation has unique properties, bonuses, and counter-relationships.
type FormationType string

const (
	FormationLine     FormationType = "line"     // Balanced front-back arrangement
	FormationBox      FormationType = "box"      // Defensive all-around protection
	FormationVanguard FormationType = "vanguard" // Aggressive forward deployment
	FormationSkirmish FormationType = "skirmish" // Mobile flanking focus
	FormationEchelon  FormationType = "echelon"  // Diagonal staggered lines
	FormationPhalanx  FormationType = "phalanx"  // Heavy frontal concentration
	FormationSwarm    FormationType = "swarm"    // Dispersed anti-AoE formation
)

// FormationPosition defines where a ship bucket is positioned within the formation.
type FormationPosition string

const (
	PositionFront   FormationPosition = "front"   // Primary damage absorption
	PositionFlank   FormationPosition = "flank"   // Mobile strike forces
	PositionBack    FormationPosition = "back"    // Ranged/support units
	PositionSupport FormationPosition = "support" // Utility/healer positions
)

// AttackDirection defines the angle of attack in battle.
type AttackDirection string

const (
	DirectionFrontal     AttackDirection = "frontal"
	DirectionFlanking    AttackDirection = "flanking"
	DirectionRear        AttackDirection = "rear"
	DirectionEnvelopment AttackDirection = "envelopment"
)

// FormationAssignment maps a specific HP bucket to a position in the formation.
// This enables granular control where different ship types can occupy any position.
type FormationAssignment struct {
	Position    FormationPosition `bson:"position" json:"position"`
	Layer       int               `bson:"layer" json:"layer"` // 0=frontline, 1=mid, 2=backline
	ShipType    ShipType          `bson:"shipType" json:"shipType"`
	BucketIndex int               `bson:"bucketIndex" json:"bucketIndex"` // Index in ship type's HP buckets
	Count       int               `bson:"count" json:"count"`             // Ships from this bucket
	AssignedHP  int               `bson:"assignedHP" json:"assignedHP"`   // Current HP of assigned ships
}

// Formation represents the tactical arrangement of a ship stack.
type Formation struct {
	Type        FormationType         `bson:"formationType" json:"formationType"`
	Facing      string                `bson:"facing" json:"facing"` // "north", "south", "east", "west"
	Assignments []FormationAssignment `bson:"assignments" json:"assignments"`
	Modifiers   FormationMods         `bson:"modifiers" json:"modifiers"`
	CreatedAt   time.Time             `bson:"createdAt" json:"createdAt"`
	Version     int                   `bson:"version" json:"version"`
}

// FormationMods contains the modifiers applied by the formation type.
type FormationMods struct {
	SpeedMultiplier   float64                        `bson:"speedMultiplier" json:"speedMultiplier"`
	ReconfigureTime   int                            `bson:"reconfigureTime" json:"reconfigureTime"` // seconds
	PositionBonuses   map[FormationPosition]StatMods `bson:"positionBonuses" json:"positionBonuses"`
	SpecialProperties []string                       `bson:"specialProperties" json:"specialProperties"`
}

// FormationSpec defines the characteristics and bonuses of a formation type.
type FormationSpec struct {
	Type              FormationType
	Name              string
	Description       string
	SpeedMultiplier   float64
	ReconfigureTime   int // seconds
	PositionBonuses   map[FormationPosition]StatMods
	SpecialProperties []string
}

// FormationCatalog contains all formation type definitions.
var FormationCatalog = map[FormationType]FormationSpec{
	FormationLine: {
		Type:            FormationLine,
		Name:            "Line Formation",
		Description:     "Balanced front-back arrangement. Strong vs frontal attacks, weak to flanking.",
		SpeedMultiplier: 1.0,
		ReconfigureTime: 120,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				LaserShieldDelta: 1,
				Damage:           DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
			},
			PositionFlank: {
				SpeedDelta: 1,
				CritPct:    0.05,
			},
			PositionBack: {
				AttackRangeDelta: 1,
				VisibilityDelta:  1,
			},
		},
		SpecialProperties: []string{"frontal_strength", "flank_vulnerable"},
	},
	FormationBox: {
		Type:            FormationBox,
		Name:            "Box Formation",
		Description:     "Defensive all-around protection. Even damage distribution, excellent vs siege.",
		SpeedMultiplier: 0.75,
		ReconfigureTime: 150,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				LaserShieldDelta:      1,
				NuclearShieldDelta:    1,
				AntimatterShieldDelta: 1,
			},
			PositionFlank: {
				LaserShieldDelta:      1,
				NuclearShieldDelta:    1,
				AntimatterShieldDelta: 1,
			},
			PositionBack: {
				LaserShieldDelta:      1,
				NuclearShieldDelta:    1,
				AntimatterShieldDelta: 1,
			},
			PositionSupport: {
				LaserShieldDelta:      1,
				NuclearShieldDelta:    1,
				AntimatterShieldDelta: 1,
			},
		},
		SpecialProperties: []string{"even_distribution", "siege_resistant"},
	},
	FormationVanguard: {
		Type:            FormationVanguard,
		Name:            "Vanguard Formation",
		Description:     "Aggressive forward deployment. Fast reconfiguration, high front damage.",
		SpeedMultiplier: 1.1,
		ReconfigureTime: 60,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				Damage:             DamageMods{LaserPct: 0.25, NuclearPct: 0.25, AntimatterPct: 0.25},
				NuclearShieldDelta: 1,
			},
			PositionSupport: {
				BucketHPPct:        -0.20,
				AbilityCooldownPct: -0.30,
			},
		},
		SpecialProperties: []string{"fast_reconfig", "aggressive"},
	},
	FormationSkirmish: {
		Type:            FormationSkirmish,
		Name:            "Skirmish Formation",
		Description:     "Mobile flanking focus. Excellent for hit-and-run tactics.",
		SpeedMultiplier: 1.2,
		ReconfigureTime: 90,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFlank: {
				SpeedDelta:  2,
				AccuracyPct: 0.15,
				Damage:      DamageMods{LaserPct: 0.20, NuclearPct: 0.20, AntimatterPct: 0.20},
			},
			PositionFront: {
				Damage:      DamageMods{LaserPct: 0.10, NuclearPct: 0.10, AntimatterPct: 0.10},
				BucketHPPct: -0.10,
			},
		},
		SpecialProperties: []string{"mobile", "hit_and_run"},
	},
	FormationEchelon: {
		Type:            FormationEchelon,
		Name:            "Echelon Formation",
		Description:     "Diagonal staggered lines. Asymmetric flank exposure, good vs concentrated attacks.",
		SpeedMultiplier: 0.95,
		ReconfigureTime: 120,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				LaserShieldDelta: 1,
				Damage:           DamageMods{LaserPct: 0.12, NuclearPct: 0.12, AntimatterPct: 0.12},
			},
			PositionFlank: {
				SpeedDelta: 1,
				CritPct:    0.08,
			},
			PositionBack: {
				AttackRangeDelta: 1,
				AccuracyPct:      0.05,
			},
		},
		SpecialProperties: []string{"asymmetric", "concentrated_defense"},
	},
	FormationPhalanx: {
		Type:            FormationPhalanx,
		Name:            "Phalanx Formation",
		Description:     "Heavy frontal concentration. Very wide front, extreme flank exposure.",
		SpeedMultiplier: 0.8,
		ReconfigureTime: 180,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				LaserShieldDelta:      2,
				NuclearShieldDelta:    2,
				AntimatterShieldDelta: 1,
				BucketHPPct:           0.15,
				Damage:                DamageMods{LaserPct: 0.15, NuclearPct: 0.15, AntimatterPct: 0.15},
			},
			PositionBack: {
				AttackRangeDelta: 2,
			},
		},
		SpecialProperties: []string{"frontal_fortress", "extreme_flank_weakness"},
	},
	FormationSwarm: {
		Type:            FormationSwarm,
		Name:            "Swarm Formation",
		Description:     "Dispersed anti-AoE formation. Reduces splash damage effectiveness.",
		SpeedMultiplier: 1.05,
		ReconfigureTime: 100,
		PositionBonuses: map[FormationPosition]StatMods{
			PositionFront: {
				SpeedDelta: 1,
			},
			PositionFlank: {
				SpeedDelta: 1,
			},
			PositionBack: {
				SpeedDelta: 1,
			},
		},
		SpecialProperties: []string{"dispersed", "anti_aoe", "splash_resistant"},
	},
}

// FormationCounterMatrix defines the rock-paper-scissors relationships between formations.
// Values are damage multipliers: attacker formation vs defender formation.
var FormationCounterMatrix = map[FormationType]map[FormationType]float64{
	FormationLine: {
		FormationLine:     1.0,
		FormationBox:      0.8,
		FormationVanguard: 1.3,
		FormationSkirmish: 0.9,
		FormationEchelon:  1.1,
		FormationPhalanx:  0.85,
		FormationSwarm:    1.0,
	},
	FormationBox: {
		FormationLine:     1.2,
		FormationBox:      1.0,
		FormationVanguard: 0.7,
		FormationSkirmish: 1.1,
		FormationEchelon:  0.9,
		FormationPhalanx:  1.15,
		FormationSwarm:    1.05,
	},
	FormationVanguard: {
		FormationLine:     0.7,
		FormationBox:      1.3,
		FormationVanguard: 1.0,
		FormationSkirmish: 1.4,
		FormationEchelon:  0.8,
		FormationPhalanx:  0.75,
		FormationSwarm:    1.2,
	},
	FormationSkirmish: {
		FormationLine:     1.1,
		FormationBox:      0.9,
		FormationVanguard: 0.6,
		FormationSkirmish: 1.0,
		FormationEchelon:  1.2,
		FormationPhalanx:  1.3,
		FormationSwarm:    0.95,
	},
	FormationEchelon: {
		FormationLine:     0.9,
		FormationBox:      1.1,
		FormationVanguard: 1.2,
		FormationSkirmish: 0.8,
		FormationEchelon:  1.0,
		FormationPhalanx:  0.9,
		FormationSwarm:    1.05,
	},
	FormationPhalanx: {
		FormationLine:     1.15,
		FormationBox:      0.85,
		FormationVanguard: 1.25,
		FormationSkirmish: 0.7,
		FormationEchelon:  1.1,
		FormationPhalanx:  1.0,
		FormationSwarm:    0.8,
	},
	FormationSwarm: {
		FormationLine:     1.0,
		FormationBox:      0.95,
		FormationVanguard: 0.8,
		FormationSkirmish: 1.05,
		FormationEchelon:  0.95,
		FormationPhalanx:  1.2,
		FormationSwarm:    1.0,
	},
}

// DirectionalDamageWeights defines how damage is distributed across positions based on attack direction.
var DirectionalDamageWeights = map[AttackDirection]map[FormationPosition]float64{
	DirectionFrontal: {
		PositionFront:   0.60,
		PositionFlank:   0.20,
		PositionBack:    0.10,
		PositionSupport: 0.10,
	},
	DirectionFlanking: {
		PositionFront:   0.30,
		PositionFlank:   0.40,
		PositionBack:    0.20,
		PositionSupport: 0.10,
	},
	DirectionRear: {
		PositionFront:   0.10,
		PositionFlank:   0.30,
		PositionBack:    0.50,
		PositionSupport: 0.10,
	},
	DirectionEnvelopment: {
		PositionFront:   0.25,
		PositionFlank:   0.25,
		PositionBack:    0.25,
		PositionSupport: 0.25,
	},
}

// GetFormationCounterMultiplier returns the damage multiplier for attacker vs defender formations.
func GetFormationCounterMultiplier(attacker, defender FormationType) float64 {
	if matrix, ok := FormationCounterMatrix[attacker]; ok {
		if mult, ok := matrix[defender]; ok {
			return mult
		}
	}
	return 1.0 // default: no bonus/penalty
}

// CalculateDamageDistribution computes how incoming damage is distributed across formation positions.
func (f *Formation) CalculateDamageDistribution(incomingDamage int, direction AttackDirection) map[FormationPosition]int {
	distribution := make(map[FormationPosition]int)

	weights, ok := DirectionalDamageWeights[direction]
	if !ok {
		weights = DirectionalDamageWeights[DirectionFrontal] // default
	}

	// Distribute damage to positions based on directional weights
	for position, weight := range weights {
		positionDamage := int(float64(incomingDamage) * weight)
		if positionDamage > 0 {
			distribution[position] = positionDamage
		}
	}

	return distribution
}

// GetAssignmentsByPosition returns all assignments for a given position.
func (f *Formation) GetAssignmentsByPosition(position FormationPosition) []FormationAssignment {
	var assignments []FormationAssignment
	for _, assignment := range f.Assignments {
		if assignment.Position == position {
			assignments = append(assignments, assignment)
		}
	}
	return assignments
}

// CalculateAssignmentDamage distributes position damage across assignments in that position.
// Distribution is proportional to ship count and HP.
func CalculateAssignmentDamage(positionDamage int, assignment FormationAssignment, allAssignments []FormationAssignment) int {
	if len(allAssignments) == 0 {
		return 0
	}

	// Calculate total HP weight in this position
	totalWeight := 0
	for _, a := range allAssignments {
		totalWeight += a.AssignedHP
	}

	if totalWeight == 0 {
		return 0
	}

	// Distribute proportionally by HP
	weight := float64(assignment.AssignedHP) / float64(totalWeight)
	return int(float64(positionDamage) * weight)
}

// ApplyPositionBonusesToShip applies formation position bonuses to a ship's stat mods.
func (f *Formation) ApplyPositionBonusesToShip(position FormationPosition, baseMods StatMods) StatMods {
	spec, ok := FormationCatalog[f.Type]
	if !ok {
		return baseMods
	}

	if posBonus, ok := spec.PositionBonuses[position]; ok {
		return CombineMods(baseMods, posBonus)
	}

	return baseMods
}

// GetEffectiveSpeed returns the formation's effective speed multiplier.
func (f *Formation) GetEffectiveSpeed(baseSpeed int) int {
	spec, ok := FormationCatalog[f.Type]
	if !ok {
		return baseSpeed
	}
	return int(math.Round(float64(baseSpeed) * spec.SpeedMultiplier))
}

// AutoAssignFormation automatically assigns ship buckets to optimal positions in the formation.
func AutoAssignFormation(ships map[ShipType][]HPBucket, formationType FormationType, now time.Time) Formation {
	formation := Formation{
		Type:        formationType,
		Facing:      "north",
		Assignments: []FormationAssignment{},
		CreatedAt:   now,
		Version:     1,
	}

	spec, ok := FormationCatalog[formationType]
	if ok {
		formation.Modifiers = FormationMods{
			SpeedMultiplier:   spec.SpeedMultiplier,
			ReconfigureTime:   spec.ReconfigureTime,
			PositionBonuses:   spec.PositionBonuses,
			SpecialProperties: spec.SpecialProperties,
		}
	}

	// Auto-assignment logic: assign ship types to positions based on their characteristics
	for shipType, buckets := range ships {
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 || bucket.HP == 0 {
				continue
			}

			position := DetermineOptimalPosition(shipType, formationType)
			layer := DetermineLayer(position, shipType)

			formation.Assignments = append(formation.Assignments, FormationAssignment{
				Position:    position,
				Layer:       layer,
				ShipType:    shipType,
				BucketIndex: bucketIndex,
				Count:       bucket.Count,
				AssignedHP:  bucket.HP * bucket.Count,
			})
		}
	}

	return formation
}

// DetermineOptimalPosition assigns a ship type to the most suitable formation position.
func DetermineOptimalPosition(shipType ShipType, formationType FormationType) FormationPosition {
	blueprint, ok := ShipBlueprints[shipType]
	if !ok {
		return PositionFront // default
	}

	// Position assignment based on ship characteristics
	switch shipType {
	case Drone:
		return PositionSupport // Drones are support/economic units

	case Scout:
		// Scouts are fast and good for flanking/recon
		if formationType == FormationSkirmish || formationType == FormationSwarm {
			return PositionFlank
		}
		return PositionFlank

	case Fighter:
		// Fighters are versatile front-line combatants
		if formationType == FormationVanguard {
			return PositionFront
		}
		return PositionFront

	case Bomber:
		// Bombers are long-range siege units, belong in back
		return PositionBack

	case Carrier:
		// Carriers are tanky support platforms
		if formationType == FormationBox {
			return PositionFront // Use their bulk in defensive formations
		}
		return PositionSupport

	case Destroyer:
		// Destroyers are heavy hitters, can be front or flank
		if formationType == FormationVanguard || formationType == FormationPhalanx {
			return PositionFront
		}
		// High speed makes them good flankers
		if blueprint.Speed >= 6 {
			return PositionFlank
		}
		return PositionFront
	}

	return PositionFront // default fallback
}

// DetermineLayer assigns a layer (depth) based on position and ship type.
func DetermineLayer(position FormationPosition, shipType ShipType) int {
	switch position {
	case PositionFront:
		return 0 // frontline
	case PositionFlank:
		return 1 // mid-line
	case PositionBack:
		return 2 // backline
	case PositionSupport:
		return 1 // mid-line
	default:
		return 0
	}
}

// RoleModeFormationBonus applies role-specific modifiers to formation reconfiguration and effectiveness.
func RoleModeFormationBonus(role RoleMode, reconfigTime int) int {
	switch role {
	case RoleTactical:
		// Tactical mode: -30% reconfiguration time
		return int(float64(reconfigTime) * 0.70)
	case RoleEconomic:
		// Economic mode: +50% reconfiguration time
		return int(float64(reconfigTime) * 1.50)
	case RoleRecon:
		// Recon mode: faster spotting, no penalty
		return reconfigTime
	case RoleScientific:
		// Scientific mode: normal speed
		return reconfigTime
	default:
		return reconfigTime
	}
}

// ApplyFormationRoleModifiers combines formation position bonuses with role mode bonuses.
func ApplyFormationRoleModifiers(baseMods StatMods, formation *Formation, position FormationPosition, role RoleMode) StatMods {
	// Apply formation position bonus
	mods := formation.ApplyPositionBonusesToShip(position, baseMods)

	// Apply role-specific enhancements to formation bonuses
	switch role {
	case RoleTactical:
		// +10% effectiveness to position bonuses
		enhancementMods := StatMods{
			Damage: DamageMods{
				LaserPct:      mods.Damage.LaserPct * 0.10,
				NuclearPct:    mods.Damage.NuclearPct * 0.10,
				AntimatterPct: mods.Damage.AntimatterPct * 0.10,
			},
		}
		mods = CombineMods(mods, enhancementMods)

	case RoleEconomic:
		// Defensive position bonuses enhanced when in economic mode
		if position == PositionSupport || position == PositionBack {
			defenseMods := StatMods{
				LaserShieldDelta:   1,
				NuclearShieldDelta: 1,
			}
			mods = CombineMods(mods, defenseMods)
		}

	case RoleRecon:
		// Enhanced visibility and detection in all positions
		reconMods := StatMods{
			VisibilityDelta: 1,
		}
		mods = CombineMods(mods, reconMods)

	case RoleScientific:
		// No specific formation bonuses in scientific mode
	}

	return mods
}
