package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type AddressModeMethod func(programCounter types.Address) (finalAddress types.Address, operand [3]byte, cycles int, pageCrossed bool)

// Cpu6502 Represents a CPU 6502
type Cpu6502 struct {
	registers cpu.Registers
	memory    Memory

	instructions [256]cpu.Instruction
	opCyclesLeft byte // How many cycles left to finish execution of current cycle
	cycle        uint32

	addressEvaluators [13]AddressModeMethod

	// Debug parameters
	debugger    *cpu.Debugger
	cyclesLimit uint16
}

func CreateCPU(memory Memory, debugger *cpu.Debugger) *Cpu6502 {
	cpu6502 := Cpu6502{
		memory:    memory,
		registers: cpu.CreateRegisters(),
		debugger:  debugger,
	}

	cpu6502.initInstructionsTable()
	cpu6502.initAddressModeEvaluators()

	return &cpu6502
}

func (cpu6502 *Cpu6502) ProgramCounter() types.Address {
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
	cpu6502.instructions = [256]cpu.Instruction{
		cpu.CreateInstruction("BRK", cpu.Implicit, cpu6502.brk, 7, 1),
		cpu.CreateInstruction("ORA", cpu.IndirectX, cpu6502.ora, 6, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("ORA", cpu.ZeroPage, cpu6502.ora, 3, 2),
		cpu.CreateInstruction("ASL", cpu.ZeroPage, cpu6502.asl, 5, 2),
		{},
		cpu.CreateInstruction("PHP", cpu.Implicit, cpu6502.php, 3, 1),
		cpu.CreateInstruction("ORA", cpu.Immediate, cpu6502.ora, 2, 2),
		cpu.CreateInstruction("ASL", cpu.Implicit, cpu6502.asl, 2, 1),
		{},
		{},
		cpu.CreateInstruction("ORA", cpu.Absolute, cpu6502.ora, 4, 3),
		cpu.CreateInstruction("ASL", cpu.Absolute, cpu6502.asl, 6, 3),
		{},

		// 0x10
		cpu.CreateInstruction("BPL", cpu.Relative, cpu6502.bpl, 2, 2),
		cpu.CreateInstruction("ORA", cpu.IndirectY, cpu6502.ora, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("ORA", cpu.ZeroPageX, cpu6502.ora, 4, 2),
		cpu.CreateInstruction("ASL", cpu.ZeroPageX, cpu6502.asl, 6, 2),
		{},
		cpu.CreateInstruction("CLC", cpu.Implicit, cpu6502.clc, 2, 1),
		cpu.CreateInstruction("ORA", cpu.AbsoluteYIndexed, cpu6502.ora, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("ORA", cpu.AbsoluteXIndexed, cpu6502.ora, 4, 3),
		cpu.CreateInstruction("ASL", cpu.AbsoluteXIndexed, cpu6502.asl, 7, 3),
		{},

		// 0x20
		cpu.CreateInstruction("JSR", cpu.Absolute, cpu6502.jsr, 6, 3),
		cpu.CreateInstruction("AND", cpu.IndirectX, cpu6502.and, 6, 2),
		{},
		{},
		cpu.CreateInstruction("BIT", cpu.ZeroPage, cpu6502.bit, 3, 2),
		cpu.CreateInstruction("AND", cpu.ZeroPage, cpu6502.and, 3, 2),
		cpu.CreateInstruction("ROL", cpu.ZeroPage, cpu6502.rol, 5, 2),
		{},
		cpu.CreateInstruction("PLP", cpu.Implicit, cpu6502.plp, 4, 1),
		cpu.CreateInstruction("AND", cpu.Immediate, cpu6502.and, 2, 2),
		cpu.CreateInstruction("ROL", cpu.Implicit, cpu6502.rol, 2, 1),
		{},
		cpu.CreateInstruction("BIT", cpu.Absolute, cpu6502.bit, 4, 3),
		cpu.CreateInstruction("AND", cpu.Absolute, cpu6502.and, 4, 3),
		cpu.CreateInstruction("ROL", cpu.Absolute, cpu6502.rol, 6, 3),
		{},

		// 0x30
		cpu.CreateInstruction("BMI", cpu.Relative, cpu6502.bmi, 2, 2),
		cpu.CreateInstruction("AND", cpu.IndirectY, cpu6502.and, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("AND", cpu.ZeroPageX, cpu6502.and, 4, 2),
		cpu.CreateInstruction("ROL", cpu.ZeroPageX, cpu6502.rol, 6, 2),
		{},
		cpu.CreateInstruction("SEC", cpu.Implicit, cpu6502.sec, 2, 1),
		cpu.CreateInstruction("AND", cpu.AbsoluteYIndexed, cpu6502.and, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("AND", cpu.AbsoluteXIndexed, cpu6502.and, 4, 3),
		cpu.CreateInstruction("ROL", cpu.AbsoluteXIndexed, cpu6502.rol, 7, 3),
		{},

		// 0x40
		cpu.CreateInstruction("RTI", cpu.Implicit, cpu6502.rti, 6, 1),
		cpu.CreateInstruction("EOR", cpu.IndirectX, cpu6502.eor, 6, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("EOR", cpu.ZeroPage, cpu6502.eor, 3, 2),
		cpu.CreateInstruction("LSR", cpu.ZeroPage, cpu6502.lsr, 5, 2),
		{},
		cpu.CreateInstruction("PHA", cpu.Implicit, cpu6502.pha, 3, 1),
		cpu.CreateInstruction("EOR", cpu.Immediate, cpu6502.eor, 2, 2),
		cpu.CreateInstruction("LSR", cpu.Implicit, cpu6502.lsr, 2, 1),
		{},
		cpu.CreateInstruction("JMP", cpu.Absolute, cpu6502.jmp, 3, 3),
		cpu.CreateInstruction("EOR", cpu.Absolute, cpu6502.eor, 4, 3),
		cpu.CreateInstruction("LSR", cpu.Absolute, cpu6502.lsr, 6, 3),
		{},

		// 0x50
		cpu.CreateInstruction("BVC", cpu.Relative, cpu6502.bvc, 2, 2),
		cpu.CreateInstruction("EOR", cpu.IndirectY, cpu6502.eor, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("EOR", cpu.ZeroPageX, cpu6502.eor, 4, 2),
		cpu.CreateInstruction("LSR", cpu.ZeroPageX, cpu6502.lsr, 6, 2),
		{},
		cpu.CreateInstruction("CLI", cpu.Implicit, cpu6502.cli, 2, 1),
		cpu.CreateInstruction("EOR", cpu.AbsoluteYIndexed, cpu6502.eor, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("EOR", cpu.AbsoluteXIndexed, cpu6502.eor, 4, 3),
		cpu.CreateInstruction("LSR", cpu.AbsoluteXIndexed, cpu6502.lsr, 7, 3),
		{},

		// 0x60
		cpu.CreateInstruction("RTS", cpu.Implicit, cpu6502.rts, 6, 1),
		cpu.CreateInstruction("ADC", cpu.IndirectX, cpu6502.adc, 6, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("ADC", cpu.ZeroPage, cpu6502.adc, 3, 2),
		cpu.CreateInstruction("ROR", cpu.ZeroPage, cpu6502.ror, 5, 2),
		{},
		cpu.CreateInstruction("PLA", cpu.Implicit, cpu6502.pla, 4, 1),
		cpu.CreateInstruction("ADC", cpu.Immediate, cpu6502.adc, 2, 2),
		cpu.CreateInstruction("ROR", cpu.Implicit, cpu6502.ror, 2, 1),
		{},
		cpu.CreateInstruction("JMP", cpu.Indirect, cpu6502.jmp, 5, 3),
		cpu.CreateInstruction("ADC", cpu.Absolute, cpu6502.adc, 4, 3),
		cpu.CreateInstruction("ROR", cpu.Absolute, cpu6502.ror, 6, 3),
		{},

		// 0x70
		cpu.CreateInstruction("BVS", cpu.Relative, cpu6502.bvs, 2, 2),
		cpu.CreateInstruction("ADC", cpu.IndirectY, cpu6502.adc, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("ADC", cpu.ZeroPageX, cpu6502.adc, 4, 2),
		cpu.CreateInstruction("ROR", cpu.ZeroPageX, cpu6502.ror, 6, 2),
		{},
		cpu.CreateInstruction("SEI", cpu.Implicit, cpu6502.sei, 2, 1),
		cpu.CreateInstruction("ADC", cpu.AbsoluteYIndexed, cpu6502.adc, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("ADC", cpu.AbsoluteXIndexed, cpu6502.adc, 4, 3),
		cpu.CreateInstruction("ROR", cpu.AbsoluteXIndexed, cpu6502.ror, 7, 3),
		{},

		// 0x80
		{},
		cpu.CreateInstruction("STA", cpu.IndirectX, cpu6502.sta, 6, 2),
		{},
		{},
		cpu.CreateInstruction("STY", cpu.ZeroPage, cpu6502.sty, 3, 2),
		cpu.CreateInstruction("STA", cpu.ZeroPage, cpu6502.sta, 3, 2),
		cpu.CreateInstruction("STX", cpu.ZeroPage, cpu6502.stx, 3, 2),
		{},
		cpu.CreateInstruction("DEY", cpu.Implicit, cpu6502.dey, 2, 1),
		{},
		cpu.CreateInstruction("TXA", cpu.Implicit, cpu6502.txa, 2, 1),
		{},
		cpu.CreateInstruction("STY", cpu.Absolute, cpu6502.sty, 4, 3),
		cpu.CreateInstruction("STA", cpu.Absolute, cpu6502.sta, 4, 3),
		cpu.CreateInstruction("STX", cpu.Absolute, cpu6502.stx, 4, 3),
		{},

		// 0x90
		cpu.CreateInstruction("BCC", cpu.Relative, cpu6502.bcc, 2, 2),
		cpu.CreateInstruction("STA", cpu.IndirectY, cpu6502.sta, 6, 2),
		{},
		{},
		cpu.CreateInstruction("STY", cpu.ZeroPageX, cpu6502.sty, 4, 2),
		cpu.CreateInstruction("STA", cpu.ZeroPageX, cpu6502.sta, 4, 2),
		cpu.CreateInstruction("STX", cpu.ZeroPageY, cpu6502.stx, 4, 2),
		{},
		cpu.CreateInstruction("TYA", cpu.Implicit, cpu6502.tya, 2, 1),
		cpu.CreateInstruction("STA", cpu.AbsoluteYIndexed, cpu6502.sta, 5, 3),
		cpu.CreateInstruction("TXS", cpu.Implicit, cpu6502.txs, 2, 1),
		{},
		{},
		cpu.CreateInstruction("STA", cpu.AbsoluteXIndexed, cpu6502.sta, 5, 3),
		{},
		{},

		// 0xA0
		cpu.CreateInstruction("LDY", cpu.Immediate, cpu6502.ldy, 2, 2),
		cpu.CreateInstruction("LDA", cpu.IndirectX, cpu6502.lda, 6, 2),
		cpu.CreateInstruction("LDX", cpu.Immediate, cpu6502.ldx, 2, 2),
		{},
		cpu.CreateInstruction("LDY", cpu.ZeroPage, cpu6502.ldy, 3, 2),
		cpu.CreateInstruction("LDA", cpu.ZeroPage, cpu6502.lda, 3, 2),
		cpu.CreateInstruction("LDX", cpu.ZeroPage, cpu6502.ldx, 3, 2),
		{},
		cpu.CreateInstruction("TAY", cpu.Implicit, cpu6502.tay, 2, 1),
		cpu.CreateInstruction("LDA", cpu.Immediate, cpu6502.lda, 2, 2),
		cpu.CreateInstruction("TAX", cpu.Implicit, cpu6502.tax, 2, 1),
		{},
		cpu.CreateInstruction("LDY", cpu.Absolute, cpu6502.ldy, 4, 3),
		cpu.CreateInstruction("LDA", cpu.Absolute, cpu6502.lda, 4, 3),
		cpu.CreateInstruction("LDX", cpu.Absolute, cpu6502.ldx, 4, 3),
		{},

		// 0xB0
		cpu.CreateInstruction("BCS", cpu.Relative, cpu6502.bcs, 2, 2),
		cpu.CreateInstruction("LDA", cpu.IndirectY, cpu6502.lda, 5, 2),
		{},
		{},
		cpu.CreateInstruction("LDY", cpu.ZeroPageX, cpu6502.ldy, 4, 2),
		cpu.CreateInstruction("LDA", cpu.ZeroPageX, cpu6502.lda, 4, 2),
		cpu.CreateInstruction("LDX", cpu.ZeroPageY, cpu6502.ldx, 4, 2),
		{},
		cpu.CreateInstruction("CLV", cpu.Implicit, cpu6502.clv, 2, 1),
		cpu.CreateInstruction("LDA", cpu.AbsoluteYIndexed, cpu6502.lda, 4, 3),
		cpu.CreateInstruction("TSX", cpu.Implicit, cpu6502.tsx, 2, 1),
		{},
		cpu.CreateInstruction("LDY", cpu.AbsoluteXIndexed, cpu6502.ldy, 4, 3),
		cpu.CreateInstruction("LDA", cpu.AbsoluteXIndexed, cpu6502.lda, 4, 3),
		cpu.CreateInstruction("LDX", cpu.AbsoluteYIndexed, cpu6502.ldx, 4, 3),
		{},

		// 0xC0
		cpu.CreateInstruction("CPY", cpu.Immediate, cpu6502.cpy, 2, 2),
		cpu.CreateInstruction("CMP", cpu.IndirectX, cpu6502.cmp, 6, 2),
		{},
		{},
		cpu.CreateInstruction("CPY", cpu.ZeroPage, cpu6502.cpy, 3, 2),
		cpu.CreateInstruction("CMP", cpu.ZeroPage, cpu6502.cmp, 3, 2),
		cpu.CreateInstruction("DEC", cpu.ZeroPage, cpu6502.dec, 5, 2),
		{},
		cpu.CreateInstruction("INY", cpu.Implicit, cpu6502.iny, 2, 1),
		cpu.CreateInstruction("CMP", cpu.Immediate, cpu6502.cmp, 2, 2),
		cpu.CreateInstruction("DEX", cpu.Implicit, cpu6502.dex, 2, 1),
		{},
		cpu.CreateInstruction("CPY", cpu.Absolute, cpu6502.cpy, 4, 3),
		cpu.CreateInstruction("CMP", cpu.Absolute, cpu6502.cmp, 4, 3),
		cpu.CreateInstruction("DEC", cpu.Absolute, cpu6502.dec, 6, 3),
		{},

		// 0xD0
		cpu.CreateInstruction("BNE", cpu.Relative, cpu6502.bne, 2, 2),
		cpu.CreateInstruction("CMP", cpu.IndirectY, cpu6502.cmp, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("CMP", cpu.ZeroPageX, cpu6502.cmp, 4, 2),
		cpu.CreateInstruction("DEC", cpu.ZeroPageX, cpu6502.dec, 6, 2),
		{},
		cpu.CreateInstruction("CLD", cpu.Implicit, cpu6502.cld, 2, 1),
		cpu.CreateInstruction("CMP", cpu.AbsoluteYIndexed, cpu6502.cmp, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("CMP", cpu.AbsoluteXIndexed, cpu6502.cmp, 4, 3),
		cpu.CreateInstruction("DEC", cpu.AbsoluteXIndexed, cpu6502.dec, 7, 3),
		{},

		// 0xE0
		cpu.CreateInstruction("CPX", cpu.Immediate, cpu6502.cpx, 2, 2),
		cpu.CreateInstruction("SBC", cpu.IndirectX, cpu6502.sbc, 6, 2),
		{},
		{},
		cpu.CreateInstruction("CPX", cpu.ZeroPage, cpu6502.cpx, 3, 2),
		cpu.CreateInstruction("SBC", cpu.ZeroPage, cpu6502.sbc, 3, 2),
		cpu.CreateInstruction("INC", cpu.ZeroPage, cpu6502.inc, 5, 2),
		{},
		cpu.CreateInstruction("INX", cpu.Implicit, cpu6502.inx, 2, 1),
		cpu.CreateInstruction("SBC", cpu.Immediate, cpu6502.sbc, 2, 2),
		cpu.CreateInstruction("NOP", cpu.Implicit, cpu6502.nop, 2, 1),
		{},
		cpu.CreateInstruction("CPX", cpu.Absolute, cpu6502.cpx, 4, 3),
		cpu.CreateInstruction("SBC", cpu.Absolute, cpu6502.sbc, 4, 3),
		cpu.CreateInstruction("INC", cpu.Absolute, cpu6502.inc, 6, 3),
		{},

		// 0xF0
		cpu.CreateInstruction("BEQ", cpu.Relative, cpu6502.beq, 2, 2),
		cpu.CreateInstruction("SBC", cpu.IndirectY, cpu6502.sbc, 5, 2),
		{},
		{},
		{},
		cpu.CreateInstruction("SBC", cpu.ZeroPageX, cpu6502.sbc, 4, 2),
		cpu.CreateInstruction("INC", cpu.ZeroPageX, cpu6502.inc, 6, 2),
		{},
		cpu.CreateInstruction("SED", cpu.Implicit, cpu6502.sed, 2, 1),
		cpu.CreateInstruction("SBC", cpu.AbsoluteYIndexed, cpu6502.sbc, 4, 3),
		{},
		{},
		{},
		cpu.CreateInstruction("SBC", cpu.AbsoluteXIndexed, cpu6502.sbc, 4, 3),
		cpu.CreateInstruction("INC", cpu.AbsoluteXIndexed, cpu6502.inc, 7, 3),
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

func (cpu6502 *Cpu6502) evaluateOperandAddress(addressMode cpu.AddressMode, pc types.Address) (finalAddress types.Address, operand [3]byte, pageCrossed bool) {
	if addressMode == cpu.Implicit {
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
