# Layered Modifier System

## Overview

The new modifier system provides **fine-grained control** over stat modifications from multiple sources while maintaining the delta/pct semantics of the original `StatMods` approach.

## Key Features

### 1. **Source Tracking**
Every modifier layer tracks its source, allowing you to:
- Debug which bonuses are active
- Display modifier breakdowns in UI
- Remove modifiers by source when conditions change

### 2. **Lifetime Management**
Modifiers have different lifetimes:
- **Permanent**: Gems, role modes (until changed)
- **Conditional**: Formation bonuses (while formation active)
- **Temporary**: Abilities, buffs, debuffs (duration-based)

### 3. **Priority System**
Modifiers are applied in priority order:
```
Gems (100) → GemWords (150) → Role Mode (200) → Formation (300) 
→ Composition (350) → Synergies (400) → Environment (500) 
→ Abilities (600) → Buffs (700) → Debuffs (800)
```

### 4. **Context-Aware Resolution**
Modifiers can be conditional on:
- Combat state (in combat / out of combat)
- Formation presence
- Time (expiration)

## Architecture

### Core Types

```go
// ModifierLayer - A single layer of modifiers
type ModifierLayer struct {
    Source      ModifierSource  // Where it comes from
    SourceID    string          // Specific identifier
    Description string          // Human-readable
    Mods        StatMods        // The actual modifiers
    AppliedAt   time.Time       // When applied
    ExpiresAt   *time.Time      // When it expires (nil = permanent)
    Priority    int             // Resolution order
    ActiveInCombat    *bool     // Combat requirement
    RequiresFormation *bool     // Formation requirement
}

// ModifierStack - Collection of layers
type ModifierStack struct {
    Layers []ModifierLayer
}
```

### Builder Pattern

Use `ModifierBuilder` for fluent construction:

```go
builder := NewModifierBuilder(now)
builder.
    AddGemsFromLoadout(loadout).
    AddRoleMode(RoleTactical).
    AddFormationPosition(formation, PositionFront).
    AddCompositionBonus(ships)

stack := builder.Build()
```

## Usage Examples

### Example 1: Basic Ship Modifiers

```go
// Get all modifiers for a ship in a stack
modStack, finalMods := ComputeStackModifiers(
    stack,           // ShipStack
    Fighter,         // ShipType
    0,               // bucketIndex
    time.Now(),      // now
    true,            // inCombat
    FormationBox,    // enemyFormation
)

// Apply to ship blueprint
effectiveShip := ApplyStatModsToShip(blueprint, finalMods)
```

### Example 2: Debugging Modifiers

```go
// Get a breakdown of all active modifiers
breakdown := GetModifierBreakdown(
    stack,
    Fighter,
    0,
    time.Now(),
    true,
    FormationBox,
)

// Display in UI
for _, summary := range breakdown {
    fmt.Printf("%s: %s (Active: %v)\n", 
        summary.Source, 
        summary.Description, 
        summary.IsActive,
    )
    if summary.ExpiresIn != nil {
        fmt.Printf("  Expires in: %.1fs\n", *summary.ExpiresIn)
    }
}
```

### Example 3: Temporary Ability Modifiers

```go
builder := NewModifierBuilder(now)

// Add base modifiers
builder.
    AddGemsFromLoadout(loadout).
    AddRoleMode(stack.Role)

// Add temporary ability buff
abilityMods := StatMods{
    Damage: DamageMods{LaserPct: 0.25},
    AccuracyPct: 0.15,
}
builder.AddAbility(
    AbilityFocusFire,
    abilityMods,
    30 * time.Second, // 30 second duration
)

stack := builder.Build()
```

### Example 4: Formation Counter Modifiers

```go
builder := NewModifierBuilder(now)

// Add formation counter bonus (dynamic, determined at engagement)
builder.AddFormationCounter(
    FormationVanguard,  // Your formation
    FormationBox,       // Enemy formation
    true,               // inCombat
)

// This adds a damage multiplier based on formation matchup
// Only active during combat
```

### Example 5: Comparing Equipment Changes

```go
// Before: Current loadout
beforeStack, beforeMods := ComputeStackModifiers(
    stack, Fighter, 0, now, false, "",
)

// Simulate socketing a new gem
newLoadout := stack.GetOrInitLoadout(Fighter)
newLoadout.Sockets = append(newLoadout.Sockets, newGem)

// After: New loadout
afterStack, afterMods := ComputeStackModifiers(
    stack, Fighter, 0, now, false, "",
)

// Compare
diff := DiffModifierStacks(beforeStack, afterStack)
fmt.Printf("Added: %d, Removed: %d, Changed: %d\n",
    len(diff.Added), len(diff.Removed), len(diff.Changed))
```

## Modifier Sources

### Permanent Sources
- **`SourceGem`**: Socketed gems (while equipped)
- **`SourceGemWord`**: GemWord pattern bonuses
- **`SourceRoleMode`**: Active role mode bonuses

