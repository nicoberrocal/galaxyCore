package ships

// Ship defines a ship TYPE blueprint (not a single unit instance) with its
// base attributes and static equipment/metadata. Runtime state for stacks,
// HP buckets, cooldowns, etc. lives in stack-related structures.
type Ship struct {
	ShipType         string
	AttackType       string
	LaserShield      int
	NuclearShield    int
	AntimatterShield int
	Speed            int
	VisibilityRange  int
	AttackRange      int
	HP               int
	AttackDamage     int
	AttackInterval   float64

	// Abilities defines this ship type's available tactical/strategic tools.
	// The full catalog is described in abilities.go. A ship may have up to 3.
	Abilities []Ability

	// RoleMode lets a ship type adopt a posture: Tactical/Economic/Recon/Scientific.
	// This is a soft-reconfiguration with tradeoffs, not a full role swap.
	// See roles.go for details and ApplyRoleMode() helper to compute modifiers.
	RoleMode RoleMode

	// Sockets are now managed per-ship-type in the stack's ShipLoadout
	// Construction costs
	MetalCost   int
	CrystalCost int
	PlasmaCost  int
	// Carrier specific
	TransportCapacity int
	CanTransport      []string
}

// Ability is a static definition used by ship blueprints and catalogs.
// Active abilities use CooldownSeconds/DurationSeconds. Passives typically have 0.
// The concrete runtime state (timers/activation) is tracked separately in AbilityState.
type Ability struct {
	// ID is the stable identifier (e.g. "LightSpeed", "AlphaStrike").
	ID AbilityID
	// Name is a human-readable label for UI.
	Name string
	// Kind defines the interaction model (passive/active/toggle/aura/travel/conditional).
	Kind AbilityKind
	// CooldownSeconds is the cooldown after use. 0 for passives or always-on toggles.
	CooldownSeconds int
	// DurationSeconds is the active window for temporary effects. 0 for instant or toggle.
	DurationSeconds int
	// Description documents tactical/strategic intent and effect summary.
	Description string
}
