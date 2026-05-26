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
	glID   uint32
}

func NewTexture(width, height uint32) (out *Texture) {
	id := idSys.GetID()
	colors := blankTexture(width, height)
	glID := renderer.CreateTexture(int32(width), int32(height), flatData(colors))
	out = new(Texture)
	out.Id = id
	out.Width = width
	out.Height = height
	out.Colors = colors
	out.glID = glID
	return
}

func (s *Texture) Get(pos [2]int32) (ok bool, color [4]float32) {
	index := int(pos[1]*int32(s.Height) + pos[0])
	if index < 0 || index > len(s.Colors)-1 {
		ok = false
	} else {
		ok = true
		color = s.Colors[index]
	}
	return
}

func (s *Texture) Set(pos [2]int32, color [4]float32) (ok bool) {
	index := int(pos[1]*int32(s.Height) + pos[0])
	if index < 0 || index > len(s.Colors)-1 {
		ok = false
	} else {
		ok = true
		s.Colors[index] = color
	}
	return
}

func (s *Texture) Clear() {
	s.Colors = blankTexture(s.Width, s.Height)
	s.update()
}

func (s *Texture) update() {
	data := flatData(s.Colors)
	renderer.WriteTexture(s.glID, int32(s.Width), int32(s.Height), data)
}

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type pixelsChange []pixelChange

type TextureEdit struct {
	Id        int32
	Width     uint32
	Height    uint32
	Aspect    float32
	texture   *Texture
	preview   *Texture
	changes   []pixelsChange
	undoLevel int32
}

func (s *TextureEdit) unchange(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.before)
	}
	s.texture.update()
}

func (s *TextureEdit) change(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.after)
	}
	s.texture.update()
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
	return
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

func (s *TextureEdit) ChangePreview(pixels [][2]int32, color [4]float32) {
	if len(pixels) > 0 {
		changed := false
		for _, pixel := range pixels {
			if pixel[0] < 0 || pixel[0] >= int32(s.texture.Width) || pixel[1] < 0 || pixel[1] >= int32(s.texture.Height) {
				continue
			}
			s.preview.Set(pixel, color)
			changed = true
		}
		if changed {
			s.preview.update()
		}
	}
}

func (s *TextureEdit) ChangeTexture(pixels [][2]int32, color [4]float32) {
	if len(pixels) > 0 {
		changes := []pixelChange{}
		for _, pixel := range pixels {
			if ok, beforeColor := s.texture.Get(pixel); ok {
				var change pixelChange
				change.pos = pixel
				change.before = beforeColor
				change.after = color
				changes = append(changes, change)
			}
		}
		if len(changes) > 0 {
			s.applyChanges(changes)
		}
	}
}

func (s *TextureEdit) ResetPreview() {
	s.preview.Clear()
}

func (s *TextureEdit) GlID() (uint32, uint32) {
	return s.texture.glID, s.preview.glID
}

func (s *TextureEdit) SaveTextureAsFile(fileName, path string) bool {
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

var idSys *util.IdSystem

func Init() {
	idSys = util.NewIdSystem()
}
