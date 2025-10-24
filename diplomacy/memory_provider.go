package diplomacy

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PeaceDaysLookup allows the provider to respect map start and peace window
// without depending on the maps package directly.
type PeaceDaysLookup func(mapID bson.ObjectID) (start time.Time, peaceDays int, ok bool)

type MemoryProvider struct {
	states      map[bson.ObjectID]*State // keyed by MapID
	peaceLookup PeaceDaysLookup          // optional
}

func NewMemoryProvider(peaceLookup PeaceDaysLookup) *MemoryProvider {
	return &MemoryProvider{states: make(map[bson.ObjectID]*State), peaceLookup: peaceLookup}
}

// EnsureState makes sure a State exists for mapID and returns it.
func (p *MemoryProvider) EnsureState(mapID bson.ObjectID) *State {
	st, ok := p.states[mapID]
	if !ok {
		st = &State{MapID: mapID, Relations: make(map[Pair]Entry)}
		p.states[mapID] = st
	}
	return st
}

func (p *MemoryProvider) AreAllies(mapID, a, b bson.ObjectID, now time.Time) (bool, error) {
	pair := normalizePair(a, b)
	st := p.EnsureState(mapID)
	if e, ok := st.Relations[pair]; ok {
		if e.Relation == RelationAlly {
			if e.Until.IsZero() || now.Before(e.Until) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (p *MemoryProvider) AreEnemies(mapID, a, b bson.ObjectID, now time.Time) (bool, error) {
	// Respect PeaceDays if available
	if p.peaceLookup != nil {
		if start, days, ok := p.peaceLookup(mapID); ok && days > 0 {
			if now.Before(start.Add(time.Duration(days) * 24 * time.Hour)) {
				// During peace period, not enemies unless explicitly marked as enemy
				pair := normalizePair(a, b)
				st := p.EnsureState(mapID)
				if e, ok := st.Relations[pair]; ok {
					if e.Relation == RelationEnemy {
						if e.Until.IsZero() || now.Before(e.Until) {
							return true, nil
						}
					}
				}
				return false, nil
			}
		}
	}

	// Outside peace, default to enemies unless ally or valid ceasefire.
	pair := normalizePair(a, b)
	st := p.EnsureState(mapID)
	if e, ok := st.Relations[pair]; ok {
		switch e.Relation {
		case RelationAlly:
			if e.Until.IsZero() || now.Before(e.Until) {
				return false, nil
			}
		case RelationCeasefire:
			if e.Until.After(now) {
				return false, nil
			}
		case RelationEnemy:
			if e.Until.IsZero() || now.Before(e.Until) {
				return true, nil
			}
		}
	}
	return true, nil
}

// Manager methods for convenience in tests/tools.
func (p *MemoryProvider) FormAlliance(mapID, a, b bson.ObjectID, until *time.Time) {
	st := p.EnsureState(mapID)
	entry := Entry{Relation: RelationAlly}
	if until != nil {
		entry.Until = *until
	}
	st.Relations[normalizePair(a, b)] = entry
}

func (p *MemoryProvider) BreakAlliance(mapID, a, b bson.ObjectID) {
	st := p.EnsureState(mapID)
	delete(st.Relations, normalizePair(a, b))
}

func (p *MemoryProvider) SetCeasefire(mapID, a, b bson.ObjectID, until time.Time) {
	st := p.EnsureState(mapID)
	st.Relations[normalizePair(a, b)] = Entry{Relation: RelationCeasefire, Until: until}
}

func (p *MemoryProvider) SetEnemy(mapID, a, b bson.ObjectID, until *time.Time) {
	st := p.EnsureState(mapID)
	entry := Entry{Relation: RelationEnemy}
	if until != nil {
		entry.Until = *until
	}
	st.Relations[normalizePair(a, b)] = entry
}
