package presentation

import (
	"roguelike/domain/entity"
)

type MapState struct {
	Map             [][]rune
	ColorMap        [][]int
	VisibleRooms    [entity.ROOMS_NUM]bool
	VisiblePassages []bool
	MapHeight       int
	MapWidth        int
}

func NewMapState(height, width int) *MapState {
	mapState := &MapState{
		MapHeight:       height,
		MapWidth:        width,
		VisiblePassages: make([]bool, entity.MAX_PASSAGES_NUM),
	}
	mapState.Map = make([][]rune, height)
	mapState.ColorMap = make([][]int, height)
	for i := range mapState.Map {
		mapState.Map[i] = make([]rune, width)
		mapState.ColorMap[i] = make([]int, width)
		for j := range mapState.Map[i] {
			mapState.Map[i][j] = TileFog
			mapState.ColorMap[i][j] = WhiteFont
		}
	}
	return mapState
}

func (ms *MapState) ClearMap() {
	for i := 0; i < ms.MapHeight; i++ {
		for j := 0; j < ms.MapWidth; j++ {
			ms.Map[i][j] = TileFog
			ms.ColorMap[i][j] = WhiteFont
		}
	}
}

type Renderer struct {
	mapState *MapState
}

func NewRenderer(mapHeight, mapWidth int) *Renderer {
	return &Renderer{
		mapState: NewMapState(mapHeight, mapWidth),
	}
}

func (r *Renderer) GetMap() [][]rune {
	return r.mapState.Map
}

func (r *Renderer) GetColorMap() [][]int {
	return r.mapState.ColorMap
}

func (r *Renderer) ClearMap() {
	r.mapState.ClearMap()
	for i := 0; i < entity.ROOMS_NUM; i++ {
		r.mapState.VisibleRooms[i] = false
	}
	for i := 0; i < len(r.mapState.VisiblePassages); i++ {
		r.mapState.VisiblePassages[i] = false
	}
}

func (r *Renderer) GetVisibleRooms() [entity.ROOMS_NUM]bool {
	return r.mapState.VisibleRooms
}

func (r *Renderer) GetVisiblePassages() []bool {
	return r.mapState.VisiblePassages
}

func (r *Renderer) CreateNewMap(level *entity.Level, player *entity.Player) {
	r.mapState.ClearMap()
	r.roomsToMap(level, player)
	r.passagesToMap(level, player)
	r.consumablesToMap(level)
	r.ApplyFogOfWar(level, player)
	r.exitToMap(level)
	r.keysToMap(level)
	r.monstersToMap(level, player)
	r.doorsToMap(level)
	r.playerToMap(player)
}

func (r *Renderer) getRoomByCoord(coords *entity.Object, rooms []entity.Room) int {
	x := coords.XYcoords.X
	y := coords.XYcoords.Y

	for i := 0; i < len(rooms); i++ {
		room := &rooms[i]
		xRoom := room.Coordinates.XYcoords.X
		yRoom := room.Coordinates.XYcoords.Y
		xSize := room.Coordinates.W
		ySize := room.Coordinates.H

		if x >= xRoom && x < xRoom+xSize && y >= yRoom && y < yRoom+ySize {
			return i
		}
	}
	return -1
}

func (r *Renderer) isInsideRoom(coords *entity.Object, roomCoords *entity.Object) bool {
	x := coords.XYcoords.X
	y := coords.XYcoords.Y
	xRoom := roomCoords.XYcoords.X
	yRoom := roomCoords.XYcoords.Y
	xSize := roomCoords.W
	ySize := roomCoords.H

	return x >= xRoom && x < xRoom+xSize && y >= yRoom && y < yRoom+ySize
}

func (r *Renderer) isOutsideRoom(coords *entity.Object, roomCoords *entity.Object) bool {
	return !r.isInsideRoom(coords, roomCoords)
}

func (r *Renderer) isPassageAdjacentToRoom(passage *entity.Object, room *entity.Room) bool {
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

func (r *Renderer) roomsToMap(level *entity.Level, player *entity.Player) {
	playerRoom := r.getRoomByCoord(&player.BaseStats.Pos, level.Rooms[:])

	for i := 0; i < entity.ROOMS_NUM; i++ {
		if !r.mapState.VisibleRooms[i] && playerRoom != i {
			continue
		}

		room := &level.Rooms[i]
		x1 := room.Coordinates.XYcoords.X
		y1 := room.Coordinates.XYcoords.Y
		xSize := room.Coordinates.W
		ySize := room.Coordinates.H

		if xSize <= 0 || ySize <= 0 {
			continue
		}

		for y := 0; y < r.mapState.MapHeight; y++ {
			for x := 0; x < r.mapState.MapWidth; x++ {
				checkY := (y == y1 || y == y1+ySize-1) && (x1 <= x && x < x1+xSize)
				if checkY {
					r.mapState.Map[y][x] = TileWallHorizontal
				}
				checkX := (x == x1 || x == x1+xSize-1) && (y1 <= y && y < y1+ySize)
				if checkX {
					r.mapState.Map[y][x] = TileWallVertical
				}
			}
		}

		if playerRoom == i {
			for y := y1 + 1; y < y1+ySize-1; y++ {
				for x := x1 + 1; x < x1+xSize-1; x++ {
					if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
						if r.mapState.Map[y][x] != TileWallVertical && r.mapState.Map[y][x] != TileWallHorizontal {
							r.mapState.Map[y][x] = TileFloor
						}
					}
				}
			}
		}

		r.mapState.VisibleRooms[i] = true
	}
}

