package nes

import (
	"fmt"
)

// AddressMode is an enum of the available Addressing Modes in this cpu
type AddressMode int
type AddressModeMethod func(programCounter Address) (finalAddress Address, operand [3]byte, cycles int, pageCrossed bool)

const (
	Implicit AddressMode = iota
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteXIndexed
	AbsoluteYIndexed
	Indirect
	IndirectX
	IndirectY
	Relative
)

// Cpu6502 Represents a CPU 6502
type Cpu6502 struct {
	registers Cpu6502Registers
	memory    Memory

	instructions [256]Instruction
	opCyclesLeft byte // How many cycles left to finish execution of current cycle
	cycle        uint32

	addressEvaluators [13]AddressModeMethod

	// Debug parameters
	debug       bool
	Logger      cpu6502Logger
	cyclesLimit uint16
}

type Cpu6502DebugOptions struct {
	enabled       bool
	outputLogPath string
}

func CreateCPU(memory Memory, debug Cpu6502DebugOptions) *Cpu6502 {
	cpu := Cpu6502{
		memory:    memory,
		registers: CreateRegisters(),
		debug:     debug.enabled,
	}

	if debug.enabled {
		cpu.Logger = createCPULogger(debug.outputLogPath)
	}

	cpu.initInstructionsTable()
	cpu.initAddressModeEvaluators()

	return &cpu
}

func (cpu Cpu6502) ProgramCounter() Address {
	return cpu.Registers().Pc
}

func (cpu *Cpu6502) Registers() *Cpu6502Registers {
	return &cpu.registers
}

func (cpu *Cpu6502) pushStack(value byte) {
	address := cpu.registers.stackPointerAddress()
	cpu.memory.Write(
		address,
		value,
	)

	cpu.registers.stackPointerPushed()
}

func (cpu *Cpu6502) popStack() byte {
	cpu.registers.stackPointerPopped()
	address := cpu.registers.stackPointerAddress()
	return cpu.memory.Read(address)
}

// Reads value located at Program Counter and increments it
func (cpu *Cpu6502) fetch() byte {
	value := cpu.memory.Read(cpu.registers.Pc)
	cpu.registers.Pc++

	return value
}

/*
func (cpu *Cpu6502) Read(address defs.Address) byte {
	return cpu.memory.Read(address)
}
*/
func (cpu *Cpu6502) read16(address Address) Word {
	low := cpu.memory.Read(address)
	high := cpu.memory.Read(address + 1)

	return CreateWord(low, high)
}

func (cpu *Cpu6502) read16Bugged(address Address) Word {
	lsb := address
	msb := (lsb & 0xFF00) | Address(byte(lsb)+1)

	low := cpu.memory.Read(lsb)
	high := cpu.memory.Read(msb)

	return CreateWord(low, high)
}

