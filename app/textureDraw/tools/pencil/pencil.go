package pencil

type Pencil struct{}

var painting bool
var painted [][2]int32

func (_ Pencil) SendTexture(colors [][4]float32, width, height uint32) {

}

func (_ Pencil) ButtonPress(pos [2]int32) {
	painted = append(painted, pos)
	painting = true
}

func (_ Pencil) ButtonRelease(pos [2]int32) {
	painting = false
}

func (_ Pencil) Move(pos1, pos2 [2]int32) {
	if painting {
		painted = append(painted, pos2)
	}
}

func (_ Pencil) Visualize() (toVisualize [][2]int32) {
	toVisualize = painted
	return
}

func (_ Pencil) Change() (toChange [][2]int32) {
	toChange = painted
	painted = make([][2]int32, 0)
	return
}

func (_ Pencil) Reset() {
	painting = false
	painted = make([][2]int32, 0)
}
