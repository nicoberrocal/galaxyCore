# Cosmic Synthesis Gem System

## Overview

A sophisticated gem system combining Diablo-style upgrading with a rich synthesis system. Gems represent crystallized cosmic energies that can be socketed into ships to enhance their capabilities.

## Core Concepts

### Gem Families

Thematic groups that define a gem's core identity and effects:

- **Laser**: Focused energy weapons and precision systems
- **Nuclear**: High-yield explosive and reactor systems
- **Antimatter**: Exotic matter manipulation and energy conversion
- **Kinetic**: Physical impact and projectile systems
- **Sensor**: Detection, scanning, and targeting systems
- **Warp**: FTL travel and spatial manipulation
- **Engineering**: Ship systems and structural integrity
- **Logistics**: Support, repair, and resource management

### Tiers

Gems range from Tier 1 (weakest) to Tier 5 (most powerful). Higher tiers provide stronger effects but are harder to obtain.

### Kinds

- **Pure**: Base family gems with straightforward effects
- **Hybrid**: Combinations of 2 families with unique blended effects
- **Relic**: Rare combinations of 3+ families with powerful, specialized effects

### Origins

Astrophysical sources that influence gem properties:

- **Pulsar**: Precision and energy focus
- **Supernova**: Raw power and explosive potential
- **Wormhole**: Spatial manipulation and FTL effects
- **Singularity**: Exotic matter and extreme physics
- **Nebula**: System stability and efficiency
- **Quasar**: High-energy emissions and detection
- **Big Bang Remnant**: Primordial cosmic energies

## Gem Synthesis Tree

### Pure Gem Progression

```text
Tier 1: [Laser I]     [Nuclear I]    [Antimatter I]  [Kinetic I]    [Sensor I]     [Warp I]       [Engineering I] [Logistics I]
           |             |              |               |              |              |               |                |
Tier 2: [Laser II]    [Nuclear II]   [Antimatter II] [Kinetic II]   [Sensor II]    [Warp II]      [Engineering II][Logistics II]
           |             |              |               |              |              |               |                |
Tier 3: [Laser III]   [Nuclear III]  [Antimatter III][Kinetic III]  [Sensor III]   [Warp III]     [Engineering III][Logistics III]
           |             |              |               |              |              |               |                |
Tier 4: [Laser IV]    [Nuclear IV]   [Antimatter IV] [Kinetic IV]   [Sensor IV]    [Warp IV]      [Engineering IV][Logistics IV]
           |             |              |               |              |              |               |                |
Tier 5: [Laser V]     [Nuclear V]    [Antimatter V]  [Kinetic V]    [Sensor V]     [Warp V]       [Engineering V] [Logistics V]
```

### Hybrid Synthesis Paths

```text
Laser I + Nuclear I     → [Photon Core I] (Laser/Nuclear hybrid)
Laser I + Sensor I      → [Optical Matrix I] (Laser/Sensor hybrid)
Nuclear I + Kinetic I   → [Fission Driver I] (Nuclear/Kinetic hybrid)
Sensor I + Warp I       → [Quantum Echo I] (Sensor/Warp hybrid)
Engineering I + Logistics I → [Harmonic Core I] (Engineering/Logistics hybrid)
Antimatter I + Laser I  → [Singularity Focus I] (Antimatter/Laser hybrid)
Warp I + Engineering I  → [Phase Regulator I] (Warp/Engineering hybrid)
```

### Advanced Hybrid Evolution (Tier 3+)

```text
Photon Core III + Fission Driver III → [Plasma Annihilator I] (Laser/Nuclear/Kinetic)
Optical Matrix III + Quantum Echo III → [Clairvoyance Engine I] (Laser/Sensor/Warp)
Singularity Focus III + Phase Regulator III → [Void Rift Generator I] (Antimatter/Laser/Warp)
Harmonic Core III + Quantum Echo III → [Omniversal Matrix I] (Engineering/Logistics/Sensor)
```

