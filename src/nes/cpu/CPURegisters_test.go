package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAFlag(t *testing.T) {
	registers := CreateRegisters()

	registers.updateFlag(carryFlag, 1)
}

func TestNegativeFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.updateNegativeFlag(0x80)

	assert.Equal(t, byte(1), registers.negativeFlag())
}

func TestNegativeFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.updateNegativeFlag(0x7F)

	assert.Equal(t, byte(0), registers.negativeFlag())
}

func TestZeroFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.updateZeroFlag(0x00)

	assert.Equal(t, byte(1), registers.zeroFlag())
}

func TestZeroFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.updateZeroFlag(0x01)

	assert.Equal(t, byte(0), registers.zeroFlag())
}
