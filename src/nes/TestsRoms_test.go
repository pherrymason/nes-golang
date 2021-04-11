package nes

import (
	"testing"

	"github.com/raulferras/nes-golang/src/log"
)

func TestNestest(t *testing.T) {
	gamePak := readRom("./../../tests/roms/nestest/nestest.nes")

	//logger := log.MemoryLogger{}
	outputLogPath := "./../../var/nestest.log"
	logger := log.CreateFileLogger(outputLogPath)

	nes := CreateDebuggableNes(&logger)
	nes.InsertCartridge(&gamePak)
	nes.cpu.ResetToAddress(0xC000)
	nes.Start()

	// Compare logs
	compareLogs(outputLogPath)
}

func compareLogs(logPath string) {

}

func TestCPUDummyReads(t *testing.T) {
	t.Skip()
	cartridge := readRom("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes()
	nes.InsertCartridge(&cartridge)
	nes.Start()
}
