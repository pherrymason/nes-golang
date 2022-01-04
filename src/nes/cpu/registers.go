package cpu

import "github.com/raulferras/nes-golang/src/nes/types"

type StatusRegisterFlag int

const (
	CarryFlag StatusRegisterFlag = iota
	ZeroFlag
	InterruptFlag
	DecimalFlag
	BreakCommandFlag
	UnusedFlag
	OverflowFlag
	NegativeFlag
)

type Registers struct {
	// Accumulator
	// and along with the arithmetic logic unit (ALU), supports using the status register for carrying, overflow
	// detection, and so on.
	A byte

	// Indexes X Y
	// used for several addressing modes. They can be used as loop counters easily, using INC/DEC and branch
	// instructions. Not being the Accumulator, they have limited addressing modes themselves when loading and saving.
	X byte
	Y byte

	// Program Counter:
	// supports 65536 direct (unbanked) memory locations, however not all values are sent to the gamePak.
	// It can be accessed either by allowing Cpu6502's internal fetch logic increment the address bus, an interrupt
	// (NMI, Reset, IRQ/BRQ), and using the RTS/JMP/JSR/Branch instructions.
	Pc types.Address

	// Stack Pointer:
	// The NMOS 65xx processors have 256 bytes of stack memory, ranging
	// from $0100 to $01FF. The S register is a 8-bit offset to the stack
	// page. In other words, whenever anything is being pushed on the
	// stack, it will be stored to the address $0100+S.
	//
	// The Stack pointer can be read and written by transfering its value
	// to or from the index register X (see below) with the TSX and TXS
	// this register is decremented every time a byte is pushed onto the stack,
	// and incremented when a byte is popped off the stack.
	Sp byte

	// Status Processor [NV-BDIZC]
	// N: negative flag
	// V: Overflow flag
	// -: Unused
	// B: Break flag.
	// D:
	// I:
	// Z: zero flag
	// C: carry flag
	Status byte
}

func (registers *Registers) Reset() {
	registers.A = 0x00
	registers.X = 0x00
	registers.Y = 0x00
	registers.Sp = 0xFD
	registers.Pc = types.Address(ResetVectorAddress)
	registers.Status = 0x24
}

func (registers *Registers) SetStackPointer(stackPointer byte) {
	registers.Sp = stackPointer
}

func (registers Registers) StackPointerAddress() types.Address {
	return types.Address(0x100 + uint16(registers.Sp))
}

func (registers *Registers) StackPointerPushed() {
	if registers.Sp > 0x00 {
		registers.Sp--
	} else {
		registers.Sp = 0xFF
	}
}

func (registers *Registers) StackPointerPopped() {
	if registers.Sp < 0xFF {
		registers.Sp++
	}
}

// Status register getters
func (registers *Registers) CarryFlag() byte {
	return registers.Status & 0x01
}

func (registers *Registers) SetCarryFlag(set bool) {
	if set {
		registers.SetFlag(CarryFlag)
	} else {
		registers.unsetFlag(CarryFlag)
	}
}

func (registers *Registers) UnsetCarryFlag() {
	registers.unsetFlag(CarryFlag)
}

func (registers *Registers) ZeroFlag() byte {
	return registers.Status & 0x02 >> 1
}

func (registers *Registers) SetZeroFlag(set bool) {
	if set {
		registers.SetFlag(ZeroFlag)
	} else {
		registers.unsetFlag(ZeroFlag)
	}
}

func (registers *Registers) InterruptFlag() byte {
	return registers.Status & 0x04 >> 2
}

func (registers *Registers) SetInterruptFlag(set bool) {
	if set {
		registers.SetFlag(InterruptFlag)
	} else {
		registers.unsetFlag(InterruptFlag)
	}
}

func (registers *Registers) DecimalFlag() byte {
	return registers.Status & 0x08 >> 3
}

func (registers *Registers) SetDecimalFlag(set bool) {
	if set {
		registers.SetFlag(DecimalFlag)
	} else {
		registers.unsetFlag(DecimalFlag)
	}
}

func (registers *Registers) BreakFlag() byte {
	return registers.Status & 0x10 >> 4
}

func (registers *Registers) OverflowFlag() byte {
	return registers.Status & 0x40 >> 6
}

func (registers *Registers) SetOverflowFlag(set bool) {
	if set {
		registers.SetFlag(OverflowFlag)
	} else {
		registers.unsetFlag(OverflowFlag)
	}
}

func (registers *Registers) NegativeFlag() byte {
	return registers.Status & 0x80 >> 7
}

func (registers *Registers) SetNegativeFlag(set bool) {
	if set {
		registers.SetFlag(NegativeFlag)
	} else {
		registers.unsetFlag(NegativeFlag)
	}
}

func (registers *Registers) UpdateNegativeFlag(value byte) {
	//Registers.NegativeFlag = value&0x80 == 0x80
	if value&0x80 == 0x80 {
		registers.Status |= 1 << NegativeFlag
	} else {
		registers.Status &= 0b01111111
	}
}

func (registers *Registers) UpdateZeroFlag(value byte) {
	if value == 0x00 {
		registers.Status |= 1 << ZeroFlag
	} else {
		registers.Status &= 0b11111101
	}
}

func (registers *Registers) UpdateFlag(flag StatusRegisterFlag, state byte) {
	if state == 1 {
		registers.Status |= 1 << flag
	} else {
		registers.Status &= ^(1 << flag)
	}
}

func (registers *Registers) SetFlag(flag StatusRegisterFlag) {
	registers.Status |= 1 << flag
}

func (registers *Registers) unsetFlag(flag StatusRegisterFlag) {
	registers.Status &= ^(1 << flag)
}

func (registers *Registers) StatusRegister() byte {
	// Bit 5 is always read as set
	return registers.Status | 0x20
}

func (registers *Registers) LoadStatusRegister(value byte) {
	// From http://nesdev.com/the%20%27B%27%20flag%20&%20BRK%20instruction.txt
	// ...when the flags are restored (via PLP or RTI), the Break flag (4 bit) is discarded.
	registers.Status = value & 0b11101111
}

// CreateRegisters creates a properly initialized Cpu6502 Register
func CreateRegisters() Registers {
	return Registers{
		0x00,   // A
		0x00,   // X
		0x00,   // Y
		0xFFFC, // Program Counter
		0xFF,   // Stack Pointer
		0x20,
	}
}
