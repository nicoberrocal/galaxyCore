package diplomacy

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Relation string

const (
	RelationUnknown   Relation = "unknown"
	RelationEnemy     Relation = "enemy"
	RelationAlly      Relation = "ally"
	RelationCeasefire Relation = "ceasefire"
)

type RelationDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	MapID     bson.ObjectID `bson:"mapId"`
	Player1   bson.ObjectID `bson:"player1"`
	Player2   bson.ObjectID `bson:"player2"`
	Relation  Relation      `bson:"relation"`
	Until     time.Time     `bson:"until,omitempty"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}
