package texture

import (
	"slices"
)

type rect struct {
	x, y, w, h int32
}

type packet struct {
	x, y    int32
	surface *MultiQuad
	r1, r2  rect
	area    int32
}

type Packer struct {
	surfaces []*packet
	packed   []*packet
	w, h     int32
}

func toPacket(mq *MultiQuad) (pa *packet) {
	pa = new(packet)
	pa.surface = mq
	pa.r1.w = mq.width
	pa.r1.h = mq.cutHeight
	pa.r2.w = mq.cutWidth
	pa.r2.h = mq.height
	pa.area = mq.width*mq.height - (mq.width-mq.cutWidth)*(mq.height-mq.cutHeight)
	return
}

func (p *Packer) Add(mq *MultiQuad) {
	p.surfaces = append(p.surfaces, toPacket(mq))
}

func (p *Packer) reset() {
	p.packed = nil
	for _, s := range p.surfaces {
		p.move(s, 0, 0)
	}
}

func overlapRect(r1, r2 rect) bool {
	if r1.x >= r2.x+r2.w || r2.x >= r1.x+r1.w {
		return false
	}
	if r1.y >= r2.y+r2.h || r2.y >= r1.y+r1.h {
		return false
	}
	return true
}

func overlapSurface(s1, s2 *packet) bool {
	return overlapRect(s1.r1, s2.r1) || overlapRect(s1.r1, s2.r2) || overlapRect(s1.r2, s2.r1) || overlapRect(s1.r2, s2.r2)
}

func (p *Packer) check(s1 *packet) bool {
	if s1.x < 0 || s1.y < 0 || s1.x+s1.r1.w > p.w || s1.y+s1.r2.h > p.h {
		return false
	}
	for _, s := range p.packed {
		if overlapSurface(s1, s) {
			return false
		}
	}
	return true
}

func (p *Packer) move(s *packet, x, y int32) bool {
	s.x = x
	s.y = y
	s.r1.x = s.x
	s.r1.y = s.y
	s.r2.x = s.x
	s.r2.y = s.y
	return p.check(s)
}

func (p *Packer) toBottomLeft(s *packet) bool {
	x := p.w - s.r1.w
	y := p.h - s.r2.h
	return p.move(s, x, y)
}

func (p *Packer) pack(s *packet) bool {
	if !p.toBottomLeft(s) {
		return false
	}

	var directionX int32 = 0
	var directionY int32 = -1
	lastX := s.x
	lastY := s.y
	count := 0

	swap := func() {
		if directionX == 0 {
			directionX = -1
			directionY = 0
		} else {
			directionX = 0
			directionY = -1
		}
	}

	for {
		if p.move(s, lastX+directionX, lastY+directionY) {
			lastX = s.x
			lastY = s.y
			count = 0
		} else {
			p.move(s, lastX, lastY)
			swap()
			count += 1
			if count >= 2 {
				break
			}
		}
	}

	result := p.check(s)
	if result {
		p.packed = append(p.packed, s)
	}

	return result
}

func (p *Packer) tryPack() bool {
	for _, s := range p.surfaces {
		if !p.pack(s) {
			return false
		}
	}
	return true
}

func (p *Packer) PackAll() {
	slices.SortFunc(p.surfaces, func(s1, s2 *packet) int {
		return int(s2.area - s1.area)
	})
	for {
		p.reset()
		if p.tryPack() {
			break
		} else {
			p.h *= 2
			p.w *= 2
		}
	}
	for _, s := range p.packed {
		s.surface.pos[0] = s.x
		s.surface.pos[1] = s.y
	}
}

func newPacker(w, h int32) (p *Packer) {
	p = new(Packer)
	p.w = w
	p.h = h
	return
}

func packSurfaces(surfaces []*MultiQuad) int32 {
	p := newPacker(16, 16)
	for _, s := range surfaces {
		p.Add(s)
	}
	p.PackAll()
	return p.w
}
