package characters

import (
	"roguelike/domain/entity"
)

type ConsumableType int

const (
	NoneType ConsumableType = iota
	FoodType
	ElixirType
	WeaponType
	ScrollType
)

func EqualWeapons(weapon1, weapon2 entity.Weapon) bool {
	return weapon1.Strength == weapon2.Strength && weapon1.Name == weapon2.Name
}

func CheckConsumables(player *entity.Player, room *entity.Room) {
	pos := player.BaseStats.Pos.XYcoords

	// Elixirs
	for i := 0; i < room.Consumables.ElixirNumber; i++ {
		if player.Backpack.ElixirNumber >= entity.CONSUMABLES_TYPE_MAX_NUM {
			break
		}
		if CheckEqualCoords(room.Consumables.Elixirs[i].Geometry.XYcoords, pos) {
			TakeElixir(&player.Backpack, room, i)
			return
		}
	}

	// Scrolls
	for i := 0; i < room.Consumables.ScrollNumber; i++ {
		if player.Backpack.ScrollNumber >= entity.CONSUMABLES_TYPE_MAX_NUM {
			break
		}
		if CheckEqualCoords(room.Consumables.Scrolls[i].Geometry.XYcoords, pos) {
			TakeScroll(&player.Backpack, room, i)
			return
		}
	}

	// Food
	for i := 0; i < room.Consumables.FoodNumber; i++ {
		if player.Backpack.FoodNumber >= entity.CONSUMABLES_TYPE_MAX_NUM {
			break
		}
		if CheckEqualCoords(room.Consumables.Food[i].Geometry.XYcoords, pos) {
			TakeFood(&player.Backpack, room, i)
			return
		}
	}

	// Weapons
	for i := 0; i < room.Consumables.WeaponNumber; i++ {
		if player.Backpack.WeaponNumber >= entity.CONSUMABLES_TYPE_MAX_NUM {
			break
		}
		if CheckEqualCoords(room.Consumables.Weapons[i].Geometry.XYcoords, pos) {
			TakeWeapon(player, room, i)
			return
		}
	}
}


func TakeScroll(backpack *entity.Backpack, room *entity.Room, index int) {
	scroll := room.Consumables.Scrolls[index].Scroll

	DeleteFromRoom(room, index, ScrollType)

	backpack.Scrolls[backpack.ScrollNumber] = scroll
	backpack.ScrollNumber++
	backpack.CurrentSize++
}


func TakeElixir(backpack *entity.Backpack, room *entity.Room, index int) {
	elixir := room.Consumables.Elixirs[index].Elixir

	DeleteFromRoom(room, index, ElixirType)

	backpack.Elixirs[backpack.ElixirNumber] = elixir
	backpack.ElixirNumber++
	backpack.CurrentSize++
}


func TakeFood(backpack *entity.Backpack, room *entity.Room, index int) {
	food := room.Consumables.Food[index].Food

	DeleteFromRoom(room, index, FoodType)

	backpack.Foods[backpack.FoodNumber] = food
	backpack.FoodNumber++
	backpack.CurrentSize++
}


func TakeWeapon(player *entity.Player, room *entity.Room, index int) {
	weapon := room.Consumables.Weapons[index].Weapon

	DeleteFromRoom(room, index, WeaponType)

	player.Backpack.Weapons[player.Backpack.WeaponNumber] = weapon
	player.Backpack.WeaponNumber++
	player.Backpack.CurrentSize++
}


func DeleteFromRoom(room *entity.Room, index int, consumableType ConsumableType) {
	switch consumableType {

	case ElixirType:
		if index < room.Consumables.ElixirNumber {
			room.Consumables.Elixirs[index] =
				room.Consumables.Elixirs[room.Consumables.ElixirNumber-1]
			room.Consumables.ElixirNumber--
			room.ConsumablesNumber--
		}

	case ScrollType:
		if index < room.Consumables.ScrollNumber {
			room.Consumables.Scrolls[index] =
				room.Consumables.Scrolls[room.Consumables.ScrollNumber-1]
			room.Consumables.ScrollNumber--
			room.ConsumablesNumber--
		}

	case FoodType:
		if index < room.Consumables.FoodNumber {
			room.Consumables.Food[index] =
				room.Consumables.Food[room.Consumables.FoodNumber-1]
			room.Consumables.FoodNumber--
			room.ConsumablesNumber--
		}

	case WeaponType:
		if index < room.Consumables.WeaponNumber {
			room.Consumables.Weapons[index] =
				room.Consumables.Weapons[room.Consumables.WeaponNumber-1]
			room.Consumables.WeaponNumber--
			room.ConsumablesNumber--
		}
	}
}


