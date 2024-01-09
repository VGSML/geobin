package geobin

import (
	"context"

	"github.com/VGSML/geobin/h3b"
	"github.com/VGSML/geobin/h3f"
	"github.com/VGSML/geobin/orbf"
	"github.com/VGSML/geobin/orbf/planar"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/project"
	"github.com/uber/h3-go/v4"
)

type Geometry interface {
	Geom() orb.Geometry
}

type Item interface {
	Index() int
	indexedCells() []h3.Cell

	Intersects(context.Context, Item) bool
	ContainsIn(context.Context, Item) bool
}

func newBoundIndexedItem(idx int, geom orb.Geometry, res int, proj Projection) Item {
	bound := geom.Bound()
	if proj == Mercator {
		bound = project.Bound(bound, project.Mercator.ToWGS84)
	}
	cell1 := h3.LatLngToCell(h3.LatLng{Lat: bound.Min.Lat(), Lng: bound.Min.Lon()}, res)
	cell2 := h3.LatLngToCell(h3.LatLng{Lat: bound.Max.Lat(), Lng: bound.Max.Lon()}, res)
	pc := []h3.Cell{cell1}
	cell := h3f.ParentIndex(
		cell1,
		cell2,
	)
	if cell != 0 {
		pc[0] = cell
	} else {
		pc = append(pc, cell2)
	}

	return &BoundIndexedItem{
		proj:      proj,
		idx:       idx,
		geom:      geom,
		baseCells: pc,
	}
}

type BoundIndexedItem struct {
	proj      Projection
	idx       int
	baseCells []h3.Cell
	geom      orb.Geometry
}

func (item *BoundIndexedItem) Index() int {
	return item.idx
}

func (item *BoundIndexedItem) Geom() orb.Geometry {
	return item.geom
}

func (item *BoundIndexedItem) indexedCells() []h3.Cell {
	return item.baseCells
}

func (item *BoundIndexedItem) Intersects(ctx context.Context, in Item) bool {
	switch in := in.(type) {
	case *BoundIndexedItem:
		if item.proj == WGS84 && item.proj == in.proj {
			return orbf.Intersects(item.geom, in.geom)
		}
		if item.proj == Mercator && item.proj == in.proj {
			return planar.Intersects(item.geom, in.geom)
		}
		if item.proj == Mercator {
			return planar.Intersects(item.geom,
				project.Geometry(orb.Clone(in.geom), project.WGS84.ToMercator),
			)
		}
		return planar.Intersects(
			project.Geometry(orb.Clone(item.geom), project.WGS84.ToMercator),
			in.geom,
		)
	case Geometry:
		if item.proj == Mercator {
			return planar.Intersects(item.geom, in.Geom())
		}
		return orbf.Intersects(item.geom, in.Geom())
	default:
		return false
	}
}

func (item *BoundIndexedItem) ContainsIn(ctx context.Context, in Item) bool {
	switch in := in.(type) {
	case *BoundIndexedItem:
		if item.proj == WGS84 && item.proj == in.proj {
			return orbf.Contains(item.geom, in.geom)
		}
		if item.proj == Mercator && item.proj == in.proj {
			return planar.Contains(item.geom, in.geom)
		}
		if item.proj == Mercator {
			return planar.Contains(item.geom,
				project.Geometry(orb.Clone(in.geom), project.WGS84.ToMercator),
			)
		}
		return planar.Contains(
			project.Geometry(orb.Clone(item.geom), project.WGS84.ToMercator),
			in.geom,
		)
	case Geometry:
		if item.proj == Mercator {
			return planar.Contains(item.geom, in.Geom())
		}
		return orbf.Contains(item.geom, in.Geom())
	default:
		return false
	}
}

type IndexedItem struct {
	idx   int
	index *h3b.Index // for single geometry item - point, line, polygon
}

func newIndexedItem(idx int, geom orb.Geometry, res int, compact bool, proj Projection) Item {
	if proj == Mercator {
		geom = project.Geometry(orb.Clone(geom), project.Mercator.ToWGS84)
	}
	bm := h3b.New(res)
	for i, cell := range h3f.GeometryCells(geom, res, compact) {
		bm.Insert(uint64(i), cell)
	}
	return &IndexedItem{
		idx:   idx,
		index: bm,
	}
}

func (item *IndexedItem) Index() int {
	return item.idx
}

func (item *IndexedItem) indexedCells() []h3.Cell {
	return item.index.ParentCells()
}

// Intersects return true if item intersects with given geometry.
func (item *IndexedItem) Intersects(ctx context.Context, in Item) bool {
	switch in := in.(type) {
	case *IndexedItem:
		return h3b.CheckIntersection(item.index, in.index)
	case Geometry:
		cells := h3f.GeometryCells(in.Geom(), int(item.index.Res()), true)
		check := h3b.New(int(item.index.Res()))
		for i, c := range cells {
			check.Insert(uint64(i), c)
		}
		return h3b.CheckIntersection(item.index, check)
	default:
		return false
	}
}

// ContainsIn return true if item inside given geometry.
func (item IndexedItem) ContainsIn(ctx context.Context, in Item) bool {
	switch in := in.(type) {
	case *IndexedItem:
		return h3b.CheckContainsIn(item.index, in.index)
	case Geometry:
		cells := h3f.GeometryCells(in.Geom(), int(item.index.Res()), true)
		check := h3b.New(int(item.index.Res()))
		for i, c := range cells {
			check.Insert(uint64(i), c)
		}
		return h3b.CheckContainsIn(item.index, check)
	default:
		return false
	}
}
