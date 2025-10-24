# Combat System Implementation - Deterministic Mechanics

## Overview
Implemented a comprehensive deterministic combat system with type-specific weighted shield mitigation, counter-based mechanics, and bio trait integration for hourly turn-based battles.

## Key Changes

### 1. Combat Counters (stack.go)
- **Added `CombatCounters` struct** to track deterministic combat state
  - `AttackCount`: Total attacks made (for crit timing)
  - `DefenseCount`: Total attacks received (for evasion timing)
  - `LastCritAttack`: Attack number of last crit
- **Integrated into `BattleState`** for persistence across combat rounds

### 2. Enhanced CombatContext (formation_combat.go)
- **Added timestamp field** (`Now time.Time`) for accurate stat calculations
- **Added damage composition tracking** (`AttackerDamageByType map[string]int`)
- **Pre-calculates damage by attack type** for weighted shield application
- Uses `EffectiveShipInCombat()` for all stat lookups

### 3. Asymptotic Shield Mitigation
- **Formula**: `damage / (1 + shieldValue * 0.15)`
- **Characteristics**:
  - Diminishing returns (never reaches 100% mitigation)
  - Shield 3: ~69% damage (31% reduction)
  - Shield 10: ~40% damage (60% reduction)
  - Bio debuffs can reduce shields below 0 (capped at 0 for calculations)

### 4. Weighted Type-Specific Shields
- **Replaced averaged shield mitigation** with proper type-vs-type resolution
- **Each attack type** (Laser/Nuclear/Antimatter) mitigated by corresponding shield
- **Weighted by damage composition** for mixed fleets
- **Example**:
  ```
  Attacker: 200 Laser + 350 Nuclear damage
  Defender: LaserShield=3, NuclearShield=0
  
  Laser: 200 / (1 + 3*0.15) = 138 damage
  Nuclear: 350 / (1 + 0*0.15) = 350 damage (full damage!)
  Total: 488 damage
  ```

### 5. Deterministic Crit System
- **Counter-based, not RNG**
- `CritPct = 0.33` → crit every 3rd attack (interval = 1/0.33)
- `CritPct = 0.50` → crit every 2nd attack
- **Crit damage**: base damage * 1.5 (+50%)
- **Tracked per stack** via `BattleState.Counters.AttackCount`

### 6. Deterministic Evasion System
- **Flat damage reduction, not dodge chance**
- `EvasionPct = 0.35` → 35% damage reduction on ALL incoming damage
- **Capped at 75% reduction** (EvasionPct = 0.75)
- **Stacks additively** from bio traits, formations, gems
- **Applied after shields** in damage pipeline

### 7. First Strike Bonus
- **Triggers on attack counter == 1**
- `FirstVolleyPct = 0.30` → +30% damage on first attack
- **Resets** when battle ends or stack enters cooldown
- **Works with bio traits** like "Predator Pounce"

### 8. Bio Debuff Integration
- **Applied post-combat** via `applyBioDebuffsPostCombat()`
- **Affects next hourly round**
- **Stacks over multiple rounds**
- **Examples**:
  - Shield reduction debuffs (e.g., -1 NuclearShield per stack)
  - Damage over time effects
  - Stat penalties

