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

var first = true

func Loop() {
	im.ClearSizeCallbackPool()
	dockID = im.IDStr("My Dockspace")
	if first {
		first = false
		// im.InternalDockBuilderAddNodeV(dockID, im.DockNodeFlags(im.DockNodeFlagsDockSpace))
		// im.InternalDockBuilderSetNodeSize(dockID, im.MainViewport().Size())
		// var down im.ID = 0
		// var main im.ID = dockID
		// im.InternalDockBuilderSplitNode(main, im.DirLeft, 0.5, &down, &main)
		// im.InternalDockBuilderDockWindow("Window 1", main)
		// im.InternalDockBuilderDockWindow("Image", down)
		// im.InternalDockBuilderFinish(dockID)
	}
	im.DockSpaceOverViewportV(dockID, im.MainViewport(), im.DockNodeFlagsNone, im.NewEmptyWindowClass())

	ShowWidgetsDemo()
}

func ShowWidgetsDemo() {
	if showDemoWindow {
		im.ShowDemoWindowV(&showDemoWindow)
	}

	im.Begin("Window 1")

	if im.ButtonV("Click Me", im.NewVec2(80, 20)) {

	}
	im.TextUnformatted("Unformatted text")
	im.Checkbox("Show demo window", &showDemoWindow)
	if im.BeginCombo("Combo", "Combo preview") {
		im.SelectableBoolPtr("Item 1", &selected)
		im.SelectableBool("Item 2")
		im.SelectableBool("Item 3")
		im.EndCombo()
	}

	if im.RadioButtonBool("Radio button1", selected) {
		selected = true
	}

	im.SameLine()

	if im.RadioButtonBool("Radio button2", !selected) {
		selected = false
	}

	im.InputTextWithHint("Name", "write your name here", &content, 0, nil)
	im.Text(content)
	im.SliderInt("Slider int", &value3, 0, 100)
	im.DragInt("Drag int", &value1)
	im.DragInt2("Drag int2", &values)
	value1 = values[0]
	im.ColorEdit4("Color Edit3", &color4)

	im.End()
}

func Nuke() {
	renderer.Nuke()
}
