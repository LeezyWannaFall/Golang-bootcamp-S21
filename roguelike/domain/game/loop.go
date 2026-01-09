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
	logic.GenerateConsumables(gs.CurrentLevel, gs.CurrentRoom, gs.Player, gs.CurrentLevel.LevelNumber)
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
	characters.MoveCharacterByDirectionObj(direction, &gs.Player.BaseStats.Pos)

	characters.CheckTempEffectEnd(gs.Player)

	monster := gs.FindMonsterAtPosition(&gs.Player.BaseStats.Pos)
	if monster != nil {
		gs.Player.BaseStats.Pos = prevPos
		gs.InitiateBattle(monster)
	}

	if prevPos.XYcoords.X != gs.Player.BaseStats.Pos.XYcoords.X || prevPos.XYcoords.Y != gs.Player.BaseStats.Pos.XYcoords.Y {
		gs.IncrementTilesTraveled()
	}

	gs.UpdateCurrentRoom()
	currentRoom := gs.GetCurrentRoom()
	if currentRoom != nil {
		gs.CheckAndPickupItems(currentRoom)
	}
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
			if monster.Stats.Pos.XYcoords.X == pos.XYcoords.X && monster.Stats.Pos.XYcoords.Y == pos.XYcoords.Y {
				return monster
			}
		}
	}
	return nil
}

func (gs *GameSession) IsPlayerAtExit() bool {
	playerPos := &gs.Player.BaseStats.Pos
	exit := &gs.CurrentLevel.EndOfLevel
	return playerPos.XYcoords.X == exit.XYcoords.X && playerPos.XYcoords.Y == exit.XYcoords.Y
}

func (gs *GameSession) CheckGameOver() bool {
	return gs.Player.BaseStats.Health <= 0
}

func (gs *GameSession) NextLevel() {
	gs.CurrentLevel.LevelNumber++
	gs.Statistics.DeepestLevel = gs.CurrentLevel.LevelNumber
	gs.InitLevel()
}

func (gs *GameSession) CheckAndPickupItems(room *entity.Room) {
	playerPos := &gs.Player.BaseStats.Pos
	wasConsumed := false

	for i := 0; i < room.Consumables.ElixirNumber && !wasConsumed && gs.Player.Backpack.ElixirNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		itemPos := &room.Consumables.RoomElixir[i].Geometry
		if playerPos.XYcoords.X == itemPos.XYcoords.X && playerPos.XYcoords.Y == itemPos.XYcoords.Y {
			gs.Player.Backpack.Elixirs[gs.Player.Backpack.ElixirNumber] = room.Consumables.RoomElixir[i].Elixir
			gs.Player.Backpack.ElixirNumber++
			gs.Player.Backpack.CurrentSize++
			gs.removeElixirFromRoom(room, i)
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.ScrollNumber && !wasConsumed && gs.Player.Backpack.ScrollNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		itemPos := &room.Consumables.RoomScroll[i].Geometry
		if playerPos.XYcoords.X == itemPos.XYcoords.X && playerPos.XYcoords.Y == itemPos.XYcoords.Y {
			gs.Player.Backpack.Scrolls[gs.Player.Backpack.ScrollNumber] = room.Consumables.RoomScroll[i].Scroll
			gs.Player.Backpack.ScrollNumber++
			gs.Player.Backpack.CurrentSize++
			gs.removeScrollFromRoom(room, i)
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.FoodNumber && !wasConsumed && gs.Player.Backpack.FoodNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		itemPos := &room.Consumables.RoomFood[i].Geometry
		if playerPos.XYcoords.X == itemPos.XYcoords.X && playerPos.XYcoords.Y == itemPos.XYcoords.Y {
			gs.Player.Backpack.Foods[gs.Player.Backpack.FoodNumber] = room.Consumables.RoomFood[i].Food
			gs.Player.Backpack.FoodNumber++
			gs.Player.Backpack.CurrentSize++
			gs.removeFoodFromRoom(room, i)
			wasConsumed = true
		}
	}

	for i := 0; i < room.Consumables.WeaponNumber && !wasConsumed && gs.Player.Backpack.WeaponNumber < entity.CONSUMABLES_TYPE_MAX_NUM; i++ {
		itemPos := &room.Consumables.WeaponRoom[i].Geometry
		if playerPos.XYcoords.X == itemPos.XYcoords.X && playerPos.XYcoords.Y == itemPos.XYcoords.Y {
			gs.Player.Backpack.Weapons[gs.Player.Backpack.WeaponNumber] = room.Consumables.WeaponRoom[i].Weapon
			gs.Player.Backpack.WeaponNumber++
			gs.Player.Backpack.CurrentSize++
			gs.removeWeaponFromRoom(room, i)
			wasConsumed = true
		}
	}
}

func (gs *GameSession) removeElixirFromRoom(room *entity.Room, index int) {
	if index < room.Consumables.ElixirNumber-1 {
		room.Consumables.RoomElixir[index] = room.Consumables.RoomElixir[room.Consumables.ElixirNumber-1]
	}
	room.Consumables.ElixirNumber--
}

func (gs *GameSession) removeScrollFromRoom(room *entity.Room, index int) {
	if index < room.Consumables.ScrollNumber-1 {
		room.Consumables.RoomScroll[index] = room.Consumables.RoomScroll[room.Consumables.ScrollNumber-1]
	}
	room.Consumables.ScrollNumber--
}

func (gs *GameSession) removeFoodFromRoom(room *entity.Room, index int) {
	if index < room.Consumables.FoodNumber-1 {
		room.Consumables.RoomFood[index] = room.Consumables.RoomFood[room.Consumables.FoodNumber-1]
	}
	room.Consumables.FoodNumber--
}

func (gs *GameSession) removeWeaponFromRoom(room *entity.Room, index int) {
	if index < room.Consumables.WeaponNumber-1 {
		room.Consumables.WeaponRoom[index] = room.Consumables.WeaponRoom[room.Consumables.WeaponNumber-1]
	}
	room.Consumables.WeaponNumber--
}
