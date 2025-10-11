package ships

import (
	"encoding/json"
	"fmt"
)

// ExampleFormationLayoutUsage demonstrates how to use the formation layout system.
func ExampleFormationLayoutUsage() {
	// Example 1: Get initial slots for a formation position
	fmt.Println("=== Example 1: Initial Slots ===")
	phalanxFrontSlots := GetInitialSlots(FormationPhalanx, PositionFront)
	fmt.Printf("Phalanx Front initial slots: %d\n", len(phalanxFrontSlots))
	for i, slot := range phalanxFrontSlots {
		fmt.Printf("  Slot %d: (%.1f, %.1f)\n", i, slot.X, slot.Y)
	}

	// Example 2: Get the next slot when expanding
	fmt.Println("\n=== Example 2: Dynamic Expansion ===")
	// Phalanx starts with 3 front slots, let's add more
	for i := 3; i < 8; i++ {
		coord, ok := GetNextSlotCoordinate(FormationPhalanx, PositionFront, i)
		if ok {
			fmt.Printf("Slot %d (expansion): (%.1f, %.1f)\n", i, coord.X, coord.Y)
		}
	}

	// Example 3: Get all slots for a position
	fmt.Println("\n=== Example 3: All Slots for Position ===")
	vanguardFrontSlots := GetAllSlotsForPosition(FormationVanguard, PositionFront, 5)
	fmt.Printf("Vanguard Front with 5 ships:\n")
	for i, slot := range vanguardFrontSlots {
		fmt.Printf("  Ship %d: (%.2f, %.2f)\n", i, slot.X, slot.Y)
	}

	// Example 4: Generate a complete formation snapshot
	fmt.Println("\n=== Example 4: Complete Formation Snapshot ===")
	// Simulate a fleet with ships in different positions
	slotCounts := map[FormationPosition]int{
		PositionFront:   5,
		PositionFlank:   4,
		PositionBack:    2,
		PositionSupport: 3,
	}

	snapshot := GenerateFormationLayoutSnapshot(FormationVanguard, slotCounts)
	fmt.Printf("Formation: %s\n", snapshot.FormationType)
	fmt.Printf("Bounds: X[%.2f, %.2f], Y[%.2f, %.2f]\n",
		snapshot.Bounds.MinX, snapshot.Bounds.MaxX,
		snapshot.Bounds.MinY, snapshot.Bounds.MaxY)

	for position, slots := range snapshot.Positions {
		fmt.Printf("\n%s (%d ships):\n", position, len(slots))
		for _, slot := range slots {
			initialTag := ""
			if slot.IsInitial {
				initialTag = " [initial]"
			}
			fmt.Printf("  Slot %d: (%.2f, %.2f)%s\n",
				slot.SlotIndex, slot.Coordinate.X, slot.Coordinate.Y, initialTag)
		}
	}

	// Example 5: JSON output for frontend
	fmt.Println("\n=== Example 5: JSON for Frontend ===")
	jsonData, _ := json.MarshalIndent(snapshot, "", "  ")
	fmt.Println(string(jsonData))
}

// ExampleFormationComparison shows how different formations arrange ships differently.
func ExampleFormationComparison() {
	fmt.Println("=== Formation Visual Comparison ===")
	fmt.Println("Showing front position with 5 ships for each formation type:")
	fmt.Println()

	formations := []FormationType{
		FormationLine,
		FormationBox,
		FormationVanguard,
		FormationSkirmish,
		FormationEchelon,
		FormationPhalanx,
		FormationSwarm,
	}

	for _, formation := range formations {
		fmt.Printf("%s:\n", formation)
		slots := GetAllSlotsForPosition(formation, PositionFront, 5)
		for i, slot := range slots {
			fmt.Printf("  Slot %d: (%.2f, %.2f)\n", i+1, slot.X, slot.Y)
		}
		fmt.Println()
		fmt.Println("---")
		fmt.Println()
	}
}

// ExampleFormationGrowth demonstrates how a formation grows as ships are added.
func ExampleFormationGrowth() {
	fmt.Println("=== Formation Growth Simulation ===")
	fmt.Println("Watching Phalanx formation grow from 1 to 10 front ships:")
	fmt.Println("")

	for shipCount := 1; shipCount <= 10; shipCount++ {
		slots := GetAllSlotsForPosition(FormationPhalanx, PositionFront, shipCount)
		fmt.Printf("With %2d ships: ", shipCount)

		// Calculate width
		minX, maxX := 0.0, 0.0
		for i, slot := range slots {
			if i == 0 || slot.X < minX {
				minX = slot.X
			}
			if i == 0 || slot.X > maxX {
				maxX = slot.X
			}
		}
		width := maxX - minX
		fmt.Printf("Width = %.2f units (from %.2f to %.2f)\n", width, minX, maxX)
	}

	fmt.Println("\nNote: Phalanx front line gets progressively wider,")
	fmt.Println("reflecting its 'extreme_flank_weakness' special property.")
}

// ExampleSwarmHexPattern demonstrates the hexagonal dispersion pattern of Swarm formation.
func ExampleSwarmHexPattern() {
	fmt.Println("=== Swarm Hexagonal Pattern ===")
	fmt.Println("Swarm formation uses hexagonal rings for anti-AoE dispersion:")
	fmt.Println()

	// Show first 18 ships (3 complete hexagonal rings)
	slots := GetAllSlotsForPosition(FormationSwarm, PositionFront, 18)

	for i, slot := range slots {
		ring := i / 6
		posInRing := i % 6
		fmt.Printf("Ship %2d (Ring %d, Pos %d): (X: %6.2f, Y: %6.2f)\n",
			i, ring, posInRing, slot.X, slot.Y)
	}

	fmt.Println("\nEach ring has 6 positions at 60° intervals,")
	fmt.Println("creating maximum dispersion to counter AoE attacks.")
}
