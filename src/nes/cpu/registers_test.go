package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAFlag(t *testing.T) {
	registers := CreateRegisters()

	registers.UpdateFlag(CarryFlag, 1)
}

func TestNegativeFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.UpdateNegativeFlag(0x80)

	assert.Equal(t, byte(1), registers.NegativeFlag())
}

func TestNegativeFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.UpdateNegativeFlag(0x7F)

	assert.Equal(t, byte(0), registers.NegativeFlag())
}

func TestZeroFlagIsSet(t *testing.T) {
	registers := CreateRegisters()
	registers.UpdateZeroFlag(0x00)

	assert.Equal(t, byte(1), registers.ZeroFlag())
}

func TestZeroFlagIsUnset(t *testing.T) {
	registers := CreateRegisters()
	registers.UpdateZeroFlag(0x01)

	assert.Equal(t, byte(0), registers.ZeroFlag())
}
