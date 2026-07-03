package image

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/util"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
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

func NewTexture(width, height uint32) (out *Texture) {
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

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type TextureEdit struct {
	Id          int32
	Width       uint32
	Height      uint32
	Aspect      float32
	texture     *Texture
	textureGlID uint32
	previewGlID uint32
	preview     *Texture
	changes     [][]pixelChange
	undoLevel   int32
}

func (s *TextureEdit) unchange(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.before)
	}
	s.updateTexture()
}

func (s *TextureEdit) change(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.after)
	}
	s.updateTexture()
}

func (s *TextureEdit) applyChanges(changes []pixelChange) {
	s.change(changes)
	s.changes = s.changes[:len(s.changes)-int(s.undoLevel)]
	s.changes = append(s.changes, changes)
	s.undoLevel = 0
}

func NewTextureEdit(tex *Texture) (out *TextureEdit) {
	out = new(TextureEdit)
	out.Width = tex.Width
	out.Height = tex.Height
	out.Aspect = float32(tex.Width) / float32(tex.Height)
	out.texture = tex
	out.preview = NewTexture(tex.Width, tex.Height)
	out.textureGlID = renderer.CreateTexture(int32(tex.Width), int32(tex.Height), flatData(tex.Colors))
	out.previewGlID = renderer.CreateTexture(int32(tex.Width), int32(tex.Height), flatData(out.preview.Colors))
	return
}

func (s *TextureEdit) updateTexture() {
	renderer.WriteTexture(s.textureGlID, int32(s.Width), int32(s.Height), flatData(s.texture.Colors))
}

func (s *TextureEdit) updatePreview() {
	renderer.WriteTexture(s.previewGlID, int32(s.Width), int32(s.Height), flatData(s.preview.Colors))
}

func (s *TextureEdit) Undo() bool {
	changesIdx := len(s.changes) - 1 - int(s.undoLevel)
	if changesIdx >= 0 {
		lastChanges := s.changes[changesIdx]
		s.undoLevel++
		s.unchange(lastChanges)
		return true
	}
	return false
}

func (s *TextureEdit) Redo() {
	if s.undoLevel > 0 {
		changesIdx := len(s.changes) - int(s.undoLevel)
		if changesIdx >= 0 {
			lastChanges := s.changes[changesIdx]
			s.undoLevel--
			s.change(lastChanges)
		}
	}
}

func (s *TextureEdit) Colors() [][4]float32 {
	return s.texture.Colors
}

func (s *TextureEdit) Change(pixels []PixelEdit) {
	if len(pixels) > 0 {
		changes := []pixelChange{}
		for _, pixel := range pixels {
			if ok, beforeColor := s.texture.Get(pixel.Pos); ok {
				var change pixelChange
				change.pos = pixel.Pos
				change.before = beforeColor
				change.after = pixel.Color
				changes = append(changes, change)
			}
		}
		if len(changes) > 0 {
			s.applyChanges(changes)
		}
	}
}

func (s *TextureEdit) ChangePreview(pixels []PixelEdit) {
	if len(pixels) > 0 {
		if changed := s.preview.BulkSet(pixels); changed {
			s.updatePreview()
		}
	}
}

func (s *TextureEdit) ResetPreview() {
	s.preview.Clear()
	s.updatePreview()
}

func (s *TextureEdit) GlID() (uint32, uint32) {
	return s.textureGlID, s.previewGlID
}

func (s *TextureEdit) SaveTextureAsFile(fileName, path string) bool {
	if path == "" {
		path = "./UserData/textures"
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println(err)
		return false
	}
	file, err := os.Create(path + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, int(s.texture.Width), int(s.texture.Height)))
	for x := int32(0); x < int32(s.texture.Width); x++ {
		for y := int32(0); y < int32(s.texture.Height); y++ {
			_, rgba := s.texture.Get([2]int32{x, y})
			alpha := rgba[3]
			red := uint8(255 * rgba[0] * alpha)
			green := uint8(255 * rgba[1] * alpha)
			blue := uint8(255 * rgba[2] * alpha)
			img.SetRGBA(int(x), int(y), color.RGBA{red, green, blue, uint8(255 * rgba[3])})
		}
	}
	if err := png.Encode(file, img); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

var idSys = util.NewIdSystem()
