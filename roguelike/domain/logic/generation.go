package logic

import (
	"roguelike/domain/datastructs"
	"roguelike/domain/entity"
	"time"
)

func ClearData(level *entity.Level) {
	for room := 0; room < entity.ROOMS_NUM; room++ {
		level.Rooms[room].MonsterNumbers = 0
		level.Rooms[room].Consumables.FoodNumber = 0
		level.Rooms[room].Consumables.WeaponNumber = 0
		level.Rooms[room].Consumables.ElixirNumber = 0
		level.Rooms[room].Consumables.ScrollNumber = 0
	}
}

func GenerateNextRoom(level *entity.Level, player *entity.Player) {
	ClearData(level)
	level.LevelNumber++
	GenerateRooms(level.Rooms[:])
}

func GenerateRooms(room []entity.Room) {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		WidthRoom := GetRandomInRange(entity.MIN_ROOM_WIDTH, entity.MAX_ROOM_WIDTH)
		HeightRoom := GetRandomInRange(entity.MIN_ROOM_HEIGHT, entity.MAX_ROOM_HEIGHT)

		regionX := (i % entity.ROOMS_IN_WIDTH)
		regionY := (i / entity.ROOMS_IN_WIDTH)

		LeftRangeCoord := regionX*entity.REGION_WIDTH + 1
		RightRangeCoord := (regionX+1)*entity.REGION_WIDTH - WidthRoom - 1
		XCoord := GetRandomInRange(LeftRangeCoord, RightRangeCoord)

		UpRangeCoord := regionY*entity.REGION_HEIGHT + 1
		BottomRangeCoord := regionY*entity.REGION_HEIGHT - HeightRoom - 1
		YCoord := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

		room[i].Coordinates.W = WidthRoom
		room[i].Coordinates.H = HeightRoom

		room[i].Coordinates.XYcoords.X = XCoord
		room[i].Coordinates.XYcoords.Y = YCoord
	}
}

func GenerateEdgesForRooms(Edges []datastructs.Edge, EdgesCount *int) {
	*EdgesCount = 0

	for i := 0; i < entity.ROOMS_IN_HEIGHT; i++ {
		for j := 0; j+1 < entity.ROOMS_IN_WIDTH; j++ {
			CurrentRoom := i*entity.ROOMS_IN_WIDTH + j

			Edges[*EdgesCount].U = CurrentRoom
			Edges[*EdgesCount].V = CurrentRoom + 1

			*EdgesCount++
		}
	}

	for i := 0; i+1 < entity.ROOMS_IN_HEIGHT; i++ {
		for j := 0; j < entity.ROOMS_IN_WIDTH; j++ {
			CurrentRoom := i*entity.ROOMS_IN_WIDTH + j

			Edges[*EdgesCount].U = CurrentRoom
			Edges[*EdgesCount].V = CurrentRoom + entity.ROOMS_IN_WIDTH

			*EdgesCount++
		}
	}
}

func CreatePassage(XCoord, YCoord, Width, Height int, Passages *entity.Passages) {
	Passages.Passages = append(Passages.Passages, entity.Passage{})

	PassageCounter := Passages.PassagesNumber

	Passages.Passages[PassageCounter].XYcoords.X = XCoord - 1
	Passages.Passages[PassageCounter].XYcoords.Y = YCoord - 1

	Passages.Passages[PassageCounter].W = Width + 2
	Passages.Passages[PassageCounter].H = Height + 2

	Passages.PassagesNumber++
}