func (r *Renderer) passagesToMap(level *entity.Level, player *entity.Player) {
	for i := 0; i < level.Passages.PassagesNumber; i++ {
		passage := &level.Passages.Passages[i]
		passageObj := (*entity.Object)(passage)

		visible := true
		if !r.mapState.VisiblePassages[i] && r.isOutsideRoom(&player.BaseStats.Pos, passageObj) {
			visible = false
		}

		x1 := passage.XYcoords.X
		y1 := passage.XYcoords.Y
		xSize := passage.W
		ySize := passage.H

		for y := 0; y < r.mapState.MapHeight; y++ {
			for x := 0; x < r.mapState.MapWidth; x++ {
				checkX := (x1 < x && x < x1+xSize-1) && (y1 < y && y < y1+ySize-1)
				if checkX {
					coords := entity.Object{XYcoords: entity.Pos{X: x, Y: y}, W: 1, H: 1}
					roomIdx := r.getRoomByCoord(&coords, level.Rooms[:])
					if visible {
						if roomIdx == -1 {
							r.mapState.Map[y][x] = TileCorridor
						} else {
							r.mapState.Map[y][x] = TileRoomOpening
						}
					} else if roomIdx != -1 && r.mapState.VisibleRooms[roomIdx] {
						currentTile := r.mapState.Map[y][x]
						if currentTile != TileWallVertical && currentTile != TileWallHorizontal {
							r.mapState.Map[y][x] = TileRoomOpening
						}
					}
				}
			}
		}

		r.mapState.VisiblePassages[i] = visible
	}
}

func (r *Renderer) monstersToMap(level *entity.Level, player *entity.Player) {
	playerRoom := r.getRoomByCoord(&player.BaseStats.Pos, level.Rooms[:])

	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]

		for j := 0; j < room.MonsterNumbers; j++ {
			monster := &room.Monsters[j]

			if monster.Stats.Health <= 0 {
				continue
			}

			monsterInPlayerRoom := (playerRoom != -1 && playerRoom == i)
			monsterInSameRoomOrPassage := r.onTheSameRoomOrPassage(level, &player.BaseStats.Pos, &monster.Stats.Pos)

			if !monsterInPlayerRoom && !monsterInSameRoomOrPassage {
				continue
			}

			x := monster.Stats.Pos.XYcoords.X
			y := monster.Stats.Pos.XYcoords.Y

			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				currentTile := r.mapState.Map[y][x]
				if currentTile != TileWallVertical && currentTile != TileWallHorizontal &&
					currentTile != TileDoor {
					r.mapState.Map[y][x] = GetMonsterTile(monster.Type)
					r.mapState.ColorMap[y][x] = getMonsterColor(monster.Type)
				}
			}
		}
	}
}

func (r *Renderer) onTheSameRoomOrPassage(level *entity.Level, charCoords *entity.Object, monsterCoords *entity.Object) bool {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]
		charInside := r.isInsideRoom(charCoords, &room.Coordinates)
		monsterInside := r.isInsideRoom(monsterCoords, &room.Coordinates)
		if charInside && monsterInside {
			return true
		}
	}
	for i := 0; i < level.Passages.PassagesNumber; i++ {
		passage := (*entity.Object)(&level.Passages.Passages[i])
		charInside := r.isInsideRoom(charCoords, passage)
		monsterInside := r.isInsideRoom(monsterCoords, passage)
		if charInside && monsterInside {
			return true
		}
	}
	return false
}

func getMonsterColor(monsterType entity.MonsterType) int {
	switch monsterType {
	case entity.Zombie:
		return GreenFont
	case entity.Vampire:
		return RedFont
	case entity.Ghost:
		return WhiteFont
	case entity.Ogre:
		return YellowFont
	case entity.Snake:
		return WhiteFont
	case entity.Mimic:
		return WhiteFont
	default:
		return WhiteFont
	}
}

