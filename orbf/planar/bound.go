package planar

import "github.com/paulmach/orb"

// boundContains check if bound completely contain geometry.
func boundContains(b orb.Bound, geom orb.Geometry) bool {
	if !b.Contains(geom.Bound().Min) {
		return false
	}
	return b.Contains(geom.Bound().Max)
}
