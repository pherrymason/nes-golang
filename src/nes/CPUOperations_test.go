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
		expectedNegativeFlag byte
		expectedZeroFlag     byte
	}

	var dataProviders [3]dataProvider
	dataProviders[0] = dataProvider{0b01001100, 0b01000101, 0b01000100, 0, 0}
	dataProviders[1] = dataProvider{0b10000000, 0b10000000, 0b10000000, 1, 0}
	dataProviders[2] = dataProvider{0b10000000, 0b01000000, 0b00000000, 0, 1}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.write(0x100, dp.operand)
		cpu.registers.A = dp.A

		cpu.and(0x100)

		assert.Equal(t, byte(dp.expectedA), cpu.registers.A, fmt.Sprintf("Iteration %d failed, unexpected register A result", i))
		assert.Equal(t, dp.expectedNegativeFlag, cpu.registers.negativeFlag(), fmt.Sprintf("Iteration %d failed, unexpected NegativeFlag result", i))
		assert.Equal(t, dp.expectedZeroFlag, cpu.registers.zeroFlag(), fmt.Sprintf("Iteration %d failed, unexpected ZeroFlag result", i))
	}
}

func TestASL_Accumulator(t *testing.T) {
	type dataProvider struct {
		accumulator      byte
		expectedRegister CPURegisters
	}

	expectedRegisters := func(negative byte, zero byte, carry byte) CPURegisters {
		registers := CPURegisters{}
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(carryFlag, carry)

		return registers
	}

	var dataProviders [4]dataProvider
	dataProviders[0] = dataProvider{0b00000001, expectedRegisters(0, 0, 0)}
	dataProviders[1] = dataProvider{0b10000001, expectedRegisters(0, 0, 1)}
	dataProviders[2] = dataProvider{0b10000000, expectedRegisters(0, 1, 1)}
	dataProviders[3] = dataProvider{0b11000000, expectedRegisters(1, 0, 1)}

	for i := 0; i < len(dataProviders); i++ {
		cpu := CreateCPU()
		cpu.registers.A = dataProviders[i].accumulator

		cpu.asl(operation{accumulator, 0x0000})

		assert.Equal(t, dataProviders[i].accumulator<<1, cpu.registers.A, fmt.Sprintf("Iteration %d failed @ expected accumulator", i))
		assert.Equal(t, dataProviders[i].expectedRegister.carryFlag(), cpu.registers.carryFlag(), fmt.Sprintf("Iteration %d failed @ expected CarryFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.zeroFlag(), cpu.registers.zeroFlag(), fmt.Sprintf("Iteration %d failed @ expected ZeroFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.negativeFlag(), cpu.registers.negativeFlag(), fmt.Sprintf("Iteration %d failed @ expected NegativeFlag", i))
	}
}

func TestASL_Memory(t *testing.T) {
	type dataProvider struct {
		operand          byte
		expectedRegister CPURegisters
	}

	expectedRegisters := func(negative byte, zero byte, carry byte) CPURegisters {
		registers := CPURegisters{}
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(carryFlag, carry)

		return registers
	}

	var dataProviders [4]dataProvider
	dataProviders[0] = dataProvider{0b00000001, expectedRegisters(0, 0, 0)}
	dataProviders[1] = dataProvider{0b10000001, expectedRegisters(0, 0, 1)}
	dataProviders[2] = dataProvider{0b10000000, expectedRegisters(0, 1, 1)}
	dataProviders[3] = dataProvider{0b11000000, expectedRegisters(1, 0, 1)}

	for i := 0; i < len(dataProviders); i++ {
		cpu := CreateCPU()
		cpu.registers.A = dataProviders[i].operand

		cpu.write(0x0000, dataProviders[i].operand)
		cpu.asl(operation{zeroPage, 0x0000})

		assert.Equal(t, dataProviders[i].operand<<1, cpu.read(0x0000), fmt.Sprintf("Iteration %d failed @ expected operand", i))
		assert.Equal(t, dataProviders[i].expectedRegister.carryFlag(), cpu.registers.carryFlag(), fmt.Sprintf("Iteration %d failed @ expected CarryFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.zeroFlag(), cpu.registers.zeroFlag(), fmt.Sprintf("Iteration %d failed @ expected ZeroFlag", i))
		assert.Equal(t, dataProviders[i].expectedRegister.negativeFlag(), cpu.registers.negativeFlag(), fmt.Sprintf("Iteration %d failed @ expected NegativeFlag", i))
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
	expectedRegisters := func(accumulator byte, negative byte, zero byte, carry byte, overflow byte) CPURegisters {
		registers := CPURegisters{}
		registers.A = accumulator
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(carryFlag, carry)
		registers.updateFlag(overflowFlag, overflow)

		return registers
	}
	dataProviders := [...]dataProvider{
		{0x05, 0, 0x10, expectedRegisters(0x15, 0, 0, 0, 0)},
		{0x05, 1, 0x10, expectedRegisters(0x16, 0, 0, 0, 0)},
		{0b10000000, 0, 0b10000001, expectedRegisters(0x01, 0, 0, 1, 1)},
		{80, 1, 80, expectedRegisters(161, 1, 0, 0, 1)},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu.registers.A = dp.accumulator
		cpu.registers.updateFlag(carryFlag, dp.carryFlag)
		cpu.write(0x0000, dp.operand)
		cpu.adc(operation{immediate, 0x0000})

		assert.Equal(t, dp.expectedRegister.A, cpu.registers.A, fmt.Sprintf("Iteration %d failed, unexpected A", i))
		assert.Equal(t, dp.expectedRegister.negativeFlag(), cpu.registers.negativeFlag(), fmt.Sprintf("Iteration %d failed, unexpected NegativeFlag", i))
		assert.Equal(t, dp.expectedRegister.zeroFlag(), cpu.registers.zeroFlag(), fmt.Sprintf("Iteration %d failed, unexpected ZeroFlag", i))
		assert.Equal(t, dp.expectedRegister.carryFlag(), cpu.registers.carryFlag(), fmt.Sprintf("Iteration %d failed, unexpected CarryFlag", i))
		assert.Equal(t, dp.expectedRegister.overflowFlag(), cpu.registers.overflowFlag(), fmt.Sprintf("Iteration %d failed, unexpected OverflowFlag", i))
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
		cpu.registers.updateFlag(carryFlag, dp.carryFlag)
		cpu.registers.Pc = dp.pc

		cpu.bcc(operation{relative, Address(dp.branchAddress)})

		assert.Equal(t, Address(dp.expectedPc), cpu.registers.Pc)
	}
}

func TestBCS(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(carryFlag, 1)

	cpu.bcs(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bcs(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0000), cpu.registers.Pc)
}

func TestBEQ(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(zeroFlag, 1)

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
	expectedRegisters := func(accumulator byte, negative byte, zero byte, overflow byte) CPURegisters {
		registers := CreateRegisters()
		registers.A = accumulator
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(overflowFlag, overflow)
		return registers
	}

	dataProviders := [...]dataProvider{
		// All disabled
		{0b00001111, 0b00001111, expectedRegisters(0b00001111, 0, 0, 0)},
		// Negative true
		{0b10001111, 0b10001111, expectedRegisters(0b10001111, 1, 0, 0)},
		// Overflow true
		{0b01001111, 0b01001111, expectedRegisters(0b01001111, 0, 0, 1)},
		// Zero Flag true
		{0, 0, expectedRegisters(0, 0, 1, 0)},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.accumulator
		cpu.write(0x0001, dp.operand)

		cpu.bit(operation{zeroPage, 0x0001})

		assert.Equal(t, dp.expectedRegister, cpu.registers, fmt.Sprintf("Iteration %d failed", i))
	}
}

func TestBMI(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(negativeFlag, 1)

	cpu.bmi(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bmi(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0000), cpu.registers.Pc)
}

func TestBNE(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(zeroFlag, 1)

	cpu.bne(operation{relative, 0x0010})

	assert.Equal(t, Address(0), cpu.registers.Pc)

	cpu.registers.reset()
	cpu.bne(operation{relative, 0x0010})

	assert.Equal(t, Address(0x0010), cpu.registers.Pc)
}

func TestBPL(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(negativeFlag, 1)

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
	cpu.registers.Status = 0b11100011
	cpu.write(Address(0xFFFE), 0x99)
	cpu.write(Address(0xFFFF), 0x99)

	cpu.brk(operation{implicit, 0x0000})

	assert.Equal(t, byte(1), cpu.registers.breakFlag())
	assert.Equal(t, Word(programCounter), cpu.read16(0x1FE))
	// Stored status registers in stack should be...
	assert.Equal(t, byte(0b11110011), cpu.read(0x1FD))

	assert.Equal(t, byte(1), cpu.registers.interruptFlag())
	assert.Equal(t, byte(1), cpu.registers.breakFlag())

	assert.Equal(t, expectedPc, cpu.registers.Pc)
}

func TestBVC_overflow_is_not_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(overflowFlag, 1)
	cpu.bvc(operation{relative, 0x5})

	assert.Equal(t, Address(0), cpu.registers.Pc)
}

func TestBVC_overflow_is_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(overflowFlag, 0)
	cpu.bvc(operation{relative, 0x5})

	assert.Equal(t, Address(0x5), cpu.registers.Pc)
}

func TestBVS_overflow_is_clear(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(overflowFlag, 0)
	cpu.bvs(operation{relative, 0x5})

	assert.Equal(t, Address(0x0), cpu.registers.Pc)
}

func TestBVS_overflow_is_set(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(overflowFlag, 1)
	cpu.bvs(operation{relative, 0x5})

	assert.Equal(t, Address(0x5), cpu.registers.Pc)
}

func TestCLC(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(carryFlag, 1)

	cpu.clc(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.carryFlag())
}

func TestCLD(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(decimalFlag, 1)

	cpu.cld(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.decimalFlag())
}

func TestCLI(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(interruptFlag, 1)

	cpu.cli(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.interruptFlag())
}

func TestCLV(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.updateFlag(overflowFlag, 1)

	cpu.clv(operation{implicit, 0x00})

	assert.Zero(t, cpu.registers.overflowFlag())
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
		expectedZero     byte
		expectedNegative byte
	}

	dps := [...]dataProvider{
		{"A>M", byte(0x09), cpu.cmp, 1, 0, 0},
		{"A<M", byte(0x15), cpu.cmp, 0, 0, 1},
		{"A=M", byte(0x10), cpu.cmp, 1, 1, 0},
		{"X>M", byte(0x09), cpu.cpx, 1, 0, 0},
		{"X<M", byte(0x15), cpu.cpx, 0, 0, 1},
		{"X=M", byte(0x10), cpu.cpx, 1, 1, 0},
		{"Y>M", byte(0x09), cpu.cpy, 1, 0, 0},
		{"Y<M", byte(0x15), cpu.cpy, 0, 0, 1},
		{"Y=M", byte(0x10), cpu.cpy, 1, 1, 0},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]

		cpu.write(0x00, dp.operand)

		dp.op(operation{zeroPage, Address(0x00)})

		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag(), dp.title+": Carry")
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), dp.title+": Zero")
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag(), dp.title+": Negative")
	}
}

