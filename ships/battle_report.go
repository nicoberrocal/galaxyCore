package ships

import (
	"time"

	bson "go.mongodb.org/mongo-driver/v2/bson"
)

// BattleReport is a comprehensive record of a battle between two stacks.
// It serves as both a live-updating document and historical record.
// Each stack can have multiple reports if attacked by multiple enemies.
type BattleReport struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// Battle Identity
	BattleID       string        `bson:"battleId" json:"battleId"`                             // Unique battle identifier
	AttackerStackID bson.ObjectID `bson:"attackerStackId" json:"attackerStackId"`               // Attacker stack ID
	DefenderStackID bson.ObjectID `bson:"defenderStackId" json:"defenderStackId"`               // Defender stack ID
	AttackerPlayerID bson.ObjectID `bson:"attackerPlayerId" json:"attackerPlayerId"`            // Attacker player ID
	DefenderPlayerID bson.ObjectID `bson:"defenderPlayerId" json:"defenderPlayerId"`            // Defender player ID
	
	// Battle Metadata
	StartedAt time.Time `bson:"startedAt" json:"startedAt"`                                   // When battle began
	EndedAt   *time.Time `bson:"endedAt,omitempty" json:"endedAt,omitempty"`                  // When battle ended (nil if ongoing)
	Location  BattleLocation `bson:"location" json:"location"`                                 // Where battle occurred
	Status    BattleStatus `bson:"status" json:"status"`                                       // Current battle status
	Outcome   *BattleOutcome `bson:"outcome,omitempty" json:"outcome,omitempty"`               // Final outcome (nil if ongoing)
	
	// Initial State Snapshots
	AttackerInitial StackSnapshot `bson:"attackerInitial" json:"attackerInitial"`               // Attacker state at battle start
	DefenderInitial StackSnapshot `bson:"defenderInitial" json:"defenderInitial"`               // Defender state at battle start
	
	// Current State (live updates)
	AttackerCurrent StackSnapshot `bson:"attackerCurrent" json:"attackerCurrent"`               // Current attacker state
	DefenderCurrent StackSnapshot `bson:"defenderCurrent" json:"defenderCurrent"`               // Current defender state
	
	// Round-by-Round Timeline
	Rounds []BattleRound `bson:"rounds" json:"rounds"`                                         // Chronological battle rounds
	
	// Aggregate Statistics
	TotalRounds         int `bson:"totalRounds" json:"totalRounds"`                            // Total rounds fought
	AttackerTotalDamage int `bson:"attackerTotalDamage" json:"attackerTotalDamage"`            // Total damage dealt by attacker
	DefenderTotalDamage int `bson:"defenderTotalDamage" json:"defenderTotalDamage"`            // Total damage dealt by defender
	AttackerShipsLost   map[ShipType]int `bson:"attackerShipsLost" json:"attackerShipsLost"`   // Ships lost by attacker
	DefenderShipsLost   map[ShipType]int `bson:"defenderShipsLost" json:"defenderShipsLost"`   // Ships lost by defender
	
	// Timestamps
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// BattleLocation describes where the battle took place
type BattleLocation struct {
	Type       string        `bson:"type" json:"type"`                                         // "empty_space", "asteroid", "nebula", "planet_orbit"
	LocationID bson.ObjectID `bson:"locationId,omitempty" json:"locationId,omitempty"`         // ID of location if applicable
	X          float64       `bson:"x" json:"x"`                                               // Coordinates
	Y          float64       `bson:"y" json:"y"`
	SectorID   bson.ObjectID `bson:"sectorId,omitempty" json:"sectorId,omitempty"`             // Sector ID
}

// BattleStatus represents the current state of the battle
type BattleStatus string

const (
	BattleStatusOngoing  BattleStatus = "ongoing"   // Battle is currently active
	BattleStatusEnded    BattleStatus = "ended"     // Battle has concluded
	BattleStatusRetreat  BattleStatus = "retreat"   // One side retreated
	BattleStatusStalemate BattleStatus = "stalemate" // Battle ended in stalemate
)

