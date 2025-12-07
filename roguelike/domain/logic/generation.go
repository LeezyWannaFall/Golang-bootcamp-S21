package logic

import (
	"math/rand/v2"
	"roguelike/domain/entity"
)

type Edge struct {
	u int // first point of edge
	v int // second point of edge
}

/* ---------------Math Functions--------------------------*/
func GetRandomInRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
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

func GenerateHorizontalPassage(FirstRoom, SecondRoom int, Room []entity.Room, Passages *entity.Passages) {
	FirstCoords := Room[FirstRoom].Coordinates
	SecondCoords := Room[SecondRoom].Coordinates

    // Для первой комнаты фиксируется правая стена, поэтому X координата определяется однозначно
    // Для Y координаты определяются возможные значения и среди них выбирается случайное 
	FirstX := FirstCoords.X + FirstCoords.W - 1
	UpRangeCoord := FirstCoords.Y + 1
	BottomRangeCoord := FirstCoords.Y + FirstCoords.H - 2
	FirstY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	// Аналогично для второй комнаты с левой стеной
	SecondX := SecondCoords.X
	UpRangeCoord = SecondCoords.Y + 1
	BottomRangeCoord = SecondCoords.Y + SecondCoords.H - 2
	SecondY := GetRandomInRange(UpRangeCoord, BottomRangeCoord)

	// Если Y координаты равны, то строится прямой коридор, иначе с изгибом,
    // место которого выбирается случайно
	if FirstY == SecondY {
		CreatePassage(FirstX, FirstY, Abs(SecondX - FirstX) + 1, 1, Passages)
	} else {
		Vertical := GetRandomInRange(min(FirstX, SecondX) + 1, max(FirstX, SecondX) - 1)
		CreatePassage(FirstX, FirstY, Abs(Vertical - FirstX) + 1, 1, Passages)
		CreatePassage(Vertical, min(FirstY, SecondY), 1, Abs(SecondY - FirstY) + 1, Passages)
		CreatePassage(Vertical, SecondY, Abs(SecondX - Vertical) + 1, 1, Passages)
	}
}
