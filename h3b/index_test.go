package h3b

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/VGSML/geobin/h3f"
	"github.com/VGSML/geobin/internal/fixture"
	"github.com/uber/h3-go/v4"
)

func Test_baseCellNum(t *testing.T) {
	cell := h3.Cell(0x851205b7fffffff)

	t.Logf("%b\n", cell)
	baseCell := h3f.BaseCellNum(cell)
	t.Logf("basecell: %d, %[1]b\n", h3f.BaseCellNum(cell))
	t.Logf("res: %d, %[1]b\n", h3f.H3Res(cell))

	var compoents []int
	for res := 1; res < 16; res++ {
		cellNum := h3f.CellIndexInRes(cell, res)
		compoents = append(compoents, cellNum)
		t.Logf("num res %d: %d, %[2]b\n", res, cellNum)

	}

	res := 0
	part := 0
	for i := 0; i < 15; i++ {
		cellNum := 7
		if len(compoents) > i {
			cellNum = compoents[i]
		}
		if cellNum != 7 {
			res++
		}
		part |= cellNum << (45 - (i+1)*3)
	}
	newCell := h3.Cell(1<<59 | (res&15)<<52 | (baseCell&127)<<45 | part)
	t.Logf("%b\n", cell)
	t.Logf("%b\n", newCell)
	t.Logf("basecell: %d, %[1]b\n", h3f.BaseCellNum(newCell))
	t.Logf("res: %d, %d, %[2]b\n", res, h3f.H3Res(newCell))

	t.Logf("basecell: %d %[1]b", (newCell>>45)&127)

}

func Test_baseCellChildren(t *testing.T) {
	for i := 0; i < 122; i++ {
		bn := h3f.BuildH3Cell(i)
		if !bn.IsValid() {
			t.Error("cell is not valid", i, bn.String())
		}
		t.Logf("BaseCell:%d, isPentagon: %v: %s", i, bn.IsPentagon(), bn.String())
	}
}

func Test_nodeIndexes(t *testing.T) {
	dist := h3.GridDistance(
		h3.NewLatLng(40.7128, -74.0060).Cell(15),
		h3.NewLatLng(51.5074, -30.1278).Cell(15),
		//h3.NewLatLng(51.5074, -0.1278).Cell(15),
	)
	if dist == 0 {
		t.Fatal("dist == 0")
	}

	cc := h3.GridPath(
		h3.NewLatLng(40.7128, -74.0060).Cell(15),
		h3.NewLatLng(51.5074, -30.1278).Cell(15),
		//h3.NewLatLng(51.5074, -0.1278).Cell(15),
	)
	t.Log(len(cc))
}

func TestBitmapIndex_Insert(t *testing.T) {
	lines := fixture.LineString()
	b := New(15)
	cells := h3f.GeometryCells(lines[0], 13, true)

	for idx, cell := range cells {
		b.Insert(uint64(idx), cell)
	}

	for bn, baseMap := range b.baseCellMap {
		t.Logf("BaseCell: %d, exists %v", bn, baseMap != nil)
		if baseMap != nil {
			t.Log(len(baseMap.ToArray()))
		}
	}
}

