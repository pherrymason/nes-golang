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
		return CPURegisters{0, 0, 0, 0, 0, carryFlag, zeroFlag, false, false, false, 0, negativeFlag}
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
		return CPURegisters{0, 0, 0, 0, 0,
			carryFlag, zeroFlag, false, false, false, 0, negativeFlag}
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
		return CPURegisters{accumulator, 0, 0, 0, 0,
			carryFlag, zeroFlag, false, false, false, overflowFlag, negativeFlag}
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

func TestBEQ(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.ZeroFlag = true

	cpu.beq(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.beq(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0000), cpu.registers.Pc)
}

func TestBIT(t *testing.T) {
	type dataProvider struct {
		accumulator      byte
		operand          byte
		expectedRegister CPURegisters
	}
	expectedRegisters := func(accumulator byte, negativeFlag bool, zeroFlag bool, overflowFlag byte) CPURegisters {
		return CPURegisters{accumulator, 0, 0, 0, 0xFF,
			0, zeroFlag, false, false, false, overflowFlag, negativeFlag}
	}

	dataProviders := [...]dataProvider{
		// All disabled
		{0b00001111, 0b00001111, expectedRegisters(0b00001111, false, false, 0)},
		// Negative true
		{0b10001111, 0b10001111, expectedRegisters(0b10001111, true, false, 0)},
		// Overflow true
		{0b01001111, 0b01001111, expectedRegisters(0b01001111, false, false, 1)},
		// Zero Flag true
		{0, 0, expectedRegisters(0, false, true, 0)},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.accumulator
		cpu.ram.write(0x0001, dp.operand)

		cpu.bit(operation{zeroPage, 0x0001})

		assert.Equal(t, dp.expectedRegister, cpu.registers)
	}
}

func TestBMI(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.NegativeFlag = true

	cpu.bmi(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bmi(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0000), cpu.registers.Pc)
}

func TestBNE(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.ZeroFlag = true

	cpu.bne(operation{relative, 0x0010})

	assert.Equal(t, Address(0), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bne(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)
}

func TestBPL(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.NegativeFlag = true

	cpu.bpl(operation{relative, 0x0010})

	assert.Equal(t, Address(0), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bpl(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)
}

func TestBRK(t *testing.T) {
	programCounter := Address(0x2020)
	expectedPc := Address(0x9999)
	cpu := CreateCPU()
	cpu.registers.Pc = programCounter
	cpu.registers.CarryFlag = 1
	cpu.registers.ZeroFlag = true
	cpu.registers.InterruptDisable = false
	cpu.registers.OverflowFlag = 1
	cpu.registers.NegativeFlag = true
	cpu.ram.write(Address(0xFFFE), 0x99)
	cpu.ram.write(Address(0xFFFF), 0x99)

	cpu.brk(operation{implicit, 0x0000})

	assert.Equal(t, true, cpu.registers.BreakCommand)
	assert.Equal(t, Word(programCounter), cpu.ram.read16(0x1FE))
	assert.Equal(t, byte(0b11110011), cpu.ram.read(0x1FD))

	assert.Equal(t, true, cpu.registers.InterruptDisable)
	assert.Equal(t, true, cpu.registers.BreakCommand)

	assert.Equal(t, expectedPc, cpu.registers.Pc)
}

func TestBVC_overflow_is_not_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.OverflowFlag = 1
	cpu.bvc(operation{relative, 0x5})

	assert.Equal(t, Address(0), cpu.registers.Pc)
}

func TestBVC_overflow_is_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.OverflowFlag = 0
	cpu.bvc(operation{relative, 0x5})

	assert.Equal(t, Address(0x5), cpu.registers.Pc)
}

func TestBVS_overflow_is_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.OverflowFlag = 0
	cpu.bvs(operation{relative, 0x5})

	assert.Equal(t, Address(0x0), cpu.registers.Pc)
}

func TestBVS_overflow_is_set(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.OverflowFlag = 1
	cpu.bvs(operation{relative, 0x5})

	assert.Equal(t, Address(0x5), cpu.registers.Pc)
}

func TestCLC(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.CarryFlag = 1

	cpu.clc(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.CarryFlag)
}

func TestCLD(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.DecimalFlag = true

	cpu.cld(operation{implicit, 0x00})

	assert.False(t, cpu.registers.DecimalFlag)
}

func TestCLI(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.InterruptDisable = true

	cpu.cli(operation{implicit, 0x00})

	assert.False(t, cpu.registers.InterruptDisable)
}

func TestCLV(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.OverflowFlag = 1

	cpu.clv(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.OverflowFlag)
}

