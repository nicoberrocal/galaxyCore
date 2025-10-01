package ships

import (
	"time"
)

// ComputeLoadoutV2 is the new version that uses the layered modifier system.
// It builds a complete modifier stack from all contributing sources and resolves it.
// This replaces the simpler ComputeLoadout function with full transparency and control.
func ComputeLoadoutV2(
	ship Ship,
	role RoleMode,
	loadout ShipLoadout,
	formation *Formation,
	position FormationPosition,
	ships map[ShipType][]HPBucket,
	now time.Time,
	inCombat bool,
) (*ModifierStack, StatMods, []AbilityID) {
	builder := NewModifierBuilder(now)
	
	// 1. Add gem modifiers (permanent while socketed)
	builder.AddGemsFromLoadout(loadout)
	
	// 2. Add role mode modifiers (semi-permanent until switched)
	builder.AddRoleMode(role)
	
	// 3. Add formation position bonuses (fixed while formation active)
	if formation != nil {
		builder.AddFormationPosition(formation, position)
		builder.AddFormationRoleSynergy(formation, position, role)
		builder.AddGemPositionSynergy(loadout.Sockets, position)
	}
	
	// 4. Add composition bonuses (based on fleet makeup)
	if ships != nil {
		builder.AddCompositionBonus(ships)
	}
	
	// 5. Add anchored penalty if applicable
	builder.AddAnchoredPenalty(loadout.Anchored)
	
	// Build the stack
	stack := builder.Build()
	
	// Resolve to final mods
	ctx := ResolveContext{
		Now:          now,
		InCombat:     inCombat,
		HasFormation: formation != nil,
	}
	if formation != nil {
		ctx.FormationType = formation.Type
	}
	
	finalMods := stack.Resolve(ctx)
	
	// Collect granted abilities
	_, grants, _ := EvaluateGemSockets(loadout.Sockets)
	
	return stack, finalMods, grants
}

// ComputeStackModifiers computes the complete modifier stack for a ship stack.
// This is the primary entry point for getting all modifiers affecting a stack.
func ComputeStackModifiers(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
	inCombat bool,
	enemyFormation FormationType,
) (*ModifierStack, StatMods) {
	loadout := stack.GetOrInitLoadout(shipType)
	position := stack.GetFormationPosition(shipType, bucketIndex)
	
	builder := NewModifierBuilder(now)
	
	// 1. Gems (permanent)
	builder.AddGemsFromLoadout(loadout)
	
	// 2. Role mode (semi-permanent)
	builder.AddRoleMode(stack.Role)
	
	// 3. Formation (conditional on formation being active)
	if stack.Formation != nil {
		builder.AddFormationPosition(stack.Formation, position)
		builder.AddFormationRoleSynergy(stack.Formation, position, stack.Role)
		builder.AddGemPositionSynergy(loadout.Sockets, position)
		
		// Formation counter (only in combat)
		if inCombat && enemyFormation != "" {
			builder.AddFormationCounter(stack.Formation.Type, enemyFormation, inCombat)
		}
	}
	
	// 4. Composition bonuses
	builder.AddCompositionBonus(stack.Ships)
	
	// 5. Anchored penalty
	builder.AddAnchoredPenalty(loadout.Anchored)
	
	// 6. Active abilities (temporary)
	if stack.Ability != nil {
		for _, abilityState := range *stack.Ability {
			if abilityState.IsActive && abilityState.ShipType == shipType {
				// Convert ability bonus to StatMods
				// This is a simplified conversion - you may want to expand this
				mods := abilityBonusToStatMods(abilityState.Bonus)
				duration := time.Duration(abilityState.Duration) * time.Second
				builder.AddAbility(AbilityID(abilityState.Ability), mods, duration)
			}
		}
	}
	
	modStack := builder.Build()
	
	// Resolve context
	ctx := ResolveContext{
		Now:          now,
		InCombat:     inCombat,
		HasFormation: stack.Formation != nil,
	}
	if stack.Formation != nil {
		ctx.FormationType = stack.Formation.Type
		ctx.EnemyFormation = enemyFormation
	}
	
	finalMods := modStack.Resolve(ctx)
	
	return modStack, finalMods
}

// ComputeEffectiveShipV2 computes effective ship stats using the new modifier system.
// This is the V2 replacement for EffectiveShipInFormation.
func ComputeEffectiveShipV2(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
	inCombat bool,
	enemyFormation FormationType,
) (Ship, []Ability, *ModifierStack) {
	blueprint, ok := ShipBlueprints[shipType]
	if !ok {
		return Ship{}, nil, NewModifierStack()
	}
	
	// Get modifier stack and final mods
	modStack, finalMods := ComputeStackModifiers(stack, shipType, bucketIndex, now, inCombat, enemyFormation)
	
	// Apply mods to ship
	effectiveShip := ApplyStatModsToShip(blueprint, finalMods)
	
	// Get abilities
	loadout := stack.GetOrInitLoadout(shipType)
	_, grants, _ := EvaluateGemSockets(loadout.Sockets)
	abilities := FilterAbilitiesForMode(effectiveShip, stack.Role, grants)
	
	return effectiveShip, abilities, modStack
}

