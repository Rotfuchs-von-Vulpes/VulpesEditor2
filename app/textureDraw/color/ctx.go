package color

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/color/palette"
)

type data struct {
	color1 [4]float32
	color2 [4]float32
}

func (s *data) Use() {
	currentCtx = s
}

var currentCtx *data
var ctxManager = context.New()

func Begin(id int32) {
	ctxManager.Begin(id)
	palette.Begin(id)
}

func End() {
	ctxManager.End()
	palette.End()
}

func New(id int32) {
	ctx := new(data)
	palette.New(id)
	ctxManager.Add(id, ctx)
}
