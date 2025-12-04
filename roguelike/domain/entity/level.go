package entity

type Level struct {
	Coordinates Object
	Rooms Room
	Passages Passages
	LevelNumber int
	EndOfLevel Object
}