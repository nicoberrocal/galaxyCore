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

// ExampleFormationCombat demonstrates combat with formations using V2 compute system.
func ExampleFormationCombat() {
	// Attacker with Vanguard formation
	attacker := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 15}},
			Destroyer: {{HP: 600, Count: 5}},
		},
		Role: RoleTactical,
	}
	now := time.Now()
	attacker.SetFormation(FormationVanguard, now)
	
	// Defender with Box formation
	defender := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 12}},
			Carrier: {{HP: 900, Count: 2}},
		},
		Role: RoleTactical,
	}
	defender.SetFormation(FormationBox, now)
	
	// Setup combat context
	combatCtx := NewCombatContext(attacker, defender)
	fmt.Printf("Formation Counter: %.2fx (Vanguard vs Box)\n", combatCtx.FormationCounter)
	fmt.Printf("Attack Direction: %s\n", combatCtx.AttackDirection)
	
	// Calculate effective combat stats using V2
	attackerFighter, _, _ := ComputeEffectiveShipV2(
		attacker, Fighter, 0, now, true, defender.Formation.Type)
	attackerDestroyer, _, _ := ComputeEffectiveShipV2(
		attacker, Destroyer, 0, now, true, defender.Formation.Type)
	
	defenderFighter, _, _ := ComputeEffectiveShipV2(
		defender, Fighter, 0, now, true, attacker.Formation.Type)
	defenderCarrier, _, _ := ComputeEffectiveShipV2(
		defender, Carrier, 0, now, true, attacker.Formation.Type)
	
	fmt.Printf("\nAttacker effective stats (in combat vs %s):\n", defender.Formation.Type)
	fmt.Printf("  Fighter: %d damage, %d HP\n", attackerFighter.AttackDamage, attackerFighter.HP)
	fmt.Printf("  Destroyer: %d damage, %d HP\n", attackerDestroyer.AttackDamage, attackerDestroyer.HP)
	
	fmt.Printf("\nDefender effective stats (in combat vs %s):\n", attacker.Formation.Type)
	fmt.Printf("  Fighter: %d damage, %d HP\n", defenderFighter.AttackDamage, defenderFighter.HP)
	fmt.Printf("  Carrier: %d damage, %d HP\n", defenderCarrier.AttackDamage, defenderCarrier.HP)
	
	// Demonstrate damage distribution
	totalAttackerDamage := int(float64(attackerFighter.AttackDamage*15+attackerDestroyer.AttackDamage*5) * combatCtx.FormationCounter)
	damageDistribution := defender.Formation.CalculateDamageDistribution(totalAttackerDamage, combatCtx.AttackDirection)
	
	fmt.Printf("\nDamage distribution to defender (%d total):\n", totalAttackerDamage)
	for position, damage := range damageDistribution {
		percentage := float64(damage) / float64(totalAttackerDamage) * 100
		fmt.Printf("  %s: %d damage (%.1f%%)\n", position, damage, percentage)
	}
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

// ExampleGemPositionSynergy demonstrates gem-position bonuses using V2 compute system.
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
	now := time.Now()
	
	// Get effective stats for Fighter in front position using V2
	fighterStats, fighterAbilities, fighterModStack := ComputeEffectiveShipV2(
		stack, Fighter, 0, now, false, "")
	fmt.Printf("Fighter in Front position (V2):\n")
	fmt.Printf("  Attack Damage: %d\n", fighterStats.AttackDamage)
	fmt.Printf("  Laser Shield: %d\n", fighterStats.LaserShield)
	fmt.Printf("  HP: %d\n", fighterStats.HP)
	fmt.Printf("  Active Abilities: %d\n", len(fighterAbilities))
	
	// Show modifier breakdown for Fighter
	fighterBreakdown := GetModifierBreakdown(stack, Fighter, 0, now, false, "")
	fmt.Printf("  Modifier layers: %d\n", len(fighterBreakdown))
	for _, mod := range fighterBreakdown {
		if mod.IsActive {
			fmt.Printf("    - %s: %s\n", mod.Source, mod.Description)
		}
	}
	
	// Get effective stats for Bomber in back position using V2
	bomberStats, bomberAbilities, _ := ComputeEffectiveShipV2(
		stack, Bomber, 0, now, false, "")
	fmt.Printf("\nBomber in Back position (V2):\n")
	fmt.Printf("  Attack Damage: %d\n", bomberStats.AttackDamage)
	fmt.Printf("  Attack Range: %d\n", bomberStats.AttackRange)
	fmt.Printf("  Visibility: %d\n", bomberStats.VisibilityRange)
	fmt.Printf("  Active Abilities: %d\n", len(bomberAbilities))
	
	_ = fighterModStack // Available for further analysis if needed
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

