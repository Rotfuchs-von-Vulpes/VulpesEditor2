package tools

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/tools/bucket"
	"VulpesEditor/app/textureDraw/tools/colorPicker"
	"VulpesEditor/app/textureDraw/tools/eraser"
	"VulpesEditor/app/textureDraw/tools/line"
	"VulpesEditor/app/textureDraw/tools/pencil"
	"VulpesEditor/app/textureDraw/tools/rectangle"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type tool interface {
	ButtonPress(pos [2]int32, firstButton bool)
	ButtonRelease(pos [2]int32)
	Move(pos1, pos2 [2]int32)
	Reset()
}

type drawingTool interface {
	tool
	Visualize() []texture.PixelEdit
	Change() []texture.PixelEdit
}

type toolsData struct {
	selectedTool tool
	texture      *texture.Texture
}

func Resize(width, height uint32) {
	ctx.texture.Resize(width, height)
}

func SetColors(c [][4]float32) {
	ctx.texture.Colors = c
}

func ButtonPress(pos [2]int32, firstButton bool) {
	if _, ok := ctx.selectedTool.(bucket.Bucket); ok {
		bucket.Texture = ctx.texture
	} else if _, ok := ctx.selectedTool.(colorPicker.ColorPicker); ok {
		colorPicker.Texture = ctx.texture
	}
	ctx.selectedTool.ButtonPress(pos, firstButton)
}

func ButtonRelease(pos [2]int32) {
	ctx.selectedTool.ButtonRelease(pos)
}

func Move(pos1, pos2 [2]int32) {
	ctx.selectedTool.Move(pos1, pos2)
}

func Visualize() []texture.PixelEdit {
	if s, ok := ctx.selectedTool.(drawingTool); ok {
		return s.Visualize()
	}
	return nil
}

func Change() []texture.PixelEdit {
	if s, ok := ctx.selectedTool.(drawingTool); ok {
		return s.Change()
	}
	return nil
}

func Show() {
	im.Begin("Tools")
	if im.Button("Pencil") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = pencil.Pencil{}
	}
	if im.Button("Bucket") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = bucket.Bucket{}
	}
	if im.Button("Line") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = line.Line{}
	}
	if im.Button("Rect") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = rectangle.Rectangle{}
	}
	if im.Button("Eraser") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = eraser.Eraser{}
	}
	if im.Button("Color Picker") {
		ctx.selectedTool.Reset()
		ctx.selectedTool = colorPicker.ColorPicker{}
	}
	im.End()
}
