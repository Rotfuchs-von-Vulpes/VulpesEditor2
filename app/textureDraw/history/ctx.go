package history

import (
	"VulpesEditor/app/context"
)

func (s *historyData) Use() {
	ctx = s
}

func (s *historyData) Reset() {
	ctx = nil
}

var ctx *historyData
var ctxManager = context.New()

func Begin(id int32) {
	ctxManager.Begin(id)
}

func End() {
	ctxManager.End()
}

func New(id int32) {
	c := new(historyData)
	c.history = make([]change, 0)
	ctxManager.Add(id, c)
}
