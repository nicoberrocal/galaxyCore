# Virtual HP Bucket Splitting - Visual Guide

## The Concept

**Virtual splitting** allows users to spread ships visually across multiple slots **without actually splitting the HP bucket**.

## How It Works

### Before Split

```text
HP Bucket 0 (10 Fighters, 1000 HP)
         ↓
    Assignment 0
    ┌─────────────────┐
    │ BucketIndex: 0  │
    │ Count: 10       │
    │ AssignedHP: 1000│
    │ SlotIndex: 0    │
    └─────────────────┘
         ↓
    Visual Slot 0
    [F][F][F][F][F][F][F][F][F][F]
```

### After Virtual Split (4 ships to slot 2)

```text
HP Bucket 0 (10 Fighters, 1000 HP) ← Still ONE bucket!
         ↓                    ↓
    Assignment 0         Assignment 1
    ┌─────────────────┐  ┌─────────────────┐
    │ BucketIndex: 0  │  │ BucketIndex: 0  │ ← Same bucket!
    │ Count: 6        │  │ Count: 4        │
    │ AssignedHP: 600 │  │ AssignedHP: 400 │
    │ SlotIndex: 0    │  │ SlotIndex: 2    │ ← Different slot!
    └─────────────────┘  └─────────────────┘
         ↓                        ↓
    Visual Slot 0           Visual Slot 2
    [F][F][F][F][F][F]      [F][F][F][F]
```

**Key Point**: Both assignments point to `BucketIndex: 0` - they're from the **same underlying HP bucket**, just displayed in different visual positions!

## Visual Formation Example

### Phalanx Formation - Before Split

```text
Front Line (Slot 0, 1, 2):
    [F][F][F][F][F][F][F][F][F][F]
    ↑
    All 10 fighters in slot 0
```

### Phalanx Formation - After Split

```text
Front Line (Slots 0, 1, 2):
    [F][F][F][F][F][F]  [ ]  [F][F][F][F]
    ↑                         ↑
    Slot 0 (6 ships)         Slot 2 (4 ships)
    BucketIndex: 0           BucketIndex: 0 (same!)
```

## MongoDB Representation

### Single Assignment (Before Split)

```json
{
  "slotAssignments": [
    {
      "position": "front",
      "shipType": "fighter",
      "bucketIndex": 0,
      "count": 10,
      "assignedHP": 1000,
      "slotIndex": 0,
      "isManuallyPlaced": false
    }
  ]
}
```

### Virtual Split (After Split)

```json
{
  "slotAssignments": [
    {
      "position": "front",
      "shipType": "fighter",
      "bucketIndex": 0,        // ← Same bucket
      "count": 6,
      "assignedHP": 600,
      "slotIndex": 0,
      "isManuallyPlaced": true
    },
    {
      "position": "front",
      "shipType": "fighter",
      "bucketIndex": 0,        // ← Same bucket!
      "count": 4,
      "assignedHP": 400,
      "slotIndex": 2,          // ← Different slot
      "isManuallyPlaced": true
    }
  ]
}
```

## Combat System Integration

### When Converting to Combat Format

```go
fws := FormationWithSlots{
    SlotAssignments: []FormationSlotAssignment{
        {BucketIndex: 0, Count: 6, SlotIndex: 0},  // Visual slot 0
        {BucketIndex: 0, Count: 4, SlotIndex: 2},  // Visual slot 2
    },
}

combatFormation := fws.ToFormation()
// Returns Formation with 2 assignments:
// - Assignment 0: BucketIndex 0, Count 6
// - Assignment 1: BucketIndex 0, Count 4
```

### Combat Damage Application

When taking damage:

```text
Incoming Damage: 300 HP
         ↓
Distributed to Front Position
         ↓
    Assignment 0        Assignment 1
    (600 HP)           (400 HP)
         ↓                  ↓
    Takes 180 HP       Takes 120 HP
    (proportional)     (proportional)
         ↓                  ↓
    Now 420 HP         Now 280 HP
         ↓                  ↓
    Both still reference BucketIndex 0
```

## Why This Works

### Performance Benefits

✅ **No bucket duplication**: HP buckets remain intact
✅ **Minimal overhead**: Only adds SlotIndex + IsManuallyPlaced
✅ **Efficient combat**: Standard combat calculations work unchanged

### User Experience Benefits

✅ **Visual control**: Users can arrange ships as they want
✅ **Tactical clarity**: Formation arrangement is immediately evident
✅ **Reversible**: Can merge assignments back together

### Technical Benefits

✅ **MongoDB efficient**: Small document size increase
✅ **Backward compatible**: Converts to standard Formation
✅ **Frontend ready**: Coordinates calculated automatically

## Code Example

### Splitting

```go
// Start with 10 fighters in one slot
fws := FormationWithSlots{
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

// User wants to split 4 fighters to slot 2
err := fws.SplitAssignmentToSlot(0, 4, 2)

// Result: 2 assignments, same bucket
// Assignment 0: BucketIndex 0, Count 6, SlotIndex 0
// Assignment 1: BucketIndex 0, Count 4, SlotIndex 2
```

### Merging Back

```go
// Merge them back together
err := fws.MergeAssignments(0, 1)

// Result: 1 assignment
// Assignment 0: BucketIndex 0, Count 10, SlotIndex 0
```

## Visual Comparison

### Traditional Bucket Split (NOT USED)

```text
❌ Actual bucket split (performance cost):

HP Bucket 0 (6 Fighters)    HP Bucket 1 (4 Fighters)
         ↓                           ↓
    Assignment 0                Assignment 1
    BucketIndex: 0              BucketIndex: 1
    Count: 6                    Count: 4

Problem: Creates new HP bucket, increases memory/processing
```

### Virtual Split (OUR APPROACH)

```text
✅ Virtual split (no performance cost):

HP Bucket 0 (10 Fighters) ← Single bucket
         ↓           ↓
    Assignment 0  Assignment 1
    BucketIndex: 0  BucketIndex: 0  ← Same bucket!
    Count: 6        Count: 4
    SlotIndex: 0    SlotIndex: 2    ← Different visual slots

Benefit: No new bucket, just different visual representation
```

## Real-World Scenario

### User Story

1. **User has**: 10 Fighters in Phalanx front line
2. **User wants**: Spread them across the wide front
3. **User action**: Drag 4 fighters to slot 2
4. **System does**: Virtual split (same bucket, different slot)
5. **Result**: Ships appear spread out, but remain in same HP bucket

### Frontend Display

```text
Before:
    [F][F][F][F][F][F][F][F][F][F]  [ ]  [ ]
    ↑
    All in slot 0

After drag:
    [F][F][F][F][F][F]  [ ]  [F][F][F][F]
    ↑                         ↑
    Slot 0                   Slot 2
    (Same bucket, different visual positions)
```

### MongoDB Storage

```json
// Efficient storage - only 2 assignments
{
  "slotAssignments": [
    {"bucketIndex": 0, "count": 6, "slotIndex": 0},
    {"bucketIndex": 0, "count": 4, "slotIndex": 2}
  ]
}
```

## Summary

**Virtual HP Bucket Splitting** allows:

✅ **Visual customization** without performance cost
✅ **Same bucket reference** across multiple visual slots
✅ **Proportional damage** distribution in combat
✅ **Reversible operations** (split/merge)
✅ **MongoDB efficient** storage
✅ **Combat compatible** via ToFormation()

The key insight: **SlotIndex is visual only, BucketIndex is the source of truth for HP**.
