package web

import (
	"TicTacToe/internal/domain/model"
)

func FromRequest(dto RequestDTO) *model.Game {
	return &model.Game{
		Field: model.GameField{
			Cells: dto.GameField,
		},
	} 
}

func ToResponse(game *model.Game) ResponseDTO {
	return ResponseDTO{
		GameID: game.ID,
		GameField: game.Field.Cells,
		Winner: game.Winner,
	}
}