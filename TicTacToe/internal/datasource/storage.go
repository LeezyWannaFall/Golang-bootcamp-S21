package datasource

import (
	"sync"
)

type GameStorage struct {
	data sync.Map
}

func NewStorage() *GameStorage {
	return &GameStorage{}
}
