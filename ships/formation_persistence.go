package ships

import (
    "fmt"
    "sort"
    "time"
)

// FormationSlotAssignment extends FormationAssignment with visual slot information.
// This allows saving user-arranged formations to MongoDB while preserving visual layout.
// The SlotIndex maps to specific coordinates from the formation layout system.
type FormationSlotAssignment struct {
	FormationAssignment        // Embedded: Position, ShipType, BucketIndex, Count, AssignedHP
	SlotIndex           int    `bson:"slotIndex" json:"slotIndex"`               // Which visual slot in this position (0-based)
	SlotKey             string `bson:"slotKey,omitempty" json:"slotKey,omitempty"` // Stable id: "x:y" from predefined Vector2
	IsManuallyPlaced    bool   `bson:"isManuallyPlaced" json:"isManuallyPlaced"` // User placed vs auto-assigned
}

// GetLayoutOccupancy returns the predefined layout with occupancy flags set from SlotAssignments.
func (fws *FormationWithSlots) GetLayoutOccupancy() map[string][]FormationLayoutPosition {
    mp, ok := getPredefinedMap(fws.Type)
    if !ok {
        return map[string][]FormationLayoutPosition{}
    }

    // Prepare result keyed by predefined position keys ("front","flank","back","support")
    result := make(map[string][]FormationLayoutPosition, len(mp))

    // Backfill SlotKey and group assignments by position key
    byPos := make(map[string][]FormationSlotAssignment)
    for i := range fws.SlotAssignments {
        a := fws.SlotAssignments[i]
        if a.SlotKey == "" {
            if coord, ok := GetNextSlotCoordinate(fws.Type, a.Position, a.SlotIndex); ok {
                fws.SlotAssignments[i].SlotKey = fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
                a.SlotKey = fws.SlotAssignments[i].SlotKey
            }
        }
        key := getPredefinedPositionKey(a.Position)
        byPos[key] = append(byPos[key], a)
    }

    for key, list := range mp {
        // Copy and sort by Order
        items := make([]FormationLayoutPosition, len(list))
        copy(items, list)
        sort.Slice(items, func(i, j int) bool { return items[i].Order < items[j].Order })

        // Clear occupancy fields
        for i := range items {
            items[i].Filled = false
            items[i].BucketIndex = nil
            items[i].Quantity = 0
        }

        // Apply assignments (resolve by SlotKey first, then SlotIndex)
        if assigns, ok := byPos[key]; ok {
            // Build coord-key -> index map
            indexByKey := make(map[string]int, len(items))
            for idx := range items {
                k := fmt.Sprintf("%.6f:%.6f", items[idx].Position.X, items[idx].Position.Y)
                indexByKey[k] = idx
            }
            for _, a := range assigns {
                idx := -1
                if a.SlotKey != "" {
                    if v, ok := indexByKey[a.SlotKey]; ok {
                        idx = v
                    }
                }
                if idx == -1 {
                    idx = a.SlotIndex
                }
                if idx >= 0 && idx < len(items) {
                    items[idx].Filled = a.Count > 0
                    b := a.BucketIndex
                    items[idx].BucketIndex = &b
                    items[idx].Quantity = a.Count
                }
            }
        }

        result[key] = items
    }

    return result
}

// FormationWithSlots extends Formation with visual slot assignments.
// Use this for saving/loading user-arranged formations.
type FormationWithSlots struct {
	Type             FormationType             `bson:"formationType" json:"formationType"`
	Facing           string                    `bson:"facing" json:"facing"`
	SlotAssignments  []FormationSlotAssignment `bson:"slotAssignments" json:"slotAssignments"`
	Modifiers        FormationMods             `bson:"modifiers" json:"modifiers"`
	CreatedAt        time.Time                 `bson:"createdAt" json:"createdAt"`
	Version          int                       `bson:"version" json:"version"`
}

// ToFormation converts FormationWithSlots to standard Formation for combat calculations.
// This strips visual information and returns the combat-relevant data.
func (fws *FormationWithSlots) ToFormation() Formation {
	assignments := make([]FormationAssignment, len(fws.SlotAssignments))
	for i, slotAssignment := range fws.SlotAssignments {
		assignments[i] = slotAssignment.FormationAssignment
	}

	return Formation{
		Type:        fws.Type,
		Facing:      fws.Facing,
		Assignments: assignments,
		Modifiers:   fws.Modifiers,
		CreatedAt:   fws.CreatedAt,
		Version:     fws.Version,
	}
}

