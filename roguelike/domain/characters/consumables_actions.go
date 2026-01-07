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

func CheckConsumables(player *entity.Player, room entity.Room) {
	WasConsumed := false
	for i := 0; i < room.Consumables.ElixirNumber && !WasConsumed && player.Backpack.ElixirNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		if CheckEqualCoords()
	}
}

func TakeScroll(backpack entity.Backpack, scroll entity.Scroll, room entity.Room) {

}

func TakeElixir(backpack entity.Backpack, elixir entity.Elixir, room entity.Room) {

}

func TakeFood(backpack entity.Backpack, food entity.Food, room entity.Room) {

}

func TakeWeapon(backpack entity.Backpack, weapon entity.Weapon, room entity.Room) {

}

func DeleteFromRoom(room entity.Room, consumableCoords entity.Object, consumableType entity.ConsumableType) {

}

func ThrowOnGround(player entity.Player, room *entity.Room, weapon entity.Weapon) {

}

func CheckUnoccupiedRoom(room *entity.Room, coords entity.Object) bool {

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