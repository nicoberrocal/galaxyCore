package ships

import (
	"time"
)

/*
Simplified V2 Compute Helpers

This file provides convenient wrapper functions around the V2 modifier system
for common use cases. These helpers eliminate boilerplate while maintaining
full transparency and debuggability of the underlying modifier stack.

Use these functions when you:
- Don't need to inspect individual modifier layers
- Want quick effective ship calculations
- Are migrating from the old compute system

For advanced scenarios requiring layer inspection, use the full V2 API:
- NewModifierBuilder() for manual stack construction
- ComputeStackModifiers() for full control
- GetModifierBreakdown() for debugging
*/

// QuickEffectiveShip computes effective ship stats using V2 system with sensible defaults.
// This is the simplest way to get effective stats for a ship in a stack.
//
// Use this for:
// - UI display of ship stats
// - Quick calculations without needing the modifier stack
// - Simple scenarios outside of combat
//
// For combat calculations or when you need the modifier stack, use ComputeEffectiveShipV2.
func QuickEffectiveShip(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
) (Ship, []Ability) {
	ship, abilities, _ := ComputeEffectiveShipV2(
		stack,
		shipType,
		bucketIndex,
		now,
		false, // not in combat
		"",    // no enemy formation
	)
	return ship, abilities
}

// QuickEffectiveShipInCombat computes effective ship stats in combat context.
// This includes formation counter bonuses and combat-only modifiers.
func QuickEffectiveShipInCombat(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	enemyFormation FormationType,
	now time.Time,
) (Ship, []Ability) {
	ship, abilities, _ := ComputeEffectiveShipV2(
		stack,
		shipType,
		bucketIndex,
		now,
		true,           // in combat
		enemyFormation, // for counter bonuses
	)
	return ship, abilities
}

// QuickModifierStack builds a complete modifier stack for a ship without computing final stats.
// Returns just the stack for inspection or later resolution.
func QuickModifierStack(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
	inCombat bool,
	enemyFormation FormationType,
) *ModifierStack {
	modStack, _ := ComputeStackModifiers(stack, shipType, bucketIndex, now, inCombat, enemyFormation)
	return modStack
}

// CompareLoadoutChange shows the stat difference when changing a ship's loadout.
// Useful for UI to preview gem socketing/unsocketing.
func CompareLoadoutChange(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	newLoadout ShipLoadout,
	now time.Time,
) (before, after StatMods, diff ModifierDiff) {
	// Before: current loadout
	beforeStack, beforeMods := ComputeStackModifiers(
		stack,
		shipType,
		bucketIndex,
		now,
		false, // out of combat for comparison
		"",
	)

	// Temporarily swap loadout
	originalLoadout := stack.GetOrInitLoadout(shipType)
	if stack.Loadouts == nil {
		stack.Loadouts = make(map[ShipType]ShipLoadout)
	}
	stack.Loadouts[shipType] = newLoadout

	// After: new loadout
	afterStack, afterMods := ComputeStackModifiers(
		stack,
		shipType,
		bucketIndex,
		now,
		false,
		"",
	)

	// Restore original
	stack.Loadouts[shipType] = originalLoadout

	diff = DiffModifierStacks(beforeStack, afterStack)
	return beforeMods, afterMods, diff
}

// CompareFormationChange shows the stat difference when changing formations.
// Useful for formation selection UI.
func CompareFormationChange(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	newFormationType FormationType,
	now time.Time,
) (before, after StatMods, diff ModifierDiff) {
	// Before: current formation
	beforeStack, beforeMods := ComputeStackModifiers(
		stack,
		shipType,
		bucketIndex,
		now,
		false,
		"",
	)

	// Temporarily swap formation
	originalFormation := stack.Formation
	newFormation := AutoAssignFormation(stack.Ships, newFormationType, now)
	stack.Formation = &newFormation

	// After: new formation
	afterStack, afterMods := ComputeStackModifiers(
		stack,
		shipType,
		bucketIndex,
		now,
		false,
		"",
	)

	// Restore original
	stack.Formation = originalFormation

	diff = DiffModifierStacks(beforeStack, afterStack)
	return beforeMods, afterMods, diff
}

