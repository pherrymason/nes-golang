package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newNotWarmedUpPPU() *Ppu2c02 {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	ppu := CreatePPU(cartridge, false, "")
	ppu.warmup = false

	return ppu
}

func aPPU() *Ppu2c02 {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	ppu := CreatePPU(cartridge, false, "")
	ppu.warmup = true

	return ppu
}

// Render cycles tests
func TestPPU_Render_Cycles_should_increment_scanline_after_341_cycles(t *testing.T) {
	ppu := aPPU()
	ppu.renderCycle = 340

	ppu.Tick()

	assert.Equal(t, int16(1), ppu.currentScanline)
}

func TestPPU_Render_Cycles_should_reset_scanline_after_261_scanlines(t *testing.T) {
	ppu := aPPU()
	ppu.renderCycle = 340
	ppu.currentScanline = 261

	ppu.Tick()

	assert.Equal(t, int16(0), ppu.currentScanline)
	assert.Equal(t, uint16(0), ppu.renderCycle)
}

func TestPPU_Render_Cycles_should_trigger_vblank_on_second_cycle_of_scanline_241(t *testing.T) {
	ppu := aPPU()
	ppu.currentScanline = 241
	ppu.renderCycle = 0

	ppu.Tick()
	assert.False(t, ppu.ppuStatus.verticalBlankStarted, "VBlank has been generated too early")

	ppu.Tick()
	assert.True(t, ppu.ppuStatus.verticalBlankStarted, "Should have generated VBlank on second cycle of 241 scanline")
}

func TestPPU_VBlank_should_return_disable_vblank_on_second_cycle_of_scanline_261(t *testing.T) {
	ppu := aPPU()
	ppu.ppuStatus.verticalBlankStarted = true
	ppu.currentScanline = 261
	ppu.renderCycle = 0

	ppu.Tick()
	assert.True(t, ppu.ppuStatus.verticalBlankStarted, "VBlank has been disabled too early")

	ppu.Tick()
	assert.False(t, ppu.ppuStatus.verticalBlankStarted, "Should have been disabled VBlank on second cycle of 261 scanline")
}

func Test_should_trigger_NMI_on_vBlank(t *testing.T) {
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
			ppu.renderCycle = 1
			ppu.currentScanline = 241

			ppu.Tick()

			assert.Equal(t, tt.allowNMI, ppu.nmi, "Unexpected NMI behaviour")
		})
	}
}

// Testing Render lifecycle
//func Test_should_load_next_tileId(t *testing.T) {
//	ppu := aPPU()
//
//}
