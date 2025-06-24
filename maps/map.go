package maps

import "go.mongodb.org/mongo-driver/v2/bson"

type MongoMap struct {
	ID         bson.ObjectID  `bson:"_id,omitempty"`
	ReadableId int64          `bson:"readableId,omitempty"`
	CreatorID  bson.ObjectID  `bson:"creatorId,omitempty"`
	Players    []PlayerConfig `bson:"players"`
	GameName   string         `bson:"gameName"`
	QPlayers   int8           `bson:"qPlayers"`
	PeaceDays  int8           `bson:"peaceDays"`
}
type PlayerConfig struct {
	PlayerID bson.ObjectID `bson:"playerId,omitempty"`
	SetID    bson.ObjectID `bson:"shipSettings"`
}
type Set struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
