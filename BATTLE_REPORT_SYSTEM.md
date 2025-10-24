# Battle Report System

## Overview

A comprehensive battle reporting system that captures every detail of combat between ship stacks. Serves as both a live-updating combat tracker and historical battle record.

## Key Features

### ✅ **Complete State Snapshots**
- Initial state of both combatants
- Current state (live updates)
- Pre/post-round states for each combat round
- Base stats + effective stats with all modifiers applied

### ✅ **Full Modifier Breakdown**
- Formation bonuses (position, counter, tree nodes)
- Bio trait effects (passive, triggered, accumulative)
- Gem bonuses
- Active buffs and debuffs
- Shows exact source of each modifier

### ✅ **Round-by-Round Timeline**
- Chronological record of all combat rounds
- Attacker phase and defender phase details
- Damage calculation breakdown
- Shield mitigation by attack type
- Evasion reduction
- Ships destroyed per round
- Special events (crits, first strikes, debuffs)

### ✅ **Multi-Enemy Support**
- Each attacker-defender pair gets its own report
- Stack can have multiple reports if attacked by multiple enemies
- Query by stack ID to get all active battles

### ✅ **Comprehensive Statistics**
- Total rounds fought
- Total damage dealt by each side
- Ships lost by type
- Formation effectiveness
- Bio trait impact

## Data Structures

### BattleReport (Main Document)
```go
type BattleReport struct {
    // Identity
    BattleID         string
    AttackerStackID  bson.ObjectID
    DefenderStackID  bson.ObjectID
    
    // Metadata
    StartedAt        time.Time
    EndedAt          *time.Time
    Location         BattleLocation
    Status           BattleStatus  // "ongoing", "ended", "retreat", "stalemate"
    Outcome          *BattleOutcome
    
    // Snapshots
    AttackerInitial  StackSnapshot
    DefenderInitial  StackSnapshot
    AttackerCurrent  StackSnapshot
    DefenderCurrent  StackSnapshot
    
    // Timeline
    Rounds           []BattleRound
    
    // Aggregates
    TotalRounds         int
    AttackerTotalDamage int
    DefenderTotalDamage int
    AttackerShipsLost   map[ShipType]int
    DefenderShipsLost   map[ShipType]int
}
```

### StackSnapshot
```go
type StackSnapshot struct {
    // Fleet Composition
    Ships      map[ShipType][]HPBucket
    TotalShips int
    TotalHP    int
    
    // Formation
    Formation  *FormationSnapshot
    
    // Bio State
    BioPath        string
    ActiveBioNodes []string
    BioDebuffs     []BioDebuffSnapshot
    
    // Combat Counters
    AttackCount  int
    DefenseCount int
    
    // Effective Stats (per ship type)
    EffectiveStats map[ShipType]EffectiveShipStats
}
```

### EffectiveShipStats
```go
type EffectiveShipStats struct {
    // Base Stats
    BaseAttackDamage     int
    BaseLaserShield      int
    BaseNuclearShield    int
    BaseAntimatterShield int
    BaseHP               int
    BaseSpeed            int
    
    // Effective Stats (with modifiers)
    EffectiveAttackDamage     int
    EffectiveLaserShield      int
    EffectiveNuclearShield    int
    EffectiveAntimatterShield int
    EffectiveHP               int
    EffectiveSpeed            int
    
    // Modifier Breakdown
    Modifiers ModifierBreakdown
}
```

### ModifierBreakdown
```go
type ModifierBreakdown struct {
    Formation []ModifierSourceDetail  // Formation bonuses
    Bio       []ModifierSourceDetail  // Bio trait bonuses
    Gems      []ModifierSourceDetail  // Gem bonuses
    Buffs     []ModifierSourceDetail  // Active buffs
    Debuffs   []ModifierSourceDetail  // Active debuffs
}
```

### BattleRound
```go
type BattleRound struct {
    RoundNumber int
    Timestamp   time.Time
    
    // Pre-Round State
    AttackerPreRound CombatantState
    DefenderPreRound CombatantState
    
    // Combat Events
    AttackerPhase CombatPhase
    DefenderPhase CombatPhase
    
    // Post-Round State
    AttackerPostRound CombatantState
    DefenderPostRound CombatantState
    
    // Summary
    AttackerDamageDealt int
    DefenderDamageDealt int
    AttackerShipsLost   map[ShipType]int
    DefenderShipsLost   map[ShipType]int
    
    // Special Events
    Events []RoundEvent
}
```

