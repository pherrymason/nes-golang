package nes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleMapper struct {
	memory [0x10000]byte
}

func (s SimpleMapper) prgBanks() byte {
	return 1
}

func (s SimpleMapper) chrBanks() byte {
	return 0
}

func (s SimpleMapper) Read(address Address) byte {
	return s.memory[address]
}

func (s *SimpleMapper) Write(address Address, value byte) {
	s.memory[address] = value
}

func CreateCPUMemoryWithSimpleMapper() *CPUMemory {
	mapper := &SimpleMapper{}

	return &CPUMemory{gamePak: nil, mapper: mapper, ppu: nil}
}

// CreateCPUWithGamePak creates a Cpu6502 with a Bus, Useful for tests
func CreateCPUWithGamePak() *Cpu6502 {
	cpu := CreateCPU(
		CreateCPUMemoryWithSimpleMapper(),
		Cpu6502DebugOptions{false, ""},
	)
	return cpu
}

// Decode Operation Address Mode

func TestAND(t *testing.T) {
	cases := []struct {
		name                 string
		operand              byte
		A                    byte
		expectedA            byte
		expectedNegativeFlag byte
		expectedZeroFlag     byte
	}{
		{"result is > 0", 0b01001100, 0b01000101, 0b01000100, 0, 0},
		{"result is < 0", 0b10000000, 0b10000000, 0b10000000, 1, 0},
		{"result is 0", 0b10000000, 0b01000000, 0b00000000, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(0x100, tt.operand)
			cpu.registers.A = tt.A

			extraCycle := cpu.and(OperationMethodArgument{Immediate, 0x100})

			assert.Equal(t, tt.expectedA, cpu.registers.A, fmt.Sprintf("unexpected register A result"))
			assert.Equal(t, tt.expectedNegativeFlag, cpu.registers.negativeFlag(), fmt.Sprintf("unexpected NegativeFlag result"))
			assert.Equal(t, tt.expectedZeroFlag, cpu.registers.zeroFlag(), fmt.Sprintf("unexpected ZeroFlag result"))
			assert.True(t, extraCycle)
		})
	}
}

func TestASL_Accumulator(t *testing.T) {
	expectedRegisters := func(negative byte, zero byte, carry byte) Cpu6502Registers {
		registers := Cpu6502Registers{}
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(carryFlag, carry)

		return registers
	}

	cases := []struct {
		name             string
		input            byte
		addressMode      AddressMode
		expectedRegister Cpu6502Registers
	}{
		// Acumulator
		{"Shift left Acumulator, result is > 0 without carry", 0b00000001, Implicit, expectedRegisters(0, 0, 0)},
		{"Shift left Acumulator,result is > 0 with carry", 0b10000001, Implicit, expectedRegisters(0, 0, 1)},
		{"Shift left Acumulator,result is 0 with carry", 0b10000000, Implicit, expectedRegisters(0, 1, 1)},
		{"Shift left Acumulator,result is < 0 with carry", 0b11000000, Implicit, expectedRegisters(1, 0, 1)},
		{"Shift left Acumulator,result is < 0 without carry", 0b01000000, Implicit, expectedRegisters(1, 0, 0)},
		// Over memory
		{"Shift left $, result > 0 without carry", 0b00000001, ZeroPage, expectedRegisters(0, 0, 0)},
		{"Shift left $, result is > 0 with carry", 0b10000001, ZeroPage, expectedRegisters(0, 0, 1)},
		{"Shift left $, result is 0 with carry", 0b10000000, ZeroPage, expectedRegisters(0, 1, 1)},
		{"Shift left $, result is < 0, with carry", 0b11000000, ZeroPage, expectedRegisters(1, 0, 1)},
		{"Shift left $, result is < 0, without carry", 0b01000000, ZeroPage, expectedRegisters(1, 0, 0)},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.input
			cpu.memory.Write(0x0000, tt.input)

			cpu.asl(OperationMethodArgument{tt.addressMode, 0x0000})

			if tt.addressMode == Implicit {
				assert.Equal(t, tt.input<<1, cpu.registers.A, "unexpected Accumulator")
			} else {
				assert.Equal(t, tt.input<<1, cpu.memory.Read(0x0000), "unexpected result")
			}
			assert.Equal(t, tt.expectedRegister.carryFlag(), cpu.registers.carryFlag(), "unexpected CarryFlag")
			assert.Equal(t, tt.expectedRegister.zeroFlag(), cpu.registers.zeroFlag(), "unexpected ZeroFlag")
			assert.Equal(t, tt.expectedRegister.negativeFlag(), cpu.registers.negativeFlag(), "unexpected NegativeFlag")
		})
	}
}

