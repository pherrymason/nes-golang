package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_loopyRegister_push_first_byte(t *testing.T) {
	register := loopyRegister{}

	register.push(0xFF)

	expected := types.Address(0x3f00)
	assert.Equal(t, expected, register.address())
	assert.Equal(t, uint8(1), register.latch)
}

func Test_loopyRegister_push_second_byte(t *testing.T) {
	register := loopyRegister{
		latch:       1,
		_coarseY:    0b11000,
		_nameTableX: 1,
		_nameTableY: 1,
		_fineY:      0b111,
	}

	register.push(0xFF)

	expected := types.Address(0x3fFF)
	assert.Equal(t, expected, register.address())
	assert.Equal(t, uint8(0), register.latch)
}

func Test_loopyRegister_increment(t *testing.T) {
	register := loopyRegister{
		_coarseX: 0b11111,
	}

	register.increment(0)

	assert.Equal(t, types.Word(0b100000), register.address())
}
