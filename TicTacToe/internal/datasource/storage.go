package datasource

import (
	"sync"
)

type GameStore struct {
	data sync.Map
}
