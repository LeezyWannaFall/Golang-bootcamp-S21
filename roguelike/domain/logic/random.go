package logic

import 	"math/rand/v2"

func GetRandomInRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func ShuffleEdges(edges []Edge) {
	rand.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
}
