package h3f

import (
	"github.com/VGSML/geo-index/geo"
	"github.com/paulmach/orb"
	"github.com/uber/h3-go/v4"
)

// GeometryCells returns cells inside given geometry.
func GeometryCells(geom orb.Geometry, res int, compact bool) []h3.Cell {
	switch g := geom.(type) {
	case orb.Bound:
		poly := h3.GeoPolygon{
			GeoLoop: pointsToLatLngArray(g.ToRing()),
		}
		return poly.Cells(res)
	case orb.Point:
		return []h3.Cell{h3.NewLatLng(g.Lat(), g.Lon()).Cell(res)}
	case orb.MultiPoint:
		pp := pointsToLatLngArray(g)
		cells := make([]h3.Cell, 0, len(pp))
		for _, p := range pp {
			cells = append(cells, p.Cell(res))
		}
		return cells
	case orb.LineString:
		return lineCells(
			pointsToLatLngArray(g),
			res,
		)
	case orb.MultiLineString:
		var out []h3.Cell
		for _, l := range g {
			out = append(out,
				lineCells(
					pointsToLatLngArray(l),
					res,
				)...,
			)
		}
		return out
	case orb.Ring:
		return lineCells(
			pointsToLatLngArray(g),
			res,
		)
	case orb.Polygon:
		p := h3.GeoPolygon{}
		for i, r := range g {
			l := pointsToLatLngArray(r)
			if i == 0 {
				p.GeoLoop = l
				continue
			}
			p.Holes = append(p.Holes, l)
		}
		if !compact {
			return p.Cells(res)
		}
		cells := p.Cells(res)
		for len(cells) == 0 && res < 14 {
			res++
			cells = p.Cells(res)
		}
		if len(cells) < 100 {
			return cells
		}
		return h3.CompactCells(cells)
	case orb.MultiPolygon:
		var out []h3.Cell
		for _, p := range g {
			out = append(out,
				GeometryCells(p, res, compact)...,
			)
		}
		return out
	case orb.Collection:
		var out []h3.Cell
		for _, g := range g {
			out = append(out,
				GeometryCells(g, res, compact)...,
			)
		}
		return out
	}
	return nil
}

// GeometryPointCells returns cells corresponded for points of the given geometry.
func GeometryPointCells(geom orb.Geometry, res int) []h3.Cell {
	switch g := geom.(type) {
	case orb.Bound:
		return pointsToCellsArray(g.ToRing(), res)
	case orb.Point:
		return []h3.Cell{h3.NewLatLng(g.Lat(), g.Lon()).Cell(res)}
	case orb.MultiPoint:
		pp := pointsToLatLngArray(g)
		cells := make([]h3.Cell, 0, len(pp))
		for _, p := range pp {
			cells = append(cells, p.Cell(res))
		}
		return cells
	case orb.LineString:
		return pointsToCellsArray(g, res)
	case orb.MultiLineString:
		var out []h3.Cell
		for _, l := range g {
			out = append(out,
				pointsToCellsArray(l, res)...,
			)
		}
		return out
	case orb.Ring:
		return pointsToCellsArray(g, res)
	case orb.Polygon:
		if len(g) == 0 {
			return nil
		}
		return pointsToCellsArray(g[0], res)
	case orb.MultiPolygon:
		var out []h3.Cell
		for _, p := range g {
			out = append(out,
				GeometryPointCells(p, res)...,
			)
		}
		return out
	case orb.Collection:
		var out []h3.Cell
		for _, g := range g {
			out = append(out,
				GeometryPointCells(g, res)...,
			)
		}
		return out
	}
	return nil
}

// pointsToLatLngArray returns slice of h3.LatLng from []orb.Point
func pointsToLatLngArray(pp []orb.Point) []h3.LatLng {
	l := make([]h3.LatLng, 0, len(pp))
	for _, p := range pp {
		l = append(l, h3.NewLatLng(p.Lat(), p.Lon()))
	}
	return l
}

// pointsToCellsArray returns cells for given res from []orb.Point
func pointsToCellsArray(pp []orb.Point, res int) []h3.Cell {
	l := make([]h3.Cell, 0, len(pp))
	for _, p := range pp {
		l = append(l, h3.NewLatLng(p.Lat(), p.Lon()).Cell(res))
	}
	return l
}

// lineCells returns grid path for a long line.
func lineCells(line []h3.LatLng, res int) []h3.Cell {
	if len(line) == 0 {
		return nil
	}
	var out []h3.Cell
	first := line[0].Cell(res)
	for i := 1; i < len(line); i++ {
		next := line[i].Cell(res)
		if next == first { // skip equal point
			continue
		}
		if h3.GridDistance(first, next) == 0 {
			// distance too much split distance using planar geometry
			l := orb.LineString{
				{line[i-1].Lng, line[i-1].Lat},
				{line[i].Lng, line[i].Lat},
			}
			frac := 0.5
			splitPoint := geo.LineInterpolatePoint(l, 0.5)
			splitCell := h3.NewLatLng(splitPoint.Lat(), splitPoint.Lon()).Cell(res)
			for frac > 0.001 {
				if h3.GridDistance(first, splitCell) == 0 {
					frac /= 2
					splitPoint = geo.LineInterpolatePoint(l, frac)
					splitCell = h3.NewLatLng(splitPoint.Lat(), splitPoint.Lon()).Cell(res)
					continue
				}
				// add cells from split line
				for ff := frac; ff < 1; ff += frac {
					path := h3.GridPath(first, splitCell)
					out = append(out, path[:len(path)-1]...)
					first = splitCell
					splitPoint = geo.LineInterpolatePoint(l, ff+frac)
					splitCell = h3.NewLatLng(splitPoint.Lat(), splitPoint.Lon()).Cell(res)
				}
				break
			}
		}
		path := h3.GridPath(first, next)
		out = append(out, path[:len(path)-1]...)
		first = next
	}
	return append(out, line[len(line)-1].Cell(res))
}
