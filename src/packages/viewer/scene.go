package viewer

import (
	"fmt"
	"math"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"

	octree "gemilang/src/packages/octree"
	parser "gemilang/src/packages/parser"
)

// Launch opens a 3D viewer window displaying all voxel leaf nodes.
func Launch(leaves []*octree.Octree, bbMin, bbMax parser.Vec3) {
	a := app.App()
	a.IWindow.(*window.GlfwWindow).SetTitle("OctoVox — 3D Voxel Viewer")

	scene := core.NewNode()

	// Camera
	cam := camera.New(1)
	cam.SetProjection(camera.Perspective)
	scene.Add(cam)

	// Position camera to see the whole model
	cx := float32((bbMin.X + bbMax.X) / 2)
	cy := float32((bbMin.Y + bbMax.Y) / 2)
	cz := float32((bbMin.Z + bbMax.Z) / 2)
	dx := bbMax.X - bbMin.X
	dy := bbMax.Y - bbMin.Y
	dz := bbMax.Z - bbMin.Z
	diagonal := float32(math.Sqrt(dx*dx+dy*dy+dz*dz)) * 1.5
	cam.SetPosition(cx, cy+float32(diagonal*0.3), cz+diagonal)
	cam.LookAt(&math32.Vector3{X: cx, Y: cy, Z: cz}, &math32.Vector3{X: 0, Y: 1, Z: 0})

	// Orbit controls
	camera.NewOrbitControl(cam)

	// Lighting
	ambientLight := light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.5)
	scene.Add(ambientLight)

	dirLight := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.8)
	dirLight.SetPosition(1, 2, 3)
	scene.Add(dirLight)

	// Add voxels
	AddVoxels(scene, leaves, bbMin, bbMax)

	fmt.Printf("Viewer: %d voxel ditampilkan. Tutup window untuk keluar.\n", len(leaves))

	// Render loop
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
