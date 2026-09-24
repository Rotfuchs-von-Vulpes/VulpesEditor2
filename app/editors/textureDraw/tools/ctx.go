package tools

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/editors/textureDraw/canvas/texture"
	"VulpesEditor/app/editors/textureDraw/tools/pencil"
)

func (s *toolsData) Use() {
	ctx = s
}

var ctx *toolsData
var ctxManager = context.New()

func New(id string, w, h int32) {
	c := new(toolsData)
	c.selectedTool = pencil.Pencil{}
	c.texture = texture.New(w, h)
	ctxManager.Add(id, c)
}
