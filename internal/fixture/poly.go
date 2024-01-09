package fixture

import (
	_ "embed"

	"github.com/paulmach/orb/geojson"
)

//go:embed poly.geojson
var polyJSON []byte

func PolyGeoJSON() (*geojson.FeatureCollection, error) {
	return geojson.UnmarshalFeatureCollection(polyJSON)
}
