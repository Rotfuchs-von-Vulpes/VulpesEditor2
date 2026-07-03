package texture

import (
	"VulpesEditor/app/util"
)

func blankTexture(width, height uint32) (data [][4]float32) {
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			data = append(data, [4]float32{0, 0, 0, 0})
		}
	}
	return
}

func flatData(colors [][4]float32) (data []float32) {
	for _, color := range colors {
		data = append(data, color[0], color[1], color[2], color[3])
	}
	return
}

type Texture struct {
	Id     int32
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
	id := idSys.GetID()
	colors := blankTexture(width, height)
	out = new(Texture)
	out.Id = id
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
	index := int(pos[1]*int32(s.Height) + pos[0])
	if pos[0] < 0 || pos[0] >= int32(s.Width) || pos[1] < 0 || pos[1] >= int32(s.Height) {
		ok = false
	} else {
		ok = true
		color = s.Colors[index]
	}
	return
}

func (s *Texture) Set(pos [2]int32, color [4]float32) (ok bool) {
	index := int(pos[1]*int32(s.Height) + pos[0])
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

func (s *Texture) FlatColors() []float32 {
	return flatData(s.Colors)
}

var idSys = util.NewIdSystem()
