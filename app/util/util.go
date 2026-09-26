package util

import (
	"os"
	"path/filepath"
)

var AppDir string

var paths = [][]string{
	{"projects", "textures"},
	{"projects", "models"},
	{"projects", "blocks"},
	{"assets", "textures"},
}

func Init() {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	AppDir = filepath.Join(path, "VulpesEditor")
	for _, p := range paths {
		if err := os.MkdirAll(filepath.Join(AppDir, filepath.Join(p...)), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
