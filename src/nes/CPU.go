package nes

import (
	"fmt"

	"github.com/raulferras/nes-golang/src/log"
)

// CPU Represents a NES cpu
type CPU struct {
	registers CPURegisters
	bus       *Bus

	instructions [256]instruction
	nose         [13]func(programCounter Address) (pc Address, address Address, cycles int)
	debug        bool
	logger       log.Logger
}

func CreateCPU(bus *Bus) CPU {
	registers := CreateRegisters()
	return CPU{
		registers: registers,
		bus:       bus,
		debug:     false,
	}
}

func CreateCPUDebuggable(bus *Bus, logger log.Logger) CPU {
	registers := CreateRegisters()
	return CPU{
		registers: registers,
		bus:       bus,
		debug:     true,
		logger:    logger,
	}
}

// CreateCPUWithBus creates a CPU with a Bus, Useful for tests
func CreateCPUWithBus() CPU {
	registers := CreateRegisters()

	ram := RAM{}
	gamepak := GamePak{
		header: Header{},
		prgROM: make([]byte, 0xFFFF),
	}

	bus := CreateBus(&ram)
	bus.attachCartridge(&gamepak)

	return CPU{
		registers: registers,
		bus:       &bus,
		debug:     false,
	}
}

func (cpu *CPU) reset() {
	cpu.registers.reset()

	// Read Reset Vector
	address := cpu.bus.read16(0xFFFC)
	cpu.registers.Pc = Address(address)
}

func (cpu *CPU) tick() {

	// Read opcode
	if cpu.debug {
		cpu.printStep()
	}

	opcode := cpu.read(cpu.registers.Pc)
	cpu.registers.Pc++

	instruction := cpu.instructions[opcode]
	if instruction.method == nil {
		msg := fmt.Errorf("Error: Opcode 0x%X not implemented!", opcode)
		panic(msg)
	}

	_, operandAddress, _ := cpu.nose[instruction.addressMode](cpu.registers.Pc)
	instruction.method(InfoStep{
		instruction.addressMode,
		operandAddress,
	})

	// -analyze opcode:
	//	-address mode
	//  -get operand
	//  - update PC accordingly
	//  - run InfoStep

}

func (cpu *CPU) printStep() {

	pc := cpu.registers.Pc
	opcode := cpu.read(pc)
	pc++
	instruction := cpu.instructions[opcode]

	_, evaluatedAddress, _ := cpu.nose[instruction.addressMode](pc)

	var msg string
	msg += fmt.Sprintf("%X", pc) + "  "
	msg += fmt.Sprintf("%X ", opcode) + " "

	for i := byte(0); i < (instruction.size - 1); i++ {
		msg += fmt.Sprintf("%X ", cpu.read(pc+Address(i)))
	}

	for i := len(msg); i <= 16; i++ {
		msg += " "
	}

	msg += instruction.name + " "

	if instruction.addressMode == immediate {
		msg += "#"
	} else {
		msg += fmt.Sprintf("$%X", evaluatedAddress)
	}

	for i := len(msg); i <= 48; i++ {
		msg += " "
	}

	msg += fmt.Sprintf(
		"A:%X X:%X Y:%X P:%X SP:%X PPU:___,___ CYC:%d",
		cpu.registers.A,
		cpu.registers.X,
		cpu.registers.Y,
		cpu.registers.Status,
		cpu.registers.Sp,
		0,
	)

	cpu.logger.Log(msg)
}

func (cpu *CPU) pushStack(value byte) {
	address := cpu.registers.stackPointerAddress()
	cpu.bus.write(
		address,
		value,
	)

	cpu.registers.stackPointerPushed()
}

func (cpu *CPU) popStack() byte {
	cpu.registers.stackPointerPopped()
	address := cpu.registers.stackPointerAddress()
	return cpu.bus.read(address)
}

// Reads value located at Program Counter and increments it
func (cpu *CPU) fetch() byte {
	value := cpu.bus.read(cpu.registers.Pc)
	cpu.registers.Pc++

	return value
}

