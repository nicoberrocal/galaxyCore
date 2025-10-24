package diplomacy

import (
	"bytes"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Relation int

const (
	RelationUnknown Relation = iota
	RelationEnemy
	RelationAlly
	RelationCeasefire
)

type Pair struct {
	A bson.ObjectID
	B bson.ObjectID
}

func normalizePair(a, b bson.ObjectID) Pair {
	if bytes.Compare(a[:], b[:]) <= 0 {
		return Pair{A: a, B: b}
	}
	return Pair{A: b, B: a}
}

type Entry struct {
	Relation Relation
	// Until is zero for permanent relations
	Until time.Time
}

type State struct {
	MapID bson.ObjectID
	// Relations keyed by normalized player pair
	Relations map[Pair]Entry
}
