package textureEdit

import (
	"fmt"

	im "github.com/AllenDang/cimgui-go/imgui"
)

var selectedLayers = []bool{}

func (s *TextureEdit) ShowLayers() {
	im.Begin("Layers")

	if im.Button("Add") {
		s.AppendLayer()
	}
	im.SameLine()
	if im.Button("Remove") {
		selectedLayers = make([]bool, len(s.layers))
		im.OpenPopupStr("Remove Layers")
	}
	im.SameLine()
	if im.Button("Merge") {
		selectedLayers = make([]bool, len(s.layers))
		im.OpenPopupStr("Merge Layers")
	}

	for idx := range s.layers {
		i := len(s.layers) - idx - 1
		layer := s.layers[i]
		cannotDown := i == 0
		cannotUp := i == len(s.layers)-1
		selected := layer.Id == s.layer.Id
		str := fmt.Sprintf("Layer #%d", i)
		if selected {
			// #86BDFFFF
			im.PushStyleColorVec4(im.ColText, im.NewVec4(0.52, 0.74, 1, 1))
		}
		im.PushIDStr(str)
		if im.ImageButton("Set", layer.Image.Tex.ID, im.NewVec2(20, 20)) {
			s.layer = s.layers[i]
			s.SetLayer(i)
		}
		im.SameLine()
		im.Text(str)
		im.SameLine()
		if cannotDown {
			im.BeginDisabled()
		}
		if im.Button("down") {
			s.Swap(i, i-1)
		}
		if cannotDown {
			im.EndDisabled()
		}
		im.SameLine()
		if cannotUp {
			im.BeginDisabled()
		}
		if im.Button("up") {
			s.Swap(i, i+1)
		}
		if cannotUp {
			im.EndDisabled()
		}
		im.SameLine()
		if im.Checkbox("Show", &s.layers[i].Show) {
			s.UpdateTexture()
		}
		im.PopID()
		if selected {
			im.PopStyleColor()
		}
	}

	if im.BeginPopupModal("Remove Layers") {
		im.Text("Select Layers to remove")
		for i := range s.layers {
			str := fmt.Sprintf("Layer #%d", i)
			im.Checkbox(str, &selectedLayers[i])
			im.SameLine()
			im.ImageWithBg(s.layers[i].Image.Tex.ID, im.NewVec2(15, 15))
		}
		if im.Button("Remove") {
			s.Remove(selectedLayers)
			im.CloseCurrentPopup()
		}
		im.SameLine()
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	if im.BeginPopupModal("Merge Layers") {
		im.Text("Select Layers to merge:")
		for i := range s.layers {
			str := fmt.Sprintf("Layer #%d", i)
			im.Checkbox(str, &selectedLayers[i])
			im.SameLine()
			im.ImageWithBg(s.layers[i].Image.Tex.ID, im.NewVec2(15, 15))
		}
		if im.Button("Merge") {
			s.Merge(selectedLayers)
			im.CloseCurrentPopup()
		}
		im.SameLine()
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	im.End()
}
