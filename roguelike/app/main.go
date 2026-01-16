package main

import (
	"fmt"
	"os"
	"path/filepath"
	"roguelike/datalayer/jsonrepo"
	"roguelike/domain/entity"
	"roguelike/domain/game"
	"roguelike/presentation"
	"time"

	"github.com/rthornton128/goncurses"
)

func main() {
	saveRepo := jsonrepo.NewJSONGameSaveRepository("./save.json")
	statRepo := jsonrepo.NewJSONStatisticRepository("./stats.json")

	stdscr, err := presentation.InitCurses()
	if err != nil {
		panic(err)
	}
	defer goncurses.End()

	presentation.StartScreen(stdscr)

	topRuns, _ := statRepo.LoadTop(10)

	var menuChoice int
	var hasSave bool
	runningMenu := true
	for runningMenu {
		hasSave, _ = saveRepo.HasSave()
		menuChoice = presentation.MenuScreen(stdscr, 0, hasSave)

		if hasSave {
			switch menuChoice {
			case 0:
				runningMenu = false
			case 1:
				runningMenu = false
			case 2:
				if len(topRuns) > 0 {
					presentation.DisplayScoreboard(stdscr, topRuns)
				}
			case 3:
				return
			}
		} else {
			switch menuChoice {
			case 0:
				runningMenu = false
			case 1:
				if len(topRuns) > 0 {
					presentation.DisplayScoreboard(stdscr, topRuns)
				}
			case 2:
				return
			}
		}
	}

	gs := game.NewGameSession()

	hasSave, _ = saveRepo.HasSave()
	if hasSave && menuChoice == 1 {
		state, err := saveRepo.LoadGame()
		if err != nil {
			stdscr.Clear()
			stdscr.MovePrint(0, 0, "Error loading game: "+err.Error())
			stdscr.MovePrint(1, 0, "Starting new game...")
			stdscr.MovePrint(2, 0, "Press any key to continue...")
			stdscr.Refresh()
			_ = stdscr.GetChar()
			gs.Start()
			gs.InitLevel()
		} else {
			defer func() {
				if r := recover(); r != nil {
					stdscr.Clear()
					stdscr.MovePrint(0, 0, fmt.Sprintf("Error restoring game: %v", r))
					stdscr.MovePrint(1, 0, "Starting new game...")
					stdscr.MovePrint(2, 0, "Press any key to continue...")
					stdscr.Refresh()
					_ = stdscr.GetChar()
					gs.Start()
					gs.InitLevel()
				}
			}()
			gs.Restore(state)
			gs.Start()
		}
	} else {
		gs.Start()
		gs.InitLevel()
	}

	mapHeight := entity.ROOMS_IN_HEIGHT * entity.REGION_HEIGHT
	mapWidth := entity.ROOMS_IN_WIDTH * entity.REGION_WIDTH
	renderer := presentation.NewRenderer(mapHeight, mapWidth)

	stdscr.Clear()
	stdscr.Refresh()

	running := true
	gameEnded := false
	for running {
		presentation.DisplayMap(stdscr, renderer, gs.CurrentLevel, gs.Player)
		stdscr.Refresh()

		quit := presentation.ProcessInput(stdscr, gs)
		if quit {
			state := gs.ExportState()
			if err := saveRepo.SaveGame(state); err != nil {
				stdscr.Clear()
				stdscr.MovePrint(0, 0, "Error saving game: "+err.Error())
				wd, _ := os.Getwd()
				stdscr.MovePrint(1, 0, fmt.Sprintf("Current dir: %s", wd))
				stdscr.MovePrint(2, 0, "Press any key to exit...")
				stdscr.Refresh()
				_ = stdscr.GetChar()
			} else {
				savePath, _ := filepath.Abs("./save.json")
				stdscr.Clear()
				stdscr.MovePrint(0, 0, "Game saved successfully!")
				stdscr.MovePrint(1, 0, fmt.Sprintf("Saved to: %s", savePath))
				stdscr.MovePrint(2, 0, "Press any key to exit...")
				stdscr.Refresh()
				_ = stdscr.GetChar()
			}
			running = false
			gameEnded = false
		}

		if gs.IsPlayerAtExit() {
			renderer.ClearMap()
			gs.NextLevel()
			state := gs.ExportState()
			if err := saveRepo.SaveGame(state); err != nil {
				stdscr.MovePrint(0, 0, "Error saving game: "+err.Error())
				stdscr.Refresh()
			}

			run := game.RunResult{
				Statistics:    gs.Statistics,
				FinalLevel:    gs.Statistics.DeepestLevel,
				IsGameRunning: gs.IsRunning,
				Timestamp:     time.Now(),
			}
			_ = statRepo.SaveRun(run)
		}

		if gs.CheckGameOver() {
			running = false
			gameEnded = true
			presentation.DeadScreen(stdscr)
		}

		if gs.CurrentLevel.LevelNumber >= entity.LEVEL_NUM {
			running = false
			gameEnded = true
			presentation.EndgameScreen(stdscr)
		}
	}

	run := game.RunResult{
		Statistics:    gs.Statistics,
		FinalLevel:    gs.Statistics.DeepestLevel,
		IsGameRunning: gs.IsRunning,
		Timestamp:     time.Now(),
	}

	_ = statRepo.SaveRun(run)

	if gameEnded {
		_ = saveRepo.DeleteSave()
	}

	topRuns, _ = statRepo.LoadTop(10)
	if len(topRuns) > 0 {
		presentation.DisplayScoreboard(stdscr, topRuns)
	}
}
