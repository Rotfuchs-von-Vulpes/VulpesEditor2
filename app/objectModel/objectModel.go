package objectModel

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/tabs"
	"VulpesEditor/app/history"
	"VulpesEditor/app/objectModel/model"
	"VulpesEditor/app/objectModel/texture"
	"VulpesEditor/app/objectModel/tools"
	"VulpesEditor/app/objectModel/view"
	"VulpesEditor/app/textureDraw/canvas"
	"VulpesEditor/app/util"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	im "github.com/AllenDang/cimgui-go/imgui"
)

var AllModels []file.Project

func Init() {
	AllModels = file.GetAllProjects("models")
}

type creationData struct {
	name string
}

var isNewModelOpen = false
var nameInput string

func OpenNewModelWindow() {
	isNewModelOpen = true
}

func closeNewModelWindow() {
	isNewModelOpen = false
	nameInput = ""
	im.CloseCurrentPopup()
}

func newModelWindow() {
	if isNewModelOpen {
		if !im.IsPopupOpenStr("New Model") {
			im.OpenPopupStr("New Model")
		}

		if im.BeginPopupModal("New Model") {
			im.InputTextWithHint("name", "", &nameInput, im.InputTextFlagsNone, nil)
			if im.Button("Create") {
				var c creationData
				if nameInput == "" {
					nameInput = "unnamed_model"
				}
				c.name = strings.Clone(nameInput)
				openNew(c)
				closeNewModelWindow()
			}
			im.SameLine()
			if im.Button("Cancel") {
				closeNewModelWindow()
			}
			im.EndPopup()
		}
	}
}

var selected int = -1
var isOpenModelOpen = false

func OpenOpenModelWindow() {
	isOpenModelOpen = true
}

func closeOpenModelWindow() {
	isOpenModelOpen = false
	selected = -1
	im.CloseCurrentPopup()
}

func openModelWindow() {
	if isOpenModelOpen {
		if !im.IsPopupOpenStr("Open Model") {
			im.OpenPopupStr("Open Model")
		}

		if im.BeginPopupModal("Open Model") {
			im.BeginListBox("Select")
			for i, p := range AllModels {
				if im.SelectableBoolV(p.Name, i == selected, im.SelectableFlagsAllowDoubleClick, im.NewVec2(0, 0)) {
					if im.IsMouseDoubleClicked(0) {
						OpenModel(p.Path)
						closeOpenModelWindow()
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
				OpenModel(AllModels[selected].Path)
				closeOpenModelWindow()
			}
			if dis {
				im.EndDisabled()
			}
			im.SameLine()
			if im.Button("Cancel") {
				closeOpenModelWindow()
			}
			im.EndPopup()
		}
	}
}

func Show() {
	newModelWindow()
	openModelWindow()
}

var IdSys *util.IdSystem

type instance struct {
	name string

	model *model.Model

	id    int32
	focus bool
}

func (s *instance) init() {
	history.New(s.id)
	view.New(s.id)
	s.model = view.Model()
	tools.New(s.id, s.model)
	texture.New(s.id, s.model)
	s.name = "Model #" + strconv.FormatInt(int64(s.id), 10)
}

func (s *instance) Focus() bool {
	if s.focus {
		s.focus = false
		return true
	}
	return false
}

func (s *instance) Name() string {
	return s.name
}

func (s *instance) Show() {
	history.Loop(s.id)
	tools.Show(s.id)
	texture.Show(s.id)
	view.Show(s.id)
	s.model.Reset()
}

func (s *instance) Save() {
	w, err := file.NewArchive(filepath.Join(util.AppDir, "projects", "models"), s.name)
	if err != nil {
		fmt.Println(err)
		return
	}
	// canvas.Save(w)
	b := strings.Builder{}
	b.WriteString("model")
	b.WriteRune('\n')
	b.WriteString(s.name)
	w.Write("metaData.txt", []byte(b.String()))
	w.Save()
}

func OpenModel(path string) {
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
	if len(field) < 2 {
		fmt.Println("Incomplete data")
		return
	}
	if field[0] != "model" {
		fmt.Println("Wrong project type")
		return
	}
	itc := new(instance)
	itc.name = field[1]
	itc.id = IdSys.GetID()
	itc.focus = true
	if err := canvas.Open(itc.id, r); err != nil {
		fmt.Println(err)
		return
	}
	itc.init()
	tabs.Push(itc)
}

func openNew(c creationData) {
	itc := new(instance)
	itc.name = c.name
	itc.id = IdSys.GetID()
	itc.focus = true
	itc.init()
	tabs.Push(itc)
}
