package nes

import (
	nesCpu "github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPushIntoStack(t *testing.T) {
	cpu := CreateCPUWithGamePak()

	cpu.pushStack(0x20)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.memory.Read(0x1FF))
}

func TestPushIntoStackWrapsAround(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Sp = 0x00
	cpu.pushStack(0x20)
	cpu.pushStack(0x21)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.memory.Read(0x100))
	assert.Equal(t, byte(0x21), cpu.memory.Read(0x1FF))
}

func TestGetStatusRegister(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Status = 0

	cpu.registers.SetFlag(nesCpu.CarryFlag)
	assert.Equal(t, byte(0x1), cpu.registers.Status)

	cpu.registers.SetFlag(nesCpu.ZeroFlag)
	assert.Equal(t, byte(0x3), cpu.registers.Status)

	cpu.registers.SetFlag(nesCpu.InterruptFlag)
	assert.Equal(t, byte(0x7), cpu.registers.Status)

	cpu.registers.SetFlag(nesCpu.DecimalFlag)
	assert.Equal(t, byte(0xF), cpu.registers.Status)

	//cpu.registers.SetFlag(nesCpu.breakCommandFlag)
	//assert.Equal(t, byte(0x3F), cpu.registers.Status)

	cpu.registers.SetFlag(nesCpu.OverflowFlag)
	assert.Equal(t, byte(0x4F), cpu.registers.Status)

	cpu.registers.SetFlag(nesCpu.NegativeFlag)
	assert.Equal(t, byte(0xCF), cpu.registers.Status)
}
