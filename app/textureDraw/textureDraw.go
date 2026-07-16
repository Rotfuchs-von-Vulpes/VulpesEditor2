package textureDraw

import (
	"VulpesEditor/app/textureDraw/canvas"
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/history"
	"VulpesEditor/app/textureDraw/tools"
)

func Init() {
	color.Init()
	tools.Init()
}

func AfterCreateContext() {
	canvas.AddTexture(16, 16)
}

func Loop() {
	color.Loop()
	tools.Show()
	for _, c := range canvas.AllCtx {
		c.Show()
	}
	canvas.ShowLayers()
	history.Loop()
}
