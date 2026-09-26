package blocks

import "VulpesEditor/app/context"

var ctx *BlocksContext

func (c *BlocksContext) Use() {
	ctx = c
}

var ctxManager = context.New()

func New(id string) {
	ctx := new(BlocksContext)

	ctx.top = new(texButton)
	ctx.top.label = "Top"
	ctx.bottom = new(texButton)
	ctx.bottom.label = "Bottom"
	ctx.front = new(texButton)
	ctx.front.label = "Front"
	ctx.back = new(texButton)
	ctx.back.label = "Back"
	ctx.left = new(texButton)
	ctx.left.label = "Left"
	ctx.right = new(texButton)
	ctx.right.label = "Right"

	ctxManager.Add(id, ctx)
}