func TestDEC(t *testing.T) {
	cpu := CreateCPU()

	cpu.write(0x0000, 0x02)

	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0x01), cpu.read(0))
	assert.Equal(t, byte(0), cpu.registers.zeroFlag())
	assert.Equal(t, byte(0), cpu.registers.negativeFlag())

	// Zero result
	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0x00), cpu.read(0))
	assert.Equal(t, byte(0), cpu.registers.negativeFlag())
	assert.Equal(t, byte(1), cpu.registers.zeroFlag())

	// Negative result
	cpu.dec(operation{zeroPage, Address(0x0000)})

	assert.Equal(t, byte(0xFF), cpu.read(0))
	assert.Equal(t, byte(1), cpu.registers.negativeFlag())
	assert.Equal(t, byte(0), cpu.registers.zeroFlag())
}

func TestDECXY(t *testing.T) {
	type dataProvider struct {
		title            string
		op               func(operation)
		expectedX        byte
		expectedY        byte
		expectedZero     byte
		expectedNegative byte
	}

	cpu := CreateCPU()
	cpu.registers.X = 2
	cpu.registers.Y = 2

	dps := [...]dataProvider{
		{"X=2", cpu.dex, 1, 2, 0, 0},
		{"X=1", cpu.dex, 0, 2, 1, 0},
		{"X=0", cpu.dex, 0xFF, 2, 0, 1},
		{"Y=2", cpu.dey, 0xFF, 1, 0, 0},
		{"Y=1", cpu.dey, 0xFF, 0, 1, 0},
		{"Y=0", cpu.dey, 0xFF, 0xFF, 0, 1},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		msg := fmt.Sprintf("%s: Unexpected when value is X:%X Y:%X", dp.title, cpu.registers.X, cpu.registers.Y)

		dp.op(operation{implicit, 0})
		assert.Equal(t, dp.expectedX, cpu.registers.X, dp.title)
		assert.Equal(t, dp.expectedY, cpu.registers.Y, dp.title)
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag(), msg)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), msg)
	}
}

