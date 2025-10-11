package ships

import "math"

// SlotCoordinate defines a 2D position for a ship slot within a formation's visual layout.
// The coordinate system is centered at (0,0).
// X represents the horizontal axis (positive is right, negative is left).
// Y represents the vertical axis (positive is forward/front, negative is back/rear).
// Units are abstract and should be scaled by the frontend based on ship size.
type SlotCoordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// FormationLayouts defines the initial visual slot coordinates for each position within a formation.
// These are the base positions shown when a formation is first created.
// Design principles:
// - Y-axis: positive=front, negative=back (matches directional damage weights)
// - Visual density reflects tactical role (e.g., Phalanx has wide front, Vanguard has concentrated spearhead)
// - Layouts support the formation counter matrix relationships
// - Each position (Front/Flank/Back/Support) has distinct visual placement
var FormationLayouts = map[FormationType]map[FormationPosition][]SlotCoordinate{
	// LINE: Balanced front-back arrangement, strong vs frontal, weak to flanking
	// Front gets primary position, flanks are exposed
	FormationLine: {
		PositionFront:   {{X: 0, Y: 1.5}},
		PositionFlank:   {{X: -1.5, Y: 0}, {X: 1.5, Y: 0}},
		PositionBack:    {{X: 0, Y: -1.5}},
		PositionSupport: {{X: 0, Y: 0}},
	},
	// BOX: Defensive all-around protection, even distribution
	// All positions get equal representation forming a box perimeter
	FormationBox: {
		PositionFront:   {{X: -1, Y: 1.5}, {X: 1, Y: 1.5}},
		PositionFlank:   {{X: -2, Y: 0}, {X: 2, Y: 0}},
		PositionBack:    {{X: -1, Y: -1.5}, {X: 1, Y: -1.5}},
		PositionSupport: {{X: 0, Y: 0}},
	},
	// VANGUARD: Aggressive spearhead, concentrated front, fast reconfiguration
	// Front forms a sharp point, flanks create the 'V' shape
	FormationVanguard: {
		PositionFront:   {{X: 0, Y: 2}},
		PositionFlank:   {{X: -1.2, Y: 0.8}, {X: 1.2, Y: 0.8}},
		PositionBack:    {{X: -0.8, Y: -1}, {X: 0.8, Y: -1}},
		PositionSupport: {{X: 0, Y: -0.3}},
	},
	// SKIRMISH: Mobile flanking focus, wide flanks for hit-and-run
	// Flanks are emphasized and spread wide, front is minimal
	FormationSkirmish: {
		PositionFront:   {{X: -0.6, Y: 0.8}, {X: 0.6, Y: 0.8}},
		PositionFlank:   {{X: -2.5, Y: 0.2}, {X: 2.5, Y: 0.2}},
		PositionBack:    {{X: 0, Y: -1.2}},
		PositionSupport: {{X: 0, Y: -0.2}},
	},
	// ECHELON: Diagonal staggered lines, asymmetric
	// Creates a diagonal line from front-right to back-left (can be mirrored)
	FormationEchelon: {
		PositionFront:   {{X: 1.5, Y: 1.5}},
		PositionFlank:   {{X: 0.5, Y: 0.5}, {X: -0.5, Y: -0.5}},
		PositionBack:    {{X: -1.5, Y: -1.5}},
		PositionSupport: {{X: -0.8, Y: -1}},
	},
	// PHALANX: Heavy frontal concentration, very wide front line
	// Front is massively wide, flanks are far out, back is minimal
	FormationPhalanx: {
		PositionFront:   {{X: -1.5, Y: 1.5}, {X: 0, Y: 1.5}, {X: 1.5, Y: 1.5}},
		PositionFlank:   {{X: -3, Y: 0.8}, {X: 3, Y: 0.8}},
		PositionBack:    {{X: 0, Y: -1.2}},
		PositionSupport: {{X: -0.7, Y: 0.2}, {X: 0.7, Y: 0.2}},
	},
	// SWARM: Dispersed anti-AoE, hexagonal spread pattern
	// Ships are spread out in a loose hexagonal pattern for maximum dispersion
	FormationSwarm: {
		PositionFront:   {{X: 0, Y: 1.5}},
		PositionFlank:   {{X: -1.3, Y: 0.75}, {X: 1.3, Y: 0.75}},
		PositionBack:    {{X: -1.3, Y: -0.75}, {X: 1.3, Y: -0.75}},
		PositionSupport: {{X: 0, Y: -1.5}},
	},
}

