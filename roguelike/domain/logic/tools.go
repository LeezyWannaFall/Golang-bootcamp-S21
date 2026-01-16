package logic

import "roguelike/domain/entity"

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

func CheckUnoccupiedRoom(coordinates *entity.Object, room *entity.Room) bool {
	for i := 0; i < room.Consumables.ElixirNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomElixir[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomElixir[i].Geometry.XYcoords.Y {
			return false
		}
	}

	for i := 0; i < room.Consumables.FoodNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomFood[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomFood[i].Geometry.XYcoords.Y {
			return false
		}
	}

	for i := 0; i < room.Consumables.ScrollNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomScroll[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomScroll[i].Geometry.XYcoords.Y {
			return false
		}
	}

	for i := 0; i < room.Consumables.WeaponNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.WeaponRoom[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.WeaponRoom[i].Geometry.XYcoords.Y {
			return false
		}
	}

	for i := 0; i < room.MonsterNumbers; i++ {
		if coordinates.XYcoords.X == room.Monsters[i].Stats.Pos.XYcoords.X &&
			coordinates.XYcoords.Y == room.Monsters[i].Stats.Pos.XYcoords.Y {
			return false
		}
	}

	return true
}

// CopySlice создает копию слайса любого типа
func CopySlice[T any](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}
