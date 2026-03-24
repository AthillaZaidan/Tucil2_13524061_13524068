package viewer

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	octree "gemilang/src/packages/octree"
	parser "gemilang/src/packages/parser"
)

func AddVoxels(scene *core.Node, leaves []*octree.Octree, bbMin, bbMax parser.Vec3) {
	rangeX := bbMax.X - bbMin.X
	rangeY := bbMax.Y - bbMin.Y
	rangeZ := bbMax.Z - bbMin.Z

	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}
	if rangeZ == 0 {
		rangeZ = 1
	}

	for _, leaf := range leaves {
		cx := float32((leaf.Min.X+leaf.Max.X)/2)
		cy := float32((leaf.Min.Y+leaf.Max.Y)/2)
		cz := float32((leaf.Min.Z+leaf.Max.Z)/2)

		sx := float32(leaf.Max.X - leaf.Min.X)
		sy := float32(leaf.Max.Y - leaf.Min.Y)
		sz := float32(leaf.Max.Z - leaf.Min.Z)

		r := float32((leaf.Min.X - bbMin.X) / rangeX)
		g := float32((leaf.Min.Y - bbMin.Y) / rangeY)
		b := float32((leaf.Min.Z - bbMin.Z) / rangeZ)

		geom := geometry.NewBox(sx, sy, sz)
		mat := material.NewStandard(&math32.Color{R: r, G: g, B: b})
		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(cx, cy, cz)
		scene.Add(mesh)
	}
}
