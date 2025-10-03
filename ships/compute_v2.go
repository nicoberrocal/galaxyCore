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
	
	// 1. Gems: provide their own StatMods from gem properties
	builder.AddGemsFromLoadout(loadout)
	
	// 2. Role Mode: provides its own StatMods
	builder.AddRoleMode(role)
	
	// 3. Formation: provides StatMods from FormationCatalog position bonuses only
	if formation != nil {
		builder.AddFormationPosition(formation, position)
	}
	
	// 4. Anchored state: provides penalty mods
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
	
	// 1. Gems: provide their own StatMods
	builder.AddGemsFromLoadout(loadout)
	
	// 2. Role Mode: provides its own StatMods
	builder.AddRoleMode(stack.Role)
	
	// 3. Formation: provides StatMods from FormationCatalog + tree nodes
	if stack.Formation != nil {
		builder.AddFormationPosition(stack.Formation, position)
		
		// Formation counter (only in combat)
		if inCombat && enemyFormation != "" {
			builder.AddFormationCounter(stack.Formation.Type, enemyFormation, inCombat)
		}
	}
	
	// 4. Anchored state: provides penalty mods
	builder.AddAnchoredPenalty(loadout.Anchored)
	
	// 5. Abilities: provide their own StatMods when active
	if stack.Ability != nil {
		for _, abilityState := range *stack.Ability {
			if abilityState.IsActive && abilityState.ShipType == shipType {
				// Get ability mods from catalog
				mods := GetAbilityMods(AbilityID(abilityState.Ability))
				if !isZeroMods(mods) {
					duration := time.Duration(abilityState.Duration) * time.Second
					builder.AddAbility(AbilityID(abilityState.Ability), mods, duration)
				}
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

// ===========================================
// Stat Application Functions (from compute.go)
// ===========================================

// DamageMultiplierFor returns the damage multiplier for the ship's current attack type.
// Caller applies this to base AttackDamage when displaying effective damage.
func DamageMultiplierFor(s Ship, mods StatMods) float64 {
	switch s.AttackType {
	case "Laser":
		return 1.0 + mods.Damage.LaserPct
	case "Nuclear":
		return 1.0 + mods.Damage.NuclearPct
	case "Antimatter":
		return 1.0 + mods.Damage.AntimatterPct
	default:
		return 1.0
	}
}

// EffectiveAttackInterval applies AttackIntervalPct to the base interval and returns the result.
func EffectiveAttackInterval(s Ship, mods StatMods) float64 {
	v := s.AttackInterval * (1.0 + mods.AttackIntervalPct)
	if v < 0.1 { // safety clamp
		v = 0.1
	}
	return v
}

// ApplyStatModsToShip computes a presentational "effective" Ship snapshot by applying StatMods.
// Note: This does not persist or mutate runtime state; it's for UI calculations.
func ApplyStatModsToShip(base Ship, mods StatMods) Ship {
	s := base
	s.Speed += mods.SpeedDelta
	s.VisibilityRange += mods.VisibilityDelta
	s.AttackRange += mods.AttackRangeDelta

	s.LaserShield += mods.LaserShieldDelta
	s.NuclearShield += mods.NuclearShieldDelta
	s.AntimatterShield += mods.AntimatterShieldDelta

	// Damage is multiplicative and type-dependent; update AttackDamage accordingly
	s.AttackDamage = int(float64(s.AttackDamage) * DamageMultiplierFor(base, mods))
	s.AttackInterval = EffectiveAttackInterval(base, mods)
	// BucketHPPct modifies per-bucket HP; we reflect on base HP for preview purposes only
	s.HP = int(float64(s.HP) * (1.0 + mods.BucketHPPct))
	// Transport capacity percentage
	s.TransportCapacity = int(float64(s.TransportCapacity) * (1.0 + mods.TransportCapacityPct))
	return s
}

// FilterAbilitiesForMode returns the abilities usable in the stack's current RoleMode.
// It takes the ship's built-in abilities, adds GemWord-granted abilities, then
// applies Disabled/Enabled lists from RoleModesCatalog.
func FilterAbilitiesForMode(s Ship, role RoleMode, runewordGrants []AbilityID) []Ability {
	spec, ok := RoleModesCatalog[role]
	if !ok {
		// Unknown mode, return baseline abilities only
		base := make([]Ability, 0, len(s.Abilities))
		base = append(base, s.Abilities...)
		return base
	}

	// Build a set for disabled IDs for quick lookup
	disabled := make(map[AbilityID]struct{}, len(spec.DisabledAbilities))
	for _, id := range spec.DisabledAbilities {
		disabled[id] = struct{}{}
	}

	// De-dup base abilities by ID
	out := make([]Ability, 0, len(s.Abilities)+len(runewordGrants)+len(spec.EnabledAbilities))
	seen := map[AbilityID]struct{}{}

	// Helper to append if not disabled and not seen
	appendIfAllowed := func(a Ability) {
		if _, isDisabled := disabled[a.ID]; isDisabled {
			return
		}
		if _, exists := seen[a.ID]; exists {
			return
		}
		out = append(out, a)
		seen[a.ID] = struct{}{}
	}

	// Base abilities
	for _, a := range s.Abilities {
		appendIfAllowed(a)
	}
	// Runeword-granted abilities
	for _, id := range runewordGrants {
		appendIfAllowed(abilityByID(id))
	}
	// Mode-enabled abilities (ensure included while in this mode)
	for _, id := range spec.EnabledAbilities {
		appendIfAllowed(abilityByID(id))
	}
	return out
}

// Internal: fetch ability from catalog with a safe fallback for missing data.
func abilityByID(id AbilityID) Ability {
	if a, ok := AbilitiesCatalog[id]; ok {
		return a
	}
	return Ability{ID: id, Name: string(id), Kind: AbilityPassive, Description: "(missing from catalog)"}
}
