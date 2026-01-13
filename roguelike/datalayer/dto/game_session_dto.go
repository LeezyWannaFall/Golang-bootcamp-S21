package dto

type GameSessionDTO struct {
	LevelNumber int              `json:"level_number"`
	IsRunning   bool             `json:"is_running"`
	Player      PlayerDTO        `json:"player"`
	Level       LevelDTO         `json:"level"`
	Statistics  SessionStatisticsDTO  `json:"statistics"`
}
