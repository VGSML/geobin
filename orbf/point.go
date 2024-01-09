package orbf

import (
	"github.com/VGSML/geobin/orbf/planar"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/project"
)

// closestPoint calculates the closes point on the line between two points.
func closestPoint(p1, p2, find orb.Point) (orb.Point, float64) {
	a1 := geo.Bearing(p1, p2)
	a2 := geo.Bearing(p1, find)

	da := a2 - a1
	if da > 180 {
		da -= 360
	} else if da < -180 {
		da += 360
	}

	p1 = project.Point(p1, project.WGS84.ToMercator)
	p2 = project.Point(p2, project.WGS84.ToMercator)
	find2D := project.Point(find, project.WGS84.ToMercator)

	newPoint, _ := planar.ClosestPoint(p1, p2, find2D)
	direction := 0.
	if da > 0 {
		direction = 1
	}
	if da < 0 {
		direction = -1
	}

	newPoint = project.Point(newPoint, project.Mercator.ToWGS84)
	return newPoint, direction * geo.Distance(find, newPoint)
}
