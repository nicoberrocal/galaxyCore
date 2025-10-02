package ships

import (
	"fmt"
	"testing"
	"time"
)

// TestDamageDistributionRedistribution tests that damage is properly redistributed
// when some formation positions are empty.
func TestDamageDistributionRedistribution(t *testing.T) {
	// Create a formation with only front and back positions filled
	formation := Formation{
		Type:   FormationLine,
		Facing: "north",
		Assignments: []FormationAssignment{
			{
				Position:   PositionFront,
				ShipType:   Fighter,
				Count:      10,
				AssignedHP: 1000, // 10 ships * 100 HP each
			},
			{
				Position:   PositionBack,
				ShipType:   Bomber,
				Count:      5,
				AssignedHP: 500, // 5 ships * 100 HP each
			},
			// Note: PositionFlank and PositionSupport are empty
		},
		CreatedAt: time.Now(),
		Version:   1,
	}

	incomingDamage := 1000
	direction := DirectionFrontal

	// Get the original weights for frontal attack
	originalWeights := DirectionalDamageWeights[DirectionFrontal]
	fmt.Printf("Original frontal attack weights: %+v\n", originalWeights)

	// Calculate damage distribution
	distribution := formation.CalculateDamageDistribution(incomingDamage, direction)

	fmt.Printf("Damage distribution result: %+v\n", distribution)

	// Verify that only filled positions receive damage
	if _, hasFlank := distribution[PositionFlank]; hasFlank {
		t.Errorf("Expected no damage to empty flank position, but got %d", distribution[PositionFlank])
	}
	if _, hasSupport := distribution[PositionSupport]; hasSupport {
		t.Errorf("Expected no damage to empty support position, but got %d", distribution[PositionSupport])
	}

	// Verify that filled positions receive damage
	frontDamage, hasFront := distribution[PositionFront]
	if !hasFront || frontDamage <= 0 {
		t.Errorf("Expected front position to receive damage, got %d", frontDamage)
	}

	backDamage, hasBack := distribution[PositionBack]
	if !hasBack || backDamage <= 0 {
		t.Errorf("Expected back position to receive damage, got %d", backDamage)
	}

	// Calculate total distributed damage
	totalDistributed := 0
	for _, damage := range distribution {
		totalDistributed += damage
	}

	// Verify that total distributed damage equals incoming damage (within rounding)
	if totalDistributed < incomingDamage-5 || totalDistributed > incomingDamage+5 {
		t.Errorf("Expected total distributed damage to be close to %d, got %d", incomingDamage, totalDistributed)
	}

	// Verify proportional redistribution
	// Original: front=0.6, flank=0.2, back=0.1, support=0.1
	// Filled: front=0.6, back=0.1 (total=0.7)
	// Redistributed: front=0.6/0.7≈0.857, back=0.1/0.7≈0.143
	expectedFrontRatio := 0.6 / 0.7
	expectedBackRatio := 0.1 / 0.7

	actualFrontRatio := float64(frontDamage) / float64(incomingDamage)
	actualBackRatio := float64(backDamage) / float64(incomingDamage)

	if abs(actualFrontRatio-expectedFrontRatio) > 0.01 {
		t.Errorf("Expected front ratio ~%.3f, got %.3f", expectedFrontRatio, actualFrontRatio)
	}
	if abs(actualBackRatio-expectedBackRatio) > 0.01 {
		t.Errorf("Expected back ratio ~%.3f, got %.3f", expectedBackRatio, actualBackRatio)
	}

	fmt.Printf("✓ Test passed: Damage properly redistributed from empty positions\n")
	fmt.Printf("  Front: %d damage (%.1f%% of total)\n", frontDamage, actualFrontRatio*100)
	fmt.Printf("  Back: %d damage (%.1f%% of total)\n", backDamage, actualBackRatio*100)
}

// TestEmptyFormationDamageDistribution tests behavior when formation is completely empty
func TestEmptyFormationDamageDistribution(t *testing.T) {
	formation := Formation{
		Type:        FormationLine,
		Facing:      "north",
		Assignments: []FormationAssignment{}, // No assignments
		CreatedAt:   time.Now(),
		Version:     1,
	}

	distribution := formation.CalculateDamageDistribution(1000, DirectionFrontal)

	if len(distribution) != 0 {
		t.Errorf("Expected empty distribution for empty formation, got %+v", distribution)
	}

	fmt.Printf("✓ Test passed: Empty formation returns empty distribution\n")
}

// TestSinglePositionFormation tests when only one position is filled
func TestSinglePositionFormation(t *testing.T) {
	formation := Formation{
		Type:   FormationLine,
		Facing: "north",
		Assignments: []FormationAssignment{
			{
				Position:   PositionFront,
				ShipType:   Fighter,
				Count:      10,
				AssignedHP: 1000,
			},
		},
		CreatedAt: time.Now(),
		Version:   1,
	}

	incomingDamage := 1000
	distribution := formation.CalculateDamageDistribution(incomingDamage, DirectionFrontal)

	// All damage should go to the single filled position
	if len(distribution) != 1 {
		t.Errorf("Expected exactly 1 position in distribution, got %d", len(distribution))
	}

	frontDamage, hasFront := distribution[PositionFront]
	if !hasFront {
		t.Errorf("Expected front position to receive damage")
	}

	if frontDamage != incomingDamage {
		t.Errorf("Expected front to receive all %d damage, got %d", incomingDamage, frontDamage)
	}

	fmt.Printf("✓ Test passed: Single position receives all damage (%d)\n", frontDamage)
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
