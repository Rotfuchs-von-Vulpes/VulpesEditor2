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
	points   [][5]float32
	vertices []float32
	indices  []uint32
}

func (s *buffer) addFace(face [4][5]float32) {
	count := length(s.points)
	s.points = append(s.points, face[0])
	s.points = append(s.points, face[1])
	s.points = append(s.points, face[2])
	s.points = append(s.points, face[3])
	s.indices = append(s.indices, count, count+1, count+2)
	s.indices = append(s.indices, count+2, count+3, count)
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
	b.addFace(upFace)
	b.addFace(frontFace)
	b.addFace(rightFace)
	b.addFace(backFace)
	b.addFace(leftFace)
	b.addFace(downFace)
	for _, point := range b.points {
		b.vertices = append(b.vertices, point[0]*s.size[0]/2+s.pos[0])
		b.vertices = append(b.vertices, point[1]*s.size[1]/2+s.pos[1])
		b.vertices = append(b.vertices, point[2]*s.size[2]/2+s.pos[2])
		b.vertices = append(b.vertices, point[3])
		b.vertices = append(b.vertices, point[4])
	}
}

type Model struct {
	units []*cubeUnit
}

func (s *Model) AddUnit(unit *cubeUnit) {
	s.units = append(s.units, unit)
}

func (s Model) ToBuffer() (vertices []float32, indexBuffer []uint32) {
	b := new(buffer)
	for _, unit := range s.units {
		unit.toBuffer(b)
	}
	return b.vertices, b.indices
}

func NewModel() (s *Model) {
	s = new(Model)
	s.AddUnit(NewUnit([3]float32{0, 0, 0}, [3]float32{1, 1, 1}))
	return
}
