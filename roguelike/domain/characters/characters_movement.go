package characters

import (
	"roguelike/domain/datastructs"
	"roguelike/domain/entity"
	"roguelike/domain/logic"
)


func CharacterOutsideBorder(characterCoords, room *entity.Object) bool {
	return (characterCoords.X+characterCoords.W-1 >= room.X+room.W-1) ||
		(characterCoords.X <= room.X) ||
		(characterCoords.Y <= room.Y) ||
		(characterCoords.Y+characterCoords.H-1 >= room.Y+room.H-1)
}

func MoveCharacterByDirection(direction entity.Direction, characterGeometry *entity.Object) {
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


func MoveMonster(monster *entity.Monster, level *entity.Level) {
	switch monster.Type {
	case entity.Zombie:
		if IsPlayerNear() {
			FindPathToPlayer()
		} else {
			patternZombie(monster, level)
		}
	case entity.Vampire:
		if IsPlayerNear() {
			FindPathToPlayer()
		} else {
			patternVampire(monster, level)			
		}		
	case entity.Ghost:
		if IsPlayerNear() {
			FindPathToPlayer()
		} else {
			patternGhost(monster, level)
		}
	case entity.Ogre:
		if IsPlayerNear() {
			FindPathToPlayer()
		} else {
			patternOgre(monster, level)
		}
	case entity.Snake:
		if IsPlayerNear() {
			FindPathToPlayer()
		} else {
			patternSnake(monster, level)
		}
	}
}

func PatternMonsters(monster *entity.Monster, level *entity.Level) {
	for try := 0; try < MAX_TRIES_TO_MOVE; try++ {
		coords := monster.Stats.Pos
		direction := entity.Direction(logic.GetRandomInRange(0, SIMPLE_DIRECTIONS))
		MoveCharacterByDirection(direction, &coords)

		if !CheckOutsideBorder(&coords, level) && CheckUnoccupiedLevel(&coords, level) {
            // Если ход допустим, обновляем координаты монстра
            monster.Stats.Pos = coords
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

func IsPlayerNear(playerCoordinates *entity.Object, monster *entity.Monster) bool {

}

func FindPathToPlayer(monster *entity.Monster, level *entity.Level, player entity.Player) {
	queue := []entity.Pos{}
	visited := make(map[entity.Pos]bool)
	parent := make(map[entity.Pos]entity.Pos)

	start := entity.Pos{
		X: monster.Stats.Pos.X,
		Y: monster.Stats.Pos.Y,
	}

	target := entity.Pos{
		X: player.BaseStats.Pos.X,
		Y: player.BaseStats.Pos.Y,
	}

	visited[start] = true
}