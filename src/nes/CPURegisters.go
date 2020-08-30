package nes

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

	NegativeFlag bool
	ZeroFlag     bool
	CarryFlag    bool
}

func (registers *CPURegisters) reset() {
	registers.A = 0x00
	registers.X = 0x00
	registers.Y = 0x00
	registers.Sp = 0xFF
	registers.Pc = Address(0x0000)
}

func (registers *CPURegisters) updateNegativeFlag(value byte) {
	registers.NegativeFlag = value&0x80 == 0x80
}

func (registers *CPURegisters) updateZeroFlag(value byte) {
	registers.ZeroFlag = value == 0x00
}

// CreateRegisters creates a properly initialized CPU Register
func CreateRegisters() CPURegisters {
	return CPURegisters{0x00, 0x00, 0x00, 0x0000, 0xFF, false, false, false}
}
