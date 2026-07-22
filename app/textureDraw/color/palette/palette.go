package palette

import (
	"VulpesEditor/app/front"
	pallete_file "VulpesEditor/app/textureDraw/color/palette/paleteFile"
	"VulpesEditor/app/util"
	"strconv"
	"strings"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type color struct {
	id    int32
	value [4]float32
	mark  im.Vec4
}

type palette struct {
	id      int32
	name    string
	creator string
	colors  []color
}

type entry struct {
	palette *palette
	show    bool
}

type paletteData struct {
	palettes map[int32]bool
	color1id int32
	color2id int32
}

var palettes []*palette

var idSys = util.NewIdSystem()

func highContrast(rgba [4]float32) im.Vec4 {
	hsv := [4]float32{}
	im.ColorConvertRGBtoHSV(rgba[0], rgba[1], rgba[2], &hsv[0], &hsv[1], &hsv[2])
	if hsv[0] = hsv[0] + 0.5; hsv[0] > 1 {
		hsv[0] -= 1
	}
	hsv[1] = 1
	hsv[2] = 1
	final := [4]float32{}
	im.ColorConvertHSVtoRGB(hsv[0], hsv[1], hsv[2], &final[0], &final[1], &final[2])
	final[3] = 1
	return newVec4(final)
}

func newVec4(vec [4]float32) im.Vec4 {
	return im.NewVec4(vec[0], vec[1], vec[2], vec[3])
}

func addPalette(data pallete_file.PaletteData) (p *palette) {
	p = new(palette)
	p.id = idSys.GetID()
	p.name = data.Name
	p.creator = data.Creator
	for _, c := range data.Colors {
		rgba := [4]float32{c[0], c[1], c[2], 1}
		p.colors = append(p.colors, color{idSys.GetID(), rgba, highContrast(rgba)})
	}
	palettes = append(palettes, p)
	return
}

func addLospecByName(name string) bool {
	if ok, data := pallete_file.GetPaletteFromLospec(name); ok {
		p := addPalette(data)
		ctx.palettes[p.id] = true
		return true
	}
	return false
}

func addLospecByLink(link string) bool {
	if ok, data := pallete_file.GetPaletteFromLospecLink(link); ok {
		p := addPalette(data)
		ctx.palettes[p.id] = true
		return true
	}
	return false
}

func Reset(change [3]bool) {
	if change[0] {
		ctx.color1id = -1
	}
	if change[1] {
		ctx.color2id = -1
	}
	if change[2] {
		temp := ctx.color1id
		ctx.color1id = ctx.color2id
		ctx.color2id = temp
	}
}

func Init() {
	var step float32
	pData := pallete_file.PaletteData{}
	pData.Creator = "This Program"
	pData.Name = "Default"
	greyCount := 10
	step = 1 / float32(greyCount)
	var value float32 = 0
	for range greyCount {
		var rgb [3]float32
		im.ColorConvertHSVtoRGB(1, 0, value, &rgb[0], &rgb[1], &rgb[2])
		value += step
		pData.Colors = append(pData.Colors, rgb)
	}
	colorCount := 20
	step = 1 / float32(colorCount)
	var hue float32 = 0
	for range colorCount {
		var rgb [3]float32
		im.ColorConvertHSVtoRGB(hue, 1, 1, &rgb[0], &rgb[1], &rgb[2])
		hue += step
		pData.Colors = append(pData.Colors, rgb)
	}
	addPalette(pData)

	palettesData := pallete_file.GetAllPalettes()
	for _, data := range palettesData {
		addPalette(data)
	}
}

var lospecInput string

func Loop(id int32, color1, color2 *[4]float32) {
	ctxManager.Check(id)

	var toPop string
	im.BeginV("Color Palette", nil, im.WindowFlagsMenuBar)
	if im.BeginMenuBar() {
		if im.BeginMenu("Import") {
			if im.MenuItemBool("Import from Lospec") {
				toPop = "Lospec"
			}
			if im.MenuItemBool("Import from File") {
				toPop = "Not Implement"
			}
			im.EndMenu()
		}
		if im.BeginMenu("View") {
			for _, p := range palettes {
				show := ctx.palettes[p.id]
				im.MenuItemBoolPtr(p.name, "", &show)
				ctx.palettes[p.id] = show
			}
			im.EndMenu()
		}
		im.EndMenuBar()
	}
	if toPop != "" {
		im.OpenPopupStr(toPop)
		toPop = ""
	}
	if im.BeginPopupModal("Lospec") {
		im.InputTextWithHint("Lospec Palette", "Lospec Palette name or link", &lospecInput, im.InputTextFlagsNone, nil)
		if im.Button("Add") {
			if words := strings.Split(lospecInput, "/"); len(words) > 1 {
				if addLospecByLink(lospecInput) {
					lospecInput = ""
					im.CloseCurrentPopup()
				}
			} else {
				if addLospecByName(lospecInput) {
					lospecInput = ""
					im.CloseCurrentPopup()
				}
			}
		}
		im.SameLine()
		if im.Button("Cancel") {
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	front.NotImplementPopUp()

	var width float32 = 46
	for _, palette := range palettes {
		if !ctx.palettes[palette.id] {
			continue
		}
		im.SeparatorText(palette.name)
		for i, color := range palette.colors {
			id := color.id
			if ctx.color1id == id || ctx.color2id == id {
				im.PushStyleColorVec4(im.ColFrameBg, color.mark)
			}
			availableSpace := im.ContentRegionAvail().X
			im.PushIDInt(id)
			im.ColorButton("color #"+strconv.FormatInt(int64(i), 10), newVec4(color.value))
			im.PopID()
			if ctx.color1id == id || ctx.color2id == id {
				im.PopStyleColor()
			}
			if im.IsItemHovered() {
				io := im.CurrentContext().IO()
				mouseRelease := io.MouseReleased()
				if mouseRelease[0] {
					*color1 = color.value
					ctx.color1id = color.id
				}
				if mouseRelease[1] {
					*color2 = color.value
					ctx.color2id = color.id
				}
			}
			if i != len(palette.colors)-1 && availableSpace-width > 0 {
				im.SameLine()
			}
		}
	}
	im.End()
}
