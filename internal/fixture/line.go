package fixture

import (
	_ "embed"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func LineString() []orb.LineString {
	return []orb.LineString{
		{
			{-122.4194, 37.7749}, // San Francisco, CA
			{-118.2437, 34.0522}, // Los Angeles, CA
			{-87.6298, 41.8781},  // Chicago, IL
			{-95.3698, 29.7604},  // Houston, TX
			{-75.1652, 39.9526},  // Philadelphia, PA
			{-74.0060, 40.7128},  // New York, NY
			{-0.1278, 51.5074},   // London, UK
			{2.3522, 48.8566},    // Paris, France
			{139.6917, 35.6895},  // Tokyo, Japan
			{151.2093, -33.8688}, // Sydney, Australia
		},
		{
			{2.3522, 48.8566},  // Paris, France
			{-0.1278, 51.5074}, // London, UK
			{4.9041, 52.3676},  // Amsterdam, Netherlands
			{13.4050, 52.5200}, // Berlin, Germany
			{18.6435, 60.1282}, // Stockholm, Sweden
			{12.4964, 41.9028}, // Rome, Italy
			{23.7275, 37.9838}, // Athens, Greece
			{31.2357, 30.0444}, // Istanbul, Turkey
		},
		{
			{13.4050, 52.5200}, // Berlin
			{9.9937, 53.5511},  // Hamburg
			{11.5810, 48.1351}, // Munich
			{6.9613, 50.9352},  // Cologne
			{9.1800, 48.7761},  // Stuttgart
			{8.6821, 50.1109},  // Frankfurt
			{10.5267, 52.3759}, // Hannover
		},
		{
			{13.4050, 52.5200}, // Berlin
			{9.9937, 53.5511},  // Hamburg
			{11.5810, 48.1351}, // Munich
		},
		{
			{-0.1223, 51.5073}, // Near Hungerford Bridge
			{-0.1233, 51.5072},
			{-0.1243, 51.5071},
			{-0.1253, 51.5069},
			{-0.1263, 51.5068}, // Near intersection with Northumberland Avenue
		},
	}
}

//go:embed osm_roads.geojson
var osmRoadsJSON []byte

func RoadsGeoJSON() (*geojson.FeatureCollection, error) {
	return geojson.UnmarshalFeatureCollection(osmRoadsJSON)
}

//go:embed lines.geojson
var linesJSON []byte

func LinesGeoJSON() (*geojson.FeatureCollection, error) {
	return geojson.UnmarshalFeatureCollection(linesJSON)
}
