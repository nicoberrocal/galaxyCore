package ships

import (
	"time"
)

// AddFormationTreeNodes adds all unlocked formation tree node bonuses to the modifier stack.
// This integrates the formation mastery system with the V2 compute architecture.
func (mb *ModifierBuilder) AddFormationTreeNodes(treeState *FormationTreeState, formation FormationType) *ModifierBuilder {
	if treeState == nil {
		return mb
	}
	
	// Get unlocked nodes for this formation
	unlockedNodes := treeState.GetUnlockedNodesInTree(formation)
	
	// Also get global tree nodes
	globalNodes := treeState.GetUnlockedNodesInTree("")
	
	// Combine them
	allNodes := append(unlockedNodes, globalNodes...)
	
	// Apply each node's effects
	for _, node := range allNodes {
		mb.applyNodeEffects(node)
	}
	
	return mb
}

// applyNodeEffects applies a single node's effects to the modifier builder.
func (mb *ModifierBuilder) applyNodeEffects(node FormationTreeNode) {
	// Apply position-specific modifiers
	for position, mods := range node.Effects.PositionMods {
		if !isZeroMods(mods) {
			requiresFormation := true
			mb.stack.AddConditional(
				ModifierSource("formation_tree"),
				node.ID,
				node.Name+" ("+string(position)+")",
				mods,
				PriorityFormation+50, // Slightly higher than base formation
				mb.now,
				nil,
				&requiresFormation,
			)
		}
	}
	
	// Apply formation-wide modifiers
	if !isZeroMods(node.Effects.FormationMods) {
		requiresFormation := true
		mb.stack.AddConditional(
			ModifierSource("formation_tree"),
			node.ID,
			node.Name+" (Formation)",
			node.Effects.FormationMods,
			PriorityFormation+50,
			mb.now,
			nil,
			&requiresFormation,
		)
	}
	
	// Apply global modifiers (always active)
	if !isZeroMods(node.Effects.GlobalMods) {
		mb.stack.AddPermanent(
			ModifierSource("formation_tree"),
			node.ID,
			node.Name+" (Global)",
			node.Effects.GlobalMods,
			PriorityFormation+50,
			mb.now,
		)
	}
	
	// Meta modifiers (these are handled by custom effect system)
	// Custom effects need to be handled by combat/game logic layer
}

// ComputeLoadoutV2WithTree extends ComputeLoadoutV2 to include formation tree bonuses.
func ComputeLoadoutV2WithTree(
	ship Ship,
	role RoleMode,
	loadout ShipLoadout,
	formation *Formation,
	position FormationPosition,
	ships map[ShipType][]HPBucket,
	treeState *FormationTreeState,
	now time.Time,
	inCombat bool,
) (*ModifierStack, StatMods, []AbilityID) {
	builder := NewModifierBuilder(now)
	
	// 1. Gems: provide their own StatMods
	builder.AddGemsFromLoadout(loadout)
	
	// 2. Role Mode: provides its own StatMods
	builder.AddRoleMode(role)
	
	// 3. Formation Tree: provides StatMods from unlocked nodes
	if formation != nil && treeState != nil {
		builder.AddFormationTreeNodes(treeState, formation.Type)
	}
	
	// 4. Formation: provides StatMods from FormationCatalog position bonuses
	if formation != nil {
		builder.AddFormationPosition(formation, position)
	}
	
	// 5. Anchored state: provides penalty mods
	builder.AddAnchoredPenalty(loadout.Anchored)
	
	// Build and resolve
	stack := builder.Build()
	
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
	
	// Add abilities from tree nodes
	if formation != nil && treeState != nil {
		treeGrants := GetTreeGrantedAbilities(treeState, formation.Type)
		grants = append(grants, treeGrants...)
	}
	
	return stack, finalMods, grants
}

// GetTreeGrantedAbilities returns all abilities granted by unlocked tree nodes.
func GetTreeGrantedAbilities(treeState *FormationTreeState, formation FormationType) []AbilityID {
	if treeState == nil {
		return []AbilityID{}
	}
	
	abilities := []AbilityID{}
	
	// Check all unlocked nodes
	unlockedNodes := treeState.GetUnlockedNodesInTree(formation)
	globalNodes := treeState.GetUnlockedNodesInTree("")
	allNodes := append(unlockedNodes, globalNodes...)
	
	for _, node := range allNodes {
		if node.Effects.UnlocksAbility != "" {
			abilities = append(abilities, node.Effects.UnlocksAbility)
		}
	}
	
	return abilities
}

