# Compute Architecture - V2 System

## Overview

The compute system has been consolidated into a clean, scalable architecture where all stat modifications flow through the **compute_v2** modifier system. The old `compute.go` and unclear formation synergies have been removed.

## Core Principle

**Everything feeds into compute_v2:**
- Gems create their own `StatMods` via `Gem.Mods`
- Formation tree nodes define `StatMods` in their effects
- Abilities map to `StatMods` via `AbilityEffectsCatalog`
- All sources are combined through the `ModifierBuilder` → `ModifierStack` → `Resolve()` pipeline

## File Structure

### Core Compute Files

- **`compute_v2.go`** - Main compute system with modifier stack resolution
  - `ComputeLoadoutV2()` - Build modifier stack from all sources
  - `ComputeStackModifiers()` - Get modifiers for a ship stack
  - `ComputeEffectiveShipV2()` - Primary entry point for effective stats
  - `ApplyStatModsToShip()` - Apply final mods to ship (migrated from old compute.go)
  - `FilterAbilitiesForMode()` - Filter abilities by role mode (migrated from old compute.go)

- **`modifier_builder.go`** - Fluent API for building modifier stacks
  - `AddGems()` / `AddGemWords()` - Gem system integration
  - `AddRoleMode()` - Role mode bonuses
  - `AddFormationPosition()` - Formation position bonuses
  - `AddFormationCounter()` - Formation matchup bonuses
  - `AddCompositionBonus()` - Fleet composition bonuses
  - `AddGemPositionSynergy()` - Gem-position synergies
  - `AddActiveAbilities()` - Active ability modifiers

- **`modifier_stack.go`** - Modifier layer management and resolution
  - Tracks source, priority, conditions (combat/formation/time)
  - Resolves all layers into final `StatMods`

### Data Catalogs

- **`ability_effects.go`** - Maps abilities to their `StatMods` effects
  - `AbilityEffectsCatalog` - Ability → StatMods mapping
  - `GetAbilityMods()` - Retrieve mods for an ability
  - `AddActiveAbilities()` - Builder method for active abilities

- **`formation_composition.go`** - Fleet composition and gem-position synergies
  - `GemPositionEffectsCatalog` - Gem family + position → bonus mods
  - `CompositionBonusesCatalog` - Fleet composition requirements → bonus mods
  - `FormationTemplatesCatalog` - Pre-configured formation setups
  - Helper functions: `ApplyGemPositionEffects()`, `EvaluateCompositionBonuses()`, `FindBestTemplate()`

- **`formation_tree_compute.go`** - Formation tree integration
  - `AddFormationTreeNodes()` - Apply unlocked tree node bonuses
  - `ComputeLoadoutV2WithTree()` - Extended compute with tree support
  - Tree-enhanced helpers for reconfiguration time, counter multipliers, custom effects

### Supporting Files

- **`gems.go`** - Gems already create their own `Mods` field
- **`formation_tree.go`** - Tree nodes define `StatMods` in `NodeEffects`
- **`stack.go`** - Stack methods now delegate to V2 system

## Data Flow

```
1. Gems → Gem.Mods (defined in gems.go)
2. Formation Tree → NodeEffects.PositionMods/FormationMods/GlobalMods
3. Abilities → AbilityEffectsCatalog[abilityID]
4. Role Mode → RoleModesCatalog[role].BaseMods
5. Formation Position → FormationCatalog[type].PositionBonuses[position]
6. Gem-Position Synergy → GemPositionEffectsCatalog
7. Fleet Composition → CompositionBonusesCatalog

↓ All feed into ↓

ModifierBuilder
  → AddGems()
  → AddRoleMode()
  → AddFormationTreeNodes()
  → AddFormationPosition()
  → AddFormationRoleSynergy()
  → AddGemPositionSynergy()
  → AddCompositionBonus()
  → AddActiveAbilities()
  → Build()

↓

ModifierStack
  → Resolve(context)

↓

Final StatMods
  → ApplyStatModsToShip()

↓

Effective Ship Stats
```

