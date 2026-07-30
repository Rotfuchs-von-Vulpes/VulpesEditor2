package canvas

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/textureDraw/canvas/texture"
	"VulpesEditor/app/textureDraw/canvas/textureEdit"
)

func (s *TextureContext) Use() {
	ctx = s
}

var ctx *TextureContext
var ctxManager = context.New()

func New(id int32, w, h uint32) {
	OpenTexture(id, textureEdit.New(texture.New(w, h)))
}

func OpenTexture(id int32, tex *textureEdit.TextureEdit) {
	ctx = new(TextureContext)
	ctx.zoom = 0.9
	ctx.textureViewer = renderer.CreateFramebuffer(500, 500)
	viwerSize = [2]float32{500, 500}
	ctx.texture = tex
	ctxManager.Add(id, ctx)
}

func Save(w *file.ArchiveWriter) {
	ctx.texture.Save(w)
}

func Open(id int32, r *file.ArchiveReader) (err error) {
	tex, err := textureEdit.Open(r)
	if err != nil {
		return err
	}
	OpenTexture(id, tex)
	return
}
