package nes

import (
	"bufio"
	"fmt"
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	"os"
	"testing"
)

func TestNestest(t *testing.T) {
	gamePak := gamePak2.CreateGamePakFromROMFile("./../../assets/roms/nestest/nestest.nes")
	outputLogPath := "./../../var/nestest.log"

	var limitCycles uint32 = 5004

	nes := CreateNes(
		&gamePak,
		&NesDebugger{true, nil, nil, outputLogPath, nil},
	)

	nes.StartAt(0xC000)

	var i uint32 = 1
	for {
		nes.Tick()
		opCyclesLeft := nes.cpu.opCyclesLeft
		if opCyclesLeft == 0 {
			i++
		}
		if limitCycles > 0 && i >= limitCycles {
			break
		}
	}

	nes.cpu.Logger.Close()
	// Compare logs
	compareLogs(t, nes.cpu.Logger.Snapshots())
}

func compareLogs(t *testing.T, snapshots []CpuState) {
	fmt.Println("Comparing state")

	file, err := os.Open("./../../assets/roms/nestest/nestest.log")
	if err != nil {
		fmt.Println(fmt.Errorf("could not find file nestest.log"))
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i, state := range snapshots {
		scanner.Scan()
		nesTestLine := scanner.Text()
		nesTestState := CreateStateFromNesTestLine(nesTestLine)
		if !state.RegistersEquals(nesTestState) {
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
	gamePak := gamePak2.CreateGamePakFromROMFile("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes(&gamePak, &NesDebugger{})
	nes.Start()
}
