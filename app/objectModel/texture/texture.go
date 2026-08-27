package texture

import (
	"VulpesEditor/app/objectModel/model"
	"VulpesEditor/app/textureDraw"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"bytes"
	"fmt"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type quad struct {
	pos    [2]int32
	width  int32
	height int32
}

type MultiQuad struct {
	id        int32
	pos       [2]int32
	width     int32
	height    int32
	cutWidth  int32
	cutHeight int32
	quads     [6]quad
}

func primitive(c *model.CubeUnit) (unity *MultiQuad) {
	unity = new(MultiQuad)
	_, size := c.Data()

	x := int32(size[0] * 16)
	y := int32(size[1] * 16)
	z := int32(size[2] * 16)

	unity.quads[0] = quad{[2]int32{0, 0}, x, z}
	unity.quads[1] = quad{[2]int32{x, 0}, x, z}
	unity.quads[2] = quad{[2]int32{0, z}, x, y}
	unity.quads[3] = quad{[2]int32{x, z}, x, y}
	unity.quads[4] = quad{[2]int32{2 * x, 0}, y, z}
	unity.quads[5] = quad{[2]int32{2*x + y, 0}, y, z}

	unity.id = c.Id
	unity.pos = [2]int32{0, 0}
	unity.width = 2 * (x + y)
	unity.height = z + y
	unity.cutWidth = 2 * x
	unity.cutHeight = z

	return
}

func getPositions(v1, offset [2]int32) (v2 [2]int32) {
	v2[0] = v1[0] + offset[0]
	v2[1] = v1[1] + offset[1]
	return
}

func toFloat(v1, offset [2]int32, size float32) (v2 [2]float32) {
	v2[0] = float32(v1[0]+offset[0]) / size
	v2[1] = float32(v1[1]+offset[1]) / size
	return
}

func setUv(unit *model.CubeUnit, surface MultiQuad, size float32) {
	for i := range 6 {
		q := surface.quads[i]
		init := toFloat(q.pos, surface.pos, size)
		end := toFloat([2]int32{q.pos[0] + q.width, q.pos[1] + q.height}, surface.pos, size)
		unit.SetUV(i, init, end)
	}
}

func draw(surfaces []*MultiQuad, size uint32) {
	if ctx.texture == nil {
		ctx.texture = texture.New(size, size)
	} else {
		ctx.texture.Clear()
	}
	if ctx.texture.Width != size || ctx.texture.Height != size {
		ctx.texture.Resize(size, size)
	}
	colors := [12][4]float32{
		{0.25, 0.25, 0.75, 1},
		{0.5, 0.5, 1, 1},
		{0.75, 0.75, 0.25, 1},
		{1, 1, 0.5, 1},
		{0.75, 0.75, 0.75, 1},
		{1, 1, 1, 1},
		{0.25, 0.25, 0.25, 1},
		{0.5, 0.5, 0.5, 1},
		{0.25, 0.75, 0.25, 1},
		{0.5, 1, 0.5, 1},
		{0.75, 0.25, 0.25, 1},
		{1, 0.5, 0.5, 1},
	}
	for _, s := range surfaces {
		for i := range 6 {
			q := s.quads[i]
			init := getPositions(q.pos, s.pos)
			end := getPositions([2]int32{q.pos[0] + q.width, q.pos[1] + q.height}, s.pos)
			var pixels []texture.PixelEdit
			for x := init[0]; x < end[0]; x++ {
				for y := init[1]; y < end[1]; y++ {
					var p texture.PixelEdit
					p.Pos = [2]int32{x, y}
					if x == init[0] || x == end[0]-1 || y == init[1] || y == end[1]-1 {
						p.Color = colors[2*i+1]
					} else {
						p.Color = colors[2*i]
					}
					pixels = append(pixels, p)
				}
			}
			ctx.texture.BulkSet(pixels)
		}
	}
	ctx.changed = true
}

func generateTexture() {
	all := []*MultiQuad{}
	for _, u := range ctx.model.Units {
		all = append(all, primitive(u))
	}
	ctx.size = packSurfaces(all)
	for _, s := range all {
		var u *model.CubeUnit
		for _, unit := range ctx.model.Units {
			if unit.Id == s.id {
				u = unit
				break
			}
		}
		if u == nil {
			fmt.Printf("Unit %d not found.", s.id)
			continue
		}
		setUv(u, *s, float32(ctx.size))
	}
	draw(all, uint32(ctx.size))
}

func GetData() (width, height int32, data []float32) {
	return ctx.size, ctx.size, ctx.texture.FlatColors()
}

type TextureContext struct {
	model     *model.Model
	texture   *texture.Texture
	callbacks []func()
	changed   bool
	size      int32
}

func OnChange(callback func()) {
	ctx.callbacks = append(ctx.callbacks, callback)
}

func call() {
	if ctx.changed {
		ctx.changed = false
		for _, f := range ctx.callbacks {
			f()
		}
	}
}

func Show(id int32) {
	ctxManager.Check(id)
	if ctx.model.Changed() {
		generateTexture()
	}
	call()
	if im.Begin("Texture Manager") {
		if im.Button("Open Edit") {
			b := bytes.Buffer{}
			ctx.texture.ToPNG(&b)
			textureDraw.OpenImage(&b, "Model Test")
		}
	}
	im.End()
}
