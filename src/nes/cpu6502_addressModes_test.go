package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmediate(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = types.Address(0x100)

	address, _, _, pageCrossed := cpu.evalImmediate(cpu.registers.Pc)

	assert.Equal(t, types.Address(0x100), address, "Immediate address mode failed to evaluate address")
	assert.False(t, pageCrossed)

}

func TestZeroPage(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x00, 0x40)

	result, _, _, pageCrossed := cpu.evalZeroPage(cpu.registers.Pc)
	expected := types.Address(0x040)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPage address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.False(t, pageCrossed)
}

func TestZeroPageX(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x00, 0x05)
	cpu.registers.X = 0x10

	state := cpu.registers.Pc
	result, _, _, pageCrossed := cpu.evalZeroPageX(state)

	expected := types.Address(0x15)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageX address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.False(t, pageCrossed)
}

func TestZeroPageY(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Y = 0x10
	cpu.registers.Pc = types.Address(0x0000)
	cpu.memory.Write(cpu.registers.Pc, 0xF0)

	result, _, _, pageCrossed := cpu.evalZeroPageY(cpu.registers.Pc)

	expected := types.Address(0x00)

	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageY address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.False(t, pageCrossed)
}

func TestAbsolute(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x00
	cpu.memory.Write(0x0000, 0x30)
	cpu.memory.Write(0x0001, 0x01)

	result, _, _, pageCrossed := cpu.evalAbsolute(cpu.registers.Pc)

	assert.Equal(t, types.Address(0x0130), result, "Error")
	assert.False(t, pageCrossed)
}

func TestAbsoluteXIndexed(t *testing.T) {
	cases := []struct {
		name                string
		lsb                 byte
		hsb                 byte
		expectedAddress     types.Address
		expectedPageCrossed bool
	}{
		{"page not crossed", 0x01, 0x01, 0x106, false},
		{"page crossed", 0xFF, 0x00, 0x104, true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.X = 5
			cpu.registers.Pc = 0x00
			cpu.memory.Write(0x0000, tt.lsb)
			cpu.memory.Write(0x0001, tt.hsb)

			result, _, _, pageCrossed := cpu.evalAbsoluteXIndexed(cpu.registers.Pc)

			assert.Equal(t, tt.expectedAddress, result)
			assert.Equal(t, tt.expectedPageCrossed, pageCrossed)
		})
	}
}

func TestAbsoluteYIndexed(t *testing.T) {
	cases := []struct {
		test                string
		lsb                 byte
		hsb                 byte
		expectedAddress     types.Address
		expectedPageCrossed bool
	}{
		{"page not crossed", 0x01, 0x01, 0x106, false},
		{"page crossed", 0xFF, 0x00, 0x104, true},
	}

	for _, tt := range cases {
		t.Run(tt.test, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = 0x00
			cpu.registers.Y = 5

			cpu.memory.Write(0x0000, tt.lsb)
			cpu.memory.Write(0x0001, tt.hsb)

			result, _, _, pageCrossed := cpu.evalAbsoluteYIndexed(cpu.registers.Pc)

			assert.Equal(t, tt.expectedAddress, result)
			assert.Equal(t, tt.expectedPageCrossed, pageCrossed)
		})
	}
}

func TestIndirect(t *testing.T) {
	cases := []struct {
		name           string
		pointerAddress types.Address
		finalAddress   types.Address
	}{
		{"normal indirect jump", types.Address(0x0134), types.Address(0x200)},
		/*
			6502 has indirect jump operation. It's expressed as JMP (address), so JMP ($C000) would first fetch the target address from $C000 (low byte) and $C001 (high byte), then jump to that address.

			However, 6502 is buggy what comes to page boundary: if we were to do JMP ($C0FF), the high byte will not be fetched from $C100 as expected, but from $C000!
		*/
		{"indirect jump across page", types.Address(0x00FF), types.Address(0x155)},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = 0x10        // Arbitrary pc to not collide with wraparoun bug @
			cpu.memory.Write(0x0000, 0x01) // This is to pass jump with bug

			cpu.memory.Write(0x10, byte(tt.pointerAddress))
			cpu.memory.Write(0x11, byte(tt.pointerAddress>>8))

			cpu.memory.Write(tt.pointerAddress, byte(tt.finalAddress))
			cpu.memory.Write(tt.pointerAddress+1, byte(tt.finalAddress>>8))

			result, _, _, _ := cpu.evalIndirect(cpu.registers.Pc)
			assert.Equal(t, tt.finalAddress, result, "Indirect error")
		})
	}
}

func TestIndexed_indirect(t *testing.T) {
	cases := []struct {
		name                string
		x                   byte
		operand             byte
		lsbPtr              types.Address
		hsbPtr              types.Address
		expectedAddress     types.Address
		expectedPageCrossed bool
	}{
		{"page not crossed", 0x04, 0x10, 0x0014, 0x0015, 0x1025, false},
		{"page crossed", 0x03, 0xFE, 0x01, 0x02, 0x510, false},
		{"page crossed", 0x01, 0xFE, 0xFF, 0x00, 0x510, false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = 0x00
			cpu.registers.X = 4

			// WritePrgROM Operand
			cpu.memory.Write(0, 0x10)

			// WritePrgROM Offset Table
			cpu.memory.Write(0x0014, byte(tt.expectedAddress))
			cpu.memory.Write(0x0015, byte(tt.expectedAddress>>8))

			result, _, _, pageCrossed := cpu.evalIndirectX(cpu.registers.Pc)

			assert.Equal(t, tt.expectedAddress, result, "unexpected address")
			assert.Equal(t, tt.expectedPageCrossed, pageCrossed, "page crossed")
		})
	}
}

func TestIndirectY(t *testing.T) {
	cases := []struct {
		name                string
		y                   byte
		operand             byte
		lsbPtr              types.Address
		hsbPtr              types.Address
		expectedAddress     types.Address
		expectedPageCrossed bool
	}{
		{"page not crossed", 0x10, 0x86, 0x86, 0x87, 0x4038, false},
		{"page crossed", 0xFF, 0xFF, 0xFF, 0x00, 0x245, true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = 0x50
			cpu.registers.Y = tt.y

			// Operand
			cpu.memory.Write(0x50, tt.operand)

			// Indexed Table Pointers
			cpu.memory.Write(types.Address(tt.operand), byte(tt.expectedAddress-types.Address(tt.y)))
			cpu.memory.Write(types.Address(tt.operand+1), byte((tt.expectedAddress-types.Address(tt.y))>>8))

			result, _, _, pageCrossed := cpu.evalIndirectY(cpu.registers.Pc)

			assert.Equal(t, tt.expectedAddress, result, "unexpected address")
			assert.Equal(t, tt.expectedPageCrossed, pageCrossed, "page crossed")
		})
	}
}

func TestRelativeAddressMode(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x10

	// WritePrgROM Operand
	cpu.memory.Write(0x09, 0xFF) // OpCode
	cpu.memory.Write(0x10, 0x04) // Operand

	result, _, _, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, types.Address(0x15), result)
}

func TestRelativeAddressModeNegative(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = 0x10

	// WritePrgROM Operand
	cpu.memory.Write(0x10, 0x100-4)

	result, _, _, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, types.Address(0x0D), result)
}
