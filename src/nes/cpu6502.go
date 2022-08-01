package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/types"
)

// AddressMode is an enum of the available Addressing Modes in this cpu
type AddressMode int
type AddressModeMethod func(programCounter types.Address) (finalAddress types.Address, operand [3]byte, cycles int, pageCrossed bool)

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
	registers cpu.Registers
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
	cpu6502 := Cpu6502{
		memory:    memory,
		registers: cpu.CreateRegisters(),
		debug:     debug.enabled,
	}

	if debug.enabled {
		cpu6502.Logger = createCPULogger(debug.outputLogPath)
	}

	cpu6502.initInstructionsTable()
	cpu6502.initAddressModeEvaluators()

	return &cpu6502
}

func (cpu6502 Cpu6502) ProgramCounter() types.Address {
	return cpu6502.Registers().Pc
}

func (cpu6502 *Cpu6502) Registers() *cpu.Registers {
	return &cpu6502.registers
}

func (cpu6502 *Cpu6502) pushStack(value byte) {
	address := cpu6502.registers.StackPointerAddress()
	cpu6502.memory.Write(
		address,
		value,
	)

	cpu6502.registers.StackPointerPushed()
}

func (cpu6502 *Cpu6502) popStack() byte {
	cpu6502.registers.StackPointerPopped()
	address := cpu6502.registers.StackPointerAddress()
	return cpu6502.memory.Read(address)
}

// Reads value located at Program Counter and increments it
func (cpu6502 *Cpu6502) fetch() byte {
	value := cpu6502.memory.Read(cpu6502.registers.Pc)
	cpu6502.registers.Pc++

	return value
}

/*
func (cpu *Cpu6502) Read(address defs.Address) byte {
	return cpu.memory.Read(address)
}
*/
func (cpu6502 *Cpu6502) read16(address types.Address) types.Word {
	low := cpu6502.memory.Read(address)
	high := cpu6502.memory.Read(address + 1)

	return types.CreateWord(low, high)
}

func (cpu6502 *Cpu6502) read16Bugged(address types.Address) types.Word {
	lsb := address
	msb := (lsb & 0xFF00) | types.Address(byte(lsb)+1)

	low := cpu6502.memory.Read(lsb)
	high := cpu6502.memory.Read(msb)

	return types.CreateWord(low, high)
}

