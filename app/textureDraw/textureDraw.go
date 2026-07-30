package textureDraw

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/history"
	"VulpesEditor/app/textureDraw/canvas"
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/tools"
	"VulpesEditor/app/util"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	im "github.com/AllenDang/cimgui-go/imgui"
)

func Init() {
	color.Init()
}

func AfterCreateContext() {
	// canvas.AddTexture(16, 16)
}

var standardTexSize = [2]int32{16, 16}
var textureSize = standardTexSize
var isNewTextureOpen = false
var isOpenTextureOpen = false
var nameInput string

func OpenNewTextureWindow() {
	isNewTextureOpen = true
}

func OpenOpenTextureWindow() {
	isOpenTextureOpen = true
}

func closeNewTextureWindow() {
	isNewTextureOpen = false
	textureSize = standardTexSize
	nameInput = ""
	im.CloseCurrentPopup()
}

func closeOpenTextureWindow() {
	isOpenTextureOpen = false
	im.CloseCurrentPopup()
}

func Show() {
	newTextureWindow()
	openTextureWindow()
}

func newTextureWindow() {
	if isNewTextureOpen {
		if !im.IsPopupOpenStr("New Texture") {
			im.OpenPopupStr("New Texture")
		}

		if im.BeginPopupModal("New Texture") {
			im.InputTextWithHint("name", "", &nameInput, im.InputTextFlagsNone, nil)
			im.DragInt2V("Texture Size", &textureSize, 1, 1, 1024, "%d", im.SliderFlagsNone)
			if im.Button("Create") {
				var c creationData
				if nameInput == "" {
					nameInput = "unnamed_texture"
				}
				c.name = strings.Clone(nameInput)
				c.width = uint32(textureSize[0])
				c.height = uint32(textureSize[1])
				openWindow(c)
				closeNewTextureWindow()
			}
			if im.Button("Cancel") {
				closeNewTextureWindow()
			}
			im.EndPopup()
		}
	}
}

func openTextureWindow() {
	if isOpenTextureOpen {
		if !im.IsPopupOpenStr("Open Texture") {
			im.OpenPopupStr("Open Texture")
		}

		if im.BeginPopupModal("Open Texture") {
			im.BeginDisabled()
			if im.Button("Open") {
				var c creationData
				if nameInput == "" {
					nameInput = "unnamed_texture"
				}
				c.name = strings.Clone(nameInput)
				c.width = uint32(textureSize[0])
				c.height = uint32(textureSize[1])
				openWindow(c)
				closeOpenTextureWindow()
			}
			im.EndDisabled()
			if im.Button("Cancel") {
				closeOpenTextureWindow()
			}
			im.EndPopup()
		}
	}
}

var count int32 = 0

type creationData struct {
	name   string
	width  uint32
	height uint32
}

type instance struct {
	name   string
	width  uint32
	height uint32

	id    int32
	focus bool
}

var Instances []*instance

func (s *instance) init() {
	history.New(s.id)
	color.New(s.id)
	tools.New(s.id, s.width, s.height)
}

func (s *instance) Focus() bool {
	if s.focus {
		s.focus = false
		return true
	}
	return false
}

func (s *instance) Name() string {
	return "Texture #" + strconv.FormatInt(int64(s.id), 10)
}

func (s *instance) Show() {
	history.Loop(s.id)
	color.Show(s.id)
	tools.Show(s.id)
	canvas.Show(s.id)
}

func (s *instance) Save() {
	w, err := file.NewArchive(filepath.Join(util.AppDir, "projects", "textures"), s.name)
	if err != nil {
		fmt.Println(err)
		return
	}
	canvas.Save(w)
	w.Save()
}

func OpenTexture(name string) {
	r, err := file.Load(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	itc := new(instance)
	itc.id = count
	itc.focus = true
	if err := canvas.Open(itc.id, r); err != nil {
		fmt.Println(err)
		return
	}
	itc.width, itc.height = canvas.Size()
	itc.init()
	Instances = append(Instances, itc)
	count += 1
}

func openWindow(c creationData) {
	itc := new(instance)
	itc.name = c.name
	itc.width = c.width
	itc.height = c.height
	itc.id = count
	itc.focus = true
	itc.init()
	canvas.New(itc.id, itc.width, itc.height)
	Instances = append(Instances, itc)
	count += 1
}
