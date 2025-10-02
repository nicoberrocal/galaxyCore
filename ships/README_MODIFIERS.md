# Ship Modifier System Documentation

## Overview

This directory contains a comprehensive **layered modifier system** for managing stat modifications across multiple sources (formations, gems, abilities, roles, etc.). The system prioritizes **transparency**, **scalability**, and **maintainability**.

## 📚 Documentation Files

### Quick Start
- **[MODIFIER_SYSTEM_QUICKSTART.md](MODIFIER_SYSTEM_QUICKSTART.md)** - Start here! 30-second migration guide and common patterns

### Migration
- **[MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)** - Complete migration guide from old system to V2 with examples

### Reference
- **[MODIFIER_SYSTEM.md](MODIFIER_SYSTEM.md)** - Complete technical documentation of the V2 system

### Implementation
- **[SHIP_SYSTEM_OVERVIEW.md](../docs/SHIP_SYSTEM_OVERVIEW.md)** - High-level overview of the entire ship system

## 🎯 Problem Solved

### Before: The Maintenance Nightmare
```go
// Modifiers scattered everywhere
roleMods := RoleModeMods(role, shipType)
gemMods, _, _ := EvaluateGemSockets(gems)
formationMods := formation.ApplyPositionBonusesToShip(position, mods)
mods = ApplyFormationRoleModifiers(mods, formation, position, role)
gemPosMods := ApplyGemPositionEffects(gems, position)
compositionMods, _ := EvaluateCompositionBonuses(ships)

// Manual combining (easy to miss one!)
mods = CombineMods(roleMods, gemMods)
mods = CombineMods(mods, formationMods)
mods = CombineMods(mods, gemPosMods)
mods = CombineMods(mods, compositionMods)

// No idea where each modifier came from
// Can't debug why stats are X
// Hard to add new modifier sources
```

### After: Scalable & Transparent
```go
// Everything automatic
ship, abilities, modStack := stack.EffectiveShipInFormationV2(Fighter, 0, time.Now())

// Full transparency - see exactly what's active
breakdown := modStack.GetSummary(ctx)
for _, mod := range breakdown {
    fmt.Printf("%s: %s (Active: %v)\n", mod.Source, mod.Description, mod.IsActive)
}

// Easy to extend - just add a new layer
builder.AddCustomModifier(...)
```

## 🏗️ Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────┐
│  modifier_stack.go - Core modifier layer types          │
│  - ModifierLayer: Single source of modifiers            │
│  - ModifierStack: Collection of layers                  │
│  - Priority system & conditional activation             │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│  modifier_builder.go - Fluent API for construction      │
│  - NewModifierBuilder()                                 │
│  - AddGems(), AddRoleMode(), AddFormation()...          │
│  - Build() returns complete ModifierStack               │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│  compute_v2.go - High-level compute functions           │
│  - ComputeEffectiveShipV2()                            │
│  - ComputeStackModifiers()                             │
│  - GetModifierBreakdown()                              │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│  compute_helpers.go - Convenient wrapper functions      │
│  - QuickEffectiveShip() - Simple use cases             │
│  - CompareLoadoutChange() - UI previews                │
│  - RecommendFormation() - AI helpers                   │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│  stack.go - ShipStack methods (V2 versions)            │
│  - EffectiveShipV2()                                   │
│  - EffectiveShipInFormationV2()                        │
│  - EffectiveShipInCombat()                             │
└─────────────────────────────────────────────────────────┘
```

### Modifier Sources

| Source | Priority | Lifetime | File |
|--------|----------|----------|------|
| Gems | 100 | Permanent (while equipped) | `gems.go` |
| GemWords | 150 | Permanent (while equipped) | `gems.go` |
| Role Mode | 200 | Semi-permanent (until switched) | `roles.go` |
| Formation Position | 300 | Conditional (while in formation) | `formation.go` |
| Composition | 350 | Dynamic (based on fleet) | `formation_synergy.go` |
| Gem+Position Synergy | 400 | Conditional (gems + formation) | `formation_synergy.go` |
| Formation Counter | 400 | Combat-only | `formation.go` |
| Environment | 500 | Situational | `modifier_builder.go` |
| Abilities | 600 | Temporary (duration) | `abilities.go` |
| Buffs | 700 | Temporary (duration) | Custom |
| Debuffs | 800 | Temporary (duration) | Custom |

## 📖 Quick Reference

### Common Operations

| Task | Function | File |
|------|----------|------|
| Get ship stats for UI | `QuickEffectiveShip()` | `compute_helpers.go` |
| Calculate combat damage | `EffectiveShipInCombat()` | `stack.go` |
| Debug modifiers | `GetModifierBreakdown()` | `compute_v2.go` |
| Preview gem change | `CompareLoadoutChange()` | `compute_helpers.go` |
| Suggest formation | `RecommendFormation()` | `compute_helpers.go` |
| Build custom modifiers | `NewModifierBuilder()` | `modifier_builder.go` |

### File Structure

```
ships/
├── Core Modifier System
│   ├── modifier_stack.go          - Layer types & stack operations
│   ├── modifier_builder.go        - Fluent API for construction
│   ├── compute_v2.go              - V2 compute functions
│   └── compute_helpers.go         - Convenient wrappers
│
├── Legacy System (Deprecated)
│   ├── compute.go                 - Old compute functions (⚠️ deprecated)
│   └── modifiers.go               - StatMods type (still used)
│
├── Modifier Sources
│   ├── gems.go                    - Gem system & synthesis
│   ├── roles.go                   - Role modes
│   ├── formation.go               - Formation types & bonuses
│   ├── formation_synergy.go       - Cross-system synergies
│   └── abilities.go               - Ship abilities
│
├── Ship Types & Combat
│   ├── stack.go                   - ShipStack with V2 methods
│   ├── ship.go                    - Ship struct
│   ├── blueprints.go              - Ship blueprints
│   ├── formation_combat.go        - Combat integration
│   └── formation_examples.go      - Working examples
│
└── Documentation
    ├── README_MODIFIERS.md        - This file
    ├── MODIFIER_SYSTEM_QUICKSTART.md  - Quick start guide
    ├── MIGRATION_GUIDE.md         - Migration from old system
    └── MODIFIER_SYSTEM.md         - Complete technical docs
