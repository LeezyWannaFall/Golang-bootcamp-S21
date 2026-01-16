package logic

import (
	"roguelike/domain/datastructs"
	"roguelike/domain/entity"
	"time"
)

func ClearData(level *entity.Level) {
	level.Coordinates = entity.Object{
		XYcoords: entity.Pos{X: 0, Y: 0},
		W:        entity.ROOMS_IN_WIDTH * entity.REGION_WIDTH,
		H:        entity.ROOMS_IN_HEIGHT * entity.REGION_HEIGHT,
	}

	for room := 0; room < entity.ROOMS_NUM; room++ {
		level.Rooms[room].MonsterNumbers = 0
		level.Rooms[room].Consumables.FoodNumber = 0
		level.Rooms[room].Consumables.WeaponNumber = 0
		level.Rooms[room].Consumables.ElixirNumber = 0
		level.Rooms[room].Consumables.ScrollNumber = 0
		level.Rooms[room].Consumables.KeyNumber = 0
	}
	level.Passages.Passages = make([]entity.Passage, 0)
	level.Passages.PassagesNumber = 0
	level.Doors = make([]entity.Door, 0, entity.MAX_DOORS_PER_LEVEL)
	level.DoorNumber = 0
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
		if RightRangeCoord <= LeftRangeCoord {
			RightRangeCoord = LeftRangeCoord + 1
		}
		if RightRangeCoord < LeftRangeCoord {
			RightRangeCoord = LeftRangeCoord
		}
		XCoord := GetRandomInRange(LeftRangeCoord, RightRangeCoord)

		UpRangeCoord := regionY*entity.REGION_HEIGHT + 1
		BottomRangeCoord := (regionY+1)*entity.REGION_HEIGHT - HeightRoom - 1
		if BottomRangeCoord <= UpRangeCoord {
			BottomRangeCoord = UpRangeCoord + 1
		}
		if BottomRangeCoord < UpRangeCoord {
			BottomRangeCoord = UpRangeCoord
		}
		YCoord := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

		if WidthRoom < entity.MIN_ROOM_WIDTH {
			WidthRoom = entity.MIN_ROOM_WIDTH
		}
		if HeightRoom < entity.MIN_ROOM_HEIGHT {
			HeightRoom = entity.MIN_ROOM_HEIGHT
		}

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

	FirstX := FirstCoords.XYcoords.X + FirstCoords.W - 1
	UpRangeCoord := FirstCoords.XYcoords.Y + 1
	BottomRangeCoord := FirstCoords.XYcoords.Y + FirstCoords.H - 2
	if BottomRangeCoord < UpRangeCoord {
		BottomRangeCoord = UpRangeCoord
	}
	FirstY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	SecondX := SecondCoords.XYcoords.X
	UpRangeCoord = SecondCoords.XYcoords.Y + 1
	BottomRangeCoord = SecondCoords.XYcoords.Y + SecondCoords.H - 2
	if BottomRangeCoord < UpRangeCoord {
		BottomRangeCoord = UpRangeCoord
	}
	SecondY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstY == SecondY {
		CreatePassage(FirstX, FirstY, Abs(SecondX-FirstX)+1, 1, passages)
	} else {
		Vertical := GetRandomInRange(Min(FirstX, SecondX)+1, Max(FirstX, SecondX)-1)
		CreatePassage(FirstX, FirstY, Abs(Vertical-FirstX)+1, 1, passages)
		CreatePassage(Vertical, Min(FirstY, SecondY), 1, Abs(SecondY-FirstY)+1, passages)
		CreatePassage(Vertical, SecondY, Abs(SecondX-Vertical)+1, 1, passages)
	}
}

