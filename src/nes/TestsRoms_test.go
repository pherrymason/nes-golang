package nes

import (
	"bufio"
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"os"
	"testing"
)

func TestNestest(t *testing.T) {
	gamePak := readRom("./../../tests/roms/nestest/nestest.nes")

	outputLogPath := "./../../var/nestest.log"
	logger := cpu.CreateCPULogger(outputLogPath)

	nes := CreateDebuggableNes(DebuggableNes{true, &logger, 1000})
	nes.InsertCartridge(&gamePak)
	nes.cpu.ResetToAddress(0xC000)
	nes.Start()

	// Compare logs
	compareLogs(t, nes.cpu.Logger.Snapshots())
}

func compareLogs(t *testing.T, snapshots []cpu.State) {
	fmt.Println("Comparing state")

	file, err := os.Open("./../../tests/roms/nestest/nestest.log")
	if err != nil {
		fmt.Println(fmt.Errorf("could not find file nestest.log"))
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i, state := range snapshots {
		scanner.Scan()
		nesTestLine := scanner.Text()
		nesTestState := cpu.CreateStateFromNesTestLine(nesTestLine)

		if !state.Equals(nesTestState) {
			msg := fmt.Sprintf("Error in iteration %d\n", i)
			msg += fmt.Sprintf("Expected: %s\n", nesTestState.ToString())
			msg += fmt.Sprintf("Actual: %s\n", state.ToString())

			t.Errorf(msg)
			t.FailNow()
		}
	}
}

func TestCPUDummyReads(t *testing.T) {
	t.Skip()
	cartridge := readRom("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes()
	nes.InsertCartridge(&cartridge)
	nes.Start()
}
