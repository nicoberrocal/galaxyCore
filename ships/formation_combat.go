package ships

// Formation combat integration helpers
// These functions demonstrate how formations integrate with the turn-based combat system.

// CombatContext holds the state for a formation-aware combat encounter.
type CombatContext struct {
	Attacker         *ShipStack
	Defender         *ShipStack
	AttackDirection  AttackDirection
	FormationCounter float64 // Attacker's formation advantage multiplier
}

// NewCombatContext initializes a combat context between two stacks.
func NewCombatContext(attacker, defender *ShipStack) *CombatContext {
	ctx := &CombatContext{
		Attacker:        attacker,
		Defender:        defender,
		AttackDirection: DetermineAttackDirection(attacker, defender),
	}

	// Calculate formation counter multiplier
	if attacker.Formation != nil && defender.Formation != nil {
		ctx.FormationCounter = GetFormationCounterMultiplier(
			attacker.Formation.Type,
			defender.Formation.Type,
		)
	} else {
		ctx.FormationCounter = 1.0
	}

	return ctx
}

// DetermineAttackDirection calculates the attack angle based on stack positioning and formations.
func DetermineAttackDirection(attacker, defender *ShipStack) AttackDirection {
	// Default to frontal if no formations
	if attacker.Formation == nil || defender.Formation == nil {
		return DirectionFrontal
	}

	// Simple direction logic based on formation types
	// In a full implementation, this would consider map positioning, facing, and movement vectors
	attackerType := attacker.Formation.Type
	defenderType := defender.Formation.Type

	// Skirmish formations tend to flank
	if attackerType == FormationSkirmish {
		return DirectionFlanking
	}

	// Vanguard charges frontally
	if attackerType == FormationVanguard {
		return DirectionFrontal
	}

	// Swarm formations envelop
	if attackerType == FormationSwarm && defenderType != FormationSwarm {
		return DirectionEnvelopment
	}

	// Default to frontal assault
	return DirectionFrontal
}

// CalculateFormationDamage computes damage output with formation bonuses applied.
func (ctx *CombatContext) CalculateFormationDamage(baseAttackerDamage int, attackerShipType ShipType, attackerBucketIndex int) int {
	damage := float64(baseAttackerDamage)

	// Apply formation counter multiplier
	damage *= ctx.FormationCounter

	// Apply attacker's formation position bonuses
	if ctx.Attacker.Formation != nil {
		// Get effective ship stats with formation bonuses
		effectiveShip, _ := ctx.Attacker.EffectiveShipInFormation(attackerShipType, attackerBucketIndex)

		// Use the effective ship's attack damage (already includes all bonuses)
		damage = float64(effectiveShip.AttackDamage)
	}

	return int(damage)
}

// DistributeDamageToDefender distributes incoming damage across the defender's formation.
func (ctx *CombatContext) DistributeDamageToDefender(totalDamage int) map[ShipType]map[int]int {
	damageMap := make(map[ShipType]map[int]int)

	// If defender has no formation, distribute evenly
	if ctx.Defender.Formation == nil {
		return ctx.distributeEvenlyToDefender(totalDamage)
	}
	formation := ctx.Defender.Formation.ToFormation()
	// Calculate positional damage distribution
	positionDamage := formation.CalculateDamageDistribution(totalDamage, ctx.AttackDirection)

	// Distribute damage within each position to specific buckets
	for position, damage := range positionDamage {
		assignments := formation.GetAssignmentsByPosition(position)

		for _, assignment := range assignments {
			if assignment.Count == 0 || assignment.AssignedHP == 0 {
				continue
			}

			// Calculate how much damage this assignment takes
			assignmentDamage := CalculateAssignmentDamage(damage, assignment, assignments)

			// Apply defender's shield effectiveness
			finalDamage := ctx.applyShieldMitigation(assignmentDamage, assignment.ShipType)

			// Record damage for this ship type and bucket
			if damageMap[assignment.ShipType] == nil {
				damageMap[assignment.ShipType] = make(map[int]int)
			}
			damageMap[assignment.ShipType][assignment.BucketIndex] += finalDamage
		}
	}

	return damageMap
}

// distributeEvenlyToDefender distributes damage evenly when no formation is present.
func (ctx *CombatContext) distributeEvenlyToDefender(totalDamage int) map[ShipType]map[int]int {
	damageMap := make(map[ShipType]map[int]int)

	// Count total ships
	totalShips := 0
	for _, buckets := range ctx.Defender.Ships {
		for _, bucket := range buckets {
			totalShips += bucket.Count
		}
	}

	if totalShips == 0 {
		return damageMap
	}

	// Distribute proportionally
	for shipType, buckets := range ctx.Defender.Ships {
		damageMap[shipType] = make(map[int]int)
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}
			bucketProportion := float64(bucket.Count) / float64(totalShips)
			bucketDamage := int(float64(totalDamage) * bucketProportion)
			damageMap[shipType][bucketIndex] = ctx.applyShieldMitigation(bucketDamage, shipType)
		}
	}

	return damageMap
}

// applyShieldMitigation applies shield effectiveness based on attack type.
func (ctx *CombatContext) applyShieldMitigation(damage int, defenderShipType ShipType) int {
	// Get effective ship with formation bonuses
	effectiveShip, _ := ctx.Defender.EffectiveShip(defenderShipType)

	// Shield effectiveness reduces damage
	// This is a simplified model - full implementation would check attack type vs shield type
	avgShield := (effectiveShip.LaserShield + effectiveShip.NuclearShield + effectiveShip.AntimatterShield) / 3
	mitigation := float64(avgShield) * 0.05 // 5% reduction per shield point

	finalDamage := float64(damage) * (1.0 - mitigation)
	if finalDamage < 0 {
		finalDamage = 0
	}

	return int(finalDamage)
}

