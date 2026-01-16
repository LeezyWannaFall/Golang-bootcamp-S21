package presentation

import (
	"roguelike/domain/entity"
)

func (r *Renderer) ApplyFogOfWar(level *entity.Level, player *entity.Player) {
	playerRoom := r.getRoomByCoord(&player.BaseStats.Pos, level.Rooms[:])

	for roomIdx := 0; roomIdx < entity.ROOMS_NUM; roomIdx++ {
		room := &level.Rooms[roomIdx]

		if int(roomIdx) != playerRoom {
			if r.mapState.VisibleRooms[roomIdx] {
				r.fillRoomByFog(room)
			}
		}

		if int(roomIdx) == playerRoom && playerRoom != -1 && r.isOutsideRoom(&player.BaseStats.Pos, &room.Coordinates) {
			if r.isPlayerNearRoom(player, room) {
				r.fillRoomByPartialFog(level, player, room)
			}
		}
	}

	for i := 0; i < level.Passages.PassagesNumber; i++ {
		if !r.mapState.VisiblePassages[i] {
			passage := &level.Passages.Passages[i]
			r.fillPassageByFog(passage)
		}
	}
}

func (r *Renderer) fillRoomByFog(room *entity.Room) {
	xRoom := room.Coordinates.XYcoords.X
	yRoom := room.Coordinates.XYcoords.Y
	xSize := room.Coordinates.W
	ySize := room.Coordinates.H

	for x := xRoom + 1; x < xRoom+xSize-1; x++ {
		for y := yRoom + 1; y < yRoom+ySize-1; y++ {
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				currentTile := r.mapState.Map[y][x]
				if currentTile != TileWallVertical && currentTile != TileWallHorizontal &&
					currentTile != TileDoor && currentTile != TileKey {
					r.mapState.Map[y][x] = TileFloor
				}
			}
		}
	}
}

func (r *Renderer) isPlayerNearRoom(player *entity.Player, room *entity.Room) bool {
	playerX := player.BaseStats.Pos.XYcoords.X
	playerY := player.BaseStats.Pos.XYcoords.Y

	xRoom := room.Coordinates.XYcoords.X
	yRoom := room.Coordinates.XYcoords.Y
	xSize := room.Coordinates.W
	ySize := room.Coordinates.H

	return (playerX >= xRoom-1 && playerX <= xRoom+xSize) &&
		(playerY >= yRoom-1 && playerY <= yRoom+ySize) &&
		r.isOutsideRoom(&player.BaseStats.Pos, &room.Coordinates)
}

func (r *Renderer) isVerticalDirectionFog(coords *entity.Object, roomCoords *entity.Object) bool {
	newCoords := *coords
	newCoords.XYcoords.X++
	if !r.isOutsideRoom(&newCoords, roomCoords) {
		return false
	}
	newCoords.XYcoords.X -= 2
	if !r.isOutsideRoom(&newCoords, roomCoords) {
		return false
	}
	return true
}

func (r *Renderer) fillRoomByPartialFog(level *entity.Level, player *entity.Player, room *entity.Room) {
	playerX := player.BaseStats.Pos.XYcoords.X
	playerY := player.BaseStats.Pos.XYcoords.Y

	xRoom := room.Coordinates.XYcoords.X
	yRoom := room.Coordinates.XYcoords.Y
	xSize := room.Coordinates.W
	ySize := room.Coordinates.H

	for y := yRoom + 1; y < yRoom+ySize-1; y++ {
		for x := xRoom + 1; x < xRoom+xSize-1; x++ {
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				if r.mapState.Map[y][x] == TileWallVertical || r.mapState.Map[y][x] == TileWallHorizontal {
					continue
				}

				if !r.CastRay(playerX, playerY, x, y, level) {
					if r.mapState.Map[y][x] != TileWallVertical && r.mapState.Map[y][x] != TileWallHorizontal &&
						r.mapState.Map[y][x] != TileDoor && r.mapState.Map[y][x] != TileKey {
						r.mapState.Map[y][x] = TileFloor
					}
				}
			}
		}
	}
}

func (r *Renderer) fillPassageByFog(passage *entity.Passage) {
	x1 := passage.XYcoords.X
	y1 := passage.XYcoords.Y
	xSize := passage.W
	ySize := passage.H

	for x := x1; x < x1+xSize; x++ {
		for y := y1; y < y1+ySize; y++ {
			if y >= 0 && y < r.mapState.MapHeight && x >= 0 && x < r.mapState.MapWidth {
				currentTile := r.mapState.Map[y][x]
				if currentTile != TileWallVertical && currentTile != TileWallHorizontal &&
					currentTile != TileDoor && currentTile != TileKey && currentTile != TileRoomOpening {
					r.mapState.Map[y][x] = TileFog
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (r *Renderer) BresenhamLine(x0, y0, x1, y1 int) []entity.Pos {
	var points []entity.Pos

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		points = append(points, entity.Pos{X: x, Y: y})

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}

	return points
}

func (r *Renderer) CastRay(playerX, playerY, targetX, targetY int, level *entity.Level) bool {
	points := r.BresenhamLine(playerX, playerY, targetX, targetY)

	for _, point := range points {
		if point.Y >= 0 && point.Y < r.mapState.MapHeight &&
			point.X >= 0 && point.X < r.mapState.MapWidth {
			tile := r.mapState.Map[point.Y][point.X]
			if tile == TileWallVertical || tile == TileWallHorizontal {
				return false
			}
		}
	}

	return true
}
