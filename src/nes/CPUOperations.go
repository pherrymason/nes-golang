package nes

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

	http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	https://forums.nesdev.com/viewtopic.php?t=6331
*/
func (cpu *CPU) adc(info operation) {
	carryIn := cpu.registers.carryFlag()
	a := cpu.registers.A
	value := cpu.read(info.operandAddress)
	adc := uint16(a) + uint16(value) + uint16(carryIn)
	adc8 := cpu.registers.A + value + cpu.registers.carryFlag()

	cpu.registers.A = adc8
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)

	if (adc) > 0xFF {
		cpu.registers.setFlag(carryFlag)
	} else {
		cpu.registers.unsetFlag(carryFlag)
	}

	// The exclusive-or bitwise operator is a neat little tool to check if the sign of two numbers is the same
	// If the sign of the sum matches either the sign of A or the sign of v, then you don't overflow
	if ((uint16(a) ^ adc) & (uint16(value) ^ adc) & 0x80) > 0 {
		cpu.registers.setFlag(overflowFlag)
	} else {
		cpu.registers.unsetFlag(overflowFlag)
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
	cpu.registers.A &= cpu.read(operandAddress)
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
		cpu.registers.updateFlag(carryFlag, cpu.registers.A>>7&0x01)
		cpu.registers.A = cpu.registers.A << 1
		cpu.registers.updateNegativeFlag(cpu.registers.A)
		cpu.registers.updateZeroFlag(cpu.registers.A)
	} else {
		value := cpu.read(info.operandAddress)
		cpu.registers.updateFlag(carryFlag, value>>7&0x01)
		value = value << 1
		cpu.write(info.operandAddress, value)
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
	if cpu.registers.carryFlag() == 0 {
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
	if cpu.registers.carryFlag() == 1 {
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
	if cpu.registers.zeroFlag() == 1 {
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
	value := cpu.read(info.operandAddress)
	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateFlag(overflowFlag, (value>>6)&0x01)
	cpu.registers.updateZeroFlag(value & cpu.registers.A)
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
	if cpu.registers.negativeFlag() == 1 {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BNE  Branch on Result not Zero

	branch on Z = 0                  N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BNE oper      D0    2     2**
*/
func (cpu *CPU) bne(info operation) {
	// CHeck how to negate a bit and apply it here
	//if !cpu.registers.zeroFlag() == 1 {
	if cpu.registers.zeroFlag() == 0 {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BPL  Branch on Result Plus

	branch on N = 0             N Z C I D V
								- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BPL oper      10    2     2**
*/
func (cpu *CPU) bpl(info operation) {
	//if !cpu.registers.NegativeFlag {
	if cpu.registers.negativeFlag() == 0 {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BRK Force Break
	The BRK instruction forces the generation of an interrupt request.
    The program counter and processor status are pushed on the stack then
    the IRQ interrupt vector at $FFFE/F is loaded into the PC and the break
    flag in the status set to one.

	interrupt,                       N Z C I D V
	push PC+2, push SR               - - - 1 - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       BRK           00    1     7
*/
func (cpu *CPU) brk(info operation) {
	// Store PC in stack
	cpu.pushStack(byte(cpu.registers.Pc & 0xFF))
	cpu.pushStack(byte(cpu.registers.Pc >> 8))
	//
	cpu.registers.updateFlag(breakCommandFlag, 1)
	cpu.pushStack(cpu.registers.Status)

	cpu.registers.updateFlag(interruptFlag, 1)

	cpu.registers.Pc = Address(cpu.read16(0xFFFE))
}

/*
	BVC  Branch on Overflow Clear
	branch on V = 0               N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BVC oper      50    2     2**
*/
func (cpu *CPU) bvc(info operation) {
	if cpu.registers.overflowFlag() == byte(1) {
		return
	}

	cpu.registers.Pc = info.operandAddress
}

/*
	BVS  Branch on Overflow Set
	branch on V = 1               N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BVC oper      70    2     2**
*/
func (cpu *CPU) bvs(info operation) {
	if cpu.registers.overflowFlag() == 0 {
		return
	}

	cpu.registers.Pc = info.operandAddress
}

/*
	CLC  Clear Carry Flag
	0 -> C                        N Z C I D V
								  - - 0 - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       CLC           18    1     2
*/
func (cpu *CPU) clc(info operation) {
	cpu.registers.updateFlag(carryFlag, 0)
}

/*
	CLD  Clear Decimal Mode
	0 -> D                        N Z C I D V
								  - - - - 0 -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       CLD           D8    1     2
*/
func (cpu *CPU) cld(info operation) {
	cpu.registers.updateFlag(decimalFlag, 0)
}

/*
	CLI  Clear Interrupt Disable Bit
	0 -> I                        N Z C I D V
								  - - - 0 - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       CLI           58    1     2
*/
func (cpu *CPU) cli(info operation) {
	cpu.registers.updateFlag(interruptFlag, 0)
}

/*
	CLV  Clear Overflow Flag
	0 -> V                        N Z C I D V
								  - - - - - 0

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       CLV           B8    1     2
*/
func (cpu *CPU) clv(info operation) {
	cpu.registers.updateFlag(overflowFlag, 0)
}

/*
	CMP (CoMPare accumulator)

	Affects Flags: S Z C

	MODE           SYNTAX       HEX LEN TIM
	Immediate     CMP #$44      $C9  2   2
	Zero Page     CMP $44       $C5  2   3
	Zero Page,X   CMP $44,X     $D5  2   4
	Absolute      CMP $4400     $CD  3   4
	Absolute,X    CMP $4400,X   $DD  3   4+
	Absolute,Y    CMP $4400,Y   $D9  3   4+
	Indirect,X    CMP ($44,X)   $C1  2   6
	Indirect,Y    CMP ($44),Y   $D1  2   5+

	+ add 1 cycle if page boundary crossed

	Compare sets flags as if a subtraction had been carried out. If the value in the accumulator is equal or greater than the compared value, the Carry will be set. The equal (Z) and sign (S) flags will be set based on equality or lack thereof and the sign (i.e. A>=$80) of the accumulator.
*/
func (cpu *CPU) cmp(info operation) {
	operand := cpu.read(info.operandAddress)
	cpu.compare(cpu.registers.A, operand)
}

/*
	CPX (ComPare X register)

	Affects Flags: S Z C

	MODE           SYNTAX       HEX LEN TIM
	Immediate     CPX #$44      $E0  2   2
	Zero Page     CPX $44       $E4  2   3
	Absolute      CPX $4400     $EC  3   4
*/
func (cpu *CPU) cpx(info operation) {
	operand := cpu.read(info.operandAddress)
	cpu.compare(cpu.registers.X, operand)
}

/*
	CPY (ComPare Y register)

	Affects Flags: S Z C

	MODE           SYNTAX       HEX LEN TIM
	Immediate     CPY #$44      $C0  2   2
	Zero Page     CPY $44       $C4  2   3
	Absolute      CPY $4400     $CC  3   4
*/
func (cpu *CPU) cpy(info operation) {
	operand := cpu.read(info.operandAddress)
	cpu.compare(cpu.registers.Y, operand)
}

func (cpu *CPU) compare(register byte, operand byte) {
	substraction := register - operand

	cpu.registers.updateFlag(zeroFlag, 0)
	cpu.registers.updateFlag(carryFlag, 0)
	cpu.registers.updateFlag(negativeFlag, 0)

	if register >= operand {
		cpu.registers.updateFlag(carryFlag, 1)
	}

	if register == operand {
		cpu.registers.updateFlag(zeroFlag, 1)
	}

	//if substraction&0x80 == 0x80 {
	cpu.registers.updateNegativeFlag(substraction)
	//	}
}

func (cpu *CPU) dec(info operation) {
	address := info.operandAddress
	operand := cpu.read(address)

	operand--
	cpu.write(address, operand)

	cpu.registers.updateZeroFlag(operand)
	//if operand == 0 {
	//	cpu.registers.ZeroFlag = true
	//} else {
	//	cpu.registers.ZeroFlag = false
	//}
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0xFF {
	//	cpu.registers.updateFlag(negativeFlag, 1)
	//} else {
	//	cpu.registers.updateFlag(negativeFlag, 0)
	//}
}

func (cpu *CPU) dex(info operation) {
	cpu.registers.X--
	operand := cpu.registers.X

	cpu.registers.updateZeroFlag(operand)
	//if operand == 0 {
	//	cpu.registers.ZeroFlag = true
	//} else {
	//	cpu.registers.ZeroFlag = false
	//}
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0xFF {
	//	cpu.registers.updateFlag(negativeFlag, 1)
	//} else {
	//	cpu.registers.updateFlag(negativeFlag, )NegativeFlag = false
	//}
}

func (cpu *CPU) dey(info operation) {
	operand := cpu.registers.Y

	operand--
	cpu.registers.Y = operand

	cpu.registers.updateZeroFlag(operand)
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0 {
	//	cpu.registers.ZeroFlag = true
	//} else {
	//	cpu.registers.ZeroFlag = false
	//}
	//
	//if operand == 0xFF {
	//	cpu.registers.NegativeFlag = true
	//} else {
	//	cpu.registers.NegativeFlag = false
	//}
}

/*
	EOR (bitwise Exclusive OR)
	Affects Flags: S Z
	A EOR M -> A                     N Z C I D V
                                     + + - - - -

	MODE           SYNTAX       HEX LEN TIM
	Immediate     EOR #$44      $49  2   2
	Zero Page     EOR $44       $45  2   3
	Zero Page,X   EOR $44,X     $55  2   4
	Absolute      EOR $4400     $4D  3   4
	Absolute,X    EOR $4400,X   $5D  3   4+
	Absolute,Y    EOR $4400,Y   $59  3   4+
	Indirect,X    EOR ($44,X)   $41  2   6
	Indirect,Y    EOR ($44),Y   $51  2   5+

	+ add 1 cycle if page boundary crossed
*/
func (cpu *CPU) eor(info operation) {
	value := cpu.read(info.operandAddress)

	cpu.registers.A = cpu.registers.A ^ value
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	INC  Increment Memory by One
	M + 1 -> M                    N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	zeropage      INC oper      E6    2     5
	zeropage,X    INC oper,X    F6    2     6
	absolute      INC oper      EE    3     6
	absolute,X    INC oper,X    FE    3     7
*/
func (cpu *CPU) inc(info operation) {
	value := cpu.read(info.operandAddress)
	value += 1

	cpu.write(info.operandAddress, value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateNegativeFlag(value)
}

/*
	INX  Increment Index X by One
	X + 1 -> X                N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       INX           E8    1     2
*/
func (cpu *CPU) inx(info operation) {
	cpu.registers.X += 1
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)

}

/*
	INY  Increment Index Y by One
	Y + 1 -> Y                    N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       INY           C8    1     2
*/
func (cpu *CPU) iny(info operation) {
	cpu.registers.Y += 1
	cpu.registers.updateZeroFlag(cpu.registers.Y)
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
}

/*
	JMP  Jump to New Location
	(PC+1) -> PCL                    N Z C I D V
	(PC+2) -> PCH                    - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	absolute      JMP oper      4C    3     3
	indirect      JMP (oper)    6C    3     5
*/
func (cpu *CPU) jmp(info operation) {
	cpu.registers.Pc = info.operandAddress
}

/*
	JSR  Jump to New Location Saving Return Address
	push (PC+2),                     N Z C I D V
	(PC+1) -> PCL                    - - - - - -
	(PC+2) -> PCH

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	absolute      JSR oper      20    3     6
*/
func (cpu *CPU) jsr(info operation) {
	value := cpu.read16(info.operandAddress)

	// TODO CHECK HERE because ProgramCounter should point to Opcode
	cpu.registers.Pc -= 3
	cpu.pushStack(byte(cpu.registers.Pc & 0xFF))
	cpu.pushStack(byte(cpu.registers.Pc >> 8))

	cpu.registers.Pc = Address(value)
}

/*
	LDA  Load Accumulator with Memory
	M -> A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	immediate     LDA #oper     A9    2     2
	zeropage      LDA oper      A5    2     3
	zeropage,X    LDA oper,X    B5    2     4
	absolute      LDA oper      AD    3     4
	absolute,X    LDA oper,X    BD    3     4*
	absolute,Y    LDA oper,Y    B9    3     4*
	(indirect,X)  LDA (oper,X)  A1    2     6
	(indirect),Y  LDA (oper),Y  B1    2     5*
*/
func (cpu *CPU) lda(info operation) {
	cpu.registers.A = cpu.read(info.operandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	LDX  Load Index X with Memory
	M -> X                    N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	immediate     LDX #oper     A2    2     2
	zeropage      LDX oper      A6    2     3
	zeropage,Y    LDX oper,Y    B6    2     4
	absolute      LDX oper      AE    3     4
	absolute,Y    LDX oper,Y    BE    3     4*
*/
func (cpu *CPU) ldx(info operation) {
	cpu.registers.X = cpu.read(info.operandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)
}

/*
	LDY  Load Index Y with Memory
	M -> Y                N Z C I D V
						  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	immidiate     LDY #oper     A0    2     2
	zeropage      LDY oper      A4    2     3
	zeropage,X    LDY oper,X    B4    2     4
	absolute      LDY oper      AC    3     4
	absolute,X    LDY oper,X    BC    3     4*
*/
func (cpu *CPU) ldy(info operation) {
	cpu.registers.Y = cpu.read(info.operandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.Y)
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
}

/*
	LSR  Shift One Bit Right (Memory or Accumulator)
	0 -> [76543210] -> C      N Z C I D V
							  0 + + - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	accumulator   LSR A         4A    1     2
	zeropage      LSR oper      46    2     5
	zeropage,X    LSR oper,X    56    2     6
	absolute      LSR oper      4E    3     6
	absolute,X    LSR oper,X    5E    3     7
*/
func (cpu *CPU) lsr(info operation) {
	var value byte
	if info.addressMode == accumulator {
		value = cpu.registers.A
	} else {
		value = cpu.read(info.operandAddress)
	}

	//cpu.registers.CarryFlag = value & 0x01
	cpu.registers.updateFlag(carryFlag, value&0x01)

	value >>= 1
	cpu.registers.updateZeroFlag(value)

	if info.addressMode == accumulator {
		cpu.registers.A = value
	} else {
		cpu.write(info.operandAddress, value)
	}
}

/*
	NOP  No Operation
	---                           N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       NOP           EA    1     2
*/
func (cpu *CPU) nop(info operation) {

}

/*
	ORA  OR Memory with Accumulator
	A OR M -> A               N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	immidiate     ORA #oper     09    2     2
	zeropage      ORA oper      05    2     3
	zeropage,X    ORA oper,X    15    2     4
	absolute      ORA oper      0D    3     4
	absolute,X    ORA oper,X    1D    3     4*
	absolute,Y    ORA oper,Y    19    3     4*
	(indirect,X)  ORA (oper,X)  01    2     6
	(indirect),Y  ORA (oper),Y  11    2     5*
*/
func (cpu *CPU) ora(info operation) {
	value := cpu.read(info.operandAddress)
	cpu.registers.A |= value
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	PHA  Push Accumulator on Stack
	push A                        N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       PHA           48    1     3
*/
func (cpu *CPU) pha(info operation) {
	cpu.pushStack(cpu.registers.A)
}

/*
	PHP  Push Processor Status on Stack
	push SR                       N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       PHP           08    1     3
*/
func (cpu *CPU) php(info operation) {
	value := cpu.registers.statusRegister()
	cpu.pushStack(value)
}

/*
	PLA  Pull Accumulator from Stack
	pull A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       PLA           68    1     4
*/
func (cpu *CPU) pla(info operation) {
	cpu.registers.A = cpu.popStack()
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)
}

/*
	PLP  Pull Processor Status from Stack
	pull SR                       N Z C I D V
								  from stack

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       PLP           28    1     4
*/
func (cpu *CPU) plp(info operation) {
	value := cpu.popStack()

	cpu.registers.loadStatusRegister(value)
}

/*
	ROL  Rotate One Bit Left (Memory or Accumulator)
	C <- [76543210] <- C          N Z C I D V
								  + + + - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	accumulator   ROL A         2A    1     2
	zeropage      ROL oper      26    2     5
	zeropage,X    ROL oper,X    36    2     6
	absolute      ROL oper      2E    3     6
	absolute,X    ROL oper,X    3E    3     7
*/
func (cpu *CPU) rol(info operation) {
	var newCarry byte
	var value byte
	if info.addressMode == accumulator {
		newCarry = cpu.registers.A & 0x80 >> 7
		cpu.registers.A <<= 1
		cpu.registers.A |= cpu.registers.carryFlag()
		value = cpu.registers.A
	} else {
		value = cpu.read(info.operandAddress)
		newCarry = value & 0x80 >> 7
		value <<= 1
		value |= cpu.registers.carryFlag()
		cpu.write(info.operandAddress, value)
	}

	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateFlag(carryFlag, newCarry)
}

/*
	ROR  Rotate One Bit Right (Memory or Accumulator)
	C -> [76543210] -> C          N Z C I D V
								  + + + - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	accumulator   ROR A         6A    1     2
	zeropage      ROR oper      66    2     5
	zeropage,X    ROR oper,X    76    2     6
	absolute      ROR oper      6E    3     6
	absolute,X    ROR oper,X    7E    3     7
*/
func (cpu *CPU) ror(info operation) {
	var newCarry byte
	var value byte
	if info.addressMode == accumulator {
		newCarry = cpu.registers.A & 0x01
		cpu.registers.A >>= 1
		cpu.registers.A |= cpu.registers.carryFlag() << 7
		value = cpu.registers.A
	} else {
		value = cpu.read(info.operandAddress)
		newCarry = value & 0x01
		value >>= 1
		value |= cpu.registers.carryFlag() << 7
		cpu.write(info.operandAddress, value)
	}

	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateFlag(carryFlag, newCarry)
}

/*
	RTI  Return from Interrupt

	pull SR, pull PC              N Z C I D V
								  from stack

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       RTI           40    1     6
*/
func (cpu *CPU) rti(info operation) {
	statusRegister := cpu.popStack()
	cpu.registers.loadStatusRegister(statusRegister)

	msb := cpu.popStack()
	lsb := cpu.popStack()
	cpu.registers.Pc = CreateAddress(lsb, msb)
}

/*
	RTS  Return from Subroutine
	pull PC, PC+1 -> PC           N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       RTS           60    1     6
*/
func (cpu *CPU) rts(info operation) {
	msb := cpu.popStack()
	lsb := cpu.popStack()
	cpu.registers.Pc = CreateAddress(lsb, msb)
}

/*
	SBC  Subtract Memory from Accumulator with Borrow
	A - M - C -> A                N Z C I D V
								  + + + - - +

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	immidiate     SBC #oper     E9    2     2
	zeropage      SBC oper      E5    2     3
	zeropage,X    SBC oper,X    F5    2     4
	absolute      SBC oper      ED    3     4
	absolute,X    SBC oper,X    FD    3     4*
	absolute,Y    SBC oper,Y    F9    3     4*
	(indirect,X)  SBC (oper,X)  E1    2     6
	(indirect),Y  SBC (oper),Y  F1    2     5*
*/
func (cpu *CPU) sbc(info operation) {
	value := cpu.read(info.operandAddress)
	borrow := (1 - cpu.registers.carryFlag()) & 0x01 // == !CarryFlag
	a := cpu.registers.A
	result := a - value - borrow
	cpu.registers.A = result

	cpu.registers.updateZeroFlag(byte(result))
	cpu.registers.updateNegativeFlag(byte(result))

	// Set overflow flag
	if (a^cpu.registers.A)&0x80 != 0 && (a^value)&0x80 != 0 {
		cpu.registers.updateFlag(overflowFlag, 1)
	} else {
		cpu.registers.updateFlag(overflowFlag, 0)
	}

	if int(a)-int(value)-int(borrow) < 0 {
		cpu.registers.updateFlag(carryFlag, 0)
	} else {
		cpu.registers.updateFlag(carryFlag, 1)
	}
}

/*
	SEC  Set Carry Flag
	1 -> C                        N Z C I D V
								  - - 1 - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       SEC           38    1     2
*/
func (cpu *CPU) sec(info operation) {
	cpu.registers.updateFlag(carryFlag, 1)
}

/*
	SED  Set Decimal Flag
	1 -> D                    N Z C I D V
							  - - - - 1 -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       SED           F8    1     2
*/
func (cpu *CPU) sed(info operation) {
	cpu.registers.updateFlag(decimalFlag, 1)
}

func (cpu *CPU) sei(info operation) {
	cpu.registers.updateFlag(interruptFlag, 1)
}

func (cpu *CPU) sta(info operation) {
	cpu.write(info.operandAddress, cpu.registers.A)
}

func (cpu *CPU) stx(info operation) {
	cpu.write(info.operandAddress, cpu.registers.X)
}

func (cpu *CPU) sty(info operation) {
	cpu.write(info.operandAddress, cpu.registers.Y)
}

/*
	TAX  Transfer Accumulator to Index X
	A -> X                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       TAX           AA    1     2
*/
func (cpu *CPU) tax(info operation) {
	cpu.registers.X = cpu.registers.A
	cpu.registers.updateNegativeFlag(cpu.registers.X)
	cpu.registers.updateZeroFlag(cpu.registers.X)
}

/*
	TAY  Transfer Accumulator to Index Y
	A -> Y                    N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       TAY           A8    1     2
*/
func (cpu *CPU) tay(info operation) {
	cpu.registers.Y = cpu.registers.A
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
	cpu.registers.updateZeroFlag(cpu.registers.Y)
}

/*
	TSX  Transfer Stack Pointer to Index X
	SP -> X                       N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       TSX           BA    1     2
*/
func (cpu *CPU) tsx(info operation) {
	cpu.registers.X = cpu.popStack()
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)
}

/*
	TXA  Transfer Index X to Accumulator
	X -> A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       TXA           8A    1     2
*/
func (cpu *CPU) txa(info operation) {
	cpu.registers.A = cpu.registers.X
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	TXS  Transfer Index X to Stack Pointer
	X -> SP                       N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       TXS           9A    1     2
*/
func (cpu *CPU) txs(info operation) {
	cpu.registers.Sp = cpu.registers.X
}

/*
	TYA  Transfer Index Y to Accumulator
	 Y -> A                           N Z C I D V
									  + + - - - -
	 addressing    assembler    opc  bytes  cyles
	 --------------------------------------------
	 implied       TYA           98    1     2

	*  add 1 to cycles if page boundery is crossed
	** add 1 to cycles if branch occurs on same page
	 add 2 to cycles if branch occurs to different page


	 Legend to Flags:  + .... modified
					   - .... not modified
					   1 .... set
					   0 .... cleared
					  M6 .... memory bit 6
					  M7 .... memory bit 7


	Note on assembler syntax:
	Most assemblers employ "OPC *oper" for forced zeropage addressing.
*/
func (cpu *CPU) tya(info operation) {
	cpu.registers.A = cpu.registers.Y
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}
