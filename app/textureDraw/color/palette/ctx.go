package palette

import (
	"VulpesEditor/app/context"
)

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
	ctx.palettes = make(map[int32]bool)
	for i, p := range palettes {
		if i == 0 {
			ctx.palettes[p.id] = true
			break
		}
	}
	ctxManager.Add(id, ctx)
}
