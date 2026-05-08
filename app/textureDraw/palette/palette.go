package palette

import (
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
)

var color1 [4]float32
var color2 [4]float32
var color1Idx int32 = -1
var color2Idx int32 = -1
var color1Mark imgui.Vec4
var color2Mark imgui.Vec4
var colorCount int32
var colors [][4]float32

func highContrast(rgba [4]float32) imgui.Vec4 {
	// l := 0.2126*rgb[0] + 0.7152*rgb[1] + 0.0722*rgb[2]
	hsv := [4]float32{}
	imgui.ColorConvertRGBtoHSV(rgba[0], rgba[1], rgba[2], &hsv[0], &hsv[1], &hsv[2])
	if hsv[0] = hsv[0] + 0.5; hsv[0] > 1 {
		hsv[0] -= 1
	}
	final := [4]float32{}
	imgui.ColorConvertHSVtoRGB(hsv[0], hsv[1], hsv[2], &final[0], &final[1], &final[2])
	final[3] = 1
	return newVec4(final)
}

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

func Init() {
	colorCount = 20
	var step float32 = 1 / float32(colorCount)
	var hue float32 = 0
	for i := int32(0); i < colorCount; i++ {
		var rgb [3]float32
		imgui.ColorConvertHSVtoRGB(hue, 1, 1, &rgb[0], &rgb[1], &rgb[2])
		hue += step
		rgba := [4]float32{rgb[0], rgb[1], rgb[2], 1}
		colors = append(colors, rgba)
	}
	color1 = [4]float32{1, 1, 1, 1}
	color2 = [4]float32{0, 0, 0, 1}
}

func Loop() {
	imgui.Begin("Color Palette")
	if color1Idx == color2Idx {
		color1 = [4]float32{0, 0, 0, 1}
		color1Idx = -1
	}
	for i, color := range colors {
		idx := int32(i)
		if color1Idx == idx {
			imgui.PushStyleColorVec4(imgui.ColFrameBg, color1Mark)
		}
		if color2Idx == idx {
			imgui.PushStyleColorVec4(imgui.ColFrameBg, color2Mark)
		}
		imgui.ColorButton("color #"+strconv.FormatInt(int64(i), 10), newVec4(color))
		if color1Idx == idx || color2Idx == idx {
			imgui.PopStyleColor()
		}
		if imgui.IsItemHovered() {
			io := imgui.CurrentContext().IO()
			mouseRelease := io.MouseReleased()
			if mouseRelease[0] {
				if color2Idx == color1Idx {
					color2 = color1
					color2Idx = color1Idx
					color2Mark = color1Mark
				}
				color1 = color
				color1Idx = idx
				color1Mark = highContrast(color)
			}
			if mouseRelease[1] {
				if color2Idx == color1Idx {
					color1 = color2
					color1Idx = color2Idx
					color1Mark = color2Mark
				}
				color2 = color
				color2Idx = idx
				color2Mark = highContrast(color)
			}
		}
		imgui.SameLine()
	}
	imgui.End()
}

func SelectedColors() ([4]float32, [4]float32) {
	return color1, color2
}
