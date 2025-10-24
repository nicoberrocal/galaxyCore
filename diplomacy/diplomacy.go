package diplomacy

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RelationDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	MapID     bson.ObjectID `bson:"mapId"`
	Player1   bson.ObjectID `bson:"player1"`
	Player2   bson.ObjectID `bson:"player2"`
	Relation  int32         `bson:"relation"`
	Until     time.Time     `bson:"until,omitempty"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}
