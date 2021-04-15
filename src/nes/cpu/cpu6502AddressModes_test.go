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
	expected := defs.Address(0x040)
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

func TestIndexed_indirect(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.X = 4

	// Write Operand
	cpu.bus.Write(0, 0x10)

	// Write Offset Table
	// write low
	cpu.bus.Write(0x0014, 0x25)
	// write high
	cpu.bus.Write(0x0015, 0x10)

	_, result, _ := cpu.evalIndirectX(cpu.registers.Pc)

	expected := defs.Address(0x1025)
	assert.Equal(t, expected, result)
	//assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestIndexed_Indirect_X_wraps(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x0055
	cpu.registers.X = 3
	// Operand
	cpu.bus.Write(0x0055, 0xFE)

	cpu.bus.Write(0x01, 0x10)
	cpu.bus.Write(0x02, 0x05)

	_, result, _ := cpu.evalIndirectX(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x0510), result)
}

func TestIndexed_indirect_with_wrap_around(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x55
	cpu.registers.X = 1
	// Operand
	cpu.bus.Write(0x0055, 0xFE)

	cpu.bus.Write(0xFF, 0x10)
	cpu.bus.Write(0x100, 0x99)
	cpu.bus.Write(0x00, 0x05)

	_, result, _ := cpu.evalIndirectX(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x0510), result)
}

func TestIndirect_indexed(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Pc = 0x00
	cpu.registers.Y = 0x10

	// Operand
	cpu.bus.Write(0x00, 0x86)

	// Indexed Table Pointers
	cpu.bus.Write(0x86, 0x28)
	cpu.bus.Write(0x87, 0x40)
	cpu.bus.Write(0x4038, 0x99)

	pc, result, _ := cpu.evalIndirectY(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x4038), result)
	assert.Equal(t, cpu.registers.Pc+1, pc)
}

func TestIndexed_Indirect_With_Wrap_Around_at_zero_page(t *testing.T) {
	cpu := CreateCPUWithBus()

	cpu.registers.Pc = 0x100
	cpu.registers.Y = 0xFF

	// Opcode Operand
	cpu.bus.Write(0x0100, 0xFF)

	// Indexed Table Pointers
	cpu.bus.Write(0xFF, 0x46)
	cpu.bus.Write(0x00, 0x01)

	pc, result, _ := cpu.evalIndirectY(cpu.registers.Pc)

	assert.Equal(t, defs.Address(0x245), result)
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
