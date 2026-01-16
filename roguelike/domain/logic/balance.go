package logic

const (
	MAX_BALANCE_ADJUSTMENT = 50
	MIN_BALANCE_ADJUSTMENT = -50
)

type BalanceAdjustment struct {
	MonsterDifficulty int
	MonsterCount      int
	ConsumableCount   int
	FoodBonus         int
}

type SessionStatistics struct {
	TreasuresCollected int
	DeepestLevel       int
	EnemiesDefeated    int
	FoodConsumed       int
	ElixirsDrunk       int
	ScrollsRead        int
	AttacksDealt       int
	AttacksMissed      int
	TilesTraveled      int
}

func CalculateBalanceAdjustment(statistics SessionStatistics, levelNum int) BalanceAdjustment {
	adjustment := BalanceAdjustment{
		MonsterDifficulty: 0,
		MonsterCount:      0,
		ConsumableCount:   0,
		FoodBonus:         0,
	}

	if levelNum <= 1 {
		return adjustment
	}

	avgMovesPerLevel := 0
	if levelNum > 0 {
		avgMovesPerLevel = statistics.TilesTraveled / levelNum
	}

	avgFoodPerLevel := 0
	if levelNum > 0 {
		avgFoodPerLevel = statistics.FoodConsumed / levelNum
	}

	avgElixirsPerLevel := 0
	if levelNum > 0 {
		avgElixirsPerLevel = statistics.ElixirsDrunk / levelNum
	}

	avgEnemiesPerLevel := 0
	if levelNum > 0 {
		avgEnemiesPerLevel = statistics.EnemiesDefeated / levelNum
	}

	missRate := 0
	if statistics.AttacksDealt > 0 {
		missRate = statistics.AttacksMissed * 100 / statistics.AttacksDealt
	}

	if avgMovesPerLevel > 200 {
		adjustment.MonsterDifficulty -= 10
		adjustment.MonsterCount -= 1
		adjustment.ConsumableCount += 1
		adjustment.FoodBonus += 20
	} else if avgMovesPerLevel < 100 {
		adjustment.MonsterDifficulty += 10
		adjustment.MonsterCount += 1
		adjustment.ConsumableCount -= 1
	}

	if avgFoodPerLevel > 2 {
		adjustment.MonsterDifficulty -= 15
		adjustment.FoodBonus += 30
		adjustment.ConsumableCount += 1
	} else if avgFoodPerLevel == 0 && levelNum > 3 {
		adjustment.MonsterDifficulty += 15
		adjustment.ConsumableCount -= 1
	}

	if avgElixirsPerLevel > 1 {
		adjustment.MonsterDifficulty -= 10
		adjustment.ConsumableCount += 1
	} else if avgElixirsPerLevel == 0 && levelNum > 5 {
		adjustment.MonsterDifficulty += 10
	}

	if missRate > 30 {
		adjustment.MonsterDifficulty -= 10
		adjustment.MonsterCount -= 1
	} else if missRate < 10 && statistics.AttacksDealt > 20 {
		adjustment.MonsterDifficulty += 10
		adjustment.MonsterCount += 1
	}

	if avgEnemiesPerLevel > 5 && levelNum > 3 {
		adjustment.MonsterDifficulty += 5
	} else if avgEnemiesPerLevel < 2 && levelNum > 3 {
		adjustment.MonsterDifficulty -= 5
		adjustment.MonsterCount -= 1
	}

	if adjustment.MonsterDifficulty > MAX_BALANCE_ADJUSTMENT {
		adjustment.MonsterDifficulty = MAX_BALANCE_ADJUSTMENT
	}
	if adjustment.MonsterDifficulty < MIN_BALANCE_ADJUSTMENT {
		adjustment.MonsterDifficulty = MIN_BALANCE_ADJUSTMENT
	}

	if adjustment.MonsterCount > 2 {
		adjustment.MonsterCount = 2
	}
	if adjustment.MonsterCount < -2 {
		adjustment.MonsterCount = -2
	}

	if adjustment.ConsumableCount > 2 {
		adjustment.ConsumableCount = 2
	}
	if adjustment.ConsumableCount < -2 {
		adjustment.ConsumableCount = -2
	}

	if adjustment.FoodBonus > 50 {
		adjustment.FoodBonus = 50
	}
	if adjustment.FoodBonus < 0 {
		adjustment.FoodBonus = 0
	}

	return adjustment
}
