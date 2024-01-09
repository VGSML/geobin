package planar

import (
	"math"
	"testing"

	"github.com/VGSML/geobin/internal/fixture"
	"github.com/VGSML/geobin/orbf/vector"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

func TestLineInterpolatePointMercator(t *testing.T) {
	line := fixture.LineString()[3]
	tests := []struct {
		name string
		frac float64
		want orb.Point
	}{
		{
			name: "Negative fraction",
			frac: -0.1,
			want: orb.Point{},
		},
		{
			name: "Fraction exceeds 1",
			frac: 1.1,
			want: orb.Point{},
		},
		{
			name: "Start of the line",
			frac: 0,
			want: line[0],
		},
		{
			name: "End of the line",
			frac: 1,
			want: line[len(line)-1],
		},
		{
			name: "17% of the line",
			frac: 0.17,
			want: orb.Point{11.4982673, 53.099368}, // Approximate value
		},
		{
			name: "Quarter of the line",
			frac: 0.25,
			want: orb.Point{10.6009814, 53.369345108}, // Approximate value
		},
		{
			name: "Half of the line",
			frac: 0.5,
			want: orb.Point{10.4404625, 52.090204}, // Approximate value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := fixture.LineString()[3]
			lt := project.LineString(line, project.WGS84.ToMercator)
			got := LineInterpolatePoint(lt, tt.frac)
			got = project.Point(got, project.Mercator.ToWGS84)
			if !vector.PointEqual(got, tt.want) {
				t.Errorf("LineInterpolatePointMercator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineLocatePoint(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name         string
		line         orb.LineString
		point        orb.Point
		expectedFrac float64
		expectedDist float64
	}{
		{
			name:         "Simple case 1",
			line:         orb.LineString{{0, 0}, {10, 10}},
			point:        orb.Point{5, 5},
			expectedFrac: 0.5,
			expectedDist: 0,
		},
		{
			name:         "Simple case 2",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{5, 5},
			expectedFrac: 0.5,
			expectedDist: -5,
		},
		{
			name:         "Simple case 3",
			line:         orb.LineString{{0, 0}, {0, 10}},
			point:        orb.Point{5, 5},
			expectedFrac: 0.5,
			expectedDist: 5,
		},
		{
			name:         "Simple case 4",
			line:         orb.LineString{{0, 0}, {0, 10}, {5, 10}, {15, 10}},
			point:        orb.Point{5, 5},
			expectedFrac: 0.2,
			expectedDist: 5,
		},
		{
			name:         "Simple case 5",
			line:         orb.LineString{{0, 0}, {0, 10}, {5, 10}, {15, 10}},
			point:        orb.Point{9, 12},
			expectedFrac: 0.76,
			expectedDist: -2,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultFrac, resultDist := LineLocatePoint(tc.line, tc.point)
			if !vector.FloatEqual(resultFrac, tc.expectedFrac, 4) || !vector.FloatEqual(resultDist, tc.expectedDist, 2) {
				t.Errorf("expected: (%v, %v), got: (%v, %v)", tc.expectedFrac, tc.expectedDist, resultFrac, resultDist)
			}
			pntOnLine := LineInterpolatePoint(tc.line, resultFrac)
			dist := planar.Distance(pntOnLine, tc.point)
			if !vector.FloatEqual(math.Abs(resultDist), dist, 2) {
				t.Errorf("expected: (%v, %v), got: (%v, %v)", tc.expectedFrac, dist, resultFrac, resultDist)
			}
		})
	}
}

func TestLinesIntersectPoint(t *testing.T) {
	testCases := []struct {
		name         string
		line1, line2 orb.LineString
		expected     orb.Point
		ok           bool
	}{
		{
			line1:    orb.LineString{{0, 0}, {1, 1}},
			line2:    orb.LineString{{0, 1}, {1, 0}},
			expected: orb.Point{0.5, 0.5},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {2, 2}},
			line2:    orb.LineString{{0, 2}, {2, 0}},
			expected: orb.Point{1, 1},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {3, 3}},
			line2:    orb.LineString{{0, 1}, {1, 0}},
			expected: orb.Point{0.5, 0.5},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {5, 5}},
			line2:    orb.LineString{{0, 5}, {5, 5}},
			expected: orb.Point{5, 5},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {5, 5}},
			line2:    orb.LineString{{3, 3}, {6, 6}},
			expected: orb.Point{3, 3},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {5, 5}, {10, 10}},
			line2:    orb.LineString{{3, 3}, {6, 6}},
			expected: orb.Point{3, 3},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {5, 5}, {10, 10}},
			line2:    orb.LineString{{3, 3}, {6, 6}, {12, 12}},
			expected: orb.Point{3, 3},
			ok:       true,
		},
		{
			line1:    orb.LineString{{0, 0}, {5, 5}, {10, 10}},
			line2:    orb.LineString{{3, 3}, {6, 6}, {20, 12}},
			expected: orb.Point{3, 3},
			ok:       true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := LinesIntersectPoint(tt.line1, tt.line2)
			if !vector.PointEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
			if ok != tt.ok {
				t.Errorf("expected %v, got %v", tt.ok, ok)
			}
		})
	}
}

