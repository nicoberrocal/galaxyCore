package orbitables

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResourceExtraction struct {
	IsBeingMined   bool             `bson:"isBeingMined"`
	ResourceType   string           `bson:"resourceType"`   // "metal", "crystal", "hydrogen"
	MiningFleets   []bson.ObjectID  `bson:"miningFleets"`   // Stack IDs currently mining
	ExtractionRate map[string]int64 `bson:"extractionRate"` // Resources per tick per fleet
	RemainingRes   map[string]int64 `bson:"remainingRes"`   // Remaining resources
	LastExtracted  time.Time        `bson:"lastExtracted"`  // Last time resources were extracted
}

type Nebula struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	MapID   bson.ObjectID `bson:"mapId,omitempty"`
	X       float64       `bson:"x"`
	Y       float64       `bson:"y"`
	Name    string        `bson:"name"`
	Texture string        `bson:"texture"`

	// Hydrogen extraction - separate stacks can extract simultaneously
	ResourceExtraction ResourceExtraction `bson:"resourceExtraction"`

	// Collision detection for extraction
	CollisionRadius float64 `bson:"collisionRadius"` // Radius for determining extraction range

	// Available resources
	TotalHydrogen int64 `bson:"totalHydrogen"` // Total hydrogen available
}
