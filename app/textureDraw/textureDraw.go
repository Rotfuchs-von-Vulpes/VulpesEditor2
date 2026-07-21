package textureDraw

import (
	"VulpesEditor/app/textureDraw/canvas"
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/history"
	"VulpesEditor/app/textureDraw/tools"
	"strconv"

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
var isOpen = false

func OpenNewTextureWindow() {
	isOpen = true
}

func CloseNewTextureWindow() {
	isOpen = false
	textureSize = standardTexSize
	im.CloseCurrentPopup()
}

func ShowNewTextureWindow() {
	if isOpen {
		if !im.IsPopupOpenStr("New Texture") {
			im.OpenPopupStr("New Texture")
		}

		if im.BeginPopupModal("New Texture") {
			im.DragInt2V("Texture Size", &textureSize, 1, 1, 1024, "%d", im.SliderFlagsNone)
			if im.Button("Create") {
				var c creationData
				c.width = uint32(textureSize[0])
				c.height = uint32(textureSize[1])
				open(c)
				CloseNewTextureWindow()
			}
			if im.Button("Cancel") {
				CloseNewTextureWindow()
			}
			im.EndPopup()
		}
	}
}

var count int32 = 0

type creationData struct {
	width  uint32
	height uint32
}

type instance struct {
	id    int32
	focus bool
}

var Instances []*instance

func (s *instance) init(c creationData) {
	color.New(s.id)
	tools.New(s.id, c.width, c.height)
	history.New(s.id)
	canvas.New(s.id, c.width, c.height)
}

func (s *instance) begin() {
	color.Begin(s.id)
	tools.Begin(s.id)
	history.Begin(s.id)
	canvas.Begin(s.id)
}

func (s *instance) end() {
	color.End()
	tools.End()
	history.End()
	canvas.End()
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
	s.begin()
	color.Loop()
	tools.Show()
	canvas.Show()
	canvas.ShowLayers()
	history.Loop()
	s.end()
}

func open(c creationData) {
	itc := new(instance)
	itc.id = count
	itc.focus = true
	itc.init(c)
	Instances = append(Instances, itc)
	count += 1
}
