package game

import (
	"roguelike/domain/characters"
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)

func (gs *GameSession) InitLevel() {
	logic.ClearData(gs.CurrentLevel)
	logic.GenerateRooms(gs.CurrentLevel.Rooms[:])
	logic.GeneratePassages(&gs.CurrentLevel.Passages, gs.CurrentLevel.Rooms[:])
	gs.CurrentRoom = logic.GeneratePlayer(gs.CurrentLevel.Rooms[:], gs.Player)
	logic.GenerateExit(gs.CurrentLevel, gs.CurrentRoom)
	logic.GenerateMonsters(gs.CurrentLevel, gs.CurrentRoom)
	// TODO: добавить генерацию всех типов предметов
}

func (gs *GameSession) GameLoop() {
	gs.Start()
	gs.InitLevel()

	for gs.IsRunning {
		// TODO: интеграция с presentation layer для получения ввода
		gs.ProcessMonstersTurn()

		if gs.CheckGameOver() {
			gs.Stop()
			break
		}

		if gs.CurrentLevel.LevelNumber >= entity.LEVEL_NUM {
			gs.Stop()
			break
		}

		if gs.IsPlayerAtExit() {
			gs.NextLevel()
		}
	}
}

func (gs *GameSession) ProcessPlayerTurn(direction entity.Direction) {
	prevPos := gs.Player.BaseStats.Pos
	characters.MoveCharacterByDirectionObj(direction, &gs.Player.BaseStats.Pos) // move using Object (костыль, но пока что так)

	// TODO: добавить проверку через CheckOutsideBorder

	monster := gs.FindMonsterAtPosition(&gs.Player.BaseStats.Pos)
	if monster != nil {
		gs.Player.BaseStats.Pos = prevPos
		gs.InitiateBattle(monster)
	}

	if prevPos.X != gs.Player.BaseStats.Pos.X || prevPos.Y != gs.Player.BaseStats.Pos.Y {
		gs.IncrementTilesTraveled()
	}

	// TODO: добавить логику подбора предметов
}

func (gs *GameSession) ProcessMonstersTurn() {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &gs.CurrentLevel.Rooms[i]
		for j := 0; j < room.MonsterNumbers; j++ {
			monster := &room.Monsters[j]
			if monster.Stats.Health <= 0 {
				continue
			}
			characters.MoveMonster(monster, gs.CurrentLevel, gs.Player) // add gs.Player
		}
	}
}

func (gs *GameSession) InitiateBattle(monster *entity.Monster) {
	battleInfo := &characters.BattleInfo{
		Enemy:              monster,
		VampireFirstAttack: monster.Type == entity.Vampire,
		PlayerAsleep:       false,
		OgreCooldown:       false,
		IsFighting:         true,
	}

	characters.Attack(gs.Player, battleInfo, characters.PlayerTurn)
	gs.IncrementAttacksDealt()

	if monster.Stats.Health > 0 {
		characters.Attack(gs.Player, battleInfo, characters.MonsterTurn)
	} else {
		gs.IncrementEnemiesDefeated()
	}
}

func (gs *GameSession) FindMonsterAtPosition(pos *entity.Object) *entity.Monster {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &gs.CurrentLevel.Rooms[i]
		for j := 0; j < room.MonsterNumbers; j++ {
			monster := &room.Monsters[j]
			if monster.Stats.Pos.X == pos.X && monster.Stats.Pos.Y == pos.Y {
				return monster
			}
		}
	}
	return nil
}

func (gs *GameSession) IsPlayerAtExit() bool {
	playerPos := &gs.Player.BaseStats.Pos
	exit := &gs.CurrentLevel.EndOfLevel
	return playerPos.X == exit.X && playerPos.Y == exit.Y
}

func (gs *GameSession) CheckGameOver() bool {
	return gs.Player.BaseStats.Health <= 0
}

func (gs *GameSession) NextLevel() {
	gs.CurrentLevel.LevelNumber++
	gs.Statistics.DeepestLevel = gs.CurrentLevel.LevelNumber
	gs.InitLevel()
}
