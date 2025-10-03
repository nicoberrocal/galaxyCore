package ships

import (
	"time"
)

// BattleResult contains information about a completed battle for XP calculation.
type BattleResult struct {
	Victory           bool
	FlawlessVictory   bool // No ships lost
	OutnumberedWin    bool // Won with fewer ships
	EnemyShipsDestroyed int
	DamageDone        int
	DamageTaken       int
	FormationUsed     FormationType
	BattleDuration    time.Duration
	EnemyFormation    FormationType
	CounterAdvantage  bool // Used formation with advantage
}

// CalculateExperienceGain calculates XP earned from a battle.
func CalculateExperienceGain(result BattleResult, now time.Time) ExperienceGain {
	xp := 0
	
	// Base XP from enemy ships destroyed
	xp += result.EnemyShipsDestroyed * 10
	
	// Damage contribution
	xp += result.DamageDone / 100
	
	// Victory bonus
	if result.Victory {
		xp += 50
	}
	
	// Flawless victory (no losses)
	if result.FlawlessVictory {
		xp += 100
	}
	
	// Outnumbered bonus
	if result.OutnumberedWin {
		xp *= 2
	}
	
	// Counter advantage bonus (smart formation choice)
	if result.CounterAdvantage {
		xp += 25
	}
	
	// Duration bonus (longer battles give more XP)
	if result.BattleDuration > 5*time.Minute {
		xp += 20
	}
	
	// Minimum XP (even losses grant something)
	if xp < 5 {
		xp = 5
	}
	
	description := "Battle experience"
	if result.Victory {
		description = "Victory against enemy fleet"
	}
	if result.FlawlessVictory {
		description = "Flawless victory"
	}
	
	return ExperienceGain{
		Source:       "combat",
		Amount:       xp,
		FormationType: result.FormationUsed,
		Timestamp:    now,
		Description:  description,
	}
}

// CalculateDailyLoginXP awards daily login bonus.
func CalculateDailyLoginXP(consecutiveDays int, now time.Time) ExperienceGain {
	baseXP := 10
	bonus := consecutiveDays * 2
	
	if consecutiveDays >= 7 {
		bonus += 20 // Weekly streak bonus
	}
	if consecutiveDays >= 30 {
		bonus += 50 // Monthly streak bonus
	}
	
	return ExperienceGain{
		Source:      "daily_login",
		Amount:      baseXP + bonus,
		Timestamp:   now,
		Description: "Daily login bonus",
	}
}

// CalculateQuestXP awards XP for completing quests/achievements.
func CalculateQuestXP(questType string, difficulty int, now time.Time) ExperienceGain {
	baseXP := 20
	difficultyMult := float64(difficulty)
	
	xp := int(float64(baseXP) * difficultyMult)
	
	return ExperienceGain{
		Source:      "quest",
		Amount:      xp,
		Timestamp:   now,
		Description: "Quest: " + questType,
	}
}

// CalculateFormationMasteryXP awards XP for using a formation extensively.
func CalculateFormationMasteryXP(formation FormationType, battlesWithFormation int, now time.Time) ExperienceGain {
	// Diminishing returns: first 10 battles give 5 XP each, then less
	xp := 0
	
	if battlesWithFormation <= 10 {
		xp = battlesWithFormation * 5
	} else if battlesWithFormation <= 50 {
		xp = 50 + (battlesWithFormation-10)*2
	} else {
		xp = 50 + 80 + (battlesWithFormation-50)
	}
	
	return ExperienceGain{
		Source:       "formation_mastery",
		Amount:       xp,
		FormationType: formation,
		Timestamp:    now,
		Description:  "Formation mastery progress",
	}
}

// AwardExperienceForBattle is the primary entry point for battle XP.
func AwardExperienceForBattle(treeState *FormationTreeState, result BattleResult, now time.Time) {
	if treeState == nil {
		return
	}
	
	gain := CalculateExperienceGain(result, now)
	treeState.AwardExperience(gain)
}

// GrantMonthlyFreeReset grants a free reset if enough time has passed.
func GrantMonthlyFreeReset(treeState *FormationTreeState, now time.Time) bool {
	if treeState == nil {
		return false
	}
	
	maxFreeResets := 3 // Cap at 3 stored
	
	// Check if 30 days have passed since last free reset grant
	if treeState.NextFreeResetAt.IsZero() || now.After(treeState.NextFreeResetAt) {
		treeState.GrantFreeReset(now, maxFreeResets)
		return true
	}
	
	return false
}

// GetXPRequiredForNextRank calculates XP needed for next admiral rank.
func GetXPRequiredForNextRank(currentRank int) int {
	// Exponential scaling: Rank 1 needs 100 XP, each rank doubles
	baseXP := 100
	return baseXP * (1 << uint(currentRank)) // 2^rank
}

// CalculateAdmiralRank determines rank based on total XP.
func CalculateAdmiralRank(totalXP int) int {
	rank := 0
	xpNeeded := 0
	
	for rank < 20 { // Cap at rank 20
		xpNeeded += GetXPRequiredForNextRank(rank)
		if totalXP < xpNeeded {
			break
		}
		rank++
	}
	
	return rank
}

// GetRankTitle returns the admiral rank title.
func GetRankTitle(rank int) string {
	titles := []string{
		"Ensign",           // 0
		"Lieutenant",       // 1
		"Commander",        // 2
		"Captain",          // 3
		"Commodore",        // 4
		"Rear Admiral",     // 5
		"Vice Admiral",     // 6
		"Admiral",          // 7
		"Fleet Admiral",    // 8
		"Grand Admiral",    // 9
		"Supreme Admiral",  // 10
	}
	
	if rank < len(titles) {
		return titles[rank]
	}
	
	// Beyond rank 10
	return "Legendary Admiral"
}

// XPProgress represents XP progress towards next rank.
type XPProgress struct {
	CurrentRank     int    `json:"currentRank"`
	RankTitle       string `json:"rankTitle"`
	TotalXP         int    `json:"totalXP"`
	CurrentRankXP   int    `json:"currentRankXP"`   // XP earned in current rank
	XPForNextRank   int    `json:"xpForNextRank"`   // XP needed for next rank
	ProgressPercent float64 `json:"progressPercent"` // 0.0 to 1.0
}

// GetXPProgress calculates detailed XP progression.
func GetXPProgress(treeState *FormationTreeState) XPProgress {
	if treeState == nil {
		return XPProgress{}
	}
	
	rank := CalculateAdmiralRank(treeState.TotalXP)
	
	// Calculate XP at start of current rank
	xpAtRankStart := 0
	for r := 0; r < rank; r++ {
		xpAtRankStart += GetXPRequiredForNextRank(r)
	}
	
	currentRankXP := treeState.TotalXP - xpAtRankStart
	xpForNext := GetXPRequiredForNextRank(rank)
	
	progress := 0.0
	if xpForNext > 0 {
		progress = float64(currentRankXP) / float64(xpForNext)
	}
	
	return XPProgress{
		CurrentRank:     rank,
		RankTitle:       GetRankTitle(rank),
		TotalXP:         treeState.TotalXP,
		CurrentRankXP:   currentRankXP,
		XPForNextRank:   xpForNext,
		ProgressPercent: progress,
	}
}
