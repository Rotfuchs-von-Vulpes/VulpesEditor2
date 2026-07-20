package canvas

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/canvas/textureEdit"
	"VulpesEditor/app/textureDraw/tools"
	"VulpesEditor/app/util"
	"fmt"
	"math"
	"slices"
	"strconv"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type TextureContext struct {
	id            int32
	windowName    string
	zoom          float32
	mousePos      [2]float32
	lastMousePos  [2]float32
	mousePressPos [2]float32
	mouseCanDrag  bool
	painting      bool
	firstButton   bool
	pos           [2]float32
	viwerSize     [2]float32
	aspect        float32
	textureViewer *renderer.FrameBuffer
	texture       *textureEdit.TextureEdit
	layer         *textureEdit.LayerEdit
	pixelPos      [2]int32
	lastPixelPos  [2]int32
}

func (s *TextureContext) resetPreview() {
	s.texture.ResetPreview()
}

func (s *TextureContext) reset() {
	s.zoom = 1
	s.pos = [2]float32{0, 0}
}

func (s *TextureContext) pixelPosMouse() (pixelPos [2]int32) {
	var pos1 [2]float32
	pos1[0] = s.viwerSize[0] - s.mousePos[0] - s.pos[0]*s.aspect
	pos1[1] = s.viwerSize[1] - s.mousePos[1] - s.pos[1]
	pos1[0] = 2 * ((pos1[0] / s.viwerSize[0]) - 0.5)
	pos1[1] = 2 * ((pos1[1] / s.viwerSize[1]) - 0.5)
	pos1[0] = pos1[0] / (s.zoom * s.aspect * s.texture.Aspect)
	pos1[1] = pos1[1] / s.zoom
	pixelPos[0] = int32(math.Floor(float64(s.texture.Width) * float64(pos1[0]/2+0.5)))
	pixelPos[1] = int32(math.Floor(float64(s.texture.Height) * (1 - float64(pos1[1]/2+0.5))))
	return
}

func (s *TextureContext) scroll(yoffset float32) {
	mouse := [2]float32{s.mousePos[0] - 0.5*s.viwerSize[0], s.mousePos[1] - 0.5*s.viwerSize[1]}
	position := [2]float32{(s.pos[0] + mouse[0]) / s.zoom, (s.pos[1] + mouse[1]) / s.zoom}

	if yoffset < 0 {
		s.zoom *= 0.9
	} else if yoffset > 0 {
		s.zoom *= 1.1
	}

	s.pos = [2]float32{s.zoom*position[0] - mouse[0], s.zoom*position[1] - mouse[1]}
}

func (s *TextureContext) move(pos im.Vec2) {
	s.mousePos = [2]float32{s.viwerSize[0] - pos.X, pos.Y}

	pixel := s.pixelPosMouse()
	if (pixel[0] != s.lastPixelPos[0] || pixel[1] != s.lastPixelPos[1]) && s.painting {
		tools.Move(s.lastPixelPos, pixel)
		s.texture.ResetPreview()
		pixels := tools.Visualize()
		s.texture.ChangePreview(pixels)
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
		toFocus = true
		lastEditCtx = s
		lastEditId = s.id
	}
	if buttons[0] || buttons[1] {
		s.painting = true
		s.firstButton = buttons[0]
		tools.Texture.Resize(s.texture.Width, s.texture.Height)
		tools.Texture.Colors = slices.Clone(s.layer.Texture.Colors)
		tools.ButtonPress(s.pixelPosMouse(), buttons[0])
		toFocus = true
		lastEditCtx = s
		lastEditId = s.id
	}
}

func (s *TextureContext) buttonRelease(buttons [5]bool) {
	if buttons[2] {
		s.mouseCanDrag = false
	}
	if buttons[0] || buttons[1] {
		s.painting = false
		tools.ButtonRelease(s.pixelPosMouse())
		pixels := tools.Change()
		s.layer.Change(pixels)
		s.texture.ResetPreview()
	}
}

var textureFileName string
var textureFilePath string
var lastEditId int32
var toFocus bool

