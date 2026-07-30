package app

import (
	"VulpesEditor/app/front"
	"VulpesEditor/app/textureDraw"
	"VulpesEditor/app/util"
	"os"
	"path/filepath"
	"strconv"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type Tab interface {
	Name() string
	Show()
	Focus() bool
	Save()
}

type project struct {
	name string
	path string
}

var allTextures []project

func Init() {
	util.Init()
	textureDraw.Init()

	projectsDir := filepath.Join(util.AppDir, "projects", "textures")

	if files, err := os.ReadDir(projectsDir); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				var p project
				p.name = file.Name()
				p.path = filepath.Join(projectsDir, file.Name())
				allTextures = append(allTextures, p)
			}
		}
	}
}

func AfterCreateContext() {
	front.Init()
	textureDraw.AfterCreateContext()
}

func BeforeDestroyContext() {
	front.Nuke()
}

var first = true
var save = false

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
					for _, project := range allTextures {
						if im.MenuItemBool(project.name) {
							textureDraw.OpenTexture(project.path)
						}
					}
					im.EndMenu()
				}
				if im.MenuItemBool("Save") {
					save = true
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

					if save {
						t.Save()
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
}
