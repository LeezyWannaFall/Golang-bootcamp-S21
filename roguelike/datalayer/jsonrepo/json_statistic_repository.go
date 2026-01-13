package jsonrepo

import (
	"encoding/json"
	"os"
	"sort"

	"roguelike/domain/game"
	"roguelike/domain/ports"
)

type JSONStatisticRepository struct {
	path string
}

func NewJSONStatisticRepository(path string) ports.StatisticRepository {
	return &JSONStatisticRepository{
		path: path,
	}
}

func (r *JSONStatisticRepository) loadAll() ([]game.RunResult, error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return []game.RunResult{}, nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	var runs []game.RunResult
	if err := json.Unmarshal(data, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}

func (r *JSONStatisticRepository) SaveRun(run game.RunResult) error {
	runs, err := r.loadAll()
	if err != nil {
		return err
	}

	runs = append(runs, run)

	data, err := json.MarshalIndent(runs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, data, 0644)
}

func (r *JSONStatisticRepository) LoadTop(limit int) ([]game.RunResult, error) {
	runs, err := r.loadAll()
	if err != nil {
		return nil, err
	}
	
	sort.Slice(runs, func(i, j int) bool {
		if runs[i].IsGameRunning != runs[j].IsGameRunning {
			return runs[i].IsGameRunning
		}
		return runs[i].FinalLevel > runs[j].FinalLevel
	})

	if limit > 0 && len(runs) > limit {
		runs = runs[:limit]
	}

	return runs, nil
}
