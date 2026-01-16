package characters

import (
	"roguelike/domain/entity"
	"roguelike/domain/logic"
	"time"
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
	wasConsumed := false
	for i := 0; i < room.Consumables.ElixirNumber && !wasConsumed && player.Backpack.ElixirNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		if CheckEqualCoords(room.Consumables.RoomElixir[i].Geometry.XYcoords, player.BaseStats.Pos.XYcoords) {
			TakeElixir(&player.Backpack, room, &room.Consumables.RoomElixir[i])
			player.Backpack.CurrentSize++
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.ScrollNumber && !wasConsumed && player.Backpack.ScrollNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		if CheckEqualCoords(room.Consumables.RoomScroll[i].Geometry.XYcoords, player.BaseStats.Pos.XYcoords) {
			TakeScroll(&player.Backpack, room, &room.Consumables.RoomScroll[i])
			player.Backpack.CurrentSize++
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.FoodNumber && !wasConsumed && player.Backpack.FoodNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		if CheckEqualCoords(room.Consumables.RoomFood[i].Geometry.XYcoords, player.BaseStats.Pos.XYcoords) {
			TakeFood(&player.Backpack, room, &room.Consumables.RoomFood[i])
			player.Backpack.CurrentSize++
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.WeaponNumber && !wasConsumed && player.Backpack.WeaponNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		if CheckEqualCoords(room.Consumables.WeaponRoom[i].Geometry.XYcoords, player.BaseStats.Pos.XYcoords) {
			TakeWeapon(player, room, &room.Consumables.WeaponRoom[i])
			player.Backpack.CurrentSize++
			wasConsumed = true
		}
	}
}

func TakeScroll(backpack *entity.Backpack, room *entity.Room, scroll *entity.ScrollRoom) {
	DeleteFromRoom(room, &scroll.Geometry, ScrollType)
	backpack.Scrolls[backpack.ScrollNumber] = scroll.Scroll
	backpack.ScrollNumber++
}

func TakeElixir(backpack *entity.Backpack, room *entity.Room, elixir *entity.ElixirRoom) {
	DeleteFromRoom(room, &elixir.Geometry, ElixirType)
	backpack.Elixirs[backpack.ElixirNumber] = elixir.Elixir
	backpack.ElixirNumber++
}

func TakeFood(backpack *entity.Backpack, room *entity.Room, food *entity.FoodRoom) {
	DeleteFromRoom(room, &food.Geometry, FoodType)
	backpack.Foods[backpack.FoodNumber] = food.Food
	backpack.FoodNumber++
}

func TakeWeapon(player *entity.Player, room *entity.Room, weapon *entity.WeaponRoom) {
	DeleteFromRoom(room, &weapon.Geometry, WeaponType)
	player.Backpack.Weapons[player.Backpack.WeaponNumber] = weapon.Weapon
	player.Backpack.WeaponNumber++
}

func DeleteFromRoom(room *entity.Room, consumableCoords *entity.Object, consumableType ConsumableType) {
	switch consumableType {
	case ElixirType:
		for i := 0; i < room.Consumables.ElixirNumber; i++ {
			if CheckEqualCoords(consumableCoords.XYcoords, room.Consumables.RoomElixir[i].Geometry.XYcoords) {
				room.Consumables.RoomElixir[i] = room.Consumables.RoomElixir[room.Consumables.ElixirNumber-1]
				room.Consumables.ElixirNumber--
				break
			}
		}
	case FoodType:
		for i := 0; i < room.Consumables.FoodNumber; i++ {
			if CheckEqualCoords(consumableCoords.XYcoords, room.Consumables.RoomFood[i].Geometry.XYcoords) {
				room.Consumables.RoomFood[i] = room.Consumables.RoomFood[room.Consumables.FoodNumber-1]
				room.Consumables.FoodNumber--
				break
			}
		}
	case ScrollType:
		for i := 0; i < room.Consumables.ScrollNumber; i++ {
			if CheckEqualCoords(consumableCoords.XYcoords, room.Consumables.RoomScroll[i].Geometry.XYcoords) {
				room.Consumables.RoomScroll[i] = room.Consumables.RoomScroll[room.Consumables.ScrollNumber-1]
				room.Consumables.ScrollNumber--
				break
			}
		}
	case WeaponType:
		for i := 0; i < room.Consumables.WeaponNumber; i++ {
			if CheckEqualCoords(consumableCoords.XYcoords, room.Consumables.WeaponRoom[i].Geometry.XYcoords) {
				room.Consumables.WeaponRoom[i] = room.Consumables.WeaponRoom[room.Consumables.WeaponNumber-1]
				room.Consumables.WeaponNumber--
				break
			}
		}
	}
	room.ConsumablesNumber--
}

func ThrowOnGround(player *entity.Player, room *entity.Room, weapon entity.Weapon) {
	room.Consumables.WeaponRoom[room.Consumables.WeaponNumber].Geometry = player.BaseStats.Pos
	room.Consumables.WeaponRoom[room.Consumables.WeaponNumber].Weapon = weapon
	roomWeapon := room.Consumables.WeaponRoom[room.Consumables.WeaponNumber]
	MoveCharacterByDirection(entity.Right, &roomWeapon.Geometry.XYcoords)

	for direction := entity.Direction(0); IsOutsideRoom(roomWeapon.Geometry, *room) || !logic.CheckUnoccupiedRoom(&roomWeapon.Geometry, room); direction++ {
		roomWeapon = room.Consumables.WeaponRoom[room.Consumables.WeaponNumber]
		MoveCharacterByDirection(direction, &roomWeapon.Geometry.XYcoords)
		if direction >= entity.Stop {
			break
		}
	}
	room.Consumables.WeaponRoom[room.Consumables.WeaponNumber].Geometry = roomWeapon.Geometry
	room.ConsumablesNumber++
	room.Consumables.WeaponNumber++
}

func CheckUnoccupiedLevel(level *entity.Level, coords entity.Object) bool {
	unoccupied := true
	for i := 0; unoccupied && i < entity.ROOMS_NUM; i++ {
		unoccupied = logic.CheckUnoccupiedRoom(&coords, &level.Rooms[i])
	}
	if CheckEqualCoords(coords.XYcoords, level.EndOfLevel.XYcoords) {
		unoccupied = false
	}
	return unoccupied
}

func UseConsumable(player *entity.Player, consumablesType ConsumableType, room *entity.Room, index int) {
	var oldWeapon entity.Weapon
	switch consumablesType {
	case ScrollType:
		if index >= 0 && index < player.Backpack.ScrollNumber {
			ReadScroll(player, player.Backpack.Scrolls[index])
			RemoveFromBackpack(&player.Backpack, ScrollType, index)
			player.Backpack.CurrentSize--
		}
	case ElixirType:
		if index >= 0 && index < player.Backpack.ElixirNumber {
			DrinkElixir(player, player.Backpack.Elixirs[index])
			RemoveFromBackpack(&player.Backpack, ElixirType, index)
			player.Backpack.CurrentSize--
		}
	case FoodType:
		if index >= 0 && index < player.Backpack.FoodNumber {
			EatFood(player, player.Backpack.Foods[index])
			RemoveFromBackpack(&player.Backpack, FoodType, index)
			player.Backpack.CurrentSize--
		}
	case WeaponType:
		if index == -1 {
			player.Weapon.Strength = entity.NO_WEAPON
		} else if room != nil && index >= 0 && index < player.Backpack.WeaponNumber {
			if !EqualWeapons(player.Backpack.Weapons[index], player.Weapon) && player.Weapon.Strength != entity.NO_WEAPON {
				oldWeapon = player.Weapon
				player.Weapon = player.Backpack.Weapons[index]
				ThrowOnGround(player, room, oldWeapon)
			} else if player.Weapon.Strength == entity.NO_WEAPON {
				player.Weapon = player.Backpack.Weapons[index]
			} else {
				player.Weapon = player.Backpack.Weapons[0]
			}
		}
	default:
		break
	}
}

func EatFood(player *entity.Player, food entity.Food) {
	newHealth := player.BaseStats.Health + float64(food.ToRegen)
	if newHealth > float64(player.RegenLimit) {
		player.BaseStats.Health = float64(player.RegenLimit)
	} else {
		player.BaseStats.Health = newHealth
	}
}

func DrinkElixir(player *entity.Player, elixir entity.Elixir) {
	currentTime := time.Now().Unix()
	durationSeconds := int64(elixir.Duration.Seconds())
	switch elixir.Stat {
	case entity.Health:
		idx := player.ElixirBuffs.CurrentHealthBuffNumber
		player.ElixirBuffs.MaxHealth[idx].StatIncrease += elixir.Increase
		player.ElixirBuffs.MaxHealth[idx].EffectEnd = currentTime + durationSeconds
		player.RegenLimit += elixir.Increase
		player.BaseStats.Health += float64(elixir.Increase)
		player.ElixirBuffs.CurrentHealthBuffNumber++
	case entity.Agility:
		idx := player.ElixirBuffs.CurrentAgilityBuffNumber
		player.ElixirBuffs.Agility[idx].StatIncrease += elixir.Increase
		player.ElixirBuffs.Agility[idx].EffectEnd = currentTime + durationSeconds
		player.BaseStats.Agility += elixir.Increase
		player.ElixirBuffs.CurrentAgilityBuffNumber++
	case entity.Strength:
		idx := player.ElixirBuffs.CurrentStrengthBuffNumber
		player.ElixirBuffs.Strength[idx].StatIncrease += elixir.Increase
		player.ElixirBuffs.Strength[idx].EffectEnd = currentTime + durationSeconds
		player.BaseStats.Strength += elixir.Increase
		player.ElixirBuffs.CurrentStrengthBuffNumber++
	}
}

func ReadScroll(player *entity.Player, scroll entity.Scroll) {
	switch scroll.Stat {
	case entity.Health:
		player.RegenLimit += scroll.Increase
		player.BaseStats.Health += float64(scroll.Increase)
	case entity.Agility:
		player.BaseStats.Agility += scroll.Increase
	case entity.Strength:
		player.BaseStats.Strength += scroll.Increase
	}
}

func RemoveFromBackpack(backpack *entity.Backpack, consumableType ConsumableType, index int) {
	switch consumableType {
	case ScrollType:
		backpack.Scrolls[index] = backpack.Scrolls[backpack.ScrollNumber-1]
		backpack.ScrollNumber--
	case ElixirType:
		backpack.Elixirs[index] = backpack.Elixirs[backpack.ElixirNumber-1]
		backpack.ElixirNumber--
	case FoodType:
		backpack.Foods[index] = backpack.Foods[backpack.FoodNumber-1]
		backpack.FoodNumber--
	default:
		break
	}
}

func CheckTempEffectEnd(player *entity.Player) {
	currentTime := time.Now().Unix()

	for i := 0; i < player.ElixirBuffs.CurrentHealthBuffNumber; {
		if player.ElixirBuffs.MaxHealth[i].EffectEnd > currentTime {
			i++
		} else {
			player.RegenLimit -= player.ElixirBuffs.MaxHealth[i].StatIncrease
			player.BaseStats.Health -= float64(player.ElixirBuffs.MaxHealth[i].StatIncrease)
			if player.BaseStats.Health <= 0 {
				player.BaseStats.Health = 1
			}
			player.ElixirBuffs.MaxHealth[i] = player.ElixirBuffs.MaxHealth[player.ElixirBuffs.CurrentHealthBuffNumber-1]
			player.ElixirBuffs.CurrentHealthBuffNumber--
		}
	}

	for i := 0; i < player.ElixirBuffs.CurrentAgilityBuffNumber; {
		if player.ElixirBuffs.Agility[i].EffectEnd > currentTime {
			i++
		} else {
			player.BaseStats.Agility -= player.ElixirBuffs.Agility[i].StatIncrease
			player.ElixirBuffs.Agility[i] = player.ElixirBuffs.Agility[player.ElixirBuffs.CurrentAgilityBuffNumber-1]
			player.ElixirBuffs.CurrentAgilityBuffNumber--
		}
	}

	for i := 0; i < player.ElixirBuffs.CurrentStrengthBuffNumber; {
		if player.ElixirBuffs.Strength[i].EffectEnd > currentTime {
			i++
		} else {
			player.BaseStats.Strength -= player.ElixirBuffs.Strength[i].StatIncrease
			player.ElixirBuffs.Strength[i] = player.ElixirBuffs.Strength[player.ElixirBuffs.CurrentStrengthBuffNumber-1]
			player.ElixirBuffs.CurrentStrengthBuffNumber--
		}
	}
}
