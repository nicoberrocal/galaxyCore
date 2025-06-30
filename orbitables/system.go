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
	ID                  bson.ObjectID         `bson:"_id,omitempty"`
	SystemID            bson.ObjectID         `bson:"systemId,omitempty"`
	Name                string                `bson:"name"`
	NorthPole           b.EnergyBuilding      `bson:"northPole"`
	Left                b.EnergyBuilding      `bson:"left"`
	Right               b.EnergyBuilding      `bson:"right"`
	Back                b.EnergyBuilding      `bson:"back"`
	Front               b.EnergyBuilding      `bson:"front"`
	ShipYard            b.ShipYard            `bson:"shipyard"`
	ParticleAccelerator b.ParticleAccelerator `bson:"particleAccelerator"`
	FusionReactor       b.FusionReactor       `bson:"fusionReactor"`
	Metals              int64                 `bson:"metals"`
	Crystals            int64                 `bson:"crystals"`
	Hydrogen            int64                 `bson:"hydrogen"`
	Plasma              int64                 `bson:"plasma"`
}
