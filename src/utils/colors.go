package utils

import (
	"github.com/lachee/raylib-goplus/raylib"
	"image/color"
)

func NewColorRGB(r uint8, g uint8, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func RGBA2raylibColor(pixelColor color.Color) raylib.Color {
	r, g, b, _ := pixelColor.RGBA()

	return raylib.NewColor(
		uint8(r),
		uint8(g),
		uint8(b),
		255,
	)
}
