package repository

import "TicTacToe/domain/model"

type GameRepository interface {
	Save (game* model.Game) error
	Get (id string) (model.Game, error)
}