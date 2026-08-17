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

func New(id int32, m *model.Model) {
	ctx = new(ToolContext)
	ctx.model = m
	ctxManager.Add(id, ctx)
}