func GenerateHorizontalPassage(FirstRoom, SecondRoom int, room []entity.Room, passages *entity.Passages) {
	FirstCoords := room[FirstRoom].Coordinates
	SecondCoords := room[SecondRoom].Coordinates

	// правая стена первой комнаты
	FirstX := FirstCoords.XYcoords.X + FirstCoords.W - 1
	UpRangeCoord := FirstCoords.XYcoords.Y + 1
	BottomRangeCoord := FirstCoords.XYcoords.Y + FirstCoords.H - 2
	FirstY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	// левая стена второй комнаты
	SecondX := SecondCoords.XYcoords.X
	UpRangeCoord = SecondCoords.XYcoords.Y + 1
	BottomRangeCoord = SecondCoords.XYcoords.Y + SecondCoords.H - 2
	SecondY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstY == SecondY {
		// прямой коридор
		CreatePassage(FirstX, FirstY, Abs(SecondX-FirstX)+1, 1, passages)
	} else {
		Vertical := GetRandomInRange(Min(FirstX, SecondX)+1, Max(FirstX, SecondX)-1)
		// коридор с изгибом
		CreatePassage(FirstX, FirstY, Abs(Vertical-FirstX)+1, 1, passages)
		CreatePassage(Vertical, Min(FirstY, SecondY), 1, Abs(SecondY-FirstY)+1, passages)
		CreatePassage(Vertical, SecondY, Abs(SecondX-Vertical)+1, 1, passages)
	}
}

func GenerateVerticalPassages(FirstRoom, SecondRoom int, room []entity.Room, passages *entity.Passages) {
	FirstCoords := room[FirstRoom].Coordinates
	SecondCoords := room[SecondRoom].Coordinates

	FirstY := FirstCoords.XYcoords.Y + FirstCoords.H
	UpRangeCoord := FirstCoords.XYcoords.X + 1
	BottomRangeCoord := FirstCoords.XYcoords.X + FirstCoords.W - 2
	FirstX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	SecondY := SecondCoords.XYcoords.Y
	UpRangeCoord = SecondCoords.XYcoords.X + 1
	BottomRangeCoord = SecondCoords.XYcoords.X + SecondCoords.W - 2
	SecondX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstX == SecondX {
		// прямой коридор
		CreatePassage(FirstX, FirstY, 1, Abs(SecondY-FirstY)+1, passages)
	} else {
		Horizont := GetRandomInRange(Min(FirstY, SecondY)+1, Max(FirstY, SecondY)-1)
		// коридор с изгибом
		CreatePassage(FirstX, FirstY, 1, Abs(Horizont-FirstY)+1, passages)
		CreatePassage(min(FirstX, SecondX), Horizont, Abs(SecondX-FirstX)+1, 1, passages)
		CreatePassage(SecondX, Horizont, 1, Abs(SecondY-Horizont)+1, passages)
	}
}

func GeneratePassages(passages *entity.Passages, rooms []entity.Room) {
	// Создание массива ребер и получение его случайной перестановки
	passages.PassagesNumber = 0
	var countPassages int
	edges := make([]datastructs.Edge, entity.MAX_PASSAGES_NUM)
	GenerateEdgesForRooms(edges, &countPassages)
	ShuffleEdges(edges[:countPassages])

	// Коридоры между комнатами будут создаваться при помощи системы непересекающихся множеств
	parent := make([]int, entity.ROOMS_NUM)
	rank := make([]int, entity.ROOMS_NUM)
	datastructs.MakeSets(parent, rank)

	for i := 0; i < countPassages; i++ {
		if datastructs.FindSet(edges[i].U, parent) != datastructs.FindSet(edges[i].V, parent) {
			datastructs.UnionSets(edges[i].U, edges[i].V, parent, rank)

			if edges[i].U/entity.ROOMS_IN_WIDTH == edges[i].V/entity.ROOMS_IN_WIDTH {
				GenerateHorizontalPassage(edges[i].U, edges[i].V, rooms, passages)
			} else {
				GenerateVerticalPassages(edges[i].U, edges[i].V, rooms, passages)
			}
		}
	}
}

func GenerateCoordsOfEntity(room *entity.Room, coords *entity.Object) {
	UpperLeftX := room.Coordinates.XYcoords.X + 1
	UpperLeftY := room.Coordinates.XYcoords.Y + 1

	LowerRightX := room.Coordinates.XYcoords.X + room.Coordinates.W - 3
	LowerRightY := room.Coordinates.XYcoords.Y + room.Coordinates.H - 3

	coords.XYcoords.X = GetRandomInRange(UpperLeftX, LowerRightX)
	coords.XYcoords.Y = GetRandomInRange(UpperLeftY, LowerRightY)

	coords.W = 1
	coords.H = 1
}

