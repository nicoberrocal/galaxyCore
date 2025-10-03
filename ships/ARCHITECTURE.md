# Ships System Architecture

## Clean System Separation Principle

The ships system follows a **clean separation of concerns** where each subsystem generates its own `StatMods` independently and feeds them to the modifier stack. There are **no cross-system synergies**.

### Core Principle

```
Each System → StatMods → ModifierStack → Resolved Stats
```

**No synergies like:**
- ❌ Gem + Position synergy
- ❌ Formation + Role synergy  
- ❌ Entity1 + Entity2 synergy

**Instead:**
- ✅ Each system provides its own bonuses
- ✅ Player decisions are clear and predictable
- ✅ Systems can be balanced independently

---

## System Breakdown

### 1. **Gem System** (`gems.go`)
**Provides:** StatMods from gem properties only

- Each gem has intrinsic `Mods` based on family and tier
- GemWords activate when socket sequences match
- Hybrid/Relic gems have combined properties
- **No position-based bonuses**

```go
// Gems generate their own mods
mods, grants, gemWords := EvaluateGemSockets(loadout.Sockets)
builder.AddGems(loadout.Sockets)
builder.AddGemWords(gemWords)
```

---

### 2. **Role Mode System** (`roles.go`)
**Provides:** StatMods from active role mode

- Tactical, Economic, Recon, Scientific modes
- Each role has `BaseMods` that apply globally
- Roles enable/disable certain abilities
- **No formation interaction**

```go
// Role provides its own mods
spec := RoleModesCatalog[role]
builder.AddRoleMode(role) // Adds spec.BaseMods
```

---

### 3. **Formation System** (`formation.go`, `formation_tree.go`)
**Provides:** StatMods from FormationCatalog + active tree nodes

#### Base Formation Bonuses
- Each `FormationType` has position-specific bonuses in `FormationCatalog`
- Example: Line formation gives front position +10% damage
- **No role or gem interaction**

```go
// Formation provides position bonuses from catalog
spec := FormationCatalog[formation.Type]
posBonus := spec.PositionBonuses[position]
builder.AddFormationPosition(formation, position)
```

#### Formation Tree Bonuses
- Global tree: applies to all formations
- Formation-specific trees: apply when using that formation
- Nodes can provide:
  - `PositionMods`: bonuses to specific positions
  - `FormationMods`: bonuses to entire formation
  - `GlobalMods`: always-active bonuses
  - Meta modifiers: reconfig time, counter bonuses, etc.

```go
// Tree nodes provide additional mods
if treeState != nil {
    builder.AddFormationTreeNodes(treeState, formation.Type)
}
```

---

### 4. **Ability System** (`abilities.go`, `ability_effects.go`)
**Provides:** StatMods when abilities are active

- Each `AbilityID` maps to `StatMods` in `AbilityEffectsCatalog`
- Passive abilities: always active
- Active abilities: temporary duration
- Toggle abilities: active until toggled off

```go
// Abilities provide their own mods when active
mods := GetAbilityMods(abilityID)
builder.AddAbility(abilityID, mods, duration)
```

---

### 5. **Environmental/State Systems**
**Provides:** StatMods from game state

- **Anchored state**: penalty when mining
- **Formation counters**: rock-paper-scissors matchups
- **Environmental effects**: nebula, asteroid fields, etc.

```go
builder.AddAnchoredPenalty(anchored)
builder.AddFormationCounter(attackerFormation, defenderFormation, inCombat)
builder.AddEnvironment(envID, description, mods)
```

---

## Modifier Stack Flow

### Building the Stack

```go
builder := NewModifierBuilder(now)

// 1. Gems
builder.AddGemsFromLoadout(loadout)

// 2. Role Mode
builder.AddRoleMode(role)

// 3. Formation Tree (if unlocked nodes exist)
if treeState != nil {
    builder.AddFormationTreeNodes(treeState, formation.Type)
}

// 4. Formation Position
if formation != nil {
    builder.AddFormationPosition(formation, position)
}

// 5. Abilities (active only)
for _, ability := range activeAbilities {
    mods := GetAbilityMods(ability.ID)
    builder.AddAbility(ability.ID, mods, duration)
}

// 7. State modifiers
builder.AddAnchoredPenalty(anchored)

stack := builder.Build()
```

### Resolving the Stack

```go
ctx := ResolveContext{
    Now:          time.Now(),
    InCombat:     true,
    HasFormation: true,
    FormationType: FormationLine,
}

finalMods := stack.Resolve(ctx)
effectiveShip := ApplyStatModsToShip(blueprint, finalMods)
```

---

## Priority System

Modifiers are applied in priority order:

