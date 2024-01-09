package vector

import (
	"math"

	"github.com/paulmach/orb"
)

// epsilon is the tolerance used for float comparison
const epsilon = 1e-5

// FloatEqual checks if two float64 numbers are equal within a tolerance.
func FloatEqual(a, b float64, digits ...int) bool {
	if len(digits) != 0 {
		return math.Abs(a-b) < 1/math.Pow10(-digits[0])
	}
	return math.Abs(a-b) < epsilon
}

// PointEqual checks if two points are equal within a tolerance.
func PointEqual(p1, p2 orb.Point) bool {
	return FloatEqual(p1[0], p2[0]) && FloatEqual(p1[1], p2[1])
}
