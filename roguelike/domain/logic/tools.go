package logic

// import (
// 	"roguelike/domain/entity"
// )

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func ObjCoordsToPos(object entity.Object) entity.Pos {
// 	pos := entity.Pos{
// 		X: object.X,
// 		Y: object.Y,
// 	}

// 	return pos
// }