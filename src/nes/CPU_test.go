package nes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Decode Operation Address Mode

func TestAND(t *testing.T) {
	type dataProvider struct {
		operand              byte
		A                    byte
		expectedA            byte
		expectedNegativeFlag bool
		expectedZeroFlag     bool
	}

	var dataProviders [3]dataProvider
	dataProviders[0] = dataProvider{0b01001100, 0b01000101, 0b01000100, false, false}
	dataProviders[1] = dataProvider{0b10000000, 0b10000000, 0b10000000, true, false}
	dataProviders[2] = dataProvider{0b10000000, 0b01000000, 0b00000000, false, true}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.ram.write(0x100, dp.operand)
		cpu.registers.A = dp.A

		cpu.and(0x100)

		assert.Equal(t, byte(dp.expectedA), cpu.registers.A, fmt.Sprintf("Iteration %d failed, unexpected register A result", i))
		assert.Equal(t, dp.expectedNegativeFlag, cpu.registers.NegativeFlag, fmt.Sprintf("Iteration %d failed, unexpected NegativeFlag result", i))
		assert.Equal(t, dp.expectedZeroFlag, cpu.registers.ZeroFlag, fmt.Sprintf("Iteration %d failed, unexpected ZeroFlag result", i))
	}
}
