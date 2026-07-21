package color

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/color/palette"
)

type colorData struct {
	color1 [4]float32
	color2 [4]float32
}

func (s *colorData) Use() {
	ctx = s
}

func (s *colorData) Reset() {
	ctx = nil
}

var ctx *colorData
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
	c := new(colorData)
	palette.New(id)
	ctxManager.Add(id, c)
}
