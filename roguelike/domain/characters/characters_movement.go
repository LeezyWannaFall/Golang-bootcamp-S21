package characters

import (
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)

const OGRE_STEP = 2
const SIMPLE_DIRECTIONS = 4
const DIAGONAL_DIRECTIONS = 4
const ALL_DIRECTIONS = 8
const SIMPLE_TO_DIAGONAL_SHIFT = 4
const MAX_TRIES_TO_MOVE = 16

var DirectionDeltas = map[entity.Direction]entity.Pos{
	entity.Forward: {X: 0, Y: -1},
	entity.Back:    {X: 0, Y: 1},
	entity.Left:    {X: -1, Y: 0},
	entity.Right:   {X: 1, Y: 0},

	entity.DiagonallyForwardLeft:  {X: -1, Y: -1},
	entity.DiagonallyForwardRight: {X: 1, Y: -1},
	entity.DiagonallyBackLeft:     {X: -1, Y: 1},
	entity.DiagonallyBackRight:    {X: 1, Y: 1},
}

func IsOutsideRoom(char entity.Object, room entity.Room) bool {
	return char.XYcoords.X < room.Coordinates.XYcoords.X ||
		char.XYcoords.X+char.W-1 >= room.Coordinates.XYcoords.X+room.Coordinates.W ||
		char.XYcoords.Y < room.Coordinates.XYcoords.Y ||
		char.XYcoords.Y+char.H-1 >= room.Coordinates.XYcoords.Y+room.Coordinates.H
}

func IsOutsideLevel(pos entity.Pos, level *entity.Level) bool {
	bounds := level.Coordinates
	return pos.X < bounds.XYcoords.X ||
		pos.X >= bounds.XYcoords.X+bounds.W ||
		pos.Y < bounds.XYcoords.Y ||
		pos.Y >= bounds.XYcoords.Y+bounds.H
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
	}
}

func MoveCharacterByDirectionObj(direction entity.Direction, characterGeometry *entity.Object) {
	switch direction {
	case entity.Forward:
		characterGeometry.XYcoords.Y--
	case entity.Left:
		characterGeometry.XYcoords.X--
	case entity.Right:
		characterGeometry.XYcoords.X++
	case entity.Back:
		characterGeometry.XYcoords.Y++
	case entity.DiagonallyForwardLeft:
		characterGeometry.XYcoords.X--
		characterGeometry.XYcoords.Y--
	case entity.DiagonallyForwardRight:
		characterGeometry.XYcoords.X++
		characterGeometry.XYcoords.Y--
	case entity.DiagonallyBackLeft:
		characterGeometry.XYcoords.X--
		characterGeometry.XYcoords.Y++
	case entity.DiagonallyBackRight:
		characterGeometry.XYcoords.X++
		characterGeometry.XYcoords.Y++
	case entity.Stop:
	}
}

func MoveMonster(monster *entity.Monster, level *entity.Level, player *entity.Player) {
	var path []entity.Pos

	if IsPlayerNear(monster, player) {
		path = FindPathToPlayer(monster, level, player)
	}

	if len(path) == 0 {
		switch monster.Type {
		case entity.Zombie:
			patternZombie(monster, level)
		case entity.Vampire:
			patternVampire(monster, level)
		case entity.Ghost:
			patternGhost(monster, level)
		case entity.Ogre:
			patternOgre(monster, level)
		case entity.Snake:
			patternSnake(monster, level)
		case entity.Mimic:
			patternZombie(monster, level)
		}
		return
	}

	testPos := entity.Pos{
		X: monster.Stats.Pos.XYcoords.X,
		Y: monster.Stats.Pos.XYcoords.Y,
	}

	if len(path) > 1 {
		next := path[len(path)-2]
		testPos = next
	}

	if testPos.X != player.BaseStats.Pos.XYcoords.X || testPos.Y != player.BaseStats.Pos.XYcoords.Y {
		if len(path) > 1 {
			next := path[len(path)-2]
			monster.Stats.Pos.XYcoords.X = next.X
			monster.Stats.Pos.XYcoords.Y = next.Y

			if len(path) >= 2 {
				prev := path[len(path)-1]
				next := path[len(path)-2]
				dx := next.X - prev.X
				dy := next.Y - prev.Y

				if dx == 0 && dy == -1 {
					monster.Dir = entity.Forward
				} else if dx == 0 && dy == 1 {
					monster.Dir = entity.Back
				} else if dx == -1 && dy == 0 {
					monster.Dir = entity.Left
				} else if dx == 1 && dy == 0 {
					monster.Dir = entity.Right
				} else if dx == -1 && dy == -1 {
					monster.Dir = entity.DiagonallyForwardLeft
				} else if dx == 1 && dy == -1 {
					monster.Dir = entity.DiagonallyForwardRight
				} else if dx == -1 && dy == 1 {
					monster.Dir = entity.DiagonallyBackLeft
				} else if dx == 1 && dy == 1 {
					monster.Dir = entity.DiagonallyBackRight
				}
			}
		}
	}
}

