package nes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createBus() Bus {
	ram := RAM{}
	return Bus{Ram: &ram}
}

func TestImmediate(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = Address(0x100)
	state := CreateAddressModeState(&cpu)
	result := cpu.evalImmediate(state)

	assert.Equal(t, Address(0x100), result, "Immediate address mode failed to evaluate address")
}

func TestZeroPage(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.write(0x00, 0x40)
	state := CreateAddressModeState(&cpu)

	result := cpu.evalZeroPage(state)
	expected := Address(0x4000)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPage address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestZeroPageX(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.write(0x00, 0x05)
	cpu.registers.X = 0x10

	state := CreateAddressModeState(&cpu)
	result := cpu.evalZeroPageX(state)

	expected := Address(0x15)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageX address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestZeroPageY(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Y = 0x10
	cpu.registers.Pc = Address(0x0000)
	cpu.bus.write(cpu.registers.Pc, 0xF0)

	result := cpu.evalZeroPageY(CreateAddressModeState(&cpu))

	expected := Address(0x00)

	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageY address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestAbsolute(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.write(0x0000, 0x30)
	cpu.bus.write(0x0001, 0x01)

	result := cpu.evalAbsolute(CreateAddressModeState(&cpu))

	assert.Equal(t, Address(0x0130), result, "Error")
}

func TestAbsoluteXIndexed(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.X = 5
	cpu.registers.Pc = 0x00
	cpu.bus.write(0x0000, 0x01)
	cpu.bus.write(0x0001, 0x01)

	result := cpu.evalAbsoluteXIndexed(CreateAddressModeState(&cpu))

	assert.Equal(t, Address(0x106), result)
}

func TestAbsoluteYIndexed(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.Y = 5

	cpu.bus.write(0x0000, 0x01)
	cpu.bus.write(0x0001, 0x01)

	result := cpu.evalAbsoluteYIndexed(CreateAddressModeState(&cpu))

	assert.Equal(t, Address(0x106), result)
}

func TestIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00

	// Write Pointer to address 0x0134 in Bus
	cpu.bus.write(0, 0x34)
	cpu.bus.write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	cpu.bus.write(Address(0x134), 0x00)
	cpu.bus.write(Address(0x135), 0x02)

	result := cpu.evalIndirect(CreateAddressModeState(&cpu))
	expected := Address(0x200)

	assert.Equal(t, expected, result, "Indirect error")
}

func TestIndirectBug(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	// Write Pointer to address 0x01FF in Bus
	cpu.bus.write(0, 0xFF)
	cpu.bus.write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	cpu.bus.write(Address(0x1FF), 0x32)
	cpu.bus.write(Address(0x200), 0x04)

	result := cpu.evalIndirect(CreateAddressModeState(&cpu))
	expected := Address(0x432)

	assert.Equal(t, expected, result, "Indirect error")
}

func TestPreIndexedIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.X = 4

	// Write Operand
	cpu.bus.write(0, 0x10)

	// Write Offset Table
	cpu.bus.write(0x0014, 0x25)

	result := cpu.evalIndirectX(CreateAddressModeState(&cpu))

	expected := Address(0x0025)
	assert.Equal(t, expected, result)
}

func TestPreIndexedIndirectWithWrapAround(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.X = 21
	// Write Operan
	cpu.bus.write(0x0000, 250)

	// Write Offset Table
	cpu.bus.write(0x000F, 0x10)

	result := cpu.evalIndirectX(CreateAddressModeState(&cpu))

	expected := Address(0x0010)
	assert.Equal(t, expected, result)
}

func TestPostIndexedIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00

	expected := Address(0xF0)

	cpu.registers.Y = 0x05

	// Opcode Operand
	cpu.bus.write(0x0000, 0x05)

	// Indexed Table Pointers
	cpu.bus.write(0x05, 0x20)
	cpu.bus.write(0x06, 0x00)

	// Offset pointer
	cpu.bus.write(0x0025, byte(expected))

	result := cpu.evalIndirectY(CreateAddressModeState(&cpu))

	assert.Equal(t, expected, result)
}

func TestPostIndexedIndirectWithWrapAround(t *testing.T) {
	t.Skip()
	cpu := CreateCPUWithBus()

	expected := Address(0xF0)
	cpu.registers.Pc = 0x00
	cpu.registers.Y = 15

	// Opcode Operand
	cpu.bus.write(0x0000, 0x05)

	// Indexed Table Pointers
	cpu.bus.write(0x05, 0xFB)
	cpu.bus.write(0x06, 0x00)

	// Offset pointer
	cpu.bus.write(0x000A, byte(expected))

	result := cpu.evalIndirectY(CreateAddressModeState(&cpu))

	assert.Equal(t, expected, result)
}

func TestRelativeAddressMode(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.bus.write(0x09, 0xFF) // OpCode
	cpu.bus.write(0x10, 0x04) // Operand

	result := cpu.evalRelative(CreateAddressModeState(&cpu))

	assert.Equal(t, Address(0x15), result)
}

func TestRelativeAddressModeNegative(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.bus.write(0x10, 0x100-4)

	result := cpu.evalRelative(CreateAddressModeState(&cpu))

	assert.Equal(t, Address(0x0D), result)
}
