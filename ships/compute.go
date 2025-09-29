package ships

// ComputeLoadout resolves the cumulative StatMods and GemWord-granted abilities
// for a ship blueprint based on the stack's role and the ship's loadout. It returns:
// - combined StatMods from RoleMode + sockets + GemWords
// - abilities granted by GemWords (as AbilityIDs)
// - matched GemWords (for UI/debug)
func ComputeLoadout(s Ship, role RoleMode, loadout ShipLoadout) (StatMods, []AbilityID, []GemWord) {
    roleMods := RoleModeMods(role, s.ShipType)
    socketMods, grants, matched := EvaluateGemSockets(loadout.Sockets)
    combined := CombineMods(roleMods, socketMods)
    return combined, grants, matched
}

// DamageMultiplierFor returns the damage multiplier for the ship's current attack type.
// Caller applies this to base AttackDamage when displaying effective damage.
func DamageMultiplierFor(s Ship, mods StatMods) float64 {
    switch s.AttackType {
    case "Laser":
        return 1.0 + mods.Damage.LaserPct
    case "Nuclear":
        return 1.0 + mods.Damage.NuclearPct
    case "Antimatter":
        return 1.0 + mods.Damage.AntimatterPct
    default:
        return 1.0
    }
}

// EffectiveAttackInterval applies AttackIntervalPct to the base interval and returns the result.
func EffectiveAttackInterval(s Ship, mods StatMods) float64 {
    v := s.AttackInterval * (1.0 + mods.AttackIntervalPct)
    if v < 0.1 { // safety clamp
        v = 0.1
    }
    return v
}

// FilterAbilitiesForMode returns the abilities usable in the stack's current RoleMode.
// It takes the ship's built-in abilities, adds GemWord-granted abilities, then
// applies Disabled/Enabled lists from RoleModesCatalog.
func FilterAbilitiesForMode(s Ship, role RoleMode, runewordGrants []AbilityID) []Ability {
    spec, ok := RoleModesCatalog[role]
    if !ok {
        // Unknown mode, return baseline abilities only
        base := make([]Ability, 0, len(s.Abilities))
        base = append(base, s.Abilities...)
        return base
    }

    // Build a set for disabled IDs for quick lookup
    disabled := make(map[AbilityID]struct{}, len(spec.DisabledAbilities))
    for _, id := range spec.DisabledAbilities {
        disabled[id] = struct{}{}
    }

    // De-dup base abilities by ID
    out := make([]Ability, 0, len(s.Abilities)+len(runewordGrants)+len(spec.EnabledAbilities))
    seen := map[AbilityID]struct{}{}

    // Helper to append if not disabled and not seen
    appendIfAllowed := func(a Ability) {
        if _, isDisabled := disabled[a.ID]; isDisabled {
            return
        }
        if _, exists := seen[a.ID]; exists {
            return
        }
        out = append(out, a)
        seen[a.ID] = struct{}{}
    }

    // Base abilities
    for _, a := range s.Abilities {
        appendIfAllowed(a)
    }
    // Runeword-granted abilities
    for _, id := range runewordGrants {
        appendIfAllowed(abilityByID(id))
    }
    // Mode-enabled abilities (ensure included while in this mode)
    for _, id := range spec.EnabledAbilities {
        appendIfAllowed(abilityByID(id))
    }
    return out
}

// ApplyStatModsToShip computes a presentational "effective" Ship snapshot by applying StatMods.
// Note: This does not persist or mutate runtime state; it's for UI calculations.
func ApplyStatModsToShip(base Ship, mods StatMods) Ship {
    s := base
    s.Speed += mods.SpeedDelta
    s.VisibilityRange += mods.VisibilityDelta
    s.AttackRange += mods.AttackRangeDelta

    s.LaserShield += mods.LaserShieldDelta
    s.NuclearShield += mods.NuclearShieldDelta
    s.AntimatterShield += mods.AntimatterShieldDelta

    // Damage is multiplicative and type-dependent; update AttackDamage accordingly
    s.AttackDamage = int(float64(s.AttackDamage) * DamageMultiplierFor(base, mods))
    s.AttackInterval = EffectiveAttackInterval(base, mods)
    // BucketHPPct modifies per-bucket HP; we reflect on base HP for preview purposes only
    s.HP = int(float64(s.HP) * (1.0 + mods.BucketHPPct))
    // Transport capacity percentage
    s.TransportCapacity = int(float64(s.TransportCapacity) * (1.0 + mods.TransportCapacityPct))
    return s
}

// Internal: fetch ability from catalog with a safe fallback for missing data.
func abilityByID(id AbilityID) Ability {
    if a, ok := AbilitiesCatalog[id]; ok {
        return a
    }
    return Ability{ID: id, Name: string(id), Kind: AbilityPassive, Description: "(missing from catalog)"}
}
