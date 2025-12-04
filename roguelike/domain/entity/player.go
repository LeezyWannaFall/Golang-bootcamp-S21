package entity

import "roguelike/domain/game"

type Player struct {
	BaseStats Character
	RegenLimit int
	Backpack game.Backpack
	Weapon Weapon
	ElixirBuffs Buffs
}