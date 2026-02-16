package datasource

import (
	"TicTacToe/internal/domain/model"
)

func ToDTO(game *model.Game) *GameDTO {
	return &GameDTO{
		ID:          game.ID,
		Field:       ToFieldDTO(&game.Field),
		CurrentTurn: game.CurrentTurn,
		Winner:      game.Winner,
	}
}

func FromDTO(dto *GameDTO) *model.Game {
	return &model.Game{
		ID:          dto.ID,
		Field:       *FromFieldDTO(dto.Field),
		CurrentTurn: dto.CurrentTurn,
		Winner:      dto.Winner,
	}
}

func ToFieldDTO(field *model.GameField) *FieldDTO {
	return &FieldDTO{
		Cells: field.Cells,
	}
}

func FromFieldDTO(dto *FieldDTO) *model.GameField {
	return &model.GameField{
		Cells: dto.Cells,
	}
}