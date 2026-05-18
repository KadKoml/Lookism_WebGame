package usecase

import (
	"math"
)

type ScaledEnemy struct {
	Level int32
	HP    int32
	STR   int32
}

func CalculateAverageSquadLevel(squadLevels []int32) int32 {
	if len(squadLevels) == 0 {
		return 1
	}
	var total int32
	for _, l := range squadLevels {
		total += l
	}
	return int32(math.Round(float64(total) / float64(len(squadLevels))))
}

// ScaleEnemy dynamically scales an enemy's base stats to match the target level
// using the same +3% per level formula used by players.
func ScaleEnemy(baseHp, baseStr int32, targetLevel int32) ScaledEnemy {
	if targetLevel < 1 {
		targetLevel = 1
	}
	levelMultiplier := 1.0 + (float64(targetLevel-1) * 0.03)

	return ScaledEnemy{
		Level: targetLevel,
		HP:    int32(math.Round(float64(baseHp) * levelMultiplier)),
		STR:   int32(math.Round(float64(baseStr) * levelMultiplier)),
	}
}