### CombatPhase
```go
type CombatPhase struct {
    // Damage Calculation
    BaseDamage              int
    FormationMultiplier     float64
    FirstStrikeBonus        bool
    CriticalHit             bool
    FinalDamage             int
    
    // Damage Distribution
    DamageByType     map[string]int              // By attack type
    DamageByShipType map[ShipType]map[int]int    // To each ship/bucket
    
    // Shield Mitigation
    ShieldMitigation map[string]ShieldMitigationDetail
    
    // Evasion
    EvasionReduction float64
    
    // Casualties
    ShipsDestroyed map[ShipType]int
    
    // Bio Effects
    DebuffsApplied []string
}
```

### RoundEvent
```go
type RoundEvent struct {
    Timestamp   time.Time
    EventType   string  // "crit", "first_strike", "debuff_applied", "ship_destroyed"
    ActorID     bson.ObjectID
    TargetID    bson.ObjectID
    Description string
    Data        map[string]interface{}
}
```

## Usage

### 1. Initiate Battle

```go
// When combat begins
location := BattleLocation{
    Type: "empty_space",
    X:    attacker.X,
    Y:    attacker.Y,
}

report := InitiateBattle(attacker, defender, location, time.Now())

// Save to database
SaveBattleReport(report)
```

### 2. Process Combat Round

```go
// In your hourly tick system
report, result := ProcessCombatWithReporting(attacker, defender, report, time.Now())

// Report is automatically updated with:
// - Round details
// - Damage breakdown
// - Shield mitigation
// - Events (crits, first strikes, debuffs)
// - Updated snapshots

// Save updated report
SaveBattleReport(report)
```

### 3. End Battle

```go
// When battle concludes
if isStackDestroyed(defender) {
    report.EndBattle(BattleOutcome{
        Victor:        "attacker",
        VictorStackID: attacker.ID,
        Reason:        "total_destruction",
        EndedAt:       time.Now(),
    }, time.Now())
}

SaveBattleReport(report)
```

### 4. Query Reports

```go
// Get all active battles for a stack
reports := GetBattleReportForStack(stackID, time.Now())

// Get specific battle
report := GetBattleReport(battleID)

// Generate summary
summary := CreateBattleReportSummary(report)
```

## Integration with Tick System

```go
func ProcessHourlyCombat(gameState *GameState, now time.Time) {
    // 1. Find all active battles
    activeBattles := gameState.GetActiveBattles()
    
    // 2. Process each battle
    for _, battle := range activeBattles {
        attacker := GetStack(battle.AttackerStackID)
        defender := GetStack(battle.DefenderStackID)
        
        // Get or create report
        report := GetBattleReport(battle.BattleID)
        if report == nil {
            location := BattleLocation{
                Type: battle.Location,
                X:    attacker.X,
                Y:    attacker.Y,
            }
            report = InitiateBattle(attacker, defender, location, now)
        }
        
        // Execute combat with reporting
        report, result := ProcessCombatWithReporting(attacker, defender, report, now)
        
        // Save everything
        SaveBattleReport(report)
        SaveStack(attacker)
        SaveStack(defender)
        
        // Clean up if battle ended
        if report.Status == BattleStatusEnded {
            CleanupBattle(battle)
        }
    }
}
```

## What Gets Captured

### Per Round
- ✅ Ship counts and HP before/after
- ✅ Combat counters (attack count, defense count)
- ✅ Active buffs and debuffs
- ✅ Base damage calculation
- ✅ Formation multiplier
- ✅ First strike bonus (if applicable)
- ✅ Critical hit (if applicable)
- ✅ Damage by attack type (Laser/Nuclear/Antimatter)
- ✅ Damage to each ship type and bucket
- ✅ Shield mitigation per attack type
- ✅ Evasion reduction
- ✅ Ships destroyed
- ✅ Bio debuffs applied

### Per Battle
- ✅ Initial fleet composition
- ✅ Formation configuration with ship assignments
- ✅ Active bio nodes
- ✅ Active debuffs with stacks
- ✅ Base stats vs effective stats
- ✅ Modifier breakdown by source
- ✅ Total damage dealt
- ✅ Total ships lost
- ✅ Battle outcome

## Example Report Structure

