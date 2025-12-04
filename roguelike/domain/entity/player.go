package entity

type Player struct {
	BaseStats Character
	RegenLimit int
	Backpack Backpack
	Weapon Weapon
	ElixirBuffs Buffs
}