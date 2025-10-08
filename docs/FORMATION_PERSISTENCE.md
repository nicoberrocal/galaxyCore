# Formation Persistence System

## Overview

The formation persistence system allows saving user-arranged formations to MongoDB while maintaining compatibility with the combat system. It extends `FormationAssignment` with visual slot information, enabling users to customize ship placement without affecting combat mechanics.

## Key Concepts

### HP Bucket Splitting (Virtual)

Users can **visually split** HP buckets across multiple slots **without actually splitting the bucket**:

- **Same BucketIndex**: Multiple assignments can reference the same HP bucket
- **Visual Only**: Ships appear in different slots but remain in the same bucket for combat
- **Reversible**: Assignments can be merged back together

### Data Structures

#### FormationSlotAssignment

Extends `FormationAssignment` with visual information:

```go
type FormationSlotAssignment struct {
    FormationAssignment        // Position, ShipType, BucketIndex, Count, AssignedHP
    SlotIndex           int    // Which visual slot (0-based)
    IsManuallyPlaced    bool   // User placed vs auto-assigned
}
```

#### FormationWithSlots

MongoDB-ready structure with visual layout:

```go
type FormationWithSlots struct {
    Type             FormationType
    Facing           string
    SlotAssignments  []FormationSlotAssignment
    Modifiers        FormationMods
    CreatedAt        time.Time
    Version          int
}
```

## MongoDB Schema

### Document Structure

```json
{
  "_id": ObjectId("..."),
  "formationType": "vanguard",
  "facing": "north",
  "slotAssignments": [
    {
      "position": "front",
      "layer": 0,
      "shipType": "fighter",
      "bucketIndex": 0,
      "count": 3,
      "assignedHP": 300,
      "slotIndex": 0,
      "isManuallyPlaced": false
    },
    {
      "position": "front",
      "layer": 0,
      "shipType": "fighter",
      "bucketIndex": 0,
      "count": 2,
      "assignedHP": 200,
      "slotIndex": 2,
      "isManuallyPlaced": true
    }
  ],
  "modifiers": { ... },
  "createdAt": ISODate("..."),
  "version": 1
}
```

**Note**: Both assignments above reference `bucketIndex: 0` - they're from the same HP bucket, just visually split!

### Indexes

Recommended MongoDB indexes:

```javascript
db.formations.createIndex({ "formationType": 1 })
db.formations.createIndex({ "slotAssignments.shipType": 1 })
db.formations.createIndex({ "createdAt": -1 })
```

## API Reference

### Conversion Functions

#### `FromFormation(formation) FormationWithSlots`

Converts standard `Formation` to `FormationWithSlots` with auto-assigned slots.

```go
formation := Formation{ /* ... */ }
fws := FromFormation(formation)
// Auto-assigns slots: 0, 1, 2, ... per position
```

#### `ToFormation() Formation`

Converts `FormationWithSlots` back to standard `Formation` for combat.

```go
combatFormation := fws.ToFormation()
// Strips visual info, ready for combat calculations
```

### Visual Layout Functions

#### `GenerateVisualSnapshot() FormationVisualSnapshot`

**Recommended for frontend** - Generates complete visual layout with coordinates.

```go
snapshot := fws.GenerateVisualSnapshot()
// Returns: assignments with coordinates, bounds, etc.
```

#### `GetSlotCoordinates() map[int]SlotCoordinate`

Returns coordinates for all assignments.

```go
coords := fws.GetSlotCoordinates()
coord := coords[0] // Coordinate for assignment 0
```

### User Arrangement Functions

#### `MoveAssignmentToSlot(assignmentIndex, newSlotIndex) error`

Moves an assignment to a different visual slot.

```go
err := fws.MoveAssignmentToSlot(0, 3)
// Moves assignment 0 to slot 3 in its position
```

#### `SwapAssignmentSlots(index1, index2) error`

Swaps visual positions of two assignments (must be same position).

```go
err := fws.SwapAssignmentSlots(0, 1)
// Swaps slots of assignments 0 and 1
```

#### `SplitAssignmentToSlot(sourceIndex, splitCount, targetSlotIndex) error`

Splits ships from a bucket to a new visual slot **without splitting the bucket**.

