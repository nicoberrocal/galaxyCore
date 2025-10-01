# Formation Combat System Design Document

## 1. Overview

### 1.1 Purpose
Enhance tactical combat depth by introducing formation mechanics that leverage the existing HP bucket system, allowing mixed-ship type positioning at the bucket level for sophisticated battle strategies.

### 1.2 Design Philosophy
- **Granular Control**: Formation assignments at HP bucket level
- **Hybrid Stacking**: Maintain performance benefits of stacking while enabling tactical positioning
- **Seamless Integration**: Build upon existing gem, role mode, and ability systems
- **Strategic Depth**: Multiple layers of counter-play and composition optimization

## 2. Core Concepts

### 2.1 Formation Components

#### Formation Types
```go
type FormationType string

const (
    FormationLine       FormationType = "line"      // Balanced front-back arrangement
    FormationBox        FormationType = "box"       // Defensive all-around protection
    FormationVanguard   FormationType = "vanguard"  // Aggressive forward deployment
    FormationSkirmish   FormationType = "skirmish"  // Mobile flanking focus
    FormationEchelon    FormationType = "echelon"   // Diagonal staggered lines
    FormationPhalanx    FormationType = "phalanx"   // Heavy frontal concentration
    FormationSwarm      FormationType = "swarm"     // Dispersed anti-AoE formation
)
```
Geometric Properties
Formation	Shape	Front Width	Depth	Flank Exposure
Line	Rectangle	Wide	Shallow	High
Box	Square	Medium	Medium	Low
Vanguard	Triangle	Narrow	Deep	Medium
Skirmish	Circular	Variable	Variable	None
Echelon	Parallelogram	Medium	Deep	Asymmetric
Phalanx	Trapezoid	Very Wide	Deep	Extreme
Swarm	Scattered	N/A	N/A	N/A
Arrowhead	Triangle	Narrow	Deep	Medium

#### Formation Positions
```go
type FormationPosition string

const (
    PositionFront   FormationPosition = "front"    // Primary damage absorption
    PositionFlank   FormationPosition = "flank"    // Mobile strike forces
    PositionBack    FormationPosition = "back"     // Ranged/support units
    PositionSupport FormationPosition = "support"  // Utility/healer positions
)
```

### 2.2 Mixed HP Bucket Assignments

#### Core Data Structure
```go
type FormationAssignment struct {
    Position    FormationPosition `bson:"position" json:"position"`
    Layer       int               `bson:"layer" json:"layer"`           // 0=frontline, 1=mid, 2=backline
    ShipType    ShipType          `bson:"shipType" json:"shipType"`
    BucketIndex int               `bson:"bucketIndex" json:"bucketIndex"` // Index in ship type's HP buckets
    Count       int               `bson:"count" json:"count"`             // Ships from this bucket
    AssignedHP  int               `bson:"assignedHP" json:"assignedHP"`   // Current HP of assigned ships
}

type Formation struct {
    Type        FormationType              `bson:"formationType" json:"formationType"`
    Facing      string                     `bson:"facing" json:"facing"`   // "north", "south", "east", "west"
    Assignments []FormationAssignment      `bson:"assignments" json:"assignments"`
    Modifiers   FormationMods              `bson:"modifiers" json:"modifiers"`
    CreatedAt   time.Time                  `bson:"createdAt" json:"createdAt"`
    Version     int                        `bson:"version" json:"version"`
}
```

## 3. Formation Properties & Bonuses

### 3.1 Formation Type Modifiers

```go
type FormationMods struct {
    SpeedMultiplier    float64            `bson:"speedMultiplier" json:"speedMultiplier"`
    ReconfigureTime    int                `bson:"reconfigureTime" json:"reconfigureTime"`
    PositionBonuses    map[FormationPosition]FormationBonus `bson:"positionBonuses" json:"positionBonuses"`
    SpecialProperties  []string           `bson:"specialProperties" json:"specialProperties"`
}

type FormationBonus struct {
    DefenseBonus    map[string]float64 `bson:"defenseBonus" json:"defenseBonus"`       // shield type -> multiplier
    OffenseBonus    map[string]float64 `bson:"offenseBonus" json:"offenseBonus"`       // stat -> multiplier
    AbilityEnhance  []AbilityEnhance   `bson:"abilityEnhance" json:"abilityEnhance"`   // ability modifications
}
```

