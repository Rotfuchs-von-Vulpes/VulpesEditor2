package picker

import (
	im "github.com/AllenDang/cimgui-go/imgui"
)

func newVec4(vec [4]float32) im.Vec4 {
	return im.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

func Reset(change [2]bool) {

}

func Loop(color1, color2 *[4]float32) (change [3]bool) {
	im.Begin("Color Picker")
	if im.ColorPicker4V("Color", color1, im.ColorEditFlagsPickerHueWheel, nil) {
		change[0] = true
	}
	im.ColorButton("Color 1", newVec4(*color1))
	im.SameLine()
	im.Text("First Color")
	im.ColorButton("Color 2", newVec4(*color2))
	im.SameLine()
	im.Text("Second Color")
	if im.Button("Swap") {
		change[2] = true
		temp := *color1
		*color1 = *color2
		*color2 = temp
	}
	im.End()
	return
}
