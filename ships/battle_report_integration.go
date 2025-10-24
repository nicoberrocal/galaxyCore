package ships

import (
	"time"

	bson "go.mongodb.org/mongo-driver/v2/bson"
)

// Example integration with your tick system

// ProcessCombatWithReporting executes a combat round and updates the battle report
func ProcessCombatWithReporting(
	attacker, defender *ShipStack,
	report *BattleReport,
	now time.Time,
) (*BattleReport, FormationBattleResult) {
	// Track events that occur during this round
	events := make([]RoundEvent, 0)
	
	// Execute the combat round
	result := ExecuteFormationBattleRound(attacker, defender, now)
	
	// Create combat context for detailed tracking
	ctx := NewCombatContext(attacker, defender, now)
	
	// Track first strike event
	if attacker.Battle.Counters.AttackCount == 1 {
		events = append(events, CreateRoundEvent(
			"first_strike",
			attacker.ID,
			defender.ID,
			"Attacker unleashes first strike bonus",
			map[string]interface{}{
				"damageBonus": "30%",
			},
			now,
		))
	}
	
	// Track crit events (you'd need to track this in ExecuteFormationBattleRound)
	// if critOccurred {
	// 	events = append(events, CreateRoundEvent(
	// 		"critical_hit",
	// 		attacker.ID,
	// 		defender.ID,
	// 		"Critical hit! +50% damage",
	// 		map[string]interface{}{
	// 			"damageMultiplier": 1.5,
	// 		},
	// 		now,
	// 	))
	// }
	
	// Track ship destruction events
	for shipType, lost := range result.DefenderShipsLost {
		if lost > 0 {
			events = append(events, CreateRoundEvent(
				"ships_destroyed",
				attacker.ID,
				defender.ID,
				"Ships destroyed",
				map[string]interface{}{
					"shipType": string(shipType),
					"count":    lost,
				},
				now,
			))
		}
	}
	
	// Track bio debuff applications
	debuffsApplied := make([]string, 0)
	if attacker.Bio != nil {
		for _, node := range attacker.Bio.Nodes {
			if node.OutgoingDebuffID != "" && (node.Stage == BioStageTriggered || node.Stage == BioStageCompositeActive) {
				debuffsApplied = append(debuffsApplied, node.OutgoingDebuffID)
				events = append(events, CreateRoundEvent(
					"debuff_applied",
					attacker.ID,
					defender.ID,
					"Bio debuff applied: "+node.OutgoingDebuffID,
					map[string]interface{}{
						"debuffId": node.OutgoingDebuffID,
						"nodeId":   node.ID,
					},
					now,
				))
			}
		}
	}
	
	// Create combat phases (simplified - you'd track more details in actual combat)
	attackerPhase := CreateCombatPhase(
		attacker.ID,
		defender.ID,
		attacker,
		defender,
		result.AttackerDamageDealt,
		nil, // You'd pass the actual damage map here
		ctx,
		result.DefenderShipsLost,
		debuffsApplied,
		now,
	)
	
	defenderPhase := CombatPhase{
		AttackerID:      defender.ID,
		DefenderID:      attacker.ID,
		FinalDamage:     result.DefenderDamageDealt,
		ShipsDestroyed:  result.AttackerShipsLost,
	}
	
	// Add round to report
	report.AddBattleRound(attacker, defender, result, attackerPhase, defenderPhase, events, now)
	
	// Check if battle should end
	if isStackDestroyed(defender) {
		report.EndBattle(BattleOutcome{
			Victor:        "attacker",
			VictorStackID: attacker.ID,
			Reason:        "total_destruction",
			EndedAt:       now,
		}, now)
	} else if isStackDestroyed(attacker) {
		report.EndBattle(BattleOutcome{
			Victor:        "defender",
			VictorStackID: defender.ID,
			Reason:        "total_destruction",
			EndedAt:       now,
		}, now)
	}
	
	return report, result
}

