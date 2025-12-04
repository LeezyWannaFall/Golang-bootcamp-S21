package entity

import "time"

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