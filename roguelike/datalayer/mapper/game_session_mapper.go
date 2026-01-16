package mapper

import (
	"roguelike/datalayer/dto"
	"roguelike/domain/entity"
	"roguelike/domain/game"
	"time"
)

func ToDTO(state game.GameSessionState) dto.GameSessionDTO {
	return dto.GameSessionDTO{
		LevelNumber: state.LevelNumber,
		IsRunning:   state.IsRunning,
		Player:      playerToDTO(state.Player),
		Level:       levelToDTO(state.Level),
		Statistics:  sessionStatsToDTO(state.Statistics),
	}
}

func playerToDTO(ps game.PlayerState) dto.PlayerDTO {
	return dto.PlayerDTO{
		Position: dto.PositionDTO{
			X: ps.Position.XYcoords.X,
			Y: ps.Position.XYcoords.Y,
		},
		Health:     ps.Health,
		Agility:    ps.Agility,
		Strength:   ps.Strength,
		RegenLimit: ps.RegenLimit,
		Weapon:     weaponToDTO(ps.Weapon),
		Backpack:   backpackToDTO(ps.Backpack),
		Buffs:      buffsToDTO(ps.Buffs),
	}
}

func backpackToDTO(bs game.BackpackState) dto.BackpackDTO {
	keys := make([]int, 0)
	for i := 0; i < int(entity.KeyColorCount); i++ {
		if bs.Keys[i] {
			keys = append(keys, i)
		}
	}
	return dto.BackpackDTO{
		Foods:     foodsToDTO(bs.Foods),
		Elixirs:   elixirsToDTO(bs.Elixirs),
		Scrolls:   scrollsToDTO(bs.Scrolls),
		Weapons:   weaponsToDTO(bs.Weapons),
		Treasures: bs.Treasures.Value,
		Keys:      keys,
	}
}

func buffsToDTO(bs game.BuffsState) dto.BuffsDTO {
	return dto.BuffsDTO{
		MaxHealth: buffsSliceToDTO(bs.MaxHealth),
		Agility:   buffsSliceToDTO(bs.Agility),
		Strength:  buffsSliceToDTO(bs.Strength),
	}
}

func buffsSliceToDTO(buffs []game.BuffState) []dto.BuffDTO {
	result := make([]dto.BuffDTO, len(buffs))
	for i, b := range buffs {
		result[i] = dto.BuffDTO{
			StatIncrease: b.StatIncrease,
			EffectEnd:    b.EffectEnd,
		}
	}
	return result
}

func sessionStatsToDTO(s game.SessionStatistics) dto.SessionStatisticsDTO {
	return dto.SessionStatisticsDTO{
		TreasuresCollected: s.TreasuresCollected,
		DeepestLevel:       s.DeepestLevel,
		EnemiesDefeated:    s.EnemiesDefeated,
		FoodConsumed:       s.FoodConsumed,
		ElixirsDrunk:       s.ElixirsDrunk,
		ScrollsRead:        s.ScrollsRead,
		AttacksDealt:       s.AttacksDealt,
		AttacksMissed:      s.AttacksMissed,
		TilesTraveled:      s.TilesTraveled,
	}
}

func FromDTO(d dto.GameSessionDTO) game.GameSessionState {
	return game.GameSessionState{
		LevelNumber: d.LevelNumber,
		IsRunning:   d.IsRunning,
		Player:      playerFromDTO(d.Player),
		Level:       levelFromDTO(d.Level),
		Statistics:  sessionStatsFromDTO(d.Statistics),
	}
}

func playerFromDTO(d dto.PlayerDTO) game.PlayerState {
	return game.PlayerState{
		Position: entity.Object{
			XYcoords: entity.Pos{
				X: d.Position.X,
				Y: d.Position.Y,
			},
		},
		Health:     d.Health,
		Agility:    d.Agility,
		Strength:   d.Strength,
		RegenLimit: d.RegenLimit,
		Weapon:     weaponFromDTO(d.Weapon),
		Backpack:   backpackFromDTO(d.Backpack),
		Buffs:      buffsFromDTO(d.Buffs),
	}
}