// ExampleRoleModeFormationSynergy demonstrates role-formation interactions using V2 compute system.
func ExampleRoleModeFormationSynergy() {
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter: {{HP: 200, Count: 10}},
			Bomber:  {{HP: 500, Count: 5}},
		},
		Role: RoleTactical,
	}
	now := time.Now()
	
	// Tactical mode gets -30% reconfiguration time
	stack.SetFormation(FormationVanguard, now)
	fmt.Printf("Tactical mode reconfiguration: %d seconds\n", 
		RoleModeFormationBonus(RoleTactical, FormationCatalog[FormationVanguard].ReconfigureTime))
	
	// Get stats with Tactical+Vanguard using V2
	fighterTactical, _, tacticalModStack := ComputeEffectiveShipV2(
		stack, Fighter, 0, now, false, "")
	fmt.Printf("\nFighter stats with Tactical+Vanguard (V2):\n")
	fmt.Printf("  Attack Damage: %d\n", fighterTactical.AttackDamage)
	fmt.Printf("  Shields: L%d N%d A%d\n", 
		fighterTactical.LaserShield, fighterTactical.NuclearShield, fighterTactical.AntimatterShield)
	
	// Switch to Economic mode
	stack.Role = RoleEconomic
	stack.SetFormation(FormationBox, now)
	fmt.Printf("\nEconomic mode reconfiguration: %d seconds\n",
		RoleModeFormationBonus(RoleEconomic, FormationCatalog[FormationBox].ReconfigureTime))
	
	// Get stats with Economic+Box using V2
	fighterEconomic, _, economicModStack := ComputeEffectiveShipV2(
		stack, Fighter, 0, now, false, "")
	fmt.Printf("\nFighter stats with Economic+Box (V2):\n")
	fmt.Printf("  Attack Damage: %d\n", fighterEconomic.AttackDamage)
	fmt.Printf("  Shields: L%d N%d A%d\n", 
		fighterEconomic.LaserShield, fighterEconomic.NuclearShield, fighterEconomic.AntimatterShield)
	
	// Compare the two configurations
	diff := DiffModifierStacks(tacticalModStack, economicModStack)
	fmt.Printf("\nConfiguration differences:\n")
	fmt.Printf("  Added modifiers: %d\n", len(diff.Added))
	fmt.Printf("  Removed modifiers: %d\n", len(diff.Removed))
	fmt.Printf("  Changed modifiers: %d\n", len(diff.Changed))
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

// ExampleV2ComputeWorkflow demonstrates the complete V2 compute workflow.
func ExampleV2ComputeWorkflow() {
	// Create a complex stack with gems, role, and formation
	stack := &ShipStack{
		Ships: map[ShipType][]HPBucket{
			Fighter:   {{HP: 200, Count: 12}},
			Destroyer: {{HP: 600, Count: 4}},
			Bomber:    {{HP: 500, Count: 3}},
		},
		Role: RoleTactical,
		Loadouts: map[ShipType]ShipLoadout{
			Fighter: {
				Sockets: []Gem{
					GemCatalog["laser-3"],
					GemCatalog["kinetic-2"],
				},
				Anchored: false,
			},
			Destroyer: {
				Sockets: []Gem{
					GemCatalog["nuclear-3"],
					GemCatalog["armor-2"],
				},
				Anchored: false,
			},
		},
	}
	
	now := time.Now()
	stack.SetFormation(FormationVanguard, now)
	
	fmt.Printf("=== V2 Compute Workflow Example ===\n")
	
	// 1. Get baseline stats (no formation, no combat)
	baselineFighter, _, _ := ComputeEffectiveShipV2(stack, Fighter, 0, now, false, "")
	fmt.Printf("\n1. Baseline Fighter (gems + role only):\n")
	fmt.Printf("   Damage: %d, HP: %d, Shields: L%d N%d A%d\n",
		baselineFighter.AttackDamage, baselineFighter.HP,
		baselineFighter.LaserShield, baselineFighter.NuclearShield, baselineFighter.AntimatterShield)
	
	// 2. Add formation bonuses (out of combat)
	formationFighter, _, formationModStack := ComputeEffectiveShipV2(stack, Fighter, 0, now, false, "")
	fmt.Printf("\n2. With Vanguard Formation (out of combat):\n")
	fmt.Printf("   Damage: %d, HP: %d, Shields: L%d N%d A%d\n",
		formationFighter.AttackDamage, formationFighter.HP,
		formationFighter.LaserShield, formationFighter.NuclearShield, formationFighter.AntimatterShield)
	
	// 3. In combat against enemy formation
	combatFighter, _, combatModStack := ComputeEffectiveShipV2(stack, Fighter, 0, now, true, FormationBox)
	fmt.Printf("\n3. In Combat vs Box Formation:\n")
	fmt.Printf("   Damage: %d, HP: %d, Shields: L%d N%d A%d\n",
		combatFighter.AttackDamage, combatFighter.HP,
		combatFighter.LaserShield, combatFighter.NuclearShield, combatFighter.AntimatterShield)
	
	// 4. Show modifier breakdown
	breakdown := GetModifierBreakdown(stack, Fighter, 0, now, true, FormationBox)
	fmt.Printf("\n4. Active Modifier Layers (%d total):\n", len(breakdown))
	for i, mod := range breakdown {
		if mod.IsActive {
			status := "Active"
			if mod.ExpiresIn != nil {
				status = fmt.Sprintf("Expires in %.1fs", *mod.ExpiresIn)
			}
			fmt.Printf("   %d. %s: %s [%s]\n", i+1, mod.Source, mod.Description, status)
		}
	}
	
	// 5. Compare different scenarios
	diff := DiffModifierStacks(formationModStack, combatModStack)
	fmt.Printf("\n5. Formation vs Combat Differences:\n")
	fmt.Printf("   Added in combat: %d modifiers\n", len(diff.Added))
	fmt.Printf("   Removed in combat: %d modifiers\n", len(diff.Removed))
	fmt.Printf("   Changed in combat: %d modifiers\n", len(diff.Changed))
	
	// Show specific changes
	for _, added := range diff.Added {
		fmt.Printf("   + %s: %s\n", added.Source, added.Description)
	}
	
	// 6. Demonstrate manual modifier building
	fmt.Printf("\n6. Manual Modifier Building:\n")
	builder := NewModifierBuilder(now)
	builder.AddRoleMode(RoleTactical)
	builder.AddFormationPosition(stack.Formation, PositionFront)
	builder.AddFormationCounter(FormationVanguard, FormationBox, true)
	
	manualStack := builder.Build()
	ctx := ResolveContext{
		Now:             now,
		InCombat:        true,
		HasFormation:    true,
		FormationType:   FormationVanguard,
		EnemyFormation:  FormationBox,
	}
	manualMods := manualStack.Resolve(ctx)
	
	fmt.Printf("   Manual stack layers: %d\n", len(manualStack.Layers))
	fmt.Printf("   Final damage bonus: %.1f%%\n", manualMods.Damage.LaserPct*100)
}

