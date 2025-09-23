package orbitables

import (
	"time"

	b "github.com/nicoberrocal/galaxyCore/buildings"
	"github.com/nicoberrocal/galaxyCore/ships"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Colonization struct {
	IsColonized       bool          `bson:"isColonized" json:"isColonized"`
	ColonizedBy       bson.ObjectID `bson:"colonizedBy,omitempty" json:"colonizedBy,omitempty"`             // FK to player who owns the system
	ColonizingFleetID bson.ObjectID `bson:"colonizingFleetId,omitempty" json:"colonizingFleetId,omitempty"` // FK to original stack that colonized
	ColonizedAt       time.Time     `bson:"colonizedAt,omitempty" json:"colonizedAt,omitempty"`             // When system was first colonized
}

// DefendingFleet represents the fleet currently controlling/defending a system
// This is embedded in the system when a stack colonizes it
// NOTE: If DefendingFleet exists, system MUST be colonized (Colonization.IsColonized = true)
type DefendingFleet struct {
	OriginalStackID bson.ObjectID                       `bson:"originalStackId" json:"originalStackId"` // Reference to stack that originally colonized
	PlayerID        bson.ObjectID                       `bson:"playerId" json:"playerId"`               // Current controlling player (may differ from ColonizedBy if allies took over)
	Ships           map[ships.ShipType][]ships.HPBucket `bson:"ships" json:"ships"`                     // Current fleet composition
	ArrivedAt       time.Time                           `bson:"arrivedAt" json:"arrivedAt"`
	Activity        string                              `bson:"activity" json:"activity"` // "defending", "colonizing", "building"

	// Allies (for symbolic merging in battle resolution)
	// These are fleets that merged with the main defending fleet
	AlliedFleets []AlliedFleet `bson:"alliedFleets,omitempty" json:"alliedFleets,omitempty"`
}

// AlliedFleet tracks fleets that merged symbolically for battle purposes
type AlliedFleet struct {
	OriginalStackID bson.ObjectID `bson:"originalStackId" json:"originalStackId"`
	PlayerID        bson.ObjectID `bson:"playerId" json:"playerId"`
	ArrivedAt       time.Time     `bson:"arrivedAt" json:"arrivedAt"`
}

// System represents a colonizable star system
// IMPORTANT: System state consistency rules:
// 1. If DefendingFleet exists → Colonization.IsColonized MUST be true
// 2. If Colonization.IsColonized is true → DefendingFleet MUST exist
// 3. If DefendingFleet is nil → System is unclaimed/orphaned
// System is more militaristic than planets, it has a colonization state and a defending fleet
type System struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	X             float64       `bson:"x" json:"x"`
	Y             float64       `bson:"y" json:"y"`
	MapID         bson.ObjectID `bson:"mapId,omitempty" json:"mapId,omitempty"`
	Name          string        `bson:"name" json:"name"`
	Texture       string        `bson:"texture" json:"texture"`
	Constellation string        `bson:"constellation,omitempty" json:"constellation,omitempty"`
	Planet        *Planet       `bson:"planet,omitempty" json:"planet,omitempty"`
	// System control and defense - uses hybrid approach
	Colonization   *Colonization   `bson:"colonization,omitempty" json:"colonization,omitempty"`
	DefendingFleet *DefendingFleet `bson:"defendingFleet,omitempty" json:"defendingFleet,omitempty"` // Embedded fleet when colonized

	// Collision detection for system entry
	CollisionRadius float64 `bson:"collisionRadius" json:"collisionRadius"` // Radius for determining system entry
	Version         int64   `bson:"version" json:"version"`                 // For optimistic locking
}

// Planet represents a colonizable planet
// IMPORTANT: Planet state consistency rules:
// 1. If Planet.SystemID exists → Planet.SystemID MUST be a valid system ID
// 2. If Planet.SystemID is nil → Planet is unclaimed/orphaned
// 3. If Planet.SystemID exists → Planet.SystemID MUST be a valid system ID
// 4. If Planet.SystemID is nil → Planet is unclaimed/orphaned
// Planet is more economic than systems, it has energy and material buildings, logistic hubs, and a shipyard
type Planet struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SystemID            bson.ObjectID `bson:"systemId,omitempty" json:"systemId,omitempty"`
	MapID               bson.ObjectID `bson:"mapId,omitempty" json:"mapId,omitempty"`
	Name                string        `bson:"name" json:"name"`
	NorthPole           *bson.M       `bson:"northPole,omitempty" json:"northPole,omitempty"`
	Left                *bson.M       `bson:"left,omitempty" json:"left,omitempty"`
	Right               *bson.M       `bson:"right,omitempty" json:"right,omitempty"`
	Back                *bson.M       `bson:"back,omitempty" json:"back,omitempty"`
	Front               *bson.M       `bson:"front,omitempty" json:"front,omitempty"`
	ShipYard            *bson.M       `bson:"shipyard,omitempty" json:"shipyard,omitempty"`
	ParticleAccelerator *bson.M       `bson:"particleAccelerator,omitempty" json:"particleAccelerator,omitempty"`
	FusionReactor       *bson.M       `bson:"fusionReactor,omitempty" json:"fusionReactor,omitempty"`
	Metals              int64         `bson:"metals" json:"metals"`
	Crystals            int64         `bson:"crystals" json:"crystals"`
	Hydrogen            int64         `bson:"hydrogen" json:"hydrogen"`
	Plasma              int64         `bson:"plasma" json:"plasma"`
	Version             int64         `bson:"version" json:"version"` // For optimistic locking
}

// Helper methods to work with buildings as interfaces
func (p *Planet) GetNorthPole() (b.Building, error) {
	if p.NorthPole == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.NorthPole)
}

func (p *Planet) SetNorthPole(building b.Building) {
	result := b.BuildingToBSON(building)
	p.NorthPole = &result
}

func (p *Planet) GetLeft() (b.Building, error) {
	if p.Left == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.Left)
}

func (p *Planet) SetLeft(building b.Building) {
	result := b.BuildingToBSON(building)
	p.Left = &result
}

func (p *Planet) GetRight() (b.Building, error) {
	if p.Right == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.Right)
}

func (p *Planet) SetRight(building b.Building) {
	result := b.BuildingToBSON(building)
	p.Right = &result
}

func (p *Planet) GetBack() (b.Building, error) {
	if p.Back == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.Back)
}

func (p *Planet) SetBack(building b.Building) {
	result := b.BuildingToBSON(building)
	p.Back = &result
}

func (p *Planet) GetFront() (b.Building, error) {
	if p.Front == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.Front)
}

func (p *Planet) SetFront(building b.Building) {
	result := b.BuildingToBSON(building)
	p.Front = &result
}

func (p *Planet) GetShipYard() (b.Building, error) {
	if p.ShipYard == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.ShipYard)
}

func (p *Planet) SetShipYard(building b.Building) {
	result := b.BuildingToBSON(building)
	p.ShipYard = &result
}

func (p *Planet) GetParticleAccelerator() (b.Building, error) {
	if p.ParticleAccelerator == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.ParticleAccelerator)
}

func (p *Planet) SetParticleAccelerator(building b.Building) {
	result := b.BuildingToBSON(building)
	p.ParticleAccelerator = &result
}

func (p *Planet) GetFusionReactor() (b.Building, error) {
	if p.FusionReactor == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(*p.FusionReactor)
}

func (p *Planet) SetFusionReactor(building b.Building) {
	result := b.BuildingToBSON(building)
	p.FusionReactor = &result
}