// InitiateBattle creates a new battle and report when combat begins
func InitiateBattle(
	attacker, defender *ShipStack,
	location BattleLocation,
	now time.Time,
) *BattleReport {
	// Initialize battle state on both stacks
	if attacker.Battle == nil {
		attacker.Battle = &BattleState{
			Counters: &CombatCounters{},
		}
	}
	if defender.Battle == nil {
		defender.Battle = &BattleState{
			Counters: &CombatCounters{},
		}
	}
	
	attacker.Battle.IsInCombat = true
	attacker.Battle.EnemyStackID = []bson.ObjectID{defender.ID}
	attacker.Battle.EnemyPlayerID = []bson.ObjectID{defender.PlayerID}
	attacker.Battle.BattleStartedAt = now
	attacker.Battle.BattleLocation = location.Type
	attacker.Battle.LocationID = location.LocationID
	
	defender.Battle.IsInCombat = true
	defender.Battle.EnemyStackID = []bson.ObjectID{attacker.ID}
	defender.Battle.EnemyPlayerID = []bson.ObjectID{attacker.PlayerID}
	defender.Battle.BattleStartedAt = now
	defender.Battle.BattleLocation = location.Type
	defender.Battle.LocationID = location.LocationID
	
	// Create battle report
	report := NewBattleReport(attacker, defender, location, now)
	
	return report
}

// Example: Your main tick system integration
func ExampleTickSystemIntegration(gameState interface{}, now time.Time) {
	// Pseudocode for your tick system
	
	// 1. Find all active battles
	// activeBattles := gameState.GetActiveBattles()
	
	// 2. Process each battle
	// for _, battle := range activeBattles {
	// 	attacker := GetStack(battle.AttackerStackID)
	// 	defender := GetStack(battle.DefenderStackID)
	// 	
	// 	// Get or create battle report
	// 	report := GetBattleReport(battle.BattleID)
	// 	if report == nil {
	// 		location := BattleLocation{
	// 			Type: battle.Location,
	// 			X:    attacker.X,
	// 			Y:    attacker.Y,
	// 		}
	// 		report = InitiateBattle(attacker, defender, location, now)
	// 	}
	// 	
	// 	// Execute combat round with reporting
	// 	report, result := ProcessCombatWithReporting(attacker, defender, report, now)
	// 	
	// 	// Save updated report to database
	// 	SaveBattleReport(report)
	// 	
	// 	// Save updated stacks
	// 	SaveStack(attacker)
	// 	SaveStack(defender)
	// 	
	// 	// If battle ended, clean up
	// 	if report.Status == BattleStatusEnded {
	// 		CleanupBattle(battle)
	// 	}
	// }
}

// GetBattleReportForStack retrieves all battle reports for a stack
// This handles the case where a stack is attacked by multiple enemies
func GetBattleReportForStack(stackID bson.ObjectID, now time.Time) []*BattleReport {
	// Query your database for all reports where:
	// (AttackerStackID == stackID OR DefenderStackID == stackID) AND Status == "ongoing"
	
	// Pseudocode:
	// reports := database.Find(BattleReport{
	// 	$or: [
	// 		{AttackerStackID: stackID},
	// 		{DefenderStackID: stackID},
	// 	],
	// 	Status: "ongoing",
	// })
	
	// return reports
	return nil // Placeholder
}

// CreateBattleReportSummary generates a human-readable summary
func CreateBattleReportSummary(report *BattleReport) string {
	// Example summary generation
	summary := "Battle Report\n"
	summary += "=============\n\n"
	
	summary += "Duration: " + report.StartedAt.Format("2006-01-02 15:04:05") + "\n"
	if report.EndedAt != nil {
		summary += "Ended: " + report.EndedAt.Format("2006-01-02 15:04:05") + "\n"
	}
	summary += "\n"
	
	summary += "Initial Forces:\n"
	summary += "  Attacker: " + string(report.AttackerInitial.TotalShips) + " ships\n"
	summary += "  Defender: " + string(report.DefenderInitial.TotalShips) + " ships\n"
	summary += "\n"
	
	summary += "Total Rounds: " + string(report.TotalRounds) + "\n"
	summary += "Total Damage Dealt:\n"
	summary += "  Attacker: " + string(report.AttackerTotalDamage) + "\n"
	summary += "  Defender: " + string(report.DefenderTotalDamage) + "\n"
	summary += "\n"
	
	if report.Outcome != nil {
		summary += "Outcome: " + report.Outcome.Victor + " victory\n"
		summary += "Reason: " + report.Outcome.Reason + "\n"
	}
	
	return summary
}
