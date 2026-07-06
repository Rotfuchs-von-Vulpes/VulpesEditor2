package textureEdit

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/util"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"slices"
)

var idSys = util.NewIdSystem()

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type LayerEdit struct {
	parent    *TextureEdit
	Id        int32
	width     uint32
	height    uint32
	texture   *texture.Texture
	changes   [][]pixelChange
	undoLevel int32
	Show      bool
}

func (s *LayerEdit) unchange(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.before)
	}
	s.parent.UpdateTexture()
}

func (s *LayerEdit) change(changes []pixelChange) {
	for _, change := range changes {
		s.texture.Set(change.pos, change.after)
	}
	s.parent.UpdateTexture()
}

func (s *LayerEdit) applyChanges(changes []pixelChange) {
	s.change(changes)
	s.changes = s.changes[:len(s.changes)-int(s.undoLevel)]
	s.changes = append(s.changes, changes)
	s.undoLevel = 0
}

func (s *LayerEdit) Undo() bool {
	changesIdx := len(s.changes) - 1 - int(s.undoLevel)
	if changesIdx >= 0 {
		lastChanges := s.changes[changesIdx]
		s.undoLevel++
		s.unchange(lastChanges)
		return true
	}
	return false
}

func (s *LayerEdit) Redo() {
	if s.undoLevel > 0 {
		changesIdx := len(s.changes) - int(s.undoLevel)
		if changesIdx >= 0 {
			lastChanges := s.changes[changesIdx]
			s.undoLevel--
			s.change(lastChanges)
		}
	}
}

func (s *LayerEdit) Change(pixels []texture.PixelEdit) {
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

type preview struct {
	layerIdx int
	pixels   []texture.PixelEdit
}

func (s *preview) clear() {
	s.pixels = make([]texture.PixelEdit, 0)
}

type TextureEdit struct {
	Id      int32
	Width   uint32
	Height  uint32
	Aspect  float32
	GlID    uint32
	Layers  []*LayerEdit
	texture *texture.Texture
	preview *preview
}

func (s *TextureEdit) SetLayer(idx int) {
	if idx >= len(s.Layers) || idx < 0 {
		panic(fmt.Sprintf("Illegal layer index: %d of length %d", idx, len(s.Layers)))
	}
	s.preview.clear()
	s.preview.layerIdx = idx
}

func (s *TextureEdit) AddLayer() {
	newLayer := new(LayerEdit)
	newLayer.Id = idSys.GetID()
	newLayer.parent = s
	newLayer.width = s.Width
	newLayer.height = s.Height
	newLayer.texture = texture.New(s.Width, s.Height)
	newLayer.Show = true
	s.Layers = append(s.Layers, newLayer)
}

func New(tex *texture.Texture) (out *TextureEdit) {
	out = new(TextureEdit)
	out.Id = idSys.GetID()
	out.Width = tex.Width
	out.Height = tex.Height
	out.Aspect = float32(tex.Width) / float32(tex.Height)
	out.AddLayer()
	out.Layers[0].texture = tex
	out.texture = texture.New(tex.Width, tex.Height)
	out.GlID = renderer.CreateTexture(int32(tex.Width), int32(tex.Height), tex.FlatColors())
	out.preview = new(preview)
	out.preview.layerIdx = 0
	return
}

func (s *TextureEdit) Remove(toDelete []bool) {
	if len(toDelete) != len(s.Layers) {
		panic(fmt.Sprintf("Wrong list length: %d merge indexes, %d layers cout", len(toDelete), len(s.Layers)))
	}
	final := []*LayerEdit{}
	for i := range s.Layers {
		if !toDelete[i] {
			final = append(final, s.Layers[i])
		}
	}
	if len(final) == 0 {
		return
	}
	s.Layers = final
}

func (s *TextureEdit) Merge(merge []bool) {
	if len(merge) != len(s.Layers) {
		panic(fmt.Sprintf("Wrong list length: %d merge indexes, %d layers cout", len(merge), len(s.Layers)))
	}
	count := 0
	for _, b := range merge {
		if b {
			count += 1
		}
	}
	if count < 2 {
		return
	}
	tempTex := texture.New(s.Width, s.Height)
	resultIdx := 0
	first := true
	toDelete := make([]bool, len(s.Layers))
	for i := range s.Layers {
		if merge[i] {
			if first {
				first = false
				resultIdx = i
			} else {
				toDelete[i] = true
			}
			tempTex.Colors = texture.Merge(tempTex, s.Layers[i].texture)
		}
	}
	s.Layers[resultIdx].texture.Colors = tempTex.Colors
	s.Remove(toDelete)
}

func (s *TextureEdit) UpdateTexture() {
	s.texture.Clear()
	for i, layer := range s.Layers {
		if layer.Show {
			if s.preview.layerIdx == i {
				tex := texture.New(s.Width, s.Height)
				tex.Colors = slices.Clone(layer.texture.Colors)
				tex.BulkSet(s.preview.pixels)
				s.texture.Colors = texture.Merge(s.texture, tex)
			} else {
				s.texture.Colors = texture.Merge(s.texture, layer.texture)
			}
		}
	}
	renderer.WriteTexture(s.GlID, int32(s.Width), int32(s.Height), s.texture.FlatColors())
}

func (s *TextureEdit) Colors() [][4]float32 {
	return s.texture.Colors
}

func (s *TextureEdit) ChangePreview(pixels []texture.PixelEdit) {
	s.preview.clear()
	s.preview.pixels = pixels
}

func (s *TextureEdit) ResetPreview() {
	s.preview.clear()
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

	s.texture.Clear()
	for _, layer := range s.Layers {
		s.texture.Colors = texture.Merge(s.texture, layer.texture)
	}

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
