package textureDraw

import "VulpesEditor/app/textureDraw/texture"

func Init() {
	texture.Init()
}

func AfterCreateContext() {
	texture.AddTexture(16, 16)
}

func Loop() {
	for _, c := range texture.AllCtx {
		c.Show()
	}
}
