package textureEdit

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/history"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/util"
	"bytes"
	"fmt"
	"image"
	"os"
	"slices"

	"github.com/AllenDang/cimgui-go/backend"
)

var idSys = util.NewIdSystem()

type preview struct {
	pixels []texture.PixelEdit
}

func (s *preview) clear() {
	s.pixels = make([]texture.PixelEdit, 0)
}

type Image struct {
	Img *image.RGBA
	Tex *backend.Texture
}

type TextureEdit struct {
	Id      int32
	Width   uint32
	Height  uint32
	Aspect  float32
	GlID    uint32
	layers  []*layerEdit
	layer   *layerEdit
	texture *texture.Texture
	preview *preview
}

func (s *TextureEdit) addLayer(idx int, tex *texture.Texture) {
	layer := new(layerEdit)
	layer.Id = idSys.GetID()
	layer.parent = s
	layer.width = s.Width
	layer.height = s.Height
	layer.Texture = tex
	layer.Show = true
	layer.Image = new(Image)
	layer.Image.Img = image.NewRGBA(image.Rect(0, 0, int(s.Width), int(s.Height)))
	layer.Image.Tex = backend.NewTextureFromRgba(layer.Image.Img)
	s.layers = slices.Insert(s.layers, idx, layer)
}

func (s *TextureEdit) AppendLayer() {
	c := new(layersChange)
	c.parent = s
	c.before = slices.Clone(s.layers)
	s.addLayer(len(s.layers), texture.New(s.Width, s.Height))
	s.layer = s.layers[len(s.layers)-1]
	c.after = slices.Clone(s.layers)
	history.Append(c)
}

func New(tex *texture.Texture) (out *TextureEdit) {
	out = new(TextureEdit)
	out.Id = idSys.GetID()
	out.Width = tex.Width
	out.Height = tex.Height
	out.Aspect = float32(tex.Width) / float32(tex.Height)
	out.addLayer(0, tex)
	out.layer = out.layers[0]
	out.texture = texture.New(tex.Width, tex.Height)
	out.GlID = renderer.CreateTexture(int32(tex.Width), int32(tex.Height), tex.FlatColors())
	out.preview = new(preview)
	return
}

var isEditing = false
var layerChange *layersChange

func (s *TextureEdit) onLayerEdit() {
	found := false
	for _, layer := range s.layers {
		if layer.Id == s.layer.Id {
			found = true
			break
		}
	}
	if !found {
		s.layer = s.layers[len(s.layers)-1]
	}
}

func (s *TextureEdit) beginLayerEdit() {
	if isEditing {
		panic("Too much begin edit call")
	}
	isEditing = true
	layerChange = new(layersChange)
	layerChange.parent = s
	layerChange.before = slices.Clone(s.layers)
}

func (s *TextureEdit) endLayerEdit() {
	if !isEditing {
		panic("Too much end edit call")
	}
	isEditing = false
	layerChange.after = slices.Clone(s.layers)
	c := *layerChange
	history.Append(&c)
	layerChange = nil
	s.update()
}

type layersChange struct {
	parent *TextureEdit
	before []*layerEdit
	after  []*layerEdit
}

func (s *layersChange) Undo() {
	s.parent.layers = slices.Clone(s.before)
	s.parent.update()
	s.parent.onLayerEdit()
}

func (s *layersChange) Redo() {
	s.parent.layers = slices.Clone(s.after)
	s.parent.update()
	s.parent.onLayerEdit()
}

func (s *TextureEdit) update() {
	for _, layer := range s.layers {
		layer.updatePreview()
	}
	s.UpdateTexture()
}

func (s *TextureEdit) delete(toDelete []bool) {
	if len(toDelete) != len(s.layers) {
		panic(fmt.Sprintf("Wrong list length: %d remove indices, %d layers count", len(toDelete), len(s.layers)))
	}
	final := []*layerEdit{}
	for i := range s.layers {
		if !toDelete[i] {
			final = append(final, s.layers[i])
		}
	}
	s.layers = final
}

