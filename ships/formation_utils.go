package ships

import "fmt"

// Formation utility functions for common operations and debugging

// GetFormationInfo returns a human-readable summary of a formation.
func GetFormationInfo(formation *Formation) string {
	if formation == nil {
		return "No formation set"
	}
	
	spec, ok := FormationCatalog[formation.Type]
	if !ok {
		return fmt.Sprintf("Unknown formation: %s", formation.Type)
	}
	
	info := fmt.Sprintf("Formation: %s\n", spec.Name)
	info += fmt.Sprintf("Description: %s\n", spec.Description)
	info += fmt.Sprintf("Speed Multiplier: %.2fx\n", spec.SpeedMultiplier)
	info += fmt.Sprintf("Reconfigure Time: %ds\n", spec.ReconfigureTime)
	info += fmt.Sprintf("Special Properties: %v\n", spec.SpecialProperties)
	info += fmt.Sprintf("Assignments: %d\n", len(formation.Assignments))
	
	// Count ships per position
	positionCounts := make(map[FormationPosition]int)
	for _, assignment := range formation.Assignments {
		positionCounts[assignment.Position] += assignment.Count
	}
	
	info += "\nShip Distribution:\n"
	for position, count := range positionCounts {
		info += fmt.Sprintf("  %s: %d ships\n", position, count)
	}
	
	return info
}

// ValidateFormation checks if a formation is properly configured.
func ValidateFormation(formation *Formation) []string {
	var errors []string
	
	if formation == nil {
		return []string{"Formation is nil"}
	}
	
	// Check if formation type exists
	if _, ok := FormationCatalog[formation.Type]; !ok {
		errors = append(errors, fmt.Sprintf("Unknown formation type: %s", formation.Type))
	}
	
	// Check if there are assignments
	if len(formation.Assignments) == 0 {
		errors = append(errors, "Formation has no assignments")
	}
	
	// Check for empty assignments
	for i, assignment := range formation.Assignments {
		if assignment.Count <= 0 {
			errors = append(errors, fmt.Sprintf("Assignment %d has count <= 0", i))
		}
		if assignment.AssignedHP <= 0 {
			errors = append(errors, fmt.Sprintf("Assignment %d has HP <= 0", i))
		}
	}
	
	return errors
}

// CompareFormations returns tactical analysis comparing two formations.
func CompareFormations(attacker, defender FormationType) string {
	counterMult := GetFormationCounterMultiplier(attacker, defender)
	
	analysis := fmt.Sprintf("Formation Matchup: %s vs %s\n", attacker, defender)
	analysis += fmt.Sprintf("Counter Multiplier: %.2fx\n", counterMult)
	
	if counterMult > 1.1 {
		analysis += "Result: Strong advantage for attacker\n"
	} else if counterMult < 0.9 {
		analysis += "Result: Strong advantage for defender\n"
	} else {
		analysis += "Result: Even matchup\n"
	}
	
	attackerSpec := FormationCatalog[attacker]
	defenderSpec := FormationCatalog[defender]
	
	analysis += fmt.Sprintf("\n%s Properties:\n", attacker)
	analysis += fmt.Sprintf("  Speed: %.2fx\n", attackerSpec.SpeedMultiplier)
	analysis += fmt.Sprintf("  Reconfig: %ds\n", attackerSpec.ReconfigureTime)
	analysis += fmt.Sprintf("  Special: %v\n", attackerSpec.SpecialProperties)
	
	analysis += fmt.Sprintf("\n%s Properties:\n", defender)
	analysis += fmt.Sprintf("  Speed: %.2fx\n", defenderSpec.SpeedMultiplier)
	analysis += fmt.Sprintf("  Reconfig: %ds\n", defenderSpec.ReconfigureTime)
	analysis += fmt.Sprintf("  Special: %v\n", defenderSpec.SpecialProperties)
	
	return analysis
}

// GetCounterFormations returns formations that counter the given formation.
func GetCounterFormations(target FormationType) []FormationType {
	var counters []FormationType
	
	for formation := range FormationCatalog {
		mult := GetFormationCounterMultiplier(formation, target)
		if mult > 1.15 { // 15% or more advantage
			counters = append(counters, formation)
		}
	}
	
	return counters
}

// GetCounteredByFormations returns formations that are strong against the given formation.
func GetCounteredByFormations(target FormationType) []FormationType {
	var counteredBy []FormationType
	
	for formation := range FormationCatalog {
		mult := GetFormationCounterMultiplier(target, formation)
		if mult < 0.85 { // 15% or more disadvantage
			counteredBy = append(counteredBy, formation)
		}
	}
	
	return counteredBy
}

