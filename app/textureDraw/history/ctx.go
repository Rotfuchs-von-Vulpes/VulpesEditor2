package history

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
	ctx.history = make([]change, 0)
	ctxManager.Add(id, ctx)
}
