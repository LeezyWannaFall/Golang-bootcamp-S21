package presentation

import "roguelike/domain/entity"

const (
	TileWallHorizontal = '-'
	TileWallVertical   = '|'
	TileFloor          = '.'
	TileCorridor       = '#'
	TileRoomOpening    = '+'
	TileExit           = 'E'
	TileFog            = ' '
)

const (
	TilePlayer = '@'
)

const (
	TileMonsterZombie  = 'z'
	TileMonsterVampire = 'v'
	TileMonsterGhost   = 'g'
	TileMonsterOgre    = 'O'
	TileMonsterSnake   = 's'
	TileMonsterMimic   = 'm'
)

const (
	TileFood     = 'f'
	TileElixir   = 'e'
	TileScroll   = 'S'
	TileWeapon   = 'w'
	TileTreasure = '$'
	TileDoor     = 'D'
	TileKey      = 'k'
)

func GetMonsterTile(monsterType entity.MonsterType) rune {
	switch monsterType {
	case entity.Zombie:
		return TileMonsterZombie
	case entity.Vampire:
		return TileMonsterVampire
	case entity.Ghost:
		return TileMonsterGhost
	case entity.Ogre:
		return TileMonsterOgre
	case entity.Snake:
		return TileMonsterSnake
	case entity.Mimic:
		return TileMonsterMimic
	default:
		return '?'
	}
}

const (
	MapHeight = entity.ROOMS_IN_HEIGHT * entity.REGION_HEIGHT
	MapWidth  = entity.ROOMS_IN_WIDTH * entity.REGION_WIDTH
)
