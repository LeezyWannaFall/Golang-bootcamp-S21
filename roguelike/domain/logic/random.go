package logic

import (
	"math/rand/v2"
	"roguelike/domain/datastructs"
)

func GetRandomInRange(min, max int) int {
	if max < min {
		min, max = max, min
	}
	if max == min {
		return min
	}
	return rand.IntN(max-min+1) + min
}

func ShuffleEdges(edges []datastructs.Edge) {
	rand.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
}
