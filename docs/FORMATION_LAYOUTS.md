# Formation Layout System

## Overview

The formation layout system provides visual coordinates for rendering ship formations in the browser-based MMORTS. While ships are **stacked for performance**, the system allows displaying **one ship per formation position** to make the formation arrangement immediately evident.

## Key Concepts

### Coordinate System
- **Origin**: (0, 0) at the center
- **X-axis**: Positive = right, Negative = left
- **Y-axis**: Positive = forward/front, Negative = back/rear
- **Units**: Abstract (scale in frontend based on ship size)

### Formation Positions
Every formation has 4 position types:
- **Front**: Primary damage absorption (60% damage in frontal attacks)
- **Flank**: Mobile strike forces (20% damage in frontal attacks)
- **Back**: Ranged/support units (10% damage in frontal attacks)
- **Support**: Utility/healer positions (10% damage in frontal attacks)

### Layout Components

1. **Initial Layout** (`FormationLayouts`): Base coordinates when formation is first created
2. **Expansion Rules** (`FormationExpansionRules`): Functions that calculate new slot positions as ships are added
3. **Helper Functions**: Utilities to work with the layout system

## Formation Types

### 1. LINE Formation
**Tactical Role**: Balanced front-back arrangement
- **Strengths**: Strong vs frontal attacks
- **Weaknesses**: Vulnerable to flanking
- **Visual**: Single-file columns on front/back axes, flanks spread horizontally
- **Expansion**: Front/back extend linearly, flanks alternate left-right

```
    Front
      |
Flank-+-Flank
      |
    Back
```

### 2. BOX Formation
**Tactical Role**: Defensive all-around protection
- **Strengths**: Even damage distribution, siege resistant
- **Weaknesses**: Slower speed (0.75x)
- **Visual**: Ships form a box perimeter
- **Expansion**: Perimeter expands outward, support fills center

```
Front Front
|         |
Flank   Flank
|         |
Back  Back
```

### 3. VANGUARD Formation
**Tactical Role**: Aggressive forward deployment
- **Strengths**: Fast reconfiguration (60s), high front damage (+25%)
- **Weaknesses**: Weak support/back
- **Visual**: Sharp spearhead pointing forward
- **Expansion**: Reinforces tip, widens V-shape

```
    Front
   /   \
Flank Flank
  |     |
 Back Back
```

### 4. SKIRMISH Formation
**Tactical Role**: Mobile flanking focus
- **Strengths**: Fastest speed (1.2x), excellent hit-and-run
- **Weaknesses**: Minimal front presence
- **Visual**: Very wide flanks, small front
- **Expansion**: Flanks extend extremely wide

```
Flank ---- Flank
  \         /
   Front Front
      |
     Back
```

### 5. ECHELON Formation
**Tactical Role**: Diagonal staggered lines
- **Strengths**: Asymmetric defense, good vs concentrated attacks
- **Weaknesses**: One flank more exposed
- **Visual**: Diagonal line from front-right to back-left
- **Expansion**: Extends along diagonal

```
        Front
       /
    Flank
   /
Flank
 /
Back
```

### 6. PHALANX Formation
**Tactical Role**: Heavy frontal concentration
- **Strengths**: Massive front line (+15% damage, +2 shields), long range back
- **Weaknesses**: Extreme flank vulnerability, slow (0.8x)
- **Visual**: Very wide front line, flanks far out
- **Expansion**: Front gets progressively wider

```
Front Front Front
  |     |     |
Flank       Flank
      |
    Back
```

### 7. SWARM Formation
**Tactical Role**: Dispersed anti-AoE
- **Strengths**: Reduces splash damage, dispersed positioning
- **Weaknesses**: Less concentrated firepower
- **Visual**: Hexagonal dispersion pattern
- **Expansion**: Hexagonal rings (6 ships per ring)

```
    Front
   /    \
Flank  Flank
  |      |
Back    Back
   \    /
   Support
```

## API Reference

### Core Functions

#### `GetInitialSlots(formationType, position) []SlotCoordinate`
Returns the initial slot coordinates for a formation position.

```go
slots := GetInitialSlots(FormationPhalanx, PositionFront)
// Returns: [{-1.5, 1.5}, {0, 1.5}, {1.5, 1.5}]
```

#### `GetNextSlotCoordinate(formationType, position, existingSlotCount) (SlotCoordinate, bool)`
Returns the coordinate for the next slot when expanding.

```go
// Phalanx has 3 initial front slots, add a 4th
coord, ok := GetNextSlotCoordinate(FormationPhalanx, PositionFront, 3)
// Returns: {-2.4, 1.5}, true
```

#### `GetAllSlotsForPosition(formationType, position, totalSlotCount) []SlotCoordinate`
Returns all slots (initial + expanded) up to the specified count.

```go
slots := GetAllSlotsForPosition(FormationVanguard, PositionFront, 5)
// Returns 5 coordinates: 1 initial + 4 expanded
```

