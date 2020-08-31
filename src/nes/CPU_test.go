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

func TestASL_Accumulator(t *testing.T) {
	type dataProvider struct {
		accumulator      byte
		expectedRegister CPURegisters
	}

	expectedRegisters := func(negativeFlag bool, zeroFlag bool, carryFlag byte) CPURegisters {
		return CPURegisters{0, 0, 0, 0, 0, negativeFlag, zeroFlag, carryFlag, 0}
	}

	var dataProviders [4]dataProvider
	dataProviders[0] = dataProvider{0b00000001, expectedRegisters(false, false, 0)}
	dataProviders[1] = dataProvider{0b10000001, expectedRegisters(false, false, 1)}
	dataProviders[2] = dataProvider{0b10000000, expectedRegisters(false, true, 1)}
	dataProviders[3] = dataProvider{0b11000000, expectedRegisters(true, false, 1)}

	for i := 0; i < len(dataProviders); i++ {
		cpu := CreateCPU()
		cpu.registers.A = dataProviders[i].accumulator

		cpu.asl(operation{accumulator, 0x0000})

		assert.Equal(t, dataProviders[i].accumulator<<1, cpu.registers.A, fmt.Sprintf("Iteration %d failed @ expected accumulator", i))
		assert.Equal(t, dataProviders[i].expectedRegister.CarryFlag, cpu.registers.CarryFlag, fmt.Sprintf("Iteration %d failed @ expected CarryFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.ZeroFlag, cpu.registers.ZeroFlag, fmt.Sprintf("Iteration %d failed @ expected ZeroFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.NegativeFlag, cpu.registers.NegativeFlag, fmt.Sprintf("Iteration %d failed @ expected NegativeFlag", i))
	}
}

func TestASL_Memory(t *testing.T) {
	type dataProvider struct {
		operand          byte
		expectedRegister CPURegisters
	}

	expectedRegisters := func(negativeFlag bool, zeroFlag bool, carryFlag byte) CPURegisters {
		return CPURegisters{0, 0, 0, 0, 0, negativeFlag, zeroFlag, carryFlag, 0}
	}

	var dataProviders [4]dataProvider
	dataProviders[0] = dataProvider{0b00000001, expectedRegisters(false, false, 0)}
	dataProviders[1] = dataProvider{0b10000001, expectedRegisters(false, false, 1)}
	dataProviders[2] = dataProvider{0b10000000, expectedRegisters(false, true, 1)}
	dataProviders[3] = dataProvider{0b11000000, expectedRegisters(true, false, 1)}

	for i := 0; i < len(dataProviders); i++ {
		cpu := CreateCPU()
		cpu.registers.A = dataProviders[i].operand

		cpu.ram.write(0x0000, dataProviders[i].operand)
		cpu.asl(operation{zeroPage, 0x0000})

		assert.Equal(t, dataProviders[i].operand<<1, cpu.ram.read(0x0000), fmt.Sprintf("Iteration %d failed @ expected operand", i))
		assert.Equal(t, dataProviders[i].expectedRegister.CarryFlag, cpu.registers.CarryFlag, fmt.Sprintf("Iteration %d failed @ expected CarryFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.ZeroFlag, cpu.registers.ZeroFlag, fmt.Sprintf("Iteration %d failed @ expected ZeroFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.NegativeFlag, cpu.registers.NegativeFlag, fmt.Sprintf("Iteration %d failed @ expected NegativeFlag", i))
	}
}

func TestADC(t *testing.T) {

	cpu := CreateCPU()

	type dataProvider struct {
		accumulator      byte
		carryFlag        byte
		operand          byte
		expectedRegister CPURegisters
	}
	expectedRegisters := func(accumulator byte, negativeFlag bool, zeroFlag bool, carryFlag byte, overflowFlag byte) CPURegisters {
		return CPURegisters{accumulator, 0, 0, 0, 0, negativeFlag, zeroFlag, carryFlag, overflowFlag}
	}
	dataProviders := [...]dataProvider{
		{0x05, 0, 0x10, expectedRegisters(0x15, false, false, 0, 0)},
		{0x05, 1, 0x10, expectedRegisters(0x16, false, false, 0, 0)},
		{0b10000000, 0, 0b10000001, expectedRegisters(0x01, false, false, 1, 1)},
		{80, 1, 80, expectedRegisters(161, true, false, 0, 1)},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu.registers.A = dp.accumulator
		cpu.registers.CarryFlag = dp.carryFlag
		cpu.ram.write(0x0000, dp.operand)
		cpu.adc(operation{immediate, 0x0000})

		assert.Equal(t, dp.expectedRegister.A, cpu.registers.A, fmt.Sprintf("Iteration %d failed, unexpected A", i))
		assert.Equal(t, dp.expectedRegister.NegativeFlag, cpu.registers.NegativeFlag, fmt.Sprintf("Iteration %d failed, unexpected NegativeFlag", i))
		assert.Equal(t, dp.expectedRegister.ZeroFlag, cpu.registers.ZeroFlag, fmt.Sprintf("Iteration %d failed, unexpected ZeroFlag", i))
		assert.Equal(t, dp.expectedRegister.CarryFlag, cpu.registers.CarryFlag, fmt.Sprintf("Iteration %d failed, unexpected CarryFlag", i))
		assert.Equal(t, dp.expectedRegister.OverflowFlag, cpu.registers.OverflowFlag, fmt.Sprintf("Iteration %d failed, unexpected OverflowFlag", i))
	}
}

func TestBCC(t *testing.T) {
	type dataProvider struct {
		carryFlag     byte
		pc            Address
		branchAddress Address
		expectedPc    Address
	}

	dataProviders := [...]dataProvider{
		{0, 0x02, 0x0004, 0x04},
		{1, 0x02, 0x00FF, 0x02},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]

		cpu := CreateCPU()
		cpu.registers.CarryFlag = dp.carryFlag
		cpu.registers.Pc = dp.pc

		cpu.bcc(operation{relative, Address(dp.branchAddress)})

		assert.Equal(t, Address(dp.expectedPc), cpu.registers.Pc)
	}
}

func TestBCS(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.CarryFlag = 1

	cpu.bcs(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bcs(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0000), cpu.registers.Pc)
}
