package pencil

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/color"
)

type Pencil struct{}

var painting bool
var painted []texture.PixelEdit
var drawingColor [4]float32

func (_ Pencil) ButtonPress(pos [2]int32, firstButton bool) {
	c1, c2 := color.GetColors()
	if firstButton {
		drawingColor = c1
	} else {
		drawingColor = c2
	}
	painted = append(painted, texture.PixelEdit{Pos: pos, Color: drawingColor})
	painting = true
}

func (_ Pencil) ButtonRelease(pos [2]int32) {
	painting = false
}

func (_ Pencil) Move(pos1, pos2 [2]int32) {
	if painting {
		painted = append(painted, texture.PixelEdit{Pos: pos2, Color: drawingColor})
	}
}

func (_ Pencil) Visualize() (toVisualize []texture.PixelEdit) {
	toVisualize = painted
	return
}

func (_ Pencil) Change() (toChange []texture.PixelEdit) {
	toChange = painted
	painted = make([]texture.PixelEdit, 0)
	return
}

func (_ Pencil) Reset() {
	painting = false
	painted = make([]texture.PixelEdit, 0)
}
