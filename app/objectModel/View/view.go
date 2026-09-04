package view

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
	"math"

	im "github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelContext struct {
	model       *model.Model
	modelViewer *renderer.FrameBuffer
	camera      *renderer.Camera
	mesh        *renderer.Mesh

	pos          [2]float32
	accumulation [2]float32
	zoom         float32
}

var (
	viwerSize [2]float32
	aspect    float32 = 1

	mousePressedPos [2]float32
	mousePos        [2]float32
	pressedPos      [2]float32
	mouseCanDrag    bool

	firstButton bool
	toFocus     bool
)

func moveCamera() {
	ctx.pos[0] = ctx.accumulation[0] + (mousePos[0]-mousePressedPos[0])/(200*aspect)
	ctx.pos[1] = ctx.accumulation[1] + (mousePos[1]-mousePressedPos[1])/200

	if ctx.pos[1] > math.Pi/2-.01 {
		ctx.pos[1] = math.Pi/2 - .01
	} else if ctx.pos[1] < -math.Pi/2+.01 {
		ctx.pos[1] = -math.Pi/2 + .01
	}

	_, f1 := math.Modf(float64(ctx.pos[0]) / (2 * math.Pi))
	_, f2 := math.Modf(float64(ctx.pos[1]) / (2 * math.Pi))
	ctx.pos[0] = float32(2 * math.Pi * f1)
	ctx.pos[1] = float32(2 * math.Pi * f2)

	calcCamera()
}

func calcCamera() {
	x := ctx.zoom * float32(math.Cos(float64(ctx.pos[0]))*math.Cos(float64(ctx.pos[1])))
	y := ctx.zoom * float32(math.Sin(float64(ctx.pos[1])))
	z := ctx.zoom * float32(math.Sin(float64(ctx.pos[0]))*math.Cos(float64(ctx.pos[1])))

	cameraPos := mgl32.Vec3{x, y, z}
	cameraFront := cameraPos.Mul(-1)

	ctx.camera.Move(cameraPos)
	ctx.camera.Turn(cameraFront)
}

func scroll(yoffset float32) {
	if yoffset < 0 {
		ctx.zoom *= 1.1
	} else if yoffset > 0 {
		ctx.zoom *= 0.9
	}

	calcCamera()
}

func move(pos im.Vec2) {
	mousePos = [2]float32{pos.X, pos.Y}

	if mouseCanDrag {
		moveCamera()
	}
}

var secondButton bool

func buttonPress(buttons [5]bool) {
	if buttons[2] {
		mousePressedPos = mousePos
		mouseCanDrag = true
		toFocus = true
	}
	if buttons[0] || buttons[1] {
		firstButton = buttons[0]
		toFocus = true
	}
}

func buttonRelease(buttons [5]bool) {
	if buttons[2] {
		mouseCanDrag = false
		ctx.accumulation = ctx.pos

		mousePos = [2]float32{0, 0}
		mousePressedPos = [2]float32{0, 0}
	}
	if buttons[0] || buttons[1] {

	}
}

func Show(id int32) {
	ctxManager.Check(id)

	if ctx.model.Changed {
		f, e := ctx.model.ToBuffer()
		ctx.mesh.SetVertices(f, e)
	}

	if toFocus {
		im.SetNextWindowFocus()
		toFocus = false
	}
	if im.Begin("Model") {
		wSize := im.ContentRegionAvail()
		width := int32(wSize.X)
		height := int32(wSize.Y)

		if viwerSize[0] != wSize.X || viwerSize[1] != wSize.Y {
			ctx.modelViewer.Resize(width, height)
			viwerSize[0] = wSize.X
			viwerSize[1] = wSize.Y
			aspect = wSize.Y / wSize.X
			calcCamera()
		}

		im.ImageV(
			*im.NewTextureRefTextureID(im.TextureID(ctx.modelViewer.Image())),
			ctx.modelViewer.Size(),
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
		} else {
			buttonRelease([5]bool{true, true, true, false, false})
		}
	}

	ctx.modelViewer.RenderModel(ctx.camera, ctx.mesh)

	im.End()
}

func Model() *model.Model {
	return ctx.model
}
