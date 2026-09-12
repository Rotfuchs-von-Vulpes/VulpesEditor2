package model

import (
	"VulpesEditor/app/file"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
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
	rot      [3]float32
	uvs      [6][2][2]float32
	Changed  bool
	Resized  bool
	QuadPos  [2]float32
	boneID   uint32
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
				s.Source.AppendUnit(v)
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

func (s *CubeUnit) Data() ([3]float32, [3]float32, [3]float32) {
	return s.pos, s.size, s.rot
}

func (s *CubeUnit) Edit(pos, size, rot [3]float32) {
	if s.size != size {
		s.Resized = true
	}
	s.pos = pos
	s.size = size
	s.rot = rot
	s.Changed = true
	if s.Source != nil {
		s.Source.Changed = true
	}
}

func (s CubeUnit) absolutePos() (pos [3]float32) {
	if s.parent != nil {
		pos = s.parent.absolutePos()
	}
	pos[0] += s.pos[0]
	pos[1] += s.pos[1]
	pos[2] += s.pos[2]
	return
}

func (s CubeUnit) toBuffer(b *buffer) {
	boneId := math.Float32frombits(s.boneID)
	// pos X, pos Y, pos Z, UV x, UV y
	upFace := [4][6]float32{
		{-1, 1, -1, s.uvs[0][0][0], s.uvs[0][0][1], boneId},
		{-1, 1, 1, s.uvs[0][0][0], s.uvs[0][1][1], boneId},
		{1, 1, 1, s.uvs[0][1][0], s.uvs[0][1][1], boneId},
		{1, 1, -1, s.uvs[0][1][0], s.uvs[0][0][1], boneId},
	}
	downFace := [4][6]float32{
		{-1, -1, 1, s.uvs[1][1][0], s.uvs[1][0][1], boneId},
		{-1, -1, -1, s.uvs[1][1][0], s.uvs[1][1][1], boneId},
		{1, -1, -1, s.uvs[1][0][0], s.uvs[1][1][1], boneId},
		{1, -1, 1, s.uvs[1][0][0], s.uvs[1][0][1], boneId},
	}
	frontFace := [4][6]float32{
		{1, -1, 1, s.uvs[2][1][0], s.uvs[2][1][1], boneId},
		{1, 1, 1, s.uvs[2][1][0], s.uvs[2][0][1], boneId},
		{-1, 1, 1, s.uvs[2][0][0], s.uvs[2][0][1], boneId},
		{-1, -1, 1, s.uvs[2][0][0], s.uvs[2][1][1], boneId},
	}
	backFace := [4][6]float32{
		{-1, -1, -1, s.uvs[3][1][0], s.uvs[3][1][1], boneId},
		{-1, 1, -1, s.uvs[3][1][0], s.uvs[3][0][1], boneId},
		{1, 1, -1, s.uvs[3][0][0], s.uvs[3][0][1], boneId},
		{1, -1, -1, s.uvs[3][0][0], s.uvs[3][1][1], boneId},
	}
	rightFace := [4][6]float32{
		{1, -1, -1, s.uvs[4][1][0], s.uvs[4][1][1], boneId},
		{1, 1, -1, s.uvs[4][1][0], s.uvs[4][0][1], boneId},
		{1, 1, 1, s.uvs[4][0][0], s.uvs[4][0][1], boneId},
		{1, -1, 1, s.uvs[4][0][0], s.uvs[4][1][1], boneId},
	}
	leftFace := [4][6]float32{
		{-1, -1, 1, s.uvs[5][1][0], s.uvs[5][1][1], boneId},
		{-1, 1, 1, s.uvs[5][1][0], s.uvs[5][0][1], boneId},
		{-1, 1, -1, s.uvs[5][0][0], s.uvs[5][0][1], boneId},
		{-1, -1, -1, s.uvs[5][0][0], s.uvs[5][1][1], boneId},
	}
	pos := s.absolutePos()
	if s.size[0] != 0 && s.size[2] != 0 {
		b.addFace(upFace, s.size, pos, s.rot)
		b.addFace(downFace, s.size, pos, s.rot)
	}
	if s.size[0] != 0 && s.size[1] != 0 {
		b.addFace(frontFace, s.size, pos, s.rot)
		b.addFace(backFace, s.size, pos, s.rot)
	}
	if s.size[1] != 0 && s.size[2] != 0 {
		b.addFace(rightFace, s.size, pos, s.rot)
		b.addFace(leftFace, s.size, pos, s.rot)
	}
}

func NewUnit(pos, size, rot [3]float32) (u *CubeUnit) {
	u = new(CubeUnit)
	u.Id = uuid.New().String()
	u.pos = pos
	u.size = size
	u.rot = rot
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

func rotatePointX(p1 [3]float64, angle float64) (p2 [3]float64) {
	p2[0] = p1[0]
	p2[1] = p1[1]*math.Cos(angle) - p1[2]*math.Sin(angle)
	p2[2] = p1[1]*math.Sin(angle) + p1[2]*math.Cos(angle)
	return
}

func rotatePointY(p1 [3]float64, angle float64) (p2 [3]float64) {
	p2[0] = p1[0]*math.Cos(angle) + p1[2]*math.Sin(angle)
	p2[1] = p1[1]
	p2[2] = -p1[0]*math.Sin(angle) + p1[2]*math.Cos(angle)
	return
}

func rotatePointZ(p1 [3]float64, angle float64) (p2 [3]float64) {
	p2[0] = p1[0]*math.Cos(angle) - p1[1]*math.Sin(angle)
	p2[1] = p1[0]*math.Sin(angle) + p1[1]*math.Cos(angle)
	p2[2] = p1[2]
	return
}

func toVec(x, y, z float32) (p1 [3]float64) {
	p1[0] = float64(x)
	p1[1] = float64(y)
	p1[2] = float64(z)
	return
}

func fromVec(p1 [3]float64) (x, y, z float32) {
	x = float32(p1[0])
	y = float32(p1[1])
	z = float32(p1[2])
	return
}

func (s *buffer) addFace(face [4][6]float32, size, pos, rot [3]float32) {
	points := [][6]float32{}
	points = append(points, face[0])
	points = append(points, face[1])
	points = append(points, face[2])
	points = append(points, face[3])
	s.indices = append(s.indices, s.count, s.count+1, s.count+2)
	s.indices = append(s.indices, s.count+2, s.count+3, s.count)
	for _, point := range points {
		angleX := float64(rot[0])
		angleY := float64(rot[1])
		angleZ := float64(rot[2])
		vec := toVec(point[0]*size[0]/2, point[1]*size[1]/2, point[2]*size[2]/2)
		vec = rotatePointX(vec, angleX)
		vec = rotatePointY(vec, angleY)
		vec = rotatePointZ(vec, angleZ)
		point[0], point[1], point[2] = fromVec(vec)
		s.vertices = append(s.vertices, point[0]+pos[0])
		s.vertices = append(s.vertices, point[1]+pos[1])
		s.vertices = append(s.vertices, point[2]+pos[2])
		s.vertices = append(s.vertices, point[3])
		s.vertices = append(s.vertices, point[4])
		s.vertices = append(s.vertices, point[5])
	}
	s.count += 4
}

type Model struct {
	Id      string
	Units   []*CubeUnit
	Changed bool
	count   uint32
}

func NewModel() (s *Model) {
	s = new(Model)
	s.Id = uuid.New().String()
	return
}

func (s *Model) AppendUnit(unit *CubeUnit) bool {
	if unit.Source != nil {
		unit.Source.Remove(unit)
	}
	unit.Source = s
	s.Units = append(s.Units, unit)
	s.Changed = true
	unit.parent = nil
	unit.boneID = s.count
	if unit.Name == "" {
		unit.Name = fmt.Sprintf("Unit #%d", s.count)
	}
	s.count++
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
		u.Resized = false
	}
	s.Changed = false
}

func (s *Model) GetId() string {
	return s.Id
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
		return s.AppendUnit(v)
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

type cube struct {
	Id       string           `json:"id"`
	Position [3]float32       `json:"position"`
	Size     [3]float32       `json:"size"`
	Rotation [3]float32       `json:"rotation"`
	Uv       [6][2][2]float32 `json:"uv"`
}

func (s *Model) saveModel(w *file.ArchiveWriter) (err error) {
	cubes := []cube{}
	for _, u := range s.Units {
		c := cube{}
		c.Id = u.Id
		c.Position = u.pos
		c.Size = u.size
		c.Rotation = u.rot
		c.Uv = u.uvs
		cubes = append(cubes, c)
	}
	buff := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cubes)
	if err != nil {
		return
	}
	w.Write("model.json", buff.Bytes())
	return
}

func (s *Model) saveHierarchy(w *file.ArchiveWriter) (err error) {
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
	encoder := csv.NewWriter(buff)
	if err := encoder.Write([]string{"unit", "parent"}); err != nil {
		return err
	}
	for _, u := range s.Units {
		if u.parent == nil {
			err = encoder.Write([]string{u.Id, s.Id})
		} else {
			err = encoder.Write([]string{u.Id, u.parent.Id})
		}
		if err != nil {
			return
		}
	}
	encoder.Flush()
	w.Write("hierarchy.csv", buff.Bytes())
	return
}

func (s *Model) Save(w *file.ArchiveWriter) (err error) {
	if err := s.saveModel(w); err != nil {
		return err
	}
	if err := s.saveHierarchy(w); err != nil {
		return err
	}
	return
}

func (m *Model) readModel(r *file.ArchiveReader) (err error) {
	f, err := r.Open("model.json")
	if err != nil {
		return
	}
	defer f.Close()
	var cubes []cube
	err = json.NewDecoder(f).Decode(&cubes)
	if err != nil {
		return
	}
	for _, c := range cubes {
		u := NewUnit(c.Position, c.Size, c.Rotation)
		u.Id = c.Id
		u.SetUVs(c.Uv)
		m.AppendUnit(u)
	}
	return
}

func (m *Model) readHierarchy(r *file.ArchiveReader) (err error) {
	f, err := r.Open("hierarchy.csv")
	if err != nil {
		return
	}
	defer f.Close()
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return
	}
	for i, line := range data {
		if i == 0 {
			continue
		}
		if len(line) != 2 {
			err = fmt.Errorf("Wrong Number of columns in line %d.", i)
			return
		}
		childId := line[0]
		parentId := line[1]
		var child *CubeUnit
		var parent *CubeUnit
		for _, u := range m.Units {
			if u.Id == childId {
				child = u
				break
			}
		}
		if child == nil {
			err = fmt.Errorf("Unit %s does not exist.", childId)
			return
		}
		if m.Id == parentId {
			continue
		}
		for _, u := range m.Units {
			if u.Id == parentId {
				parent = u
				break
			}
		}
		if parent == nil {
			err = fmt.Errorf("Unit %s does not exist.", childId)
			return
		}
		parent.Append(child)
	}
	return
}

func OpenModel(r *file.ArchiveReader, id string) (m *Model, err error) {
	m = NewModel()
	m.Id = id
	if err := m.readModel(r); err != nil {
		return nil, err
	}
	if err := m.readHierarchy(r); err != nil {
		return nil, err
	}
	return
}
