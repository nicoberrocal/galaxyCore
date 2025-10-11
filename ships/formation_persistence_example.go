package ships

import (
	"encoding/json"
	"fmt"
	"time"
)

// ExampleFormationPersistence demonstrates how to save and load formations with visual layouts.
func ExampleFormationPersistence() {
	fmt.Println("=== Formation Persistence Example ===")
	fmt.Println()

	// Step 1: Create a standard formation (from combat system)
	now := time.Now()
	formation := Formation{
		Type:   FormationVanguard,
		Facing: "north",
		Assignments: []FormationAssignment{
			{
				Position:    PositionFront,
				ShipType:    Fighter,
				BucketIndex: 0,
				Count:       3,
				AssignedHP:  300,
			},
			{
				Position:    PositionFront,
				ShipType:    Destroyer,
				BucketIndex: 0,
				Count:       2,
				AssignedHP:  400,
			},
			{
				Position:    PositionFlank,
				ShipType:    Scout,
				BucketIndex: 0,
				Count:       5,
				AssignedHP:  250,
			},
			{
				Position:    PositionBack,
				ShipType:    Bomber,
				BucketIndex: 0,
				Count:       2,
				AssignedHP:  200,
			},
		},
		CreatedAt: now,
		Version:   1,
	}

	// Step 2: Convert to FormationWithSlots for visual layout
	formationWithSlots := FromFormation(formation)
	fmt.Println("Converted to FormationWithSlots:")
	fmt.Printf("  Type: %s\n", formationWithSlots.Type)
	fmt.Printf("  Assignments: %d\n", len(formationWithSlots.SlotAssignments))
	for i, assignment := range formationWithSlots.SlotAssignments {
		fmt.Printf("    [%d] %s at %s, Slot %d: %d ships\n",
			i, assignment.ShipType, assignment.Position, assignment.SlotIndex, assignment.Count)
	}

	// Step 3: Generate visual snapshot for frontend
	fmt.Println("\n=== Visual Snapshot ===")
	snapshot := formationWithSlots.GenerateVisualSnapshot()
	for _, assignment := range snapshot.Assignments {
		fmt.Printf("%s %s (Slot %d): %d ships at (%.2f, %.2f)\n",
			assignment.Position, assignment.ShipType, assignment.SlotIndex,
			assignment.Count, assignment.Coordinate.X, assignment.Coordinate.Y)
	}
	fmt.Printf("Bounds: X[%.2f, %.2f], Y[%.2f, %.2f]\n\n",
		snapshot.Bounds.MinX, snapshot.Bounds.MaxX,
		snapshot.Bounds.MinY, snapshot.Bounds.MaxY)

	// Step 4: Save to MongoDB (as JSON for this example)
	jsonData, _ := json.MarshalIndent(formationWithSlots, "", "  ")
	fmt.Println("=== MongoDB Document (JSON) ===")
	fmt.Println(string(jsonData))
}

// ExampleFormationUserArrangement shows how users can rearrange ships visually.
func ExampleFormationUserArrangement() {
	fmt.Println()
	fmt.Println("=== User Arrangement Example ===")
	fmt.Println()

	// Start with a formation
	formation := Formation{
		Type:   FormationPhalanx,
		Facing: "north",
		Assignments: []FormationAssignment{
			{Position: PositionFront, ShipType: Fighter, BucketIndex: 0, Count: 5, AssignedHP: 500},
			{Position: PositionFront, ShipType: Destroyer, BucketIndex: 0, Count: 3, AssignedHP: 600},
			{Position: PositionFlank, ShipType: Scout, BucketIndex: 0, Count: 4, AssignedHP: 200},
		},
		CreatedAt: time.Now(),
		Version:   1,
	}

	fws := FromFormation(formation)
	fmt.Println("Initial Layout:")
	printFormationLayout(&fws)

	// User wants to swap the Fighters and Destroyers in front
	fmt.Println("\nUser swaps Fighters and Destroyers...")
	err := fws.SwapAssignmentSlots(0, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Swapped successfully!")
		printFormationLayout(&fws)
	}

	// User wants to move Scouts to a different slot
	fmt.Println("\nUser moves Scouts to slot 3...")
	err = fws.MoveAssignmentToSlot(2, 3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Moved successfully!")
		printFormationLayout(&fws)
	}
}

