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

/*
	ADC  Add Memory to Accumulator with Carry

     A + M + C -> A, C                N Z C I D V
                                      + + + - - +

     addressing    assembler    opc  bytes  cyles
     --------------------------------------------
     immidiate     ADC #oper     69    2     2
     zeropage      ADC oper      65    2     3
     zeropage,X    ADC oper,X    75    2     4
     absolute      ADC oper      6D    3     4
     absolute,X    ADC oper,X    7D    3     4*
     absolute,Y    ADC oper,Y    79    3     4*
     (indirect,X)  ADC (oper,X)  61    2     6
	 (indirect),Y  ADC (oper),Y  71    2     5*
*/
func (cpu *CPU) adc(info operation) {
	carryIn := cpu.registers.CarryFlag
	a := cpu.registers.A
	value := cpu.ram.read(info.operandAddress)
	adc := uint16(a) + uint16(value) + uint16(carryIn)
	adc8 := cpu.registers.A + value + cpu.registers.CarryFlag

	cpu.registers.A = adc8
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)

	if (adc) > 0xFF {
		cpu.registers.CarryFlag = 1
	} else {
		cpu.registers.CarryFlag = 0
	}

	// If the sign of the sum matches either the sign of A or the sign of v, then you don't overflow
	if ((uint16(a) ^ adc) & (uint16(value) ^ adc) & 0x80) > 0 {
		cpu.registers.OverflowFlag = 1
	} else {
		cpu.registers.OverflowFlag = 0
	}
}

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
		cpu.registers.CarryFlag = cpu.registers.A >> 7 & 0x01
		cpu.registers.A = cpu.registers.A << 1
		cpu.registers.updateNegativeFlag(cpu.registers.A)
		cpu.registers.updateZeroFlag(cpu.registers.A)
	} else {
		value := cpu.ram.read(info.operandAddress)
		cpu.registers.CarryFlag = value >> 7 & 0x01
		value = value << 1
		cpu.ram.write(info.operandAddress, value)
		cpu.registers.updateNegativeFlag(value)
		cpu.registers.updateZeroFlag(value)
	}
}

/*
	BCC  Branch on Carry Clear

	branch on C = 0                  N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BCC oper      90    2     2**
*/
func (cpu *CPU) bcc(info operation) {
	if cpu.registers.CarryFlag == 0 {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BCS  Branch on Carry Set

	branch on C = 1                 N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BCS oper      B0    2     2**
*/
func (cpu *CPU) bcs(info operation) {
	if cpu.registers.CarryFlag == 1 {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BEQ  Branch on Result Zero

	branch on Z = 1             N Z C I D V
								- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BEQ oper      F0    2     2**
*/
func (cpu *CPU) beq(info operation) {
	if cpu.registers.ZeroFlag {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BIT  Test Bits in Memory with Accumulator

	bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
	the zeroflag is set to the result of operand AND accumulator.
	The result is not kept.

	A AND M, M7 -> N, M6 -> V        N Z C I D V
									M7 + - - - M6

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	zeropage      BIT oper      24    2     3
	absolute      BIT oper      2C    3     4
*/
func (cpu *CPU) bit(info operation) {
	value := cpu.ram.read(info.operandAddress)
	cpu.registers.NegativeFlag = value&0x80 == 0x80
	cpu.registers.OverflowFlag = (value >> 6) & 0x01
	cpu.registers.ZeroFlag = value&cpu.registers.A == 0
}

/*
	BMI  Branch on Result Minus

	branch on N = 1                 N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BMI oper      30    2     2**
*/
func (cpu *CPU) bmi(info operation) {
	if cpu.registers.NegativeFlag {
		cpu.registers.Pc = info.operandAddress
	}
}

// CreateCPU a CPU
func CreateCPU() CPU {

	registers := CreateRegisters()
	ram := RAM{}

	return CPU{registers, &ram}
}
