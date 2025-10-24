package ships

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShipType string

const (
	Drone     ShipType = "drone"
	Scout     ShipType = "scout"
	Fighter   ShipType = "fighter"
	Bomber    ShipType = "bomber"
	Carrier   ShipType = "carrier"
	Destroyer ShipType = "destroyer"
	Cruiser   ShipType = "cruiser"
	Corvette  ShipType = "corvette"
	Ballista  ShipType = "ballista"
	Ghost     ShipType = "ghost"
	Frigate   ShipType = "frigate"
)

type BioTreePath string

const (
	Cephalopod     BioTreePath = "cephalopod"
	Chondrichthyan BioTreePath = "chondrichthyan"
	Cetacean       BioTreePath = "cetacean"
	Carnivora      BioTreePath = "carnivora"
	Arbor          BioTreePath = "arbor"
	VerdantBloom   BioTreePath = "verdant_bloom"
	Sporeform      BioTreePath = "sporeform"
	Cordyceps      BioTreePath = "cordyceps"
	Mycorrhiza     BioTreePath = "mycorrhiza"
	Apex           BioTreePath = "apex"
	PackHunter     BioTreePath = "pack_hunter"
	Scavengers     BioTreePath = "scavengers"
)

// ShipStack represents a fleet that is NOT currently defending a system
// When a stack colonizes a system, it gets embedded in the system's DefendingFleet
// and this document is deleted (hybrid approach)
// Battles in free space/mining locations are handled by updating this document directly
type ShipStack struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	PlayerID  bson.ObjectID `bson:"playerId"`
	MapID     bson.ObjectID `bson:"mapId"`
	PositionX float64       `bson:"x"`
	PositionY float64       `bson:"y"`

	// Fleet composition
	Ships     map[ShipType][]HPBucket `bson:"ships"`     // HP bucketed ships
	CreatedAt time.Time               `bson:"createdAt"` // tick timestamp

	// Role represents the tactical intent of the entire stack (tactical/economic/recon/scientific)
	// Role             RoleMode  `bson:"role" json:"role"`
	// ReconfigureUntil time.Time `bson:"reconfigureUntil,omitempty" json:"reconfigureUntil,omitempty"`

	// Loadouts track per-ship-type socket configurations for this particular stack.
	// This allows two stacks to field the same ship type with different gem setups.
	Loadouts map[ShipType]ShipLoadout `bson:"loadouts,omitempty" json:"loadouts,omitempty"`

	// Formation defines the tactical positioning of ships within the stack
	Formation              *FormationWithSlots                  `bson:"formation,omitempty" json:"formation,omitempty"`
	FormationReconfigUntil time.Time                            `bson:"formationReconfigUntil,omitempty" json:"formationReconfigUntil,omitempty"`
	SavedFormations        map[FormationType]FormationWithSlots `bson:"savedFormations,omitempty" json:"savedFormations,omitempty"`

	// Current activity and movement
	Movement    []*MovementState `bson:"movement,omitempty" json:"movement,omitempty"`
	Battle      *BattleState     `bson:"battle,omitempty" json:"battle,omitempty"`       // Combat state for free space battles
	Ability     *[]AbilityState  `bson:"ability,omitempty" json:"ability,omitempty"`     // Active ship ability state
	Gathering   *GatheringState  `bson:"gathering,omitempty" json:"gathering,omitempty"` // Active gathering state
	BioTreePath BioTreePath      `bson:"bioTreePath,omitempty" json:"bioTreePath,omitempty"`
	Bio         *BioMachine      `bson:"bio,omitempty" json:"bio,omitempty"` // Biology node runtime state machine
	Version     int64            `bson:"version" json:"version"`             // For optimistic locking
}

// EnsureBio initializes the bio runtime machine if missing and returns it.
func (s *ShipStack) EnsureBio(now time.Time) *BioMachine {
	if s.Bio == nil {
		s.Bio = NewBioMachine(now)
	}
	return s.Bio
}

func (s *ShipStack) BuildBioFromCurrentPath(now time.Time) {
	s.EnsureBio(now)
	if BioPopulateFromPath != nil {
		BioPopulateFromPath(s, now)
	}
}

