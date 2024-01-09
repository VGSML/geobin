package h3b

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/VGSML/geobin/h3f"
	"github.com/uber/h3-go/v4"
)

// Index bitmap index for h3.Cells.
type Index struct {
	res          int8
	minRes       int8
	baseCellsLen int8
	maxItemIndex uint64

	baseCellsMask *roaring64.Bitmap
	baseCellMap   [122]*roaring64.Bitmap   // bitmaps for each base cell
	resMaps       [15][8]*roaring64.Bitmap // set of bitmaps (0-6 - num of cell in parent, 7 if all cells) for each resolution
}

// New creates new bitmap index for h3.Cells
func New(res int) *Index {
	return &Index{
		baseCellsMask: roaring64.New(),
		res:           int8(res),
	}
}

// Res returns maximum resolution for the index.
func (i *Index) Res() int8 {
	return i.res
}

// MaxItemIndex returns maximum item index.
func (i *Index) MaxItemIndex() uint64 {
	return i.maxItemIndex
}

// SetMaxItemIndex SetMaxItemIndex maximum item index.
// Sets and returns true if given index greatest or equal current maximum index.
func (i *Index) SetMaxItemIndex(idx uint64) bool {
	i.maxItemIndex = max(i.maxItemIndex, idx)
	return i.maxItemIndex == idx
}

// Insert adds new cell to the index.
func (i *Index) Insert(idx uint64, cell h3.Cell) {
	bn := h3f.BaseCellNum(cell)
	if i.baseCellMap[bn] == nil {
		i.baseCellMap[bn] = roaring64.New()
		i.baseCellsLen += 1
		i.baseCellsMask.Add(uint64(bn))
	}
	i.baseCellMap[bn].Add(idx)
	cellRes := h3f.H3Res(cell)
	for res := 0; res < int(i.res); res++ {
		if res+1 > cellRes {
			if i.resMaps[res][7] == nil {
				i.resMaps[res][7] = roaring64.New()
			}
			i.resMaps[res][7].Add(idx)
			break
		}
		crn := h3f.CellIndexInRes(cell, res+1)
		if i.resMaps[res][crn] == nil {
			i.resMaps[res][crn] = roaring64.New()
		}
		i.resMaps[res][crn].Add(idx)
	}
	i.minRes = min(i.minRes, int8(cellRes))
	if i.minRes == 0 {
		i.minRes = 1
	}
	i.maxItemIndex = max(i.maxItemIndex, idx)
}

// Remove delete indexed element.
func (i *Index) Remove(idx uint64) bool {
	rm := false
	for bn, r := range i.baseCellMap {
		if r == nil {
			continue
		}
		if r.Contains(idx) {
			r.Remove(idx)
			rm = true
			if r.IsEmpty() {
				i.baseCellsMask.Remove(uint64(bn))
			}
		}
	}
	if !rm {
		return false
	}
	for _, rm := range i.resMaps {
		for _, r := range rm {
			r.Remove(idx)
		}
	}
	return true
}

// BaseCellsCount returns count of different base cells
func (i *Index) BaseCellsCount() int {
	return int(i.baseCellsLen)
}

// BaseCellsNum returns slice of all base cells nums.
func (i *Index) BaseCellsNum() []int8 {
	var nn []int8
	for n, m := range i.baseCellMap {
		if m != nil {
			nn = append(nn, int8(n))
		}
	}
	return nn
}

// HasCell checks index has intersects point with cell.
func (i *Index) HasCell(cell h3.Cell) bool {
	bn := h3f.BaseCellNum(cell)
	if i.baseCellMap[bn] == nil {
		return false
	}
	cellRes := h3f.H3Res(cell)
	base := i.baseCellMap[bn].Clone()

	for r := i.minRes - 1; r < i.res && r < int8(cellRes); r++ {
		crn := h3f.CellIndexInRes(cell, int(r+1))
		resMap := i.resMaps[r][crn]
		if resMap == nil {
			resMap = roaring64.New()
		}
		if i.resMaps[r][7] != nil {
			resMap.Or(i.resMaps[r][7])
		}
		if resMap.IsEmpty() {
			return false
		}
		base.And(resMap)
	}

	return !base.IsEmpty()
}

// ItemCells returns h3 cells for item.
func (i *Index) ItemCells(idx uint64) []h3.Cell {
	out := make([]h3.Cell, 0, 2)
	filter := roaring64.New()
	filter.Add(idx)
	for bn, bm := range i.baseCellMap {
		if bm == nil {
			continue
		}
		base := roaring64.And(bm, filter)
		if base.IsEmpty() {
			continue
		}
		var levels [][]int
		for _, rm := range i.resMaps {
			level := make([][]int, 0, len(levels))
			for cn, cm := range rm {
				if !cm.Contains(idx) {
					continue
				}
				if cn == 7 {
					for _, item := range levels {
						out = append(out, h3f.BuildH3Cell(bn, item...))
					}
					continue
				}
				for i := range levels {
					level = append(level, append(levels[i], cn))
				}
			}
			levels = level
		}
	}
	return out
}