```go
err := fws.SplitAssignmentToSlot(0, 3, 2)
// Takes 3 ships from assignment 0, places in slot 2
// Both assignments keep same bucketIndex!
```

#### `MergeAssignments(index1, index2) error`

Merges two assignments from the same bucket.

```go
err := fws.MergeAssignments(0, 1)
// Combines assignments 0 and 1 (must be same bucket)
```

## Usage Examples

### Example 1: Basic Save/Load

```go
// Create formation from combat system
formation := Formation{
    Type: FormationVanguard,
    Assignments: []FormationAssignment{
        {Position: PositionFront, ShipType: Fighter, BucketIndex: 0, Count: 5, AssignedHP: 500},
    },
}

// Convert for MongoDB
fws := FromFormation(formation)

// Save to MongoDB
collection.InsertOne(ctx, fws)

// Load from MongoDB
var loaded FormationWithSlots
collection.FindOne(ctx, filter).Decode(&loaded)

// Use in combat
combatFormation := loaded.ToFormation()
```

### Example 2: User Customization

```go
// User wants to rearrange ships
fws := FromFormation(formation)

// Move fighters to slot 2
fws.MoveAssignmentToSlot(0, 2)

// Split 3 fighters to slot 4
fws.SplitAssignmentToSlot(0, 3, 4)

// Save customized layout
collection.UpdateOne(ctx, filter, bson.M{"$set": fws})
```

### Example 3: Frontend Integration

```go
// API endpoint handler
func GetFormationLayout(w http.ResponseWriter, r *http.Request) {
    var fws FormationWithSlots
    collection.FindOne(ctx, filter).Decode(&fws)
    
    // Generate visual snapshot
    snapshot := fws.GenerateVisualSnapshot()
    
    // Send to frontend
    json.NewEncoder(w).Encode(snapshot)
}
```

Frontend receives:

```json
{
  "formationType": "vanguard",
  "facing": "north",
  "assignments": [
    {
      "position": "front",
      "shipType": "fighter",
      "bucketIndex": 0,
      "count": 3,
      "assignedHP": 300,
      "slotIndex": 0,
      "coordinate": {"x": 0, "y": 2},
      "isManuallyPlaced": false
    },
    {
      "position": "front",
      "shipType": "fighter",
      "bucketIndex": 0,
      "count": 2,
      "assignedHP": 200,
      "slotIndex": 2,
      "coordinate": {"x": 0.3, "y": 2.6},
      "isManuallyPlaced": true
    }
  ],
  "bounds": {
    "minX": -1.2,
    "maxX": 1.2,
    "minY": -1.0,
    "maxY": 2.6
  }
}
```

## HP Bucket Splitting Explained

### The Problem

Users want to spread ships visually, but HP buckets shouldn't be split for performance.

### The Solution

**Virtual splitting** - Multiple assignments reference the same bucket:

```go
// Original: 10 fighters in bucket 0
{Position: PositionFront, BucketIndex: 0, Count: 10, SlotIndex: 0}

// After split: Still bucket 0, but in 2 visual slots
{Position: PositionFront, BucketIndex: 0, Count: 6, SlotIndex: 0}
{Position: PositionFront, BucketIndex: 0, Count: 4, SlotIndex: 2}
```

### How It Works

1. **Same BucketIndex**: Both assignments point to the same HP bucket
2. **Different SlotIndex**: Ships appear in different visual positions
3. **Count Split**: Total count is preserved (6 + 4 = 10)
4. **HP Split**: HP is proportionally divided (600 + 400 = 1000)

### Combat Integration

When converting to `Formation` for combat:

```go
combatFormation := fws.ToFormation()
// Returns 2 assignments, both with bucketIndex: 0
// Combat system treats them as separate targets
// But they reference the same underlying HP bucket
```

## Workflow Diagrams

### Save Flow

```text
Combat System
    ↓
Formation (standard)
    ↓
FromFormation()
    ↓
FormationWithSlots (auto-assigned slots)
    ↓
User customization (move/split/swap)
    ↓
MongoDB save
```

### Load Flow

```text
MongoDB load
    ↓
FormationWithSlots (with user layout)
    ↓
GenerateVisualSnapshot() → Frontend
    ↓
ToFormation() → Combat System
```

