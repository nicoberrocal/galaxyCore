package ships

// DEPRECATED: This file has been replaced by the Formation Tree system.
// See formation_tree.go, formation_tree_catalog.go, and related files.
// This file is kept only for package compilation and will be removed in future cleanup.
//
// The old hardcoded synergy system (ability-formation, gem-position, composition bonuses)
// has been replaced by a player-driven progression tree where players choose their bonuses
// by unlocking nodes with experience points.
//
// Migration path:
// 1. Old system: Hardcoded bonuses in catalogs
// 2. New system: Player unlocks nodes in formation trees to gain bonuses
// 3. Bonuses are applied through ModifierBuilder.AddFormationTreeNodes()