func (s *TextureEdit) Remove(toDelete []bool) {
	if len(toDelete) != len(s.layers) {
		panic(fmt.Sprintf("Wrong list length: %d remove indices, %d layers count", len(toDelete), len(s.layers)))
	}
	count := 0
	for _, b := range toDelete {
		if b {
			count += 1
		}
	}
	length := len(s.layers)
	s.beginLayerEdit()
	s.delete(toDelete)
	if count == length {
		s.addLayer(len(s.layers), texture.New(s.Width, s.Height))
	}
	s.endLayerEdit()
	s.onLayerEdit()
}

func (s *TextureEdit) Swap(idx1, idx2 int) {
	if idx1 < 0 || idx1 >= len(s.layers) {
		panic(fmt.Sprintf("Out of bounds: %d of length %d", idx1, len(s.layers)))
	}
	if idx2 < 0 || idx2 >= len(s.layers) {
		panic(fmt.Sprintf("Out of bounds: %d of length %d", idx2, len(s.layers)))
	}
	if idx1 == idx2 {
		panic(fmt.Sprintf("Same indices: %d and %d", idx1, idx2))
	}
	s.beginLayerEdit()
	temp := s.layers[idx1]
	s.layers[idx1] = s.layers[idx2]
	s.layers[idx2] = temp
	s.endLayerEdit()
	s.onLayerEdit()
}

func (s *TextureEdit) Merge(merge []bool) {
	if len(merge) != len(s.layers) {
		panic(fmt.Sprintf("Wrong list length: %d merge indices, %d layers count", len(merge), len(s.layers)))
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
	toDelete := make([]bool, len(s.layers))
	for i := range s.layers {
		if merge[i] {
			if first {
				first = false
				resultIdx = i
			}
			toDelete[i] = true
			tempTex.Colors = texture.Merge(tempTex, s.layers[i].Texture)
		}
	}
	s.beginLayerEdit()
	s.delete(toDelete)
	s.addLayer(resultIdx, tempTex)
	s.endLayerEdit()
	s.onLayerEdit()
}

func (s *TextureEdit) UpdateTexture() {
	s.texture.Clear()
	for _, layer := range s.layers {
		if layer.Show {
			if s.layer.Id == layer.Id {
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

func (s *TextureEdit) LayerColors() [][4]float32 {
	return slices.Clone(s.layer.Texture.Colors)
}

func (s *TextureEdit) Colors() [][4]float32 {
	return slices.Clone(s.texture.Colors)
}

func (s *TextureEdit) LayerChange(pixels []texture.PixelEdit) {
	s.layer.ApplyChange(pixels)
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
	for _, layer := range s.layers {
		s.texture.Colors = texture.Merge(s.texture, layer.Texture)
	}

	if err := s.texture.ToPNG(file); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (s *TextureEdit) Save(r *file.ArchiveWriter) {
	for i, layer := range s.layers {
		buff := bytes.NewBuffer(nil)
		layer.Texture.ToPNG(buff)
		r.Write(fmt.Sprintf("layers/layer%d.png", i), buff.Bytes())
	}
}

func Open(r *file.ArchiveReader) (out *TextureEdit, err error) {
	var width uint32
	var height uint32
	var layers []*texture.Texture
	count := 0
	for {
		name := fmt.Sprintf("layers/layer%d.png", count)
		f, err := r.Open(name)
		if err != nil {
			break
		}
		defer f.Close()
		tex, err := texture.DecodePNG(f)
		if err != nil {
			return nil, err
		}
		width = tex.Width
		height = tex.Height
		layers = append(layers, tex)
		count += 1
	}
	if count == 0 {
		return nil, fmt.Errorf("No texture")
	}
	out = new(TextureEdit)
	out.Id = idSys.GetID()
	out.Width = width
	out.Height = height
	out.Aspect = float32(width) / float32(height)
	for i, layer := range layers {
		out.addLayer(i, layer)
	}
	out.layer = out.layers[0]
	out.texture = texture.New(width, height)
	out.GlID = renderer.CreateTexture(int32(width), int32(height), out.layer.Texture.FlatColors())
	out.preview = new(preview)
	out.update()
	return
}
