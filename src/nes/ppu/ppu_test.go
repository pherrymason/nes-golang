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
	ppu := CreatePPU(cartridge)
	ppu.warmup = false

	return ppu
}

func aPPU() *Ppu2c02 {
	cartridge := gamePak.NewDummyGamePak(
		gamePak.NewEmptyCHRROM(),
	)
	ppu := CreatePPU(cartridge)
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

func TestPPU_Render_Cycles_should_not_trigger_vblank_before_second_cycle_before_start_of_scanline_241(t *testing.T) {
	ppu := aPPU()
	cyclesPerScanline := 341
	firstCycleOfScanline241 := 1

	for i := 0; i < (cyclesPerScanline*241 + firstCycleOfScanline241); i++ {
		ppu.Tick()

		if ppu.ppuStatus.verticalBlankStarted == true {
			assert.FailNowf(t, "VerticalBlank generated", "generated at scanline:%d, cycle:%d", ppu.currentScanline, ppu.renderCycle)
		}
	}
}

func TestPPU_Render_Cycles_should_trigger_vblank_from_scanline_241_to_261(t *testing.T) {
	ppu := aPPU()
	ppu.currentScanline = 241
	ppu.renderCycle = 1
	cyclesPerScanline := 341

	for i := 0; i < cyclesPerScanline*(261-241); i++ {
		ppu.Tick()

		if ppu.ppuStatus.verticalBlankStarted == false {
			assert.FailNowf(t, "VerticalBlank not triggered", "disabled at scanline:%d, cycle:%d", ppu.currentScanline, ppu.renderCycle)
		}
	}
}

func TestPPU_VBlank_should_return_false_when_current_scanline_is_below_241(t *testing.T) {
	ppu := CreatePPU(gamePak.NewDummyGamePak(gamePak.NewEmptyCHRROM()))
	ppu.currentScanline = 240

	assert.False(t, ppu.VBlank())
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

			//assert.True(t, ppu.ppuStatus.verticalBlankStarted)
			assert.Equal(t, tt.allowNMI, ppu.nmi, "Unexpected NMI behaviour")
		})
	}
}
