package entity

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