package h3b

import (
	"context"
	"testing"

	"github.com/VGSML/geobin/bjoin"
	"github.com/uber/h3-go/v4"
	"golang.org/x/exp/slices"
)

func TestCheckIntersects(t *testing.T) {
	tests := []struct {
		name string
		a, b []h3.Cell
		want bool
	}{
		{
			name: "single equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			want: true,
		},
		{
			name: "single not equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff},
			want: false,
		},
		{
			name: "single inside single",
			a:    []h3.Cell{0x822baffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			want: true,
		},
		{
			name: "single inside several",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single inside several",
			a:    []h3.Cell{0x812b3ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several inside several",
			a:    []h3.Cell{0x822b87fffffffff, 0x822b8ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several equal",
			a:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several not equal",
			a:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x812e7ffffffffff},
			want: true,
		},
		{
			name: "several base several inside several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: true,
		},
		{
			name: "several base several intersects several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832673fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: true,
		},
		{
			name: "several base several intersects several different level",
			a:    []h3.Cell{0x8426635ffffffff, 0x8426449ffffffff, 0x842644dffffffff, 0x8426713ffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: true,
		},
		{
			name: "several base not equal",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff, 0x833baafffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: true,
		},
		{
			name: "several base not intersects same level",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff},
			b:    []h3.Cell{0x832aa3fffffffff, 0x832aa1fffffffff, 0x832aa5fffffffff, 0x832aa0fffffffff, 0x832aa4fffffffff, 0x832b99fffffffff},
			want: false,
		},
		{
			name: "several base not intersects different level",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff},
			b:    []h3.Cell{0x842b8c1ffffffff, 0x842b889ffffffff, 0x842b8ddffffffff, 0x842b8c3ffffffff, 0x842b8cbffffffff, 0x842b8c7ffffffff},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(15)
			for i, c := range tt.a {
				a.Insert(uint64(i), c)
			}
			b := New(15)
			for i, c := range tt.b {
				b.Insert(uint64(i), c)
			}
			got := CheckIntersection(a, b)
			if got != tt.want {
				t.Errorf("CheckIntersects returned unexpected result: got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestCheckContainsIn(t *testing.T) {
	tests := []struct {
		name string
		a, b []h3.Cell
		want bool
	}{
		{
			name: "single equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			want: true,
		},
		{
			name: "single not equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff},
			want: false,
		},
		{
			name: "single inside single",
			a:    []h3.Cell{0x822baffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			want: true,
		},
		{
			name: "single inside several",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single inside several",
			a:    []h3.Cell{0x812b3ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several inside several",
			a:    []h3.Cell{0x822b87fffffffff, 0x822b8ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several equal",
			a:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			b:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			want: true,
		},
		{
			name: "single base several not equal",
			a:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x812e7ffffffffff},
			want: false,
		},
		{
			name: "several base several inside several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: true,
		},
		{
			name: "several base several not inside several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832673fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: false,
		},
		{
			name: "several base not equal",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff, 0x833baafffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(15)
			for i, c := range tt.a {
				a.Insert(uint64(i), c)
			}
			b := New(15)
			for i, c := range tt.b {
				b.Insert(uint64(i), c)
			}
			got := CheckContainsIn(a, b)
			if got != tt.want {
				t.Errorf("CheckContainsIn returned unexpected result: got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestJoinIntersects(t *testing.T) {
	tests := []struct {
		name string
		a, b []h3.Cell
		left bool
		want []bjoin.Pair
	}{
		{
			name: "single equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			left: false,
			want: []bjoin.Pair{{A: 0, B: []uint64{0}}},
		},
		{
			name: "single equal left",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812bbffffffffff},
			left: true,
			want: []bjoin.Pair{{A: 0, B: []uint64{0}}},
		},
		{
			name: "single not equal",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff},
			left: false,
			want: []bjoin.Pair{},
		},
		{
			name: "single not equal left",
			a:    []h3.Cell{0x812bbffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff},
			left: true,
			want: []bjoin.Pair{{A: 0, B: nil}},
		},
		{
			name: "two items indexes intersects",
			a:    []h3.Cell{0x812bbffffffffff, 0x812b3ffffffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x8426713ffffffff},
			left: false,
			want: []bjoin.Pair{{A: 1, B: []uint64{0}}},
		},
		{
			name: "two items indexes intersects children",
			a:    []h3.Cell{0x812bbffffffffff, 0x85267123fffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x8426713ffffffff},
			left: false,
			want: []bjoin.Pair{{A: 1, B: []uint64{1}}},
		},
		{
			name: "two items indexes intersects children",
			a:    []h3.Cell{0x812bbffffffffff, 0x8426713ffffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x85267123fffffff},
			left: false,
			want: []bjoin.Pair{{A: 1, B: []uint64{1}}},
		},
		{
			name: "two items indexes intersects children left",
			a:    []h3.Cell{0x812bbffffffffff, 0x8426713ffffffff},
			b:    []h3.Cell{0x812b3ffffffffff, 0x85267123fffffff},
			left: true,
			want: []bjoin.Pair{{A: 0, B: nil}, {A: 1, B: []uint64{1}}},
		},
		{
			name: "several base several inside several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{1}},
				{A: 1, B: []uint64{1}},
				{A: 2, B: []uint64{2}},
				{A: 3, B: []uint64{2}},
			},
		},
		{
			name: "several base several intersects several",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832673fffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{1}},
				{A: 1, B: []uint64{1}},
				{A: 2, B: []uint64{2}},
			},
		},
		{
			name: "several base several intersects several different level",
			a:    []h3.Cell{0x8426635ffffffff, 0x8426449ffffffff, 0x842644dffffffff, 0x8426713ffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{0}},
			},
		},
		{
			name: "several base not equal",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff, 0x833baafffffffff},
			b:    []h3.Cell{0x822667fffffffff, 0x82274ffffffffff, 0x822ab7fffffffff},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{1}},
				{A: 1, B: []uint64{1}},
				{A: 2, B: []uint64{2}},
				{A: 3, B: []uint64{2}},
				{A: 4, B: []uint64{2}},
			},
		},
		{
			name: "several base not intersects same level",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff},
			b:    []h3.Cell{0x832aa3fffffffff, 0x832aa1fffffffff, 0x832aa5fffffffff, 0x832aa0fffffffff, 0x832aa4fffffffff, 0x832b99fffffffff},
			want: []bjoin.Pair{},
		},
		{
			name: "several base not intersects different level",
			a:    []h3.Cell{0x832748fffffffff, 0x83274dfffffffff, 0x832ab6fffffffff, 0x832ab2fffffffff, 0x832ab4fffffffff, 0x832b9bfffffffff},
			b:    []h3.Cell{0x842b8c1ffffffff, 0x842b889ffffffff, 0x842b8ddffffffff, 0x842b8c3ffffffff, 0x842b8cbffffffff, 0x842b8c7ffffffff},
			want: []bjoin.Pair{},
		},
		{
			name: "parent of res 4 many base with two cell",
			a: []h3.Cell{
				0x851205a3fffffff,
				0x851205a7fffffff,
				0x851205abfffffff,
				0x851205affffffff,
				0x851205b3fffffff,
				0x851205b7fffffff,
				0x851205bbfffffff,
				0x85136a43fffffff,
				0x85136a47fffffff,
				0x85136a4bfffffff,
				0x85136a4ffffffff,
				0x85136a53fffffff,
				0x85136a57fffffff,
				0x85136a5bfffffff,
				0x8448c43ffffffff,
				0x8448c45ffffffff,
				0x8448c47ffffffff,
				0x8448c49ffffffff,
				0x8448c4bffffffff,
				0x8448c4dffffffff,
			},
			b: []h3.Cell{0x831205fffffffff, 0x842aac1ffffffff},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{0}},
				{A: 1, B: []uint64{0}},
				{A: 2, B: []uint64{0}},
				{A: 3, B: []uint64{0}},
				{A: 4, B: []uint64{0}},
				{A: 5, B: []uint64{0}},
				{A: 6, B: []uint64{0}},
			},
		},
		{
			name: "two different cells set in one base",
			a: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			b: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			want: []bjoin.Pair{
				{A: 20, B: []uint64{0, 1, 2, 3}},
				{A: 21, B: []uint64{0, 1, 2, 3}},
			},
		},
		{
			name: "two different cells set in one base",
			a: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			b: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			left: true,
			want: []bjoin.Pair{
				{A: 0, B: nil},
				{A: 1, B: nil},
				{A: 2, B: nil},
				{A: 3, B: nil},
				{A: 4, B: nil},
				{A: 5, B: nil},
				{A: 6, B: nil},
				{A: 7, B: nil},
				{A: 8, B: nil},
				{A: 9, B: nil},
				{A: 10, B: nil},
				{A: 11, B: nil},
				{A: 12, B: nil},
				{A: 13, B: nil},
				{A: 14, B: nil},
				{A: 15, B: nil},
				{A: 16, B: nil},
				{A: 17, B: nil},
				{A: 18, B: nil},
				{A: 19, B: nil},
				{A: 20, B: []uint64{0, 1, 2, 3}},
				{A: 21, B: []uint64{0, 1, 2, 3}},
			},
		},
		{
			name: "two different cells set in one base",
			a: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			b: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			want: []bjoin.Pair{
				{A: 0, B: []uint64{20, 21}},
				{A: 1, B: []uint64{20, 21}},
				{A: 2, B: []uint64{20, 21}},
				{A: 3, B: []uint64{20, 21}},
			},
		},
		{
			name: "two different cells set in one base left",
			a: []h3.Cell{
				0x852aae17fffffff, 0x852aae07fffffff, 0x852aae3bfffffff, 0x852aae2bfffffff, 0x852aae6ffffffff,
				0x852aae67fffffff, 0x852aaad3fffffff, 0x852aaadbfffffff, 0x852a8523fffffff, 0x852a8537fffffff,
				0x852a852bfffffff, 0x852a8567fffffff, 0x852a856ffffffff, 0x852a81d3fffffff, 0x852a81dbfffffff,
				0x852a8037fffffff,
			},
			b: []h3.Cell{
				0x842aac1ffffffff, 0x842aacdffffffff, 0x842aac9ffffffff, 0x842a127ffffffff, 0x842a123ffffffff,
				0x842a135ffffffff, 0x842aacbffffffff, 0x842a13dffffffff, 0x842a12bffffffff, 0x842a121ffffffff,
				0x842a125ffffffff, 0x842aa1bffffffff, 0x842aa13ffffffff, 0x842aac5ffffffff, 0x842aae9ffffffff,
				0x842aaebffffffff, 0x842aac7ffffffff, 0x842aac3ffffffff, 0x842aa89ffffffff, 0x842aa8dffffffff,
				0x842aae3ffffffff, 0x842aae1ffffffff,
			},
			left: true,
			want: []bjoin.Pair{
				{A: 0, B: []uint64{20, 21}},
				{A: 1, B: []uint64{20, 21}},
				{A: 2, B: []uint64{20, 21}},
				{A: 3, B: []uint64{20, 21}},
				{A: 4, B: nil},
				{A: 5, B: nil},
				{A: 6, B: nil},
				{A: 7, B: nil},
				{A: 8, B: nil},
				{A: 9, B: nil},
				{A: 10, B: nil},
				{A: 11, B: nil},
				{A: 12, B: nil},
				{A: 13, B: nil},
				{A: 14, B: nil},
				{A: 15, B: nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(15)
			for i, c := range tt.a {
				a.Insert(uint64(i), c)
			}
			b := New(15)
			for i, c := range tt.b {
				b.Insert(uint64(i), c)
			}
			got := JoinIntersects(a, b, tt.left)
			testCheckJoinResult(got, tt.want, t.Errorf)
		})
	}
}

func testCheckJoinResult(got *bjoin.Index, want []bjoin.Pair, errorOut func(format string, args ...any)) bool {
	res := true
	var findPairs []bjoin.Pair
	for got := range got.PairsGen(context.Background()) {
		findPairs = append(findPairs, got)
		found := false
		var foundPair bjoin.Pair
		for _, check := range want {
			if check.A == got.A {
				found = true
				foundPair = check
				break
			}
		}
		if !found {
			errorOut("pair for item A (%d) not found", got.A)
			res = false
			continue
		}
		if slices.Compare(got.B, foundPair.B) != 0 {
			errorOut("pair for item A (%d) contains different elements of B, got %v, want %v", got.A, got.B, foundPair.B)
			res = false
		}
	}
	if len(findPairs) != len(want) {
		errorOut("len of pairs doesn't match, got %v, want %v", findPairs, want)
		res = true
	}
	return res
}
