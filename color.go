package heligo

import (
	"image/color"
	"math"
)

type defaultColorScheme struct{}

func (h *defaultColorScheme) GetColor(groups, index int) color.Color {
	const (
		baseSaturation = 0.7
		baseLightness  = 0.6
	)

	hue := float64(index) / 2 * 360.0
	hue = math.Mod(hue+30, 360)
	c := (1 - math.Abs(2*baseLightness-1)) * baseSaturation
	x := c * (1 - math.Abs(math.Mod(hue/60, 2)-1))
	m := baseLightness - c/2

	var r, g, b float64
	if hue >= 0 && hue < 60 {
		r, g, b = c, x, 0
	} else if hue >= 60 && hue < 120 {
		r, g, b = x, c, 0
	} else if hue >= 120 && hue < 180 {
		r, g, b = 0, c, x
	} else if hue >= 180 && hue < 240 {
		r, g, b = 0, x, c
	} else if hue >= 240 && hue < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}
