package vector

import (
	"strconv"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/project"
)

func TestClosestPointOnLineSegment(t *testing.T) {
	testCases := []struct {
		name string
		p1   orb.Point
		p2   orb.Point
		find orb.Point
		want orb.Point
	}{
		{
			name: "right corner",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 0},
			find: orb.Point{1, 1},
			want: orb.Point{1, 0},
		},
		{
			name: "left corner",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 0},
			find: orb.Point{0, 1},
			want: orb.Point{0, 0},
		},
		{
			name: "middle corner",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 0},
			find: orb.Point{0.5, 1},
			want: orb.Point{0.5, 0},
		},
		{
			name: "middle center 2",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 0},
			find: orb.Point{0.5, 0.5},
			want: orb.Point{0.5, 0},
		},
		{
			name: "45 degree center 2",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 1},
			find: orb.Point{0, 1},
			want: orb.Point{0.5, 0.5},
		},
		{
			name: "45 degree quarter",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 1},
			find: orb.Point{0, 0.5},
			want: orb.Point{0.25, 0.25},
		},
		{
			name: "45 degree 3 quarter",
			p1:   orb.Point{0, 0},
			p2:   orb.Point{1, 1},
			find: orb.Point{0, 1.5},
			want: orb.Point{0.75, 0.75},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClosestPointOnLineSegment(tt.p1, tt.p2, tt.find); !PointEqual(got, tt.want) {
				t.Errorf("ClosestPointOnLineSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClosestPoint(t *testing.T) {
	testCases := []struct {
		name   string
		p1, p2 orb.Point
		point  orb.Point
		want   orb.Point
		dist   float64
	}{
		{
			name:  "Closest point between Berlin and Hamburg left",
			p1:    orb.Point{13.4050, 52.5200}, // Berlin, Germany
			p2:    orb.Point{9.9937, 53.5511},  // Hamburg, Germany
			point: orb.Point{10.6866, 52.8124},
			want:  orb.Point{11.0414872453, 53.23701764}, // Approximate value
			dist:  -87959.517,                            // Approximate distance in meters
		},
		{
			name:  "Closest point between Berlin and Hamburg right",
			p1:    orb.Point{13.4050, 52.5200}, // Berlin, Germany
			p2:    orb.Point{9.9937, 53.5511},  // Hamburg, Germany
			point: orb.Point{10.6866, 53.8124},
			want:  orb.Point{10.369783227, 53.438632758}, // Approximate value
			dist:  78523.674,                             // Approximate distance in meters
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p1 := project.Point(tt.p1, project.WGS84.ToMercator)
			p2 := project.Point(tt.p2, project.WGS84.ToMercator)
			point := project.Point(tt.point, project.WGS84.ToMercator)
			got, dist := ClosestPoint(p1, p2, point)
			got = project.Point(got, project.Mercator.ToWGS84)
			if !PointEqual(got, tt.want) || !FloatEqual(dist, tt.dist, 2) {
				t.Errorf("closestPointWGS84() = %v - %v, want %v - %v", got, dist, tt.want, tt.dist)
			}
		})
	}
}

func TestIntersectionPoint(t *testing.T) {
	testCases := []struct {
		p1, p2, p3, p4  orb.Point
		expected        orb.Point
		hasIntersection bool
	}{
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{4, 4},
			p3:              orb.Point{1, 4},
			p4:              orb.Point{4, 1},
			expected:        orb.Point{2.5, 2.5},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{1, 4},
			p3:              orb.Point{1, 4},
			p4:              orb.Point{4, 1},
			expected:        orb.Point{1, 4},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{4, 4},
			p3:              orb.Point{6, 6},
			p4:              orb.Point{9, 9},
			expected:        orb.Point{},
			hasIntersection: false,
		},
		{
			p1:              orb.Point{2, 2},
			p2:              orb.Point{4, 4},
			p3:              orb.Point{1, 4},
			p4:              orb.Point{3, 2},
			expected:        orb.Point{2.5, 2.5},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 2},
			p2:              orb.Point{5, 2},
			p3:              orb.Point{3, 4},
			p4:              orb.Point{3, 1},
			expected:        orb.Point{3, 2},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 2},
			p2:              orb.Point{5, 2},
			p3:              orb.Point{5, 2},
			p4:              orb.Point{7, 2},
			expected:        orb.Point{5, 2},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{5, 5},
			p4:              orb.Point{7, 2},
			expected:        orb.Point{5, 5},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{6, 6},
			p4:              orb.Point{10, 10},
			expected:        orb.Point{},
			hasIntersection: false,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{5, 5},
			p4:              orb.Point{10, 10},
			expected:        orb.Point{5, 5},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{4, 4},
			p4:              orb.Point{10, 10},
			expected:        orb.Point{4, 4},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 5},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{4, 5},
			p4:              orb.Point{10, 5},
			expected:        orb.Point{4, 5},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{5, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{5, 4},
			p4:              orb.Point{5, 10},
			expected:        orb.Point{5, 4},
			hasIntersection: true,
		},
		{
			p1:              orb.Point{1, 1},
			p2:              orb.Point{5, 5},
			p3:              orb.Point{2, 2},
			p4:              orb.Point{7, 2},
			expected:        orb.Point{2, 2},
			hasIntersection: true,
		},
	}

	for i, tt := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, hasIntersection := IntersectionPoint(tt.p1, tt.p2, tt.p3, tt.p4)
			if got != tt.expected {
				t.Errorf("intersectionPoint(%v, %v, %v, %v) == %v, want %v", tt.p1, tt.p2, tt.p3, tt.p4, got, tt.expected)
			}
			if hasIntersection != tt.hasIntersection {
				t.Errorf("Expected hasIntersection %v, got %v", tt.hasIntersection, hasIntersection)
			}
		})
	}
}
