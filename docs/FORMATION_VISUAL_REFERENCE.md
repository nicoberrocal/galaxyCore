# Formation Visual Reference

Quick visual guide showing how each formation arranges ships. Each diagram shows the initial layout with position labels.

**Legend:**
- `F` = Front position
- `L` = Flank position
- `B` = Back position
- `S` = Support position

---

## LINE Formation
**Tactical**: Balanced front-back, strong vs frontal, weak to flanking
**Speed**: 1.0x | **Reconfig**: 120s
**Slot Limits**: Front=15, Flank=10, Back=15, Support=8 (Total: 48)

```text
        F
        |
    L---S---L
        |
        B
```

**Expansion**: Front/Back extend in columns, Flanks alternate outward

---

## BOX Formation
**Tactical**: Defensive all-around, even distribution, siege resistant
**Speed**: 0.75x | **Reconfig**: 150s
**Slot Limits**: Front=12, Flank=10, Back=12, Support=10 (Total: 44)

```text
    F-------F
    |       |
    |       |
    L   S   L
    |       |
    |       |
    B-------B
```

**Expansion**: Perimeter expands outward, Support fills center

---

## VANGUARD Formation
**Tactical**: Aggressive spearhead, fast reconfig, high front damage
**Speed**: 1.1x | **Reconfig**: 60s
**Slot Limits**: Front=20, Flank=8, Back=10, Support=6 (Total: 44)

```text
        F
       / \
      /   \
     L     L
    /       \
   B    S    B
```

**Expansion**: Reinforces tip, widens V-shape

---

## SKIRMISH Formation
**Tactical**: Mobile flanking, hit-and-run, very wide
**Speed**: 1.2x | **Reconfig**: 90s
**Slot Limits**: Front=8, Flank=20, Back=12, Support=8 (Total: 48)

```text
L-----------L
 \         /
  \       /
   F  S  F
      |
      B
```

**Expansion**: Flanks extend extremely wide

---

## ECHELON Formation
**Tactical**: Diagonal staggered, asymmetric, concentrated defense
**Speed**: 0.95x | **Reconfig**: 120s
**Slot Limits**: Front=10, Flank=12, Back=10, Support=8 (Total: 40)

```text
            F
           /
          L
         /
        /
       L
      /
     B
    /
   S
```

**Expansion**: Extends along diagonal

---

## PHALANX Formation
**Tactical**: Heavy frontal concentration, extreme flank weakness
**Speed**: 0.8x | **Reconfig**: 180s
**Slot Limits**: Front=25, Flank=6, Back=8, Support=10 (Total: 49)

```text
    F---F---F
    |   |   |
    S       S
    |       |
L---+-------+---L
        |
        B
```

**Expansion**: Front line gets progressively wider

---

## SWARM Formation
**Tactical**: Dispersed anti-AoE, splash resistant
**Speed**: 1.05x | **Reconfig**: 100s
**Slot Limits**: Front=12, Flank=12, Back=12, Support=12 (Total: 48)

```text
        F
       / \
      L   L
     /     \
    B   S   B
```

**Expansion**: Hexagonal rings (6 ships per ring)

---

## Coordinate System Reference

```text
        +Y (Forward/Front)
         |
         |
         |
-X ------+------ +X
(Left)   |      (Right)
         |
         |
         -Y (Back/Rear)
```

---

## Position Slot Limits

Each formation has maximum slot limits per position to maintain:

- **Visual clarity**: Formations remain recognizable and readable
- **Performance**: Frontend can efficiently render ~40-50 total slots
- **Tactical meaning**: Position assignments remain strategically significant

### Limits by Formation

| Formation | Front | Flank | Back | Support | Total |
|-----------|-------|-------|------|---------|-------|
| **Line**     | 15    | 10    | 15   | 8       | 48    |
| **Box**      | 12    | 10    | 12   | 10      | 44    |
| **Vanguard** | 20    | 8     | 10   | 6       | 44    |
| **Skirmish** | 8     | 20    | 12   | 8       | 48    |
| **Echelon**  | 10    | 12    | 10   | 8       | 40    |
| **Phalanx**  | 25    | 6     | 8    | 10      | 49    |
| **Swarm**    | 12    | 12    | 12   | 12      | 48    |

**Note**: Limits reflect tactical focus (e.g., Phalanx emphasizes front with 25 slots, Skirmish emphasizes flanks with 20 slots).

### API Behavior

- `GetNextSlotCoordinate()` returns `false` when position is full
- `GetAllSlotsForPosition()` automatically caps at position limit
- `IsPositionFull()` checks if a position has reached capacity
- `GetMaxSlotsForPosition()` returns the limit for a specific position

---