func TestADC(t *testing.T) {

	cpu := CreateCPUWithGamePak()

	expectedRegisters := func(accumulator byte, negative byte, zero byte, carry byte, overflow byte) Cpu6502Registers {
		registers := Cpu6502Registers{}
		registers.A = accumulator
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(carryFlag, carry)
		registers.updateFlag(overflowFlag, overflow)

		return registers
	}

	cases := []struct {
		name             string
		accumulator      byte
		carryFlag        byte
		operand          byte
		expectedRegister Cpu6502Registers
	}{
		{"result is > 0 w/o C and O", 0x05, 0, 0x10, expectedRegisters(0x15, 0, 0, 0, 0)},
		{"result is > 0 w/o C and O", 0x05, 1, 0x10, expectedRegisters(0x16, 0, 0, 0, 0)},
		{"result is > 0 with C and O", 0b10000000, 0, 0b10000001, expectedRegisters(0x01, 0, 0, 1, 1)},
		{"result is < 0 w/o C and with O", 80, 1, 80, expectedRegisters(161, 1, 0, 0, 1)},
		//TODO {"result is < 0 with C and w/o O", ?, ?, ?, expectedRegisters()},
		//TODO {"result is 0 w/o C and O", ?, ?, ?, expectedRegisters()},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu.registers.A = tt.accumulator
			cpu.registers.updateFlag(carryFlag, tt.carryFlag)
			cpu.memory.Write(0x0000, tt.operand)

			extraCycle := cpu.adc(OperationMethodArgument{Immediate, 0x0000})

			assert.Equal(t, tt.expectedRegister.A, cpu.registers.A, "unexpected A")
			assert.Equal(t, tt.expectedRegister.negativeFlag(), cpu.registers.negativeFlag(), "unexpected NegativeFlag")
			assert.Equal(t, tt.expectedRegister.zeroFlag(), cpu.registers.zeroFlag(), "unexpected ZeroFlag")
			assert.Equal(t, tt.expectedRegister.carryFlag(), cpu.registers.carryFlag(), "unexpected CarryFlag")
			assert.Equal(t, tt.expectedRegister.overflowFlag(), cpu.registers.overflowFlag(), "unexpected OverflowFlag")
			assert.True(t, extraCycle)
		})
	}
}

