package nes

import (
	"bufio"
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	ppu2 "github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/raulferras/nes-golang/src/utils"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestNestest(t *testing.T) {
	gamePak := gamePak2.CreateGamePakFromROMFile("./../../assets/roms/tests/nestest/nestest.nes")
	outputLogPath := "./../../var"

	var limitCycles uint32 = 5004

	nes := CreateNes(
		&gamePak,
		CreateNesDebugger(outputLogPath, true, false),
	)

	nes.StartAt(0xC000)

	var i uint32 = 1
	for {
		nes.Tick()
		opCyclesLeft := nes.Cpu.opCyclesLeft
		if opCyclesLeft == 0 {
			i++
		}
		if limitCycles > 0 && i >= limitCycles {
			break
		}
	}

	nes.Cpu.debugger.Stop()
	// Compare logs
	compareLogs(t, nes.Cpu.debugger.Logger.Snapshots())
}

func compareLogs(t *testing.T, snapshots []cpu.Snapshot) {
	fmt.Println("Comparing state")

	file, err := os.Open("./../../assets/roms/tests/nestest/nestest.log")
	if err != nil {
		fmt.Println(fmt.Errorf("could not find file nestest.log"))
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i, snapshot := range snapshots {
		scanner.Scan()
		nesTestLine := scanner.Text()
		nesTestSnapshot := CreateSnapshotFromNesTestLine(nesTestLine)
		if !snapshot.CpuState.RegistersEquals(nesTestSnapshot.CpuState) {
			msg := fmt.Sprintf("Error in iteration %d\n", i+1)
			msg += fmt.Sprintf("Expected: %s\n", nesTestSnapshot.CpuState.ToString())
			msg += fmt.Sprintf("Actual: %s\n", snapshot.CpuState.ToString())

			t.Errorf(msg)
			t.FailNow()
		}

		if nesTestSnapshot.CpuState.CyclesSinceReset != snapshot.CpuState.CyclesSinceReset {
			t.Errorf("Error in iteration %d\nCPU Cycles don't match. Expected: %d. Actual: %d", i+1, nesTestSnapshot.CpuState.CyclesSinceReset, snapshot.CpuState.CyclesSinceReset)
			t.FailNow()
		}

		if nesTestSnapshot.PpuState.RenderCycle != snapshot.PpuState.RenderCycle {
			t.Errorf("Error in iteration %d\nPPU X doesn't match. Expected: %d. Actual: %d", i+1, nesTestSnapshot.PpuState.RenderCycle, snapshot.PpuState.RenderCycle)
			t.FailNow()
		}
	}
}

func TestCPUDummyReads(t *testing.T) {
	t.Skip()
	gamePak := gamePak2.CreateGamePakFromROMFile("./../../tests/roms/cpu_dummy_reads.nes")

	nes := CreateNes(&gamePak, &Debugger{})
	nes.Start()
}

func CreateSnapshotFromNesTestLine(nesTestLine string) cpu.Snapshot {
	tokens := strings.Fields(nesTestLine)
	//_ = opCodeTokens

	blocks := utils.StringSplitByRegex(nesTestLine)

	result := utils.HexStringToByteArray(blocks[0])
	pc := types.CreateAddress(result[1], result[0])

	opCodeTokens := strings.Fields(blocks[1])
	opcode := [3]byte{utils.HexStringToByteArray(opCodeTokens[0])[0]}

	flagFields := strings.Fields(blocks[3])

	r, _ := regexp.Compile("CYC:([0-9]+)$")
	cpuCyclesString := r.FindStringSubmatch(nesTestLine)

	cpuCycles, _ := strconv.ParseUint(cpuCyclesString[1], 10, 16)

	cpuState := cpu.CreateState(
		cpu.Registers{
			utils.NestestDecodeRegisterFlag(flagFields[0]),
			utils.NestestDecodeRegisterFlag(flagFields[1]),
			utils.NestestDecodeRegisterFlag(flagFields[2]),
			pc,
			utils.NestestDecodeRegisterFlag(flagFields[4]),
			utils.NestestDecodeRegisterFlag(flagFields[3]),
		},
		opcode,
		cpu.CreateInstruction(
			strings.Fields(blocks[2])[0],
			cpu.Implicit,
			nil,
			0,
			0,
		),
		cpu.OperationMethodArgument{0, 0},
		uint32(cpuCycles),
	)

	ppuXIndex := len(tokens) - 2
	var ppuX string
	if strings.Contains(tokens[ppuXIndex], ",") {
		s := tokens[ppuXIndex]
		ppuX = strings.Split(s, ",")[1]
	} else {
		ppuX = tokens[ppuXIndex]
	}
	nesTestPpuX, _ := strconv.ParseUint(ppuX, 10, 16)
	ppu := ppu2.SimplePPUState{0, uint16(nesTestPpuX), 0}

	return cpu.Snapshot{cpuState, ppu}
}