func backpackFromDTO(d dto.BackpackDTO) game.BackpackState {
	keys := [entity.KeyColorCount]bool{}
	for _, keyIdx := range d.Keys {
		if keyIdx >= 0 && keyIdx < int(entity.KeyColorCount) {
			keys[keyIdx] = true
		}
	}
	return game.BackpackState{
		Foods:     foodsFromDTO(d.Foods),
		Elixirs:   elixirsFromDTO(d.Elixirs),
		Scrolls:   scrollsFromDTO(d.Scrolls),
		Weapons:   weaponsFromDTO(d.Weapons),
		Treasures: entity.Treasure{Value: d.Treasures},
		Keys:      keys,
	}
}

func sessionStatsFromDTO(d dto.SessionStatisticsDTO) game.SessionStatistics {
	return game.SessionStatistics{
		TreasuresCollected: d.TreasuresCollected,
		DeepestLevel:       d.DeepestLevel,
		EnemiesDefeated:    d.EnemiesDefeated,
		FoodConsumed:       d.FoodConsumed,
		ElixirsDrunk:       d.ElixirsDrunk,
		ScrollsRead:        d.ScrollsRead,
		AttacksDealt:       d.AttacksDealt,
		AttacksMissed:      d.AttacksMissed,
		TilesTraveled:      d.TilesTraveled,
	}
}

func levelToDTO(ls game.LevelState) dto.LevelDTO {
	rooms := make([]dto.RoomDTO, len(ls.Rooms))
	for i, r := range ls.Rooms {
		rooms[i] = roomToDTO(r)
	}

	passages := make([]dto.PassageDTO, len(ls.Passages))
	for i, p := range ls.Passages {
		passages[i] = passageToDTO(p)
	}

	doors := make([]dto.DoorDTO, len(ls.Doors))
	for i, door := range ls.Doors {
		doors[i] = dto.DoorDTO{
			Position: objectToDTO(door.Position),
			Color:    int(door.Color),
			IsOpen:   door.IsOpen,
		}
	}

	return dto.LevelDTO{
		Rooms:      rooms,
		Passages:   passages,
		EndOfLevel: objectToDTO(ls.EndOfLevel),
		Doors:      doors,
	}
}

func levelFromDTO(d dto.LevelDTO) game.LevelState {
	rooms := make([]game.RoomState, len(d.Rooms))
	for i, r := range d.Rooms {
		rooms[i] = roomFromDTO(r)
	}

	passages := make([]entity.Passage, len(d.Passages))
	for i, p := range d.Passages {
		passages[i] = passageFromDTO(p)
	}

	doors := make([]entity.Door, len(d.Doors))
	for i, door := range d.Doors {
		doors[i] = entity.Door{
			Position: objectFromDTO(door.Position),
			Color:    entity.KeyColor(door.Color),
			IsOpen:   door.IsOpen,
		}
	}

	return game.LevelState{
		Rooms:      rooms,
		Passages:   passages,
		EndOfLevel: objectFromDTO(d.EndOfLevel),
		Doors:      doors,
	}
}

func objectToDTO(o entity.Object) dto.ObjectDTO {
	return dto.ObjectDTO{
		Position: dto.PositionDTO{
			X: o.XYcoords.X,
			Y: o.XYcoords.Y,
		},
		Width:  o.W,
		Height: o.H,
	}
}

func objectFromDTO(d dto.ObjectDTO) entity.Object {
	return entity.Object{
		XYcoords: entity.Pos{
			X: d.Position.X,
			Y: d.Position.Y,
		},
		W: d.Width,
		H: d.Height,
	}
}

