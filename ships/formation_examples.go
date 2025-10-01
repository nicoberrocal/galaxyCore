package ships

import (
	"fmt"
	"time"
)

// Example functions demonstrating formation system usage
// These are for documentation and testing purposes

// ExampleBasicFormationSetup demonstrates setting up a formation for a stack.
func ExampleBasicFormationSetup() {
	// Create a ship stack
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 10}},
			Bomber:    {{HP: 500, Count: 5}},
			Scout:     {{HP: 100, Count: 8}},
			Destroyer: {{HP: 600, Count: 3}},
		},
		Role: RoleTactical,
	}
	
	// Set formation
	now := time.Now()
	eta := stack.SetFormation(FormationVanguard, now)
	
	fmt.Printf("Formation set to: %s\n", stack.Formation.Type)
	fmt.Printf("Will be active at: %s\n", eta)
	fmt.Printf("Reconfiguration time: %s\n", eta.Sub(now))
	
	// Print formation info
	fmt.Println(GetFormationInfo(stack.Formation))
}

// ExampleFormationCombat demonstrates combat with formations.
func ExampleFormationCombat() {
	// Attacker with Vanguard formation
	attacker := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 15}},
			Destroyer: {{HP: 600, Count: 5}},
		},
		Role: RoleTactical,
	}
	attacker.SetFormation(FormationVanguard, time.Now())
	
	// Defender with Box formation
	defender := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 12}},
			Carrier: {{HP: 900, Count: 2}},
		},
		Role: RoleTactical,
	}
	defender.SetFormation(FormationBox, time.Now())
	
	// Execute battle round
	result := ExecuteFormationBattleRound(attacker, defender)
	
	fmt.Printf("Formation Counter: %.2fx (Vanguard vs Box)\n", result.FormationAdvantage)
	fmt.Printf("Attacker dealt: %d damage\n", result.AttackerDamageDealt)
	fmt.Printf("Defender dealt: %d damage\n", result.DefenderDamageDealt)
	fmt.Printf("Attacker ships lost: %v\n", result.AttackerShipsLost)
	fmt.Printf("Defender ships lost: %v\n", result.DefenderShipsLost)
}

// ExampleFormationCounters demonstrates formation rock-paper-scissors.
func ExampleFormationCounters() {
	// Check what counters Box formation
	counters := GetCounterFormations(FormationBox)
	fmt.Printf("Formations that counter Box: %v\n", counters)
	
	// Check what Box is weak against
	weakAgainst := GetCounteredByFormations(FormationBox)
	fmt.Printf("Box is weak against: %v\n", weakAgainst)
	
	// Compare two formations
	analysis := CompareFormations(FormationVanguard, FormationBox)
	fmt.Println(analysis)
}

// ExampleGemPositionSynergy demonstrates gem-position bonuses.
func ExampleGemPositionSynergy() {
	// Create a stack with gems
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 10}},
			Bomber:  {{HP: 500, Count: 5}},
		},
		Role: RoleTactical,
		Loadouts: map[ShipType]ShipLoadout{
			Fighter: {
				Sockets: []Gem{
					GemCatalog["laser-3"],  // Laser gem tier 3
					GemCatalog["kinetic-2"], // Kinetic gem tier 2
				},
			},
			Bomber: {
				Sockets: []Gem{
					GemCatalog["nuclear-3"], // Nuclear gem tier 3
					GemCatalog["sensor-2"],  // Sensor gem tier 2
				},
			},
		},
	}
	
	// Set formation (Fighters go to Front, Bombers to Back)
	stack.SetFormation(FormationLine, time.Now())
	
	// Get effective stats for Fighter in front position
	fighterStats, _ := stack.EffectiveShipInFormation(Fighter, 0)
	fmt.Printf("Fighter in Front position:\n")
	fmt.Printf("  Attack Damage: %d\n", fighterStats.AttackDamage)
	fmt.Printf("  Laser Shield: %d\n", fighterStats.LaserShield)
	fmt.Printf("  HP: %d\n", fighterStats.HP)
	
	// Get effective stats for Bomber in back position
	bomberStats, _ := stack.EffectiveShipInFormation(Bomber, 0)
	fmt.Printf("Bomber in Back position:\n")
	fmt.Printf("  Attack Damage: %d\n", bomberStats.AttackDamage)
	fmt.Printf("  Attack Range: %d\n", bomberStats.AttackRange)
	fmt.Printf("  Visibility: %d\n", bomberStats.VisibilityRange)
}

