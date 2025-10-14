package ships

import (
 "math"
 "sort"
)

// SlotCoordinate defines a 2D position for a ship slot within a formation's visual layout.
// The coordinate system is centered at (0,0).
// X represents the horizontal axis (positive is right, negative is left).
// Y represents the vertical axis (positive is forward/front, negative is back/rear).
// Units are abstract and should be scaled by the frontend based on ship size.
type SlotCoordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// FormationLayoutSnapshot represents a complete visual layout of a formation at a point in time.
// This is useful for sending to the frontend for rendering.
type FormationLayoutSnapshot struct {
    FormationType FormationType                        `json:"formationType"`
    Positions     map[FormationPosition][]SlotPosition `json:"positions"`
    Bounds        LayoutBounds                         `json:"bounds"`
}

// SlotPosition represents a single ship slot with its coordinate and metadata.
type SlotPosition struct {
    Coordinate SlotCoordinate `json:"coordinate"`
    SlotIndex  int            `json:"slotIndex"` // Index within the position (0-based)
    IsInitial  bool           `json:"isInitial"` // True if this is from initial layout, false if expanded
}

// LayoutBounds represents the bounding box of a formation layout.
type LayoutBounds struct {
    MinX float64 `json:"minX"`
    MaxX float64 `json:"maxX"`
    MinY float64 `json:"minY"`
    MaxY float64 `json:"maxY"`
}

// getPredefinedPositionKey maps FormationPosition to the key used in predefined layout maps.
func getPredefinedPositionKey(position FormationPosition) string {
	switch position {
	case PositionFront:
		return "front"
	case PositionFlank:
		return "flank"
	case PositionBack:
		return "back"
	case PositionSupport:
		return "support"
	default:
		return "front"
	}
}

// getPredefinedMap returns the predefined layout map for a formation type.
func getPredefinedMap(formationType FormationType) (map[string][]FormationLayoutPosition, bool) {
	switch formationType {
	case FormationLine:
		return LineFormation, true
	case FormationBox:
		return BoxFormation, true
	case FormationVanguard:
		return VanguardFormation, true
	case FormationSkirmish:
		return SkirmishFormation, true
	case FormationEchelon:
		return EchelonFormation, true
	case FormationPhalanx:
		return PhalanxFormation, true
	case FormationSwarm:
		return SwarmFormation, true
	default:
		return nil, false
	}
}

// getPredefinedSlots returns all predefined coordinates for a given formation and position,
// sorted by Order. If not found, returns an empty slice.
func getPredefinedSlots(formationType FormationType, position FormationPosition) []SlotCoordinate {
	if mp, ok := getPredefinedMap(formationType); ok {
		key := getPredefinedPositionKey(position)
		if list, ok := mp[key]; ok {
			// Ensure stable ordering by Order field
			items := make([]FormationLayoutPosition, len(list))
			copy(items, list)
			sort.Slice(items, func(i, j int) bool { return items[i].Order < items[j].Order })

			coords := make([]SlotCoordinate, 0, len(items))
			for _, it := range items {
				coords = append(coords, SlotCoordinate{X: it.Position.X, Y: it.Position.Y})
			}
			return coords
		}
	}
	return []SlotCoordinate{}
}

// FormationLayouts defines the initial visual slot coordinates for each position within a formation.
// These are the base positions shown when a formation is first created.
// Design principles:
// - Y-axis: positive=front, negative=back (matches directional damage weights)
// - Visual density reflects tactical role (e.g., Phalanx has wide front, Vanguard has concentrated spearhead)
// - Layouts support the formation counter matrix relationships
// - Each position (Front/Flank/Back/Support) has distinct visual placement
var FormationLayouts = map[FormationType]map[FormationPosition][]SlotCoordinate{}

// GetInitialSlots returns the initial slot coordinates for a given formation type and position.
// Returns an empty slice if the formation type or position is not found.
func GetInitialSlots(formationType FormationType, position FormationPosition) []SlotCoordinate {
	// Use predefined static slots only.
	slots := getPredefinedSlots(formationType, position)
	if len(slots) == 0 {
		return []SlotCoordinate{}
	}
	// Return a copy to prevent external modification
	result := make([]SlotCoordinate, len(slots))
	copy(result, slots)
	return result
}

