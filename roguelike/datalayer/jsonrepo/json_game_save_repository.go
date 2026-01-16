package jsonrepo

import (
	"encoding/json"
	"errors"
	"os"

	"roguelike/datalayer/dto"
	"roguelike/datalayer/mapper"
	"roguelike/domain/game"
	"roguelike/domain/ports"
)

type JSONGameSaveRepository struct {
	path string
}

func NewJSONGameSaveRepository(path string) ports.GameSaveRepository {
	return &JSONGameSaveRepository{
		path: path,
	}
}

func (r *JSONGameSaveRepository) SaveGame(state game.GameSessionState) error {
	dto := mapper.ToDTO(state)
	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(r.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONGameSaveRepository) LoadGame() (game.GameSessionState, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return game.GameSessionState{}, err
	}

	var dto dto.GameSessionDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return game.GameSessionState{}, err
	}

	state := mapper.FromDTO(dto)
	return state, nil
}

func (r *JSONGameSaveRepository) HasSave() (bool, error) {
	_, err := os.Stat(r.path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (r *JSONGameSaveRepository) DeleteSave() error {
	if _, err := os.Stat(r.path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(r.path)
}

