package model

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

// func NewGame() *Game {
// 	return &Game{
// 		ID:          uuid.New(),
// 		IsFinished:  false,
// 		CurrentTurn: Cross,
// 		Winner:      0,
// 	}
// }


// func (g *Game) SwitchTurn() {
// 	if g.CurrentTurn == Cross {
// 		g.CurrentTurn = Zero
// 	} else {
// 		g.CurrentTurn = Cross
// 	}
// }

// func (g* Game) MakeMove(x, y int) bool {
// 	if g.IsFinished {
// 		return false
// 	}

// 	if !g.Field.PlaceSymbolOnField(x, y, g.CurrentTurn) {
// 		return false
// 	}

// 	if g.Field.CheckWin(g.CurrentTurn) {
// 		g.Winner = g.CurrentTurn
// 		g.IsFinished = true
// 		return false
// 	}

// 	g.SwitchTurn()
// 	return true
// }