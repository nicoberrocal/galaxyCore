package ships

import (
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
	Role             RoleMode  `bson:"role" json:"role"`
	ReconfigureUntil time.Time `bson:"reconfigureUntil,omitempty" json:"reconfigureUntil,omitempty"`

	// Loadouts track per-ship-type socket configurations for this particular stack.
	// This allows two stacks to field the same ship type with different gem setups.
	Loadouts map[ShipType]ShipLoadout `bson:"loadouts,omitempty" json:"loadouts,omitempty"`

	// Formation defines the tactical positioning of ships within the stack
	Formation              *Formation `bson:"formation,omitempty" json:"formation,omitempty"`
	FormationReconfigUntil time.Time  `bson:"formationReconfigUntil,omitempty" json:"formationReconfigUntil,omitempty"`

	// Current activity and movement
	Movement  []*MovementState `bson:"movement,omitempty" json:"movement,omitempty"`
	Battle    *BattleState     `bson:"battle,omitempty" json:"battle,omitempty"`       // Combat state for free space battles
	Ability   *[]AbilityState  `bson:"ability,omitempty" json:"ability,omitempty"`     // Active ship ability state
	Gathering *GatheringState  `bson:"gathering,omitempty" json:"gathering,omitempty"` // Active gathering state
	Version   int64            `bson:"version" json:"version"`                         // For optimistic locking
}

// StartModeSwitch sets a new RoleMode for the entire stack and applies a
// reconfiguration timer based on RoleModesCatalog. Returns the ETA.
func (s *ShipStack) StartModeSwitch(newRole RoleMode, now time.Time) time.Time {
	spec, ok := RoleModesCatalog[newRole]
	reconfig := 180
	if ok && spec.ReconfigureSeconds > 0 {
		reconfig = spec.ReconfigureSeconds
	}
	s.Role = newRole
	s.ReconfigureUntil = now.Add(time.Duration(reconfig) * time.Second)
	return s.ReconfigureUntil
}

