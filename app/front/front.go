package front

import (
	"VulpesEditor/app/front/renderer"

	im "github.com/AllenDang/cimgui-go/imgui"
)

var showDemoWindow bool

func Init() {
	renderer.Init()
}

func Nuke() {
	renderer.Nuke()
}

func Loop() {
	if showDemoWindow {
		im.ShowDemoWindowV(&showDemoWindow)
	}
	im.Begin("Debug")
	im.Checkbox("Show demo window", &showDemoWindow)
	im.End()
}

func NotImplementPopUp() {
	if im.BeginPopupModal("Not Implement") {
		im.Text("Not implement yet!")
		if im.Button("Ok") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}
}
