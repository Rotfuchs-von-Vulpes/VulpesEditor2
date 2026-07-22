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

var ctx *colorData
var ctxManager = context.New()

func New(id int32) {
	c := new(colorData)
	c.color1 = [4]float32{1, 1, 1, 1}
	c.color2 = [4]float32{0, 0, 0, 1}
	palette.New(id)
	ctxManager.Add(id, c)
}
