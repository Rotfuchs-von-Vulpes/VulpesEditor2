package picker

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

func Reset(change [2]bool) {

}

func Loop(color1, color2 *[4]float32) (change [3]bool) {
	imgui.Begin("Color Picker")
	if imgui.ColorPicker4("Color", color1) {
		change[0] = true
	}
	imgui.ColorButton("Color 1", newVec4(*color1))
	imgui.SameLine()
	imgui.Text("First Color")
	imgui.ColorButton("Color 2", newVec4(*color2))
	imgui.SameLine()
	imgui.Text("Second Color")
	if imgui.Button("Swap") {
		change[2] = true
		temp := *color1
		*color1 = *color2
		*color2 = temp
	}
	imgui.End()
	return
}
