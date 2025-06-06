package players

import "go.mongodb.org/mongo-driver/v2/bson"

type Player struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
}
