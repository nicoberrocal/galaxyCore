# Compute System Quick Reference

## Getting Effective Ship Stats

### Simple Usage (Recommended)
```go
import "time"

// Get effective stats for a ship in a stack
ship, abilities, modStack := stack.EffectiveShipV2(shipType, bucketIndex, time.Now())

// Even simpler (uses bucket 0, no combat context)
ship, abilities := stack.EffectiveShipV2Simple(shipType, time.Now())
```

### With Combat Context
```go
// In combat with formation counters
ship, abilities, modStack := ComputeEffectiveShipV2(
    stack,
    shipType,
    bucketIndex,
    time.Now(),
    true,              // inCombat = true
    enemyFormationType, // for counter bonuses
)
```

### With Formation Tree
```go
// Include formation tree bonuses
modStack, finalMods, grants := ComputeLoadoutV2WithTree(
    ship,
    role,
    loadout,
    formation,
    position,
    ships,
    treeState,  // player's tree progress
    time.Now(),
    inCombat,
)
```

## Debugging Modifiers

### View All Active Modifiers
```go
breakdown := GetModifierBreakdown(stack, shipType, bucketIndex, time.Now(), inCombat, enemyFormation)

for _, summary := range breakdown {
    fmt.Printf("%s: %s (Priority: %d, Active: %v)\n", 
        summary.Source, summary.Description, summary.Priority, summary.IsActive)
}
```

### Compare Configurations
```go
// Before changing gems
_, _, stackBefore := stack.EffectiveShipV2(shipType, 0, time.Now())

// After changing gems
stack.Loadouts[shipType].Sockets = newGems
_, _, stackAfter := stack.EffectiveShipV2(shipType, 0, time.Now())

// See what changed
diff := DiffModifierStacks(stackBefore, stackAfter)
fmt.Printf("Added: %v\n", diff.Added)
fmt.Printf("Removed: %v\n", diff.Removed)
fmt.Printf("Changed: %v\n", diff.Changed)
```

## Adding Modifiers Manually

### Using ModifierBuilder
```go
builder := NewModifierBuilder(time.Now())

// Add from various sources
builder.AddGems(gems)
builder.AddRoleMode(RoleTactical)
builder.AddFormationPosition(formation, PositionFront)
builder.AddActiveAbilities([]AbilityID{AbilityOverload}, durations)

// Build and resolve
stack := builder.Build()
finalMods := stack.Resolve(ResolveContext{
    Now:      time.Now(),
    InCombat: true,
})

// Apply to ship
effectiveShip := ApplyStatModsToShip(baseShip, finalMods)
```

## Common Patterns

### Check If Ability Is Active
```go
mods := GetAbilityMods(AbilityOverload)
if !isZeroMods(mods) {
    // Ability has stat effects
}
```

### Get Composition Bonuses
```go
mods, bonuses := EvaluateCompositionBonuses(stack.Ships)
for _, bonus := range bonuses {
    fmt.Printf("Active: %s - %s\n", bonus.Type, bonus.Description)
}
```

### Get Gem-Position Synergies
```go
synergyMods := ApplyGemPositionEffects(gems, PositionFront)
```

### Find Best Formation Template
```go
template := FindBestTemplate(stack.Ships, stack.Role, enemyFormation)
if template != nil {
    fmt.Printf("Recommended: %s\n", template.Name)
}
```

## Modifier Sources (Priority Order)

1. **Gems** (Priority: 100) - Permanent while socketed
2. **GemWords** (Priority: 110) - Permanent when sequence matches
3. **Role Mode** (Priority: 200) - Semi-permanent until switched
4. **Formation** (Priority: 300) - Active while formation set
5. **Formation Tree** (Priority: 350) - Permanent once unlocked
6. **Synergies** (Priority: 400) - Conditional bonuses
7. **Composition** (Priority: 500) - Fleet-wide bonuses
8. **Abilities** (Priority: 600) - Temporary when activated
9. **Buffs** (Priority: 700) - Temporary effects
10. **Debuffs** (Priority: 800) - Negative temporary effects
11. **Environment** (Priority: 900) - Location-based

## Key Types

### StatMods
```go
type StatMods struct {
    // Damage
    Damage DamageMods // LaserPct, NuclearPct, AntimatterPct
    
    // Combat
    AttackIntervalPct float64
    AttackRangeDelta  int
    AccuracyPct       float64
    CritPct           float64
    
    // Defense
    LaserShieldDelta      int
    NuclearShieldDelta    int
    AntimatterShieldDelta int
    BucketHPPct           float64
    EvasionPct            float64
    
    // Mobility
    SpeedDelta int
    
    // Utility
    VisibilityDelta int
    // ... and more
}
```

### ResolveContext
```go
type ResolveContext struct {
    Now            time.Time
    InCombat       bool
    HasFormation   bool
    FormationType  FormationType
    EnemyFormation FormationType
}
```

## File Locations

- **Core System**: `compute_v2.go`, `modifier_builder.go`, `modifier_stack.go`
- **Ability Effects**: `ability_effects.go`
- **Composition/Gems**: `formation_composition.go`
- **Formation Tree**: `formation_tree_compute.go`
- **Examples**: `formation_examples.go`
- **Architecture Doc**: `COMPUTE_ARCHITECTURE.md`
