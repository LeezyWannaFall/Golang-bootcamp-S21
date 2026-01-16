package entity

import "time"

type Dimension int

const (
	X Dimension = iota
	Y
)

type MonsterType int

const (
	Zombie MonsterType = iota
	Vampire
	Ghost
	Ogre
	Snake
	Mimic
)

type HostilityType int

const (
	Low HostilityType = iota
	Medium
	High
)

type StatType int

const (
	Health StatType = iota
	Agility
	Strength
)

type KeyColor int

const (
	RedKey KeyColor = iota
	BlueKey
	YellowKey
	GreenKey
	KeyColorCount
)

type Direction int

const (
	Forward Direction = iota
	Back
	Left
	Right
	DiagonallyForwardLeft
	DiagonallyForwardRight
	DiagonallyBackLeft
	DiagonallyBackRight
	Stop
)

type Pos struct {
	X int
	Y int
}

type Object struct {
	XYcoords Pos
	W        int
	H        int
}

type Character struct {
	Pos      Object
	Health   float64
	Agility  int
	Strength int
}

type Monster struct {
	Stats     Character
	Type      MonsterType
	Hostility HostilityType
	IsChasing bool
	Dir       Direction
	IsDead    bool
}

type Treasure struct {
	Value int
}

type Food struct {
	ToRegen int
	Name    string
}

type Elixir struct {
	Duration time.Duration
	Stat     StatType
	Increase int
	Name     string
}

type Scroll struct {
	Stat     StatType
	Increase int
	Name     string
}

type Weapon struct {
	Strength int
	Name     string
}

type Key struct {
	Color KeyColor
}

type KeyRoom struct {
	Geometry Object
	Key      Key
}

type Door struct {
	Position Object
	Color    KeyColor
	IsOpen   bool
}

type FoodRoom struct {
	Geometry Object
	Food     Food
}

type ElixirRoom struct {
	Geometry Object
	Elixir   Elixir
}

type ScrollRoom struct {
	Geometry Object
	Scroll   Scroll
}

type WeaponRoom struct {
	Geometry Object
	Weapon   Weapon
}

type ConsumablesRoom struct {
	RoomFood   [MAX_CONSUMABLES_PER_ROOM]FoodRoom
	FoodNumber int

	RoomElixir   [MAX_CONSUMABLES_PER_ROOM]ElixirRoom
	ElixirNumber int

	RoomScroll   [MAX_CONSUMABLES_PER_ROOM]ScrollRoom
	ScrollNumber int

	WeaponRoom   [MAX_CONSUMABLES_PER_ROOM]WeaponRoom
	WeaponNumber int

	RoomKeys   [KeyColorCount]KeyRoom
	KeyNumber  int
}

type Room struct {
	Coordinates       Object
	Consumables       ConsumablesRoom
	ConsumablesNumber int
	Monsters          [MAX_MONSTERS_PER_ROOM]Monster
	MonsterNumbers    int
}

type Backpack struct {
	CurrentSize int

	Foods      [CONSUMABLES_TYPE_MAX_NUM]Food
	FoodNumber int

	Elixirs      [CONSUMABLES_TYPE_MAX_NUM]Elixir
	ElixirNumber int

	Scrolls      [CONSUMABLES_TYPE_MAX_NUM]Scroll
	ScrollNumber int

	Treasures Treasure

	Weapons      [CONSUMABLES_TYPE_MAX_NUM]Weapon
	WeaponNumber int

	Keys [KeyColorCount]bool
}

type Buff struct {
	StatIncrease int
	EffectEnd    int64 // Unix timestamp в секундах
}

type Buffs struct {
	MaxHealth               [CONSUMABLES_TYPE_MAX_NUM]Buff
	CurrentHealthBuffNumber int

	Agility                  [CONSUMABLES_TYPE_MAX_NUM]Buff
	CurrentAgilityBuffNumber int

	Strength                  [CONSUMABLES_TYPE_MAX_NUM]Buff
	CurrentStrengthBuffNumber int
}

type Player struct {
	BaseStats   Character
	RegenLimit  int
	Backpack    Backpack
	Weapon      Weapon
	ElixirBuffs Buffs
}

type Passage Object

type Passages struct {
	Passages       []Passage
	PassagesNumber int
}

type Level struct {
	Coordinates Object
	Rooms       [ROOMS_NUM]Room
	Passages    Passages
	LevelNumber int
	EndOfLevel  Object
	Doors       []Door
	DoorNumber  int
}
