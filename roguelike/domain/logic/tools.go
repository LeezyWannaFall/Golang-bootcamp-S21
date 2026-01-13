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
	// Проверка эликсиров
	for i := 0; i < room.Consumables.ElixirNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomElixir[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomElixir[i].Geometry.XYcoords.Y {
			return false
		}
	}

	// Проверка еды
	for i := 0; i < room.Consumables.FoodNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomFood[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomFood[i].Geometry.XYcoords.Y {
			return false
		}
	}

	// Проверка свитков
	for i := 0; i < room.Consumables.ScrollNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.RoomScroll[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.RoomScroll[i].Geometry.XYcoords.Y {
			return false
		}
	}

	// Проверка оружия
	for i := 0; i < room.Consumables.WeaponNumber; i++ {
		if coordinates.XYcoords.X == room.Consumables.WeaponRoom[i].Geometry.XYcoords.X &&
			coordinates.XYcoords.Y == room.Consumables.WeaponRoom[i].Geometry.XYcoords.Y {
			return false
		}
	}

	// Проверка монстров
	for i := 0; i < room.MonsterNumbers; i++ {
		if coordinates.XYcoords.X == room.Monsters[i].Stats.Pos.XYcoords.X &&
			coordinates.XYcoords.Y == room.Monsters[i].Stats.Pos.XYcoords.Y {
			return false
		}
	}

	return true
}

func CopySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func FoodRoomsToFoods(src []entity.FoodRoom) []entity.Food {
	dst := make([]entity.Food, len(src))
	for i, f := range src {
		dst[i] = entity.Food{
			Name:    f.Food.Name,
			ToRegen: f.Food.ToRegen,
		}
	}
	return dst
}

func ElixirRoomsToElixirs(src []entity.ElixirRoom) []entity.Elixir {
	dst := make([]entity.Elixir, len(src))
	for i, e := range src {
		dst[i] = entity.Elixir{
			Name:     e.Elixir.Name,
			Stat:     e.Elixir.Stat,
			Increase: e.Elixir.Increase,
			Duration: e.Elixir.Duration,
		}
	}
	return dst
}

func ScrollsRoomsToScrolls(src []entity.ScrollRoom) []entity.Scroll {
	dst := make([]entity.Scroll, len(src))
	for i, s := range src {
		dst[i] = entity.Scroll{
			Name:     s.Scroll.Name,
			Stat:     s.Scroll.Stat,
			Increase: s.Scroll.Increase,
		}
	}
	return dst
}

func WeaponRoomsToWeapon(src []entity.WeaponRoom) []entity.Weapon {
	dst := make([]entity.Weapon, len(src))
	for i, w := range src {
		dst[i] = entity.Weapon{
			Name:     w.Weapon.Name,
			Strength: w.Weapon.Strength,
		}
	}
	return dst
}

func FoodsToFoodRooms(src []entity.Food) []entity.FoodRoom {
	dst := make([]entity.FoodRoom, len(src))
	for i, f := range src {
		dst[i] = entity.FoodRoom{
			Food: entity.Food{
				Name:    f.Name,
				ToRegen: f.ToRegen,
			},
		}
	}
	return dst
}

func ElixirsToElixirRooms(src []entity.Elixir) []entity.ElixirRoom {
	dst := make([]entity.ElixirRoom, len(src))
	for i, e := range src {
		dst[i] = entity.ElixirRoom{
			Elixir: entity.Elixir{
				Name:     e.Name,
				Stat:     e.Stat,
				Increase: e.Increase,
				Duration: e.Duration,
			},
		}
	}
	return dst
}

func ScrollsToScrollsRooms(src []entity.Scroll) []entity.ScrollRoom {
	dst := make([]entity.ScrollRoom, len(src))
	for i, s := range src {
		dst[i] = entity.ScrollRoom{
			Scroll: entity.Scroll{
				Name:     s.Name,
				Stat:     s.Stat,
				Increase: s.Increase,
			},
		}
	}
	return dst
}

func WeaponsToWeaponRooms(src []entity.Weapon) []entity.WeaponRoom {
	dst := make([]entity.WeaponRoom, len(src))
	for i, w := range src {
		dst[i] = entity.WeaponRoom{
			Weapon: entity.Weapon{
				Name:     w.Name,
				Strength: w.Strength,
			},
		}
	}
	return dst
}