func TestBCC(t *testing.T) {
	cases := []struct {
		carryFlag      byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{0, 0x02, 0x04, 1},
		{0, 0x02, 0x0200, 2},
		{1, 0x02, 0x02, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.updateFlag(carryFlag, tt.carryFlag)
			cpu.registers.Pc = tt.pc

			cpu.bcc(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBCS(t *testing.T) {
	cases := []struct {
		name           string
		carry          byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when carry is set", 1, 0x0000, 0x0010, 1},
		{"branches when carry is set, crossing page", 1, 0x0000, 0x0110, 2},
		{"does not branch when carry is unset", 0, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.updateFlag(carryFlag, tt.carry)
			cpu.registers.Pc = tt.pc
			cpu.bcs(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBEQ(t *testing.T) {
	cases := []struct {
		name           string
		zero           byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when zero is set", 1, 0x0000, 0x0010, 1},
		{"branches when zero is set, crossing page", 1, 0x0000, 0x0110, 2},
		{"does not branch when zero is unset", 0, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.updateFlag(zeroFlag, tt.zero)
			cpu.registers.Pc = tt.pc

			cpu.beq(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBIT(t *testing.T) {
	expectedRegisters := func(accumulator byte, negative byte, zero byte, overflow byte) Cpu6502Registers {
		registers := CreateRegisters()
		registers.A = accumulator
		registers.updateFlag(negativeFlag, negative)
		registers.updateFlag(zeroFlag, zero)
		registers.updateFlag(overflowFlag, overflow)
		return registers
	}

	cases := []struct {
		name             string
		accumulator      byte
		operand          byte
		expectedRegister Cpu6502Registers
	}{
		// All disabled
		{"", 0b00001111, 0b00001111, expectedRegisters(0b00001111, 0, 0, 0)},
		// Negative true
		{"", 0b10001111, 0b10001111, expectedRegisters(0b10001111, 1, 0, 0)},
		// Overflow true
		{"", 0b01001111, 0b01001111, expectedRegisters(0b01001111, 0, 0, 1)},
		// Zero Flag true
		{"", 0, 0, expectedRegisters(0, 0, 1, 0)},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.accumulator
			cpu.memory.Write(0x0001, tt.operand)

			cpu.bit(OperationMethodArgument{ZeroPage, 0x0001})

			assert.Equal(t, tt.expectedRegister, cpu.registers, "invalid registers")
		})
	}
}

func TestBMI(t *testing.T) {
	cases := []struct {
		failError      string
		negative       byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when negative is set", 1, 0x0000, 0x0010, 1},
		{"branches when negative is set, crossing page", 1, 0x0000, 0x0110, 2},
		{"does not branch when negative is unset", 0, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.failError, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.updateFlag(negativeFlag, tt.negative)
			cpu.registers.Pc = tt.pc

			cpu.bmi(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBNE(t *testing.T) {
	type dataProvider struct {
		failError      string
		zero           byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}

	dataProviders := [...]dataProvider{
		{"branches when zero is set", 0, 0x0000, 0x0010, 1},
		{"branches when zero is set, crossing page", 0, 0x0000, 0x0110, 2},
		{"does not branch when zero is unset", 1, 0x0000, 0x0000, 0},
	}

	for _, dp := range dataProviders {
		cpu := CreateCPUWithGamePak()
		cpu.registers.Pc = dp.pc
		cpu.registers.updateFlag(zeroFlag, dp.zero)
		cpu.bne(OperationMethodArgument{Relative, dp.expectedPc})

		assert.Equal(t, dp.expectedPc, cpu.registers.Pc, dp.failError+":invalid pc")
		assert.Equal(t, dp.expectedCycles, cpu.opCyclesLeft, dp.failError+":invalid cycles")
	}
}

func TestBPL(t *testing.T) {
	cases := []struct {
		name           string
		negative       byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when negative is set", 0, 0x0000, 0x0010, 1},
		{"branches when negative is set, crossing page", 0, 0x0000, 0x0110, 2},
		{"does not branch when negative is unset", 1, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = tt.pc
			cpu.registers.updateFlag(negativeFlag, tt.negative)

			cpu.bpl(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBRK(t *testing.T) {
	programCounter := Address(0x2030)
	expectedPc := Address(0x9999)
	cpu := CreateCPUWithGamePak()
	cpu.registers.Pc = programCounter
	cpu.registers.Status = 0b11100011
	cpu.memory.Write(Address(0xFFFE), LowNibble(expectedPc))
	cpu.memory.Write(Address(0xFFFF), HighNibble(expectedPc))

	cpu.brk(OperationMethodArgument{Implicit, 0x0000})

	assert.Equal(t, programCounter, cpu.read16(0x1FE))
	// Stored status Registers in stack should be...
	assert.Equal(t, byte(0b11110011), cpu.memory.Read(0x1FD))
	assert.Equal(t, byte(1), cpu.registers.interruptFlag())
	assert.Equal(t, byte(0xF3), cpu.popStack(), "unexpected StatusRegister pushed in stack")
	assert.Equal(t, LowNibble(programCounter), cpu.popStack(), "unexpected low nibble in stack pointer")
	assert.Equal(t, HighNibble(programCounter), cpu.popStack(), "unexpected high nibble in stack pointer")

	assert.Equal(t, expectedPc, cpu.registers.Pc)
}

func TestBVC(t *testing.T) {
	cases := []struct {
		name           string
		overflow       byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when overflow is set", 0, 0x0000, 0x0010, 1},
		{"branches when overflow is set, crossing page", 0, 0x0000, 0x0110, 2},
		{"does not branch when overflow is unset", 1, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = tt.pc
			cpu.registers.updateFlag(overflowFlag, tt.overflow)
			cpu.bvc(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestBVS(t *testing.T) {
	type dataProvider struct {
	}

	cases := []struct {
		name           string
		overflow       byte
		pc             Address
		expectedPc     Address
		expectedCycles byte
	}{
		{"branches when overflow is set", 1, 0x0000, 0x0010, 1},
		{"branches when overflow is set, crossing page", 1, 0x0000, 0x0110, 2},
		{"does not branch when overflow is unset", 0, 0x0000, 0x0000, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Pc = tt.pc
			cpu.registers.updateFlag(overflowFlag, tt.overflow)
			cpu.bvs(OperationMethodArgument{Relative, tt.expectedPc})

			assert.Equal(t, tt.expectedPc, cpu.registers.Pc, "invalid pc")
			assert.Equal(t, tt.expectedCycles, cpu.opCyclesLeft, "invalid cycles")
		})
	}
}

func TestCLC(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.updateFlag(carryFlag, 1)

	cpu.clc(OperationMethodArgument{Implicit, 0x00})

	assert.Zero(t, cpu.registers.carryFlag())
}

func TestCLD(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.updateFlag(decimalFlag, 1)

	cpu.cld(OperationMethodArgument{Implicit, 0x00})

	assert.Zero(t, cpu.registers.decimalFlag())
}

func TestCLI(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.updateFlag(interruptFlag, 1)

	cpu.cli(OperationMethodArgument{Implicit, 0x00})

	assert.Zero(t, cpu.registers.interruptFlag())
}

func TestCLV(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.updateFlag(overflowFlag, 1)

	cpu.clv(OperationMethodArgument{Implicit, 0x00})

	assert.Zero(t, cpu.registers.overflowFlag())
}

func TestCompareOperations(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.X = 0x10
	cpu.registers.A = 0x10
	cpu.registers.Y = 0x10

	type dataProvider struct {
		title            string
		operand          byte
		op               OperationMethod
		expectedCarry    byte
		expectedZero     byte
		expectedNegative byte
	}

	cases := []struct {
		title            string
		operand          byte
		op               OperationMethod
		expectedCarry    byte
		expectedZero     byte
		expectedNegative byte
	}{
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

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			cpu.memory.Write(0x00, tt.operand)

			tt.op(OperationMethodArgument{ZeroPage, Address(0x00)})

			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag(), "unexpected Carry")
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected Zero")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "unexpected Negative")
		})
	}
}

func TestDEC(t *testing.T) {
	cases := []struct {
		name             string
		initialValue     byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"result > 0", 0x02, 0x01, 0, 0},
		{"result is 0", 0x01, 0x00, 1, 0},
		{"result < 0", 0x00, 0xFF, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(0x0000, tt.initialValue)

			cpu.dec(OperationMethodArgument{ZeroPage, Address(0x0000)})

			assert.Equal(t, tt.expectedValue, cpu.memory.Read(0))
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
		})
	}
}

func TestDECXY(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.X = 2
	cpu.registers.Y = 2

	cases := []struct {
		title            string
		op               OperationMethod
		expectedX        byte
		expectedY        byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"X=2", cpu.dex, 1, 2, 0, 0},
		{"X=1", cpu.dex, 0, 2, 1, 0},
		{"X=0", cpu.dex, 0xFF, 2, 0, 1},
		{"Y=2", cpu.dey, 0xFF, 1, 0, 0},
		{"Y=1", cpu.dey, 0xFF, 0, 1, 0},
		{"Y=0", cpu.dey, 0xFF, 0xFF, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			tt.op(OperationMethodArgument{Implicit, 0})
			assert.Equal(t, tt.expectedX, cpu.registers.X, "unexpected X")
			assert.Equal(t, tt.expectedY, cpu.registers.Y, "unexpected Y")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "unexpected Negative flag")
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected Zero flag")
		})
	}
}

func TestEOR(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		a                byte
		expectedA        byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x00, 0x00, 0x00, 1, 0},
		{"", 0x01, 0x00, 0x01, 0, 0},
		{"", 0x80, 0x00, 0x80, 0, 1},
	}

	cpu := CreateCPUWithGamePak()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu.registers.A = tt.a
			cpu.memory.Write(0x05, tt.value)

			extraCycle := cpu.eor(OperationMethodArgument{Immediate, 0x05})

			assert.Equal(t, tt.expectedA, cpu.registers.A)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.True(t, extraCycle)
		})
	}
}

func TestINC(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x00, 0x01, 0, 0},
		{"", 0x7F, 0x80, 0, 1},
		{"", 0xFF, 0x00, 1, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(0x00, tt.value)

			cpu.inc(OperationMethodArgument{ZeroPage, 0x00})

			assert.Equal(t, tt.expectedValue, cpu.memory.Read(0x00))
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
		})
	}
}

func TestINX(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x00, 0x01, 0, 0},
		{"", 0x7F, 0x80, 0, 1},
		{"", 0xFF, 0x00, 1, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.X = tt.value

			cpu.inx(OperationMethodArgument{ZeroPage, 0x00})

			assert.Equal(t, tt.expectedValue, cpu.registers.X)
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
		})
	}
}

func TestINY(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedValue    byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x00, 0x01, 0, 0},
		{"", 0x7F, 0x80, 0, 1},
		{"", 0xFF, 0x00, 1, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Y = tt.value

			cpu.iny(OperationMethodArgument{ZeroPage, 0x00})

			assert.Equal(t, tt.expectedValue, cpu.registers.Y)
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
		})
	}
}

