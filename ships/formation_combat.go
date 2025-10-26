package ships

import "time"

// Formation combat integration helpers
// These functions demonstrate how formations integrate with the turn-based combat system.

// CombatContext holds the state for a formation-aware combat encounter.
type CombatContext struct {
	Attacker             *ShipStack
	Defender             *ShipStack
	AttackDirection      AttackDirection
	FormationCounter     float64        // Attacker's formation advantage multiplier
	Now                  time.Time      // Combat timestamp for stat calculations
	AttackerDamageByType map[string]int // Damage composition by attack type (Laser/Nuclear/Antimatter)
	AttackerShieldPierce float64        // Average shield pierce across attacker's fleet
}

// NewCombatContext initializes a combat context between two stacks.
func NewCombatContext(attacker, defender *ShipStack, now time.Time) *CombatContext {
	ctx := &CombatContext{
		Attacker:             attacker,
		Defender:             defender,
		AttackDirection:      DetermineAttackDirection(attacker, defender),
		Now:                  now,
		AttackerDamageByType: make(map[string]int),
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

	// Pre-calculate damage composition by attack type for weighted shield application
	ctx.calculateDamageComposition()

	return ctx
}

// calculateDamageComposition pre-calculates the attacker's damage by attack type and average shield pierce.
// This enables weighted shield application where each attack type is mitigated by the corresponding shield.
func (ctx *CombatContext) calculateDamageComposition() {
	totalDamage := 0
	weightedShieldPierce := 0.0

	for shipType, buckets := range ctx.Attacker.Ships {
		blueprint := ShipBlueprints[shipType]
		attackType := blueprint.AttackType

		for bucketIdx, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}

			// Get effective stats with all modifiers (formation + bio + gems)
			_, finalMods := ComputeStackModifiers(
				ctx.Attacker, shipType, bucketIdx, ctx.Now, true, ctx.Defender.Formation.Type,
			)
			effectiveShip := ApplyStatModsToShip(blueprint, finalMods)

			// Calculate damage (already includes type-specific bonuses from modifiers)
			damage := effectiveShip.AttackDamage * bucket.Count
			ctx.AttackerDamageByType[attackType] += damage
			totalDamage += damage

			// Weight shield pierce by damage contribution
			weightedShieldPierce += finalMods.ShieldPiercePct * float64(damage)
		}
	}

	// Calculate average shield pierce weighted by damage
	if totalDamage > 0 {
		ctx.AttackerShieldPierce = weightedShieldPierce / float64(totalDamage)
	}
}

// applyAsymptoticShieldMitigation applies asymptotic damage reduction based on shield value.
// Formula: damage / (1 + shieldValue * scalingFactor)
// This creates diminishing returns where shields never reach 100% mitigation.
// Example with scalingFactor=0.15:
//
//	Shield 0:  100% damage (no mitigation)
//	Shield 3:  ~69% damage (31% reduction)
//	Shield 5:  ~57% damage (43% reduction)
//	Shield 10: ~40% damage (60% reduction)
//	Shield 20: ~25% damage (75% reduction)
func applyAsymptoticShieldMitigation(damage int, shieldValue int) int {
	if shieldValue < 0 {
		shieldValue = 0 // Bio debuffs can reduce shields below 0, cap at 0
	}

	scalingFactor := 0.15 // Adjust for desired shield strength
	mitigation := 1.0 / (1.0 + float64(shieldValue)*scalingFactor)
	finalDamage := float64(damage) * mitigation

	return int(finalDamage)
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

			// Apply defender's type-specific weighted shield effectiveness
			finalDamage := ctx.applyWeightedShieldMitigation(assignmentDamage, assignment.ShipType, assignment.BucketIndex, ctx.AttackerShieldPierce)

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
			damageMap[shipType][bucketIndex] = ctx.applyWeightedShieldMitigation(bucketDamage, shipType, bucketIndex, ctx.AttackerShieldPierce)
		}
	}

	return damageMap
}

