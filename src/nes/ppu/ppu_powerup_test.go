package ppu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPU_writing_to_registers_are_ignored_first_29658_CPU_clocks(t *testing.T) {
	//PPUCTRL, PPUMASK, PPUSCROLL, PPUADDR
	ppu := aPPU()

	for cpuCycles := 0; cpuCycles < 29658; cpuCycles += 3 {
		ppu.WriteRegister(PPUCTRL, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.ppuControl.value(), "Writes to PPUCTRL should be ignored first 30000 cycles")

		ppu.WriteRegister(PPUMASK, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.mask, "Writes to PPUMASK should be ignored first 30000 cycles")

		ppu.WriteRegister(PPUSCROLL, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.ppuScroll.scrollX, "Writes to PPUSCROLL should be ignored first 30000 cycles, scrollX was modified")
		assert.NotEqual(t, 0xFF, ppu.ppuScroll.scrollY, "Writes to PPUSCROLL should be ignored first 30000 cycles, scrollY was modified")
		assert.NotEqual(t, 0x00, ppu.ppuScroll.latch, "Writes to PPUSCROLL should be ignored first 30000 cycles, scroll latch was modified")

		ppu.WriteRegister(PPUADDR, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.ppuDataAddr, "Writes to PPUADDR should be ignored first 30000 cycles")
	}
}

func TestPPU_writing_to_registers_are_ready_first_29658_CPU_clocks(t *testing.T) {
	//PPUSTATUS, OAMADDR, OAMDATA ($2004), PPUDATA, and OAMDMA
	ppu := aPPU()

	for cpuCycles := 0; cpuCycles < 29658; cpuCycles += 3 {
		//ppu.WriteRegister(PPUSTATUS, 0xFF)
		//assert.NotEqual(t, 0xFF, ppu.registers.status)

		ppu.WriteRegister(OAMADDR, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.oamAddr, fmt.Sprintf("OAMAddr was not 0xFF at cycle %d", cpuCycles))

		ppu.WriteRegister(OAMDATA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.ReadRegister(OAMDATA), fmt.Sprintf("OAMData was not 0xFF at cycle %d", cpuCycles))

		ppu.registers.ppuDataAddr = 0x2000
		ppu.WriteRegister(PPUDATA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.Read(0x00), fmt.Sprintf("PPUDATA was not 0xFF at cycle %d", cpuCycles))

		ppu.WriteRegister(OAMDMA, 0xFF)
		assert.NotEqual(t, 0xFF, ppu.registers.ppuDataAddr, fmt.Sprintf("OAMDMA was not 0xFF at cycle %d", cpuCycles))
	}
}