func MovePlayer(player *entity.Player, level *entity.Level, direction entity.Direction) {
	switch direction {
	case entity.Forward:
		player.BaseStats.Pos.XYcoords.Y--
	case entity.Back:
		player.BaseStats.Pos.XYcoords.Y++
	case entity.Left:
		player.BaseStats.Pos.XYcoords.X--
	case entity.Right:
		player.BaseStats.Pos.XYcoords.X++
	default:
	}

	newcoords := entity.Pos{
		X: player.BaseStats.Pos.XYcoords.X,
		Y: player.BaseStats.Pos.XYcoords.Y,
	}

	if !IsOutsideLevel(newcoords, level) && IsPassable(newcoords, level) {
		player.BaseStats.Pos.XYcoords.X = newcoords.X
		player.BaseStats.Pos.XYcoords.Y = newcoords.Y
	}
}

func patternZombie(monster *entity.Monster, level *entity.Level) {
	monsterObj := entity.Object{
		XYcoords: entity.Pos{X: monster.Stats.Pos.XYcoords.X, Y: monster.Stats.Pos.XYcoords.Y},
		W:        1,
		H:        1,
	}
	currentRoom := FindCurrentRoom(monsterObj, level)
	if currentRoom == nil {
		return
	}

	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		XY := entity.Pos{
			X: monster.Stats.Pos.XYcoords.X,
			Y: monster.Stats.Pos.XYcoords.Y,
		}

		direction := entity.Direction(logic.GetRandomInRange(0, SIMPLE_DIRECTIONS-1))
		MoveCharacterByDirection(direction, &XY)

		if !IsOutsideLevel(XY, level) && IsPassable(XY, level) {
			XYObj := entity.Object{XYcoords: XY, W: 1, H: 1}
			if !IsOutsideRoom(XYObj, *currentRoom) {
				monster.Stats.Pos.XYcoords.X = XY.X
				monster.Stats.Pos.XYcoords.Y = XY.Y
				monster.Dir = direction
				return
			}
		}
	}
}

func patternVampire(monster *entity.Monster, level *entity.Level) {
	monsterObj := entity.Object{
		XYcoords: entity.Pos{X: monster.Stats.Pos.XYcoords.X, Y: monster.Stats.Pos.XYcoords.Y},
		W:        1,
		H:        1,
	}
	currentRoom := FindCurrentRoom(monsterObj, level)
	if currentRoom == nil {
		return
	}

	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		XY := entity.Pos{
			X: monster.Stats.Pos.XYcoords.X,
			Y: monster.Stats.Pos.XYcoords.Y,
		}

		direction := entity.Direction(logic.GetRandomInRange(0, ALL_DIRECTIONS-1))
		MoveCharacterByDirection(direction, &XY)

		if !IsOutsideLevel(XY, level) && IsPassable(XY, level) {
			XYObj := entity.Object{XYcoords: XY, W: 1, H: 1}
			if !IsOutsideRoom(XYObj, *currentRoom) {
				monster.Stats.Pos.XYcoords.X = XY.X
				monster.Stats.Pos.XYcoords.Y = XY.Y
				monster.Dir = direction
				return
			}
		}
	}
}

