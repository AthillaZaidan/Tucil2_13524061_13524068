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

func mergeBounds(minA, maxA, minB, maxB parser.Vec3) (parser.Vec3, parser.Vec3) {
	min := parser.Vec3{
		X: minA.X,
		Y: minA.Y,
		Z: minA.Z,
	}
	max := parser.Vec3{
		X: maxA.X,
		Y: maxA.Y,
		Z: maxA.Z,
	}

	if minB.X < min.X {
		min.X = minB.X
	}
	if minB.Y < min.Y {
		min.Y = minB.Y
	}
	if minB.Z < min.Z {
		min.Z = minB.Z
	}

	if maxB.X > max.X {
		max.X = maxB.X
	}
	if maxB.Y > max.Y {
		max.Y = maxB.Y
	}
	if maxB.Z > max.Z {
		max.Z = maxB.Z
	}

	return min, max
}

func boundingBoxDivideConquer(verts []parser.Vec3, left, right int) (parser.Vec3, parser.Vec3) {
	if left == right {
		return verts[left], verts[left]
	}

	mid := left + (right-left)/2
	leftMin, leftMax := boundingBoxDivideConquer(verts, left, mid)
	rightMin, rightMax := boundingBoxDivideConquer(verts, mid+1, right)

	return mergeBounds(leftMin, leftMax, rightMin, rightMax)
}

func BoundingBox(verts []parser.Vec3) (parser.Vec3, parser.Vec3) {
	if len(verts) == 0 {
		return parser.Vec3{}, parser.Vec3{}
	}

	return boundingBoxDivideConquer(verts, 0, len(verts)-1)
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
	yUp := i&2 != 0
	zBack := i&4 != 0

	xMin, xMax := pickHalf(xRight, min.X, mid.X, max.X)
	yMin, yMax := pickHalf(yUp, min.Y, mid.Y, max.Y)
	zMin, zMax := pickHalf(zBack, min.Z, mid.Z, max.Z)

	return &Octree{
		Min: parser.Vec3{X: xMin, Y: yMin, Z: zMin},
		Max: parser.Vec3{X: xMax, Y: yMax, Z: zMax},
	}
}

// CollectLeaves recursively collects all leaf nodes from the octree.
func CollectLeaves(node *Octree, leaves *[]*Octree) {
	if node == nil {
		return
	}
	if node.IsLeaf {
		*leaves = append(*leaves, node)
		return
	}
	for _, child := range node.Children {
		CollectLeaves(child, leaves)
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
			go func(idx int, c *Octree) {
				defer wg.Done()
				localPruned := map[int]int{}

				Build(c, verts, faces, depth+1, maxDepth, localPruned)

				mu.Lock()
				for x, y := range localPruned {
					prunedCounts[x] += y
				}
				mu.Unlock()
			}(i, child)
		} else {
			mu.Lock()
			prunedCounts[depth+1]++
			mu.Unlock()
		}
	}
	wg.Wait()
}
