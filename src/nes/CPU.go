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

func (cpu *CPU) pushStack(value byte) {
	address := cpu.registers.spAddress()
	cpu.ram.write(
		address,
		value,
	)

	cpu.registers.spPushed()
}

func (cpu *CPU) popStack() byte {
	cpu.registers.spPopped()
	address := cpu.registers.spAddress()
	return cpu.ram.read(address)
}

type operation struct {
	addressMode    AddressMode
	operandAddress Address
}

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()
	ram := RAM{}

	return CPU{registers, &ram}
}