## Expansion Examples: From Minimal to Full Fleet

### Example 1: LINE Formation Growth

**Stage 1: Minimal (4 ships)**

```text
        F₁
        |
    L₁--S₁--L₂
        |
        B₁
```

**Stage 2: Small Fleet (8 ships)**

```text
        F₂
        |
        F₁
        |
    L₁--S₁--L₂
        |
        B₁
        |
        B₂
```

**Stage 3: Medium Fleet (12 ships)**

```text
        F₃
        |
        F₂
        |
        F₁
        |
L₃--L₁--S₁--L₂--L₄
        |
        B₁
        |
        B₂
        |
        B₃
```

**Stage 4: Large Fleet (16 ships)**

```text
        F₄
        |
        F₃
        |
        F₂
        |
        F₁
        |
L₃--L₁--S₁--L₂--L₄
        |
        B₁
        |
        B₂
        |
        B₃
        |
        B₄
```

**Pattern**: Front/Back extend vertically in columns, Flanks alternate left-right outward, Support stays central.

---

### Example 2: VANGUARD Formation Growth

**Stage 1: Minimal (4 ships)**

```text
        F₁
       / \
      /   \
     L₁   L₂
    /       \
   B₁   S₁   B₂
```

**Stage 2: Small Fleet (8 ships)**

```text
        F₂
        |
        F₁
       / \
      /   \
     L₁   L₂
    /       \
   B₁   S₁   B₂
   |         |
   B₃        B₄
```

**Stage 3: Medium Fleet (12 ships)**

```text
        F₃
        |
        F₂
        |
        F₁
       / \
      /   \
     L₃   L₄
    /       \
   L₁       L₂
  /           \
 B₁   S₁ S₂   B₂
 |             |
 B₃            B₄
```

**Stage 4: Large Fleet (16 ships)**

```text
        F₄
        |
        F₃
        |
        F₂
        |
        F₁
       / \
      /   \
     L₃   L₄
    /       \
   L₁       L₂
  /           \
 B₅           B₆
 |             |
 B₃   S₁ S₂   B₄
 |             |
 B₁            B₂
```

**Pattern**: Front stacks vertically at tip, Flanks widen the V-shape, Back extends along outer edges, Support fills center rear.

---

### Example 3: BOX Formation Growth

**Stage 1: Minimal (4 ships)**

```text
    F₁------F₂
    |       |
    |       |
    L₁  S₁  L₂
    |       |
    |       |
    B₁------B₂
```

**Stage 2: Small Fleet (8 ships)**

```text
    F₂------F₃
    |       |
    F₁      F₄
    |       |
    L₁  S₁  L₂
    |       |
    B₁      B₃
    |       |
    B₂------B₄
```

**Stage 3: Medium Fleet (12 ships)**

```text
    F₃------F₅
    |       |
    F₂      F₆
    |       |
    F₁      F₄
    |       |
    L₁  S₁  L₂
    |   S₂  |
    B₁      B₄
    |       |
    B₂      B₅
    |       |
    B₃------B₆
```

**Stage 4: Large Fleet (16 ships)**

```text
    F₄------F₇
    |       |
    F₃      F₈
    |       |
    F₂      F₆
    |       |
    F₁      F₅
    |       |
    L₁  S₁  L₂
    |   S₂  |
    B₁      B₅
    |       |
    B₂      B₆
    |       |
    B₃      B₇
    |       |
    B₄------B₈
```

**Pattern**: Perimeter expands outward maintaining rectangular shape, Support fills center grid, corners remain anchored.

---

### Example 4: SKIRMISH Formation Growth

**Stage 1: Minimal (4 ships)**

```text
L₁-----------L₂
 \           /
  \         /
   F₁  S₁  F₂
       |
       B₁
```

**Stage 2: Small Fleet (8 ships)**

```text
L₁---------------L₂
 \               /
  \             /
   F₁  S₁ S₂  F₂
       |   |
       B₁  B₂
```

**Stage 3: Medium Fleet (12 ships)**

```text
L₃---L₁-----------L₂---L₄
     \             /
      \           /
       F₁  S₁   F₂
       |   S₂   |
       F₃       F₄
           |
           B₁
           |
           B₂
```

**Stage 4: Large Fleet (16 ships)**

```text
L₅---L₃---L₁-----------L₂---L₄---L₆
          \             /
           \           /
            F₁  S₁   F₂
            |   S₂   |
            F₃       F₄
                |
                B₁
                |
                B₂
                |
                B₃
                |
                B₄
```

**Pattern**: Flanks extend extremely wide horizontally, Front positions spread along forward arc, Back stacks vertically behind Support.

---

## Formation Counter Matrix (Quick Reference)

