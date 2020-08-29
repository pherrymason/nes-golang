package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegativeFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.updateNegativeFlag(0x80)

	assert.True(t, registers.NegativeFlag)
}

func TestNegativeFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.updateNegativeFlag(0x7F)

	assert.False(t, registers.NegativeFlag)
}

func TestZeroFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.updateZeroFlag(0x00)

	assert.True(t, registers.ZeroFlag)
}

func TestZeroFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.updateZeroFlag(0x01)

	assert.False(t, registers.ZeroFlag)
}