// GetModifierBreakdown returns a detailed breakdown of all modifiers for debugging/UI.
func GetModifierBreakdown(
	stack *ShipStack,
	shipType ShipType,
	bucketIndex int,
	now time.Time,
	inCombat bool,
	enemyFormation FormationType,
) []ModifierSummary {
	modStack, _ := ComputeStackModifiers(stack, shipType, bucketIndex, now, inCombat, enemyFormation)
	
	ctx := ResolveContext{
		Now:          now,
		InCombat:     inCombat,
		HasFormation: stack.Formation != nil,
	}
	if stack.Formation != nil {
		ctx.FormationType = stack.Formation.Type
		ctx.EnemyFormation = enemyFormation
	}
	
	return modStack.GetSummary(ctx)
}

// abilityBonusToStatMods converts ability bonus map to StatMods.
// This is a helper to bridge the existing AbilityState.Bonus format.
func abilityBonusToStatMods(bonus map[string]int) StatMods {
	mods := ZeroMods()
	
	for key, value := range bonus {
		floatVal := float64(value) / 100.0 // Assuming bonus is in percentage points
		
		switch key {
		case "speed":
			mods.SpeedDelta = value
		case "defense":
			mods.LaserShieldDelta = value / 3
			mods.NuclearShieldDelta = value / 3
			mods.AntimatterShieldDelta = value / 3
		case "attack":
			mods.Damage.LaserPct = floatVal
			mods.Damage.NuclearPct = floatVal
			mods.Damage.AntimatterPct = floatVal
		case "visibility":
			mods.VisibilityDelta = value
		case "accuracy":
			mods.AccuracyPct = floatVal
		case "evasion":
			mods.EvasionPct = floatVal
		}
	}
	
	return mods
}

// UpdateStackModifiers is a helper to refresh a stack's modifier state.
// Call this when gems are socketed/unsocketed, formation changes, etc.
func UpdateStackModifiers(stack *ShipStack, now time.Time) {
	// Remove expired temporary modifiers
	// In the future, you might store the modifier stack on the ShipStack itself
	// For now, this is a placeholder for when you integrate it into the data model
	
	// This function would be expanded when you add a ModifierStack field to ShipStack
}

// CompareModifierStacks compares two modifier stacks and returns the differences.
// Useful for showing "before/after" when changing equipment, formations, etc.
type ModifierDiff struct {
	Added   []ModifierLayer
	Removed []ModifierLayer
	Changed []ModifierLayer
}

// DiffModifierStacks computes the difference between two modifier stacks.
func DiffModifierStacks(before, after *ModifierStack) ModifierDiff {
	diff := ModifierDiff{
		Added:   []ModifierLayer{},
		Removed: []ModifierLayer{},
		Changed: []ModifierLayer{},
	}
	
	// Build maps for comparison
	beforeMap := make(map[string]ModifierLayer)
	afterMap := make(map[string]ModifierLayer)
	
	for _, layer := range before.Layers {
		key := string(layer.Source) + ":" + layer.SourceID
		beforeMap[key] = layer
	}
	
	for _, layer := range after.Layers {
		key := string(layer.Source) + ":" + layer.SourceID
		afterMap[key] = layer
	}
	
	// Find added and changed
	for key, afterLayer := range afterMap {
		if beforeLayer, exists := beforeMap[key]; exists {
			// Check if changed (simplified comparison)
			if !modsEqual(beforeLayer.Mods, afterLayer.Mods) {
				diff.Changed = append(diff.Changed, afterLayer)
			}
		} else {
			diff.Added = append(diff.Added, afterLayer)
		}
	}
	
	// Find removed
	for key, beforeLayer := range beforeMap {
		if _, exists := afterMap[key]; !exists {
			diff.Removed = append(diff.Removed, beforeLayer)
		}
	}
	
	return diff
}

// modsEqual checks if two StatMods are equal (simplified).
func modsEqual(a, b StatMods) bool {
	// This is a simplified comparison - you might want a more thorough check
	return a.Damage.LaserPct == b.Damage.LaserPct &&
		a.Damage.NuclearPct == b.Damage.NuclearPct &&
		a.Damage.AntimatterPct == b.Damage.AntimatterPct &&
		a.SpeedDelta == b.SpeedDelta &&
		a.AttackIntervalPct == b.AttackIntervalPct
	// Add more fields as needed
}
