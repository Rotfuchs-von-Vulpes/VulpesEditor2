package palette

import (
	"VulpesEditor/app/context"
)

func (s *paletteData) Use() {
	ctx = s
}

func (s *paletteData) Reset() {
	ctx = s
}

var ctx *paletteData
var ctxManager = context.New()

func Begin(id int32) {
	ctxManager.Begin(id)
}

func End() {
	ctxManager.End()
}

func New(id int32) {
	c := new(paletteData)
	c.palettes = make(map[int32]bool)
	for i, p := range palettes {
		if i == 0 {
			c.palettes[p.id] = true
			break
		}
	}
	ctxManager.Add(id, c)
}
