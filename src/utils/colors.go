package utils

import (
	"image/color"
)

func NewColorRGB(r uint8, g uint8, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
