# Modifier System Migration Guide

## Overview

This guide helps you migrate from the old StatMods system to the new V2 layered modifier system. The V2 system provides:

- **Source Tracking**: Know exactly where each modifier comes from
- **Transparency**: Debug and inspect all active modifiers
- **Scalability**: Easy to add new modifier sources without touching existing code
- **Context Awareness**: Conditional modifiers based on combat state, formations, etc.

## Quick Reference: Old vs New

### Computing Effective Ship Stats

**Old Way (DEPRECATED):**
```go
// Out of formation
effectiveShip, abilities := stack.EffectiveShip(Fighter)

// In formation
effectiveShip, abilities := stack.EffectiveShipInFormation(Fighter, 0)
```

**New Way (RECOMMENDED):**
```go
import "time"

now := time.Now()

// Simple out-of-combat calculation
effectiveShip, abilities := stack.EffectiveShipV2Simple(Fighter, now)

// With full modifier stack for debugging
ship, abilities, modStack := stack.EffectiveShipV2(Fighter, 0, now)

// In formation (recommended)
ship, abilities, modStack := stack.EffectiveShipInFormationV2(Fighter, 0, now)

// In combat with formation counter
ship, abilities, modStack := stack.EffectiveShipInCombat(Fighter, 0, EnemyFormationType, now)
```

### Using Helper Functions

**Quick calculations without modifier stack:**
```go
// Simplest way - just effective ship stats
ship, abilities := QuickEffectiveShip(stack, Fighter, 0, now)

// For combat
ship, abilities := QuickEffectiveShipInCombat(stack, Fighter, 0, FormationBox, now)
```

**Get just the modifier stack:**
```go
modStack := QuickModifierStack(stack, Fighter, 0, now, false, "")
```

## Migration Examples

### Example 1: Basic Ship Stats Display

**Before:**
```go
func DisplayShipStats(stack *ShipStack, shipType ShipType) {
    ship, abilities := stack.EffectiveShip(shipType)
    
    fmt.Printf("Attack: %d\n", ship.AttackDamage)
    fmt.Printf("HP: %d\n", ship.HP)
    fmt.Printf("Speed: %d\n", ship.Speed)
}
```

**After:**
```go
func DisplayShipStats(stack *ShipStack, shipType ShipType) {
    ship, abilities := stack.EffectiveShipV2Simple(shipType, time.Now())
    
    fmt.Printf("Attack: %d\n", ship.AttackDamage)
    fmt.Printf("HP: %d\n", ship.HP)
    fmt.Printf("Speed: %d\n", ship.Speed)
}
```

### Example 2: Combat Calculations

**Before:**
```go
func CalculateDamage(attacker *ShipStack, defender *ShipStack) int {
    ship, _ := attacker.EffectiveShipInFormation(Fighter, 0)
    
    damage := ship.AttackDamage
    
    // Manually apply formation counter
    if attacker.Formation != nil && defender.Formation != nil {
        mult := GetFormationCounterMultiplier(
            attacker.Formation.Type,
            defender.Formation.Type,
        )
        damage = int(float64(damage) * mult)
    }
    
    return damage
}
```

**After:**
```go
func CalculateDamage(attacker *ShipStack, defender *ShipStack) int {
    var enemyFormation FormationType
    if defender.Formation != nil {
        enemyFormation = defender.Formation.Type
    }
    
    // Formation counter is automatically included!
    ship, _, _ := attacker.EffectiveShipInCombat(
        Fighter, 
        0, 
        enemyFormation, 
        time.Now(),
    )
    
    return ship.AttackDamage
}
```

### Example 3: Debugging Modifiers

**Before (not possible):**
```go
// Can't easily see which modifiers are active or why
```

**After:**
```go
func DebugShipModifiers(stack *ShipStack, shipType ShipType) {
    breakdown := stack.GetModifierBreakdownForShip(
        shipType, 
        0, 
        time.Now(), 
        false,
    )
    
    for _, summary := range breakdown {
        if summary.IsActive {
            fmt.Printf("✓ %s: %s\n", summary.Source, summary.Description)
            if summary.ExpiresIn != nil {
                fmt.Printf("  Expires in: %.1fs\n", *summary.ExpiresIn)
            }
        } else {
            fmt.Printf("✗ %s: %s (Inactive)\n", summary.Source, summary.Description)
        }
    }
}
```

