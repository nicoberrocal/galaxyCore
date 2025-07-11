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
	Movement []*MovementState `bson:"movement,omitempty" json:"movement,omitempty"`
	Battle   *BattleState     `bson:"battle,omitempty" json:"battle,omitempty"` // Combat state for free space battles
}

type HPBucket struct {
	HP    int `bson:"hp"`    // full HP of that ship
	Count int `bson:"count"` // how many ships at this HP
}

// BattleState tracks combat information for stacks in free space or mining locations
type BattleState struct {
	IsInCombat      bool            `bson:"isInCombat"`                // Currently engaged in battle
	EnemyStackID    []bson.ObjectID `bson:"enemyStackId,omitempty"`    // Opponent stack ID
	EnemyPlayerID   []bson.ObjectID `bson:"enemyPlayerId,omitempty"`   // Opponent player ID
	BattleStartedAt time.Time       `bson:"battleStartedAt,omitempty"` // When battle began
	BattleLocation  string          `bson:"battleLocation,omitempty"`  // "empty_space", "asteroid", "nebula"
	LocationID      bson.ObjectID   `bson:"locationId,omitempty"`      // ID of asteroid/nebula if applicable
}

// MovementState tracks what the stack is currently doing in free space or at mining locations
type MovementState struct {
	State      string        `bson:"state" json:"state,omitempty"`                     // "traveling", "mining", "idle", "in_combat"
	TargetID   bson.ObjectID `bson:"targetId,omitempty" json:"targetId,omitempty"`     // Target for movement/mining
	TargetType string        `bson:"targetType,omitempty" json:"targetType,omitempty"` // "asteroid", "nebula", "coordinate"
	StartX     float64       `bson:"startX,omitempty" json:"startX,omitempty"`         // Starting coordinates
	StartY     float64       `bson:"startY,omitempty" json:"startY"`                   // Starting coordinates
	Speed      int           `bson:"speed,omitempty" json:"speed"`                     // Speed of the stack
	TargetX    float64       `bson:"targetX,omitempty" json:"targetX"`                 // Target coordinates
	TargetY    float64       `bson:"targetY,omitempty" json:"targetY"`                 // Target coordinates
	StartTime  time.Time     `bson:"startTime,omitempty" json:"startTime"`             // When current action started
	EndTime    time.Time     `bson:"endTime,omitempty" json:"endTime"`                 // When current action ends
	Activity   string        `bson:"activity,omitempty" json:"activity"`               // "mining_metal", "mining_crystal", "mining_hydrogen"
}
