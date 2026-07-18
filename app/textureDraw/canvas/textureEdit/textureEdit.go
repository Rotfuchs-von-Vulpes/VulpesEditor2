package textureEdit

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/history"
	"VulpesEditor/app/util"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"slices"

	"github.com/AllenDang/cimgui-go/backend"
)

var idSys = util.NewIdSystem()

type pixelChange struct {
	pos    [2]int32
	before [4]float32
	after  [4]float32
}

type Image struct {
	Img *image.RGBA
	Tex *backend.Texture
}

type LayerEdit struct {
	parent  *TextureEdit
	Id      int32
	width   uint32
	height  uint32
	Texture *texture.Texture
	Show    bool
	Image   *Image
}

func (s *LayerEdit) updatePreview() {
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

func (s *LayerEdit) unchange(changes []pixelChange) {
	for _, change := range changes {
		s.Texture.Set(change.pos, change.before)
	}
	s.updatePreview()
	s.parent.UpdateTexture()
}

func (s *LayerEdit) change(changes []pixelChange) {
	for _, change := range changes {
		s.Texture.Set(change.pos, change.after)
	}
	s.updatePreview()
	s.parent.UpdateTexture()
}

type TextureChange struct {
	parent  *LayerEdit
	changes []pixelChange
}

func (s *TextureChange) Undo() {
	s.parent.unchange(s.changes)
}

func (s *TextureChange) Redo() {
	s.parent.change(s.changes)
}

func (s *LayerEdit) Change(pixels []texture.PixelEdit) {
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

type preview struct {
	layerIdx int
	pixels   []texture.PixelEdit
}

func (s *preview) clear() {
	s.pixels = make([]texture.PixelEdit, 0)
}

type TextureEdit struct {
	Id          int32
	Width       uint32
	Height      uint32
	Aspect      float32
	GlID        uint32
	Layers      []*LayerEdit
	texture     *texture.Texture
	preview     *preview
	OnLayerUndo func()
}

func (s *TextureEdit) SetLayer(idx int) {
	if idx >= len(s.Layers) || idx < 0 {
		panic(fmt.Sprintf("Illegal layer index: %d of length %d", idx, len(s.Layers)))
	}
	s.preview.clear()
	s.preview.layerIdx = idx
}

func (s *TextureEdit) addLayer(idx int, tex *texture.Texture) {
	layer := new(LayerEdit)
	layer.Id = idSys.GetID()
	layer.parent = s
	layer.width = s.Width
	layer.height = s.Height
	layer.Texture = tex
	layer.Show = true
	layer.Image = new(Image)
	layer.Image.Img = image.NewRGBA(image.Rect(0, 0, int(s.Width), int(s.Height)))
	layer.Image.Tex = backend.NewTextureFromRgba(layer.Image.Img)
	s.Layers = slices.Insert(s.Layers, idx, layer)
}

func (s *TextureEdit) AppendLayer() {
	c := new(LayersChange)
	c.parent = s
	c.before = slices.Clone(s.Layers)
	s.addLayer(len(s.Layers), texture.New(s.Width, s.Height))
	c.after = slices.Clone(s.Layers)
	history.Append(c)
}

func New(tex *texture.Texture) (out *TextureEdit) {
	out = new(TextureEdit)
	out.Id = idSys.GetID()
	out.Width = tex.Width
	out.Height = tex.Height
	out.Aspect = float32(tex.Width) / float32(tex.Height)
	out.addLayer(0, tex)
	out.texture = texture.New(tex.Width, tex.Height)
	out.GlID = renderer.CreateTexture(int32(tex.Width), int32(tex.Height), tex.FlatColors())
	out.preview = new(preview)
	out.preview.layerIdx = 0
	return
}

var isEditing = false
var layerChange *LayersChange

func (s *TextureEdit) beginLayerEdit() {
	if isEditing {
		panic("Too much begin edit call")
	}
	isEditing = true
	layerChange = new(LayersChange)
	layerChange.parent = s
	layerChange.before = slices.Clone(s.Layers)
}

func (s *TextureEdit) endLayerEdit() {
	if !isEditing {
		panic("Too much end edit call")
	}
	isEditing = false
	layerChange.after = slices.Clone(s.Layers)
	c := *layerChange
	history.Append(&c)
	layerChange = nil
	s.update()
}

type LayersChange struct {
	parent *TextureEdit
	before []*LayerEdit
	after  []*LayerEdit
}

func (s *LayersChange) Undo() {
	s.parent.Layers = slices.Clone(s.before)
	s.parent.update()
	s.parent.OnLayerUndo()
}

func (s *LayersChange) Redo() {
	s.parent.Layers = slices.Clone(s.after)
	s.parent.update()
	s.parent.OnLayerUndo()
}

func (s *TextureEdit) update() {
	for _, layer := range s.Layers {
		layer.updatePreview()
	}
	s.UpdateTexture()
}

func (s *TextureEdit) delete(toDelete []bool) {
	if len(toDelete) != len(s.Layers) {
		panic(fmt.Sprintf("Wrong list length: %d remove indices, %d layers count", len(toDelete), len(s.Layers)))
	}
	final := []*LayerEdit{}
	for i := range s.Layers {
		if !toDelete[i] {
			final = append(final, s.Layers[i])
		}
	}
	s.Layers = final
}

func (s *TextureEdit) Remove(toDelete []bool) {
	if len(toDelete) != len(s.Layers) {
		panic(fmt.Sprintf("Wrong list length: %d remove indices, %d layers count", len(toDelete), len(s.Layers)))
	}
	count := 0
	for _, b := range toDelete {
		if b {
			count += 1
		}
	}
	length := len(s.Layers)
	s.beginLayerEdit()
	s.delete(toDelete)
	if count == length {
		s.addLayer(len(s.Layers), texture.New(s.Width, s.Height))
	}
	s.endLayerEdit()
}

func (s *TextureEdit) Swap(idx1, idx2 int) {
	if idx1 < 0 || idx1 >= len(s.Layers) {
		panic(fmt.Sprintf("Out of bounds: %d of length %d", idx1, len(s.Layers)))
	}
	if idx2 < 0 || idx2 >= len(s.Layers) {
		panic(fmt.Sprintf("Out of bounds: %d of length %d", idx2, len(s.Layers)))
	}
	if idx1 == idx2 {
		panic(fmt.Sprintf("Same indices: %d and %d", idx1, idx2))
	}
	s.beginLayerEdit()
	temp := s.Layers[idx1]
	s.Layers[idx1] = s.Layers[idx2]
	s.Layers[idx2] = temp
	s.endLayerEdit()
}

func (s *TextureEdit) Merge(merge []bool) {
	if len(merge) != len(s.Layers) {
		panic(fmt.Sprintf("Wrong list length: %d merge indices, %d layers count", len(merge), len(s.Layers)))
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
			}
			toDelete[i] = true
			tempTex.Colors = texture.Merge(tempTex, s.Layers[i].Texture)
		}
	}
	s.beginLayerEdit()
	s.delete(toDelete)
	s.addLayer(resultIdx, tempTex)
	s.endLayerEdit()
}

func (s *TextureEdit) UpdateTexture() {
	s.texture.Clear()
	for i, layer := range s.Layers {
		if layer.Show {
			if s.preview.layerIdx == i {
				tex := texture.New(s.Width, s.Height)
				tex.Colors = slices.Clone(layer.Texture.Colors)
				tex.BulkSet(s.preview.pixels)
				s.texture.Colors = texture.Merge(s.texture, tex)
			} else {
				s.texture.Colors = texture.Merge(s.texture, layer.Texture)
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
	s.UpdateTexture()
}

func (s *TextureEdit) ResetPreview() {
	s.preview.clear()
	s.UpdateTexture()
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
		s.texture.Colors = texture.Merge(s.texture, layer.Texture)
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
