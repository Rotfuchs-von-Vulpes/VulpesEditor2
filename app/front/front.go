package front

import (
	"VulpesEditor/app/front/renderer"

	im "github.com/AllenDang/cimgui-go/imgui"
)

var (
	dockID         im.ID
	showDemoWindow bool
	value1         int32
	value2         int32
	value3         int32
	values         [2]int32 = [2]int32{value1, value2}
	content        string   = "Let me try"
	r              float32
	g              float32
	b              float32
	a              float32
	color4         [4]float32 = [4]float32{r, g, b, a}
	selected       bool
)

var spectrumRange float32 = 1
var intensityRange float32 = 0
var sigmaRange float32 = 1.5

func Init() {
	renderer.Init()
}

func Nuke() {
	renderer.Nuke()
}

var first = true

func Loop() {
	if showDemoWindow {
		im.ShowDemoWindowV(&showDemoWindow)
	}
	im.Begin("Debug")
	im.Checkbox("Show demo window", &showDemoWindow)
	im.End()
}
