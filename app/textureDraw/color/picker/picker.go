package picker

import "github.com/AllenDang/cimgui-go/imgui"

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

var selectedColor = [4]float32{0, 0, 0, 1}

func Loop(color1, color2 *[4]float32) (change [2]bool) {
	imgui.Begin("Color Picker")
	imgui.ColorPicker4("Color", &selectedColor)
	if imgui.ColorButton("Color 1", newVec4(*color1)) {
		*color1 = selectedColor
	}
	imgui.SameLine()
	imgui.Text("First Color")
	if imgui.ColorButton("Color 2", newVec4(*color2)) {
		*color2 = selectedColor
	}
	imgui.SameLine()
	imgui.Text("Second Color")
	imgui.End()
	return
}
