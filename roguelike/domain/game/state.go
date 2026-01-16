package game

import (
	"roguelike/domain/characters"
	"roguelike/domain/entity"
	"roguelike/domain/logic"
	"time"
)

type RunResult struct {
	Statistics   SessionStatistics
	FinalLevel	 int
	IsGameRunning    bool
	Timestamp    time.Time
}


type GameSessionState struct {
	LevelNumber int
	IsRunning   bool

	Player PlayerState
	Level  LevelState

	Statistics SessionStatistics
}

type PlayerState struct {
	Position    entity.Object
	Health      float64
	Agility     int
	Strength    int

	RegenLimit  int
	Weapon      entity.Weapon
	Backpack    BackpackState
	Buffs       BuffsState
}

type BackpackState struct {
	Foods      []entity.Food
	Elixirs    []entity.Elixir
	Scrolls    []entity.Scroll
	Weapons    []entity.Weapon
	Treasures  entity.Treasure
	Keys       [entity.KeyColorCount]bool
}

type BuffState struct {
	StatIncrease int
	EffectEnd    int64
}

type BuffsState struct {
	MaxHealth []BuffState
	Agility   []BuffState
	Strength  []BuffState
}

type MonsterState struct {
	Position   entity.Object
	Health     float64
	Type       entity.MonsterType
	Hostility  entity.HostilityType
	IsChasing  bool
	Direction  entity.Direction
}

type RoomState struct {
	Coordinates entity.Object
	Monsters    []MonsterState
	Consumables ConsumablesState
}

type ConsumablesState struct {
	Foods   []entity.FoodRoom
	Elixirs []entity.ElixirRoom
	Scrolls []entity.ScrollRoom
	Weapons []entity.WeaponRoom
	Keys    []entity.KeyRoom
}

type LevelState struct {
	Coordinates entity.Object
	Rooms       []RoomState
	Passages    []entity.Passage
	EndOfLevel  entity.Object
	Doors       []entity.Door
}


func (s *GameSession) ExportState() GameSessionState {
    return GameSessionState{
        LevelNumber: s.CurrentLevel.LevelNumber,
        IsRunning:   s.IsRunning,
        Player:     exportPlayerState(s.Player),
        Level:      exportLevelState(s.CurrentLevel),
        Statistics: s.Statistics,
    }
}