### 3.2 Formation Catalog

#### Line Formation
- **Front**: +15% Laser Shield, +10% Attack Damage
- **Flank**: +10% Speed, +5% Critical Chance  # Why +% speed? To make slow ships faster, which increases the speed movement of the entire stack, since the slowest speed is the one applied for movement in the movement system. Not a battle trait, but logistic/tactical one.
- **Back**: +20% Attack Range, +15% Visibility
- **Special**: Strong vs frontal attacks, weak to flanking

#### Box Formation
- **All Positions**: +10% All Shields, -25% Speed
- **Special**: Even damage distribution, excellent vs siege

#### Vanguard Formation
- **Front**: +25% Attack Damage, +15% Nuclear Shield
- **Support**: -20% HP, +30% Ability Effectiveness
- **Special**: Fast reconfiguration (60s), aggressive

#### Skirmish Formation
- **Flank**: +20% Speed, +15% Evasion Chance
- **Front**: +10% Attack Damage, -10% HP
- **Special**: Mobile, excellent for hit-and-run

## 4. Damage Distribution System

### 4.1 Position-Based Damage Routing

```go
func (s *ShipStack) CalculateDamageDistribution(incomingDamage int, attackType string, direction string) map[FormationAssignment]int {
    distribution := make(map[FormationAssignment]int)
    
    // Step 1: Calculate position damage weights
    positionWeights := s.calculatePositionWeights(direction)
    
    // Step 2: Distribute to positions
    for position, weight := range positionWeights {
        positionDamage := int(float64(incomingDamage) * weight)
        assignments := s.getAssignmentsByPosition(position)
        
        // Step 3: Distribute within position by ship count and type
        for _, assignment := range assignments {
            assignmentDamage := s.calculateAssignmentDamage(positionDamage, assignment, assignments)
            distribution[assignment] = assignmentDamage
        }
    }
    
    return distribution
}
```

### 4.2 Directional Attack Modifiers

| Attack Direction | Front | Flank | Back | Support |
|------------------|-------|-------|------|---------|
| Frontal          | 60%   | 20%   | 10%  | 10%     |
| Flanking         | 30%   | 40%   | 20%  | 10%     |
| Rear             | 10%   | 30%   | 50%  | 10%     |
| Envelopment      | 25%   | 25%   | 25%  | 25%     |

### 4.3 Bucket-Level Damage Application

```go
func (s *ShipStack) ApplyFormationDamage(damageMap map[FormationAssignment]int) {
    // Group damage by ship type and bucket
    bucketDamage := make(map[ShipType]map[int]int)
    
    for assignment, damage := range damageMap {
        if _, exists := bucketDamage[assignment.ShipType]; !exists {
            bucketDamage[assignment.ShipType] = make(map[int]int)
        }
        bucketDamage[assignment.ShipType][assignment.BucketIndex] += damage
    }
    
    // Apply damage to buckets
    for shipType, buckets := range bucketDamage {
        for bucketIndex, damage := range buckets {
            s.applyDamageToBucket(shipType, bucketIndex, damage)
        }
    }
    
    // Clean up destroyed assignments
    s.cleanupFormationAssignments()
}
```

## 5. Formation Counters & Interactions

### 5.1 Rock-Paper-Scissors Matrix

