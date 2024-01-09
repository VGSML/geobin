package geobin

import (
	"context"

	"github.com/RoaringBitmap/roaring"
	"github.com/VGSML/geobin/bjoin"
	"github.com/VGSML/geobin/h3b"
	"github.com/paulmach/orb"
	"github.com/uber/h3-go/v4"
)

// AggregationFunc is a function that will be called for each item that has cells in a given resolution.
type AggregationFunc func(cell []h3.Cell, idx int)

type ClusterFunc func(clusterIdx, idx int)

type ItemFunc func(ctx context.Context, item Item) error

// Projection data projection type.
type Projection int

const (
	WGS84 Projection = iota + 1
	Mercator
)

// Index provide method to searching and joining spatial data.
type Index struct {
	proj        Projection
	bitmap      *h3b.Index
	res         int
	newItemFunc func(idx int, geom orb.Geometry, res int, proj Projection) Item
	items       map[int]Item
}

type IndexOptions func(index *Index)

// Use full h3 indexed items instead only contains cells indexed item.
func WithIndexedItems(compact bool) IndexOptions {
	return func(index *Index) {
		index.newItemFunc = func(idx int, geom orb.Geometry, res int, proj Projection) Item {
			return newIndexedItem(idx, geom, res, compact, index.proj)
		}

	}
}

// Use custom h3 indexed items.
func WithCustomIndexedItems(newItemFunc func(idx int, geom orb.Geometry, res int, proj Projection) Item) IndexOptions {
	return func(index *Index) {
		index.newItemFunc = newItemFunc
	}
}

// Sets maximal h3 resolution (1..14).
func WithMaxResolution(res int) IndexOptions {
	return func(index *Index) {
		index.res = res
	}
}

// WithMercatorProjection data in Mercator projection.
func WithMercatorProjection() IndexOptions {
	return func(index *Index) {
		index.proj = Mercator
	}
}

// New creates new index with options.
func NewIndex(options ...IndexOptions) *Index {
	index := &Index{
		proj:        WGS84,
		res:         15,
		newItemFunc: newBoundIndexedItem,
		items:       map[int]Item{},
	}
	for _, opt := range options {
		opt(index)
	}
	index.bitmap = h3b.New(index.res)
	return index
}

// Insert adds element to index.
func (i *Index) Insert(idx int, item orb.Geometry) {
	indexItem := i.newItemFunc(idx, item, int(i.bitmap.Res()), i.proj)
	for _, cell := range indexItem.indexedCells() {
		i.bitmap.Insert(uint64(idx), cell)
	}
	i.items[idx] = indexItem
}

func (i *Index) Projection() Projection {
	return i.proj
}

// MaxItemIndex returns maximum item index.
func (i *Index) MaxItemIndex() int {
	return int(i.bitmap.MaxItemIndex())
}

// SetMaxItemIndex SetMaxItemIndex maximum item index.
// Sets and returns true if given index greatest or equal current maximum index.
func (i *Index) SetMaxItemIndex(idx int) bool {
	return i.bitmap.SetMaxItemIndex(uint64(idx))
}

// Remove delete indexed element.
func (i *Index) Remove(idx int) bool {
	delete(i.items, idx)
	return i.bitmap.Remove(uint64(idx))
}

// ContainsInItems returns items contains in given geometry
func (i *Index) ContainsInItems(ctx context.Context, in orb.Geometry) []int {
	inItem := i.newItemFunc(0, in, i.res, i.proj)
	m := i.bitmap.ContainsInItems(inItem.indexedCells())
	it := m.Iterator()
	var out []int
	for it.HasNext() {
		id := int(it.Next())
		item, ok := i.items[id]
		if !ok {
			continue
		}
		if item.ContainsIn(ctx, inItem) {
			out = append(out, id)
		}
	}
	return out
}

// IntersectionWith returns items that intersects with given geometry
func (i *Index) IntersectionWith(ctx context.Context, in orb.Geometry) []int {
	inItem := i.newItemFunc(0, in, i.res, i.proj)
	m := i.bitmap.Intersection(inItem.indexedCells())
	it := m.Iterator()
	var out []int
	for it.HasNext() {
		id := int(it.Next())
		item, ok := i.items[id]
		if !ok {
			continue
		}
		if item.Intersects(ctx, inItem) {
			out = append(out, id)
		}
	}
	return out
}

// JoinIntersects perform intersection join operations of two indexes.
func (i *Index) JoinIntersects(ctx context.Context, right *Index, left bool) *bjoin.Index {
	return h3b.JoinIntersects(i.bitmap, right.bitmap, left)
}

func (i *Index) JoinContains(ctx context.Context, right *Index, left bool) *bjoin.Index {

	return nil
}

// AggregateH3 perform aggregation of indexed items by h3 cells for given resolution.
// for each cell in given resolution, call given aggregation function with cell, item index and count of item cells in given resolution.
func (i *Index) AggregateByRes(ctx context.Context, res int, aggFunc AggregationFunc) error {
	return nil
}

// ClusterIntersects perform clustering of indexed items by intersection of given geometry.
func (i *Index) ClusterIntersects(ctx context.Context, clusterFunc ClusterFunc) error {
	return nil
}

// ClusterKMeans perform clustering of indexed items by k-means algorithm.
func (i *Index) ClusterKMeans(ctx context.Context, clusters int, clusterFunc ClusterFunc) error {
	return nil
}

// ClusterKMeansWithSeeds perform clustering of indexed items by k-means algorithm with given seeds.
func (i *Index) ClusterKMeansWithSeeds(ctx context.Context, seeds []h3.Cell, clusterFunc ClusterFunc) error {
	return nil
}

// ClusterDBSCAN perform clustering of indexed items by DBSCAN algorithm.
func (i *Index) ClusterDBSCAN(ctx context.Context, eps float64, minPts int, clusterFunc ClusterFunc) error {
	return nil
}

// FilterBitmap returns new index with items that intersects with given bitmap.
func (i *Index) FilterBitmap(ctx context.Context, bitmap *roaring.Bitmap) *Index {
	return nil
}

// ApplyItemFunc applies given function to each item.
// For example, this function can be: Centroid, Calculating length of geometry, Buffer and so on.
func (i Index) ApplyItemFunc(ctx context.Context, itemFunc ItemFunc) error {
	return nil
}