func TestJMP(t *testing.T) {
	cpu := CreateCPUWithGamePak()

	cpu.jmp(OperationMethodArgument{Absolute, 0x100})

	assert.Equal(t, Address(0x100), cpu.registers.Pc)
}

func TestJSR(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.memory.Write(Address(0x201), 0x20) // JSR Opcode
	cpu.memory.Write(Address(0x202), 0x55) // LSB
	cpu.memory.Write(Address(0x203), 0x05) // MSB

	cpu.registers.Pc = 0x0203
	cpu.jsr(OperationMethodArgument{Absolute, 0x202})

	assert.Equal(t, Address(0x0202), cpu.registers.Pc)
	assert.Equal(t, byte(0x02), cpu.popStack())
	assert.Equal(t, byte(0x02), cpu.popStack())
}

func TestLDA(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x20, 0, 0},
		{"", 0x00, 1, 0},
		{"", 0x80, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(Address(0x00), tt.value)

			extraCycle := cpu.lda(OperationMethodArgument{Immediate, 0x00})

			assert.Equal(t, tt.value, cpu.registers.A)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.True(t, extraCycle)
		})
	}
}

func TestLDX(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x20, 0, 0},
		{"", 0x00, 1, 0},
		{"", 0x80, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(Address(0x00), tt.value)

			extraCycle := cpu.ldx(OperationMethodArgument{Immediate, 0x00})

			assert.Equal(t, tt.value, cpu.registers.X)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.True(t, extraCycle)
		})
	}
}