| Attacker → Defender | Line | Box | Vanguard | Skirmish | Echelon |
|---------------------|------|-----|----------|----------|---------|
| **Line**            | 1.0x | 0.8x| 1.3x     | 0.9x     | 1.1x    |
| **Box**             | 1.2x | 1.0x| 0.7x     | 1.1x     | 0.9x    |
| **Vanguard**        | 0.7x | 1.3x| 1.0x     | 1.4x     | 0.8x    |
| **Skirmish**        | 1.1x | 0.9x| 0.6x     | 1.0x     | 1.2x    |
| **Echelon**         | 0.9x | 1.1x| 1.2x     | 0.8x     | 1.0x    |

### 5.2 Composition-Based Modifiers

```go
type CompositionBonus struct {
    Type          string  `bson:"type" json:"type"`
    Requirement   map[ShipType]int `bson:"requirement" json:"requirement"` // min ships of type
    Bonus         map[string]float64 `bson:"bonus" json:"bonus"`
    Description   string  `bson:"description" json:"description"`
}

var CompositionBonuses = []CompositionBonus{
    {
        Type: "Balanced Fleet",
        Requirement: map[ShipType]int{Scout: 1, Fighter: 2, Bomber: 1},
        Bonus: map[string]float64{"Speed": 0.1, "AttackDamage": 0.05},
        Description: "Well-rounded fleet composition",
    },
    {
        Type: "Strike Force", 
        Requirement: map[ShipType]int{Fighter: 3, Destroyer: 1},
        Bonus: map[string]float64{"AttackDamage": 0.15, "CriticalChance": 0.1},
        Description: "Heavy assault specialization",
    },
    // ... more composition bonuses
}
```

## 6. Integration with Existing Systems

### 6.1 Role Mode Interactions

**Tactical Mode**
- Formation changes: -30% reconfiguration time
- Position bonuses: +10% effectiveness
- Counter bonuses: +0.1x multiplier

**Economic Mode** 
- Formation changes: +50% reconfiguration time
- Mining formations available
- Defensive position bonuses enhanced

**Recon Mode**
- Enemy formation visibility
- Flanking detection bonuses
- Faster formation spotting

**Scientific Mode**
- Experimental formations available
- Ability to copy enemy formations
- Research bonuses for formation development

### 6.2 Ability Integration

#### Formation-Enhanced Abilities
```go
// Ability modifications based on formation position
type AbilityFormationMod struct {
    AbilityID     AbilityID          `bson:"abilityId" json:"abilityId"`
    Position      FormationPosition  `bson:"position" json:"position"`
    Modifications map[string]float64 `bson:"modifications" json:"modifications"`
}

var AbilityFormationMods = []AbilityFormationMod{
    {
        AbilityID: AbilityFocusFire,
        Position:  PositionFront,
        Modifications: map[string]float64{
            "CooldownSeconds": -0.5,
            "DamageMultiplier": 1.2,
        },
    },
    {
        AbilityID: AbilityCloakWhileAnchored, 
        Position:  PositionFlank,
        Modifications: map[string]float64{
            "DurationSeconds": 2.0,
            "DetectionResistance": 0.3,
        },
    },
}
```

### 6.3 Gem Socket Synergies

#### Position-Specific Gem Effects
```go
type GemPositionEffect struct {
    GemType      GemType            `bson:"gemType" json:"gemType"`
    Position     FormationPosition  `bson:"position" json:"position"`
    Bonus        map[string]float64 `bson:"bonus" json:"bonus"`
}

var GemPositionEffects = []GemPositionEffect{
    {
        GemType:  GemRuby,
        Position: PositionFront,
        Bonus:    map[string]float64{"LaserShield": 0.15, "HP": 0.1},
    },
    {
        GemType:  GemSapphire,
        Position: PositionBack, 
        Bonus:    map[string]float64{"AttackRange": 0.2, "VisibilityRange": 0.15},
    },
}
```

## 7. Management & UI Systems

### 7.1 Formation Assignment Interface

