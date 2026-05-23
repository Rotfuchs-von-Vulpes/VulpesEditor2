package front

import (
	"VulpesEditor/app/front/renderer"

	"github.com/AllenDang/cimgui-go/imgui"
)

var (
	dockID         imgui.ID
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

var first = true

func Loop() {
	imgui.ClearSizeCallbackPool()
	dockID = imgui.IDStr("My Dockspace")
	if first {
		first = false
		// imgui.InternalDockBuilderAddNodeV(dockID, imgui.DockNodeFlags(imgui.DockNodeFlagsDockSpace))
		// imgui.InternalDockBuilderSetNodeSize(dockID, imgui.MainViewport().Size())
		// var down imgui.ID = 0
		// var main imgui.ID = dockID
		// imgui.InternalDockBuilderSplitNode(main, imgui.DirLeft, 0.5, &down, &main)
		// imgui.InternalDockBuilderDockWindow("Window 1", main)
		// imgui.InternalDockBuilderDockWindow("Image", down)
		// imgui.InternalDockBuilderFinish(dockID)
	}
	imgui.DockSpaceOverViewportV(dockID, imgui.MainViewport(), imgui.DockNodeFlagsNone, imgui.NewEmptyWindowClass())

	ShowWidgetsDemo()
}

func ShowWidgetsDemo() {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.Begin("Window 1")

	if imgui.ButtonV("Click Me", imgui.NewVec2(80, 20)) {

	}
	imgui.TextUnformatted("Unformatted text")
	imgui.Checkbox("Show demo window", &showDemoWindow)
	if imgui.BeginCombo("Combo", "Combo preview") {
		imgui.SelectableBoolPtr("Item 1", &selected)
		imgui.SelectableBool("Item 2")
		imgui.SelectableBool("Item 3")
		imgui.EndCombo()
	}

	if imgui.RadioButtonBool("Radio button1", selected) {
		selected = true
	}

	imgui.SameLine()

	if imgui.RadioButtonBool("Radio button2", !selected) {
		selected = false
	}

	imgui.InputTextWithHint("Name", "write your name here", &content, 0, nil)
	imgui.Text(content)
	imgui.SliderInt("Slider int", &value3, 0, 100)
	imgui.DragInt("Drag int", &value1)
	imgui.DragInt2("Drag int2", &values)
	value1 = values[0]
	imgui.ColorEdit4("Color Edit3", &color4)

	imgui.End()
}

func Nuke() {
	renderer.Nuke()
}
