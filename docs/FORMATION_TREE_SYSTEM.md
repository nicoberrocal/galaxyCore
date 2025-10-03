# Formation Mastery Tree System

## Overview

The Formation Mastery Tree system replaces hardcoded formation synergies with a player-driven progression system where players unlock bonuses by spending experience points earned through combat.

## Architecture

### Core Files

1. **formation_tree.go** - Core data structures and state management
   - `FormationTreeNode` - Individual node in the tree
   - `FormationTreeState` - Player's progression state
   - `NodeEffects` - What each node does
   - Experience and reset mechanics

2. **formation_tree_catalog.go** - Central catalog initialization
   - Calls all tree init functions
   - Registry of all formation trees

3. **formation_tree_global.go** - Universal skills tree
   - Applies to all formations
   - Fleet-wide bonuses
   - Strategic flexibility nodes

4. **formation_tree_[formation].go** - Formation-specific trees
   - `formation_tree_line.go`
   - `formation_tree_box.go`
   - `formation_tree_vanguard.go`
   - `formation_tree_skirmish.go`
   - `formation_tree_echelon.go`
   - `formation_tree_phalanx.go`
   - `formation_tree_swarm.go`

5. **formation_tree_compute.go** - V2 compute integration
   - `AddFormationTreeNodes()` - Modifier builder extension
   - `ComputeLoadoutV2WithTree()` - Enhanced compute with trees
   - Counter and reconfig time calculations with tree bonuses

6. **formation_tree_experience.go** - XP system
   - Battle XP calculation
   - Admiral ranks
   - Daily login bonuses
   - Quest rewards

## Key Features

### Strategic Layer (Permanent)
- Players unlock nodes with experience points
- 4 tiers per tree (Basic → Specialization → Advanced → Mastery)
- Mutually exclusive choices create build diversity
- Free resets (1 per month) + paid resets with escalating cost

### Tactical Layer (Dynamic)
- Formation choice still flexible
- Position assignments still customizable
- Role modes still switchable
- Gem loadouts still modifiable

### Integration with V2 Compute

```go
// Old way
mods, grants, _ := ComputeLoadout(ship, role, loadout)

// New way with tree
stack, mods, grants := ComputeLoadoutV2WithTree(
    ship, role, loadout, formation, position, ships,
    treeState, // ← Player's tree progression
    now, inCombat,
)
```

### Experience Gain

```go
result := BattleResult{
    Victory: true,
    EnemyShipsDestroyed: 15,
    DamageDone: 5000,
    FormationUsed: FormationLine,
}

gain := CalculateExperienceGain(result, now)
treeState.AwardExperience(gain)
```

### Unlocking Nodes

```go
node := &FormationTreeCatalog[FormationLine].Nodes[0]
result := treeState.UnlockNode(node, now)

if result.Success {
    // Node unlocked, bonuses now apply
    // Deducted XP automatically
}
```

## Node Effects

### Direct Stat Modifications
- `PositionMods` - Bonuses for specific formation positions
- `FormationMods` - Bonuses for entire formation
- `GlobalMods` - Bonuses for all formations

### Meta Modifiers
- `ReconfigTimeMultiplier` - Speed up formation switching
- `CounterBonusMultiplier` - Enhance formation counters
- `CounterResistMultiplier` - Reduce counter damage taken
- `CompositionBonusMultiplier` - Amplify composition synergies

### Custom Effects
- `CustomEffect` - Special mechanics (e.g., "split_merge", "phantom_decoys")
- `CustomParams` - Parameters for custom effects
- Handled by game logic layer

### Unlocks
- `UnlocksAbility` - Grant new abilities
- `UnlocksFormation` - Make formation types available

## Sample Trees

### Global Tree (Fleet Command Mastery)
- **Tier 1**: Tactical Awareness, Veteran Training, Rapid Deployment
- **Tier 2**: Enhanced Communications, Strategic Vision, Superior Logistics
- **Tier 3**: Supreme Commander, Versatile Genius, Formation Specialist

### Line Formation Tree
- **Tier 1**: Defensive Stance / Offensive Posture / Balanced Deployment
- **Tier 2**: Long-Range Barrage, Shield Wall, Breakthrough Assault
- **Tier 3**: Enfilade Fire, Unbreakable Line, Hammer and Anvil
- **Tier 4**: Supreme Overlord / Master of Combined Arms

### Vanguard Formation Tree
- Focus: Aggressive advance, alpha strikes, overwhelming force
- Ultimate: Unstoppable Assault (glass cannon)

### Box Formation Tree
- Focus: All-around defense, siege resistance
- Ultimate: Impregnable Defense (-40% damage taken)

