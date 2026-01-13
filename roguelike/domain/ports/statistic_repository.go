package ports

import "roguelike/domain/game"

type StatisticRepository interface {
    SaveRun(run game.RunResult) error
    LoadTop(limit int) ([]game.RunResult, error)
}