### Example 4: Comparing Loadouts

**Before (manual calculation):**
```go
func CompareGems(stack *ShipStack, shipType ShipType, newGem Gem) {
    // Get current stats
    before, _ := stack.EffectiveShipInFormation(shipType, 0)
    
    // Manually modify loadout
    loadout := stack.GetOrInitLoadout(shipType)
    oldGems := loadout.Sockets
    loadout.Sockets = append(loadout.Sockets, newGem)
    stack.Loadouts[shipType] = loadout
    
    // Get new stats
    after, _ := stack.EffectiveShipInFormation(shipType, 0)
    
    // Restore
    loadout.Sockets = oldGems
    stack.Loadouts[shipType] = loadout
    
    // Compare manually
    damageChange := after.AttackDamage - before.AttackDamage
    fmt.Printf("Damage change: %+d\n", damageChange)
}
```

**After (automatic comparison):**
```go
func CompareGems(stack *ShipStack, shipType ShipType, newGem Gem) {
    loadout := stack.GetOrInitLoadout(shipType)
    newLoadout := loadout
    newLoadout.Sockets = append(newLoadout.Sockets, newGem)
    
    beforeMods, afterMods, diff := CompareLoadoutChange(
        stack,
        shipType,
        0,
        newLoadout,
        time.Now(),
    )
    
    fmt.Printf("Added modifiers:\n")
    for _, layer := range diff.Added {
        fmt.Printf("  + %s\n", layer.Description)
    }
    
    fmt.Printf("\nRemoved modifiers:\n")
    for _, layer := range diff.Removed {
        fmt.Printf("  - %s\n", layer.Description)
    }
}
```

### Example 5: Manual Modifier Construction

**For advanced scenarios, build modifier stacks manually:**

```go
func ApplyCustomBuff(stack *ShipStack, shipType ShipType) *ModifierStack {
    now := time.Now()
    
    builder := NewModifierBuilder(now)
    
    // Add base modifiers from ship configuration
    loadout := stack.GetOrInitLoadout(shipType)
    builder.AddGemsFromLoadout(loadout)
    builder.AddRoleMode(stack.Role)
    
    // Add formation if present
    if stack.Formation != nil {
        position := stack.GetFormationPosition(shipType, 0)
        builder.
            AddFormationPosition(stack.Formation, position).
            AddFormationRoleSynergy(stack.Formation, position, stack.Role).
            AddGemPositionSynergy(loadout.Sockets, position)
    }
    
    // Add custom temporary buff
    customBuff := StatMods{
        Damage: DamageMods{
            LaserPct:      0.25,
            NuclearPct:    0.25,
            AntimatterPct: 0.25,
        },
        AccuracyPct: 0.15,
    }
    builder.AddBuff(
        "commander_presence",
        "Commander's Presence",
        customBuff,
        5 * time.Minute, // 5 minute duration
    )
    
    return builder.Build()
}
```

## Helper Function Reference

### Quick Calculations

| Function | Use Case |
|----------|----------|
| `QuickEffectiveShip()` | Basic out-of-combat ship stats |
| `QuickEffectiveShipInCombat()` | Combat stats with formation counter |
| `QuickModifierStack()` | Get modifier stack without computing ship |

### Comparison Utilities

| Function | Use Case |
|----------|----------|
| `CompareLoadoutChange()` | Preview gem socketing changes |
| `CompareFormationChange()` | Preview formation switch |
| `DiffModifierStacks()` | Compare two modifier stacks |

### Debugging & Analysis

| Function | Use Case |
|----------|----------|
| `GetModifierBreakdown()` | Detailed modifier analysis |
| `GetActiveModifierSources()` | List active modifier sources |
| `GetStackPowerRating()` | Simple power level calculation |

### Formation Utilities

| Function | Use Case |
|----------|----------|
| `GetFormationEffectiveness()` | Check formation matchup multiplier |
| `RecommendFormation()` | Suggest best formation vs enemy |
| `SimulateCombatModifiers()` | Preview combat modifiers |

