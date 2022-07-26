package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func aPPU() *Ppu2c02 {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	ppu := CreatePPU(cartridge)
	return ppu
}

func Test_should_trigger_vBlank_on_scanline_240(t *testing.T) {
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
				ppu.ppuControl.generateNMIAtVBlank = true
			} else {
				ppu.ppuControl.generateNMIAtVBlank = false
			}
			ppu.cycle = PPU_VBLANK_START_CYCLE
			ppu.currentScanline = PPU_SCREEN_SPACE_SCANLINES

			ppu.Tick()

			assert.True(t, ppu.ppuStatus.verticalBlankStarted)
			assert.Equal(t, tt.allowNMI, ppu.nmi, "Unexpected NMI behaviour")
		})
	}
}

func TestPPU_should_end_vblank_on_end_of_scanline_261(t *testing.T) {
	ppu := aPPU()
	ppu.renderCycle = PPU_CYCLES_BY_SCANLINE - 1
	ppu.currentScanline = VBLANK_END_SCNALINE
	ppu.ppuStatus.verticalBlankStarted = true

	ppu.Tick()

	assert.False(t, ppu.ppuStatus.verticalBlankStarted)
}

func TestPPU_VBlank_should_return_true_when_current_scanline_is_above_241(t *testing.T) {
	ppu := CreatePPU(gamePak.NewDummyGamePak(gamePak.NewEmptyCHRROM()))
	ppu.currentScanline = 241

	assert.True(t, ppu.VBlank())
}

func TestPPU_VBlank_should_return_false_when_current_scanline_is_below_241(t *testing.T) {
	ppu := CreatePPU(gamePak.NewDummyGamePak(gamePak.NewEmptyCHRROM()))
	ppu.currentScanline = 240

	assert.False(t, ppu.VBlank())
}
