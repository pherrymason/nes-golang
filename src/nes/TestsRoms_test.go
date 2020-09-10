package nes

import (
	"testing"

	"github.com/raulferras/nes-golang/src/log"
)

func TestNestest(t *testing.T) {
	gamepak := readRom("./../../tests/roms/nestest/nestest.nes")

	logger := log.MemoryLogger{}

	nes := CreateDebugableNes(&logger)
	nes.InsertCartridge(&gamepak)
	nes.cpu.reset()
	nes.cpu.registers.Pc = 0xC000
	nes.Start()
}

func TestCPUDummyReads(t *testing.T) {
	t.Skip()
	cartridge := readRom("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes()
	nes.InsertCartridge(&cartridge)
	nes.Start()
}