// applyWeightedShieldMitigation applies type-specific asymptotic shield mitigation.
// Each attack type (Laser/Nuclear/Antimatter) is mitigated by the corresponding shield value.
// Damage is weighted by the attacker's damage composition to properly handle mixed fleets.
// Shield pierce from attacker reduces effective shield values.
func (ctx *CombatContext) applyWeightedShieldMitigation(
	assignmentDamage int,
	defenderShipType ShipType,
	defenderBucketIndex int,
	attackerShieldPierce float64,
) int {
	// Get defender's effective shields with all modifiers
	defenderShip, _, _ := ctx.Defender.EffectiveShipInCombat(
		defenderShipType,
		defenderBucketIndex,
		ctx.Attacker.Formation.Type,
		ctx.Now,
	)

	// Calculate total attacker damage for weighting
	totalAttackerDamage := 0
	for _, dmg := range ctx.AttackerDamageByType {
		totalAttackerDamage += dmg
	}

	if totalAttackerDamage == 0 {
		return 0 // No damage to apply
	}

	// Apply type-specific shields to each damage component
	finalDamage := 0

	for attackType, typeDamage := range ctx.AttackerDamageByType {
		// Calculate this type's proportion of total damage
		proportion := float64(typeDamage) / float64(totalAttackerDamage)
		typeAssignmentDamage := int(float64(assignmentDamage) * proportion)

		// Get the appropriate shield value for this attack type
		var shieldValue int
		switch attackType {
		case "Laser":
			shieldValue = defenderShip.LaserShield
		case "Nuclear":
			shieldValue = defenderShip.NuclearShield
		case "Antimatter":
			shieldValue = defenderShip.AntimatterShield
		default:
			shieldValue = 0
		}

		// Apply shield pierce: reduce effective shield value
		// ShieldPiercePct ranges from 0.0 to 1.0 (0% to 100% pierce)
		effectiveShield := int(float64(shieldValue) * (1.0 - attackerShieldPierce))
		if effectiveShield < 0 {
			effectiveShield = 0
		}

		// Apply asymptotic mitigation
		mitigatedDamage := applyAsymptoticShieldMitigation(typeAssignmentDamage, effectiveShield)
		finalDamage += mitigatedDamage
	}

	return finalDamage
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
// This version uses deterministic mechanics (counter-based crits, evasion as damage reduction)
// and type-specific weighted shield mitigation with full cross-stack modifier support.
func ExecuteFormationBattleRound(attacker, defender *ShipStack, now time.Time) FormationBattleResult {
	result := FormationBattleResult{
		AttackerShipsLost:     make(map[ShipType]int),
		DefenderShipsLost:     make(map[ShipType]int),
		PositionEffectiveness: make(map[FormationPosition]float64),
	}

	// Initialize battle counters
	if attacker.Battle == nil {
		attacker.Battle = &BattleState{Counters: &CombatCounters{}}
	}
	if attacker.Battle.Counters == nil {
		attacker.Battle.Counters = &CombatCounters{}
	}
	if defender.Battle == nil {
		defender.Battle = &BattleState{Counters: &CombatCounters{}}
	}
	if defender.Battle.Counters == nil {
		defender.Battle.Counters = &CombatCounters{}
	}

	// Tick bio machines before combat
	attacker.TickBio(now)
	defender.TickBio(now)

	// Increment attack/defense counters
	attacker.Battle.Counters.AttackCount++
	defender.Battle.Counters.DefenseCount++

	ctx := NewCombatContext(attacker, defender, now)
	result.FormationAdvantage = ctx.FormationCounter

	// Phase 1: Attacker deals damage with deterministic mechanics
	attackerTotalDamage := calculateStackDamage(attacker, defender, now, attacker.Battle.Counters.AttackCount, ctx.FormationCounter)
	result.AttackerDamageDealt = attackerTotalDamage

	// Distribute and apply damage to defender (with weighted shields, shield pierce, and evasion)
	defenderDamageMap := ctx.DistributeDamageToDefender(attackerTotalDamage)

	// Apply cross-stack modifiers: accuracy vs evasion (flat damage reduction)
	applyAccuracyVsEvasion(defenderDamageMap, attacker, defender, now)

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
		defender.Battle.Counters.AttackCount++
		attacker.Battle.Counters.DefenseCount++

		ctxReverse := NewCombatContext(defender, attacker, now)

		defenderTotalDamage := calculateStackDamage(defender, attacker, now, defender.Battle.Counters.AttackCount, ctxReverse.FormationCounter)
		result.DefenderDamageDealt = defenderTotalDamage

		attackerDamageMap := ctxReverse.DistributeDamageToDefender(defenderTotalDamage)

		// Apply cross-stack modifiers: accuracy vs evasion
		applyAccuracyVsEvasion(attackerDamageMap, defender, attacker, now)

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

	// Phase 3: Apply bio debuffs post-combat for next round
	applyBioDebuffsPostCombat(attacker, defender, now)

	return result
}

// calculateStackDamage computes total damage output for a stack with all modifiers applied.
// Handles first strike, crit (with configurable crit damage), and formation counter.
func calculateStackDamage(
	attacker *ShipStack,
	defender *ShipStack,
	now time.Time,
	attackCount int,
	formationCounter float64,
) int {
	totalDamage := 0

	for shipType, buckets := range attacker.Ships {
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}

			// Get effective stats and modifiers
			_, finalMods := ComputeStackModifiers(
				attacker, shipType, bucketIndex, now, true, defender.Formation.Type,
			)
			blueprint := ShipBlueprints[shipType]
			effectiveShip := ApplyStatModsToShip(blueprint, finalMods)

			baseDamage := effectiveShip.AttackDamage * bucket.Count

			// Apply deterministic first strike bonus
			if attackCount == 1 && finalMods.FirstVolleyPct > 0 {
				baseDamage = int(float64(baseDamage) * (1.0 + finalMods.FirstVolleyPct))
			}

			// Apply deterministic crit (counter-based) with configurable crit damage
			if finalMods.CritPct > 0 {
				critInterval := int(1.0 / finalMods.CritPct)
				if critInterval > 0 && attackCount%critInterval == 0 {
					// Use CritDamagePct from mods
					critDamage := finalMods.CritDamagePct

					baseDamage = int(float64(baseDamage) * (1.0 + critDamage))
				}
			}

			// Apply formation counter multiplier
			damage := int(float64(baseDamage) * formationCounter)
			totalDamage += damage
		}
	}

	return totalDamage
}

