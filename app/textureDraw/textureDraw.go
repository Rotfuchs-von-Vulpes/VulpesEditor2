package textureDraw

import (
	"VulpesEditor/app/textureDraw/palette"
	"VulpesEditor/app/textureDraw/texture"
)

func Init() {
	texture.Init()
	palette.Init()
}

func AfterCreateContext() {
	texture.AddTexture(16, 16)
}

func Loop() {
	palette.Loop()
	texture.SetColors(palette.SelectedColors())
	for _, c := range texture.AllCtx {
		c.Show()
	}
}
