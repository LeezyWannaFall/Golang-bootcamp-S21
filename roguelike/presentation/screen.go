package presentation

import (
	"roguelike/domain/game"
	"roguelike/domain/characters"
	"github.com/rthornton128/goncurses"
	"roguelike/domain/entity"
)

func StartScreen(stdscr *goncurses.Window) {
	stdscr.Clear()
	maxY, maxX := stdscr.MaxYX()

	title := "ROGUELIKE GAME"
	titleX := (maxX - len(title)) / 2
	titleY := maxY/2 - 2

	stdscr.Move(titleY, titleX)
	stdscr.Print(title)

	pressKey := "Press any key to start..."
	pressKeyX := (maxX - len(pressKey)) / 2
	pressKeyY := titleY + 2

	stdscr.Move(pressKeyY, pressKeyX)
	stdscr.Print(pressKey)

	stdscr.Refresh()
	stdscr.GetChar()
}

func MenuScreen(stdscr *goncurses.Window, currentLine int, hasSave bool) int {
	stdscr.Clear()
	maxY, maxX := stdscr.MaxYX()

	menuItems := []string{
		"New Game",
		"Continue Game",
		"Scoreboard",
		"Exit",
	}

	if !hasSave {
		menuItems = []string{
			"New Game",
			"Scoreboard",
			"Exit",
		}
		if currentLine >= len(menuItems) {
			currentLine = 0
		}
		if currentLine == 1 {
			currentLine = 0
		}
	}

	menuY := maxY/2 - len(menuItems)/2
	menuX := maxX/2 - 10

	for i, item := range menuItems {
		y := menuY + i
		stdscr.Move(y, menuX)
		if i == currentLine {
			stdscr.AttrOn(goncurses.A_REVERSE)
		}
		stdscr.Printf("  %s  ", item)
		if i == currentLine {
			stdscr.AttrOff(goncurses.A_REVERSE)
		}
	}

	instructionY := menuY + len(menuItems) + 2
	instructionX := maxX/2 - 15
	stdscr.Move(instructionY, instructionX)
	stdscr.Print("Use W/S to navigate, ENTER to select")

	stdscr.Refresh()

	for {
		key := stdscr.GetChar()
		switch key {
		case 'w', 'W':
			currentLine--
			if currentLine < 0 {
				currentLine = len(menuItems) - 1
			}
			return MenuScreen(stdscr, currentLine, hasSave)
		case 's', 'S':
			currentLine++
			if currentLine >= len(menuItems) {
				currentLine = 0
			}
			return MenuScreen(stdscr, currentLine, hasSave)
		case '\n', '\r':
			return currentLine
		case 27:
			return len(menuItems) - 1
		}
	}
}

func DeadScreen(stdscr *goncurses.Window) {
	stdscr.Clear()
	maxY, maxX := stdscr.MaxYX()

	message := "YOU DIED"
	messageX := (maxX - len(message)) / 2
	messageY := maxY / 2

	stdscr.Move(messageY, messageX)
	stdscr.Print(message)

	pressKey := "Press any key to continue..."
	pressKeyX := (maxX - len(pressKey)) / 2
	pressKeyY := messageY + 2

	stdscr.Move(pressKeyY, pressKeyX)
	stdscr.Print(pressKey)

	stdscr.Refresh()
	stdscr.GetChar()
}

func EndgameScreen(stdscr *goncurses.Window) {
	stdscr.Clear()
	maxY, maxX := stdscr.MaxYX()

	message := "CONGRATULATIONS! YOU COMPLETED THE GAME!"
	messageX := (maxX - len(message)) / 2
	messageY := maxY / 2

	stdscr.Move(messageY, messageX)
	stdscr.Print(message)

	pressKey := "Press any key to continue..."
	pressKeyX := (maxX - len(pressKey)) / 2
	pressKeyY := messageY + 2

	stdscr.Move(pressKeyY, pressKeyX)
	stdscr.Print(pressKey)

	stdscr.Refresh()
	stdscr.GetChar()
}

