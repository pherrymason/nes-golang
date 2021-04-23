package nes

import (
	"bufio"
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"os"
	"testing"
)

func TestNestest(t *testing.T) {
	gamePak := ReadRom("./../../roms/nestest/nestest.nes")

	outputLogPath := "./../../var/nestest.log"
	logger := cpu.CreateCPULogger(outputLogPath)

	nes := CreateNes(NesDebugger{true, nil, &logger, 5004, nil})
	nes.InsertGamePak(&gamePak)
	nes.StartAt(0xC000)

	var i uint16 = 1
	for {
		opCyclesLeft := nes.Tick()
		if opCyclesLeft == 0 {
			i++
		}
		if nes.debug.cyclesLimit > 0 && i >= nes.debug.cyclesLimit {
			break
		}
	}

	// Compare logs
	compareLogs(t, nes.cpu.Logger.Snapshots())
}

func compareLogs(t *testing.T, snapshots []cpu.State) {
	fmt.Println("Comparing state")

	file, err := os.Open("./../../roms/nestest/nestest.log")
	if err != nil {
		fmt.Println(fmt.Errorf("could not find file nestest.log"))
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i, state := range snapshots {
		scanner.Scan()
		nesTestLine := scanner.Text()
		nesTestState := cpu.CreateStateFromNesTestLine(nesTestLine)
		if !state.Equals(nesTestState) {
			msg := fmt.Sprintf("Error in iteration %d\n", i+1)
			msg += fmt.Sprintf("Expected: %s\n", nesTestState.ToString())
			msg += fmt.Sprintf("Actual: %s\n", state.ToString())

			t.Errorf(msg)
			t.FailNow()
		}
	}
}

func TestCPUDummyReads(t *testing.T) {
	t.Skip()
	cartridge := ReadRom("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes(NesDebugger{})
	nes.InsertGamePak(&cartridge)
	nes.Start()
}
