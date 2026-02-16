package service

import (
    "TicTacToe/internal/domain/model"
    "errors"
)

type DefaultGameService struct {
	service DomainInterface	
}

func NewDefaultGameService(service DomainInterface) *DefaultGameService {
	return &DefaultGameService{service: service}
}

func (s *DefaultGameService) NextMove(game *model.Game) (int, int) {
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

    return bestMoveX, bestMoveY
}

func (s *DefaultGameService) Validate(oldGame, newGame *model.Game) error {
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

func (s *DefaultGameService) IsGameOver(game *model.Game) bool {
	return game.IsFinished
}