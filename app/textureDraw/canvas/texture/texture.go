package texture

import (
	"VulpesEditor/app/util"
	"fmt"
)

func blankTexture(width, height uint32) (data [][4]float32) {
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			data = append(data, [4]float32{0, 0, 0, 0})
		}
	}
	return
}

type Texture struct {
	Width  uint32
	Height uint32
	Colors [][4]float32
}

type PixelEdit struct {
	Pos   [2]int32
	Color [4]float32
}

func SetEditColor(pixels [][2]int32, color [4]float32) (out []PixelEdit) {
	for _, pos := range pixels {
		out = append(out, PixelEdit{pos, color})
	}
	return
}

func New(width, height uint32) (out *Texture) {
	colors := blankTexture(width, height)
	out = new(Texture)
	out.Width = width
	out.Height = height
	out.Colors = colors
	return
}

func (s *Texture) Resize(width, height uint32) {
	if width == s.Width && height == s.Height {
		return
	}
	s.Width = width
	s.Height = height
	s.Colors = blankTexture(width, height)
}

func (s *Texture) Get(pos [2]int32) (ok bool, color [4]float32) {
	index := int(pos[1]*int32(s.Width) + pos[0])
	if pos[0] < 0 || pos[0] >= int32(s.Width) || pos[1] < 0 || pos[1] >= int32(s.Height) {
		ok = false
	} else {
		ok = true
		color = s.Colors[index]
	}
	return
}

func (s *Texture) Set(pos [2]int32, color [4]float32) (ok bool) {
	index := int(pos[1]*int32(s.Width) + pos[0])
	if pos[0] < 0 || pos[0] >= int32(s.Width) || pos[1] < 0 || pos[1] >= int32(s.Height) {
		ok = false
	} else {
		ok = true
		s.Colors[index] = color
	}
	return
}

func (s *Texture) BulkSet(pixels []PixelEdit) (ok bool) {
	for _, pixel := range pixels {
		newOk := s.Set(pixel.Pos, pixel.Color)
		ok = newOk || ok
	}
	return
}

func (s *Texture) BulkSetColor(pixels [][2]int32, color [4]float32) (ok bool) {
	for _, pos := range pixels {
		ok = ok || s.Set(pos, color)
	}
	return
}

func (s *Texture) Clear() {
	s.Colors = blankTexture(s.Width, s.Height)
}

func (s *Texture) FlatColors() (data []float32) {
	for _, color := range s.Colors {
		data = append(data, color[0], color[1], color[2], color[3])
	}
	return
}

func Merge(end, front *Texture) (colors [][4]float32) {
	if end.Height != front.Height || end.Width != front.Width {
		panic(fmt.Sprintf("Wrong Size, end size: %d, %d; front size: %d, %d", end.Width, end.Height, front.Width, front.Height))
	}
	tempTex := New(end.Width, end.Height)
	for x := range end.Width {
		for y := range end.Height {
			pos := [2]int32{int32(x), int32(y)}
			_, c1 := end.Get(pos)
			_, c2 := front.Get(pos)
			if c2[3] >= 1 {
				tempTex.Set(pos, c2)
			} else if c2[3] <= 0 {
				tempTex.Set(pos, c1)
			} else {
				c1[0] *= c1[3]
				c1[1] *= c1[3]
				c1[2] *= c1[3]
				c2[0] *= c2[3]
				c2[1] *= c2[3]
				c2[2] *= c2[3]
				red := (1-c2[3])*c1[0] + c2[0]
				green := (1-c2[3])*c1[1] + c2[1]
				blue := (1-c2[3])*c1[2] + c2[2]
				alpha := 1 - (1-c2[3])*(1-c1[3])
				c3 := [4]float32{red, green, blue, alpha}
				tempTex.Set(pos, c3)
			}
		}
	}
	return tempTex.Colors
}

var idSys = util.NewIdSystem()
