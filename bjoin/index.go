package bjoin

import (
	"context"
	"errors"

	"github.com/RoaringBitmap/roaring/roaring64"
)

// The package bjoin provides a data structure to collect and retrieve results of join two bitmap index.

// Index represent cross product matrix of two bitmap index
type Index struct {
	cp     *roaring64.Bitmap // cross product matrix when each number represent pair of elements a and b
	offset uint64            // len of b
}

// Pair represent pair of elements.
type Pair struct {
	A uint64
	B []uint64
}

// New returns new BitmapJoinIndex with given offset.
func New(offset uint64) *Index {
	return &Index{
		cp:     roaring64.New(),
		offset: offset,
	}
}

// IsEmpty returns true if index have no elements.
func (j *Index) IsEmpty() bool {
	return j.cp.IsEmpty()
}

//go:inline
func (j *Index) idxA(a uint64) uint64 {
	return a * (j.offset + 1)
}

//go:inline
func (j *Index) idx(a, b uint64) uint64 {
	return a*(j.offset+1) + b + 1
}

//go:inline
func (j *Index) a(idx uint64) uint64 {
	return idx / (j.offset + 1)
}

//go:inline
func (j *Index) b(idx uint64) (uint64, bool) {
	b := idx % (j.offset + 1)
	if b == 0 {
		return 0, false
	}
	return b - 1, true
}

//go:inline
func (j *Index) hasB(idx uint64) bool {
	return idx%(j.offset+1) != 0
}

// AddPairs add pairs of two bitmaps.
func (j *Index) AddPairs(a, b *roaring64.Bitmap) {
	itA := a.Iterator()
	for itA.HasNext() {
		a := itA.Next()
		j.cp.RemoveRange(j.idxA(a), j.idxA(a)+j.offset)
		if b == nil || b.IsEmpty() {
			j.cp.Add(j.idxA(a))
			continue
		}
		itB := b.Iterator()
		for itB.HasNext() {
			b := itB.Next()
			if b < j.offset {
				j.cp.Add(j.idx(a, b))
			}
		}
	}
}

// Clone() makes copy of bitmap join index.
func (j *Index) Clone() *Index {
	return &Index{
		offset: j.offset,
		cp:     j.cp.Clone(),
	}
}

// PairsGen returns channel of pairs.
func (j *Index) PairsGen(ctx context.Context) <-chan Pair {
	out := make(chan Pair)
	go func(ctx context.Context) {
		defer close(out)
		it := j.cp.Iterator()
		last := Pair{}
		first := true
		for it.HasNext() {
			idx := it.Next()
			a := j.a(idx)
			if !first && a != last.A {
				select {
				case <-ctx.Done():
					return
				case out <- last:
				}
				last.B = nil
			}
			last.A = a

			if b, ok := j.b(idx); ok {
				last.B = append(last.B, b)
			}
			if first {
				first = false
			}
		}
		if !first {
			select {
			case <-ctx.Done():
				return
			case out <- last:
			}
		}
	}(ctx)
	return out
}

// SingleGen returns channel with contained a elements that doesn't have join values.
func (j *Index) SingleGen(ctx context.Context) <-chan uint64 {
	out := make(chan uint64)
	go func(ctx context.Context) {
		defer close(out)
		it := j.cp.Iterator()
		for it.HasNext() {
			idx := it.Next()
			a := j.a(idx)
			if j.hasB(idx) {
				it.AdvanceIfNeeded((a+1)*j.offset + 1)
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- a:
			}
			it.AdvanceIfNeeded(idx + j.offset + 1)
		}
	}(ctx)
	return out
}

// ABGen returns channel with single pair of joined elements.
func (j *Index) ABGen(ctx context.Context) <-chan [2]uint64 {
	out := make(chan [2]uint64)
	go func(ctx context.Context) {
		defer close(out)
		it := j.cp.Iterator()
		for it.HasNext() {
			idx := it.Next()
			a := idx / (j.offset + 1)
			b := idx - a*(j.offset+1)
			select {
			case <-ctx.Done():
				return
			case out <- [2]uint64{a, b}:
			}
		}
	}(ctx)
	return out
}

var ErrDifferentOffset = errors.New("indexes have different offset")

// And computes the intersection between two bitmaps and stores the result in the current bitmap
func (j *Index) And(in *Index) error {
	if in == nil {
		j.cp = roaring64.New()
		return nil
	}
	if j.offset != in.offset {
		return ErrDifferentOffset
	}
	j.cp.And(in.cp)
	return nil
}

// Or computes the union between two bitmaps and stores the result in the current bitmap.
func (j *Index) Or(in *Index) error {
	if in == nil {
		return nil
	}
	if j.offset != in.offset {
		return ErrDifferentOffset
	}
	j.cp.Or(in.cp)
	return nil
}

// AndNot computes the difference between two bitmaps and stores the result in the current bitmap.
func (j *Index) AndNot(in *Index) error {
	if in == nil {
		return nil
	}
	if j.offset != in.offset {
		return ErrDifferentOffset
	}
	j.cp.AndNot(in.cp)
	return nil
}

// CrossJoin makes BitmapJoinIndex with pairs.
func CrossJoin(a, b *roaring64.Bitmap) *Index {
	j := &Index{
		cp:     roaring64.New(),
		offset: b.Maximum() + 1,
	}
	j.AddPairs(a, b)
	return j
}
