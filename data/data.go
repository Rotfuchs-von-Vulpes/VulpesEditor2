package data

import (
	"os"
)

func Init() {
	folderList := []string{
		"UserData/palettes",
	}

	for _, folder := range folderList {
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
