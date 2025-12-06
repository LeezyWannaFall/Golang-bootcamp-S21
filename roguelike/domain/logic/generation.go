package logic

import (
	"math/rand/v2"
	"roguelike/domain/entity"
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

func GetRandomInRange(min, max int) int {
	return rand.IntN(max-min+1) + min
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

func GenerateEdgesForRooms() {

}
