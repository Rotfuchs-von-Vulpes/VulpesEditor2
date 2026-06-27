package rectangle

import (
	"VulpesEditor/app/textureDraw/texture/image"
)

type Rectangle struct{}

var initPos [2]int32
var endPos [2]int32
var color [4]float32
var painted = []image.PixelEdit{}

func abs(n int32) int32 {
	if n >= 0 {
		return n
	} else {
		return -n
	}
}

func line(init, step [2]int32, length int32) (out [][2]int32) {
	pos := init
	for i := 0; i < int(length); i++ {
		pos[0] += step[0]
		pos[1] += step[1]
		out = append(out, pos)
	}
	return
}

func rectangle(init, end [2]int32) (out [][2]int32) {
	top := max(init[1], end[1])
	down := min(init[1], end[1])
	left := min(init[0], end[0])
	right := max(init[0], end[0])
	out = append(out, line([2]int32{left, top}, [2]int32{1, 0}, right-left)...)
	out = append(out, line([2]int32{right, top}, [2]int32{0, -1}, top-down)...)
	out = append(out, line([2]int32{right, down}, [2]int32{-1, 0}, right-left)...)
	out = append(out, line([2]int32{left, down}, [2]int32{0, 1}, top-down)...)
	return
}

func (_ Rectangle) SendTexture(colors [][4]float32, width, height uint32) {

}

func (_ Rectangle) ButtonPress(pos [2]int32, col [4]float32) {
	color = col
	initPos = pos
}

func (_ Rectangle) ButtonRelease(pos [2]int32) {

}

func (_ Rectangle) Move(pos1, pos2 [2]int32) {
	endPos = pos2

	painted = make([]image.PixelEdit, 0)
	painted = image.SetEditColor(rectangle(initPos, endPos), color)
}

func (_ Rectangle) Visualize() (toVisualize []image.PixelEdit) {
	toVisualize = painted
	return
}

func (_ Rectangle) Change() (toChange []image.PixelEdit) {
	toChange = painted
	painted = make([]image.PixelEdit, 0)
	return
}

func (_ Rectangle) Reset() {
	initPos = [2]int32{}
	endPos = [2]int32{}
	painted = make([]image.PixelEdit, 0)
}