func (cpu *CPU) read(address Address) byte {
	return cpu.bus.read(address)
}

func (cpu *CPU) read16(address Address) Word {
	low := cpu.bus.read(address)
	high := cpu.bus.read(address + 1)

	return CreateWord(low, high)
}

func (cpu *CPU) write(address Address, value byte) {
	cpu.bus.write(address, value)
}

func (cpu *CPU) initInstructionsTable() {
	cpu.instructions = [256]instruction{
		{"BRK", implicit, cpu.brk, 7, 1},
		{"ORA", indirectX, cpu.ora, 6, 2},
		{},
		{},
		{},
		{"ORA", zeroPage, cpu.ora, 3, 2},
		{"ASL", zeroPage, cpu.asl, 5, 2},
		{},
		{"PHP", implicit, cpu.php, 3, 1},
		{"ORA", immediate, cpu.ora, 2, 2},
		{"ASL", accumulator, cpu.asl, 2, 1},
		{},
		{},
		{"ORA", absolute, cpu.ora, 4, 3},
		{"ASL", absolute, cpu.asl, 6, 3},
		{},

		// 0x10
		{"BPL", relative, cpu.bpl, 2, 2},
		{"ORA", indirectY, cpu.ora, 5, 2},
		{},
		{},
		{},
		{"ORA", zeroPageX, cpu.ora, 4, 2},
		{"ASL", zeroPageX, cpu.asl, 6, 2},
		{},
		{"CLC", implicit, cpu.clc, 2, 1},
		{"ORA", absoluteYIndexed, cpu.ora, 4, 3},
		{},
		{},
		{},
		{"ORA", absoluteXIndexed, cpu.ora, 4, 3},
		{"ASL", absoluteXIndexed, cpu.asl, 7, 3},
		{},

		// 0x20
		{"JSR", absolute, cpu.jsr, 6, 3},
		{"AND", indirectX, cpu.and, 6, 2},
		{},
		{},
		{"BIT", zeroPage, cpu.bit, 3, 2},
		{"AND", zeroPage, cpu.and, 3, 2},
		{"ROL", zeroPage, cpu.rol, 5, 2},
		{},
		{"PLP", implicit, cpu.plp, 4, 1},
		{"AND", immediate, cpu.and, 2, 2},
		{"ROL", accumulator, cpu.rol, 2, 1},
		{},
		{"BIT", absolute, cpu.bit, 4, 3},
		{"AND", absolute, cpu.and, 4, 3},
		{"ROL", absolute, cpu.rol, 6, 3},
		{},

		// 0x30
		{"BMI", relative, cpu.bmi, 2, 2},
		{"AND", indirectY, cpu.and, 5, 2},
		{},
		{},
		{},
		{"AND", zeroPageX, cpu.and, 4, 2},
		{"ROL", zeroPageX, cpu.rol, 6, 2},
		{},
		{"SEC", implicit, cpu.sec, 2, 1},
		{"AND", absoluteYIndexed, cpu.and, 4, 3},
		{},
		{},
		{},
		{"AND", absoluteXIndexed, cpu.and, 4, 3},
		{"ROL", absoluteXIndexed, cpu.rol, 7, 3},
		{},

		// 0x40
		{"RTI", implicit, cpu.bmi, 6, 1},
		{"EOR", indirectX, cpu.eor, 6, 2},
		{},
		{},
		{},
		{"EOR", zeroPage, cpu.eor, 3, 2},
		{"LSR", zeroPage, cpu.lsr, 5, 2},
		{},
		{"PHA", implicit, cpu.pha, 3, 1},
		{"EOR", immediate, cpu.eor, 2, 2},
		{"LSR", accumulator, cpu.lsr, 2, 1},
		{},
		{"JMP", absolute, cpu.jmp, 3, 3},
		{"EOR", absolute, cpu.eor, 4, 3},
		{"LSR", absolute, cpu.lsr, 6, 3},
		{},

		// 0x50
		{"BVC", relative, cpu.bvc, 2, 2},
		{"EOR", indirectY, cpu.eor, 5, 2},
		{},
		{},
		{},
		{"EOR", zeroPageX, cpu.eor, 4, 2},
		{"LSR", zeroPageX, cpu.lsr, 6, 2},
		{},
		{"CLI", implicit, cpu.cli, 2, 1},
		{"EOR", absoluteYIndexed, cpu.eor, 4, 3},
		{},
		{},
		{},
		{"EOR", absoluteXIndexed, cpu.eor, 4, 3},
		{"LSR", absoluteXIndexed, cpu.lsr, 7, 3},
		{},

		// 0x60
		{"RTS", implicit, cpu.rts, 6, 1},
		{"ADC", indirectX, cpu.adc, 6, 2},
		{},
		{},
		{},
		{"ADC", zeroPage, cpu.adc, 3, 2},
		{"ROR", zeroPage, cpu.ror, 5, 2},
		{},
		{"PLA", implicit, cpu.pla, 4, 1},
		{"ADC", immediate, cpu.adc, 2, 2},
		{"ROR", accumulator, cpu.ror, 2, 1},
		{},
		{"JMP", indirect, cpu.jmp, 5, 3},
		{"ADC", absolute, cpu.adc, 4, 3},
		{"ROR", absolute, cpu.ror, 6, 3},
		{},

		// 0x70
		{"BVS", relative, cpu.bvs, 2, 2},
		{"ADC", indirectY, cpu.adc, 5, 2},
		{},
		{},
		{},
		{"ADC", zeroPageX, cpu.adc, 4, 2},
		{"ROR", zeroPageX, cpu.ror, 6, 2},
		{},
		{"SEI", implicit, cpu.sei, 2, 1},
		{"ADC", absoluteYIndexed, cpu.adc, 4, 3},
		{},
		{},
		{},
		{"ADC", absoluteXIndexed, cpu.adc, 4, 3},
		{"ROR", absoluteXIndexed, cpu.ror, 7, 3},
		{},

		// 0x80
		{},
		{"STA", indirectX, cpu.sta, 6, 2},
		{},
		{},
		{"STY", zeroPage, cpu.sty, 3, 2},
		{"STA", zeroPage, cpu.sta, 3, 2},
		{"STX", zeroPage, cpu.stx, 3, 2},
		{},
		{"DEY", implicit, cpu.dey, 2, 1},
		{},
		{"TXA", implicit, cpu.txa, 2, 1},
		{},
		{"STY", absolute, cpu.sty, 4, 3},
		{"STA", absolute, cpu.sta, 4, 3},
		{"STX", absolute, cpu.stx, 4, 3},
		{},

		// 0x90
		{"BCC", relative, cpu.bcc, 2, 2},
		{"STA", indirectY, cpu.sta, 6, 2},
		{},
		{},
		{"STY", zeroPageX, cpu.sty, 4, 2},
		{"STA", zeroPageX, cpu.sta, 4, 2},
		{"STX", zeroPageY, cpu.stx, 4, 2},
		{},
		{"TYA", implicit, cpu.tya, 2, 1},
		{"STA", absoluteYIndexed, cpu.sta, 5, 3},
		{"TXS", implicit, cpu.txs, 2, 1},
		{},
		{},
		{"STA", absoluteXIndexed, cpu.sta, 5, 3},
		{},
		{},

		// 0xA0
		{"LDY", immediate, cpu.ldy, 2, 2},
		{"LDA", indirectX, cpu.lda, 6, 2},
		{"LDX", immediate, cpu.ldx, 2, 2},
		{},
		{"LDY", zeroPage, cpu.ldy, 3, 2},
		{"LDA", zeroPage, cpu.lda, 3, 2},
		{"LDX", zeroPage, cpu.ldx, 3, 2},
		{},
		{"TAY", implicit, cpu.tay, 2, 1},
		{"LDA", immediate, cpu.lda, 2, 2},
		{"TAX", implicit, cpu.tax, 2, 1},
		{},
		{"LDY", absolute, cpu.ldy, 4, 3},
		{"LDA", absolute, cpu.lda, 4, 3},
		{"LDX", absolute, cpu.ldx, 4, 3},
		{},

		// 0xB0
		{"BCS", relative, cpu.bcs, 2, 2},
		{"LDA", indirectY, cpu.lda, 5, 2},
		{},
		{},
		{"LDY", zeroPageX, cpu.ldy, 4, 2},
		{"LDA", zeroPageX, cpu.lda, 4, 2},
		{"LDX", zeroPageY, cpu.ldx, 4, 2},
		{},
		{"CLV", implicit, cpu.clv, 2, 1},
		{"LDA", absoluteYIndexed, cpu.lda, 4, 3},
		{"TSX", implicit, cpu.tsx, 2, 1},
		{},
		{"LDY", absoluteXIndexed, cpu.ldy, 4, 3},
		{"LDA", absoluteXIndexed, cpu.lda, 4, 3},
		{"LDX", absoluteYIndexed, cpu.ldx, 4, 3},
		{},

		// 0xC0
		{"CPY", immediate, cpu.cpy, 2, 2},
		{"CMP", indirectX, cpu.cmp, 6, 2},
		{},
		{},
		{"CPY", zeroPage, cpu.cpy, 3, 2},
		{"CMP", zeroPage, cpu.cmp, 3, 2},
		{"DEC", zeroPage, cpu.dec, 5, 2},
		{},
		{"INY", implicit, cpu.iny, 2, 1},
		{"CMP", immediate, cpu.cmp, 2, 2},
		{"DEX", implicit, cpu.dex, 2, 1},
		{},
		{"CPY", absolute, cpu.cpy, 4, 3},
		{"CMP", absolute, cpu.cmp, 4, 3},
		{"DEC", absolute, cpu.dec, 6, 3},
		{},

		// 0xD0
		{"BNE", relative, cpu.bne, 2, 2},
		{"CMP", indirectY, cpu.cmp, 5, 2},
		{},
		{},
		{},
		{"CMP", zeroPageX, cpu.cmp, 4, 2},
		{"DEC", zeroPageX, cpu.dec, 6, 2},
		{},
		{"CLD", implicit, cpu.cld, 2, 1},
		{"CMP", absoluteYIndexed, cpu.cmp, 4, 3},
		{},
		{},
		{},
		{"CMP", absoluteXIndexed, cpu.cmp, 4, 3},
		{"DEC", absoluteXIndexed, cpu.dec, 7, 3},
		{},

		// 0xE0
		{"CPX", immediate, cpu.cpx, 2, 2},
		{"SBC", indirectX, cpu.sbc, 6, 2},
		{},
		{},
		{"CPX", zeroPage, cpu.cpx, 3, 2},
		{"SBC", zeroPage, cpu.sbc, 3, 2},
		{"INC", zeroPage, cpu.inc, 5, 2},
		{},
		{"INX", implicit, cpu.inx, 2, 1},
		{"SBC", immediate, cpu.sbc, 2, 2},
		{"NOP", implicit, cpu.nop, 2, 1},
		{},
		{"CPX", absolute, cpu.cpx, 4, 3},
		{"SBC", absolute, cpu.sbc, 4, 3},
		{"INC", absolute, cpu.inc, 6, 3},
		{},

		// 0xF0
		{"BEQ", relative, cpu.cpx, 2, 2},
		{"SBC", indirectY, cpu.sbc, 5, 2},
		{},
		{},
		{},
		{"SBC", zeroPageX, cpu.cpx, 4, 2},
		{"INC", zeroPageX, cpu.inc, 6, 2},
		{},
		{"SED", implicit, cpu.inx, 2, 1},
		{"SBC", absoluteYIndexed, cpu.sbc, 4, 3},
		{},
		{},
		{},
		{"SBC", absoluteXIndexed, cpu.sbc, 4, 3},
		{"INC", absoluteXIndexed, cpu.inc, 7, 3},
		{},
	}
}

func (cpu *CPU) initAddressModeEvaluators() {
	cpu.nose = [13]func(programCounter Address) (pc Address, address Address, cycles int){
		nil,
		nil,
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
	}
}
