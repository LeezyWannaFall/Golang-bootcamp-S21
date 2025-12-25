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

type Object struct {
    X    int
    Y    int
    W    int
    H    int
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
}

type Treasure struct {
	Value int
}

type Food struct {
	ToRegen int
	Name   	string
}

type Elixir struct {
	duration time.Duration
	stat 	 StatType
	increase int
	Name     string
}

type Scroll struct {
	stat 	 StatType
	increase int
	Name     string
}

type Weapon struct {
	Strength int
	Name   string
}

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
	RoomFood  [MAX_CONSUMABLES_PER_ROOM]FoodRoom
	FoodNumber int

	RoomElixir [MAX_CONSUMABLES_PER_ROOM]ElixirRoom
	ElixirNumber int

	RoomScroll [MAX_CONSUMABLES_PER_ROOM]ScrollRoom
	ScrollNumber int

	WeaponRoom [MAX_CONSUMABLES_PER_ROOM]WeaponRoom
	WeaponNumber int
}

type Room struct {
	Coordinates Object
	Consumables ConsumablesRoom
	ConsumablesNumber int
	Monsters [MAX_MONSTERS_PER_ROOM]Monster
	MonsterNumbers int
}

type Backpack struct {
	CurrentSize int

	Foods Food
	FoodNumber int

	Elixirs Elixir
	ElixirNumber int

	Scrols Scroll
	ScrollNumber int

	Treasures Treasure

	Weapons Weapon
	WeaponNumber int
}

type Buff struct {
	StatIncrease int
	EffectEnd time.Duration
}

type Buffs struct {
	MaxHealth Buff
	CurrentHealthBuffNumber int

	Agility Buff
	CurrentAgilityBuffNumber int

	Strength Buff
	CurrentStrengthBuffNumber int
}

type Player struct {
	BaseStats Character
	RegenLimit int
	Backpack Backpack
	Weapon Weapon
	ElixirBuffs Buffs
}

type Passage Object

type Passages struct {
	Passages []Passage
	PassagesNumber int
}

type Level struct {
	Coordinates Object
	Rooms [ROOMS_NUM]Room
	Passages Passages
	LevelNumber int
	EndOfLevel Object
}