func GeneratePlayer(rooms []entity.Room, player *entity.Player) int {
	PlayerRoom := GetRandomInRange(0, entity.ROOMS_NUM)
	GenerateCoordsOfEntity(&rooms[PlayerRoom], &player.BaseStats.Pos)
	return PlayerRoom
}

func GenerateExit(level *entity.Level, playerRoom int) {
	var exitRoom int

	for {
		// выбираем случайную комнату
		exitRoom = GetRandomInRange(0, entity.ROOMS_NUM-1)

		// нельзя в комнате игрока
		for exitRoom == playerRoom {
			exitRoom = GetRandomInRange(0, entity.ROOMS_NUM-1)
		}

		room := level.Rooms[exitRoom]

		// отступаем от стен
		upperLeftX := room.Coordinates.XYcoords.X + 2
		upperLeftY := room.Coordinates.XYcoords.Y + 2

		bottomRightX := upperLeftX + room.Coordinates.W - 5
		bottomRightY := upperLeftY + room.Coordinates.H - 5

		level.EndOfLevel.XYcoords.X = GetRandomInRange(upperLeftX, bottomRightX)
		level.EndOfLevel.XYcoords.Y = GetRandomInRange(upperLeftY, bottomRightY)

		level.EndOfLevel.W = 1
		level.EndOfLevel.H = 1

		// проверка, что место свободно
		if CheckUnoccupiedRoom(&level.EndOfLevel, &level.Rooms[exitRoom]) {
			break
		}
	}
}

func GenerateMonsterData(monster *entity.Monster, levelNum int) {
	monster.Type = entity.MonsterType(GetRandomInRange(0, 4))

	switch monster.Type {
	case entity.Zombie:
		monster.Hostility = entity.Medium
		monster.Stats.Agility = 25
		monster.Stats.Strength = 125
		monster.Stats.Health = 50

	case entity.Vampire:
		monster.Hostility = entity.High
		monster.Stats.Agility = 75
		monster.Stats.Strength = 125
		monster.Stats.Health = 50

	case entity.Ghost:
		monster.Hostility = entity.Low
		monster.Stats.Agility = 75
		monster.Stats.Strength = 25
		monster.Stats.Health = 75

	case entity.Ogre:
		monster.Hostility = entity.Medium
		monster.Stats.Agility = 25
		monster.Stats.Strength = 100
		monster.Stats.Health = 150

	case entity.Snake:
		monster.Hostility = entity.High
		monster.Stats.Agility = 100
		monster.Stats.Strength = 30
		monster.Stats.Health = 100
	}

	// Масштабирование сложности от уровня
	percentsUpdate := entity.PERCENTS_UPDATE_DIFFICULTY_MONSTERS * levelNum

	monster.Stats.Agility += monster.Stats.Agility * percentsUpdate / 100
	monster.Stats.Strength += monster.Stats.Strength * percentsUpdate / 100
	monster.Stats.Health += monster.Stats.Health * float64(percentsUpdate) / 100

	monster.IsChasing = false
	monster.Dir = entity.Stop
}

func GenerateMonsters(level *entity.Level, playerRoom int) {
	// Максимум монстров растёт с уровнем
	maxMonsters := entity.MAX_MONSTERS_PER_ROOM + level.LevelNumber/entity.LEVEL_UPDATE_DIFFICULTY

	for room := 0; room < entity.ROOMS_NUM; room++ {
		if room == playerRoom {
			continue
		}

		countMonsters := GetRandomInRange(0, maxMonsters)

		for i := 0; i < countMonsters; i++ {
			coords := &level.Rooms[room].Monsters[i].Stats.Pos

			for {
				GenerateCoordsOfEntity(&level.Rooms[room], coords)
				if CheckUnoccupiedRoom(coords, &level.Rooms[room]) {
					break
				}
			}

			GenerateMonsterData(&level.Rooms[room].Monsters[i], level.LevelNumber)
			level.Rooms[room].MonsterNumbers++
		}
	}
}

