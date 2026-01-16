package dto

type PlayerDTO struct {
	Position   PositionDTO `json:"position"`
	Health     float64     `json:"health"`
	Agility    int         `json:"agility"`
	Strength   int         `json:"strength"`
	RegenLimit int         `json:"regen_limit"`
	Weapon     WeaponDTO   `json:"weapon"`
	Backpack   BackpackDTO `json:"backpack"`
	Buffs      BuffsDTO    `json:"buffs"`
}

type PositionDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type CharacterStatsDTO struct {
	Position PositionDTO `json:"position"`
	Health   float64     `json:"health"`
	Agility  int         `json:"agility"`
	Strength int         `json:"strength"`
}

type WeaponDTO struct {
	Name     string `json:"name"`
	Strength int    `json:"strength"`
}

type BackpackDTO struct {
	Foods     []FoodDTO   `json:"foods"`
	Elixirs   []ElixirDTO `json:"elixirs"`
	Scrolls   []ScrollDTO `json:"scrolls"`
	Weapons   []WeaponDTO `json:"weapons"`
	Treasures int         `json:"treasures"`
	Keys      []int       `json:"keys"`
}

type BuffDTO struct {
	StatIncrease int   `json:"stat_increase"`
	EffectEnd    int64 `json:"effect_end"`
}

type BuffsDTO struct {
	MaxHealth []BuffDTO `json:"max_health"`
	Agility   []BuffDTO `json:"agility"`
	Strength  []BuffDTO `json:"strength"`
}
