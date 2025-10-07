# Galaxy Core Systems Documentation

## Table of Contents
1. [Ship Types](#ship-types)
2. [Formations](#formations)
3. [Formation Trees](#formation-trees)
4. [Gems System](#gems-system)
5. [Role Modes](#role-modes)

## Ship Types

### Drone
- **Role**: Economic unit for resource gathering
- **Attack Type**: Laser
- **Shields**: Laser(2), Nuclear(1), Antimatter(0)
- **Abilities**:
  - Resource Harvester: Gather resources while anchored
  - Self-Repair: Regenerate HP when out of combat
  - Cloak While Anchored: Become invisible while gathering

### Scout
- **Role**: Recon/light skirmisher
- **Attack Type**: Laser
- **Shields**: Laser(2), Nuclear(0), Antimatter(1)
- **Abilities**:
  - Long-Range Sensors: Increased visibility
  - Ping: Reveal and mark targets
  - Decoy Beacon: Create phantom contacts

### Fighter
- **Role**: Versatile combat unit
- **Attack Type**: Laser
- **Shields**: Laser(3), Nuclear(1), Antimatter(1)
- **Abilities**:
  - Adaptive Targeting: Change attack type
  - Focus Fire: Bonus damage to marked targets
  - Evasive Maneuvers: Temporary defense boost

### Bomber
- **Role**: Siege platform
- **Attack Type**: Nuclear
- **Shields**: Laser(1), Nuclear(3), Antimatter(2)
- **Abilities**:
  - Light-Speed: Enable warp travel
  - Siege Payload: Bonus vs structures
  - Standoff Pattern: Increased range, slower rate of fire

### Carrier
- **Role**: Mobile support hub
- **Attack Type**: Nuclear
- **Shields**: Laser(2), Nuclear(4), Antimatter(3)
- **Abilities**:
  - Hangar Launch: Deploy/recall escorts
  - Point-Defense Screen: Area defense against lasers
  - Can transport smaller ships

### Destroyer
- **Role**: Heavy combatant
- **Attack Type**: Antimatter
- **Shields**: Laser(2), Nuclear(2), Antimatter(4)
- **Abilities**:
  - Alpha Strike: High-damage opening attack
  - Interdictor Pulse: Block enemy warp
  - Heavy weapons platform

### Ballista (formerly Artillery)
- **Role**: Long-range artillery
- **Attack Type**: Nuclear
- **Shields**: Laser(2), Nuclear(4), Antimatter(1)
- **Speed**: 3
- **Visibility Range**: 7
- **Attack Range**: 7
- **HP**: 350
- **Attack Damage**: 50
- **Attack Interval**: 3.5s
- **Abilities**:
  - Light-Speed: Enable warp travel
  - Cluster Munitions: Area damage to multiple targets
  - Barrage Mode: Increased rate of fire
- **Combos**:
  - Cluster Munitions + Barrage Mode = Massive AoE coverage
  - Suppressive Fire + Standoff Pattern = Area lockdown
  - Cluster + Swarm formation = Anti-cluster defense

### Ghost (formerly Stealth Frigate)
- **Role**: Stealth assassin
- **Attack Type**: Laser
- **Shields**: Laser(2), Nuclear(0), Antimatter(2)
- **Speed**: 7
- **Visibility Range**: 5
- **Attack Range**: 3
- **HP**: 180
- **Attack Damage**: 35
- **Attack Interval**: 1.5s
- **Abilities**:
  - Active Camo: Temporary stealth
  - Backstab: Bonus damage from behind
  - Smoke Screen: Area stealth field
- **Combos**:
  - Active Camo + Backstab = Devastating backline attacks
  - Smoke Screen + Active Camo = Team stealth
  - Backstab + Flank position = Position-based assassination

### Cruiser
- **Role**: Versatile combatant
- **Attack Type**: Laser
- **Shields**: Laser(3), Nuclear(2), Antimatter(2)
- **Speed**: 4
- **Visibility Range**: 6
- **Attack Range**: 4
- **HP**: 450
- **Attack Damage**: 40
- **Attack Interval**: 2.0s
- **Abilities**:
  - Light-Speed: Enable warp travel
  - Point Defense: Intercept incoming projectiles
  - Tactical Maneuvering: Temporary speed and evasion boost

### Corvette
- **Role**: Fast attack
- **Attack Type**: Laser
- **Shields**: Laser(2), Nuclear(1), Antimatter(1)
- **Speed**: 6
- **Visibility Range**: 5
- **Attack Range**: 3
- **HP**: 200
- **Attack Damage**: 25
- **Attack Interval**: 1.2s
- **Abilities**:
  - Light-Speed: Enable warp travel
  - Hit and Run: Bonus damage on first strike
  - Evasive Pattern: Increased dodge chance

### Ghost
- **Role**: Stealth assassin
- **Attack Type**: Laser
- **Shields**: Laser(2), Nuclear(0), Antimatter(2)
- **Speed**: 7
- **Visibility Range**: 5
- **Attack Range**: 3
- **HP**: 180
- **Attack Damage**: 35
- **Attack Interval**: 1.5s
- **Abilities**:
  - Active Camo: Temporary stealth
  - Backstab: Bonus damage from behind
  - Smoke Screen: Area stealth field
- **Combos**:
  - Active Camo + Backstab: Devastating backline attacks
  - Smoke Screen + Active Camo: Team stealth
  - Backstab + Flank position: Position-based assassination

## Formations

### Line Formation
- **Type**: Balanced front-back arrangement
- **Speed Multiplier**: 1.0x
- **Bonuses**:
  - Front: +1 Laser Shield, +10% damage
  - Flank: +1 Speed, +5% crit chance
  - Back: +1 Attack Range, +1 Visibility
- **Special**: Strong vs frontal attacks, weak to flanking

### Box Formation
- **Type**: Defensive all-around
- **Speed Multiplier**: 0.75x
- **Bonuses**:
  - All positions: +1 to all shields
- **Special**: Even damage distribution, excellent vs siege

### Vanguard Formation
- **Type**: Aggressive forward deployment
- **Speed Multiplier**: 1.1x
- **Bonuses**:
  - Front: +25% damage, +1 Nuclear Shield
  - Support: -20% HP, -30% ability cooldown
- **Special**: Fast reconfiguration, high front damage

### Skirmish Formation
- **Type**: Mobile flanking
- **Speed Multiplier**: 1.2x
- **Bonuses**:
  - Flank: +2 Speed, +15% accuracy, +20% damage
  - Front: +10% damage, -10% HP
- **Special**: Hit-and-run tactics, mobile combat

### Phalanx Formation
- **Type**: Heavy frontal assault
- **Speed Multiplier**: 0.8x
- **Bonuses**:
  - Front: +3 to all shields, +15% HP, +10% damage
  - Back: +1 range, +5% accuracy
- **Special**: Excels in frontal engagements, weak to flanking

### Echelon Formation
- **Type**: Diagonal staggered lines
- **Speed Multiplier**: 1.0x
- **Bonuses**:
  - Front/Flank: +1 shield, +10% damage
  - Diagonal attacks: +15% accuracy
- **Special**: Flexible positioning, strong against single targets

### Swarm Formation
- **Type**: Overwhelming numbers
- **Speed Multiplier**: 1.3x
- **Bonuses**:
  - +5% damage per adjacent friendly ship
  - +10% dodge chance
- **Special**: Strength in numbers, weak to area attacks

## Formation Trees

### Fleet Command Mastery (Global Tree)
```
Tier 1 (Choose 1):
  [Tactical Awareness] - +1 visibility to all ships
  [Veteran Training] - +5% HP, +3% accuracy
  [Rapid Deployment] - -20% formation reconfig time

Tier 2 (Requires 1 from Tier 1):
  [Enhanced Communications] - +10% composition bonus, +1 ability range
    └─ Requires: Tactical Awareness
  [Strategic Vision] - Reveal enemy composition
    └─ Requires: Tactical Awareness
  [Superior Logistics] - -5% upkeep, +10% transport capacity
    └─ Requires: Veteran Training
  [Adaptive Tactics] - Reposition in combat once
    └─ Requires: Rapid Deployment

Tier 3 (Requires 2 from Tier 2, min 5 nodes total):
  [Supreme Commander] - +15% formation bonuses, +5% damage
    └─ Requires: Enhanced Communications, Strategic Vision
  [Versatile Genius] - Run 2 formations simultaneously
    └─ Requires: Strategic Vision, Adaptive Tactics
    └─ Mutually Exclusive with: Formation Specialist
  [Formation Specialist] - +50% to one formation, -15% to others
    └─ Mutually Exclusive with: Versatile Genius
```

### Line Formation Tree
```
Tier 1 (Choose 1):
  [Defensive Stance] - +2 to all shields in front position
  [Offensive Posture] - +15% damage, -1 to all shields in front
  [Balanced Deployment] - +5% HP and damage to all positions
  └─ Mutually Exclusive with: Defensive Stance, Offensive Posture

Tier 2 (Requires 1 from Tier 1):
  [Long-Range Barrage] - +2 range, +10% accuracy in back position
  [Shield Wall] - +1 shield to all positions
  [Volley Fire] - +10% damage from back row
  [Tactical Withdrawal] - -20% damage when retreating

Tier 3 (Requires 2 from Tier 2):
  [Concentrated Fire] - +15% damage to focused targets
  [Phalanx Training] - +10% damage resistance in front
  [Mobile Reserves] - Faster unit repositioning
  [Master of the Line] - Front: +20% damage, Back: +15% accuracy

Tier 4 (Requires 3 from Tier 3):
  [Unbreakable Phalanx] - Front: +2 to all shields, +10% damage resistance
  [Perfect Volley] - Back: +25% damage on first strike
  └─ Mutually Exclusive with: Unbreakable Phalanx
```

### Box Formation Tree
```
Tier 1 (Choose 1):
  [Reinforced Hulls] - +10% HP to all positions
  [Turtle Formation] - +15% shield effectiveness
  [Mobile Fortress] - +5% speed, -5% damage

Tier 2 (Requires 1 from Tier 1):
  [Counter-Barrage] - +15% damage when attacked
  [Shield Harmonization] - +1 shield to all positions
  [Rapid Response] - -15% ability cooldowns

Tier 3 (Requires 2 from Tier 2):
  [Impervious Defense] - +20% shield HP, +10% resistance
  [Shockwave Pulse] - 10% chance to stun attackers
  [Adaptive Shielding] - +1 to all shields when below 50% HP
```

### Vanguard Formation Tree
```
Tier 1 (Choose 1):
  [Shock Tactics] - +10% damage on first strike
  [Armor Piercing] - +10% armor penetration
  [Rapid Assault] - +5% speed, +5% attack speed

Tier 2 (Requires 1 from Tier 1):
  [Overwhelm Defenses] - +15% damage to shields
  [Precision Strikes] - +10% critical hit chance
  [Blitzkrieg] - +10% speed for first 3 turns

Tier 3 (Requires 2 from Tier 2):
  [Decapitation Strike] - +25% damage to capital ships
  [Shock and Awe] - 15% chance to disable random system
  [Relentless Assault] - +10% damage for each consecutive attack
```

### Skirmish Formation Tree
```
Tier 1 (Choose 1):
  [Hit and Run] - +10% damage when attacking first
  [Evasive Maneuvers] - +10% dodge chance
  [Guerrilla Tactics] - +5% damage per enemy ship

Tier 2 (Requires 1 from Tier 1):
  [Flanking Speed] - +10% speed when moving to flank
  [Precision Strikes] - +10% critical hit chance
  [Tactical Withdrawal] - -20% damage when retreating

Tier 3 (Requires 2 from Tier 2):
  [Master of the Skirmish] - +15% damage when outnumbered
  [Perfect Flank] - +25% damage from flanking
  [Vanishing Act] - 10% chance to disengage when hit
```

### Phalanx Formation Tree
```
Tier 1 (Choose 1):
  [Frontal Fortress] - Front: +3 to all shields, +15% HP, +10% damage
  [Shield Wall] - Front: +2 to all shields, +10% damage resistance
  [Spearhead] - Front: +20% damage, -1 to all shields

Tier 2 (Requires 1 from Tier 1):
  [Shield Bash] - Front: 20% chance to stun for 1 turn on hit
  [Extended Line] - Back: +3 range, Front: +50% ship capacity
  [Reinforced Front] - Front: +2 to all shields, +10% HP

Tier 3 (Requires 2 from Tier 2):
  [Unbreakable Wall] - Front: +30% damage resistance when not moving
  [Shock Assault] - Front: +25% damage on first attack
  [Phalanx Discipline] - All: +5% accuracy, +5% damage

Tier 4 (Requires 3 from Tier 3):
  [Titan's Wall] - Front: +4 to all shields, +20% HP
  [Juggernaut] - Front: +40% damage, +2 movement
  └─ Mutually Exclusive with: Titan's Wall
```

### Echelon Formation Tree
```
Tier 1 (Choose 1):
  [Staggered Advance] - Front/Flank: +1 shield, +10% damage
  [Diagonal Assault] - Flank: +15% damage, +10% accuracy
  [Echelon Defense] - Front: +2 shields, +10% damage resistance

Tier 2 (Requires 1 from Tier 1):
  [Crossfire] - +15% damage to targets between multiple ships
  [Enveloping Maneuver] - Flank: +2 movement, +10% damage
  [Reflexive Response] - +10% dodge, +10% accuracy against flanking

Tier 3 (Requires 2 from Tier 2):
  [Master Tactician] - +1 action point, +10% damage
  [Concentrated Fire] - +20% damage to focused targets
  [Adaptive Formation] - +10% to all stats when changing formation

Tier 4 (Requires 3 from Tier 3):
  [Decisive Engagement] - +30% damage when attacking first
  [Perfect Echelon] - All: +1 to all stats, +10% accuracy
  └─ Mutually Exclusive with: Decisive Engagement
```

### Swarm Formation Tree
```
Tier 1 (Choose 1):
  [Overwhelm] - +5% damage per adjacent friendly ship
  [Coordinated Strike] - +10% accuracy, +10% critical chance
  [Disruptive Tactics] - 10% chance to disable random system on hit

Tier 2 (Requires 1 from Tier 1):
  [Swarm Tactics] - +1 movement, +10% damage when outnumbering
  [Hit and Fade] - +15% damage on first strike, can disengage
  [Feint and Strike] - 15% dodge, +10% damage after dodging

Tier 3 (Requires 2 from Tier 2):
  [Overrun] - +20% damage to damaged targets
  [Swarm Intelligence] - +10% to all stats per adjacent friendly
  [Adaptive Swarm] - +1 movement after killing an enemy

Tier 4 (Requires 3 from Tier 3):
  [Endless Swarm] - Resurrect 1 ship per turn (max 3)
  [Perfect Swarm] - +30% damage, +2 movement, +20% dodge
  └─ Mutually Exclusive with: Endless Swarm
```
## Gems System

### Gem Families
1. **Laser**: Enhances energy weapons and precision
2. **Nuclear**: Boosts explosive damage and area effects
3. **Antimatter**: Powers exotic weapons and abilities
4. **Kinetic**: Improves projectile weapons and armor
5. **Sensor**: Enhances detection and targeting
6. **Warp**: Affects mobility and positioning
7. **Engineering**: Improves systems and repairs
8. **Logistics**: Enhances support and supply

### Gem Tiers
- **Tier 1-5**: Increasing power and effects
- **Pure Gems**: Single-family, reliable effects
- **Hybrid Gems**: Combine 2 families, unique effects
- **Relics**: Combine 3+ families, powerful bonuses

### Gem Words
Special combinations that unlock powerful effects:
- **Photon Overcharge**: Laser + Sensor - Increased energy weapon range
- **Singularity Core**: Antimatter + Warp - Creates gravity wells
- **Quantum Shielding**: Engineering + Logistics - Damage reduction

## Role Modes

### Tactical (Default)
- **Effects**:
  - +10% attack speed
  - +10% damage
  - No economy/science abilities
- **Best For**: Combat situations

### Economic
- **Effects**:
  - Enables resource gathering
  - -25% damage while active
  - Can anchor to asteroids/nebulas
- **Best For**: Resource collection

### Recon
- **Effects**:
  +3 visibility
  -1 attack range
  -15% damage
  +25% ping range
- **Best For**: Scouting and detection

### Scientific
- **Effects**:
  +50% out-of-combat regen
  -10% ability cooldowns
  +10% attack interval
  -20% damage
- **Best For**: Support and exploration