// ApplyFormationTreeModifiers applies tree bonuses to an existing modifier stack (helper).
func ApplyFormationTreeModifiers(
	stack *ModifierStack,
	treeState *FormationTreeState,
	formation FormationType,
	now time.Time,
) {
	if treeState == nil || formation == "" {
		return
	}
	
	builder := &ModifierBuilder{
		stack: stack,
		now:   now,
	}
	
	builder.AddFormationTreeNodes(treeState, formation)
}

// CalculateEffectiveReconfigTime calculates formation reconfiguration time with tree bonuses.
func CalculateEffectiveReconfigTime(
	baseReconfigTime int,
	treeState *FormationTreeState,
	formation FormationType,
) int {
	if treeState == nil {
		return baseReconfigTime
	}
	
	multiplier := 1.0
	
	// Get all unlocked nodes
	unlockedNodes := treeState.GetUnlockedNodesInTree(formation)
	globalNodes := treeState.GetUnlockedNodesInTree("")
	allNodes := append(unlockedNodes, globalNodes...)
	
	// Apply all reconfiguration multipliers
	for _, node := range allNodes {
		if node.Effects.ReconfigTimeMultiplier != 0 {
			multiplier += node.Effects.ReconfigTimeMultiplier
		}
	}
	
	// Ensure minimum is at least 5 seconds
	result := int(float64(baseReconfigTime) * multiplier)
	if result < 5 {
		result = 5
	}
	
	return result
}

// CalculateEffectiveCounterMultiplier calculates formation counter with tree bonuses.
func CalculateEffectiveCounterMultiplier(
	attackerFormation, defenderFormation FormationType,
	attackerTreeState *FormationTreeState,
) float64 {
	// Get base counter
	baseCounter := GetFormationCounterMultiplier(attackerFormation, defenderFormation)
	
	if attackerTreeState == nil {
		return baseCounter
	}
	
	bonusMultiplier := 0.0
	resistMultiplier := 0.0
	
	// Get unlocked nodes
	unlockedNodes := attackerTreeState.GetUnlockedNodesInTree(attackerFormation)
	
	for _, node := range unlockedNodes {
		bonusMultiplier += node.Effects.CounterBonusMultiplier
		resistMultiplier += node.Effects.CounterResistMultiplier
	}
	
	// Apply bonus to counter advantage
	if baseCounter > 1.0 {
		// We have an advantage, enhance it
		advantage := baseCounter - 1.0
		enhancedAdvantage := advantage * (1.0 + bonusMultiplier)
		return 1.0 + enhancedAdvantage
	} else if baseCounter < 1.0 {
		// We have a disadvantage, reduce it with resist
		disadvantage := 1.0 - baseCounter
		reducedDisadvantage := disadvantage * (1.0 - resistMultiplier)
		return 1.0 - reducedDisadvantage
	}
	
	return baseCounter
}

// GetTreeCustomEffects returns all active custom effects from unlocked nodes.
// This is for the game logic layer to handle special mechanics.
func GetTreeCustomEffects(treeState *FormationTreeState, formation FormationType) map[string]map[string]interface{} {
	if treeState == nil {
		return map[string]map[string]interface{}{}
	}
	
	effects := make(map[string]map[string]interface{})
	
	unlockedNodes := treeState.GetUnlockedNodesInTree(formation)
	globalNodes := treeState.GetUnlockedNodesInTree("")
	allNodes := append(unlockedNodes, globalNodes...)
	
	for _, node := range allNodes {
		if node.Effects.CustomEffect != "" {
			effects[node.Effects.CustomEffect] = node.Effects.CustomParams
		}
	}
	
	return effects
}

// HasTreeCustomEffect checks if a specific custom effect is active.
func HasTreeCustomEffect(treeState *FormationTreeState, formation FormationType, effectName string) bool {
	effects := GetTreeCustomEffects(treeState, formation)
	_, exists := effects[effectName]
	return exists
}

// GetTreeCustomEffectParams gets the parameters for a custom effect.
func GetTreeCustomEffectParams(treeState *FormationTreeState, formation FormationType, effectName string) map[string]interface{} {
	effects := GetTreeCustomEffects(treeState, formation)
	if params, exists := effects[effectName]; exists {
		return params
	}
	return nil
}