func ThrowOnGround(player *entity.Player, room *entity.Room, weapon entity.Weapon) {
	room.Consumables.Weapons[room.Consumables.WeaponNumber].Geometry = player.BaseStats.Pos
	room.Consumables.Weapons[room.Consumables.WeaponNumber].Weapon = weapon

	roomWeapon := room.Consumables.Weapons[room.Consumables.WeaponNumber]
	MoveCharacterByDirection(entity.Right, &roomWeapon.Geometry.XYcoords)

	for direction := entity.Direction(0); IsOutsideRoom(roomWeapon.Geometry, *room) || !CheckUnoccupiedRoom(room, roomWeapon.Geometry); direction++ {
		roomWeapon = room.Consumables.Weapons[room.Consumables.WeaponNumber]
		if direction > entity.DiagonallyBackRight {
			direction = entity.Forward
		}
		MoveCharacterByDirection(direction, &roomWeapon.Geometry.XYcoords)
	}

	room.Consumables.Weapons[room.Consumables.WeaponNumber].Geometry = roomWeapon.Geometry
	room.ConsumablesNumber++
	room.Consumables.WeaponNumber++
}

func CheckUnoccupiedRoom(room *entity.Room, coords entity.Object) bool {
	unoccupied := true

	// Check elixirs
	for i := 0; i < room.Consumables.ElixirNumber && unoccupied; i++ {
		if coords.XYcoords.X == room.Consumables.Elixirs[i].Geometry.XYcoords.X &&
			coords.XYcoords.Y == room.Consumables.Elixirs[i].Geometry.XYcoords.Y {
			unoccupied = false
		}
	}

	// Check food
	for i := 0; i < room.Consumables.FoodNumber && unoccupied; i++ {
		if coords.XYcoords.X == room.Consumables.Food[i].Geometry.XYcoords.X &&
			coords.XYcoords.Y == room.Consumables.Food[i].Geometry.XYcoords.Y {
			unoccupied = false
		}
	}

	// Check scrolls
	for i := 0; i < room.Consumables.ScrollNumber && unoccupied; i++ {
		if coords.XYcoords.X == room.Consumables.Scrolls[i].Geometry.XYcoords.X &&
			coords.XYcoords.Y == room.Consumables.Scrolls[i].Geometry.XYcoords.Y {
			unoccupied = false
		}
	}

	// Check weapons
	for i := 0; i < room.Consumables.WeaponNumber && unoccupied; i++ {
		if coords.XYcoords.X == room.Consumables.Weapons[i].Geometry.XYcoords.X &&
			coords.XYcoords.Y == room.Consumables.Weapons[i].Geometry.XYcoords.Y {
			unoccupied = false
		}
	}

	// Check monsters
	for i := 0; i < room.MonsterNumbers && unoccupied; i++ {
		if coords.XYcoords.X == room.Monsters[i].Stats.Pos.XYcoords.X &&
			coords.XYcoords.Y == room.Monsters[i].Stats.Pos.XYcoords.Y {
			unoccupied = false
		}
	}

	return unoccupied
}

func CheckUnoccupiedLevel(level *entity.Level, coords entity.Object) bool {

}

func UseConsumable(player *entity.Player, consumablesType ConsumableType, room *entity.Room, index int) {

}

func EatFood(player *entity.Player, food entity.Food) {

}

func DrinkElixir(player *entity.Player, elixir entity.Elixir) {

}

func ReadScroll(player *entity.Player, scroll entity.Scroll) {

}

func RemoveFromBackpack(backpack entity.Backpack, consumableType ConsumableType, index int) {

}

func CheckTempEffectEnd(player *entity.Player) {

}
