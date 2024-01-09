package planar

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestPolyContains(t *testing.T) {
	polygon1 := orb.Polygon{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}, // outer ring
		{{3, 3}, {3, 6}, {6, 6}, {6, 3}, {3, 3}},     // hole
	}

	tests := []struct {
		name     string
		geom     orb.Geometry
		expected bool
	}{
		{
			name:     "Point inside  polygon",
			geom:     orb.Point{1, 1},
			expected: true,
		},
		{
			name:     "Point inside hole polygon",
			geom:     orb.Point{5, 5},
			expected: false,
		},
		{
			name:     "Point outside polygon",
			geom:     orb.Point{11, 11},
			expected: false,
		},
		{
			name:     "Point inside hole of polygon",
			geom:     orb.Point{4, 4},
			expected: false,
		},
		{
			name:     "MultiPoint inside of polygon",
			geom:     orb.MultiPoint{{1, 1}, {9, 9}},
			expected: true,
		},
		{
			name:     "MultiPoint inside of polygon, and point inside of hole",
			geom:     orb.MultiPoint{{1, 1}, {9, 9}, {5, 5}},
			expected: false,
		},
		{
			name:     "Polygon contained within input LineString",
			geom:     orb.LineString{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
			expected: true,
		},
		{
			name:     "Polygon intersects within input LineString",
			geom:     orb.LineString{{1, 1}, {2, 2}, {5, 5}},
			expected: false,
		},
		{
			name: "Polygon contained within input MultiLineString",
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{8, 8}, {9, 9}},
			},
			expected: true,
		},
		{
			name: "Polygon intersects within input MultiLineString",
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{8, 8}, {9, 9}, {11, 11}},
			},
			expected: false,
		},
		{
			name:     "Polygon intersecting with input polygon",
			geom:     orb.Polygon{{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}},
			expected: false,
		},
		{
			name:     "Polygon contained within input polygon",
			geom:     orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
			expected: true,
		},
		{
			name:     "Polygon contained within input polygon but holes inside",
			geom:     orb.Polygon{{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}}},
			expected: false,
		},
		{
			name: "Polygon contained within input polygon and holes inside",
			geom: orb.Polygon{
				{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
				{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
			},
			expected: true,
		},
		{
			name: "MultiPolygon contained within input polygon",
			geom: orb.MultiPolygon{
				{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
				},
			},
			expected: true,
		},
		{
			name: "MultiPolygon contained within input polygon",
			geom: orb.MultiPolygon{
				{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 5}, {5, 5}, {5, 2}, {2, 2}},
				},
			},
			expected: false,
		},
		{
			name: "Collection contained within input polygon",
			geom: orb.Collection{
				orb.Point{1, 1},
				orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				orb.Polygon{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
				},
			},
			expected: true,
		},
		{
			name: "Collection contained within input polygon",
			geom: orb.Collection{
				orb.Point{1, 1},
				orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				orb.Polygon{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 5}, {5, 5}, {5, 2}, {2, 2}},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := polyContains(polygon1, tt.geom)
			if tt.expected != got {
				t.Errorf("polyContains return wrong value got %v want %v", got, tt.expected)
			}
		})
	}
}

func TestPolyIntersects(t *testing.T) {
	polygon1 := orb.Polygon{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}, // outer ring
		{{3, 3}, {3, 6}, {6, 6}, {6, 3}, {3, 3}},     // hole
	}

	tests := []struct {
		name     string
		geom     orb.Geometry
		expected bool
	}{
		{
			name:     "Point inside  polygon",
			geom:     orb.Point{1, 1},
			expected: true,
		},
		{
			name:     "Point inside hole polygon",
			geom:     orb.Point{5, 5},
			expected: false,
		},
		{
			name:     "Point outside polygon",
			geom:     orb.Point{11, 11},
			expected: false,
		},
		{
			name:     "Point inside hole of polygon",
			geom:     orb.Point{4, 4},
			expected: false,
		},
		{
			name:     "MultiPoint inside of polygon",
			geom:     orb.MultiPoint{{1, 1}, {9, 9}},
			expected: true,
		},
		{
			name:     "MultiPoint inside of polygon, and point inside of hole",
			geom:     orb.MultiPoint{{5, 5}, {1, 1}, {9, 9}},
			expected: true,
		},
		{
			name:     "Polygon contained within input LineString",
			geom:     orb.LineString{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
			expected: true,
		},
		{
			name:     "Polygon intersects within input LineString",
			geom:     orb.LineString{{1, 1}, {2, 2}, {5, 5}},
			expected: true,
		},
		{
			name:     "Polygon does not intersects within input LineString",
			geom:     orb.LineString{{11, 11}, {13, 13}, {18, 25}},
			expected: false,
		},
		{
			name: "Polygon contained within input MultiLineString",
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{8, 8}, {9, 9}},
			},
			expected: true,
		},
		{
			name: "Polygon intersects within input MultiLineString",
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{8, 8}, {9, 9}, {11, 11}},
			},
			expected: true,
		},
		{
			name:     "Polygon intersecting with input polygon",
			geom:     orb.Polygon{{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}},
			expected: true,
		},
		{
			name:     "Polygon is inside in hole and intersects",
			geom:     orb.Polygon{{{4, 4}, {7, 4}, {7, 5}, {4, 5}, {4, 4}}},
			expected: true,
		},
		{
			name:     "Polygon is inside in hole",
			geom:     orb.Polygon{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}},
			expected: false,
		},
		{
			name:     "Polygon is outside polygon",
			geom:     orb.Polygon{{{20, 20}, {30, 20}, {30, 30}, {20, 30}, {30, 30}}},
			expected: false,
		},
		{
			name:     "Polygon contained within input polygon",
			geom:     orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
			expected: true,
		},
		{
			name:     "Polygon contained within input polygon but holes inside",
			geom:     orb.Polygon{{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}}},
			expected: true,
		},
		{
			name: "Polygon contained within input polygon and holes inside",
			geom: orb.Polygon{
				{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
				{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
			},
			expected: true,
		},
		{
			name:     "Polygon intersects within input polygon",
			geom:     orb.Polygon{{{5, 15}, {15, 5}, {20, 20}}},
			expected: true,
		},
		{
			name: "MultiPolygon contained within input polygon",
			geom: orb.MultiPolygon{
				{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
				},
			},
			expected: true,
		},
		{
			name: "MultiPolygon contained within input polygon",
			geom: orb.MultiPolygon{
				{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 5}, {5, 5}, {5, 2}, {2, 2}},
				},
			},
			expected: true,
		},
		{
			name: "Collection contained within input polygon",
			geom: orb.Collection{
				orb.Point{1, 1},
				orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				orb.Polygon{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 7}, {7, 7}, {7, 2}, {2, 2}},
				},
			},
			expected: true,
		},
		{
			name: "Collection contained within input polygon",
			geom: orb.Collection{
				orb.Point{1, 1},
				orb.Polygon{{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}},
				orb.Polygon{
					{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
					{{2, 2}, {2, 5}, {5, 5}, {5, 2}, {2, 2}},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := polyIntersects(polygon1, tt.geom)
			if tt.expected != got {
				t.Errorf("polyIntersects return wrong value got %v want %v", got, tt.expected)
			}
		})
	}
}
