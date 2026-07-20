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

type data struct {
	selectedTool tool
	texture      *texture.Texture
}

func Resize(width, height uint32) {
	currentCtx.texture.Resize(width, height)
}

func SetColors(c [][4]float32) {
	currentCtx.texture.Colors = c
}

func ButtonPress(pos [2]int32, firstButton bool) {
	if _, ok := currentCtx.selectedTool.(bucket.Bucket); ok {
		bucket.Texture = currentCtx.texture
	} else if _, ok := currentCtx.selectedTool.(colorPicker.ColorPicker); ok {
		colorPicker.Texture = currentCtx.texture
	}
	currentCtx.selectedTool.ButtonPress(pos, firstButton)
}

func ButtonRelease(pos [2]int32) {
	currentCtx.selectedTool.ButtonRelease(pos)
}

func Move(pos1, pos2 [2]int32) {
	currentCtx.selectedTool.Move(pos1, pos2)
}

func Visualize() []texture.PixelEdit {
	if s, ok := currentCtx.selectedTool.(drawingTool); ok {
		return s.Visualize()
	}
	return nil
}

func Change() []texture.PixelEdit {
	if s, ok := currentCtx.selectedTool.(drawingTool); ok {
		return s.Change()
	}
	return nil
}

func Show() {
	im.Begin("Tools")
	if im.Button("Pencil") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = pencil.Pencil{}
	}
	if im.Button("Bucket") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = bucket.Bucket{}
	}
	if im.Button("Line") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = line.Line{}
	}
	if im.Button("Rect") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = rectangle.Rectangle{}
	}
	if im.Button("Eraser") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = eraser.Eraser{}
	}
	if im.Button("Color Picker") {
		currentCtx.selectedTool.Reset()
		currentCtx.selectedTool = colorPicker.ColorPicker{}
	}
	im.End()
}