func Test_lineContains(t *testing.T) {
	tests := []struct {
		name   string
		line   orb.LineString
		geom   orb.Geometry
		expect bool
	}{
		{
			name:   "Line contains point",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{5, 5},
			expect: true,
		},
		{
			name:   "Line does not contain point",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{20, 20},
			expect: false,
		},
		{
			name:   "Line contains multipoint",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.MultiPoint{{1, 1}, {2, 2}, {3, 3}},
			expect: true,
		},
		{
			name:   "Line does not contain multipoint",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{20, 20},
			expect: false,
		},
		{
			name:   "Line contains line",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.LineString{{5, 5}, {10, 10}},
			expect: true,
		},
		{
			name:   "Complex line contains line",
			line:   orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom:   orb.LineString{{5, 5}, {10, 10}},
			expect: true,
		},
		{
			name:   "Complex line contains complex line",
			line:   orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom:   orb.LineString{{5, 5}, {9, 9}, {10, 10}, {12, 10}},
			expect: true,
		},
		{
			name: "Complex line contains multi line",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.MultiLineString{
				{{5, 5}, {9, 9}},
				{{11, 10}, {12, 10}},
			},
			expect: true,
		},
		{
			name: "Complex line does not contain multi line",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.MultiLineString{
				{{5, 5}, {9, 9}},
				{{11, 10}, {12, 12}},
			},
			expect: false,
		},
		{
			name: "Complex line contains collection",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.Collection{
				orb.LineString{{5, 5}, {9, 9}},
				orb.LineString{{11, 10}, {12, 10}},
				orb.Point{3, 3},
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lineContains(tt.line, tt.geom)
			if result != tt.expect {
				t.Errorf("Expected %v but got %v", tt.expect, result)
			}
		})
	}
}

func Test_multiLineContains(t *testing.T) {
	tests := []struct {
		name      string
		multiLine orb.MultiLineString
		geom      orb.Geometry
		expect    bool
	}{
		{
			name: "MultiLine contains multiPoint",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom:   orb.MultiPoint{{1, 1}, {2, 2}, {21, 21}, {29, 29}},
			expect: true,
		},
		{
			name: "MultiLine does not contain multiPoint",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom:   orb.MultiPoint{{1, 1}, {20, 20}, {50, 50}},
			expect: false,
		},
		{
			name: "MultiLine contains LineString",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom:   orb.LineString{{21, 21}, {29, 29}},
			expect: true,
		},
		{
			name: "MultiLine does not contain LineString",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom:   orb.LineString{{1, 1}, {2, 2}, {21, 21}, {22, 23}},
			expect: false,
		},
		{
			name: "MultiLine contains LineString",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{21, 21}, {22, 22}},
			},
			expect: true,
		},
		{
			name: "MultiLine does not contain LineString",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom: orb.MultiLineString{
				{{1, 1}, {2, 2}},
				{{21, 21}, {15, 20}},
			},
			expect: false,
		},
		{
			name: "MultiLine contains Collection",
			multiLine: orb.MultiLineString{
				{{0, 0}, {10, 10}},
				{{20, 20}, {30, 30}},
			},
			geom: orb.Collection{
				orb.MultiPoint{{1, 1}, {2, 2}},
				orb.LineString{{21, 21}, {22, 22}},
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := multiLineContains(tt.multiLine, tt.geom)
			if result != tt.expect {
				t.Errorf("Expected %v but got %v", tt.expect, result)
			}
		})
	}
}

func Test_lineIntersects(t *testing.T) {
	tests := []struct {
		name   string
		line   orb.LineString
		geom   orb.Geometry
		expect bool
	}{
		{
			name:   "Line contains point",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{5, 5},
			expect: true,
		},
		{
			name:   "Line does not contain point",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{20, 20},
			expect: false,
		},
		{
			name:   "Line contains multipoint",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.MultiPoint{{1, 1}, {2, 2}, {3, 3}},
			expect: true,
		},
		{
			name:   "Line does not contain multipoint",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.Point{20, 20},
			expect: false,
		},
		{
			name:   "Line contains line",
			line:   orb.LineString{{0, 0}, {10, 10}},
			geom:   orb.LineString{{5, 5}, {10, 10}},
			expect: true,
		},
		{
			name:   "Complex line contains line",
			line:   orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom:   orb.LineString{{5, 5}, {10, 10}},
			expect: true,
		},
		{
			name:   "Complex line contains complex line",
			line:   orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom:   orb.LineString{{5, 5}, {9, 9}, {10, 10}, {12, 10}},
			expect: true,
		},
		{
			name:   "Complex line does not intersect line",
			line:   orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom:   orb.LineString{{1, 5}, {3, 9}, {5, 15}, {12, 25}},
			expect: false,
		},
		{
			name: "Complex line contains multi line",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.MultiLineString{
				{{5, 5}, {9, 9}},
				{{11, 10}, {12, 10}},
			},
			expect: true,
		},
		{
			name: "Complex line does not contain multi line",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.MultiLineString{
				{{5, 5}, {9, 9}},
				{{11, 10}, {12, 12}},
			},
			expect: true,
		},
		{
			name: "Complex line contains collection",
			line: orb.LineString{{0, 0}, {10, 10}, {15, 10}},
			geom: orb.Collection{
				orb.LineString{{5, 5}, {9, 9}},
				orb.LineString{{11, 10}, {12, 10}},
				orb.Point{3, 3},
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lineIntersects(tt.line, tt.geom)
			if result != tt.expect {
				t.Errorf("Expected %v but got %v", tt.expect, result)
			}
		})
	}
}
