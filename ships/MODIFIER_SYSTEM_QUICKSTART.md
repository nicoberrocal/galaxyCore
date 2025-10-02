# Modifier System Quick Start

**Problem**: StatMods scattered everywhere, hard to maintain, no transparency

**Solution**: Layered V2 modifier system with source tracking and context awareness

## TL;DR

### Before (Old, Messy)
```go
// Modifiers from everywhere, manually combined
roleMods := RoleModeMods(role, shipType)
gemMods, _, _ := EvaluateGemSockets(gems)
formationMods := formation.ApplyPositionBonusesToShip(position, baseMods)
finalMods := CombineMods(CombineMods(roleMods, gemMods), formationMods)
// ... more manual combining
```

### After (New, Clean)
```go
// Everything automatic, fully transparent
ship, abilities := QuickEffectiveShip(stack, Fighter, 0, time.Now())

// Or with full debugging
ship, abilities, modStack := stack.EffectiveShipV2(Fighter, 0, time.Now())
breakdown := modStack.GetSummary(ctx)
```

## 30-Second Migration

### Replace This
```go
ship, abilities := stack.EffectiveShip(Fighter)
ship, abilities := stack.EffectiveShipInFormation(Fighter, 0)
```

### With This
```go
ship, abilities := stack.EffectiveShipV2Simple(Fighter, time.Now())
ship, _, _ := stack.EffectiveShipInFormationV2(Fighter, 0, time.Now())
```

Done! Everything else just works.

## Common Tasks

### 1. Get Ship Stats for UI
```go
ship, abilities := QuickEffectiveShip(stack, Fighter, 0, time.Now())
fmt.Printf("Attack: %d, HP: %d\n", ship.AttackDamage, ship.HP)
```

### 2. Calculate Combat Damage
```go
ship, _, _ := stack.EffectiveShipInCombat(
    Fighter, 
    0, 
    enemyFormation, 
    time.Now(),
)
damage := ship.AttackDamage // Includes formation counter automatically!
```

### 3. Debug Why Stats Are Wrong
```go
breakdown := stack.GetModifierBreakdownForShip(Fighter, 0, time.Now(), false)
for _, mod := range breakdown {
    if mod.IsActive {
        fmt.Printf("✓ %s\n", mod.Description)
    }
}
```

### 4. Preview Gem Change
```go
newLoadout := currentLoadout
newLoadout.Sockets = append(newLoadout.Sockets, newGem)

before, after, diff := CompareLoadoutChange(
    stack, Fighter, 0, newLoadout, time.Now(),
)

fmt.Printf("Damage: %d → %d\n", 
    before.Damage.LaserPct, 
    after.Damage.LaserPct,
)
```

### 5. Find Best Formation
```go
recommended, score := RecommendFormation(stack, enemyFormation, time.Now())
fmt.Printf("Use %s for %.0f%% advantage\n", recommended, score*100)
```

## Key Benefits

### 1. **Source Tracking**
Know exactly where each modifier comes from:
- Gems
- Role Mode
- Formation Position
- Formation Counter
- Composition Bonuses
- Abilities
- Buffs/Debuffs

### 2. **No More Manual Combining**
Old way:
```go
mods := CombineMods(a, b)
mods = CombineMods(mods, c)
mods = CombineMods(mods, d)
// Easy to forget one!
```

New way:
```go
ship, _, _ := stack.EffectiveShipInFormationV2(Fighter, 0, now)
// Everything included automatically
```

### 3. **Context-Aware**
Modifiers activate only when appropriate:
- Combat-only bonuses (formation counter)
- Formation-dependent bonuses
- Temporary abilities with expiration
- Out-of-combat regeneration

### 4. **Easy Debugging**
```go
breakdown := GetModifierBreakdown(stack, Fighter, 0, now, true, "")
// See EXACTLY what's active and why
```

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│            ModifierStack                        │
│  ┌──────────────────────────────────────────┐   │
│  │  Layer 1: Gem (Photon III)               │   │
│  │  Layer 2: GemWord (Photon Overcharge)    │   │
│  │  Layer 3: Role Mode (Tactical)           │   │
│  │  Layer 4: Formation (Line/Front)         │   │
│  │  Layer 5: Composition (Strike Force)     │   │
│  │  Layer 6: Gem+Position Synergy           │   │
│  │  Layer 7: Formation Counter (+25%)       │   │
│  │  Layer 8: Active Ability (Focus Fire)    │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
                      ↓
              Resolve(context)
                      ↓
            ┌─────────────────┐
            │  Final StatMods │
            └─────────────────┘
                      ↓
          ApplyStatModsToShip(ship)
                      ↓
            ┌─────────────────┐
            │ Effective Ship  │
            └─────────────────┘
