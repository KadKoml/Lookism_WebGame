package usecase

import (
	"math/rand"
)

type CharacterStats struct {
	HP         int32
	STR        int32
	Durability int32
	Power      float64
	Agility    int32
	Speed      float64
	Reaction   int32
	Technique  float64
}

func CalculateBaseAttack(attacker CharacterStats) int32 {
	return attacker.STR
}

func CalculateBlockEffectiveness(defender CharacterStats, attackerPower float64, attackerSTR int32) int32 {
	// Block = (HP * (Durability / 100) / Enemy Power) - Enemy Strength
	// We use float64 for intermediate math to avoid integer division truncation
	durabilityPercent := float64(defender.Durability) / 100.0
	blockAmount := (float64(defender.HP) * durabilityPercent / attackerPower) - float64(attackerSTR)
	if blockAmount < 0 {
		return 0
	}
	return int32(blockAmount)
}

func CalculateDodgeChance(defenderAgility int32, attackerSpeed float64) float64 {
	// Dodge Chance = Agility / Enemy Speed
	return float64(defenderAgility) / attackerSpeed
}

func CalculateCounterChance(defenderReaction int32, attackerTechnique float64) float64 {
	// Counter Chance = Reaction / Enemy Technique
	return float64(defenderReaction) / attackerTechnique
}

func CalculateActiveAbilityDamage(attackerSTR int32, abilityMultiplier float64, attackerTechnique float64) int32 {
	// Active Ability Dmg = Strength * Ability_Multiplier * Technique
	return int32(float64(attackerSTR) * abilityMultiplier * attackerTechnique)
}

func IsActionSuccessful(chancePercentage float64) bool {
	// chancePercentage could be something like 0.45 for 45%
	return rand.Float64() <= chancePercentage
}
