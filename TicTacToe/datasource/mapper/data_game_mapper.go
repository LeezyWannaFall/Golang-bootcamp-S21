package mapper

import (
	"TicTacToe/datasource/repository"
	"TicTacToe/domain/model"
)

type GameModule struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) *GameModule {
	return &GameModule{repo: repo}
}

func (g* GameModule) SaveRun(game model.Game) error {

}

func (g* GameModule) GetRun() (model.Game, error) {

}