func TestBitmapIndex_HasCell(t *testing.T) {
	tests := []struct {
		name       string
		indexCells []h3.Cell
		inputCell  h3.Cell
		want       bool
	}{
		{
			name:       "empty",
			indexCells: []h3.Cell{},
			inputCell:  0x801dfffffffffff,
			want:       false,
		},
		{
			name:       "not intersection base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCell:  0x801dfffffffffff,
			want:       false,
		},
		{
			name:       "equal base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCell:  0x8013fffffffffff,
			want:       true,
		},
		{
			name: "equal res 4",
			indexCells: []h3.Cell{
				0x851205b7fffffff,
			},
			inputCell: 0x851205b7fffffff,
			want:      true,
		},
		{
			name: "children of res 4",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			inputCell: 0x861205bb7ffffff,
			want:      true,
		},
		{
			name: "parent of res 4",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			inputCell: 0x831205fffffffff,
			want:      true,
		},
		{
			name: "parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCell: 0x831205fffffffff,
			want:      true,
		},
		{
			name: "different parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCell: 0x842aac1ffffffff,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(15)
			for i, cell := range tt.indexCells {
				b.Insert(uint64(i), cell)
			}
			got := b.HasCell(tt.inputCell)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmapIndex_Intersection(t *testing.T) {
	tests := []struct {
		name       string
		indexCells []h3.Cell
		inputCells []h3.Cell
		want       []uint64
	}{
		{
			name:       "empty",
			indexCells: []h3.Cell{},
			inputCells: []h3.Cell{0x801dfffffffffff},
			want:       []uint64{},
		},
		{
			name:       "not intersection base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCells: []h3.Cell{0x801dfffffffffff},
			want:       []uint64{},
		},
		{
			name:       "equal base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCells: []h3.Cell{0x8013fffffffffff},
			want:       []uint64{0},
		},
		{
			name: "equal res 4",
			indexCells: []h3.Cell{
				0x851205b7fffffff,
			},
			inputCells: []h3.Cell{0x851205b7fffffff},
			want:       []uint64{0},
		},
		{
			name: "children of res 4",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			inputCells: []h3.Cell{0x861205bb7ffffff},
			want:       []uint64{6},
		},
		{
			name: "parent of res 4",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "different parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x842aac1ffffffff},
			want:       []uint64{},
		},
		{
			name: "parent of res 4 many base with two cell",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff, 0x842aac1ffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "two different cells set in one base",
			indexCells: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			inputCells: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			want: []uint64{20, 21},
		},
		{
			name: "two different cells set in one base",
			indexCells: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			inputCells: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			want: []uint64{0, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(15)
			for i, cell := range tt.indexCells {
				b.Insert(uint64(i), cell)
			}
			//testPrintBitmap(b)
			has := b.Intersects(tt.inputCells)
			if (len(tt.want) != 0) != has {
				t.Errorf("Intersects: got %v, want %v", has, len(tt.want) != 0)
			}
			got := b.Intersection(tt.inputCells)
			if !reflect.DeepEqual(got.ToArray(), tt.want) {
				t.Errorf("Intersection: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmapIndex_ContainsIn(t *testing.T) {
	tests := []struct {
		name       string
		indexCells []h3.Cell
		inputCells []h3.Cell
		want       []uint64
	}{
		{
			name:       "empty",
			indexCells: []h3.Cell{},
			inputCells: []h3.Cell{0x801dfffffffffff},
			want:       []uint64{},
		},
		{
			name:       "not intersection base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCells: []h3.Cell{0x801dfffffffffff},
			want:       []uint64{},
		},
		{
			name:       "equal base",
			indexCells: []h3.Cell{0x8013fffffffffff},
			inputCells: []h3.Cell{0x8013fffffffffff},
			want:       []uint64{0},
		},
		{
			name: "equal res 4",
			indexCells: []h3.Cell{
				0x851205b7fffffff,
			},
			inputCells: []h3.Cell{0x851205b7fffffff},
			want:       []uint64{0},
		},
		{
			name: "children of res 4",
			indexCells: []h3.Cell{
				0x861205bb7ffffff,
			},
			inputCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			want: []uint64{0},
		},
		{
			name: "parent of res 4",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "different parent of res 4 many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x842aac1ffffffff},
			want:       []uint64{},
		},
		{
			name: "parent of res 4 many base with two cell",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			inputCells: []h3.Cell{0x831205fffffffff, 0x842aac1ffffffff},
			want:       []uint64{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "two different cells set in one base",
			indexCells: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			inputCells: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			want: []uint64{},
		},
		{
			name: "two different cells set in one base",
			indexCells: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			inputCells: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			want: []uint64{0, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(15)
			for i, cell := range tt.indexCells {
				b.Insert(uint64(i), cell)
			}
			//testPrintBitmap(b)
			got := b.ContainsInItems(tt.inputCells)
			if !reflect.DeepEqual(got.ToArray(), tt.want) {
				t.Errorf("Intersection: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmapIndex_ParentCells(t *testing.T) {
	tests := []struct {
		name       string
		indexCells []h3.Cell
		want       []h3.Cell
	}{
		{
			name:       "empty",
			indexCells: []h3.Cell{},
			want:       []h3.Cell{},
		},
		{
			name:       "equal",
			indexCells: []h3.Cell{0x8013fffffffffff},
			want:       []h3.Cell{0x8013fffffffffff},
		},
		{
			name: "one base cell",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
			},
			want: []h3.Cell{0x841205bffffffff},
		},
		{
			name: "one base cell diff res",
			indexCells: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x861205b47ffffff,
				0x871205b85ffffff,
			},
			want: []h3.Cell{0x841205bffffffff},
		},
		{
			name: "many base",
			indexCells: []h3.Cell{
				0x851205a3fffffff, 0x851205a7fffffff, 0x851205abfffffff, 0x851205affffffff,
				0x851205b3fffffff, 0x851205b7fffffff, 0x851205bbfffffff, 0x85136a43fffffff,
				0x85136a47fffffff, 0x85136a4bfffffff, 0x85136a4ffffffff, 0x85136a53fffffff,
				0x85136a57fffffff, 0x85136a5bfffffff, 0x8448c43ffffffff, 0x8448c45ffffffff,
				0x8448c47ffffffff, 0x8448c49ffffffff, 0x8448c4bffffffff, 0x8448c4dffffffff,
			},
			want: []h3.Cell{0x8013fffffffffff, 0x8348c4fffffffff},
		},
		{
			name: "many base",
			indexCells: []h3.Cell{
				0x8448c47ffffffff, 0x8448c49ffffffff, 0x8448c4bffffffff,
				0x8448c43ffffffff, 0x8448c5bffffffff, 0x8448c15ffffffff,
				0x8448eb3ffffffff, 0x8429b47ffffffff, 0x8429b55ffffffff,
				0x8429b67ffffffff, 0x8429b69ffffffff, 0x8448cc3ffffffff,
				0x8448cebffffffff, 0x84299cdffffffff, 0x842995dffffffff, 0x8448c95ffffffff,
			},
			want: []h3.Cell{0x8129bffffffffff, 0x8148fffffffffff},
		},
		{
			name: "many base",
			indexCells: []h3.Cell{
				0x8448ca5ffffffff, 0x84268a3ffffffff, 0x8448dedffffffff, 0x8426d61ffffffff,
				0x448d43ffffffff, 0x8426dedffffffff, 0x8426c67ffffffff, 0x8448c3bffffffff,
			},
			want: []h3.Cell{0x8027fffffffffff, 0x8148fffffffffff},
		},
		{
			name: "one base",
			indexCells: []h3.Cell{
				0x8448ca5ffffffff, 0x8448dedffffffff,
				0x448d43ffffffff, 0x8448c3bffffffff,
			},
			want: []h3.Cell{0x8148fffffffffff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(15)
			for i, cell := range tt.indexCells {
				b.Insert(uint64(i), cell)
			}
			got := b.ParentCells()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func testPrintBitmap(b *Index) {
	for bn, bm := range b.baseCellMap {
		if bm == nil || bm.IsEmpty() {
			continue
		}
		base := bm.Clone()
		full := roaring64.New()
		for rn, rm := range b.resMaps {
			resMap := roaring64.New()
			for i := 0; i < 7; i++ {
				if pm := rm[i]; pm != nil {
					resMap.Or(pm)
				}
			}
			if rm[7] != nil {
				fmt.Printf("\t\tfull index: %d, %d. Items: %v\n", bn, rn,
					roaring64.And(base, rm[7]).ToArray(),
				)
				full.Or(roaring64.And(base, rm[7]))
			}
			if resMap.IsEmpty() {
				break
			}
			base.And(resMap)
		}
		// getting full list of elements for each base cell
		base.Or(full)
		fmt.Println("\t", bn, ":", base.ToArray())
	}
}
