package maps

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PlayerAction represents an action initiated by a player
type PlayerAction struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	PlayerID    bson.ObjectID `bson:"playerId"`
	MapID       bson.ObjectID `bson:"mapId"`
	Type        string        `bson:"type"` // ship_attack, ship_construction, building_construction, ship_ability
	TargetID    bson.ObjectID `bson:"targetId,omitempty"`
	SourceID    bson.ObjectID `bson:"sourceId,omitempty"`
	X           float64       `bson:"x,omitempty"`
	Y           float64       `bson:"y,omitempty"`
	Finised     time.Time     `bson:"finished"`              // When the action should be processed
	CreatedAt   time.Time     `bson:"createdAt"`             // When the action was created
	ProcessedAt time.Time     `bson:"processedAt,omitempty"` // When the action was processed
	Version     int64         `bson:"version"`               // For optimistic locking
	Payload     bson.D        `bson:"payload,omitempty"`     // Additional action-specific data
}
