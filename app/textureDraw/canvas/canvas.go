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

	"github.com/AllenDang/cimgui-go/imgui"
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
	viewrSize     [2]float32
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
	pos1[0] = s.viewrSize[0] - s.mousePos[0] - s.pos[0]*s.aspect
	pos1[1] = s.viewrSize[1] - s.mousePos[1] - s.pos[1]
	pos1[0] = 2 * ((pos1[0] / s.viewrSize[0]) - 0.5)
	pos1[1] = 2 * ((pos1[1] / s.viewrSize[1]) - 0.5)
	pos1[0] = pos1[0] / (s.zoom * s.aspect * s.texture.Aspect)
	pos1[1] = pos1[1] / s.zoom
	pixelPos[0] = int32(math.Floor(float64(s.texture.Width) * float64(pos1[0]/2+0.5)))
	pixelPos[1] = int32(math.Floor(float64(s.texture.Height) * (1 - float64(pos1[1]/2+0.5))))
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
		imgui.SetNextWindowFocus()
		toFocus = false
	}
	imgui.BeginV(s.windowName, nil, imgui.WindowFlagsMenuBar)
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Texture") {
			if imgui.MenuItemBool("Save Texture") {
				toPop = "Not Implement"
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Export") {
			if imgui.MenuItemBool("Export as PNG Image") {
				toPop = "Export PNG"
			}
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
	if toPop != "" {
		imgui.OpenPopupStr(toPop)
		toPop = ""
	}
	if imgui.BeginPopupModal("Export PNG") {
		imgui.InputTextWithHint("File Name", "texure_"+strconv.FormatInt(int64(s.texture.Id), 10)+".png", &textureFileName, imgui.InputTextFlagsNone, nil)
		imgui.InputTextWithHint("File Path", "", &textureFilePath, imgui.InputTextFlagsNone, nil)
		if imgui.Button("Save") {
			if ok := s.texture.SaveTextureAsFile(textureFileName, textureFilePath); ok {
				textureFileName = ""
				imgui.CloseCurrentPopup()
			}
		}
		if imgui.Button("Cancel") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}
	if imgui.BeginPopupModal("Not Implement") {
		imgui.Text("Not implement yet!")
		if imgui.Button("OK") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	wSize := imgui.ContentRegionAvail()
	width := int32(wSize.X)
	height := int32(wSize.Y)

	if s.viewrSize[0] != wSize.X || s.viewrSize[1] != wSize.Y {
		s.textureViewer.Resize(width, height)
		s.viewrSize[0] = wSize.X
		s.viewrSize[1] = wSize.Y
		s.aspect = wSize.Y / wSize.X
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
		s.buttonPress(io.MouseClicked())
		s.buttonRelease(io.MouseReleased())
	} else if s.painting {
		s.buttonRelease([5]bool{true, true, true, false, false})
	}
	if imgui.IsWindowFocused() && lastEditId == s.id {
		io := imgui.CurrentContext().IO()
		if io.KeyCtrl() && imgui.IsKeyPressedBoolV(imgui.KeyZ, true) {
			s.layer.Undo()
		}
		if io.KeyCtrl() && imgui.IsKeyPressedBoolV(imgui.KeyY, true) {
			s.layer.Redo()
		}
		buttons := io.MouseClicked()
		if buttons[3] {
			s.layer.Undo()
		}
		if buttons[4] {
			s.layer.Redo()
		}
	}
	imgui.End()
	s.textureViewer.RenderTexture(s.texture.GlID, s.zoom, s.pos, float32(s.texture.Width), float32(s.texture.Height))
}

// var AllTextures []*texture.Texture
var AllCtx []*TextureContext
var lastEditCtx *TextureContext

func createCtx(tex *texture.Texture) (ctx TextureContext) {
	ctx.id = windowIdSys.GetID()
	ctx.windowName = "Texture #" + strconv.FormatUint(uint64(ctx.id), 10)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.CreateFramebuffer(500, 500)
	ctx.viewrSize = [2]float32{500, 500}
	ctx.texture = textureEdit.New(tex)
	ctx.layer = ctx.texture.Layers[0]
	lastEditId = ctx.id
	return
}

func AddTexture(width, height uint32) {
	tex := texture.New(width, height)
	ctx := createCtx(tex)
	// AllTextures = append(AllTextures, tex)
	AllCtx = append(AllCtx, &ctx)
	lastEditCtx = AllCtx[len(AllCtx)-1]
}

func OpenTexture(tex *texture.Texture) {
	createCtx(tex)
}

var selectedLayers = []bool{}

func updateLayerId() {
	id := lastEditCtx.layer.Id
	found := false
	for _, layer := range lastEditCtx.texture.Layers {
		if layer.Id == id {
			found = true
			break
		}
	}
	if !found {
		lastEditCtx.layer = lastEditCtx.texture.Layers[len(lastEditCtx.texture.Layers)-1]
	}
}

func ShowLayers() {
	var ctx = lastEditCtx
	imgui.Begin("Layers")
	if imgui.Button("Add") {
		ctx.texture.AddLayer()
	}
	imgui.SameLine()
	if imgui.Button("Remove") {
		imgui.OpenPopupStr("Remove Layers")
	}
	imgui.SameLine()
	if imgui.Button("Merge") {
		imgui.OpenPopupStr("Merge Layers")
	}

	for i, layer := range ctx.texture.Layers {
		downDis := i == 0
		upDis := i == len(ctx.texture.Layers)-1
		selected := layer.Id == ctx.layer.Id
		str := fmt.Sprintf("Layer #%d", i)
		if selected {
			// #86BDFFFF
			imgui.PushStyleColorVec4(imgui.ColText, imgui.NewVec4(0.52, 0.74, 1, 1))
		}
		imgui.PushIDStr(str)
		if imgui.Button("Set") {
			ctx.layer = ctx.texture.Layers[i]
			ctx.texture.SetLayer(i)
		}
		imgui.SameLine()
		imgui.Text(str)
		imgui.SameLine()
		if downDis {
			imgui.BeginDisabled()
		}
		if imgui.Button("down") {
			ctx.texture.Swap(i, i-1)
		}
		if downDis {
			imgui.EndDisabled()
		}
		imgui.SameLine()
		if upDis {
			imgui.BeginDisabled()
		}
		if imgui.Button("up") {
			ctx.texture.Swap(i, i+1)
		}
		if upDis {
			imgui.EndDisabled()
		}
		imgui.SameLine()
		if imgui.Checkbox("Show", &ctx.texture.Layers[i].Show) {
			ctx.texture.UpdateTexture()
		}
		imgui.PopID()
		if selected {
			imgui.PopStyleColor()
		}
	}

	if imgui.BeginPopupModal("Remove Layers") {
		imgui.Text("Select Layers to remove")
		for i := range ctx.texture.Layers {
			str := fmt.Sprintf("Layer #%d", i)
			imgui.Checkbox(str, &selectedLayers[i])
		}
		if imgui.Button("Remove") {
			ctx.texture.Remove(selectedLayers)
			updateLayerId()
			imgui.CloseCurrentPopup()
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	if imgui.BeginPopupModal("Merge Layers") {
		imgui.Text("Select Layers to merge:")
		for i := range ctx.texture.Layers {
			str := fmt.Sprintf("Layer #%d", i)
			imgui.Checkbox(str, &selectedLayers[i])
		}
		if imgui.Button("Merge") {
			ctx.texture.Merge(selectedLayers)
			updateLayerId()
			imgui.CloseCurrentPopup()
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	imgui.End()
}

var windowIdSys = util.NewIdSystem()