```

## 🚀 Getting Started

### 1. Read the Quick Start (5 minutes)
```bash
open MODIFIER_SYSTEM_QUICKSTART.md
```
Learn the basics and see common patterns.

### 2. Try the Examples (10 minutes)
```bash
go run formation_examples.go
```
See working code demonstrating V2 usage.

### 3. Migrate Your Code (30 minutes)
```bash
open MIGRATION_GUIDE.md
```
Follow step-by-step migration instructions.

### 4. Deep Dive (optional)
```bash
open MODIFIER_SYSTEM.md
```
Complete technical reference.

## 💡 Key Concepts

### 1. Layered Modifiers
Instead of manually combining modifiers, build layers:
```go
builder := NewModifierBuilder(now)
builder.
    AddGems(gems).
    AddRoleMode(RoleTactical).
    AddFormationPosition(formation, PositionFront)

stack := builder.Build()
finalMods := stack.Resolve(ctx)
```

### 2. Source Tracking
Every modifier knows where it came from:
```go
breakdown := GetModifierBreakdown(...)
for _, mod := range breakdown {
    fmt.Printf("%s: %s\n", mod.Source, mod.Description)
}
// Output:
// gem: Photon III (Socket 1)
// rolemode: Role: Tactical
// formation_position: Line Formation - Front Position
```

### 3. Context-Aware Resolution
Modifiers activate based on context:
```go
// Out of combat - no formation counter
ctx := ResolveContext{InCombat: false}
mods := stack.Resolve(ctx)

