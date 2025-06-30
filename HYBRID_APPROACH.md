# Hybrid Data Approach for Galaxy Core

This document explains the hybrid data modeling approach used to handle different fleet states and game mechanics.

## Overview

The game has two distinct fleet behaviors that require different data storage patterns:

1. **System Defense**: Fleets that colonize systems become "embedded" in the system
2. **Free Movement**: Fleets traveling in space or mining resources remain as separate documents

## Data Flow

### Stack Enters System (Colonization)
```
1. Stack approaches system (within CollisionRadius)
2. Stack document deleted from `stacks` collection
3. Stack data embedded in System.DefendingFleet
4. System.Colonization updated (IsColonized = true, ColonizedBy = player)
5. PlayerGameState.ColonizedSystems updated
```

### Fleet Merging (Same/Allied Player)
```
1. Incoming stack ships merged into System.DefendingFleet.Ships
2. Allied fleet info added to DefendingFleet.AlliedFleets (if ally)
3. Incoming stack document deleted
4. System ownership remains with original colonizer
```

### Combat Resolution (Enemy Attack)
```
1. All combat data available in single System document
2. Ranged combat reduces DefendingFleet.Ships
3. If fleet completely destroyed:
   - System becomes "orphaned" (DefendingFleet = nil, IsColonized = false)
   - Enemy does NOT automatically colonize
   - Enemy must physically enter system to claim it
4. If allies survive: Original owner keeps system
```

### Stack Leaves System
```
1. Create new ShipStack document from DefendingFleet data
2. Clear System.DefendingFleet
3. Update System.Colonization (IsColonized = false)
4. Update PlayerGameState references
```

## Collections Structure

### Systems (`systems`)
- Embeds defending fleet when colonized
- Atomic operations for combat/colonization
- Single source of truth for territorial control
- **Consistency Rule**: DefendingFleet exists ↔ IsColonized = true

### Stacks (`stacks`)
- Only contains fleets NOT defending systems
- Used for movement, mining, and traveling
- References mining targets (asteroids/nebulas)

### Player Game State (`player_game_states`)
- Tracks resources and assets per player per map
- Denormalized references for performance
- **Mutually Exclusive Arrays**: ColonizedSystems vs ActiveStacks

### Asteroids/Nebulas (`asteroids`, `nebulas`)
- References mining stacks (not embedded)
- Multiple stacks can mine simultaneously
- Tracks extraction rates and remaining resources

## System State Consistency Rules

### **Rule 1: Colonization ↔ DefendingFleet**
```go
// ALWAYS TRUE: These must be consistent
if system.DefendingFleet != nil {
    system.Colonization.IsColonized == true
}

if system.Colonization.IsColonized == true {
    system.DefendingFleet != nil
}
```

### **Rule 2: PlayerGameState Arrays**
```go
// These arrays should be mutually exclusive:
// - ColonizedSystems: Systems with embedded defending fleets
// - ActiveStacks: Stacks in free movement (NOT in systems)
// - MiningOperations: Resources being mined by active stacks
```

## Benefits

### Performance
- System combat resolution: Single document read/write
- Territory queries: Direct system queries
- Mining operations: Separate collection, no interference

### Consistency
- Colonization is atomic (one document update)
- No orphaned stacks or systems
- Clear ownership model
- **Victory ≠ Automatic Colonization**

### Scalability
- No array growth issues (only one defending fleet per system)
- Mining operations scale independently
- Player state normalized across collections

## Query Patterns

### Common Queries
```javascript
// Find player's territories
db.systems.find({"colonization.colonizedBy": playerId})

// Find player's active fleets
db.stacks.find({"playerId": playerId})

// Find orphaned systems (ready for claiming)
db.systems.find({"defendingFleet": {$exists: false}, "colonization.isColonized": false})

// Find systems under allied defense
db.systems.find({"defendingFleet.alliedFleets": {$exists: true, $ne: []}})
```

### Mining Operations
```javascript
// Find player's mining operations
db.stacks.find({"playerId": playerId, "movement.state": "mining"})

// Find available mining spots
db.asteroids.find({"resourceExtraction.miningFleets": {$size: 0}})
```

## Game Mechanic Support

### ✅ Supported Mechanics
- Single stack per system colonization
- Automatic fleet merging (same player)
- Symbolic fleet merging (allies)
- **Ranged combat victory ≠ automatic colonization**
- **Orphaned systems require physical entry to claim**
- **Allied survivors protect original owner**
- Multi-stack mining operations
- Free movement in space
- Resource extraction from asteroids/nebulas

### 🔄 State Transitions
- **Traveling → Colonizing**: Stack deleted, embedded in system
- **Colonizing → Defeated**: DefendingFleet removed, system orphaned
- **Defeated → Orphaned**: System becomes unclaimed territory
- **Orphaned → Colonized**: New stack enters and claims system
- **Mining → Traveling**: Stack document updated, references cleaned

## Key Combat Principles

### **1. Victory ≠ Colonization**
- Ranged combat victory does NOT automatically transfer system ownership
- Enemy must physically enter system to claim it
- Prevents "remote colonization" exploits

### **2. Original Owner Protection**
- System.Colonization.ColonizedBy tracks the original colonizer
- Allied fleets can defend but don't change ownership
- Only complete defeat + physical entry transfers ownership

### **3. Orphaned Systems**
- Systems without DefendingFleet are unclaimed territory
- Any player can claim orphaned systems by entering them
- Creates strategic opportunities for "land grabs"

This approach provides optimal performance for the specific game mechanics while maintaining data consistency and preventing exploits. 