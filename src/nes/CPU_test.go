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

	cpu.registers.setFlag(carryFlag)
	assert.Equal(t, byte(0x21), cpu.registers.Status)

	cpu.registers.setFlag(zeroFlag)
	assert.Equal(t, byte(0x23), cpu.registers.Status)

	cpu.registers.setFlag(interruptFlag)
	assert.Equal(t, byte(0x27), cpu.registers.Status)

	cpu.registers.setFlag(decimalFlag)
	assert.Equal(t, byte(0x2F), cpu.registers.Status)

	cpu.registers.setFlag(breakCommandFlag)
	assert.Equal(t, byte(0x3F), cpu.registers.Status)

	cpu.registers.setFlag(overflowFlag)
	assert.Equal(t, byte(0x7F), cpu.registers.Status)

	cpu.registers.setFlag(negativeFlag)
	assert.Equal(t, byte(0xFF), cpu.registers.Status)
}
