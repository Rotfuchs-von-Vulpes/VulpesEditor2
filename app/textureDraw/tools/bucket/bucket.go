package bucket

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/color"
	"slices"
)

type Bucket struct {
}

var Texture *texture.Texture
var height, width uint32
var painted []texture.PixelEdit
var drawingColor [4]float32

func (_ Bucket) ButtonPress(pos [2]int32, firstButton bool) {
	c1, c2 := color.GetColors()
	if firstButton {
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
	painted = texture.SetEditColor(toFill, drawingColor)
}

func (_ Bucket) ButtonRelease(pos [2]int32) {

}

func (_ Bucket) Move(pos1, pos2 [2]int32) {

}

func (_ Bucket) Visualize() []texture.PixelEdit {
	return painted
}

func (_ Bucket) Change() (toChange []texture.PixelEdit) {
	toChange = painted
	painted = make([]texture.PixelEdit, 0)
	Texture.Clear()
	return
}

func (_ Bucket) Reset() {
	painted = make([]texture.PixelEdit, 0)
	Texture.Clear()
}
