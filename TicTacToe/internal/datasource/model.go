package datasource

import "github.com/google/uuid"

type GameDTO struct {
	ID          uuid.UUID
	Field       *FieldDTO
	CurrentTurn int
	Winner      int
}

type FieldDTO struct {
	Cells [3][3]int
}