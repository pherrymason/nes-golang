package cpu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/component"
	"github.com/raulferras/nes-golang/src/nes/defs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createBus() component.Bus {
	ram := component.RAM{}
	return component.Bus{Ram: &ram}
}

func TestImmediate(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = defs.Address(0x100)

	pc, address, _ := cpu.evalImmediate(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x100), address, "Immediate address mode failed to evaluate address")
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestZeroPage(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.Write(0x00, 0x40)

	pc, result, _ := cpu.evalZeroPage(cpu.registers.Pc)
	expected := defs.Address(0x4000)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPage address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestZeroPageX(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.Write(0x00, 0x05)
	cpu.registers.X = 0x10

	state := cpu.registers.Pc
	pc, result, _ := cpu.evalZeroPageX(state)

	expected := defs.Address(0x15)
	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageX address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestZeroPageY(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Y = 0x10
	cpu.registers.Pc = defs.Address(0x0000)
	cpu.bus.Write(cpu.registers.Pc, 0xF0)

	pc, result, _ := cpu.evalZeroPageY(cpu.registers.Pc)

	expected := defs.Address(0x00)

	assert.Equal(t, expected, result, fmt.Sprintf("ZeroPageY address mode decoded wrongly, expected %d, got %d", expected, result))
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestAbsolute(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.bus.Write(0x0000, 0x30)
	cpu.bus.Write(0x0001, 0x01)

	pc, result, _ := cpu.evalAbsolute(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x0130), result, "Error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestAbsoluteXIndexed(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.X = 5
	cpu.registers.Pc = 0x00
	cpu.bus.Write(0x0000, 0x01)
	cpu.bus.Write(0x0001, 0x01)

	pc, result, _ := cpu.evalAbsoluteXIndexed(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x106), result)
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestAbsoluteYIndexed(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.Y = 5

	cpu.bus.Write(0x0000, 0x01)
	cpu.bus.Write(0x0001, 0x01)

	pc, result, _ := cpu.evalAbsoluteYIndexed(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x106), result)
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00

	// Write Pointer to address 0x0134 in Bus
	cpu.bus.Write(0, 0x34)
	cpu.bus.Write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	cpu.bus.Write(defs.Address(0x134), 0x00)
	cpu.bus.Write(defs.Address(0x135), 0x02)

	pc, result, _ := cpu.evalIndirect(cpu.registers.Pc)
	expected := defs.Address(0x200)

	assert.Equal(t, expected, result, "Indirect error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestIndirectBug(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	// Write Pointer to address 0x01FF in Bus
	cpu.bus.Write(0, 0xFF)
	cpu.bus.Write(1, 0x01)

	// Write 0x0134 with final Address(0x200)
	cpu.bus.Write(defs.Address(0x1FF), 0x32)
	cpu.bus.Write(defs.Address(0x200), 0x04)

	pc, result, _ := cpu.evalIndirect(cpu.registers.Pc)
	expected := defs.Address(0x432)

	assert.Equal(t, expected, result, "Indirect error")
	assert.Equal(t, cpu.registers.Pc+2, pc)
}

func TestPreIndexedIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.X = 4

	// Write Operand
	cpu.bus.Write(0, 0x10)

	// Write Offset Table
	cpu.bus.Write(0x0014, 0x25)

	pc, result, _ := cpu.evalIndirectX(cpu.registers.Pc)

	expected := defs.Address(0x0025)
	assert.Equal(t, expected, result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestPreIndexedIndirectWithWrapAround(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.X = 21
	// Write Operan
	cpu.bus.Write(0x0000, 250)

	// Write Offset Table
	cpu.bus.Write(0x000F, 0x10)

	pc, result, _ := cpu.evalIndirectX(cpu.registers.Pc)

	expected := defs.Address(0x0010)
	assert.Equal(t, expected, result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestPostIndexedIndirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00

	expected := defs.Address(0xF0)

	cpu.registers.Y = 0x05

	// Opcode Operand
	cpu.bus.Write(0x0000, 0x05)

	// Indexed Table Pointers
	cpu.bus.Write(0x05, 0x20)
	cpu.bus.Write(0x06, 0x00)

	// Offset pointer
	cpu.bus.Write(0x0025, byte(expected))

	pc, result, _ := cpu.evalIndirectY(cpu.registers.Pc)

	assert.Equal(t, expected, result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestPostIndexedIndirectWithWrapAround(t *testing.T) {
	t.Skip()
	cpu := CreateCPUWithBus()

	expected := defs.Address(0xF0)
	cpu.registers.Pc = 0x00
	cpu.registers.Y = 15

	// Opcode Operand
	cpu.bus.Write(0x0000, 0x05)

	// Indexed Table Pointers
	cpu.bus.Write(0x05, 0xFB)
	cpu.bus.Write(0x06, 0x00)

	// Offset pointer
	cpu.bus.Write(0x000A, byte(expected))

	pc, result, _ := cpu.evalIndirectY(cpu.registers.Pc)

	assert.Equal(t, expected, result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestRelativeAddressMode(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.bus.Write(0x09, 0xFF) // OpCode
	cpu.bus.Write(0x10, 0x04) // Operand

	pc, result, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x15), result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestRelativeAddressModeNegative(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x10

	// Write Operand
	cpu.bus.Write(0x10, 0x100-4)

	pc, result, _ := cpu.evalRelative(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x0D), result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}
