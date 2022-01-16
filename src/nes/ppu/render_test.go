package ppu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gets_attribute_table_address_from_tile(t *testing.T) {
	var vram [2048]byte

	// Setup attribute table
	vram[0x03C0] = 0b11011000
	// 0b00 background palette 0
	// 0b01 background palette 1
	// 0b10 background palette 2
	// 0b11 background palette 3
	// Tiles and its palettes --
	//	[00][00] [10][10]
	//  [00][00] [10][10]
	//  [01][01] [11][11]
	//  [01][01] [11][11]

	// First meta tile should
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			palette := backgroundPalette(x, y, vram)
			assert.Equal(t, byte(0b00), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b00", x, y))
		}
	}

	// Second meta tile should
	for x := 2; x < 4; x++ {
		for y := 0; y < 2; y++ {
			palette := backgroundPalette(x, y, vram)
			assert.Equal(t, byte(0b10), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}

	// Third meta tile should
	for x := 0; x < 2; x++ {
		for y := 2; y < 4; y++ {
			palette := backgroundPalette(x, y, vram)
			assert.Equal(t, byte(0b01), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}

	// Fourth meta tile should
	for x := 2; x < 4; x++ {
		for y := 2; y < 4; y++ {
			palette := backgroundPalette(x, y, vram)
			assert.Equal(t, byte(0b11), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}
}