### Skirmish Formation Tree
- Focus: Mobile warfare, hit-and-run, flanking
- Ultimate: Phantom Fleet (creates decoys)

### Swarm Formation Tree
- Focus: Anti-AoE, overwhelming numbers, death by thousand cuts
- Ultimate: Locust Cloud (untargetable, splash immunity)

## Admiral Ranks

Experience also grants admiral ranks:

| Rank | Title | Total XP Required |
|------|-------|-------------------|
| 0 | Ensign | 0 |
| 1 | Lieutenant | 100 |
| 2 | Commander | 300 |
| 3 | Captain | 700 |
| 4 | Commodore | 1,500 |
| 5 | Rear Admiral | 3,100 |
| 6 | Vice Admiral | 6,300 |
| 7 | Admiral | 12,700 |
| 8 | Fleet Admiral | 25,500 |
| 9 | Grand Admiral | 51,100 |
| 10+ | Supreme Admiral / Legendary Admiral | ... |

## Reset Mechanics

### Free Resets
- 1 free reset granted per month
- Max 3 stored
- Refunds all spent XP

### Paid Resets
- Cost: 1,000 credits × 2^(resets)
- Escalates: 1,000 → 2,000 → 4,000 → 8,000 → ...
- Caps at 1,000,000 credits

## Migration from Old System

### Deprecated Files
- `formation_synergy.go` - Hardcoded ability-position synergies
- Parts of `compute.go` - Replaced by compute_v2.go

### Functions Kept for Compatibility
- `ApplyGemPositionEffects()` - Still used by modifier builder
- `EvaluateCompositionBonuses()` - Still used for fleet synergies
- `FormationTemplatesCatalog` - Pre-made formation configs

### New Workflow

```
Old: Hardcoded synergies → Applied automatically
New: Player earns XP → Unlocks nodes → Bonuses apply
```

## UI Considerations

### Tree Visualization
- Node states: Locked (gray), Available (yellow), Unlocked (green)
- Path highlighting to desired nodes
- Tooltips with full stat breakdowns
- "Preview Build" before committing

### XP Progress
- Current rank and title
- Progress bar to next rank
- Recent XP gains with sources
- Available XP display

### Reset Interface
- Show reset cost (free vs paid)
- Countdown to next free reset
- Confirmation dialog with refund preview

## Balance Considerations

### Power Budget
- Tier 1: ~5% power increase per node
- Tier 2: ~8% power increase
- Tier 3: ~12% power increase
- Tier 4: ~20% power increase
- Max specialized: ~180% baseline
- Max generalist: ~120% baseline

### Anti-Meta Design
- Mutually exclusive nodes prevent "one true build"
- Different formations excel in different scenarios
- Counter system still applies
- No node makes formation invincible (except ultimate choices)

## Example Usage

```go
// Create new player tree state
treeState := NewFormationTreeState(playerID, now)

// Award XP from battles
result := BattleResult{
    Victory: true,
    EnemyShipsDestroyed: 10,
    DamageDone: 3000,
    FormationUsed: FormationLine,
}
gain := CalculateExperienceGain(result, now)
treeState.AwardExperience(gain)

// Unlock a node
lineTree := FormationTreeCatalog[FormationLine]
node := lineTree.Nodes[0] // "Defensive Stance"
result := treeState.UnlockNode(&node, now)

// Compute ship stats with tree bonuses
stack, finalMods, abilities := ComputeLoadoutV2WithTree(
    ship, role, loadout, formation, position, ships,
    treeState, now, inCombat,
)

// Get custom effects for game logic
effects := GetTreeCustomEffects(treeState, FormationSkirmish)
if effects["hit_and_run"] != nil {
    // Apply hit-and-run special mechanic
}
```

## Future Enhancements

1. **Seasonal Trees** - Temporary bonus trees for events
2. **Cross-Formation Synergies** - Unlock nodes that work across formations
3. **Prestiging** - Reset all trees for permanent small bonuses
4. **Leaderboards** - Track unique builds and effectiveness
5. **Balance Patches** - Hot-patch node values without save wipes

## Testing

See `formation_examples.go` for:
- V2 compute workflow examples
- Modifier breakdown examples
- Before/after comparisons
- Integration tests

## Summary

The Formation Tree system transforms formation bonuses from:
- **Static** → **Dynamic**
- **Automatic** → **Player choice**
- **Invisible** → **Visible progression**
- **Uniform** → **Build diversity**

This creates:
- Long-term progression goals
- Meaningful strategic choices
- Emotional investment in builds
- Replayability through different specializations
