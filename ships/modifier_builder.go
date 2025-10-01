package ships

import (
	"fmt"
	"time"
)

// ModifierBuilder provides a fluent API for constructing modifier stacks from game entities.
// This bridges the gap between your existing systems (gems, formations, abilities, etc.)
// and the new layered modifier system.
type ModifierBuilder struct {
	stack *ModifierStack
	now   time.Time
}

// NewModifierBuilder creates a new builder with the current time.
func NewModifierBuilder(now time.Time) *ModifierBuilder {
	return &ModifierBuilder{
		stack: NewModifierStack(),
		now:   now,
	}
}

// Build returns the constructed modifier stack.
func (mb *ModifierBuilder) Build() *ModifierStack {
	return mb.stack
}

// ========================
// Gem System Integration
// ========================

// AddGems adds modifiers from socketed gems.
func (mb *ModifierBuilder) AddGems(gems []Gem) *ModifierBuilder {
	for i, gem := range gems {
		mb.stack.AddPermanent(
			SourceGem,
			string(gem.ID),
			fmt.Sprintf("%s (Socket %d)", gem.Name, i+1),
			gem.Mods,
			PriorityGem,
			mb.now,
		)
	}
	return mb
}

// AddGemWords adds modifiers from matched GemWords.
func (mb *ModifierBuilder) AddGemWords(gemWords []GemWord) *ModifierBuilder {
	for _, gw := range gemWords {
		mb.stack.AddPermanent(
			SourceGemWord,
			gw.Name,
			fmt.Sprintf("GemWord: %s", gw.Name),
			gw.Effects,
			PriorityGemWord,
			mb.now,
		)
	}
	return mb
}

// AddGemsFromLoadout is a convenience method that evaluates gems and gemwords together.
func (mb *ModifierBuilder) AddGemsFromLoadout(loadout ShipLoadout) *ModifierBuilder {
	socketMods, _, gemWords := EvaluateGemSockets(loadout.Sockets)
	
	// Add individual gem contributions
	mb.AddGems(loadout.Sockets)
	
	// Add GemWord bonuses
	mb.AddGemWords(gemWords)
	
	// Note: socketMods already includes both, but we're adding them separately
	// for transparency. If you prefer a single combined layer, use AddGemsFromLoadoutCombined.
	_ = socketMods // Already accounted for in AddGems + AddGemWords
	
	return mb
}

// AddGemsFromLoadoutCombined adds gems as a single combined layer (simpler, less transparent).
func (mb *ModifierBuilder) AddGemsFromLoadoutCombined(loadout ShipLoadout) *ModifierBuilder {
	socketMods, _, _ := EvaluateGemSockets(loadout.Sockets)
	
	mb.stack.AddPermanent(
		SourceGem,
		"all_gems",
		"All Socketed Gems",
		socketMods,
		PriorityGem,
		mb.now,
	)
	
	return mb
}

// ========================
// Role Mode Integration
// ========================

// AddRoleMode adds modifiers from the active role mode.
func (mb *ModifierBuilder) AddRoleMode(role RoleMode) *ModifierBuilder {
	spec, ok := RoleModesCatalog[role]
	if !ok {
		return mb
	}
	
	mb.stack.AddPermanent(
		SourceRoleMode,
		string(role),
		fmt.Sprintf("Role: %s", spec.Name),
		spec.BaseMods,
		PriorityRoleMode,
		mb.now,
	)
	
	return mb
}

// ========================
// Formation Integration
// ========================

// AddFormationPosition adds modifiers from formation position bonuses.
func (mb *ModifierBuilder) AddFormationPosition(formation *Formation, position FormationPosition) *ModifierBuilder {
	if formation == nil {
		return mb
	}
	
	spec, ok := FormationCatalog[formation.Type]
	if !ok {
		return mb
	}
	
	if posBonus, ok := spec.PositionBonuses[position]; ok {
		requiresFormation := true
		mb.stack.AddConditional(
			SourceFormationPosition,
			fmt.Sprintf("%s_%s", formation.Type, position),
			fmt.Sprintf("%s Formation - %s Position", spec.Name, position),
			posBonus,
			PriorityFormation,
			mb.now,
			nil,
			&requiresFormation,
		)
	}
	
	return mb
}

// AddFormationRoleSynergy adds modifiers from role+formation synergy.
func (mb *ModifierBuilder) AddFormationRoleSynergy(formation *Formation, position FormationPosition, role RoleMode) *ModifierBuilder {
	if formation == nil {
		return mb
	}
	
	// Calculate the synergy bonus
	baseMods := ZeroMods()
	synergyMods := ApplyFormationRoleModifiers(baseMods, formation, position, role)
	
	// Only add if there's an actual bonus
	if !isZeroMods(synergyMods) {
		requiresFormation := true
		mb.stack.AddConditional(
			SourceFormationRole,
			fmt.Sprintf("%s_%s_%s", formation.Type, position, role),
			fmt.Sprintf("%s + %s Synergy", formation.Type, role),
			synergyMods,
			PrioritySynergy,
			mb.now,
			nil,
			&requiresFormation,
		)
	}
	
	return mb
}

// AddFormationCounter adds modifiers from formation matchup (attacker vs defender).
func (mb *ModifierBuilder) AddFormationCounter(attackerFormation, defenderFormation FormationType, inCombat bool) *ModifierBuilder {
	multiplier := GetFormationCounterMultiplier(attackerFormation, defenderFormation)
	
	// Convert multiplier to damage bonus/penalty
	if multiplier != 1.0 {
		damageBonus := multiplier - 1.0
		mods := StatMods{
			Damage: DamageMods{
				LaserPct:      damageBonus,
				NuclearPct:    damageBonus,
				AntimatterPct: damageBonus,
			},
		}
		
		combatOnly := true
		requiresFormation := true
		mb.stack.AddConditional(
			SourceFormationCounter,
			fmt.Sprintf("%s_vs_%s", attackerFormation, defenderFormation),
			fmt.Sprintf("%s vs %s (%.0f%%)", attackerFormation, defenderFormation, damageBonus*100),
			mods,
			PrioritySynergy,
			mb.now,
			&combatOnly,
			&requiresFormation,
		)
	}
	
	return mb
}

