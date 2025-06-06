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
}

var Ships = []Ship{
	{
		ShipType:         "Drone",
		AttackType:       "Laser",
		LaserShield:      5,
		NuclearShield:    15,
		AntimatterShield: 10,
		Speed:            80,
		VisibilityRange:  50,
		AttackRange:      30,
		HP:               50,
		AttackDamage:     5,
		AttackInterval:   1.0,
		SpecialAbility:   "Increased Resource Gathering",
	},
	{
		ShipType:         "Fighter",
		AttackType:       "Nuclear",
		LaserShield:      30,
		NuclearShield:    15,
		AntimatterShield: 25,
		Speed:            70,
		VisibilityRange:  60,
		AttackRange:      50,
		HP:               100,
		AttackDamage:     20,
		AttackInterval:   1.5,
		SpecialAbility:   "Invisibility",
	},
	{
		ShipType:         "Scout",
		AttackType:       "Laser",
		LaserShield:      10,
		NuclearShield:    20,
		AntimatterShield: 15,
		Speed:            90,
		VisibilityRange:  100,
		AttackRange:      40,
		HP:               80,
		AttackDamage:     10,
		AttackInterval:   1.2,
		SpecialAbility:   "Scan Map Areas",
	},
	{
		ShipType:         "Carrier",
		AttackType:       "Antimatter",
		LaserShield:      45,
		NuclearShield:    45,
		AntimatterShield: 30,
		Speed:            50,
		VisibilityRange:  70,
		AttackRange:      60,
		HP:               300,
		AttackDamage:     50,
		AttackInterval:   3.0,
		SpecialAbility:   "Faster-Than-Light Transport",
	},
	{
		ShipType:         "Destroyer",
		AttackType:       "Nuclear",
		LaserShield:      30,
		NuclearShield:    40,
		AntimatterShield: 35,
		Speed:            60,
		VisibilityRange:  60,
		AttackRange:      70,
		HP:               250,
		AttackDamage:     100,
		AttackInterval:   5.0,
		SpecialAbility:   "Burst Damage, Increased Attack Interval",
	},
}