## Best Practices

### ✅ Do

- **Use FormationWithSlots** for all user-facing formation data
- **Call GenerateVisualSnapshot()** when sending to frontend
- **Call ToFormation()** when entering combat
- **Validate slot indices** before moving/splitting
- **Mark manual placements** with `IsManuallyPlaced: true`

### ❌ Don't

- **Don't modify Formation.Assignments** directly for visual changes
- **Don't split actual HP buckets** - use virtual splitting
- **Don't forget to update IsManuallyPlaced** when user customizes
- **Don't send FormationWithSlots** to combat calculations

## Error Handling

### Common Errors

```go
// Invalid assignment index
err := fws.MoveAssignmentToSlot(999, 0)
// Returns: ErrInvalidAssignmentIndex

// Invalid slot index
err := fws.MoveAssignmentToSlot(0, 999)
// Returns: ErrInvalidSlotIndex

// Different positions
err := fws.SwapAssignmentSlots(0, 1) // 0=front, 1=flank
// Returns: ErrDifferentPositions

// Invalid split count
err := fws.SplitAssignmentToSlot(0, 0, 1) // count must be > 0
// Returns: ErrInvalidSplitCount

// Cannot merge
err := fws.MergeAssignments(0, 1) // different buckets
// Returns: ErrCannotMergeAssignments
```

### Error Types

```go
var (
    ErrInvalidAssignmentIndex  = &FormationError{...}
    ErrInvalidSlotIndex        = &FormationError{...}
    ErrDifferentPositions      = &FormationError{...}
    ErrInvalidSplitCount       = &FormationError{...}
    ErrCannotMergeAssignments  = &FormationError{...}
)
```

## Performance Considerations

### MongoDB

- **Document size**: ~200-500 bytes per assignment
- **Typical formation**: 5-20 assignments = 1-10 KB
- **Indexes**: Add indexes on frequently queried fields

### Memory

- **No bucket duplication**: Virtual splitting doesn't duplicate HP buckets
- **Minimal overhead**: Only adds SlotIndex + IsManuallyPlaced per assignment
- **Efficient conversion**: ToFormation() is O(n) where n = assignments

### Frontend

- **One API call**: GenerateVisualSnapshot() includes everything
- **Pre-calculated coordinates**: No client-side layout calculation needed
- **Cached bounds**: Bounding box included for centering

## Migration Guide

### Existing Formations

If you have existing `Formation` documents:

```go
// Load old format
var oldFormation Formation
collection.FindOne(ctx, filter).Decode(&oldFormation)

// Convert to new format
newFormat := FromFormation(oldFormation)

// Save new format
collection.ReplaceOne(ctx, filter, newFormat)
```

### Backward Compatibility

`FormationWithSlots` can always convert to `Formation`:

```go
// New system → Old system
combatFormation := fws.ToFormation()
// Works with existing combat code
```

## Testing

### Unit Tests

```go
func TestSplitAssignment(t *testing.T) {
    fws := FormationWithSlots{
        Type: FormationLine,
        SlotAssignments: []FormationSlotAssignment{
            {
                FormationAssignment: FormationAssignment{
                    Position: PositionFront,
                    BucketIndex: 0,
                    Count: 10,
                    AssignedHP: 1000,
                },
                SlotIndex: 0,
            },
        },
    }
    
    err := fws.SplitAssignmentToSlot(0, 4, 1)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(fws.SlotAssignments))
    assert.Equal(t, 6, fws.SlotAssignments[0].Count)
    assert.Equal(t, 4, fws.SlotAssignments[1].Count)
    assert.Equal(t, 0, fws.SlotAssignments[0].BucketIndex)
    assert.Equal(t, 0, fws.SlotAssignments[1].BucketIndex) // Same bucket!
}
```

## Summary

The formation persistence system provides:

✅ **MongoDB-ready** structure with BSON tags
✅ **Virtual HP bucket splitting** for visual arrangement
✅ **User customization** (move, swap, split, merge)
✅ **Combat compatibility** via ToFormation()
✅ **Frontend integration** via GenerateVisualSnapshot()
✅ **Backward compatibility** with existing Formation system

All while maintaining the performance benefits of stacked ships!