func patternGhost(monster *entity.Monster, level *entity.Level) {
	ghostObj := entity.Object{
		XYcoords: entity.Pos{X: monster.Stats.Pos.XYcoords.X, Y: monster.Stats.Pos.XYcoords.Y},
		W:        1,
		H:        1,
	}

	currentRoom := FindCurrentRoom(ghostObj, level)
	if currentRoom == nil {
		return
	}

	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		randomX := logic.GetRandomInRange(currentRoom.Coordinates.XYcoords.X + 1, currentRoom.Coordinates.XYcoords.X + currentRoom.Coordinates.W - 2)
		randomY := logic.GetRandomInRange(currentRoom.Coordinates.XYcoords.Y + 1, currentRoom.Coordinates.XYcoords.Y + currentRoom.Coordinates.H - 2)

		newPos := entity.Pos{
			X: randomX,
			Y: randomY,
		}

		if IsPassable(newPos, level) {
			monster.Stats.Pos.XYcoords.X = newPos.X
			monster.Stats.Pos.XYcoords.Y = newPos.Y
			return
		}
	}
}


func patternOgre(monster *entity.Monster, level *entity.Level) {
	monsterObj := entity.Object{
		XYcoords: entity.Pos{X: monster.Stats.Pos.XYcoords.X, Y: monster.Stats.Pos.XYcoords.Y},
		W:        1,
		H:        1,
	}
	currentRoom := FindCurrentRoom(monsterObj, level)
	if currentRoom == nil {
		return
	}

	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		direction := entity.Direction(logic.GetRandomInRange(0, SIMPLE_DIRECTIONS-1))

		canMove := true
		testPos := entity.Pos{
			X: monster.Stats.Pos.XYcoords.X,
			Y: monster.Stats.Pos.XYcoords.Y,
		}

		for step := 0; step < OGRE_STEP; step++ {
			MoveCharacterByDirection(direction, &testPos)
			if IsOutsideLevel(testPos, level) || !IsPassable(testPos, level) {
				canMove = false
				break
			}
			testPosObj := entity.Object{XYcoords: testPos, W: 1, H: 1}
			if IsOutsideRoom(testPosObj, *currentRoom) {
				canMove = false
				break
			}
		}

		if canMove {
			for step := 0; step < OGRE_STEP; step++ {
				MoveCharacterByDirection(direction, &monster.Stats.Pos.XYcoords)
			}
			monster.Dir = direction
			return
		}
	}
}

func patternSnake(monster *entity.Monster, level *entity.Level) {
	monsterObj := entity.Object{
		XYcoords: entity.Pos{X: monster.Stats.Pos.XYcoords.X, Y: monster.Stats.Pos.XYcoords.Y},
		W:        1,
		H:        1,
	}
	currentRoom := FindCurrentRoom(monsterObj, level)
	if currentRoom == nil {
		return
	}

	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		XY := entity.Pos{
			X: monster.Stats.Pos.XYcoords.X,
			Y: monster.Stats.Pos.XYcoords.Y,
		}

		direction := entity.Direction(logic.GetRandomInRange(0, DIAGONAL_DIRECTIONS-1) + SIMPLE_TO_DIAGONAL_SHIFT)

		if direction == monster.Dir {
			continue
		}

		MoveCharacterByDirection(direction, &XY)

		if !IsOutsideLevel(XY, level) && IsPassable(XY, level) {
			XYObj := entity.Object{XYcoords: XY, W: 1, H: 1}
			if !IsOutsideRoom(XYObj, *currentRoom) {
				monster.Stats.Pos.XYcoords.X = XY.X
				monster.Stats.Pos.XYcoords.Y = XY.Y
				monster.Dir = direction
				return
			}
		}
	}
}

func getAggroRadius(hostility entity.HostilityType) int {
	switch hostility {
	case entity.Low:
		return entity.LOW_HOSTILITY_RADIUS
	case entity.Medium:
		return entity.AVERAGE_HOSTILITY_RADIUS
	case entity.High:
		return entity.HIGH_HOSTILITY_RADIUS
	default:
		return 4
	}
}

