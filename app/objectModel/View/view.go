package view

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
	"math"

	im "github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelContext struct {
	zoom        float32
	model       *model.Model
	modelViewer *renderer.FrameBuffer
	camera      *renderer.Camera
	mesh        *renderer.Mesh
}

var viwerSize [2]float32
var aspect float32

var times float64

func Show(id int32) {
	ctxManager.Check(id)

	cameraPos := mgl32.Vec3{3 * float32(math.Cos(times)), 2, 3 * float32(math.Sin(times))}
	// cameraPos := mgl32.Vec3{3, 2, 3}

	ctx.camera.Move(cameraPos)
	ctx.camera.Turn(cameraPos.Mul(-1))

	times += 0.01

	if im.Begin("Model") {
		wSize := im.ContentRegionAvail()
		width := int32(wSize.X)
		height := int32(wSize.Y)

		if viwerSize[0] != wSize.X || viwerSize[1] != wSize.Y {
			ctx.modelViewer.Resize(width, height)
			viwerSize[0] = wSize.X
			viwerSize[1] = wSize.Y
			aspect = wSize.Y / wSize.X
		}

		im.ImageV(
			*im.NewTextureRefTextureID(im.TextureID(ctx.modelViewer.Image())),
			ctx.modelViewer.Size(),
			im.NewVec2(0, 1),
			im.NewVec2(1, 0),
		)
	}

	ctx.modelViewer.RenderModel(ctx.camera, ctx.mesh)

	im.End()
}
