package types

import (
	"fmt"
)

const SCREEN_WIDTH = 256
const SCREEN_HEIGHT = 240

type Frame struct {
	Pixels [SCREEN_WIDTH * SCREEN_HEIGHT]Color
}

func (frame *Frame) PushTile(tile Tile, x int, y int) {
	baseY := y * 256
	baseX := x
	for i := 0; i < 8*8; i++ {
		calculatedY := baseY + (i/8)*SCREEN_WIDTH
		calculatedX := baseX + i%8
		arrayIndex := calculatedX + calculatedY
		frame.Pixels[arrayIndex] = tile.Pixels[i]
	}
}

func (frame *Frame) SetPixel(x int, y int, rgb Color) {
	pos := x + SCREEN_WIDTH*y
	if pos >= SCREEN_WIDTH*SCREEN_HEIGHT {
		panic(fmt.Sprintf("Trying to render out of screen: %d", pos))
		//pos = (SCREEN_WIDTH * SCREEN_HEIGHT) - 1
	}
	frame.Pixels[pos] = rgb
}

func CoordinatesToArrayIndex(x int, y int, canvasWidth int) int {
	return x + canvasWidth*y
}

func LinearToXCoordinate(index int, canvasWidth int) int32 {
	return int32(index % canvasWidth)
}

func LinearToYCoordinate(index int, canvasWidth int) int32 {
	return int32(index / canvasWidth)
}

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