```json
{
  "battleId": "stack1_vs_stack2_1234567890",
  "status": "ongoing",
  "startedAt": "2025-10-24T10:00:00Z",
  "totalRounds": 5,
  "attackerInitial": {
    "totalShips": 100,
    "totalHp": 50000,
    "formation": {
      "type": "wedge",
      "positions": {
        "front": [
          {"shipType": "Fighter", "count": 30, "hp": 100}
        ],
        "flanks": [
          {"shipType": "Corvette", "count": 40, "hp": 150}
        ]
      }
    },
    "activeBioNodes": ["pack_mentality", "hunting_cry"],
    "effectiveStats": {
      "Fighter": {
        "baseAttackDamage": 10,
        "effectiveAttackDamage": 15,
        "modifiers": {
          "formation": [
            {
              "sourceId": "wedge_front_bonus",
              "description": "Wedge Formation: Front Position",
              "mods": {"damage": {"laserPct": 0.20}}
            }
          ],
          "bio": [
            {
              "sourceId": "pack_mentality",
              "description": "Pack Mentality: +10% damage per 10 allies",
              "mods": {"damage": {"laserPct": 0.30}}
            }
          ]
        }
      }
    }
  },
  "rounds": [
    {
      "roundNumber": 1,
      "timestamp": "2025-10-24T10:00:00Z",
      "attackerPhase": {
        "baseDamage": 1500,
        "formationMultiplier": 1.2,
        "firstStrikeBonus": true,
        "criticalHit": false,
        "finalDamage": 1950,
        "damageByType": {
          "Laser": 1200,
          "Nuclear": 750
        },
        "shieldMitigation": {
          "Laser": {
            "rawDamage": 1200,
            "shieldValue": 5,
            "mitigatedDamage": 828,
            "mitigationPercent": 31.0
          }
        },
        "shipsDestroyed": {
          "Fighter": 5
        }
      },
      "events": [
        {
          "eventType": "first_strike",
          "description": "Attacker unleashes first strike bonus",
          "data": {"damageBonus": "30%"}
        },
        {
          "eventType": "ships_destroyed",
          "description": "Ships destroyed",
          "data": {"shipType": "Fighter", "count": 5}
        }
      ]
    }
  ]
}
```

## Database Considerations

### Indexes
```javascript
// MongoDB indexes for efficient queries
db.battle_reports.createIndex({ "battleId": 1 })
db.battle_reports.createIndex({ "attackerStackId": 1, "status": 1 })
db.battle_reports.createIndex({ "defenderStackId": 1, "status": 1 })
db.battle_reports.createIndex({ "status": 1, "startedAt": -1 })
db.battle_reports.createIndex({ "attackerPlayerId": 1, "status": 1 })
db.battle_reports.createIndex({ "defenderPlayerId": 1, "status": 1 })
```

### Storage Optimization
- Use sparse encoding for modifiers (omit zero values)
- Archive completed battles after 30 days
- Compress round data for long battles (>100 rounds)
- Consider separate collection for round details if battles get very long

## UI/Display Recommendations

### Live Battle View
- Show current round number
- Display current HP/ships for both sides
- Show active buffs/debuffs
- Highlight recent events (last 3 rounds)
- Progress bar showing relative strength

### Historical Battle View
- Timeline scrubber to view any round
- Damage graph over time
- Formation visualization
- Modifier breakdown tooltips
- Event log with filters

### Battle Summary
- Victor and reason
- Total duration (rounds)
- Damage dealt comparison
- Ships lost comparison
- Key moments (first crit, major casualties)

## Future Enhancements

1. **Replay System**: Store enough data to replay battles visually
2. **Statistics Aggregation**: Player win/loss ratios, favorite formations, etc.
3. **Battle Predictions**: Use historical data to predict outcomes
4. **Achievement Tracking**: "First blood", "Perfect victory", etc.
5. **Spectator Mode**: Watch ongoing battles in real-time

## Files Created

1. **`battle_report.go`**: Core data structures
2. **`battle_report_builder.go`**: Report creation and update functions
3. **`battle_report_integration.go`**: Integration examples with tick system
4. **`BATTLE_REPORT_SYSTEM.md`**: This documentation

## Summary

The battle report system provides:
- ✅ Complete combat transparency
- ✅ Historical battle records
- ✅ Live battle tracking
- ✅ Multi-enemy support
- ✅ Full modifier breakdown
- ✅ Round-by-round timeline
- ✅ Event tracking
- ✅ Easy integration with existing tick system

All combat details are captured automatically when you use `ProcessCombatWithReporting()` instead of direct `ExecuteFormationBattleRound()` calls.
