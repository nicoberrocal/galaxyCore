package ships

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BioNodeStage represents the runtime stage of a bio node.
// Stages are designed to cover passive, triggered-with-duration, cooldown gating,
// as well as composite states like accumulating counters and periodic tick effects.
type BioNodeStage string

const (
	BioStagePassive          BioNodeStage = "passive"          // always-on, no cooldown
	BioStageTriggered        BioNodeStage = "triggered"        // active for a duration, then transitions to cooldown
	BioStageCooldown         BioNodeStage = "cooldown"         // locked until cooldown ends
	BioStageAccumulating     BioNodeStage = "accumulating"     // builds up an accumulator/counter over ticks
	BioStageTicking          BioNodeStage = "ticking"          // applies periodic tick effects on schedule
	BioStageCompositeActive  BioNodeStage = "composite_active" // mix of active/ticking/accumulating
	BioStageCompositeCooloff BioNodeStage = "composite_cd"     // composite cooldown
)

// AbilityCastRef captures the ability that triggered a node, with the cast start time.
type AbilityCastRef struct {
	Ability   AbilityID `bson:"ability" json:"ability"`
	ShipType  ShipType  `bson:"shipType" json:"shipType"`
	StartTime time.Time `bson:"startTime" json:"startTime"`
}

// BioDebuffState is an inbound debuff applied by enemy bio/traits to this stack.
// These are consumed by the compute pipeline as debuff modifier layers.
type BioDebuffState struct {
	ID           string        `bson:"id" json:"id"`
	SourceStack  bson.ObjectID `bson:"sourceStack" json:"sourceStack"`
	SourceNodeID string        `bson:"sourceNodeId" json:"sourceNodeId"`
	Mods         StatMods      `bson:"mods,omitempty" json:"mods,omitempty"`
	Stacks       int           `bson:"stacks" json:"stacks"`
	MaxStacks    int           `bson:"maxStacks" json:"maxStacks"`
	AppliedAt    time.Time     `bson:"appliedAt" json:"appliedAt"`
	ExpiresAt    time.Time     `bson:"expiresAt" json:"expiresAt"`
}

// BioBuffState is an inbound ally buff applied by allied traits/nodes to this stack.
// These are consumed by the compute pipeline as buff modifier layers.
type BioBuffState struct {
	ID           string        `bson:"id" json:"id"`
	SourceStack  bson.ObjectID `bson:"sourceStack" json:"sourceStack"`
	SourceNodeID string        `bson:"sourceNodeId" json:"sourceNodeId"`
	Mods         StatMods      `bson:"mods,omitempty" json:"mods,omitempty"`
	Stacks       int           `bson:"stacks" json:"stacks"`
	MaxStacks    int           `bson:"maxStacks" json:"maxStacks"`
	AppliedAt    time.Time     `bson:"appliedAt" json:"appliedAt"`
	ExpiresAt    time.Time     `bson:"expiresAt" json:"expiresAt"`
	// Optional target scoping. If set, buff is intended for actions against this target stack.
	TargetStack bson.ObjectID `bson:"targetStack,omitempty" json:"targetStack,omitempty"`
	// Optional scope hint (e.g., "movement", "movement_attack_target", "combat").
	Scope string `bson:"scope,omitempty" json:"scope,omitempty"`
}

// BioActiveLayer represents a ready-to-apply modifier layer produced by the bio machine
// for the current snapshot. Builder will convert these into modifier stack layers.
type BioActiveLayer struct {
	Source    ModifierSource
	SourceID  string
	Desc      string
	Mods      StatMods   `bson:"mods,omitempty" json:"mods,omitempty"`
	ExpiresAt *time.Time // nil => permanent for the snapshot
	Priority  int
}

