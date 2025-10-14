package ships

import "testing"

// TestPositionSlotLimits verifies that slot limits are enforced correctly.
func TestPositionSlotLimits(t *testing.T) {
	tests := []struct {
		formation FormationType
		position  FormationPosition
	}{
		{FormationLine, PositionFront},
		{FormationLine, PositionFlank},
		{FormationLine, PositionBack},
		{FormationLine, PositionSupport},
		{FormationPhalanx, PositionFront},
		{FormationPhalanx, PositionFlank},
		{FormationSkirmish, PositionFlank},
		{FormationVanguard, PositionFront},
	}

	for _, tt := range tests {
		expected := len(GetInitialSlots(tt.formation, tt.position))
		got := GetMaxSlotsForPosition(tt.formation, tt.position)
		if got != expected {
			t.Errorf("GetMaxSlotsForPosition(%v, %v) = %d, want %d",
				tt.formation, tt.position, got, expected)
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
	maxSlots := GetMaxSlotsForPosition(formationType, position)

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
	maxSlots := GetMaxSlotsForPosition(formationType, position)

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
	formations := []FormationType{
		FormationLine,
		FormationBox,
		FormationVanguard,
		FormationSkirmish,
		FormationEchelon,
		FormationPhalanx,
		FormationSwarm,
	}

	for _, f := range formations {
		expected := GetMaxSlotsForPosition(f, PositionFront) +
			GetMaxSlotsForPosition(f, PositionFlank) +
			GetMaxSlotsForPosition(f, PositionBack) +
			GetMaxSlotsForPosition(f, PositionSupport)
		got := GetTotalMaxSlots(f)
		if got != expected {
			t.Errorf("GetTotalMaxSlots(%v) = %d, want %d", f, got, expected)
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

	// Verify each position is capped at its limit derived from predefined slots
	if max := GetMaxSlotsForPosition(formationType, PositionFront); len(snapshot.Positions[PositionFront]) > max {
		t.Errorf("Front slots %d exceed limit %d", len(snapshot.Positions[PositionFront]), max)
	}
	if max := GetMaxSlotsForPosition(formationType, PositionFlank); len(snapshot.Positions[PositionFlank]) > max {
		t.Errorf("Flank slots %d exceed limit %d", len(snapshot.Positions[PositionFlank]), max)
	}
	if max := GetMaxSlotsForPosition(formationType, PositionBack); len(snapshot.Positions[PositionBack]) > max {
		t.Errorf("Back slots %d exceed limit %d", len(snapshot.Positions[PositionBack]), max)
	}
	if max := GetMaxSlotsForPosition(formationType, PositionSupport); len(snapshot.Positions[PositionSupport]) > max {
		t.Errorf("Support slots %d exceed limit %d", len(snapshot.Positions[PositionSupport]), max)
	}
}
