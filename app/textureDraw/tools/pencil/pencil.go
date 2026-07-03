package pencil

import (
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/texture/image"
)

type Pencil struct{}

var painting bool
var painted []image.PixelEdit
var drawingColor [4]float32

func (_ Pencil) ButtonPress(pos [2]int32, secondButton bool) {
	c1, c2 := color.GetColors()
	if !secondButton {
		drawingColor = c1
	} else {
		drawingColor = c2
	}
	painted = append(painted, image.PixelEdit{Pos: pos, Color: drawingColor})
	painting = true
}

func (_ Pencil) ButtonRelease(pos [2]int32) {
	painting = false
}

func (_ Pencil) Move(pos1, pos2 [2]int32) {
	if painting {
		painted = append(painted, image.PixelEdit{Pos: pos2, Color: drawingColor})
	}
}

func (_ Pencil) Visualize() (toVisualize []image.PixelEdit) {
	toVisualize = painted
	return
}

func (_ Pencil) Change() (toChange []image.PixelEdit) {
	toChange = painted
	painted = make([]image.PixelEdit, 0)
	return
}

func (_ Pencil) Reset() {
	painting = false
	painted = make([]image.PixelEdit, 0)
}
