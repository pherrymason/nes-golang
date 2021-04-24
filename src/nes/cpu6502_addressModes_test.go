package nes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmediate(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = Address(0x100)

	pc, address, _, pageCrossed := cpu.evalImmediate(cpu.registers.Pc)

	assert.Equal(t, Address(0x100), address, "Immediate address mode failed to evaluate address")
	assert.Equal(t, cpu.registers.Pc+1, pc)
	assert.False(t, pageCrossed)

}

func TestZeroPage(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x00, 0x40)

	pc, result, _, pageCrossed := cpu.evalZeroPage(cpu.registers.Pc)
	expected := Address(0x040)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPage address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
	assert.False(t, pageCrossed)
}

func TestZeroPageX(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x00, 0x05)
	cpu.registers.X = 0x10

	state := cpu.registers.Pc
	pc, result, _, pageCrossed := cpu.evalZeroPageX(state)

	expected := Address(0x15)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageX address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
	assert.False(t, pageCrossed)
}

func TestZeroPageY(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Y = 0x10
	cpu.registers.Pc = Address(0x0000)
	cpu.memory.Write(cpu.registers.Pc, 0xF0)

	pc, result, _, pageCrossed := cpu.evalZeroPageY(cpu.registers.Pc)

	expected := Address(0x00)

	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageY address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
	assert.False(t, pageCrossed)
}

func TestAbsolute(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x0000, 0x30)
	cpu.memory.Write(0x0001, 0x01)

	pc, result, _, pageCrossed := cpu.evalAbsolute(cpu.registers.Pc)

	assert.Equal(t, Address(0x0130), result, "Error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
	assert.False(t, pageCrossed)
}

func TestAbsoluteXIndexed(t *testing.T) {
	type dataProvider struct {
		test                string
		lsb                 byte
		hsb                 byte
		expectedAddress     Address
		expectedPageCrossed bool
	}

	dataProviders := [...]dataProvider{
		{"page not crossed", 0x01, 0x01, 0x106, false},
		{"page crossed", 0xFF, 0x00, 0x104, true},
	}

	for _, dp := range dataProviders {
		cpu := CreateCPUWithGamePak()
		cpu.registers.X = 5
		cpu.registers.Pc = 0x00
		//cpu.memory.Write(0x0000, 0x01)
		//cpu.memory.Write(0x0001, 0x01)
		cpu.memory.Write(0x0000, dp.lsb)
		cpu.memory.Write(0x0001, dp.hsb)

		pc, result, _, pageCrossed := cpu.evalAbsoluteXIndexed(cpu.registers.Pc)

		assert.Equal(t, dp.expectedAddress, result)
		assert.Equal(t, cpu.registers.Pc+2, pc)
		assert.Equal(t, dp.expectedPageCrossed, pageCrossed)
	}
}

func TestAbsoluteYIndexed(t *testing.T) {
	type dataProvider struct {
		test                string
		lsb                 byte
		hsb                 byte
		expectedAddress     Address
		expectedPageCrossed bool
	}

	dataProviders := [...]dataProvider{
		{"page not crossed", 0x01, 0x01, 0x106, false},
		{"page crossed", 0xFF, 0x00, 0x104, true},
	}

	for _, dp := range dataProviders {
		cpu := CreateCPUWithGamePak()
		cpu.registers.Pc = 0x00
		cpu.registers.Y = 5

		//cpu.memory.Write(0x0000, 0x01)
		//cpu.memory.Write(0x0001, 0x01)

		cpu.memory.Write(0x0000, dp.lsb)
		cpu.memory.Write(0x0001, dp.hsb)

		pc, result, _, pageCrossed := cpu.evalAbsoluteYIndexed(cpu.registers.Pc)

		assert.Equal(t, dp.expectedAddress, result)
		assert.Equal(t, dp.expectedPageCrossed, pageCrossed)
		assert.Equal(t, cpu.registers.Pc+2, pc)
	}
}

