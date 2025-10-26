package essences

import (
	"time"

	"github.com/nicoberrocal/galaxyCore/ships"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/////////////////////
// Biology (BioTree)
/////////////////////

type BioTreeType string

const (
	Aquatica BioTreeType = "aquatica"
	Flora    BioTreeType = "flora"
	Fauna    BioTreeType = "fauna"
	Mycelia  BioTreeType = "mycelia"
)

// AoETarget defines area of effect targeting for abilities
// This allows effects to target multiple stacks within a radius
type AoETraitTarget struct {
	Radius     float64       // radius in units
	TargetType AoETargetType // allies, enemies, all, specific ship types
	MaxTargets int           // maximum number of targets (0 = unlimited)
	Origin     AoEOrigin     // center point of the AoE
}

type AoETargetType string

const (
	AoEAllies   AoETargetType = "allies"
	AoEEnemies  AoETargetType = "enemies"
	AoEAll      AoETargetType = "all"
	AoESpecific AoETargetType = "specific"
)

type AoEOrigin string

const (
	AoESelf     AoEOrigin = "self"
	AoETarget   AoEOrigin = "target"
	AoEPosition AoEOrigin = "position"
	AoEGlobal   AoEOrigin = "global"
)

// StatusEffect represents a debuff/buff that can be applied to stacks
// These are more complex than simple stat modifiers and have duration/stacking mechanics
type StatusEffect struct {
	ID            string
	Name          string
	Description   string
	Duration      int // ticks
	MaxStacks     int
	CurrentStacks int
	AppliedAt     time.Time
	IsBeneficial  bool // true for buffs, false for debuffs
	EffectType    StatusEffectType
}

type StatusEffectType string

const (
	StatusStun         StatusEffectType = "stun"
	StatusRoot         StatusEffectType = "root"
	StatusBlind        StatusEffectType = "blind"
	StatusConfusion    StatusEffectType = "confusion"
	StatusFear         StatusEffectType = "fear"
	StatusInfection    StatusEffectType = "infection"
	StatusSlow         StatusEffectType = "slow"
	StatusWeaken       StatusEffectType = "weaken"
	StatusVulnerable   StatusEffectType = "vulnerable"
	StatusShielded     StatusEffectType = "shielded"
	StatusRegenerating StatusEffectType = "regenerating"
	StatusEnraged      StatusEffectType = "enraged"
	StatusInvisible    StatusEffectType = "invisible"
)

// SpawnEffect represents the creation of new entities (drones, husks, micro-stacks, etc.)
type SpawnEffect struct {
	SpawnType     SpawnType
	SpawnTemplate interface{} // Template for the spawned entity
	Duration      int         // How long the spawned entity lasts (ticks)
	SpawnCount    int         // How many to spawn
	SpawnRadius   float64     // Radius around origin to spawn
}

type SpawnType string

const (
	SpawnDecoyDrone SpawnType = "decoy_drone"
	SpawnSporeHusk  SpawnType = "spore_husk"
	SpawnMicroStack SpawnType = "micro_stack"
	SpawnShockwave  SpawnType = "shockwave"
	SpawnSporeCloud SpawnType = "spore_cloud"
	SpawnOvergrowth SpawnType = "overgrowth_field"
	SpawnAcidEffect SpawnType = "acid_effect"
	SpawnFearEffect SpawnType = "fear_effect"
)

// ComplexEffect represents effects that go beyond simple stat modifications
// This includes conditional logic, delayed effects, and multi-stage effects
type ComplexEffect struct {
	EffectType      ComplexEffectType
	Trigger         Trigger         // Event that triggers this effect
	Conditions      []Condition     // Conditions that must be met for effect to trigger
	PrimaryEffect   *ships.StatMods `bson:"primaryEffect,omitempty" json:"primaryEffect,omitempty"`     // Main stat modifications
	SecondaryEffect *ships.StatMods `bson:"secondaryEffect,omitempty" json:"secondaryEffect,omitempty"` // Secondary effects (e.g., on-death effects)
	AoE             *AoETraitTarget // AoE targeting if applicable
	Spawn           *SpawnEffect    // Spawn mechanics if applicable
	StatusEffects   []StatusEffect  // Status effects to apply
	Duration        int             // Effect duration in ticks
	Cooldown        int             // Cooldown between activations
	MaxActivations  int             // Maximum number of times this can trigger (0 = unlimited)
	ActivationCount int             // Current activation count
	IsActive        bool            // Whether effect is currently active
	ActivatedAt     time.Time       // When effect was last activated
}

type ComplexEffectType string

const (
	ComplexConditional   ComplexEffectType = "conditional"
	ComplexDelayed       ComplexEffectType = "delayed"
	ComplexMultiStage    ComplexEffectType = "multi_stage"
	ComplexOnDeath       ComplexEffectType = "on_death"
	ComplexPeriodic      ComplexEffectType = "periodic"
	ComplexChainReaction ComplexEffectType = "chain_reaction"
)

// Condition represents a requirement that must be met for an effect to trigger
type Condition struct {
	ConditionType ConditionType
	Value         interface{}     // The value to check against
	CompareOp     ComparisonOp    // How to compare
	Target        ConditionTarget // What to check the condition on
}

type ConditionType string

const (
	ConditionHPPercent             ConditionType = "hp_percent"
	ConditionDistance              ConditionType = "distance"
	ConditionStackCount            ConditionType = "stack_count"
	ConditionAllyCount             ConditionType = "ally_count"
	ConditionEnemyCount            ConditionType = "enemy_count"
	ConditionResourceNear          ConditionType = "resource_near"
	ConditionTerrainNear           ConditionType = "terrain_near"
	ConditionFormationType         ConditionType = "formation_type"
	ConditionAttackType            ConditionType = "attack_type"
	ConditionHasStatus             ConditionType = "has_status"
	ConditionCombatState           ConditionType = "combat_state"
	ConditionTickCount             ConditionType = "tick_count"
	ConditionAbilityUsed           ConditionType = "ability_used"
	ConditionDamageReceived        ConditionType = "damage_received"
	ConditionKillCount             ConditionType = "kill_count"
	ConditionMovementState         ConditionType = "movement_state"
	ConditionAttackFromBehind      ConditionType = "attack_from_behind"
	ConditionCriticalHit           ConditionType = "critical_hit"
	ConditionConsecutiveAttacks    ConditionType = "consecutive_attacks"
	ConditionSystemLost            ConditionType = "system_lost"
	ConditionAllyNearby            ConditionType = "ally_nearby"
	ConditionIsAttacked            ConditionType = "is_attacked"
	ConditionStationary            ConditionType = "stationary"
	ConditionTargetInfected        ConditionType = "target_infected"
	ConditionBuildingInfected      ConditionType = "building_infected"
	ConditionAllyInNetwork         ConditionType = "ally_in_network"
	ConditionTargetIsAttackingAlly ConditionType = "target_is_attacking_ally"
)

type ComparisonOp string

const (
	CompareEqual     ComparisonOp = "equal"
	CompareNotEqual  ComparisonOp = "not_equal"
	CompareGreater   ComparisonOp = "greater"
	CompareLess      ComparisonOp = "less"
	CompareGreaterEq ComparisonOp = "greater_equal"
	CompareLessEq    ComparisonOp = "less_equal"
	CompareContains  ComparisonOp = "contains"
	CompareInRange   ComparisonOp = "in_range"
)

type ConditionTarget string

const (
	TargetSelf     ConditionTarget = "self"
	TargetAttacker ConditionTarget = "attacker"
	TargetTarget   ConditionTarget = "target"
	TargetArea     ConditionTarget = "area"
	TargetGlobal   ConditionTarget = "global"
)

// BioNode is a single selectable node in a biology tree.
type BioNode struct {
	ID          string
	Title       string
	Description string
	Path        string
	Effect      ships.StatMods `bson:"effect,omitempty" json:"effect,omitempty"`
	Effects     *BioNodeEffects
	Triggers    *[]Trigger
	Tradeoff    *ships.StatMods `bson:"tradeoff,omitempty" json:"tradeoff,omitempty"`
	// Complex effects for advanced biotree mechanics
	ComplexEffects []ComplexEffect
	// StatDelta for compatibility with existing essences system
	StatDelta StatDelta
	// Mutations is an optional map where an EssenceType may replace or alter this node.
	// If an essence key exists, it provides the mutated node instance.

}

type Trigger string

const (
	TriggerOnFormationChange      Trigger = "during_formation_change"
	TriggerOnAbilityCast          Trigger = "after_ability_cast"
	TriggerOnCombatStart          Trigger = "at_combat_start"
	TriggerOnCombatEnd            Trigger = "at_combat_end"
	TriggerOnVisibilityRangeEnter Trigger = "when_visibility_range_entered"
	TriggerOnMovementSprint       Trigger = "on_movement_sprint"
	TriggerOnStackDestroyed       Trigger = "on_stack_destroyed"
	TriggerOnSystemLost           Trigger = "on_system_lost"
	TriggerOnDamageReceived       Trigger = "on_damage_received"
	TriggerOnCriticalHit          Trigger = "on_critical_hit"
	TriggerOnAttackFromBehind     Trigger = "on_attack_from_behind"
	TriggerOnEnemyEnterRange      Trigger = "on_enemy_enter_range"
	TriggerOnStationary           Trigger = "on_stationary"
	TriggerOnSuccessfulHit        Trigger = "on_successful_hit"
	TriggerOnStackDeath           Trigger = "on_stack_death"
	TriggerOnKill                 Trigger = "on_kill"
	TriggerOnAllyDeath            Trigger = "on_ally_death"
	TriggerOnEnemyDeath           Trigger = "on_enemy_death"
	TriggerOnDeath                Trigger = "on_death"
	TriggerOnResourceExtract      Trigger = "on_resource_extract"
	TriggerOnSystemEngage         Trigger = "on_system_engage"
	TriggerOnTick                 Trigger = "on_tick"
	TriggerOnFormationChangeComplete Trigger = "on_formation_change_complete"
	TriggerOnFirstStrike          Trigger = "on_first_strike"
	TriggerOnAllyDamaged          Trigger = "on_ally_damaged"
	TriggerOnNearAsteroid         Trigger = "on_near_asteroid"
	TriggerOnNearStar             Trigger = "on_near_star"
	TriggerOnActiveAbility        Trigger = "on_active_ability"
	TriggerOnConsecutiveAttacks   Trigger = "on_consecutive_attacks"
	TriggerOnLowHP                Trigger = "on_low_hp"
	TriggerOnAllyNearby           Trigger = "on_ally_nearby"
	TriggerOnTerrainNearby        Trigger = "on_terrain_nearby"
	TriggerOnSystemEngaged        Trigger = "on_system_engaged"
	TriggerOnCombatEnter          Trigger = "on_combat_enter"
	TriggerOnAllyWouldDie         Trigger = "on_ally_would_die"
	TriggerOnShipDeathInArea      Trigger = "on_ship_death_in_area"
	TriggerOnAbilityCooldown      Trigger = "on_ability_cooldown"
)

type NodeTrigger struct {
	Trigger       Trigger
	TriggerRange  float64
	Conditions    []Condition
	Cooldown      int // ticks
	LastTriggered time.Time
}

type BioNodeEffects struct {
	Target           string
	TargetRange      bool
	TargetRangeValue float64
	Effect           *[]ships.StatMods `bson:"effect,omitempty" json:"effect,omitempty"`
	ComplexEffects   []ComplexEffect
	StatusEffects    []StatusEffect
}

// BioTree contains the full structured tree.
type BioTree struct {
	Name        string
	Description string
	// Tiers is an ordered list of tiers; each tier is a set of nodes where some are
	// mutually exclusive. We model tiers as slices of nodes (choices happen in-game).
	Tiers [][]*BioNode
}

type BioTreeState struct {
	PlayerID    bson.ObjectID `bson:"playerId" json:"playerId"`
	BioTreeType BioTreeType   `bson:"bioTreeType" json:"bioTreeType"`
	// Experience tracking
	TotalXP     int `bson:"totalXP" json:"totalXP"`         // Lifetime earned
	SpentXP     int `bson:"spentXP" json:"spentXP"`         // Used for nodes
	AvailableXP int `bson:"availableXP" json:"availableXP"` // Ready to spend

	// Active nodes
	UnlockedNodes []string `bson:"unlockedNodes" json:"unlockedNodes"` // Node IDs

	// Meta info
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	Version   int       `bson:"version" json:"version"` // For optimistic locking
}

/////////////////////////
// Biology trees (examples)
/////////////////////////

// For compactness we create helper node constructors.
func NewNode(id, title, desc, path string, eff, trade ships.StatMods) *BioNode {
	return &BioNode{
		ID:          id,
		Title:       title,
		Description: desc,
		Effect:      eff,
		Path:        path,
		Tradeoff:    &trade,
	}
}

// EvaluateTriggerAndCondition checks if both trigger and conditions are met for a bio effect
// This provides a general trigger + condition evaluation system
func EvaluateTriggerAndCondition(
	nodeID string,
	trigger Trigger,
	conditions []Condition,
	eventData interface{},
) bool {
	// First check if trigger is activated by the current event
	if !IsTriggerActive(trigger, eventData) {
		return false
	}
	
	// Then check all conditions
	if !AreConditionsMet(conditions, eventData) {
		return false
	}
	
	// Both trigger and conditions are met
	return true
}

// IsTriggerActive checks if the given trigger is activated by the current event
func IsTriggerActive(trigger Trigger, eventData interface{}) bool {
	switch trigger {
	case TriggerOnFirstStrike:
		// This would be triggered by first attack events
		return true // Placeholder - would need actual first strike detection
	case TriggerOnSuccessfulHit:
		// This would be triggered by attack events
		return true // Placeholder - would need actual hit detection
	case TriggerOnCriticalHit:
		// This would be triggered by critical hit events
		return true // Placeholder - would need actual crit detection
	case TriggerOnStationary:
		// This would be triggered by movement state events
		return true // Placeholder - would need actual stationary detection
	case TriggerOnTick:
		// This is triggered every tick
		return true
	case TriggerOnDeath:
		// This would be triggered by death events
		return true // Placeholder - would need actual death detection
	case TriggerOnKill:
		// This would be triggered by kill events
		return true // Placeholder - would need actual kill detection
	case TriggerOnAbilityCast:
		// This would be triggered by ability cast events
		return true // Placeholder - would need actual ability cast detection
	case TriggerOnActiveAbility:
		// This would be triggered by active ability usage
		return true // Placeholder - would need actual active ability detection
	case TriggerOnFormationChangeComplete:
		// This would be triggered by formation change completion
		return true // Placeholder - would need actual formation change detection
	case TriggerOnNearAsteroid:
		// This would be triggered by terrain proximity events
		return true // Placeholder - would need actual terrain detection
	case TriggerOnNearStar:
		// This would be triggered by terrain proximity events
		return true // Placeholder - would need actual terrain detection
	case TriggerOnEnemyEnterRange:
		// This would be triggered by enemy proximity events
		return true // Placeholder - would need actual proximity detection
	case TriggerOnAllyNearby:
		// This would be triggered by ally proximity events
		return true // Placeholder - would need actual proximity detection
	case TriggerOnAllyDeath:
		// This would be triggered by ally death events
		return true // Placeholder - would need actual ally death detection
	case TriggerOnEnemyDeath:
		// This would be triggered by enemy death events
		return true // Placeholder - would need actual enemy death detection
	case TriggerOnSystemEngage:
		// This would be triggered by system engagement events
		return true // Placeholder - would need actual system engagement detection
	case "":
		// No trigger specified - always active
		return true
	default:
		// Unknown trigger - default to not active
		return false
	}
}

// AreConditionsMet checks if all conditions are satisfied
func AreConditionsMet(conditions []Condition, eventData interface{}) bool {
	for _, condition := range conditions {
		if !IsConditionMet(condition, eventData) {
			return false
		}
	}
	return true
}

// IsConditionMet checks if a single condition is satisfied
func IsConditionMet(condition Condition, eventData interface{}) bool {
	switch condition.ConditionType {
	case ConditionCombatState:
		if state, ok := condition.Value.(string); ok {
			// Placeholder - would need actual combat state detection
			// For now, return true if state is "engaging" as that's commonly used
			return state == "engaging"
		}
	case ConditionCriticalHit:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual critical hit detection
			return value
		}
	case ConditionStationary:
		if value, ok := condition.Value.(int); ok {
			// Placeholder - would need actual stationary detection with tick count
			// For now, return true if value is >= 3 as commonly used
			return value >= 3
		}
	case ConditionTargetInfected:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual infection status detection
			return value
		}
	case ConditionKillCount:
		if value, ok := condition.Value.(int); ok {
			// Placeholder - would need actual kill count detection
			// For now, return true if value > 0 as commonly used
			return value > 0
		}
	case ConditionAllyCount:
		if value, ok := condition.Value.(int); ok {
			// Placeholder - would need actual ally count detection
			// For now, return true if value >= 3 as commonly used
			return value >= 3
		}
	case ConditionTerrainNear:
		if value, ok := condition.Value.(string); ok {
			// Placeholder - would need actual terrain detection
			// For now, return true for any terrain type
			return value == "asteroid" || value == "star"
		}
	case ConditionFormationType:
		if value, ok := condition.Value.(string); ok {
			// Placeholder - would need actual formation type detection
			// For now, return true for common formation types
			return value == "aggressive" || value == "defensive" || value == "Box"
		}
	case ConditionAbilityUsed:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual ability usage detection
			return value
		}
	case ConditionAllyInNetwork:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual network status detection
			return value
		}
	case ConditionBuildingInfected:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual building infection detection
			return value
		}
	case ConditionIsAttacked:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual attack status detection
			return value
		}
	case ConditionAttackFromBehind:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual attack direction detection
			return value
		}
	case ConditionTargetIsAttackingAlly:
		if value, ok := condition.Value.(bool); ok {
			// Placeholder - would need actual target attack detection
			return value
		}
	}
	return false
}
