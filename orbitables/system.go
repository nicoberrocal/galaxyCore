package orbitables

import (
	"time"

	b "github.com/nicoberrocal/galaxyCore/buildings"
	"github.com/nicoberrocal/galaxyCore/ships"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Colonization struct {
	IsColonized       bool          `bson:"isColonized"`
	ColonizedBy       bson.ObjectID `bson:"colonizedBy,omitempty"`       // FK to player who owns the system
	ColonizingFleetID bson.ObjectID `bson:"colonizingFleetId,omitempty"` // FK to original stack that colonized
	ColonizedAt       time.Time     `bson:"colonizedAt,omitempty"`       // When system was first colonized
}

// DefendingFleet represents the fleet currently controlling/defending a system
// This is embedded in the system when a stack colonizes it
// NOTE: If DefendingFleet exists, system MUST be colonized (Colonization.IsColonized = true)
type DefendingFleet struct {
	OriginalStackID bson.ObjectID                       `bson:"originalStackId"` // Reference to stack that originally colonized
	PlayerID        bson.ObjectID                       `bson:"playerId"`        // Current controlling player (may differ from ColonizedBy if allies took over)
	Ships           map[ships.ShipType][]ships.HPBucket `bson:"ships"`           // Current fleet composition
	ArrivedAt       time.Time                           `bson:"arrivedAt"`
	Activity        string                              `bson:"activity"` // "defending", "colonizing", "building"

	// Allies (for symbolic merging in battle resolution)
	// These are fleets that merged with the main defending fleet
	AlliedFleets []AlliedFleet `bson:"alliedFleets,omitempty"`
}

// AlliedFleet tracks fleets that merged symbolically for battle purposes
type AlliedFleet struct {
	OriginalStackID bson.ObjectID `bson:"originalStackId"`
	PlayerID        bson.ObjectID `bson:"playerId"`
	ArrivedAt       time.Time     `bson:"arrivedAt"`
}

// System represents a colonizable star system
// IMPORTANT: System state consistency rules:
// 1. If DefendingFleet exists → Colonization.IsColonized MUST be true
// 2. If Colonization.IsColonized is true → DefendingFleet MUST exist
// 3. If DefendingFleet is nil → System is unclaimed/orphaned
// System is more militaristic than planets, it has a colonization state and a defending fleet
type System struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	X       float64       `bson:"x"`
	Y       float64       `bson:"y"`
	MapID   bson.ObjectID `bson:"mapId,omitempty"`
	Name    string        `bson:"name"`
	Texture string        `bson:"texture"`
	Planet  Planet        `bson:"planet,omitempty"`
	// System control and defense - uses hybrid approach
	Colonization   Colonization    `bson:"colonization"`
	DefendingFleet *DefendingFleet `bson:"defendingFleet,omitempty"` // Embedded fleet when colonized

	// Collision detection for system entry
	CollisionRadius float64 `bson:"collisionRadius"` // Radius for determining system entry
}

// Planet represents a colonizable planet
// IMPORTANT: Planet state consistency rules:
// 1. If Planet.SystemID exists → Planet.SystemID MUST be a valid system ID
// 2. If Planet.SystemID is nil → Planet is unclaimed/orphaned
// 3. If Planet.SystemID exists → Planet.SystemID MUST be a valid system ID
// 4. If Planet.SystemID is nil → Planet is unclaimed/orphaned
// Planet is more economic than systems, it has energy and material buildings, logistic hubs, and a shipyard
type Planet struct {
	ID                  bson.ObjectID `bson:"_id,omitempty"`
	SystemID            bson.ObjectID `bson:"systemId,omitempty"`
	MapID               bson.ObjectID `bson:"mapId,omitempty"`
	Name                string        `bson:"name"`
	NorthPole           bson.M        `bson:"northPole"`
	Left                bson.M        `bson:"left"`
	Right               bson.M        `bson:"right"`
	Back                bson.M        `bson:"back"`
	Front               bson.M        `bson:"front"`
	ShipYard            bson.M        `bson:"shipyard"`
	ParticleAccelerator bson.M        `bson:"particleAccelerator"`
	FusionReactor       bson.M        `bson:"fusionReactor"`
	Metals              int64         `bson:"metals"`
	Crystals            int64         `bson:"crystals"`
	Hydrogen            int64         `bson:"hydrogen"`
	Plasma              int64         `bson:"plasma"`
}

// Helper methods to work with buildings as interfaces
func (p *Planet) GetNorthPole() (b.Building, error) {
	if p.NorthPole == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.NorthPole)
}

func (p *Planet) SetNorthPole(building b.Building) {
	p.NorthPole = b.BuildingToBSON(building)
}

func (p *Planet) GetLeft() (b.Building, error) {
	if p.Left == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.Left)
}

func (p *Planet) SetLeft(building b.Building) {
	p.Left = b.BuildingToBSON(building)
}

func (p *Planet) GetRight() (b.Building, error) {
	if p.Right == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.Right)
}

func (p *Planet) SetRight(building b.Building) {
	p.Right = b.BuildingToBSON(building)
}

func (p *Planet) GetBack() (b.Building, error) {
	if p.Back == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.Back)
}

func (p *Planet) SetBack(building b.Building) {
	p.Back = b.BuildingToBSON(building)
}

func (p *Planet) GetFront() (b.Building, error) {
	if p.Front == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.Front)
}

func (p *Planet) SetFront(building b.Building) {
	p.Front = b.BuildingToBSON(building)
}

func (p *Planet) GetShipYard() (b.Building, error) {
	if p.ShipYard == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.ShipYard)
}

func (p *Planet) SetShipYard(building b.Building) {
	p.ShipYard = b.BuildingToBSON(building)
}

func (p *Planet) GetParticleAccelerator() (b.Building, error) {
	if p.ParticleAccelerator == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.ParticleAccelerator)
}

func (p *Planet) SetParticleAccelerator(building b.Building) {
	p.ParticleAccelerator = b.BuildingToBSON(building)
}

func (p *Planet) GetFusionReactor() (b.Building, error) {
	if p.FusionReactor == nil {
		return nil, nil
	}
	return b.CreateBuildingFromMongoDB(p.FusionReactor)
}

func (p *Planet) SetFusionReactor(building b.Building) {
	p.FusionReactor = b.BuildingToBSON(building)
}
