package datasource

import (
	"sync"
)

type GameStorage struct {
	data sync.Map
}
