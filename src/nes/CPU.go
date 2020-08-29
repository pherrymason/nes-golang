package nes

// CPU Represents a NES cpu
type CPU struct {
	registers CPURegisters
	ram       *RAM
}

func (cpu *CPU) tick() {
	// Read opcode
	// -analyze opcode:
	//	-address mode
	//  -get operand
	//  - update PC accordingly
	//  - run operation
}

// --- Operations

/*
	Performs a logical AND on the operand and the accumulator and stores the result in the accumulator
*/
func (cpu *CPU) and(operandAddress Address) {
	cpu.registers.A &= cpu.ram.read(operandAddress)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)
}

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()
	ram := RAM{}

	return CPU{registers, &ram}
}
