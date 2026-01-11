package ports

import "roguelike/domain/game"

type StatisticRepository interface {
    SaveRun(stat game.SessionStatistics) error
    LoadTop(limit int) ([]game.SessionStatistics, error)
}
