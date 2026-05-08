package palette

import (
	"VulpesEditor/app/util"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
)

type color struct {
	id    uint32
	color [4]float32
	mark  imgui.Vec4
}

var color1 *color
var color2 *color
var colors []color

func highContrast(rgba [4]float32) imgui.Vec4 {
	hsv := [4]float32{}
	imgui.ColorConvertRGBtoHSV(rgba[0], rgba[1], rgba[2], &hsv[0], &hsv[1], &hsv[2])
	if hsv[0] = hsv[0] + 0.5; hsv[0] > 1 {
		hsv[0] -= 1
	}
	hsv[1] = 1
	hsv[2] = 1
	final := [4]float32{}
	imgui.ColorConvertHSVtoRGB(hsv[0], hsv[1], hsv[2], &final[0], &final[1], &final[2])
	final[3] = 1
	return newVec4(final)
}

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

var idSys util.IdSystem

func Init() {
	idSys = util.NewIdSystem()

	var first color
	var last color
	{
		firstIdx := len(colors)
		colorCount := 10
		var step float32 = 1 / float32(colorCount)
		var value float32 = 0
		for range colorCount {
			var rgb [3]float32
			imgui.ColorConvertHSVtoRGB(1, 0, value, &rgb[0], &rgb[1], &rgb[2])
			value += step
			rgba := [4]float32{rgb[0], rgb[1], rgb[2], 1}
			colors = append(colors, color{idSys.GetID(), rgba, highContrast(rgba)})
		}
		lastIdx := len(colors) - 1
		first = colors[firstIdx]
		last = colors[lastIdx]
	}
	color1 = &first
	color2 = &last

	{
		colorCount := 20
		var step float32 = 1 / float32(colorCount)
		var hue float32 = 0
		for range colorCount {
			var rgb [3]float32
			imgui.ColorConvertHSVtoRGB(hue, 1, 1, &rgb[0], &rgb[1], &rgb[2])
			hue += step
			rgba := [4]float32{rgb[0], rgb[1], rgb[2], 1}
			colors = append(colors, color{idSys.GetID(), rgba, highContrast(rgba)})
		}
	}
}

func Loop() {
	imgui.Begin("Color Palette")
	for i, color := range colors {
		id := color.id
		if color1.id == id || color2.id == id {
			imgui.PushStyleColorVec4(imgui.ColFrameBg, color.mark)
		}
		imgui.ColorButton("color #"+strconv.FormatInt(int64(id), 10), newVec4(color.color))
		if color1.id == id || color2.id == id {
			imgui.PopStyleColor()
		}
		if imgui.IsItemHovered() {
			io := imgui.CurrentContext().IO()
			mouseRelease := io.MouseReleased()
			if mouseRelease[0] {
				color1 = &color
			}
			if mouseRelease[1] {
				color2 = &color
			}
		}
		if i != len(colors)-1 {
			imgui.SameLine()
		}
	}
	imgui.End()
}

func SelectedColors() ([4]float32, [4]float32) {
	return color1.color, color2.color
}
