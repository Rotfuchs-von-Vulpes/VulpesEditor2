package bucket

import (
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/texture/image"
	"slices"
)

type Bucket struct {
}

var Texture *image.Texture
var height, width uint32
var painted []image.PixelEdit
var drawingColor [4]float32

func (_ Bucket) ButtonPress(pos [2]int32, secondButton bool) {
	c1, c2 := color.GetColors()
	if !secondButton {
		drawingColor = c1
	} else {
		drawingColor = c2
	}
	ok, targetColor := Texture.Get(pos)
	if !ok {
		return
	}
	toFill := [][2]int32{pos}

	canSpread := func(newPos [2]int32) bool {
		ok, c := Texture.Get(newPos)
		if !ok {
			return false
		}
		if c == targetColor && !slices.Contains(toFill, newPos) {
			return true
		}
		return false
	}

	idx := 0
	for {
		if idx > len(toFill)-1 {
			break
		}
		pos := toFill[idx]
		idx += 1
		if canSpread([2]int32{pos[0] - 1, pos[1]}) {
			toFill = append(toFill, [2]int32{pos[0] - 1, pos[1]})
		}
		if canSpread([2]int32{pos[0] + 1, pos[1]}) {
			toFill = append(toFill, [2]int32{pos[0] + 1, pos[1]})
		}
		if canSpread([2]int32{pos[0], pos[1] - 1}) {
			toFill = append(toFill, [2]int32{pos[0], pos[1] - 1})
		}
		if canSpread([2]int32{pos[0], pos[1] + 1}) {
			toFill = append(toFill, [2]int32{pos[0], pos[1] + 1})
		}
	}
	painted = image.SetEditColor(toFill, drawingColor)
}

func (_ Bucket) ButtonRelease(pos [2]int32) {

}

func (_ Bucket) Move(pos1, pos2 [2]int32) {

}

func (_ Bucket) Visualize() []image.PixelEdit {
	return painted
}

func (_ Bucket) Change() (toChange []image.PixelEdit) {
	toChange = painted
	painted = make([]image.PixelEdit, 0)
	Texture.Clear()
	return
}

func (_ Bucket) Reset() {
	painted = make([]image.PixelEdit, 0)
	Texture.Clear()
}
