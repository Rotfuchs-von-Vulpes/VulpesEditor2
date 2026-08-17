package model

type cubeUnit struct {
	pos  [3]float32
	size [3]float32
}

func NewUnit(pos, size [3]float32) (u *cubeUnit) {
	u = new(cubeUnit)
	u.pos = pos
	u.size = size
	return
}

func length[T any](b []T) uint32 {
	return uint32(len(b))
}

type buffer struct {
	count    uint32
	vertices []float32
	indices  []uint32
}

func (s *buffer) addFace(face [4][5]float32, size, pos [3]float32) {
	points := [][5]float32{}
	points = append(points, face[0])
	points = append(points, face[1])
	points = append(points, face[2])
	points = append(points, face[3])
	s.indices = append(s.indices, s.count, s.count+1, s.count+2)
	s.indices = append(s.indices, s.count+2, s.count+3, s.count)
	for _, point := range points {
		s.vertices = append(s.vertices, point[0]*size[0]/2+pos[0])
		s.vertices = append(s.vertices, point[1]*size[1]/2+pos[1])
		s.vertices = append(s.vertices, point[2]*size[2]/2+pos[2])
		s.vertices = append(s.vertices, point[3])
		s.vertices = append(s.vertices, point[4])
	}
	s.count += 4
}

func (s cubeUnit) toBuffer(b *buffer) {
	// pos X, pos Y, pos Z, UV x, UV y
	upFace := [4][5]float32{
		{-1, 1, -1, 0, 0},
		{-1, 1, 1, 0, 1},
		{1, 1, 1, 1, 1},
		{1, 1, -1, 1, 0},
	}
	downFace := [4][5]float32{
		{-1, -1, 1, 1, 0},
		{-1, -1, -1, 1, 1},
		{1, -1, -1, 0, 1},
		{1, -1, 1, 0, 0},
	}
	frontFace := [4][5]float32{
		{1, -1, 1, 1, 1},
		{1, 1, 1, 1, 0},
		{-1, 1, 1, 0, 0},
		{-1, -1, 1, 0, 1},
	}
	backFace := [4][5]float32{
		{-1, -1, -1, 1, 1},
		{-1, 1, -1, 1, 0},
		{1, 1, -1, 0, 0},
		{1, -1, -1, 0, 1},
	}
	rightFace := [4][5]float32{
		{1, -1, -1, 1, 1},
		{1, 1, -1, 1, 0},
		{1, 1, 1, 0, 0},
		{1, -1, 1, 0, 1},
	}
	leftFace := [4][5]float32{
		{-1, -1, 1, 1, 1},
		{-1, 1, 1, 1, 0},
		{-1, 1, -1, 0, 0},
		{-1, -1, -1, 0, 1},
	}
	b.addFace(upFace, s.size, s.pos)
	b.addFace(frontFace, s.size, s.pos)
	b.addFace(rightFace, s.size, s.pos)
	b.addFace(backFace, s.size, s.pos)
	b.addFace(leftFace, s.size, s.pos)
	b.addFace(downFace, s.size, s.pos)
}

type Model struct {
	units   []*cubeUnit
	changed bool
}

func (s *Model) AddUnit(unit *cubeUnit) {
	s.units = append(s.units, unit)
	s.changed = true
}

func (s Model) ToBuffer() (vertices []float32, indexBuffer []uint32) {
	b := new(buffer)
	for _, unit := range s.units {
		unit.toBuffer(b)
	}

	return b.vertices, b.indices
}

func (s *Model) Changed() (r bool) {
	r = s.changed
	s.changed = false
	return
}

func NewModel() (s *Model) {
	s = new(Model)
	s.AddUnit(NewUnit([3]float32{0, 0, 0}, [3]float32{1, 1, 1}))
	return
}
