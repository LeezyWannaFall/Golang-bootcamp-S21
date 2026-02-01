package model

const FieldSize = 3

type GameField struct {
	Cells [FieldSize][FieldSize]int
}

func (f* GameField) ClearField() {
	for i := 0; i < FieldSize; i++ {
		for j := 0; j < FieldSize; j++ {
			f.Cells[i][j] = 0
		}
	}
}

func (f* GameField) PlaceSymbolOnField(x, y int, symbol int) bool {
	if x < 0 || x >= FieldSize || y < 0 || y >= FieldSize {
		return false
	}

	if f.Cells[y][x] != 0 {
		return false
	}

	f.Cells[y][x] = symbol
	return true
}

