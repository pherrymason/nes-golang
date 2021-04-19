package cpu

import "github.com/raulferras/nes-golang/src/nes/defs"

func (cpu *Cpu6502) evalImplicit(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter
	address = 0
	cycles = 0
	return
}

/**
 * Immediate addressing allows the programmer to directly specify an 8 bit constant within the instruction.
 * It is indicated by a '#' symbol followed by an numeric expression.
 */
func (cpu *Cpu6502) evalImmediate(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter
	address = programCounter
	pc++
	cycles = 0
	return
}

func (cpu *Cpu6502) evalZeroPage(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	// 2 bytes
	var low = cpu.bus.Read(programCounter)

	address = defs.Address(low)
	pc = programCounter + 1

	return
}

func (cpu *Cpu6502) evalZeroPageX(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	registers := cpu.registers
	var low = cpu.bus.Read(programCounter) + registers.X

	address = defs.Address(low) & 0xFF
	pc = programCounter + 1

	return
}

func (cpu *Cpu6502) evalZeroPageY(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	registers := cpu.registers
	var low = cpu.bus.Read(programCounter) + registers.Y

	address = defs.Address(low) & 0xFF
	pc = programCounter + 1
	return
}

func (cpu *Cpu6502) evalAbsolute(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc += 1

	// Bug: Missing incrementing programCounter
	high := cpu.bus.Read(pc)
	pc += 1

	address = defs.CreateAddress(low, high)

	return
}

func (cpu *Cpu6502) evalAbsoluteXIndexed(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc++

	high := cpu.bus.Read(pc)
	pc++

	address = defs.CreateAddress(low, high)
	address += defs.Address(cpu.registers.X)

	pageCrossed = memoryPageDiffer(address-defs.Address(cpu.registers.X), address)

	return
}

func (cpu *Cpu6502) evalAbsoluteYIndexed(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter
	low := cpu.bus.Read(pc)
	pc++

	high := cpu.bus.Read(pc)
	pc++

	address = defs.CreateAddress(low, high)
	address += defs.Address(cpu.registers.Y)

	pageCrossed = memoryPageDiffer(address-defs.Address(cpu.registers.Y), address)

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
func (cpu *Cpu6502) evalIndirect(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter

	// Get Pointer Address
	ptrLow := cpu.bus.Read(pc)
	pc++

	ptrHigh := cpu.bus.Read(pc)
	pc++

	ptrAddress := defs.CreateAddress(ptrLow, ptrHigh)
	address = cpu.bus.Read16Bugged(ptrAddress)

	return
}

func (cpu *Cpu6502) evalIndirectX(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter

	operand := cpu.bus.Read(pc)
	operand += cpu.registers.X
	operand &= 0xFF

	effectiveLow := cpu.bus.Read(defs.Address(operand))
	effectiveHigh := cpu.bus.Read(defs.Address(operand + 1)) // automatic warp around

	address = defs.CreateAddress(effectiveLow, effectiveHigh)

	return
}

func (cpu *Cpu6502) evalIndirectY(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter

	operand := cpu.bus.Read(pc)
	pc++

	lo := cpu.bus.Read(defs.Address(operand))
	hi := cpu.bus.Read(defs.Address(operand + 1)) // automatic warp around

	address = defs.CreateAddress(lo, hi)
	address += defs.Word(cpu.registers.Y)

	pageCrossed = address&0xFF00 != defs.Address(hi)<<8
	return
}

func (cpu *Cpu6502) evalRelative(programCounter defs.Address) (pc defs.Address, address defs.Address, cycles int, pageCrossed bool) {
	pc = programCounter

	opcodeOperand := cpu.bus.Read(pc)
	pc++

	address = pc
	if opcodeOperand < 0x80 {
		address += defs.Address(opcodeOperand)
	} else {
		address += defs.Address(opcodeOperand) - 0x100
	}

	return
}