func passageToDTO(p entity.Passage) dto.PassageDTO {
	return dto.PassageDTO{
		Position: dto.PositionDTO{
			X: p.XYcoords.X,
			Y: p.XYcoords.Y,
		},
		Width:  p.W,
		Height: p.H,
	}
}

func passageFromDTO(d dto.PassageDTO) entity.Passage {
	return entity.Passage{
		XYcoords: entity.Pos{
			X: d.Position.X,
			Y: d.Position.Y,
		},
		W: d.Width,
		H: d.Height,
	}
}

func weaponToDTO(w entity.Weapon) dto.WeaponDTO {
	return dto.WeaponDTO{
		Name:     w.Name,
		Strength: w.Strength,
	}
}

func weaponFromDTO(d dto.WeaponDTO) entity.Weapon {
	return entity.Weapon{
		Name:     d.Name,
		Strength: d.Strength,
	}
}

func foodsToDTO(src []entity.Food) []dto.FoodDTO {
	dst := make([]dto.FoodDTO, len(src))
	for i, f := range src {
		dst[i] = dto.FoodDTO{
			Name:    f.Name,
			ToRegen: f.ToRegen,
		}
	}
	return dst
}

func foodsFromDTO(src []dto.FoodDTO) []entity.Food {
	dst := make([]entity.Food, len(src))
	for i, f := range src {
		dst[i] = entity.Food{
			Name:    f.Name,
			ToRegen: f.ToRegen,
		}
	}
	return dst
}

func elixirsToDTO(src []entity.Elixir) []dto.ElixirDTO {
	dst := make([]dto.ElixirDTO, len(src))
	for i, e := range src {
		dst[i] = dto.ElixirDTO{
			Name:     e.Name,
			Stat:     int(e.Stat),
			Increase: e.Increase,
			Duration: int64(e.Duration),
		}
	}
	return dst
}


func elixirsFromDTO(src []dto.ElixirDTO) []entity.Elixir {
	dst := make([]entity.Elixir, len(src))
	for i, e := range src {
		dst[i] = entity.Elixir{
			Name:     e.Name,
			Stat:     entity.StatType(e.Stat),
			Increase: e.Increase,
			Duration: time.Duration(e.Duration),
		}
	}
	return dst
}

func scrollsToDTO(src []entity.Scroll) []dto.ScrollDTO {
	dst := make([]dto.ScrollDTO, len(src))
	for i, s := range src {
		dst[i] = dto.ScrollDTO{
			Name:     s.Name,
			Stat:     int(s.Stat),
			Increase: s.Increase,
		}
	}
	return dst
}

func scrollsFromDTO(src []dto.ScrollDTO) []entity.Scroll {
	dst := make([]entity.Scroll, len(src))
	for i, s := range src {
		dst[i] = entity.Scroll{
			Name:     s.Name,
			Stat:     entity.StatType(s.Stat),
			Increase: s.Increase,
		}
	}
	return dst
}

func buffsFromDTO(d dto.BuffsDTO) game.BuffsState {
	return game.BuffsState{
		MaxHealth: buffsSliceFromDTO(d.MaxHealth),
		Agility:   buffsSliceFromDTO(d.Agility),
		Strength:  buffsSliceFromDTO(d.Strength),
	}
}

func buffsSliceFromDTO(src []dto.BuffDTO) []game.BuffState {
	dst := make([]game.BuffState, len(src))
	for i, b := range src {
		dst[i] = game.BuffState{
			StatIncrease: b.StatIncrease,
			EffectEnd:    b.EffectEnd,
		}
	}
	return dst
}

