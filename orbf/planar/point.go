package planar

import (
	"github.com/VGSML/geobin/orbf/vector"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// PointTouch checks geometry has point.
func PointTouch(p orb.Point, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Point:
		return g.Equal(p)
	case orb.MultiPoint:
		for _, cp := range g {
			if p.Equal(cp) {
				return true
			}
		}
		return false
	case orb.LineString:
		for _, cp := range g {
			if p.Equal(cp) {
				return true
			}
		}
		return false
	case orb.MultiLineString:
		found := false
		for _, l := range g {
			if PointTouch(p, l) {
				if found {
					return false
				}
				found = true
			}
		}
		return found
	case orb.Ring:
		return PointTouch(p, orb.LineString(g))
	case orb.Polygon:
		found := false
		for _, r := range g {
			if PointTouch(p, r) {
				if found {
					return false
				}
				found = true
			}
		}
		return found
	case orb.MultiPolygon:
		found := false
		for _, poly := range g {
			if PointTouch(p, poly) {
				if found {
					return false
				}
				found = true
			}
		}
		return found
	default:
		return false
	}
}

// ClosestPoint calculates the closes point on the line between two points.
func ClosestPoint(p1, p2, find orb.Point) (orb.Point, float64) {
	newPoint := vector.ClosestPointOnLineSegment(p1, p2, find)
	direction := 0.
	if vector.PointSide(p1, p2, find) > 0 {
		direction = 1
	} else {
		direction = -1
	}

	return newPoint, planar.Distance(find, newPoint) * direction
}

func pointIntersects(pnt orb.Point, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Bound:
		return g.IsEmpty()
	case orb.Point:
		return vector.PointEqual(pnt, g)
	case orb.MultiPoint:
		for _, p2 := range g {
			if vector.PointEqual(pnt, p2) {
				return true
			}
		}
		return false
	case orb.LineString:
		return lineIntersects(g, pnt)
	case orb.MultiLineString:
		return multiLineIntersects(g, pnt)
	case orb.Ring:
		return lineIntersects(orb.LineString(g), pnt)
	case orb.Polygon:
		return planar.PolygonContains(g, pnt)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(g, pnt)
	case orb.Collection:
		for _, g := range g {
			if pointIntersects(pnt, g) {
				return true
			}
		}
	}
	return false
}

// pointContain check if point contains geometry.
// Point can contain only Point and MultiPoint,
// in case geom is MultiPoint all contained points should be equal point.
func pointContain(p orb.Point, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Bound:
		return g.IsEmpty()
	case orb.Point:
		return vector.PointEqual(p, g)
	case orb.MultiPoint:
		for _, p2 := range g {
			if !vector.PointEqual(p, p2) {
				return false
			}
		}
		return true
	}
	return false
}
