package presentation

import (
	"roguelike/domain/entity"
	"roguelike/domain/game"

	"github.com/rthornton128/goncurses"
)

func DisplayMap(stdscr *goncurses.Window, renderer *Renderer, level *entity.Level, player *entity.Player) {
	renderer.CreateNewMap(level, player)
	mapData := renderer.GetMap()

	maxY, maxX := stdscr.MaxYX()
	mapHeight := len(mapData)
	mapWidth := 0
	if mapHeight > 0 {
		mapWidth = len(mapData[0])
	}

	shiftX := (maxX - mapWidth) / 2
	shiftY := (maxY - mapHeight) / 2

	clearLine := ""
	for i := 0; i < mapWidth && i < 100; i++ {
		clearLine += " "
	}
	if shiftY-2 >= 0 {
		stdscr.Move(shiftY-2, shiftX)
		stdscr.Print(clearLine)
	}
	if shiftY-1 >= 0 {
		stdscr.Move(shiftY-1, shiftX)
		stdscr.Print(clearLine)
	}

	colorMap := renderer.GetColorMap()
	for i := 0; i < mapHeight; i++ {
		stdscr.Move(shiftY+i, shiftX)
		for j := 0; j < mapWidth; j++ {
			colorPair := int16(colorMap[i][j])
			stdscr.AttrOn(goncurses.ColorPair(colorPair))
			stdscr.AddChar(goncurses.Char(mapData[i][j]))
			stdscr.AttrOff(goncurses.ColorPair(colorPair))
		}
	}

	statsY := shiftY + mapHeight
	if statsY >= 0 && statsY < maxY {
		stdscr.Move(statsY, shiftX)
		for i := 0; i < mapWidth && i < maxX-shiftX; i++ {
			stdscr.AddChar(' ')
		}
		stdscr.Move(statsY, shiftX)
		weaponStrength := 0
		if player.Weapon.Strength != entity.NO_WEAPON {
			weaponStrength = player.Weapon.Strength
		}
		stdscr.Printf("Level: %-8d ", level.LevelNumber)
		stdscr.Printf("Gold: %-8d ", player.Backpack.Treasures.Value)
		stdscr.Printf("Health: %.2f/%-8d ", player.BaseStats.Health, player.RegenLimit)
		stdscr.Printf("Agility: %-6d ", player.BaseStats.Agility)
		stdscr.Printf("Strength: %d(+%d) ", player.BaseStats.Strength, weaponStrength)
	}

	stdscr.Move(maxY-1, maxX-1)
	stdscr.Refresh()
}

func DisplayScoreboard(stdscr *goncurses.Window, runs []game.RunResult) {
	maxY, maxX := stdscr.MaxYX()
	stdscr.Clear()

	fieldSize := 10
	width := fieldSize * 9
	height := len(runs)*2 + 3
	if height > maxY-4 {
		height = maxY - 4
	}

	shiftX := (maxX - width) / 2
	shiftY := (maxY - height) / 2

	stdscr.Move(shiftY-2, shiftX)
	for i := 0; i < width; i++ {
		stdscr.AddChar('-')
	}

	stdscr.Move(shiftY-1, shiftX)
	stdscr.Printf("|%-*s", fieldSize, "treasures")
	stdscr.Printf("|%-*s", fieldSize, "level")
	stdscr.Printf("|%-*s", fieldSize, "enemies")
	stdscr.Printf("|%-*s", fieldSize, "food")
	stdscr.Printf("|%-*s", fieldSize, "elixirs")
	stdscr.Printf("|%-*s", fieldSize, "scrolls")
	stdscr.Printf("|%-*s", fieldSize, "attacks")
	stdscr.Printf("|%-*s", fieldSize, "missed")
	stdscr.Printf("|%-*s", fieldSize, "moves")
	stdscr.AddChar('|')
	stdscr.AddChar('\n')

	for i, run := range runs {
		if i >= (height-3)/2 {
			break
		}
		stdscr.Move(shiftY+2*i, shiftX)
		for j := 0; j < width; j++ {
			stdscr.AddChar('-')
		}

		stdscr.Move(shiftY+2*i+1, shiftX)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.TreasuresCollected)
		stdscr.Printf("|%*d", fieldSize, run.FinalLevel)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.EnemiesDefeated)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.FoodConsumed)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.ElixirsDrunk)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.ScrollsRead)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.AttacksDealt)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.AttacksMissed)
		stdscr.Printf("|%*d", fieldSize, run.Statistics.TilesTraveled)
		stdscr.AddChar('|')
		stdscr.AddChar('\n')
	}

	stdscr.Move(shiftY+2*len(runs), shiftX)
	for i := 0; i < width; i++ {
		stdscr.AddChar('-')
	}

	stdscr.Move(shiftY+2*(len(runs)+1), (maxX-20)/2)
	stdscr.Print("Press ESCAPE to exit.")

	stdscr.Refresh()

	for {
		key := GetInput(stdscr)
		if key == 27 {
			break
		}
	}
}
