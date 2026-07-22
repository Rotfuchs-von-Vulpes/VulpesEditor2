package tools

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/tools/pencil"
)

func (s *toolsData) Use() {
	ctx = s
}

var ctx *toolsData
var ctxManager = context.New()

func New(id int32, w, h uint32) {
	c := new(toolsData)
	c.selectedTool = pencil.Pencil{}
	c.texture = texture.New(w, h)
	ctxManager.Add(id, c)
}
