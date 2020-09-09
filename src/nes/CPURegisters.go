package nes

type StatusRegisterFlag int

const (
	carryFlag StatusRegisterFlag = iota
	zeroFlag
	interruptFlag
	decimalFlag
	breakCommandFlag
	unusedFlag
	overflowFlag
	negativeFlag
)

// CPURegisters is a representation of the registers of the NES cpu
type CPURegisters struct {
	// Accumulator
	// and along with the arithmetic logic unit (ALU), supports using the status register for carrying, overflow
	// detection, and so on.
	A byte

	// Indexes X Y
	// used for several addressing modes. They can be used as loop counters easily, using INC/DEC and branch
	// instructions. Not being the accumulator, they have limited addressing modes themselves when loading and saving.
	X byte
	Y byte

	// Program Counter:
	// supports 65536 direct (unbanked) memory locations, however not all values are sent to the cartridge.
	// It can be accessed either by allowing CPU's internal fetch logic increment the address bus, an interrupt
	// (NMI, Reset, IRQ/BRQ), and using the RTS/JMP/JSR/Branch instructions.
	Pc Address

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

	// Status Processor
	Status byte
}

func (registers *CPURegisters) reset() {
	registers.A = 0x00
	registers.X = 0x00
	registers.Y = 0x00
	registers.Sp = 0xFF
	registers.Pc = Address(0x0000)
	registers.Status = 0x20
}

func (registers *CPURegisters) stackPointerAddress() Address {
	return Address(0x100 + uint16(registers.Sp))
}

func (registers *CPURegisters) stackPointerPushed() {
	if registers.Sp > 0x00 {
		registers.Sp--
	} else {
		registers.Sp = 0xFF
	}
}

func (registers *CPURegisters) stackPointerPopped() {
	if registers.Sp < 0xFF {
		registers.Sp++
	}
}

// Status register getters
func (registers *CPURegisters) carryFlag() byte {
	return registers.Status & 0x01
}

func (registers *CPURegisters) zeroFlag() byte {
	return registers.Status & 0x02 >> 1
}

func (registers *CPURegisters) interruptFlag() byte {
	return registers.Status & 0x04 >> 2
}

func (registers *CPURegisters) decimalFlag() byte {
	return registers.Status & 0x08 >> 3
}

func (registers *CPURegisters) breakFlag() byte {
	return registers.Status & 0x10 >> 4
}

func (registers *CPURegisters) overflowFlag() byte {
	return registers.Status & 0x40 >> 6
}

func (registers *CPURegisters) negativeFlag() byte {
	return registers.Status & 0x80 >> 7
}

func (registers *CPURegisters) updateNegativeFlag(value byte) {
	//registers.NegativeFlag = value&0x80 == 0x80
	if value&0x80 == 0x80 {
		registers.Status |= 1 << negativeFlag
	} else {
		registers.Status &= 0b01111111
	}
}

func (registers *CPURegisters) updateZeroFlag(value byte) {
	if value == 0x00 {
		registers.Status |= 1 << zeroFlag
	} else {
		registers.Status &= 0b11111101
	}
}

func (registers *CPURegisters) updateFlag(flag StatusRegisterFlag, state byte) {
	if state == 1 {
		registers.Status |= 1 << flag
	} else {
		registers.Status &= ^(1 << flag)
	}
}

func (registers *CPURegisters) setFlag(flag StatusRegisterFlag) {
	registers.Status |= 1 << flag
}

func (registers *CPURegisters) unsetFlag(flag StatusRegisterFlag) {
	registers.Status &= ^(1 << flag)
}

func (registers *CPURegisters) statusRegister() byte {
	var value byte = 0x00

	if registers.carryFlag() == 1 {
		value |= 0x01
	}

	if registers.zeroFlag() == 1 {
		value |= 0x02
	}

	if registers.interruptFlag() == 1 {
		value |= 0x04
	}

	// Decimal mode
	if registers.decimalFlag() == 1 {
		value |= 0x08
	}

	if registers.breakFlag() == 1 {
		value |= 0x10
	}

	// Alway 1 flag
	value |= 0x20

	// Signed overflow
	if registers.overflowFlag() == 1 {
		value |= 0x40
	}

	// Processor Status flag
	if registers.negativeFlag() == 1 {
		value |= 0x80
	}

	return value
}

func (registers *CPURegisters) loadStatusRegister(value byte) {
	registers.Status = value
}

// CreateRegisters creates a properly initialized CPU Register
func CreateRegisters() CPURegisters {
	return CPURegisters{
		0x00,   // A
		0x00,   // X
		0x00,   // Y
		0x0000, // Program Counter
		0xFF,   // Stack Pointer

		0x20,
	}
}
