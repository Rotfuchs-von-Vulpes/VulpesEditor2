package bucket

import "slices"

type Bucket struct {
}

var texture [][][4]float32
var height, width uint32
var painted [][2]int32

func (_ Bucket) SendTexture(colors [][][4]float32, w, h uint32) {
	texture = colors
	width = w
	height = h
}

func (_ Bucket) ButtonPress(pos [2]int32) {
	color := texture[pos[0]][pos[1]]
	toFill := [][2]int32{pos}

	canSpread := func(newPos [2]int32) bool {
		if newPos[0] < 0 || newPos[0] >= int32(width) || newPos[1] < 0 || newPos[1] >= int32(height) {
			return false
		}
		c := texture[newPos[0]][newPos[1]]
		if c == color && !slices.Contains(toFill, newPos) {
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
	painted = toFill
}

func (_ Bucket) ButtonRelease(pos [2]int32) {

}

func (_ Bucket) Move(pos1, pos2 [2]int32) {

}

func (_ Bucket) Visualize() [][2]int32 {
	return painted
}

func (_ Bucket) Change() (toChange [][2]int32) {
	toChange = painted
	painted = make([][2]int32, 0)
	texture = make([][][4]float32, 0)
	return
}

func (_ Bucket) Reset() {
	painted = make([][2]int32, 0)
	texture = make([][][4]float32, 0)
}
