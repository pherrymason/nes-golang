package types

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
)

const TILE_PIXELS = 8 * 8
const WIDTH = 256
const HEIGHT = 240

type Frame struct {
	Pixels [WIDTH * HEIGHT]Color
}

func (frame *Frame) PushTile(tile [TILE_PIXELS]Pixel, colX int, colY int) {
	//for i := 0; i < TILE_PIXELS; i++ {
	//	frame.Pixels[colX*8*colY*8] = tile[i]
	//}
}

func (frame *Frame) SetPixel(x int, y int, rgb Color) {
	pos := x + WIDTH*y
	if pos >= WIDTH*HEIGHT {
		panic(fmt.Sprintf("Trying to render out of screen: %d", pos))
		//pos = (WIDTH * HEIGHT) - 1
	}
	frame.Pixels[pos] = rgb
}

type Color struct {
	R byte
	G byte
	B byte
}

type Pixel struct {
	X     int
	Y     int
	Color Color
}

func PixelColor2RaylibColor(pixelColor Color) raylib.Color {
	return raylib.NewColor(
		pixelColor.R,
		pixelColor.G,
		pixelColor.B,
		255,
	)
}
