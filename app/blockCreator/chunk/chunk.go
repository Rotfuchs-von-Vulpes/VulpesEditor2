package chunk

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/go-gl/mathgl/mgl32"
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

	Width  int32
	Depth  int32
	Height int32

	blocks []uint32
}

func New(width, height, depth int32) (c *Chunk) {
	c = new(Chunk)
	c.Width = width
	c.Height = height
	c.Depth = depth
	for range width * height * depth {
		c.blocks = append(c.blocks, 0)
	}
	return
}

func (c *Chunk) index(pos PositionInt) int32 {
	return pos[0] + int32(c.Height)*pos[1] + int32(c.Height*c.Depth)*pos[2]
}

func (c *Chunk) isOutside(pos PositionInt) bool {
	x := pos[0]
	y := pos[1]
	z := pos[2]
	if x >= int32(c.Width) || x < 0 {
		return true
	}
	if y >= int32(c.Height) || y < 0 {
		return true
	}
	if z >= int32(c.Depth) || z < 0 {
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
	for x := range c.Width {
		for y := range c.Height {
			for z := range c.Depth {
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

func Test() (c *Chunk) {
	c = New(16, 16, 16)
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
	return c
}

const MAX_DISTANCE = 40

func sign(n float32) int32 {
	if n < 0 {
		return -1
	} else if n == 0 {
		return 0
	} else {
		return 1
	}
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func floor(n float32) float32 {
	return float32(math.Floor(float64(n)))
}

type RaycastResult struct {
	Hit      bool
	BlockPos PositionInt
	Face     mgl32.Vec3
}

func (c *Chunk) getHovered(camPos, camDir mgl32.Vec3) (ok bool, side, blockPos PositionInt) {
	current := PositionInt{
		int32(floor(camPos[0])),
		int32(floor(camPos[1])),
		int32(floor(camPos[2])),
	}

	var stepX, stepY, stepZ int32 = 1, 1, 1
	if camDir[0] < 0 {
		stepX = -1
	}
	if camDir[1] < 0 {
		stepY = -1
	}
	if camDir[2] < 0 {
		stepZ = -1
	}

	tDeltaX := abs(1.0 / camDir[0])
	tDeltaY := abs(1.0 / camDir[1])
	tDeltaZ := abs(1.0 / camDir[2])

	var tMaxX, tMaxY, tMaxZ float32

	if camDir[0] > 0 {
		tMaxX = (floor(camPos[0]) + 1.0 - camPos[0]) * tDeltaX
	} else {
		tMaxX = (camPos[0] - floor(camPos[0])) * tDeltaX
	}

	if camDir[1] > 0 {
		tMaxY = (floor(camPos[1]) + 1.0 - camPos[1]) * tDeltaY
	} else {
		tMaxY = (camPos[1] - floor(camPos[1])) * tDeltaY
	}

	if camDir[2] > 0 {
		tMaxZ = (floor(camPos[2]) + 1.0 - camPos[2]) * tDeltaZ
	} else {
		tMaxZ = (camPos[2] - floor(camPos[2])) * tDeltaZ
	}

	for {
		if !c.isAir(current) {
			ok = true
			blockPos = current
			return
		}

		side = current
		if tMaxX < tMaxY {
			if tMaxX < tMaxZ {
				if tMaxX > MAX_DISTANCE {
					break
				}
				current[0] += stepX
				tMaxX += tDeltaX
			} else {
				if tMaxZ > MAX_DISTANCE {
					break
				}
				current[2] += stepZ
				tMaxZ += tDeltaZ
			}
		} else {
			if tMaxY < tMaxZ {
				if tMaxY > MAX_DISTANCE {
					break
				}
				current[1] += stepY
				tMaxY += tDeltaY
			} else {
				if tMaxZ > MAX_DISTANCE {
					break
				}
				current[2] += stepZ
				tMaxZ += tDeltaZ
			}
		}
	}

	return
}

func (c *Chunk) Destruct(pos, direction mgl32.Vec3) bool {
	ok, _, blockPos := c.getHovered(pos, direction)
	if ok {
		c.setBlock(blockPos, 0)
	}
	return ok
}

func (c *Chunk) Hovered(pos, direction mgl32.Vec3) (ok bool, p mgl32.Vec3) {
	ok, _, bp := c.getHovered(pos, direction)
	if ok {
		ok = true
		p = mgl32.Vec3{float32(bp[0]), float32(bp[1]), float32(bp[2])}
	}
	return ok, p
}

func (c *Chunk) Construct(pos, direction mgl32.Vec3) bool {
	ok, blockPos, _ := c.getHovered(pos, direction)
	if ok {
		c.setBlock(blockPos, 1)
	}
	return ok
}