func TestLDY(t *testing.T) {
	cases := []struct {
		name             string
		value            byte
		expectedZero     byte
		expectedNegative byte
	}{
		{"", 0x20, 0, 0},
		{"", 0x00, 1, 0},
		{"", 0x80, 0, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.memory.Write(Address(0x00), tt.value)

			extraCycle := cpu.ldy(OperationMethodArgument{Immediate, 0x00})

			assert.Equal(t, tt.value, cpu.registers.Y)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.True(t, extraCycle)
		})
	}
}

func TestLSR(t *testing.T) {
	cases := []struct {
		name           string
		addressingMode AddressMode
		value          byte
		expectedResult byte
		expectedZero   byte
		expectedCarry  byte
	}{
		{"", Implicit, 0b00000010, 0b00000001, 0, 0},
		{"", Implicit, 0b00000011, 0b00000001, 0, 1},
		{"", Implicit, 0b00000001, 0b00000000, 1, 1},
		{"", ZeroPage, 0b00000010, 0b00000001, 0, 0},
		{"", ZeroPage, 0b00000011, 0b00000001, 0, 1},
		{"", ZeroPage, 0b00000001, 0b00000000, 1, 1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Status = 0xFF
			cpu.registers.A = tt.value
			cpu.memory.Write(Address(0x00), tt.value)

			cpu.lsr(OperationMethodArgument{tt.addressingMode, 0x00})

			if tt.addressingMode == Implicit {
				assert.Equal(t, tt.expectedResult, cpu.registers.A)
			} else {
				assert.Equal(t, tt.expectedResult, cpu.memory.Read(0x00))
			}
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected ZeroFlag")
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag())
			assert.Zero(t, cpu.registers.negativeFlag())
		})
	}
}

