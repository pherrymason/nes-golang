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

type operation struct {
	addressMode    AddressMode
	operandAddress Address
}

// --- Operations

//	Performs a logical AND on the operand and the accumulator and stores the result in the accumulator
//
// 	Addressing Mode 	Assembly Language Form 	Opcode 	# Bytes 	# Cycles
// 	Immediate 			AND #Operand 			29 		2 			2
//	Zero Page 			AND Operand 			25 		2 			3
//	Zero Page, X 		AND Operand, X 			35 		2 			4
//	Absolute 			AND Operand 			2D 		3 			4
//	Absolute, X 		AND Operand, X 			3D 		3 			4*
//	Absolute, Y 		AND Operand, Y 			39 		3 			4*
//	(Indirect, X) 		AND (Operand, X)	 	21 		2 			6
//	(Indirect), Y 		AND (Operand), Y 		31 		2 			5*
//	* Add 1 if page boundary is crossed.
func (cpu *CPU) and(operandAddress Address) {
	cpu.registers.A &= cpu.ram.read(operandAddress)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)
}

/*
	ASL  Shift Left One Bit (Memory or Accumulator)

     C <- [76543210] <- 0             N Z C I D V
                                      + + + - - -

     addressing    assembler    opc  bytes  cyles
     --------------------------------------------
     accumulator   ASL A         0A    1     2
     zeropage      ASL oper      06    2     5
     zeropage,X    ASL oper,X    16    2     6
	 absolute      ASL oper      0E    3     6
*/
func (cpu *CPU) asl(info operation) {
	if info.addressMode == accumulator {
		cpu.registers.CarryFlag = cpu.registers.A>>7 == 1
		cpu.registers.A = cpu.registers.A << 1
		cpu.registers.updateNegativeFlag(cpu.registers.A)
		cpu.registers.updateZeroFlag(cpu.registers.A)
	} else {
		value := cpu.ram.read(info.operandAddress)
		cpu.registers.CarryFlag = value>>7 == 1
		value = value << 1
		cpu.ram.write(info.operandAddress, value)
		cpu.registers.updateNegativeFlag(value)
		cpu.registers.updateZeroFlag(value)
	}
}

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()
	ram := RAM{}

	return CPU{registers, &ram}
}
