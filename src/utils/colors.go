package utils

import (
	"github.com/lachee/raylib-goplus/raylib"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image/color"
)

func NewColorRGB(r uint8, g uint8, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func RGBA2raylibColor(pixelColor types.Color) raylib.Color {
	return raylib.NewColor(
		pixelColor.R,
		pixelColor.G,
		pixelColor.B,
		255,
	)
}
