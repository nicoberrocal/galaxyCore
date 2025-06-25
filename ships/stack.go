package ships

import "go.mongodb.org/mongo-driver/v2/bson"

type ShipType string

const (
	Drone     ShipType = "drone"
	Scout     ShipType = "scout"
	Fighter   ShipType = "fighter"
	Bomber    ShipType = "bomber"
	Carrier   ShipType = "carrier"
	Destroyer ShipType = "destroyer"
)

type HPBucket struct {
	HP    int `bson:"hp"`    // full HP of that ship
	Count int `bson:"count"` // how many ships at this HP
}

type ShipStack struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	PlayerID  bson.ObjectID `bson:"playerId"`
	MapID     bson.ObjectID `bson:"mapId"`
	PositionX float64       `bson:"x"`
	PositionY float64       `bson:"y"`
	// One entry per ship type
	Ships     map[ShipType][]HPBucket `bson:"ships"`     // HP bucketed ships
	CreatedAt int64                   `bson:"createdAt"` // tick timestamp
}
