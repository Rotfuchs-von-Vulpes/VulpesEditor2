package tools

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/tools/pencil"
)

func (s *data) Use() {
	currentCtx = s
}

var currentCtx *data
var ctxManager = context.New()

func Begin(id int32) {
	ctxManager.Begin(id)
}

func End() {
	ctxManager.End()
}

func New(id int32) {
	ctx := new(data)
	ctx.selectedTool = pencil.Pencil{}
	ctx.texture = texture.New(1, 1)
	ctxManager.Add(id, ctx)
}
