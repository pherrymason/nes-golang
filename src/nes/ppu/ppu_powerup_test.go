package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPU_writing_to_registers_are_ignored_first_29658_CPU_clocks(t *testing.T) {
	//PPUCTRL, PPUMASK, PPUSCROLL, PPUADDR
	ppu := newNotWarmedUpPPU()

	for cpuCycles := 0; cpuCycles < 29658; cpuCycles += 3 {
		ppu.WriteRegister(PPUCTRL, 0xFF)
		if byte(0xFF) == ppu.ppuControl.value() {
			assert.FailNowf(t, "ppuctrl write was not ignored", "Writes to PPUCTRL should be ignored first 30000 cycles. (Cycle :%d)", cpuCycles)
		}

		ppu.WriteRegister(PPUMASK, 0xFF)
		if byte(0xFF) == ppu.ppuMask.value() {
			assert.FailNowf(t, "", "Writes to PPUMASK should be ignored first 30000 cycles. (Cycle :%d)", cpuCycles)
		}

		ppu.WriteRegister(PPUSCROLL, 0xFF)
		if types.Address(0xFF) == ppu.tRam.value() {
			assert.FailNowf(t, "", "Writes to PPUSCROLL should be ignored first 30000 cycles, scrollX was modified (Cycle :%d)", cpuCycles)
		}
		if byte(0xFF) == ppu.fineX {
			assert.FailNowf(t, "", "Writes to PPUSCROLL should be ignored first 30000 cycles, scrollY was modified. (Cycle :%d)", cpuCycles)
		}
		if byte(0) != ppu.tRam.latch {
			assert.FailNowf(t, "", "Writes to PPUSCROLL should be ignored first 30000 cycles, scroll latch was modified. (Cycle :%d)", cpuCycles)
		}

		ppu.WriteRegister(PPUADDR, 0xFF)
		if types.Address(0xFF) == ppu.tRam.address() {
			assert.FailNowf(t, "", "Writes to PPUADDR should be ignored first 30000 cycles, address changed. (Cycle :%d)", cpuCycles)
		}
		if byte(0x0) != ppu.tRam.latch {
			assert.FailNowf(t, "", "rites to PPUADDR should be ignored first 30000 cycles, latch changed. (Cycle :%d)", cpuCycles)
		}

		ppu.Tick()
	}
}

func TestPPU_writing_to_registers_are_ready_first_29658_CPU_clocks(t *testing.T) {
	//PPUSTATUS, OAMADDR, OAMDATA ($2004), PPUDATA, and OAMDMA
	ppu := newNotWarmedUpPPU()
	ppu.cycle = 29658

	//ppu.WriteRegister(PPUSTATUS, 0xFF)
	//assert.NotEqual(t, 0xFF, ppu.registers.status)

	ppu.WriteRegister(OAMADDR, 0xFF)
	assert.NotEqual(t, 0xFF, ppu.oamAddr, "OAMAddr was not 0xFF")

	ppu.WriteRegister(OAMDATA, 0xFF)
	assert.NotEqual(t, 0xFF, ppu.ReadRegister(OAMDATA), "OAMData was not 0xFF")

	ppu.vRam.setValue(0x2000)
	ppu.WriteRegister(PPUDATA, 0xFF)
	assert.NotEqual(t, 0xFF, ppu.Read(0x00), "PPUDATA was not 0xFF")

	// Not implemented!
	//ppu.WriteRegister(OAMDMA, 0xFF)
	//assert.NotEqual(t, 0xFF, ppu.oamAddr, fmt.Sprintf("OAMDMA was not 0xFF at cycle %d", cpuCycles))
}