// BattleOutcome describes how the battle ended
type BattleOutcome struct {
	Victor       string        `bson:"victor" json:"victor"`                                   // "attacker", "defender", "draw"
	VictorStackID bson.ObjectID `bson:"victorStackId,omitempty" json:"victorStackId,omitempty"` // Winner's stack ID
	Reason       string        `bson:"reason" json:"reason"`                                   // "total_destruction", "retreat", "timeout"
	EndedAt      time.Time     `bson:"endedAt" json:"endedAt"`                                 // When battle concluded
}

// StackSnapshot captures the complete state of a stack at a point in time
type StackSnapshot struct {
	StackID   bson.ObjectID `bson:"stackId" json:"stackId"`
	PlayerID  bson.ObjectID `bson:"playerId" json:"playerId"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	
	// Fleet Composition
	Ships      map[ShipType][]HPBucket `bson:"ships" json:"ships"`                             // Current ships and HP
	TotalShips int                     `bson:"totalShips" json:"totalShips"`                   // Total ship count
	TotalHP    int                     `bson:"totalHp" json:"totalHp"`                         // Total HP across all ships
	
	// Formation
	Formation *FormationSnapshot `bson:"formation,omitempty" json:"formation,omitempty"`       // Formation configuration
	
	// Bio State
	BioPath        string                  `bson:"bioPath,omitempty" json:"bioPath,omitempty"`   // Active bio tree path
	ActiveBioNodes []string                `bson:"activeBioNodes,omitempty" json:"activeBioNodes,omitempty"` // Active bio node IDs
	BioDebuffs     []BioDebuffSnapshot     `bson:"bioDebuffs,omitempty" json:"bioDebuffs,omitempty"` // Active debuffs
	
	// Combat Counters
	AttackCount  int `bson:"attackCount" json:"attackCount"`                                   // Total attacks made
	DefenseCount int `bson:"defenseCount" json:"defenseCount"`                                 // Total attacks received
	
	// Effective Stats (per ship type, with all modifiers applied)
	EffectiveStats map[ShipType]EffectiveShipStats `bson:"effectiveStats" json:"effectiveStats"`
}

// FormationSnapshot captures formation state at a point in time
type FormationSnapshot struct {
	Type      FormationType                       `bson:"type" json:"type"`
	Level     int                                 `bson:"level" json:"level"`
	Positions map[FormationPosition][]ShipAssignment `bson:"positions" json:"positions"` // Ships assigned to each position
	TreeNodes []string                            `bson:"treeNodes,omitempty" json:"treeNodes,omitempty"` // Unlocked tree node IDs
}

// ShipAssignment describes which ships are in which formation position
type ShipAssignment struct {
	ShipType    ShipType `bson:"shipType" json:"shipType"`
	BucketIndex int      `bson:"bucketIndex" json:"bucketIndex"`
	Count       int      `bson:"count" json:"count"`
	HP          int      `bson:"hp" json:"hp"`
}

// EffectiveShipStats captures the final computed stats for a ship type
type EffectiveShipStats struct {
	ShipType ShipType `bson:"shipType" json:"shipType"`
	
	// Base Stats
	BaseAttackDamage     int `bson:"baseAttackDamage" json:"baseAttackDamage"`
	BaseLaserShield      int `bson:"baseLaserShield" json:"baseLaserShield"`
	BaseNuclearShield    int `bson:"baseNuclearShield" json:"baseNuclearShield"`
	BaseAntimatterShield int `bson:"baseAntimatterShield" json:"baseAntimatterShield"`
	BaseHP               int `bson:"baseHp" json:"baseHp"`
	BaseSpeed            int `bson:"baseSpeed" json:"baseSpeed"`
	
	// Effective Stats (with all modifiers)
	EffectiveAttackDamage     int `bson:"effectiveAttackDamage" json:"effectiveAttackDamage"`
	EffectiveLaserShield      int `bson:"effectiveLaserShield" json:"effectiveLaserShield"`
	EffectiveNuclearShield    int `bson:"effectiveNuclearShield" json:"effectiveNuclearShield"`
	EffectiveAntimatterShield int `bson:"effectiveAntimatterShield" json:"effectiveAntimatterShield"`
	EffectiveHP               int `bson:"effectiveHp" json:"effectiveHp"`
	EffectiveSpeed            int `bson:"effectiveSpeed" json:"effectiveSpeed"`
	
	// Modifier Breakdown
	Modifiers ModifierBreakdown `bson:"modifiers" json:"modifiers"`
}

// ModifierBreakdown shows where stat bonuses come from
type ModifierBreakdown struct {
	Formation []ModifierSourceDetail `bson:"formation,omitempty" json:"formation,omitempty"`         // Formation bonuses
	Bio       []ModifierSourceDetail `bson:"bio,omitempty" json:"bio,omitempty"`                     // Bio trait bonuses
	Gems      []ModifierSourceDetail `bson:"gems,omitempty" json:"gems,omitempty"`                   // Gem bonuses
	Buffs     []ModifierSourceDetail `bson:"buffs,omitempty" json:"buffs,omitempty"`                 // Active buffs
	Debuffs   []ModifierSourceDetail `bson:"debuffs,omitempty" json:"debuffs,omitempty"`             // Active debuffs
}

// ModifierSourceDetail describes a single modifier contribution
type ModifierSourceDetail struct {
	SourceID    string    `bson:"sourceId" json:"sourceId"`                                   // ID of the source
	Description string    `bson:"description" json:"description"`                             // Human-readable description
	Mods        StatMods  `bson:"mods" json:"mods"`                                           // The modifiers applied
}

// BioDebuffSnapshot captures active bio debuff state
type BioDebuffSnapshot struct {
	DebuffID   string        `bson:"debuffId" json:"debuffId"`
	SourceID   bson.ObjectID `bson:"sourceId" json:"sourceId"`                                 // Who applied it
	Stacks     int           `bson:"stacks" json:"stacks"`
	MaxStacks  int           `bson:"maxStacks" json:"maxStacks"`
	AppliedAt  time.Time     `bson:"appliedAt" json:"appliedAt"`
	ExpiresAt  time.Time     `bson:"expiresAt" json:"expiresAt"`
	Mods       StatMods      `bson:"mods" json:"mods"`
}

// BattleRound represents a single hourly combat round
type BattleRound struct {
	RoundNumber int       `bson:"roundNumber" json:"roundNumber"`                             // Round sequence number
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`                                 // When this round occurred
	
	// Pre-Round State
	AttackerPreRound CombatantState `bson:"attackerPreRound" json:"attackerPreRound"`         // Attacker state before round
	DefenderPreRound CombatantState `bson:"defenderPreRound" json:"defenderPreRound"`         // Defender state before round
	
	// Combat Events
	AttackerPhase CombatPhase `bson:"attackerPhase" json:"attackerPhase"`                     // Attacker's attack phase
	DefenderPhase CombatPhase `bson:"defenderPhase" json:"defenderPhase"`                     // Defender's return fire phase
	
	// Post-Round State
	AttackerPostRound CombatantState `bson:"attackerPostRound" json:"attackerPostRound"`       // Attacker state after round
	DefenderPostRound CombatantState `bson:"defenderPostRound" json:"defenderPostRound"`       // Defender state after round
	
	// Round Summary
	AttackerDamageDealt int              `bson:"attackerDamageDealt" json:"attackerDamageDealt"` // Total damage by attacker
	DefenderDamageDealt int              `bson:"defenderDamageDealt" json:"defenderDamageDealt"` // Total damage by defender
	AttackerShipsLost   map[ShipType]int `bson:"attackerShipsLost" json:"attackerShipsLost"`     // Ships lost this round
	DefenderShipsLost   map[ShipType]int `bson:"defenderShipsLost" json:"defenderShipsLost"`     // Ships lost this round
	
	// Special Events
	Events []RoundEvent `bson:"events,omitempty" json:"events,omitempty"`                     // Special events (crits, debuffs, etc.)
}

