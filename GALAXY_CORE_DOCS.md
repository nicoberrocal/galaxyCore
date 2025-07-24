# Galaxy Core - Game Documentation

## Overview
Galaxy Core is a browser-based MMO RTS (Massive Multiplayer Online Real-Time Strategy) game built in Go. The game features space exploration, resource management, colony building, and fleet combat mechanics.

## Core Game Elements

### 1. Player Management
```go
type Player struct {
    ID         ObjectID
    Username   string
    Email      string
    IsVerified bool
    IsPremium  bool
}
```
Players are the core users of the game, with support for premium accounts and email verification.

### 2. Game State Management
```go
type PlayerGameState struct {
    PlayerID         ObjectID
    MapID           ObjectID
    Energy          int64
    EnergyProduction int64
    
    // Territory tracking
    ColonizedSystems []ObjectID
    ActiveStacks     []ObjectID
    MiningOperations []ObjectID
}
```

### 3. Ship Types and Combat

#### Available Ship Classes
- **Drone**: Mining specialist with Mining Overdrive ability
- **Scout**: Fast reconnaissance ship with Deep Scan Pulse
- **Fighter**: Combat specialist with Stealth Cloak
- **Bomber**: Heavy damage dealer with Siege Mode
- **Carrier**: Transport ship with FTL Jump capability
- **Destroyer**: Heavy combat ship with Hunter Protocol

#### Ship Attributes
- Attack Types: Laser, Nuclear, Antimatter
- Shield Systems: Laser, Nuclear, Antimatter shields
- Combat Stats: HP, Attack Damage, Attack Range
- Movement: Speed, Visibility Range
- Special Abilities with cooldowns
- Resource Costs: Metal, Crystal, Plasma

### 4. Building System

#### Building Types
1. **Resource Production**
   - Metal Mine
   - Crystal Mine
   - Solar Farm
   - Wind Farm
   - Hydro Electric Dam
   - Balloon Energy Collector

2. **Special Facilities**
   - Shipyard
   - Particle Accelerator
   - Fusion Reactor

#### Building Characteristics
- Level-based progression
- Production rates
- Energy upkeep costs
- Construction queues
- Planet-specific suitability

### 5. Territory and Resource Management

#### Colonizable Objects
1. **Systems**
   - Colonization state tracking
   - Defending fleet management
   - Building placement
   - Resource generation

2. **Planets**
   - Building slots
   - Resource deposits
   - Environmental factors
   - Building suitability modifiers

3. **Resource Locations**
   - **Asteroids**: Minable resources
   - **Nebulas**: Special resource extraction

#### Resource Types
- Metal
- Crystal
- Hydrogen
- Plasma
- Energy

### 6. Combat System

#### Battle Scenarios
1. **System Battles**
   - Colonization-focused
   - Defending fleet mechanics
   - Allied fleet support

2. **Free Space Battles**
   - Stack vs Stack combat
   - Movement interruption
   - Fleet damage resolution

3. **Mining Location Battles**
   - Resource control
   - Mining operation disruption

#### Combat Resolution Rules
1. **Victory Conditions**
   - Physical entry required for colonization
   - Allied survivors protect ownership
   - Stack destruction mechanics

2. **Fleet Mechanics**
   - HP bucketing system
   - Ship type effectiveness
   - Shield type interactions

### 7. Movement and Actions

#### Queue System
```go
type Queue struct {
    Type      string        // ship_attack, ship_construction, building_construction
    StartTime time.Time
    EndTime   time.Time
    Payload   interface{}
}
```

#### Movement States
- Traveling
- Mining
- Idle
- In Combat

### 8. Map and Game Configuration

#### Map Properties
```go
type MongoMap struct {
    GameName  string
    QPlayers  int8
    PeaceDays int8
    StartTime time.Time
    Ranked    bool
}
```

#### Game Settings
- Player limits
- Peace period duration
- Ranked match options
- Start conditions

## Technical Implementation

### Data Model Approach
The game uses a hybrid data model that combines:
1. Embedded documents for system defense
2. Separate collections for active fleets
3. Denormalized references for performance

### Consistency Rules
1. **System State**
   - DefendingFleet existence ↔ IsColonized state
   - Original owner protection mechanics
   - Allied fleet integration

2. **Player Assets**
   - Mutually exclusive arrays for tracking
   - Clear ownership hierarchies
   - Resource balance management

### Performance Considerations
1. **Combat Resolution**
   - Single document updates for system battles
   - Parallel processing for free space combat
   - Atomic operations for resource extraction

2. **Resource Management**
   - Cached production rates
   - Periodic state updates
   - Optimistic locking for conflicts

### Planet-Specific Mechanics

#### Resource Suitability
```go
Mercury:
    Metals:   0.8
    Crystals: 1.2
Venus:
    Metals:   0.6
    Crystals: 0.7
Earth:
    Metals:   1.0
    Crystals: 0.9
Mars:
    Metals:   1.2
    Crystals: 0.5
```

#### Energy Production Efficiency
```go
Base Output:
    Hydro:   100
    Solar:   80
    Wind:    60
    Balloon: 50
```

## Game Mechanics Deep Dive

### 1. Fleet Management
- Ships are organized into stacks for movement and combat
- HP bucketing system for efficient damage tracking
- Special abilities with strategic cooldowns

### 2. Resource Generation
- Planet-specific extraction rates
- Building level progression
- Environmental modifiers
- Production efficiency calculations

### 3. Combat Resolution
- Multi-phase battle resolution
- Shield type effectiveness
- Stack-based damage distribution
- Territory control mechanics

### 4. Construction System
- Queue-based building construction
- Resource cost scaling
- Planet suitability factors
- Energy balance requirements
