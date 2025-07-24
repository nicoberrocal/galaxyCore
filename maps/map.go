package maps

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoMap struct {
	ID         bson.ObjectID  `bson:"_id,omitempty"`
	ReadableId int64          `bson:"readableId,omitempty"`
	CreatorID  bson.ObjectID  `bson:"creatorId,omitempty"`
	Players    []PlayerConfig `bson:"players"`
	GameName   string         `bson:"gameName"`
	QPlayers   int8           `bson:"qPlayers"`
	PeaceDays  int8           `bson:"peaceDays"`
	StartTime  time.Time      `bson:"startTime"`
	Started    bool           `bson:"started"`
	Finished   bool           `bson:"finished"`
	Ranked     bool           `bson:"ranked"`
}
type PlayerConfig struct {
	PlayerID bson.ObjectID `bson:"playerId,omitempty"`
	SetID    bson.ObjectID `bson:"shipSettings"`
}
type Set struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

type PlayerGameState struct {
	PlayerID         bson.ObjectID   `bson:"playerId"`         // References players collection
	MapID            bson.ObjectID   `bson:"mapId"`            // References maps collection
	ColonizedSystems []bson.ObjectID `bson:"colonizedSystems"` // References systems the player owns
	StackIDs         []bson.ObjectID `bson:"stackIds"`         // References all stacks owned by player
	Energy           int64           `bson:"energy"`
	EnergyProduction int64           `bson:"energyProduction"`
	LastUpdated      time.Time       `bson:"lastUpdate"` // Timestamp of last update
}

type ShipVisibility struct {
	StackID       bson.ObjectID   `bson:"stackId"`          // References the stack this visibility is for
	Systems       []bson.ObjectID `bson:"visibleSystems"`   // References systems that are visible to the player
	Stacks        []bson.ObjectID `bson:"visibleShips"`     // References ships that are visible to the player
	Nebulas       []bson.ObjectID `bson:"visibleNebulas"`   // References nebulas that are visible to the player
	Asteroids     []bson.ObjectID `bson:"visibleAsteroids"` // References asteroids that are visible to the player
	LastUpdated   time.Time       `bson:"lastUpdated"`      // Timestamp of last visibility update
	LastProcessed time.Time       `bson:"lastProcessed"`    // Last time this visibility was processed
	Version       int64           `bson:"version"`          // For optimistic locking
}
