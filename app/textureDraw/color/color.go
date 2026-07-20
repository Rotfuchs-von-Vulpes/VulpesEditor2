package color

import (
	"VulpesEditor/app/textureDraw/color/palette"
	"VulpesEditor/app/textureDraw/color/picker"
)

func Setcolor1(color [4]float32) {
	palette.Reset([3]bool{true, false, false})
	currentCtx.color1 = color
}

func SetColor2(color [4]float32) {
	palette.Reset([3]bool{false, true, false})
	currentCtx.color2 = color
}

func GetColors() ([4]float32, [4]float32) {
	return currentCtx.color1, currentCtx.color2
}

func Init() {
	palette.Init()
}

func Loop() {
	change := picker.Loop(&currentCtx.color1, &currentCtx.color2)
	palette.Reset(change)
	palette.Loop(&currentCtx.color1, &currentCtx.color2)
}
