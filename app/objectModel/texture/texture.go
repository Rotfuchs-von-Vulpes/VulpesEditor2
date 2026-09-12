package texture

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/objectModel/model"
	"VulpesEditor/app/textureDraw"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"bytes"
	"slices"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type quad struct {
	pos     [2]int32
	width   int32
	height  int32
	texture *texture.Texture
}

type MultiQuad struct {
	id        string
	pos       [2]int32
	width     int32
	height    int32
	cutWidth  int32
	cutHeight int32
	quads     [6]quad
}

func toInt(v1 [2]float32, size int32) (v2 [2]int32) {
	v2[0] = int32(v1[0] * float32(size))
	v2[1] = int32(v1[1] * float32(size))
	return
}

func writeTexture(s *MultiQuad, tex *texture.Texture) {
	for _, q := range s.quads {
		var pixels []texture.PixelEdit
		for x := range q.width {
			for y := range q.height {
				pos := getPositions([2]int32{q.pos[0] + x, q.pos[1] + y}, s.pos)
				if ok, color := tex.Get(pos); ok {
					var p texture.PixelEdit
					p.Pos = [2]int32{x, y}
					p.Color = color
					pixels = append(pixels, p)
				}
			}
		}
		q.texture.BulkSet(pixels)
	}
}

func open(c *model.CubeUnit) (unity *MultiQuad) {
	unity = new(MultiQuad)
	_, size, _ := c.Data()

	x := int32(size[0] * 16)
	y := int32(size[1] * 16)
	z := int32(size[2] * 16)

	unity.quads[0] = quad{[2]int32{0, y}, x, z, texture.New(x, z)}
	unity.quads[1] = quad{[2]int32{x, y}, x, z, texture.New(x, z)}
	unity.quads[2] = quad{[2]int32{0, 0}, x, y, texture.New(x, y)}
	unity.quads[3] = quad{[2]int32{x, 0}, x, y, texture.New(x, y)}
	unity.quads[4] = quad{[2]int32{2 * x, 0}, z, y, texture.New(z, y)}
	unity.quads[5] = quad{[2]int32{2*x + z, 0}, z, y, texture.New(z, y)}

	unity.id = c.Id
	unity.pos = toInt(c.QuadPos, ctx.size)
	unity.width = 2 * (x + z)
	unity.height = y + z
	unity.cutWidth = 2 * x
	unity.cutHeight = y

	writeTexture(unity, ctx.texture)

	return
}

func primitive(c *model.CubeUnit) (unity *MultiQuad) {
	unity = new(MultiQuad)
	_, size, _ := c.Data()

	x := int32(size[0] * 16)
	y := int32(size[1] * 16)
	z := int32(size[2] * 16)

	unity.quads[0] = quad{[2]int32{0, y}, x, z, texture.New(x, z)}
	unity.quads[1] = quad{[2]int32{x, y}, x, z, texture.New(x, z)}
	unity.quads[2] = quad{[2]int32{0, 0}, x, y, texture.New(x, y)}
	unity.quads[3] = quad{[2]int32{x, 0}, x, y, texture.New(x, y)}
	unity.quads[4] = quad{[2]int32{2 * x, 0}, z, y, texture.New(z, y)}
	unity.quads[5] = quad{[2]int32{2*x + z, 0}, z, y, texture.New(z, y)}

	unity.id = c.Id
	unity.pos = [2]int32{0, 0}
	unity.width = 2 * (x + z)
	unity.height = y + z
	unity.cutWidth = 2 * x
	unity.cutHeight = y

	colors := [12][4]float32{
		{0.75, 0.75, 0.75, 1}, // white
		{1, 1, 1, 1},
		{0.75, 0.75, 0.25, 1}, // yellow
		{1, 1, 0.5, 1},
		{0.75, 0.25, 0.25, 1}, // red
		{1, 0.5, 0.5, 1},
		{0.7, 0.3, 0.2, 1}, // orange
		{1, 0.5, 0.25, 1},
		{0.25, 0.25, 0.75, 1}, // blue
		{0.5, 0.5, 1, 1},
		{0.25, 0.75, 0.25, 1}, // green
		{0.5, 1, 0.5, 1},
	}
	for i := range 6 {
		q := &unity.quads[i]
		var pixels []texture.PixelEdit
		for x := range q.width {
			for y := range q.height {
				var p texture.PixelEdit
				p.Pos = [2]int32{x, y}
				if x == 0 || x == q.width-1 || y == 0 || y == q.height-1 {
					p.Color = colors[2*i+1]
				} else {
					p.Color = colors[2*i]
				}
				pixels = append(pixels, p)
			}
		}
		q.texture.BulkSet(pixels)
	}

	return
}

