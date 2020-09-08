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

/*
	BNE  Branch on Result not Zero

	branch on Z = 0                  N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	relative      BNE oper      D0    2     2**
*/
func (cpu *CPU) bne(info operation) {
	if !cpu.registers.ZeroFlag {
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
	if !cpu.registers.NegativeFlag {
		cpu.registers.Pc = info.operandAddress
	}
}

/*
	BRK Force Break
	The BRK instruction forces the generation of an interrupt request. The program counter and processor status are pushed on the stack then the IRQ interrupt vector at $FFFE/F is loaded into the PC and the break flag in the status set to one.

	interrupt,                       N Z C I D V
	push PC+2, push SR               - - - 1 - -

	addressing    assembler    opc  bytes  cyles
	--------------------------------------------
	implied       BRK           00    1     7
*/
func (cpu *CPU) brk(info operation) {

	cpu.pushStack(byte(cpu.registers.Pc & 0xFF))
	cpu.pushStack(byte(cpu.registers.Pc >> 8))

	cpu.registers.BreakCommand = true
	cpu.pushStack(cpu.registers.statusRegister())

	cpu.registers.InterruptDisable = true

	cpu.registers.Pc = Address(cpu.ram.read16(0xFFFE))
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
	if cpu.registers.OverflowFlag == byte(1) {
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
	if cpu.registers.OverflowFlag == 0 {
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
	cpu.registers.CarryFlag = 0
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
	cpu.registers.DecimalFlag = false
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
	cpu.registers.InterruptDisable = false
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
	cpu.registers.OverflowFlag = 0
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
	operand := cpu.ram.read(info.operandAddress)
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
	operand := cpu.ram.read(info.operandAddress)
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
	operand := cpu.ram.read(info.operandAddress)
	cpu.compare(cpu.registers.Y, operand)
}

func (cpu *CPU) compare(register byte, operand byte) {
	substraction := register - operand

	cpu.registers.ZeroFlag = false
	cpu.registers.CarryFlag = 0
	cpu.registers.NegativeFlag = false

	if register >= operand {
		cpu.registers.CarryFlag = 1
	}

	if register == operand {
		cpu.registers.ZeroFlag = true
	}

	if substraction&0x80 == 0x80 {
		cpu.registers.NegativeFlag = true
	}
}

func (cpu *CPU) dec(info operation) {
	address := info.operandAddress
	operand := cpu.ram.read(address)

	operand--
	cpu.ram.write(address, operand)

	if operand == 0 {
		cpu.registers.ZeroFlag = true
	} else {
		cpu.registers.ZeroFlag = false
	}

	if operand == 0xFF {
		cpu.registers.NegativeFlag = true
	} else {
		cpu.registers.NegativeFlag = false
	}
}

func (cpu *CPU) dex(info operation) {
	cpu.registers.X--
	operand := cpu.registers.X

	if operand == 0 {
		cpu.registers.ZeroFlag = true
	} else {
		cpu.registers.ZeroFlag = false
	}

	if operand == 0xFF {
		cpu.registers.NegativeFlag = true
	} else {
		cpu.registers.NegativeFlag = false
	}
}

func (cpu *CPU) dey(info operation) {
	operand := cpu.registers.Y

	operand--
	cpu.registers.Y = operand

	if operand == 0 {
		cpu.registers.ZeroFlag = true
	} else {
		cpu.registers.ZeroFlag = false
	}

	if operand == 0xFF {
		cpu.registers.NegativeFlag = true
	} else {
		cpu.registers.NegativeFlag = false
	}
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
	value := cpu.ram.read(info.operandAddress)

	cpu.registers.A = cpu.registers.A ^ value
	cpu.registers.ZeroFlag = cpu.registers.A == 0
	cpu.registers.NegativeFlag = cpu.registers.A&0x80 == 0x80
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
	value := cpu.ram.read(info.operandAddress)
	value += 1

	cpu.ram.write(info.operandAddress, value)
	cpu.registers.NegativeFlag = value&0x80 == 0x80
	cpu.registers.ZeroFlag = value == 0
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
	cpu.registers.NegativeFlag = cpu.registers.X&0x80 == 0x80
	cpu.registers.ZeroFlag = cpu.registers.X == 0
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
	cpu.registers.NegativeFlag = cpu.registers.Y&0x80 == 0x80
	cpu.registers.ZeroFlag = cpu.registers.Y == 0
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
	value := cpu.ram.read16(info.operandAddress)

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
	cpu.registers.A = cpu.ram.read(info.operandAddress)
	cpu.registers.ZeroFlag = cpu.registers.A == 0
	cpu.registers.NegativeFlag = cpu.registers.A&0x80 == 0x80
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
	cpu.registers.X = cpu.ram.read(info.operandAddress)
	cpu.registers.ZeroFlag = cpu.registers.X == 0
	cpu.registers.NegativeFlag = cpu.registers.X&0x80 == 0x80
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
	cpu.registers.Y = cpu.ram.read(info.operandAddress)
	cpu.registers.ZeroFlag = cpu.registers.Y == 0
	cpu.registers.NegativeFlag = cpu.registers.Y&0x80 == 0x80
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
		value = cpu.ram.read(info.operandAddress)
	}

	cpu.registers.CarryFlag = value & 0x01

	value >>= 1
	cpu.registers.ZeroFlag = value == 0

	if info.addressMode == accumulator {
		cpu.registers.A = value
	} else {
		cpu.ram.write(info.operandAddress, value)
	}
}

func (cpu *CPU) nop(info operation) {

}

func (cpu *CPU) ora(info operation) {

}

func (cpu *CPU) pha(info operation) {

}

func (cpu *CPU) php(info operation) {

}

func (cpu *CPU) pla(info operation) {

}

func (cpu *CPU) plp(info operation) {

}

func (cpu *CPU) rol(info operation) {

}

func (cpu *CPU) ror(info operation) {

}

func (cpu *CPU) rti(info operation) {

}

func (cpu *CPU) rts(info operation) {

}

func (cpu *CPU) sbc(info operation) {

}

func (cpu *CPU) sec(info operation) {

}

func (cpu *CPU) sed(info operation) {

}

func (cpu *CPU) sei(info operation) {

}

func (cpu *CPU) sta(info operation) {

}

func (cpu *CPU) stx(info operation) {

}

func (cpu *CPU) sty(info operation) {

}

func (cpu *CPU) tax(info operation) {

}

func (cpu *CPU) tay(info operation) {

}

func (cpu *CPU) tsx(info operation) {

}

func (cpu *CPU) txa(info operation) {

}

func (cpu *CPU) txs(info operation) {

}

func (cpu *CPU) tya(info operation) {

}
