package cpu

import "github.com/raulferras/nes-golang/src/nes/defs"

func (cpu *Cpu6502) evalImplicit(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter
	address = 0
	cycles = 0
	return
}

func (cpu *Cpu6502) evalImmediate(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter
	address = programCounter
	pc++
	cycles = 0
	return
}

func (cpu *Cpu6502) evalZeroPage(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	// 2 bytes
	var low = cpu.bus.Read(programCounter)

	address = defs.Address(low) << 8
	pc = cpu.registers.Pc + 1

	return
}

func (cpu *Cpu6502) evalZeroPageX(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	registers := cpu.registers
	var low = cpu.bus.Read(programCounter) + registers.X

	address = defs.Address(low) & 0xFF
	pc = programCounter + 1

	return
}

func (cpu *Cpu6502) evalZeroPageY(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	registers := cpu.registers
	var low = cpu.bus.Read(programCounter) + registers.Y

	address = defs.Address(low) & 0xFF
	pc = programCounter + 1
	return
}

func (cpu *Cpu6502) evalAbsolute(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc += 1

	// Bug: Missing incrementing programCounter
	high := cpu.bus.Read(pc)
	pc += 1

	address = defs.CreateAddress(low, high)

	return
}

func (cpu *Cpu6502) evalAbsoluteXIndexed(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc++

	high := cpu.bus.Read(pc)
	pc++

	address = defs.CreateAddress(low, high)
	address += defs.Address(cpu.registers.X)

	return
}

func (cpu *Cpu6502) evalAbsoluteYIndexed(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc++

	high := cpu.bus.Read(pc)
	pc++

	address = defs.CreateAddress(low, high)
	address += defs.Address(cpu.registers.Y)

	return
}

// Address Mode: Indirect
// The supplied 16-bit address is read to get the actual 16-bit address.
// This is Instruction is unusual in that it has a bug in the hardware! To emulate its
// function accurately, we also need to emulate this bug. If the low byte of the
// supplied address is 0xFF, then to read the high byte of the actual address
// we need to cross a page boundary. This doesnt actually work on the chip as
// designed, instead it wraps back around in the same page, yielding an
// invalid actual address
// Example: supplied address is (0x1FF), LSB will be 0x00 and MSB will be 0x01 instead of 0x02.

// If the 16-bit argument of an Indirect JMP is located between 2 pages (0x01FF and 0x0200 for example),
// then the LSB will be read from 0x01FF and the MSB will be read from 0x0100.
// This is an actual hardware bug in early revisions of the 6502 which happen to be present
// in the 2A03 used by the NES.
func (cpu *Cpu6502) evalIndirect(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter

	// Get Pointer Address
	ptrLow := cpu.bus.Read(pc)
	pc++

	ptrHigh := cpu.bus.Read(pc)
	pc++

	ptrAddress := defs.CreateAddress(ptrLow, ptrHigh)
	finalAddress := cpu.bus.Read16Bugged(ptrAddress)
	address = defs.Address(finalAddress)

	return
}

func (cpu *Cpu6502) evalIndirectX(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter

	low := cpu.bus.Read(pc)
	pc++

	ptrAddress := (uint16(low) + uint16(cpu.registers.X)) & 0xFF

	address = defs.Address(cpu.bus.Read16(defs.Address(ptrAddress)))

	return
}

func (cpu *Cpu6502) evalIndirectY(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter

	lo := cpu.bus.Read(pc)
	pc++

	opcodeOperand := defs.CreateAddress(lo, 0x00)

	offsetAddress := cpu.bus.Read16(opcodeOperand)
	offsetAddress += defs.Word(cpu.registers.Y)

	// Todo: Not sure if there is wrap around in adding Y

	address = defs.Address(cpu.bus.Read16(defs.Address(offsetAddress)))

	return
}

func (cpu *Cpu6502) evalRelative(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int) {
	pc = programCounter

	opcodeOperand := cpu.bus.Read(pc)
	pc++

	address = cpu.registers.Pc + 1
	if opcodeOperand < 0x80 {
		address += defs.Address(opcodeOperand)
	} else {
		address += defs.Address(opcodeOperand) - 0x100
	}

	return
}