### 9. Updated ExecuteFormationBattleRound
- **Signature changed**: now requires `time.Time` parameter
- **Initializes combat counters** for both stacks
- **Ticks bio machines** before combat
- **Increments counters** each phase
- **Uses `ComputeStackModifiers()`** for all stat calculations
- **Applies deterministic mechanics**:
  1. First strike bonus (if attack #1)
  2. Crit check (if attack % interval == 0)
  3. Formation counter multiplier
  4. Weighted shield mitigation
  5. Evasion damage reduction
- **Applies bio debuffs** after both phases

## Combat Flow

```
1. Initialize/increment counters
2. Tick bio machines
3. ATTACKER PHASE:
   - For each ship type/bucket:
     * Get effective stats (formation + bio + gems)
     * Calculate base damage
     * Apply first strike bonus (if attack #1)
     * Apply deterministic crit (if attack % interval == 0)
     * Apply formation counter
   - Distribute damage by position
   - Apply weighted shields (type-specific asymptotic)
   - Apply evasion (flat reduction)
   - Update HP buckets
4. DEFENDER RETURNS FIRE (if alive):
   - Same process as attacker
5. Apply bio debuffs for next round
```

## Modified Files

1. **ships/stack.go**
   - Added `CombatCounters` struct
   - Extended `BattleState` with `Counters` field

2. **ships/formation_combat.go**
   - Extended `CombatContext` with `Now` and `AttackerDamageByType`
   - Added `calculateDamageComposition()` method
   - Added `applyAsymptoticShieldMitigation()` helper
   - Replaced `applyShieldMitigation()` with `applyWeightedShieldMitigation()`
   - Completely rewrote `ExecuteFormationBattleRound()` with deterministic mechanics
   - Added `applyBioDebuffsPostCombat()` function

3. **ships/modifiers.go**
   - Added comprehensive documentation header explaining deterministic system
   - Updated `CritPct` documentation: "DETERMINISTIC: crit interval = 1/CritPct"
   - Updated `FirstVolleyPct` documentation: "bonus on attack counter == 1"
   - Updated `EvasionPct` documentation: "DETERMINISTIC: flat % damage reduction (not dodge chance)"

## Breaking Changes

### API Changes
- `ExecuteFormationBattleRound(attacker, defender *ShipStack)` 
  → `ExecuteFormationBattleRound(attacker, defender *ShipStack, now time.Time)`
- `NewCombatContext(attacker, defender *ShipStack)` 
  → `NewCombatContext(attacker, defender *ShipStack, now time.Time)`

### Behavior Changes
- **Crits are now predictable** (every Nth attack, not random)
- **Evasion is damage reduction** (not dodge chance)
- **Shields use asymptotic formula** (not linear percentage)
- **Type-specific shield mitigation** (not averaged)
- **Bio debuffs stack across rounds** (persistent effects)

## Migration Guide

### For Existing Combat Code
```go
// OLD
result := ExecuteFormationBattleRound(attacker, defender)

// NEW
result := ExecuteFormationBattleRound(attacker, defender, time.Now())
```

### For Bio Trait Design
```go
// Traits now work deterministically:

// OLD: "20% chance to crit"
CritPct: 0.20 // Random 20% chance

// NEW: "Crit every 5th attack"
CritPct: 0.20 // Deterministic: 1/0.20 = every 5 attacks

// OLD: "35% dodge chance"
EvasionPct: 0.35 // Random dodge

// NEW: "35% damage reduction"
EvasionPct: 0.35 // Always reduces damage by 35%
```

## Testing Recommendations

1. **Counter persistence**: Verify counters persist across hourly rounds
2. **Crit timing**: Test that crits occur exactly every Nth attack
3. **Evasion stacking**: Verify multiple evasion sources stack additively
4. **Shield type matching**: Confirm Laser damage uses LaserShield, etc.
5. **Bio debuff stacking**: Test multi-round debuff accumulation
6. **First strike**: Verify bonus only applies on attack #1
7. **Asymptotic shields**: Test shield values 0, 3, 5, 10, 20 for correct mitigation

## Performance Considerations

- **Damage composition pre-calculation**: O(N) where N = total ships in attacker stack
- **Weighted shield application**: O(M * T) where M = defender buckets, T = attack types (3)
- **Evasion calculation**: O(M) where M = defender buckets
- **Bio debuff application**: O(B) where B = active bio nodes with outgoing debuffs

All operations are linear and suitable for hourly tick processing.

## Future Enhancements

1. **Per-ship-type counters**: Track crit timing per ship type for more granular control
2. **Counter reset conditions**: Define when counters reset (battle end, cooldown, etc.)
3. **UI indicators**: Show "Next crit in 2 attacks" for player feedback
4. **Bio trait synergies**: Design traits that interact with counter thresholds
5. **Formation-specific counter modifiers**: Formations that alter crit intervals

## Compatibility

- ✅ **Backward compatible** with existing ship blueprints
- ✅ **Works with V2 compute system** (uses `ComputeStackModifiers`)
- ✅ **Integrates with bio machine** (ticks + debuffs)
- ✅ **Compatible with formation system** (position bonuses, counters)
- ✅ **Works with gem system** (modifiers applied via compute pipeline)

## Documentation

All deterministic mechanics are documented in:
- `ships/modifiers.go` (header comment)
- `ships/formation_combat.go` (function comments)
- Individual `StatMods` field comments
