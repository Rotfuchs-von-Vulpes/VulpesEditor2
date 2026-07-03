package line

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/color"
)

type Line struct {
}

var initPos [2]int32
var endPos [2]int32
var painted []texture.PixelEdit
var drawingColor [4]float32

func abs(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

func line(start, end [2]int32) (out [][2]int32) {
	x0, y0 := start[0], start[1]
	x1, y1 := end[0], end[1]
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := int32(1)
	if x0 >= x1 {
		sx = -1
	}
	sy := int32(1)
	if y0 >= y1 {
		sy = -1
	}
	err := dx + dy
	for {
		out = append(out, [2]int32{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
	return
}

func (_ Line) ButtonPress(pos [2]int32, secondButton bool) {
	c1, c2 := color.GetColors()
	if !secondButton {
		drawingColor = c1
	} else {
		drawingColor = c2
	}
	initPos = pos
}

func (_ Line) ButtonRelease(pos [2]int32) {
}

func (_ Line) Move(pos1, pos2 [2]int32) {
	endPos = pos2
	painted = make([]texture.PixelEdit, 0)
	painted = texture.SetEditColor(line(initPos, endPos), drawingColor)
}

func (_ Line) Visualize() (toVisualize []texture.PixelEdit) {
	return painted
}

func (_ Line) Change() (toChange []texture.PixelEdit) {
	toChange = painted
	painted = make([]texture.PixelEdit, 0)
	return
}

func (_ Line) Reset() {
	initPos = [2]int32{}
	endPos = [2]int32{}
	painted = make([]texture.PixelEdit, 0)
}
