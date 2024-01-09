package planar

import (
	"github.com/VGSML/geobin/orbf/vector"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// multiPolyIntersects checks if multi polygon intersects geom.
func multiPolyIntersects(mp orb.MultiPolygon, geom orb.Geometry) bool {
	for _, p := range mp {
		if polyIntersects(p, geom) {
			return true
		}
	}
	return false
}

// polyIntersects checks if polygon intersects geometry.
// TODO: write tests
func polyIntersects(poly orb.Polygon, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Bound:
		return polyIntersects(poly, g.ToPolygon())
	case orb.Point:
		return planar.PolygonContains(poly, g)
	case orb.MultiPoint:
		for _, pnt := range g {
			if planar.PolygonContains(poly, pnt) {
				return true
			}
		}
		return false
	case orb.LineString:
		for i := 1; i < len(g); i++ {
			if planar.PolygonContains(poly, g[i-1]) {
				return true
			}
			if planar.PolygonContains(poly, g[i]) {
				return true
			}
			for _, r := range poly {
				for j := 1; j < len(r); j++ {
					if _, ok := vector.IntersectionPoint(g[i-1], g[i], r[j-1], r[j]); ok {
						return true
					}
				}
			}
		}
		return false
	case orb.MultiLineString:
		for _, l := range g {
			if polyIntersects(poly, l) {
				return true
			}
		}
		return false
	case orb.Ring:
		return polyIntersects(poly, orb.LineString(g))
	case orb.Polygon:
		// check borders intersects or some point of geometry inside of polygon
		for _, r2 := range g {
			for i := 1; i < len(r2); i++ {
				if planar.PolygonContains(poly, r2[i-1]) {
					return true
				}
				if planar.PolygonContains(poly, r2[i]) {
					return true
				}
				for _, r1 := range poly {
					for j := 1; j < len(r1); j++ {
						if _, ok := vector.IntersectionPoint(r2[i-1], r2[i], r1[j-1], r1[j]); ok {
							return true
						}
					}
				}
			}
		}
		return false
	case orb.MultiPolygon:
		for _, p2 := range g {
			if polyIntersects(poly, p2) {
				return true
			}
		}
		return false
	case orb.Collection:
		for _, g := range g {
			if polyIntersects(poly, g) {
				return true
			}
		}
		return false
	}
	return false
}

// multiPolyContains checks if MultiPolygon contains geometry.
// If geometry is multi part (MultiPoint, MultiLineString, MultiPolygon or Collection),
// the function checks if polygons in given MultiPolygon contains all parts of geometry.
func multiPolyContains(mp orb.MultiPolygon, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.MultiPoint:
		for _, g := range g {
			if !multiPolyContains(mp, g) {
				return false
			}
		}
		return true
	case orb.MultiLineString:
		for _, g := range g {
			if !multiPolyContains(mp, g) {
				return false
			}
		}
		return true
	case orb.MultiPolygon:
		for _, g := range g {
			if !multiPolyContains(mp, g) {
				return false
			}
		}
		return true
	case orb.Collection:
		for _, g := range g {
			if !multiPolyContains(mp, g) {
				return false
			}
		}
		return true
	}

	for _, p := range mp {
		if polyContains(p, geom) {
			return true
		}
	}
	return false
}

// polyContains checks if polygon completely contains geometry.
func polyContains(poly orb.Polygon, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Bound:
		return polyIntersects(poly, g.ToPolygon())
	case orb.Point:
		return planar.PolygonContains(poly, g)
	case orb.MultiPoint:
		for _, pnt := range g {
			if !planar.PolygonContains(poly, pnt) {
				return false
			}
		}
		return true
	case orb.LineString:
		// check all point of line inside polygon
		for _, pnt := range g {
			if !planar.PolygonContains(poly, pnt) {
				return false
			}
		}
		// check holes do not intersect with line
		for i, r := range poly {
			if i == 0 {
				continue
			}
			if lineIntersects(g, orb.LineString(r)) {
				return false
			}
		}
		return true
	case orb.MultiLineString:
		// check all lines inside polygon
		for _, l := range g {
			if !polyContains(poly, l) {
				return false
			}
		}
		return true
	case orb.Ring:
		return polyContains(poly, orb.LineString(g))
	case orb.Polygon:
		// Check all points of polygon2 (g) are inside polygon1 (p) outer ring and not in its holes.
		if len(poly) == 0 {
			return false
		}
		for _, r := range g {
			for _, pnt := range r {
				if !planar.RingContains(poly[0], pnt) { // Not in outer ring
					return false
				}
				for _, h := range poly[1:] { // In a hole
					if planar.RingContains(h, pnt) {
						return false
					}
				}
			}
		}
		for _, h := range poly[1:] { // In a hole
			for _, pnt := range h {
				if planar.PolygonContains(g, pnt) {
					return false
				}
			}
		}
		return true
	case orb.MultiPolygon:
		for _, p2 := range g {
			if !polyContains(poly, p2) {
				return false
			}
		}
		return true
	case orb.Collection:
		for _, g := range g {
			if !polyContains(poly, g) {
				return false
			}
		}
		return true
	}

	return false
}