func IsPlayerNear(monster *entity.Monster, player *entity.Player) bool {
	dx := logic.Abs(monster.Stats.Pos.XYcoords.X - player.BaseStats.Pos.XYcoords.X)
	dy := logic.Abs(monster.Stats.Pos.XYcoords.Y - player.BaseStats.Pos.XYcoords.Y)
	dist := dx + dy
	return dist <= getAggroRadius(monster.Hostility)
}

func FindPathToPlayer(monster *entity.Monster, level *entity.Level, player *entity.Player) []entity.Pos {
	start := entity.Pos{
		X: monster.Stats.Pos.XYcoords.X,
		Y: monster.Stats.Pos.XYcoords.Y,
	}

	target := entity.Pos{
		X: player.BaseStats.Pos.XYcoords.X,
		Y: player.BaseStats.Pos.XYcoords.Y,
	}

	aggroRadius := getAggroRadius(monster.Hostility)
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

			if DistanceChebyshev(start, next) > aggroRadius {
				continue
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
	if pos.X < level.Coordinates.XYcoords.X || pos.X >= level.Coordinates.XYcoords.X+level.Coordinates.W ||
		pos.Y < level.Coordinates.XYcoords.Y || pos.Y >= level.Coordinates.XYcoords.Y+level.Coordinates.H {
		return false
	}

	for i := 0; i < level.DoorNumber; i++ {
		door := &level.Doors[i]
		if pos.X == door.Position.XYcoords.X && pos.Y == door.Position.XYcoords.Y {
			return door.IsOpen
		}
	}

	for _, room := range level.Rooms {
		r := room.Coordinates
		if pos.X > r.XYcoords.X && pos.X < r.XYcoords.X+r.W-1 &&
			pos.Y > r.XYcoords.Y && pos.Y < r.XYcoords.Y+r.H-1 {
			return true
		}
	}

	for i := 0; i < level.Passages.PassagesNumber; i++ {
		p := level.Passages.Passages[i]

		if pos.X > p.XYcoords.X && pos.X < p.XYcoords.X+p.W-1 &&
			pos.Y > p.XYcoords.Y && pos.Y < p.XYcoords.Y+p.H-1 {

			for j := 0; j < level.DoorNumber; j++ {
				door := &level.Doors[j]

				if door.Position.XYcoords.X >= p.XYcoords.X && door.Position.XYcoords.X < p.XYcoords.X+p.W &&
					door.Position.XYcoords.Y >= p.XYcoords.Y && door.Position.XYcoords.Y < p.XYcoords.Y+p.H {

					if !door.IsOpen {
						doorX := door.Position.XYcoords.X
						doorY := door.Position.XYcoords.Y

						if pos.X == doorX && pos.Y == doorY {
							return false
						}

						if p.H == 1 {
							if pos.Y == doorY {
								return false
							}
						}

						if p.W == 1 {
							if pos.X == doorX {
								return false
							}
						}
					}
				}
			}
			return true
		}
	}

	return false
}

func SkipNext(next entity.Pos, level *entity.Level, visited map[entity.Pos]bool) bool {
	if next.X < level.Coordinates.XYcoords.X || next.X >= level.Coordinates.XYcoords.X+level.Coordinates.W ||
		next.Y < level.Coordinates.XYcoords.Y || next.Y >= level.Coordinates.XYcoords.Y+level.Coordinates.H {
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

func DistanceChebyshev(first, second entity.Pos) int {
	dx := logic.Abs(first.X - second.X)
	dy := logic.Abs(first.Y - second.Y)

	if dx > dy {
		return dx
	}
	return dy
}

func RemoveMonsterFromRoom(level *entity.Level, monster *entity.Monster) {
    for i := range level.Rooms {
        room := &level.Rooms[i]
        for j := 0; j < room.MonsterNumbers; j++ {
            if &room.Monsters[j] == monster {
                // сдвигаем массив
                room.Monsters[j] = room.Monsters[room.MonsterNumbers-1]
                room.MonsterNumbers--
                return
            }
        }
    }
}