// BioNodeRuntimeState stores runtime, per-node state and its generated StatMods for each stage.
// To avoid import cycles, this struct carries StatMods directly and is fully self-contained.
type BioNodeRuntimeState struct {
	// Identity
	ID    string       `bson:"id" json:"id"`
	Stage BioNodeStage `bson:"stage" json:"stage"`

	// Applicability
	AllShips  bool              `bson:"allShips" json:"allShips"`
	ShipTypes map[ShipType]bool `bson:"shipTypes" json:"shipTypes"`

	// Targeting context for AoE or directed effects (ally/enemy selections or ship-type subsets).
	AllyTargets    []bson.ObjectID `bson:"allyTargets,omitempty" json:"allyTargets,omitempty"`
	EnemyTargets   []bson.ObjectID `bson:"enemyTargets,omitempty" json:"enemyTargets,omitempty"`
	AllyShipTypes  []ShipType      `bson:"allyShipTypes,omitempty" json:"allyShipTypes,omitempty"`
	EnemyShipTypes []ShipType      `bson:"enemyShipTypes,omitempty" json:"enemyShipTypes,omitempty"`

	// Core stage timing
	StartTime      time.Time     `bson:"startTime" json:"startTime"`
	EndTime        time.Time     `bson:"endTime" json:"endTime"`
	Duration       time.Duration `bson:"duration" json:"duration"`
	Cooldown       time.Duration `bson:"cooldown" json:"cooldown"`
	CooldownEndsAt time.Time     `bson:"cooldownEndsAt" json:"cooldownEndsAt"`
	LastTick       time.Time     `bson:"lastTick" json:"lastTick"`
	TickPeriod     time.Duration `bson:"tickPeriod" json:"tickPeriod"`

	// Counters and accumulators
	ActivationCount   int     `bson:"activationCount" json:"activationCount"`
	MaxActivations    int     `bson:"maxActivations" json:"maxActivations"`
	StackCount        int     `bson:"stackCount" json:"stackCount"`
	MaxStacks         int     `bson:"maxStacks" json:"maxStacks"`
	Accumulator       float64 `bson:"accumulator" json:"accumulator"`
	AccumulatePerTick float64 `bson:"accumulatePerTick" json:"accumulatePerTick"`
	AccumulateCap     float64 `bson:"accumulateCap" json:"accumulateCap"`

	// Trigger provenance (e.g. ability cast that triggered this node)
	TriggeredBy *AbilityCastRef `bson:"triggeredBy,omitempty" json:"triggeredBy,omitempty"`

	// Stage-based StatMods (kept small and cache-friendly)
	ModsPassive     StatMods `bson:"modsPassive,omitempty" json:"modsPassive,omitempty"`
	ModsTriggered   StatMods `bson:"modsTriggered,omitempty" json:"modsTriggered,omitempty"`
	ModsTick        StatMods `bson:"modsTick,omitempty" json:"modsTick,omitempty"`
	ModsAccumulated StatMods `bson:"modsAccumulated,omitempty" json:"modsAccumulated,omitempty"`

	// Outgoing debuff prototype (applied to enemies when node logic triggers).
	OutgoingDebuffID        string        `bson:"outgoingDebuffId,omitempty" json:"outgoingDebuffId,omitempty"`
	OutgoingDebuffMods      StatMods      `bson:"outgoingDebuffMods,omitempty" json:"outgoingDebuffMods,omitempty"`
	OutgoingDebuffDuration  time.Duration `bson:"outgoingDebuffDuration,omitempty" json:"outgoingDebuffDuration,omitempty"`
	OutgoingDebuffMaxStacks int           `bson:"outgoingDebuffMaxStacks,omitempty" json:"outgoingDebuffMaxStacks,omitempty"`

	// Internal linkage for fluent API
	parent *BioMachine `bson:"-" json:"-"`
}

// Fluent API: node scoping and configuration
func (n *BioNodeRuntimeState) ForAllShips() *BioNodeRuntimeState {
	n.AllShips = true
	return n
}
func (n *BioNodeRuntimeState) ForShip(t ShipType) *BioNodeRuntimeState {
	if n.ShipTypes == nil {
		n.ShipTypes = make(map[ShipType]bool)
	}
	n.ShipTypes[t] = true
	if n.parent != nil {
		n.parent.IndexNode(n)
	}
	return n
}
func (n *BioNodeRuntimeState) WithPassive(mods StatMods) *BioNodeRuntimeState {
	n.ModsPassive = mods
	return n
}
func (n *BioNodeRuntimeState) WithTriggered(mods StatMods, dur time.Duration, cd time.Duration) *BioNodeRuntimeState {
	n.ModsTriggered = mods
	n.Duration = dur
	n.Cooldown = cd
	return n
}
func (n *BioNodeRuntimeState) WithTick(mods StatMods, period time.Duration) *BioNodeRuntimeState {
	n.ModsTick = mods
	n.TickPeriod = period
	return n
}
func (n *BioNodeRuntimeState) WithAccumulate(perTick float64, cap float64, mods StatMods) *BioNodeRuntimeState {
	n.AccumulatePerTick = perTick
	n.AccumulateCap = cap
	n.ModsAccumulated = mods
	return n
}
func (n *BioNodeRuntimeState) TargetsAllies(ids ...bson.ObjectID) *BioNodeRuntimeState {
	n.AllyTargets = append(n.AllyTargets, ids...)
	return n
}
func (n *BioNodeRuntimeState) TargetsEnemies(ids ...bson.ObjectID) *BioNodeRuntimeState {
	n.EnemyTargets = append(n.EnemyTargets, ids...)
	return n
}
func (n *BioNodeRuntimeState) WithOutgoingDebuff(id string, mods StatMods, dur time.Duration, maxStacks int) *BioNodeRuntimeState {
	n.OutgoingDebuffID = id
	n.OutgoingDebuffMods = mods
	n.OutgoingDebuffDuration = dur
	n.OutgoingDebuffMaxStacks = maxStacks
	return n
}
func (n *BioNodeRuntimeState) Done() *BioMachine { return n.parent }

