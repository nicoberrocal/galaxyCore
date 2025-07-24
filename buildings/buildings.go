package buildings

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Shared interface — useful for matching types, even if logic is reimplemented
type Building interface {
	GetName() string
	GetLevel() int
	GetType() string
	GetProduction() int
	GetUpkeep() int
	GetConstructionTime() time.Time
	GetQueue() []Queue
}

// Core building structs with methods to implement the interface
type BaseBuilding struct {
	Name            string    `bson:"name"`
	Level           int       `bson:"level"`
	ConstuctionTime time.Time `bson:"constuctionTime,omitempty" json:"constuctionTime,omitempty"`
	Queue           []Queue   `bson:"queue,omitempty" json:"queue,omitempty"`
	Upkeep          int       `bson:"upkeep,omitempty" json:"upkeep,omitempty"`
	LastUpdated     time.Time `bson:"lastUpdated,omitempty" json:"lastUpdated,omitempty"`     // Last time this building was updated
	LastProcessed   time.Time `bson:"lastProcessed,omitempty" json:"lastProcessed,omitempty"` // Last time this building was processed
}

// Implement Building interface for BaseBuilding
func (b BaseBuilding) GetName() string {
	return b.Name
}

func (b BaseBuilding) GetLevel() int {
	return b.Level
}

func (b BaseBuilding) GetType() string {
	return "base"
}

func (b BaseBuilding) GetProduction() int {
	return 0 // Base buildings don't produce
}

func (b BaseBuilding) GetUpkeep() int {
	return b.Upkeep
}

func (b BaseBuilding) GetConstructionTime() time.Time {
	return b.ConstuctionTime
}

func (b BaseBuilding) GetQueue() []Queue {
	return b.Queue
}

type MineBuilding struct {
	BaseBuilding
	Production int `bson:"production"`
}

// Implement Building interface for MineBuilding
func (m MineBuilding) GetType() string {
	return "mine"
}

func (m MineBuilding) GetProduction() int {
	return m.Production
}

type EnergyBuilding struct {
	BaseBuilding
	Production int `bson:"production"`
}

// Implement Building interface for EnergyBuilding
func (e EnergyBuilding) GetType() string {
	return "energy"
}

func (e EnergyBuilding) GetProduction() int {
	return e.Production
}

type Queue struct {
	Action   string        `bson:"action"`
	Start    bson.DateTime `bson:"start"`
	Duration int           `bson:"duration"`
}

// Specialized types with interface implementation
type ShipYard struct {
	BaseBuilding
}

func (s ShipYard) GetType() string {
	return "shipyard"
}

type ParticleAccelerator struct {
	BaseBuilding
}

func (p ParticleAccelerator) GetType() string {
	return "particle_accelerator"
}

type FusionReactor struct {
	BaseBuilding
}

func (f FusionReactor) GetType() string {
	return "fusion_reactor"
}

type SolarFarm struct {
	EnergyBuilding
}

func (s SolarFarm) GetType() string {
	return "solar_farm"
}

type WindFarm struct {
	EnergyBuilding
}

func (w WindFarm) GetType() string {
	return "wind_farm"
}

type HydroElectricDam struct {
	EnergyBuilding
}

func (h HydroElectricDam) GetType() string {
	return "hydro_electric_dam"
}

type Balloon struct {
	EnergyBuilding
}

func (b Balloon) GetType() string {
	return "balloon"
}

type CrystalMine struct {
	MineBuilding
}

func (c CrystalMine) GetType() string {
	return "crystal_mine"
}

type MetalMine struct {
	MineBuilding
}

func (m MetalMine) GetType() string {
	return "metal_mine"
}

// MongoDB Helper Functions for Interface Polymorphism

// Helper function to create buildings from MongoDB data
func CreateBuildingFromMongoDB(data bson.M) (Building, error) {
	buildingType, ok := data["type"].(string)
	if !ok {
		return nil, errors.New("missing or invalid building type")
	}

	switch buildingType {
	case "base":
		var building BaseBuilding
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "solar_farm":
		var building SolarFarm
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "wind_farm":
		var building WindFarm
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "crystal_mine":
		var building CrystalMine
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "metal_mine":
		var building MetalMine
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "hydro_electric_dam":
		var building HydroElectricDam
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "balloon":
		var building Balloon
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "shipyard":
		var building ShipYard
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "particle_accelerator":
		var building ParticleAccelerator
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	case "fusion_reactor":
		var building FusionReactor
		dataBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := bson.Unmarshal(dataBytes, &building); err != nil {
			return nil, err
		}
		return building, nil
	default:
		return nil, fmt.Errorf("unknown building type: %s", buildingType)
	}
}

// Helper function to convert Building interface to BSON document
func BuildingToBSON(building Building) bson.M {
	if building == nil {
		return nil
	}

	return bson.M{
		"type":             building.GetType(),
		"name":             building.GetName(),
		"level":            building.GetLevel(),
		"upkeep":           building.GetUpkeep(),
		"production":       building.GetProduction(),
		"constructionTime": building.GetConstructionTime(),
		"queue":            building.GetQueue(),
	}
}
