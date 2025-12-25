package characters

import (
	"roguelike/domain/datastructs"
	"roguelike/domain/entity"
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

func MoveCharacterByPath(path *datastructs.Vector, characterGeometry *entity.Object) {
	if path == nil {
		return
	}
	for _, direction := range path.Data {
		MoveCharacterByDirection(direction.(entity.Direction), characterGeometry)
	}
}
