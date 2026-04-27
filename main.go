package main

import (
	"runtime"

	"VulpesEditor/app"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
)

var currentBackend backend.Backend[sdlbackend.SDLWindowFlags]

func main() {
	if !app.Init() {
		return
	}
	runtime.LockOSThread()
	currentBackend, _ = backend.CreateBackend(sdlbackend.NewSDLBackend())
	currentBackend.SetAfterCreateContextHook(app.AfterCreateContext)
	currentBackend.SetBeforeDestroyContextHook(app.BeforeDestroyContext)

	currentBackend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))

	currentBackend.CreateWindow("Vulpes Editor", 1200, 900)

	currentBackend.Run(app.Loop)
}
