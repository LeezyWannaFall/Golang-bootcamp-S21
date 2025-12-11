package logic

import 	"math/rand/v2"

func GetRandomInRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}