package ships

// Economy helpers for mining/salvage throughput.
// These are simple, declarative baselines intended to be used by the gathering
// system (e.g., when resolving income ticks for a stack anchored on an asteroid/nebula).
//
// Design intent:
// - Drones are the specialist baseline (1.0x).
// - Other hulls can contribute in Economic mode at a fraction of a Drone's rate.
// - Anchoring is required for mining (handled by gameplay system).
// - Gem interactions are intentionally not hard-coded here; gameplay layer may
//   choose to apply additional bonuses based on EvaluateGemSockets(loadout.Sockets).

// EconomicCap is the max fraction of a Drone's throughput a hull can achieve
// while in Economic role mode.
var EconomicCap = map[ShipType]float64{
	Drone:     1.00, // 100%
	Scout:     0.25, // 25%
	Fighter:   0.25, // 25%
	Carrier:   0.40, // 40% when anchored; attracts raids
	Destroyer: 0.10, // 10% (salvage-focused)
	Bomber:    0.10, // 10% (salvage-focused)
}

// EconomicThroughputMultiplier returns [0..1] relative to a Drone, based on
// the stack's role and the ship's loadout. Returns 0 if not in Economic mode or
// if not anchored.
func (s *ShipStack) EconomicThroughputMultiplier(t ShipType) float64 {
	load := s.GetOrInitLoadout(t)
	if !load.Anchored {
		return 0.0
	}
	cap, ok := EconomicCap[t]
	if !ok {
		return 0.0
	}
	return cap
}
