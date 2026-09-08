package tools

import (
	"VulpesEditor/app/objectModel/model"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type ToolContext struct {
	name      string
	model     *model.Model
	sizeInput [3]float32
}

var pos [3]float32
var size = [3]float32{1, 1, 1}
var posInput [3]float32
var originalPos [3]float32
var originalSize [3]float32
var editingUnit *model.CubeUnit = model.NewUnit(pos, size)

func toAbsolute(v1 [3]float32) (v2 [3]float32) {
	v2[0] = v1[0] / 16
	v2[1] = v1[1] / 16
	v2[2] = v1[2] / 16
	return
}

func toRelative(v1 [3]float32) (v2 [3]float32) {
	v2[0] = v1[0] * 16
	v2[1] = v1[1] * 16
	v2[2] = v1[2] * 16
	return
}

func reset() {
	// texture.GenerateTexture()
	posInput = [3]float32{}
	editingUnit = model.NewUnit(pos, size)
}

func setToEdit(unit *model.CubeUnit) {
	editingUnit = unit
	pos, size = unit.Data()
	originalPos = pos
	originalSize = size
	posInput = toRelative(pos)
	ctx.sizeInput = toRelative(size)
}

func Write() {
	editingUnit.Edit(pos, size)
	reset()
}

func BeginDragSource() {

}

var draging *model.CubeUnit

func DragSource(cube *model.CubeUnit) {
	if draging == nil {
		draging = cube
	}
}

func DragTarget(target model.NodeTree) {
	if draging != nil {
		unitSource := draging
		if unitSource.GetId() != target.GetId() {
			target.Append(unitSource)
		}
		draging = nil
	}
}

type Payload struct {
	kind string
	data any
}

var selected string
var hovered bool

func ShowUnits(node model.NodeTree) {
	if node == nil {
		return
	}

	im.PushIDStr(node.GetId())

	if unit, ok := node.(*model.CubeUnit); ok {
		im.SetNextItemOpenV(true, im.CondOnce)
		if im.TreeNodeStr(unit.Name) {

			if draging != nil {
				cant := !unit.CanReceive(draging)
				if cant {
					im.BeginDisabled()
				}
				im.Button("Drop")
				if cant {
					im.EndDisabled()
				}
			} else {
				im.Button("Drag")
			}
			if im.IsItemHovered() {
				hovered = true
				selected = node.GetId()

				if im.IsMouseDown(im.MouseButtonLeft) && ok {
					DragSource(unit)
				}

				if im.IsMouseReleased(im.MouseButtonLeft) {
					DragTarget(node)
				}
			}

			im.SameLine()
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

			for _, u := range unit.Children() {
				ShowUnits(u)
			}

			im.TreePop()
		}
	}
	if m, ok := node.(*model.Model); ok {
		im.SetNextItemOpenV(true, im.CondOnce)
		if im.TreeNodeStr(ctx.name) {

			if draging != nil {
				cant := !m.CanReceive(draging)
				if cant {
					im.BeginDisabled()
				}
				im.Button("Drop")
				if cant {
					im.EndDisabled()
				}
			} else {
				im.BeginDisabled()
				im.Button("Drag")
				im.EndDisabled()
			}
			if im.IsItemHovered() {
				hovered = true
				selected = node.GetId()

				if im.IsMouseReleased(im.MouseButtonLeft) {
					DragTarget(node)
				}
			}

			for _, u := range m.Children() {
				ShowUnits(u)
			}

			im.TreePop()
		}
	}

	im.PopID()
}

func Show(id string) {
	ctxManager.Check(id)

	if im.Begin("Add Unit") {
		var c1 bool
		if im.Button("Reflect") {
			posInput[0] = -posInput[0]
			c1 = true
		}

		c2 := im.InputFloat3("Position", &posInput)
		c3 := im.InputFloat3("Size", &ctx.sizeInput)

		if c3 {
			if ctx.sizeInput[0] < 0 {
				ctx.sizeInput[0] = 0
			}
			if ctx.sizeInput[1] < 0 {
				ctx.sizeInput[1] = 0
			}
			if ctx.sizeInput[2] < 0 {
				ctx.sizeInput[2] = 0
			}
		}

		if c1 || c2 || c3 {
			pos = toAbsolute(posInput)
			size = toAbsolute(ctx.sizeInput)
			editingUnit.Edit(pos, size)
		}

		if editingUnit.Source == nil {
			if im.Button("Add") {
				ctx.model.AddUnit(editingUnit)
				reset()
			}
		} else {
			if im.Button("Set") {
				reset()
			}
			im.SameLine()
			if im.Button("Delete") {
				ctx.model.Remove(editingUnit)
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
		hovered = false

		if draging != nil {
			if im.BeginTooltip() {
				im.Text("Moving " + draging.Name)
				im.EndTooltip()
			}
		}

		ShowUnits(ctx.model)

		if im.IsMouseReleased(im.MouseButtonLeft) {
			draging = nil
		}

		if !hovered {
			selected = ""
		}

		im.End()
	}
}