func exportPlayerState(p *entity.Player) PlayerState {
	backpackState := BackpackState{
		Foods:     logic.CopySlice(p.Backpack.Foods[:p.Backpack.FoodNumber]),
		Elixirs:   logic.CopySlice(p.Backpack.Elixirs[:p.Backpack.ElixirNumber]),
		Scrolls:   logic.CopySlice(p.Backpack.Scrolls[:p.Backpack.ScrollNumber]),
		Weapons:   logic.CopySlice(p.Backpack.Weapons[:p.Backpack.WeaponNumber]),
		Treasures: p.Backpack.Treasures,
		Keys:      p.Backpack.Keys,
	}


	buffsState := BuffsState{
		MaxHealth: make([]BuffState, p.ElixirBuffs.CurrentHealthBuffNumber),
		Agility:   make([]BuffState, p.ElixirBuffs.CurrentAgilityBuffNumber),
		Strength:  make([]BuffState, p.ElixirBuffs.CurrentStrengthBuffNumber),
	}


	for i := 0; i < p.ElixirBuffs.CurrentHealthBuffNumber; i++ {
		buff := p.ElixirBuffs.MaxHealth[i]
		buffsState.MaxHealth[i] = BuffState{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	for i := 0; i < p.ElixirBuffs.CurrentAgilityBuffNumber; i++ {
		buff := p.ElixirBuffs.Agility[i]
		buffsState.Agility[i] = BuffState{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	for i := 0; i < p.ElixirBuffs.CurrentStrengthBuffNumber; i++ {
		buff := p.ElixirBuffs.Strength[i]
		buffsState.Strength[i] = BuffState{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	return PlayerState{
		Position:   p.BaseStats.Pos,
		Health:     p.BaseStats.Health,
		Agility:    p.BaseStats.Agility,
		Strength:   p.BaseStats.Strength,
		RegenLimit: p.RegenLimit,
		Weapon:     p.Weapon,
		Backpack:   backpackState,
		Buffs:      buffsState,
	}
}

func exportLevelState(l *entity.Level) LevelState {
	roomStates := make([]RoomState, entity.ROOMS_NUM)
	for i, room := range l.Rooms {
		monsterStates := make([]MonsterState, room.MonsterNumbers)
		for j := 0; j < room.MonsterNumbers; j++ {
			monster := room.Monsters[j]
			monsterStates[j] = MonsterState{
				Position:  monster.Stats.Pos,
				Health:    monster.Stats.Health,
				Type:      monster.Type,
				Hostility: monster.Hostility,
				IsChasing: monster.IsChasing,
				Direction: monster.Dir,
			}
		}

		consumablesState := ConsumablesState{
			Foods:   logic.CopySlice(room.Consumables.RoomFood[:room.Consumables.FoodNumber]),
			Elixirs: logic.CopySlice(room.Consumables.RoomElixir[:room.Consumables.ElixirNumber]),
			Scrolls: logic.CopySlice(room.Consumables.RoomScroll[:room.Consumables.ScrollNumber]),
			Weapons: logic.CopySlice(room.Consumables.WeaponRoom[:room.Consumables.WeaponNumber]),
			Keys:    logic.CopySlice(room.Consumables.RoomKeys[:room.Consumables.KeyNumber]),
		}

		roomStates[i] = RoomState{
			Coordinates: room.Coordinates,
			Monsters:    monsterStates,
			Consumables: consumablesState,
		}
	}

	return LevelState{
		Coordinates: l.Coordinates,
		Rooms:       roomStates,
		Passages:    logic.CopySlice(l.Passages.Passages),
		EndOfLevel:  l.EndOfLevel,
		Doors:       logic.CopySlice(l.Doors),
	}
}

func (s *GameSession) Restore(state GameSessionState) {
	s.IsRunning = state.IsRunning
	s.Statistics = state.Statistics

	s.Player = restorePlayer(state.Player)
	s.CurrentLevel = restoreLevel(state.Level)
	s.CurrentLevel.LevelNumber = state.LevelNumber
	
	if s.CurrentLevel.Coordinates.W == 0 || s.CurrentLevel.Coordinates.H == 0 {
		s.CurrentLevel.Coordinates = entity.Object{
			XYcoords: entity.Pos{X: 0, Y: 0},
			W:        entity.ROOMS_IN_WIDTH * entity.REGION_WIDTH,
			H:        entity.ROOMS_IN_HEIGHT * entity.REGION_HEIGHT,
		}
	}
	
	s.UpdateCurrentRoom()
	
	if s.CurrentRoom < 0 || s.CurrentRoom >= entity.ROOMS_NUM {
		s.CurrentRoom = 0
		for i := 0; i < entity.ROOMS_NUM; i++ {
			room := &s.CurrentLevel.Rooms[i]
			if s.Player.BaseStats.Pos.XYcoords.X >= room.Coordinates.XYcoords.X &&
				s.Player.BaseStats.Pos.XYcoords.X < room.Coordinates.XYcoords.X+room.Coordinates.W &&
				s.Player.BaseStats.Pos.XYcoords.Y >= room.Coordinates.XYcoords.Y &&
				s.Player.BaseStats.Pos.XYcoords.Y < room.Coordinates.XYcoords.Y+room.Coordinates.H {
				s.CurrentRoom = i
				break
			}
		}
	}
	
	playerPos := entity.Pos{
		X: s.Player.BaseStats.Pos.XYcoords.X,
		Y: s.Player.BaseStats.Pos.XYcoords.Y,
	}
	
	if characters.IsOutsideLevel(playerPos, s.CurrentLevel) {
		if s.CurrentRoom >= 0 && s.CurrentRoom < entity.ROOMS_NUM {
			room := &s.CurrentLevel.Rooms[s.CurrentRoom]
			if room.Coordinates.W > 0 && room.Coordinates.H > 0 {
				s.Player.BaseStats.Pos.XYcoords.X = room.Coordinates.XYcoords.X + room.Coordinates.W/2
				s.Player.BaseStats.Pos.XYcoords.Y = room.Coordinates.XYcoords.Y + room.Coordinates.H/2
			}
		}
	}
}

func restorePlayer(ps PlayerState) *entity.Player {
	player := &entity.Player{}

	player.BaseStats = entity.Character{
		Pos:      ps.Position,
		Health:   ps.Health,
		Agility:  ps.Agility,
		Strength: ps.Strength,
	}

	player.RegenLimit = ps.RegenLimit
	player.Weapon = ps.Weapon

	player.Backpack = entity.Backpack{
		FoodNumber:    len(ps.Backpack.Foods),
		ElixirNumber:  len(ps.Backpack.Elixirs),
		ScrollNumber:  len(ps.Backpack.Scrolls),
		WeaponNumber:  len(ps.Backpack.Weapons),
		Treasures:    ps.Backpack.Treasures,
		CurrentSize: len(ps.Backpack.Foods) + len(ps.Backpack.Elixirs) + len(ps.Backpack.Scrolls) + len(ps.Backpack.Weapons),
		Keys:         ps.Backpack.Keys,
	}

	copy(player.Backpack.Foods[:], ps.Backpack.Foods)
	copy(player.Backpack.Elixirs[:], ps.Backpack.Elixirs)
	copy(player.Backpack.Scrolls[:], ps.Backpack.Scrolls)
	copy(player.Backpack.Weapons[:], ps.Backpack.Weapons)

	player.ElixirBuffs = entity.Buffs{}

	player.ElixirBuffs.CurrentHealthBuffNumber = len(ps.Buffs.MaxHealth)
	for i, buff := range ps.Buffs.MaxHealth {
		player.ElixirBuffs.MaxHealth[i] = entity.Buff{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	player.ElixirBuffs.CurrentAgilityBuffNumber = len(ps.Buffs.Agility)
	for i, buff := range ps.Buffs.Agility {
		player.ElixirBuffs.Agility[i] = entity.Buff{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	player.ElixirBuffs.CurrentStrengthBuffNumber = len(ps.Buffs.Strength)
	for i, buff := range ps.Buffs.Strength {
		player.ElixirBuffs.Strength[i] = entity.Buff{
			StatIncrease: buff.StatIncrease,
			EffectEnd:    buff.EffectEnd,
		}
	}

	return player
}

func restoreLevel(ls LevelState) *entity.Level {
	level := &entity.Level{
		Coordinates: ls.Coordinates,
		EndOfLevel:  ls.EndOfLevel,
		Doors:       logic.CopySlice(ls.Doors),
		DoorNumber:  len(ls.Doors),
	}

	level.Passages = entity.Passages{
		Passages:       logic.CopySlice(ls.Passages),
		PassagesNumber: len(ls.Passages),
	}

	for i, rs := range ls.Rooms {
		room := &level.Rooms[i]
		room.Coordinates = rs.Coordinates
		room.MonsterNumbers = len(rs.Monsters)

		for j, ms := range rs.Monsters {
			room.Monsters[j] = entity.Monster{
				Stats: entity.Character{
					Pos:    ms.Position,
					Health: ms.Health,
				},
				Type:      ms.Type,
				Hostility: ms.Hostility,
				IsChasing: ms.IsChasing,
				Dir:       ms.Direction,
			}
		}

		room.Consumables.FoodNumber = len(rs.Consumables.Foods)
		copy(room.Consumables.RoomFood[:], rs.Consumables.Foods)

		room.Consumables.ElixirNumber = len(rs.Consumables.Elixirs)
		copy(room.Consumables.RoomElixir[:], rs.Consumables.Elixirs)

		room.Consumables.ScrollNumber = len(rs.Consumables.Scrolls)
		copy(room.Consumables.RoomScroll[:], rs.Consumables.Scrolls)

		room.Consumables.WeaponNumber = len(rs.Consumables.Weapons)
		copy(room.Consumables.WeaponRoom[:], rs.Consumables.Weapons)

		room.Consumables.KeyNumber = len(rs.Consumables.Keys)
		copy(room.Consumables.RoomKeys[:], rs.Consumables.Keys)
	}

	return level
}

