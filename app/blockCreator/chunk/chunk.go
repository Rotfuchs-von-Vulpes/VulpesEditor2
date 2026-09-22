package chunk

import (
	"bytes"
	"encoding/binary"
)

type PositionInt [3]int32

func (p PositionInt) Add(p2 PositionInt) (p3 PositionInt) {
	p3[0] = p[0] + p2[0]
	p3[1] = p[1] + p2[1]
	p3[2] = p[2] + p2[2]
	return
}

type PositionFloat [3]float32

func (p PositionFloat) Add(p2 PositionFloat) (p3 PositionFloat) {
	p3[0] = p[0] + p2[0]
	p3[1] = p[1] + p2[1]
	p3[2] = p[2] + p2[2]
	return
}

type block struct {
	kindId uint32
}

type Chunk struct {
	id string

	width  int32
	depth  int32
	height int32

	blocks []uint32
}

func New(width, height, depth int32) (c *Chunk) {
	c = new(Chunk)
	c.width = width
	c.height = height
	c.depth = depth
	for range width * height * depth {
		c.blocks = append(c.blocks, 0)
	}
	return
}

func (c *Chunk) index(pos PositionInt) int32 {
	return pos[0] + int32(c.height)*pos[1] + int32(c.height*c.depth)*pos[2]
}

func (c *Chunk) isOutside(pos PositionInt) bool {
	x := pos[0]
	y := pos[1]
	z := pos[2]
	if x >= int32(c.width) || x < 0 {
		return true
	}
	if y >= int32(c.height) || y < 0 {
		return true
	}
	if z >= int32(c.depth) || z < 0 {
		return true
	}
	return false
}

func (c *Chunk) setBlock(pos PositionInt, blockID uint32) {
	if c.isOutside(pos) {
		return
	}
	c.blocks[c.index(pos)] = blockID
}

func (c *Chunk) getBlock(pos PositionInt) uint32 {
	if c.isOutside(pos) {
		return 0
	}
	return c.blocks[c.index(pos)]
}

func (c *Chunk) isAir(pos PositionInt) bool {
	b := c.getBlock(pos)
	return b == 0
}

type Side uint32

const (
	sideTop Side = iota
	sideBottom
	sideEast
	sideWest
	sideNorth
	sideSouth
)

var TOP PositionInt = [3]int32{0, 1, 0}
var BOTTOM PositionInt = [3]int32{0, -1, 0}
var EAST PositionInt = [3]int32{1, 0, 0}
var WEST PositionInt = [3]int32{-1, 0, 0}
var NORTH PositionInt = [3]int32{0, 0, 1}
var SOUTH PositionInt = [3]int32{0, 0, -1}

type faceUnit struct {
	id        uint32
	direction Side
	pos       PositionFloat
	uvs       [4][2]float32
}

func (f *faceUnit) getVertices() [4]PositionFloat {
	var p1, p2, p3, p4 PositionFloat
	switch f.direction {
	case sideTop:
		p1 = f.pos.Add(PositionFloat{0, 0, 1})
		p2 = f.pos.Add(PositionFloat{1, 0, 1})
		p3 = f.pos.Add(PositionFloat{1, 0, 0})
		p4 = f.pos
	case sideBottom:
		p1 = f.pos
		p2 = f.pos.Add(PositionFloat{1, 0, 0})
		p3 = f.pos.Add(PositionFloat{1, 0, 1})
		p4 = f.pos.Add(PositionFloat{0, 0, 1})
	case sideEast:
		p1 = f.pos.Add(PositionFloat{0, 1, 0})
		p2 = f.pos.Add(PositionFloat{0, 1, 1})
		p3 = f.pos.Add(PositionFloat{0, 0, 1})
		p4 = f.pos
	case sideWest:
		p1 = f.pos.Add(PositionFloat{0, 1, 0})
		p2 = f.pos
		p3 = f.pos.Add(PositionFloat{0, 0, 1})
		p4 = f.pos.Add(PositionFloat{0, 1, 1})
	case sideNorth:
		p1 = f.pos.Add(PositionFloat{0, 1, 0})
		p2 = f.pos
		p3 = f.pos.Add(PositionFloat{1, 0, 0})
		p4 = f.pos.Add(PositionFloat{1, 1, 0})
	case sideSouth:
		p1 = f.pos.Add(PositionFloat{0, 1, 0})
		p2 = f.pos.Add(PositionFloat{1, 1, 0})
		p3 = f.pos.Add(PositionFloat{1, 0, 0})
		p4 = f.pos
	}
	return [4]PositionFloat{p1, p2, p3, p4}
}

