package essences

import (
	"time"

	"github.com/nicoberrocal/galaxyCore/ships"
)

// CollectNodeIDsForPath returns all BioNode IDs across trees matching the given ShipStack BioTreePath.
// It's more efficient than CollectNodesForPath when you only need the node IDs.
func CollectNodeIDsForPath(path ships.BioTreePath) []string {
	out := make([]string, 0, 32)
	// Build trees
	trees := []*BioTree{
		BuildAquatica(),
		BuildFlora(),
		BuildFauna(),
		BuildMycelia(),
	}
	match := string(path)
	for _, tree := range trees {
		for _, tier := range tree.Tiers {
			for _, node := range tier {
				if node == nil {
					continue
				}
				if node.Path == match {
					out = append(out, node.ID)
				}
			}
		}
	}
	return out
}

// CollectNodesForPath returns all BioNodes across trees matching the given ShipStack BioTreePath.
// It scans all known trees (Aquatica, Flora, Fauna, Mycelia) and filters nodes by BioNode.Path.
func CollectNodesForPath(path ships.BioTreePath) []*BioNode {
	out := make([]*BioNode, 0, 32)
	// Build trees
	trees := []*BioTree{
		BuildAquatica(),
		BuildFlora(),
		BuildFauna(),
		BuildMycelia(),
	}
	match := string(path)
	for _, tree := range trees {
		for _, tier := range tree.Tiers {
			for _, node := range tier {
				if node == nil {
					continue
				}
				if node.Path == match {
					out = append(out, node)
				}
			}
		}
	}
	return out
}

func PopulateStackBioForPath(stack *ships.ShipStack, path ships.BioTreePath, now time.Time) {
	if stack == nil {
		return
	}
	bm := stack.EnsureBio(now)
	nodes := CollectNodesForPath(path)
	pathStr := string(path)
	if bm.ActivePath != pathStr {
		bm.Nodes = make(map[string]*ships.BioNodeRuntimeState, len(nodes))
		bm.ByShipType = make(map[ships.ShipType]map[string]*ships.BioNodeRuntimeState, 4)
		bm.ActivePath = pathStr
	}
	for _, bn := range nodes {
		if bn == nil {
			continue
		}
		rn := bm.Node(bn.ID).ForAllShips()
		passive := bn.Effect
		if bn.Tradeoff != nil {
			passive = ships.CombineMods(passive, *bn.Tradeoff)
		}
		if !isZeroMods(passive) {
			rn.WithPassive(passive)
		}
		for _, ce := range bn.ComplexEffects {
			if ce.PrimaryEffect != nil && ce.Duration > 0 {
				dur := time.Duration(ce.Duration) * time.Second
				cd := time.Duration(ce.Cooldown) * time.Second
				rn.WithTriggered(*ce.PrimaryEffect, dur, cd)
			}
		}
	}
}

func init() {
	ships.BioPopulateFromPath = PopulateStackBioFromPath
	ships.BioPopulateFromExplicitPath = PopulateStackBioForPath
}

// PopulateStackBioFromPath ensures the stack's BioMachine exists and configures nodes matching its BioTreePath.
// All matching nodes are considered unlocked and set up as passive by default; simple triggered durations from
// ComplexEffects are attached for future event-driven activation. This avoids ad-hoc add/remove by using upsert semantics.
func PopulateStackBioFromPath(stack *ships.ShipStack, now time.Time) {
	if stack == nil {
		return
	}
	bm := stack.EnsureBio(now)

	nodes := CollectNodesForPath(stack.BioTreePath)
	// If the active path changed, reset nodes to match the new path
	pathStr := string(stack.BioTreePath)
	if bm.ActivePath != pathStr {
		bm.Nodes = make(map[string]*ships.BioNodeRuntimeState, len(nodes))
		bm.ByShipType = make(map[ships.ShipType]map[string]*ships.BioNodeRuntimeState, 4)
		bm.ActivePath = pathStr
	}
	for _, bn := range nodes {
		if bn == nil {
			continue
		}
		rn := bm.Node(bn.ID).ForAllShips()

		// Passive effect: apply base Effect plus any Tradeoff mods
		passive := bn.Effect
		if bn.Tradeoff != nil {
			passive = ships.CombineMods(passive, *bn.Tradeoff)
		}
		if !isZeroMods(passive) { // local helper to avoid importing builder's isZeroMods
			rn.WithPassive(passive)
		}

		// Map simple ComplexEffects with duration into a triggered configuration using PrimaryEffect
		for _, ce := range bn.ComplexEffects {
			if ce.PrimaryEffect != nil && ce.Duration > 0 {
				// Interpret tree durations in seconds for now
				dur := time.Duration(ce.Duration) * time.Second
				cd := time.Duration(ce.Cooldown) * time.Second
				rn.WithTriggered(*ce.PrimaryEffect, dur, cd)
			}
		}

		// Note: More advanced condition/tick/accumulate mappings can be wired later from ComplexEffects
	}
}

// isZeroMods is a local copy to avoid exporting internals from ships.
func isZeroMods(m ships.StatMods) bool {
	if m.Damage.LaserPct != 0 || m.Damage.NuclearPct != 0 || m.Damage.AntimatterPct != 0 {
		return false
	}
	if m.AttackIntervalPct != 0 || m.SpeedDelta != 0 || m.VisibilityDelta != 0 || m.AttackRangeDelta != 0 {
		return false
	}
	if m.LaserShieldDelta != 0 || m.NuclearShieldDelta != 0 || m.AntimatterShieldDelta != 0 {
		return false
	}
	if m.BucketHPPct != 0 || m.OutOfCombatRegenPct != 0 || m.AbilityCooldownPct != 0 || m.AtCombatRegenPct != 0 {
		return false
	}
	if m.TransportCapacityPct != 0 || m.WarpChargePct != 0 || m.WarpScatterPct != 0 || m.InterdictionResistPct != 0 {
		return false
	}
	if m.StructureDamagePct != 0 || m.SplashRadiusDelta != 0 || m.AccuracyPct != 0 || m.CritPct != 0 || m.FirstVolleyPct != 0 || m.ShieldPiercePct != 0 {
		return false
	}
	if m.CloakDetect || m.PingRangePct != 0 || m.EvasionPct != 0 || m.FormationSyncBonus != 0 || m.PositionFlexibility != 0 {
		return false
	}
	if m.GlobalDefensePct != 0 || m.HPPct != 0 {
		return false
	}
	return true
}
