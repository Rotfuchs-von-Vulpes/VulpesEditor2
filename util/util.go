package util

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Data struct {
	ClearColor  [3]float32
	ObjectColor [3]float32
}

func GetData(clearColor, objectColor [3]float32) Data {
	return Data{clearColor, objectColor}
}

func Str(str string) *uint8 {
	return gl.Str(str + "\x00")
}
