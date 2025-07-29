package ships

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShipType string

const (
	Drone     ShipType = "drone"
	Scout     ShipType = "scout"
	Fighter   ShipType = "fighter"
	Bomber    ShipType = "bomber"
	Carrier   ShipType = "carrier"
	Destroyer ShipType = "destroyer"
)

// ShipStack represents a fleet that is NOT currently defending a system
// When a stack colonizes a system, it gets embedded in the system's DefendingFleet
// and this document is deleted (hybrid approach)
// Battles in free space/mining locations are handled by updating this document directly
type ShipStack struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	PlayerID  bson.ObjectID `bson:"playerId"`
	MapID     bson.ObjectID `bson:"mapId"`
	PositionX float64       `bson:"x"`
	PositionY float64       `bson:"y"`

	// Fleet composition
	Ships     map[ShipType][]HPBucket `bson:"ships"`     // HP bucketed ships
	CreatedAt time.Time               `bson:"createdAt"` // tick timestamp

	// Current activity and movement
	Movement  []*MovementState `bson:"movement,omitempty" json:"movement,omitempty"`
	Battle    *BattleState     `bson:"battle,omitempty" json:"battle,omitempty"`       // Combat state for free space battles
	Ability   *[]AbilityState  `bson:"ability,omitempty" json:"ability,omitempty"`     // Active ship ability state
	Gathering *GatheringState  `bson:"gathering,omitempty" json:"gathering,omitempty"` // Active gathering state
	Version   int64            `bson:"version" json:"version"`                         // For optimistic locking
}

type HPBucket struct {
	HP    int `bson:"hp"`    // full HP of that ship
	Count int `bson:"count"` // how many ships at this HP
}

// BattleState tracks combat information for stacks in free space or mining locations
type BattleState struct {
	IsInCombat      bool            `bson:"isInCombat"`                               // Currently engaged in battle
	EnemyStackID    []bson.ObjectID `bson:"enemyStackId,omitempty"`                   // Opponent stack ID
	EnemyPlayerID   []bson.ObjectID `bson:"enemyPlayerId,omitempty"`                  // Opponent player ID
	BattleStartedAt time.Time       `bson:"battleStartedAt,omitempty"`                // When battle began
	BattleLocation  string          `bson:"battleLocation,omitempty"`                 // "empty_space", "asteroid", "nebula"
	LocationID      bson.ObjectID   `bson:"locationId,omitempty"`                     // ID of asteroid/nebula if applicable
	ProcessedAt     time.Time       `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}

// MovementState tracks what the stack is currently doing in free space or at mining locations
type MovementState struct {
	State       string        `bson:"state" json:"state,omitempty"`                     // "traveling", "mining", "idle", "in_combat"
	TargetID    bson.ObjectID `bson:"targetId,omitempty" json:"targetId,omitempty"`     // Target for movement/mining
	TargetType  string        `bson:"targetType,omitempty" json:"targetType,omitempty"` // "asteroid", "nebula", "coordinate"
	StartX      float64       `bson:"startX,omitempty" json:"startX,omitempty"`         // Starting coordinates
	StartY      float64       `bson:"startY,omitempty" json:"startY"`                   // Starting coordinates
	Speed       int           `bson:"speed,omitempty" json:"speed"`                     // Speed of the stack
	TargetX     float64       `bson:"targetX,omitempty" json:"targetX"`                 // Target coordinates
	TargetY     float64       `bson:"targetY,omitempty" json:"targetY"`                 // Target coordinates
	StartTime   time.Time     `bson:"startTime,omitempty" json:"startTime"`             // When current action started
	EndTime     time.Time     `bson:"endTime,omitempty" json:"endTime"`                 // When current action ends
	Activity    string        `bson:"activity,omitempty" json:"activity"`               // "mining_metal", "mining_crystal", "mining_hydrogen"
	LastUpdated time.Time     `bson:"lastUpdated,omitempty" json:"lastUpdated"`         // Last time this state was updated
	ProcessedAt time.Time     `bson:"ProcessedAt,omitempty" json:"ProcessedAt"`         // Last time this state was processed
}

type AbilityState struct {
	IsActive    bool           `bson:"isActive" json:"isActive"`                 // Whether the ability is currently active
	Description string         `bson:"description" json:"description"`           // Description of the ability
	Icon        string         `bson:"icon" json:"icon"`                         // Icon representing the ability
	ShipType    ShipType       `bson:"shipType" json:"shipType"`                 // Type of ship using the ability
	Ability     string         `bson:"ability" json:"ability"`                   // Ability name (e.g., "cloak", "boost")
	Bonus       map[string]int `bson:"bonus" json:"bonus"`                       // Bonus type (e.g., "speed", "defense" or "attack") and its value
	StartTime   time.Time      `bson:"startTime" json:"startTime"`               // When the ability was activated
	EndTime     time.Time      `bson:"endTime" json:"endTime"`                   // When the ability will end
	Duration    int64          `bson:"duration" json:"duration"`                 // Duration of the ability in seconds
	LastUpdated time.Time      `bson:"lastUpdated" json:"lastUpdated"`           // Last time the ability state was updated
	ProcessedAt time.Time      `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}

type GatheringState struct {
	IsMining            bool          `bson:"isMining" json:"isMining"`                 // Whether the stack is currently mining
	TargetID            bson.ObjectID `bson:"targetId" json:"targetId"`                 // ID of the target being gathered from
	TargetType          string        `bson:"targetType" json:"targetType"`             // Type of target (e.g., "asteroid", "nebula")
	ResourcePerTimeUnit int64         `bson:"resourcePerTime" json:"resourcePerTime"`   // Amount of resource gathered per time unit
	TimeUnit            int64         `bson:"timeUnit" json:"timeUnit"`                 // Time unit in seconds for gathering
	TotalGathered       int64         `bson:"totalGathered" json:"totalGathered"`       // Total resources gathered
	StartTime           time.Time     `bson:"startTime" json:"startTime"`               // When gathering started
	ProcessedAt         time.Time     `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}