### Legendary Relics (Tier 5)

```text
[Plasma Annihilator V] + [Clairvoyance Engine V] + [Void Rift Generator V]
                     \          |          /
                      \         |         /
                       [Cosmic Singularity Core] (Legendary Relic)
```

### Naming Conventions

- **Pure Gems**: [Family Name] [Roman Numeral] (e.g., "Laser III")
- **Hybrid Gems**: [Descriptive Name] [Roman Numeral] (e.g., "Photon Core II")
- **Relics**: [Unique Name] [Roman Numeral] (e.g., "Cosmic Singularity Core")

## Core Systems

### 1. Pure Upgrading

- Combine 3 identical gems → 1 gem of next tier (up to tier 5)
- Preserves family and origin
- 100% success rate
- Example: 3x Laser I → 1x Laser II

### 2. Hybrid Synthesis

- Combine 2 different families → New hybrid gem
- Success depends on affinity between families
- Output tier based on input energy and efficiency
- May accumulate instability
- Example: Laser I + Nuclear I → Hybrid Laser-Nuclear I

### 3. Relic Creation

- Combine 3 specific high-tier hybrids → Powerful relic
- High risk/reward with special effects
- Example: 3x T5 hybrids → "Singularity Core"

### 4. Energy Model

- Each gem has an energy value based on tier and origin
- Output energy = sum(input energy) * efficiency
- Excess energy increases success chance
- Energy requirements scale non-linearly with tier

### 5. Affinity System

Families have natural affinities that affect synthesis:

| Family 1    | Family 2    | Affinity | Base Success | Efficiency | Suggested Origin  |
|-------------|-------------|----------|--------------|------------|-------------------|
| Laser      | Nuclear     | 0.7      | 85%          | 80%        | Supernova         |
| Nuclear    | Kinetic     | 0.6      | 85%          | 80%        | Supernova         |
| Laser      | Sensor      | 0.55     | 80%          | 75%        | Quasar            |
| Sensor     | Warp        | 0.5      | 75%          | 70%        | Wormhole          |
| Engineering| Logistics   | 0.65     | 85%          | 80%        | Nebula            |
| Antimatter | Kinetic     | -0.5     | 50%          | 45%        | Singularity        |
| Antimatter | Laser       | 0.3      | 70%          | 65%        | Singularity        |
| Sensor     | Antimatter  | -0.2     | 60%          | 60%        | Quasar            |
| Warp       | Engineering | 0.4      | 75%          | 70%        | Wormhole          |
| Warp       | Logistics   | 0.35     | 75%          | 70%        | Nebula            |

### 6. Gem Words

Special sequences that unlock bonus effects when socketed in order:

1. **Photon Overcharge**

   - Sequence: [Laser, Laser, Nuclear]
   - Minimum Tier: 3
   - Effects: +30% Laser Damage, +15% Crit Chance
   - Grants: "Photon Surge" ability

2. **Quantum Entanglement**
   - Sequence: [Sensor, Warp, Antimatter]
   - Minimum Tier: 4
   - Effects: +25% Scan Range, +20% Warp Speed
   - Grants: "Quantum Echo" ability

3. **Singularity Core**
   - Sequence: [Hybrid T5, Hybrid T5, Hybrid T5]
   - Minimum Tier: 5
   - Effects: +50% All Damage, -30% Ability Cooldown
   - Grants: "Singularity Pulse" ability

## Stat Modifications

### Per-Tier Scaling
Each gem family provides different stat bonuses that scale with tier:

#### Laser
- Tier 1: +5% Laser Damage
- Tier 2: +10% Laser Damage, +2% Crit Chance
- Tier 3: +15% Laser Damage, +4% Crit Chance, +5% Range
- Tier 4: +20% Laser Damage, +6% Crit Chance, +10% Range
- Tier 5: +25% Laser Damage, +8% Crit Chance, +15% Range, +5% Crit Damage

