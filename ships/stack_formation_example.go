package ships

import (
	"fmt"
	"sort"
	"time"
)

// ExampleStackFormationLifecycle demonstrates:
// - creating a stack with default line formation
// - saving alternative formation layouts for quick switching
// - switching to a saved layout and to an auto-built layout when none exists
func ExampleStackFormationLifecycle() {
	// 1) Create a sample fleet composition (HP buckets per ship type)
	ships := map[ShipType][]HPBucket{
		Fighter: {
			{HP: 100, Count: 5},
			{HP: 60, Count: 3},
		},
		Bomber: {
			{HP: 120, Count: 2},
		},
		Scout: {
			{HP: 40, Count: 4},
		},
	}

	now := time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC)
	s := &ShipStack{Ships: ships, CreatedAt: now}

	// 2) Ensure default formation (line) is initialized and active
	s.EnsureFormationInitialized(now)
	fmt.Println("Active:", s.Formation.Type)

	// 3) Save alternative formation layouts without activating them
	s.BuildAndSaveFormationLayout(FormationBox, now)
	s.BuildAndSaveFormationLayout(FormationPhalanx, now)

	// Show which formation layouts are saved
	saved := make([]string, 0, len(s.SavedFormations))
	for k := range s.SavedFormations {
		saved = append(saved, string(k))
	}
	sort.Strings(saved)
	fmt.Println("Saved:", saved)

	// 4) Switch to an existing saved layout (uses saved Box)
	s.SetFormation(FormationBox, now.Add(1*time.Minute))
	fmt.Println("Switched:", s.Formation.Type)

	// 5) Switch to a formation with no saved layout (auto-builds Vanguard, then activates it)
	s.SetFormation(FormationVanguard, now.Add(2*time.Minute))
	fmt.Println("AutoBuild switch:", s.Formation.Type)

	// Output:
	// Active: line
	// Saved: [box line phalanx]
	// Switched: box
	// AutoBuild switch: vanguard
}