// GetActiveModifierSources returns a list of all active modifier sources for debugging.
// Groups modifiers by source type for easy interpretation.
func GetActiveModifierSources(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
	inCombat bool,
) map[ModifierSource][]ModifierSummary {
	breakdown := GetModifierBreakdown(stack, shipType, bucketIndex, now, inCombat, "")

	grouped := make(map[ModifierSource][]ModifierSummary)
	for _, summary := range breakdown {
		if summary.IsActive {
			grouped[summary.Source] = append(grouped[summary.Source], summary)
		}
	}

	return grouped
}

// BatchComputeEffectiveShips computes effective stats for all ships in a stack efficiently.
// Returns a map of shipType -> effective ship.
func BatchComputeEffectiveShips(
	stack *ShipStack,
	now time.Time,
	inCombat bool,
	enemyFormation FormationType,
) map[ShipType]Ship {
	result := make(map[ShipType]Ship)

	for shipType := range stack.Ships {
		ship, _, _ := ComputeEffectiveShipV2(
			stack,
			shipType,
			0, // primary bucket
			now,
			inCombat,
			enemyFormation,
		)
		result[shipType] = ship
	}

	return result
}

// GetFormationEffectiveness calculates how effective the current formation is vs an enemy.
// Returns a score from 0-2 where 1.0 is neutral, >1.0 is advantageous, <1.0 is disadvantageous.
func GetFormationEffectiveness(
	attackerFormation FormationType,
	defenderFormation FormationType,
) float64 {
	return GetFormationCounterMultiplier(attackerFormation, defenderFormation)
}

// RecommendFormation suggests the best formation type vs a specific enemy formation.
// Returns the recommended formation and its effectiveness score.
func RecommendFormation(
	stack *ShipStack,
	enemyFormation FormationType,
	now time.Time,
) (recommended FormationType, score float64) {
	bestScore := 0.0
	bestFormation := FormationLine // default

	// Try each formation type
	for formationType := range FormationCatalog {
		mult := GetFormationCounterMultiplier(formationType, enemyFormation)
		if mult > bestScore {
			bestScore = mult
			bestFormation = formationType
		}
	}

	return bestFormation, bestScore
}

// SimulateCombatModifiers previews what modifiers will be active in combat.
// Useful for pre-combat planning UI.
func SimulateCombatModifiers(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	enemyFormation FormationType,
	now time.Time,
) []ModifierSummary {
	return GetModifierBreakdown(stack, shipType, bucketIndex, now, true, enemyFormation)
}

// GetStackPowerRating calculates a rough "power level" of a stack for matchmaking.
// This is a simplified heuristic combining ship count, HP, and modifiers.
func GetStackPowerRating(stack *ShipStack, now time.Time) float64 {
	rating := 0.0

	for shipType, buckets := range stack.Ships {
		for bucketIndex, bucket := range buckets {
			if bucket.Count == 0 {
				continue
			}

			// Get effective ship stats
			effectiveShip, _ := QuickEffectiveShip(stack, shipType, bucketIndex, now)

			// Simple power formula: (HP * Count * AttackDamage) / AttackInterval
			shipPower := float64(effectiveShip.HP * bucket.Count * effectiveShip.AttackDamage)
			if effectiveShip.AttackInterval > 0 {
				shipPower /= effectiveShip.AttackInterval
			}

			rating += shipPower
		}
	}

	return rating
}

// ValidateLoadout checks if a loadout is valid (socket count, gem compatibility, etc.).
// Returns error message if invalid, empty string if valid.
func ValidateLoadout(loadout ShipLoadout, shipType ShipType) string {
	// Check socket count (max 3)
	if len(loadout.Sockets) > 3 {
		return "maximum 3 sockets allowed"
	}

	// Check for duplicate gems (same ID)
	seen := make(map[GemID]bool)
	for _, gem := range loadout.Sockets {
		if seen[gem.ID] {
			return "duplicate gem detected: " + string(gem.ID)
		}
		seen[gem.ID] = true
	}

	// Could add more validation here:
	// - Ship type restrictions
	// - Gem tier requirements
	// - etc.

	return "" // valid
}

// CleanupExpiredModifiers removes expired temporary modifiers from a modifier stack.
// This should be called periodically to prevent memory bloat.
func CleanupExpiredModifiers(modStack *ModifierStack, now time.Time) {
	modStack.RemoveExpired(now)
}
