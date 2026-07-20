package tools

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/tools/pencil"
)

type data struct {
	selectedTool tool
}

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
	ctxManager.Add(id, ctx)
}