// ExampleCompositionBonuses demonstrates fleet composition bonuses.
func ExampleCompositionBonuses() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Scout:     {{HP: 100, Count: 1}},
			Fighter:   {{HP: 200, Count: 2}},
			Bomber:    {{HP: 500, Count: 1}},
			Destroyer: {{HP: 600, Count: 1}},
		},
		Role: RoleTactical,
	}
	
	// Evaluate composition bonuses
	bonusMods, activeBonuses := EvaluateCompositionBonuses(stack.Ships)
	
	fmt.Printf("Active Composition Bonuses:\n")
	for _, bonus := range activeBonuses {
		fmt.Printf("  - %s: %s\n", bonus.Type, bonus.Description)
	}
	
	fmt.Printf("\nTotal Bonus Stats:\n")
	fmt.Printf("  Speed: +%d\n", bonusMods.SpeedDelta)
	fmt.Printf("  Damage bonus: +%.1f%%\n", bonusMods.Damage.LaserPct*100)
}

// ExampleAbilityFormationSynergy demonstrates ability bonuses from position.
func ExampleAbilityFormationSynergy() {
	// Check optimal position for Focus Fire
	position := GetOptimalPositionForAbility(AbilityFocusFire)
	fmt.Printf("Optimal position for Focus Fire: %s\n", position)
	
	// Get modifications when used from Front
	mods := GetAbilityFormationMods(AbilityFocusFire, PositionFront)
	fmt.Printf("Focus Fire from Front position:\n")
	for stat, value := range mods {
		fmt.Printf("  %s: %.2f\n", stat, value)
	}
	
	// List all abilities that benefit from Back position
	abilities := GetAbilitiesForPosition(PositionBack)
	fmt.Printf("\nAbilities that benefit from Back position: %v\n", abilities)
}

// ExampleAutoFormationSelection demonstrates automatic formation selection.
func ExampleAutoFormationSelection() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Scout:   {{HP: 100, Count: 3}},
			Fighter: {{HP: 200, Count: 5}},
		},
		Role: RoleTactical,
	}
	
	// Enemy is using Phalanx
	enemyFormation := FormationPhalanx
	
	// Find best counter
	template := FindBestTemplate(stack.Ships, stack.Role, enemyFormation)
	if template != nil {
		fmt.Printf("Best formation vs %s: %s\n", enemyFormation, template.Formation)
		fmt.Printf("Description: %s\n", template.Description)
		
		// Apply the template
		stack.SetFormation(template.Formation, time.Now())
	}
	
	// Get recommendations
	recommendations := GetFormationRecommendations(stack.Ships, stack.Role)
	fmt.Printf("Recommended formations for this fleet: %v\n", recommendations)
}

// ExampleFormationAnalysis demonstrates analyzing formation effectiveness.
func ExampleFormationAnalysis() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 10}},
			Bomber:    {{HP: 500, Count: 5}},
			Scout:     {{HP: 100, Count: 8}},
			Destroyer: {{HP: 600, Count: 3}},
		},
		Role: RoleTactical,
	}
	
	stack.SetFormation(FormationVanguard, time.Now())
	
	// Analyze position effectiveness
	effectiveness := AnalyzePositionEffectiveness(stack)
	fmt.Printf("Position Effectiveness:\n")
	for position, score := range effectiveness {
		fmt.Printf("  %s: %.2f\n", position, score)
	}
	
	// Get suggestions for improvement
	enemyFormation := FormationBox
	suggestions := SuggestFormationChanges(stack, enemyFormation)
	fmt.Printf("\nSuggestions:\n")
	for _, suggestion := range suggestions {
		fmt.Printf("  - %s\n", suggestion)
	}
	
	// Validate formation
	errors := ValidateFormation(stack.Formation)
	if len(errors) > 0 {
		fmt.Printf("\nFormation errors: %v\n", errors)
	} else {
		fmt.Println("\nFormation is valid!")
	}
}

// ExampleRoleModeFormationSynergy demonstrates role-formation interactions.
func ExampleRoleModeFormationSynergy() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 10}},
			Bomber:  {{HP: 500, Count: 5}},
		},
		Role: RoleTactical,
	}
	
	// Tactical mode gets -30% reconfiguration time
	stack.SetFormation(FormationVanguard, time.Now())
	fmt.Printf("Tactical mode reconfiguration: %d seconds\n", 
		RoleModeFormationBonus(RoleTactical, FormationCatalog[FormationVanguard].ReconfigureTime))
	
	// Switch to Economic mode
	stack.Role = RoleEconomic
	stack.SetFormation(FormationBox, time.Now())
	fmt.Printf("Economic mode reconfiguration: %d seconds\n",
		RoleModeFormationBonus(RoleEconomic, FormationCatalog[FormationBox].ReconfigureTime))
	
	// Get effective stats with role+formation bonuses
	fighterStats, _ := stack.EffectiveShipInFormation(Fighter, 0)
	fmt.Printf("\nFighter stats with Economic+Box:\n")
	fmt.Printf("  Attack Damage: %d\n", fighterStats.AttackDamage)
	fmt.Printf("  Shields: L%d N%d A%d\n", 
		fighterStats.LaserShield, fighterStats.NuclearShield, fighterStats.AntimatterShield)
}

