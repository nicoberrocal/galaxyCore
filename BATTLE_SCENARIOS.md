# Battle Scenarios and Data Handling

This document explains how different battle scenarios are handled in the data model, including system battles, free space battles, and mining location battles.

## Battle Types Overview

### **1. System Battles** (Hybrid Approach)
- **Location**: Star systems
- **Data Storage**: Embedded in `System.DefendingFleet`
- **Outcome**: Can result in colonization changes
- **Atomic**: Single document updates

### **2. Free Space Battles** (Direct Stack Updates)
- **Location**: Empty space between objects
- **Data Storage**: Direct updates to `ShipStack` documents
- **Outcome**: Fleet damage only, no territory changes
- **Atomic**: Individual stack updates

### **3. Mining Location Battles** (Direct Stack Updates)
- **Location**: Asteroids and nebulas
- **Data Storage**: Direct updates to `ShipStack` documents
- **Outcome**: Fleet damage + potential mining interruption
- **Atomic**: Individual stack updates

## System Battles (Colonization Context)

### **Combat Resolution**
```go
// Atomic system battle resolution
func resolveSystemBattle(systemID, attackerID, defenderID) {
    db.systems.updateOne(
        {"_id": systemID},
        {
            "$set": {
                "defendingFleet.ships": updatedFleetComposition,
                "defendingFleet.alliedFleets": updatedAllies
            }
        }
    )
}
```

### **Key Principles**
- **Victory ≠ Colonization**: Ranged combat doesn't automatically transfer ownership
- **Physical Entry Required**: Enemy must enter system to claim it
- **Allied Survivors**: Original owner keeps system if allies survive
- **Orphaned Systems**: Complete defeat creates unclaimed territory

## Free Space Battles

### **Battle Initiation**
```go
// Two stacks encounter each other in empty space
func initiateFreeSpaceBattle(stack1ID, stack2ID) {
    // Update both stacks to battle state
    db.stacks.updateOne(
        {"_id": stack1ID},
        {
            "$set": {
                "battle.isInCombat": true,
                "battle.enemyStackId": stack2ID,
                "battle.enemyPlayerId": stack2PlayerID,
                "battle.battleLocation": "empty_space",
                "movement.state": "in_combat"
            }
        }
    )
}
```

### **Combat Resolution**
```go
// Apply damage to both stacks
func resolveFreeSpaceBattle(stack1ID, stack2ID, damage1, damage2) {
    // Update stack 1
    db.stacks.updateOne(
        {"_id": stack1ID},
        {
            "$set": {"ships": updatedShips1},
            "$unset": {"battle": ""},
            "$set": {"movement.state": "idle"}
        }
    )
    
    // Update stack 2
    db.stacks.updateOne(
        {"_id": stack2ID},
        {
            "$set": {"ships": updatedShips2},
            "$unset": {"battle": ""},
            "$set": {"movement.state": "idle"}
        }
    )
}
```

### **Key Characteristics**
- **No Territory Impact**: Only fleet damage
- **No Resource Loss**: Resources remain with players
- **Stack Survival**: Stacks continue to exist after battle
- **Movement Interruption**: Battle pauses movement orders

## Mining Location Battles

### **Battle Initiation at Asteroid/Nebula**
```go
// Stack encounters enemy while mining
func initiateMiningBattle(miningStackID, enemyStackID, locationID, locationType) {
    db.stacks.updateOne(
        {"_id": miningStackID},
        {
            "$set": {
                "battle.isInCombat": true,
                "battle.enemyStackId": enemyStackID,
                "battle.enemyPlayerId": enemyPlayerID,
                "battle.battleLocation": locationType, // "asteroid" or "nebula"
                "battle.locationId": locationID,
                "movement.state": "in_combat"
            }
        }
    )
}
```

### **Combat Resolution with Mining Impact**
```go
// Battle affects both fleet and mining operation
func resolveMiningBattle(miningStackID, enemyStackID, locationID) {
    // Update mining stack
    db.stacks.updateOne(
        {"_id": miningStackID},
        {
            "$set": {"ships": updatedShips},
            "$unset": {"battle": ""},
            "$set": {"movement.state": "mining"} // Resume mining if survived
        }
    )
    
    // Update mining location if stack was defeated
    if miningStackDefeated {
        db.asteroids.updateOne(
            {"_id": locationID},
            {
                "$pull": {"resourceExtraction.miningFleets": miningStackID}
            }
        )
    }
}
```