// GetNextSlotCoordinate returns the coordinate for the next slot in a formation position.
// existingSlotCount is the number of slots currently occupied (including initial layout).
// For example, if initial layout has 2 slots and you want to add a 3rd, pass existingSlotCount=2.
// Returns false if the position has reached its maximum slot limit.
func GetNextSlotCoordinate(formationType FormationType, position FormationPosition, existingSlotCount int) (SlotCoordinate, bool) {
	// Predefined slots are the only source and the hard cap.
	slots := GetInitialSlots(formationType, position)
	if existingSlotCount < len(slots) {
		return slots[existingSlotCount], true
	}
	return SlotCoordinate{}, false
}

// GetAllSlotsForPosition returns all slot coordinates for a position up to the specified count.
// This includes both initial slots and expanded slots.
// The count is automatically capped at the formation's position limit.
func GetAllSlotsForPosition(formationType FormationType, position FormationPosition, totalSlotCount int) []SlotCoordinate {
	if totalSlotCount <= 0 {
		return []SlotCoordinate{}
	}

	// Predefined slots are the cap
	predefined := GetInitialSlots(formationType, position)
	if totalSlotCount > len(predefined) {
		totalSlotCount = len(predefined)
	}
	result := make([]SlotCoordinate, totalSlotCount)
	copy(result, predefined[:totalSlotCount])
	return result
}

// CalculateFormationBounds returns the min/max X and Y coordinates for a formation layout.
// Useful for centering or scaling the formation display.
func CalculateFormationBounds(formationType FormationType, slotCounts map[FormationPosition]int) (minX, maxX, minY, maxY float64) {
	minX, minY = math.MaxFloat64, math.MaxFloat64
	maxX, maxY = -math.MaxFloat64, -math.MaxFloat64

	for position, count := range slotCounts {
		slots := GetAllSlotsForPosition(formationType, position, count)
		for _, slot := range slots {
			if slot.X < minX {
				minX = slot.X
			}
			if slot.X > maxX {
				maxX = slot.X
			}
			if slot.Y < minY {
				minY = slot.Y
			}
			if slot.Y > maxY {
				maxY = slot.Y
			}
		}
	}

	// Handle empty formation
	if minX == math.MaxFloat64 {
		return 0, 0, 0, 0
	}

	return minX, maxX, minY, maxY
}

// GetMaxSlotsForPosition returns the maximum number of slots allowed for a position in a formation.
func GetMaxSlotsForPosition(formationType FormationType, position FormationPosition) int {
	// Use the count of predefined slots as the authoritative limit.
	return len(getPredefinedSlots(formationType, position))
}

// IsPositionFull checks if a position has reached its maximum slot capacity.
func IsPositionFull(formationType FormationType, position FormationPosition, currentSlotCount int) bool {
	maxSlots := GetMaxSlotsForPosition(formationType, position)
	return currentSlotCount >= maxSlots
}

// GetTotalMaxSlots returns the total maximum slots across all positions for a formation.
func GetTotalMaxSlots(formationType FormationType) int {
    return GetMaxSlotsForPosition(formationType, PositionFront) +
        GetMaxSlotsForPosition(formationType, PositionFlank) +
        GetMaxSlotsForPosition(formationType, PositionBack) +
        GetMaxSlotsForPosition(formationType, PositionSupport)
}

// GenerateFormationLayoutSnapshot creates a complete snapshot of a formation's visual layout.
// This is the recommended function to use when sending formation data to the frontend.
// Slot counts are automatically capped at position limits.
func GenerateFormationLayoutSnapshot(formationType FormationType, slotCounts map[FormationPosition]int) FormationLayoutSnapshot {
	snapshot := FormationLayoutSnapshot{
		FormationType: formationType,
		Positions:     make(map[FormationPosition][]SlotPosition),
	}

	// Generate slot positions for each position type
	for position, count := range slotCounts {
		if count <= 0 {
			continue
		}

		// Cap count at position limit
		maxSlots := GetMaxSlotsForPosition(formationType, position)
		if count > maxSlots && maxSlots > 0 {
			count = maxSlots
		}

		initialSlots := GetInitialSlots(formationType, position)
		initialCount := len(initialSlots)
		slotPositions := make([]SlotPosition, 0, count)

		for i := 0; i < count; i++ {
			coord, ok := GetNextSlotCoordinate(formationType, position, i)
			if !ok {
				// Position is full, stop adding slots
				break
			}

			slotPositions = append(slotPositions, SlotPosition{
				Coordinate: coord,
				SlotIndex:  i,
				IsInitial:  i < initialCount,
			})
		}

		snapshot.Positions[position] = slotPositions
	}

	// Calculate bounds
	minX, maxX, minY, maxY := CalculateFormationBounds(formationType, slotCounts)
	snapshot.Bounds = LayoutBounds{
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
	}

	return snapshot
}
