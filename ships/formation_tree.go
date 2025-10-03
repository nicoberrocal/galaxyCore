package ships

import (
	"time"
)

// FormationTreeNode represents a single node in a formation mastery tree.
// Nodes are unlocked by spending experience points and grant permanent bonuses.
type FormationTreeNode struct {
	ID          string         `bson:"id" json:"id"`
	Name        string         `bson:"name" json:"name"`
	Description string         `bson:"description" json:"description"`
	Formation   FormationType  `bson:"formation" json:"formation"` // "" = global tree
	Tier        int            `bson:"tier" json:"tier"`           // 1-4
	Row         int            `bson:"row" json:"row"`             // Visual positioning (0-based)
	Column      int            `bson:"column" json:"column"`       // Visual positioning (0-based)
	Cost        NodeCost       `bson:"cost" json:"cost"`
	Requirements NodeRequirements `bson:"requirements" json:"requirements"`
	Effects     NodeEffects    `bson:"effects" json:"effects"`
	Tags        []string       `bson:"tags,omitempty" json:"tags,omitempty"` // For filtering/searching
}

// NodeCost defines what's required to unlock a node.
type NodeCost struct {
	ExperiencePoints int `bson:"xp" json:"xp"`           // Primary currency
	Credits          int `bson:"credits,omitempty" json:"credits,omitempty"` // Optional resource cost
	RequiredRank     int `bson:"rank,omitempty" json:"rank,omitempty"`       // Fleet admiral rank
}

// NodeRequirements defines prerequisites for unlocking a node.
type NodeRequirements struct {
	RequiredNodes      []string      `bson:"requiredNodes,omitempty" json:"requiredNodes,omitempty"`           // Must have these nodes
	MutuallyExclusive  []string      `bson:"mutuallyExclusive,omitempty" json:"mutuallyExclusive,omitempty"`   // Can't have with these
	MinNodesInTree     int           `bson:"minNodesInTree,omitempty" json:"minNodesInTree,omitempty"`         // Total nodes in this tree
	MinTierCompleted   int           `bson:"minTierCompleted,omitempty" json:"minTierCompleted,omitempty"`     // Must complete tier X
	RequiredFormation  FormationType `bson:"requiredFormation,omitempty" json:"requiredFormation,omitempty"`   // For cross-tree synergies
}

// NodeEffects defines what a node does when unlocked.
type NodeEffects struct {
	// Direct stat modifications
	PositionMods  map[FormationPosition]StatMods `bson:"positionMods,omitempty" json:"positionMods,omitempty"`     // Apply to specific positions
	FormationMods StatMods                       `bson:"formationMods,omitempty" json:"formationMods,omitempty"`   // Apply to whole formation
	GlobalMods    StatMods                       `bson:"globalMods,omitempty" json:"globalMods,omitempty"`         // Apply to all formations
	
	// Unlocks and grants
	UnlocksAbility   AbilityID `bson:"unlocksAbility,omitempty" json:"unlocksAbility,omitempty"`     // Grant new ability
	UnlocksFormation FormationType `bson:"unlocksFormation,omitempty" json:"unlocksFormation,omitempty"` // Unlock formation type
	
	// Meta modifiers
	ReconfigTimeMultiplier    float64                   `bson:"reconfigMult,omitempty" json:"reconfigMult,omitempty"`         // Formation switch speed
	CounterBonusMultiplier    float64                   `bson:"counterMult,omitempty" json:"counterMult,omitempty"`           // Enhance formation counters
	CounterResistMultiplier   float64                   `bson:"counterResist,omitempty" json:"counterResist,omitempty"`       // Reduce counter damage taken
	CompositionBonusMultiplier float64                  `bson:"compositionMult,omitempty" json:"compositionMult,omitempty"`   // Enhance composition bonuses
	
	// Special mechanics (custom implementations)
	CustomEffect string `bson:"customEffect,omitempty" json:"customEffect,omitempty"` // For unique node abilities
	CustomParams map[string]interface{} `bson:"customParams,omitempty" json:"customParams,omitempty"` // Parameters for custom effects
}