// CurrentLayers returns the set of active layers (if any) produced by this node for the given shipType.
// It infers the correct source and lifetime based on node stage and timers.
func (n *BioNodeRuntimeState) CurrentLayers(shipType ShipType, now time.Time) []BioActiveLayer {
	if !(n.AllShips || (n.ShipTypes != nil && n.ShipTypes[shipType])) {
		return nil
	}

	layers := make([]BioActiveLayer, 0, 2)
	switch n.Stage {
	case BioStagePassive:
		if !isZeroMods(n.ModsPassive) {
			layers = append(layers, BioActiveLayer{
				Source:    SourceBioPassive,
				SourceID:  n.ID,
				Desc:      "Bio Passive: " + n.ID,
				Mods:      n.ModsPassive,
				ExpiresAt: nil,
				Priority:  PriorityBioPassive,
			})
		}
	case BioStageTriggered, BioStageCompositeActive:
		if !isZeroMods(n.ModsTriggered) {
			var exp *time.Time
			if !n.EndTime.IsZero() && now.Before(n.EndTime) {
				exp = &n.EndTime
			}
			layers = append(layers, BioActiveLayer{
				Source:    SourceBioTriggered,
				SourceID:  n.ID,
				Desc:      "Bio Triggered: " + n.ID,
				Mods:      n.ModsTriggered,
				ExpiresAt: exp,
				Priority:  PriorityBioTriggered,
			})
		}
		// Composite may also tick while active
		fallthrough
	case BioStageTicking:
		if !isZeroMods(n.ModsTick) && n.TickPeriod > 0 {
			// Tick-based contributions are treated as temporary pulses for this snapshot if on schedule
			if n.LastTick.IsZero() || now.Sub(n.LastTick) >= n.TickPeriod {
				layers = append(layers, BioActiveLayer{
					Source:    SourceBioTick,
					SourceID:  n.ID,
					Desc:      "Bio Tick: " + n.ID,
					Mods:      n.ModsTick,
					ExpiresAt: nil,
					Priority:  PriorityBioTriggered,
				})
			}
		}
	case BioStageAccumulating:
		if !isZeroMods(n.ModsAccumulated) && n.Accumulator > 0 {
			layers = append(layers, BioActiveLayer{
				Source:    SourceBioAccum,
				SourceID:  n.ID,
				Desc:      "Bio Accumulated: " + n.ID,
				Mods:      n.ModsAccumulated,
				ExpiresAt: nil,
				Priority:  PriorityBioPassive,
			})
		}
	case BioStageCooldown, BioStageCompositeCooloff:
		// no active layers while cooling down
	}
	return layers
}

type BioMachine struct {
	Nodes          map[string]*BioNodeRuntimeState              `bson:"nodes,omitempty" json:"nodes,omitempty"`
	InboundDebuffs map[string]*BioDebuffState                   `bson:"inboundDebuffs,omitempty" json:"inboundDebuffs,omitempty"`
	InboundBuffs   map[string]*BioBuffState                     `bson:"inboundBuffs,omitempty" json:"inboundBuffs,omitempty"`
	ByShipType     map[ShipType]map[string]*BioNodeRuntimeState `bson:"-" json:"-"`
	LastProcessed  time.Time                                    `bson:"lastProcessed" json:"lastProcessed"`
	ActivePath     string                                       `bson:"activePath,omitempty" json:"activePath,omitempty"`
	UnlockAll      bool                                         `bson:"unlockAll" json:"unlockAll"` // if true, treat all configured nodes as unlocked
}

