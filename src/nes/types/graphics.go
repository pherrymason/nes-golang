package types

import (
	"fmt"
	"github.com/lachee/raylib-goplus/raylib"
)

const WIDTH = 256
const HEIGHT = 240

type Frame struct {
	Pixels [WIDTH * HEIGHT]Color
}

func (frame *Frame) PushTile(tile Tile, x int, y int) {
	baseY := y * 256
	baseX := x
	for i := 0; i < TILE_PIXELS; i++ {
		calculatedY := baseY + (i/8)*WIDTH
		calculatedX := baseX + i%8
		arrayIndex := calculatedX + calculatedY
		frame.Pixels[arrayIndex] = tile.Pixels[i]
	}
}

func (frame *Frame) SetPixel(x int, y int, rgb Color) {
	pos := x + WIDTH*y
	if pos >= WIDTH*HEIGHT {
		panic(fmt.Sprintf("Trying to render out of screen: %d", pos))
		//pos = (WIDTH * HEIGHT) - 1
	}
	frame.Pixels[pos] = rgb
}

func CoordinatesToArrayIndex(x int, y int, canvasWidth int) int {
	return x + canvasWidth*y
}

func LinearToXCoordinate(index int, canvasWidth int) int {
	return index % canvasWidth
}

func LinearToYCoordinate(index int, canvasWidth int) int {
	return index / canvasWidth
}

const TILE_WIDTH = 8
const TILE_HEIGHT = 8
const TILE_PIXELS = 8 * 8

type Tile struct {
	Pixels [8 * 8]Color
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
