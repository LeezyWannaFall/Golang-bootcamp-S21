package characters

import (
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)

const OGRE_STEP =                2
const SIMPLE_DIRECTIONS =        4
const DIAGONAL_DIRECTIONS =      4
const ALL_DIRECTIONS =           8
const SIMPLE_TO_DIAGONAL_SHIFT = 4
const MAX_TRIES_TO_MOVE =       16

var DirectionDeltas = map[entity.Direction]entity.Pos{
	entity.Forward:  {X: 0, Y: -1},
	entity.Back:     {X: 0, Y: 1},
	entity.Left:     {X: -1, Y: 0},
	entity.Right:    {X: 1, Y: 0},

	entity.DiagonallyForwardLeft:  {X: -1, Y: -1},
	entity.DiagonallyForwardRight: {X: 1, Y: -1},
	entity.DiagonallyBackLeft:     {X: -1, Y: 1},
	entity.DiagonallyBackRight:    {X: 1, Y: 1},
}


func IsOutsideRoom(char entity.Object, room entity.Room) bool {
	return char.X < room.Coordinates.X ||
           char.X + char.W - 1 >= room.Coordinates.X+room.Coordinates.W ||
           char.Y < room.Coordinates.Y ||
           char.Y + char.H - 1 >= room.Coordinates.Y+room.Coordinates.H
}

func IsOutsideLevel(pos entity.Pos, level *entity.Level) bool {
	bounds := level.Coordinates
	return pos.X < bounds.X ||
        pos.X >= bounds.X + bounds.W ||
        pos.Y < bounds.Y ||
        pos.Y >= bounds.Y + bounds.H
}

func MoveCharacterByDirection(direction entity.Direction, characterGeometry *entity.Pos) {
	switch direction {
	case entity.Forward:
		characterGeometry.Y--
	case entity.Left:
		characterGeometry.X--
	case entity.Right:
		characterGeometry.X++
	case entity.Back:
		characterGeometry.Y++
	case entity.DiagonallyForwardLeft:
		characterGeometry.X--
		characterGeometry.Y--
	case entity.DiagonallyForwardRight:
		characterGeometry.X++
		characterGeometry.Y--
	case entity.DiagonallyBackLeft:
		characterGeometry.X--
		characterGeometry.Y++
	case entity.DiagonallyBackRight:
		characterGeometry.X++
		characterGeometry.Y++
	case entity.Stop:
		// Do nothing
	}
}


func MoveMonster(monster *entity.Monster, level *entity.Level, player *entity.Player) {
	switch monster.Type {
	case entity.Zombie:
		if IsPlayerNear(monster, player) {
			path := FindPathToPlayer(monster, level, *player)
			if len(path) > 1 {
				next := path[len(path) - 2]
				monster.Stats.Pos.X = next.X
				monster.Stats.Pos.Y = next.Y
			}
		} else {
			patternZombie(monster, level)
		}
	case entity.Vampire:
		if IsPlayerNear(monster, player) {
			path := FindPathToPlayer(monster, level, *player)
			if len(path) > 1 {
				next := path[len(path) - 2]
				monster.Stats.Pos.X = next.X
				monster.Stats.Pos.Y = next.Y
			}
		} else {
			patternVampire(monster, level)			
		}		
	case entity.Ghost:
		if IsPlayerNear(monster, player) {
			path := FindPathToPlayer(monster, level, *player)
			if len(path) > 1 {
				next := path[len(path) - 2]
				monster.Stats.Pos.X = next.X
				monster.Stats.Pos.Y = next.Y
			}
		} else {
			patternGhost(monster, level)
		}
	case entity.Ogre:
		if IsPlayerNear(monster, player) {
			path := FindPathToPlayer(monster, level, *player)
			if len(path) > 1 {
				next := path[len(path) - 2]
				monster.Stats.Pos.X = next.X
				monster.Stats.Pos.Y = next.Y
			}
		} else {
			patternOgre(monster, level)
		}
	case entity.Snake:
		if IsPlayerNear(monster, player) {
			path := FindPathToPlayer(monster, level, *player)
			if len(path) > 1 {
				next := path[len(path) - 2]
				monster.Stats.Pos.X = next.X
				monster.Stats.Pos.Y = next.Y
			}
		} else {
			patternSnake(monster, level)
		}
	}
}

func MovePlayer(player *entity.Player, level *entity.Level, direction entity.Direction) {
	coords := entity.Pos{
		X: player.BaseStats.Pos.X,
		Y: player.BaseStats.Pos.Y,
	}

	switch direction {
	case entity.Forward:
		player.BaseStats.Pos.Y--
	case entity.Back:
		player.BaseStats.Pos.Y++
	case entity.Left:
		player.BaseStats.Pos.X--
	case entity.Right:
		player.BaseStats.Pos.X++
	default:
		// Do nothing
	}

	newcoords := entity.Pos{
		X: player.BaseStats.Pos.X,
		Y: player.BaseStats.Pos.Y,
	}

	if IsOutsideLevel(newcoords, level) || !IsPassable(newcoords, level) {
		player.BaseStats.Pos.X = coords.X
		player.BaseStats.Pos.Y = coords.Y
	}
}

