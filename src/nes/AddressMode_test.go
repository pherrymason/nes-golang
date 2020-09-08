package nes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmediate(t *testing.T) {
	state := AddressModeState{CreateRegisters(), &RAM{}}
	result := evalImmediate(state)

	assert.Equal(t, Address(0), result, "Immediate address mode failed to evaluate address")
}

func TestZeroPage(t *testing.T) {
	ram := RAM{}
	ram.write(0x00, 0x40)
	state := AddressModeState{CreateRegisters(), &ram}

	result := evalZeroPage(state)
	expected := Address(0x4000)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPage address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestZeroPageX(t *testing.T) {
	ram := RAM{}
	ram.write(0x00, 0x05)
	registers := CreateRegisters()
	registers.X = 0x10

	state := AddressModeState{registers, &ram}
	result := evalZeroPageX(state)

	expected := Address(0x15)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageX address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestZeroPageY(t *testing.T) {
	ram := RAM{}
	registers := CreateRegisters()
	registers.Y = 0x10
	registers.Pc = Address(0x0000)
	ram.write(registers.Pc, 0xF0)

	result := evalZeroPageY(AddressModeState{registers, &ram})

	expected := Address(0x00)

	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageY address mode decoded wrongly, expected %d, got %d", expected, result))
}

func TestAbsolute(t *testing.T) {
	ram := RAM{}
	ram.write(0x0000, 0x30)
	ram.write(0x0001, 0x01)

	registers := CreateRegisters()

	result := evalAbsolute(AddressModeState{registers, &ram})

	assert.Equal(t, Address(0x0130), result, "Error")
}

func TestAbsoluteXIndexed(t *testing.T) {
	registers := CreateRegisters()
	registers.X = 5

	ram := RAM{}
	ram.write(0x0000, 0x01)
	ram.write(0x0001, 0x01)

	result := evalAbsoluteXIndexed(AddressModeState{registers, &ram})

	assert.Equal(t, Address(0x106), result)
}

func TestAbsoluteYIndexed(t *testing.T) {
	registers := CreateRegisters()
	registers.Y = 5

	ram := RAM{}
	ram.write(0x0000, 0x01)
	ram.write(0x0001, 0x01)

	result := evalAbsoluteYIndexed(AddressModeState{registers, &ram})

	assert.Equal(t, Address(0x106), result)
}

func TestIndirect(t *testing.T) {
	registers := CreateRegisters()
	ram := RAM{}
	// Write Pointer to address 0x0134 in RAM
	ram.write(0, 0x34)
	ram.write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	ram.write(Address(0x134), 0x00)
	ram.write(Address(0x135), 0x02)

	result := evalIndirect(AddressModeState{registers, &ram})
	expected := Address(0x200)

	assert.Equal(t, expected, result, "Indirect error")
}

func TestIndirectBug(t *testing.T) {
	registers := CreateRegisters()
	ram := RAM{}
	// Write Pointer to address 0x01FF in RAM
	ram.write(0, 0xFF)
	ram.write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	ram.write(Address(0x1FF), 0x32)
	ram.write(Address(0x200), 0x04)

	result := evalIndirect(AddressModeState{registers, &ram})
	expected := Address(0x432)

	assert.Equal(t, expected, result, "Indirect error")
}

func TestPreIndexedIndirect(t *testing.T) {
	registers := CreateRegisters()
	registers.X = 4

	ram := RAM{}
	// Write Operand
	ram.write(0, 0x10)

	// Write Offset Table
	ram.write(0x0014, 0x25)

	result := evalPreIndexedIndirect(AddressModeState{registers, &ram})

	expected := Address(0x0025)
	assert.Equal(t, expected, result)
}

func TestPreIndexedIndirectWithWrapAround(t *testing.T) {
	registers := CreateRegisters()
	registers.Pc = 0x00
	registers.X = 21
	ram := RAM{}
	// Write Operan
	ram.write(0x0000, 250)

	// Write Offset Table
	ram.write(0x000F, 0x10)

	result := evalPreIndexedIndirect(AddressModeState{registers, &ram})

	expected := Address(0x0010)
	assert.Equal(t, expected, result)
}

func TestPostIndexedIndirect(t *testing.T) {
	registers := CreateRegisters()

	expected := Address(0xF0)

	registers.Y = 0x05

	ram := RAM{}
	// Opcode Operand
	ram.write(0x0000, 0x05)

	// Indexed Table Pointers
	ram.write(0x05, 0x20)
	ram.write(0x06, 0x00)

	// Offset pointer
	ram.write(0x0025, byte(expected))

	result := evalPostIndexedIndirect(AddressModeState{registers, &ram})

	assert.Equal(t, expected, result)
}

func TestPostIndexedIndirectWithWrapAround(t *testing.T) {
	t.Skip()
	registers := CreateRegisters()

	expected := Address(0xF0)

	registers.Y = 15

	ram := RAM{}
	// Opcode Operand
	ram.write(0x0000, 0x05)

	// Indexed Table Pointers
	ram.write(0x05, 0xFB)
	ram.write(0x06, 0x00)

	// Offset pointer
	ram.write(0x000A, byte(expected))

	result := evalPostIndexedIndirect(AddressModeState{registers, &ram})

	assert.Equal(t, expected, result)
}

func TestRelativeAddressMode(t *testing.T) {
	registers := CreateRegisters()
	registers.Pc = 0x10
	ram := RAM{}

	// Write Operand
	ram.write(0x09, 0xFF) // OpCode
	ram.write(0x10, 0x04) // Operand

	result := evalRelative(AddressModeState{registers, &ram})

	assert.Equal(t, Address(0x15), result)
}

func TestRelativeAddressModeNegative(t *testing.T) {
	registers := CreateRegisters()
	registers.Pc = 0x10
	ram := RAM{}

	// Write Operand
	ram.write(0x10, 0x100-4)

	result := evalRelative(AddressModeState{registers, &ram})

	assert.Equal(t, Address(0x0D), result)
}