func getColorByStatType(stat entity.StatType) int {
	switch stat {
	case entity.Health:
		return RedFont
	case entity.Strength:
		return BlueFont
	case entity.Agility:
		return GreenFont
	default:
		return WhiteFont
	}
}

func (r *Renderer) consumablesToMap(level *entity.Level) {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		if !r.mapState.VisibleRooms[i] {
			continue
		}

		room := &level.Rooms[i]

		for j := 0; j < room.Consumables.FoodNumber; j++ {
			food := room.Consumables.RoomFood[j]
			x := food.Geometry.XYcoords.X
			y := food.Geometry.XYcoords.Y
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				r.mapState.Map[y][x] = TileFood
				r.mapState.ColorMap[y][x] = RedFont
			}
		}

		for j := 0; j < room.Consumables.ElixirNumber; j++ {
			elixir := room.Consumables.RoomElixir[j]
			x := elixir.Geometry.XYcoords.X
			y := elixir.Geometry.XYcoords.Y
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				r.mapState.Map[y][x] = TileElixir
				r.mapState.ColorMap[y][x] = getColorByStatType(elixir.Elixir.Stat)
			}
		}

		for j := 0; j < room.Consumables.ScrollNumber; j++ {
			scroll := room.Consumables.RoomScroll[j]
			x := scroll.Geometry.XYcoords.X
			y := scroll.Geometry.XYcoords.Y
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				r.mapState.Map[y][x] = TileScroll
				r.mapState.ColorMap[y][x] = getColorByStatType(scroll.Scroll.Stat)
			}
		}

		for j := 0; j < room.Consumables.WeaponNumber; j++ {
			weapon := room.Consumables.WeaponRoom[j]
			x := weapon.Geometry.XYcoords.X
			y := weapon.Geometry.XYcoords.Y
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				r.mapState.Map[y][x] = TileWeapon
				r.mapState.ColorMap[y][x] = BlueFont
			}
		}
	}
}

func (r *Renderer) exitToMap(level *entity.Level) {
	exitRoom := r.getRoomByCoord(&level.EndOfLevel, level.Rooms[:])
	x := level.EndOfLevel.XYcoords.X
	y := level.EndOfLevel.XYcoords.Y

	if exitRoom < 0 {
		return
	}

	if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
		room := &level.Rooms[exitRoom]
		if r.isInsideRoom(&level.EndOfLevel, &room.Coordinates) {
			if r.mapState.VisibleRooms[exitRoom] {
				if r.mapState.Map[y][x] != TileFog {
					r.mapState.Map[y][x] = TileExit
					r.mapState.ColorMap[y][x] = YellowFont
				}
			}
		}
	}
}

func (r *Renderer) playerToMap(player *entity.Player) {
	x := player.BaseStats.Pos.XYcoords.X
	y := player.BaseStats.Pos.XYcoords.Y
	if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
		r.mapState.Map[y][x] = TilePlayer
		r.mapState.ColorMap[y][x] = CyanFont
	}
}

func (r *Renderer) keysToMap(level *entity.Level) {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]

		for j := 0; j < room.Consumables.KeyNumber; j++ {
			key := room.Consumables.RoomKeys[j]
			x := key.Geometry.XYcoords.X
			y := key.Geometry.XYcoords.Y
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				if r.mapState.Map[y][x] != TileFog {
					r.mapState.Map[y][x] = TileKey
					r.mapState.ColorMap[y][x] = getKeyColor(key.Key.Color)
				}
			}
		}
	}
}

func (r *Renderer) doorsToMap(level *entity.Level) {
	for i := 0; i < level.DoorNumber; i++ {
		door := &level.Doors[i]
		x := door.Position.XYcoords.X
		y := door.Position.XYcoords.Y

		if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
			if !door.IsOpen {
				currentTile := r.mapState.Map[y][x]
				doorRoom := r.getRoomByCoord(&door.Position, level.Rooms[:])
				isVisible := doorRoom >= 0 && r.mapState.VisibleRooms[doorRoom]

				if currentTile != TilePlayer &&
					(currentTile == TileCorridor || currentTile == TileRoomOpening || currentTile == TileFloor || currentTile == TileDoor || (currentTile == TileFog && isVisible)) {
					r.mapState.Map[y][x] = TileDoor
					r.mapState.ColorMap[y][x] = getKeyColor(door.Color)
				}
			}
		}
	}
}

func getKeyColor(color entity.KeyColor) int {
	switch color {
	case entity.RedKey:
		return RedFont
	case entity.BlueKey:
		return BlueFont
	case entity.YellowKey:
		return YellowFont
	case entity.GreenKey:
		return GreenFont
	default:
		return WhiteFont
	}
}