func GenerateVerticalPassages(FirstRoom, SecondRoom int, room []entity.Room, passages *entity.Passages) {
	FirstCoords := room[FirstRoom].Coordinates
	SecondCoords := room[SecondRoom].Coordinates

	FirstY := FirstCoords.XYcoords.Y + FirstCoords.H - 1
	UpRangeCoord := FirstCoords.XYcoords.X + 1
	BottomRangeCoord := FirstCoords.XYcoords.X + FirstCoords.W - 2
	if BottomRangeCoord < UpRangeCoord {
		BottomRangeCoord = UpRangeCoord
	}
	FirstX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	SecondY := SecondCoords.XYcoords.Y
	UpRangeCoord = SecondCoords.XYcoords.X + 1
	BottomRangeCoord = SecondCoords.XYcoords.X + SecondCoords.W - 2
	if BottomRangeCoord < UpRangeCoord {
		BottomRangeCoord = UpRangeCoord
	}
	SecondX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstX == SecondX {
		CreatePassage(FirstX, FirstY, 1, Abs(SecondY-FirstY)+1, passages)
	} else {
		Horizont := GetRandomInRange(Min(FirstY, SecondY)+1, Max(FirstY, SecondY)-1)
		CreatePassage(FirstX, FirstY, 1, Abs(Horizont-FirstY)+1, passages)
		CreatePassage(Min(FirstX, SecondX), Horizont, Abs(SecondX-FirstX)+1, 1, passages)
		CreatePassage(SecondX, Horizont, 1, Abs(SecondY-Horizont)+1, passages)
	}
}

func GeneratePassages(passages *entity.Passages, rooms []entity.Room) {
	passages.PassagesNumber = 0
	var countPassages int
	edges := make([]datastructs.Edge, entity.MAX_PASSAGES_NUM)
	GenerateEdgesForRooms(edges, &countPassages)
	ShuffleEdges(edges[:countPassages])

	// Используется DSU для гарантии связности графа
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

	LowerRightX := UpperLeftX + room.Coordinates.W - 3
	LowerRightY := UpperLeftY + room.Coordinates.H - 3

	if LowerRightX < UpperLeftX {
		LowerRightX = UpperLeftX
	}
	if LowerRightY < UpperLeftY {
		LowerRightY = UpperLeftY
	}

	coords.XYcoords.X = GetRandomInRange(UpperLeftX, LowerRightX)
	coords.XYcoords.Y = GetRandomInRange(UpperLeftY, LowerRightY)

	coords.W = 1
	coords.H = 1
}

func GeneratePlayer(rooms []entity.Room, player *entity.Player) int {
	PlayerRoom := GetRandomInRange(0, entity.ROOMS_NUM-1)
	GenerateCoordsOfEntity(&rooms[PlayerRoom], &player.BaseStats.Pos)
	return PlayerRoom
}

func GenerateExit(level *entity.Level, playerRoom int) {
	var exitRoom int

	for {
		exitRoom = GetRandomInRange(0, entity.ROOMS_NUM-1)

		for exitRoom == playerRoom {
			exitRoom = GetRandomInRange(0, entity.ROOMS_NUM-1)
		}

		room := level.Rooms[exitRoom]

		upperLeftX := room.Coordinates.XYcoords.X + 2
		upperLeftY := room.Coordinates.XYcoords.Y + 2

		bottomRightX := upperLeftX + room.Coordinates.W - 5
		bottomRightY := upperLeftY + room.Coordinates.H - 5

		level.EndOfLevel.XYcoords.X = GetRandomInRange(upperLeftX, bottomRightX)
		level.EndOfLevel.XYcoords.Y = GetRandomInRange(upperLeftY, bottomRightY)

		level.EndOfLevel.W = 1
		level.EndOfLevel.H = 1

		if CheckUnoccupiedRoom(&level.EndOfLevel, &level.Rooms[exitRoom]) {
			break
		}
	}
}

func GenerateMonsterData(monster *entity.Monster, levelNum int, balanceAdjustment int) {
	monster.Type = entity.MonsterType(GetRandomInRange(0, 5))

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

	case entity.Mimic:
		monster.Hostility = entity.Low
		monster.Stats.Agility = 100
		monster.Stats.Strength = 30
		monster.Stats.Health = 120
	}

	percentsUpdate := entity.PERCENTS_UPDATE_DIFFICULTY_MONSTERS * levelNum
	percentsUpdate += balanceAdjustment

	if percentsUpdate < 0 {
		percentsUpdate = 0
	}

	monster.Stats.Agility += monster.Stats.Agility * percentsUpdate / 100
	monster.Stats.Strength += monster.Stats.Strength * percentsUpdate / 100
	monster.Stats.Health += monster.Stats.Health * float64(percentsUpdate) / 100

	monster.IsChasing = false
	monster.Dir = entity.Stop
}

