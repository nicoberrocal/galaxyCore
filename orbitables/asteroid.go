package orbitables

import "go.mongodb.org/mongo-driver/v2/bson"

type Asteroid struct {
	ID      bson.ObjectID `bson:"id,omitempty"`
	MapID   bson.ObjectID `bson:"mapId,omitempty"`
	X       float64       `bson:"x"`
	Y       float64       `bson:"y"`
	Name    string        `bson:"name"`
	Texture string        `bson:"texture"`
}