// FromFormation creates FormationWithSlots from a standard Formation.
// Auto-assigns slots based on the formation layout system.
func FromFormation(formation Formation) FormationWithSlots {
	// Count assignments per position to determine slot indices
	positionCounts := make(map[FormationPosition]int)
	
	slotAssignments := make([]FormationSlotAssignment, len(formation.Assignments))
	
	for i, assignment := range formation.Assignments {
		// Get the next available slot index for this position
		slotIndex := positionCounts[assignment.Position]
		positionCounts[assignment.Position]++

		var slotKey string
		if coord, ok := GetNextSlotCoordinate(formation.Type, assignment.Position, slotIndex); ok {
			slotKey = fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
		}
		
		slotAssignments[i] = FormationSlotAssignment{
			FormationAssignment: assignment,
			SlotIndex:           slotIndex,
			SlotKey:             slotKey,
			IsManuallyPlaced:    false, // Auto-assigned
		}
	}

	return FormationWithSlots{
		Type:            formation.Type,
		Facing:          formation.Facing,
		SlotAssignments: slotAssignments,
		Modifiers:       formation.Modifiers,
		CreatedAt:       formation.CreatedAt,
		Version:         formation.Version,
	}
}

// GetSlotCoordinates returns the visual coordinates for all slot assignments.
// Returns a map of assignment index to coordinate.
func (fws *FormationWithSlots) GetSlotCoordinates() map[int]SlotCoordinate {
	coordinates := make(map[int]SlotCoordinate)
	
	for i, slotAssignment := range fws.SlotAssignments {
		coord, ok := GetNextSlotCoordinate(
			fws.Type,
			slotAssignment.Position,
			slotAssignment.SlotIndex,
		)
		if ok {
			coordinates[i] = coord
		}
	}
	
	return coordinates
}

// GenerateVisualSnapshot creates a complete visual layout with assignment metadata.
// This is the recommended function for sending formation data to the frontend.
func (fws *FormationWithSlots) GenerateVisualSnapshot() FormationVisualSnapshot {
	snapshot := FormationVisualSnapshot{
		FormationType: fws.Type,
		Facing:        fws.Facing,
		Assignments:   make([]AssignmentWithCoordinate, 0, len(fws.SlotAssignments)),
	}

	// Group assignments by position for bounds calculation
	slotCounts := make(map[FormationPosition]int)
	for _, assignment := range fws.SlotAssignments {
		if assignment.SlotIndex+1 > slotCounts[assignment.Position] {
			slotCounts[assignment.Position] = assignment.SlotIndex + 1
		}
	}

	// Generate assignment coordinates
	for _, slotAssignment := range fws.SlotAssignments {
		coord, ok := GetNextSlotCoordinate(
			fws.Type,
			slotAssignment.Position,
			slotAssignment.SlotIndex,
		)
		if !ok {
			coord = SlotCoordinate{X: 0, Y: 0} // Fallback
		}

		snapshot.Assignments = append(snapshot.Assignments, AssignmentWithCoordinate{
			Position:         slotAssignment.Position,
			ShipType:         slotAssignment.ShipType,
			BucketIndex:      slotAssignment.BucketIndex,
			Count:            slotAssignment.Count,
			AssignedHP:       slotAssignment.AssignedHP,
			SlotIndex:        slotAssignment.SlotIndex,
			Coordinate:       coord,
			IsManuallyPlaced: slotAssignment.IsManuallyPlaced,
		})
	}

	// Calculate bounds
	minX, maxX, minY, maxY := CalculateFormationBounds(fws.Type, slotCounts)
	snapshot.Bounds = LayoutBounds{
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
	}

	return snapshot
}

// FormationVisualSnapshot represents the complete visual state of a formation.
// This includes both combat data and visual layout information.
type FormationVisualSnapshot struct {
	FormationType FormationType              `json:"formationType"`
	Facing        string                     `json:"facing"`
	Assignments   []AssignmentWithCoordinate `json:"assignments"`
	Bounds        LayoutBounds               `json:"bounds"`
}

// AssignmentWithCoordinate combines assignment data with visual coordinates.
type AssignmentWithCoordinate struct {
	Position         FormationPosition `json:"position"`
	ShipType         ShipType          `json:"shipType"`
	BucketIndex      int               `json:"bucketIndex"`
	Count            int               `json:"count"`
	AssignedHP       int               `json:"assignedHP"`
	SlotIndex        int               `json:"slotIndex"`
	Coordinate       SlotCoordinate    `json:"coordinate"`
	IsManuallyPlaced bool              `json:"isManuallyPlaced"`
}

