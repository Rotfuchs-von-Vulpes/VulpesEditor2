package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/textureDraw"

	im "github.com/AllenDang/cimgui-go/imgui"
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
	im.ClearSizeCallbackPool()

	workerAreaFlags := im.WindowFlagsNoTitleBar |
		im.WindowFlagsNoCollapse |
		im.WindowFlagsNoDecoration |
		im.WindowFlagsNoResize |
		im.WindowFlagsNoBringToFrontOnFocus |
		im.WindowFlagsNoMove

	im.SetNextWindowPos(im.MainViewport().WorkPos())
	im.SetNextWindowSize(im.MainViewport().WorkSize())

	if im.BeginV("Work Area", nil, workerAreaFlags) {
		dockspaceId := im.IDStr("Dockspace")
		if im.BeginTabBar("AAA") {
			if im.BeginTabItem("Home") {
				// var childSize = new Vector2(RainedLogo.Width, RainedLogo.Height - 100f + ImGui.GetFrameHeight() * 16f);
				im.BeginChildStr("Contents")
				im.Text("Teste")
				im.EndChild()
				im.EndTabItem()
			}
			if im.BeginTabItem("Debug") {
				im.DockSpace(dockspaceId)
				front.Loop()
				im.EndTabItem()
			}
			if im.BeginTabItem("Texture Draw") {
				im.DockSpace(dockspaceId)
				textureDraw.Loop()
				im.EndTabItem()
			}
			im.EndTabBar()
		}
	}
	im.End()
}
