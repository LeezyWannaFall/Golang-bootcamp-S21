package dto

type ObjectDTO struct {
	Position PositionDTO `json:"position"`
	Width    int         `json:"width"`
	Height   int         `json:"height"`
}

type MonsterDTO struct {
	Position  PositionDTO `json:"position"`
	Health    float64     `json:"health"`
	Type      int         `json:"type"`
	Hostility int         `json:"hostility"`
	IsChasing bool        `json:"is_chasing"`
	Direction int         `json:"direction"`
}

type FoodDTO struct {
	Name       string  `json:"name"`
	ToRegen	   int `json:"to_regen"`
}

type ElixirDTO struct {
	Name          string  `json:"name"`
	Stat		  int     `json:"stat_increase"`
	Increase	  int     `json:"increase"`
	Duration      int64   `json:"duration"`
}

type ScrollDTO struct {
	Name        string `json:"name"`
	Stat		int    `json:"stat_increase"`
	Increase	int    `json:"increase"`
}

type ConsumablesDTO struct {
	Foods   []FoodDTO   `json:"foods"`
	Elixirs []ElixirDTO `json:"elixirs"`
	Scrolls []ScrollDTO `json:"scrolls"`
	Weapons []WeaponDTO `json:"weapons"`
}

type RoomDTO struct {
	Coordinates ObjectDTO     `json:"coordinates"`
	Monsters    []MonsterDTO    `json:"monsters"`
	Consumables ConsumablesDTO  `json:"consumables"`
}

type PassageDTO struct {
	Position PositionDTO `json:"position"`
	Width    int         `json:"width"`
	Height   int         `json:"height"`
}

type DoorDTO struct {
	Position ObjectDTO `json:"position"`
	Color    int       `json:"color"`
	IsOpen   bool      `json:"is_open"`
}

type LevelDTO struct {
	Rooms      []RoomDTO    `json:"rooms"`
	Passages   []PassageDTO `json:"passages"`
	EndOfLevel ObjectDTO    `json:"end_of_level"`
	Doors      []DoorDTO     `json:"doors"`
}
