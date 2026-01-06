package game

import (
	"roguelike/domain/entity"
)

type SessionStatistics struct {
	TreasuresCollected int // Общее количество собранных сокровищ
	DeepestLevel       int // Самый глубокий достигнутый уровень
	EnemiesDefeated    int // Количество побежденных врагов
	FoodConsumed       int // Количество съеденной еды
	ElixirsDrunk       int // Количество выпитых эликсиров
	ScrollsRead        int // Количество прочитанных свитков
	AttacksDealt       int // Общее количество нанесенных ударов
	AttacksMissed      int // Количество промахов
	TilesTraveled      int // Количество пройденных клеток
}

type GameSession struct {
	Player       *entity.Player
	CurrentLevel *entity.Level
	Statistics   SessionStatistics
	CurrentRoom  int  // Индекс текущей комнаты игрока
	IsRunning    bool // Флаг активности игровой сессии
}

// NewGameSession создает новую игровую сессию с начальными параметрами
func NewGameSession() *GameSession {
	return &GameSession{
		Player: &entity.Player{
			BaseStats: entity.Character{
				Pos: entity.Object{
					X: 0,
					Y: 0,
					W: 1,
					H: 1,
				},
				Health:    500,
				// MaxHealth: 500, не нашел в структуре Character
				Agility:   70,
				Strength:  70,
			},
			RegenLimit: 500,
			Backpack: entity.Backpack{
				CurrentSize:  0,
				FoodNumber:   0,
				ElixirNumber: 0,
				ScrollNumber: 0,
				WeaponNumber: 0,
				Treasures: entity.Treasure{
					Value: 0,
				},
			},
			Weapon: entity.Weapon{
				Strength: entity.NO_WEAPON,
				Name:     "",
			},
			ElixirBuffs: entity.Buffs{
				CurrentHealthBuffNumber:   0,
				CurrentAgilityBuffNumber:  0,
				CurrentStrengthBuffNumber: 0,
			},
		},
		CurrentLevel: &entity.Level{
			LevelNumber: 0,
			Rooms:       [entity.ROOMS_NUM]entity.Room{},
			Passages: entity.Passages{
				Passages:       []entity.Passage{},
				PassagesNumber: 0,
			},
		},
		Statistics: SessionStatistics{
			TreasuresCollected: 0,
			DeepestLevel:       0,
			EnemiesDefeated:    0,
			FoodConsumed:       0,
			ElixirsDrunk:       0,
			ScrollsRead:        0,
			AttacksDealt:       0,
			AttacksMissed:      0,
			TilesTraveled:      0,
		},
		CurrentRoom: 0,
		IsRunning:   false,
	}
}

func (gs *GameSession) Start() {
	gs.IsRunning = true
}

func (gs *GameSession) Stop() {
	gs.IsRunning = false
}

func (gs *GameSession) UpdateStatistics(stat SessionStatistics) {
	gs.Statistics = stat
}

func (gs *GameSession) IncrementEnemiesDefeated() {
	gs.Statistics.EnemiesDefeated++
}

func (gs *GameSession) IncrementFoodConsumed() {
	gs.Statistics.FoodConsumed++
}

func (gs *GameSession) IncrementElixirsDrunk() {
	gs.Statistics.ElixirsDrunk++
}

func (gs *GameSession) IncrementScrollsRead() {
	gs.Statistics.ScrollsRead++
}

func (gs *GameSession) IncrementAttacksDealt() {
	gs.Statistics.AttacksDealt++
}

func (gs *GameSession) IncrementAttacksMissed() {
	gs.Statistics.AttacksMissed++
}

func (gs *GameSession) IncrementTilesTraveled() {
	gs.Statistics.TilesTraveled++
}

// GetCurrentRoom возвращает указатель на текущую комнату игрока
func (gs *GameSession) GetCurrentRoom() *entity.Room {
	if gs.CurrentRoom >= 0 && gs.CurrentRoom < entity.ROOMS_NUM {
		return &gs.CurrentLevel.Rooms[gs.CurrentRoom]
	}
	return nil
}

// UpdateCurrentRoom обновляет индекс текущей комнаты на основе позиции игрока
func (gs *GameSession) UpdateCurrentRoom() {
	playerPos := &gs.Player.BaseStats.Pos
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &gs.CurrentLevel.Rooms[i]
		if isInsideRoom(playerPos, &room.Coordinates) {
			gs.CurrentRoom = i
			return
		}
	}
}

// isInsideRoom проверяет, находится ли позиция внутри комнаты
func isInsideRoom(pos *entity.Object, roomCoords *entity.Object) bool {
	return pos.X >= roomCoords.X && pos.X < roomCoords.X+roomCoords.W &&
		pos.Y >= roomCoords.Y && pos.Y < roomCoords.Y+roomCoords.H
}
