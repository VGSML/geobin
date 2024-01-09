package bjoin

import (
	"context"
	"testing"

	"github.com/RoaringBitmap/roaring/roaring64"
	"golang.org/x/exp/slices"
)

func TestIndex(t *testing.T) {
	a := roaring64.New()
	a.AddMany([]uint64{0, 1, 2, 3, 4, 5})
	join := New(11)
	join.AddPairs(a, nil)
	got := join.cp.ToArray()
	want := []uint64{0, 12, 24, 36, 48, 60}
	if slices.Compare(got, want) != 0 {
		t.Errorf("AddPairs() add different data to bitmap, got %v, want %v", got, want)
	}

	i := uint64(0)
	for pair := range join.PairsGen(context.Background()) {
		if pair.A != i {
			t.Errorf("pair has wrong a, got %d, want %d", pair.A, i)
		}
		if len(pair.B) != 0 {
			t.Errorf("pair for a=%d has wrong b, got %v, want %v", pair.A, pair.B, []uint64{})
		}
		i++
	}

	join = New(11)
	b := roaring64.New()
	bb := []uint64{0, 1, 2, 3, 4, 5}
	b.AddMany(bb)
	join.AddPairs(a, b)
	got = join.cp.ToArray()
	want = []uint64{1, 2, 3, 4, 5, 6, 13, 14, 15, 16, 17, 18, 25, 26, 27, 28, 29, 30, 37, 38, 39, 40, 41, 42, 49, 50, 51, 52, 53, 54, 61, 62, 63, 64, 65, 66}
	if slices.Compare(got, want) != 0 {
		t.Errorf("AddPairs() add different data to bitmap, got %v, want %v", got, want)
	}
	i = uint64(0)
	for pair := range join.PairsGen(context.Background()) {
		if pair.A != i {
			t.Errorf("pair has wrong a, got %d, want %d", pair.A, i)
		}
		if slices.Compare(pair.B, bb) != 0 {
			t.Errorf("pair for a=%d has wrong b, got %v, want %v", pair.A, pair.B, bb)
		}
		i++
	}

	join.AddPairs(a, nil)
	got = join.cp.ToArray()
	want = []uint64{0, 12, 24, 36, 48, 60}
	if slices.Compare(got, want) != 0 {
		t.Errorf("AddPairs() add different data to bitmap, got %v, want %v", got, want)
	}
	i = uint64(0)
	for pair := range join.PairsGen(context.Background()) {
		if pair.A != i {
			t.Errorf("pair has wrong a, got %d, want %d", pair.A, i)
		}
		if slices.Compare(pair.B, nil) != 0 {
			t.Errorf("pair for a=%d has wrong b, got %v, want %v", pair.A, pair.B, nil)
		}
		i++
	}

	join.AddPairs(a, b)
	a = roaring64.New()
	a.AddMany([]uint64{6, 7, 8, 9, 10})
	b = roaring64.New()
	b.AddMany([]uint64{6, 7, 8, 9, 10})
	bb2 := []uint64{6, 7, 8, 9, 10}
	join.AddPairs(a, b)
	i = uint64(0)
	for pair := range join.PairsGen(context.Background()) {
		if pair.A != i {
			t.Errorf("pair has wrong a, got %d, want %d", pair.A, i)
		}
		if i < 6 && slices.Compare(pair.B, bb) != 0 {
			t.Errorf("pair for a=%d has wrong b, got %v, want %v", pair.A, pair.B, bb)
		}
		if i > 5 && slices.Compare(pair.B, bb2) != 0 {
			t.Errorf("pair for a=%d has wrong b, got %v, want %v", pair.A, pair.B, bb2)
		}
		i++
	}

	if slices.Compare(join.cp.ToArray(), join.Clone().cp.ToArray()) != 0 {
		t.Errorf("Clone() add different data to bitmap, got %v, want %v", got, want)
	}

	join.cp.Clear()
	join.AddPairs(a, nil)
	a = roaring64.New()
	a.AddMany([]uint64{11, 12, 13, 14, 15})
	join.AddPairs(a, b)
	want = []uint64{6, 7, 8, 9, 10}
	i = 0
	for a := range join.SingleGen(context.Background()) {
		if a != want[i] {
			t.Errorf("pair has wrong a, got %d, want %d", a, want[i])
		}
		i++
	}
}
