package view

import (
	"VulpesEditor/app/editors/blockCreator/chunk"
	"VulpesEditor/app/front/renderer"
	"math"

	im "github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl32"
)

type viewContext struct {
	camera *renderer.Camera
	mesh   *renderer.BlocksMesh
	viewer *renderer.FrameBuffer

	chunk *chunk.Chunk

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

	rayOrigin      mgl32.Vec3
	rayVersor      mgl32.Vec3
	hoverdBlockPos mgl32.Vec3
	isHovering     bool
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
	diff := mgl32.Vec3{float32(ctx.chunk.Width) / 2, float32(ctx.chunk.Height) / 2, float32(ctx.chunk.Depth) / 2}
	x := ctx.zoom * float32(math.Cos(float64(ctx.pos[0]))*math.Cos(float64(ctx.pos[1])))
	y := ctx.zoom * float32(math.Sin(float64(ctx.pos[1])))
	z := ctx.zoom * float32(math.Sin(float64(ctx.pos[0]))*math.Cos(float64(ctx.pos[1])))

	cameraPos := mgl32.Vec3{x, y, z}
	cameraFront := cameraPos.Mul(-1)
	cameraPos = cameraPos.Add(diff)

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

func linear(p1, p2, v float32) (u float32) {
	if v > 1 {
		v = 1
	} else if v < -1 {
		v = -1
	}
	v = v/2 + 0.5
	u = p1*v + p2*(1-v)
	return
}

func CalculateMouseRay(screenX, screenY float32) (rayOrigin, rayDirection mgl32.Vec3) {
	ndcX := (2.0*screenX)/viwerSize[0] - 1.0
	ndcY := 1.0 - (2.0*screenY)/viwerSize[1]

	clipCoords := mgl32.Vec4{ndcX, ndcY, -1.0, 1.0}

	invProj := ctx.camera.Proj.Inv()
	eyeCoords := invProj.Mul4x1(clipCoords)
	eyeCoords = mgl32.Vec4{eyeCoords.X(), eyeCoords.Y(), -1.0, 0.0}

	invView := ctx.camera.View.Inv()
	worldCoords4 := invView.Mul4x1(eyeCoords)

	rayDirection = mgl32.Vec3{worldCoords4.X(), worldCoords4.Y(), worldCoords4.Z()}.Normalize()

	rayOrigin = mgl32.Vec3{invView.At(0, 3), invView.At(1, 3), invView.At(2, 3)}

	return rayOrigin, rayDirection
}

func move(pos im.Vec2) {
	mousePos = [2]float32{pos.X, pos.Y}

	if mouseCanDrag {
		moveCamera()
	}

	rayOrigin, rayVersor = CalculateMouseRay(pos.X, pos.Y)

	isHovering, hoverdBlockPos = ctx.chunk.Hovered(rayOrigin, rayVersor)
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
		if buttons[0] && ctx.chunk.Destruct(rayOrigin, rayVersor) {
			ctx.mesh.SetVertices(ctx.chunk.ToBuffer())
		}
		if buttons[1] && ctx.chunk.Construct(rayOrigin, rayVersor) {
			ctx.mesh.SetVertices(ctx.chunk.ToBuffer())
		}
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

func Show(id string) {
	ctxManager.Check(id)

	// if ctx.modelComm.Changed(model.ChangeTexture | model.ChangeUnits) {
	// 	f, e := ctx.model.ToBuffer()
	// 	ctx.mesh.SetVertices(f, e)
	// }

	if toFocus {
		im.SetNextWindowFocus()
		toFocus = false
	}
	if im.Begin("Model") {
		wSize := im.ContentRegionAvail()
		width := int32(wSize.X)
		height := int32(wSize.Y)

		if viwerSize[0] != wSize.X || viwerSize[1] != wSize.Y {
			ctx.viewer.Resize(width, height)
			viwerSize[0] = wSize.X
			viwerSize[1] = wSize.Y
			aspect = wSize.Y / wSize.X
			calcCamera()
		}

		im.ImageV(
			*im.NewTextureRefTextureID(im.TextureID(ctx.viewer.Image())),
			ctx.viewer.Size(),
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

	ctx.viewer.RenderBlocks(ctx.camera, ctx.mesh)
	if isHovering {
		ctx.viewer.RenderHoveredFrame(ctx.camera, hoverdBlockPos)
	}

	im.End()
}