func TestORA(t *testing.T) {
	cases := []struct {
		a                byte
		value            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
	}{
		{0x00, 0x00, 0x00, 1, 0},
		{0x80, 0x00, 0x80, 0, 1},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.a
			cpu.memory.Write(0x00, tt.value)

			extraCycle := cpu.ora(OperationMethodArgument{Immediate, 0x00})

			assert.Equal(t, tt.expectedResult, cpu.registers.A)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.True(t, extraCycle)
		})
	}
}

func TestPHA(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.A = 0x30
	cpu.pha(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(0x30), cpu.popStack())
}

func TestPHP(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Status = 0b11001111

	cpu.php(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(0xFF), cpu.popStack())
}

func TestPLA(t *testing.T) {
	cases := []struct {
		pulledValue      byte
		expectedNegative byte
		expectedZero     byte
	}{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.pushStack(tt.pulledValue)

			cpu.pla(OperationMethodArgument{Implicit, 0x00})

			assert.Equal(t, tt.pulledValue, cpu.registers.A)
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, byte(0xFF), cpu.registers.Sp)
		})
	}
}

func TestPLP(t *testing.T) {
	initialPointerStack := byte(0xFF)
	cpu := CreateCPUWithGamePak()
	cpu.pushStack(0xFF)

	cpu.plp(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(0), (cpu.registers.Status>>4)&0x01, "PLP must set B flag to 0")
	assert.Equal(t, byte(0xEF), cpu.registers.Status)
	assert.Equal(t, initialPointerStack, cpu.registers.Sp)
}

func TestROL(t *testing.T) {
	cases := []struct {
		addressingMode   AddressMode
		value            byte
		carry            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
	}{
		{Implicit, 0b00000000, 0, 0, 1, 0, 0},
		{Implicit, 0b00000000, 1, 1, 0, 0, 0},
		{Implicit, 0b00000001, 0, 0b10, 0, 0, 0},
		{Implicit, 0b10000000, 0, 0, 1, 0, 1},
		{Implicit, 0b01000000, 0, 0x80, 0, 1, 0},

		{ZeroPage, 0b00000000, 0, 0, 1, 0, 0},
		{ZeroPage, 0b00000000, 1, 1, 0, 0, 0},
		{ZeroPage, 0b00000001, 0, 0b10, 0, 0, 0},
		{ZeroPage, 0b10000000, 0, 0, 1, 0, 1},
		{ZeroPage, 0b01000000, 0, 0x80, 0, 1, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.value
			cpu.registers.updateFlag(carryFlag, tt.carry)
			cpu.memory.Write(Address(0x00), tt.value)

			cpu.rol(OperationMethodArgument{tt.addressingMode, 0x00})

			if tt.addressingMode == Implicit {
				assert.Equal(t, tt.expectedResult, cpu.registers.A)
			} else {
				assert.Equal(t, tt.expectedResult, cpu.memory.Read(0x00))
			}
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected ZeroFlag")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "unexpected Negative")
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag())
		})
	}
}