func (cpu6502 *Cpu6502) initInstructionsTable() {
	cpu6502.instructions = [256]Instruction{
		CreateInstruction("BRK", Implicit, cpu6502.brk, 7, 1),
		CreateInstruction("ORA", IndirectX, cpu6502.ora, 6, 2),
		{},
		{},
		{},
		CreateInstruction("ORA", ZeroPage, cpu6502.ora, 3, 2),
		CreateInstruction("ASL", ZeroPage, cpu6502.asl, 5, 2),
		{},
		CreateInstruction("PHP", Implicit, cpu6502.php, 3, 1),
		CreateInstruction("ORA", Immediate, cpu6502.ora, 2, 2),
		CreateInstruction("ASL", Implicit, cpu6502.asl, 2, 1),
		{},
		{},
		CreateInstruction("ORA", Absolute, cpu6502.ora, 4, 3),
		CreateInstruction("ASL", Absolute, cpu6502.asl, 6, 3),
		{},

		// 0x10
		CreateInstruction("BPL", Relative, cpu6502.bpl, 2, 2),
		CreateInstruction("ORA", IndirectY, cpu6502.ora, 5, 2),
		{},
		{},
		{},
		CreateInstruction("ORA", ZeroPageX, cpu6502.ora, 4, 2),
		CreateInstruction("ASL", ZeroPageX, cpu6502.asl, 6, 2),
		{},
		CreateInstruction("CLC", Implicit, cpu6502.clc, 2, 1),
		CreateInstruction("ORA", AbsoluteYIndexed, cpu6502.ora, 4, 3),
		{},
		{},
		{},
		CreateInstruction("ORA", AbsoluteXIndexed, cpu6502.ora, 4, 3),
		CreateInstruction("ASL", AbsoluteXIndexed, cpu6502.asl, 7, 3),
		{},

		// 0x20
		CreateInstruction("JSR", Absolute, cpu6502.jsr, 6, 3),
		CreateInstruction("AND", IndirectX, cpu6502.and, 6, 2),
		{},
		{},
		CreateInstruction("BIT", ZeroPage, cpu6502.bit, 3, 2),
		CreateInstruction("AND", ZeroPage, cpu6502.and, 3, 2),
		CreateInstruction("ROL", ZeroPage, cpu6502.rol, 5, 2),
		{},
		CreateInstruction("PLP", Implicit, cpu6502.plp, 4, 1),
		CreateInstruction("AND", Immediate, cpu6502.and, 2, 2),
		CreateInstruction("ROL", Implicit, cpu6502.rol, 2, 1),
		{},
		CreateInstruction("BIT", Absolute, cpu6502.bit, 4, 3),
		CreateInstruction("AND", Absolute, cpu6502.and, 4, 3),
		CreateInstruction("ROL", Absolute, cpu6502.rol, 6, 3),
		{},

		// 0x30
		CreateInstruction("BMI", Relative, cpu6502.bmi, 2, 2),
		CreateInstruction("AND", IndirectY, cpu6502.and, 5, 2),
		{},
		{},
		{},
		CreateInstruction("AND", ZeroPageX, cpu6502.and, 4, 2),
		CreateInstruction("ROL", ZeroPageX, cpu6502.rol, 6, 2),
		{},
		CreateInstruction("SEC", Implicit, cpu6502.sec, 2, 1),
		CreateInstruction("AND", AbsoluteYIndexed, cpu6502.and, 4, 3),
		{},
		{},
		{},
		CreateInstruction("AND", AbsoluteXIndexed, cpu6502.and, 4, 3),
		CreateInstruction("ROL", AbsoluteXIndexed, cpu6502.rol, 7, 3),
		{},

		// 0x40
		CreateInstruction("RTI", Implicit, cpu6502.rti, 6, 1),
		CreateInstruction("EOR", IndirectX, cpu6502.eor, 6, 2),
		{},
		{},
		{},
		CreateInstruction("EOR", ZeroPage, cpu6502.eor, 3, 2),
		CreateInstruction("LSR", ZeroPage, cpu6502.lsr, 5, 2),
		{},
		CreateInstruction("PHA", Implicit, cpu6502.pha, 3, 1),
		CreateInstruction("EOR", Immediate, cpu6502.eor, 2, 2),
		CreateInstruction("LSR", Implicit, cpu6502.lsr, 2, 1),
		{},
		CreateInstruction("JMP", Absolute, cpu6502.jmp, 3, 3),
		CreateInstruction("EOR", Absolute, cpu6502.eor, 4, 3),
		CreateInstruction("LSR", Absolute, cpu6502.lsr, 6, 3),
		{},

		// 0x50
		CreateInstruction("BVC", Relative, cpu6502.bvc, 2, 2),
		CreateInstruction("EOR", IndirectY, cpu6502.eor, 5, 2),
		{},
		{},
		{},
		CreateInstruction("EOR", ZeroPageX, cpu6502.eor, 4, 2),
		CreateInstruction("LSR", ZeroPageX, cpu6502.lsr, 6, 2),
		{},
		CreateInstruction("CLI", Implicit, cpu6502.cli, 2, 1),
		CreateInstruction("EOR", AbsoluteYIndexed, cpu6502.eor, 4, 3),
		{},
		{},
		{},
		CreateInstruction("EOR", AbsoluteXIndexed, cpu6502.eor, 4, 3),
		CreateInstruction("LSR", AbsoluteXIndexed, cpu6502.lsr, 7, 3),
		{},

		// 0x60
		CreateInstruction("RTS", Implicit, cpu6502.rts, 6, 1),
		CreateInstruction("ADC", IndirectX, cpu6502.adc, 6, 2),
		{},
		{},
		{},
		CreateInstruction("ADC", ZeroPage, cpu6502.adc, 3, 2),
		CreateInstruction("ROR", ZeroPage, cpu6502.ror, 5, 2),
		{},
		CreateInstruction("PLA", Implicit, cpu6502.pla, 4, 1),
		CreateInstruction("ADC", Immediate, cpu6502.adc, 2, 2),
		CreateInstruction("ROR", Implicit, cpu6502.ror, 2, 1),
		{},
		CreateInstruction("JMP", Indirect, cpu6502.jmp, 5, 3),
		CreateInstruction("ADC", Absolute, cpu6502.adc, 4, 3),
		CreateInstruction("ROR", Absolute, cpu6502.ror, 6, 3),
		{},

		// 0x70
		CreateInstruction("BVS", Relative, cpu6502.bvs, 2, 2),
		CreateInstruction("ADC", IndirectY, cpu6502.adc, 5, 2),
		{},
		{},
		{},
		CreateInstruction("ADC", ZeroPageX, cpu6502.adc, 4, 2),
		CreateInstruction("ROR", ZeroPageX, cpu6502.ror, 6, 2),
		{},
		CreateInstruction("SEI", Implicit, cpu6502.sei, 2, 1),
		CreateInstruction("ADC", AbsoluteYIndexed, cpu6502.adc, 4, 3),
		{},
		{},
		{},
		CreateInstruction("ADC", AbsoluteXIndexed, cpu6502.adc, 4, 3),
		CreateInstruction("ROR", AbsoluteXIndexed, cpu6502.ror, 7, 3),
		{},

		// 0x80
		{},
		CreateInstruction("STA", IndirectX, cpu6502.sta, 6, 2),
		{},
		{},
		CreateInstruction("STY", ZeroPage, cpu6502.sty, 3, 2),
		CreateInstruction("STA", ZeroPage, cpu6502.sta, 3, 2),
		CreateInstruction("STX", ZeroPage, cpu6502.stx, 3, 2),
		{},
		CreateInstruction("DEY", Implicit, cpu6502.dey, 2, 1),
		{},
		CreateInstruction("TXA", Implicit, cpu6502.txa, 2, 1),
		{},
		CreateInstruction("STY", Absolute, cpu6502.sty, 4, 3),
		CreateInstruction("STA", Absolute, cpu6502.sta, 4, 3),
		CreateInstruction("STX", Absolute, cpu6502.stx, 4, 3),
		{},

		// 0x90
		CreateInstruction("BCC", Relative, cpu6502.bcc, 2, 2),
		CreateInstruction("STA", IndirectY, cpu6502.sta, 6, 2),
		{},
		{},
		CreateInstruction("STY", ZeroPageX, cpu6502.sty, 4, 2),
		CreateInstruction("STA", ZeroPageX, cpu6502.sta, 4, 2),
		CreateInstruction("STX", ZeroPageY, cpu6502.stx, 4, 2),
		{},
		CreateInstruction("TYA", Implicit, cpu6502.tya, 2, 1),
		CreateInstruction("STA", AbsoluteYIndexed, cpu6502.sta, 5, 3),
		CreateInstruction("TXS", Implicit, cpu6502.txs, 2, 1),
		{},
		{},
		CreateInstruction("STA", AbsoluteXIndexed, cpu6502.sta, 5, 3),
		{},
		{},

		// 0xA0
		CreateInstruction("LDY", Immediate, cpu6502.ldy, 2, 2),
		CreateInstruction("LDA", IndirectX, cpu6502.lda, 6, 2),
		CreateInstruction("LDX", Immediate, cpu6502.ldx, 2, 2),
		{},
		CreateInstruction("LDY", ZeroPage, cpu6502.ldy, 3, 2),
		CreateInstruction("LDA", ZeroPage, cpu6502.lda, 3, 2),
		CreateInstruction("LDX", ZeroPage, cpu6502.ldx, 3, 2),
		{},
		CreateInstruction("TAY", Implicit, cpu6502.tay, 2, 1),
		CreateInstruction("LDA", Immediate, cpu6502.lda, 2, 2),
		CreateInstruction("TAX", Implicit, cpu6502.tax, 2, 1),
		{},
		CreateInstruction("LDY", Absolute, cpu6502.ldy, 4, 3),
		CreateInstruction("LDA", Absolute, cpu6502.lda, 4, 3),
		CreateInstruction("LDX", Absolute, cpu6502.ldx, 4, 3),
		{},

		// 0xB0
		CreateInstruction("BCS", Relative, cpu6502.bcs, 2, 2),
		CreateInstruction("LDA", IndirectY, cpu6502.lda, 5, 2),
		{},
		{},
		CreateInstruction("LDY", ZeroPageX, cpu6502.ldy, 4, 2),
		CreateInstruction("LDA", ZeroPageX, cpu6502.lda, 4, 2),
		CreateInstruction("LDX", ZeroPageY, cpu6502.ldx, 4, 2),
		{},
		CreateInstruction("CLV", Implicit, cpu6502.clv, 2, 1),
		CreateInstruction("LDA", AbsoluteYIndexed, cpu6502.lda, 4, 3),
		CreateInstruction("TSX", Implicit, cpu6502.tsx, 2, 1),
		{},
		CreateInstruction("LDY", AbsoluteXIndexed, cpu6502.ldy, 4, 3),
		CreateInstruction("LDA", AbsoluteXIndexed, cpu6502.lda, 4, 3),
		CreateInstruction("LDX", AbsoluteYIndexed, cpu6502.ldx, 4, 3),
		{},

		// 0xC0
		CreateInstruction("CPY", Immediate, cpu6502.cpy, 2, 2),
		CreateInstruction("CMP", IndirectX, cpu6502.cmp, 6, 2),
		{},
		{},
		CreateInstruction("CPY", ZeroPage, cpu6502.cpy, 3, 2),
		CreateInstruction("CMP", ZeroPage, cpu6502.cmp, 3, 2),
		CreateInstruction("DEC", ZeroPage, cpu6502.dec, 5, 2),
		{},
		CreateInstruction("INY", Implicit, cpu6502.iny, 2, 1),
		CreateInstruction("CMP", Immediate, cpu6502.cmp, 2, 2),
		CreateInstruction("DEX", Implicit, cpu6502.dex, 2, 1),
		{},
		CreateInstruction("CPY", Absolute, cpu6502.cpy, 4, 3),
		CreateInstruction("CMP", Absolute, cpu6502.cmp, 4, 3),
		CreateInstruction("DEC", Absolute, cpu6502.dec, 6, 3),
		{},

		// 0xD0
		CreateInstruction("BNE", Relative, cpu6502.bne, 2, 2),
		CreateInstruction("CMP", IndirectY, cpu6502.cmp, 5, 2),
		{},
		{},
		{},
		CreateInstruction("CMP", ZeroPageX, cpu6502.cmp, 4, 2),
		CreateInstruction("DEC", ZeroPageX, cpu6502.dec, 6, 2),
		{},
		CreateInstruction("CLD", Implicit, cpu6502.cld, 2, 1),
		CreateInstruction("CMP", AbsoluteYIndexed, cpu6502.cmp, 4, 3),
		{},
		{},
		{},
		CreateInstruction("CMP", AbsoluteXIndexed, cpu6502.cmp, 4, 3),
		CreateInstruction("DEC", AbsoluteXIndexed, cpu6502.dec, 7, 3),
		{},

		// 0xE0
		CreateInstruction("CPX", Immediate, cpu6502.cpx, 2, 2),
		CreateInstruction("SBC", IndirectX, cpu6502.sbc, 6, 2),
		{},
		{},
		CreateInstruction("CPX", ZeroPage, cpu6502.cpx, 3, 2),
		CreateInstruction("SBC", ZeroPage, cpu6502.sbc, 3, 2),
		CreateInstruction("INC", ZeroPage, cpu6502.inc, 5, 2),
		{},
		CreateInstruction("INX", Implicit, cpu6502.inx, 2, 1),
		CreateInstruction("SBC", Immediate, cpu6502.sbc, 2, 2),
		CreateInstruction("NOP", Implicit, cpu6502.nop, 2, 1),
		{},
		CreateInstruction("CPX", Absolute, cpu6502.cpx, 4, 3),
		CreateInstruction("SBC", Absolute, cpu6502.sbc, 4, 3),
		CreateInstruction("INC", Absolute, cpu6502.inc, 6, 3),
		{},

		// 0xF0
		CreateInstruction("BEQ", Relative, cpu6502.beq, 2, 2),
		CreateInstruction("SBC", IndirectY, cpu6502.sbc, 5, 2),
		{},
		{},
		{},
		CreateInstruction("SBC", ZeroPageX, cpu6502.sbc, 4, 2),
		CreateInstruction("INC", ZeroPageX, cpu6502.inc, 6, 2),
		{},
		CreateInstruction("SED", Implicit, cpu6502.sed, 2, 1),
		CreateInstruction("SBC", AbsoluteYIndexed, cpu6502.sbc, 4, 3),
		{},
		{},
		{},
		CreateInstruction("SBC", AbsoluteXIndexed, cpu6502.sbc, 4, 3),
		CreateInstruction("INC", AbsoluteXIndexed, cpu6502.inc, 7, 3),
		{},
	}
}

