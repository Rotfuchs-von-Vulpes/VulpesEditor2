package tools

import (
	"VulpesEditor/app/textureDraw/tools/bucket"
	"VulpesEditor/app/textureDraw/tools/pencil"

	"github.com/AllenDang/cimgui-go/imgui"
)

var selectedTool Tool

func Init() {
	selectedTool = pencil.Pencil{}
}

type Tool interface {
	SendTexture(colors [][][4]float32, width, height uint32)
	ButtonPress(pos [2]int32)
	ButtonRelease(pos [2]int32)
	Move(pos1, pos2 [2]int32)
	Visualize() [][2]int32
	Change() [][2]int32
	Reset()
}

func ButtonPress(pos [2]int32, colors [][][4]float32, width, height uint32) {
	selectedTool.SendTexture(colors, width, height)
	selectedTool.ButtonPress(pos)
}

func ButtonRelease(pos [2]int32) {
	selectedTool.ButtonRelease(pos)
}

func Move(pos1, pos2 [2]int32) {
	selectedTool.Move(pos1, pos2)
}

func Visualize() [][2]int32 {
	return selectedTool.Visualize()
}

func Change() [][2]int32 {
	return selectedTool.Change()
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
	imgui.End()
}
