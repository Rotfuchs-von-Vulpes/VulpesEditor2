package objectModel

import (
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/tabs"
	"VulpesEditor/app/history"
	"VulpesEditor/app/objectModel/model"
	"VulpesEditor/app/objectModel/texture"
	"VulpesEditor/app/objectModel/tools"
	"VulpesEditor/app/objectModel/view"
	"VulpesEditor/app/util"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"uuid"

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

type instance struct {
	name  string
	title string

	model *model.Model

	id    string
	focus bool
}

func (s *instance) init(m *model.Model) {
	history.New(s.id)
	s.model = m
	tools.New(s.id, s.name, s.model)
	texture.New(s.id, s.name, s.model)
	view.New(s.id, s.model)
	s.title = s.name + " Model"
}

func (s *instance) Focus() bool {
	if s.focus {
		s.focus = false
		return true
	}
	return false
}

func (s *instance) Name() string {
	return s.title
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
	if err := s.model.Save(w); err != nil {
		fmt.Println(err)
		return
	}
	texture.Save(s.id, w)
	b := strings.Builder{}
	b.WriteString("model")
	b.WriteRune('\n')
	b.WriteString(s.name)
	b.WriteRune('\n')
	b.WriteString(s.model.Id)
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
	if len(field) < 3 {
		fmt.Println("Incomplete data")
		return
	}
	if field[0] != "model" {
		fmt.Println("Wrong project type")
		return
	}
	itc := new(instance)
	itc.name = field[1]
	modelId := field[2]
	itc.id = uuid.New().String()
	itc.focus = true
	m, err := model.OpenModel(r, modelId)
	if err != nil {
		fmt.Println(err)
		return
	}
	itc.init(m)
	if err := texture.Open(itc.id, r); err != nil {
		fmt.Println(err)
		return
	}
	tabs.Push(itc)
}

func NewTest() {
	itc := new(instance)
	itc.name = "Test"
	itc.id = uuid.New().String()
	itc.focus = true
	m := model.NewModel()
	m.AddUnit(model.NewUnit([3]float32{0, 0, 0}, [3]float32{1, 1, 1}, [3]float32{0, 0, 0}))
	itc.init(m)
	tabs.Push(itc)
}

func openNew(c creationData) {
	itc := new(instance)
	itc.name = c.name
	itc.id = uuid.New().String()
	itc.focus = true
	m := model.NewModel()
	m.AddUnit(model.NewUnit([3]float32{0, 0, 0}, [3]float32{1, 1, 1}, [3]float32{0, 0, 0}))
	itc.init(m)
	tabs.Push(itc)
}
