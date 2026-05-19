package pencil

type Pencil struct{}

var painting bool
var painted [][2]int32

func (s Pencil) SendTexture(colors [][4]float32, width, height uint32) {

}

func (s Pencil) ButtonPress(pos [2]int32) {
	painting = true
}

func (s Pencil) ButtonRelease(pos [2]int32) {
	painting = false
}

func (s Pencil) Move(pos1, pos2 [2]int32) {
	if painting {
		painted = append(painted, pos2)
	}
}

func (s Pencil) Visualize() (toPaint [][2]int32) {
	return painted
}

func (s Pencil) Change() (toChange [][2]int32) {
	toChange = make([][2]int32, len(painted))
	copy(toChange, painted)
	painted = make([][2]int32, 0)
	return
}

func (s Pencil) Reset() {
	painting = false
	painted = make([][2]int32, 0)
}
