package textureDraw

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/util"
	"math"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
)

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type pixelsChange []pixelChange

type Texture struct {
	id        uint32
	width     uint32
	height    uint32
	colors    [][4]float32
	glID      uint32
	changes   []pixelsChange
	undoLevel int32
}

func blankTexture(width, height uint32) (data [][4]float32) {
	length := int(width * height)
	for i := 0; i < int(length); i++ {
		data = append(data, [4]float32{0, 0, 0, 0})
	}
	return
}

func flatData(colors [][4]float32) (data []float32) {
	for _, color := range colors {
		data = append(data, color[0], color[1], color[2], color[3])
	}
	return
}

func newTexture(width, height uint32) (tex Texture) {
	id := idSys.GetID()
	colors := blankTexture(width, height)
	glID := renderer.CreateTexture(int32(width), int32(height), flatData(colors))
	tex.id = id
	tex.width = width
	tex.height = height
	tex.colors = blankTexture(width, height)
	tex.glID = glID
	return
}

func isSameChanges(last, after []pixelChange) bool {
	if len(last) != len(after) {
		return false
	}
loop:
	for _, change2 := range last {
		for _, change1 := range after {
			if change1.pos == change2.pos && change1.after == change2.after {
				continue loop
			}
		}
		return false
	}
	return true
}

func (s *Texture) change(changes []pixelChange) {
	for _, change := range changes {
		index := change.pos[1]*int32(s.width) + change.pos[0]
		s.colors[index] = change.after
	}
}

func (s *Texture) applyChanges(changes []pixelChange) {
	if len(s.changes) > 0 {
		lastChange := s.changes[len(s.changes)-1]
		if isSameChanges(lastChange, changes) {
			return
		}
	}
	s.change(changes)
	s.update()
	s.changes = s.changes[:len(s.changes)-int(s.undoLevel)]
	s.changes = append(s.changes, changes)
	s.undoLevel = 0
}

func (s *Texture) undo() {
	changesIdx := len(s.changes) - 1 - int(s.undoLevel)
	if changesIdx >= 0 {
		lastChanges := s.changes[changesIdx]
		var undoChanges []pixelChange
		s.undoLevel++
		for _, change := range lastChanges {
			var undoChange pixelChange
			undoChange.pos = change.pos
			undoChange.after = change.before
			undoChange.before = change.after
			undoChanges = append(undoChanges, undoChange)
		}
		s.change(undoChanges)
		s.update()
	}
}

func (s *Texture) redo() {
	if s.undoLevel > 0 {
		changesIdx := len(s.changes) - int(s.undoLevel)
		if changesIdx >= 0 {
			lastChanges := s.changes[changesIdx]
			s.undoLevel--
			s.change(lastChanges)
			s.update()
		}
	}
}

func (s *Texture) update() {
	data := flatData(s.colors)
	renderer.WriteTexture(s.glID, int32(s.width), int32(s.height), data)
}

type TextureContext struct {
	windowName    string
	zoom          float32
	mousePos      [2]float32
	mouseTarget   [2]float32
	mousePressPos [2]float32
	mouseCanDrag  bool
	pos           [2]float32
	size          [2]float32
	textureViewer renderer.FrameBuffer
	texture       *Texture
}

func (s *TextureContext) change_pixel() {
	var pos [2]float32
	pos[0] = s.size[0] - s.mousePos[0] - s.pos[0] - s.size[0]/2 + s.size[0]*s.zoom/2
	pos[1] = s.mousePos[1] + s.pos[1] - s.size[1]/2 + s.size[0]*s.zoom/2
	pos[0] = pos[0] / (s.size[0] * s.zoom) * float32(s.texture.width)
	pos[1] = pos[1] / (s.size[0] * s.zoom) * float32(s.texture.height)
	pixel_pos := [2]int32{int32(math.Floor(float64(pos[0]))), int32(math.Floor(float64(pos[1])))}
	if pixel_pos[0] < 0 || pixel_pos[0] >= int32(s.texture.width) || pixel_pos[1] < 0 || pixel_pos[1] >= int32(s.texture.height) {
		return
	}
	//col := &pixels[pixel_pos.y][pixel_pos.x]
	var change pixelChange
	index := pixel_pos[1]*int32(s.texture.width) + pixel_pos[0]
	change.pos = pixel_pos
	change.before = s.texture.colors[index]
	change.after = [4]float32{1, 1, 1, 1}
	s.texture.applyChanges([]pixelChange{change})
}

