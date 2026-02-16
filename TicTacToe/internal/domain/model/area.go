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

func (f* GameField) CheckWin(symbol int) bool {
	// строки
	for i := 0; i < FieldSize; i++ {
		if f.Cells[i][0] == symbol &&
			f.Cells[i][1] == symbol &&
			f.Cells[i][2] == symbol {
			return true
		}
	}

	// колонки
	for j := 0; j < FieldSize; j++ {
		if f.Cells[0][j] == symbol &&
			f.Cells[1][j] == symbol &&
			f.Cells[2][j] == symbol {
			return true
		}
	}

	// диагонали
	if f.Cells[0][0] == symbol &&
		f.Cells[1][1] == symbol &&
		f.Cells[2][2] == symbol {
		return true
	}

	if f.Cells[0][2] == symbol &&
		f.Cells[1][1] == symbol &&
		f.Cells[2][0] == symbol {
		return true
	}

	return false
}

func CheckAllCellsFilled(f GameField) bool {
	for i := 0; i < FieldSize; i++ {
		for j := 0; j < FieldSize; j++ {
			if f.Cells[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

