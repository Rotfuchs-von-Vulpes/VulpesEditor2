package blocks

import "VulpesEditor/app/context"

var ctx *BlocksContext

func (c *BlocksContext) Use() {
	ctx = c
}

var ctxManager = context.New()

func New(id string) {
	ctx := new(BlocksContext)

	ctx.textures[0] = new(texButton)
	ctx.textures[0].label = "Top"
	ctx.textures[1] = new(texButton)
	ctx.textures[1].label = "Bottom"
	ctx.textures[2] = new(texButton)
	ctx.textures[2].label = "Front"
	ctx.textures[3] = new(texButton)
	ctx.textures[3].label = "Back"
	ctx.textures[4] = new(texButton)
	ctx.textures[4].label = "Left"
	ctx.textures[5] = new(texButton)
	ctx.textures[5].label = "Right"

	ctxManager.Add(id, ctx)
}