// IsReconfiguring reports whether the stack is still switching modes at the given time.
func (s *ShipStack) IsReconfiguring(now time.Time) bool {
	return now.Before(s.ReconfigureUntil)
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

// CanWarp checks RoleMode rules and anchoring to determine if warping is allowed.
func (s *ShipStack) CanWarp(t ShipType) bool {
	load := s.GetOrInitLoadout(t)
	spec, ok := RoleModesCatalog[s.Role]
	if !ok {
		return !load.Anchored
	}
	// Economic mode can be allowed to warp when not anchored; anchoring always disables warp.
	if load.Anchored {
		return false
	}
	return spec.WarpAllowed
}

type HPBucket struct {
	HP    int `bson:"hp"`    // current HP of ships in this bucket
	Count int `bson:"count"` // how many ships at this HP
}

// BattleState tracks combat information for stacks in free space or mining locations
type BattleState struct {
	IsInCombat      bool            `bson:"isInCombat"`                               // Currently engaged in battle
	EnemyStackID    []bson.ObjectID `bson:"enemyStackId,omitempty"`                   // Opponent stack ID
	EnemyPlayerID   []bson.ObjectID `bson:"enemyPlayerId,omitempty"`                  // Opponent player ID
	BattleStartedAt time.Time       `bson:"battleStartedAt,omitempty"`                // When battle began
	BattleLocation  string          `bson:"battleLocation,omitempty"`                 // "empty_space", "asteroid", "nebula"
	LocationID      bson.ObjectID   `bson:"locationId,omitempty"`                     // ID of asteroid/nebula if applicable
	ProcessedAt     time.Time       `bson:"ProcessedAt,omitempty" json:"ProcessedAt"` // Last time this state was processed
}

// MovementState tracks what the stack is currently doing in free space or at mining locations
type MovementState struct {
	State       string        `bson:"state" json:"state,omitempty"`                     // "traveling", "mining", "idle", "in_combat"
	TargetID    bson.ObjectID `bson:"targetId,omitempty" json:"targetId,omitempty"`     // Target for movement/mining
	TargetType  string        `bson:"targetType,omitempty" json:"targetType,omitempty"` // "asteroid", "nebula", "coordinate"
	StartX      float64       `bson:"startX,omitempty" json:"startX,omitempty"`         // Starting coordinates
	StartY      float64       `bson:"startY,omitempty" json:"startY"`                   // Starting coordinates
	Speed       int           `bson:"speed,omitempty" json:"speed"`                     // Speed of the stack
	TargetX     float64       `bson:"targetX,omitempty" json:"targetX"`                 // Target coordinates
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
	IsMining            bool          `bson:"isMining" json:"isMining"`                 // Whether the stack is currently mining
	TargetID            bson.ObjectID `bson:"targetId" json:"targetId"`                 // ID of the target being gathered from
	TargetType          string        `bson:"targetType" json:"targetType"`             // Type of target (e.g., "asteroid", "nebula")
	ResourcePerTimeUnit int64         `bson:"resourcePerTime" json:"resourcePerTime"`   // Amount of resource gathered per time unit
	TimeUnit            int64         `bson:"timeUnit" json:"timeUnit"`                 // Time unit in seconds for gathering
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
func (s *ShipStack) EffectiveShip(t ShipType) (Ship, []Ability) {
	bp, ok := ShipBlueprints[t]
	if !ok {
		return Ship{}, nil
	}
	loadout := s.GetOrInitLoadout(t)
	// Compute mods using the stack's role and the ship's loadout
	mods, grants, _ := ComputeLoadout(bp, s.Role, loadout)
	eff := ApplyStatModsToShip(bp, mods)

	// Filter abilities based on the stack's role
	abilities := FilterAbilitiesForMode(eff, s.Role, grants)
	return eff, abilities
}

// SetFormation changes the stack's formation and applies reconfiguration time.
func (s *ShipStack) SetFormation(formationType FormationType, now time.Time) time.Time {
	formation := AutoAssignFormation(s.Ships, formationType, now)
	s.Formation = &formation

	// Apply role mode bonus to reconfiguration time
	reconfigTime := formation.Modifiers.ReconfigureTime
	reconfigTime = RoleModeFormationBonus(s.Role, reconfigTime)

	s.FormationReconfigUntil = now.Add(time.Duration(reconfigTime) * time.Second)
	return s.FormationReconfigUntil
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

	for _, assignment := range s.Formation.Assignments {
		if assignment.ShipType == shipType && assignment.BucketIndex == bucketIndex {
			return assignment.Position
		}
	}

	return PositionFront // default fallback
}

// EffectiveShipInFormation computes effective stats including formation position bonuses.
func (s *ShipStack) EffectiveShipInFormation(t ShipType, bucketIndex int) (Ship, []Ability) {
	bp, ok := ShipBlueprints[t]
	if !ok {
		return Ship{}, nil
	}

	loadout := s.GetOrInitLoadout(t)

	// Base mods from role and loadout
	mods, grants, _ := ComputeLoadout(bp, s.Role, loadout)

	// Add formation position bonuses if formation is set
	if s.Formation != nil {
		position := s.GetFormationPosition(t, bucketIndex)

		// Apply formation position bonuses
		mods = s.Formation.ApplyPositionBonusesToShip(position, mods)

		// Apply role-specific formation enhancements
		mods = ApplyFormationRoleModifiers(mods, s.Formation, position, s.Role)

		// Apply gem-position synergy bonuses
		gemPositionMods := ApplyGemPositionEffects(loadout.Sockets, position)
		mods = CombineMods(mods, gemPositionMods)
	}

	// Apply composition bonuses
	compositionMods, _ := EvaluateCompositionBonuses(s.Ships)
	mods = CombineMods(mods, compositionMods)

	eff := ApplyStatModsToShip(bp, mods)
	abilities := FilterAbilitiesForMode(eff, s.Role, grants)

	return eff, abilities
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
		return s.Formation.GetEffectiveSpeed(slowest)
	}

	return slowest
}

// UpdateFormationAssignments refreshes formation assignments after combat or bucket changes.
func (s *ShipStack) UpdateFormationAssignments() {
	if s.Formation == nil {
		return
	}

	// Update HP values for each assignment based on current buckets
	for i := range s.Formation.Assignments {
		assignment := &s.Formation.Assignments[i]

		if buckets, ok := s.Ships[assignment.ShipType]; ok {
			if assignment.BucketIndex < len(buckets) {
				bucket := buckets[assignment.BucketIndex]
				assignment.Count = bucket.Count
				assignment.AssignedHP = bucket.HP * bucket.Count
			} else {
				// Bucket no longer exists, mark for removal
				assignment.Count = 0
				assignment.AssignedHP = 0
			}
		} else {
			// Ship type no longer in stack
			assignment.Count = 0
			assignment.AssignedHP = 0
		}
	}

	// Remove dead assignments
	var activeAssignments []FormationAssignment
	for _, assignment := range s.Formation.Assignments {
		if assignment.Count > 0 && assignment.AssignedHP > 0 {
			activeAssignments = append(activeAssignments, assignment)
		}
	}
	s.Formation.Assignments = activeAssignments
}
