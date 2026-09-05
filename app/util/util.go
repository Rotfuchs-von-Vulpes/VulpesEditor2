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
