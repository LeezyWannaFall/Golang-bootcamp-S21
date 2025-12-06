package entity

type Level struct {
	Coordinates Object
	Rooms [ROOMS_NUM]Room
	Passages Passages
	LevelNumber int
	EndOfLevel Object
}