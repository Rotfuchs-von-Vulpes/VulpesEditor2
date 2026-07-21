package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/textureDraw"
	"strconv"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type Tab interface {
	Name() string
	Show()
	Focus() bool
}

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
	var allTabs []Tab
	for _, itc := range textureDraw.Instances {
		allTabs = append(allTabs, itc)
	}

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
				childSize := im.NewVec2(200, 100+im.FrameHeight()*16)
				im.SetCursorPos(im.WindowSize().Sub(childSize).Div(2))
				im.BeginChildStrV("Contents", childSize, im.ChildFlagsNone, im.WindowFlagsNone)
				size := im.NewVec2(-0.00001, 0)
				if im.ButtonV("New Texture...", size) {
					textureDraw.OpenNewTextureWindow()
				}
				if im.ButtonV("Open Texture...", size) {
					im.OpenPopupStr("Not Implement")
				}

				front.NotImplementPopUp()

				im.EndChild()
				im.EndTabItem()
			}

			if im.BeginTabItem("Debug") {
				im.DockSpace(dockspaceId)
				front.Loop()
				im.EndTabItem()
			}

			for i, t := range allTabs {
				f := im.TabItemFlagsNone
				if t.Focus() {
					f |= im.TabItemFlagsSetSelected
				}
				tabID := t.Name() + "###" + strconv.FormatInt(int64(i), 10)
				im.PushIDStr(tabID)
				if im.BeginTabItemV(t.Name(), nil, f) {
					im.DockSpace(dockspaceId)
					t.Show()
					im.EndTabItem()
				}
				im.PopID()
			}

			im.EndTabBar()
		}
	}
	im.End()

	textureDraw.ShowNewTextureWindow()
}
