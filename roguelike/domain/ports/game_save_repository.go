package ports

import "roguelike/domain/game"

type GameSaveRepository interface {
    SaveGame(state game.GameSessionState) error
    LoadGame() (game.GameSessionState, error)
    HasSave() (bool, error)
    DeleteSave() error
}
