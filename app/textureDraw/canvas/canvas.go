package canvas

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/canvas/textureEdit"
	"VulpesEditor/app/textureDraw/tools"
	"VulpesEditor/app/util"
	"math"
	"strconv"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type TextureContext struct {
	zoom float32
	pos  [2]float32

	textureViewer *renderer.FrameBuffer

	texture *textureEdit.TextureEdit
}

var (
	viwerSize     [2]float32
	aspect        float32
	mousePos      [2]float32
	lastMousePos  [2]float32
	mousePressPos [2]float32
	mouseCanDrag  bool
	painting      bool
	firstButton   bool
	pixelPos      [2]int32
	lastPixelPos  [2]int32
)

func resetPreview() {
	ctx.texture.ResetPreview()
}

func reset() {
	ctx.zoom = 1
	ctx.pos = [2]float32{0, 0}
}

func pixelPosMouse() (pixelPos [2]int32) {
	var pos1 [2]float32
	pos1[0] = viwerSize[0] - mousePos[0] - ctx.pos[0]*aspect
	pos1[1] = viwerSize[1] - mousePos[1] - ctx.pos[1]
	pos1[0] = 2 * ((pos1[0] / viwerSize[0]) - 0.5)
	pos1[1] = 2 * ((pos1[1] / viwerSize[1]) - 0.5)
	pos1[0] = pos1[0] / (ctx.zoom * aspect * ctx.texture.Aspect)
	pos1[1] = pos1[1] / ctx.zoom
	pixelPos[0] = int32(math.Floor(float64(ctx.texture.Width) * float64(pos1[0]/2+0.5)))
	pixelPos[1] = int32(math.Floor(float64(ctx.texture.Height) * (1 - float64(pos1[1]/2+0.5))))
	return
}

func scroll(yoffset float32) {
	mouse := [2]float32{mousePos[0] - 0.5*viwerSize[0], mousePos[1] - 0.5*viwerSize[1]}
	position := [2]float32{(ctx.pos[0] + mouse[0]) / ctx.zoom, (ctx.pos[1] + mouse[1]) / ctx.zoom}

	if yoffset < 0 {
		ctx.zoom *= 0.9
	} else if yoffset > 0 {
		ctx.zoom *= 1.1
	}

	ctx.pos = [2]float32{ctx.zoom*position[0] - mouse[0], ctx.zoom*position[1] - mouse[1]}
}

func move(pos im.Vec2) {
	mousePos = [2]float32{viwerSize[0] - pos.X, pos.Y}

	pixel := pixelPosMouse()
	if (pixel[0] != lastPixelPos[0] || pixel[1] != lastPixelPos[1]) && painting {
		tools.Move(lastPixelPos, pixel)
		ctx.texture.ResetPreview()
		pixels := tools.Visualize()
		ctx.texture.ChangePreview(pixels)
		lastPixelPos = pixel
	}

	if mouseCanDrag {
		ctx.pos[0] = (lastMousePos[0]-mousePos[0])/aspect + mousePressPos[0]
		ctx.pos[1] = lastMousePos[1] - mousePos[1] + mousePressPos[1]
	}
}

var secondButton bool

func buttonPress(buttons [5]bool) {
	if buttons[2] {
		lastMousePos = mousePos
		mousePressPos = ctx.pos
		mouseCanDrag = true
		toFocus = true
	}
	if buttons[0] || buttons[1] {
		painting = true
		firstButton = buttons[0]
		tools.Resize(ctx.texture.Width, ctx.texture.Height)
		tools.SetColors(ctx.texture.LayerColors())
		tools.ButtonPress(pixelPosMouse(), buttons[0])
		toFocus = true
	}
}

func buttonRelease(buttons [5]bool) {
	if buttons[2] {
		mouseCanDrag = false
	}
	if buttons[0] || buttons[1] {
		painting = false
		tools.ButtonRelease(pixelPosMouse())
		pixels := tools.Change()
		ctx.texture.LayerChange(pixels)
		ctx.texture.ResetPreview()
	}
}

var textureFileName string
var textureFilePath string
var toFocus bool

func Show() {
	var toPop string
	if toFocus {
		im.SetNextWindowFocus()
		toFocus = false
	}
	im.BeginV("Texture", nil, im.WindowFlagsMenuBar)
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
		im.InputTextWithHint("File Name", "texure_"+strconv.FormatInt(int64(ctx.texture.Id), 10)+".png", &textureFileName, im.InputTextFlagsNone, nil)
		im.InputTextWithHint("File Path", "", &textureFilePath, im.InputTextFlagsNone, nil)
		if im.Button("Save") {
			if ok := ctx.texture.SaveTextureAsFile(textureFileName, textureFilePath); ok {
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

	if viwerSize[0] != wSize.X || viwerSize[1] != wSize.Y {
		ctx.textureViewer.Resize(width, height)
		viwerSize[0] = wSize.X
		viwerSize[1] = wSize.Y
		aspect = wSize.Y / wSize.X
	}

	im.ImageV(
		*im.NewTextureRefTextureID(im.TextureID(ctx.textureViewer.Image())),
		ctx.textureViewer.Size(),
		im.NewVec2(0, 1),
		im.NewVec2(1, 0),
	)
	if im.IsItemHovered() {
		io := im.CurrentContext().IO()
		if io.MouseWheel() != 0 {
			scroll(io.MouseWheel())
		}
		mouse_pos_abs := io.MousePos()
		screen_pos_abs := im.ItemRectMin()
		var mouse_pos_rel im.Vec2
		mouse_pos_rel.X = mouse_pos_abs.X - screen_pos_abs.X
		mouse_pos_rel.Y = mouse_pos_abs.Y - screen_pos_abs.Y
		move(mouse_pos_rel)
		buttonPress(io.MouseClicked())
		buttonRelease(io.MouseReleased())
	} else if painting {
		buttonRelease([5]bool{true, true, true, false, false})
	}
	im.End()
	ctx.textureViewer.RenderTexture(ctx.texture.GlID, ctx.zoom, ctx.pos, float32(ctx.texture.Width), float32(ctx.texture.Height))
}

func ShowLayers() {
	ctx.texture.ShowLayers()
}

func createCtx(tex *texture.Texture) (ctx *TextureContext) {
	ctx = new(TextureContext)
	// ctx.id = windowIdSys.GetID()
	// ctx.windowName = "Texture #" + strconv.FormatUint(uint64(ctx.id), 10)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.CreateFramebuffer(500, 500)
	viwerSize = [2]float32{500, 500}
	ctx.texture = textureEdit.New(tex)
	//ctx.layer = ctx.texture.Layers[0]

	return
}

func OpenTexture(tex *texture.Texture) {
	createCtx(tex)
}

var windowIdSys = util.NewIdSystem()