func BackpackScreen(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 40) / 2
	shiftY := (maxY - 20) / 2

	stdscr.Clear()
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
			case 0:
				statName = "Health"
			case 1:
				statName = "Agility"
			case 2:
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
			case 0:
				statName = "Health"
			case 1:
				statName = "Agility"
			case 2:
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
	for i := 0; i < 4; i++ {
		if gs.Player.Backpack.Keys[i] {
			hasKeys = true
			colorName := "Red"
			switch i {
			case 1:
				colorName = "Blue"
			case 2:
				colorName = "Yellow"
			case 3:
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

func FoodScreen(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Clear()
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
		return
	}

	stdscr.Move(shiftY+countFood+1, shiftX)
	stdscr.Printf("Press 1-%d key to eat or any key to continue", countFood)
	stdscr.Refresh()
	key := GetInput(stdscr)
	if key >= '1' && key <= '0'+countFood {
		idx := key - '1'
		if idx >= 0 && int(idx) < countFood {
			food := gs.Player.Backpack.Foods[idx]
			// Применить еду
			gs.Player.BaseStats.Health += float64(food.ToRegen)
			if gs.Player.BaseStats.Health > float64(gs.Player.RegenLimit) {
				gs.Player.BaseStats.Health = float64(gs.Player.RegenLimit)
			}
			// Удалить из рюкзака
			for j := int(idx); j < countFood-1; j++ {
				gs.Player.Backpack.Foods[j] = gs.Player.Backpack.Foods[j+1]
			}
			gs.Player.Backpack.FoodNumber--
		}
	}
}

func ElixirScreen(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Clear()
	stdscr.Move(shiftY, shiftX)
	stdscr.Print("Choose elixir:")

	countElixir := gs.Player.Backpack.ElixirNumber
	statNames := []string{"health", "agility", "strength"}
	for i := 1; i <= countElixir; i++ {
		elixir := gs.Player.Backpack.Elixirs[i-1]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d %s", i, elixir.Name, elixir.Increase, statNames[elixir.Stat])
	}

	if countElixir == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't elixirs!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
		return
	}

	stdscr.Move(shiftY+countElixir+1, shiftX)
	stdscr.Printf("Press 1-%d key to drink or any key to continue", countElixir)
	stdscr.Refresh()
	key := GetInput(stdscr)
	if key >= '1' && key <= '0'+countElixir {
		idx := key - '1'
		if idx >= 0 && int(idx) < countElixir {
			// Применить эликсир
			currentRoom := gs.GetCurrentRoom()
			characters.UseConsumable(gs.Player, characters.ElixirType, currentRoom, int(idx))
		}
	}
}

func WeaponScreen(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Clear()
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
		return
	}

	stdscr.Move(shiftY+countWeapon+1, shiftX)
	stdscr.Printf("Press 0-%d key to choose weapon or any key to continue", countWeapon)
	stdscr.Refresh()
	key := GetInput(stdscr)
	if key >= 48 && key <= 48+countWeapon {
		choice := key - 48
		if choice == 0 {
			gs.Player.Weapon = entity.Weapon{Strength: entity.NO_WEAPON, Name: ""}
		} else if choice > 0 && int(choice) <= countWeapon {
			gs.Player.Weapon = gs.Player.Backpack.Weapons[choice-1]
		}
	}
}

func ScrollScreen(stdscr *goncurses.Window, gs *game.GameSession) {
	maxY, maxX := stdscr.MaxYX()
	shiftX := (maxX - 30) / 2
	shiftY := (maxY - 10) / 2

	stdscr.Clear()
	stdscr.Move(shiftY, shiftX)
	stdscr.Print("Choose scroll:")

	countScroll := gs.Player.Backpack.ScrollNumber
	statNames := []string{"health", "agility", "strength"}
	for i := 1; i <= countScroll; i++ {
		scroll := gs.Player.Backpack.Scrolls[i-1]
		stdscr.Move(shiftY+i, shiftX)
		stdscr.Printf("%d. %s +%d %s", i, scroll.Name, scroll.Increase, statNames[scroll.Stat])
	}

	if countScroll == 0 {
		stdscr.Move(shiftY+1, shiftX)
		stdscr.Print("You haven't scrolls!")
		stdscr.Move(shiftY+2, shiftX)
		stdscr.Print("Press any key to continue...")
		stdscr.GetChar()
		return
	}

	stdscr.Move(shiftY+countScroll+1, shiftX)
	stdscr.Printf("Press 1-%d key to read or any key to continue", countScroll)
	stdscr.Refresh()
	key := GetInput(stdscr)
	if key >= '1' && key <= '0'+countScroll {
		idx := key - '1'
		if idx >= 0 && int(idx) < countScroll {
			// Применить свиток
			currentRoom := gs.GetCurrentRoom()
			characters.UseConsumable(gs.Player, characters.ScrollType, currentRoom, int(idx))
		}
	}
}
