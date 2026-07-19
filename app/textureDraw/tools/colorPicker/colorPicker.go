package colorPicker

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/color"
)

type ColorPicker struct{}

var Texture *texture.Texture
var FirstButton bool

func (_ ColorPicker) ButtonPress(pos [2]int32, firstButton bool) {
	FirstButton = firstButton
}

func (_ ColorPicker) ButtonRelease(pos [2]int32) {
	ok, targetColor := Texture.Get(pos)
	if ok {
		if FirstButton {
			color.Setcolor1(targetColor)
		} else {
			color.SetColor2(targetColor)
		}
	}
}

func (_ ColorPicker) Move(pos1, pos2 [2]int32) {

}

func (_ ColorPicker) Reset() {

}
