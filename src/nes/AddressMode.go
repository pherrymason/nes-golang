package nes

// AddressMode is an enum of the available Addressing Modes in this cpu
type AddressMode int

const (
	implicit AddressMode = iota
	accumulator
	immediate
	zeroPage
	zeroPageX
	zeroPageY
	absolute
	absoluteXIndexed
	absoluteYIndexed
	indirect
	indirectX
	indirectY
	relative
)

func (cpu *CPU) evalImmediate(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter
	address = programCounter
	pc++
	cycles = 0
	return
}

func (cpu *CPU) evalZeroPage(programCounter Address) (pc Address, address Address, cycles int) {
	// 2 bytes
	var low = cpu.bus.read(programCounter)

	address = Address(low) << 8
	pc = cpu.registers.Pc + 1

	return
}

func (cpu *CPU) evalZeroPageX(programCounter Address) (pc Address, address Address, cycles int) {
	registers := cpu.registers
	var low = cpu.bus.read(programCounter) + registers.X

	address = Address(low) & 0xFF
	pc = programCounter + 1

	return
}

func (cpu *CPU) evalZeroPageY(programCounter Address) (pc Address, address Address, cycles int) {
	registers := cpu.registers
	var low = cpu.bus.read(programCounter) + registers.Y

	address = Address(low) & 0xFF
	pc = programCounter + 1
	return
}

func (cpu *CPU) evalAbsolute(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter
	low := cpu.bus.read(pc)
	pc += 1

	// Bug: Missing incrementing programCounter
	high := cpu.bus.read(pc)
	pc += 1

	address = CreateAddress(low, high)

	return
}

func (cpu *CPU) evalAbsoluteXIndexed(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter
	low := cpu.bus.read(pc)
	pc++

	high := cpu.bus.read(pc)
	pc++

	address = CreateAddress(low, high)
	address += Address(cpu.registers.X)

	return
}

func (cpu *CPU) evalAbsoluteYIndexed(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter
	low := cpu.bus.read(pc)
	pc++

	high := cpu.bus.read(pc)
	pc++

	address = CreateAddress(low, high)
	address += Address(cpu.registers.Y)

	return
}

// Address Mode: Indirect
// The supplied 16-bit address is read to get the actual 16-bit address.
// This is instruction is unusual in that it has a bug in the hardware! To emulate its
// function accurately, we also need to emulate this bug. If the low byte of the
// supplied address is 0xFF, then to read the high byte of the actual address
// we need to cross a page boundary. This doesnt actually work on the chip as
// designed, instead it wraps back around in the same page, yielding an
// invalid actual address
// Example: supplied address is (0x1FF), LSB will be 0x00 and MSB will be 0x01 instead of 0x02.

// If the 16-bit argument of an indirect JMP is located between 2 pages (0x01FF and 0x0200 for example),
// then the LSB will be read from 0x01FF and the MSB will be read from 0x0100.
// This is an actual hardware bug in early revisions of the 6502 which happen to be present
// in the 2A03 used by the NES.
func (cpu *CPU) evalIndirect(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter

	// Get Pointer Address
	ptrLow := cpu.bus.read(pc)
	pc++

	ptrHigh := cpu.bus.read(pc)
	pc++

	ptrAddress := CreateAddress(ptrLow, ptrHigh)
	finalAddress := cpu.bus.read16Bugged(ptrAddress)
	address = Address(finalAddress)

	return
}

func (cpu *CPU) evalIndirectX(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter

	low := cpu.bus.read(pc)
	pc++

	ptrAddress := (uint16(low) + uint16(cpu.registers.X)) & 0xFF

	address = Address(cpu.bus.read16(Address(ptrAddress)))

	return
}

func (cpu *CPU) evalIndirectY(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter

	lo := cpu.bus.read(pc)
	pc++
	//hi := bus.read(Registers.Pc + 1)

	opcodeOperand := CreateAddress(lo, 0x00)

	offsetAddress := cpu.bus.read16(opcodeOperand)
	offsetAddress += Word(cpu.registers.Y)

	// Todo: Not sure if there is wrap around in adding Y

	address = Address(cpu.bus.read16(Address(offsetAddress)))

	return
}

func (cpu *CPU) evalRelative(programCounter Address) (pc Address, address Address, cycles int) {
	pc = programCounter

	opcodeOperand := cpu.bus.read(pc)
	pc++

	address = cpu.registers.Pc + 1
	if opcodeOperand < 0x80 {
		address += Address(opcodeOperand)
	} else {
		address += Address(opcodeOperand) - 0x100
	}

	return
}
