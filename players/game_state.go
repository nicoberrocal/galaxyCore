package players

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PlayerGameState tracks a player's state within a specific game/map
// This replaces the PlayerConfig embedded in maps for better normalization
type PlayerGameState struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	PlayerID bson.ObjectID `bson:"playerId"` // FK to players.Player
	MapID    bson.ObjectID `bson:"mapId"`    // FK to maps.MongoMap

	// Resources
	Energy           int64 `bson:"energy"`
	EnergyProduction int64 `bson:"energyProduction"`

	// Territory and assets - denormalized for performance
	// IMPORTANT: These arrays should be mutually exclusive:
	// - ColonizedSystems: Systems where player has embedded defending fleet
	// - ActiveStacks: Stacks in free movement (space, mining, traveling)
	// - MiningOperations: Asteroid/Nebula IDs where player has mining stacks
	ColonizedSystems []bson.ObjectID `bson:"colonizedSystems"` // System IDs where player has DefendingFleet
	ActiveStacks     []bson.ObjectID `bson:"activeStacks"`     // Stack IDs in free movement/mining (NOT in systems)
	MiningOperations []bson.ObjectID `bson:"miningOperations"` // Asteroid/Nebula IDs being mined by player's stacks

	// Game metadata
	IsAlive    bool      `bson:"isAlive"`
	JoinedAt   time.Time `bson:"joinedAt"`   // When player joined this game
	LastActive time.Time `bson:"lastActive"` // Last tick player was active
}
