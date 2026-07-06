package tools

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/tools/bucket"
	"VulpesEditor/app/textureDraw/tools/colorPicker"
	"VulpesEditor/app/textureDraw/tools/eraser"
	"VulpesEditor/app/textureDraw/tools/line"
	"VulpesEditor/app/textureDraw/tools/pencil"
	"VulpesEditor/app/textureDraw/tools/rectangle"

	"github.com/AllenDang/cimgui-go/imgui"
)

type tool interface {
	ButtonPress(pos [2]int32, secondButton bool)
	ButtonRelease(pos [2]int32)
	Move(pos1, pos2 [2]int32)
	Reset()
}

type drawingTool interface {
	tool
	Visualize() []texture.PixelEdit
	Change() []texture.PixelEdit
}

var selectedTool tool

var Color1 *[4]float32
var Color2 *[4]float32

func Init() {
	selectedTool = pencil.Pencil{}
}

var Texture *texture.Texture = texture.New(1, 1)

func ButtonPress(pos [2]int32, secondButton bool) {
	if _, ok := selectedTool.(bucket.Bucket); ok {
		bucket.Texture = Texture
	} else if _, ok := selectedTool.(colorPicker.ColorPicker); ok {
		colorPicker.Texture = Texture
	}
	selectedTool.ButtonPress(pos, secondButton)
}

func ButtonRelease(pos [2]int32) {
	selectedTool.ButtonRelease(pos)
}

func Move(pos1, pos2 [2]int32) {
	selectedTool.Move(pos1, pos2)
}

func Visualize() []texture.PixelEdit {
	if s, ok := selectedTool.(drawingTool); ok {
		return s.Visualize()
	}
	return nil
}

func Change() []texture.PixelEdit {
	if s, ok := selectedTool.(drawingTool); ok {
		return s.Change()
	}
	return nil
}

func Show() {
	imgui.Begin("Tools")
	if imgui.Button("Pencil") {
		selectedTool.Reset()
		selectedTool = pencil.Pencil{}
	}
	if imgui.Button("Bucket") {
		selectedTool.Reset()
		selectedTool = bucket.Bucket{}
	}
	if imgui.Button("Line") {
		selectedTool.Reset()
		selectedTool = line.Line{}
	}
	if imgui.Button("Rect") {
		selectedTool.Reset()
		selectedTool = rectangle.Rectangle{}
	}
	if imgui.Button("Eraser") {
		selectedTool.Reset()
		selectedTool = eraser.Eraser{}
	}
	if imgui.Button("Color Picker") {
		selectedTool.Reset()
		selectedTool = colorPicker.ColorPicker{}
	}
	imgui.End()
}
