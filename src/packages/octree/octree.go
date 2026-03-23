package octree

import (
	intersect "gemilang/src/packages/intersect"
	parser "gemilang/src/packages/parser"
	"sync"
)

type Octree struct {
	Min      parser.Vec3
	Max      parser.Vec3
	Children [8]*Octree
	IsLeaf   bool
}

func BoundingBox(verts []parser.Vec3) (parser.Vec3, parser.Vec3) {
	min := verts[0]
	max := verts[0]

	for _, v := range verts {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.Z < min.Z {
			min.Z = v.Z
		}

		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
		if v.Z > max.Z {
			max.Z = v.Z
		}

	}
	return min, max
}

func MidPoint(a, b parser.Vec3) parser.Vec3 {
	var mid parser.Vec3
	mid.X = (a.X + b.X) / 2
	mid.Y = (a.Y + b.Y) / 2
	mid.Z = (a.Z + b.Z) / 2

	return mid
}

// pickHalf memilih separuh bawah atau atas dari satu sumbu.
// Kalau useUpperHalf = true -> ambil [mid, hi], sebaliknya -> [lo, mid].
func pickHalf(useUpperHalf bool, lo, mid, hi float64) (float64, float64) {
	if useUpperHalf {
		return mid, hi
	}
	return lo, mid
}

// MakeOctant membuat salah satu dari 8 octant anak.
// Index i (0–7) dibaca sebagai 3 bit: bit0=X, bit1=Y, bit2=Z.
// Bit 0 (bernilai 1) -> separuh kanan X, bit 1 (bernilai 2) -> separuh atas Y,
// bit 2 (bernilai 4) -> separuh belakang Z.
func MakeOctant(min, max, mid parser.Vec3, i int) *Octree {
	xRight := i&1 != 0
	yUp    := i&2 != 0
	zBack  := i&4 != 0

	xMin, xMax := pickHalf(xRight, min.X, mid.X, max.X)
	yMin, yMax := pickHalf(yUp,    min.Y, mid.Y, max.Y)
	zMin, zMax := pickHalf(zBack,  min.Z, mid.Z, max.Z)

	return &Octree{
		Min: parser.Vec3{X: xMin, Y: yMin, Z: zMin},
		Max: parser.Vec3{X: xMax, Y: yMax, Z: zMax},
	}
}

func Build(node *Octree, verts []parser.Vec3, faces []parser.Face, depth, maxDepth int, prunedCounts map[int]int) {
	if depth == maxDepth {
		node.IsLeaf = true
		return
	}
	mid := MidPoint(node.Min, node.Max)

	// implmentasi concurrency
	var wg sync.WaitGroup
    var mu sync.Mutex
	for i := 0; i < 8; i++ {
		child := MakeOctant(node.Min, node.Max, mid, i)
		if intersect.TriBoxIntersect(child.Min, child.Max, verts, faces) {
			node.Children[i] = child
			// Build(child, verts, faces, depth+1, maxDepth, prunedCounts)
			wg.Add(1)
			go func (idx int, c *Octree)  {
				defer wg.Done()
				localPruned := map[int]int{}

				Build(c, verts, faces, depth+1, maxDepth, localPruned)

				mu.Lock()
				for x, y := range localPruned{
					prunedCounts[x] += y;
				}
				mu.Unlock()
			}(i, child)
		} else {
			mu.Lock();
			prunedCounts[depth+1]++
			mu.Unlock();
		}
	}
	wg.Wait()
}