### Batch Operations

| Function | Use Case |
|----------|----------|
| `BatchComputeEffectiveShips()` | Calculate all ships in stack |
| `MergeStacksModifiers()` | Combined modifiers for merged fleets |

## Common Patterns

### Pattern 1: UI Display with Modifier Breakdown

```go
func DisplayShipWithModifiers(stack *ShipStack, shipType ShipType) {
    now := time.Now()
    
    // Get effective stats
    ship, abilities, _ := stack.EffectiveShipInFormationV2(shipType, 0, now)
    
    // Get modifier breakdown for UI
    breakdown := stack.GetModifierBreakdownForShip(shipType, 0, now, false)
    
    // Display base stats
    fmt.Printf("=== %s Stats ===\n", shipType)
    fmt.Printf("Attack: %d\n", ship.AttackDamage)
    fmt.Printf("HP: %d\n", ship.HP)
    fmt.Printf("Speed: %d\n", ship.Speed)
    
    // Display modifiers
    fmt.Printf("\n=== Active Modifiers ===\n")
    for _, summary := range breakdown {
        if summary.IsActive {
            fmt.Printf("• %s\n", summary.Description)
        }
    }
}
```

### Pattern 2: Formation Selection UI

```go
func ShowFormationOptions(stack *ShipStack, enemyFormation FormationType) {
    shipType := Fighter
    now := time.Now()
    
    // Get current stats
    currentShip, _, _ := stack.EffectiveShipInFormationV2(shipType, 0, now)
    
    fmt.Printf("Current Formation: %s\n", stack.Formation.Type)
    fmt.Printf("Current Damage: %d\n\n", currentShip.AttackDamage)
    
    // Show all formation options
    for formationType := range FormationCatalog {
        beforeMods, afterMods, _ := CompareFormationChange(
            stack,
            shipType,
            0,
            formationType,
            now,
        )
        
        // Calculate damage change (simplified)
        effectiveness := GetFormationEffectiveness(formationType, enemyFormation)
        
        fmt.Printf("%s: %.0f%% effectiveness vs %s\n",
            formationType,
            effectiveness*100,
            enemyFormation,
        )
    }
    
    // Recommend best formation
    recommended, score := RecommendFormation(stack, enemyFormation, now)
    fmt.Printf("\nRecommended: %s (%.0f%% advantage)\n", recommended, score*100)
}
```

### Pattern 3: Combat Preview

```go
func PreviewCombat(attacker, defender *ShipStack) {
    now := time.Now()
    
    var defenderFormation FormationType
    if defender.Formation != nil {
        defenderFormation = defender.Formation.Type
    }
    
    // Get attacker's effective stats
    ship, _, _ := attacker.EffectiveShipInCombat(
        Fighter,
        0,
        defenderFormation,
        now,
    )
    
    // Show what modifiers will be active in combat
    combatMods := SimulateCombatModifiers(
        attacker,
        Fighter,
        0,
        defenderFormation,
        now,
    )
    
    fmt.Printf("=== Combat Preview ===\n")
    fmt.Printf("Effective Damage: %d\n", ship.AttackDamage)
    fmt.Printf("\nCombat Modifiers:\n")
    for _, mod := range combatMods {
        if mod.IsActive {
            fmt.Printf("  ✓ %s\n", mod.Description)
        }
    }
}
```

## Best Practices

### ✅ DO

1. **Use helper functions for common cases**
   ```go
   ship, abilities := QuickEffectiveShip(stack, shipType, 0, time.Now())
   ```

2. **Pass `time.Now()` for current calculations**
   ```go
   now := time.Now()
   ship, _, _ := stack.EffectiveShipV2(Fighter, 0, now)
   ```

3. **Use `GetModifierBreakdown()` for debugging**
   ```go
   breakdown := GetModifierBreakdown(stack, Fighter, 0, now, true, "")
   ```

4. **Leverage comparison functions for UI**
   ```go
   before, after, diff := CompareLoadoutChange(...)
   ```

5. **Use context-appropriate functions**
   ```go
   // Out of combat
   ship, _ := QuickEffectiveShip(stack, Fighter, 0, now)
   
   // In combat
   ship, _ := QuickEffectiveShipInCombat(stack, Fighter, 0, enemyFormation, now)
   ```