func TestEOR(t *testing.T) {
	type dataProvider struct {
		value            byte
		a                byte
		expectedA        byte
		expectedZero     byte
		expectedNegative byte
	}
	dps := [...]dataProvider{
		{0x00, 0x00, 0x00, 1, 0},
		{0x01, 0x00, 0x01, 0, 0},
		{0x80, 0x00, 0x80, 0, 1},
	}

	cpu := CreateCPU()
	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu.registers.A = dp.a
		cpu.write(0x05, dp.value)
		cpu.eor(operation{immediate, 0x05})

		assert.Equal(t, dp.expectedA, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestINC(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}
	dps := [...]dataProvider{
		{0x00, 0x01, 0, 0},
		{0x7F, 0x80, 0, 1},
		{0xFF, 0x00, 1, 0},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.write(0x00, dp.value)

		cpu.inc(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.read(0x00))
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
	}
}

func TestINX(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}
	dps := [...]dataProvider{
		{0x00, 0x01, 0, 0},
		{0x7F, 0x80, 0, 1},
		{0xFF, 0x00, 1, 0},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.registers.X = dp.value

		cpu.inx(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.registers.X)
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
	}
}
func TestINY(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}
	dps := [...]dataProvider{
		{0x00, 0x01, 0, 0},
		{0x7F, 0x80, 0, 1},
		{0xFF, 0x00, 1, 0},
	}

	for i := 0; i < len(dps); i++ {
		dp := dps[i]
		cpu := CreateCPU()
		cpu.registers.Y = dp.value

		cpu.iny(operation{zeroPage, 0x00})
		assert.Equal(t, dp.expectedValue, cpu.registers.Y)
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
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
	cpu.write(Address(0x201), 0x20) // Opcode
	cpu.write(Address(0x202), 0x55) // LSB
	cpu.write(Address(0x203), 0x05) // MSB
	cpu.jsr(operation{absolute, 0x202})

	assert.Equal(t, Address(0x0555), cpu.registers.Pc)
	assert.Equal(t, byte(0x02), cpu.popStack())
	assert.Equal(t, byte(0x01), cpu.popStack())
}

