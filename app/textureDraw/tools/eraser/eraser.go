package eraser

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
)

type Eraser struct{}

var erasing bool
var erased []texture.PixelEdit
var trasnparent = [4]float32{0, 0, 0, 0}

func (_ Eraser) ButtonPress(pos [2]int32, secondButton bool) {
	erased = append(erased, texture.PixelEdit{Pos: pos, Color: trasnparent})
	erasing = true
}

func (_ Eraser) ButtonRelease(pos [2]int32) {
	erasing = false
}

func (_ Eraser) Move(pos1, pos2 [2]int32) {
	if erasing {
		erased = append(erased, texture.PixelEdit{Pos: pos2, Color: trasnparent})
	}
}

func (_ Eraser) Visualize() (toVisualize []texture.PixelEdit) {
	toVisualize = erased
	return
}

func (_ Eraser) Change() (toChange []texture.PixelEdit) {
	toChange = erased
	erased = make([]texture.PixelEdit, 0)
	return
}

func (_ Eraser) Reset() {
	erasing = false
	erased = make([]texture.PixelEdit, 0)
}
