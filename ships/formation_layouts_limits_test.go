package ships

import "testing"

// TestPositionSlotLimits verifies that slot limits are enforced correctly.
func TestPositionSlotLimits(t *testing.T) {
	tests := []struct {
		formation FormationType
		position  FormationPosition
		expected  int
	}{
		{FormationLine, PositionFront, 15},
		{FormationLine, PositionFlank, 10},
		{FormationLine, PositionBack, 15},
		{FormationLine, PositionSupport, 8},
		{FormationPhalanx, PositionFront, 25},
		{FormationPhalanx, PositionFlank, 6},
		{FormationSkirmish, PositionFlank, 20},
		{FormationVanguard, PositionFront, 20},
	}

	for _, tt := range tests {
		got := GetMaxSlotsForPosition(tt.formation, tt.position)
		if got != tt.expected {
			t.Errorf("GetMaxSlotsForPosition(%v, %v) = %d, want %d",
				tt.formation, tt.position, got, tt.expected)
		}
	}
}

// TestGetNextSlotCoordinateEnforcesLimit verifies that GetNextSlotCoordinate
// returns false when the position limit is reached.
func TestGetNextSlotCoordinateEnforcesLimit(t *testing.T) {
	formationType := FormationLine
	position := PositionFront
	maxSlots := GetMaxSlotsForPosition(formationType, position)

	// Should succeed up to the limit
	for i := 0; i < maxSlots; i++ {
		_, ok := GetNextSlotCoordinate(formationType, position, i)
		if !ok {
			t.Errorf("GetNextSlotCoordinate failed at slot %d/%d", i, maxSlots)
		}
	}

	// Should fail at the limit
	_, ok := GetNextSlotCoordinate(formationType, position, maxSlots)
	if ok {
		t.Errorf("GetNextSlotCoordinate should fail at limit %d", maxSlots)
	}

	// Should fail beyond the limit
	_, ok = GetNextSlotCoordinate(formationType, position, maxSlots+10)
	if ok {
		t.Errorf("GetNextSlotCoordinate should fail beyond limit %d", maxSlots)
	}
}

// TestGetAllSlotsForPositionCapsAtLimit verifies that GetAllSlotsForPosition
// automatically caps the returned slots at the position limit.
func TestGetAllSlotsForPositionCapsAtLimit(t *testing.T) {
	formationType := FormationPhalanx
	position := PositionFront
	maxSlots := GetMaxSlotsForPosition(formationType, position) // 25

	// Request more than the limit
	slots := GetAllSlotsForPosition(formationType, position, maxSlots+100)

	if len(slots) != maxSlots {
		t.Errorf("GetAllSlotsForPosition returned %d slots, want %d (capped at limit)",
			len(slots), maxSlots)
	}

	// Request exactly the limit
	slots = GetAllSlotsForPosition(formationType, position, maxSlots)
	if len(slots) != maxSlots {
		t.Errorf("GetAllSlotsForPosition returned %d slots, want %d",
			len(slots), maxSlots)
	}

	// Request less than the limit
	slots = GetAllSlotsForPosition(formationType, position, 5)
	if len(slots) != 5 {
		t.Errorf("GetAllSlotsForPosition returned %d slots, want 5",
			len(slots))
	}
}

// TestIsPositionFull verifies the position full check.
func TestIsPositionFull(t *testing.T) {
	formationType := FormationSkirmish
	position := PositionFlank
	maxSlots := GetMaxSlotsForPosition(formationType, position) // 20

	// Not full
	if IsPositionFull(formationType, position, maxSlots-1) {
		t.Errorf("IsPositionFull should return false at %d/%d", maxSlots-1, maxSlots)
	}

	// Exactly full
	if !IsPositionFull(formationType, position, maxSlots) {
		t.Errorf("IsPositionFull should return true at %d/%d", maxSlots, maxSlots)
	}

	// Over full
	if !IsPositionFull(formationType, position, maxSlots+5) {
		t.Errorf("IsPositionFull should return true at %d/%d", maxSlots+5, maxSlots)
	}
}

// TestGetTotalMaxSlots verifies total slot calculation.
func TestGetTotalMaxSlots(t *testing.T) {
	tests := []struct {
		formation FormationType
		expected  int
	}{
		{FormationLine, 48},     // 15+10+15+8
		{FormationBox, 44},      // 12+10+12+10
		{FormationVanguard, 44}, // 20+8+10+6
		{FormationSkirmish, 48}, // 8+20+12+8
		{FormationEchelon, 40},  // 10+12+10+8
		{FormationPhalanx, 49},  // 25+6+8+10
		{FormationSwarm, 48},    // 12+12+12+12
	}

	for _, tt := range tests {
		got := GetTotalMaxSlots(tt.formation)
		if got != tt.expected {
			t.Errorf("GetTotalMaxSlots(%v) = %d, want %d",
				tt.formation, got, tt.expected)
		}
	}
}

// TestFormationLayoutSnapshotRespectsLimits verifies that the snapshot
// generation respects position limits.
func TestFormationLayoutSnapshotRespectsLimits(t *testing.T) {
	formationType := FormationVanguard

	// Request excessive slots for each position
	slotCounts := map[FormationPosition]int{
		PositionFront:   100,
		PositionFlank:   100,
		PositionBack:    100,
		PositionSupport: 100,
	}

	snapshot := GenerateFormationLayoutSnapshot(formationType, slotCounts)

	// Verify each position is capped at its limit
	limits := FormationSlotLimits[formationType]

	if len(snapshot.Positions[PositionFront]) > limits.Front {
		t.Errorf("Front slots %d exceed limit %d",
			len(snapshot.Positions[PositionFront]), limits.Front)
	}
	if len(snapshot.Positions[PositionFlank]) > limits.Flank {
		t.Errorf("Flank slots %d exceed limit %d",
			len(snapshot.Positions[PositionFlank]), limits.Flank)
	}
	if len(snapshot.Positions[PositionBack]) > limits.Back {
		t.Errorf("Back slots %d exceed limit %d",
			len(snapshot.Positions[PositionBack]), limits.Back)
	}
	if len(snapshot.Positions[PositionSupport]) > limits.Support {
		t.Errorf("Support slots %d exceed limit %d",
			len(snapshot.Positions[PositionSupport]), limits.Support)
	}
}