func TestLDA(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedZero     byte
		expectedNegative byte
	}
	dataProviders := [...]dataProvider{
		{0x20, 0, 0},
		{0x00, 1, 0},
		{0x80, 0, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.write(Address(0x00), dp.value)

		cpu.lda(operation{immediate, 0x00})

		assert.Equal(t, dp.value, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestLDX(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedZero     byte
		expectedNegative byte
	}
	dataProviders := [...]dataProvider{
		{0x20, 0, 0},
		{0x00, 1, 0},
		{0x80, 0, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.write(Address(0x00), dp.value)

		cpu.ldx(operation{immediate, 0x00})

		assert.Equal(t, dp.value, cpu.registers.X)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestLDY(t *testing.T) {
	type dataProvider struct {
		value            byte
		expectedZero     byte
		expectedNegative byte
	}
	dataProviders := [...]dataProvider{
		{0x20, 0, 0},
		{0x00, 1, 0},
		{0x80, 0, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.write(Address(0x00), dp.value)

		cpu.ldy(operation{immediate, 0x00})

		assert.Equal(t, dp.value, cpu.registers.Y)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestLSR(t *testing.T) {
	type dataProvider struct {
		addressingMode AddressMode
		value          byte
		expectedResult byte
		expectedZero   byte
		expectedCarry  byte
	}
	dataProviders := [...]dataProvider{
		{accumulator, 0b00000010, 0b00000001, 0, 0},
		{accumulator, 0b00000011, 0b00000001, 0, 1},
		{accumulator, 0b00000001, 0b00000000, 1, 1},
		{zeroPage, 0b00000010, 0b00000001, 0, 0},
		{zeroPage, 0b00000011, 0b00000001, 0, 1},
		{zeroPage, 0b00000001, 0b00000000, 1, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.value
		cpu.write(Address(0x00), dp.value)

		cpu.lsr(operation{dp.addressingMode, 0x00})

		if dp.addressingMode == accumulator {
			assert.Equal(t, dp.expectedResult, cpu.registers.A)
		} else {
			assert.Equal(t, dp.expectedResult, cpu.read(0x00))
		}
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), fmt.Sprintf("Iteration[%d] unexpected ZeroFlag", i))
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag())
		assert.Zero(t, cpu.registers.negativeFlag())
	}
}

func TestORA(t *testing.T) {
	type dataProvider struct {
		a                byte
		value            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
	}
	dataProviders := [...]dataProvider{
		{0x00, 0x00, 0x00, 1, 0},
		{0x80, 0x00, 0x80, 0, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.a
		cpu.write(0x00, dp.value)

		cpu.ora(operation{immediate, 0x00})

		assert.Equal(t, dp.expectedResult, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestPHA(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.A = 0x30
	cpu.pha(operation{implicit, 0x00})

	assert.Equal(t, byte(0x30), cpu.popStack())
}

func TestPHP(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.Status = 0xFF

	cpu.php(operation{implicit, 0x00})

	assert.Equal(t, byte(0xFF), cpu.popStack())
}

func TestPLA(t *testing.T) {
	type dataProvider struct {
		pulledValue      byte
		expectedNegative byte
		expectedZero     byte
	}

	dataProviders := [...]dataProvider{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.pushStack(dp.pulledValue)

		cpu.pla(operation{implicit, 0x00})

		assert.Equal(t, dp.pulledValue, cpu.registers.A)
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, byte(0xFF), cpu.registers.Sp)
	}
}

func TestPLP(t *testing.T) {
	cpu := CreateCPU()
	cpu.pushStack(0xFF)

	cpu.plp(operation{implicit, 0x00})

	assert.Equal(t, byte(0xFF), cpu.registers.Status)
	assert.Equal(t, byte(0xFF), cpu.registers.Sp)
}

func TestROL(t *testing.T) {
	type dataProvider struct {
		addressingMode   AddressMode
		value            byte
		carry            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
	}
	dataProviders := [...]dataProvider{
		{accumulator, 0b00000000, 0, 0, 1, 0, 0},
		{accumulator, 0b00000000, 1, 1, 0, 0, 0},
		{accumulator, 0b00000001, 0, 0b10, 0, 0, 0},
		{accumulator, 0b10000000, 0, 0, 1, 0, 1},
		{accumulator, 0b01000000, 0, 0x80, 0, 1, 0},

		{zeroPage, 0b00000000, 0, 0, 1, 0, 0},
		{zeroPage, 0b00000000, 1, 1, 0, 0, 0},
		{zeroPage, 0b00000001, 0, 0b10, 0, 0, 0},
		{zeroPage, 0b10000000, 0, 0, 1, 0, 1},
		{zeroPage, 0b01000000, 0, 0x80, 0, 1, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.value
		cpu.registers.updateFlag(carryFlag, dp.carry)
		cpu.write(Address(0x00), dp.value)

		cpu.rol(operation{dp.addressingMode, 0x00})

		if dp.addressingMode == accumulator {
			assert.Equal(t, dp.expectedResult, cpu.registers.A)
		} else {
			assert.Equal(t, dp.expectedResult, cpu.read(0x00))
		}
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), fmt.Sprintf("Iteration[%d] unexpected ZeroFlag", i))
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag(), fmt.Sprintf("Iteration[%d] unexpected Negative", i))
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag())
	}
}

func TestROR(t *testing.T) {
	type dataProvider struct {
		addressingMode   AddressMode
		value            byte
		carry            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
	}
	dataProviders := [...]dataProvider{
		{accumulator, 0b00000000, 0, 0, 1, 0, 0},
		{accumulator, 0b00000001, 0, 0, 1, 0, 1},
		{accumulator, 0b00000000, 1, 0x80, 0, 1, 0},
		{accumulator, 0b10000000, 0, 0x40, 0, 0, 0},
		{accumulator, 0b10000001, 1, 0xC0, 0, 1, 1},

		{zeroPage, 0b00000000, 0, 0, 1, 0, 0},
		{zeroPage, 0b00000001, 0, 0, 1, 0, 1},
		{zeroPage, 0b00000000, 1, 0x80, 0, 1, 0},
		{zeroPage, 0b10000000, 0, 0x40, 0, 0, 0},
		{zeroPage, 0b10000001, 1, 0xC0, 0, 1, 1},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.value
		cpu.registers.updateFlag(carryFlag, dp.carry)
		cpu.write(Address(0x00), dp.value)

		cpu.ror(operation{dp.addressingMode, 0x00})

		if dp.addressingMode == accumulator {
			assert.Equal(t, dp.expectedResult, cpu.registers.A)
		} else {
			assert.Equal(t, dp.expectedResult, cpu.read(0x00))
		}
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), fmt.Sprintf("Iteration[%d] unexpected ZeroFlag", i))
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag(), fmt.Sprintf("Iteration[%d] unexpected Negative", i))
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag())
	}
}

func TestRTI(t *testing.T) {
	cpu := CreateCPU()
	// Push an Address into Stack
	expectedProgramCounter := Address(0x532)
	cpu.pushStack(byte(expectedProgramCounter & 0xFF))
	cpu.pushStack(byte(expectedProgramCounter >> 8))
	// Push a StatusRegister into stack
	cpu.pushStack(0xFF)

	cpu.rti(operation{implicit, 0xFF})

	assert.Equal(t, expectedProgramCounter, cpu.registers.Pc)
	assert.Equal(t, byte(0xFF), cpu.registers.Status)
}

func TestRTS(t *testing.T) {
	cpu := CreateCPU()
	// Push an Address into Stack
	expectedProgramCounter := Address(0x532)
	cpu.pushStack(byte(expectedProgramCounter & 0xFF))
	cpu.pushStack(byte(expectedProgramCounter >> 8))

	cpu.rts(operation{implicit, 0x00})

	assert.Equal(t, expectedProgramCounter, cpu.registers.Pc)
}

func TestSBC(t *testing.T) {
	type dataProvider struct {
		a     byte
		value byte
		carry byte

		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
		expectedOverflow byte
	}
	// Fixtures taken from http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html (section SBC)
	// TODO: Improve these fixtures by adding more (if really needed)
	dataProviders := [...]dataProvider{
		{0x01, 1, 1, 0, 1, 0, 1, 0},
		{0x50, 0xF0, 1, 0x60, 0, 0, 0, 0},
		{0x50, 0xB0, 1, 0xA0, 0, 1, 0, 1},
		{0x50, 0x70, 1, 0xE0, 0, 1, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.A = dp.a
		cpu.registers.updateFlag(carryFlag, dp.carry)
		cpu.write(0x00, dp.value)

		cpu.sbc(operation{immediate, 0x00})

		assert.Equal(t, dp.expectedResult, cpu.registers.A, "Invalid subtraction result")
		assert.Equal(t, dp.expectedCarry, cpu.registers.carryFlag(), "Invalid CarryFlag")
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag(), "Invalid zeroflag")
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag(), "Invalid negative Flag")
		assert.Equal(t, dp.expectedOverflow, cpu.registers.overflowFlag(), "Invalid Overflow Flag")
	}
}

func TestSEC(t *testing.T) {
	cpu := CreateCPU()

	cpu.sec(operation{implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.carryFlag())
}

func TestSED(t *testing.T) {
	cpu := CreateCPU()

	cpu.sed(operation{implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.decimalFlag())
}

func TestSEI(t *testing.T) {
	cpu := CreateCPU()

	cpu.sei(operation{implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.interruptFlag())
}

func TestSTA(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.A = 0xFF

	cpu.sta(operation{implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.read(0x522))
}

func TestSTX(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.X = 0xFF

	cpu.stx(operation{implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.read(0x522))
}

func TestSTY(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.Y = 0xFF

	cpu.sty(operation{implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.read(0x522))
}

func TestTAX_TAY(t *testing.T) {
	type dataProvider struct {
		name             string
		op               func(info operation)
		a                byte
		expectedNegative byte
		expectedZero     byte
	}

	cpu := CreateCPU()
	dataProviders := [...]dataProvider{
		{"tax", cpu.tax, 0x00, 0, 1},
		{"tax", cpu.tax, 0x80, 1, 0},
		{"tax", cpu.tax, 0x20, 0, 0},
		{"tay", cpu.tay, 0x00, 0, 1},
		{"tay", cpu.tay, 0x80, 1, 0},
		{"tay", cpu.tay, 0x20, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu.registers.reset()
		cpu.registers.A = dp.a

		dp.op(operation{implicit, 0x00})

		if dp.name == "tax" {
			assert.Equal(t, cpu.registers.A, cpu.registers.X)
		} else {
			assert.Equal(t, cpu.registers.A, cpu.registers.Y)
		}
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
	}
}

func TestTSX(t *testing.T) {
	type dataProvider struct {
		sp               byte
		expectedNegative byte
		expectedZero     byte
	}

	dataProviders := [...]dataProvider{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.pushStack(dp.sp)

		cpu.tsx(operation{implicit, 0x00})

		assert.Equal(t, dp.sp, cpu.registers.X)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestTXA(t *testing.T) {
	type dataProvider struct {
		x                byte
		expectedNegative byte
		expectedZero     byte
	}

	dataProviders := [...]dataProvider{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.X = dp.x

		cpu.txa(operation{implicit, 0x00})

		assert.Equal(t, dp.x, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}

func TestTXS(t *testing.T) {
	cpu := CreateCPU()
	cpu.registers.X = 0xFF

	cpu.txs(operation{implicit, 0x00})

	assert.Equal(t, byte(0xFF), cpu.registers.Sp)
}

func TestTYA(t *testing.T) {
	type dataProvider struct {
		y                byte
		expectedNegative byte
		expectedZero     byte
	}

	dataProviders := [...]dataProvider{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for i := 0; i < len(dataProviders); i++ {
		dp := dataProviders[i]
		cpu := CreateCPU()
		cpu.registers.Y = dp.y

		cpu.tya(operation{implicit, 0x00})

		assert.Equal(t, cpu.registers.Y, cpu.registers.A)
		assert.Equal(t, dp.expectedZero, cpu.registers.zeroFlag())
		assert.Equal(t, dp.expectedNegative, cpu.registers.negativeFlag())
	}
}
