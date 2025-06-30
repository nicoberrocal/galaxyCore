package buildings

var BaseEnergyOutput = map[string]int{
	"Hydro":   100,
	"Solar":   80,
	"Wind":    60,
	"Balloon": 50,
}
var PlanetSuitability = map[string]map[string]float64{
	"Mercury": {
		"Hydro":   0,
		"Solar":   1,
		"Wind":    0.1,
		"Balloon": 0.8,
	},
	"Venus": {
		"Hydro":   0,
		"Solar":   0.4,
		"Wind":    0.9,
		"Balloon": 1.0,
	},
	"Earth": {
		"Hydro":   1,
		"Solar":   0.8,
		"Wind":    0.8,
		"Balloon": 0.6,
	},
	"Mars": {
		"Hydro":   0,
		"Solar":   0.6,
		"Wind":    0.3,
		"Balloon": 0.4,
	},
}
var GrowthRate = map[string]map[string]float64{
	"Mercury": {
		"Hydro":   0,
		"Solar":   0.15,
		"Wind":    0.05,
		"Balloon": 0.1,
	},
	"Venus": {
		"Hydro":   0.01,
		"Solar":   0.05,
		"Wind":    0.12,
		"Balloon": 0.15,
	},
	"Earth": {
		"Hydro":   0.1,
		"Solar":   0.08,
		"Wind":    0.12,
		"Balloon": 0.06,
	},
	"Mars": {
		"Hydro":   0.01,
		"Solar":   0.10,
		"Wind":    0.05,
		"Balloon": 0.08,
	},
}

// BaseExtractionRate defines the base extraction rate for each mine type (units per time).
var BaseExtractionRate = map[string]int{
	"MetalMine":   100, // Base rate for metal extraction
	"CrystalMine": 80,  // Base rate for crystal extraction
}

// ResourceSuitability defines the suitability multiplier for each resource on each planet.
var ResourceSuitability = map[string]map[string]float64{
	"Mercury": {
		"Metals":   0.8, // Moderate metal availability
		"Crystals": 1.2, // High crystal availability due to solar proximity
	},
	"Venus": {
		"Metals":   0.6, // Low metal availability
		"Crystals": 0.7, // Moderate crystal availability
	},
	"Earth": {
		"Metals":   1.0, // Balanced metal availability
		"Crystals": 0.9, // Slightly below average crystal availability
	},
	"Mars": {
		"Metals":   1.2, // High metal availability (e.g., iron oxide)
		"Crystals": 0.5, // Low crystal availability
	},
}

// ExtractionGrowthRate defines the growth rate per level for each mine type on each planet.
var ExtractionGrowthRate = map[string]map[string]float64{
	"Mercury": {
		"MetalMine":   0.10, // 10% growth per level for metal mines
		"CrystalMine": 0.15, // 15% growth per level for crystal mines
	},
	"Venus": {
		"MetalMine":   0.05, // 5% growth per level for metal mines
		"CrystalMine": 0.08, // 8% growth per level for crystal mines
	},
	"Earth": {
		"MetalMine":   0.12, // 12% growth per level for metal mines
		"CrystalMine": 0.10, // 10% growth per level for crystal mines
	},
	"Mars": {
		"MetalMine":   0.15, // 15% growth per level for metal mines
		"CrystalMine": 0.05, // 5% growth per level for crystal mines
	},
}
