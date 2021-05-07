package nes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateDummyGamePak() *GamePak {
	return &GamePak{
		Header{1, 1, 0, 0, 0, 0, 0},
		make([]byte, 100),
		make([]byte, 0x01FFF),
	}
}

func aPPU() *Ppu2c02 {
	gamePak := CreateDummyGamePak()
	memory := CreatePPUMemory(gamePak)

	// 0x3F00
	memory.paletteTable[0] = 0x21 // light blue
	// 0x3F01
	memory.paletteTable[1] = 0x01 // Dark Blue
	memory.paletteTable[2] = 0x02 // Blue-Purple
	memory.paletteTable[3] = 0x03 // Dark Purple

	// Mirrored addresses
	memory.paletteTable[4] = 0x04
	memory.paletteTable[8] = 0x08
	memory.paletteTable[0x0C] = 0x0C

	ppu := CreatePPU(memory)
	return ppu
}

func TestPPU_should_get_propper_color_for_a_given_pixel_color_and_palette(t *testing.T) {
	ppu := aPPU()
	backgroundColor := [3]byte{236, 88, 180}
	cases := []struct {
		name          string
		palette       byte
		pixelColor    byte
		expectedColor [3]byte
	}{
		{"", 0, 0, backgroundColor},
		{"", 0, 1, [3]byte{0, 30, 116}},
		{"", 0, 2, [3]byte{8, 16, 144}},
		{"", 0, 3, [3]byte{48, 0, 136}},
		{"mirroring $0x3F10", 4, 0, backgroundColor},
		{"mirroring $0x3F14", 5, 0, [3]byte{68, 0, 100}},
		{"mirroring $0x3F18", 6, 0, [3]byte{32, 42, 0}},
		{"mirroring $0x3F1C", 7, 0, [3]byte{0, 50, 60}},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			color := ppu.getColorFromPaletteRam(tt.palette, tt.pixelColor)
			assert.Equal(t, tt.expectedColor, color)
		})
	}
}
