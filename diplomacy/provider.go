package diplomacy

import (
	"time"

	"github.com/nicoberrocal/galaxyCore/ships"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Provider interface {
	AreAllies(mapID, a, b bson.ObjectID, now time.Time) (bool, error)
	AreEnemies(mapID, a, b bson.ObjectID, now time.Time) (bool, error)
}

func AreStacksEnemies(p Provider, s1, s2 *ships.ShipStack, now time.Time) (bool, error) {
	return p.AreEnemies(s1.MapID, s1.PlayerID, s2.PlayerID, now)
}

func AreStacksAllies(p Provider, s1, s2 *ships.ShipStack, now time.Time) (bool, error) {
	return p.AreAllies(s1.MapID, s1.PlayerID, s2.PlayerID, now)
}
