package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/textureDraw"
)

func Init() bool {
	textureDraw.Init()
	return true
}

func AfterCreateContext() {
	front.Init()
	textureDraw.AddTexture(16, 16)
}

func BeforeDestroyContext() {
	front.Nuke()
}

var first = true

func Loop() {
	front.Loop()
	textureDraw.Loop()
}
