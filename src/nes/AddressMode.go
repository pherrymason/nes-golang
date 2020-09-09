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

// AddressModeState is
type AddressModeState struct {
	registers CPURegisters
	bus       *Bus
}

func evalImmediate(state AddressModeState) Address {
	return state.registers.Pc
}

func evalZeroPage(state AddressModeState) Address {
	// 2 bytes
	var low = state.bus.read(state.registers.Pc)

	return Address(low) << 8
}

func evalZeroPageX(state AddressModeState) Address {
	registers := state.registers
	var low = state.bus.read(registers.Pc) + registers.X

	return Address(low) & 0xFF
}

func evalZeroPageY(state AddressModeState) Address {
	registers := state.registers
	var low = state.bus.read(registers.Pc) + registers.Y

	return Address(low) & 0xFF
}

func evalAbsolute(state AddressModeState) Address {
	registers := state.registers
	low := state.bus.read(registers.Pc)

	// Bug: Missing incrementing programCounter
	high := state.bus.read(registers.Pc + 1)

	return CreateAddress(low, high)
}

func evalAbsoluteXIndexed(state AddressModeState) Address {
	registers := state.registers
	low := state.bus.read(registers.Pc)
	high := state.bus.read(registers.Pc + 1)

	address := CreateAddress(low, high)
	address += Address(registers.X)

	return address
}

func evalAbsoluteYIndexed(state AddressModeState) Address {
	registers := state.registers
	low := state.bus.read(registers.Pc)
	high := state.bus.read(registers.Pc + 1)

	address := CreateAddress(low, high)
	address += Address(registers.Y)

	return address
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
func evalIndirect(state AddressModeState) Address {
	registers := state.registers
	bus := state.bus

	// Get Pointer Address
	ptrLow := bus.read(registers.Pc)
	ptrHigh := bus.read(registers.Pc + 1)
	ptrAddress := CreateAddress(ptrLow, ptrHigh)

	finalAddress := bus.read16Bugged(ptrAddress)

	return Address(finalAddress)
}

func evalIndirectX(state AddressModeState) Address {
	registers := state.registers
	bus := state.bus

	low := state.bus.read(registers.Pc)
	address := (uint16(low) + uint16(registers.X)) & 0xFF

	finalAddress := bus.read16(Address(address))

	return Address(finalAddress)
}

func evalIndirectY(state AddressModeState) Address {
	registers := state.registers
	bus := state.bus

	lo := bus.read(registers.Pc)
	//hi := bus.read(registers.Pc + 1)

	opcodeOperand := CreateAddress(lo, 0x00)

	offsetAddress := bus.read16(opcodeOperand)
	offsetAddress += Word(registers.Y)

	// Todo: Not sure if there is wrap around in adding Y

	return Address(bus.read16(Address(offsetAddress)))
}

func evalRelative(state AddressModeState) Address {
	registers := state.registers
	bus := state.bus

	opcodeOperand := bus.read(registers.Pc)

	address := registers.Pc + 1
	if opcodeOperand < 0x80 {
		address += Address(opcodeOperand)
	} else {
		address += Address(opcodeOperand) - 0x100
	}

	return address
}