func TestROR(t *testing.T) {
	cases := []struct {
		addressingMode   AddressMode
		value            byte
		carry            byte
		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
	}{
		{Implicit, 0b00000000, 0, 0, 1, 0, 0},
		{Implicit, 0b00000001, 0, 0, 1, 0, 1},
		{Implicit, 0b00000000, 1, 0x80, 0, 1, 0},
		{Implicit, 0b10000000, 0, 0x40, 0, 0, 0},
		{Implicit, 0b10000001, 1, 0xC0, 0, 1, 1},

		{ZeroPage, 0b00000000, 0, 0, 1, 0, 0},
		{ZeroPage, 0b00000001, 0, 0, 1, 0, 1},
		{ZeroPage, 0b00000000, 1, 0x80, 0, 1, 0},
		{ZeroPage, 0b10000000, 0, 0x40, 0, 0, 0},
		{ZeroPage, 0b10000001, 1, 0xC0, 0, 1, 1},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.value
			cpu.registers.updateFlag(carryFlag, tt.carry)
			cpu.memory.Write(Address(0x00), tt.value)

			cpu.ror(OperationMethodArgument{tt.addressingMode, 0x00})

			if tt.addressingMode == Implicit {
				assert.Equal(t, tt.expectedResult, cpu.registers.A)
			} else {
				assert.Equal(t, tt.expectedResult, cpu.memory.Read(0x00))
			}
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag())
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected ZeroFlag")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "unexpected Negative")
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag())
		})
	}
}

func TestRTI(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	// Push an Address into Stack
	pc := Address(0x532)
	cpu.pushStack(HighNibble(pc))
	cpu.pushStack(LowNibble(pc))
	// Push a StatusRegister into stack
	cpu.pushStack(0xFF)

	cpu.rti(OperationMethodArgument{AddressMode: Implicit, OperandAddress: 0xFF})

	assert.Equal(t, pc, cpu.registers.Pc)
	assert.Equal(t, byte(0xeF), cpu.registers.Status)
}

func TestRTS(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	// Push an Address into Stack
	pc := Address(0x532)
	cpu.pushStack(HighNibble(pc))
	cpu.pushStack(LowNibble(pc))

	cpu.rts(OperationMethodArgument{Implicit, 0x00})

	expectedProgramCounter := Address(0x533)
	assert.Equal(t, expectedProgramCounter, cpu.registers.Pc)
}

func TestSBC(t *testing.T) {
	// Fixtures taken from http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html (section SBC)
	// TODO: Improve these fixtures by adding more (if really needed)
	cases := []struct {
		a     byte
		value byte
		carry byte

		expectedResult   byte
		expectedZero     byte
		expectedNegative byte
		expectedCarry    byte
		expectedOverflow byte
	}{
		{0x01, 1, 1, 0, 1, 0, 1, 0},
		{0x50, 0xF0, 1, 0x60, 0, 0, 0, 0},
		{0x50, 0xB0, 1, 0xA0, 0, 1, 0, 1},
		{0x50, 0x70, 1, 0xE0, 0, 1, 0, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.A = tt.a
			cpu.registers.updateFlag(carryFlag, tt.carry)
			cpu.memory.Write(0x00, tt.value)

			extraCycle := cpu.sbc(OperationMethodArgument{Immediate, 0x00})

			assert.Equal(t, tt.expectedResult, cpu.registers.A, "Invalid subtraction result")
			assert.Equal(t, tt.expectedCarry, cpu.registers.carryFlag(), "Invalid CarryFlag")
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "Invalid zeroflag")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "Invalid negative Flag")
			assert.Equal(t, tt.expectedOverflow, cpu.registers.overflowFlag(), "Invalid Overflow Flag")
			assert.True(t, extraCycle)
		})
	}
}