// CombatantState captures a combatant's state at a specific moment
type CombatantState struct {
	TotalShips int                     `bson:"totalShips" json:"totalShips"`
	TotalHP    int                     `bson:"totalHp" json:"totalHp"`
	Ships      map[ShipType][]HPBucket `bson:"ships" json:"ships"`
	
	// Combat Counters
	AttackCount  int `bson:"attackCount" json:"attackCount"`
	DefenseCount int `bson:"defenseCount" json:"defenseCount"`
	
	// Active Effects
	ActiveBuffs   []string `bson:"activeBuffs,omitempty" json:"activeBuffs,omitempty"`
	ActiveDebuffs []string `bson:"activeDebuffs,omitempty" json:"activeDebuffs,omitempty"`
}

// CombatPhase describes what happened during an attack phase
type CombatPhase struct {
	AttackerID bson.ObjectID `bson:"attackerId" json:"attackerId"`
	DefenderID bson.ObjectID `bson:"defenderId" json:"defenderId"`
	
	// Damage Calculation
	BaseDamage              int                     `bson:"baseDamage" json:"baseDamage"`                 // Raw damage before modifiers
	FormationMultiplier     float64                 `bson:"formationMultiplier" json:"formationMultiplier"` // Formation counter bonus
	FirstStrikeBonus        bool                    `bson:"firstStrikeBonus" json:"firstStrikeBonus"`     // First strike applied?
	CriticalHit             bool                    `bson:"criticalHit" json:"criticalHit"`               // Crit applied?
	FinalDamage             int                     `bson:"finalDamage" json:"finalDamage"`               // Damage after all modifiers
	
	// Damage Distribution
	DamageByType     map[string]int              `bson:"damageByType" json:"damageByType"`             // Damage by attack type
	DamageByShipType map[ShipType]map[int]int    `bson:"damageByShipType" json:"damageByShipType"`     // Damage to each ship type/bucket
	
	// Shield Mitigation
	ShieldMitigation map[string]ShieldMitigationDetail `bson:"shieldMitigation" json:"shieldMitigation"` // Per-type shield mitigation
	
	// Evasion
	EvasionReduction float64 `bson:"evasionReduction" json:"evasionReduction"`                 // % damage reduced by evasion
	
	// Casualties
	ShipsDestroyed map[ShipType]int `bson:"shipsDestroyed" json:"shipsDestroyed"`               // Ships destroyed this phase
	
	// Bio Effects Applied
	DebuffsApplied []string `bson:"debuffsApplied,omitempty" json:"debuffsApplied,omitempty"` // Debuffs applied this phase
}

// ShieldMitigationDetail shows how shields mitigated damage
type ShieldMitigationDetail struct {
	AttackType       string  `bson:"attackType" json:"attackType"`                             // "Laser", "Nuclear", "Antimatter"
	RawDamage        int     `bson:"rawDamage" json:"rawDamage"`                               // Damage before shields
	ShieldValue      int     `bson:"shieldValue" json:"shieldValue"`                           // Shield strength
	MitigatedDamage  int     `bson:"mitigatedDamage" json:"mitigatedDamage"`                   // Damage after shields
	MitigationPercent float64 `bson:"mitigationPercent" json:"mitigationPercent"`               // % damage blocked
}

// RoundEvent describes special events that occurred during a round
type RoundEvent struct {
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	EventType   string    `bson:"eventType" json:"eventType"`                                 // "crit", "first_strike", "debuff_applied", "ship_destroyed", etc.
	ActorID     bson.ObjectID `bson:"actorId" json:"actorId"`                                 // Who caused the event
	TargetID    bson.ObjectID `bson:"targetId,omitempty" json:"targetId,omitempty"`           // Who was affected
	Description string    `bson:"description" json:"description"`                             // Human-readable description
	Data        map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`          // Additional event data
}
