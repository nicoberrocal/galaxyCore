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

// AddFormationRoleSynergy is DEPRECATED - removed for clean system separation.
// Formation bonuses come from FormationCatalog + tree nodes only.
// Role bonuses come from RoleMode only.
func (mb *ModifierBuilder) AddFormationRoleSynergy(formation *Formation, position FormationPosition, role RoleMode) *ModifierBuilder {
	// This method is deprecated and does nothing.
	// Kept for backward compatibility during transition.
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

// AddCompositionBonus is DEPRECATED - removed for clean system separation.
// Fleet composition bonuses create implicit synergies between ship types.
// Each ship should contribute independently.
func (mb *ModifierBuilder) AddCompositionBonus(ships map[ShipType][]HPBucket) *ModifierBuilder {
	// This method is deprecated and does nothing.
	// Kept for backward compatibility during transition.
	return mb
}

// AddGemPositionSynergy is DEPRECATED - removed for clean system separation.
// Gems provide their own StatMods only.
// Formation provides its own StatMods only.
func (mb *ModifierBuilder) AddGemPositionSynergy(gems []Gem, position FormationPosition) *ModifierBuilder {
	// This method is deprecated and does nothing.
	// Kept for backward compatibility during transition.
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
// Biology Runtime Integration
// ========================

// AddBioFromMachine adds active bio node layers for the given ship type at builder time.
func (mb *ModifierBuilder) AddBioFromMachine(bio *BioMachine, shipType ShipType) *ModifierBuilder {
	if bio == nil {
		return mb
	}
	layers := bio.CollectActiveLayersForShip(shipType, mb.now)
	for _, bl := range layers {
		if !isZeroMods(bl.Mods) {
			if bl.ExpiresAt != nil && mb.now.Before(*bl.ExpiresAt) {
				mb.stack.AddTemporary(bl.Source, bl.SourceID, bl.Desc, bl.Mods, bl.Priority, mb.now, bl.ExpiresAt.Sub(mb.now))
			} else {
				// treat as permanent for this snapshot
				mb.stack.AddPermanent(bl.Source, bl.SourceID, bl.Desc, bl.Mods, bl.Priority, mb.now)
			}
		}
	}
	return mb
}

// AddInboundBioDebuffs adds enemy-applied debuffs captured by the bio machine.
func (mb *ModifierBuilder) AddInboundBioDebuffs(bio *BioMachine) *ModifierBuilder {
	if bio == nil {
		return mb
	}
	for _, d := range bio.CollectInboundDebuffs(mb.now) {
		mods := scaleMods(d.Mods, float64(max(1, d.Stacks)))
		dur := d.ExpiresAt.Sub(mb.now)
		if dur <= 0 {
			// if already expired by clock skew, skip
			continue
		}
		mb.stack.AddTemporary(SourceBioDebuff, d.ID, "Bio Debuff: "+d.ID, mods, PriorityBioDebuff, mb.now, dur)
	}
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
    return m.IsZero()
}