var BioPopulateFromPath func(stack *ShipStack, now time.Time)
var BioPopulateFromExplicitPath func(stack *ShipStack, path BioTreePath, now time.Time)

func NewBioMachine(now time.Time) *BioMachine {
	return &BioMachine{
		Nodes:          make(map[string]*BioNodeRuntimeState, 16),
		InboundDebuffs: make(map[string]*BioDebuffState, 4),
		InboundBuffs:   make(map[string]*BioBuffState, 4),
		ByShipType:     make(map[ShipType]map[string]*BioNodeRuntimeState, 4),
		LastProcessed:  now,
		UnlockAll:      true,
	}
}

// Node returns a node builder for the given ID, creating it if missing.
func (bm *BioMachine) Node(id string) *BioNodeRuntimeState {
	n, ok := bm.Nodes[id]
	if !ok {
		n = &BioNodeRuntimeState{ID: id, Stage: BioStagePassive, parent: bm}
		bm.Nodes[id] = n
	}
	return n
}

// IndexNode adds the node to a per-shipType index for fast lookup.
func (bm *BioMachine) IndexNode(n *BioNodeRuntimeState) {
	if n.AllShips {
		return
	}
	for st := range n.ShipTypes {
		m := bm.ByShipType[st]
		if m == nil {
			m = make(map[string]*BioNodeRuntimeState, 4)
			bm.ByShipType[st] = m
		}
		m[n.ID] = n
	}
}

// Tick advances timers, performs periodic accumulation, and transitions stages.
// Designed to be O(#active nodes) and allocation-free on hot paths.
func (bm *BioMachine) Tick(now time.Time) {
	if now.Before(bm.LastProcessed) {
		bm.LastProcessed = now
	}
	dt := now.Sub(bm.LastProcessed)
	if dt <= 0 {
		return
	}

	for _, n := range bm.Nodes {
		switch n.Stage {
		case BioStageTriggered, BioStageCompositeActive:
			if !n.EndTime.IsZero() && !now.Before(n.EndTime) {
				// transition to cooldown
				n.Stage = BioStageCooldown
				n.CooldownEndsAt = now.Add(n.Cooldown)
			}
		case BioStageCooldown, BioStageCompositeCooloff:
			if !n.CooldownEndsAt.IsZero() && now.After(n.CooldownEndsAt) {
				// return to passive or accumulating baseline
				if n.AccumulatePerTick > 0 {
					n.Stage = BioStageAccumulating
				} else {
					n.Stage = BioStagePassive
				}
			}
		}

		// Ticking and accumulation
		if n.TickPeriod > 0 && (n.Stage == BioStageTicking || n.Stage == BioStageCompositeActive) {
			if n.LastTick.IsZero() || now.Sub(n.LastTick) >= n.TickPeriod {
				n.LastTick = now
			}
		}
		if n.AccumulatePerTick > 0 && (n.Stage == BioStageAccumulating || n.Stage == BioStageCompositeActive) {
			n.Accumulator += n.AccumulatePerTick * dt.Seconds()
			if n.Accumulator > n.AccumulateCap && n.AccumulateCap > 0 {
				n.Accumulator = n.AccumulateCap
			}
		}
	}

	// Expire inbound debuffs
	for id, d := range bm.InboundDebuffs {
		if now.After(d.ExpiresAt) {
			delete(bm.InboundDebuffs, id)
		}
	}

	// Expire inbound buffs
	for id, b := range bm.InboundBuffs {
		if now.After(b.ExpiresAt) {
			delete(bm.InboundBuffs, id)
		}
	}

	bm.LastProcessed = now
}

// OnAbilityCast wires ability usage to bio nodes that listen for ability-trigger transitions.
// Upsert semantics ensure we avoid ad-hoc add/remove across ticks.
func (bm *BioMachine) OnAbilityCast(ability AbilityID, shipType ShipType, start time.Time) {
	cast := &AbilityCastRef{Ability: ability, ShipType: shipType, StartTime: start}
	// Trigger all nodes that have a triggered stage configured for this shipType
	for _, n := range bm.Nodes {
		if !(n.AllShips || (n.ShipTypes != nil && n.ShipTypes[shipType])) {
			continue
		}
		if isZeroMods(n.ModsTriggered) || n.Duration <= 0 {
			continue
		}
		// only if not in cooldown
		if n.Stage == BioStageCooldown || n.Stage == BioStageCompositeCooloff {
			continue
		}
		n.Stage = BioStageTriggered
		n.TriggeredBy = cast
		n.StartTime = start
		n.EndTime = start.Add(n.Duration)
		n.ActivationCount++
		if n.MaxActivations > 0 && n.ActivationCount >= n.MaxActivations {
			// immediately set to cooldown once it ends
		}
	}
}