// ApplyDamageToStack applies the calculated damage to the defender's HP buckets.
func ApplyDamageToStack(defender *ShipStack, damageMap map[ShipType]map[int]int) {
	for shipType, bucketDamages := range damageMap {
		buckets, ok := defender.Ships[shipType]
		if !ok {
			continue
		}

		for bucketIndex, damage := range bucketDamages {
			if bucketIndex >= len(buckets) {
				continue
			}

			bucket := &buckets[bucketIndex]
			totalBucketHP := bucket.HP * bucket.Count

			// Apply damage
			totalBucketHP -= damage

			if totalBucketHP <= 0 {
				// Bucket destroyed
				bucket.HP = 0
				bucket.Count = 0
			} else {
				// Recalculate bucket after damage
				blueprint := ShipBlueprints[shipType]
				fullHP := blueprint.HP

				// How many full HP ships remain?
				fullShips := totalBucketHP / fullHP
				remainderHP := totalBucketHP % fullHP

				if remainderHP > 0 {
					// We have a partial HP ship
					bucket.Count = fullShips + 1
					bucket.HP = remainderHP
				} else {
					// All ships at full HP
					bucket.Count = fullShips
					bucket.HP = fullHP
				}
			}
		}

		// Update the buckets array
		defender.Ships[shipType] = buckets
	}

	// Update formation assignments to reflect new HP values
	defender.UpdateFormationAssignments()
}

// FormationBattleResult summarizes the outcome of a formation-aware battle round.
type FormationBattleResult struct {
	AttackerDamageDealt   int
	DefenderDamageDealt   int
	AttackerShipsLost     map[ShipType]int
	DefenderShipsLost     map[ShipType]int
	FormationAdvantage    float64                       // Attacker's formation counter multiplier
	PositionEffectiveness map[FormationPosition]float64 // How effective each position was
}

// ExecuteFormationBattleRound performs one round of turn-based combat with formations.
func ExecuteFormationBattleRound(attacker, defender *ShipStack) FormationBattleResult {
	result := FormationBattleResult{
		AttackerShipsLost:     make(map[ShipType]int),
		DefenderShipsLost:     make(map[ShipType]int),
		PositionEffectiveness: make(map[FormationPosition]float64),
	}

	ctx := NewCombatContext(attacker, defender)
	result.FormationAdvantage = ctx.FormationCounter

	// Phase 1: Attacker deals damage
	attackerTotalDamage := 0
	for shipType, buckets := range attacker.Ships {
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}

			// Calculate damage with formation bonuses
			baseDamage := ShipBlueprints[shipType].AttackDamage * bucket.Count
			damage := ctx.CalculateFormationDamage(baseDamage, shipType, bucketIndex)
			attackerTotalDamage += damage
		}
	}

	result.AttackerDamageDealt = attackerTotalDamage

	// Distribute and apply damage to defender
	defenderDamageMap := ctx.DistributeDamageToDefender(attackerTotalDamage)
	defenderShipsBeforeDamage := countShips(defender.Ships)
	ApplyDamageToStack(defender, defenderDamageMap)
	defenderShipsAfterDamage := countShips(defender.Ships)

	// Calculate ships lost
	for shipType := range defenderShipsBeforeDamage {
		lost := defenderShipsBeforeDamage[shipType] - defenderShipsAfterDamage[shipType]
		if lost > 0 {
			result.DefenderShipsLost[shipType] = lost
		}
	}

	// Phase 2: Defender returns fire (if still alive)
	if !isStackDestroyed(defender) {
		ctxReverse := NewCombatContext(defender, attacker)

		defenderTotalDamage := 0
		for shipType, buckets := range defender.Ships {
			for bucketIndex, bucket := range buckets {
				if bucket.Count == 0 {
					continue
				}

				baseDamage := ShipBlueprints[shipType].AttackDamage * bucket.Count
				damage := ctxReverse.CalculateFormationDamage(baseDamage, shipType, bucketIndex)
				defenderTotalDamage += damage
			}
		}

		result.DefenderDamageDealt = defenderTotalDamage

		attackerDamageMap := ctxReverse.DistributeDamageToDefender(defenderTotalDamage)
		attackerShipsBeforeDamage := countShips(attacker.Ships)
		ApplyDamageToStack(attacker, attackerDamageMap)
		attackerShipsAfterDamage := countShips(attacker.Ships)

		for shipType := range attackerShipsBeforeDamage {
			lost := attackerShipsBeforeDamage[shipType] - attackerShipsAfterDamage[shipType]
			if lost > 0 {
				result.AttackerShipsLost[shipType] = lost
			}
		}
	}

	return result
}

// Helper functions

func countShips(ships map[ShipType][]HPBucket) map[ShipType]int {
	counts := make(map[ShipType]int)
	for shipType, buckets := range ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		counts[shipType] = total
	}
	return counts
}

func isStackDestroyed(stack *ShipStack) bool {
	for _, buckets := range stack.Ships {
		for _, bucket := range buckets {
			if bucket.Count > 0 {
				return false
			}
		}
	}
	return true
}
