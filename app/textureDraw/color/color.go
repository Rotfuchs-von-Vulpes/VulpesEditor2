package color

import (
	"VulpesEditor/app/textureDraw/color/palette"
	"VulpesEditor/app/textureDraw/color/picker"
)

var color1, color2 [4]float32

func Init() {
	palette.Init(&color1, &color2)
}

func Setcolor1(color [4]float32) {
	palette.Reset([3]bool{true, false, false})
	color1 = color
}

func SetColor2(color [4]float32) {
	palette.Reset([3]bool{false, true, false})
	color2 = color
}

func GetColors() ([4]float32, [4]float32) {
	return color1, color2
}

func Loop() {
	change := picker.Loop(&color1, &color2)
	palette.Reset(change)
	palette.Loop(&color1, &color2)
}
