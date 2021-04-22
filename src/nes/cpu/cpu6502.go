package cpu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/component"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

// Cpu6502 Represents a NES cpu
type Cpu6502 struct {
	registers Cpu6502Registers
	bus       *component.Bus

	instructions     [256]defs.Instruction
	instructionCycle byte
	Cycle            uint16

	addressEvaluators [13]func(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool)
	// Debug parameters
	debug       bool
	Logger      Logger
	cyclesLimit uint16
}

func CreateCPU(bus *component.Bus) Cpu6502 {
	registers := CreateRegisters()
	cpu := Cpu6502{
		registers: registers,
		bus:       bus,
		debug:     false,
	}

	cpu.Init()

	return cpu
}

func CreateCPUDebuggable(bus *component.Bus, logger *Logger) Cpu6502 {
	registers := CreateRegisters()
	cpu := Cpu6502{
		registers: registers,
		bus:       bus,
		debug:     true,
		Logger:    *logger,
	}

	cpu.Init()

	return cpu
}

// CreateCPUWithBus creates a Cpu6502 with a Bus, Useful for tests
func CreateCPUWithBus() Cpu6502 {
	registers := CreateRegisters()

	ram := component.RAM{}
	gamepak := component.CreateDummyGamePak()

	bus := component.CreateBus(&ram)
	bus.InsertGamePak(&gamepak)

	cpu := Cpu6502{
		registers: registers,
		bus:       &bus,
		debug:     false,
	}

	cpu.Init()

	return cpu
}

func (cpu *Cpu6502) Registers() *Cpu6502Registers {
	return &cpu.registers
}

func (cpu *Cpu6502) pushStack(value byte) {
	address := cpu.registers.stackPointerAddress()
	cpu.bus.Write(
		address,
		value,
	)

	cpu.registers.stackPointerPushed()
}

func (cpu *Cpu6502) popStack() byte {
	cpu.registers.stackPointerPopped()
	address := cpu.registers.stackPointerAddress()
	return cpu.bus.Read(address)
}

// Reads value located at Program Counter and increments it
func (cpu *Cpu6502) fetch() byte {
	value := cpu.bus.Read(cpu.registers.Pc)
	cpu.registers.Pc++

	return value
}

func (cpu *Cpu6502) Read(address defs.Address) byte {
	return cpu.bus.Read(address)
}

func (cpu *Cpu6502) read16(address defs.Address) defs.Word {
	low := cpu.bus.Read(address)
	high := cpu.bus.Read(address + 1)

	return defs.CreateWord(low, high)
}

func (cpu *Cpu6502) write(address defs.Address, value byte) {
	cpu.bus.Write(address, value)
}