### Conditional Sources (Formation)
- **`SourceFormationPosition`**: Position-based bonuses (front/flank/back/support)
- **`SourceFormationRole`**: Role + formation synergy
- **`SourceFormationCounter`**: Formation matchup bonuses (determined per engagement)
- **`SourceComposition`**: Fleet composition bonuses
- **`SourceGemPosition`**: Gem + position synergy

### Temporary Sources
- **`SourceAbility`**: Active ability effects (duration-based)
- **`SourceBuff`**: Allied buffs
- **`SourceDebuff`**: Enemy debuffs

### Environmental Sources
- **`SourceEnvironment`**: Terrain, nebula effects
- **`SourceAnchored`**: Anchoring penalties/bonuses

## Integration with Existing Systems

### Gems
```go
// Old way (still works)
socketMods, grants, gemWords := EvaluateGemSockets(loadout.Sockets)

// New way (more transparent)
builder.AddGemsFromLoadout(loadout)
// OR for individual control:
builder.AddGems(loadout.Sockets)
builder.AddGemWords(gemWords)
```

### Role Modes
```go
// Old way
roleMods := RoleModeMods(role, shipType)

// New way
builder.AddRoleMode(role)
```

### Formations
```go
// Old way
mods = formation.ApplyPositionBonusesToShip(position, mods)
mods = ApplyFormationRoleModifiers(mods, formation, position, role)

// New way
builder.
    AddFormationPosition(formation, position).
    AddFormationRoleSynergy(formation, position, role).
    AddGemPositionSynergy(gems, position)
```

## Migration Path

### Phase 1: Parallel Systems (Current)
- Keep existing `ComputeLoadout` and `EffectiveShipInFormation`
- Add new `ComputeLoadoutV2` and `ComputeEffectiveShipV2`
- Use V2 functions in new code

### Phase 2: Gradual Migration
- Update combat resolution to use V2
- Update UI to display modifier breakdowns
- Add modifier stack to `ShipStack` struct

### Phase 3: Full Adoption
- Replace old functions with V2 versions
- Remove deprecated code
- Optimize performance

## Performance Considerations

### Optimization Tips
1. **Cache resolved modifiers** when context doesn't change
2. **Reuse builder instances** for batch operations
3. **Remove expired layers** periodically with `RemoveExpired()`
4. **Use combined layers** for simple cases (e.g., `AddGemsFromLoadoutCombined`)

### Memory Usage
- Each layer is ~200 bytes
- Typical stack: 10-20 layers = 2-4 KB
- Acceptable overhead for the transparency gained

## Advanced Features

### Custom Modifiers
```go
// Add custom modifier from any source
builder.stack.AddLayer(ModifierLayer{
    Source:      ModifierSource("custom"),
    SourceID:    "my_custom_effect",
    Description: "Custom Effect",
    Mods:        myCustomMods,
    AppliedAt:   now,
    Priority:    PrioritySynergy,
})
```

### Conditional Logic
```go
// Combat-only modifier
combatOnly := true
builder.stack.AddConditional(
    SourceAbility,
    "combat_boost",
    "Combat Boost",
    mods,
    PriorityAbility,
    now,
    &combatOnly,  // Only active in combat
    nil,          // No formation requirement
)
```

### Stacking Rules
By default, all modifiers stack additively (as per `CombineMods`). For custom stacking:
1. Use priority to control application order
2. Implement custom resolution logic in `Resolve()`
3. Use source filtering to prevent duplicate effects

## Future Enhancements

### Planned Features
1. **Diminishing returns**: Cap certain stats at thresholds
2. **Multiplicative stacking**: Some modifiers multiply instead of add
3. **Conditional triggers**: Activate modifiers based on events
4. **Modifier persistence**: Save/load modifier stacks from DB
5. **Modifier history**: Track modifier changes over time

### Integration Points
- **Combat system**: Apply debuffs from enemy attacks
- **Ability system**: Generate modifiers from ability definitions
- **Environment system**: Apply terrain/weather effects
- **Event system**: Trigger modifiers from game events

## Best Practices

### DO
✓ Use the builder pattern for construction  
✓ Set descriptive `Description` fields for UI  
✓ Use appropriate priorities for ordering  
✓ Clean up expired modifiers regularly  
✓ Use context-aware resolution  

### DON'T
✗ Mutate layers after adding to stack  
✗ Use string concatenation for SourceID (use consistent format)  
✗ Forget to set expiration for temporary effects  
✗ Apply modifiers twice (check for duplicates)  
✗ Ignore the priority system  

## Troubleshooting

### Modifiers Not Applying
1. Check `isLayerApplicable()` conditions
2. Verify priority order
3. Check expiration times
4. Verify context flags (inCombat, hasFormation)

### Unexpected Values
1. Use `GetSummary()` to inspect active layers
2. Check for duplicate sources
3. Verify `CombineMods` logic for your stat
4. Check for conflicting conditional flags

### Performance Issues
1. Remove expired layers regularly
2. Use combined layers for simple cases
3. Cache resolved modifiers when possible
4. Profile with `go test -bench`

## API Reference

See inline documentation in:
- `modifier_stack.go` - Core types and stack operations
- `modifier_builder.go` - Builder pattern and integration
- `compute_v2.go` - High-level compute functions
