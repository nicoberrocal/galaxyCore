package orbitables

import "go.mongodb.org/mongo-driver/v2/bson"

type Asteroid struct {
	ID                 bson.ObjectID      `bson:"_id,omitempty"`
	MapID              bson.ObjectID      `bson:"mapId,omitempty"`
	X                  float64            `bson:"x"`
	Y                  float64            `bson:"y"`
	Name               string             `bson:"name"`
	Texture            string             `bson:"texture"`
	ResourceExtraction ResourceExtraction `bson:"resourceExtraction"`
	CollisionRadius    float64            `bson:"collisionRadius"`
	TotalMetal         int64              `bson:"totalMetal"`
	TotalCrystal       int64              `bson:"totalCrystal"`
}