// CalculateFormationStrength returns a simple strength score for a formation.
func CalculateFormationStrength(formation *Formation, ships map[ShipType][]HPBucket) float64 {
	if formation == nil {
		return 0
	}
	
	strength := 0.0
	
	// Count total ships and HP
	for shipType, buckets := range ships {
		blueprint, ok := ShipBlueprints[shipType]
		if !ok {
			continue
		}
		
		for _, bucket := range buckets {
			// Base strength from HP and attack
			shipStrength := float64(bucket.Count) * float64(blueprint.HP+blueprint.AttackDamage)
			strength += shipStrength
		}
	}
	
	// Apply formation speed multiplier as mobility factor
	spec, ok := FormationCatalog[formation.Type]
	if ok {
		strength *= (1.0 + (spec.SpeedMultiplier-1.0)*0.5) // Speed contributes 50% weight
	}
	
	return strength
}

// SuggestFormationChanges analyzes current formation and suggests improvements.
func SuggestFormationChanges(stack *ShipStack, enemyFormation FormationType) []string {
	suggestions := []string{}
	
	if stack.Formation == nil {
		suggestions = append(suggestions, "No formation set - consider setting a formation for tactical bonuses")
		return suggestions
	}
	
	currentType := stack.Formation.Type
	
	// Check if we're being countered
	counterMult := GetFormationCounterMultiplier(currentType, enemyFormation)
	if counterMult < 0.9 {
		counters := GetCounterFormations(enemyFormation)
		if len(counters) > 0 {
			suggestions = append(suggestions, fmt.Sprintf("Current formation is weak vs %s (%.2fx). Consider: %v", 
				enemyFormation, counterMult, counters))
		}
	}
	
	// Check role mode synergy
	if stack.Role == RoleTactical && currentType == FormationBox {
		suggestions = append(suggestions, "Box formation is defensive - consider Vanguard or Line for Tactical mode")
	}
	if stack.Role == RoleEconomic && currentType == FormationVanguard {
		suggestions = append(suggestions, "Vanguard is aggressive - consider Box for Economic mode defense")
	}
	
	// Check composition bonuses
	_, activeBonuses := EvaluateCompositionBonuses(stack.Ships)
	if len(activeBonuses) == 0 {
		suggestions = append(suggestions, "No composition bonuses active - consider diversifying ship types")
	}
	
	// Validate current formation
	errors := ValidateFormation(stack.Formation)
	if len(errors) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Formation errors: %v", errors))
	}
	
	return suggestions
}

// GetOptimalPositionForAbility returns the best formation position to use an ability.
// This is a heuristic based on ability type.
func GetOptimalPositionForAbility(abilityID AbilityID) FormationPosition {
	_, ok := AbilitiesCatalog[abilityID]
	if !ok {
		return PositionFront
	}
	
	// Heuristic based on ability kind and purpose
	switch abilityID {
	case AbilityLongRangeSensors, AbilityPing, AbilityStandoffPattern, AbilitySiegePayload:
		return PositionBack
	case AbilityEvasiveManeuvers, AbilityInterdictorPulse, AbilityActiveCamo, AbilityBackstab:
		return PositionFlank
	case AbilityPointDefenseScreen, AbilityTargetingUplink, AbilitySelfRepair, AbilityRepairDrones:
		return PositionSupport
	default:
		return PositionFront
	}
}

// GetGemFamiliesForPosition returns gem families that synergize with a position.
func GetGemFamiliesForPosition(position FormationPosition) []GemFamily {
	familyMap := make(map[GemFamily]bool)
	
	for _, effect := range GemPositionEffectsCatalog {
		if effect.Position == position {
			familyMap[effect.GemFamily] = true
		}
	}
	
	var families []GemFamily
	for family := range familyMap {
		families = append(families, family)
	}
	
	return families
}