func TestSEC(t *testing.T) {
	cpu := CreateCPUWithGamePak()

	cpu.sec(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.carryFlag())
}

func TestSED(t *testing.T) {
	cpu := CreateCPUWithGamePak()

	cpu.sed(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.decimalFlag())
}

func TestSEI(t *testing.T) {
	cpu := CreateCPUWithGamePak()

	cpu.sei(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(1), cpu.registers.interruptFlag())
}

func TestSTA(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.A = 0xFF

	cpu.sta(OperationMethodArgument{Implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.memory.Read(0x522))
}

func TestSTX(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.X = 0xFF

	cpu.stx(OperationMethodArgument{Implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.memory.Read(0x522))
}

func TestSTY(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.Y = 0xFF

	cpu.sty(OperationMethodArgument{Implicit, 0x522})

	assert.Equal(t, byte(0xFF), cpu.memory.Read(0x522))
}

func TestTAX_TAY(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cases := []struct {
		name             string
		op               OperationMethod
		a                byte
		expectedNegative byte
		expectedZero     byte
	}{
		{"tax", cpu.tax, 0x00, 0, 1},
		{"tax", cpu.tax, 0x80, 1, 0},
		{"tax", cpu.tax, 0x20, 0, 0},
		{"tay", cpu.tay, 0x00, 0, 1},
		{"tay", cpu.tay, 0x80, 1, 0},
		{"tay", cpu.tay, 0x20, 0, 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cpu.registers.Reset()
			cpu.registers.A = tt.a

			tt.op(OperationMethodArgument{Implicit, 0x00})

			if tt.name == "tax" {
				assert.Equal(t, cpu.registers.A, cpu.registers.X, "unexpected X")
			} else {
				assert.Equal(t, cpu.registers.A, cpu.registers.Y, "unexpected Y")
			}
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "unexpected Negative flag")
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "unexpected Zero flag")
		})
	}
}

func TestTSX(t *testing.T) {
	cases := []struct {
		sp               byte
		expectedNegative byte
		expectedZero     byte
	}{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.Registers().setStackPointer(tt.sp)

			cpu.tsx(OperationMethodArgument{Implicit, 0x00})

			assert.Equal(t, tt.sp, cpu.registers.X)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag(), "Incorrect Zero Flag")
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag(), "Incorrect Negative Flag")
		})
	}
}

func TestTXA(t *testing.T) {
	cases := []struct {
		x                byte
		expectedNegative byte
		expectedZero     byte
	}{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.X = tt.x

			cpu.txa(OperationMethodArgument{Implicit, 0x00})

			assert.Equal(t, tt.x, cpu.registers.A)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
		})
	}
}

func TestTXS(t *testing.T) {
	cpu := CreateCPUWithGamePak()
	cpu.registers.X = 0xFF

	cpu.txs(OperationMethodArgument{Implicit, 0x00})

	assert.Equal(t, byte(0xFF), cpu.registers.Sp)
}

func TestTYA(t *testing.T) {
	cases := []struct {
		y                byte
		expectedNegative byte
		expectedZero     byte
	}{
		{0x00, 0, 1},
		{0x80, 1, 0},
		{0x20, 0, 0},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			cpu := CreateCPUWithGamePak()
			cpu.registers.Y = tt.y

			cpu.tya(OperationMethodArgument{Implicit, 0x00})

			assert.Equal(t, cpu.registers.Y, cpu.registers.A)
			assert.Equal(t, tt.expectedZero, cpu.registers.zeroFlag())
			assert.Equal(t, tt.expectedNegative, cpu.registers.negativeFlag())
		})
	}
}