#### Drag & Drop Bucket Management
```
[Formation Canvas]
┌────────────────────────────────────────┐
│ Front Line                             │
│ ┌─────┐ ┌─────┐ ┌─────┐               │
│ │F 50 │ │F 30 │ │D 20 │               │
│ └─────┘ └─────┘ └─────┘               │
│                                        │
│ Flank Positions                        │
│ ┌─────┐ ┌─────┐        Back Line      │
│ │S 40 │ │S 20 │        ┌─────┐ ┌─────┐│
│ └─────┘ └─────┘        │B 15 │ │C 1  ││
│                        └─────┘ └─────┘│
└────────────────────────────────────────┘

[Available Buckets]
Fighter (100/100) [Split] [Assign All]
Scout (60/60) [Split] [Assign All]  
Bomber (15/15) [Split] [Assign All]
Carrier (1/1) [Split] [Assign All]
```

### 7.2 Auto-Assignment Algorithms

```go
// Smart assignment based on ship type strengths
func (s *ShipStack) AutoAssignOptimal(formation FormationType) Formation {
    newFormation := Formation{Type: formation}
    
    // Priority queue for assignment
    assignments := s.prioritizeAssignments(formation)
    
    // Fill positions based on capacity and optimization
    for _, assignment := range assignments {
        position := s.calculateOptimalPosition(assignment.ShipType, formation)
        layer := s.getLayerForPosition(position)
        
        newFormation.Assignments = append(newFormation.Assignments, FormationAssignment{
            Position:    position,
            Layer:       layer,
            ShipType:    assignment.ShipType,
            BucketIndex: assignment.BucketIndex,
            Count:       assignment.Count,
            AssignedHP:  assignment.AssignedHP,
        })
    }
    
    return newFormation
}
```

### 7.3 Formation Templates

```go
type FormationTemplate struct {
    ID          bson.ObjectID                  `bson:"_id,omitempty"`
    Name        string                         `bson:"name" json:"name"`
    Description string                         `bson:"description" json:"description"`
    Formation   FormationType                  `bson:"formation" json:"formation"`
    Assignments map[ShipType]FormationPosition `bson:"assignments" json:"assignments"`
    Conditions  []TemplateCondition            `bson:"conditions" json:"conditions"`
}

type TemplateCondition struct {
    MinShips    map[ShipType]int `bson:"minShips" json:"minShips"`
    RoleMode    RoleMode         `bson:"roleMode" json:"roleMode"`
    Against     FormationType    `bson:"against" json:"against"`
}
```

## 8. Battle Resolution Flow

### 8.1 Combat Sequence with Formations

1. **Pre-Battle Phase**
   - Detect enemy formation and composition
   - Calculate formation counters
   - Apply pre-battle ability effects

2. **Damage Calculation Phase**
   - Determine attack direction based on positioning
   - Apply formation RPS multipliers
   - Calculate position-based damage distribution

3. **Damage Application Phase** 
   - Apply damage to specific bucket assignments
   - Handle bucket splitting and destruction
   - Update formation assignments

4. **Post-Battle Phase**
   - Experience and resource calculations
   - Formation integrity checks
   - Ability cooldown updates

### 8.2 Example Combat Scenario

**Situation**: Player A (Vanguard Formation) attacks Player B (Box Formation)

1. **Counter Calculation**: Vanguard vs Box = 1.3x damage multiplier
2. **Direction**: Frontal attack (Player A's Vanguard vs Player B's front line)
3. **Damage Distribution**: 
   - 70% to Player B's front positions
   - 20% to flank positions  
   - 10% to back positions
4. **Application**: Frontline buckets take heavy damage, may trigger bucket splitting
5. **Result**: Player A gains advantage due to formation counter

## 9. Balance Considerations



### 9.1 Strategic Balance

**Risk vs Reward**
- Specialized formations: High bonuses but predictable
- Mixed formations: Adaptive but lower peak performance
- Counter-picking: Formation scouting becomes valuable

**Economic Considerations**
- Formation changes cost time (30-180 seconds)
- Damaged formations require repair time
- Specialized formations may require research
