package color

import (
	"VulpesEditor/app/textureDraw/color/palette"
	"VulpesEditor/app/textureDraw/color/picker"
)

func Setcolor1(color [4]float32) {
	palette.Reset([3]bool{true, false, false})
	ctx.color1 = color
}

func SetColor2(color [4]float32) {
	palette.Reset([3]bool{false, true, false})
	ctx.color2 = color
}

func GetColors() ([4]float32, [4]float32) {
	return ctx.color1, ctx.color2
}

func Init() {
	palette.Init()
}

func Loop() {
	change := picker.Loop(&ctx.color1, &ctx.color2)
	palette.Reset(change)
	palette.Loop(&ctx.color1, &ctx.color2)
}
