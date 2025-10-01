package ships

import (
	"time"
)

// ModifierSource identifies where a modifier comes from and its lifetime characteristics.
type ModifierSource string

const (
	// Permanent sources (while equipped/active)
	SourceGem       ModifierSource = "gem"        // From socketed gems
	SourceGemWord   ModifierSource = "gemword"    // From GemWord patterns
	SourceRoleMode  ModifierSource = "rolemode"   // From active role mode
	
	// Formation sources (while formation active)
	SourceFormationPosition ModifierSource = "formation_position" // Fixed position bonuses
	SourceFormationRole     ModifierSource = "formation_role"     // Role+formation synergy
	SourceFormationCounter  ModifierSource = "formation_counter"  // Cross-formation matchup
	SourceComposition       ModifierSource = "composition"        // Fleet composition bonuses
	SourceGemPosition       ModifierSource = "gem_position"       // Gem+position synergy
	
	// Temporary sources (duration-based)
	SourceAbility      ModifierSource = "ability"       // Active ability effects
	SourceAbilityStack ModifierSource = "ability_stack" // Stacking ability effects
	SourceDebuff       ModifierSource = "debuff"        // Enemy-applied debuffs
	SourceBuff         ModifierSource = "buff"          // Ally-applied buffs
	
	// Environmental/situational
	SourceEnvironment ModifierSource = "environment" // Terrain, nebula effects, etc.
	SourceAnchored    ModifierSource = "anchored"    // Anchoring penalties/bonuses
)

// ModifierLayer represents a single layer of modifiers from a specific source.
// Each layer tracks its source, lifetime, and the actual stat modifications.
type ModifierLayer struct {
	Source      ModifierSource `bson:"source" json:"source"`
	SourceID    string         `bson:"sourceId" json:"sourceId"`       // Specific identifier (gem ID, ability ID, etc.)
	Description string         `bson:"description" json:"description"` // Human-readable description
	Mods        StatMods       `bson:"mods" json:"mods"`
	
	// Lifetime tracking
	AppliedAt time.Time  `bson:"appliedAt" json:"appliedAt"`
	ExpiresAt *time.Time `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"` // nil = permanent/conditional
	
	// Priority for resolution order (higher = applied later, can override)
	Priority int `bson:"priority" json:"priority"`
	
	// Conditional flags
	ActiveInCombat    *bool `bson:"activeInCombat,omitempty" json:"activeInCombat,omitempty"`       // nil = always, true = combat only, false = OOC only
	RequiresFormation *bool `bson:"requiresFormation,omitempty" json:"requiresFormation,omitempty"` // nil = always, true = formation required
}

// ModifierStack is a collection of modifier layers that can be resolved into final StatMods.
// This allows inspection, debugging, and fine-grained control over modifier application.
type ModifierStack struct {
	Layers []ModifierLayer `bson:"layers" json:"layers"`
}

// Priority constants for ordering layers
const (
	PriorityBase         = 0   // Base stats (not a modifier, but conceptual baseline)
	PriorityGem          = 100 // Gems apply first
	PriorityGemWord      = 150 // GemWords build on gems
	PriorityRoleMode     = 200 // Role mode is a strategic choice
	PriorityFormation    = 300 // Formation positioning
	PriorityComposition  = 350 // Fleet composition synergies
	PrioritySynergy      = 400 // Cross-system synergies (gem+position, formation+role)
	PriorityEnvironment  = 500 // Environmental effects
	PriorityAbility      = 600 // Active abilities
	PriorityBuff         = 700 // Allied buffs
	PriorityDebuff       = 800 // Enemy debuffs (applied last to see final stats)
)

// NewModifierStack creates an empty modifier stack.
func NewModifierStack() *ModifierStack {
	return &ModifierStack{
		Layers: []ModifierLayer{},
	}
}

// AddLayer adds a new modifier layer to the stack.
func (ms *ModifierStack) AddLayer(layer ModifierLayer) {
	ms.Layers = append(ms.Layers, layer)
}

// AddPermanent adds a permanent modifier layer (gems, role mode, etc.).
func (ms *ModifierStack) AddPermanent(source ModifierSource, sourceID, description string, mods StatMods, priority int, now time.Time) {
	ms.AddLayer(ModifierLayer{
		Source:      source,
		SourceID:    sourceID,
		Description: description,
		Mods:        mods,
		AppliedAt:   now,
		ExpiresAt:   nil,
		Priority:    priority,
	})
}

// AddTemporary adds a temporary modifier layer with expiration (abilities, buffs, debuffs).
func (ms *ModifierStack) AddTemporary(source ModifierSource, sourceID, description string, mods StatMods, priority int, now time.Time, duration time.Duration) {
	expiresAt := now.Add(duration)
	ms.AddLayer(ModifierLayer{
		Source:      source,
		SourceID:    sourceID,
		Description: description,
		Mods:        mods,
		AppliedAt:   now,
		ExpiresAt:   &expiresAt,
		Priority:    priority,
	})
}

// AddConditional adds a conditional modifier (formation, combat-only, etc.).
func (ms *ModifierStack) AddConditional(source ModifierSource, sourceID, description string, mods StatMods, priority int, now time.Time, activeInCombat, requiresFormation *bool) {
	ms.AddLayer(ModifierLayer{
		Source:            source,
		SourceID:          sourceID,
		Description:       description,
		Mods:              mods,
		AppliedAt:         now,
		ExpiresAt:         nil,
		Priority:          priority,
		ActiveInCombat:    activeInCombat,
		RequiresFormation: requiresFormation,
	})
}

