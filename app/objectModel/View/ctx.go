package view

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
	"VulpesEditor/app/objectModel/texture"
	"math"
)

func (s *ModelContext) Use() {
	ctx = s
}

var ctx *ModelContext
var ctxManager = context.New()

func New(id string, model *model.Model, mesh *renderer.ModelMesh) {
	OpenModel(id, model, mesh)
}

func OpenModel(id string, model *model.Model, mesh *renderer.ModelMesh) {
	ctx = new(ModelContext)
	ctx.accumulation = [2]float32{math.Pi / 4, math.Pi / 8}
	ctx.zoom = 3
	ctx.camera = renderer.NewCamera(500, 500)
	ctx.viewer = renderer.CreateFramebuffer(500, 500)
	ctx.model = model
	ctx.modelComm = model.NewChangeAsker()
	f, e := ctx.model.ToBuffer()
	ctx.mesh = mesh
	ctx.mesh.SetVertices(f, e)
	viwerSize = [2]float32{500, 500}
	ctxManager.Add(id, ctx)

	texture.OnChange(id, func() {
		ctx.mesh.SetTexture(texture.GetData())
	})

	moveCamera()
}

func Save(w *file.ArchiveWriter) {
	// ctx.texture.Save(w)
}

func Open(id string, model *model.Model, mesh *renderer.ModelMesh, r *file.ArchiveReader) (err error) {
	OpenModel(id, model, mesh)
	return
}