func TestCompareOperations(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.X = 0x10
	cpu.registers.A = 0x10
	cpu.registers.Y = 0x10

	type dataProvider struct {
		title            string
		operand          byte
		op               func(operation)
		expectedCarry    byte
		expectedZero     bool
		expectedNegative bool
	}

	dps := [...]dataProvider{
		{"A>M", byte(0x09), cpu.cmp, 1, false, false},
		{"A<M", byte(0x15), cpu.cmp, 0, false, true},
		{"A=M", byte(0x10), cpu.cmp, 1, true, false},
		{"X>M", byte(0x09), cpu.cpx, 1, false, false},
		{"X<M", byte(0x15), cpu.cpx, 0, false, true},
		{"X=M", byte(0x10), cpu.cpx, 1, true, false},
		{"Y>M", byte(0x09), cpu.cpy, 1, false, false},
		{"Y<M", byte(0x15), cpu.cpy, 0, false, true},
		{"Y=M", byte(0x10), cpu.cpy, 1, true, false},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]

		cpu.ram.write(0x00, dp.operand)

		dp.op(operation{zeroPage, Address(0x00)})

		assert.Equal(t, dp.expectedCarry, cpu.registers.CarryFlag, dp.title+": Carry")
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag, dp.title+": Zero")
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag, dp.title+": Negative")
	}
}

func TestDEC(t *testing.T) {
	cpu := CreateCPU()

	cpu.ram.write(0x0000, 0x02)

	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0x01), cpu.ram.read(0))
	assert.Equal(t, false, cpu.registers.NegativeFlag)
	assert.Equal(t, false, cpu.registers.ZeroFlag)

	// Zero result
	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0x00), cpu.ram.read(0))
	assert.Equal(t, false, cpu.registers.NegativeFlag)
	assert.Equal(t, true, cpu.registers.ZeroFlag)

	// Negative result
	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0xFF), cpu.ram.read(0))
	assert.Equal(t, true, cpu.registers.NegativeFlag)
	assert.Equal(t, false, cpu.registers.ZeroFlag)
}

func TestDECXY(t *testing.T) {
	type dataProvider struct {
		title            string
		op               func(operation)
		expectedX        byte
		expectedY        byte
		expectedZero     bool
		expectedNegative bool
	}

	cpu := CreateCPU()
	cpu.registers.X = 2
	cpu.registers.Y = 2

	dps := [...]dataProvider{
		{"X=2", cpu.dex, 1, 2, false, false},
		{"X=1", cpu.dex, 0, 2, true, false},
		{"X=0", cpu.dex, 0xFF, 2, false, true},
		{"Y=2", cpu.dey, 0xFF, 1, false, false},
		{"Y=1", cpu.dey, 0xFF, 0, true, false},
		{"Y=0", cpu.dey, 0xFF, 0xFF, false, true},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		msg := fmt.Sprintf("%s: Unexpected when value is X:%X Y:%X", dp.title, cpu.registers.X, cpu.registers.Y)

		dp.op(operation{implicit, 0})
		assert.Equal(t, dp.expectedX, cpu.registers.X, dp.title)
		assert.Equal(t, dp.expectedY, cpu.registers.Y, dp.title)
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag, msg)
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag, msg)
	}
}

func TestEOR(t *testing.T) {
	type dataProvider struct {
		value            byte
		a                byte
		expectedA        byte
		expectedZero     bool
		expectedNegative bool
	}
	dps := [...]dataProvider{
		{0x00, 0x00, 0x00, true, false},
		{0x01, 0x00, 0x01, false, false},
		{0x80, 0x00, 0x80, false, true},
	}

	cpu := CreateCPU()
	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu.registers.A = dp.a
		cpu.ram.write(0x05, dp.value)
		cpu.eor(operation{immediate, 0x05})

		assert.Equal(t, dp.expectedA, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag)
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag)
	}
}

func TestINC(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     bool
		expectedNegative bool
	}
	dps := [...]dataProvider{
		{0x00, 0x01, false, false},
		{0x7F, 0x80, false, true},
		{0xFF, 0x00, true, false},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.ram.write(0x00, dp.value)

		cpu.inc(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.ram.read(0x00))
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag)
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag)
	}
}

func TestINX(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     bool
		expectedNegative bool
	}
	dps := [...]dataProvider{
		{0x00, 0x01, false, false},
		{0x7F, 0x80, false, true},
		{0xFF, 0x00, true, false},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.registers.X = dp.value

		cpu.inx(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.registers.X)
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag)
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag)
	}
}
func TestINY(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     bool
		expectedNegative bool
	}
	dps := [...]dataProvider{
		{0x00, 0x01, false, false},
		{0x7F, 0x80, false, true},
		{0xFF, 0x00, true, false},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.registers.Y = dp.value

		cpu.iny(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.registers.Y)
		assert.Equal(t, dp.expectedNegative, cpu.registers.NegativeFlag)
		assert.Equal(t, dp.expectedZero, cpu.registers.ZeroFlag)
	}
}

func TestJMP(t *testing.T) {
	cpu := CreateCPU()

	cpu.jmp(operation{absolute, 0x100})

	assert.Equal(t, Address(0x100), cpu.registers.Pc)
}

func TestJSR(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.Pc = 0x0204
	cpu.ram.write(Address(0x201), 0x20) // Opcode
	cpu.ram.write(Address(0x202), 0x55) // LSB
	cpu.ram.write(Address(0x203), 0x05) // MSB
	cpu.jsr(operation{absolute, 0x202})

	assert.Equal(t, Address(0x0555), cpu.registers.Pc)
	assert.Equal(t, byte(0x02), cpu.popStack())
	assert.Equal(t, byte(0x01), cpu.popStack())
}
