# Formation System Quick Start Guide

## Basic Usage

### 1. Setting a Formation

```go
import "time"

// Create or get your ship stack
stack := &ShipStack{
    Ships: map[ShipType][]HPBucket{
        Fighter: {{HP: 200, Count: 10}},
        Bomber: {{HP: 500, Count: 5}},
    },
    Role: RoleTactical,
}

// Set formation (auto-assigns ships to optimal positions)
now := time.Now()
eta := stack.SetFormation(FormationVanguard, now)

// Formation will be active after reconfiguration time
// Tactical mode: 60s * 0.7 = 42 seconds
```

### 2. Getting Ship Stats with Formation Bonuses

```go
// Get effective stats for a ship in formation
effectiveShip, abilities := stack.EffectiveShipInFormation(Fighter, 0)

// effectiveShip includes:
// - Base stats
// - Role mode bonuses
// - Gem bonuses
// - Formation position bonuses
// - Gem-position synergies
// - Composition bonuses
```

### 3. Combat with Formations

```go
// Execute battle round
result := ExecuteFormationBattleRound(attackerStack, defenderStack)

fmt.Printf("Formation advantage: %.2fx\n", result.FormationAdvantage)
fmt.Printf("Attacker dealt: %d damage\n", result.AttackerDamageDealt)
fmt.Printf("Defender ships lost: %v\n", result.DefenderShipsLost)
```

## Formation Types at a Glance

| Formation | Best For | Counters | Weak Against | Speed |
|-----------|----------|----------|--------------|-------|
| **Line** | Balanced battles | Vanguard | Box | 1.0x |
| **Box** | Defense, sieges | Line | Vanguard | 0.75x |
| **Vanguard** | Aggressive assault | Box, Skirmish | Line | 1.1x |
| **Skirmish** | Hit-and-run | Phalanx | Vanguard | 1.2x |
| **Echelon** | Versatile | Various | Skirmish | 0.95x |
| **Phalanx** | Frontal assault | Box | Skirmish | 0.8x |
| **Swarm** | Anti-AoE | Phalanx | Vanguard | 1.05x |

## Position Bonuses Quick Reference

### Front Position
- **Best for:** Fighters, Destroyers, Carriers (in Box)
- **Bonuses:** +Shields, +Damage, +HP
- **Gem synergy:** Laser, Nuclear, Kinetic, Antimatter
- **Abilities:** Focus Fire, Alpha Strike, Overload, Adaptive Targeting

### Flank Position
- **Best for:** Scouts, fast Destroyers
- **Bonuses:** +Speed, +Crit, +Evasion
- **Gem synergy:** Warp, Laser, Sensor
- **Abilities:** Evasive Maneuvers, Interdictor Pulse, Cloak While Anchored

### Back Position
- **Best for:** Bombers, long-range ships
- **Bonuses:** +Range, +Visibility, +Accuracy
- **Gem synergy:** Sensor, Antimatter
- **Abilities:** Standoff Pattern, Ping, Long-Range Sensors, Siege Payload

### Support Position
- **Best for:** Carriers, Drones
- **Bonuses:** +Abilities, +Transport, +Regen
- **Gem synergy:** Engineering, Logistics, Kinetic
- **Abilities:** Point Defense Screen, Self-Repair, Targeting Uplink

## Common Patterns

### Aggressive Strike Force
```go
// Composition: Fighters + Destroyers
// Formation: Vanguard
// Position: All in Front
// Bonuses: +25% damage, fast reconfig (60s)
```

### Defensive Mining Fleet
```go
// Composition: Drones + Carrier + Fighters
// Formation: Box
// Role: Economic
// Position: Drones in Support, Fighters in Front
// Bonuses: +10% all shields, even damage distribution
```

### Scout Reconnaissance
```go
// Composition: Scouts + light Fighters
// Formation: Swarm or Skirmish
// Role: Recon
// Position: Scouts in Flank
// Bonuses: +Speed, +Visibility, dispersed (anti-AoE)
```

### Siege Armada
```go
// Composition: Bombers + Carrier
// Formation: Line or Echelon
// Role: Tactical
// Position: Bombers in Back, Carrier in Support
// Bonuses: +25% structure damage, +1 range, +1 splash
```

