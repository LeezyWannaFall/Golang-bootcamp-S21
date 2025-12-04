package entity

type FoodRoom struct {
	Geometry Object
	Food Food
}

type ElixirRoom struct {
	Geometry Object
	Elixir Elixir
}

type ScrollRoom struct {
	Geometry Object
	Scroll Scroll
}

type WeaponRoom struct {
	Geometry Object
	Weapon Weapon
}

type ConsumablesRoom struct {
	RoomFood  FoodRoom
	FoodNumber int

	RoomElixir ElixirRoom
	ElixirNumber int

	RoomScroll ScrollRoom
	ScrollNumber int

	WeaponRoom WeaponRoom
	WeaponNumber int
}