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

func (g *Game) SwitchTurn() {
	if g.CurrentTurn == Cross {
		g.CurrentTurn = Zero
	} else {
		g.CurrentTurn = Cross
	}
}