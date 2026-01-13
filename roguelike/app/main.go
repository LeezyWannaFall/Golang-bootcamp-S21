package main

import (
	"roguelike/datalayer/jsonrepo"
	"roguelike/domain/game"
	"time"
	"fmt"
)

func main() {
	// --- repositories ---
	saveRepo := jsonrepo.NewJSONGameSaveRepository("save.json")
	statRepo := jsonrepo.NewJSONStatisticRepository("stats.json")

	// --- game session ---
	session := game.NewGameSession()
	session.Start()

	// --- load save if exists ---
	if hasSave, _ := saveRepo.HasSave(); hasSave {
		fmt.Println("Найдено сохранение. Продолжаем игру...")
		state, err := saveRepo.LoadGame()
		if err == nil {
			session.Restore(state)
		}
	} else {
		fmt.Println("Новая игра")
	}

	// --- start game loop ---
	session.Start()

	// --- game finished ---
	run := game.RunResult{
		Statistics:   session.Statistics,
		FinalLevel: session.Statistics.DeepestLevel,
		Completed:    session.IsRunning,
		Timestamp:    time.Now(),
	}

	_ = statRepo.SaveRun(run)
	_ = saveRepo.DeleteSave()
}