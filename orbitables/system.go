package orbitables

import (
	b "github.com/nicoberocal/galaxyCore/buildings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Define the structs
type System struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	X       float64       `bson:"x"`
	Y       float64       `bson:"y"`
	MapID   bson.ObjectID `bson:"mapId,omitempty"`
	Name    string        `bson:"name"`
	Texture string        `bson:"texture"`
}

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
}
