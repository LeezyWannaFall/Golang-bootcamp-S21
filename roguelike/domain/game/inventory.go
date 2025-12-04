package game

import "roguelike/domain/entity"

type Backpack struct {
	CurrentSize int

	Foods []entity.Food
	FoodNumber int

	Elixirs []entity.Elixir
	ElixirNumber int

	Scrols []entity.Scroll
	ScrollNumber int

	Treasires []entity.Treasure

	Weapons []entity.Weapon
	WeaponNumber int
}