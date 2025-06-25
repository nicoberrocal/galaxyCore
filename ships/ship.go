package ships

type Ship struct {
	ShipType         string
	AttackType       string
	LaserShield      int
	NuclearShield    int
	AntimatterShield int
	Speed            int
	VisibilityRange  int
	AttackRange      int
	HP               int
	AttackDamage     int
	AttackInterval   float64
	SpecialAbility   string
	AbilityCooldown  float64 // in hours
	AbilityDuration  float64 // in hours
	// Construction costs
	MetalCost   int
	CrystalCost int
	PlasmaCost  int
	// Carrier specific
	TransportCapacity int
	CanTransport      []string
}

var Ships = []Ship{
	{
		ShipType:          "Drone",
		AttackType:        "Laser",
		LaserShield:       5,
		NuclearShield:     15,
		AntimatterShield:  10,
		Speed:             80,
		VisibilityRange:   50,
		AttackRange:       30,
		HP:                50,
		AttackDamage:      5,
		AttackInterval:    1.0, // Uniform across all ships
		SpecialAbility:    "Mining Overdrive",
		AbilityCooldown:   12.0, // 12 hours
		AbilityDuration:   6.0,  // 6 hours of boosted resource gathering
		MetalCost:         100,
		CrystalCost:       50,
		PlasmaCost:        0,
		TransportCapacity: 0,
		CanTransport:      []string{},
	},
	{
		ShipType:          "Fighter",
		AttackType:        "Nuclear",
		LaserShield:       30,
		NuclearShield:     15,
		AntimatterShield:  25,
		Speed:             70,
		VisibilityRange:   60,
		AttackRange:       50,
		HP:                100,
		AttackDamage:      20,
		AttackInterval:    1.0,
		SpecialAbility:    "Stealth Cloak",
		AbilityCooldown:   8.0, // 8 hours
		AbilityDuration:   2.0, // 2 hours of invisibility
		MetalCost:         200,
		CrystalCost:       100,
		PlasmaCost:        0,
		TransportCapacity: 0,
		CanTransport:      []string{},
	},
	{
		ShipType:          "Scout",
		AttackType:        "Laser",
		LaserShield:       10,
		NuclearShield:     20,
		AntimatterShield:  15,
		Speed:             90,
		VisibilityRange:   100,
		AttackRange:       40,
		HP:                80,
		AttackDamage:      10,
		AttackInterval:    1.0,
		SpecialAbility:    "Deep Scan Pulse",
		AbilityCooldown:   6.0, // 6 hours
		AbilityDuration:   4.0, // 4 hours of extended vision range
		MetalCost:         150,
		CrystalCost:       75,
		PlasmaCost:        0,
		TransportCapacity: 0,
		CanTransport:      []string{},
	},
	{
		ShipType:          "Carrier",
		AttackType:        "Laser",
		LaserShield:       40,
		NuclearShield:     35,
		AntimatterShield:  45,
		Speed:             40,
		VisibilityRange:   80,
		AttackRange:       45,
		HP:                400,
		AttackDamage:      15,
		AttackInterval:    1.0,
		SpecialAbility:    "FTL Jump",
		AbilityCooldown:   24.0, // 24 hours
		AbilityDuration:   0.0,  // Instant ability
		MetalCost:         800,
		CrystalCost:       400,
		PlasmaCost:        200,
		TransportCapacity: 20,
		CanTransport:      []string{"Drone", "Fighter", "Scout"},
	},
	{
		ShipType:          "Bomber",
		AttackType:        "Antimatter",
		LaserShield:       35,
		NuclearShield:     40,
		AntimatterShield:  50,
		Speed:             55,
		VisibilityRange:   70,
		AttackRange:       90,
		HP:                280,
		AttackDamage:      75,
		AttackInterval:    1.0,
		SpecialAbility:    "Siege Mode",
		AbilityCooldown:   48.0, // 48 hours
		AbilityDuration:   12.0, // 12 hours of extreme range (interrupted by movement)
		MetalCost:         600,
		CrystalCost:       300,
		PlasmaCost:        150,
		TransportCapacity: 0,
		CanTransport:      []string{},
	},
	{
		ShipType:          "Destroyer",
		AttackType:        "Nuclear",
		LaserShield:       45,
		NuclearShield:     55,
		AntimatterShield:  40,
		Speed:             50,
		VisibilityRange:   75,
		AttackRange:       70,
		HP:                350,
		AttackDamage:      120,
		AttackInterval:    1.0,
		SpecialAbility:    "Hunter Protocol",
		AbilityCooldown:   72.0, // 72 hours
		AbilityDuration:   8.0,  // 8 hours of enhanced damage and speed
		MetalCost:         1000,
		CrystalCost:       500,
		PlasmaCost:        300,
		TransportCapacity: 0,
		CanTransport:      []string{},
	},
}