// ExampleAdvancedModifierManagement demonstrates advanced modifier stack operations.
func ExampleAdvancedModifierManagement() {
	now := time.Now()
	stack := NewModifierStack()
	
	fmt.Printf("=== Advanced Modifier Management ===\n")
	
	// 1. Add various types of modifiers
	stack.AddPermanent(SourceGem, "laser-gem", "Laser Gem Tier 3",
		StatMods{Damage: DamageMods{LaserPct: 0.15}}, PriorityGem, now)
	
	stack.AddTemporary(SourceAbility, "shield-boost", "Shield Boost Ability",
		StatMods{LaserShieldDelta: 3}, PriorityAbility, now, 30*time.Second)
	
	combatOnly := true
	stack.AddConditional(SourceFormationCounter, "vanguard-vs-box", "Formation Advantage",
		StatMods{Damage: DamageMods{LaserPct: 0.30}}, PrioritySynergy, now, &combatOnly, nil)
	
	fmt.Printf("1. Initial stack: %d layers\n", len(stack.Layers))
	
	// 2. Resolve in different contexts
	oocCtx := ResolveContext{Now: now, InCombat: false, HasFormation: true}
	combatCtx := ResolveContext{Now: now, InCombat: true, HasFormation: true}
	
	oocMods := stack.Resolve(oocCtx)
	combatMods := stack.Resolve(combatCtx)
	
	fmt.Printf("2. Out of combat damage bonus: %.1f%%\n", oocMods.Damage.LaserPct*100)
	fmt.Printf("   In combat damage bonus: %.1f%%\n", combatMods.Damage.LaserPct*100)
	
	// 3. Remove expired modifiers
	futureTime := now.Add(45 * time.Second)
	stack.RemoveExpired(futureTime)
	fmt.Printf("3. After 45s, remaining layers: %d\n", len(stack.Layers))
	
	// 4. Remove by source
	stack.RemoveBySource(SourceGem)
	fmt.Printf("4. After removing gems: %d layers\n", len(stack.Layers))
	
	// 5. Get summary for UI
	summary := stack.GetSummary(combatCtx)
	fmt.Printf("5. Summary for UI:\n")
	for _, mod := range summary {
		status := "Inactive"
		if mod.IsActive {
			status = "Active"
			if mod.ExpiresIn != nil {
				status = fmt.Sprintf("Active (%.1fs left)", *mod.ExpiresIn)
			}
		}
		fmt.Printf("   %s: %s [%s]\n", mod.Source, mod.Description, status)
	}
}
