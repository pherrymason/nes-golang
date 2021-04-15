package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPushIntoStack(t *testing.T) {
	cpu := CreateCPUWithBus()

	cpu.pushStack(0x20)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.Read(0x1FF))
}

func TestPushIntoStackWrapsAround(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Sp = 0x00
	cpu.pushStack(0x20)
	cpu.pushStack(0x21)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.Read(0x100))
	assert.Equal(t, byte(0x21), cpu.Read(0x1FF))
}

func TestGetStatusRegister(t *testing.T) {
	cpu := CreateCPUWithBus()
	cpu.registers.Status = 0

	cpu.registers.setFlag(carryFlag)
	assert.Equal(t, byte(0x1), cpu.registers.Status)

	cpu.registers.setFlag(zeroFlag)
	assert.Equal(t, byte(0x3), cpu.registers.Status)

	cpu.registers.setFlag(interruptFlag)
	assert.Equal(t, byte(0x7), cpu.registers.Status)

	cpu.registers.setFlag(decimalFlag)
	assert.Equal(t, byte(0xF), cpu.registers.Status)

	//cpu.registers.setFlag(breakCommandFlag)
	//assert.Equal(t, byte(0x3F), cpu.registers.Status)

	cpu.registers.setFlag(overflowFlag)
	assert.Equal(t, byte(0x4F), cpu.registers.Status)

	cpu.registers.setFlag(negativeFlag)
	assert.Equal(t, byte(0xCF), cpu.registers.Status)
}