func roomToDTO(r game.RoomState) dto.RoomDTO {
	monsters := make([]dto.MonsterDTO, len(r.Monsters))
	for i, m := range r.Monsters {
		monsters[i] = dto.MonsterDTO{
			Position: dto.PositionDTO{
				X: m.Position.XYcoords.X,
				Y: m.Position.XYcoords.Y,
			},
			Health:    m.Health,
			Type:      int(m.Type),
			Hostility: int(m.Hostility),
			IsChasing: m.IsChasing,
			Direction: int(m.Direction),
		}
	}

	foods := make([]entity.Food, len(r.Consumables.Foods))
	for i, room := range r.Consumables.Foods {
		foods[i] = room.Food
	}
	
	elixirs := make([]entity.Elixir, len(r.Consumables.Elixirs))
	for i, room := range r.Consumables.Elixirs {
		elixirs[i] = room.Elixir
	}
	
	scrolls := make([]entity.Scroll, len(r.Consumables.Scrolls))
	for i, room := range r.Consumables.Scrolls {
		scrolls[i] = room.Scroll
	}
	
	weapons := make([]entity.Weapon, len(r.Consumables.Weapons))
	for i, room := range r.Consumables.Weapons {
		weapons[i] = room.Weapon
	}

	return dto.RoomDTO{
		Coordinates: objectToDTO(r.Coordinates),
		Monsters:    monsters,
		Consumables: dto.ConsumablesDTO{
			Foods:   foodsToDTO(foods),
			Elixirs: elixirsToDTO(elixirs),
			Scrolls: scrollsToDTO(scrolls),
			Weapons: weaponsToDTO(weapons),
		},
	}
}

func roomFromDTO(d dto.RoomDTO) game.RoomState {
	monsters := make([]game.MonsterState, len(d.Monsters))
	for i, m := range d.Monsters {
		monsters[i] = game.MonsterState{
			Position: entity.Object{
				XYcoords: entity.Pos{
					X: m.Position.X,
					Y: m.Position.Y,
				},
			},
			Health:    m.Health,
			Type:      entity.MonsterType(m.Type),
			Hostility: entity.HostilityType(m.Hostility),
			IsChasing: m.IsChasing,
			Direction: entity.Direction(m.Direction),
		}
	}

	foods := foodsFromDTO(d.Consumables.Foods)
	foodRooms := make([]entity.FoodRoom, len(foods))
	for i, food := range foods {
		foodRooms[i] = entity.FoodRoom{
			Geometry: entity.Object{},
			Food:     food,
		}
	}
	
	elixirs := elixirsFromDTO(d.Consumables.Elixirs)
	elixirRooms := make([]entity.ElixirRoom, len(elixirs))
	for i, elixir := range elixirs {
		elixirRooms[i] = entity.ElixirRoom{
			Geometry: entity.Object{},
			Elixir:   elixir,
		}
	}
	
	scrolls := scrollsFromDTO(d.Consumables.Scrolls)
	scrollRooms := make([]entity.ScrollRoom, len(scrolls))
	for i, scroll := range scrolls {
		scrollRooms[i] = entity.ScrollRoom{
			Geometry: entity.Object{},
			Scroll:   scroll,
		}
	}
	
	weapons := weaponsFromDTO(d.Consumables.Weapons)
	weaponRooms := make([]entity.WeaponRoom, len(weapons))
	for i, weapon := range weapons {
		weaponRooms[i] = entity.WeaponRoom{
			Geometry: entity.Object{},
			Weapon:   weapon,
		}
	}

	return game.RoomState{
		Coordinates: objectFromDTO(d.Coordinates),
		Monsters:    monsters,
		Consumables: game.ConsumablesState{
			Foods:   foodRooms,
			Elixirs: elixirRooms,
			Scrolls: scrollRooms,
			Weapons: weaponRooms,
		},
	}
}

func weaponsToDTO(src []entity.Weapon) []dto.WeaponDTO {
	dst := make([]dto.WeaponDTO, len(src))
	for i, w := range src {
		dst[i] = weaponToDTO(w)
	}
	return dst
}

func weaponsFromDTO(src []dto.WeaponDTO) []entity.Weapon {
	dst := make([]entity.Weapon, len(src))
	for i, w := range src {
		dst[i] = weaponFromDTO(w)
	}
	return dst
}