// AnalyzePositionEffectiveness evaluates how well ships are positioned in a formation.
func AnalyzePositionEffectiveness(stack *ShipStack) map[FormationPosition]float64 {
	effectiveness := make(map[FormationPosition]float64)
	
	if stack.Formation == nil {
		return effectiveness
	}
	
	for _, assignment := range stack.Formation.Assignments {
		optimal := DetermineOptimalPosition(assignment.ShipType, stack.Formation.Type)
		
		// Calculate effectiveness score
		score := 1.0
		if assignment.Position == optimal {
			score = 1.0 // Perfect positioning
		} else if assignment.Position == PositionFront && optimal == PositionFlank {
			score = 0.8 // Reasonable sub-optimal
		} else if assignment.Position == PositionFlank && optimal == PositionFront {
			score = 0.8 // Reasonable sub-optimal
		} else {
			score = 0.6 // Poor positioning
		}
		
		// Weight by ship count
		weight := float64(assignment.Count)
		effectiveness[assignment.Position] += score * weight
	}
	
	// Normalize by position
	positionCounts := make(map[FormationPosition]float64)
	for _, assignment := range stack.Formation.Assignments {
		positionCounts[assignment.Position] += float64(assignment.Count)
	}
	
	for position := range effectiveness {
		if positionCounts[position] > 0 {
			effectiveness[position] /= positionCounts[position]
		}
	}
	
	return effectiveness
}

// GetFormationRecommendations provides formation suggestions based on fleet composition.
func GetFormationRecommendations(ships map[ShipType][]HPBucket, role RoleMode) []FormationType {
	recommendations := []FormationType{}
	
	// Count ship types
	counts := make(map[ShipType]int)
	for shipType, buckets := range ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		counts[shipType] = total
	}
	
	// Scout-heavy fleets
	if counts[Scout] >= 3 {
		recommendations = append(recommendations, FormationSkirmish, FormationSwarm)
	}
	
	// Fighter-heavy fleets
	if counts[Fighter] >= 5 {
		recommendations = append(recommendations, FormationLine, FormationVanguard)
	}
	
	// Balanced fleets
	if len(counts) >= 3 {
		recommendations = append(recommendations, FormationLine, FormationEchelon)
	}
	
	// Defensive fleets (Carriers)
	if counts[Carrier] >= 1 {
		recommendations = append(recommendations, FormationBox, FormationPhalanx)
	}
	
	// Role-based recommendations
	switch role {
	case RoleTactical:
		if !contains(recommendations, FormationVanguard) {
			recommendations = append(recommendations, FormationVanguard)
		}
	case RoleEconomic:
		if !contains(recommendations, FormationBox) {
			recommendations = append(recommendations, FormationBox)
		}
	case RoleRecon:
		if !contains(recommendations, FormationSwarm) {
			recommendations = append(recommendations, FormationSwarm)
		}
	}
	
	return recommendations
}

// Helper function
func contains(slice []FormationType, item FormationType) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ExportFormationTemplate creates a template from current formation configuration.
func ExportFormationTemplate(stack *ShipStack, name, description string) FormationTemplate {
	template := FormationTemplate{
		Name:        name,
		Description: description,
		Assignments: make(map[ShipType]FormationPosition),
		Conditions:  []TemplateCondition{},
	}
	
	if stack.Formation != nil {
		template.Formation = stack.Formation.Type
		
		// Record ship type assignments
		for _, assignment := range stack.Formation.Assignments {
			template.Assignments[assignment.ShipType] = assignment.Position
		}
	}
	
	// Add minimum ship requirements
	minShips := make(map[ShipType]int)
	for shipType, buckets := range stack.Ships {
		total := 0
		for _, bucket := range buckets {
			total += bucket.Count
		}
		if total > 0 {
			minShips[shipType] = 1 // Require at least 1 of each type present
		}
	}
	
	template.Conditions = append(template.Conditions, TemplateCondition{
		MinShips: minShips,
		RoleMode: stack.Role,
	})
	
	return template
}

// CloneFormation creates a deep copy of a formation.
func CloneFormation(original *Formation) *Formation {
	if original == nil {
		return nil
	}
	
	clone := &Formation{
		Type:      original.Type,
		Facing:    original.Facing,
		CreatedAt: original.CreatedAt,
		Version:   original.Version,
		Modifiers: FormationMods{
			SpeedMultiplier:   original.Modifiers.SpeedMultiplier,
			ReconfigureTime:   original.Modifiers.ReconfigureTime,
			PositionBonuses:   make(map[FormationPosition]StatMods),
			SpecialProperties: append([]string{}, original.Modifiers.SpecialProperties...),
		},
		Assignments: make([]FormationAssignment, len(original.Assignments)),
	}
	
	// Deep copy position bonuses
	for pos, mods := range original.Modifiers.PositionBonuses {
		clone.Modifiers.PositionBonuses[pos] = mods
	}
	
	// Copy assignments
	copy(clone.Assignments, original.Assignments)
	
	return clone
}