// Intersects checks index has intersects point with cells.
func (i *Index) Intersects(cells []h3.Cell) bool {
	allItems := roaring64.New()
	for _, cell := range cells {
		bn := h3f.BaseCellNum(cell)
		if i.baseCellMap[bn] == nil {
			continue
		}
		base := i.baseCellMap[bn].Clone()
		cellRes := h3f.H3Res(cell)
		for r := i.minRes - 1; r < i.res && r < int8(cellRes); r++ {
			crn := h3f.CellIndexInRes(cell, int(r+1))
			resMap := i.resMaps[r][crn]
			if resMap == nil {
				resMap = roaring64.New()
			}
			if i.resMaps[r][7] != nil {
				resMap.Or(i.resMaps[r][7])
			}
			base.And(resMap)
		}
		allItems.Or(base)
	}

	return !allItems.IsEmpty()
}

// Intersection returns indexes that have intersects point with cells.
func (i *Index) Intersection(cells []h3.Cell) *roaring64.Bitmap {
	allItems := roaring64.New()
	for _, cell := range cells {
		bn := h3f.BaseCellNum(cell)
		if i.baseCellMap[bn] == nil {
			continue
		}
		base := i.baseCellMap[bn].Clone()
		cellRes := h3f.H3Res(cell)
		for r := i.minRes - 1; r < i.res && r < int8(cellRes); r++ {
			crn := h3f.CellIndexInRes(cell, int(r+1))
			resMap := i.resMaps[r][crn]
			if resMap == nil {
				resMap = roaring64.New()
			}
			if i.resMaps[r][7] != nil {
				resMap.Or(i.resMaps[r][7])
			}
			base.And(resMap)
		}
		allItems.Or(base)
	}

	return allItems
}

// ContainsInItems returns all indexed elements that inside cells.
func (i *Index) ContainsInItems(cells []h3.Cell) *roaring64.Bitmap {
	allItems := roaring64.New()
	for _, cell := range cells {
		bn := h3f.BaseCellNum(cell)
		if i.baseCellMap[bn] == nil {
			continue
		}
		cellRes := h3f.H3Res(cell)
		base := i.baseCellMap[bn].Clone()
		// loop over cell children, add to current cell items map full children cells or last cell
		fullChildren := roaring64.New()
		for r := 0; r < int(i.res); r++ {
			resMap := roaring64.New()
			// add all bitmaps before current cell res
			if r < cellRes {
				crn := h3f.CellIndexInRes(cell, r+1)
				if i.resMaps[r][crn] != nil {
					resMap = i.resMaps[r][crn]
				}
				base.And(resMap)
				continue
			}

			for n := 0; n < 7; n++ {
				if i.resMaps[r][n] != nil {
					resMap.Or(i.resMaps[r][n])
				}
			}
			if i.resMaps[r][7] != nil {
				fullChildren.Or(
					roaring64.And(base, i.resMaps[r][7]),
				)
			}
			if resMap.IsEmpty() {
				break
			}
			base.And(resMap)
		}
		base.Or(fullChildren)
		allItems.Or(base)
	}

	return allItems
}

// ParentCells returns cells, that contains all index element
func (i *Index) ParentCells() []h3.Cell {
	out := make([]h3.Cell, 0, i.baseCellsLen)
	for bn, bm := range i.baseCellMap {
		if bm == nil || bm.IsEmpty() {
			continue
		}
		out = append(out, i.parentCellForBaseCell(bn))
	}
	return out
}

// parentCellForBaseCell return parent cell for all items in base cell.
func (i *Index) parentCellForBaseCell(baseCellNum int) h3.Cell {
	base := i.baseCellMap[baseCellNum].Clone()
	var components []int
	for res := 0; res < int(i.res); res++ {
		cnn := make([]int, 0, 7)
		isFull := false
		for cn, rm := range i.resMaps[res] {
			if rm == nil || rm.IsEmpty() {
				continue
			}
			if !base.Intersects(rm) {
				continue
			}
			cnn = append(cnn, cn)
			if cn == 7 {
				isFull = true
			}
		}
		if isFull || len(cnn) != 1 { // if number of children cells > 1 add prev cell and exit
			return h3f.BuildH3Cell(baseCellNum, components...)
		}
		components = append(components, cnn[0])
	}
	return h3f.BuildH3Cell(baseCellNum, components...)
}

// Clone makes a copy of BitmapIndex.
func (i *Index) Clone() *Index {
	ni := &Index{
		res:          i.res,
		minRes:       i.minRes,
		baseCellsLen: i.baseCellsLen,
		maxItemIndex: i.maxItemIndex,
		baseCellMap:  [122]*roaring64.Bitmap{},
		resMaps:      [15][8]*roaring64.Bitmap{},
	}

	for bn, bm := range i.baseCellMap {
		if bm == nil {
			continue
		}
		ni.baseCellMap[bn] = bm.Clone()
	}

	for res := range i.resMaps {
		for cn, cm := range i.resMaps[res] {
			if cm == nil {
				continue
			}
			ni.resMaps[res][cn] = cm.Clone()
		}
	}

	return ni
}
