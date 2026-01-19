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
	
	stats := logic.SessionStatistics{
		TreasuresCollected: gs.Statistics.TreasuresCollected,
		DeepestLevel:       gs.Statistics.DeepestLevel,
		EnemiesDefeated:    gs.Statistics.EnemiesDefeated,
		FoodConsumed:       gs.Statistics.FoodConsumed,
		ElixirsDrunk:       gs.Statistics.ElixirsDrunk,
		ScrollsRead:        gs.Statistics.ScrollsRead,
		AttacksDealt:       gs.Statistics.AttacksDealt,
		AttacksMissed:      gs.Statistics.AttacksMissed,
		TilesTraveled:      gs.Statistics.TilesTraveled,
	}
	
	balance := logic.CalculateBalanceAdjustment(stats, gs.CurrentLevel.LevelNumber, gs.Player)
	logic.GenerateMonsters(gs.CurrentLevel, gs.CurrentRoom, balance)
	logic.GenerateConsumables(gs.CurrentLevel, gs.CurrentRoom, gs.Player, gs.CurrentLevel.LevelNumber, balance)
	logic.GenerateDoorsAndKeys(gs.CurrentLevel, gs.CurrentRoom)
}

func (gs *GameSession) GameLoop() {
	gs.Start()
	gs.InitLevel()

	for gs.IsRunning {
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
	if direction == entity.Stop {
		return
	}

	prevPos := gs.Player.BaseStats.Pos

	characters.MoveCharacterByDirectionObj(direction, &gs.Player.BaseStats.Pos)

	newCoords := entity.Pos{
		X: gs.Player.BaseStats.Pos.XYcoords.X,
		Y: gs.Player.BaseStats.Pos.XYcoords.Y,
	}

	if characters.IsOutsideLevel(newCoords, gs.CurrentLevel) {
		gs.Player.BaseStats.Pos = prevPos
		return
	}

	if !characters.IsPassable(newCoords, gs.CurrentLevel) {
		gs.TryOpenDoor(newCoords)
		if !characters.IsPassable(newCoords, gs.CurrentLevel) {
			gs.Player.BaseStats.Pos = prevPos
			return
		}
	}

	characters.CheckTempEffectEnd(gs.Player)

	monster := gs.FindMonsterAtPosition(&gs.Player.BaseStats.Pos)
	if monster != nil {
		gs.Player.BaseStats.Pos = prevPos
		gs.InitiateBattle(monster, gs.CurrentLevel)
		return
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

func (gs *GameSession) InitiateBattle(monster *entity.Monster, level *entity.Level) {
	battleInfo := &characters.BattleInfo{
		Enemy:              monster,
		VampireFirstAttack: monster.Type == entity.Vampire,
		PlayerAsleep:       false,
		OgreCooldown:       false,
		IsFighting:         true,
	}

	characters.Attack(gs.Player, battleInfo, characters.PlayerTurn, level)
	gs.IncrementAttacksDealt()

	if monster.Stats.Health > 0 {
		characters.Attack(gs.Player, battleInfo, characters.MonsterTurn, level)
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

	for i := 0; i < room.Consumables.KeyNumber; i++ {
		itemPos := &room.Consumables.RoomKeys[i].Geometry
		if playerPos.XYcoords.X == itemPos.XYcoords.X && playerPos.XYcoords.Y == itemPos.XYcoords.Y {
			keyColor := room.Consumables.RoomKeys[i].Key.Color
			gs.Player.Backpack.Keys[keyColor] = true
			gs.removeKeyFromRoom(room, i)
			break
		}
	}
}

func (gs *GameSession) removeKeyFromRoom(room *entity.Room, index int) {
	if index < room.Consumables.KeyNumber-1 {
		room.Consumables.RoomKeys[index] = room.Consumables.RoomKeys[room.Consumables.KeyNumber-1]
	}
	room.Consumables.KeyNumber--
}

func (gs *GameSession) TryOpenDoor(pos entity.Pos) {
	for i := 0; i < gs.CurrentLevel.DoorNumber; i++ {
		door := &gs.CurrentLevel.Doors[i]
		if pos.X == door.Position.XYcoords.X && pos.Y == door.Position.XYcoords.Y {
			if !door.IsOpen && gs.Player.Backpack.Keys[door.Color] {
				door.IsOpen = true
			}
			break
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
