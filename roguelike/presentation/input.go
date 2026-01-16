package presentation

import (
	"roguelike/domain/entity"
	"roguelike/domain/game"

	"github.com/rthornton128/goncurses"
)

func ProcessInput(stdscr *goncurses.Window, gs *game.GameSession) bool {
	key := GetInput(stdscr)

	maxY, maxX := stdscr.MaxYX()
	mapHeight := MapHeight
	mapWidth := MapWidth
	shiftX := (maxX - mapWidth) / 2
	shiftY := (maxY - mapHeight) / 2

	clearLine := ""
	for i := 0; i < mapWidth && i < 100; i++ {
		clearLine += " "
	}
	if shiftY-1 >= 0 {
		stdscr.Move(shiftY-1, shiftX)
		stdscr.Print(clearLine)
	}
	if shiftY-2 >= 0 {
		stdscr.Move(shiftY-2, shiftX)
		stdscr.Print(clearLine)
	}

	quit := false

	switch key {
	case 'w', 'W':
		gs.ProcessPlayerTurn(entity.Forward)
		gs.ProcessMonstersTurn()
		stdscr.Refresh()
	case 'a', 'A':
		gs.ProcessPlayerTurn(entity.Left)
		gs.ProcessMonstersTurn()
		stdscr.Refresh()
	case 's', 'S':
		gs.ProcessPlayerTurn(entity.Back)
		gs.ProcessMonstersTurn()
		stdscr.Refresh()
	case 'd', 'D':
		gs.ProcessPlayerTurn(entity.Right)
		gs.ProcessMonstersTurn()
		stdscr.Refresh()
	case 'h', 'H':
		WeaponScreen(stdscr, gs)
		stdscr.Refresh()
	case 'j', 'J':
		FoodScreen(stdscr, gs)
		stdscr.Refresh()
	case 'k', 'K':
		ElixirScreen(stdscr, gs)
		stdscr.Refresh()
	case 'e', 'E':
		ScrollScreen(stdscr, gs)
		stdscr.Refresh()
	case 'i', 'I', 'b', 'B':
		BackpackScreen(stdscr, gs)
		stdscr.Refresh()
	case 27:
		quit = true
	}

	return quit
}
