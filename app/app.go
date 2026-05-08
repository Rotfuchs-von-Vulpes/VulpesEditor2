package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/textureDraw"
)

func Init() {
	textureDraw.Init()
}

func AfterCreateContext() {
	front.Init()
	textureDraw.AfterCreateContext()
}

func BeforeDestroyContext() {
	front.Nuke()
}

var first = true

func Loop() {
	front.Loop()
	textureDraw.Loop()
}