// AddCompositionBonus adds modifiers from fleet composition bonuses.
func (mb *ModifierBuilder) AddCompositionBonus(ships map[ShipType][]HPBucket) *ModifierBuilder {
	mods, bonuses := EvaluateCompositionBonuses(ships)
	
	for _, bonus := range bonuses {
		mb.stack.AddPermanent(
			SourceComposition,
			bonus.Type,
			fmt.Sprintf("Composition: %s", bonus.Type),
			bonus.Bonus,
			PriorityComposition,
			mb.now,
		)
	}
	
	_ = mods // Already accounted for in individual bonuses
	return mb
}

// AddGemPositionSynergy adds modifiers from gem+position synergies.
func (mb *ModifierBuilder) AddGemPositionSynergy(gems []Gem, position FormationPosition) *ModifierBuilder {
	synergyMods := ApplyGemPositionEffects(gems, position)
	
	if !isZeroMods(synergyMods) {
		requiresFormation := true
		mb.stack.AddConditional(
			SourceGemPosition,
			fmt.Sprintf("gems_%s", position),
			fmt.Sprintf("Gem-Position Synergy (%s)", position),
			synergyMods,
			PrioritySynergy,
			mb.now,
			nil,
			&requiresFormation,
		)
	}
	
	return mb
}

// ========================
// Ability Integration
// ========================

// AddAbility adds temporary modifiers from an active ability.
func (mb *ModifierBuilder) AddAbility(abilityID AbilityID, mods StatMods, duration time.Duration) *ModifierBuilder {
	mb.stack.AddTemporary(
		SourceAbility,
		string(abilityID),
		fmt.Sprintf("Ability: %s", abilityID),
		mods,
		PriorityAbility,
		mb.now,
		duration,
	)
	return mb
}

// AddBuff adds a temporary buff modifier.
func (mb *ModifierBuilder) AddBuff(buffID string, description string, mods StatMods, duration time.Duration) *ModifierBuilder {
	mb.stack.AddTemporary(
		SourceBuff,
		buffID,
		description,
		mods,
		PriorityBuff,
		mb.now,
		duration,
	)
	return mb
}

// AddDebuff adds a temporary debuff modifier.
func (mb *ModifierBuilder) AddDebuff(debuffID string, description string, mods StatMods, duration time.Duration) *ModifierBuilder {
	mb.stack.AddTemporary(
		SourceDebuff,
		debuffID,
		description,
		mods,
		PriorityDebuff,
		mb.now,
		duration,
	)
	return mb
}

// ========================
// Environmental Integration
// ========================

// AddEnvironment adds modifiers from environmental effects (nebula, asteroid field, etc.).
func (mb *ModifierBuilder) AddEnvironment(envID string, description string, mods StatMods) *ModifierBuilder {
	mb.stack.AddPermanent(
		SourceEnvironment,
		envID,
		description,
		mods,
		PriorityEnvironment,
		mb.now,
	)
	return mb
}

// AddAnchoredPenalty adds the anchored state modifier.
func (mb *ModifierBuilder) AddAnchoredPenalty(anchored bool) *ModifierBuilder {
	if !anchored {
		return mb
	}
	
	// Anchored ships have reduced combat effectiveness
	mods := StatMods{
		Damage: DamageMods{
			LaserPct:      -0.15,
			NuclearPct:    -0.15,
			AntimatterPct: -0.15,
		},
		LaserShieldDelta:      -1,
		NuclearShieldDelta:    -1,
		AntimatterShieldDelta: -1,
	}
	
	mb.stack.AddPermanent(
		SourceAnchored,
		"anchored",
		"Anchored (Mining)",
		mods,
		PriorityEnvironment,
		mb.now,
	)
	
	return mb
}

// ========================
// Utility Functions
// ========================

// isZeroMods checks if a StatMods struct has any non-zero values.
func isZeroMods(m StatMods) bool {
	if m.Damage.LaserPct != 0 || m.Damage.NuclearPct != 0 || m.Damage.AntimatterPct != 0 {
		return false
	}
	if m.AttackIntervalPct != 0 || m.SpeedDelta != 0 || m.VisibilityDelta != 0 || m.AttackRangeDelta != 0 {
		return false
	}
	if m.LaserShieldDelta != 0 || m.NuclearShieldDelta != 0 || m.AntimatterShieldDelta != 0 {
		return false
	}
	if m.BucketHPPct != 0 || m.OutOfCombatRegenPct != 0 || m.AbilityCooldownPct != 0 {
		return false
	}
	if m.TransportCapacityPct != 0 || m.WarpChargePct != 0 || m.WarpScatterPct != 0 || m.InterdictionResistPct != 0 {
		return false
	}
	if m.StructureDamagePct != 0 || m.SplashRadiusDelta != 0 || m.AccuracyPct != 0 || m.CritPct != 0 {
		return false
	}
	if m.FirstVolleyPct != 0 || m.ShieldPiercePct != 0 || m.UpkeepPct != 0 || m.ConstructionCostPct != 0 {
		return false
	}
	if m.CloakDetect || m.PingRangePct != 0 || m.EvasionPct != 0 || m.FormationSyncBonus != 0 || m.PositionFlexibility != 0 {
		return false
	}
	return true
}