### **Key Characteristics**
- **Mining Interruption**: Battle pauses resource extraction
- **Location Impact**: Defeated stack removed from mining operation
- **Resource Loss**: Mining time lost during battle
- **Stack Survival**: Surviving stack can resume mining

## Data Consistency and Exploit Prevention

### **Race Condition Prevention**

#### **System Battles**
```go
// Atomic updates prevent race conditions
// Only one combat resolution can succeed per system
db.systems.updateOne(
    {"_id": systemID, "defendingFleet": currentFleetState}, // Optimistic locking
    combatUpdate
)
```

#### **Free Space Battles**
```go
// Individual stack updates prevent conflicts
// Each stack can only be in one battle at a time
db.stacks.updateOne(
    {"_id": stackID, "battle.isInCombat": false}, // Ensure not already in battle
    battleInitiation
)
```

### **Data Validation Rules**

#### **Stack Battle State**
```go
func ValidateStackBattleState(stack ShipStack) error {
    // Rule: Can't be in battle if embedded in system
    if stack.Battle.IsInCombat {
        // This stack should be in free movement, not defending a system
        return nil
    }
    
    // Rule: Battle state should be consistent
    if stack.Battle.IsInCombat && stack.Battle.EnemyStackID == "" {
        return errors.New("in combat but no enemy specified")
    }
    
    return nil
}
```

#### **Player Game State**
```go
func ValidatePlayerGameState(state PlayerGameState) error {
    // Rule: ActiveStacks should not include system-defending fleets
    // (those are tracked in ColonizedSystems)
    
    // Rule: MiningOperations should match active mining stacks
    // (stacks with movement.state = "mining")
    
    return nil
}
```

## Query Patterns for Battle Management

### **Find Stacks in Combat**
```javascript
// Find all stacks currently in battle
db.stacks.find({"battle.isInCombat": true})

// Find battles at specific locations
db.stacks.find({
    "battle.isInCombat": true,
    "battle.battleLocation": "asteroid"
})
```

### **Find Available Mining Spots**
```javascript
// Find asteroids with no active battles
db.asteroids.find({
    "resourceExtraction.miningFleets": {
        $not: {
            $elemMatch: {
                $in: db.stacks.distinct("_id", {"battle.isInCombat": true})
            }
        }
    }
})
```

### **Find Player's Combat Status**
```javascript
// Find player's stacks in combat
db.stacks.find({
    "playerId": playerId,
    "battle.isInCombat": true
})

// Find systems under attack by player
db.systems.find({
    "defendingFleet.enemyPlayerId": playerId
})
```

## Performance Considerations

### **System Battles**
- **Single Document Update**: Maximum performance
- **Atomic Operations**: No race conditions
- **Territory Queries**: Direct system lookups

### **Free Space Battles**
- **Individual Updates**: Each stack updated separately
- **Battle State Tracking**: Prevents duplicate battles
- **Movement Interruption**: Pauses other activities

### **Mining Location Battles**
- **Dual Updates**: Stack + mining location updates
- **Resource Impact**: Mining interruption tracked
- **Location Cleanup**: Defeated stacks removed from mining

## Exploit Prevention Summary

### **1. Race Conditions**
- **System Battles**: Atomic updates prevent multiple combat resolutions
- **Free Space Battles**: Battle state prevents duplicate engagements
- **Mining Battles**: Location state prevents mining conflicts

### **2. Data Inconsistencies**
- **System State**: DefendingFleet ↔ Colonization consistency enforced
- **Stack State**: Battle state validation prevents invalid states
- **Player State**: Mutually exclusive arrays prevent overlaps

### **3. Resource Exploits**
- **Mining Interruption**: Battles pause resource extraction
- **Territory Protection**: Victory ≠ automatic colonization
- **Allied Defense**: Original owner protection mechanisms

This battle system ensures realistic combat mechanics while maintaining data consistency and preventing common gaming exploits. 