// ApplyInboundDebuff upserts a debuff applied by an enemy trait/node to this stack.
func (bm *BioMachine) ApplyInboundDebuff(id string, mods StatMods, duration time.Duration, stacks int, maxStacks int, sourceStack bson.ObjectID, sourceNodeID string, now time.Time) {
	d, ok := bm.InboundDebuffs[id]
	if !ok {
		d = &BioDebuffState{ID: id, Mods: mods, Stacks: 0, MaxStacks: maxStacks, AppliedAt: now}
		bm.InboundDebuffs[id] = d
	}
	// stack with clamp
	d.Stacks += stacks
	if d.MaxStacks > 0 && d.Stacks > d.MaxStacks {
		d.Stacks = d.MaxStacks
	}
	d.SourceStack = sourceStack
	d.SourceNodeID = sourceNodeID
	d.ExpiresAt = now.Add(duration)
}

// ApplyInboundBuff upserts a buff applied by an allied trait/node to this stack.
func (bm *BioMachine) ApplyInboundBuff(
	id string,
	mods StatMods,
	duration time.Duration,
	stacks int,
	maxStacks int,
	sourceStack bson.ObjectID,
	sourceNodeID string,
	targetStack bson.ObjectID,
	scope string,
	now time.Time,
) {
	b, ok := bm.InboundBuffs[id]
	if !ok {
		b = &BioBuffState{ID: id, Mods: mods, Stacks: 0, MaxStacks: maxStacks, AppliedAt: now}
		bm.InboundBuffs[id] = b
	}
	b.Stacks += stacks
	if b.MaxStacks > 0 && b.Stacks > b.MaxStacks {
		b.Stacks = b.MaxStacks
	}
	b.SourceStack = sourceStack
	b.SourceNodeID = sourceNodeID
	b.TargetStack = targetStack
	b.Scope = scope
	b.ExpiresAt = now.Add(duration)
}

// CollectActiveLayersForShip returns a list of bio-generated layers for a specific shipType at time now.
func (bm *BioMachine) CollectActiveLayersForShip(shipType ShipType, now time.Time) []BioActiveLayer {
	layers := make([]BioActiveLayer, 0, 8)
	if idx, ok := bm.ByShipType[shipType]; ok {
		for _, n := range idx {
			layers = append(layers, n.CurrentLayers(shipType, now)...)
		}
	}
	// Also include global nodes
	for _, n := range bm.Nodes {
		if n.AllShips {
			layers = append(layers, n.CurrentLayers(shipType, now)...)
		}
	}
	return layers
}

// CollectInboundDebuffs returns the current active inbound debuffs.
func (bm *BioMachine) CollectInboundDebuffs(now time.Time) []BioDebuffState {
	res := make([]BioDebuffState, 0, len(bm.InboundDebuffs))
	for _, d := range bm.InboundDebuffs {
		if now.Before(d.ExpiresAt) {
			res = append(res, *d)
		}
	}
	return res
}

// CollectInboundBuffs returns the current active inbound ally buffs.
func (bm *BioMachine) CollectInboundBuffs(now time.Time) []BioBuffState {
	res := make([]BioBuffState, 0, len(bm.InboundBuffs))
	for _, b := range bm.InboundBuffs {
		if now.Before(b.ExpiresAt) {
			res = append(res, *b)
		}
	}
	return res
}

// CollectInboundBuffsForTarget returns active inbound buffs applicable to a specific target.
// A buff applies if it has no TargetStack set, or matches the provided target.
func (bm *BioMachine) CollectInboundBuffsForTarget(target bson.ObjectID, now time.Time) []BioBuffState {
	res := make([]BioBuffState, 0, len(bm.InboundBuffs))
	for _, b := range bm.InboundBuffs {
		if !now.Before(b.ExpiresAt) {
			continue
		}
		if b.TargetStack.IsZero() || b.TargetStack == target {
			res = append(res, *b)
		}
	}
	return res
}
