package textureDraw

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/history"
	"VulpesEditor/app/textureDraw/canvas"
	"VulpesEditor/app/textureDraw/color"
	"VulpesEditor/app/textureDraw/tools"
	"VulpesEditor/app/util"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	im "github.com/AllenDang/cimgui-go/imgui"
)

var AllTextures []file.Project

func Init() {
	color.Init()
	AllTextures = file.GetAllProjects("textures")
}

func AfterCreateContext() {
	// canvas.AddTexture(16, 16)
}

var standardTexSize = [2]int32{16, 16}
var textureSize = standardTexSize
var isNewTextureOpen = false
var nameInput string

func OpenNewTextureWindow() {
	isNewTextureOpen = true
}

func closeNewTextureWindow() {
	isNewTextureOpen = false
	textureSize = standardTexSize
	nameInput = ""
	im.CloseCurrentPopup()
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
				openNew(c)
				closeNewTextureWindow()
			}
			im.SameLine()
			if im.Button("Cancel") {
				closeNewTextureWindow()
			}
			im.EndPopup()
		}
	}
}

var selected int = -1
var isOpenTextureOpen = false

func OpenOpenTextureWindow() {
	isOpenTextureOpen = true
}

func closeOpenTextureWindow() {
	isOpenTextureOpen = false
	selected = -1
	im.CloseCurrentPopup()
}

func openTextureWindow() {
	if isOpenTextureOpen {
		if !im.IsPopupOpenStr("Open Texture") {
			im.OpenPopupStr("Open Texture")
		}

		if im.BeginPopupModal("Open Texture") {
			im.BeginListBox("Select")
			for i, p := range AllTextures {
				if im.SelectableBoolV(p.Name, i == selected, im.SelectableFlagsAllowDoubleClick, im.NewVec2(0, 0)) {
					if im.IsMouseDoubleClicked(0) {
						OpenTexture(p.Path)
						closeOpenTextureWindow()
					}
					selected = i
				}
			}
			im.EndListBox()
			dis := selected == -1
			if dis {
				im.BeginDisabled()
			}
			if im.Button("Open") {
				OpenTexture(AllTextures[selected].Path)
				closeOpenTextureWindow()
			}
			if dis {
				im.EndDisabled()
			}
			im.SameLine()
			if im.Button("Cancel") {
				closeOpenTextureWindow()
			}
			im.EndPopup()
		}
	}
}

func Show() {
	newTextureWindow()
	openTextureWindow()
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
	b := strings.Builder{}
	b.WriteString(s.name)
	b.WriteRune('\n')
	b.WriteString(strconv.FormatInt(int64(s.width), 10))
	b.WriteRune('\n')
	b.WriteString(strconv.FormatInt(int64(s.height), 10))
	w.Write("metaData.txt", []byte(b.String()))
	w.Save()
}

func OpenTexture(path string) {
	r, err := file.Load(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := r.Open("metaData.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	b := strings.Builder{}
	io.Copy(&b, f)
	file := b.String()
	f.Close()
	field := strings.Split(file, "\n")
	if len(field) < 3 {
		fmt.Println("Incomplete data")
	}
	itc := new(instance)
	itc.name = field[0]
	width, err := strconv.ParseInt(field[1], 10, 32)
	if err != nil {
		fmt.Println("can't parse width")
	}
	height, err := strconv.ParseInt(field[2], 10, 32)
	if err != nil {
		fmt.Println("can't parse height")
	}
	itc.width = uint32(width)
	itc.height = uint32(height)
	itc.id = count
	itc.focus = true
	if err := canvas.Open(itc.id, r); err != nil {
		fmt.Println(err)
		return
	}
	itc.init()
	Instances = append(Instances, itc)
	count += 1
}

func openNew(c creationData) {
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
