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
			im.DragInt2V("Texture Size", &textureSize, 1, 1, 1024, "", im.SliderFlagsNone)
			if im.Button("Create") {
				open(canvas.AddTexture(uint32(textureSize[0]), uint32(textureSize[1])))
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

type instance struct {
	id  int32
	ctx *canvas.TextureContext
}

var Instances []*instance

func (s *instance) Name() string {
	return "Texture #" + strconv.FormatInt(int64(s.id), 10)
}

func (s *instance) init() {
	color.New(s.id)
	tools.New(s.id)
	history.New(s.id)
}

func (s *instance) begin() {
	color.Begin(s.id)
	tools.Begin(s.id)
	history.Begin(s.id)
}

func (s *instance) end() {
	color.End()
	tools.End()
	history.End()
}

func (s *instance) Show() {
	s.begin()
	color.Loop()
	tools.Show()
	s.ctx.Show()
	s.ctx.ShowLayers()
	history.Loop()
	s.end()
}

func open(ctx *canvas.TextureContext) {
	itc := new(instance)
	itc.id = count
	itc.ctx = ctx
	itc.init()
	Instances = append(Instances, itc)
	count += 1
}
