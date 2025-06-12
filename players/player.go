package players

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Player struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	Username   string        `bson:"username"`
	Email      string        `bson:"email"`
	Password   string        `bson:"password"`
	IsVerified bool          `bson:"verified"`
	IsPremium  bool          `bson:"is_premium"`
	CreatedAt  time.Time     `bson:"created_at"`
}
