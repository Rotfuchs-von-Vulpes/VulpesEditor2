package util

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func Str(str string) *uint8 {
	return gl.Str(str + "\x00")
}
