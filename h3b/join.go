package h3b

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/VGSML/geobin/bjoin"
)

// perform inner or left join by operation Intersects and Contains

// CheckIntersection checks that a and b have intersections (minimum 1 same element or element inside).
func CheckIntersection(a, b *Index) bool {
	base := roaring64.And(a.baseCellsMask, b.baseCellsMask)
	if base.IsEmpty() {
		return false
	}
	it := base.Iterator() // intersection of base cells
	for it.HasNext() {
		bn := it.Next()
		bma := a.baseCellMap[bn]
		if bma == nil || bma.IsEmpty() {
			continue
		}
		bmb := b.baseCellMap[bn]
		if bmb == nil || bmb.IsEmpty() {
			continue
		}
		baseA := bma.Clone()
		baseB := bmb.Clone()
		for res := 0; res < 15; res++ {
			resA := roaring64.New()
			resB := roaring64.New()
			allA := false
			allB := false
			for cn := 7; cn >= 0; cn-- {
				cmb := b.resMaps[res][cn]
				emptyB := cmb == nil || cmb.IsEmpty()
				if !emptyB {
					cmb = roaring64.And(baseB, cmb)
					emptyB = cmb.IsEmpty()
				}
				if cn == 7 && !emptyB {
					allB = true
				}
				cma := a.resMaps[res][cn]
				emptyA := cma == nil || cma.IsEmpty()
				if !emptyA {
					cma = roaring64.And(baseA, cma)
					emptyA = cma.IsEmpty()

				}
				if cn == 7 && !emptyA {
					allA = true
				}
				if allA && !emptyB ||
					allB && !emptyA {
					return true
				}
				if !emptyA && !emptyB {
					resA.Or(cma)
					resB.Or(cmb)
				}
			}
			if resA.IsEmpty() || resB.IsEmpty() {
				break
			}
			baseA = resA
			baseB = resB
		}
	}
	return false
}

// CheckContainsIn checks that all items in a contains in b.
func CheckContainsIn(a, b *Index) bool {
	if !roaring64.Xor(a.baseCellsMask, b.baseCellsMask).IsEmpty() {
		return false
	}
	it := a.baseCellsMask.Iterator()
	for it.HasNext() {
		bn := it.Next()
		bmb := b.baseCellMap[bn]
		bma := a.baseCellMap[bn]
		emptyB := (bmb == nil || bmb.IsEmpty())
		emptyA := (bma == nil || bma.IsEmpty())
		if emptyA && emptyB {
			continue
		}
		// A has items outside base cell of B
		if emptyB && !emptyA {
			return false
		}
		if emptyA { // no check empty A and non empty B
			continue
		}

		baseA := bma.Clone()
		baseB := bmb.Clone()
		allB := false
		for res := 0; res < 15; res++ {
			resA := roaring64.New()
			resB := roaring64.New()
			for cn := 7; cn >= 0; cn-- {
				cmb := b.resMaps[res][cn]
				emptyB = cmb == nil || cmb.IsEmpty()
				if !emptyB {
					cmb = roaring64.And(baseB, cmb)
					emptyB = cmb.IsEmpty()
					resB.Or(cmb)
				}
				cma := a.resMaps[res][cn]
				emptyA = cma == nil || cma.IsEmpty()
				if !emptyA {
					cma = roaring64.And(baseA, cma)
					emptyA = cma.IsEmpty()
					resA.Or(cma)
				}
				if cn == 7 && !emptyB {
					allB = true
					if !emptyA {
						break
					}
				}
				if !allB && emptyB && !emptyA {
					return false
				}

			}
			if !allB && resB.IsEmpty() && !resA.IsEmpty() {
				return false
			}
			if resB.IsEmpty() && resA.IsEmpty() {
				break
			}
			baseA = resA
			baseB = resB
		}
	}
	return true
}

// JoinIntersects perform join of two bitmap index and return cross product matrix.
func JoinIntersects(a, b *Index, left bool) *bjoin.Index {
	// for each corresponded base cell
	join := bjoin.New(b.MaxItemIndex() + 1)

	// find elements of a that not intersections base cells
	noIntersectsA := roaring64.New()
	intersectsA := roaring64.New()

	base := roaring64.And(a.baseCellsMask, b.baseCellsMask)
	if base.IsEmpty() {
		return join
	}
	it := base.Iterator() // by intersection of base cells

	for it.HasNext() {
		bn := it.Next()
		bma := a.baseCellMap[bn]
		if bma == nil || bma.IsEmpty() {
			continue
		}
		bmb := b.baseCellMap[bn]
		if bmb == nil || bmb.IsEmpty() {
			continue
		}
		baseA := bma.Clone()
		baseB := bmb.Clone()
		// check elements of a and b by base cell
		for res := 0; res < 15; res++ {
			resA := roaring64.New()
			resB := roaring64.New()
			fullA := roaring64.New()
			fullB := roaring64.New()
			for cn := 7; cn >= 0; cn-- {
				rmA := a.resMaps[res][cn]
				rmB := b.resMaps[res][cn]
				emptyA := rmA == nil || rmA.IsEmpty()
				emptyB := rmB == nil || rmB.IsEmpty()
				if !emptyA {
					rmA = roaring64.And(rmA, baseA)
					emptyA = rmA.IsEmpty()
				}
				if cn == 7 && !emptyA {
					fullA = rmA
				}
				if !emptyB {
					rmB = roaring64.And(rmB, baseB)
					emptyB = rmB.IsEmpty()
				}
				if cn == 7 && !emptyB {
					fullB = rmB
				}
				if emptyA && emptyB || cn == 7 {
					continue
				}
				if !emptyA && emptyB && fullB.IsEmpty() {
					if left {
						noIntersectsA.Or(rmA)
					}
					continue
				}
				if !emptyB && emptyA && fullA.IsEmpty() {
					continue
				}
				if !emptyA {
					resA.Or(rmA)
				}
				if !emptyB {
					resB.Or(rmB)
				}
			}
			if !fullA.IsEmpty() && (!fullB.IsEmpty() || !resB.IsEmpty()) {
				// add all items from resB intersects with fullB
				intersectsA.Or(fullA)
				if !fullB.IsEmpty() {
					join.AddPairs(fullA, fullB)
				}
				if !resB.IsEmpty() {
					join.AddPairs(fullA, resB)
				}
			}
			if !fullB.IsEmpty() && !resA.IsEmpty() {
				// add all items from resA
				intersectsA.Or(resA)
				join.AddPairs(resA, fullB)
			}
			if left && !fullA.IsEmpty() && fullB.IsEmpty() && resB.IsEmpty() {
				noIntersectsA.Or(fullA)
			}
			if resA.IsEmpty() || resB.IsEmpty() {
				break
			}
			baseA.And(resA)
			baseB.And(resB)
		}
	}
	if left {
		// add all A's items with base cells that doesn't contain in B
		it := roaring64.AndNot(a.baseCellsMask, b.baseCellsMask).Iterator()
		for it.HasNext() {
			bn := it.Next()
			bma := a.baseCellMap[bn]
			if bma == nil || bma.IsEmpty() {
				continue
			}
			noIntersectsA.Or(bma)
		}
		noIntersectsA.AndNot(intersectsA)
		join.AddPairs(noIntersectsA, nil)
	}
	return join
}
