package nes

import (
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateDummyGamePak() *gamePak.GamePak {
	pak := gamePak.CreateGamePak(
		gamePak.CreateINes1Header(1, 1, 0, 0, 0, 0, 0),
		make([]byte, 100),
		make([]byte, 0x01FFF),
	)

	return &pak
}

func aPPU() *Ppu2c02 {
	cartridge := CreateDummyGamePak()
	memory := CreatePPUMemory(cartridge)

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

func TestPPU_tick_should_start_vblank_on_scanline_240(t *testing.T) {
	cases := []struct {
		name     string
		allowNMI bool
	}{
		{"vblank + nmi allowed", true},
		{"vblank + nmi disallowed", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ppu := aPPU()
			if tt.allowNMI {
				ppu.ppuctrlWriteFlag(generateNMIAtVBlank, 1)
			} else {
				ppu.ppuctrlWriteFlag(generateNMIAtVBlank, 0)
			}
			ppu.cycle = PPU_VBLANK_START_CYCLE

			ppu.Tick()

			assert.Equal(t, byte(1), (ppu.registers.status>>verticalBlankStarted)&0x01)
			assert.Equal(t, tt.allowNMI, ppu.nmi, "should have enabled NMI on vblank")
		})
	}
}

func TestPPU_tick_should_end_vblank_on_end_of_scanline_261(t *testing.T) {
	ppu := aPPU()
	ppu.cycle = PPU_VBLANK_END_CYCLE
	ppu.registers.status |= 1 << verticalBlankStarted

	ppu.Tick()

	assert.Equal(t, byte(0), (ppu.registers.status>>verticalBlankStarted)&0x01)
}

func TestPPU_should_get_propper_color_for_a_given_pixel_color_and_palette(t *testing.T) {
	ppu := aPPU()
	backgroundColor := graphics.Color{236, 88, 180}
	cases := []struct {
		name          string
		palette       byte
		pixelColor    byte
		expectedColor graphics.Color
	}{
		{"", 0, 0, backgroundColor},
		{"", 0, 1, graphics.Color{0, 30, 116}},
		{"", 0, 2, graphics.Color{8, 16, 144}},
		{"", 0, 3, graphics.Color{48, 0, 136}},
		{"mirroring $0x3F10", 4, 0, backgroundColor},
		{"mirroring $0x3F14", 5, 0, graphics.Color{68, 0, 100}},
		{"mirroring $0x3F18", 6, 0, graphics.Color{32, 42, 0}},
		{"mirroring $0x3F1C", 7, 0, graphics.Color{0, 50, 60}},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			color := ppu.getColorFromPaletteRam(tt.palette, tt.pixelColor)
			assert.Equal(t, tt.expectedColor, color)
		})
	}
}
