package animation

import (
	"VulpesEditor/app/context"
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
	"uuid"
)

var ctx *AnimationContext
var ctxManager = context.New()

func (s *AnimationContext) Use() {
	ctx = s
}

func New(id string, md *model.Model, ms *renderer.Mesh) {
	ctx := new(AnimationContext)
	ctx.model = md
	ctx.model.Inscribe("animation")
	ctx.mesh = ms
	ctx.animation = new(animation)
	ctx.animation.bones = map[string]*boneAnimation{}
	ctx.animation.id = uuid.New().String()
	ctxManager.Add(id, ctx)
}
