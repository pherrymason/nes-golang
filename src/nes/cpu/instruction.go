package cpu

type Instruction struct {
	name        string
	addressMode AddressMode
	method      OperationMethod
	cycles      byte
	size        byte
}

func CreateInstruction(name string, addressMode AddressMode, method OperationMethod, cycles byte, size byte) Instruction {
	return Instruction{
		name,
		addressMode,
		method,
		cycles,
		size,
	}
}

func (instruction Instruction) Name() string {
	return instruction.name
}

func (instruction Instruction) AddressMode() AddressMode {
	return instruction.addressMode
}

func (instruction Instruction) Size() byte {
	return instruction.size
}

func (instruction Instruction) Method() OperationMethod {
	return instruction.method
}

func (instruction Instruction) Cycles() byte {
	return instruction.cycles
}
