package texture

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/tools"
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
	id        int32
	width     uint32
	height    uint32
	aspect    float32
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

func newTexture(width, height uint32) (tex *Texture) {
	id := idSys.GetID()
	colors := blankTexture(width, height)
	glID := renderer.CreateTexture(int32(width), int32(height), flatData(colors))
	tex = new(Texture)
	tex.id = id
	tex.width = width
	tex.height = height
	tex.aspect = float32(width) / float32(height)
	tex.colors = colors
	tex.glID = glID
	return
}

func (s *Texture) change(changes []pixelChange) {
	for _, change := range changes {
		index := change.pos[1]*int32(s.width) + change.pos[0]
		s.colors[index] = change.after
	}
}

func (s *Texture) applyChanges(changes []pixelChange) {
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
	lastMousePos  [2]float32
	mousePressPos [2]float32
	mouseCanDrag  bool
	painting      bool
	firstButton   bool
	pos           [2]float32
	viewrSize     [2]float32
	aspect        float32
	textureViewer *renderer.FrameBuffer
	texture       *Texture
	preview       *Texture
	pixelPos      [2]int32
	lastPixelPos  [2]int32
}

func (s *TextureContext) changePixels(pixels [][2]int32, color [4]float32) {
	var changes []pixelChange
	for _, pixelPos := range pixels {
		if pixelPos[0] < 0 || pixelPos[0] >= int32(s.texture.width) || pixelPos[1] < 0 || pixelPos[1] >= int32(s.texture.height) {
			continue
		}
		var change pixelChange
		index := pixelPos[1]*int32(s.texture.width) + pixelPos[0]
		change.pos = pixelPos
		change.before = s.texture.colors[index]
		change.after = color
		changes = append(changes, change)
	}
	if len(changes) > 0 {
		s.texture.applyChanges(changes)
	}
}

func (s *TextureContext) applyPreview(pixels [][2]int32, color [4]float32) {
	var changes []pixelChange
	for _, pixelPos := range pixels {
		if pixelPos[0] < 0 || pixelPos[0] >= int32(s.preview.width) || pixelPos[1] < 0 || pixelPos[1] >= int32(s.preview.height) {
			continue
		}
		var change pixelChange
		change.pos = pixelPos
		change.before = [4]float32{0, 0, 0, 0}
		change.after = color
		changes = append(changes, change)
	}
	if len(changes) > 0 {
		s.preview.applyChanges(changes)
	}
}

func (s *TextureContext) resetPreview() {
	s.preview.colors = blankTexture(s.preview.width, s.preview.height)
	s.preview.update()
}

func (s *TextureContext) reset() {
	s.zoom = 1
	s.pos = [2]float32{0, 0}
}

func (s *TextureContext) getPixelPosMouse() (pixelPos [2]int32) {
	var pos1 [2]float32
	pos1[0] = s.viewrSize[0] - s.mousePos[0] - s.pos[0]*s.aspect
	pos1[1] = s.viewrSize[1] - s.mousePos[1] - s.pos[1]
	pos1[0] = 2 * ((pos1[0] / s.viewrSize[0]) - 0.5)
	pos1[1] = 2 * ((pos1[1] / s.viewrSize[1]) - 0.5)
	pos1[0] = pos1[0] / (s.zoom * s.aspect * s.texture.aspect)
	pos1[1] = pos1[1] / s.zoom
	pixelPos[0] = int32(math.Floor(float64(s.texture.width) * float64(pos1[0]/2+0.5)))
	pixelPos[1] = int32(math.Floor(float64(s.texture.height) * (1 - float64(pos1[1]/2+0.5))))
	return
}

