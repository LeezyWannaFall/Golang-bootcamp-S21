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

func ShowBackpack(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 40) / 2
	shiftY := (maxY - 20) / 2

	stdscr.Move(shiftY-1, shiftX)
	stdscr.Print("=== BACKPACK ===")

	line := 0

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Printf("Gold: %d", gs.Player.Backpack.Treasures.Value)
	line++

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Print("--- Weapons ---")
	line++
	countWeapon := gs.Player.Backpack.WeaponNumber
	if countWeapon > 0 {
		for i := 0; i < countWeapon; i++ {
			weapon := gs.Player.Backpack.Weapons[i]
			stdscr.Move(shiftY+line, shiftX)
			stdscr.Printf("  %s (+%d strength)", weapon.Name, weapon.Strength)
			line++
		}
	} else {
		stdscr.Move(shiftY+line, shiftX)
		stdscr.Print("  (empty)")
		line++
	}

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Print("--- Food ---")
	line++
	countFood := gs.Player.Backpack.FoodNumber
	if countFood > 0 {
		for i := 0; i < countFood; i++ {
			food := gs.Player.Backpack.Foods[i]
			stdscr.Move(shiftY+line, shiftX)
			stdscr.Printf("  %s (+%d health)", food.Name, food.ToRegen)
			line++
		}
	} else {
		stdscr.Move(shiftY+line, shiftX)
		stdscr.Print("  (empty)")
		line++
	}

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Print("--- Elixirs ---")
	line++
	countElixir := gs.Player.Backpack.ElixirNumber
	if countElixir > 0 {
		for i := 0; i < countElixir; i++ {
			elixir := gs.Player.Backpack.Elixirs[i]
			var statName string
			switch elixir.Stat {
			case entity.Health:
				statName = "Health"
			case entity.Agility:
				statName = "Agility"
			case entity.Strength:
				statName = "Strength"
			}
			stdscr.Move(shiftY+line, shiftX)
			stdscr.Printf("  %s (+%d %s, %s)", elixir.Name, elixir.Increase, statName, elixir.Duration.String())
			line++
		}
	} else {
		stdscr.Move(shiftY+line, shiftX)
		stdscr.Print("  (empty)")
		line++
	}

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Print("--- Scrolls ---")
	line++
	countScroll := gs.Player.Backpack.ScrollNumber
	if countScroll > 0 {
		for i := 0; i < countScroll; i++ {
			scroll := gs.Player.Backpack.Scrolls[i]
			var statName string
			switch scroll.Stat {
			case entity.Health:
				statName = "Health"
			case entity.Agility:
				statName = "Agility"
			case entity.Strength:
				statName = "Strength"
			}
			stdscr.Move(shiftY+line, shiftX)
			stdscr.Printf("  %s (+%d %s)", scroll.Name, scroll.Increase, statName)
			line++
		}
	} else {
		stdscr.Move(shiftY+line, shiftX)
		stdscr.Print("  (empty)")
		line++
	}

	stdscr.Move(shiftY+line, shiftX)
	stdscr.Print("--- Keys ---")
	line++
	hasKeys := false
	for i := 0; i < int(entity.KeyColorCount); i++ {
		if gs.Player.Backpack.Keys[i] {
			hasKeys = true
			colorName := "Red"
			switch entity.KeyColor(i) {
			case entity.BlueKey:
				colorName = "Blue"
			case entity.YellowKey:
				colorName = "Yellow"
			case entity.GreenKey:
				colorName = "Green"
			}
			stdscr.Move(shiftY+line, shiftX)
			stdscr.Printf("  %s key", colorName)
			line++
		}
	}
	if !hasKeys {
		stdscr.Move(shiftY+line, shiftX)
		stdscr.Print("  (empty)")
		line++
	}

	stdscr.Move(shiftY+line+1, shiftX)
	stdscr.Print("Press any key to continue...")
	stdscr.Refresh()
	GetInput(stdscr)
}

func ShowWeaponMenu(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Move(shiftY-1, shiftX)
	stdscr.Print("Choose weapon:")

	countWeapon := gs.Player.Backpack.WeaponNumber
	if countWeapon > 0 {
		stdscr.Move(shiftY, shiftX)
		stdscr.Print("0. Without weapon")
	}
	for i := 1; i <= countWeapon; i++ {
		weapon := gs.Player.Backpack.Weapons[i-1]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d strength", i, weapon.Name, weapon.Strength)
	}

	if countWeapon == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't weapon!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
	} else {
		stdscr.Move(shiftY+countWeapon+1, shiftX)
		stdscr.Printf("Press 0-%d key to choose weapon or any key to continue", countWeapon)
		stdscr.Refresh()
		key := GetInput(stdscr)
		if key >= 48 && key <= 48+countWeapon {
			choice := key - 48
			if choice == 0 {
				gs.Player.Weapon = entity.Weapon{Strength: entity.NO_WEAPON, Name: ""}
			} else if choice > 0 && choice <= countWeapon {
				gs.Player.Weapon = gs.Player.Backpack.Weapons[choice-1]
			}
		}
	}

	clearLine := "                              "
	for i := -1; i < 15; i++ {
		if shiftY+i >= 0 {
			stdscr.Move(shiftY+i, shiftX)
			stdscr.Print(clearLine)
		}
	}
}

