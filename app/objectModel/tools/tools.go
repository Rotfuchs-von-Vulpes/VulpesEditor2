package tools

import (
	"VulpesEditor/app/objectModel/model"
	"math"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type ToolContext struct {
	name          string
	model         *model.Model
	rotationInput [3]float32
	sizeInput     [3]float32
}

var posInput [3]float32
var originalPos [3]float32
var originalSize [3]float32
var originalRot [3]float32
var editingUnit *model.CubeUnit = model.NewUnit([3]float32{}, [3]float32{1, 1, 1}, [3]float32{})

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

func toRadians(v1 [3]float32) (v2 [3]float32) {
	v2[0] = v1[0] * math.Pi / 180.0
	v2[1] = v1[1] * math.Pi / 180.0
	v2[2] = v1[2] * math.Pi / 180.0
	return
}

func toDegree(v1 [3]float32) (v2 [3]float32) {
	v2[0] = v1[0] * 180.0 / math.Pi
	v2[1] = v1[1] * 180.0 / math.Pi
	v2[2] = v1[2] * 180.0 / math.Pi
	return
}

func reset() {
	// texture.GenerateTexture()
	posInput = [3]float32{}
	editingUnit = model.NewUnit([3]float32{}, [3]float32{1, 1, 1}, [3]float32{})
}

func setToEdit(unit *model.CubeUnit) {
	editingUnit = unit
	pos, size, rot := unit.Data()
	originalPos = pos
	originalSize = size
	originalRot = rot
	posInput = toRelative(pos)
	ctx.sizeInput = toRelative(size)
	ctx.rotationInput = toDegree(rot)
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
				pos, size, rot := unit.Data()
				u := model.NewUnit(pos, size, rot)
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
		c3 := im.InputFloat3("Rotation", &ctx.rotationInput)
		c4 := im.InputFloat3("Size", &ctx.sizeInput)

		if c4 {
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

		if c1 || c2 || c3 || c4 {
			pos := toAbsolute(posInput)
			size := toAbsolute(ctx.sizeInput)
			rot := toRadians(ctx.rotationInput)
			editingUnit.Edit(pos, size, rot)
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
				editingUnit.Edit(originalPos, originalSize, originalRot)
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
