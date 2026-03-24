package viewer

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

const (
	sidebarW = float32(280)
	padX     = float32(20)
	fieldW   = sidebarW - padX*2
)

// colors
var (
	colBg       = math32.Color4{R: 0.030, G: 0.040, B: 0.080, A: 1}
	colHeaderBg = math32.Color4{R: 0.040, G: 0.055, B: 0.110, A: 1}
	colBorder   = math32.Color4{R: 0.08, G: 0.18, B: 0.40, A: 1}
	colCyan     = math32.Color{R: 0.00, G: 0.83, B: 1.00}
	colCyanDim  = math32.Color{R: 0.20, G: 0.55, B: 0.75}
	colText     = math32.Color{R: 0.85, G: 0.88, B: 0.95}
	colMuted    = math32.Color{R: 0.40, G: 0.47, B: 0.60}
	colGreen    = math32.Color{R: 0.20, G: 0.90, B: 0.50}
	colOrange   = math32.Color{R: 1.00, G: 0.60, B: 0.10}
)

type Sidebar struct {
	statusLabel     *gui.Label
	voxelLabel      *gui.Label
	vertexLabel     *gui.Label
	faceLabel       *gui.Label
	timeLabel       *gui.Label
	startBtn        *gui.Button
	prunedContainer *gui.Panel
}

func (s *Sidebar) UpdateStats(result voxelResult) {
	s.statusLabel.SetText("● Done")
	s.statusLabel.SetColor(&colGreen)
	s.voxelLabel.SetText(fmt.Sprintf("%d", result.voxelCount))
	s.vertexLabel.SetText(fmt.Sprintf("%d", result.voxelCount*8))
	s.faceLabel.SetText(fmt.Sprintf("%d", result.voxelCount*12))
	s.timeLabel.SetText(result.elapsed.Round(time.Millisecond).String())
	s.startBtn.SetEnabled(true)

	// Repopulate pruned nodes per depth
	s.prunedContainer.DisposeChildren(true)
	rowY := float32(0)
	for d := 1; d <= result.maxDepth; d++ {
		count := result.pruned[d]
		keyLbl := gui.NewLabel(fmt.Sprintf("Depth %d :", d))
		keyLbl.SetColor(&colMuted)
		keyLbl.SetFontSize(12)
		keyLbl.SetPosition(0, rowY)
		s.prunedContainer.Add(keyLbl)

		valLbl := gui.NewLabel(fmt.Sprintf("%d", count))
		valLbl.SetColor(&colText)
		valLbl.SetFontSize(12)
		valLbl.SetPosition(90, rowY)
		s.prunedContainer.Add(valLbl)
		rowY += 22
	}
	s.prunedContainer.SetHeight(rowY)
}

func (s *Sidebar) SetProcessing() {
	s.statusLabel.SetText("● Processing...")
	s.statusLabel.SetColor(&colOrange)
	s.startBtn.SetEnabled(false)
}

func (s *Sidebar) SetError(msg string) {
	s.statusLabel.SetText("● " + msg)
	s.statusLabel.SetColor(&colOrange)
	s.startBtn.SetEnabled(true)
}