func GenerateFoodData(food *entity.Food, player *entity.Player) {
	names := [entity.CONSUMABLES_TYPE_MAX_NUM]string{
		"Ration of the Ironclad",
		"Crimson Berry Cluster",
		"Loaf of the Forgotten Baker",
		"Smoked Wyrm Jerky",
		"Golden Apple of Vitality",
		"Hardtack of the Endless March",
		"Spiced Venison Strips",
		"Honeyed Nectar Bread",
		"Dried Mushrooms of the Deep"}

	MaxRegen := player.BaseStats.Health * entity.MAX_PERCENT_FOOD_REGEN_FROM_HEALTH / 100
	food.ToRegen = GetRandomInRange(1, int(MaxRegen))
	food.Name = names[GetRandomInRange(0, entity.CONSUMABLES_TYPE_MAX_NUM-1)]
}

func GenerateFood(room *entity.Room, player *entity.Player) {
	CountFood := room.Consumables.FoodNumber
	Coords := &room.Consumables.RoomFood[CountFood].Geometry

	for {
		GenerateCoordsOfEntity(room, Coords)
		if CheckUnoccupiedRoom(Coords, room) {
			break
		}
	}

	GenerateFoodData(&room.Consumables.RoomFood[CountFood].Food, player)
	room.Consumables.FoodNumber++
}

func GenerateElixirData(elixir *entity.Elixir, player *entity.Player) {
	names := [entity.CONSUMABLES_TYPE_MAX_NUM]string{
		"Elixir of the Jade Serpent",
		"Potion of the Phantom's Breath",
		"Vial of Crimson Vitality",
		"Draught of the Frozen Star",
		"Elixir of the Shattered Mind",
		"Potion of the Wandering Soul",
		"Vial of Ember Essence",
		"Elixir of the Obsidian Veil",
		"Potion of the Howling Wind",
	}

	statType := entity.StatType(GetRandomInRange(0, int(entity.Strength)))
	var maxIncrease int

	switch statType {
	case entity.Health:
		maxIncrease = player.RegenLimit * entity.MAX_PERCENT_FOOD_REGEN_FROM_HEALTH / 100
	case entity.Agility:
		maxIncrease = player.BaseStats.Agility * entity.MAX_PERCENT_AGILITY_INCREASE / 100
	case entity.Strength:
		maxIncrease = player.BaseStats.Strength * entity.MAX_PERCENT_STRENGTH_INCREASE / 100
	}

	elixir.Stat = statType
	elixir.Increase = GetRandomInRange(1, maxIncrease)
	elixir.Duration = time.Duration(GetRandomInRange(entity.MIN_ELIXIR_DURATION_SECONDS, entity.MAX_ELIXIR_DURATION_SECONDS)) * time.Second
	elixir.Name = names[GetRandomInRange(0, entity.CONSUMABLES_TYPE_MAX_NUM-1)]
}

func GenerateElixir(room *entity.Room, player *entity.Player) {
	CountElixir := room.Consumables.ElixirNumber
	Coords := &room.Consumables.RoomElixir[CountElixir].Geometry

	for {
		GenerateCoordsOfEntity(room, Coords)
		if CheckUnoccupiedRoom(Coords, room) {
			break
		}
	}

	GenerateElixirData(&room.Consumables.RoomElixir[CountElixir].Elixir, player)
	room.Consumables.ElixirNumber++
}

