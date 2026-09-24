package texture

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/editors/objectModel/model"
)

func (s *TextureContext) Use() {
	ctx = s
}

var ctx *TextureContext
var ctxManager = context.New()

func New(id, name string, m *model.Model) {
	ctx = new(TextureContext)
	ctx.name = name
	ctx.model = m
	ctx.modelComm = m.NewChangeAsker()
	ctxManager.Add(id, ctx)
}