## Key Functions

### Primary Entry Points

```go
// For ship stacks (recommended)
ship, abilities, modStack := ComputeEffectiveShipV2(
    stack, shipType, bucketIndex, now, inCombat, enemyFormation)

// For loadouts (lower level)
modStack, finalMods, grants := ComputeLoadoutV2(
    ship, role, loadout, formation, position, ships, now, inCombat)

// With formation tree support
modStack, finalMods, grants := ComputeLoadoutV2WithTree(
    ship, role, loadout, formation, position, ships, treeState, now, inCombat)
```

### Stack Methods (Convenience)

```go
// V2 methods (recommended)
ship, abilities, modStack := stack.EffectiveShipV2(shipType, bucketIndex, now)
ship, abilities := stack.EffectiveShipV2Simple(shipType, now)
ship, abilities, modStack := stack.EffectiveShipInFormationV2(shipType, bucketIndex, now)

// Deprecated methods (redirect to V2)
ship, abilities := stack.EffectiveShip(shipType) // → EffectiveShipV2Simple
ship, abilities := stack.EffectiveShipInFormation(shipType, bucketIndex) // → EffectiveShipV2Simple
```

## Adding New Modifiers

### 1. For a New Ability

Add to `ability_effects.go`:

```go
AbilityNewAbility: {
    Damage: DamageMods{LaserPct: 0.25},
    SpeedDelta: 2,
},
```

### 2. For a New Gem Effect

Gems define their own mods in `gems.go`:

```go
Gem{
    ID: "new-gem",
    Mods: StatMods{
        AttackRangeDelta: 1,
        AccuracyPct: 0.10,
    },
}
```

### 3. For a New Formation Tree Node

Define in tree catalog files:

```go
FormationTreeNode{
    Effects: NodeEffects{
        PositionMods: map[FormationPosition]StatMods{
            PositionFront: {BucketHPPct: 0.15},
        },
    },
}
```

### 4. For a New Composition Bonus

Add to `formation_composition.go`:

```go
CompositionBonus{
    Type: "New Combo",
    Requirement: map[ShipType]int{Scout: 2, Fighter: 3},
    Bonus: StatMods{SpeedDelta: 2},
}
```

## Removed Systems

### ❌ Deleted Files
- `compute.go` - Replaced by compute_v2.go
- `formation_synergy.go` - Replaced by formation_composition.go

### ❌ Removed Concepts
- **Ability-Formation Position Synergies** - These were unclear and unused. Abilities now have consistent effects regardless of position. Position bonuses come from the formation itself.
- **Manual Modifier Combination** - Old code manually called `CombineMods()`. Now everything goes through the builder pattern.

## Migration Guide

### Old Code
```go
mods, grants, _ := ComputeLoadout(ship, role, loadout)
eff := ApplyStatModsToShip(ship, mods)
abilities := FilterAbilitiesForMode(eff, role, grants)
```

### New Code
```go
ship, abilities, modStack := ComputeEffectiveShipV2(
    stack, shipType, bucketIndex, now, inCombat, enemyFormation)
```

## Benefits

✅ **Single Source of Truth** - All modifiers flow through one system  
✅ **Full Transparency** - `ModifierStack` shows exactly what's active and why  
✅ **Easy Debugging** - `GetModifierBreakdown()` shows all active layers  
✅ **Scalable** - Add new modifier sources without touching core compute logic  
✅ **Type Safe** - Everything uses `StatMods` struct  
✅ **Conditional Logic** - Modifiers can be combat-only, formation-only, or time-limited  

## Testing

All examples in `formation_examples.go` have been updated to use the V2 system and demonstrate:
- Formation combat with counter mechanics
- Gem-position synergies
- Role-formation synergies
- Complete V2 workflow
- Advanced modifier management
