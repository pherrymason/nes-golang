package nes

// AddressModeState is
type AddressModeState struct {
	registers Registers
	ram       *RAM
}

func immediate(state AddressModeState) Address {
	return state.registers.Pc
}

func zeroPage(state AddressModeState) Address {
	// 2 bytes
	var low = state.ram.read(state.registers.Pc)

	return Address(low) << 8
}

func zeroPageX(state AddressModeState) Address {
	registers := state.registers
	var low = state.ram.read(registers.Pc) + registers.X

	return Address(low) & 0xFF
}

func zeroPageY(state AddressModeState) Address {
	registers := state.registers
	var low = state.ram.read(registers.Pc) + registers.Y

	return Address(low) & 0xFF
}

func absolute(state AddressModeState) Address {
	registers := state.registers
	low := state.ram.read(registers.Pc)

	// Bug: Missing incrementing programCounter
	high := state.ram.read(registers.Pc + 1)

	return CreateAddress(low, high)
}

func absoluteXIndexed(state AddressModeState) Address {
	registers := state.registers
	low := state.ram.read(registers.Pc)
	high := state.ram.read(registers.Pc + 1)

	address := CreateAddress(low, high)
	address += Address(registers.X)

	return address
}

func absoluteYIndexed(state AddressModeState) Address {
	registers := state.registers
	low := state.ram.read(registers.Pc)
	high := state.ram.read(registers.Pc + 1)

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
func indirect(state AddressModeState) Address {
	registers := state.registers
	ram := state.ram

	// Get Pointer Address
	ptrLow := ram.read(registers.Pc)
	ptrHigh := ram.read(registers.Pc + 1)
	ptrAddress := CreateAddress(ptrLow, ptrHigh)

	finalAddress := ram.read16Bugged(ptrAddress)

	return Address(finalAddress)
}

func preIndexedIndirect(state AddressModeState) Address {
	registers := state.registers
	ram := state.ram

	low := state.ram.read(registers.Pc)
	address := (uint16(low) + uint16(registers.X)) & 0xFF

	finalAddress := ram.read16(Address(address))

	return Address(finalAddress)
}

func postIndexedIndirect(state AddressModeState) Address {
	registers := state.registers
	ram := state.ram

	lo := ram.read(registers.Pc)
	//hi := ram.read(registers.Pc + 1)

	opcodeOperand := CreateAddress(lo, 0x00)

	offsetAddress := ram.read16(opcodeOperand)
	offsetAddress += Word(registers.Y)

	// Todo: Not sure if there is wrap around in adding Y

	return Address(ram.read16(Address(offsetAddress)))
}

func relative(state AddressModeState) {

}
