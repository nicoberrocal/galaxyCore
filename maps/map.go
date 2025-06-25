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
	PlayerID         bson.ObjectID `bson:"playerId,omitempty"`
	SetID            bson.ObjectID `bson:"shipSettings"`
	Energy           int64         `bson:"energy"`
	EnergyProduction int64         `bson:"energyProduction"`
}
type Set struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
