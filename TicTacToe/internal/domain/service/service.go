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

    symbol := model.Cross
    next := model.Zero
    bestScore := -2
    bestMoveX := 0
    bestMoveY := 0

    if game.CurrentTurn == model.Zero {
        symbol = model.Zero
        next = model.Cross
        bestScore = 2
    }

    for i := 0; i < model.FieldSize; i++ {
        for j := 0; j < model.FieldSize; j++ {
            if game.Field.Cells[i][j] == model.Empty {
                game.Field.Cells[i][j] = symbol
                
                score := MiniMax(game.Field, next)
                
                game.Field.Cells[i][j] = model.Empty

                if (symbol == model.Cross && score > bestScore) ||
                   (symbol == model.Zero && score < bestScore) {
                    bestScore = score
                    bestMoveX = j
                    bestMoveY = i
                }
            }
        }
    }
	
    symbolToPlace := game.CurrentTurn
    
    success := game.Field.PlaceSymbolOnField(bestMoveX, bestMoveY, symbolToPlace)
    if !success {
        return false
    }

    if game.Field.CheckWin(symbolToPlace) {
        game.Winner = symbolToPlace
        game.IsFinished = true
    } else if game.Field.CheckAllCellsFilled() {
        game.Winner = model.Empty
        game.IsFinished = true
    } else {
        game.SwitchTurn()
    }

    err = s.repo.Save(game)
    if err != nil {
        return false
    }

	return true
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