// ExampleFormationSplitting shows how to split HP buckets across visual slots.
func ExampleFormationSplitting() {
	fmt.Println()
	fmt.Println("=== HP Bucket Splitting Example ===")
	fmt.Println()

	// Formation with a large fighter group
	formation := Formation{
		Type:   FormationLine,
		Facing: "north",
		Assignments: []FormationAssignment{
			{
				Position:    PositionFront,
				ShipType:    Fighter,
				BucketIndex: 0,
				Count:       10, // 10 fighters in one bucket
				AssignedHP:  1000,
			},
		},
		CreatedAt: time.Now(),
		Version:   1,
	}

	fws := FromFormation(formation)
	fmt.Println("Initial: 10 Fighters in slot 0")
	printFormationLayout(&fws)

	// User wants to split 4 fighters to slot 1 (visual spread)
	fmt.Println("\nSplitting 4 fighters to slot 1...")
	err := fws.SplitAssignmentToSlot(0, 4, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Split successfully!")
		printFormationLayout(&fws)
		
		// Show they're still from the same bucket
		fmt.Println("\nNote: Both assignments reference BucketIndex 0:")
		for i, assignment := range fws.SlotAssignments {
			fmt.Printf("  [%d] Bucket %d: %d ships\n",
				i, assignment.BucketIndex, assignment.Count)
		}
	}

	// User can merge them back
	fmt.Println("\nMerging assignments back together...")
	err = fws.MergeAssignments(0, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Merged successfully!")
		printFormationLayout(&fws)
	}
}

// ExampleFormationRoundTrip shows saving to MongoDB and loading back.
func ExampleFormationRoundTrip() {
	fmt.Println()
	fmt.Println("=== MongoDB Round Trip Example ===")
	fmt.Println()

	// Create formation with user arrangement
	formation := Formation{
		Type:   FormationVanguard,
		Facing: "north",
		Assignments: []FormationAssignment{
			{Position: PositionFront, ShipType: Fighter, BucketIndex: 0, Count: 5, AssignedHP: 500},
			{Position: PositionFlank, ShipType: Scout, BucketIndex: 0, Count: 3, AssignedHP: 150},
		},
		CreatedAt: time.Now(),
		Version:   1,
	}

	fws := FromFormation(formation)
	
	// User customizes
	_ = fws.MoveAssignmentToSlot(0, 2) // Move fighters to slot 2
	
	fmt.Println("Before saving to MongoDB:")
	printFormationLayout(&fws)

	// Simulate MongoDB save/load
	jsonData, _ := json.Marshal(fws)
	fmt.Printf("\nMongoDB document size: %d bytes\n", len(jsonData))

	// Load back from MongoDB
	var loadedFws FormationWithSlots
	_ = json.Unmarshal(jsonData, &loadedFws)

	fmt.Println("\nAfter loading from MongoDB:")
	printFormationLayout(&loadedFws)

	// Convert back to standard Formation for combat
	combatFormation := loadedFws.ToFormation()
	fmt.Printf("\nConverted to combat Formation: %d assignments\n", len(combatFormation.Assignments))
	fmt.Println("(Visual slot info stripped, ready for combat calculations)")
}

// Helper function to print formation layout
func printFormationLayout(fws *FormationWithSlots) {
	for i, assignment := range fws.SlotAssignments {
		coord, _ := GetNextSlotCoordinate(fws.Type, assignment.Position, assignment.SlotIndex)
		manualTag := ""
		if assignment.IsManuallyPlaced {
			manualTag = " [manual]"
		}
		fmt.Printf("  [%d] %s at %s Slot %d (%.1f, %.1f): %d ships%s\n",
			i, assignment.ShipType, assignment.Position, assignment.SlotIndex,
			coord.X, coord.Y, assignment.Count, manualTag)
	}
}
