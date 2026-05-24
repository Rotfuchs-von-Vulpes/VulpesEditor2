package textureDraw

import (
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/texture"
	"VulpesEditor/app/textureDraw/tools"
)

func Init() {
	texture.Init()
	color.Init()
	tools.Init()
}

func AfterCreateContext() {
	texture.AddTexture(16, 16)
}

func Loop() {
	color.Loop()
	tools.Show()
	texture.SetColors(color.SelectedColors())
	for _, c := range texture.AllCtx {
		c.Show()
	}
}