// applyAccuracyVsEvasion applies cross-stack accuracy vs evasion mechanics.
// Attacker's accuracy reduces defender's evasion. Negative accuracy (from debuffs) increases attacker's evasion.
// Evasion is applied as flat damage reduction, capped at 75%.
func applyAccuracyVsEvasion(
	damageMap map[ShipType]map[int]int,
	attacker *ShipStack,
	defender *ShipStack,
	now time.Time,
) {
	// Calculate attacker's average accuracy weighted by damage
	attackerAccuracy := calculateAverageAccuracy(attacker, defender, now)

	// Apply evasion reduction to each defender bucket
	for shipType, bucketDamages := range damageMap {
		for bucketIndex, damage := range bucketDamages {
			// Get defender's effective evasion from modifiers
			_, defMods := ComputeStackModifiers(
				defender, shipType, bucketIndex, now, true, attacker.Formation.Type,
			)

			// Accuracy reduces target evasion
			effectiveEvasion := defMods.EvasionPct - attackerAccuracy
			if effectiveEvasion < 0 {
				effectiveEvasion = 0 // Evasion cannot go negative
			}

			// Evasion as flat damage reduction (not dodge chance)
			evasionMult := 1.0 - effectiveEvasion
			if evasionMult < 0.25 { // Cap at 75% reduction
				evasionMult = 0.25
			}

			damageMap[shipType][bucketIndex] = int(float64(damage) * evasionMult)
		}
	}
}

// calculateAverageAccuracy computes the attacker's average accuracy weighted by damage output.
func calculateAverageAccuracy(attacker *ShipStack, defender *ShipStack, now time.Time) float64 {
	totalDamage := 0
	weightedAccuracy := 0.0

	for shipType, buckets := range attacker.Ships {
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}

			_, finalMods := ComputeStackModifiers(
				attacker, shipType, bucketIndex, now, true, defender.Formation.Type,
			)
			blueprint := ShipBlueprints[shipType]
			effectiveShip := ApplyStatModsToShip(blueprint, finalMods)

			damage := effectiveShip.AttackDamage * bucket.Count
			totalDamage += damage

			// Weight accuracy by damage contribution
			weightedAccuracy += finalMods.AccuracyPct * float64(damage)
		}
	}

	if totalDamage > 0 {
		return weightedAccuracy / float64(totalDamage)
	}
	return 0.0
}

// applyBioDebuffsPostCombat applies outgoing bio debuffs from both stacks after combat.
// This affects the next hourly combat round.
func applyBioDebuffsPostCombat(attacker, defender *ShipStack, now time.Time) {
	// Apply attacker's outgoing debuffs to defender
	if attacker.Bio != nil {
		for _, node := range attacker.Bio.Nodes {
			if node.OutgoingDebuffID != "" && (node.Stage == BioStageTriggered || node.Stage == BioStageCompositeActive) {
				defender.BioApplyInboundDebuff(
					node.OutgoingDebuffID,
					node.OutgoingDebuffMods,
					node.OutgoingDebuffDuration,
					1, // Apply 1 stack per combat round
					node.OutgoingDebuffMaxStacks,
					attacker.ID,
					node.ID,
					now,
				)
			}
		}
	}

	// Apply defender's outgoing debuffs to attacker
	if defender.Bio != nil {
		for _, node := range defender.Bio.Nodes {
			if node.OutgoingDebuffID != "" && (node.Stage == BioStageTriggered || node.Stage == BioStageCompositeActive) {
				attacker.BioApplyInboundDebuff(
					node.OutgoingDebuffID,
					node.OutgoingDebuffMods,
					node.OutgoingDebuffDuration,
					1,
					node.OutgoingDebuffMaxStacks,
					defender.ID,
					node.ID,
					now,
				)
			}
		}
	}
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
