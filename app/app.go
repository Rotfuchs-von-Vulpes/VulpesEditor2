package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/front/tabs"
	"VulpesEditor/app/textureDraw"
	"VulpesEditor/app/util"
	"strconv"

	im "github.com/AllenDang/cimgui-go/imgui"
)

func Init() {
	util.Init()
	textureDraw.Init()
}

func AfterCreateContext() {
	front.Init()
}

func BeforeDestroyContext() {
	front.Nuke()
}

var first = true
var save = false
var close = false

func Loop() {
	im.ClearSizeCallbackPool()

	workerAreaFlags := im.WindowFlagsNoTitleBar |
		im.WindowFlagsNoCollapse |
		im.WindowFlagsNoDecoration |
		im.WindowFlagsNoResize |
		im.WindowFlagsNoBringToFrontOnFocus |
		im.WindowFlagsNoMove |
		im.WindowFlagsMenuBar

	im.SetNextWindowPos(im.MainViewport().WorkPos())
	im.SetNextWindowSize(im.MainViewport().WorkSize())

	if im.BeginV("Work Area", nil, workerAreaFlags) {
		if im.BeginMenuBar() {
			if im.BeginMenu("File") {
				if im.BeginMenu("New") {
					if im.MenuItemBool("Texture") {
						textureDraw.OpenNewTextureWindow()
					}
					im.EndMenu()
				}
				if im.BeginMenu("Open") {
					if im.MenuItemBool("Texture") {
						textureDraw.OpenOpenTextureWindow()
					}
					im.EndMenu()
				}
				if im.BeginMenu("Open Recent") {
					for _, project := range textureDraw.AllTextures {
						if im.MenuItemBool(project.Name) {
							textureDraw.OpenTexture(project.Path)
						}
					}
					im.EndMenu()
				}
				if im.MenuItemBool("Save") {
					save = true
				}
				if im.MenuItemBool("Close") {
					close = true
				}
				im.EndMenu()
			}
			im.EndMenuBar()
		}

		dockspaceId := im.IDStr("Dockspace")
		if im.BeginTabBar("AAA") {
			if im.BeginTabItem("Home") {
				childSize := im.NewVec2(500, 100+im.FrameHeight()*16)
				im.SetCursorPos(im.WindowSize().Sub(childSize).Div(2))
				im.BeginChildStrV("Contents", childSize, im.ChildFlagsNone, im.WindowFlagsNone)
				size := im.NewVec2(-0.00001, 0)
				if im.ButtonV("New Texture...", size) {
					textureDraw.OpenNewTextureWindow()
				}
				if im.ButtonV("Open Texture...", size) {
					textureDraw.OpenOpenTextureWindow()
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

			for i, t := range tabs.AllTabs {
				if t == nil {
					continue
				}
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

					if save {
						t.Save()
					}
					if close {
						tabs.Close(i)
					}
				}
				im.PopID()
			}

			im.EndTabBar()
		}
	}
	im.End()

	textureDraw.Show()

	save = false
	close = false
}