func GenerateScrollData(scroll *entity.Scroll, player *entity.Player) {
	names := [entity.CONSUMABLES_TYPE_MAX_NUM]string{
		"Scroll of Shadowstep",
		"Parchment of Eternal Flame",
		"Manuscript of Forgotten Truths",
		"Scroll of Iron Will",
		"Vellum of the Void",
		"Scroll of Whispers",
		"Tome of the Lost King",
		"Scroll of Unseen Paths",
		"Parchment of Thunderous Roar",
	}

	statType := entity.StatType(GetRandomInRange(0, int(entity.Strength)))
	var maxIncrease int

	switch statType {
	case entity.Health:
		maxIncrease = player.RegenLimit * entity.MAX_PERCENT_FOOD_REGEN_FROM_HEALTH / 100
	case entity.Agility:
		maxIncrease = player.BaseStats.Agility * entity.MAX_PERCENT_AGILITY_INCREASE / 100
	case entity.Strength:
		maxIncrease = player.BaseStats.Strength * entity.MAX_PERCENT_STRENGTH_INCREASE / 100
	}

	scroll.Stat = statType
	scroll.Increase = GetRandomInRange(1, maxIncrease)
	scroll.Name = names[GetRandomInRange(0, entity.CONSUMABLES_TYPE_MAX_NUM-1)]
}

func GenerateScroll(room *entity.Room, player *entity.Player) {
	CountScroll := room.Consumables.ScrollNumber
	Coords := &room.Consumables.RoomScroll[CountScroll].Geometry

	for {
		GenerateCoordsOfEntity(room, Coords)
		if CheckUnoccupiedRoom(Coords, room) {
			break
		}
	}

	GenerateScrollData(&room.Consumables.RoomScroll[CountScroll].Scroll, player)
	room.Consumables.ScrollNumber++
}

func GenerateWeaponData(weapon *entity.Weapon, player *entity.Player) {
	names := [entity.CONSUMABLES_TYPE_MAX_NUM]string{
		"Blade of the Forgotten Dawn",
		"Obsidian Reaver",
		"Fang of the Shadow Wolf",
		"Ironclad Cleaver",
		"Crimson Talon",
		"Thunderstrike Maul",
		"Serpent's Kiss Dagger",
		"Voidrend Sword",
		"Ebonheart Spear",
	}

	maxStrength := entity.MAX_WEAPON_STRENGTH
	if player.Weapon.Strength < maxStrength && player.Weapon.Strength != entity.NO_WEAPON {
		maxStrength = player.Weapon.Strength
	}
	weapon.Strength = GetRandomInRange(entity.MIN_WEAPON_STRENGTH, maxStrength)
	weapon.Name = names[GetRandomInRange(0, entity.CONSUMABLES_TYPE_MAX_NUM-1)]
}

func GenerateWeapon(room *entity.Room, player *entity.Player) {
	CountWeapon := room.Consumables.WeaponNumber
	Coords := &room.Consumables.WeaponRoom[CountWeapon].Geometry

	for {
		GenerateCoordsOfEntity(room, Coords)
		if CheckUnoccupiedRoom(Coords, room) {
			break
		}
	}

	GenerateWeaponData(&room.Consumables.WeaponRoom[CountWeapon].Weapon, player)
	room.Consumables.WeaponNumber++
}

func GenerateConsumables(level *entity.Level, playerRoom int, player *entity.Player, levelNum int) {
	generateFuncs := []func(*entity.Room, *entity.Player){
		GenerateFood,
		GenerateElixir,
		GenerateScroll,
		GenerateWeapon,
	}

	maxConsumables := entity.MAX_CONSUMABLES_PER_ROOM - levelNum/entity.LEVEL_UPDATE_DIFFICULTY
	if maxConsumables < 1 {
		maxConsumables = 1
	}

	for room := 0; room < entity.ROOMS_NUM; room++ {
		if room == playerRoom {
			continue
		}

		countConsumables := GetRandomInRange(0, maxConsumables)
		for i := 0; i < countConsumables; i++ {
			consumableType := GetRandomInRange(0, entity.CONSUMABLES_TYPES_NUM-1)
			generateFuncs[consumableType](&level.Rooms[room], player)
		}
	}
}
