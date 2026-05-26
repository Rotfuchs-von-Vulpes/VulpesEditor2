package palette

import (
	pallete_file "VulpesEditor/app/textureDraw/color/palette/paleteFile"
	"VulpesEditor/app/util"
	"strconv"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
)

type color struct {
	id    int32
	value [4]float32
	mark  imgui.Vec4
}

type palette struct {
	id      int32
	name    string
	creator string
	colors  []color
	show    bool
}

var idSys *util.IdSystem
var palettes []*palette

func highContrast(rgba [4]float32) imgui.Vec4 {
	hsv := [4]float32{}
	imgui.ColorConvertRGBtoHSV(rgba[0], rgba[1], rgba[2], &hsv[0], &hsv[1], &hsv[2])
	if hsv[0] = hsv[0] + 0.5; hsv[0] > 1 {
		hsv[0] -= 1
	}
	hsv[1] = 1
	hsv[2] = 1
	final := [4]float32{}
	imgui.ColorConvertHSVtoRGB(hsv[0], hsv[1], hsv[2], &final[0], &final[1], &final[2])
	final[3] = 1
	return newVec4(final)
}

func newVec4(vec [4]float32) imgui.Vec4 {
	return imgui.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

func addPalette(data pallete_file.PaletteData, show bool) (p *palette) {
	p = new(palette)
	p.id = idSys.GetID()
	p.name = data.Name
	p.creator = data.Creator
	p.show = show
	for _, c := range data.Colors {
		rgba := [4]float32{c[0], c[1], c[2], 1}
		p.colors = append(p.colors, color{idSys.GetID(), rgba, highContrast(rgba)})
	}
	palettes = append(palettes, p)
	return
}

func addLospecByName(name string) bool {
	if ok, data := pallete_file.GetPaletteFromLospec(name); ok {
		addPalette(data, true)
		return true
	}
	return false
}

func addLospecByLink(link string) bool {
	if ok, data := pallete_file.GetPaletteFromLospecLink(link); ok {
		addPalette(data, true)
		return true
	}
	return false
}

var color1id, color2id int32

func Reset(change [3]bool) {
	if change[0] {
		color1id = -1
	}
	if change[1] {
		color2id = -1
	}
	if change[2] {
		temp := color1id
		color1id = color2id
		color2id = temp
	}
}

func Init(color1, color2 *[4]float32) {
	idSys = util.NewIdSystem()

	var step float32
	pData := pallete_file.PaletteData{}
	pData.Creator = "This Program"
	pData.Name = "Default"
	greyCount := 10
	step = 1 / float32(greyCount)
	var value float32 = 0
	for range greyCount {
		var rgb [3]float32
		imgui.ColorConvertHSVtoRGB(1, 0, value, &rgb[0], &rgb[1], &rgb[2])
		value += step
		pData.Colors = append(pData.Colors, rgb)
	}
	colorCount := 20
	step = 1 / float32(colorCount)
	var hue float32 = 0
	for range colorCount {
		var rgb [3]float32
		imgui.ColorConvertHSVtoRGB(hue, 1, 1, &rgb[0], &rgb[1], &rgb[2])
		hue += step
		pData.Colors = append(pData.Colors, rgb)
	}
	p := addPalette(pData, true)
	*color1 = p.colors[0].value
	color1id = p.colors[0].id
	*color2 = p.colors[greyCount-1].value
	color2id = p.colors[greyCount-1].id

	palettesData := pallete_file.GetAllPalettes()
	for _, data := range palettesData {
		addPalette(data, false)
	}
}

var lospecInput string

func Loop(color1, color2 *[4]float32) {
	var toPop string
	imgui.BeginV("Color Palette", nil, imgui.WindowFlagsMenuBar)
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Import") {
			if imgui.MenuItemBool("Import from Lospec") {
				toPop = "Lospec"
			}
			if imgui.MenuItemBool("Import from File") {
				toPop = "Not Implement"
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("View") {
			for i := range palettes {
				imgui.MenuItemBoolPtr(palettes[i].name, "", &palettes[i].show)
			}
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
	if toPop != "" {
		imgui.OpenPopupStr(toPop)
		toPop = ""
	}
	if imgui.BeginPopupModal("Lospec") {
		imgui.InputTextWithHint("Lospec Palette", "Lospec Palette name or link", &lospecInput, imgui.InputTextFlagsNone, nil)
		if imgui.Button("Add") {
			if words := strings.Split(lospecInput, "/"); len(words) > 1 {
				if addLospecByLink(lospecInput) {
					lospecInput = ""
					imgui.CloseCurrentPopup()
				}
			} else {
				if addLospecByName(lospecInput) {
					lospecInput = ""
					imgui.CloseCurrentPopup()
				}
			}
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}
	if imgui.BeginPopupModal("Not Implement") {
		imgui.Text("Not implement yet!")
		if imgui.Button("OK") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}
	var width float32 = 46
	for _, palette := range palettes {
		if !palette.show {
			continue
		}
		imgui.SeparatorText(palette.name)
		for i, color := range palette.colors {
			id := color.id
			if color1id == id || color2id == id {
				imgui.PushStyleColorVec4(imgui.ColFrameBg, color.mark)
			}
			availableSpace := imgui.ContentRegionAvail().X
			imgui.PushIDInt(id)
			imgui.ColorButton("color #"+strconv.FormatInt(int64(i), 10), newVec4(color.value))
			imgui.PopID()
			if color1id == id || color2id == id {
				imgui.PopStyleColor()
			}
			if imgui.IsItemHovered() {
				io := imgui.CurrentContext().IO()
				mouseRelease := io.MouseReleased()
				if mouseRelease[0] {
					*color1 = color.value
					color1id = color.id
				}
				if mouseRelease[1] {
					*color2 = color.value
					color2id = color.id
				}
			}
			if i != len(palette.colors)-1 && availableSpace-width > 0 {
				imgui.SameLine()
			}
		}
	}
	imgui.End()
}
