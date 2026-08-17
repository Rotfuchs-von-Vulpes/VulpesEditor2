package tools

import (
	"VulpesEditor/app/objectModel/model"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type ToolContext struct {
	id    int32
	model *model.Model
}

var pos [3]float32
var size = [3]float32{1, 1, 1}

func reset() {
	pos = [3]float32{}
}

func Show(id int32) {
	ctxManager.Check(id)

	if im.Begin("Add Unit") {
		im.InputFloat3("Position", &pos)
		if im.InputFloat3("Size", &size) {
			if size[0] < 0 {
				size[0] = 0
			}
			if size[1] < 0 {
				size[1] = 0
			}
			if size[2] < 0 {
				size[2] = 0
			}
		}
		if im.Button("Add") {
			ctx.model.AddUnit(model.NewUnit(pos, size))
			reset()
		}
		im.End()
	}

	if im.Begin("Units") {

		im.End()
	}
}