func TestIndirect(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00

	// Write Pointer to address 0x0134 in Bus
	cpu.memory.Write(0, 0x34)
	cpu.memory.Write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	cpu.memory.Write(Address(0x134), 0x00)
	cpu.memory.Write(Address(0x135), 0x02)

	pc, result, _, _ := cpu.evalIndirect(cpu.registers.Pc)
	expected := Address(0x200)

	assert.Equal(t, expected, result, "Indirect error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestIndirect_bug(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	// Write Pointer to address 0xC0FF in Bus
	cpu.memory.Write(0, 0xFF)
	cpu.memory.Write(1, 0xC0)

	cpu.memory.Write(0xC0FF, 0x55)
	cpu.memory.Write(0xC100, 0x04)
	cpu.memory.Write(0xC000, 0x01)

	pc, result, _, _ := cpu.evalIndirect(cpu.registers.Pc)
	expected := Address(0x155)

	assert.Equal(t, expected, result, "Indirect error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestIndexed_indirect(t *testing.T) {
	type dataProvider struct {
		test                string
		x                   byte
		operand             byte
		lsbPtr              Address
		hsbPtr              Address
		expectedAddress     Address
		expectedPageCrossed bool
	}

	dataProviders := [...]dataProvider{
		{"page not crossed", 0x04, 0x10, 0x0014, 0x0015, 0x1025, false},
		{"page crossed", 0x03, 0xFE, 0x01, 0x02, 0x510, false},
		{"page crossed", 0x01, 0xFE, 0xFF, 0x00, 0x510, false},
	}

	for testIdx, dp := range dataProviders {
		cpu := CreateCPUWithGamePak()
		cpu.registers.Pc = 0x00
		cpu.registers.X = 4

		// Write Operand
		cpu.memory.Write(0, 0x10)

		// Write Offset Table
		cpu.memory.Write(0x0014, byte(dp.expectedAddress))
		cpu.memory.Write(0x0015, byte(dp.expectedAddress>>8))

		_, result, _, pageCrossed := cpu.evalIndirectX(cpu.registers.Pc)

		assert.Equal(t, dp.expectedAddress, result, fmt.Sprintf("%s [%d]: address", dp.test, testIdx))
		assert.Equal(t, dp.expectedPageCrossed, pageCrossed, fmt.Sprintf("%s [%d]: page crossed", dp.test, testIdx))
	}
}

func TestIndirectY(t *testing.T) {
	type dataProvider struct {
		test                string
		y                   byte
		operand             byte
		lsbPtr              Address
		hsbPtr              Address
		expectedAddress     Address
		expectedPageCrossed bool
	}

	dataProviders := [...]dataProvider{
		{"page not crossed", 0x10, 0x86, 0x86, 0x87, 0x4038, false},
		{"page crossed", 0xFF, 0xFF, 0xFF, 0x00, 0x245, true},
	}

	for testIdx, dp := range dataProviders {
		cpu := CreateCPUWithGamePak()
		cpu.registers.Pc = 0x50
		cpu.registers.Y = dp.y

		// Operand
		cpu.memory.Write(0x50, dp.operand)

		// Indexed Table Pointers
		cpu.memory.Write(Address(dp.operand), byte(dp.expectedAddress-Address(dp.y)))
		cpu.memory.Write(Address(dp.operand+1), byte((dp.expectedAddress-Address(dp.y))>>8))

		pc, result, _, pageCrossed := cpu.evalIndirectY(cpu.registers.Pc)

		assert.Equal(t, dp.expectedAddress, result, fmt.Sprintf("%s [%d]: address", dp.test, testIdx))
		assert.Equal(t, cpu.registers.Pc+1, pc, fmt.Sprintf("%s [%d]: pc", dp.test, testIdx))
		assert.Equal(t, dp.expectedPageCrossed, pageCrossed, fmt.Sprintf("%s [%d]: page crossed", dp.test, testIdx))
	}
}

func TestRelativeAddressMode(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.memory.Write(0x09, 0xFF) // OpCode
	cpu.memory.Write(0x10, 0x04) // Operand

	pc, result, _, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, Address(0x15), result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestRelativeAddressModeNegative(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.memory.Write(0x10, 0x100-4)

	pc, result, _, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, Address(0x0D), result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}