func (s *TextureContext) Show() {
	var toPop string
	if toFocus && lastEditId == s.id {
		im.SetNextWindowFocus()
		toFocus = false
	}
	im.BeginV(s.windowName, nil, im.WindowFlagsMenuBar)
	if im.BeginMenuBar() {
		if im.BeginMenu("Texture") {
			if im.MenuItemBool("Save Texture") {
				toPop = "Not Implement"
			}
			im.EndMenu()
		}
		if im.BeginMenu("Export") {
			if im.MenuItemBool("Export as PNG Image") {
				toPop = "Export PNG"
			}
			im.EndMenu()
		}
		im.EndMenuBar()
	}
	if toPop != "" {
		im.OpenPopupStr(toPop)
		toPop = ""
	}
	if im.BeginPopupModal("Export PNG") {
		im.InputTextWithHint("File Name", "texure_"+strconv.FormatInt(int64(s.texture.Id), 10)+".png", &textureFileName, im.InputTextFlagsNone, nil)
		im.InputTextWithHint("File Path", "", &textureFilePath, im.InputTextFlagsNone, nil)
		if im.Button("Save") {
			if ok := s.texture.SaveTextureAsFile(textureFileName, textureFilePath); ok {
				textureFileName = ""
				im.CloseCurrentPopup()
			}
		}
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}
	if im.BeginPopupModal("Not Implement") {
		im.Text("Not implement yet!")
		if im.Button("OK") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	wSize := im.ContentRegionAvail()
	width := int32(wSize.X)
	height := int32(wSize.Y)

	if s.viwerSize[0] != wSize.X || s.viwerSize[1] != wSize.Y {
		s.textureViewer.Resize(width, height)
		s.viwerSize[0] = wSize.X
		s.viwerSize[1] = wSize.Y
		s.aspect = wSize.Y / wSize.X
	}

	im.ImageV(
		*im.NewTextureRefTextureID(im.TextureID(s.textureViewer.Image())),
		s.textureViewer.Size(),
		im.NewVec2(0, 1),
		im.NewVec2(1, 0),
	)
	if im.IsItemHovered() {
		io := im.CurrentContext().IO()
		if io.MouseWheel() != 0 {
			s.scroll(io.MouseWheel())
		}
		mouse_pos_abs := io.MousePos()
		screen_pos_abs := im.ItemRectMin()
		var mouse_pos_rel im.Vec2
		mouse_pos_rel.X = mouse_pos_abs.X - screen_pos_abs.X
		mouse_pos_rel.Y = mouse_pos_abs.Y - screen_pos_abs.Y
		s.move(mouse_pos_rel)
		s.buttonPress(io.MouseClicked())
		s.buttonRelease(io.MouseReleased())
	} else if s.painting {
		s.buttonRelease([5]bool{true, true, true, false, false})
	}
	im.End()
	s.textureViewer.RenderTexture(s.texture.GlID, s.zoom, s.pos, float32(s.texture.Width), float32(s.texture.Height))
}

// var AllTextures []*texture.Texture
var AllCtx []*TextureContext
var lastEditCtx *TextureContext

func createCtx(tex *texture.Texture) (ctx *TextureContext) {
	ctx = new(TextureContext)
	ctx.id = windowIdSys.GetID()
	ctx.windowName = "Texture #" + strconv.FormatUint(uint64(ctx.id), 10)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.CreateFramebuffer(500, 500)
	ctx.viwerSize = [2]float32{500, 500}
	ctx.texture = textureEdit.New(tex)
	ctx.layer = ctx.texture.Layers[0]
	lastEditId = ctx.id

	ctx.texture.OnLayerEdit = func() {
		found := false
		for _, layer := range ctx.texture.Layers {
			if layer.Id == ctx.layer.Id {
				found = true
				break
			}
		}
		if !found {
			idx := len(ctx.texture.Layers) - 1
			ctx.texture.SetLayer(idx)
			ctx.layer = ctx.texture.Layers[idx]
		}
	}

	return
}

func AddTexture(width, height uint32) (ctx *TextureContext) {
	tex := texture.New(width, height)
	ctx = createCtx(tex)
	return
}

func OpenTexture(tex *texture.Texture) {
	createCtx(tex)
}

var selectedLayers = []bool{}

func (s *TextureContext) ShowLayers() {
	im.Begin("Layers")

	if im.Button("Add") {
		s.texture.AppendLayer()
	}
	im.SameLine()
	if im.Button("Remove") {
		selectedLayers = make([]bool, len(s.texture.Layers))
		im.OpenPopupStr("Remove Layers")
	}
	im.SameLine()
	if im.Button("Merge") {
		selectedLayers = make([]bool, len(s.texture.Layers))
		im.OpenPopupStr("Merge Layers")
	}

	for idx := range s.texture.Layers {
		i := len(s.texture.Layers) - idx - 1
		layer := s.texture.Layers[i]
		cannotDown := i == 0
		cannotUp := i == len(s.texture.Layers)-1
		selected := layer.Id == s.layer.Id
		str := fmt.Sprintf("Layer #%d", i)
		if selected {
			// #86BDFFFF
			im.PushStyleColorVec4(im.ColText, im.NewVec4(0.52, 0.74, 1, 1))
		}
		im.PushIDStr(str)
		if im.ImageButton("Set", layer.Image.Tex.ID, im.NewVec2(20, 20)) {
			s.layer = s.texture.Layers[i]
			s.texture.SetLayer(i)
		}
		im.SameLine()
		im.Text(str)
		im.SameLine()
		if cannotDown {
			im.BeginDisabled()
		}
		if im.Button("down") {
			s.texture.Swap(i, i-1)
		}
		if cannotDown {
			im.EndDisabled()
		}
		im.SameLine()
		if cannotUp {
			im.BeginDisabled()
		}
		if im.Button("up") {
			s.texture.Swap(i, i+1)
		}
		if cannotUp {
			im.EndDisabled()
		}
		im.SameLine()
		if im.Checkbox("Show", &s.texture.Layers[i].Show) {
			s.texture.UpdateTexture()
		}
		im.PopID()
		if selected {
			im.PopStyleColor()
		}
	}

	if im.BeginPopupModal("Remove Layers") {
		im.Text("Select Layers to remove")
		for i := range s.texture.Layers {
			str := fmt.Sprintf("Layer #%d", i)
			im.Checkbox(str, &selectedLayers[i])
			im.SameLine()
			im.ImageWithBg(s.texture.Layers[i].Image.Tex.ID, im.NewVec2(15, 15))
		}
		if im.Button("Remove") {
			s.texture.Remove(selectedLayers)
			im.CloseCurrentPopup()
		}
		im.SameLine()
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	if im.BeginPopupModal("Merge Layers") {
		im.Text("Select Layers to merge:")
		for i := range s.texture.Layers {
			str := fmt.Sprintf("Layer #%d", i)
			im.Checkbox(str, &selectedLayers[i])
			im.SameLine()
			im.ImageWithBg(s.texture.Layers[i].Image.Tex.ID, im.NewVec2(15, 15))
		}
		if im.Button("Merge") {
			s.texture.Merge(selectedLayers)
			im.CloseCurrentPopup()
		}
		im.SameLine()
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	im.End()
}

var windowIdSys = util.NewIdSystem()