```

## Priority System

Modifiers apply in order (low → high priority):

1. **Gems** (100) - Base equipment
2. **GemWords** (150) - Pattern bonuses
3. **Role Mode** (200) - Strategic choice
4. **Formation** (300) - Tactical positioning
5. **Composition** (350) - Fleet synergy
6. **Synergies** (400) - Cross-system bonuses
7. **Environment** (500) - Terrain effects
8. **Abilities** (600) - Active abilities
9. **Buffs** (700) - Allied buffs
10. **Debuffs** (800) - Enemy debuffs

Higher priority = applied later = can override earlier modifiers.

## When to Use What

### Use `QuickEffectiveShip()`
- UI display of ship stats
- Simple calculations
- Out of combat scenarios
- Don't need modifier stack

### Use `stack.EffectiveShipInFormationV2()`
- Full formation support
- Need to inspect modifiers
- Out of combat with formations

### Use `stack.EffectiveShipInCombat()`
- Combat calculations
- Formation counter important
- In-battle scenarios

### Use `GetModifierBreakdown()`
- Debugging
- UI tooltip showing active modifiers
- Understanding why stats are X

### Use `CompareLoadoutChange()`
- Gem socketing preview
- Equipment comparison UI
- Before/after analysis

### Use `NewModifierBuilder()`
- Custom modifier scenarios
- Complex temporary effects
- Non-standard bonus sources

## Helper Function Cheat Sheet

| Function | Returns | Use For |
|----------|---------|---------|
| `QuickEffectiveShip()` | Ship, Abilities | Simple UI stats |
| `QuickEffectiveShipInCombat()` | Ship, Abilities | Combat damage calc |
| `QuickModifierStack()` | ModifierStack | Just the stack |
| `CompareLoadoutChange()` | Before, After, Diff | Gem preview |
| `CompareFormationChange()` | Before, After, Diff | Formation preview |
| `GetModifierBreakdown()` | []Summary | Debugging UI |
| `GetActiveModifierSources()` | Grouped summaries | Quick active check |
| `RecommendFormation()` | Formation, Score | AI/suggestions |
| `GetStackPowerRating()` | float64 | Matchmaking |

## Examples Repository

See `formation_examples.go` for comprehensive working examples:
- Basic effective ship calculation
- Combat scenarios
- Modifier debugging
- Formation comparisons
- Advanced builder usage

## Common Mistakes

### ❌ Forgetting `time.Now()`
```go
// Won't compile
ship, _ := stack.EffectiveShipV2Simple(Fighter)
```
```go
// ✅ Correct
ship, _ := stack.EffectiveShipV2Simple(Fighter, time.Now())
```

### ❌ Using deprecated functions
```go
// ❌ Deprecated
ship, _ := stack.EffectiveShip(Fighter)
```
```go
// ✅ New
ship, _ := stack.EffectiveShipV2Simple(Fighter, time.Now())
```

### ❌ Missing formation counter in combat
```go
// ❌ No counter bonus
ship, _, _ := stack.EffectiveShipInFormationV2(Fighter, 0, now)
```
```go
// ✅ Includes counter
ship, _, _ := stack.EffectiveShipInCombat(Fighter, 0, enemyFormation, now)
```

### ❌ Ignoring modifier stack
```go
// ❌ Can't debug
ship, abilities := QuickEffectiveShip(...)
// Why is damage 150 instead of 100???
```
```go
// ✅ Can debug
ship, abilities, modStack := stack.EffectiveShipV2(...)
breakdown := modStack.GetSummary(ctx)
// Oh, formation counter is +50%!
```

## Real-World Usage Pattern

```go
// 1. Combat preparation
now := time.Now()
enemyFormation := FormationBox

// 2. Get effective stats
ship, abilities, modStack := stack.EffectiveShipInCombat(
    Fighter,
    0,
    enemyFormation,
    now,
)

// 3. Calculate damage
totalDamage := ship.AttackDamage * shipCount

// 4. Show UI breakdown
breakdown := modStack.GetSummary(ResolveContext{
    Now:            now,
    InCombat:       true,
    HasFormation:   true,
    FormationType:  stack.Formation.Type,
    EnemyFormation: enemyFormation,
})

fmt.Printf("=== Combat Stats ===\n")
fmt.Printf("Damage: %d\n", ship.AttackDamage)
fmt.Printf("\n=== Active Modifiers ===\n")
for _, mod := range breakdown {
    if mod.IsActive {
        fmt.Printf("• %s\n", mod.Description)
    }
}
```

## Next Steps

1. **Read**: `MIGRATION_GUIDE.md` for detailed migration instructions
2. **Study**: `formation_examples.go` for working examples
3. **Reference**: `MODIFIER_SYSTEM.md` for complete documentation
4. **Implement**: Start with simple helper functions, then explore advanced features

## Philosophy

**Old System**: "Combine everything manually and hope you didn't miss anything"

**New System**: "Build transparent layers, resolve automatically, inspect anytime"

The V2 system is designed to be:
- **Easy to use**: Helper functions for common cases
- **Easy to debug**: Full modifier transparency
- **Easy to extend**: Just add new layers
- **Hard to break**: Context-aware activation

You now have a **scalable, maintainable modifier system** instead of a maintenance nightmare! 🎉
