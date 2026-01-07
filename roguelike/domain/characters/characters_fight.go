package characters

import (
	"math/rand"
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)

const INITIAL_HIT_CHANCE =     70
const STANDART_AGILITY =       50
const AGILITY_FACTOR =        0.3
const INITIAL_DAMAGE =         30
const STANDART_STRENGTH =      50
const STRENGTH_FACTOR =       0.3
const STRENGTH_ADDITION =      65
const SLEEP_CHANCE =           15
const MAX_HP_PART =            10
const LOOT_AGILITY_FACTOR =   0.2
const LOOT_HP_FACTOR =        0.5
const LOOT_STRENGTH_FACTOR =  0.5
const MAXIMUM_FIGHTS =          8

type Turn int

const (
	PlayerTurn  Turn = iota
	MonsterTurn
)

type BattleInfo struct {
	Enemy             *entity.Monster
	VampireFirstAttack bool
	PlayerAsleep      bool
	OgreCooldown      bool
	IsFighting        bool
}

func Attack(player *entity.Player, battleInfo *BattleInfo, turn Turn) {
	switch turn {
	case PlayerTurn:
		if CheckHit(player, battleInfo.Enemy, PlayerTurn) {
			battleInfo.Enemy.Stats.Health -= CalculateDamage(player, battleInfo, PlayerTurn)
		}
		if battleInfo.Enemy.Stats.Health <= 0 {
			player.Backpack.Treasures.Value += int(CalculateLoot(battleInfo.Enemy))
		}
	case MonsterTurn:
		if CheckHit(player, battleInfo.Enemy, MonsterTurn) {
			player.BaseStats.Health -= CalculateDamage(player, battleInfo, MonsterTurn)
		}
	}
}

func CheckHit(player *entity.Player, monster *entity.Monster, turn Turn) bool {
	wasHit := false
	chance := INITIAL_HIT_CHANCE
	switch turn {
	case PlayerTurn:
		chance += HitChanceFormula(player.BaseStats.Agility, monster.Stats.Agility)
	case MonsterTurn:
		chance += HitChanceFormula(monster.Stats.Agility, player.BaseStats.Agility)
	}

	if rand.Intn(100) < chance || monster.Type == entity.Ogre {
		wasHit = true
	}
	return wasHit
}

func CalculateDamage(player *entity.Player, battleInfo *BattleInfo, turn Turn) float64 {
	damage := float64(INITIAL_DAMAGE)
	monsterDamageFormulas := map[entity.MonsterType]func(*BattleInfo) float64{
		entity.Zombie:  ZombieGhostDamageFormula,
		entity.Vampire: nil,
		entity.Ghost:   ZombieGhostDamageFormula,
		entity.Ogre:    OgreDamageFormula,
		entity.Snake:   SnakeDamageFormula,
	}

	switch turn {
	case PlayerTurn:
		if !(battleInfo.Enemy.Type == entity.Vampire && battleInfo.VampireFirstAttack) &&
			!(battleInfo.Enemy.Type == entity.Snake && battleInfo.PlayerAsleep) {
			if player.Weapon.Strength == entity.NO_WEAPON {
				damage += float64(player.BaseStats.Strength - STANDART_STRENGTH) * STRENGTH_FACTOR
			} else {
				damage = float64(player.Weapon.Strength) * float64(player.BaseStats.Strength + STRENGTH_ADDITION) / 100
			}
		} else if battleInfo.Enemy.Type == entity.Vampire && battleInfo.VampireFirstAttack {
			battleInfo.VampireFirstAttack = false
		} else {
			battleInfo.PlayerAsleep = false
		}
	case MonsterTurn:
		if battleInfo.Enemy.Type == entity.Vampire {
			damage = VampireDamageFormula(player)
		} else {
			damage = monsterDamageFormulas[battleInfo.Enemy.Type](battleInfo)
		}
	}
	return damage
}

func CalculateLoot(monster *entity.Monster) float64 {
	loot := float64(monster.Stats.Agility) * LOOT_AGILITY_FACTOR +
		float64(monster.Stats.Health) * LOOT_HP_FACTOR +
		float64(monster.Stats.Strength) * LOOT_STRENGTH_FACTOR +
		rand.Float64() * 20.0
	return loot
}

func ZombieGhostDamageFormula(battleInfo *BattleInfo) float64 {
	return INITIAL_DAMAGE + float64(battleInfo.Enemy.Stats.Strength - STANDART_STRENGTH) * STRENGTH_FACTOR
}

func OgreDamageFormula(battleInfo *BattleInfo) float64 {
	if !battleInfo.OgreCooldown {
		damage := float64(battleInfo.Enemy.Stats.Strength - STANDART_STRENGTH) * STRENGTH_FACTOR
		battleInfo.OgreCooldown = true
		return damage
	} else {
		battleInfo.OgreCooldown = false
		return 0
	}
}

func SnakeDamageFormula(battleInfo *BattleInfo) float64 {
	if rand.Intn(100) < SLEEP_CHANCE {
		battleInfo.PlayerAsleep = true
	}
	return ZombieGhostDamageFormula(battleInfo)
}

func VampireDamageFormula(player *entity.Player) float64 {
	return float64(player.RegenLimit) / MAX_HP_PART
}

func HitChanceFormula(attackerAgility int, defenderAgility int) int {
	return int(AGILITY_FACTOR * float64(attackerAgility - defenderAgility - STANDART_AGILITY))
}

func CheckEqualCoords(firstCoords, secondCoords entity.Pos) bool {
	return firstCoords.X == secondCoords.X && firstCoords.Y == secondCoords.Y
}

func CheckIfNeighborTile(first, second entity.Pos) bool {
	return (first.X == second.X && logic.Abs(first.Y - second.Y) == 1) || (first.Y == second.Y && logic.Abs(first.X - second.X) == 1)
}

func CheckIfDiagonallyNeighborTile(first, second entity.Pos) bool {
	return logic.Abs(first.X - second.X) == 1 && logic.Abs(first.Y - second.Y) == 1
}

func CheckUnique(monster *entity.Monster, battles_array []BattleInfo) bool {
	IsUnique := true


	for i := 0; i < MAXIMUM_FIGHTS && IsUnique; i++ {
		if battles_array[i].IsFighting && CheckEqualCoords(battles_array[i].Enemy.Stats.Pos.XYcoords, monster.Stats.Pos.XYcoords) {
			IsUnique = false
		}
	}
	return IsUnique
}