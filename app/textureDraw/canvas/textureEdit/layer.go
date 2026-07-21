package textureEdit

import (
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/history"
	"math"

	"github.com/AllenDang/cimgui-go/backend"
)

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type TextureChange struct {
	parent  *layerEdit
	changes []pixelChange
}

func (s *TextureChange) Undo() {
	s.parent.unchange(s.changes)
}

func (s *TextureChange) Redo() {
	s.parent.change(s.changes)
}

type layerEdit struct {
	parent  *TextureEdit
	Id      int32
	width   uint32
	height  uint32
	Texture *texture.Texture
	Show    bool
	Image   *Image
}

func (s *layerEdit) updatePreview() {
	floatData := s.Texture.FlatColors()

	for i := 0; i < int(s.Texture.Width*s.Texture.Height); i++ {
		offset := i * 4

		r := uint8(math.Min(255.0, float64(floatData[offset])*255.0))
		g := uint8(math.Min(255.0, float64(floatData[offset+1])*255.0))
		b := uint8(math.Min(255.0, float64(floatData[offset+2])*255.0))
		a := uint8(math.Min(255.0, float64(floatData[offset+3])*255.0))

		s.Image.Img.Pix[offset] = r
		s.Image.Img.Pix[offset+1] = g
		s.Image.Img.Pix[offset+2] = b
		s.Image.Img.Pix[offset+3] = a
	}

	s.Image.Tex.Release()
	s.Image.Tex = backend.NewTextureFromRgba(s.Image.Img)
}

func (s *layerEdit) unchange(changes []pixelChange) {
	for _, change := range changes {
		s.Texture.Set(change.pos, change.before)
	}
	s.updatePreview()
	s.parent.UpdateTexture()
}

func (s *layerEdit) change(changes []pixelChange) {
	for _, change := range changes {
		s.Texture.Set(change.pos, change.after)
	}
	s.updatePreview()
	s.parent.UpdateTexture()
}

func (s *layerEdit) ApplyChange(pixels []texture.PixelEdit) {
	if len(pixels) > 0 {
		changes := []pixelChange{}
		for _, pixel := range pixels {
			if ok, beforeColor := s.Texture.Get(pixel.Pos); ok {
				var change pixelChange
				change.pos = pixel.Pos
				change.before = beforeColor
				change.after = pixel.Color
				changes = append(changes, change)
			}
		}
		if len(changes) > 0 {
			s.change(changes)
			c := new(TextureChange)
			c.parent = s
			c.changes = changes
			history.Append(c)
		}
	}
}