func getPositions(v1, offset [2]int32) (v2 [2]int32) {
	v2[0] = v1[0] + offset[0]
	v2[1] = v1[1] + offset[1]
	return
}

func toFloat(v1, offset [2]int32, size, d float32) (v2 [2]float32) {
	v2[0] = float32(v1[0]+offset[0])/size + d
	v2[1] = float32(v1[1]+offset[1])/size + d
	return
}

func setUv(unit *model.CubeUnit, surface MultiQuad, size float32) {
	var d float32 = 0.001
	for i := range 6 {
		q := surface.quads[i]
		init := toFloat(q.pos, surface.pos, size, d)
		end := toFloat([2]int32{q.pos[0] + q.width, q.pos[1] + q.height}, surface.pos, size, -d)
		unit.SetUV(i, init, end)
	}
}

func draw() {
	if ctx.texture == nil {
		ctx.texture = texture.New(ctx.size, ctx.size)
	} else {
		ctx.texture.Clear()
	}
	if ctx.texture.Width != ctx.size || ctx.texture.Height != ctx.size {
		ctx.texture.Resize(ctx.size, ctx.size)
	}
	for _, s := range ctx.surfaces {
		for i := range 6 {
			q := s.quads[i]
			var pixels []texture.PixelEdit
			for x := range q.width {
				for y := range q.height {
					var p texture.PixelEdit
					p.Pos = getPositions([2]int32{q.pos[0] + x, q.pos[1] + y}, s.pos)
					if ok, color := q.texture.Get([2]int32{x, y}); ok {
						p.Color = color
						pixels = append(pixels, p)
					}
				}
			}
			ctx.texture.BulkSet(pixels)
		}
	}
	ctx.changed = true
}

func loadTexture() {
	for _, u := range ctx.model.Units {
		ctx.surfaces = append(ctx.surfaces, open(u))
	}
}

func generateTexture() {
	for _, u := range ctx.model.Units {
		found := false
		for i, s := range ctx.surfaces {
			if u.Id == s.id {
				if u.Resized {
					ctx.surfaces[i] = primitive(u)
				}
				found = true
				break
			}
		}
		if !found {
			ctx.surfaces = append(ctx.surfaces, primitive(u))
		}
	}
	ctx.size = packSurfaces(ctx.surfaces)
	toRemove := []int{}
	for i, s := range ctx.surfaces {
		var u *model.CubeUnit
		for _, unit := range ctx.model.Units {
			if unit.Id == s.id {
				u = unit
				break
			}
		}
		if u == nil {
			toRemove = append(toRemove, i)
			continue
		}
		setUv(u, *s, float32(ctx.size))
	}
	for _, i := range toRemove {
		ctx.surfaces = slices.Delete(ctx.surfaces, i, i+1)
	}
	draw()
}

func GetData() (width, height int32, data []float32) {
	return ctx.size, ctx.size, ctx.texture.FlatColors()
}

type TextureContext struct {
	name string

	model    *model.Model
	texture  *texture.Texture
	surfaces []*MultiQuad
	size     int32

	callbacks  []func()
	changed    bool
	subProject *textureDraw.SubProject
}

func OnChange(id string, callback func()) {
	ctxManager.Check(id)
	ctx.callbacks = append(ctx.callbacks, callback)
}

func call() {
	if ctx.changed {
		ctx.changed = false
		for _, f := range ctx.callbacks {
			f()
		}
		if ctx.subProject != nil {
			ctx.subProject.Add(ctx.texture)
		}
	}
}

func Show(id string) {
	ctxManager.Check(id)
	if ctx.model.Changed {
		generateTexture()
	}
	call()
	if im.Begin("Texture Manager") {
		if im.Button("Open Edit") {
			if ctx.subProject == nil {
				ctx.subProject = textureDraw.OpenSubProject(ctx.texture, ctx.name)
				ctx.subProject.OnChange(func(colors [][4]float32) {
					ctx.texture.Colors = colors
					ctx.changed = true

					for _, s := range ctx.surfaces {
						writeTexture(s, ctx.texture)
					}
				})
			} else {
				ctx.subProject.Focus()
			}
		}
	}
	im.End()
}

func Save(id string, w *file.ArchiveWriter) {
	ctxManager.Check(id)
	buff := bytes.NewBuffer(nil)
	ctx.texture.ToPNG(buff)
	w.Write("texture.png", buff.Bytes())
}

func Open(id string, r *file.ArchiveReader) error {
	ctxManager.Check(id)
	f, err := r.Open("texture.png")
	if err != nil {
		return err
	}
	ctx.texture, err = texture.DecodePNG(f)
	if err != nil {
		return err
	}
	ctx.size = ctx.texture.Width
	loadTexture()
	ctx.changed = true
	f.Close()
	return nil
}
