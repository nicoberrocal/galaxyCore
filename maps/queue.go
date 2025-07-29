package maps

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PlayerAction represents an action initiated by a player
type Queue struct {
	ID              bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	PlayerID        bson.ObjectID    `bson:"playerId" json:"playerId"`
	MapID           bson.ObjectID    `bson:"mapId" json:"mapId"`
	Type            string           `bson:"type" json:"type"`                             // ship_attack, ship_construction, building_construction, ship_ability
	TargetID        bson.ObjectID    `bson:"targetId,omitempty" json:"targetId,omitempty"` // Target of the action (e.g., ship, building)
	SourceID        bson.ObjectID    `bson:"sourceId,omitempty" json:"sourceId,omitempty"` // Source of the action (e.g., ship, building)
	StartX          float64          `bson:"startX,omitempty" json:"startX,omitempty"`     // Starting coordinates
	StartY          float64          `bson:"startY,omitempty" json:"startY"`
	TargetX         float64          `bson:"targetX,omitempty" json:"targetX"` // Target coordinates
	TargetY         float64          `bson:"targetY,omitempty" json:"targetY"`
	StartTime       time.Time        `bson:"startTime" json:"startTime"`
	EndTime         time.Time        `bson:"endTime,omitempty" json:"endTime"`
	CreatedAt       time.Time        `bson:"createdAt" json:"createdAt"`               // When the action was created
	ProcessedAt     time.Time        `bson:"processedAt,omitempty" json:"processedAt"` // When the action was processed
	PredictedEvents []PredictedEvent `bson:"predictedEvents,omitempty" json:"predictedEvents,omitempty"`
	Version         int64            `bson:"version" json:"version"`                     // For optimistic locking
	Payload         bson.D           `bson:"payload,omitempty" json:"payload,omitempty"` // Additional action-specific data
}

type ResponseQueue struct {
	ID                bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	MapID             bson.ObjectID    `bson:"mapId" json:"mapId"`
	PlayerID          bson.ObjectID    `bson:"playerId" json:"playerId"`
	QueueItemID       bson.ObjectID    `bson:"queueItemId" json:"queueItemId"`               // ID of the queue item this response is for
	Type              string           `bson:"type" json:"type"`                             // ship_attack, ship_construction, building_construction, ship_ability
	TargetID          bson.ObjectID    `bson:"targetId,omitempty" json:"targetId,omitempty"` // Target of the action (e.g., ship, building)
	SourceID          bson.ObjectID    `bson:"sourceId,omitempty" json:"sourceId,omitempty"` // Source of the action (e.g., ship, building)
	StartX            float64          `bson:"startX,omitempty" json:"startX,omitempty"`     // Starting coordinates
	StartY            float64          `bson:"startY,omitempty" json:"startY"`
	TargetX           float64          `bson:"targetX,omitempty" json:"targetX"` // Target coordinates
	TargetY           float64          `bson:"targetY,omitempty" json:"targetY"`
	StartTime         time.Time        `bson:"startTime" json:"startTime"`
	EndTime           time.Time        `bson:"endTime,omitempty" json:"endTime"`
	CreatedAt         time.Time        `bson:"createdAt" json:"createdAt"`                 // When the action was created
	ProcessedAt       time.Time        `bson:"processedAt,omitempty" json:"processedAt"`   // When the action was processed
	Version           int64            `bson:"version" json:"version"`                     // For optimistic locking
	Payload           bson.D           `bson:"payload,omitempty" json:"payload,omitempty"` // Additional action-specific data
	PredictedEvents   []PredictedEvent `bson:"predictedEvents,omitempty" json:"predictedEvents,omitempty"`
	EventParticipants []bson.ObjectID  `bson:"eventParticipants,omitempty" json:"eventParticipants,omitempty"` // Participants in the event
}

type PredictedEvent struct {
	Type      string        `bson:"type"` // "collision", "attack_range"
	Timestamp time.Time     `bson:"timestamp"`
	WithID    bson.ObjectID `bson:"withId"`
	WithType  string        `bson:"withType"`
	AtX       float64       `bson:"atX"`
	AtY       float64       `bson:"atY"`
}
