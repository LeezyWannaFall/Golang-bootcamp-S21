package entity

import "time"

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