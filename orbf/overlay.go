package orbf

import (
	"github.com/VGSML/geobin/orbf/planar"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/project"
)

func Intersects(a, b orb.Geometry) bool {
	geom1 := projectToMercator(a)
	if geom1 == nil {
		return false
	}
	geom2 := projectToMercator(b)
	if geom1 == nil {
		return false
	}
	return planar.Intersects(geom1, geom2)
}

func Contains(a, b orb.Geometry) bool {
	geom1 := projectToMercator(a)
	if geom1 == nil {
		return false
	}
	geom2 := projectToMercator(b)
	if geom1 == nil {
		return false
	}
	return planar.Contains(geom1, geom2)
}

func projectToMercator(geom orb.Geometry) orb.Geometry {
	switch g := geom.(type) {
	case orb.Bound, orb.Point:
		return project.Geometry(g, project.WGS84.ToMercator)
	case orb.MultiPoint:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.LineString:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.MultiLineString:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.Polygon:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.MultiPolygon:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.Collection:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	case orb.Ring:
		return project.Geometry(g.Clone(), project.WGS84.ToMercator)
	default:
		return nil
	}
}
