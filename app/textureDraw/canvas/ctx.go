package canvas

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/canvas/texture"
)

func (s *TextureContext) Use() {
	ctx = s
}

func (s *TextureContext) Reset() {
	ctx = nil
}

var ctx *TextureContext
var ctxManager = context.New()

func Begin(id int32) {
	ctxManager.Begin(id)
}

func End() {
	ctxManager.End()
}

func New(id int32, w, h uint32) {
	c := createCtx(texture.New(w, h))
	ctxManager.Add(id, c)
}
