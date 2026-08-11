package view

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type ModelContext struct {
	zoom        float32
	model       *model.Model
	modelViewer *renderer.FrameBuffer
	length      int32
}

var viwerSize [2]float32
var aspect float32

func Show(id int32) {
	ctxManager.Check(id)

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

	ctx.modelViewer.RenderModel(ctx.length)

	im.End()
}
