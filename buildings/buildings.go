package buildings

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Shared interface — useful for matching types, even if logic is reimplemented
type Building interface {
	GetName() string
}

// Core building structs with only fields and bson tags — NO METHODS
type BaseBuilding struct {
	Name            string    `bson:"name"`
	Level           int       `bson:"level"`
	ConstuctionTime time.Time `bson:"constuctionTime"`
	Queue           []Queue   `bson:"queue"`
	Upkeep          int       `bson:"upkeep"`
}

type MineBuilding struct {
	BaseBuilding
	Production int `bson:"production"`
}

type EnergyBuilding struct {
	BaseBuilding
	Production int `bson:"production"`
}

type Queue struct {
	Action   string        `bson:"action"`
	Start    bson.DateTime `bson:"start"`
	Duration int           `bson:"duration"`
}

// Specialized types (also no logic)
type ShipYard struct {
	BaseBuilding
}

type ParticleAccelerator struct {
	BaseBuilding
}

type FusionReactor struct {
	BaseBuilding
}

type SolarFarm struct {
	EnergyBuilding
}

type WindFarm struct {
	EnergyBuilding
}

type HydroElectricDam struct {
	EnergyBuilding
}

type Balloon struct {
	EnergyBuilding
}

type CrystalMine struct {
	MineBuilding
}

type MetalMine struct {
	MineBuilding
}
