package characters

import (
	"roguelike/domain/entity"
)

func CheckConsumable(player *entity.Player, room *entity.Room) {
	wasConsumable := false

	for i := 0; i < room.Consumables.ElixirNum && !wasConsumable && player.Backpack.ElixirNum < entity.ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Elixirs[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			TakeElixir(&player.Backpack, room, &room.Consumables.Elixirs[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.ScrollNum && !wasConsumable && player.Backpack.ScrollNum < entity.ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Scrolls[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			TakeScroll(&player.Backpack, room, &room.Consumables.Scrolls[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.FoodNum && !wasConsumable && player.Backpack.FoodNum < entity.ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.RoomFood[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			TakeFood(&player.Backpack, room, &room.Consumables.RoomFood[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.WeaponNum && !wasConsumable && player.Backpack.WeaponNum < entity.ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Weapons[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			TakeWeapon(player, room, &room.Consumables.Weapons[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}
}