func (s *ShipStack) BuildBioFromPath(path BioTreePath, now time.Time) {
	s.BioTreePath = path
	s.EnsureBio(now)
	if BioPopulateFromExplicitPath != nil {
		BioPopulateFromExplicitPath(s, path, now)
	} else if BioPopulateFromPath != nil {
		BioPopulateFromPath(s, now)
	}
}

func (s *ShipStack) TickBio(now time.Time) {
	if s.Bio != nil {
		s.Bio.Tick(now)
	}
}

// BioOnAbilityCast proxies an ability-cast event into the bio machine for stage transitions.
func (s *ShipStack) BioOnAbilityCast(ability AbilityID, shipType ShipType, start time.Time) {
	if s.Bio == nil {
		s.Bio = NewBioMachine(start)
	}
	s.Bio.OnAbilityCast(ability, shipType, start)
}

// BioApplyInboundDebuff upserts an enemy-applied bio debuff on this stack.
func (s *ShipStack) BioApplyInboundDebuff(id string, mods StatMods, duration time.Duration, stacks int, maxStacks int, sourceStack bson.ObjectID, now time.Time) {
	if s.Bio == nil {
		s.Bio = NewBioMachine(now)
	}
	s.Bio.ApplyInboundDebuff(id, mods, duration, stacks, maxStacks, sourceStack, now)
}

// StartModeSwitch sets a new RoleMode for the entire stack and applies a
// reconfiguration timer based on RoleModesCatalog. Returns the ETA.
func (s *ShipStack) StartModeSwitch(newRole RoleMode, now time.Time) time.Time {
	spec, ok := RoleModesCatalog[newRole]
	reconfig := 180
	if ok && spec.ReconfigureSeconds > 0 {
		reconfig = spec.ReconfigureSeconds
	}
	reconf := now.Add(time.Duration(reconfig) * time.Second)
	return reconf
}

// SetAnchored updates the anchored state for this ship type on the stack.
func (s *ShipStack) SetAnchored(t ShipType, anchored bool) {
	load := s.GetOrInitLoadout(t)
	load.Anchored = anchored
	if s.Loadouts == nil {
		s.Loadouts = make(map[ShipType]ShipLoadout)
	}
	s.Loadouts[t] = load
}

type HPBucket struct {
	HP    int `bson:"hp"`    // current HP of ships in this bucket
	Count int `bson:"count"` // how many ships at this HP
}

// CombatCounters tracks deterministic combat mechanics (crit intervals, evasion, etc.)
// These counters enable predictable combat outcomes without RNG.
type CombatCounters struct {
	AttackCount    int `bson:"attackCount" json:"attackCount"`       // Total attacks made (for crit timing)
	DefenseCount   int `bson:"defenseCount" json:"defenseCount"`     // Total attacks received (for evasion timing)
	LastCritAttack int `bson:"lastCritAttack" json:"lastCritAttack"` // Attack number of last crit
}

