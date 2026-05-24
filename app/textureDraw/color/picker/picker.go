package picker

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

var selectedColor = [4]float32{1, 0.5, 0.25, 0.5}
var seekColor1, seekColor2 bool

func Reset(change [2]bool) {
	if seekColor1 && change[0] {
		seekColor1 = false
	}
	if seekColor2 && change[1] {
		seekColor2 = false
	}
}

func Loop(color1, color2 *[4]float32) (change [2]bool) {
	imgui.Begin("Color Picker")
	imgui.ColorPicker4("Color", &selectedColor)
	if seekColor1 {
		*color1 = selectedColor
	} else if seekColor2 {
		*color2 = selectedColor
	}
	if imgui.ColorButton("Color 1", newVec4(*color1)) {
		*color1 = selectedColor
		change[0] = true
		seekColor1 = true
	}
	imgui.SameLine()
	imgui.Text("First Color")
	if imgui.ColorButton("Color 2", newVec4(*color2)) {
		*color2 = selectedColor
		change[1] = true
		seekColor2 = true
	}
	imgui.SameLine()
	imgui.Text("Second Color")

	imgui.End()
	return
}
