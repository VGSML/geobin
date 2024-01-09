package h3f

import "github.com/uber/h3-go/v4"

// ParentIndex returns parent index for both cell.
func ParentIndex(idx1, idx2 h3.Cell) h3.Cell {
	if idx1 == 0 {
		return idx2
	}
	if idx2 == 0 {
		return idx1
	}
	minRes := H3Res(idx1)
	if minRes > H3Res(idx2) {
		minRes = H3Res(idx2)
	}
	for res := minRes; res > 0; res-- {
		parentIdx := H3Parent(idx1, res)
		if parentIdx == H3Parent(idx2, res) {
			return parentIdx
		}
	}
	return 0
}

// IsParent checks the idx1 is parent of the idx2.
func IsParent(idx1, idx2 h3.Cell) bool {
	if idx1 == idx2 {
		return true
	}
	res := H3Res(idx1)
	if res > H3Res(idx2) {
		res = H3Res(idx2)
	}
	return H3Parent(idx1, res) == H3Parent(idx2, res)
}

// H3Res returns cell resolution.
func H3Res(val h3.Cell) int {
	return int((val >> 52) & 15) // int((uint64(val) & 0xf0000000000000) >> 52)
}

// H3Parent returns parent cell.
func H3Parent(idx h3.Cell, res int) h3.Cell {
	if res > 15 || res < 0 {
		return 0
	}
	idx = h3.Cell(uint64(idx)&0xff0fffffffffffff) | h3.Cell(res<<52)

	res = 15 - res
	return idx | (1<<(res*3) - 1)
}

// BaseCellNum returns base cell num
func BaseCellNum(cell h3.Cell) int {
	return int((cell >> 45) & 127)
}

// CellIndexInRes returns the cell number (0..7) in resolution.
func CellIndexInRes(idx h3.Cell, res int) int {
	switch res {
	case 1:
		return int(idx&0x1c0000000000) >> 42
	case 2:
		return int(idx&0x38000000000) >> 39
	case 3:
		return int(idx&0x7000000000) >> 36
	case 4:
		return int(idx&0xe00000000) >> 33
	case 5:
		return int(idx&0x1c0000000) >> 30
	case 6:
		return int(idx&0x38000000) >> 27
	case 7:
		return int(idx&0x7000000) >> 24
	case 8:
		return int(idx&0xe00000) >> 21
	case 9:
		return int(idx&0x1c0000) >> 18
	case 10:
		return int(idx&0x38000) >> 15
	case 11:
		return int(idx&0x7000) >> 12
	case 12:
		return int(idx&0xe00) >> 9
	case 13:
		return int(idx&0x1c0) >> 6
	case 14:
		return int(idx&0x38) >> 3
	case 15:
		return int(idx & 0x7)
	}
	return 7
}

// BuildH3Cell returns h3.Cell for given mode, baseCell, and cellNum for each resolution.
func BuildH3Cell(baseCell int, cellNum ...int) h3.Cell {
	res := 0
	part := 0
	for i := 0; i < 15; i++ {
		cn := 7
		if len(cellNum) > i {
			cn = cellNum[i]
		}
		if cn != 7 {
			res++
		}
		part |= cn << (45 - (i+1)*3)
	}

	return h3.Cell(1<<59 | (res&15)<<52 | (baseCell&127)<<45 | part)
}

// buildH3Cell returns h3.Cell for given mode, baseCell, and cellNum for each resolution.
func BuildH3CellRes(res, baseCell int, digits *[15]int) h3.Cell {
	part := 0
	for i := 0; i < 15; i++ {
		cn := 7
		if digits != nil && i < res {
			cn = digits[i]
		}
		part |= cn << (45 - (i+1)*3)
	}

	return h3.Cell(1<<59 | (res&15)<<52 | (baseCell&127)<<45 | part)
}

// ChildrenBitmapMask returns children bitmap mask.
func ChildrenBitmapMask(idx h3.Cell, res int) uint {
	cellIndex1 := CellIndexInRes(idx, res)
	if cellIndex1 == 7 {
		return 0b1111111111111111111111111111111111111111111111111
	}
	cellIndex2 := CellIndexInRes(idx, res+1)
	if cellIndex2 == 7 || res == 15 {
		bitmap := uint(1<<(cellIndex2+1) - 1)
		return bitmap << ((7 - cellIndex1) * 6)
	}
	bitmap := uint(1 << (cellIndex2))
	return bitmap << ((7 - cellIndex1) * 6)
}
