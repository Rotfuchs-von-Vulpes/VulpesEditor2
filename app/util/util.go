package util

import (
	"os"
	"path/filepath"
)

var AppDir string

func Init() {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	AppDir = filepath.Join(path, "VulpesEditor")
}

type IdSystem struct {
	lastId int32
}

func NewIdSystem() (sys *IdSystem) {
	sys = new(IdSystem)
	sys.lastId = 0
	return
}

func (s *IdSystem) GetID() (id int32) {
	id = s.lastId
	s.lastId += 1
	return
}