#### Nuclear
- Tier 1: +5% Nuclear Damage
- Tier 2: +10% Nuclear Damage, +5% AOE Radius
- Tier 3: +15% Nuclear Damage, +10% AOE Radius, +5% Armor Pen
- Tier 4: +20% Nuclear Damage, +15% AOE Radius, +10% Armor Pen
- Tier 5: +25% Nuclear Damage, +20% AOE Radius, +15% Armor Pen, +10% Shield Pen

#### Sensor
- Tier 1: +5% Scan Range
- Tier 2: +10% Scan Range, +5% Detection
- Tier 3: +15% Scan Range, +10% Detection, +5% Accuracy
- Tier 4: +20% Scan Range, +15% Detection, +10% Accuracy
- Tier 5: +25% Scan Range, +20% Detection, +15% Accuracy, +10% Crit Chance

#### Warp
- Tier 1: +5% Warp Speed
- Tier 2: +10% Warp Speed, +5% Fuel Efficiency
- Tier 3: +15% Warp Speed, +10% Fuel Efficiency, +5% Cooldown Reduction
- Tier 4: +20% Warp Speed, +15% Fuel Efficiency, +10% Cooldown Reduction
- Tier 5: +25% Warp Speed, +20% Fuel Efficiency, +15% Cooldown Reduction, +5% Evasion

## Synthesis Mechanics

### Success Chance Calculation
```
base_chance = affinity.BaseSuccess
tier_bonus = 0.05 * (avg_input_tier - output_tier)
instability_penalty = 0.5 * avg_instability
final_chance = base_chance + tier_bonus - instability_penalty
final_chance = clamp(final_chance, 0.05, 0.95)
```

### Energy Calculation
```
input_energy = sum(energy_of_each_gem)
output_energy = input_energy * affinity.Efficiency
output_tier = tier_from_energy(output_energy)
```

### Instability
- Increases with failed synthesis attempts
- Reduces success chance
- Can be mitigated by using higher-tier gems or special stabilizers
- Resets to 0 on successful synthesis

## Advanced Strategies

### Optimal Synthesis Paths
1. **Early Game**: Focus on pure upgrades for your primary weapon systems
2. **Mid Game**: Experiment with high-affinity hybrid combinations
3. **Late Game**: Aim for specific Gem Words and Relics

### Risk Management
- Always check affinity before attempting synthesis
- Balance between pushing for higher tiers and maintaining stability
- Keep backup gems in case of synthesis failure

### Gem Word Optimization
- Plan socket layouts in advance
- Consider using lower-tier gems to complete powerful Gem Words
- Some Gem Words have synergistic effects when combined

## Example Builds

### 1. Sniper Build
- **Gems**: 3x Laser V
- **Gem Word**: Photon Overcharge
- **Effects**: +75% Laser Damage, +24% Crit Chance, +15% Range, +15% Crit Damage

### 2. Scout Build
- **Gems**: Sensor V, Warp V, Engineering V
- **Gem Word**: Quantum Entanglement
- **Effects**: +25% Scan Range, +20% Warp Speed, +15% System Efficiency

### 3. Endgame Hybrid
- **Gems**: Singularity Core (3x T5 Hybrid)
- **Effects**: +50% All Damage, -30% Ability Cooldown, +25% All Resistances

## Conclusion
The Cosmic Synthesis Gem System offers deep customization and progression through its combination of pure upgrades, hybrid synthesis, and powerful Gem Words. By understanding the affinity system and carefully planning your synthesis paths, you can create powerful combinations that complement your playstyle.

Remember to balance risk versus reward, and don't be afraid to experiment with different combinations to discover powerful synergies!

### 9.1 Performance Impact Mitigation

- Batch formation updates per tick
- Cache formation calculations
- Limit formation changes per time window
- Use efficient data structures for bucket operations