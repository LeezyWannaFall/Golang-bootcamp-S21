package service

import (
	"TicTacToe/internal/domain/model"
)

const (
	AiZeroWin = -1
	PlayerCrossWin = 1
	Draw = 0
)

func MiniMax(field model.GameField, currentTurn int) int {
	if field.CheckWin(model.Cross) { return PlayerCrossWin }
	if field.CheckWin(model.Zero) { return AiZeroWin }
	if field.CheckAllCellsFilled() { return Draw }	

	if currentTurn == model.Cross {
		bestScore := -2

		for i := 0; i < model.FieldSize; i++ {
			for j := 0; j < model.FieldSize; j++ {

				if field.Cells[i][j] == model.Empty {
					field.Cells[i][j] = model.Cross
					score := MiniMax(field, model.Zero)
					field.Cells[i][j] = model.Empty

					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		
		return bestScore
	} else {
		bestScore := 2

		for i := 0; i < model.FieldSize; i++ {
			for j := 0; j < model.FieldSize; j++ {

				if field.Cells[i][j] == model.Empty {
					field.Cells[i][j] = model.Zero
					score := MiniMax(field, model.Cross)
					field.Cells[i][j] = model.Empty

					if score < bestScore {
						bestScore = score
					}
				}
			}
		}

		return bestScore
	}
}