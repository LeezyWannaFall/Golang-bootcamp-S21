package game

import (
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)

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
}

type LevelState struct {
	Rooms       []RoomState
	Passages    []entity.Passage
	EndOfLevel  entity.Object
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
		}

		roomStates[i] = RoomState{
			Coordinates: room.Coordinates,
			Monsters:    monsterStates,
			Consumables: consumablesState,
		}
	}

	return LevelState{
		Rooms:      roomStates,
		Passages:   logic.CopySlice(l.Passages.Passages),
		EndOfLevel: l.EndOfLevel,
	}
}

func (s *GameSession) Restore(state GameSessionState) {
	s.IsRunning = state.IsRunning
	s.CurrentLevel.LevelNumber = state.LevelNumber
	s.Statistics = state.Statistics

	s.Player = restorePlayer(state.Player)
	s.CurrentLevel = restoreLevel(state.Level)
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
		EndOfLevel:  ls.EndOfLevel,
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
	}

	return level
}