func (s *TextureContext) reset() {
	s.zoom = 1
	s.pos = [2]float32{0, 0}
}

func (s *TextureContext) scroll(yoffset float32) {
	mouse := [2]float32{s.mousePos[0] - 0.5*s.size[0], s.mousePos[1] - 0.5*s.size[1]}
	position := [2]float32{(s.pos[0] + mouse[0]) / s.zoom, (s.pos[1] + mouse[1]) / s.zoom}

	if yoffset < 0 {
		s.zoom *= 0.9
	} else if yoffset > 0 {
		s.zoom *= 1.1
	}

	s.pos = [2]float32{s.zoom*position[0] - mouse[0], s.zoom*position[1] - mouse[1]}
}

func (s *TextureContext) move(pos imgui.Vec2, buttons [5]bool) {
	s.mousePos = [2]float32{s.size[0] - pos.X, pos.Y}

	if buttons[0] {
		s.change_pixel()
		//fmt.Println("Pintado!")
		// change_pixel(false)
	} else if buttons[1] {
		// change_pixel(true)
	}

	if s.mouseCanDrag {
		s.pos[0] = s.mouseTarget[0] - s.mousePos[0] + s.mousePressPos[0]
		s.pos[1] = s.mouseTarget[1] - s.mousePos[1] + s.mousePressPos[1]
	}
}

func (s *TextureContext) buttonPress(buttons [5]bool) {
	if buttons[2] {
		s.mouseTarget = s.mousePos
		s.mousePressPos = s.pos
		s.mouseCanDrag = true
	}
}

func (s *TextureContext) buttonRelease(buttons [5]bool) {
	if buttons[2] {
		s.mouseCanDrag = false
	} else if buttons[3] {
		s.texture.undo()
	} else if buttons[4] {
		s.texture.redo()
	}
}

func (s *TextureContext) Show() {
	imgui.Begin(s.windowName)
	if imgui.IsWindowFocused() {
		io := imgui.CurrentContext().IO()
		if io.KeyCtrl() && imgui.IsKeyPressedBoolV(imgui.KeyZ, true) {
			s.texture.undo()
		}
		if io.KeyCtrl() && imgui.IsKeyPressedBoolV(imgui.KeyY, true) {
			s.texture.redo()
		}
	}
	imgui.ImageV(
		*imgui.NewTextureRefTextureID(imgui.TextureID(s.textureViewer.Image())),
		s.textureViewer.Size(),
		imgui.NewVec2(0, 1),
		imgui.NewVec2(1, 0),
	)
	if imgui.IsItemHovered() {
		io := imgui.CurrentContext().IO()
		if io.MouseWheel() != 0 {
			s.scroll(io.MouseWheel())
		}
		mouse_pos_abs := io.MousePos()
		screen_pos_abs := imgui.ItemRectMin()
		var mouse_pos_rel imgui.Vec2
		mouse_pos_rel.X = mouse_pos_abs.X - screen_pos_abs.X
		mouse_pos_rel.Y = mouse_pos_abs.Y - screen_pos_abs.Y
		s.move(mouse_pos_rel, io.MouseDown())
		s.buttonPress(io.MouseDown())
		s.buttonRelease(io.MouseReleased())
	}
	imgui.End()
	renderer.RenderTexture(s.textureViewer, s.texture.glID, s.zoom, s.pos[0], s.pos[1])
}

var allTextures []*Texture
var allCtx []*TextureContext

func createCtx(tex *Texture) (ctx TextureContext) {
	ctx.windowName = "Texture #" + strconv.FormatUint(uint64(tex.id), 10)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.NewFrameBuffer(500, 500)
	ctx.size = [2]float32{500, 500}
	ctx.texture = tex
	return
}

func AddTexture(width, height uint32) {
	tex := newTexture(width, height)
	ctx := createCtx(&tex)
	allTextures = append(allTextures, &tex)
	allCtx = append(allCtx, &ctx)
}

func OpenTexture(tex *Texture) {
	createCtx(tex)
}

var idSys util.IdSystem

func Init() {
	idSys = util.NewIdSystem()
}

func Loop() {
	for _, c := range allCtx {
		c.Show()
	}
}
