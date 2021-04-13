package cpu

import "github.com/raulferras/nes-golang/src/nes/defs"

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

// Cpu6502Registers is a representation of the Registers of the NES cpu
type Cpu6502Registers struct {
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
	// supports 65536 direct (unbanked) memory locations, however not all values are sent to the cartridge.
	// It can be accessed either by allowing Cpu6502's internal fetch logic increment the address bus, an interrupt
	// (NMI, Reset, IRQ/BRQ), and using the RTS/JMP/JSR/Branch instructions.
	Pc defs.Address

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

func (registers *Cpu6502Registers) reset() {
	registers.A = 0x00
	registers.X = 0x00
	registers.Y = 0x00
	registers.Sp = 0xFD
	registers.Pc = defs.Address(0xFFFC)
	registers.Status = 0x24
}

func (registers *Cpu6502Registers) stackPointerAddress() defs.Address {
	return defs.Address(0x100 + uint16(registers.Sp))
}

func (registers *Cpu6502Registers) stackPointerPushed() {
	if registers.Sp > 0x00 {
		registers.Sp--
	} else {
		registers.Sp = 0xFF
	}
}

func (registers *Cpu6502Registers) stackPointerPopped() {
	if registers.Sp < 0xFF {
		registers.Sp++
	}
}

// Status register getters
func (registers *Cpu6502Registers) carryFlag() byte {
	return registers.Status & 0x01
}

func (registers *Cpu6502Registers) zeroFlag() byte {
	return registers.Status & 0x02 >> 1
}

func (registers *Cpu6502Registers) interruptFlag() byte {
	return registers.Status & 0x04 >> 2
}

func (registers *Cpu6502Registers) decimalFlag() byte {
	return registers.Status & 0x08 >> 3
}

func (registers *Cpu6502Registers) breakFlag() byte {
	return registers.Status & 0x10 >> 4
}

func (registers *Cpu6502Registers) overflowFlag() byte {
	return registers.Status & 0x40 >> 6
}

func (registers *Cpu6502Registers) negativeFlag() byte {
	return registers.Status & 0x80 >> 7
}

func (registers *Cpu6502Registers) updateNegativeFlag(value byte) {
	//Registers.NegativeFlag = value&0x80 == 0x80
	if value&0x80 == 0x80 {
		registers.Status |= 1 << negativeFlag
	} else {
		registers.Status &= 0b01111111
	}
}

func (registers *Cpu6502Registers) updateZeroFlag(value byte) {
	if value == 0x00 {
		registers.Status |= 1 << zeroFlag
	} else {
		registers.Status &= 0b11111101
	}
}

func (registers *Cpu6502Registers) updateFlag(flag StatusRegisterFlag, state byte) {
	if state == 1 {
		registers.Status |= 1 << flag
	} else {
		registers.Status &= ^(1 << flag)
	}
}

func (registers *Cpu6502Registers) setFlag(flag StatusRegisterFlag) {
	registers.Status |= 1 << flag
}

func (registers *Cpu6502Registers) unsetFlag(flag StatusRegisterFlag) {
	registers.Status &= ^(1 << flag)
}

func (registers *Cpu6502Registers) statusRegister() byte {
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

func (registers *Cpu6502Registers) loadStatusRegister(value byte) {
	registers.Status = value
}

// CreateRegisters creates a properly initialized Cpu6502 Register
func CreateRegisters() Cpu6502Registers {
	return Cpu6502Registers{
		0x00,   // A
		0x00,   // X
		0x00,   // Y
		0xFFFC, // Program Counter
		0xFF,   // Stack Pointer

		0x20,
	}
}
