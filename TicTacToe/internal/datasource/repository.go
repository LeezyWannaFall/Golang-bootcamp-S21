package datasource

import (
	"TicTacToe/internal/domain/model"
	"errors"
)

type GameRepository struct {
	storage *GameStorage
}

func NewRepository(storage *GameStorage) *GameRepository {
	return &GameRepository{storage: storage}
}

func (r *GameRepository) Save(game *model.Game) error {
	DTOgame := ToDTO(game)
	r.storage.data.Store(DTOgame.ID.String(), DTOgame)
	
	return nil
}

func (r *GameRepository) Get(id string) (*model.Game, error) {
	value, ok := r.storage.data.Load(id)
	if !ok {
		return nil, errors.New("game not found")
	}

	DTOgame, ok := value.(*GameDTO)
	if !ok {
		return nil, errors.New("invalid game data")
	}
	return FromDTO(DTOgame), nil
}