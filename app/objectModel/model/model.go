package model

import (
	"VulpesEditor/app/file"
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"uuid"
)

type NodeTree interface {
	GetId() string
	Parent() NodeTree
	Children() []NodeTree
	Append(node NodeTree) bool
	CanReceive(node NodeTree) bool
}

type CubeUnit struct {
	Id       string
	Name     string
	Source   *Model
	parent   *CubeUnit
	children []*CubeUnit
	pos      [3]float32
	size     [3]float32
	uvs      [6][2][2]float32
	Changed  bool
	QuadPos  [2]float32
}

func (s *CubeUnit) GetId() string {
	return s.Id
}

func (s *CubeUnit) Parent() NodeTree {
	if s.parent != nil {
		return s.parent
	} else {
		return s.Source
	}
}

func (s *CubeUnit) Children() (r []NodeTree) {
	for _, u := range s.children {
		r = append(r, u)
	}
	return
}

func (s *CubeUnit) Remove(unit *CubeUnit) {
	if unit.parent == s.parent {
		unit.parent = nil
	}
	s.children = slices.DeleteFunc(s.children, func(u *CubeUnit) bool {
		return u.Id == unit.Id
	})
}

func (s *CubeUnit) CanReceive(node NodeTree) bool {
	switch v := node.(type) {
	case (*Model):
		return false
	case (*CubeUnit):
		if v.parent == s {
			return false
		}

		p := s
		for {
			p = p.parent
			if p == nil {
				break
			}
			if p.Id == v.Id {
				return false
			}
		}

		return true
	}
	return false
}

func (s *CubeUnit) Append(node NodeTree) bool {
	switch v := node.(type) {
	case (*Model):
		return false
	case (*CubeUnit):
		if s.CanReceive(v) {
			if s.Source != nil {
				s.Source.appendUnit(v)
			}
			if v.parent != nil {
				v.parent.Remove(v)
			}
			v.parent = s
			s.children = append(s.children, v)
			return true
		}
	}
	return false
}

func (s *CubeUnit) SetUV(faceCount int, init, end [2]float32) {
	if faceCount >= 6 {
		return
	}
	s.uvs[faceCount][0] = init
	s.uvs[faceCount][1] = end
	if s.Source != nil {
		s.Source.Changed = true
	}
	if faceCount == 2 {
		s.QuadPos = init
	}
}

func (s *CubeUnit) SetUVs(uvs [6][2][2]float32) {
	s.uvs = uvs
	s.QuadPos = uvs[2][0]
	if s.Source != nil {
		s.Source.Changed = true
	}
}

func (s *CubeUnit) Data() ([3]float32, [3]float32) {
	return s.pos, s.size
}

func (s *CubeUnit) Edit(pos, size [3]float32) {
	s.pos = pos
	s.size = size
	s.Changed = true
	if s.Source != nil {
		s.Source.Changed = true
	}
}

func (s CubeUnit) toBuffer(b *buffer) {
	// pos X, pos Y, pos Z, UV x, UV y
	upFace := [4][5]float32{
		{-1, 1, -1, s.uvs[0][0][0], s.uvs[0][0][1]},
		{-1, 1, 1, s.uvs[0][0][0], s.uvs[0][1][1]},
		{1, 1, 1, s.uvs[0][1][0], s.uvs[0][1][1]},
		{1, 1, -1, s.uvs[0][1][0], s.uvs[0][0][1]},
	}
	downFace := [4][5]float32{
		{-1, -1, 1, s.uvs[1][1][0], s.uvs[1][0][1]},
		{-1, -1, -1, s.uvs[1][1][0], s.uvs[1][1][1]},
		{1, -1, -1, s.uvs[1][0][0], s.uvs[1][1][1]},
		{1, -1, 1, s.uvs[1][0][0], s.uvs[1][0][1]},
	}
	frontFace := [4][5]float32{
		{1, -1, 1, s.uvs[2][1][0], s.uvs[2][1][1]},
		{1, 1, 1, s.uvs[2][1][0], s.uvs[2][0][1]},
		{-1, 1, 1, s.uvs[2][0][0], s.uvs[2][0][1]},
		{-1, -1, 1, s.uvs[2][0][0], s.uvs[2][1][1]},
	}
	backFace := [4][5]float32{
		{-1, -1, -1, s.uvs[3][1][0], s.uvs[3][1][1]},
		{-1, 1, -1, s.uvs[3][1][0], s.uvs[3][0][1]},
		{1, 1, -1, s.uvs[3][0][0], s.uvs[3][0][1]},
		{1, -1, -1, s.uvs[3][0][0], s.uvs[3][1][1]},
	}
	rightFace := [4][5]float32{
		{1, -1, -1, s.uvs[4][1][0], s.uvs[4][1][1]},
		{1, 1, -1, s.uvs[4][1][0], s.uvs[4][0][1]},
		{1, 1, 1, s.uvs[4][0][0], s.uvs[4][0][1]},
		{1, -1, 1, s.uvs[4][0][0], s.uvs[4][1][1]},
	}
	leftFace := [4][5]float32{
		{-1, -1, 1, s.uvs[5][1][0], s.uvs[5][1][1]},
		{-1, 1, 1, s.uvs[5][1][0], s.uvs[5][0][1]},
		{-1, 1, -1, s.uvs[5][0][0], s.uvs[5][0][1]},
		{-1, -1, -1, s.uvs[5][0][0], s.uvs[5][1][1]},
	}
	if s.size[0] != 0 && s.size[2] != 0 {
		b.addFace(upFace, s.size, s.pos)
		b.addFace(downFace, s.size, s.pos)
	}
	if s.size[0] != 0 && s.size[1] != 0 {
		b.addFace(frontFace, s.size, s.pos)
		b.addFace(backFace, s.size, s.pos)
	}
	if s.size[1] != 0 && s.size[2] != 0 {
		b.addFace(rightFace, s.size, s.pos)
		b.addFace(leftFace, s.size, s.pos)
	}
}

