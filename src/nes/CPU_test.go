package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushIntoStack(t *testing.T) {
	cpu := CreateCPU()

	cpu.pushStack(0x20)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.ram.read(0x1FF))
}

func TestPushIntoStackWrapsAround(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.Sp = 0x00
	cpu.pushStack(0x20)
	cpu.pushStack(0x21)

	assert.Equal(t, byte(0xff-1), cpu.registers.Sp)
	assert.Equal(t, byte(0x20), cpu.ram.read(0x100))
	assert.Equal(t, byte(0x21), cpu.ram.read(0x1FF))
}