// BattleState tracks combat information for stacks in free space or mining locations
type BattleState struct {
	IsInCombat      bool             `bson:"isInCombat"`                               // Currently engaged in battle
	EnemyStackID    []bson.ObjectID  `bson:"enemyStackId,omitempty"`                   // Opponent stack ID
	EnemyPlayerID   []bson.ObjectID  `bson:"enemyPlayerId,omitempty"`                  // Opponent player ID
	BattleStartedAt time.Time        `bson:"battleStartedAt,omitempty"`                // When battle began
	BattleLocation  string           `bson:"battleLocation,omitempty"`                 // "empty_space", "asteroid", "nebula"
	LocationID      bson.ObjectID    `bson:"locationId,omitempty"`                     // ID of asteroid/nebula if applicable
	ProcessedAt     time.Time        `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
	Counters        *CombatCounters  `bson:"counters,omitempty" json:"counters,omitempty"` // Deterministic combat counters
}

// MovementState tracks what the stack is currently doing in free space or at mining locations
type MovementState struct {
	Type        string        `bson:"type" json:"type,omitempty"`                       // "traveling", "mining", "idle", "in_combat"
	State       string        `bson:"state" json:"state,omitempty"`                     // "traveling", "mining", "idle", "in_combat"
	TargetID    bson.ObjectID `bson:"targetId,omitempty" json:"targetId,omitempty"`     // Target for movement/mining
	TargetType  string        `bson:"targetType,omitempty" json:"targetType,omitempty"` // "asteroid", "nebula", "coordinate"
	StartX      float64       `bson:"startX,omitempty" json:"startX,omitempty"`         // Starting coordinates
	StartY      float64       `bson:"startY,omitempty" json:"startY"`                   // Starting coordinates
	Speed       int           `bson:"speed,omitempty" json:"speed"`                     // Speed of the stack
	TargetX     float64       `bson:"targetX,omitempty" json:"targetX"`                 // Target coordinates
	Stage       string        `bson:"stage,omitempty" json:"stage,omitempty"`           // Stage of the gathering process
	TargetY     float64       `bson:"targetY,omitempty" json:"targetY"`                 // Target coordinates
	StartTime   time.Time     `bson:"startTime,omitempty" json:"startTime"`             // When current action started
	EndTime     time.Time     `bson:"endTime,omitempty" json:"endTime"`                 // When current action ends
	Activity    string        `bson:"activity,omitempty" json:"activity"`               // "mining_metal", "mining_crystal", "mining_hydrogen"
	LastUpdated time.Time     `bson:"lastUpdated,omitempty" json:"lastUpdated"`         // Last time this state was updated
	ProcessedAt time.Time     `bson:"ProcessedAt,omitempty" json:"ProcessedAt"`         // Last time this state was processed
}

type AbilityState struct {
	IsActive    bool           `bson:"isActive" json:"isActive"`                 // Whether the ability is currently active
	Description string         `bson:"description" json:"description"`           // Description of the ability
	Icon        string         `bson:"icon" json:"icon"`                         // Icon representing the ability
	ShipType    ShipType       `bson:"shipType" json:"shipType"`                 // Type of ship using the ability
	Ability     string         `bson:"ability" json:"ability"`                   // Ability name (e.g., "cloak", "boost")
	Bonus       map[string]int `bson:"bonus" json:"bonus"`                       // Bonus type (e.g., "speed", "defense" or "attack") and its value
	StartTime   time.Time      `bson:"startTime" json:"startTime"`               // When the ability was activated
	EndTime     time.Time      `bson:"endTime" json:"endTime"`                   // When the ability will end
	Duration    int64          `bson:"duration" json:"duration"`                 // Duration of the ability in seconds
	LastUpdated time.Time      `bson:"lastUpdated" json:"lastUpdated"`           // Last time the ability state was updated
	ProcessedAt time.Time      `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}

