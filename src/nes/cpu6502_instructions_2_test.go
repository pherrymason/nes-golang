package nes

import (
	"encoding/json"
	"fmt"
	"github.com/raulferras/nes-golang/src/mocks"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestCpuInstructions(t *testing.T) {
	for _, filename := range findTestsFiles() {
		//if i > 1 {
		//	break
		//}

		if isUnofficialOpcode(filename) {
			t.Log("Skipped unofficial " + filename)
			continue
		}

		t.Run(filename, func(t *testing.T) {
			jsonDoc := readJSONTestFile(filename)
			cpuTests := decodeJson(jsonDoc)

			for _, cpuTest := range cpuTests {
				cpu := createCPUFromCPUState(cpuTest.Initial)
				strOp := cpuTest.Name[0:2]
				op, _ := strconv.ParseUint(strOp, 16, 8)
				operation := cpu.GetOperation(byte(op))

				t.Run(operation.Name()+" "+cpuTest.Name, func(t *testing.T) {
					cpu.Tick()

					assertExpectedCpuStatus(t, cpuTest.Final, cpu, cpuTest.Name+"("+operation.Name()+")")
				})
			}
		})
	}
}

func isUnofficialOpcode(filename string) bool {
	unofficialList := []string{
		"02.json",
		"03.json",
		"04.json",
		"07.json",
		"0b.json",
		"0c.json",
		"0f.json",
		"12.json",
		"13.json",
		"14.json",
		"17.json",
		"1a.json",
		"1b.json",
		"1c.json",
		"1f.json",
		"22.json",
		"23.json",
		"27.json",
		"2a.json",
		"2f.json",
		"32.json",
		"33.json",
		"34.json",
		"37.json",
		"3a.json",
		"3b.json",
		"3c.json",
		"3f.json",
	}

	for _, banned := range unofficialList {
		if strings.HasSuffix(filename, banned) {
			return true
		}
	}

	return false
}

func createCPUFromCPUState(state CPUState) *Cpu6502 {
	cpu := CreateCPU(
		mocks.NewSimpleMemory(),
		Cpu6502DebugOptions{false, ""},
	)

	// Apply initial state
	cpu.registers.Pc = types.Address(state.Pc)
	cpu.registers.A = state.A
	cpu.registers.X = state.X
	cpu.registers.Y = state.Y
	cpu.registers.Status = state.P
	cpu.registers.Sp = state.S

	for _, ram := range state.Ram {
		cpu.memory.Write(types.Address(ram[0]), byte(ram[1]))
	}

	return cpu
}

func assertExpectedCpuStatus(t *testing.T, expected CPUState, cpu *Cpu6502, operation string) {
	assert.Equal(t, types.Address(expected.Pc), cpu.registers.Pc, operation+" Invalid Pc")
	assert.Equal(t, expected.A, cpu.registers.A, operation+" Invalid A")
	assert.Equal(t, expected.X, cpu.registers.X, operation+" Invalid X")
	assert.Equal(t, expected.Y, cpu.registers.Y, operation+" Invalid Y")
	assert.Equal(t, expected.S, cpu.registers.Sp, operation+" Invalid Stack Pointer")
	assert.Equal(t, expected.P, cpu.registers.Status, fmt.Sprintf("%s Invalid Status: %b => %b", operation, expected.P, cpu.registers.Status))

	for _, ram := range expected.Ram {
		assert.Equal(
			t,
			byte(ram[1]), cpu.memory.Read(types.Address(ram[0])),
			fmt.Sprintf(operation+" Invalid value in memory @%X", types.Address(ram[0])),
		)
	}
}

type CPUTest struct {
	Name    string
	Initial CPUState
	Final   CPUState
}

type CPUState struct {
	Pc  int
	S   byte
	A   byte
	X   byte
	Y   byte
	P   byte
	Ram [][2]int
}

func decodeJson(jsonDoc string) []CPUTest {
	var test []CPUTest
	json.Unmarshal([]byte(jsonDoc), &test)

	return test
}

func readJSONTestFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := new(strings.Builder)
	io.Copy(b, f)

	return b.String()
}

func findTestsFiles() []string {
	var tests []string
	filepath.WalkDir("../../assets/tests/tomharte-processortests/nes6502/v1/", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".json" {
			tests = append(tests, s)
		}
		return nil
	})

	return tests
}