// In combat - formation counter active!
ctx := ResolveContext{InCombat: true, EnemyFormation: FormationBox}
mods := stack.Resolve(ctx)
```

### 4. Priority System
Higher priority = applied later:
```go
PriorityGem          = 100  // Base equipment
PriorityRoleMode     = 200  // Strategic choice
PriorityFormation    = 300  // Tactical positioning
PrioritySynergy      = 400  // Cross-system bonuses
PriorityAbility      = 600  // Active abilities
PriorityDebuff       = 800  // Applied last
```

## 📋 Common Patterns

### Pattern: UI Display
```go
func ShowShipStats(stack *ShipStack, shipType ShipType) {
    ship, abilities, _ := stack.EffectiveShipInFormationV2(shipType, 0, time.Now())
    
    fmt.Printf("Attack: %d\n", ship.AttackDamage)
    fmt.Printf("HP: %d\n", ship.HP)
    fmt.Printf("Abilities: %v\n", len(abilities))
}
```

### Pattern: Combat Calculation
```go
func CalculateDamage(attacker, defender *ShipStack) int {
    enemyFormation := defender.Formation.Type
    
    ship, _, _ := attacker.EffectiveShipInCombat(
        Fighter,
        0,
        enemyFormation,
        time.Now(),
    )
    
    return ship.AttackDamage // Includes counter bonus!
}
```

### Pattern: Debugging
```go
func DebugModifiers(stack *ShipStack, shipType ShipType) {
    breakdown := stack.GetModifierBreakdownForShip(
        shipType, 0, time.Now(), false,
    )
    
    for _, mod := range breakdown {
        if mod.IsActive {
            fmt.Printf("✓ %s\n", mod.Description)
        }
    }
}
```

### Pattern: Equipment Preview
```go
func PreviewGem(stack *ShipStack, shipType ShipType, newGem Gem) {
    loadout := stack.GetOrInitLoadout(shipType)
    newLoadout := loadout
    newLoadout.Sockets = append(newLoadout.Sockets, newGem)
    
    before, after, diff := CompareLoadoutChange(
        stack, shipType, 0, newLoadout, time.Now(),
    )
    
    fmt.Printf("Damage: %.2f%% → %.2f%%\n",
        before.Damage.LaserPct * 100,
        after.Damage.LaserPct * 100,
    )
}
```

## 🔧 Extending the System

### Adding a New Modifier Source

1. **Define the source constant** (`modifier_stack.go`):
```go
const (
    // ... existing sources
    SourceTerrain ModifierSource = "terrain"
)
```

2. **Add builder method** (`modifier_builder.go`):
```go
func (mb *ModifierBuilder) AddTerrain(terrainType string, mods StatMods) *ModifierBuilder {
    mb.stack.AddPermanent(
        SourceTerrain,
        terrainType,
        fmt.Sprintf("Terrain: %s", terrainType),
        mods,
        PriorityEnvironment,
        mb.now,
    )
    return mb
}
```

3. **Use it**:
```go
builder := NewModifierBuilder(now)
builder.
    AddGems(gems).
    AddTerrain("nebula", terrainMods).
    AddRoleMode(role)
```

Done! No need to modify existing code.

## ⚠️ Migration Status

### ✅ V2 System (Current, Recommended)
- `compute_v2.go` - Full implementation
- `compute_helpers.go` - Convenient wrappers
- `modifier_stack.go` - Core types
- `modifier_builder.go` - Builder API
- `stack.go` - V2 methods added

### ⚠️ Legacy System (Deprecated)
- `compute.go` - Old functions marked deprecated
- `stack.EffectiveShip()` - Deprecated, use V2
- `stack.EffectiveShipInFormation()` - Deprecated, use V2

### 🔄 Shared Components (Still Used)
- `modifiers.go` - StatMods type (still used by V2)
- `formation.go` - Formation definitions
- `gems.go` - Gem system
- `roles.go` - Role modes

## 📊 Performance

### Typical Modifier Stack
- 10-20 layers per ship
- ~2-4 KB memory per stack
- Negligible CPU overhead
- Acceptable for real-time calculations

### Optimization Tips
1. Cache resolved modifiers when context doesn't change
2. Clean up expired layers periodically
3. Use helper functions for simple cases
4. Reuse builder instances for batch operations

## 🐛 Troubleshooting

### Stats seem wrong?
✅ Use `GetModifierBreakdown()` to see what's active

### Formation counter not applying?
✅ Use `EffectiveShipInCombat()` with enemy formation

### Too verbose?
✅ Use helper functions like `QuickEffectiveShip()`

### Need custom modifiers?
✅ Use `NewModifierBuilder()` for full control

### Migration errors?
✅ See `MIGRATION_GUIDE.md` for detailed examples

## 🎓 Learning Path

1. **Beginner**: Read `MODIFIER_SYSTEM_QUICKSTART.md`
2. **Intermediate**: Study `formation_examples.go`
3. **Advanced**: Read `MODIFIER_SYSTEM.md`
4. **Expert**: Extend with custom modifiers

## 🤝 Contributing

When adding new modifier sources:
1. Add to appropriate file (gems, formations, etc.)
2. Create builder method in `modifier_builder.go`
3. Add priority constant if needed
4. Update documentation
5. Add example to `formation_examples.go`

## 📝 Summary

The V2 modifier system transforms stat modification from a **maintenance nightmare** into a **scalable, transparent, and debuggable** system.

**Key Benefits**:
- ✅ Source tracking - know where every modifier comes from
- ✅ Easy debugging - inspect all active modifiers
- ✅ Scalable - add new sources without touching existing code
- ✅ Context-aware - modifiers activate when appropriate
- ✅ Maintainable - clean separation of concerns

**Next Steps**:
1. Start with `MODIFIER_SYSTEM_QUICKSTART.md`
2. Migrate using `MIGRATION_GUIDE.md`
3. Explore examples in `formation_examples.go`
4. Refer to `MODIFIER_SYSTEM.md` for details

Happy modifying! 🚀
