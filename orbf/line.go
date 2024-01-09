package orbf

import (
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// LineLength returns the length of LineString in WGS84 using Haversine method.
func LineLength(line orb.LineString) float64 {
	dist := 0.
	for i := 1; i < len(line); i++ {
		dist += geo.DistanceHaversine(line[i-1], line[i])
	}
	return dist
}

// LineInterpolatePoint interpolates point on the line.
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

	len := LineLength(l)
	dist := len * frac
	point, _ := geo.PointAtDistanceAlongLine(l, dist)

	return point
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

	lineLen := LineLength(line)
	if lineLen == 0 {
		return 0, 0
	}

	var find orb.Point
	index := 0
	minDist := math.MaxFloat64
	for i := 1; i < len(line); i++ {
		cp, dist := closestPoint(line[i-1], line[i], point)
		if math.Abs(dist) < math.Abs(minDist) {
			find = cp
			minDist = dist
			index = i
		}
	}

	return LineLength(append(line[0:index], find)) / lineLen, minDist
}

// IsPointOnTheLine returns true if point located on the line.
func IsPointOnTheLine(line orb.LineString, point orb.Point) bool {
	_, d := LineLocatePoint(line, point)
	return d == 0
}
