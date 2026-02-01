package service

import (
	"TicTacToe/domain/model"
)

type GameService interface {
	NextMove(game *model.Game) (int, int)
	Validate(game *model.Game) error
	IsGameOver(game *model.Game) bool
}