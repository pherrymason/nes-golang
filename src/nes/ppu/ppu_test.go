package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
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
	memory := CreateMemory(cartridge)
	ppu := CreatePPU(*memory)
	return ppu
}

func TestPPU_tick_should_start_vblank_on_scanline_240(t *testing.T) {
	cases := []struct {
		name     string
		allowNMI bool
	}{
		{"should trigger nmi on vblank", true},
		{"should not trigger nmi on vblank", false},
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
			ppu.currentScanline = PPU_SCREEN_SPACE_SCANLINES

			ppu.Tick()

			assert.Equal(t, byte(1), (ppu.registers.status>>verticalBlankStarted)&0x01)
			assert.Equal(t, tt.allowNMI, ppu.nmi, "Unexpected NMI behaviour")
		})
	}
}

func TestPPU_tick_should_end_vblank_on_end_of_scanline_261(t *testing.T) {
	ppu := aPPU()
	ppu.renderCycle = PPU_CYCLES_BY_SCANLINE - 1
	ppu.currentScanline = VBLANK_END_SCNALINE
	ppu.registers.status |= 1 << verticalBlankStarted

	ppu.Tick()

	assert.Equal(t, byte(0), (ppu.registers.status>>verticalBlankStarted)&0x01)
}

func TestPPU_writes_and_reads_into_palette(t *testing.T) {
	ppu := aPPU()

	for i := 0; i < 32; i++ {
		colorIndex := byte(i + 1)
		address := PaletteLowAddress + types.Address(i)
		ppu.Write(address, colorIndex)
		readValue := ppu.Read(address)
		assert.Equal(
			t,
			colorIndex,
			readValue,
			fmt.Sprintf("@%X has unexpected value", address),
		)
	}
}

func TestPPU_should_get_propper_color_for_a_given_pixel_color_and_palette(t *testing.T) {
	ppu := aPPU()
	backgroundColor := types.Color{236, 88, 180}
	cases := []struct {
		name          string
		palette       byte
		colorIndex    byte
		expectedColor types.Color
	}{
		{"", 0, 0, backgroundColor},
		{"", 0, 1, types.Color{0, 30, 116}},
		{"", 0, 2, types.Color{8, 16, 144}},
		{"", 0, 3, types.Color{48, 0, 136}},
		//{"mirroring $0x3F10", 4, 0, backgroundColor},
		//{"mirroring $0x3F14", 5, 0, graphics.Color{68, 0, 100}},
		//{"mirroring $0x3F18", 6, 0, graphics.Color{32, 42, 0}},
		//{"mirroring $0x3F1C", 7, 0, graphics.Color{0, 50, 60}},
	}

	// 0x3F00
	ppu.Write(PaletteLowAddress+0, 0x25) // light blue
	ppu.Write(PaletteLowAddress+1, 0x01) // Dark Blue
	ppu.Write(PaletteLowAddress+2, 0x02) // Blue-Purple
	ppu.Write(PaletteLowAddress+3, 0x03) // Dark Purple

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			color := ppu.GetColorFromPaletteRam(tt.palette, tt.colorIndex)
			assert.Equal(t, tt.expectedColor, color)
		})
	}
}
