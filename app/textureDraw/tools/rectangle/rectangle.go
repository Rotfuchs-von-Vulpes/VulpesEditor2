package rectangle

type Rectangle struct{}

var initPos [2]int32
var endPos [2]int32
var painted [][2]int32

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

func (_ Rectangle) SendTexture(colors [][][4]float32, width, height uint32) {

}

func (_ Rectangle) ButtonPress(pos [2]int32) {
	initPos = pos
}

func (_ Rectangle) ButtonRelease(pos [2]int32) {

}

func (_ Rectangle) Move(pos1, pos2 [2]int32) {
	endPos = pos2
	painted = rectangle(initPos, endPos)
}

func (_ Rectangle) Visualize() (toVisualize [][2]int32) {
	toVisualize = painted
	return
}

func (_ Rectangle) Change() (toChange [][2]int32) {
	toChange = painted
	painted = make([][2]int32, 0)
	return
}

func (_ Rectangle) Reset() {
	initPos = [2]int32{}
	endPos = [2]int32{}
	painted = make([][2]int32, 0)
}