func (cpu *Cpu6502) initInstructionsTable() {
	cpu.instructions = [256]defs.Instruction{
		defs.CreateInstruction("BRK", defs.Implicit, cpu.brk, 7, 1),
		defs.CreateInstruction("ORA", defs.IndirectX, cpu.ora, 6, 2),
		{},
		{},
		{},
		defs.CreateInstruction("ORA", defs.ZeroPage, cpu.ora, 3, 2),
		defs.CreateInstruction("ASL", defs.ZeroPage, cpu.asl, 5, 2),
		{},
		defs.CreateInstruction("PHP", defs.Implicit, cpu.php, 3, 1),
		defs.CreateInstruction("ORA", defs.Immediate, cpu.ora, 2, 2),
		defs.CreateInstruction("ASL", defs.Implicit, cpu.asl, 2, 1),
		{},
		{},
		defs.CreateInstruction("ORA", defs.Absolute, cpu.ora, 4, 3),
		defs.CreateInstruction("ASL", defs.Absolute, cpu.asl, 6, 3),
		{},

		// 0x10
		defs.CreateInstruction("BPL", defs.Relative, cpu.bpl, 2, 2),
		defs.CreateInstruction("ORA", defs.IndirectY, cpu.ora, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("ORA", defs.ZeroPageX, cpu.ora, 4, 2),
		defs.CreateInstruction("ASL", defs.ZeroPageX, cpu.asl, 6, 2),
		{},
		defs.CreateInstruction("CLC", defs.Implicit, cpu.clc, 2, 1),
		defs.CreateInstruction("ORA", defs.AbsoluteYIndexed, cpu.ora, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("ORA", defs.AbsoluteXIndexed, cpu.ora, 4, 3),
		defs.CreateInstruction("ASL", defs.AbsoluteXIndexed, cpu.asl, 7, 3),
		{},

		// 0x20
		defs.CreateInstruction("JSR", defs.Absolute, cpu.jsr, 6, 3),
		defs.CreateInstruction("AND", defs.IndirectX, cpu.and, 6, 2),
		{},
		{},
		defs.CreateInstruction("BIT", defs.ZeroPage, cpu.bit, 3, 2),
		defs.CreateInstruction("AND", defs.ZeroPage, cpu.and, 3, 2),
		defs.CreateInstruction("ROL", defs.ZeroPage, cpu.rol, 5, 2),
		{},
		defs.CreateInstruction("PLP", defs.Implicit, cpu.plp, 4, 1),
		defs.CreateInstruction("AND", defs.Immediate, cpu.and, 2, 2),
		defs.CreateInstruction("ROL", defs.Implicit, cpu.rol, 2, 1),
		{},
		defs.CreateInstruction("BIT", defs.Absolute, cpu.bit, 4, 3),
		defs.CreateInstruction("AND", defs.Absolute, cpu.and, 4, 3),
		defs.CreateInstruction("ROL", defs.Absolute, cpu.rol, 6, 3),
		{},

		// 0x30
		defs.CreateInstruction("BMI", defs.Relative, cpu.bmi, 2, 2),
		defs.CreateInstruction("AND", defs.IndirectY, cpu.and, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("AND", defs.ZeroPageX, cpu.and, 4, 2),
		defs.CreateInstruction("ROL", defs.ZeroPageX, cpu.rol, 6, 2),
		{},
		defs.CreateInstruction("SEC", defs.Implicit, cpu.sec, 2, 1),
		defs.CreateInstruction("AND", defs.AbsoluteYIndexed, cpu.and, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("AND", defs.AbsoluteXIndexed, cpu.and, 4, 3),
		defs.CreateInstruction("ROL", defs.AbsoluteXIndexed, cpu.rol, 7, 3),
		{},

		// 0x40
		defs.CreateInstruction("RTI", defs.Implicit, cpu.rti, 6, 1),
		defs.CreateInstruction("EOR", defs.IndirectX, cpu.eor, 6, 2),
		{},
		{},
		{},
		defs.CreateInstruction("EOR", defs.ZeroPage, cpu.eor, 3, 2),
		defs.CreateInstruction("LSR", defs.ZeroPage, cpu.lsr, 5, 2),
		{},
		defs.CreateInstruction("PHA", defs.Implicit, cpu.pha, 3, 1),
		defs.CreateInstruction("EOR", defs.Immediate, cpu.eor, 2, 2),
		defs.CreateInstruction("LSR", defs.Implicit, cpu.lsr, 2, 1),
		{},
		defs.CreateInstruction("JMP", defs.Absolute, cpu.jmp, 3, 3),
		defs.CreateInstruction("EOR", defs.Absolute, cpu.eor, 4, 3),
		defs.CreateInstruction("LSR", defs.Absolute, cpu.lsr, 6, 3),
		{},

		// 0x50
		defs.CreateInstruction("BVC", defs.Relative, cpu.bvc, 2, 2),
		defs.CreateInstruction("EOR", defs.IndirectY, cpu.eor, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("EOR", defs.ZeroPageX, cpu.eor, 4, 2),
		defs.CreateInstruction("LSR", defs.ZeroPageX, cpu.lsr, 6, 2),
		{},
		defs.CreateInstruction("CLI", defs.Implicit, cpu.cli, 2, 1),
		defs.CreateInstruction("EOR", defs.AbsoluteYIndexed, cpu.eor, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("EOR", defs.AbsoluteXIndexed, cpu.eor, 4, 3),
		defs.CreateInstruction("LSR", defs.AbsoluteXIndexed, cpu.lsr, 7, 3),
		{},

		// 0x60
		defs.CreateInstruction("RTS", defs.Implicit, cpu.rts, 6, 1),
		defs.CreateInstruction("ADC", defs.IndirectX, cpu.adc, 6, 2),
		{},
		{},
		{},
		defs.CreateInstruction("ADC", defs.ZeroPage, cpu.adc, 3, 2),
		defs.CreateInstruction("ROR", defs.ZeroPage, cpu.ror, 5, 2),
		{},
		defs.CreateInstruction("PLA", defs.Implicit, cpu.pla, 4, 1),
		defs.CreateInstruction("ADC", defs.Immediate, cpu.adc, 2, 2),
		defs.CreateInstruction("ROR", defs.Implicit, cpu.ror, 2, 1),
		{},
		defs.CreateInstruction("JMP", defs.Indirect, cpu.jmp, 5, 3),
		defs.CreateInstruction("ADC", defs.Absolute, cpu.adc, 4, 3),
		defs.CreateInstruction("ROR", defs.Absolute, cpu.ror, 6, 3),
		{},

		// 0x70
		defs.CreateInstruction("BVS", defs.Relative, cpu.bvs, 2, 2),
		defs.CreateInstruction("ADC", defs.IndirectY, cpu.adc, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("ADC", defs.ZeroPageX, cpu.adc, 4, 2),
		defs.CreateInstruction("ROR", defs.ZeroPageX, cpu.ror, 6, 2),
		{},
		defs.CreateInstruction("SEI", defs.Implicit, cpu.sei, 2, 1),
		defs.CreateInstruction("ADC", defs.AbsoluteYIndexed, cpu.adc, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("ADC", defs.AbsoluteXIndexed, cpu.adc, 4, 3),
		defs.CreateInstruction("ROR", defs.AbsoluteXIndexed, cpu.ror, 7, 3),
		{},

		// 0x80
		{},
		defs.CreateInstruction("STA", defs.IndirectX, cpu.sta, 6, 2),
		{},
		{},
		defs.CreateInstruction("STY", defs.ZeroPage, cpu.sty, 3, 2),
		defs.CreateInstruction("STA", defs.ZeroPage, cpu.sta, 3, 2),
		defs.CreateInstruction("STX", defs.ZeroPage, cpu.stx, 3, 2),
		{},
		defs.CreateInstruction("DEY", defs.Implicit, cpu.dey, 2, 1),
		{},
		defs.CreateInstruction("TXA", defs.Implicit, cpu.txa, 2, 1),
		{},
		defs.CreateInstruction("STY", defs.Absolute, cpu.sty, 4, 3),
		defs.CreateInstruction("STA", defs.Absolute, cpu.sta, 4, 3),
		defs.CreateInstruction("STX", defs.Absolute, cpu.stx, 4, 3),
		{},

		// 0x90
		defs.CreateInstruction("BCC", defs.Relative, cpu.bcc, 2, 2),
		defs.CreateInstruction("STA", defs.IndirectY, cpu.sta, 6, 2),
		{},
		{},
		defs.CreateInstruction("STY", defs.ZeroPageX, cpu.sty, 4, 2),
		defs.CreateInstruction("STA", defs.ZeroPageX, cpu.sta, 4, 2),
		defs.CreateInstruction("STX", defs.ZeroPageY, cpu.stx, 4, 2),
		{},
		defs.CreateInstruction("TYA", defs.Implicit, cpu.tya, 2, 1),
		defs.CreateInstruction("STA", defs.AbsoluteYIndexed, cpu.sta, 5, 3),
		defs.CreateInstruction("TXS", defs.Implicit, cpu.txs, 2, 1),
		{},
		{},
		defs.CreateInstruction("STA", defs.AbsoluteXIndexed, cpu.sta, 5, 3),
		{},
		{},

		// 0xA0
		defs.CreateInstruction("LDY", defs.Immediate, cpu.ldy, 2, 2),
		defs.CreateInstruction("LDA", defs.IndirectX, cpu.lda, 6, 2),
		defs.CreateInstruction("LDX", defs.Immediate, cpu.ldx, 2, 2),
		{},
		defs.CreateInstruction("LDY", defs.ZeroPage, cpu.ldy, 3, 2),
		defs.CreateInstruction("LDA", defs.ZeroPage, cpu.lda, 3, 2),
		defs.CreateInstruction("LDX", defs.ZeroPage, cpu.ldx, 3, 2),
		{},
		defs.CreateInstruction("TAY", defs.Implicit, cpu.tay, 2, 1),
		defs.CreateInstruction("LDA", defs.Immediate, cpu.lda, 2, 2),
		defs.CreateInstruction("TAX", defs.Implicit, cpu.tax, 2, 1),
		{},
		defs.CreateInstruction("LDY", defs.Absolute, cpu.ldy, 4, 3),
		defs.CreateInstruction("LDA", defs.Absolute, cpu.lda, 4, 3),
		defs.CreateInstruction("LDX", defs.Absolute, cpu.ldx, 4, 3),
		{},

		// 0xB0
		defs.CreateInstruction("BCS", defs.Relative, cpu.bcs, 2, 2),
		defs.CreateInstruction("LDA", defs.IndirectY, cpu.lda, 5, 2),
		{},
		{},
		defs.CreateInstruction("LDY", defs.ZeroPageX, cpu.ldy, 4, 2),
		defs.CreateInstruction("LDA", defs.ZeroPageX, cpu.lda, 4, 2),
		defs.CreateInstruction("LDX", defs.ZeroPageY, cpu.ldx, 4, 2),
		{},
		defs.CreateInstruction("CLV", defs.Implicit, cpu.clv, 2, 1),
		defs.CreateInstruction("LDA", defs.AbsoluteYIndexed, cpu.lda, 4, 3),
		defs.CreateInstruction("TSX", defs.Implicit, cpu.tsx, 2, 1),
		{},
		defs.CreateInstruction("LDY", defs.AbsoluteXIndexed, cpu.ldy, 4, 3),
		defs.CreateInstruction("LDA", defs.AbsoluteXIndexed, cpu.lda, 4, 3),
		defs.CreateInstruction("LDX", defs.AbsoluteYIndexed, cpu.ldx, 4, 3),
		{},

		// 0xC0
		defs.CreateInstruction("CPY", defs.Immediate, cpu.cpy, 2, 2),
		defs.CreateInstruction("CMP", defs.IndirectX, cpu.cmp, 6, 2),
		{},
		{},
		defs.CreateInstruction("CPY", defs.ZeroPage, cpu.cpy, 3, 2),
		defs.CreateInstruction("CMP", defs.ZeroPage, cpu.cmp, 3, 2),
		defs.CreateInstruction("DEC", defs.ZeroPage, cpu.dec, 5, 2),
		{},
		defs.CreateInstruction("INY", defs.Implicit, cpu.iny, 2, 1),
		defs.CreateInstruction("CMP", defs.Immediate, cpu.cmp, 2, 2),
		defs.CreateInstruction("DEX", defs.Implicit, cpu.dex, 2, 1),
		{},
		defs.CreateInstruction("CPY", defs.Absolute, cpu.cpy, 4, 3),
		defs.CreateInstruction("CMP", defs.Absolute, cpu.cmp, 4, 3),
		defs.CreateInstruction("DEC", defs.Absolute, cpu.dec, 6, 3),
		{},

		// 0xD0
		defs.CreateInstruction("BNE", defs.Relative, cpu.bne, 2, 2),
		defs.CreateInstruction("CMP", defs.IndirectY, cpu.cmp, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("CMP", defs.ZeroPageX, cpu.cmp, 4, 2),
		defs.CreateInstruction("DEC", defs.ZeroPageX, cpu.dec, 6, 2),
		{},
		defs.CreateInstruction("CLD", defs.Implicit, cpu.cld, 2, 1),
		defs.CreateInstruction("CMP", defs.AbsoluteYIndexed, cpu.cmp, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("CMP", defs.AbsoluteXIndexed, cpu.cmp, 4, 3),
		defs.CreateInstruction("DEC", defs.AbsoluteXIndexed, cpu.dec, 7, 3),
		{},

		// 0xE0
		defs.CreateInstruction("CPX", defs.Immediate, cpu.cpx, 2, 2),
		defs.CreateInstruction("SBC", defs.IndirectX, cpu.sbc, 6, 2),
		{},
		{},
		defs.CreateInstruction("CPX", defs.ZeroPage, cpu.cpx, 3, 2),
		defs.CreateInstruction("SBC", defs.ZeroPage, cpu.sbc, 3, 2),
		defs.CreateInstruction("INC", defs.ZeroPage, cpu.inc, 5, 2),
		{},
		defs.CreateInstruction("INX", defs.Implicit, cpu.inx, 2, 1),
		defs.CreateInstruction("SBC", defs.Immediate, cpu.sbc, 2, 2),
		defs.CreateInstruction("NOP", defs.Implicit, cpu.nop, 2, 1),
		{},
		defs.CreateInstruction("CPX", defs.Absolute, cpu.cpx, 4, 3),
		defs.CreateInstruction("SBC", defs.Absolute, cpu.sbc, 4, 3),
		defs.CreateInstruction("INC", defs.Absolute, cpu.inc, 6, 3),
		{},

		// 0xF0
		defs.CreateInstruction("BEQ", defs.Relative, cpu.beq, 2, 2),
		defs.CreateInstruction("SBC", defs.IndirectY, cpu.sbc, 5, 2),
		{},
		{},
		{},
		defs.CreateInstruction("SBC", defs.ZeroPageX, cpu.sbc, 4, 2),
		defs.CreateInstruction("INC", defs.ZeroPageX, cpu.inc, 6, 2),
		{},
		defs.CreateInstruction("SED", defs.Implicit, cpu.sed, 2, 1),
		defs.CreateInstruction("SBC", defs.AbsoluteYIndexed, cpu.sbc, 4, 3),
		{},
		{},
		{},
		defs.CreateInstruction("SBC", defs.AbsoluteXIndexed, cpu.sbc, 4, 3),
		defs.CreateInstruction("INC", defs.AbsoluteXIndexed, cpu.inc, 7, 3),
		{},
	}
}

func (cpu *Cpu6502) initAddressModeEvaluators() {
	cpu.addressEvaluators = [13]func(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool){
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

func (cpu Cpu6502) evaluateOperandAddress(addressMode defs.AddressMode, pc defs.Address) (operandAddress defs.Address, pageCrossed bool) {
	if addressMode == defs.Implicit {
		operandAddress = 0
		return
	}

	if cpu.addressEvaluators[addressMode] == nil {
		msg := fmt.Errorf("cannot find address evaluator for address mode \"%d\"", addressMode)
		panic(msg)
	}

	_, evaluatedAddress, _, pageCrossed := cpu.addressEvaluators[addressMode](pc)

	operandAddress = evaluatedAddress

	return
}

func memoryPageDiffer(address defs.Address, finalAddress defs.Address) bool {
	return address&0xFF00 != finalAddress&0xFF00
}