func (cpu6502 *Cpu6502) initAddressModeEvaluators() {
	cpu6502.addressEvaluators = [13]AddressModeMethod{
		cpu6502.evalImplicit,
		cpu6502.evalImmediate,
		cpu6502.evalZeroPage,
		cpu6502.evalZeroPageX,
		cpu6502.evalZeroPageY,
		cpu6502.evalAbsolute,
		cpu6502.evalAbsoluteXIndexed,
		cpu6502.evalAbsoluteYIndexed,
		cpu6502.evalIndirect,
		cpu6502.evalIndirectX,
		cpu6502.evalIndirectY,
		cpu6502.evalRelative,
	}
}

func (cpu6502 Cpu6502) evaluateOperandAddress(addressMode AddressMode, pc types.Address) (finalAddress types.Address, operand [3]byte, pageCrossed bool) {
	if addressMode == Implicit {
		finalAddress = 0
		return
	}

	if cpu6502.addressEvaluators[addressMode] == nil {
		msg := fmt.Errorf("cannot find address evaluator for address mode \"%d\"", addressMode)
		cpu6502.Stop()
		panic(msg)
	}

	finalAddress, operand, _, pageCrossed = cpu6502.addressEvaluators[addressMode](pc)

	return
}

func memoryPageDiffer(address types.Address, finalAddress types.Address) bool {
	return address&0xFF00 != finalAddress&0xFF00
}
