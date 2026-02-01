package domain

import "github.com/google/uuid"

const (
	Empty = 0
	Cross = 1
	Zero = 2
)

type Game struct {
	ID uuid.UUID
	IsFinished bool
	Field GameField
	CurrentTurn int
	Winner int
}

func NewGame() *Game {
	return &Game{
		ID:          uuid.New(),
		IsFinished:  false,
		CurrentTurn: Cross,
		Winner:      0,
	}
}


func (g *Game) SwitchTurn() {
	if g.CurrentTurn == Cross {
		g.CurrentTurn = Zero
	} else {
		g.CurrentTurn = Cross
	}
}

func (g* Game) MakeMove(x, y int) bool {
	if g.IsFinished {
		return false
	}

	if !g.Field.PlaceSymbolOnField(x, y, g.CurrentTurn) {
		return false
	}

	if g.CheckWin(g.CurrentTurn) {
		g.Winner = g.CurrentTurn
		g.IsFinished = true
		return false
	}

	g.SwitchTurn()
	return true
}

func (g *Game) CheckWin(symbol int) bool {
	// строки
	for i := 0; i < FieldSize; i++ {
		if g.Field.Cells[i][0] == symbol &&
			g.Field.Cells[i][1] == symbol &&
			g.Field.Cells[i][2] == symbol {
			return true
		}
	}

	// колонки
	for j := 0; j < FieldSize; j++ {
		if g.Field.Cells[0][j] == symbol &&
			g.Field.Cells[1][j] == symbol &&
			g.Field.Cells[2][j] == symbol {
			return true
		}
	}

	// диагонали
	if g.Field.Cells[0][0] == symbol &&
		g.Field.Cells[1][1] == symbol &&
		g.Field.Cells[2][2] == symbol {
		return true
	}

	if g.Field.Cells[0][2] == symbol &&
		g.Field.Cells[1][1] == symbol &&
		g.Field.Cells[2][0] == symbol {
		return true
	}

	return false
}