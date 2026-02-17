package service

import (
	"TicTacToe/internal/domain/model"
	"errors"

	"github.com/google/uuid"
)

type GameService struct {
	repo DataInterface
}

func NewGameService(repo DataInterface) DomainInterface {
    return &GameService{repo: repo}
}

/*
	NextMove(game *model.Game) (int, int)
	Validate(game *model.Game) error
	IsGameOver(game *model.Game) bool
*/

func (s *GameService) NextMove(game *model.Game) bool {
	dataGame, err := s.repo.Get(game.ID.String())
	if err != nil {
		return false
	}
	
	err = s.Validate(dataGame, game)
	if err != nil {
		return false
	}

    if game.Field.CheckWin(model.Cross) {
        game.Winner = model.Cross
        game.IsFinished = true
        s.repo.Save(game)
        return true
    }

    if game.Field.CheckAllCellsFilled() {
        game.Winner = model.Empty
        game.IsFinished = true
        s.repo.Save(game)
        return true
    }

    bestScore := 2
    bestMoveX, bestMoveY := -1, -1

    for i := 0; i < model.FieldSize; i++ {
        for j := 0; j < model.FieldSize; j++ {
            if game.Field.Cells[i][j] == model.Empty {
                game.Field.Cells[i][j] = model.Zero
                
                score := MiniMax(game.Field, model.Cross)
                
                game.Field.Cells[i][j] = model.Empty

                if score < bestScore { 
                    bestScore = score
                    bestMoveX = j
                    bestMoveY = i
                }
            }
        }
    }
    
    if bestMoveX != -1 && bestMoveY != -1 {
        game.Field.PlaceSymbolOnField(bestMoveX, bestMoveY, model.Zero)
        
        if game.Field.CheckWin(model.Zero) {
            game.Winner = model.Zero
            game.IsFinished = true
        } else if game.Field.CheckAllCellsFilled() {
            game.Winner = model.Empty
            game.IsFinished = true
        } else {
            game.CurrentTurn = model.Cross
        }
    }

    err = s.repo.Save(game)
    return err == nil
}

func (s *GameService) Validate(oldGame, newGame *model.Game) error {
	if s.IsGameOver(oldGame) {
        return errors.New("game already finished") 
    }

    changes := 0

    for i := 0; i < model.FieldSize; i++ {
        for j := 0; j < model.FieldSize; j++ {
            oldValue := oldGame.Field.Cells[i][j]
            newValue := newGame.Field.Cells[i][j]

            if oldValue != newValue {
                changes += 1

                if oldValue != model.Empty {
                    return errors.New("previous move was changed")
                }

                if newValue == model.Empty {
                    return errors.New("previous move was changed")
                }

                if newValue != oldGame.CurrentTurn {
                    return errors.New("invlaid symbol placed")
                }
            }
        }
    }

    if changes == 0 {
        return errors.New("no move made")
    }

    if changes > 1 {
        return errors.New("more than 1 move made")
    }

    return nil
}

func (s *GameService) IsGameOver(game *model.Game) bool {
	return game.IsFinished
}

func (s *GameService) StartGame() (*model.Game, error) {
    newGame := &model.Game{
        ID: uuid.New(),
        IsFinished: false,
        Field: model.GameField{},
        CurrentTurn: model.Cross,
        Winner: model.Empty,
    }
    newGame.Field.ClearField()

    err := s.repo.Save(newGame)
    if err != nil {
        return nil, err
    }

    return newGame, nil
}