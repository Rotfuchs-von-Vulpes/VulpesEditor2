package tools

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/objectModel/model"
)

func (s *ToolContext) Use() {
	ctx = s
}

var ctx *ToolContext
var ctxManager = context.New()

func New(id string, name string, m *model.Model) {
	ctx = new(ToolContext)
	ctx.name = name
	ctx.model = m
	ctx.sizeInput = [3]float32{16, 16, 16}
	ctxManager.Add(id, ctx)
}