## Auto-Selection Strategy

### Counter-Pick Enemy Formation
```go
// Find best formation vs enemy
template := FindBestTemplate(
    myStack.Ships,
    myStack.Role,
    enemyStack.Formation.Type,
)

if template != nil {
    myStack.SetFormation(template.Formation, time.Now())
}
```

### Get Recommendations
```go
// Get formation suggestions based on fleet composition
recommendations := GetFormationRecommendations(myStack.Ships, myStack.Role)

// Pick first recommendation
if len(recommendations) > 0 {
    myStack.SetFormation(recommendations[0], time.Now())
}
```

## Optimization Tips

### 1. Match Gems to Positions
```go
// Front-line fighters: Laser + Kinetic gems
// Back-line bombers: Nuclear + Sensor gems
// Flank scouts: Warp + Sensor gems
// Support carriers: Engineering + Logistics gems
```

### 2. Use Abilities from Optimal Positions
- **Focus Fire** from Front → -50% cooldown, +20% damage
- **Evasive Maneuvers** from Flank → +25% evasion, +50% duration
- **Standoff Pattern** from Back → +30% range, +15% damage

### 3. Build for Composition Bonuses
- **Balanced Fleet:** 1 Scout + 2 Fighters + 1 Bomber
- **Strike Force:** 3 Fighters + 1 Destroyer
- **Recon Squadron:** 3 Scouts

### 4. Role-Formation Synergy
- **Tactical + Vanguard:** Fast reconfig (42s), +10% position bonuses
- **Economic + Box:** Enhanced defense, -50% slower reconfig
- **Recon + Skirmish/Swarm:** +Visibility, flanking detection

## Common Mistakes to Avoid

❌ **Don't:** Put slow ships in Skirmish formation
✅ **Do:** Match formation to fleet speed profile

❌ **Don't:** Use Vanguard in Economic mode
✅ **Do:** Use Box or Line for Economic operations

❌ **Don't:** Ignore formation counters
✅ **Do:** Scout enemy and counter-pick formations

❌ **Don't:** Put Bombers in Front position
✅ **Do:** Auto-assign or manually place by ship role

❌ **Don't:** Change formation during cooldown
✅ **Do:** Plan formation before engagement

## Debugging and Analysis

### Check Formation Status
```go
// Is formation reconfiguring?
if stack.IsFormationReconfiguring(time.Now()) {
    fmt.Println("Formation still reconfiguring...")
}

// Get formation info
info := GetFormationInfo(stack.Formation)
fmt.Println(info)
```

### Validate Formation
```go
errors := ValidateFormation(stack.Formation)
if len(errors) > 0 {
    fmt.Printf("Formation issues: %v\n", errors)
}
```

### Analyze Effectiveness
```go
effectiveness := AnalyzePositionEffectiveness(stack)
for position, score := range effectiveness {
    fmt.Printf("%s: %.2f effectiveness\n", position, score)
}
```

### Get Improvement Suggestions
```go
suggestions := SuggestFormationChanges(stack, enemyFormation)
for _, suggestion := range suggestions {
    fmt.Printf("💡 %s\n", suggestion)
}
```

## Integration Checklist

- [x] Formation data structures
- [x] Formation catalog with 7 types
- [x] Rock-paper-scissors counter matrix
- [x] Position-based damage distribution
- [x] Auto-assignment algorithms
- [x] ShipStack integration
- [x] Role mode interactions
- [x] Gem-position synergies
- [x] Ability-position enhancements
- [x] Composition bonuses
- [x] Combat integration
- [x] Utility functions
- [x] Example implementations

## Next Steps

1. **Test the system:** Run example functions in `formation_examples.go`
2. **API integration:** Create REST endpoints for formation management
3. **UI design:** Build formation editor and visualizer
4. **Balance tuning:** Adjust multipliers based on gameplay
5. **Database:** Ensure MongoDB schema supports formations
6. **Documentation:** Add API docs and UI guides

## Support

See full documentation in:
- `FORMATION_IMPLEMENTATION.md` - Complete technical details
- `FORMATION.md` - Original design document
- `formation_examples.go` - Code examples
- `formation_utils.go` - Utility functions
