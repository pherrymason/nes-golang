package nes

type StatusRegister struct {
	// Unsigned overflow
	CarryFlag byte

	ZeroFlag bool

	InterruptDisable bool

	BreakCommand bool

	// Signed overflow
	OverflowFlag byte

	// Processor Status flag
	NegativeFlag bool
}

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
	// Unsigned overflow
	CarryFlag byte

	ZeroFlag bool

	InterruptDisable bool

	DecimalFlag bool

	BreakCommand bool

	// Signed overflow
	OverflowFlag byte

	// Processor Status flag
	NegativeFlag bool
}

func (registers *CPURegisters) reset() {
	registers.A = 0x00
	registers.X = 0x00
	registers.Y = 0x00
	registers.Sp = 0xFF
	registers.Pc = Address(0x0000)

	registers.CarryFlag = 0
	registers.ZeroFlag = false
	registers.InterruptDisable = false
	registers.BreakCommand = false
	registers.OverflowFlag = 0
	registers.NegativeFlag = false
}

func (registers *CPURegisters) spAddress() Address {
	return Address(0x100 + uint16(registers.Sp))
}

func (registers *CPURegisters) spPushed() {
	if registers.Sp > 0x00 {
		registers.Sp--
	} else {
		registers.Sp = 0xFF
	}
}

func (registers *CPURegisters) spPopped() {
	if registers.Sp < 0xFF {
		registers.Sp++
	}
}

func (registers *CPURegisters) updateNegativeFlag(value byte) {
	registers.NegativeFlag = value&0x80 == 0x80
}

func (registers *CPURegisters) updateZeroFlag(value byte) {
	registers.ZeroFlag = value == 0x00
}

func (registers *CPURegisters) statusRegister() byte {
	var value byte = 0x00

	if registers.CarryFlag == 1 {
		value |= 0x01
	}

	if registers.ZeroFlag {
		value |= 0x02
	}

	if registers.InterruptDisable {
		value |= 0x04
	}

	// Decimal mode
	if registers.DecimalFlag {
		value |= 0x08
	}

	if registers.BreakCommand {
		value |= 0x10
	}

	// Alway 1 flag
	value |= 0x20

	// Signed overflow
	if registers.OverflowFlag == 1 {
		value |= 0x40
	}

	// Processor Status flag
	if registers.NegativeFlag {
		value |= 0x80
	}

	return value
}

// CreateRegisters creates a properly initialized CPU Register
func CreateRegisters() CPURegisters {
	return CPURegisters{
		0x00,   // A
		0x00,   // X
		0x00,   // Y
		0x0000, // Program Counter
		0xFF,   // Stack Pointer

		0,     // Carry Flag
		false, // Zero Flag
		false, // Interrupt Disable
		false, // Break Command
		false, // Decimal Flag
		0,     // Overflow Flag
		false, // Negative Flag
	}
}
