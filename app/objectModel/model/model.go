package model

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/util"
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
)

type CubeUnit struct {
	Id      int32
	Parent  *Model
	pos     [3]float32
	size    [3]float32
	uvs     [6][2][2]float32
	Changed bool
	QuadPos [2]float32
}

func (s *CubeUnit) SetUV(faceCount int, init, end [2]float32) {
	if faceCount >= 6 {
		return
	}
	s.uvs[faceCount][0] = init
	s.uvs[faceCount][1] = end
	if s.Parent != nil {
		s.Parent.Changed = true
	}
	if faceCount == 2 {
		s.QuadPos = init
	}
}

func (s *CubeUnit) SetUVs(uvs [6][2][2]float32) {
	s.uvs = uvs
	s.QuadPos = uvs[2][0]
	if s.Parent != nil {
		s.Parent.Changed = true
	}
}

func (s *CubeUnit) Data() ([3]float32, [3]float32) {
	return s.pos, s.size
}

func (s *CubeUnit) Edit(pos, size [3]float32) {
	s.pos = pos
	s.size = size
	s.Changed = true
	if s.Parent != nil {
		s.Parent.Changed = true
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
	Units   []*CubeUnit
	Changed bool
	idSys   *util.IdSystem
}

func NewModel() (s *Model) {
	s = new(Model)
	s.idSys = util.NewIdSystem()
	return
}

func (s *Model) AddUnit(unit *CubeUnit) (err error) {
	if unit.Parent != nil {
		err = fmt.Errorf("This unit already have parent")
		return
	}
	unit.Parent = s
	unit.Id = s.idSys.GetID()
	s.Units = append(s.Units, unit)
	s.Changed = true
	return
}

func (s *Model) appendUnit(unit *CubeUnit) {
	if unit.Parent != nil {
		unit.Parent.Remove(unit)
	}
	unit.Parent = s
	s.Units = append(s.Units, unit)
	s.Changed = true
}

func (s *Model) Remove(unit *CubeUnit) (err error) {
	if s != unit.Parent {
		err = fmt.Errorf("This unit does not belongs here")
		return
	}
	idx := -1
	for i, u := range s.Units {
		if u.Id == unit.Id {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.Units = slices.Delete(s.Units, idx, idx+1)
		unit.Parent = nil
		unit.Id = 0
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
	Id       int32            `json:"id"`
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
	for _, c := range cubes {
		u := NewUnit(c.Position, c.Size)
		u.Id = c.Id
		u.SetUVs(c.Uv)
		m.appendUnit(u)
	}
	return
}
