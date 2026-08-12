package view

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
)

func (s *ModelContext) Use() {
	ctx = s
}

var ctx *ModelContext
var ctxManager = context.New()

func New(id int32) {
	OpenModel(id)
}

func OpenModel(id int32) {
	ctx = new(ModelContext)
	ctx.zoom = 0.9
	ctx.camera = renderer.NewCamera(500, 500)
	ctx.modelViewer = renderer.CreateFramebuffer(500, 500)
	ctx.model = model.NewModel()
	f, e := ctx.model.ToBuffer()
	ctx.mesh = renderer.NewMesh()
	ctx.mesh.SetVertices(f, e)
	viwerSize = [2]float32{500, 500}
	ctxManager.Add(id, ctx)
}

func Save(w *file.ArchiveWriter) {
	// ctx.texture.Save(w)
}

func Open(id int32, r *file.ArchiveReader) (err error) {
	// tex, err := textureEdit.Open(r)
	// if err != nil {
	// 	return err
	// }
	OpenModel(id)
	return
}
