package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead16Bugged(t *testing.T) {
	ram := RAM{}

	ram.write(Address(0x1FF), 0x11)
	ram.write(Address(0x200), 0x01)

	result := ram.read16Bugged(Address(0x1FF))
	expected := Word(0x111)

	assert.Equal(t, expected, result)
}