func (s *TextureContext) scroll(yoffset float32) {
	mouse := [2]float32{s.mousePos[0] - 0.5*s.viewrSize[0], s.mousePos[1] - 0.5*s.viewrSize[1]}
	position := [2]float32{(s.pos[0] + mouse[0]) / s.zoom, (s.pos[1] + mouse[1]) / s.zoom}

	if yoffset < 0 {
		s.zoom *= 0.9
	} else if yoffset > 0 {
		s.zoom *= 1.1
	}

	s.pos = [2]float32{s.zoom*position[0] - mouse[0], s.zoom*position[1] - mouse[1]}
}

func (s *TextureContext) move(pos imgui.Vec2) {
	s.mousePos = [2]float32{s.viewrSize[0] - pos.X, pos.Y}

	pixel := s.getPixelPosMouse()
	if (pixel[0] != s.lastPixelPos[0] || pixel[1] != s.lastPixelPos[1]) && s.painting {
		tools.Move(s.lastPixelPos, pixel)
		s.resetPreview()
		pixels := tools.Visualize()
		if s.firstButton {
			s.applyPreview(pixels, color1)
		} else {
			s.applyPreview(pixels, color2)
		}
		s.lastPixelPos = pixel
	}

	if s.mouseCanDrag {
		s.pos[0] = (s.lastMousePos[0]-s.mousePos[0])/s.aspect + s.mousePressPos[0]
		s.pos[1] = s.lastMousePos[1] - s.mousePos[1] + s.mousePressPos[1]
	}
}

var secondButton bool

func (s *TextureContext) buttonPress(buttons [5]bool) {
	if buttons[2] {
		s.lastMousePos = s.mousePos
		s.mousePressPos = s.pos
		s.mouseCanDrag = true
	}
	if buttons[0] || buttons[1] {
		s.painting = true
		s.firstButton = buttons[0]
		tools.ButtonPress(s.getPixelPosMouse(), s.texture.colors, s.texture.width, s.texture.height)
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
	if buttons[0] || buttons[1] {
		s.painting = false
		tools.ButtonRelease(s.getPixelPosMouse())
		pixels := tools.Change()
		if s.firstButton {
			s.changePixels(pixels, color1)
		} else {
			s.changePixels(pixels, color2)
		}
		s.resetPreview()
	}
}

func (s *TextureContext) Show() {
	imgui.Begin(s.windowName)

	wSize := imgui.ContentRegionAvail()
	width := int32(wSize.X)
	height := int32(wSize.Y)

	if s.viewrSize[0] != wSize.X || s.viewrSize[1] != wSize.Y {
		s.textureViewer.Resize(width, height)
		s.viewrSize[0] = wSize.X
		s.viewrSize[1] = wSize.Y
		s.aspect = wSize.Y / wSize.X
	}

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
		s.move(mouse_pos_rel)
		s.buttonPress(io.MouseDown())
		s.buttonRelease(io.MouseReleased())
	}
	imgui.End()
	x := float32(s.texture.width)
	y := float32(s.texture.height)
	renderer.RenderTexture(*s.textureViewer, s.texture.glID, s.preview.glID, s.zoom, s.pos[0], s.pos[1], s.texture.aspect, x, y)
}

var AllTextures []*Texture
var AllCtx []*TextureContext

func createCtx(tex *Texture) (ctx TextureContext) {
	ctx.windowName = "Texture #" + strconv.FormatUint(uint64(tex.id), 10)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.CreateFramebuffer(500, 500)
	ctx.viewrSize = [2]float32{500, 500}
	ctx.texture = tex
	ctx.preview = newTexture(tex.width, tex.height)
	return
}

func AddTexture(width, height uint32) {
	tex := newTexture(width, height)
	ctx := createCtx(tex)
	AllTextures = append(AllTextures, tex)
	AllCtx = append(AllCtx, &ctx)
}

func OpenTexture(tex *Texture) {
	createCtx(tex)
}

var idSys *util.IdSystem
var color1 [4]float32
var color2 [4]float32

func Init() {
	idSys = util.NewIdSystem()
}

func SetColors(c1, c2 [4]float32) {
	color1 = c1
	color2 = c2
}
