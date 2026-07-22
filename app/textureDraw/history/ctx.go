package history

import (
	"VulpesEditor/app/context"
)

func (s *historyData) Use() {
	ctx = s
}

var ctx *historyData
var ctxManager = context.New()

func New(id int32) {
	c := new(historyData)
	c.history = make([]change, 0)
	ctxManager.Add(id, c)
}