### ❌ DON'T

1. **Don't manually combine modifiers**
   ```go
   // ❌ Bad
   mods := CombineMods(roleMods, gemMods)
   mods = CombineMods(mods, formationMods)
   
   // ✅ Good
   ship, _, _ := stack.EffectiveShipInFormationV2(Fighter, 0, now)
   ```

2. **Don't use deprecated functions in new code**
   ```go
   // ❌ Bad (deprecated)
   ship, _ := stack.EffectiveShip(Fighter)
   
   // ✅ Good
   ship, _ := stack.EffectiveShipV2Simple(Fighter, time.Now())
   ```

3. **Don't ignore the modifier stack when debugging**
   ```go
   // ❌ Bad - can't see why stats are what they are
   ship, _ := QuickEffectiveShip(...)
   
   // ✅ Good - can inspect modifiers
   ship, _, modStack := stack.EffectiveShipV2(Fighter, 0, now)
   breakdown := modStack.GetSummary(ctx)
   ```

4. **Don't forget to pass formation enemy for combat**
   ```go
   // ❌ Bad - missing formation counter bonus
   ship, _, _ := ComputeEffectiveShipV2(stack, Fighter, 0, now, true, "")
   
   // ✅ Good
   ship, _, _ := stack.EffectiveShipInCombat(Fighter, 0, enemyFormation, now)
   ```

## Timeline

### Phase 1: Compatibility (Current)
- ✅ Old functions marked deprecated but still work
- ✅ New V2 functions available
- ✅ Helper functions added for easy migration
- ✅ Both systems run in parallel

### Phase 2: Migration (Recommended)
- Update all call sites to use V2 functions
- Use helper functions for simple cases
- Add modifier debugging to UIs
- Test thoroughly

### Phase 3: Deprecation (Future)
- Old functions removed
- Only V2 system remains
- Simplified codebase

## Troubleshooting

### "Too many return values"
**Problem:** Old function returns 2 values, new returns 3
```go
ship, abilities := stack.EffectiveShip(Fighter) // Old: 2 returns
ship, abilities, modStack := stack.EffectiveShipV2(...) // New: 3 returns
```
**Solution:** Use `EffectiveShipV2Simple()` for drop-in replacement
```go
ship, abilities := stack.EffectiveShipV2Simple(Fighter, time.Now())
```

### "Stats seem wrong"
**Problem:** Missing formation or composition bonuses
**Solution:** Check that you're using the formation-aware function
```go
// ❌ Wrong - no formation bonuses
ship, _ := stack.EffectiveShipV2Simple(Fighter, now)

// ✅ Correct - includes all bonuses
ship, _, _ := stack.EffectiveShipInFormationV2(Fighter, 0, now)
```

### "Can't debug why stat is X"
**Problem:** Can't see modifier sources
**Solution:** Use `GetModifierBreakdown()`
```go
breakdown := stack.GetModifierBreakdownForShip(Fighter, 0, now, false)
for _, summary := range breakdown {
    fmt.Printf("%s: %+v\n", summary.Description, summary.Mods)
}
```

### "Formation counter not applying"
**Problem:** Not passing enemy formation to combat function
**Solution:** Always pass enemy formation for combat calculations
```go
ship, _, _ := stack.EffectiveShipInCombat(
    Fighter,
    0,
    defender.Formation.Type, // Required!
    time.Now(),
)
```

## Additional Resources

- **MODIFIER_SYSTEM.md**: Complete V2 system documentation
- **compute_v2.go**: Core V2 implementation
- **compute_helpers.go**: Helper function implementations
- **modifier_stack.go**: Modifier layer types and operations
- **modifier_builder.go**: Builder pattern API
- **formation_examples.go**: Working examples using V2

## Questions?

The V2 system is designed to be easier to use and maintain than the old system. If you find something confusing or have questions, check:

1. The helper functions in `compute_helpers.go`
2. Examples in `formation_examples.go`
3. The MODIFIER_SYSTEM.md documentation

The key philosophy: **Transparency over simplicity**. We'd rather have a slightly more complex system that you can fully understand and debug than a simple system that feels like magic.
