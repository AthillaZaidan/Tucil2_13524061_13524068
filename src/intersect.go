package main

import "math"

func dot(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func cross(a, b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func sub(a, b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func projectTriangle(axis, v0, v1, v2 Vec3) (float64, float64) {
	p0 := dot(axis, v0)
	p1 := dot(axis, v1)
	p2 := dot(axis, v2)
	return math.Min(p0, math.Min(p1, p2)),
		math.Max(p0, math.Max(p1, p2))
}

func projectBox(axis, center, halfSize Vec3) (float64, float64) {
	r := math.Abs(dot(axis, Vec3{halfSize.X, 0, 0})) +
		math.Abs(dot(axis, Vec3{0, halfSize.Y, 0})) +
		math.Abs(dot(axis, Vec3{0, 0, halfSize.Z}))
	c := dot(axis, center)
	return c - r, c + r
}

func overlaps(minA, maxA, minB, maxB float64) bool {
	return maxA >= minB && maxB >= minA
}

func triBoxOverlap(boxMin, boxMax Vec3, v0, v1, v2 Vec3) bool {
	center := Vec3{
		X: (boxMin.X + boxMax.X) / 2,
		Y: (boxMin.Y + boxMax.Y) / 2,
		Z: (boxMin.Z + boxMax.Z) / 2,
	}
	half := Vec3{
		X: (boxMax.X - boxMin.X) / 2,
		Y: (boxMax.Y - boxMin.Y) / 2,
		Z: (boxMax.Z - boxMin.Z) / 2,
	}

	v0 = sub(v0, center)
	v1 = sub(v1, center)
	v2 = sub(v2, center)

	e0 := sub(v1, v0)
	e1 := sub(v2, v1)
	e2 := sub(v0, v2)

	axes := []Vec3{
		cross(e0, Vec3{1, 0, 0}), cross(e0, Vec3{0, 1, 0}), cross(e0, Vec3{0, 0, 1}),
		cross(e1, Vec3{1, 0, 0}), cross(e1, Vec3{0, 1, 0}), cross(e1, Vec3{0, 0, 1}),
		cross(e2, Vec3{1, 0, 0}), cross(e2, Vec3{0, 1, 0}), cross(e2, Vec3{0, 0, 1}),
	}
	for _, axis := range axes {
		if dot(axis, axis) < 1e-10 {
			continue
		}
		tMin, tMax := projectTriangle(axis, v0, v1, v2)
		bMin, bMax := projectBox(axis, Vec3{}, half)
		if !overlaps(tMin, tMax, bMin, bMax) {
			return false
		}
	}

	for _, axis := range []Vec3{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}} {
		tMin, tMax := projectTriangle(axis, v0, v1, v2)
		bMin, bMax := projectBox(axis, Vec3{}, half)
		if !overlaps(tMin, tMax, bMin, bMax) {
			return false
		}
	}

	normal := cross(e0, e1)
	tMin, tMax := projectTriangle(normal, v0, v1, v2)
	bMin, bMax := projectBox(normal, Vec3{}, half)
	if !overlaps(tMin, tMax, bMin, bMax) {
		return false
	}

	return true
}

func TriBoxIntersect(boxMin, boxMax Vec3, verts []Vec3, faces []Face) bool {
	for _, f := range faces {
		v0 := verts[f.V1]
		v1 := verts[f.V2]
		v2 := verts[f.V3]
		if triBoxOverlap(boxMin, boxMax, v0, v1, v2) {
			return true
		}
	}
	return false
}
