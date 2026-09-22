package view

import (
	"VulpesEditor/app/blockCreator/chunk"
	"VulpesEditor/app/context"
	"VulpesEditor/app/front/renderer"
	"math"
)

var ctx *viewContext

func (s *viewContext) Use() {
	ctx = s
}

var ctxManager = context.New()

func New(id string) {
	ctx := new(viewContext)
	ctx.accumulation = [2]float32{math.Pi / 4, math.Pi / 8}
	ctx.zoom = 10
	ctx.camera = renderer.NewCamera(500, 500)
	ctx.viewer = renderer.CreateFramebuffer(500, 500)
	viwerSize = [2]float32{500, 500}
	ctx.mesh = renderer.NewBlocksMesh()
	ctx.mesh.SetVertices(chunk.Test())
	ctx.mesh.Centralize(8, 8, 8)
	ctxManager.Add(id, ctx)

	moveCamera()
}