// ExpansionFunc defines a function that calculates the next slot coordinate based on existing slots.
// It takes the current slot count (number of existing slots) and returns the coordinate for the new slot.
// The slot count starts at 0 for the first expansion beyond the initial layout.
type ExpansionFunc func(slotCount int) SlotCoordinate

// FormationExpansionRules defines how formations grow when new ships are assigned to positions.
// Each position in each formation has a specific expansion pattern that maintains the formation's tactical character.
// The expansion functions receive the count of EXISTING slots beyond the initial layout.
// For example, if initial layout has 2 front slots and you add a 3rd ship, c=0 (first expansion).
var FormationExpansionRules = map[FormationType]map[FormationPosition]ExpansionFunc{
	// ===== LINE FORMATION =====
	// Expands linearly along axes, maintaining the front-back line structure
	FormationLine: {
		PositionFront: func(c int) SlotCoordinate {
			// Expands forward in a column
			return SlotCoordinate{X: 0, Y: 1.5 + 0.8*float64(c+1)}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Alternates left-right, expanding outward
			side := c / 2 // 0,0,1,1,2,2...
			dist := 1.5 + float64(side)*0.8
			if c%2 == 0 {
				return SlotCoordinate{X: -dist, Y: 0}
			}
			return SlotCoordinate{X: dist, Y: 0}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Expands backward in a column
			return SlotCoordinate{X: 0, Y: -1.5 - 0.8*float64(c+1)}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Fills in behind the front line, alternating left-right
			offset := 0.6 * float64(c+1)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: 0.3}
			}
			return SlotCoordinate{X: offset, Y: 0.3}
		},
	},

	// ===== BOX FORMATION =====
	// Expands the perimeter outward, maintaining defensive box shape
	FormationBox: {
		PositionFront: func(c int) SlotCoordinate {
			// First fills center, then expands outward
			if c == 0 {
				return SlotCoordinate{X: 0, Y: 1.5}
			}
			side := (c - 1) / 2
			dist := 1.0 + float64(side)*0.8
			if c%2 == 1 {
				return SlotCoordinate{X: -dist, Y: 1.5}
			}
			return SlotCoordinate{X: dist, Y: 1.5}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Expands flanks outward, alternating sides
			side := c / 2
			dist := 2.0 + float64(side)*0.8
			if c%2 == 0 {
				return SlotCoordinate{X: -dist, Y: 0}
			}
			return SlotCoordinate{X: dist, Y: 0}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Mirrors front expansion
			if c == 0 {
				return SlotCoordinate{X: 0, Y: -1.5}
			}
			side := (c - 1) / 2
			dist := 1.0 + float64(side)*0.8
			if c%2 == 1 {
				return SlotCoordinate{X: -dist, Y: -1.5}
			}
			return SlotCoordinate{X: dist, Y: -1.5}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Fills center in a circular pattern
			radius := 0.5 + float64(c/4)*0.4
			angle := float64(c%4) * 90.0 * (math.Pi / 180.0)
			return SlotCoordinate{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
		},
	},

	// ===== VANGUARD FORMATION =====
	// Expands the spearhead forward and widens the V-shape
	FormationVanguard: {
		PositionFront: func(c int) SlotCoordinate {
			// Reinforces the spear tip, alternating slightly left-right
			offset := 0.3 * float64((c+1)/2)
			yPos := 2.0 + 0.6*float64(c)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: yPos}
			}
			return SlotCoordinate{X: offset, Y: yPos}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Extends the V-shape outward and slightly back
			side := c / 2
			xDist := 1.2 + float64(side)*0.7
			yPos := 0.8 - float64(side)*0.3
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: yPos}
			}
			return SlotCoordinate{X: xDist, Y: yPos}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Fills the rear, alternating sides
			side := c / 2
			xDist := 0.8 + float64(side)*0.6
			yPos := -1.0 - float64(side)*0.3
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: yPos}
			}
			return SlotCoordinate{X: xDist, Y: yPos}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Stacks behind the front, slightly offset
			offset := 0.5 * float64((c+1)/2)
			yPos := -0.3 - 0.5*float64(c)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: yPos}
			}
			return SlotCoordinate{X: offset, Y: yPos}
		},
	},

	// ===== SKIRMISH FORMATION =====
	// Emphasizes wide flanks and mobility
	FormationSkirmish: {
		PositionFront: func(c int) SlotCoordinate {
			// Minimal front, fills in the small gap
			side := c / 2
			xDist := 0.6 + float64(side)*0.5
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: 0.8}
			}
			return SlotCoordinate{X: xDist, Y: 0.8}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Expands flanks very wide for mobility
			side := c / 2
			xDist := 2.5 + float64(side)*0.9
			// Slight forward-back stagger for depth
			yOffset := 0.2 - float64(side)*0.1
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: yOffset}
			}
			return SlotCoordinate{X: xDist, Y: yOffset}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Expands backward in a line
			offset := 0.6 * float64((c+1)/2)
			yPos := -1.2 - 0.5*float64(c)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: yPos}
			}
			return SlotCoordinate{X: offset, Y: yPos}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Clusters near center for quick support
			offset := 0.5 * float64((c+1)/2)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: -0.2}
			}
			return SlotCoordinate{X: offset, Y: -0.2}
		},
	},

	// ===== ECHELON FORMATION =====
	// Maintains diagonal staggered pattern
	FormationEchelon: {
		PositionFront: func(c int) SlotCoordinate {
			// Extends the diagonal forward-right
			step := float64(c + 1)
			return SlotCoordinate{X: 1.5 + step*0.8, Y: 1.5 + step*0.8}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Fills the diagonal, alternating between the two flank positions
			if c%2 == 0 {
				// Upper-right flank
				step := float64(c/2 + 1)
				return SlotCoordinate{X: 0.5 + step*0.7, Y: 0.5 + step*0.7}
			}
			// Lower-left flank
			step := float64(c/2 + 1)
			return SlotCoordinate{X: -0.5 - step*0.7, Y: -0.5 - step*0.7}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Extends the diagonal backward-left
			step := float64(c + 1)
			return SlotCoordinate{X: -1.5 - step*0.8, Y: -1.5 - step*0.8}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Fills along the back-left area
			step := float64(c + 1)
			return SlotCoordinate{X: -0.8 - step*0.5, Y: -1.0 - step*0.4}
		},
	},

	// ===== PHALANX FORMATION =====
	// Massively wide front line, vulnerable flanks
	FormationPhalanx: {
		PositionFront: func(c int) SlotCoordinate {
			// Expands the front line wider and wider
			side := c / 2
			xDist := 1.5 + float64(side+1)*0.9
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: 1.5}
			}
			return SlotCoordinate{X: xDist, Y: 1.5}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Flanks extend very far out
			side := c / 2
			xDist := 3.0 + float64(side)*1.0
			// Slightly forward to protect the wide front
			yPos := 0.8 - float64(side)*0.2
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: yPos}
			}
			return SlotCoordinate{X: xDist, Y: yPos}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Minimal back, fills in a tight column
			offset := 0.5 * float64((c+1)/2)
			yPos := -1.2 - 0.6*float64(c)
			if c%2 == 0 {
				return SlotCoordinate{X: -offset, Y: yPos}
			}
			return SlotCoordinate{X: offset, Y: yPos}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Fills behind the front line, alternating
			side := c / 2
			xDist := 0.7 + float64(side)*0.6
			if c%2 == 0 {
				return SlotCoordinate{X: -xDist, Y: 0.2}
			}
			return SlotCoordinate{X: xDist, Y: 0.2}
		},
	},

	// ===== SWARM FORMATION =====
	// Hexagonal dispersion pattern, anti-AoE
	FormationSwarm: {
		PositionFront: func(c int) SlotCoordinate {
			// Expands in a hexagonal ring pattern
			ring := c / 6
			radius := 1.5 + float64(ring)*1.2
			angle := float64(c%6)*60.0*(math.Pi/180.0) + math.Pi/2
			return SlotCoordinate{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
		},
		PositionFlank: func(c int) SlotCoordinate {
			// Offset hexagonal pattern for flanks
			ring := c / 6
			radius := 1.3 + float64(ring)*1.1
			angle := float64(c%6)*60.0*(math.Pi/180.0) + math.Pi/6 // 30° offset
			return SlotCoordinate{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
		},
		PositionBack: func(c int) SlotCoordinate {
			// Another offset hexagonal pattern
			ring := c / 6
			radius := 1.3 + float64(ring)*1.1
			angle := float64(c%6)*60.0*(math.Pi/180.0) - math.Pi/2 // Bottom-oriented
			return SlotCoordinate{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
		},
		PositionSupport: func(c int) SlotCoordinate {
			// Fills gaps in the hexagonal pattern
			ring := c / 6
			radius := 0.8 + float64(ring)*1.0
			angle := float64(c%6)*60.0*(math.Pi/180.0) + math.Pi // 180° offset
			return SlotCoordinate{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
		},
	},
}

// GetInitialSlots returns the initial slot coordinates for a given formation type and position.
// Returns an empty slice if the formation type or position is not found.
func GetInitialSlots(formationType FormationType, position FormationPosition) []SlotCoordinate {
	if formationLayout, ok := FormationLayouts[formationType]; ok {
		if slots, ok := formationLayout[position]; ok {
			// Return a copy to prevent external modification
			result := make([]SlotCoordinate, len(slots))
			copy(result, slots)
			return result
		}
	}
	return []SlotCoordinate{}
}

// GetNextSlotCoordinate returns the coordinate for the next slot in a formation position.
// existingSlotCount is the number of slots currently occupied (including initial layout).
// For example, if initial layout has 2 slots and you want to add a 3rd, pass existingSlotCount=2.
// Returns false if the position has reached its maximum slot limit.
func GetNextSlotCoordinate(formationType FormationType, position FormationPosition, existingSlotCount int) (SlotCoordinate, bool) {
	// Check if we've reached the slot limit for this position
	if limits, ok := FormationSlotLimits[formationType]; ok {
		maxSlots := GetPositionLimit(limits, position)
		if existingSlotCount >= maxSlots {
			return SlotCoordinate{}, false // Position is full
		}
	}

	// Get the initial slot count for this position
	initialSlots := GetInitialSlots(formationType, position)
	initialCount := len(initialSlots)

	// If we haven't filled the initial slots yet, return the next initial slot
	if existingSlotCount < initialCount {
		return initialSlots[existingSlotCount], true
	}

	// Otherwise, use the expansion function
	if formationRules, ok := FormationExpansionRules[formationType]; ok {
		if expansionFunc, ok := formationRules[position]; ok {
			// Calculate the expansion index (0-based, starting after initial slots)
			expansionIndex := existingSlotCount - initialCount
			return expansionFunc(expansionIndex), true
		}
	}

	// No expansion rule found
	return SlotCoordinate{}, false
}

// GetAllSlotsForPosition returns all slot coordinates for a position up to the specified count.
// This includes both initial slots and expanded slots.
// The count is automatically capped at the formation's position limit.
func GetAllSlotsForPosition(formationType FormationType, position FormationPosition, totalSlotCount int) []SlotCoordinate {
	if totalSlotCount <= 0 {
		return []SlotCoordinate{}
	}

	// Cap at the position limit
	if limits, ok := FormationSlotLimits[formationType]; ok {
		maxSlots := GetPositionLimit(limits, position)
		if totalSlotCount > maxSlots {
			totalSlotCount = maxSlots
		}
	}

	slots := make([]SlotCoordinate, 0, totalSlotCount)

	// Add initial slots
	initialSlots := GetInitialSlots(formationType, position)
	for i := 0; i < totalSlotCount && i < len(initialSlots); i++ {
		slots = append(slots, initialSlots[i])
	}

	// If we need more slots, use expansion
	if totalSlotCount > len(initialSlots) {
		if formationRules, ok := FormationExpansionRules[formationType]; ok {
			if expansionFunc, ok := formationRules[position]; ok {
				for i := len(initialSlots); i < totalSlotCount; i++ {
					expansionIndex := i - len(initialSlots)
					slots = append(slots, expansionFunc(expansionIndex))
				}
			}
		}
	}

	return slots
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

// GetPositionLimit returns the maximum slot count for a specific position.
func GetPositionLimit(limits PositionSlotLimits, position FormationPosition) int {
	switch position {
	case PositionFront:
		return limits.Front
	case PositionFlank:
		return limits.Flank
	case PositionBack:
		return limits.Back
	case PositionSupport:
		return limits.Support
	default:
		return 0
	}
}

// GetMaxSlotsForPosition returns the maximum number of slots allowed for a position in a formation.
func GetMaxSlotsForPosition(formationType FormationType, position FormationPosition) int {
	if limits, ok := FormationSlotLimits[formationType]; ok {
		return GetPositionLimit(limits, position)
	}
	return 0
}

// IsPositionFull checks if a position has reached its maximum slot capacity.
func IsPositionFull(formationType FormationType, position FormationPosition, currentSlotCount int) bool {
	maxSlots := GetMaxSlotsForPosition(formationType, position)
	return currentSlotCount >= maxSlots
}

// GetTotalMaxSlots returns the total maximum slots across all positions for a formation.
func GetTotalMaxSlots(formationType FormationType) int {
	if limits, ok := FormationSlotLimits[formationType]; ok {
		return limits.Front + limits.Flank + limits.Back + limits.Support
	}
	return 0
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