// ExampleTemplateCreation demonstrates creating custom formation templates.
func ExampleTemplateCreation() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 10}},
			Destroyer: {{HP: 600, Count: 5}},
			Bomber:    {{HP: 500, Count: 3}},
		},
		Role: RoleTactical,
	}
	
	stack.SetFormation(FormationVanguard, time.Now())
	
	// Export current configuration as template
	template := ExportFormationTemplate(stack, "Alpha Strike Setup", 
		"Aggressive formation for first strike advantage")
	
	fmt.Printf("Template: %s\n", template.Name)
	fmt.Printf("Formation: %s\n", template.Formation)
	fmt.Printf("Assignments:\n")
	for shipType, position := range template.Assignments {
		fmt.Printf("  %s -> %s\n", shipType, position)
	}
	fmt.Printf("Conditions: %v\n", template.Conditions)
}

// ExampleDamageDistribution demonstrates how damage is distributed in formations.
func ExampleDamageDistribution() {
	defender := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 10}},
			Bomber:  {{HP: 500, Count: 5}},
			Scout:   {{HP: 100, Count: 8}},
		},
		Role: RoleTactical,
	}
	
	defender.SetFormation(FormationLine, time.Now())
	
	// Simulate incoming damage
	totalDamage := 1000
	directions := []AttackDirection{DirectionFrontal, DirectionFlanking, DirectionRear}
	
	for _, direction := range directions {
		distribution := defender.Formation.CalculateDamageDistribution(totalDamage, direction)
		
		fmt.Printf("\n%s attack (%d total damage):\n", direction, totalDamage)
		for position, damage := range distribution {
			percentage := float64(damage) / float64(totalDamage) * 100
			fmt.Printf("  %s: %d damage (%.1f%%)\n", position, damage, percentage)
		}
	}
}

// ExampleStackSpeed demonstrates formation impact on movement speed.
func ExampleStackSpeed() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 10}}, // Speed: 6
			Bomber:  {{HP: 500, Count: 5}},  // Speed: 5
			Scout:   {{HP: 100, Count: 8}},  // Speed: 9
		},
		Role: RoleTactical,
	}
	
	// Base speed (slowest ship)
	baseSpeed := stack.GetEffectiveStackSpeed()
	fmt.Printf("Base stack speed (no formation): %d\n", baseSpeed)
	
	// With Skirmish formation (1.2x speed)
	stack.SetFormation(FormationSkirmish, time.Now())
	skirmishSpeed := stack.GetEffectiveStackSpeed()
	fmt.Printf("Skirmish formation speed: %d\n", skirmishSpeed)
	
	// With Box formation (0.75x speed)
	stack.SetFormation(FormationBox, time.Now())
	boxSpeed := stack.GetEffectiveStackSpeed()
	fmt.Printf("Box formation speed: %d\n", boxSpeed)
	
	// With Vanguard formation (1.1x speed)
	stack.SetFormation(FormationVanguard, time.Now())
	vanguardSpeed := stack.GetEffectiveStackSpeed()
	fmt.Printf("Vanguard formation speed: %d\n", vanguardSpeed)
}

// ExampleFormationCloning demonstrates cloning formations.
func ExampleFormationCloning() {
	// Create original formation
	original := &Formation{
		Type:   FormationVanguard,
		Facing: "north",
		Assignments: []FormationAssignment{
			{Position: PositionFront, ShipType: Fighter, Count: 10},
			{Position: PositionBack, ShipType: Bomber, Count: 5},
		},
		CreatedAt: time.Now(),
	}
	
	// Clone it
	clone := CloneFormation(original)
	
	// Modify clone
	clone.Facing = "south"
	clone.Assignments[0].Count = 15
	
	// Original remains unchanged
	fmt.Printf("Original facing: %s, Front count: %d\n", 
		original.Facing, original.Assignments[0].Count)
	fmt.Printf("Clone facing: %s, Front count: %d\n", 
		clone.Facing, clone.Assignments[0].Count)
}