// RemoveExpired removes all expired temporary modifiers.
func (ms *ModifierStack) RemoveExpired(now time.Time) {
	active := make([]ModifierLayer, 0, len(ms.Layers))
	for _, layer := range ms.Layers {
		if layer.ExpiresAt == nil || now.Before(*layer.ExpiresAt) {
			active = append(active, layer)
		}
	}
	ms.Layers = active
}

// RemoveBySource removes all layers from a specific source.
func (ms *ModifierStack) RemoveBySource(source ModifierSource) {
	filtered := make([]ModifierLayer, 0, len(ms.Layers))
	for _, layer := range ms.Layers {
		if layer.Source != source {
			filtered = append(filtered, layer)
		}
	}
	ms.Layers = filtered
}

// RemoveBySourceID removes all layers with a specific source ID.
func (ms *ModifierStack) RemoveBySourceID(sourceID string) {
	filtered := make([]ModifierLayer, 0, len(ms.Layers))
	for _, layer := range ms.Layers {
		if layer.SourceID != sourceID {
			filtered = append(filtered, layer)
		}
	}
	ms.Layers = filtered
}

// Clear removes all modifier layers.
func (ms *ModifierStack) Clear() {
	ms.Layers = []ModifierLayer{}
}

// ResolveContext provides context for resolving modifiers.
type ResolveContext struct {
	Now            time.Time
	InCombat       bool
	HasFormation   bool
	FormationType  FormationType
	EnemyFormation FormationType
}

// Resolve computes the final StatMods by combining all applicable layers.
// Layers are applied in priority order, with higher priority layers applied last.
func (ms *ModifierStack) Resolve(ctx ResolveContext) StatMods {
	// Sort layers by priority (stable sort maintains insertion order for same priority)
	sortedLayers := make([]ModifierLayer, len(ms.Layers))
	copy(sortedLayers, ms.Layers)
	
	// Simple insertion sort by priority
	for i := 1; i < len(sortedLayers); i++ {
		key := sortedLayers[i]
		j := i - 1
		for j >= 0 && sortedLayers[j].Priority > key.Priority {
			sortedLayers[j+1] = sortedLayers[j]
			j--
		}
		sortedLayers[j+1] = key
	}
	
	// Combine applicable layers
	result := ZeroMods()
	for _, layer := range sortedLayers {
		if !ms.isLayerApplicable(layer, ctx) {
			continue
		}
		result = CombineMods(result, layer.Mods)
	}
	
	return result
}

// isLayerApplicable checks if a layer should be applied given the context.
func (ms *ModifierStack) isLayerApplicable(layer ModifierLayer, ctx ResolveContext) bool {
	// Check expiration
	if layer.ExpiresAt != nil && !ctx.Now.Before(*layer.ExpiresAt) {
		return false
	}
	
	// Check combat requirement
	if layer.ActiveInCombat != nil {
		if *layer.ActiveInCombat && !ctx.InCombat {
			return false
		}
		if !*layer.ActiveInCombat && ctx.InCombat {
			return false
		}
	}
	
	// Check formation requirement
	if layer.RequiresFormation != nil && *layer.RequiresFormation && !ctx.HasFormation {
		return false
	}
	
	return true
}

// GetLayersBySource returns all layers from a specific source.
func (ms *ModifierStack) GetLayersBySource(source ModifierSource) []ModifierLayer {
	var result []ModifierLayer
	for _, layer := range ms.Layers {
		if layer.Source == source {
			result = append(result, layer)
		}
	}
	return result
}

// GetLayersBySourceID returns all layers with a specific source ID.
func (ms *ModifierStack) GetLayersBySourceID(sourceID string) []ModifierLayer {
	var result []ModifierLayer
	for _, layer := range ms.Layers {
		if layer.SourceID == sourceID {
			result = append(result, layer)
		}
	}
	return result
}

// Summary returns a breakdown of modifiers by source for debugging/UI.
type ModifierSummary struct {
	Source      ModifierSource `json:"source"`
	Description string         `json:"description"`
	Mods        StatMods       `json:"mods"`
	IsActive    bool           `json:"isActive"`
	ExpiresIn   *float64       `json:"expiresIn,omitempty"` // seconds, nil if permanent
}

// GetSummary returns a human-readable summary of all modifier layers.
func (ms *ModifierStack) GetSummary(ctx ResolveContext) []ModifierSummary {
	summaries := make([]ModifierSummary, 0, len(ms.Layers))
	
	for _, layer := range ms.Layers {
		isActive := ms.isLayerApplicable(layer, ctx)
		
		var expiresIn *float64
		if layer.ExpiresAt != nil {
			remaining := layer.ExpiresAt.Sub(ctx.Now).Seconds()
			if remaining > 0 {
				expiresIn = &remaining
			}
		}
		
		summaries = append(summaries, ModifierSummary{
			Source:      layer.Source,
			Description: layer.Description,
			Mods:        layer.Mods,
			IsActive:    isActive,
			ExpiresIn:   expiresIn,
		})
	}
	
	return summaries
}

// Clone creates a deep copy of the modifier stack.
func (ms *ModifierStack) Clone() *ModifierStack {
	clone := &ModifierStack{
		Layers: make([]ModifierLayer, len(ms.Layers)),
	}
	copy(clone.Layers, ms.Layers)
	return clone
}
