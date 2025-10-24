package ships

import (
	"fmt"
	"time"

	bson "go.mongodb.org/mongo-driver/v2/bson"
)

// NewBattleReport creates a new battle report when combat begins
func NewBattleReport(
	attacker, defender *ShipStack,
	location BattleLocation,
	now time.Time,
) *BattleReport {
	battleID := fmt.Sprintf("%s_vs_%s_%d", attacker.ID.Hex(), defender.ID.Hex(), now.Unix())
	
	report := &BattleReport{
		ID:              bson.NewObjectID(),
		BattleID:        battleID,
		AttackerStackID: attacker.ID,
		DefenderStackID: defender.ID,
		AttackerPlayerID: attacker.PlayerID,
		DefenderPlayerID: defender.PlayerID,
		StartedAt:       now,
		Location:        location,
		Status:          BattleStatusOngoing,
		
		AttackerInitial: CaptureStackSnapshot(attacker, now),
		DefenderInitial: CaptureStackSnapshot(defender, now),
		AttackerCurrent: CaptureStackSnapshot(attacker, now),
		DefenderCurrent: CaptureStackSnapshot(defender, now),
		
		Rounds:              make([]BattleRound, 0, 100), // Pre-allocate for ~100 rounds
		AttackerShipsLost:   make(map[ShipType]int),
		DefenderShipsLost:   make(map[ShipType]int),
		
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	return report
}

// CaptureStackSnapshot creates a complete snapshot of a stack's current state
func CaptureStackSnapshot(stack *ShipStack, now time.Time) StackSnapshot {
	snapshot := StackSnapshot{
		StackID:   stack.ID,
		PlayerID:  stack.PlayerID,
		Timestamp: now,
		Ships:     make(map[ShipType][]HPBucket),
		EffectiveStats: make(map[ShipType]EffectiveShipStats),
	}
	
	// Copy ships and HP buckets
	totalShips := 0
	totalHP := 0
	for shipType, buckets := range stack.Ships {
		bucketsCopy := make([]HPBucket, len(buckets))
		copy(bucketsCopy, buckets)
		snapshot.Ships[shipType] = bucketsCopy
		
		for _, bucket := range buckets {
			totalShips += bucket.Count
			totalHP += bucket.HP * bucket.Count
		}
	}
	snapshot.TotalShips = totalShips
	snapshot.TotalHP = totalHP
	
	// Capture formation state
	if stack.Formation != nil {
		snapshot.Formation = CaptureFormationSnapshot(stack)
	}
	
	// Capture bio state
	if stack.Bio != nil {
		snapshot.BioPath = stack.Bio.ActivePath
		snapshot.ActiveBioNodes = make([]string, 0, len(stack.Bio.Nodes))
		for nodeID, node := range stack.Bio.Nodes {
			if node.Stage != BioStagePassive && node.Stage != BioStageCooldown {
				snapshot.ActiveBioNodes = append(snapshot.ActiveBioNodes, nodeID)
			}
		}
		
		// Capture active debuffs
		snapshot.BioDebuffs = make([]BioDebuffSnapshot, 0, len(stack.Bio.InboundDebuffs))
		for debuffID, debuff := range stack.Bio.InboundDebuffs {
			snapshot.BioDebuffs = append(snapshot.BioDebuffs, BioDebuffSnapshot{
				DebuffID:  debuffID,
				SourceID:  debuff.SourceStack,
				Stacks:    debuff.Stacks,
				MaxStacks: debuff.MaxStacks,
				AppliedAt: debuff.AppliedAt,
				ExpiresAt: debuff.ExpiresAt,
				Mods:      debuff.Mods,
			})
		}
	}
	
	// Capture combat counters
	if stack.Battle != nil && stack.Battle.Counters != nil {
		snapshot.AttackCount = stack.Battle.Counters.AttackCount
		snapshot.DefenseCount = stack.Battle.Counters.DefenseCount
	}
	
	// Capture effective stats for each ship type
	for shipType := range stack.Ships {
		snapshot.EffectiveStats[shipType] = CaptureEffectiveStats(stack, shipType, now)
	}
	
	return snapshot
}

// CaptureFormationSnapshot captures formation configuration
func CaptureFormationSnapshot(stack *ShipStack) *FormationSnapshot {
	if stack.Formation == nil {
		return nil
	}
	
	formation := stack.Formation.ToFormation()
	snapshot := &FormationSnapshot{
		Type:      stack.Formation.Type,
		Level:     1, // TODO: Add level tracking to FormationWithSlots if needed
		Positions: make(map[FormationPosition][]ShipAssignment),
	}
	
	// Capture ship assignments to positions
	for _, assignment := range formation.Assignments {
		if assignment.Count == 0 {
			continue
		}
		
		position := assignment.Position
		if snapshot.Positions[position] == nil {
			snapshot.Positions[position] = make([]ShipAssignment, 0)
		}
		
		snapshot.Positions[position] = append(snapshot.Positions[position], ShipAssignment{
			ShipType:    assignment.ShipType,
			BucketIndex: assignment.BucketIndex,
			Count:       assignment.Count,
			HP:          assignment.AssignedHP,
		})
	}
	
	// Capture unlocked tree nodes (if you have formation trees)
	// snapshot.TreeNodes = stack.Formation.UnlockedNodes // Implement if you have this
	
	return snapshot
}

// CaptureEffectiveStats captures base and effective stats with modifier breakdown
func CaptureEffectiveStats(stack *ShipStack, shipType ShipType, now time.Time) EffectiveShipStats {
	blueprint, ok := ShipBlueprints[shipType]
	if !ok {
		return EffectiveShipStats{ShipType: shipType}
	}
	
	stats := EffectiveShipStats{
		ShipType:              shipType,
		BaseAttackDamage:      blueprint.AttackDamage,
		BaseLaserShield:       blueprint.LaserShield,
		BaseNuclearShield:     blueprint.NuclearShield,
		BaseAntimatterShield:  blueprint.AntimatterShield,
		BaseHP:                blueprint.HP,
		BaseSpeed:             blueprint.Speed,
	}
	
	// Get effective stats with all modifiers
	var enemyFormation FormationType
	if stack.Battle != nil && len(stack.Battle.EnemyStackID) > 0 {
		// You'd need to look up enemy formation here
		// For now, we'll compute without enemy formation
	}
	
	modStack, finalMods := ComputeStackModifiers(stack, shipType, 0, now, true, enemyFormation)
	effectiveShip := ApplyStatModsToShip(blueprint, finalMods)
	
	stats.EffectiveAttackDamage = effectiveShip.AttackDamage
	stats.EffectiveLaserShield = effectiveShip.LaserShield
	stats.EffectiveNuclearShield = effectiveShip.NuclearShield
	stats.EffectiveAntimatterShield = effectiveShip.AntimatterShield
	stats.EffectiveHP = effectiveShip.HP
	stats.EffectiveSpeed = effectiveShip.Speed
	
	// Build modifier breakdown
	stats.Modifiers = BuildModifierBreakdown(modStack)
	
	return stats
}

// BuildModifierBreakdown organizes modifiers by source type
func BuildModifierBreakdown(modStack *ModifierStack) ModifierBreakdown {
	breakdown := ModifierBreakdown{
		Formation: make([]ModifierSourceDetail, 0),
		Bio:       make([]ModifierSourceDetail, 0),
		Gems:      make([]ModifierSourceDetail, 0),
		Buffs:     make([]ModifierSourceDetail, 0),
		Debuffs:   make([]ModifierSourceDetail, 0),
	}
	
	for _, layer := range modStack.Layers {
		source := ModifierSourceDetail{
			SourceID:    layer.SourceID,
			Description: layer.Description,
			Mods:        layer.Mods,
		}
		
		switch layer.Source {
		case SourceFormationPosition, SourceFormationCounter:
			breakdown.Formation = append(breakdown.Formation, source)
		case SourceBioPassive, SourceBioTriggered, SourceBioTick, SourceBioAccum:
			breakdown.Bio = append(breakdown.Bio, source)
		case SourceGem, SourceGemWord:
			breakdown.Gems = append(breakdown.Gems, source)
		case SourceBuff:
			breakdown.Buffs = append(breakdown.Buffs, source)
		case SourceDebuff, SourceBioDebuff:
			breakdown.Debuffs = append(breakdown.Debuffs, source)
		}
	}
	
	return breakdown
}

// AddBattleRound adds a new round to the battle report
func (br *BattleReport) AddBattleRound(
	attacker, defender *ShipStack,
	result FormationBattleResult,
	attackerPhase, defenderPhase CombatPhase,
	events []RoundEvent,
	now time.Time,
) {
	roundNumber := len(br.Rounds) + 1
	
	// Capture pre-round state (current state before this round)
	preAttacker := CaptureCombatantState(attacker)
	preDefender := CaptureCombatantState(defender)
	
	round := BattleRound{
		RoundNumber:      roundNumber,
		Timestamp:        now,
		AttackerPreRound: preAttacker,
		DefenderPreRound: preDefender,
		AttackerPhase:    attackerPhase,
		DefenderPhase:    defenderPhase,
		AttackerDamageDealt: result.AttackerDamageDealt,
		DefenderDamageDealt: result.DefenderDamageDealt,
		AttackerShipsLost:   result.AttackerShipsLost,
		DefenderShipsLost:   result.DefenderShipsLost,
		Events:              events,
	}
	
	// Capture post-round state (after damage applied)
	round.AttackerPostRound = CaptureCombatantState(attacker)
	round.DefenderPostRound = CaptureCombatantState(defender)
	
	br.Rounds = append(br.Rounds, round)
	
	// Update aggregate statistics
	br.TotalRounds = roundNumber
	br.AttackerTotalDamage += result.AttackerDamageDealt
	br.DefenderTotalDamage += result.DefenderDamageDealt
	
	for shipType, lost := range result.AttackerShipsLost {
		br.AttackerShipsLost[shipType] += lost
	}
	for shipType, lost := range result.DefenderShipsLost {
		br.DefenderShipsLost[shipType] += lost
	}
	
	// Update current snapshots
	br.AttackerCurrent = CaptureStackSnapshot(attacker, now)
	br.DefenderCurrent = CaptureStackSnapshot(defender, now)
	br.UpdatedAt = now
}

// CaptureCombatantState creates a lightweight state snapshot for round tracking
func CaptureCombatantState(stack *ShipStack) CombatantState {
	state := CombatantState{
		Ships: make(map[ShipType][]HPBucket),
	}
	
	totalShips := 0
	totalHP := 0
	for shipType, buckets := range stack.Ships {
		bucketsCopy := make([]HPBucket, len(buckets))
		copy(bucketsCopy, buckets)
		state.Ships[shipType] = bucketsCopy
		
		for _, bucket := range buckets {
			totalShips += bucket.Count
			totalHP += bucket.HP * bucket.Count
		}
	}
	state.TotalShips = totalShips
	state.TotalHP = totalHP
	
	// Capture combat counters
	if stack.Battle != nil && stack.Battle.Counters != nil {
		state.AttackCount = stack.Battle.Counters.AttackCount
		state.DefenseCount = stack.Battle.Counters.DefenseCount
	}
	
	// Capture active effects
	if stack.Bio != nil {
		for nodeID, node := range stack.Bio.Nodes {
			if node.Stage == BioStageTriggered || node.Stage == BioStageCompositeActive {
				state.ActiveBuffs = append(state.ActiveBuffs, nodeID)
			}
		}
		
		for debuffID := range stack.Bio.InboundDebuffs {
			state.ActiveDebuffs = append(state.ActiveDebuffs, debuffID)
		}
	}
	
	return state
}

// EndBattle marks the battle as ended with an outcome
func (br *BattleReport) EndBattle(outcome BattleOutcome, now time.Time) {
	br.Status = BattleStatusEnded
	br.EndedAt = &now
	br.Outcome = &outcome
	br.UpdatedAt = now
}

// CreateCombatPhase creates a combat phase record from combat execution
func CreateCombatPhase(
	attackerID, defenderID bson.ObjectID,
	attacker, defender *ShipStack,
	totalDamage int,
	damageMap map[ShipType]map[int]int,
	ctx *CombatContext,
	shipsDestroyed map[ShipType]int,
	debuffsApplied []string,
	now time.Time,
) CombatPhase {
	phase := CombatPhase{
		AttackerID:          attackerID,
		DefenderID:          defenderID,
		BaseDamage:          totalDamage,
		FormationMultiplier: ctx.FormationCounter,
		FinalDamage:         totalDamage,
		DamageByType:        ctx.AttackerDamageByType,
		DamageByShipType:    damageMap,
		ShieldMitigation:    make(map[string]ShieldMitigationDetail),
		ShipsDestroyed:      shipsDestroyed,
		DebuffsApplied:      debuffsApplied,
	}
	
	// Check for first strike
	if attacker.Battle != nil && attacker.Battle.Counters != nil {
		phase.FirstStrikeBonus = attacker.Battle.Counters.AttackCount == 1
	}
	
	// Check for crit (you'd need to track this in ExecuteFormationBattleRound)
	// phase.CriticalHit = ... 
	
	// Calculate shield mitigation details
	for attackType, typeDamage := range ctx.AttackerDamageByType {
		// You'd need to track pre-shield and post-shield damage
		// For now, this is a placeholder
		phase.ShieldMitigation[attackType] = ShieldMitigationDetail{
			AttackType:  attackType,
			RawDamage:   typeDamage,
			ShieldValue: 0, // Would need to calculate average shield for this type
			MitigatedDamage: typeDamage, // Placeholder
			MitigationPercent: 0,
		}
	}
	
	return phase
}

// CreateRoundEvent creates an event record
func CreateRoundEvent(
	eventType string,
	actorID bson.ObjectID,
	targetID bson.ObjectID,
	description string,
	data map[string]interface{},
	now time.Time,
) RoundEvent {
	return RoundEvent{
		Timestamp:   now,
		EventType:   eventType,
		ActorID:     actorID,
		TargetID:    targetID,
		Description: description,
		Data:        data,
	}
}
