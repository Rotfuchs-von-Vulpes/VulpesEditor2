package texture

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/objectModel/model"
)

func (s *TextureContext) Use() {
	ctx = s
}

var ctx *TextureContext
var ctxManager = context.New()

func New(id int32, m *model.Model) {
	ctx = new(TextureContext)
	ctx.model = m
	// GenerateTexture()
	ctxManager.Add(id, ctx)
}
