package nes

// CPU Represents a NES cpu
type CPU struct {
	registers CPURegisters
	bus       *Bus
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
	cpu.bus.write(
		address,
		value,
	)

	cpu.registers.spPushed()
}

func (cpu *CPU) popStack() byte {
	cpu.registers.spPopped()
	address := cpu.registers.spAddress()
	return cpu.bus.read(address)
}

func (cpu *CPU) read(address Address) byte {
	return cpu.bus.read(address)
}

func (cpu *CPU) read16(address Address) Word {
	low := cpu.bus.read(address)
	high := cpu.bus.read(address + 1)

	return CreateWord(low, high)
}

func (cpu *CPU) write(address Address, value byte) {
	cpu.bus.write(address, value)
}

type operation struct {
	addressMode    AddressMode
	operandAddress Address
}

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()

	ram := RAM{}

	bus := CreateBus(&ram)

	return CPU{registers, &bus}
}
