package logic

import (
	"roguelike/domain/entity"
)

type Edge struct {
	u int // first point of edge
	v int // second point of edge
}

/* ---------------Math Functions--------------------------*/
/*----------------------------------------------------------*/

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

		LeftRangeCoord := (i % entity.ROOMS_IN_WIDTH) * entity.REGION_WIDTH + 1
		RightRangeCoord := (i % entity.ROOMS_IN_WIDTH  + 1) * entity.REGION_WIDTH - WidthRoom - 1
		XCoord := GetRandomInRange(LeftRangeCoord, RightRangeCoord)

		UpRangeCoord := (i / entity.ROOMS_IN_WIDTH) * entity.REGION_HEIGHT + 1
		BottomRangeCoord := (i / entity.ROOMS_IN_WIDTH) * entity.REGION_HEIGHT - HeightRoom - 1
		YCoord := GetRandomInRange(UpRangeCoord, BottomRangeCoord)
		
		room[i].Coordinates.W = WidthRoom
		room[i].Coordinates.H = HeightRoom

		room[i].Coordinates.X = XCoord
		room[i].Coordinates.Y = YCoord
	}
}

func GenerateEdgesForRooms(Edges []Edge, EdgesCount *int) {
	*EdgesCount = 0

	for i := 0; i < entity.ROOMS_IN_HEIGHT; i++ {
		for j := 0; j + 1 < entity.ROOMS_IN_WIDTH; j++ {
			CurrentRoom := i * entity.ROOMS_IN_WIDTH + j
			
			Edges[*EdgesCount].u = CurrentRoom
			Edges[*EdgesCount].v = CurrentRoom + 1
			
			*EdgesCount++
		}	
	}

	for i := 0; i + 1 < entity.ROOMS_IN_HEIGHT; i++ {
		for j := 0; j < entity.ROOMS_IN_WIDTH; j++ {
			CurrentRoom := i * entity.ROOMS_IN_WIDTH + j

			Edges[*EdgesCount].u = CurrentRoom
			Edges[*EdgesCount].v = CurrentRoom + entity.ROOMS_IN_WIDTH

			*EdgesCount++
		}	
	}
}

func CreatePassage(XCoord, YCoord, Width, Height int, Passages *entity.Passages) {
	Passages.Passages = append(Passages.Passages, entity.Passage{})
	
	PassageCounter := Passages.PassagesNumber

	Passages.Passages[PassageCounter].X = XCoord -1 
	Passages.Passages[PassageCounter].Y = YCoord - 1

	Passages.Passages[PassageCounter].W = Width + 2
	Passages.Passages[PassageCounter].H = Height + 2

	Passages.PassagesNumber++
}

func GenerateHorizontalPassage(FirstRoom, SecondRoom int, room []entity.Room, passages *entity.Passages) {
	FirstCoords := room[FirstRoom].Coordinates
	SecondCoords := room[SecondRoom].Coordinates

	// правая стена первой комнаты
	FirstX := FirstCoords.X + FirstCoords.W - 1
	UpRangeCoord := FirstCoords.Y + 1
	BottomRangeCoord := FirstCoords.Y + FirstCoords.H - 2
	FirstY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	// левая стена второй комнаты
	SecondX := SecondCoords.X
	UpRangeCoord = SecondCoords.Y + 1
	BottomRangeCoord = SecondCoords.Y + SecondCoords.H - 2
	SecondY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstY == SecondY {
		// прямой коридор
		CreatePassage(FirstX, FirstY, Abs(SecondX - FirstX) + 1, 1, passages)
	} else {
		Vertical := GetRandomInRange(Min(FirstX, SecondX) + 1, Max(FirstX, SecondX) - 1)
		// коридор с изгибом
		CreatePassage(FirstX, FirstY, Abs(Vertical - FirstX) + 1, 1, passages)
		CreatePassage(Vertical, Min(FirstY, SecondY), 1, Abs(SecondY - FirstY) + 1, passages)
		CreatePassage(Vertical, SecondY, Abs(SecondX - Vertical) + 1, 1, passages)
	}
}

func GenerateVerticalPassages(FirstRoom, SecondRoom int, room []entity.Room, passages *entity.Passages) {
	FirstCoords := room[FirstRoom].Coordinates
	SecondCoords := room[SecondRoom].Coordinates

	FirstY := FirstCoords.Y + FirstCoords.H
	UpRangeCoord := FirstCoords.X + 1
	BottomRangeCoord := FirstCoords.X + FirstCoords.W - 2
	FirstX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	SecondY := SecondCoords.Y
	UpRangeCoord = SecondCoords.X + 1
	BottomRangeCoord = SecondCoords.X + SecondCoords.W - 2
	SecondX := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	if FirstX == SecondX {
		// прямой коридор
		CreatePassage(FirstX, FirstY, 1, Abs(SecondY - FirstY) + 1, passages)
	} else {
		Horizont := GetRandomInRange(Min(FirstY, SecondY) + 1, Max(FirstY, SecondY) -1)
		// коридор с изгибом
		CreatePassage(FirstX, FirstY, 1, Abs(Horizont - FirstY) + 1, passages)
		CreatePassage(min(FirstX, SecondX), Horizont, Abs(SecondX - FirstX) + 1, 1, passages)
		CreatePassage(SecondX, Horizont, 1, Abs(SecondY - Horizont) + 1, passages)
	}
}
