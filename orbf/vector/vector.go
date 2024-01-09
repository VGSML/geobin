package vector

import (
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// Package to work with planar vector

func DotProduct(a, b orb.Point) float64 {
	return a.X()*b.X() + a.Y()*b.Y()
}

func SubtractPoints(a, b orb.Point) orb.Point {
	return orb.Point{a.X() - b.X(), a.Y() - b.Y()}
}

func AddPoints(a, b orb.Point) orb.Point {
	return orb.Point{a.X() + b.X(), a.Y() + b.Y()}
}

func ScalarMultiply(a orb.Point, s float64) orb.Point {
	return orb.Point{a.X() * s, a.Y() * s}
}

func ClosestPointOnLineSegment(p1, p2, find orb.Point) orb.Point {
	AP := SubtractPoints(find, p1)
	AB := SubtractPoints(p2, p1)

	// Compute the projection of AP onto AB.
	t := DotProduct(AP, AB) / DotProduct(AB, AB)

	// Check where the projection lies.
	if t < 0 {
		return p1
	}
	if t > 1 {
		return p2
	}
	return orb.Point{p1.X() + AB.X()*t, p1.Y() + AB.Y()*t}
}

// PointSide returns side the point of the line segment.
// if positiv - rights, negativ - left, zero - on the line.
func PointSide(p1, p2, find orb.Point) int {
	return int((find.X()-p1.X())*(p2.Y()-p1.Y()) - (p2.X()-p1.X())*(find.Y()-p1.Y()))
}

// PointOnSegment checks if point r lies on the line segment 'p1q1'
func PointOnSegment(p1, p2, find orb.Point) bool {
	if PointSide(p1, p2, find) != 0 {
		return false
	}
	return find.X() <= math.Max(p1.X(), p2.X()) && find.X() >= math.Min(p1.X(), p2.X()) &&
		find.Y() <= math.Max(p1.Y(), p2.Y()) && p2.Y() >= math.Min(p1.Y(), p2.Y())
}

// LineSegmentIsParallel test if line segment is parallel.
func LineSegmentIsParallel(p1, p2, p3, p4 orb.Point) bool {
	return (p4.Y()-p3.Y())*(p2.X()-p1.X()) == (p4.X()-p3.X())*(p2.Y()-p1.Y())
}

// ClosestPoint calculates the closes point on the line between two points
// in 2d geometries.
func ClosestPoint(p1, p2, find orb.Point) (orb.Point, float64) {
	newPoint := ClosestPointOnLineSegment(p1, p2, find)
	direction := 0.
	if PointSide(p1, p2, find) > 0 {
		direction = 1
	} else {
		direction = -1
	}

	return newPoint, planar.Distance(find, newPoint) * direction
}

// LineSegmentInterpolatePoint returns point between two points in 2d euclidean geometry.
func LineSegmentInterpolatePoint(p1, p2 orb.Point, frac float64) orb.Point {
	if frac == 0 {
		return p1
	}

	return orb.Point{
		p1.X() + frac*(p2.X()-p1.X()),
		p1.Y() + frac*(p2.Y()-p1.Y()),
	}
}

// IntersectionPoint finds intersection point of two line segments.
// If the line segments is parallel and has intersection return true and first intersection point.
func IntersectionPoint(p1, p2, p3, p4 orb.Point) (orb.Point, bool) {
	denom := (p4.Y()-p3.Y())*(p2.X()-p1.X()) - (p4.X()-p3.X())*(p2.Y()-p1.Y())

	// lines are parallel
	if denom == 0 {
		if pnt, dist := ClosestPoint(p1, p2, p3); dist == 0 {
			return pnt, true
		}
		if pnt, dist := ClosestPoint(p1, p2, p4); dist == 0 {
			return pnt, true
		}
		if pnt, dist := ClosestPoint(p3, p4, p1); dist == 0 {
			return pnt, true
		}
		if pnt, dist := ClosestPoint(p3, p4, p2); dist == 0 {
			return pnt, true
		}
		return orb.Point{}, false
	}

	ua := ((p4.X()-p3.X())*(p1.Y()-p3.Y()) - (p4.Y()-p3.Y())*(p1.X()-p3.X())) / denom
	ub := ((p2.X()-p1.X())*(p1.Y()-p3.Y()) - (p2.Y()-p1.Y())*(p1.X()-p3.X())) / denom

	// intersection point is not within the line segments
	if ua < 0 || ua > 1 || ub < 0 || ub > 1 {
		return orb.Point{}, false
	}

	// intersection point
	x := p1.X() + ua*(p2.X()-p1.X())
	y := p1.Y() + ua*(p2.Y()-p1.Y())

	return orb.Point{x, y}, true
}
