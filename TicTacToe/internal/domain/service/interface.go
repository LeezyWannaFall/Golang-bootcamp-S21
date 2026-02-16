package service

import (
	"TicTacToe/internal/domain/model"
)

type DataInterface interface {
	Save(game *model.Game) error
	Get(id string) (*model.Game, error)
}

type DomainInterface interface {
	NextMove(game *model.Game) (int, int)
	Validate(oldGame, newGame *model.Game) error
	IsGameOver(game *model.Game) bool
}