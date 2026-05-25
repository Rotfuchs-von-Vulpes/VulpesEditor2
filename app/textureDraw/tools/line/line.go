package line

type Line struct {
}

var initPos [2]int32
var endPos [2]int32
var painted [][2]int32

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

func (_ Line) SendTexture(colors [][4]float32, w, h uint32) {

}

func (_ Line) ButtonPress(pos [2]int32) {
	initPos = pos
}

func (_ Line) ButtonRelease(pos [2]int32) {
}

func (_ Line) Move(pos1, pos2 [2]int32) {
	endPos = pos2
	painted = make([][2]int32, 0)
	painted = line(initPos, endPos)
}

func (_ Line) Visualize() (toVisualize [][2]int32) {
	return painted
}

func (_ Line) Change() (toChange [][2]int32) {
	toChange = painted
	painted = make([][2]int32, 0)
	return
}

func (_ Line) Reset() {
	initPos = [2]int32{}
	endPos = [2]int32{}
	painted = make([][2]int32, 0)
}
