package tools

import (
	"VulpesEditor/app/objectModel/model"
	"fmt"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type ToolContext struct {
	id    int32
	model *model.Model
}

var pos [3]float32
var size = [3]float32{1, 1, 1}
var originalPos [3]float32
var originalSize [3]float32
var editingUnit *model.CubeUnit = model.NewUnit(pos, size)

func reset() {
	pos = [3]float32{}
	editingUnit = model.NewUnit(pos, size)
}

func setToEdit(unit *model.CubeUnit) {
	editingUnit = unit
	pos, size = unit.Data()
	originalPos = pos
	originalSize = size
}

func Write() {
	editingUnit.Edit(pos, size)
	reset()
}

func Show(id int32) {
	ctxManager.Check(id)

	if im.Begin("Add Unit") {
		var c1 bool
		if im.Button("Reflect") {
			pos[0] = -pos[0]
			c1 = true
		}

		c2 := im.InputFloat3("Position", &pos)
		c3 := im.InputFloat3("Size", &size)

		if c3 {
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

		if c1 || c2 || c3 {
			editingUnit.Edit(pos, size)
		}

		if editingUnit.Parent == nil {
			if im.Button("Add") {
				ctx.model.AddUnit(editingUnit)
				reset()
			}
		} else {
			if im.Button("Set") {
				reset()
			}
			im.SameLine()
			if im.Button("Cancel") {
				editingUnit.Edit(originalPos, originalSize)
				reset()
			}
		}
		im.End()
	}

	if im.Begin("Units") {
		for idx, unit := range ctx.model.Units {
			im.Text(fmt.Sprintf("Unit #%d", idx))
			im.SameLine()
			im.PushIDInt(int32(idx))
			if im.Button("Edit") {
				setToEdit(unit)
			}
			im.SameLine()
			if im.Button("Clone") {
				pos, size := unit.Data()
				u := model.NewUnit(pos, size)
				ctx.model.AddUnit(u)
				setToEdit(u)
			}
			im.PopID()
		}
		im.End()
	}
}
