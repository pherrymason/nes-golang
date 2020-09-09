package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushIntoStack(t *testing.T) {
	cpu := CreateCPU()

	cpu.pushStack(0x20)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.read(0x1FF))
}

func TestPushIntoStackWrapsAround(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.Sp = 0x00
	cpu.pushStack(0x20)
	cpu.pushStack(0x21)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.read(0x100))
	assert.Equal(t, byte(0x21), cpu.read(0x1FF))
}

func TestGetStatusRegister(t *testing.T) {
	cpu := CreateCPU()

	cpu.registers.CarryFlag = 1
	assert.Equal(t, byte(0x21), cpu.registers.statusRegister())

	cpu.registers.ZeroFlag = true
	assert.Equal(t, byte(0x23), cpu.registers.statusRegister())

	cpu.registers.InterruptDisable = true
	assert.Equal(t, byte(0x27), cpu.registers.statusRegister())

	cpu.registers.DecimalFlag = true
	assert.Equal(t, byte(0x2F), cpu.registers.statusRegister())

	cpu.registers.BreakCommand = true
	assert.Equal(t, byte(0x3F), cpu.registers.statusRegister())

	cpu.registers.OverflowFlag = 1
	assert.Equal(t, byte(0x7F), cpu.registers.statusRegister())

	cpu.registers.NegativeFlag = true
	assert.Equal(t, byte(0xFF), cpu.registers.statusRegister())
}