func GenerateMonsters(level *entity.Level, playerRoom int, balance BalanceAdjustment) {
	maxMonsters := entity.MAX_MONSTERS_PER_ROOM + level.LevelNumber/entity.LEVEL_UPDATE_DIFFICULTY
	maxMonsters += balance.MonsterCount
	if maxMonsters < 0 {
		maxMonsters = 0
	}
	if maxMonsters > entity.MAX_MONSTERS_PER_ROOM+3 {
		maxMonsters = entity.MAX_MONSTERS_PER_ROOM + 3
	}

	for room := 0; room < entity.ROOMS_NUM; room++ {
		if room == playerRoom {
			continue
		}

		countMonsters := GetRandomInRange(1, maxMonsters)
		if countMonsters > entity.MAX_MONSTERS_PER_ROOM {
			countMonsters = entity.MAX_MONSTERS_PER_ROOM
		}

		for i := 0; i < countMonsters; i++ {
			coords := &level.Rooms[room].Monsters[i].Stats.Pos

			for {
				GenerateCoordsOfEntity(&level.Rooms[room], coords)
				if CheckUnoccupiedRoom(coords, &level.Rooms[room]) {
					break
				}
			}

			GenerateMonsterData(&level.Rooms[room].Monsters[i], level.LevelNumber, balance.MonsterDifficulty)
			level.Rooms[room].MonsterNumbers++
		}
	}
}

func GenerateFoodData(food *entity.Food, player *entity.Player, foodBonus int) {
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
	MaxRegen = MaxRegen + MaxRegen*float64(foodBonus)/100
	if MaxRegen < 1 {
		MaxRegen = 1
	}
	food.ToRegen = GetRandomInRange(1, int(MaxRegen))
	food.Name = names[GetRandomInRange(0, entity.CONSUMABLES_TYPE_MAX_NUM-1)]
}

