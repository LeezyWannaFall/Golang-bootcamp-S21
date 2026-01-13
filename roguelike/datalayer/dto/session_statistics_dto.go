package dto

type SessionStatisticsDTO struct {
	TreasuresCollected int `json:"treasures_collected"`
	DeepestLevel       int `json:"deepest_level"`
	EnemiesDefeated    int `json:"enemies_defeated"`
	FoodConsumed       int `json:"food_consumed"`
	ElixirsDrunk       int `json:"elixirs_drunk"`
	ScrollsRead        int `json:"scrolls_read"`
	AttacksDealt       int `json:"attacks_dealt"`
	AttacksMissed      int `json:"attacks_missed"`
	TilesTraveled      int `json:"tiles_traveled"`
}