1. **PriorityGem** (100): Gem bonuses
2. **PriorityGemWord** (150): GemWord bonuses
3. **PriorityRoleMode** (200): Role mode bonuses
4. **PriorityFormation** (300): Formation position bonuses
5. **PriorityFormation+50** (350): Formation tree node bonuses
6. **PrioritySynergy** (400): ~~Deprecated synergies~~
7. **PriorityComposition** (500): Fleet composition bonuses
8. **PriorityAbility** (600): Active abilities
9. **PriorityBuff** (700): Temporary buffs
10. **PriorityDebuff** (800): Debuffs
11. **PriorityEnvironment** (900): Environmental effects

---

## Deprecated Systems

The following systems have been **deprecated** for clean separation:

### ❌ Gem-Position Synergy

- **Old:** Gems gave extra bonuses based on formation position
- **Problem:** Mixed gem and formation concepts
- **Solution:** Gems provide fixed mods, formations provide fixed mods

### ❌ Formation-Role Synergy

- **Old:** Role modes enhanced formation bonuses
- **Problem:** Made it unclear what each system contributed
- **Solution:** Each system provides independent bonuses

### ❌ Fleet Composition Bonuses

- **Old:** Bonuses activated when fleet had specific ship type combinations
- **Problem:** Created implicit synergies between ship types
- **Solution:** Each ship contributes independently, no "combo" bonuses

### Backward Compatibility

Deprecated functions are kept but do nothing:

- `AddFormationRoleSynergy()` → no-op
- `AddGemPositionSynergy()` → no-op
- `AddCompositionBonus()` → no-op
- `ApplyFormationRoleModifiers()` → returns only formation mods
- `ApplyGemPositionEffects()` → returns zero mods
- `EvaluateCompositionBonuses()` → returns zero mods

---

## Benefits of Clean Separation

### 1. **Clear Player Decisions**
Players can understand exactly what each choice provides:
- "This gem gives +10% laser damage"
- "This formation position gives +1 shield"
- "This role mode gives +2 speed"

### 2. **Independent Balancing**
Each system can be tuned without affecting others:
- Buff gems → doesn't affect formations
- Nerf a formation → doesn't affect gems

### 3. **Easier Testing**
Each system can be tested in isolation:
```go
// Test gems alone
mods, _, _ := EvaluateGemSockets(gems)

// Test formation alone
posBonus := spec.PositionBonuses[position]

// Test role alone
roleMods := RoleModesCatalog[role].BaseMods
```

### 4. **Simpler Code**
No complex synergy calculations or lookup tables:
```go
// Old (complex)
synergy := GemPositionEffectsCatalog.Lookup(gem, position)

// New (simple)
mods := gem.Mods
```

### 5. **Better UI/UX**
Tooltips can show exact contributions:
```
Total Damage: +35%
  Gems: +15%
  Formation: +10%
  Role Mode: +5%
  Abilities: +5%
```

---

## Future Extensions

If synergies are needed in the future, they should be **within** a system:

### ✅ Formation Tree Synergies
```go
// Node that unlocks when you have both Line and Box formations mastered
node := FormationTreeNode{
    Requirements: NodeRequirements{
        RequiredNodes: []string{"line_mastery", "box_mastery"},
    },
    Effects: NodeEffects{
        GlobalMods: StatMods{...}, // Applies to all formations
    },
}
```

### ✅ Gem Synthesis Synergies
```go
// Hybrid gems that combine two families
hybrid := SynthesizeGems(laserGem, nuclearGem)
// Result has combined properties, not position-dependent
```

### ✅ Ability Combos
```go
// Abilities that enhance each other when both active
if HasAbility(AlphaStrike) && HasAbility(TargetingUplink) {
    // Handled within ability system
}
```

---

## Migration Guide

If you have code using deprecated synergies:

### Before

```go
builder.AddFormationRoleSynergy(formation, position, role)
builder.AddGemPositionSynergy(gems, position)
builder.AddCompositionBonus(ships)
```

### After

```go
// Just remove these calls - bonuses now come from:
// - Formation tree nodes (if you want formation-specific bonuses)
// - Role mode BaseMods (for role bonuses)
// - Gem Mods (for gem bonuses)
// - Each ship contributes independently (no composition combos)
```

### Adding Formation-Specific Bonuses

Instead of synergies, use formation tree nodes:

```go
// Create a node in the formation tree
node := FormationTreeNode{
    Formation: FormationLine,
    Effects: NodeEffects{
        PositionMods: map[FormationPosition]StatMods{
            PositionFront: {
                Damage: DamageMods{LaserPct: 0.10},
            },
        },
    },
}
```

---

## Summary

**One Rule:** Each system generates its own StatMods independently.

**No cross-system synergies.** If synergies exist, they live within a single system (e.g., formation tree nodes that require multiple formations unlocked).

This creates a **predictable, testable, and maintainable** architecture where player choices are clear and systems can evolve independently.