func NewUnit(pos, size [3]float32) (u *CubeUnit) {
	u = new(CubeUnit)
	u.Id = uuid.New().String()
	u.pos = pos
	u.size = size
	for i := range u.uvs {
		u.uvs[i][0] = [2]float32{0, 0}
		u.uvs[i][1] = [2]float32{1, 1}
	}
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

type Model struct {
	id      string
	Units   []*CubeUnit
	Changed bool
}

func NewModel() (s *Model) {
	s = new(Model)
	s.id = uuid.New().String()
	return
}

func (s *Model) AddUnit(unit *CubeUnit) (err error) {
	if unit.Source != nil {
		err = fmt.Errorf("This unit already have parent")
		return
	}
	unit.Source = s
	// unit.Id = uuid.New().String()
	s.Units = append(s.Units, unit)
	s.Changed = true
	return
}

func (s *Model) appendUnit(unit *CubeUnit) bool {
	if unit.Source == nil {
		unit.Source = s
		s.Units = append(s.Units, unit)
		s.Changed = true
	}
	if unit.Source != s {
		unit.Source.Remove(unit)
		unit.Source = s
		s.Units = append(s.Units, unit)
		s.Changed = true
	}
	if unit.parent != nil {
		unit.parent.Remove(unit)
	}
	unit.parent = nil
	return true
}

func (s *Model) Remove(unit *CubeUnit) (err error) {
	if s != unit.Source {
		err = fmt.Errorf("This unit does not belongs here")
		return
	}
	idx := slices.Index(s.Units, unit)
	if idx >= 0 {
		s.Units = slices.Delete(s.Units, idx, idx+1)
		if unit.parent != nil {
			unit.parent.Remove(unit)
		}
		unit.Source = nil
		s.Changed = true
	}
	return
}

func (s Model) ToBuffer() (vertices []float32, indexBuffer []uint32) {
	b := new(buffer)
	for _, unit := range s.Units {
		unit.toBuffer(b)
	}

	return b.vertices, b.indices
}

func (s *Model) Reset() {
	for _, u := range s.Units {
		u.Changed = false
	}
	s.Changed = false
}

type cube struct {
	Id       string           `json:"id"`
	Position [3]float32       `json:"position"`
	Size     [3]float32       `json:"size"`
	Uv       [6][2][2]float32 `json:"uv"`
}

func (s *Model) Save(w *file.ArchiveWriter) {
	cubes := []cube{}
	for _, u := range s.Units {
		c := cube{}
		c.Id = u.Id
		c.Position = u.pos
		c.Size = u.size
		c.Uv = u.uvs
		cubes = append(cubes, c)
	}
	buff := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(cubes)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write("model.json", buff.Bytes())
}

func (s *Model) GetId() string {
	return s.id
}

func (s *Model) Parent() NodeTree {
	return nil
}

func (s *Model) Children() (r []NodeTree) {
	for _, u := range s.Units {
		if model, ok := u.Parent().(*Model); ok && model == s {
			r = append(r, u)
		}
	}
	return
}

func (s *Model) Append(node NodeTree) bool {
	switch v := node.(type) {
	case (*Model):
		return false
	case (*CubeUnit):
		return s.appendUnit(v)
	}
	return false
}

func (s *Model) CanReceive(node NodeTree) bool {
	switch v := node.(type) {
	case (*Model):
		return false
	case (*CubeUnit):
		if v.parent == nil && v.Source == s {
			return false
		}
		return true
	}
	return false
}

func OpenModel(r *file.ArchiveReader) (m *Model, err error) {
	f, err := r.Open("model.json")
	if err != nil {
		return nil, err
	}
	var cubes []cube
	err = json.NewDecoder(f).Decode(&cubes)
	if err != nil {
		return nil, err
	}
	f.Close()
	m = NewModel()
	for i, c := range cubes {
		u := NewUnit(c.Position, c.Size)
		u.Id = c.Id
		u.Name = fmt.Sprintf("Unit #%d", i)
		u.SetUVs(c.Uv)
		m.appendUnit(u)
	}
	return
}
