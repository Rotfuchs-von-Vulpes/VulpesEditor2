package oklab

import (
	"math"

	im "github.com/AllenDang/cimgui-go/imgui"
)

type oklab struct {
	l, a, b float32
}

type oklch struct {
	l, c, h float32
}

type rgb struct {
	r, g, b float32
}

func (self rgb) toOklab() (final oklab) {
	l := 0.4122214708*self.r + 0.5363325363*self.g + 0.0514459929*self.b
	m := 0.2119034982*self.r + 0.6806995451*self.g + 0.1073969566*self.b
	s := 0.0883024619*self.r + 0.2817188376*self.g + 0.6299787005*self.b

	l_ := float32(math.Cbrt(float64(l)))
	m_ := float32(math.Cbrt(float64(m)))
	s_ := float32(math.Cbrt(float64(s)))

	final.l = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	final.a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	final.b = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_

	return
}

func (self oklab) toRGB() (final rgb) {
	l_ := self.l + 0.3963377774*self.a + 0.2158037573*self.b
	m_ := self.l - 0.1055613458*self.a - 0.0638541728*self.b
	s_ := self.l - 0.0894841775*self.a - 1.2914855414*self.b

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	final.r = +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	final.g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	final.b = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s

	return
}

func (self oklab) toOklch() (final oklch) {
	c := float32(math.Sqrt(float64(self.a*self.a + self.b*self.b)))
	h := float32(math.Atan2(float64(self.b), float64(self.a)))
	if h < 0.0 {
		h += 2.0 * 3.1415926535
	}
	final.l = self.l
	final.c = c
	final.h = h

	return
}

func (self oklch) toOklab() (final oklab) {
	final.l = self.l
	final.a = self.c * float32(math.Cos(float64(self.h)))
	final.b = self.c * float32(math.Sin(float64(self.h)))

	return
}

func (self oklch) toRGB() rgb {
	return self.toOklab().toRGB().clamp()
}

func clamp(v float32, mx float32, mn float32) float32 {
	return max(min(v, mx), mn)
}

func (self rgb) clamp() (final rgb) {
	final.r = clamp(self.r, 1, 0)
	final.g = clamp(self.g, 1, 0)
	final.b = clamp(self.b, 1, 0)

	return
}

func OklchColorPicker4(label string, color *[4]float32) (changed bool) {
	im.PushIDStr(label)

	current_rgb := rgb{color[0], color[1], color[2]}
	lab := current_rgb.toOklab()
	lch := lab.toOklch()

	if im.SliderFloatV("Lightness (L)", &lch.l, 0.0, 1.0, "%.3f", im.SliderFlagsNone) {
		changed = true
	}
	if im.SliderFloatV("Chroma (C)", &lch.c, 0.0, 0.4, "%.3f", im.SliderFlagsNone) {
		changed = true
	}
	hue_degrees := lch.h * (180.0 / 3.1415926535)
	if im.SliderFloatV("Hue (h°)", &hue_degrees, 0.0, 360.0, "%.1f", im.SliderFlagsNone) {
		lch.h = hue_degrees * (3.1415926535 / 180.0)
		changed = true
	}
	im.SliderFloat("Alpha", &color[3], 0, 1)

	if changed {
		target := lch.toRGB()
		color[0] = target.r
		color[1] = target.g
		color[2] = target.b
	}

	im.ColorButton("##current_color", im.NewVec4(color[0], color[1], color[2], color[3]))
	im.SameLine()
	im.TextUnformatted(label)

	im.PopID()

	return
}
