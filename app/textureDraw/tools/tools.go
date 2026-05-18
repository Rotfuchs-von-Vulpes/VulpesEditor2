package tools

import "VulpesEditor/app/textureDraw/tools/pencil"

var selectedTool Tool

func Init() {
	selectedTool = pencil.Pencil{}
}

type Tool interface {
	ButtonPress(pos [2]int32)
	ButtonRelease(pos [2]int32)
	Move(pos1, pos2 [2]int32)
	Visualize() [][2]int32
	Change() [][2]int32
}

func ButtonPress(pos [2]int32) {
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
