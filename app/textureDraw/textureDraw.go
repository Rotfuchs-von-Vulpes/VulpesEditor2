package textureDraw

import (
	"VulpesEditor/app/textureDraw/palette"
	"VulpesEditor/app/textureDraw/texture"
	"VulpesEditor/app/textureDraw/tools"
)

func Init() {
	texture.Init()
	palette.Init()
	tools.Init()
}

func AfterCreateContext() {
	texture.AddTexture(16, 16)
}

func Loop() {
	palette.Loop()
	tools.Show()
	texture.SetColors(palette.SelectedColors())
	for _, c := range texture.AllCtx {
		c.Show()
	}
}
