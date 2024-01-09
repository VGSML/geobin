package planar

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// Clone returns full copy of geometry
func Clone(in orb.Geometry) orb.Geometry {
	switch g := in.(type) {
	case orb.Bound:
		return g
	case orb.Point:
		return g
	case orb.MultiPoint:
		return g.Clone()
	case orb.LineString:
		return g.Clone()
	case orb.MultiLineString:
		return g.Clone()
	case orb.Ring:
		return g.Clone()
	case orb.Polygon:
		return g.Clone()
	case orb.MultiPolygon:
		return g.Clone()
	case orb.Collection:
		return g.Clone()
	}
	return in
}

// Intersects checks a intersects b.
func Intersects(a, b orb.Geometry) bool {
	if !a.Bound().Intersects(b.Bound()) {
		return false
	}
	switch g := a.(type) {
	case orb.Bound:
		return true
	case orb.Point:
		return pointIntersects(g, b)
	case orb.MultiPoint:
		for _, p := range g {
			if pointIntersects(p, b) {
				return true
			}
		}
	case orb.LineString:
		return lineIntersects(g, b)
	case orb.MultiLineString:
		return multiLineIntersects(g, b)
	case orb.Ring:
		return lineIntersects(orb.LineString(g), b)
	case orb.Polygon:
		return polyIntersects(g, b)
	case orb.MultiPolygon:
		return multiPolyIntersects(g, b)
	case orb.Collection:
		for _, g := range g {
			if Intersects(g, b) {
				return true
			}
		}
	}
	return false
}

// Contains check a inside in b.
func Contains(a, b orb.Geometry) bool {
	if !a.Bound().Intersects(b.Bound()) {
		return false
	}
	switch g := a.(type) {
	case orb.Bound:
		return boundContains(g, b)
	case orb.Point:
		return pointContain(g, b)
	case orb.MultiPoint:
		for _, p := range g {
			if !pointContain(p, b) {
				return false
			}
		}
		return true
	case orb.LineString:
		return lineContains(g, b)
	case orb.MultiLineString:
		return multiLineContains(g, b)
	case orb.Ring:
		return lineContains(orb.LineString(g), b)
	case orb.Polygon:
		return polyContains(g, b)
	case orb.MultiPolygon:
		return multiPolyContains(g, b)
	case orb.Collection:
		for _, g := range g {
			if !Contains(g, b) {
				return false
			}
		}
		return true
	}
	return false
}

// Centroid calculate centroid of given geometry.
// For multipart geometry returned multiPoint.
func Centroid(geom orb.Geometry) orb.Geometry {
	switch geom := geom.(type) {
	case orb.Point, orb.MultiPoint:
		return geom
	case orb.LineString:
		return LineInterpolatePoint(geom, 0.5)
	case orb.MultiLineString:
		var mp orb.MultiPoint
		for _, l := range geom {
			mp = append(mp,
				LineInterpolatePoint(l, 0.5),
			)
		}
		return mp
	case orb.Polygon:
		pnt, _ := planar.CentroidArea(geom)
		return pnt
	case orb.MultiPolygon:
		var mp orb.MultiPoint
		for _, p := range geom {
			pnt, _ := planar.CentroidArea(p)
			mp = append(mp, pnt)
		}
		return mp
	case orb.Collection:
		var mp orb.Collection
		for _, geom := range geom {
			cp := Centroid(geom)
			mpp, ok := cp.(orb.MultiPoint)
			if !ok {
				mp = append(mp, cp)
			}
			if ok {
				for _, p := range mpp {
					mp = append(mp, p)
				}
			}
		}
		return mp
	case orb.Bound:
		return geom.Center()
	default:
		return nil
	}
}
