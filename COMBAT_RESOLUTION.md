# Combat Resolution and System State Logic

This document explains the corrected combat resolution mechanics and system state consistency rules.

## System State Consistency Rules

### **Rule 1: Colonization ↔ DefendingFleet Relationship**
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

## Combat Scenarios

### **Scenario 1: Stack Enters Unclaimed System**
```
Initial State: System has no DefendingFleet, Colonization.IsColonized = false
Action: Player A stack enters system
Result: 
- Stack document deleted from stacks collection
- System.DefendingFleet = stack data
- System.Colonization.IsColonized = true
- System.Colonization.ColonizedBy = Player A
- PlayerGameState.ColonizedSystems += system ID
```

### **Scenario 2: Same Player Fleet Merging**
```
Initial State: System defended by Player A fleet
Action: Another Player A stack enters system
Result:
- Incoming stack ships merged into DefendingFleet.Ships
- Incoming stack document deleted
- System remains colonized by Player A
- No change to Colonization.ColonizedBy
```

### **Scenario 3: Allied Fleet Merging**
```
Initial State: System defended by Player A fleet
Action: Player B (ally) stack enters system
Result:
- Incoming stack ships merged into DefendingFleet.Ships
- AlliedFleet info added to DefendingFleet.AlliedFleets
- Incoming stack document deleted
- System remains colonized by Player A (original owner)
- DefendingFleet.PlayerID remains Player A
```

### **Scenario 4: Enemy Combat Victory (Ranged Combat)**
```
Initial State: System defended by Player A fleet (with/without allies)
Action: Player C defeats defending fleet in ranged combat
Result:
- DefendingFleet.Ships reduced/eliminated based on combat
- If DefendingFleet.Ships completely destroyed:
  - System.DefendingFleet = nil
  - System.Colonization.IsColonized = false
  - System.Colonization.ColonizedBy = nil
  - System becomes "orphaned" (unclaimed)
- Player C does NOT automatically colonize
- Player C must physically enter system to claim it
```

### **Scenario 5: Allied Fleet Survivors**
```
Initial State: System defended by Player A + Player B allies
Action: Player C defeats main fleet but allies survive
Result:
- DefendingFleet.Ships updated (main fleet destroyed, allies remain)
- DefendingFleet.AlliedFleets updated (defeated allies removed)
- System.Colonization.ColonizedBy remains Player A (original owner)
- System remains colonized by Player A
- DefendingFleet.PlayerID remains Player A
```

### **Scenario 6: Enemy Claims Orphaned System**
```
Initial State: System is orphaned (no DefendingFleet)
Action: Player C stack enters orphaned system
Result:
- Stack document deleted from stacks collection
- System.DefendingFleet = stack data
- System.Colonization.IsColonized = true
- System.Colonization.ColonizedBy = Player C
- PlayerGameState.ColonizedSystems += system ID
```

## Key Combat Principles

### **1. Victory ≠ Colonization**
- Ranged combat victory does NOT automatically transfer system ownership
- Enemy must physically enter system to claim it
- This prevents "remote colonization" exploits

### **2. Original Owner Protection**
- System.Colonization.ColonizedBy tracks the original colonizer
- Allied fleets can defend but don't change ownership
- Only complete defeat + physical entry transfers ownership

### **3. Orphaned Systems**
- Systems without DefendingFleet are unclaimed territory
- Any player can claim orphaned systems by entering them
- Creates strategic opportunities for "land grabs"

### **4. Atomic Combat Resolution**
- All combat state changes happen in single System document update
- No intermediate states where system is inconsistent
- Prevents race conditions and exploits

## Data Validation Rules

### **System Validation**
```go
func ValidateSystem(system System) error {
    // Rule 1: DefendingFleet ↔ Colonization consistency
    if system.DefendingFleet != nil && !system.Colonization.IsColonized {
        return errors.New("system has defending fleet but is not colonized")
    }
    if system.Colonization.IsColonized && system.DefendingFleet == nil {
        return errors.New("system is colonized but has no defending fleet")
    }
    
    // Rule 2: DefendingFleet.PlayerID consistency
    if system.DefendingFleet != nil {
        if system.DefendingFleet.PlayerID != system.Colonization.ColonizedBy {
            // This is valid for allied defense scenarios
            // but should be logged for tracking
        }
    }
    
    return nil
}
```

### **PlayerGameState Validation**
```go
func ValidatePlayerGameState(state PlayerGameState) error {
    // Rule: ColonizedSystems should not overlap with ActiveStacks
    // (a stack can't be both embedded in system AND in free movement)
    
    // Rule: MiningOperations should reference valid asteroids/nebulas
    // where player has active mining stacks
    
    return nil
}
```

## Query Examples

### **Find Orphaned Systems**
```javascript
db.systems.find({
    "defendingFleet": {$exists: false},
    "colonization.isColonized": false
})
```

### **Find Systems Under Allied Defense**
```javascript
db.systems.find({
    "defendingFleet.alliedFleets": {$exists: true, $ne: []}
})
```

### **Find Player's Territories**
```javascript
db.systems.find({
    "colonization.colonizedBy": playerId,
    "defendingFleet": {$exists: true}
})
```

### **Find Systems Ready for Claiming**
```javascript
db.systems.find({
    "defendingFleet": {$exists: false},
    "colonization.isColonized": false
})
```

This combat resolution system ensures realistic territorial control while maintaining data consistency and preventing exploits. 