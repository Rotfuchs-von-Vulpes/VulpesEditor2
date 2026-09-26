package blockCreator

import (
	"VulpesEditor/app/editors/blockCreator/blocks"
	"VulpesEditor/app/editors/blockCreator/view"
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/tabs"
	"VulpesEditor/app/util"
	"fmt"
	"path/filepath"
	"strings"
	"uuid"
)

type instance struct {
	name  string
	title string
	id    string

	focus bool
}

func (s *instance) init() {
	s.title = s.name + " Blocks"
	blocks.New(s.id)
	view.New(s.id)
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
	blocks.Show(s.id)
	view.Show(s.id)
}

func (s *instance) Save() {
	w, err := file.NewArchive(filepath.Join(util.AppDir, "projects", "blocks"), s.name)
	if err != nil {
		fmt.Println(err)
		return
	}
	// canvas.Save(w)
	b := strings.Builder{}
	b.WriteString("block")
	b.WriteRune('\n')
	b.WriteString(s.name)
	w.Write("metaData.txt", []byte(b.String()))
	w.Save()
}

func (s *instance) Export() {

}

type creationData struct {
	name string
}

func NewTest() {
	itc := new(instance)
	itc.name = "Test"
	itc.id = uuid.New().String()
	itc.focus = true
	itc.init()
	tabs.Push(itc)
}

func openNew(c creationData) {
	itc := new(instance)
	itc.id = uuid.Max().String()
	itc.name = c.name
	itc.focus = true
	itc.init()
	tabs.Push(itc)
}
