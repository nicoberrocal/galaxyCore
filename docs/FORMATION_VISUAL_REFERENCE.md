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