// MoveAssignmentToSlot moves an assignment to a different visual slot within the same position.
// This allows users to rearrange ships visually without changing combat mechanics.
func (fws *FormationWithSlots) MoveAssignmentToSlot(assignmentIndex int, newSlotIndex int) error {
	if assignmentIndex < 0 || assignmentIndex >= len(fws.SlotAssignments) {
		return ErrInvalidAssignmentIndex
	}

	assignment := &fws.SlotAssignments[assignmentIndex]
	
	// Verify the slot exists for this position
	_, ok := GetNextSlotCoordinate(fws.Type, assignment.Position, newSlotIndex)
	if !ok {
		return ErrInvalidSlotIndex
	}

	assignment.SlotIndex = newSlotIndex
    if coord, ok := GetNextSlotCoordinate(fws.Type, assignment.Position, newSlotIndex); ok {
        assignment.SlotKey = fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
    }
	assignment.IsManuallyPlaced = true
	
	return nil
}

// SwapAssignmentSlots swaps the visual positions of two assignments.
// Both assignments must be in the same formation position.
func (fws *FormationWithSlots) SwapAssignmentSlots(index1, index2 int) error {
	if index1 < 0 || index1 >= len(fws.SlotAssignments) ||
		index2 < 0 || index2 >= len(fws.SlotAssignments) {
		return ErrInvalidAssignmentIndex
	}

	assignment1 := &fws.SlotAssignments[index1]
	assignment2 := &fws.SlotAssignments[index2]

	if assignment1.Position != assignment2.Position {
		return ErrDifferentPositions
	}

	// Swap slot indices
	assignment1.SlotIndex, assignment2.SlotIndex = assignment2.SlotIndex, assignment1.SlotIndex
    if coord, ok := GetNextSlotCoordinate(fws.Type, assignment1.Position, assignment1.SlotIndex); ok {
        assignment1.SlotKey = fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
    }
    if coord, ok := GetNextSlotCoordinate(fws.Type, assignment2.Position, assignment2.SlotIndex); ok {
        assignment2.SlotKey = fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
    }
	assignment1.IsManuallyPlaced = true
	assignment2.IsManuallyPlaced = true

	return nil
}

// SplitAssignmentToSlot creates a new assignment by splitting ships from an existing bucket.
// This allows users to split HP buckets across multiple visual slots without actually splitting the bucket.
func (fws *FormationWithSlots) SplitAssignmentToSlot(
    sourceIndex int,
    splitCount int,
    targetSlotIndex int,
) error {
    return ErrPartialAllocationDisabled
}

// MergeAssignments combines two assignments from the same bucket into one.
// This is the reverse of SplitAssignmentToSlot.
func (fws *FormationWithSlots) MergeAssignments(index1, index2 int) error {
	if index1 < 0 || index1 >= len(fws.SlotAssignments) ||
		index2 < 0 || index2 >= len(fws.SlotAssignments) {
		return ErrInvalidAssignmentIndex
	}

	assignment1 := &fws.SlotAssignments[index1]
	assignment2 := &fws.SlotAssignments[index2]

	// Must be same position, ship type, and bucket
	if assignment1.Position != assignment2.Position ||
		assignment1.ShipType != assignment2.ShipType ||
		assignment1.BucketIndex != assignment2.BucketIndex {
		return ErrCannotMergeAssignments
	}

	// Merge into assignment1
	assignment1.Count += assignment2.Count
	assignment1.AssignedHP += assignment2.AssignedHP
	assignment1.IsManuallyPlaced = true

	// Remove assignment2
	fws.SlotAssignments = append(
		fws.SlotAssignments[:index2],
		fws.SlotAssignments[index2+1:]...,
	)

	return nil
}

// Errors for formation persistence operations
var (
	ErrInvalidAssignmentIndex  = &FormationError{Message: "invalid assignment index"}
	ErrInvalidSlotIndex        = &FormationError{Message: "invalid slot index for this position"}
	ErrDifferentPositions      = &FormationError{Message: "assignments must be in the same position"}
	ErrInvalidSplitCount       = &FormationError{Message: "split count must be between 1 and count-1"}
	ErrCannotMergeAssignments  = &FormationError{Message: "assignments cannot be merged (different position/type/bucket)"}
	ErrPartialAllocationDisabled = &FormationError{Message: "partial allocation disabled: assignments must be bucket-wide"}
)

// FormationError represents a formation operation error.
type FormationError struct {
	Message string
}

func (e *FormationError) Error() string {
	return e.Message
}
