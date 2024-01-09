package planar

import (
	"math"

	"github.com/VGSML/geobin/orbf/vector"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

// IsPointOnTheLine returns true if point located on the line.
func IsPointOnTheLine(line orb.LineString, pnt orb.Point) bool {
	for i := 1; i < len(line); i++ {
		if vector.PointOnSegment(line[i-1], line[i], pnt) {
			return true
		}
	}
	return false
}

// LineLocatePoint returns fraction of the line length on the line
// and minimal distance between point and line, if point located on
// the right forward line direction then distance positive, else negative.
func LineLocatePoint(line orb.LineString, point orb.Point) (float64, float64) {
	if len(line) == 0 {
		return 0, 0
	}
	if len(line) == 1 {
		return 0, geo.Distance(line[0], point)
	}

	lineLen := planar.Length(line)
	if lineLen == 0 {
		return 0, 0
	}

	var cp orb.Point
	index := 0
	minDist := math.MaxFloat64
	for i := 1; i < len(line); i++ {
		pnt, dist := vector.ClosestPoint(line[i-1], line[i], point)
		if math.Abs(dist) < math.Abs(minDist) {
			cp = pnt
			minDist = dist
			index = i
		}
	}

	return planar.Length(append(line.Clone()[0:index], cp)) / lineLen, minDist
}

// LineInterpolatePoint returns the point on the line located by fraction of the line length.
// Work with 2d euclidean geometry.
func LineInterpolatePoint(l orb.LineString, frac float64) orb.Point {
	if frac < 0 || frac > 1 || len(l) == 0 {
		return orb.Point{}
	}
	if frac == 0 {
		return l[0]
	}
	if frac == 1 {
		return l[len(l)-1]
	}

	lineDist := planar.Length(l)
	traveled := 0.
	for i := 1; i < len(l); i++ {
		dist := planar.Distance(l[i-1], l[i])
		actDist := lineDist*frac - traveled
		if dist == 0 {
			continue
		}
		if dist >= actDist {
			frac = actDist / dist
			return vector.LineSegmentInterpolatePoint(l[i-1], l[i], frac)
		}
		traveled += dist
	}

	return orb.Point{}
}

func LinesIntersectPoint(line1, line2 orb.LineString) (orb.Point, bool) {
	if !line1.Bound().Intersects(line2.Bound()) {
		return orb.Point{}, false
	}
	for i := 1; i < len(line1); i++ {
		for j := 1; j < len(line2); j++ {
			if p, ok := vector.IntersectionPoint(line1[i-1], line1[i], line2[j-1], line2[j]); ok {
				return p, true
			}
		}
	}

	return orb.Point{}, false
}

// LineTouch return common points of line and geometry, if exists return true.
func LineTouch(line, geom orb.Geometry) (orb.Point, bool) {
	if !line.Bound().Intersects(geom.Bound()) {
		return orb.Point{}, false
	}
	if ml, ok := line.(*orb.MultiLineString); ok {
		return LineTouch(*ml, geom)
	}
	if ml, ok := line.(orb.MultiLineString); ok {
		found := false
		var point orb.Point
		for _, l := range ml {
			if pnt, ok := LineTouch(l, geom); ok {
				if found {
					return orb.Point{}, false
				}
				found = true
				point = pnt
			}
		}
		return point, found
	}
	if l, ok := line.(*orb.LineString); ok {
		return LineTouch(*l, geom)
	}
	l, ok := line.(orb.LineString)
	if !ok {
		return orb.Point{}, false
	}
	found := false
	var point orb.Point
	for _, p := range l {
		if PointTouch(p, geom) {
			if found {
				return orb.Point{}, false
			}
			found = true
			point = p
		}
	}
	return point, found
}

// multiLineIntersects checks if multiline intersects geometry.
func multiLineIntersects(ml orb.MultiLineString, geom orb.Geometry) bool {
	for _, l := range ml {
		if lineIntersects(l, geom) {
			return true
		}
	}
	return false
}

// lineIntersects checks if line intersects geometry.
func lineIntersects(l orb.LineString, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Bound:
		return lineIntersects(l, g.ToPolygon())
	case orb.Point:
		return IsPointOnTheLine(l, g)
	case orb.MultiPoint:
		for _, p := range g {
			if lineIntersects(l, p) {
				return true
			}
		}
		return false
	case orb.LineString:
		_, ok := LinesIntersectPoint(l, g)
		return ok
	case orb.MultiLineString:
		for _, l2 := range g {
			if _, ok := LinesIntersectPoint(l, l2); ok {
				return ok
			}
		}
	case orb.Ring:
		_, ok := LinesIntersectPoint(l, orb.LineString(g))
		return ok
	case orb.Polygon:
		return polyIntersects(g, l)
	case orb.MultiPolygon:
		return multiPolyIntersects(g, l)
	case orb.Collection:
		for _, g := range g {
			if lineIntersects(l, g) {
				return true
			}
		}
	}
	return false
}

// multiLineContains checks if multiline contains geometry.
func multiLineContains(ml orb.MultiLineString, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.MultiPoint:
		for _, g := range g {
			if !multiLineContains(ml, g) {
				return false
			}
		}
		return true
	case orb.MultiLineString:
		for _, g := range g {
			if !multiLineContains(ml, g) {
				return false
			}
		}
		return true
	case orb.Collection:
		for _, g := range g {
			if !multiLineContains(ml, g) {
				return false
			}
		}
		return true
	case orb.Bound, orb.Polygon, orb.MultiPolygon:
		return false
	}

	for _, l := range ml {
		if lineContains(l, geom) {
			return true
		}
	}
	return false
}

// lineContains checks if line contains geometry.
// Line can contain only Point, MultiPoint, LineString, Ring, MultiLineString or GeometryCollection which contains their.
func lineContains(l orb.LineString, geom orb.Geometry) bool {
	switch g := geom.(type) {
	case orb.Point:
		return IsPointOnTheLine(l, g)
	case orb.MultiPoint:
		for _, p := range g {
			if !IsPointOnTheLine(l, p) {
				return false
			}
		}
		return true
	case orb.Ring:
		return lineContains(l, orb.LineString(g))
	case orb.MultiLineString:
		for _, l2 := range g {
			if !lineContains(l, l2) {
				return false
			}
		}
		return true
	case orb.Collection:
		for _, g := range g {
			if !lineContains(l, g) {
				return false
			}
		}
		return true
	case orb.LineString:
		// check all segments line b liegt on line a
		if len(l) < 2 || len(g) < 2 {
			return false
		}
		if !IsPointOnTheLine(l, g[0]) {
			return false
		}
		if !IsPointOnTheLine(l, g[len(g)-1]) {
			return false
		}
		i := 1 // next point on line a
		for j := 1; j < len(g); j++ {
			if i == len(l) {
				return false
			}
			if !vector.LineSegmentIsParallel(l[i-1], l[i], g[j-1], g[j]) {
				i++
				continue
			}
		}
		return true
	case orb.Bound, orb.Polygon, orb.MultiPolygon:
		return false
	}
	return false
}
