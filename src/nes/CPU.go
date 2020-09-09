package nes

// CPU Represents a NES cpu
type CPU struct {
	registers CPURegisters
	bus       *Bus

	instructions [256]instruction
}

func (cpu *CPU) tick() {
	// Read opcode
	// -analyze opcode:
	//	-address mode
	//  -get operand
	//  - update PC accordingly
	//  - run operation
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

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()

	ram := RAM{}

	bus := CreateBus(&ram)

	return CPU{
		registers: registers,
		bus:       &bus,
	}
}

func (cpu *CPU) initInstructionsTable() {
	cpu.instructions = [256]instruction{
		{"BRK", implicit, cpu.brk, 0},
		{"ORA", indirectX, cpu.ora, 0},
		{},
		{},
		{},
		{"ORA", zeroPage, cpu.ora, 0},
		{"ASL", zeroPage, cpu.asl, 0},
		{},
		{"PHP", implicit, cpu.php, 0},
		{"ORA", immediate, cpu.ora, 0},
		{"ASL", accumulator, cpu.asl, 0},
		{},
		{},
		{"ORA", absolute, cpu.ora, 0},
		{"ASL", absolute, cpu.asl, 0},
		{},

		// 0x10
		{"BPL", relative, cpu.bpl, 0},
		{"ORA", indirectY, cpu.ora, 0},
		{},
		{},
		{},
		{"ORA", zeroPageX, cpu.ora, 0},
		{"ASL", zeroPageX, cpu.asl, 0},
		{},
		{"CLC", implicit, cpu.clc, 0},
		{"ORA", absoluteYIndexed, cpu.ora, 0},
		{},
		{},
		{},
		{"ORA", absoluteXIndexed, cpu.ora, 0},
		{"ASL", absoluteXIndexed, cpu.asl, 0},
		{},

		// 0x20
		{"JSR", absolute, cpu.jsr, 0},
		{"AND", indirectX, cpu.and, 0},
		{},
		{},
		{"BIT", zeroPage, cpu.bit, 0},
		{"AND", zeroPage, cpu.and, 0},
		{"ROL", zeroPage, cpu.rol, 0},
		{},
		{"PLP", implicit, cpu.plp, 0},
		{"AND", immediate, cpu.and, 0},
		{"ROL", accumulator, cpu.rol, 0},
		{},
		{"BIT", absolute, cpu.bit, 0},
		{"AND", absolute, cpu.and, 0},
		{"ROL", absolute, cpu.rol, 0},
		{},

		// 0x30
		{"BMI", relative, cpu.bmi, 0},
		{"AND", indirectY, cpu.and, 0},
		{},
		{},
		{},
		{"AND", zeroPageX, cpu.and, 0},
		{"ROL", zeroPageX, cpu.rol, 0},
		{},
		{"SEC", implicit, cpu.sec, 0},
		{"AND", absoluteYIndexed, cpu.and, 0},
		{},
		{},
		{},
		{"AND", absoluteXIndexed, cpu.and, 0},
		{"ROL", absoluteXIndexed, cpu.rol, 0},
		{},

		// 0x40
		{"RTI", implicit, cpu.bmi, 0},
		{"EOR", indirectX, cpu.and, 0},
		{},
		{},
		{},
		{"EOR", zeroPage, cpu.eor, 0},
		{"LSR", zeroPage, cpu.lsr, 0},
		{},
		{"PHA", implicit, cpu.pha, 0},
		{"EOR", immediate, cpu.eor, 0},
		{"LSR", accumulator, cpu.lsr, 0},
		{},
		{"JPM", absolute, cpu.jmp, 0},
		{"EOR", absolute, cpu.eor, 0},
		{"LSR", absolute, cpu.lsr, 0},
		{},

		// 0x50
		{"BVC", relative, cpu.bvc, 0},
		{"EOR", indirectY, cpu.eor, 0},
		{},
		{},
		{},
		{"EOR", zeroPageX, cpu.eor, 0},
		{"LSR", zeroPageX, cpu.lsr, 0},
		{},
		{"CLI", implicit, cpu.cli, 0},
		{"EOR", absoluteYIndexed, cpu.eor, 0},
		{},
		{},
		{},
		{"EOR", absoluteXIndexed, cpu.eor, 0},
		{"LSR", absoluteXIndexed, cpu.lsr, 0},
		{},

		// 0x60
		{"RTS", implicit, cpu.rts, 0},
		{"ADC", indirectX, cpu.adc, 0},
		{},
		{},
		{},
		{"ADC", zeroPage, cpu.adc, 0},
		{"ROR", zeroPage, cpu.ror, 0},
		{},
		{"PLA", implicit, cpu.pla, 0},
		{"ADC", immediate, cpu.adc, 0},
		{"ROR", accumulator, cpu.ror, 0},
		{},
		{"JMP", indirect, cpu.jmp, 0},
		{"ADC", absolute, cpu.adc, 0},
		{"ROR", absoluteXIndexed, cpu.ror, 0},
		{},

		// 0x70
		{"BVS", relative, cpu.bvs, 0},
		{"ADC", indirectY, cpu.adc, 0},
		{},
		{},
		{},
		{"ADC", zeroPageX, cpu.adc, 0},
		{"ROR", zeroPageX, cpu.ror, 0},
		{},
		{"SEI", implicit, cpu.sei, 0},
		{"ADC", absoluteYIndexed, cpu.adc, 0},
		{},
		{},
		{},
		{"ADC", absoluteXIndexed, cpu.adc, 0},
		{"ROR", absoluteXIndexed, cpu.ror, 0},
		{},

		// 0x80
		{},
		{"STA", indirectX, cpu.sta, 0},
		{},
		{},
		{"STY", zeroPage, cpu.sty, 0},
		{"STA", zeroPage, cpu.sta, 0},
		{"STX", zeroPage, cpu.stx, 0},
		{},
		{"DEY", implicit, cpu.dey, 0},
		{},
		{"TXA", implicit, cpu.txa, 0},
		{},
		{"STY", absolute, cpu.sty, 0},
		{"STA", absolute, cpu.sta, 0},
		{"STX", absolute, cpu.stx, 0},
		{},

		// 0x90
		{"BCC", relative, cpu.bcc, 0},
		{"STA", indirectY, cpu.sta, 0},
		{},
		{},
		{"STY", zeroPageX, cpu.sty, 0},
		{"STA", zeroPageX, cpu.sta, 0},
		{"STX", zeroPageY, cpu.stx, 0},
		{},
		{"TYA", implicit, cpu.tya, 0},
		{"STA", absoluteYIndexed, cpu.sta, 0},
		{"TXS", implicit, cpu.txs, 0},
		{},
		{},
		{"STA", absoluteXIndexed, cpu.sta, 0},
		{},
		{},

		// 0xA0
		{"LDY", immediate, cpu.ldy, 0},
		{"LDA", indirectX, cpu.lda, 0},
		{"LDX", immediate, cpu.ldx, 0},
		{},
		{"LDY", zeroPage, cpu.ldy, 0},
		{"LDA", zeroPage, cpu.lda, 0},
		{"LDX", zeroPage, cpu.ldx, 0},
		{},
		{"TAY", implicit, cpu.tay, 0},
		{"LDA", immediate, cpu.lda, 0},
		{"TAX", implicit, cpu.tax, 0},
		{},
		{"LDY", absolute, cpu.ldy, 0},
		{"LDA", absolute, cpu.lda, 0},
		{"LDX", absolute, cpu.ldx, 0},
		{},

		// 0xB0
		{"BCS", relative, cpu.bcs, 0},
		{"LDA", indirectY, cpu.lda, 0},
		{},
		{},
		{"LDY", zeroPageX, cpu.ldy, 0},
		{"LDA", zeroPageX, cpu.lda, 0},
		{"LDX", zeroPageY, cpu.ldx, 0},
		{},
		{"CLV", implicit, cpu.clv, 0},
		{"LDA", absoluteYIndexed, cpu.lda, 0},
		{"TSX", implicit, cpu.tsx, 0},
		{},
		{"LDY", absoluteXIndexed, cpu.ldy, 0},
		{"LDA", absoluteXIndexed, cpu.lda, 0},
		{"LDX", absoluteYIndexed, cpu.ldx, 0},
		{},

		// 0xC0
		{"CPY", immediate, cpu.cpy, 0},
		{"CMP", indirectX, cpu.cmp, 0},
		{},
		{},
		{"CPY", zeroPage, cpu.cpy, 0},
		{"CMP", zeroPage, cpu.cmp, 0},
		{"DEC", zeroPage, cpu.dec, 0},
		{},
		{"INY", implicit, cpu.iny, 0},
		{"CMP", immediate, cpu.cmp, 0},
		{"DEX", implicit, cpu.dex, 0},
		{},
		{"CPY", absolute, cpu.cpy, 0},
		{"CMP", absolute, cpu.cmp, 0},
		{"DEC", absolute, cpu.dec, 0},
		{},

		// 0xD0
		{"BNE", relative, cpu.bne, 0},
		{"CMP", indirectY, cpu.cmp, 0},
		{},
		{},
		{},
		{"CMP", zeroPageX, cpu.cmp, 0},
		{"DEC", zeroPageX, cpu.dec, 0},
		{},
		{"CLD", implicit, cpu.cld, 0},
		{"CMP", absoluteYIndexed, cpu.cmp, 0},
		{},
		{},
		{},
		{"CMP", absoluteXIndexed, cpu.cmp, 0},
		{"DEC", absoluteXIndexed, cpu.dec, 0},
		{},

		// 0xE0
		{"CPX", implicit, cpu.cpx, 0},
		{"SBC", indirectX, cpu.sbc, 0},
		{},
		{},
		{"CPX", zeroPage, cpu.cpx, 0},
		{"SBC", zeroPage, cpu.sbc, 0},
		{"INC", zeroPage, cpu.inc, 0},
		{},
		{"INX", implicit, cpu.inx, 0},
		{"SBC", immediate, cpu.sbc, 0},
		{"NOP", implicit, cpu.nop, 0},
		{},
		{"CPX", absolute, cpu.cpx, 0},
		{"SBC", absolute, cpu.sbc, 0},
		{"INC", absolute, cpu.inc, 0},
		{},

		// 0xF0
		{"BEQ", relative, cpu.cpx, 0},
		{"SBC", indirectY, cpu.sbc, 0},
		{},
		{},
		{},
		{"SBC", zeroPageX, cpu.cpx, 0},
		{"INC", zeroPageX, cpu.inc, 0},
		{},
		{"SED", implicit, cpu.inx, 0},
		{"SBC", absoluteYIndexed, cpu.sbc, 0},
		{},
		{},
		{},
		{"SBC", absoluteXIndexed, cpu.sbc, 0},
		{"INC", absoluteXIndexed, cpu.inc, 0},
		{},
	}
}