type cubeUnit struct {
	id  uint32
	pos PositionInt
}

func (c *Chunk) isSIdeExposed(pos PositionInt, offset PositionInt) bool {
	return c.isAir(pos.Add(offset))
}

func NewFace(pos PositionInt, id uint32, side Side) (face faceUnit) {
	face.id = id
	face.direction = side
	switch side {
	case sideTop:
		face.pos[0] = float32(pos[0])
		face.pos[1] = float32(pos[1] + 1)
		face.pos[2] = float32(pos[2])
	case sideBottom:
		face.pos[0] = float32(pos[0])
		face.pos[1] = float32(pos[1])
		face.pos[2] = float32(pos[2])
	case sideEast:
		face.pos[0] = float32(pos[0] + 1)
		face.pos[1] = float32(pos[1])
		face.pos[2] = float32(pos[2])
	case sideWest:
		face.pos[0] = float32(pos[0])
		face.pos[1] = float32(pos[1])
		face.pos[2] = float32(pos[2])
	case sideNorth:
		face.pos[0] = float32(pos[0])
		face.pos[1] = float32(pos[1])
		face.pos[2] = float32(pos[2] + 1)
	case sideSouth:
		face.pos[0] = float32(pos[0])
		face.pos[1] = float32(pos[1])
		face.pos[2] = float32(pos[2])
	}
	face.uvs[0] = [2]float32{1, 1}
	face.uvs[1] = [2]float32{1, 0}
	face.uvs[2] = [2]float32{0, 0}
	face.uvs[3] = [2]float32{0, 1}
	return
}

func (c *Chunk) hasSideExposed(cube *cubeUnit, faces *[]faceUnit) bool {
	final := false
	if c.isSIdeExposed(cube.pos, TOP) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideTop))
		final = true
	}
	if c.isSIdeExposed(cube.pos, BOTTOM) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideBottom))
		final = true
	}
	if c.isSIdeExposed(cube.pos, EAST) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideEast))
		final = true
	}
	if c.isSIdeExposed(cube.pos, WEST) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideWest))
		final = true
	}
	if c.isSIdeExposed(cube.pos, NORTH) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideNorth))
		final = true
	}
	if c.isSIdeExposed(cube.pos, SOUTH) {
		*faces = append(*faces, NewFace(cube.pos, cube.id, sideSouth))
		final = true
	}
	return final
}

func (c *Chunk) ToBuffer() ([]byte, []uint32) {
	var faces []faceUnit
	for x := range c.width {
		for y := range c.height {
			for z := range c.depth {
				var pos PositionInt = [3]int32{x, y, z}
				var cube cubeUnit
				cube.id = c.getBlock(pos)
				cube.pos = pos
				if !c.isAir(pos) {
					c.hasSideExposed(&cube, &faces)
				}
			}
		}
	}

	verticesBuffer := bytes.NewBuffer(nil)
	var indices []uint32
	var count uint32
	for _, face := range faces {
		indices = append(indices, count, count+1, count+2)
		indices = append(indices, count+2, count+3, count)
		vertices := face.getVertices()
		for i, v := range vertices {
			binary.Write(verticesBuffer, binary.LittleEndian, v[0])
			binary.Write(verticesBuffer, binary.LittleEndian, v[1])
			binary.Write(verticesBuffer, binary.LittleEndian, v[2])
			binary.Write(verticesBuffer, binary.LittleEndian, face.uvs[i][0])
			binary.Write(verticesBuffer, binary.LittleEndian, face.uvs[i][1])
			binary.Write(verticesBuffer, binary.LittleEndian, face.id)
		}
		count += 4
	}
	return verticesBuffer.Bytes(), indices
}

func Test() ([]byte, []uint32) {
	c := New(16, 16, 16)
	for x := range 16 {
		for y := range 16 {
			for z := range 16 {
				pos := PositionInt{int32(x), int32(y), int32(z)}
				x1 := x - 8
				y1 := y - 8
				z1 := z - 8
				if x1*x1+y1*y1+z1*z1 < 8*8 {
					c.setBlock(pos, 1)
				}
			}
		}
	}
	return c.ToBuffer()
}