type GatheringState struct {
	IsGathering         bool          `bson:"isGathering" json:"isGathering"`           // Whether the stack is currently mining
	TargetID            bson.ObjectID `bson:"targetId" json:"targetId"`                 // ID of the target being gathered from
	TargetType          string        `bson:"targetType" json:"targetType"`             // Type of target (e.g., "asteroid", "nebula")
	ResourcePerTimeUnit int64         `bson:"resourcePerTime" json:"resourcePerTime"`   // Amount of resource gathered per time unit
	TimeUnit            int64         `bson:"timeUnit" json:"timeUnit"`                 // Time unit in seconds for gathering
	Stage               string        `bson:"stage,omitempty" json:"stage,omitempty"`   // Stage of the gathering process
	MetalsGathered      int64         `bson:"metalsGathered" json:"metalsGathered"`     // Total metals gathered
	HydrogenGathered    int64         `bson:"hydrogenGathered" json:"hydrogenGathered"` // Total hydrogen gathered
	CrystalsGathered    int64         `bson:"crystalsGathered" json:"crystalsGathered"` // Total crystals gathered
	StartTime           time.Time     `bson:"startTime" json:"startTime"`               // When gathering started
	ProcessedAt         time.Time     `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}

// ShipLoadout captures per-ship-type socket configuration within a stack.
// It complements the static blueprint in ShipBlueprints by adding:
// - Sockets: up to 3 runes (see runes.go) in order
// - Anchored: whether currently anchored (e.g., for economic gathering)
// Notes:
//   - enforcement of anchoring rules and mining throughput penalties should be
//     handled at the game systems layer using these fields.
type ShipLoadout struct {
	Sockets  []Gem `bson:"sockets,omitempty" json:"sockets,omitempty"`
	Anchored bool  `bson:"anchored,omitempty" json:"anchored,omitempty"`
}

// GetOrInitLoadout returns the loadout for a ship type on this stack, creating
// a default entry if missing (no sockets, not anchored).
func (s *ShipStack) GetOrInitLoadout(t ShipType) ShipLoadout {
	if s.Loadouts == nil {
		s.Loadouts = make(map[ShipType]ShipLoadout)
	}
	if l, ok := s.Loadouts[t]; ok {
		return l
	}
	// Return default loadout with empty sockets and not anchored
	l := ShipLoadout{
		Sockets:  []Gem{},
		Anchored: false,
	}
	s.Loadouts[t] = l
	return l
}

// EffectiveShip computes the effective stats and usable abilities for a given
// ship type in this stack, taking into account the stack's role and the ship's loadout.
// Returns the effective ship snapshot and its abilities list.
//
// DEPRECATED: Use EffectiveShipV2 instead for full modifier transparency.
// This version lacks formation bonuses, composition bonuses, and modifier tracking.
func (s *ShipStack) EffectiveShip(t ShipType) (Ship, []Ability) {
	// DEPRECATED: Use EffectiveShipV2Simple instead
	return s.EffectiveShipV2Simple(t, time.Now())
}

// EffectiveShipV2 computes effective stats using the V2 modifier system.
// This is the recommended method for getting ship stats with full transparency.
// Returns the effective ship, abilities, and the complete modifier stack for debugging.
func (s *ShipStack) EffectiveShipV2(t ShipType, bucketIndex int, now time.Time) (Ship, []Ability, *ModifierStack) {
	ship, abilities, modStack := ComputeEffectiveShipV2(s, t, bucketIndex, now, false, "")
	return ship, abilities, modStack
}

// EffectiveShipV2Simple is a simplified version that matches the old EffectiveShip signature.
// Use this as a drop-in replacement for EffectiveShip.
func (s *ShipStack) EffectiveShipV2Simple(t ShipType, now time.Time) (Ship, []Ability) {
	ship, abilities := QuickEffectiveShip(s, t, 0, now)
	return ship, abilities
}

// SetFormation changes the stack's formation and applies reconfiguration time.
func (s *ShipStack) SetFormation(formationType FormationType, now time.Time) time.Time {
	if s.SavedFormations == nil {
		s.SavedFormations = make(map[FormationType]FormationWithSlots)
	}

	fws, ok := s.SavedFormations[formationType]
	if !ok {
		fws = s.BuildAndSaveFormationLayout(formationType, now)
	}
	// ensure latest buckets are reflected before activating
	s.updateFormationAssignmentsFor(&fws)
	s.SavedFormations[formationType] = fws
	s.Formation = &fws

	reconfigTime := fws.Modifiers.ReconfigureTime
	s.FormationReconfigUntil = now.Add(time.Duration(reconfigTime) * time.Second)
	return s.FormationReconfigUntil
}

func (s *ShipStack) EnsureFormationInitialized(now time.Time) {
	if s.SavedFormations == nil {
		s.SavedFormations = make(map[FormationType]FormationWithSlots)
	}
	if s.Formation == nil {
		fws := s.BuildAndSaveFormationLayout(FormationLine, now)
		s.Formation = &fws
		s.FormationReconfigUntil = now
	}
}

func (s *ShipStack) BuildAndSaveFormationLayout(formationType FormationType, now time.Time) FormationWithSlots {
	formation := AutoAssignFormation(s.Ships, formationType, now)
	fws := FromFormation(formation)
	if s.SavedFormations == nil {
		s.SavedFormations = make(map[FormationType]FormationWithSlots)
	}
	s.SavedFormations[formationType] = fws
	return fws
}

// IsFormationReconfiguring reports whether the stack is still changing formations.
func (s *ShipStack) IsFormationReconfiguring(now time.Time) bool {
	return now.Before(s.FormationReconfigUntil)
}

// GetFormationPosition returns the formation position for a specific ship type and bucket.
func (s *ShipStack) GetFormationPosition(shipType ShipType, bucketIndex int) FormationPosition {
	if s.Formation == nil {
		return PositionFront // default if no formation set
	}
	formation := s.Formation.ToFormation()
	for _, assignment := range formation.Assignments {
		if assignment.ShipType == shipType && assignment.BucketIndex == bucketIndex {
			return assignment.Position
		}
	}

	return PositionFront // default fallback
}

// EffectiveShipInFormation computes effective stats including formation position bonuses.
//
// DEPRECATED: Use EffectiveShipInFormationV2 instead.
// This version manually combines modifiers and lacks source tracking.
func (s *ShipStack) EffectiveShipInFormation(t ShipType, bucketIndex int) (Ship, []Ability) {
	// DEPRECATED: Use EffectiveShipV2 instead
	return s.EffectiveShipV2Simple(t, time.Now())
}

// EffectiveShipInFormationV2 computes effective stats using the V2 system with full formation support.
// This is the recommended replacement for EffectiveShipInFormation.
// Returns ship, abilities, and modifier stack for debugging.
func (s *ShipStack) EffectiveShipInFormationV2(t ShipType, bucketIndex int, now time.Time) (Ship, []Ability, *ModifierStack) {
	return ComputeEffectiveShipV2(s, t, bucketIndex, now, false, "")
}

// EffectiveShipInCombat computes effective stats for combat with formation counter bonuses.
// Use this when calculating damage in battle.
func (s *ShipStack) EffectiveShipInCombat(t ShipType, bucketIndex int, enemyFormation FormationType, now time.Time) (Ship, []Ability, *ModifierStack) {
	return ComputeEffectiveShipV2(s, t, bucketIndex, now, true, enemyFormation)
}

// GetModifierBreakdownForShip returns a detailed breakdown of all modifiers affecting a ship.
// Useful for debugging and UI display.
func (s *ShipStack) GetModifierBreakdownForShip(t ShipType, bucketIndex int, now time.Time, inCombat bool) []ModifierSummary {
	return GetModifierBreakdown(s, t, bucketIndex, now, inCombat, "")
}

// GetEffectiveStackSpeed returns the stack's movement speed considering formation.
func (s *ShipStack) GetEffectiveStackSpeed() int {
	// Find the slowest ship in the stack (baseline stack speed)
	slowest := 99999
	for shipType := range s.Ships {
		bp, ok := ShipBlueprints[shipType]
		if ok && bp.Speed < slowest {
			slowest = bp.Speed
		}
	}

	if slowest == 99999 {
		return 0
	}

	// Apply formation speed multiplier
	if s.Formation != nil {
		formation := s.Formation.ToFormation()
		return formation.GetEffectiveSpeed(slowest)
	}

	return slowest
}

// UpdateFormationAssignments refreshes formation assignments after combat or bucket changes.
func (s *ShipStack) updateFormationAssignmentsFor(fws *FormationWithSlots) {
	var refreshed []FormationSlotAssignment
	for _, a := range fws.SlotAssignments {
		buckets, ok := s.Ships[a.ShipType]
		if !ok || len(buckets) == 0 {
			continue // drop assignments for missing ship types
		}

		// Rebind if BucketIndex is invalid
		if a.BucketIndex < 0 || a.BucketIndex >= len(buckets) {
			// approximate per-ship HP from previous data (may be zero)
			hpPer := 0
			if a.Count > 0 {
				hpPer = a.AssignedHP / a.Count
			}
			// choose closest HP bucket
			bestIdx := 0
			bestDiff := 1<<31 - 1
			for idx, b := range buckets {
				d := intAbs(b.HP - hpPer)
				if d < bestDiff {
					bestDiff = d
					bestIdx = idx
				}
			}
			a.BucketIndex = bestIdx
		}

		// Refresh to bucket-wide values
		b := buckets[a.BucketIndex]
		if b.Count <= 0 || b.HP <= 0 {
			continue // drop empty buckets
		}
		a.Count = b.Count
		a.AssignedHP = b.HP * b.Count
		refreshed = append(refreshed, a)
	}

	// 2) Enforce uniqueness: one assignment per (ShipType, BucketIndex)
	unique := make([]FormationSlotAssignment, 0, len(refreshed))
	seen := make(map[ShipType]map[int]bool)
	for _, a := range refreshed {
		m, ok := seen[a.ShipType]
		if !ok {
			m = make(map[int]bool)
			seen[a.ShipType] = m
		}
		if m[a.BucketIndex] {
			// duplicate for the same bucket; skip (bucket-wide)
			continue
		}
		m[a.BucketIndex] = true
		unique = append(unique, a)
	}

	fws.SlotAssignments = unique

	// 3) Add missing buckets: any canonical bucket without an assignment should be placed using overflow policy
	// Build assigned set
	assigned := make(map[ShipType]map[int]bool)
	for _, a := range fws.SlotAssignments {
		m, ok := assigned[a.ShipType]
		if !ok {
			m = make(map[int]bool)
			assigned[a.ShipType] = m
		}
		m[a.BucketIndex] = true
	}

	// Build a temporary Formation and positionCounts to reuse overflow selector
	tempFormation := Formation{Type: fws.Type}
	for _, a := range fws.SlotAssignments {
		tempFormation.Assignments = append(tempFormation.Assignments, a.FormationAssignment)
	}
	posCounts := findPositionCounts(fws)

	for st, buckets := range s.Ships {
		for idx, b := range buckets {
			if b.Count <= 0 || b.HP <= 0 {
				continue
			}
			if m := assigned[st]; m != nil && m[idx] {
				continue // already assigned
			}

			// Select position: optimal or overflow fallback
			pos := DetermineOptimalPosition(st, fws.Type)
			cap := GetMaxSlotsForPosition(fws.Type, pos)
			if cap > 0 && posCounts[pos] >= cap {
				if alt, ok := chooseOverflowPosition(&tempFormation, s.Ships, st, idx, posCounts); ok {
					pos = alt
				} else {
					continue // nowhere to place
				}
			}

			// Find free slot index
			slotIdx, ok := findFreeSlotIndex(fws, pos)
			if !ok {
				continue
			}

			layer := DetermineLayer(pos, st)
			newA := FormationSlotAssignment{
				FormationAssignment: FormationAssignment{
					Position:    pos,
					Layer:       layer,
					ShipType:    st,
					BucketIndex: idx,
					Count:       b.Count,
					AssignedHP:  b.HP * b.Count,
				},
				SlotIndex: slotIdx,
				SlotKey: func() string {
					if coord, ok := GetNextSlotCoordinate(fws.Type, pos, slotIdx); ok {
						return fmt.Sprintf("%.6f:%.6f", coord.X, coord.Y)
					}
					return ""
				}(),
				IsManuallyPlaced: false,
			}
			fws.SlotAssignments = append(fws.SlotAssignments, newA)
			// Update helpers
			tempFormation.Assignments = append(tempFormation.Assignments, newA.FormationAssignment)
			posCounts[pos]++
		}
	}
}

func (s *ShipStack) UpdateFormationAssignments() {
	if s.Formation != nil {
		s.updateFormationAssignmentsFor(s.Formation)
	}
	if s.SavedFormations != nil {
		for ft, f := range s.SavedFormations {
			s.updateFormationAssignmentsFor(&f)
			s.SavedFormations[ft] = f
		}
	}
}

// intAbs returns the absolute value of an int without importing math.
func intAbs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// findPositionCounts returns the current slot usage per position.
func findPositionCounts(fws *FormationWithSlots) map[FormationPosition]int {
	counts := make(map[FormationPosition]int)
	for _, a := range fws.SlotAssignments {
		if a.Count > 0 && a.AssignedHP > 0 {
			counts[a.Position]++
		}
	}
	return counts
}

// findFreeSlotIndex finds the smallest free SlotIndex for a position under capacity.
func findFreeSlotIndex(fws *FormationWithSlots, pos FormationPosition) (int, bool) {
	max := GetMaxSlotsForPosition(fws.Type, pos)
	if max <= 0 {
		return -1, false
	}
	used := make(map[int]bool)
	for _, a := range fws.SlotAssignments {
		if a.Position == pos {
			used[a.SlotIndex] = true
		}
	}
	for i := 0; i < max; i++ {
		if !used[i] {
			return i, true
		}
	}
	return -1, false
}
