package datastructs

import "roguelike/domain/entity"

type Vector struct {
	Size, Capacity int
	Data []entity.Direction
}