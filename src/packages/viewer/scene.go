package viewer

import (
	"fmt"
	"math"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"

	octreeLib "gemilang/src/packages/octree"
	parser "gemilang/src/packages/parser"
)

type voxelResult struct {
	leaves     []*octreeLib.Octree
	bbMin      parser.Vec3
	bbMax      parser.Vec3
	elapsed    time.Duration
	voxelCount int
	pruned     map[int]int
	maxDepth   int
}

// Launch opens the OctoVox GUI: left sidebar for input/stats, right area for 3D view.
func Launch() {
	a := app.App()
	glfwWin := a.IWindow.(*window.GlfwWindow)
	glfwWin.SetTitle("OctoVox — Voxel Viewer")
	glfwWin.SetSize(1100, 700)

	a.Gls().ClearColor(0.04, 0.05, 0.09, 1.0)

	ww, wh := a.IWindow.GetSize()

	scene := core.NewNode()
	gui.Manager().Set(scene)

	// Camera: aspect accounts for the visible 3D area (right of sidebar)
	viewW := float32(ww) - sidebarW
	cam := camera.New(viewW / float32(wh))
	cam.SetProjection(camera.Perspective)
	cam.SetPosition(0, 5, 15)
	scene.Add(cam)

	orbitCtrl := camera.NewOrbitControl(cam)
	orbitCtrl.SetEnabled(camera.OrbitNone)

	// Lighting
	ambLight := light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.55)
	scene.Add(ambLight)
	dirLight := light.NewDirectional(&math32.Color{R: 1, G: 0.95, B: 0.90}, 1.0)
	dirLight.SetPosition(2, 3, 2)
	scene.Add(dirLight)

	// Handle window resize
	a.IWindow.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) {
		newW, newH := a.IWindow.GetSize()
		newViewW := float32(newW) - sidebarW
		if newViewW < 1 {
			newViewW = 1
		}
		cam.SetAspect(newViewW / float32(newH))
		a.Gls().Viewport(0, 0, int32(newW), int32(newH))
	})

	resultChan := make(chan voxelResult, 1)
	var voxelRoot *core.Node
	processing := false

	sidebar := CreateSidebar(scene, wh, func(filename string, depth int) {
		if processing {
			return
		}
		processing = true

		go func() {
			path := "obj/" + filename
			verts, faces, err := parser.ParseOBJ(path)
			if err != nil {
				fmt.Println("Error parsing OBJ:", err)
				processing = false
				return
			}
			bbMin, bbMax := octreeLib.BoundingBox(verts)
			root := &octreeLib.Octree{Min: bbMin, Max: bbMax}
			pruned := map[int]int{}
			start := time.Now()
			octreeLib.Build(root, verts, faces, 0, depth, pruned)
			elapsed := time.Since(start)

			var leaves []*octreeLib.Octree
			octreeLib.CollectLeaves(root, &leaves)

			resultChan <- voxelResult{
				leaves:     leaves,
				bbMin:      bbMin,
				bbMax:      bbMax,
				elapsed:    elapsed,
				voxelCount: len(leaves),
				pruned:     pruned,
				maxDepth:   depth,
			}
		}()
	})

	// Hint label shown in the 3D area before first voxelization
	hintLabel := gui.NewLabel("Select a file and press VOXELIZE")
	hintLabel.SetColor(&math32.Color{R: 0.25, G: 0.35, B: 0.55})
	hintLabel.SetFontSize(14)
	hintLabel.SetPosition(sidebarW+float32(ww-int(sidebarW))/2-145, float32(wh)/2)
	scene.Add(hintLabel)

	hintShown := true

	// Render loop
	a.Run(func(rend *renderer.Renderer, deltaTime time.Duration) {
		select {
		case result := <-resultChan:
			processing = false

			// Remove previous voxels and hint
			if voxelRoot != nil {
				scene.Remove(voxelRoot)
			}
			if hintShown {
				scene.Remove(hintLabel)
				hintShown = false
			}

			// Add new voxels
			voxelRoot = core.NewNode()
			scene.Add(voxelRoot)
			positionCamera(cam, result.bbMin, result.bbMax)
			orbitCtrl.SetEnabled(camera.OrbitAll)
			AddVoxels(voxelRoot, result.leaves, result.bbMin, result.bbMax)
			sidebar.UpdateStats(result)
		default:
		}

		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		rend.Render(scene, cam)
	})
}

// positionCamera moves the camera to fully see the bounding box.
func positionCamera(cam *camera.Camera, bbMin, bbMax parser.Vec3) {
	cx := float32((bbMin.X + bbMax.X) / 2)
	cy := float32((bbMin.Y + bbMax.Y) / 2)
	cz := float32((bbMin.Z + bbMax.Z) / 2)
	dx := bbMax.X - bbMin.X
	dy := bbMax.Y - bbMin.Y
	dz := bbMax.Z - bbMin.Z
	diag := float32(math.Sqrt(dx*dx+dy*dy+dz*dz)) * 1.5
	cam.SetPosition(cx+float32(sidebarW)*0.01, cy+diag*0.25, cz+diag)
	cam.LookAt(&math32.Vector3{X: cx, Y: cy, Z: cz}, &math32.Vector3{X: 0, Y: 1, Z: 0})
}