func CreateSidebar(scene *core.Node, wh int, onSubmit func(filename string, depth int, outputName string)) *Sidebar {
	sb := &Sidebar{}

	panel := gui.NewPanel(sidebarW, float32(wh))
	panel.SetColor4(&colBg)
	panel.SetBordersFrom(&gui.RectBounds{Right: 1})
	panel.SetBordersColor4(&colBorder)
	panel.SetPosition(0, 0)

	y := float32(0)

	header := gui.NewPanel(sidebarW, 70)
	header.SetColor4(&colHeaderBg)
	header.SetBordersFrom(&gui.RectBounds{Bottom: 1})
	header.SetBordersColor4(&colBorder)
	header.SetPosition(0, y)

	logoLabel := gui.NewLabel("◈ OCTOVOX")
	logoLabel.SetColor(&colCyan)
	logoLabel.SetFontSize(20)
	logoLabel.SetPosition(padX, 12)
	header.Add(logoLabel)

	subtitleLabel := gui.NewLabel("3D Voxelization Viewer")
	subtitleLabel.SetColor(&colMuted)
	subtitleLabel.SetFontSize(11)
	subtitleLabel.SetPosition(padX, 42)
	header.Add(subtitleLabel)

	panel.Add(header)
	y += 82

	sectionLabel(panel, "INPUT FILE", padX, y)
	y += 22

	objFiles := scanObjFiles("obj")
	var dd *gui.DropDown
	if len(objFiles) > 0 {
		dd = gui.NewDropDown(fieldW, gui.NewImageLabel(objFiles[0]))
		for _, f := range objFiles {
			dd.Add(gui.NewImageLabel(f))
		}
	} else {
		dd = gui.NewDropDown(fieldW, gui.NewImageLabel("(no .obj files in obj/)"))
	}
	dd.SetPosition(padX, y)
	panel.Add(dd)
	y += 42

	sectionLabel(panel, "MAX DEPTH", padX, y)
	y += 22

	depthEdit := gui.NewEdit(int(fieldW), "4")
	depthEdit.SetPosition(padX, y)
	panel.Add(depthEdit)
	y += 42

	y += 4
	startBtn := gui.NewButton("  VOXELIZE  ")
	startBtn.SetColor4(&math32.Color4{R: 0.00, G: 0.50, B: 0.72, A: 1.0})
	startBtn.Label.SetColor(&math32.Color{R: 0.92, G: 0.97, B: 1.0})
	startBtn.Label.SetFontSize(13)
	startBtn.SetWidth(fieldW)
	startBtn.SetPosition(padX, y)
	sb.startBtn = startBtn

	startBtn.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		if len(objFiles) == 0 {
			return
		}
		var filename string
		if dd.Selected() == nil {
			filename = objFiles[0]
		} else {
			filename = dd.Selected().Text()
		}
		depthStr := strings.TrimSpace(depthEdit.Text())
		if depthStr == "" {
			depthStr = "4"
		}
		depth, err := strconv.Atoi(depthStr)
		if err != nil || depth < 1 {
			depth = 4
		}

		defaultOut := strings.TrimSuffix(filename, filepath.Ext(filename)) + "_voxelized.obj"
		showOutputPrompt(scene, wh, defaultOut, func(outputName string) {
			sb.SetProcessing()
			onSubmit(filename, depth, outputName)
		})
	})
	panel.Add(startBtn)
	y += 46

	divider(panel, padX, y)
	y += 14

	statusKey := gui.NewLabel("STATUS")
	statusKey.SetColor(&colMuted)
	statusKey.SetFontSize(10)
	statusKey.SetPosition(padX, y)
	panel.Add(statusKey)
	y += 16

	sb.statusLabel = gui.NewLabel("● Idle")
	sb.statusLabel.SetColor(&colCyanDim)
	sb.statusLabel.SetFontSize(13)
	sb.statusLabel.SetPosition(padX, y)
	panel.Add(sb.statusLabel)
	y += 30

	divider(panel, padX, y)
	y += 14

	sectionLabel(panel, "STATISTICS", padX, y)
	y += 24

	sb.voxelLabel = statRow(panel, "Voxels", "—", padX, y)
	y += 22
	sb.vertexLabel = statRow(panel, "Vertices", "—", padX, y)
	y += 22
	sb.faceLabel = statRow(panel, "Faces", "—", padX, y)
	y += 22
	sb.timeLabel = statRow(panel, "Time", "—", padX, y)
	y += 30

	divider(panel, padX, y)
	y += 14

	sectionLabel(panel, "PRUNED NODES PER DEPTH", padX, y)
	y += 24

	prunedContainer := gui.NewPanel(fieldW, 0)
	prunedContainer.SetPosition(padX, y)
	sb.prunedContainer = prunedContainer
	panel.Add(prunedContainer)

	scene.Add(panel)
	return sb
}

func showOutputPrompt(scene *core.Node, wh int, defaultName string, onConfirm func(outputName string)) {
	overlayW := float32(390)
	overlayH := float32(170)
	overlay := gui.NewPanel(overlayW, overlayH)
	overlay.SetColor4(&math32.Color4{R: 0.05, G: 0.08, B: 0.13, A: 0.98})
	overlay.SetBordersFrom(&gui.RectBounds{Top: 1, Right: 1, Bottom: 1, Left: 1})
	overlay.SetBordersColor4(&colBorder)
	overlay.SetPosition(sidebarW+20, float32(wh)*0.5-overlayH*0.5)

	title := gui.NewLabel("Output file name")
	title.SetColor(&colCyan)
	title.SetFontSize(14)
	title.SetPosition(16, 14)
	overlay.Add(title)

	nameEdit := gui.NewEdit(int(overlayW-32), defaultName)
	nameEdit.SetPosition(16, 48)
	overlay.Add(nameEdit)

	hint := gui.NewLabel("Saved to obj/<name>.obj")
	hint.SetColor(&colMuted)
	hint.SetFontSize(10)
	hint.SetPosition(16, 80)
	overlay.Add(hint)

	cancelBtn := gui.NewButton("Cancel")
	cancelBtn.SetPosition(16, 108)
	cancelBtn.SetWidth(90)
	cancelBtn.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		scene.Remove(overlay)
	})
	overlay.Add(cancelBtn)

	confirmBtn := gui.NewButton("Save + Voxelize")
	confirmBtn.SetPosition(120, 108)
	confirmBtn.SetWidth(150)
	confirmBtn.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		name := filepath.Base(strings.TrimSpace(nameEdit.Text()))
		if name == "" {
			return
		}
		if filepath.Ext(strings.ToLower(name)) != ".obj" {
			name += ".obj"
		}
		scene.Remove(overlay)
		onConfirm(name)
	})
	overlay.Add(confirmBtn)

	scene.Add(overlay)
}

func sectionLabel(parent *gui.Panel, text string, x, y float32) {
	lbl := gui.NewLabel(text)
	lbl.SetColor(&colCyanDim)
	lbl.SetFontSize(10)
	lbl.SetPosition(x, y)
	parent.Add(lbl)
}

func divider(parent *gui.Panel, x, y float32) {
	d := gui.NewPanel(fieldW, 1)
	d.SetColor4(&colBorder)
	d.SetPosition(x, y)
	parent.Add(d)
}

func statRow(parent *gui.Panel, key, value string, x, y float32) *gui.Label {
	keyLbl := gui.NewLabel(key + " :")
	keyLbl.SetColor(&colMuted)
	keyLbl.SetFontSize(12)
	keyLbl.SetPosition(x, y)
	parent.Add(keyLbl)

	valLbl := gui.NewLabel(value)
	valLbl.SetColor(&colText)
	valLbl.SetFontSize(12)
	valLbl.SetPosition(x+90, y)
	parent.Add(valLbl)

	return valLbl
}

func scanObjFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".obj") {
			files = append(files, e.Name())
		}
	}
	return files
}
