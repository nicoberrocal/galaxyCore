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
	Conditions      []Condition     // Conditions that must be met for effect to trigger
	PrimaryEffect   *ships.StatMods // Main stat modifications
	SecondaryEffect *ships.StatMods // Secondary effects (e.g., on-death effects)
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
	Effect      ships.StatMods
	Effects     *BioNodeEffects
	Triggers    *[]Trigger
	Tradeoff    *ships.StatMods
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
	TriggerOnConsecutiveAttacks   Trigger = "on_consecutive_attacks"
	TriggerOnLowHP                Trigger = "on_low_hp"
	TriggerOnAllyNearby           Trigger = "on_ally_nearby"
	TriggerOnTerrainNearby        Trigger = "on_terrain_nearby"
	TriggerOnSystemEngaged        Trigger = "on_system_engaged"
	TriggerOnKill                 Trigger = "on_kill"
	TriggerOnCombatEnter          Trigger = "on_combat_enter"
	TriggerOnAllyWouldDie         Trigger = "on_ally_would_die"
	TriggerOnShipDeathInArea      Trigger = "on_ship_death_in_area"
	TriggerOnAbilityCooldown      Trigger = "on_ability_cooldown"
	TriggerOnTick                 Trigger = "on_tick"
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
	Effect           *[]ships.StatMods
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