| Attacker ↓ / Defender → | Line | Box  | Vanguard | Skirmish | Echelon | Phalanx | Swarm |
|-------------------------|------|------|----------|----------|---------|---------|-------|
| **Line**                | 1.0  | 0.8  | **1.3**  | 0.9      | 1.1     | 0.85    | 1.0   |
| **Box**                 | 1.2  | 1.0  | 0.7      | 1.1      | 0.9     | 1.15    | 1.05  |
| **Vanguard**            | 0.7  | **1.3** | 1.0   | **1.4**  | 0.8     | 0.75    | 1.2   |
| **Skirmish**            | 1.1  | 0.9  | 0.6      | 1.0      | 1.2     | **1.3** | 0.95  |
| **Echelon**             | 0.9  | 1.1  | 1.2      | 0.8      | 1.0     | 0.9     | 1.05  |
| **Phalanx**             | 1.15 | 0.85 | **1.25** | 0.7      | 1.1     | 1.0     | 0.8   |
| **Swarm**               | 1.0  | 0.95 | 0.8      | 1.05     | 0.95    | 1.2     | 1.0   |

**Bold** = Strong counter (≥1.25x damage)

---

## Directional Damage Distribution

### Frontal Attack
```text
Front:   60% ████████████
Flank:   20% ████
Back:    10% ██
Support: 10% ██
```

### Flanking Attack
```text
Front:   30% ██████
Flank:   40% ████████
Back:    20% ████
Support: 10% ██
```

### Rear Attack
```text
Front:   10% ██
Flank:   30% ██████
Back:    50% ██████████
Support: 10% ██
```

### Envelopment Attack
```text
Front:   25% █████
Flank:   25% █████
Back:    25% █████
Support: 25% █████
```

---

## Position Bonuses Summary

### LINE
- **Front**: +1 Laser Shield, +10% all damage
- **Flank**: +1 Speed, +5% Crit
- **Back**: +1 Range, +1 Visibility

### BOX
- **All Positions**: +1 to all shields (Laser, Nuclear, Antimatter)

### VANGUARD
- **Front**: +25% all damage, +1 Nuclear Shield
- **Support**: -20% HP, -30% Ability Cooldown

### SKIRMISH
- **Flank**: +2 Speed, +15% Accuracy, +20% all damage
- **Front**: +10% all damage, -10% HP

### ECHELON
- **Front**: +1 Laser Shield, +12% all damage
- **Flank**: +1 Speed, +8% Crit
- **Back**: +1 Range, +5% Accuracy

### PHALANX
- **Front**: +2 Laser/Nuclear Shield, +1 Antimatter Shield, +15% HP, +15% all damage
- **Back**: +2 Range

### SWARM
- **All Positions**: +1 Speed
- **Special**: Dispersed positioning reduces AoE effectiveness

---

## Usage in Frontend

### Scaling Example
```javascript
// Abstract coordinates → Screen pixels
const SCALE = 40; // pixels per unit
const screenX = centerX + (coordinate.x * SCALE);
const screenY = centerY - (coordinate.y * SCALE); // Flip Y
```

### Ship Sprite Selection
```javascript
function getShipSprite(position, shipType) {
  // Different sprites based on position and type
  if (position === 'front') return frontSprites[shipType];
  if (position === 'flank') return flankSprites[shipType];
  if (position === 'back') return backSprites[shipType];
  if (position === 'support') return supportSprites[shipType];
}
```

### Visual Indicators
```javascript
// Highlight initial vs expanded slots
if (slot.isInitial) {
  drawGlowEffect(x, y, 'blue');
} else {
  drawGlowEffect(x, y, 'green');
}
```

---

## Formation Selection Guide

### Offensive Scenarios
- **Against Box**: Use Vanguard (1.3x) or Phalanx (1.25x)
- **Against Phalanx**: Use Skirmish (1.3x) to exploit flanks
- **Against Vanguard**: Use Line (1.3x) or Skirmish (1.4x)

### Defensive Scenarios
- **Expect frontal assault**: Use Phalanx or Line
- **Expect flanking**: Use Box or Swarm
- **Expect siege**: Use Box (siege resistant)
- **Expect AoE**: Use Swarm (splash resistant)

### Mobility Scenarios
- **Hit-and-run**: Use Skirmish (1.2x speed)
- **Fast repositioning**: Use Vanguard (1.1x speed, 60s reconfig)
- **Pursuit**: Use Skirmish or Swarm

### Mixed Fleet Composition
- **Heavy ships (Carriers, Destroyers)**: Front or Flank
- **Fast ships (Scouts, Fighters)**: Flank
- **Long-range (Bombers)**: Back
- **Utility (Drones)**: Support