func ShowFoodMenu(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Move(shiftY, shiftX)
	stdscr.Print("Choose food:")

	countFood := gs.Player.Backpack.FoodNumber
	for i := 1; i <= countFood; i++ {
		food := gs.Player.Backpack.Foods[i-1]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d health", i, food.Name, food.ToRegen)
	}

	if countFood == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't food!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
	} else {
		stdscr.Move(shiftY+countFood+1, shiftX)
		stdscr.Printf("Press 1-%d key to choose food or any key to continue", countFood)
		key := GetInput(stdscr)
		if key >= '1' && key <= '0'+countFood {
			choice := key - '0'
			if choice <= countFood {
				food := gs.Player.Backpack.Foods[choice-1]
				gs.Player.BaseStats.Health += float64(food.ToRegen)
				if gs.Player.BaseStats.Health > float64(gs.Player.RegenLimit) {
					gs.Player.BaseStats.Health = float64(gs.Player.RegenLimit)
				}
				for i := choice - 1; i < countFood-1; i++ {
					gs.Player.Backpack.Foods[i] = gs.Player.Backpack.Foods[i+1]
				}
				gs.Player.Backpack.FoodNumber--
				gs.IncrementFoodConsumed()
			}
		}
	}

	clearLine := "                              "
	for i := 0; i < 15; i++ {
		if shiftY+i >= 0 {
			stdscr.Move(shiftY+i, shiftX)
			stdscr.Print(clearLine)
		}
	}
}

func ShowElixirMenu(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Move(shiftY, shiftX)
	stdscr.Print("Choose elixir:")

	countElixir := gs.Player.Backpack.ElixirNumber
	statNames := []string{"health", "agility", "strength"}
	for i := 1; i <= countElixir; i++ {
		elixir := gs.Player.Backpack.Elixirs[i-1]
		statName := statNames[elixir.Stat]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d %s for %d seconds", i, elixir.Name, elixir.Increase, statName, int(elixir.Duration.Seconds()))
	}

	if countElixir == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't elixir!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
	} else {
		stdscr.Move(shiftY+countElixir+1, shiftX)
		stdscr.Printf("Press 1-%d key to choose elixir or any key to continue", countElixir)
		key := GetInput(stdscr)
		if key >= '1' && key <= '0'+countElixir {
			choice := key - '0'
			if choice <= countElixir {
				for i := choice - 1; i < countElixir-1; i++ {
					gs.Player.Backpack.Elixirs[i] = gs.Player.Backpack.Elixirs[i+1]
				}
				gs.Player.Backpack.ElixirNumber--
				gs.IncrementElixirsDrunk()
			}
		}
	}

	clearLine := "                              "
	for i := 0; i < 15; i++ {
		if shiftY+i >= 0 {
			stdscr.Move(shiftY+i, shiftX)
			stdscr.Print(clearLine)
		}
	}
}

func ShowScrollMenu(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Move(shiftY, shiftX)
	stdscr.Print("Choose scroll:")

	countScroll := gs.Player.Backpack.ScrollNumber
	statNames := []string{"health", "agility", "strength"}
	for i := 1; i <= countScroll; i++ {
		scroll := gs.Player.Backpack.Scrolls[i-1]
		statName := statNames[scroll.Stat]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d %s", i, scroll.Name, scroll.Increase, statName)
	}

	if countScroll == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't scroll!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
	} else {
		stdscr.Move(shiftY+countScroll+1, shiftX)
		stdscr.Printf("Press 1-%d key to choose scroll or any key to continue", countScroll)
		key := GetInput(stdscr)
		if key >= '1' && key <= '0'+countScroll {
			choice := key - '0'
			if choice <= countScroll {
				scroll := gs.Player.Backpack.Scrolls[choice-1]
				switch scroll.Stat {
				case entity.Health:
					gs.Player.RegenLimit += scroll.Increase
					gs.Player.BaseStats.Health += float64(scroll.Increase)
				case entity.Agility:
					gs.Player.BaseStats.Agility += scroll.Increase
				case entity.Strength:
					gs.Player.BaseStats.Strength += scroll.Increase
				}
				for i := choice - 1; i < countScroll-1; i++ {
					gs.Player.Backpack.Scrolls[i] = gs.Player.Backpack.Scrolls[i+1]
				}
				gs.Player.Backpack.ScrollNumber--
				gs.IncrementScrollsRead()
			}
		}
	}

	clearLine := "                              "
	for i := 0; i < 15; i++ {
		if shiftY+i >= 0 {
			stdscr.Move(shiftY+i, shiftX)
			stdscr.Print(clearLine)
		}
	}
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
