package usecase

import (
	"math"
)

// CalculateFinalStat applies the merging and leveling formula
// Final_Stat = Base_Stat * (1 + (MergeStars * 0.05)) * (1 + ((Level - 1) * 0.03))
func CalculateFinalStat(baseStat float64, level int32, mergeStars int32) float64 {
	if level < 1 {
		level = 1
	}
	if level > 60 {
		level = 60
	}
	mergeMultiplier := 1.0 + (float64(mergeStars) * 0.05)
	levelMultiplier := 1.0 + (float64(level-1) * 0.03)

	return baseStat * mergeMultiplier * levelMultiplier
}

func CalculateFinalIntStat(baseStat int32, level int32, mergeStars int32) int32 {
	return int32(math.Round(CalculateFinalStat(float64(baseStat), level, mergeStars)))
}