func (cpu *Cpu6502) initInstructionsTable() {
	cpu.instructions = [256]Instruction{
		CreateInstruction("BRK", Implicit, cpu.brk, 7, 1),
		CreateInstruction("ORA", IndirectX, cpu.ora, 6, 2),
		{},
		{},
		{},
		CreateInstruction("ORA", ZeroPage, cpu.ora, 3, 2),
		CreateInstruction("ASL", ZeroPage, cpu.asl, 5, 2),
		{},
		CreateInstruction("PHP", Implicit, cpu.php, 3, 1),
		CreateInstruction("ORA", Immediate, cpu.ora, 2, 2),
		CreateInstruction("ASL", Implicit, cpu.asl, 2, 1),
		{},
		{},
		CreateInstruction("ORA", Absolute, cpu.ora, 4, 3),
		CreateInstruction("ASL", Absolute, cpu.asl, 6, 3),
		{},

		// 0x10
		CreateInstruction("BPL", Relative, cpu.bpl, 2, 2),
		CreateInstruction("ORA", IndirectY, cpu.ora, 5, 2),
		{},
		{},
		{},
		CreateInstruction("ORA", ZeroPageX, cpu.ora, 4, 2),
		CreateInstruction("ASL", ZeroPageX, cpu.asl, 6, 2),
		{},
		CreateInstruction("CLC", Implicit, cpu.clc, 2, 1),
		CreateInstruction("ORA", AbsoluteYIndexed, cpu.ora, 4, 3),
		{},
		{},
		{},
		CreateInstruction("ORA", AbsoluteXIndexed, cpu.ora, 4, 3),
		CreateInstruction("ASL", AbsoluteXIndexed, cpu.asl, 7, 3),
		{},

		// 0x20
		CreateInstruction("JSR", Absolute, cpu.jsr, 6, 3),
		CreateInstruction("AND", IndirectX, cpu.and, 6, 2),
		{},
		{},
		CreateInstruction("BIT", ZeroPage, cpu.bit, 3, 2),
		CreateInstruction("AND", ZeroPage, cpu.and, 3, 2),
		CreateInstruction("ROL", ZeroPage, cpu.rol, 5, 2),
		{},
		CreateInstruction("PLP", Implicit, cpu.plp, 4, 1),
		CreateInstruction("AND", Immediate, cpu.and, 2, 2),
		CreateInstruction("ROL", Implicit, cpu.rol, 2, 1),
		{},
		CreateInstruction("BIT", Absolute, cpu.bit, 4, 3),
		CreateInstruction("AND", Absolute, cpu.and, 4, 3),
		CreateInstruction("ROL", Absolute, cpu.rol, 6, 3),
		{},

		// 0x30
		CreateInstruction("BMI", Relative, cpu.bmi, 2, 2),
		CreateInstruction("AND", IndirectY, cpu.and, 5, 2),
		{},
		{},
		{},
		CreateInstruction("AND", ZeroPageX, cpu.and, 4, 2),
		CreateInstruction("ROL", ZeroPageX, cpu.rol, 6, 2),
		{},
		CreateInstruction("SEC", Implicit, cpu.sec, 2, 1),
		CreateInstruction("AND", AbsoluteYIndexed, cpu.and, 4, 3),
		{},
		{},
		{},
		CreateInstruction("AND", AbsoluteXIndexed, cpu.and, 4, 3),
		CreateInstruction("ROL", AbsoluteXIndexed, cpu.rol, 7, 3),
		{},

		// 0x40
		CreateInstruction("RTI", Implicit, cpu.rti, 6, 1),
		CreateInstruction("EOR", IndirectX, cpu.eor, 6, 2),
		{},
		{},
		{},
		CreateInstruction("EOR", ZeroPage, cpu.eor, 3, 2),
		CreateInstruction("LSR", ZeroPage, cpu.lsr, 5, 2),
		{},
		CreateInstruction("PHA", Implicit, cpu.pha, 3, 1),
		CreateInstruction("EOR", Immediate, cpu.eor, 2, 2),
		CreateInstruction("LSR", Implicit, cpu.lsr, 2, 1),
		{},
		CreateInstruction("JMP", Absolute, cpu.jmp, 3, 3),
		CreateInstruction("EOR", Absolute, cpu.eor, 4, 3),
		CreateInstruction("LSR", Absolute, cpu.lsr, 6, 3),
		{},

		// 0x50
		CreateInstruction("BVC", Relative, cpu.bvc, 2, 2),
		CreateInstruction("EOR", IndirectY, cpu.eor, 5, 2),
		{},
		{},
		{},
		CreateInstruction("EOR", ZeroPageX, cpu.eor, 4, 2),
		CreateInstruction("LSR", ZeroPageX, cpu.lsr, 6, 2),
		{},
		CreateInstruction("CLI", Implicit, cpu.cli, 2, 1),
		CreateInstruction("EOR", AbsoluteYIndexed, cpu.eor, 4, 3),
		{},
		{},
		{},
		CreateInstruction("EOR", AbsoluteXIndexed, cpu.eor, 4, 3),
		CreateInstruction("LSR", AbsoluteXIndexed, cpu.lsr, 7, 3),
		{},

		// 0x60
		CreateInstruction("RTS", Implicit, cpu.rts, 6, 1),
		CreateInstruction("ADC", IndirectX, cpu.adc, 6, 2),
		{},
		{},
		{},
		CreateInstruction("ADC", ZeroPage, cpu.adc, 3, 2),
		CreateInstruction("ROR", ZeroPage, cpu.ror, 5, 2),
		{},
		CreateInstruction("PLA", Implicit, cpu.pla, 4, 1),
		CreateInstruction("ADC", Immediate, cpu.adc, 2, 2),
		CreateInstruction("ROR", Implicit, cpu.ror, 2, 1),
		{},
		CreateInstruction("JMP", Indirect, cpu.jmp, 5, 3),
		CreateInstruction("ADC", Absolute, cpu.adc, 4, 3),
		CreateInstruction("ROR", Absolute, cpu.ror, 6, 3),
		{},

		// 0x70
		CreateInstruction("BVS", Relative, cpu.bvs, 2, 2),
		CreateInstruction("ADC", IndirectY, cpu.adc, 5, 2),
		{},
		{},
		{},
		CreateInstruction("ADC", ZeroPageX, cpu.adc, 4, 2),
		CreateInstruction("ROR", ZeroPageX, cpu.ror, 6, 2),
		{},
		CreateInstruction("SEI", Implicit, cpu.sei, 2, 1),
		CreateInstruction("ADC", AbsoluteYIndexed, cpu.adc, 4, 3),
		{},
		{},
		{},
		CreateInstruction("ADC", AbsoluteXIndexed, cpu.adc, 4, 3),
		CreateInstruction("ROR", AbsoluteXIndexed, cpu.ror, 7, 3),
		{},

		// 0x80
		{},
		CreateInstruction("STA", IndirectX, cpu.sta, 6, 2),
		{},
		{},
		CreateInstruction("STY", ZeroPage, cpu.sty, 3, 2),
		CreateInstruction("STA", ZeroPage, cpu.sta, 3, 2),
		CreateInstruction("STX", ZeroPage, cpu.stx, 3, 2),
		{},
		CreateInstruction("DEY", Implicit, cpu.dey, 2, 1),
		{},
		CreateInstruction("TXA", Implicit, cpu.txa, 2, 1),
		{},
		CreateInstruction("STY", Absolute, cpu.sty, 4, 3),
		CreateInstruction("STA", Absolute, cpu.sta, 4, 3),
		CreateInstruction("STX", Absolute, cpu.stx, 4, 3),
		{},

		// 0x90
		CreateInstruction("BCC", Relative, cpu.bcc, 2, 2),
		CreateInstruction("STA", IndirectY, cpu.sta, 6, 2),
		{},
		{},
		CreateInstruction("STY", ZeroPageX, cpu.sty, 4, 2),
		CreateInstruction("STA", ZeroPageX, cpu.sta, 4, 2),
		CreateInstruction("STX", ZeroPageY, cpu.stx, 4, 2),
		{},
		CreateInstruction("TYA", Implicit, cpu.tya, 2, 1),
		CreateInstruction("STA", AbsoluteYIndexed, cpu.sta, 5, 3),
		CreateInstruction("TXS", Implicit, cpu.txs, 2, 1),
		{},
		{},
		CreateInstruction("STA", AbsoluteXIndexed, cpu.sta, 5, 3),
		{},
		{},

		// 0xA0
		CreateInstruction("LDY", Immediate, cpu.ldy, 2, 2),
		CreateInstruction("LDA", IndirectX, cpu.lda, 6, 2),
		CreateInstruction("LDX", Immediate, cpu.ldx, 2, 2),
		{},
		CreateInstruction("LDY", ZeroPage, cpu.ldy, 3, 2),
		CreateInstruction("LDA", ZeroPage, cpu.lda, 3, 2),
		CreateInstruction("LDX", ZeroPage, cpu.ldx, 3, 2),
		{},
		CreateInstruction("TAY", Implicit, cpu.tay, 2, 1),
		CreateInstruction("LDA", Immediate, cpu.lda, 2, 2),
		CreateInstruction("TAX", Implicit, cpu.tax, 2, 1),
		{},
		CreateInstruction("LDY", Absolute, cpu.ldy, 4, 3),
		CreateInstruction("LDA", Absolute, cpu.lda, 4, 3),
		CreateInstruction("LDX", Absolute, cpu.ldx, 4, 3),
		{},

		// 0xB0
		CreateInstruction("BCS", Relative, cpu.bcs, 2, 2),
		CreateInstruction("LDA", IndirectY, cpu.lda, 5, 2),
		{},
		{},
		CreateInstruction("LDY", ZeroPageX, cpu.ldy, 4, 2),
		CreateInstruction("LDA", ZeroPageX, cpu.lda, 4, 2),
		CreateInstruction("LDX", ZeroPageY, cpu.ldx, 4, 2),
		{},
		CreateInstruction("CLV", Implicit, cpu.clv, 2, 1),
		CreateInstruction("LDA", AbsoluteYIndexed, cpu.lda, 4, 3),
		CreateInstruction("TSX", Implicit, cpu.tsx, 2, 1),
		{},
		CreateInstruction("LDY", AbsoluteXIndexed, cpu.ldy, 4, 3),
		CreateInstruction("LDA", AbsoluteXIndexed, cpu.lda, 4, 3),
		CreateInstruction("LDX", AbsoluteYIndexed, cpu.ldx, 4, 3),
		{},

		// 0xC0
		CreateInstruction("CPY", Immediate, cpu.cpy, 2, 2),
		CreateInstruction("CMP", IndirectX, cpu.cmp, 6, 2),
		{},
		{},
		CreateInstruction("CPY", ZeroPage, cpu.cpy, 3, 2),
		CreateInstruction("CMP", ZeroPage, cpu.cmp, 3, 2),
		CreateInstruction("DEC", ZeroPage, cpu.dec, 5, 2),
		{},
		CreateInstruction("INY", Implicit, cpu.iny, 2, 1),
		CreateInstruction("CMP", Immediate, cpu.cmp, 2, 2),
		CreateInstruction("DEX", Implicit, cpu.dex, 2, 1),
		{},
		CreateInstruction("CPY", Absolute, cpu.cpy, 4, 3),
		CreateInstruction("CMP", Absolute, cpu.cmp, 4, 3),
		CreateInstruction("DEC", Absolute, cpu.dec, 6, 3),
		{},

		// 0xD0
		CreateInstruction("BNE", Relative, cpu.bne, 2, 2),
		CreateInstruction("CMP", IndirectY, cpu.cmp, 5, 2),
		{},
		{},
		{},
		CreateInstruction("CMP", ZeroPageX, cpu.cmp, 4, 2),
		CreateInstruction("DEC", ZeroPageX, cpu.dec, 6, 2),
		{},
		CreateInstruction("CLD", Implicit, cpu.cld, 2, 1),
		CreateInstruction("CMP", AbsoluteYIndexed, cpu.cmp, 4, 3),
		{},
		{},
		{},
		CreateInstruction("CMP", AbsoluteXIndexed, cpu.cmp, 4, 3),
		CreateInstruction("DEC", AbsoluteXIndexed, cpu.dec, 7, 3),
		{},

		// 0xE0
		CreateInstruction("CPX", Immediate, cpu.cpx, 2, 2),
		CreateInstruction("SBC", IndirectX, cpu.sbc, 6, 2),
		{},
		{},
		CreateInstruction("CPX", ZeroPage, cpu.cpx, 3, 2),
		CreateInstruction("SBC", ZeroPage, cpu.sbc, 3, 2),
		CreateInstruction("INC", ZeroPage, cpu.inc, 5, 2),
		{},
		CreateInstruction("INX", Implicit, cpu.inx, 2, 1),
		CreateInstruction("SBC", Immediate, cpu.sbc, 2, 2),
		CreateInstruction("NOP", Implicit, cpu.nop, 2, 1),
		{},
		CreateInstruction("CPX", Absolute, cpu.cpx, 4, 3),
		CreateInstruction("SBC", Absolute, cpu.sbc, 4, 3),
		CreateInstruction("INC", Absolute, cpu.inc, 6, 3),
		{},

		// 0xF0
		CreateInstruction("BEQ", Relative, cpu.beq, 2, 2),
		CreateInstruction("SBC", IndirectY, cpu.sbc, 5, 2),
		{},
		{},
		{},
		CreateInstruction("SBC", ZeroPageX, cpu.sbc, 4, 2),
		CreateInstruction("INC", ZeroPageX, cpu.inc, 6, 2),
		{},
		CreateInstruction("SED", Implicit, cpu.sed, 2, 1),
		CreateInstruction("SBC", AbsoluteYIndexed, cpu.sbc, 4, 3),
		{},
		{},
		{},
		CreateInstruction("SBC", AbsoluteXIndexed, cpu.sbc, 4, 3),
		CreateInstruction("INC", AbsoluteXIndexed, cpu.inc, 7, 3),
		{},
	}
}

func (cpu *Cpu6502) initAddressModeEvaluators() {
	cpu.addressEvaluators = [13]AddressModeMethod{
		cpu.evalImplicit,
		cpu.evalImmediate,
		cpu.evalZeroPage,
		cpu.evalZeroPageX,
		cpu.evalZeroPageY,
		cpu.evalAbsolute,
		cpu.evalAbsoluteXIndexed,
		cpu.evalAbsoluteYIndexed,
		cpu.evalIndirect,
		cpu.evalIndirectX,
		cpu.evalIndirectY,
		cpu.evalRelative,
	}
}

func (cpu Cpu6502) evaluateOperandAddress(addressMode AddressMode, pc Address) (finalAddress Address, operand [3]byte, pageCrossed bool) {
	if addressMode == Implicit {
		finalAddress = 0
		return
	}

	if cpu.addressEvaluators[addressMode] == nil {
		msg := fmt.Errorf("cannot find address evaluator for address mode \"%d\"", addressMode)
		panic(msg)
	}

	finalAddress, operand, _, pageCrossed = cpu.addressEvaluators[addressMode](pc)

	return
}

func memoryPageDiffer(address Address, finalAddress Address) bool {
	return address&0xFF00 != finalAddress&0xFF00
}
