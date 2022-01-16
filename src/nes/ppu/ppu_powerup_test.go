package ppu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPU_writing_to_registers_are_ignored_first_29658_CPU_clocks(t *testing.T) {
	//PPUCTRL, PPUMASK, PPUSCROLL, PPUADDR
	ppu := aPPU()

	for cpuCycles := 0; cpuCycles < 29658; cpuCycles += 3 {
		ppu.WriteRegister(PPUCTRL, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.ctrl)

		ppu.WriteRegister(PPUMASK, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.mask)

		ppu.WriteRegister(PPUSCROLL, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.scrollX)
		assert.NotEqual(t, 0xFF, ppu.registers.scrollY)

		ppu.WriteRegister(PPUADDR, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.ppuAddr)
	}
}

func TestPPU_writing_to_registers_are_ready_first_29658_CPU_clocks(t *testing.T) {
	//PPUSTATUS, OAMADDR, OAMDATA ($2004), PPUDATA, and OAMDMA
	ppu := aPPU()

	for cpuCycles := 0; cpuCycles < 29658; cpuCycles += 3 {
		//ppu.WriteRegister(PPUSTATUS, 0xFF)
		//assert.NotEqual(t, 0xFF, ppu.registers.status)

		ppu.WriteRegister(OAMADDR, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.oamAddr)

		ppu.WriteRegister(OAMDATA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.ReadRegister(OAMDATA))

		ppu.WriteRegister(PPUDATA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.Read(0x00))

		ppu.WriteRegister(OAMDMA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.ppuAddr)
	}
}
