package octree

import (
	parser "gemilang/src/packages/parser"
	intersect "gemilang/src/packages/intersect"
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

func MakeOctant(min, max, mid parser.Vec3, i int) *Octree {
	var oMin, oMax parser.Vec3

	// bit 0 → sumbu X
	if i&1 != 0 {
		oMin.X = mid.X
		oMax.X = max.X
	} else {
		oMin.X = min.X
		oMax.X = mid.X
	}

	// bit 1 → sumbu Y
	if i&2 != 0 {
		oMin.Y = mid.Y
		oMax.Y = max.Y
	} else {
		oMin.Y = min.Y
		oMax.Y = mid.Y
	}

	// bit 2 → sumbu Z
	if i&4 != 0 {
		oMin.Z = mid.Z
		oMax.Z = max.Z
	} else {
		oMin.Z = min.Z
		oMax.Z = mid.Z
	}

	return &Octree{Min: oMin, Max: oMax}
}

func Build(node *Octree, verts []parser.Vec3, faces []parser.Face, depth, maxDepth int) {
	if depth == maxDepth {
		node.IsLeaf = true
		return
	}
	mid := MidPoint(node.Min, node.Max)
	for i := 0; i < 8; i++ {
		child := MakeOctant(node.Min, node.Max, mid, i)
		if intersect.TriBoxIntersect(child.Min, child.Max, verts, faces) {
			node.Children[i] = child
			Build(child, verts, faces, depth+1, maxDepth)
		}
	}
}