// FormationTreeState tracks a player's progress in the formation mastery system.
type FormationTreeState struct {
	PlayerID        string                    `bson:"playerId" json:"playerId"`
	
	// Experience tracking
	TotalXP         int                       `bson:"totalXP" json:"totalXP"`         // Lifetime earned
	SpentXP         int                       `bson:"spentXP" json:"spentXP"`         // Used for nodes
	AvailableXP     int                       `bson:"availableXP" json:"availableXP"` // Ready to spend
	
	// Formation-specific XP (optional tracking for UI)
	FormationXP     map[FormationType]int     `bson:"formationXP,omitempty" json:"formationXP,omitempty"`
	
	// Active nodes
	UnlockedNodes   []string                  `bson:"unlockedNodes" json:"unlockedNodes"`   // Node IDs
	
	// Reset tracking
	LastResetAt     time.Time                 `bson:"lastResetAt,omitempty" json:"lastResetAt,omitempty"`
	FreeResetsLeft  int                       `bson:"freeResetsLeft" json:"freeResetsLeft"`     // Regenerate over time
	TotalResets     int                       `bson:"totalResets" json:"totalResets"`           // Lifetime count
	NextFreeResetAt time.Time                 `bson:"nextFreeResetAt,omitempty" json:"nextFreeResetAt,omitempty"`
	
	// Meta info
	CreatedAt       time.Time                 `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time                 `bson:"updatedAt" json:"updatedAt"`
	Version         int                       `bson:"version" json:"version"` // For optimistic locking
}

// FormationTree represents the complete tree structure for a formation type (or global).
type FormationTree struct {
	Formation   FormationType       `bson:"formation" json:"formation"` // "" = global
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	Nodes       []FormationTreeNode `bson:"nodes" json:"nodes"`
	MaxTier     int                 `bson:"maxTier" json:"maxTier"`
}

// ExperienceGain represents XP awarded from various sources.
type ExperienceGain struct {
	Source       string    `bson:"source" json:"source"`             // "combat", "quest", "achievement"
	Amount       int       `bson:"amount" json:"amount"`
	FormationType FormationType `bson:"formation,omitempty" json:"formation,omitempty"` // If gained while using specific formation
	Timestamp    time.Time `bson:"timestamp" json:"timestamp"`
	Description  string    `bson:"description,omitempty" json:"description,omitempty"`
}

// NodeUnlockResult contains the result of attempting to unlock a node.
type NodeUnlockResult struct {
	Success       bool     `json:"success"`
	Node          *FormationTreeNode `json:"node,omitempty"`
	ErrorMessage  string   `json:"errorMessage,omitempty"`
	XPRemaining   int      `json:"xpRemaining"`
	EffectsApplied bool    `json:"effectsApplied"`
}

// TreeResetCost calculates the cost to reset the tree.
type TreeResetCost struct {
	Credits      int    `json:"credits"`
	IsFree       bool   `json:"isFree"`
	NextFreeIn   int64  `json:"nextFreeInSeconds,omitempty"` // Seconds until next free reset
	ResetCount   int    `json:"resetCount"`                  // Current total resets
}

// ===================
// Core Methods
// ===================

// HasNode checks if a node is unlocked.
func (ts *FormationTreeState) HasNode(nodeID string) bool {
	for _, id := range ts.UnlockedNodes {
		if id == nodeID {
			return true
		}
	}
	return false
}

// CanUnlockNode checks if a node can be unlocked (has requirements and XP).
func (ts *FormationTreeState) CanUnlockNode(node *FormationTreeNode) (bool, string) {
	// Already unlocked?
	if ts.HasNode(node.ID) {
		return false, "Node already unlocked"
	}
	
	// Sufficient XP?
	if ts.AvailableXP < node.Cost.ExperiencePoints {
		return false, "Insufficient experience points"
	}
	
	// Check required nodes
	for _, reqID := range node.Requirements.RequiredNodes {
		if !ts.HasNode(reqID) {
			return false, "Missing required node: " + reqID
		}
	}
	
	// Check mutually exclusive nodes
	for _, exclusiveID := range node.Requirements.MutuallyExclusive {
		if ts.HasNode(exclusiveID) {
			return false, "Cannot unlock: mutually exclusive with " + exclusiveID
		}
	}
	
	// Check minimum nodes in tree
	if node.Requirements.MinNodesInTree > 0 {
		count := ts.CountNodesInTree(node.Formation)
		if count < node.Requirements.MinNodesInTree {
			return false, "Requires more nodes in tree"
		}
	}
	
	return true, ""
}

// UnlockNode attempts to unlock a node and deduct XP.
func (ts *FormationTreeState) UnlockNode(node *FormationTreeNode, now time.Time) NodeUnlockResult {
	canUnlock, errMsg := ts.CanUnlockNode(node)
	
	if !canUnlock {
		return NodeUnlockResult{
			Success:      false,
			ErrorMessage: errMsg,
			XPRemaining:  ts.AvailableXP,
		}
	}
	
	// Deduct XP
	ts.AvailableXP -= node.Cost.ExperiencePoints
	ts.SpentXP += node.Cost.ExperiencePoints
	
	// Add node
	ts.UnlockedNodes = append(ts.UnlockedNodes, node.ID)
	ts.UpdatedAt = now
	ts.Version++
	
	return NodeUnlockResult{
		Success:        true,
		Node:           node,
		XPRemaining:    ts.AvailableXP,
		EffectsApplied: true,
	}
}

// AwardExperience adds XP to the state.
func (ts *FormationTreeState) AwardExperience(gain ExperienceGain) {
	ts.TotalXP += gain.Amount
	ts.AvailableXP += gain.Amount
	
	// Track per-formation XP if specified
	if gain.FormationType != "" {
		if ts.FormationXP == nil {
			ts.FormationXP = make(map[FormationType]int)
		}
		ts.FormationXP[gain.FormationType] += gain.Amount
	}
	
	ts.UpdatedAt = gain.Timestamp
	ts.Version++
}

// CountNodesInTree counts how many nodes are unlocked in a specific tree.
func (ts *FormationTreeState) CountNodesInTree(formation FormationType) int {
	count := 0
	
	// Get the tree
	tree := GetFormationTree(formation)
	if tree == nil {
		return 0
	}
	
	// Count matches
	for _, nodeID := range ts.UnlockedNodes {
		for _, node := range tree.Nodes {
			if node.ID == nodeID {
				count++
				break
			}
		}
	}
	
	return count
}

// GetResetCost calculates the cost to reset the tree.
func (ts *FormationTreeState) GetResetCost(now time.Time) TreeResetCost {
	// Check if free reset available
	if ts.FreeResetsLeft > 0 {
		return TreeResetCost{
			Credits:    0,
			IsFree:     true,
			ResetCount: ts.TotalResets,
		}
	}
	
	// Calculate escalating cost
	// Formula: 1000 * (2 ^ resets)
	baseCost := 1000
	cost := baseCost
	for i := 0; i < ts.TotalResets; i++ {
		cost *= 2
		if cost > 1000000 { // Cap at 1 million
			cost = 1000000
			break
		}
	}
	
	// Calculate time until next free reset
	var nextFreeIn int64
	if !ts.NextFreeResetAt.IsZero() && now.Before(ts.NextFreeResetAt) {
		nextFreeIn = int64(ts.NextFreeResetAt.Sub(now).Seconds())
	}
	
	return TreeResetCost{
		Credits:    cost,
		IsFree:     false,
		NextFreeIn: nextFreeIn,
		ResetCount: ts.TotalResets,
	}
}

// ResetTree resets all unlocked nodes and refunds XP.
func (ts *FormationTreeState) ResetTree(now time.Time, useFreeReset bool) {
	// Refund XP
	ts.AvailableXP += ts.SpentXP
	ts.SpentXP = 0
	
	// Clear nodes
	ts.UnlockedNodes = []string{}
	
	// Update reset tracking
	ts.TotalResets++
	ts.LastResetAt = now
	
	if useFreeReset && ts.FreeResetsLeft > 0 {
		ts.FreeResetsLeft--
	}
	
	ts.UpdatedAt = now
	ts.Version++
}

// GrantFreeReset adds a free reset (e.g., monthly grant).
func (ts *FormationTreeState) GrantFreeReset(now time.Time, maxFreeResets int) {
	if ts.FreeResetsLeft < maxFreeResets {
		ts.FreeResetsLeft++
		// Set next free reset 30 days from now
		ts.NextFreeResetAt = now.AddDate(0, 1, 0)
		ts.UpdatedAt = now
	}
}

// GetUnlockedNodesInTree returns all unlocked nodes for a specific tree.
func (ts *FormationTreeState) GetUnlockedNodesInTree(formation FormationType) []FormationTreeNode {
	tree := GetFormationTree(formation)
	if tree == nil {
		return []FormationTreeNode{}
	}
	
	unlocked := []FormationTreeNode{}
	for _, nodeID := range ts.UnlockedNodes {
		for _, node := range tree.Nodes {
			if node.ID == nodeID {
				unlocked = append(unlocked, node)
				break
			}
		}
	}
	
	return unlocked
}

// ===================
// Helper Functions
// ===================

// NewFormationTreeState creates a new state for a player.
func NewFormationTreeState(playerID string, now time.Time) *FormationTreeState {
	return &FormationTreeState{
		PlayerID:       playerID,
		TotalXP:        0,
		SpentXP:        0,
		AvailableXP:    0,
		UnlockedNodes:  []string{},
		FreeResetsLeft: 1, // Start with 1 free reset
		TotalResets:    0,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}
}

// GetFormationTree retrieves a tree by formation type (or global if empty).
// This is a placeholder - actual implementation will use FormationTreeCatalog.
func GetFormationTree(formation FormationType) *FormationTree {
	tree, ok := FormationTreeCatalog[formation]
	if ok {
		return &tree
	}
	// Try global tree
	tree, ok = FormationTreeCatalog[""]
	if ok {
		return &tree
	}
	return nil
}