#### `GenerateFormationLayoutSnapshot(formationType, slotCounts) FormationLayoutSnapshot`
**Recommended for frontend**: Generates complete formation layout with metadata.

```go
slotCounts := map[FormationPosition]int{
    PositionFront:   5,
    PositionFlank:   4,
    PositionBack:    2,
    PositionSupport: 3,
}
snapshot := GenerateFormationLayoutSnapshot(FormationVanguard, slotCounts)
// Returns complete layout with bounds and metadata
```

#### `CalculateFormationBounds(formationType, slotCounts) (minX, maxX, minY, maxY)`
Calculates bounding box for centering/scaling.

```go
minX, maxX, minY, maxY := CalculateFormationBounds(FormationPhalanx, slotCounts)
width := maxX - minX
height := maxY - minY
```

## Frontend Integration

### Recommended Workflow

1. **Get Formation Data**: Call `GenerateFormationLayoutSnapshot()` with ship counts per position
2. **Parse JSON**: The snapshot is JSON-serializable
3. **Render Ships**: For each position, render one ship sprite per slot coordinate
4. **Scale Coordinates**: Multiply coordinates by your ship sprite size
5. **Center Formation**: Use `Bounds` to center the formation in viewport

### Example JSON Output

```json
{
  "formationType": "vanguard",
  "positions": {
    "front": [
      {
        "coordinate": {"x": 0, "y": 2},
        "slotIndex": 0,
        "isInitial": true
      },
      {
        "coordinate": {"x": -0.3, "y": 2.6},
        "slotIndex": 1,
        "isInitial": false
      }
    ],
    "flank": [...]
  },
  "bounds": {
    "minX": -1.2,
    "maxX": 1.2,
    "minY": -1.0,
    "maxY": 2.6
  }
}
```

### Rendering Example (Pseudocode)

```javascript
const snapshot = getFormationSnapshot(); // From backend API
const shipSpriteSize = 32; // pixels
const scale = 40; // pixels per coordinate unit

// Center the formation
const centerX = canvas.width / 2;
const centerY = canvas.height / 2;

// Render each position
for (const [position, slots] of Object.entries(snapshot.positions)) {
  const shipSprite = getShipSpriteForPosition(position);
  
  for (const slot of slots) {
    const x = centerX + (slot.coordinate.x * scale);
    const y = centerY - (slot.coordinate.y * scale); // Flip Y for screen coords
    
    // Render ship
    drawSprite(shipSprite, x, y);
    
    // Optional: highlight initial vs expanded slots
    if (!slot.isInitial) {
      drawExpansionIndicator(x, y);
    }
  }
}
```

## Design Principles

### 1. Tactical Accuracy
Each formation's visual layout reflects its tactical characteristics:
- **Phalanx**: Wide front shows "frontal_fortress" property
- **Skirmish**: Wide flanks show "mobile" and "hit_and_run" properties
- **Swarm**: Hexagonal dispersion shows "anti_aoe" property

### 2. Damage Distribution Alignment
Visual density correlates with directional damage weights:
- **Frontal attacks**: 60% front, 20% flank, 10% back, 10% support
- Formations with strong fronts (Phalanx, Vanguard) have more front slots
- Formations with strong flanks (Skirmish) have more flank slots

### 3. Counter Matrix Representation
Visual arrangements support formation counter relationships:
- **Phalanx** (wide front) counters **Vanguard** (concentrated spear) = 1.25x damage
- **Skirmish** (wide flanks) counters **Phalanx** (exposed flanks) = 1.3x damage
- **Box** (even perimeter) counters **Vanguard** (single point) = 1.3x damage

### 4. Expansion Logic
Each formation expands in a way that maintains its tactical character:
- **Line**: Extends linearly along axes
- **Box**: Expands perimeter outward
- **Vanguard**: Reinforces spearhead
- **Phalanx**: Widens front line
- **Swarm**: Adds hexagonal rings

## Performance Considerations

### Backend
- Ships remain **stacked** in HP buckets for performance
- Layout system is purely visual/cosmetic
- No impact on combat calculations

### Frontend
- Display **one ship sprite per formation position slot**
- Ships move as a **unit** (same velocity vector)
- Use sprite batching for rendering efficiency
- Consider LOD (Level of Detail) for distant formations

## Integration with Combat System

The layout system integrates with:

1. **Formation Counter Matrix**: Visual arrangement reflects counter relationships
2. **Directional Damage Weights**: Position density matches damage distribution
3. **Position Bonuses**: Visual placement shows which ships get which bonuses
4. **Special Properties**: Layout reflects properties like "dispersed", "frontal_fortress"

## Future Enhancements

Potential additions:
- **Rotation**: Mirror formations (e.g., left vs right echelon)
- **Facing**: Rotate entire formation based on `Formation.Facing` field
- **Animation**: Smooth transitions when reconfiguring formations
- **Collision Detection**: Prevent formation overlap in tight spaces
- **Formation Density**: Adjust spacing based on ship size/type