func GenerateFood(room *entity.Room, player *entity.Player, foodBonus int) {
	CountFood := room.Consumables.FoodNumber
	Coords := &room.Consumables.RoomFood[CountFood].Geometry

	for {
		GenerateCoordsOfEntity(room, Coords)
		if CheckUnoccupiedRoom(Coords, room) {
			break
		}
	}

	GenerateFoodData(&room.Consumables.RoomFood[CountFood].Food, player, foodBonus)
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

func GenerateConsumables(level *entity.Level, playerRoom int, player *entity.Player, levelNum int, balance BalanceAdjustment) {
	maxConsumables := entity.MAX_CONSUMABLES_PER_ROOM - levelNum/entity.LEVEL_UPDATE_DIFFICULTY
	maxConsumables += balance.ConsumableCount
	if maxConsumables < 1 {
		maxConsumables = 1
	}
	if maxConsumables > entity.MAX_CONSUMABLES_PER_ROOM+2 {
		maxConsumables = entity.MAX_CONSUMABLES_PER_ROOM + 2
	}

	foodWeight := 30
	if balance.FoodBonus > 0 {
		foodWeight = 50
	}

	for room := 0; room < entity.ROOMS_NUM; room++ {
		if room == playerRoom {
			continue
		}

		countConsumables := GetRandomInRange(0, maxConsumables)
		for i := 0; i < countConsumables; i++ {
			consumableType := GetRandomInRange(0, entity.CONSUMABLES_TYPES_NUM-1)

			if balance.FoodBonus > 0 && GetRandomInRange(0, 100) < foodWeight {
				consumableType = 0
			}

			switch consumableType {
			case 0:
				GenerateFood(&level.Rooms[room], player, balance.FoodBonus)
			case 1:
				GenerateElixir(&level.Rooms[room], player)
			case 2:
				GenerateScroll(&level.Rooms[room], player)
			case 3:
				GenerateWeapon(&level.Rooms[room], player)
			}
		}
	}
}

func GenerateDoorsAndKeys(level *entity.Level, playerRoom int) {
	generateDoorsAndKeysRecursive(level, playerRoom, 0)
}

func generateDoorsAndKeysRecursive(level *entity.Level, playerRoom int, attempts int) {
	if attempts > 20 {
		return
	}

	doors := make([]entity.Door, 0)
	keyColors := []entity.KeyColor{entity.RedKey, entity.BlueKey, entity.YellowKey, entity.GreenKey}

	usedColors := make(map[entity.KeyColor]bool)
	roomPassageCount := make(map[int]int)
	roomDoorCount := make(map[int]int)

	for i := 0; i < level.Passages.PassagesNumber; i++ {
		passage := &level.Passages.Passages[i]
		adjacentRooms := getAdjacentRooms(passage, level)
		for _, roomIdx := range adjacentRooms {
			roomPassageCount[roomIdx]++
		}
	}

	playerRoomPassageCount := 0
	for i := 0; i < level.Passages.PassagesNumber; i++ {
		passage := &level.Passages.Passages[i]
		adjacentRooms := getAdjacentRooms(passage, level)
		for _, roomIdx := range adjacentRooms {
			if roomIdx == playerRoom {
				playerRoomPassageCount++
				break
			}
		}
	}

	playerRoomOpenPassages := 0

	for i := 0; i < level.Passages.PassagesNumber; i++ {
		passage := &level.Passages.Passages[i]

		adjacentRooms := getAdjacentRooms(passage, level)
		if len(adjacentRooms) < 2 {
			continue
		}

		canPlaceDoor := true
		hasPlayerRoom := false
		for _, roomIdx := range adjacentRooms {
			if roomIdx == playerRoom {
				hasPlayerRoom = true
				break
			}
		}

		if !hasPlayerRoom {
			for _, roomIdx := range adjacentRooms {
				if roomDoorCount[roomIdx] >= roomPassageCount[roomIdx] {
					canPlaceDoor = false
					break
				}
			}
		} else {
			if playerRoomOpenPassages == 0 {
				canPlaceDoor = false
			} else if playerRoomOpenPassages >= playerRoomPassageCount {
				canPlaceDoor = false
			} else {
				for _, roomIdx := range adjacentRooms {
					if roomDoorCount[roomIdx] >= roomPassageCount[roomIdx] {
						canPlaceDoor = false
						break
					}
				}
			}
		}

		if GetRandomInRange(0, 100) < 70 && canPlaceDoor {
			color := keyColors[GetRandomInRange(0, len(keyColors)-1)]
			usedColors[color] = true

			doorPos := findDoorPosition(passage)
			if doorPos.XYcoords.X >= 0 && doorPos.XYcoords.Y >= 0 {
				doors = append(doors, entity.Door{
					Position: doorPos,
					Color:    color,
					IsOpen:   false,
				})
				for _, roomIdx := range adjacentRooms {
					roomDoorCount[roomIdx]++
				}
				if hasPlayerRoom {
					playerRoomOpenPassages++
				}
			}
		} else if hasPlayerRoom {
			playerRoomOpenPassages++
		}
	}

	level.Doors = doors
	level.DoorNumber = len(doors)

	if level.DoorNumber == 0 && attempts < 5 {
		generateDoorsAndKeysRecursive(level, playerRoom, attempts+1)
		return
	}

	keyRooms := make([]int, 0)
	if len(usedColors) > 0 {
		accessibleRooms := getAccessibleRooms(level, playerRoom)
		accessibleRoomsList := make([]int, 0)
		for i := 0; i < entity.ROOMS_NUM; i++ {
			if i != playerRoom && accessibleRooms[i] {
				accessibleRoomsList = append(accessibleRoomsList, i)
			}
		}

		if len(accessibleRoomsList) == 0 {
			if attempts < 20 {
				level.Doors = make([]entity.Door, 0)
				level.DoorNumber = 0
				generateDoorsAndKeysRecursive(level, playerRoom, attempts+1)
				return
			}
		}

		for color := range usedColors {
			var room int
			if len(accessibleRoomsList) > 0 {
				room = accessibleRoomsList[GetRandomInRange(0, len(accessibleRoomsList)-1)]
			} else {
				room = GetRandomInRange(0, entity.ROOMS_NUM-1)
				for room == playerRoom {
					room = GetRandomInRange(0, entity.ROOMS_NUM-1)
				}
			}
			keyRooms = append(keyRooms, room)
			GenerateKey(&level.Rooms[room], color)
		}
	}

	if len(keyRooms) > 0 && !ValidateKeyAccessibility(level, playerRoom, keyRooms) {
		if attempts < 20 {
			for _, roomIdx := range keyRooms {
				level.Rooms[roomIdx].Consumables.KeyNumber = 0
			}
			level.Doors = make([]entity.Door, 0)
			level.DoorNumber = 0
			generateDoorsAndKeysRecursive(level, playerRoom, attempts+1)
			return
		}
	}
}

func findDoorPosition(passage *entity.Passage) entity.Object {
	midX := passage.XYcoords.X + passage.W/2
	midY := passage.XYcoords.Y + passage.H/2

	return entity.Object{
		XYcoords: entity.Pos{X: midX, Y: midY},
		W:        1,
		H:        1,
	}
}

func GenerateKey(room *entity.Room, color entity.KeyColor) {
	if room.Consumables.KeyNumber >= int(entity.KeyColorCount) {
		return
	}

	coords := &room.Consumables.RoomKeys[room.Consumables.KeyNumber].Geometry

	for {
		GenerateCoordsOfEntity(room, coords)
		if CheckUnoccupiedRoom(coords, room) {
			break
		}
	}

	room.Consumables.RoomKeys[room.Consumables.KeyNumber].Key.Color = color
	room.Consumables.KeyNumber++
}

func getAccessibleRooms(level *entity.Level, startRoom int) map[int]bool {
	accessible := make(map[int]bool)
	accessible[startRoom] = true

	queue := []int{startRoom}
	keys := make(map[entity.KeyColor]bool)

	for len(queue) > 0 {
		currentRoom := queue[0]
		queue = queue[1:]

		for i := 0; i < level.Rooms[currentRoom].Consumables.KeyNumber; i++ {
			keyColor := level.Rooms[currentRoom].Consumables.RoomKeys[i].Key.Color
			keys[keyColor] = true
		}

		for i := 0; i < level.Passages.PassagesNumber; i++ {
			passage := &level.Passages.Passages[i]

			adjacentRooms := getAdjacentRooms(passage, level)
			if len(adjacentRooms) < 2 {
				continue
			}

			room1, room2 := adjacentRooms[0], adjacentRooms[1]

			canPass := true
			for j := 0; j < level.DoorNumber; j++ {
				door := &level.Doors[j]
				if isDoorInPassage(door, passage) {
					canPass = door.IsOpen || keys[door.Color]
					break
				}
			}

			if !canPass {
				continue
			}

			if room1 == currentRoom && !accessible[room2] {
				accessible[room2] = true
				queue = append(queue, room2)
			} else if room2 == currentRoom && !accessible[room1] {
				accessible[room1] = true
				queue = append(queue, room1)
			}
		}
	}

	return accessible
}

func ValidateKeyAccessibility(level *entity.Level, startRoom int, keyRooms []int) bool {
	accessible := getAccessibleRooms(level, startRoom)

	for _, keyRoom := range keyRooms {
		if !accessible[keyRoom] {
			return false
		}
	}

	exitRoom := -1
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]
		if level.EndOfLevel.XYcoords.X >= room.Coordinates.XYcoords.X &&
			level.EndOfLevel.XYcoords.X < room.Coordinates.XYcoords.X+room.Coordinates.W &&
			level.EndOfLevel.XYcoords.Y >= room.Coordinates.XYcoords.Y &&
			level.EndOfLevel.XYcoords.Y < room.Coordinates.XYcoords.Y+room.Coordinates.H {
			exitRoom = i
			break
		}
	}

	if exitRoom >= 0 && !accessible[exitRoom] {
		return false
	}

	return true
}

func getAdjacentRooms(passage *entity.Passage, level *entity.Level) []int {
	rooms := make([]int, 0)

	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]
		if isPassageAdjacentToRoomForDoors(passage, room) {
			rooms = append(rooms, i)
		}
	}

	return rooms
}

func isPassageAdjacentToRoomForDoors(passage *entity.Passage, room *entity.Room) bool {
	px1 := passage.XYcoords.X
	py1 := passage.XYcoords.Y
	px2 := px1 + passage.W
	py2 := py1 + passage.H

	rx1 := room.Coordinates.XYcoords.X
	ry1 := room.Coordinates.XYcoords.Y
	rx2 := rx1 + room.Coordinates.W
	ry2 := ry1 + room.Coordinates.H

	return (px2 >= rx1 && px1 <= rx2 && py2 >= ry1 && py1 <= ry2)
}

func isDoorInPassage(door *entity.Door, passage *entity.Passage) bool {
	dx := door.Position.XYcoords.X
	dy := door.Position.XYcoords.Y

	px1 := passage.XYcoords.X
	py1 := passage.XYcoords.Y
	px2 := px1 + passage.W
	py2 := py1 + passage.H

	return dx >= px1 && dx < px2 && dy >= py1 && dy < py2
}