func PatternMonsters(monster *entity.Monster, level *entity.Level) {
	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		XYcoords := entity.Pos{
			X: monster.Stats.Pos.X,
			Y: monster.Stats.Pos.Y,
		}

		direction := entity.Direction(logic.GetRandomInRange(0, SIMPLE_DIRECTIONS - 1))
		MoveCharacterByDirection(direction, &XYcoords)

		if !IsOutsideLevel(XYcoords, level) && IsPassable(XYcoords, level) {
            // Если ход допустим, обновляем координаты монстра
            monster.Stats.Pos.X = XYcoords.X
			monster.Stats.Pos.Y = XYcoords.Y
            monster.Dir = direction
            return
        }
	}
}

func patternZombie(monster *entity.Monster, level *entity.Level) {
	PatternMonsters(monster, level)
}

func patternVampire(monster *entity.Monster, level *entity.Level) {
	PatternMonsters(monster, level)
}

func patternGhost(monster *entity.Monster, level *entity.Level) {
	PatternMonsters(monster, level)
}

func patternOgre(monster *entity.Monster, level *entity.Level) {
	for step := 0; step < OGRE_STEP; step++ {
		PatternMonsters(monster, level)
	}
}

func patternSnake(monster *entity.Monster, level *entity.Level) {
	PatternMonsters(monster, level)
}

func getAggroRadius(hostility entity.HostilityType) int {
    switch hostility {
    case entity.Low: return entity.LOW_HOSTILITY_RADIUS
    case entity.Medium: return entity.AVERAGE_HOSTILITY_RADIUS
    case entity.High: return entity.HIGH_HOSTILITY_RADIUS
    default: return 4
    }
}

func IsPlayerNear(monster *entity.Monster, player *entity.Player) bool {
    dx := logic.Abs(monster.Stats.Pos.X - player.BaseStats.Pos.X)
    dy := logic.Abs(monster.Stats.Pos.Y - player.BaseStats.Pos.Y)
    dist := logic.Max(dx, dy) // Chebyshev для 8-dir
    return dist <= getAggroRadius(monster.Hostility)
}

func FindPathToPlayer(monster *entity.Monster, level *entity.Level, player entity.Player) []entity.Pos {
	start := entity.Pos {
		X: monster.Stats.Pos.X,
		Y: monster.Stats.Pos.Y,
	}

	target := entity.Pos {
		X: player.BaseStats.Pos.X,
		Y: player.BaseStats.Pos.Y,
	}

	queue := []entity.Pos{start}
	visited := make(map[entity.Pos]bool)
	parent := make(map[entity.Pos]entity.Pos)

	visited[start] = true

	directions := []entity.Direction{
		entity.Forward,
		entity.Back,
		entity.Left,
		entity.Right,
		entity.DiagonallyForwardLeft,
		entity.DiagonallyForwardRight,
		entity.DiagonallyBackLeft,
		entity.DiagonallyBackRight,
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == target {
			break
		}

		for _, dir := range directions {
			delta := DirectionDeltas[dir]

			next := entity.Pos{
				X: current.X + delta.X,
				Y: current.Y + delta.Y,
			}

			if SkipNext(next, level, visited) {
				continue
			}

			visited[next] = true
			parent[next] = current
			queue = append(queue, next)
		}
	}

	if !visited[target] {
		return nil
	}

	path := []entity.Pos{}
	cur := target

	for cur != start {
		path = append(path, cur)
		cur = parent[cur]
	}

	path = append(path, start)

	return path
}

func IsPassable(pos entity.Pos, level *entity.Level) bool {
	if pos.X < level.Coordinates.X || pos.X >= level.Coordinates.X+level.Coordinates.W ||
    	pos.Y < level.Coordinates.Y || pos.Y >= level.Coordinates.Y+level.Coordinates.H {
    	return false
    }

	// Проверяем, внутри комнаты ли (и не стена)
    for _, room := range level.Rooms {
        r := room.Coordinates
        if pos.X > r.X && pos.X < r.X+r.W-1 && pos.Y > r.Y && pos.Y < r.Y+r.H-1 { // внутри, не граница
            return true
        }
    }

    // Проверяем, внутри прохода ли 
    for i := 0; i < level.Passages.PassagesNumber; i++ {
        p := level.Passages.Passages[i]
        if pos.X >= p.X && pos.X < p.X+p.W && pos.Y >= p.Y && pos.Y < p.Y+p.H { // внутри прохода
            return true
        }
    }

	return false
}

func SkipNext(next entity.Pos, level *entity.Level, visited map[entity.Pos]bool) bool {
	if next.X < level.Coordinates.X || next.X >= level.Coordinates.X + level.Coordinates.W ||
   		next.Y < level.Coordinates.Y || next.Y >= level.Coordinates.Y + level.Coordinates.H {
    return true
}

	if !IsPassable(next, level) {
		return true
	}

	if visited[next] {
		return true
	}

	return false
}

func FindCurrentRoom(pos entity.Object, level *entity.Level) *entity.Room {
	for i := 0; i < entity.ROOMS_NUM; i++ {
		room := &level.Rooms[i]
		if !IsOutsideRoom(pos, *room) {
			return room
		}
	}
	return nil
}