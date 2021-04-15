package cpu

import (
	"github.com/raulferras/nes-golang/src/nes/defs"
)

/*
	ADC  Add Memory to Accumulator with Carry
     A + M + C -> A, C                N Z C I D V
                                      + + + - - +

     addressing    assembler    opc  bytes  cycles
     --------------------------------------------
     immidiate     ADC #oper     69    2     2
     zeropage      ADC oper      65    2     3
     zeropage,X    ADC oper,X    75    2     4
     Absolute      ADC oper      6D    3     4
     Absolute,X    ADC oper,X    7D    3     4*
     Absolute,Y    ADC oper,Y    79    3     4*
     (Indirect,X)  ADC (oper,X)  61    2     6
	 (Indirect),Y  ADC (oper),Y  71    2     5*

	http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	https://forums.nesdev.com/viewtopic.php?t=6331
*/
func (cpu *Cpu6502) adc(info defs.InfoStep) {
	carryIn := cpu.registers.carryFlag()
	a := cpu.registers.A
	value := cpu.Read(info.OperandAddress)
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

//	Performs a logical AND on the operand and the Accumulator and stores the result in the Accumulator
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
func (cpu *Cpu6502) and(info defs.InfoStep) {
	cpu.registers.A &= cpu.Read(info.OperandAddress)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)
}

/*
	ASL  Shift Left One Bit (Memory or Accumulator)

     C <- [76543210] <- 0             N Z C I D V
                                      + + + - - -

     addressing    assembler    opc  bytes  cycles
     --------------------------------------------
     Accumulator   ASL A         0A    1     2
     zeropage      ASL oper      06    2     5
     zeropage,X    ASL oper,X    16    2     6
	 Absolute      ASL oper      0E    3     6
*/
func (cpu *Cpu6502) asl(info defs.InfoStep) {
	if info.AddressMode == defs.Implicit {
		cpu.registers.updateFlag(carryFlag, cpu.registers.A>>7&0x01)
		cpu.registers.A = cpu.registers.A << 1
		cpu.registers.updateNegativeFlag(cpu.registers.A)
		cpu.registers.updateZeroFlag(cpu.registers.A)
	} else {
		value := cpu.Read(info.OperandAddress)
		cpu.registers.updateFlag(carryFlag, value>>7&0x01)
		value = value << 1
		cpu.write(info.OperandAddress, value)
		cpu.registers.updateNegativeFlag(value)
		cpu.registers.updateZeroFlag(value)
	}
}

/*
	BCC  Branch on Carry Clear

	branch on C = 0                  N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BCC oper      90    2     2**
*/
func (cpu *Cpu6502) bcc(info defs.InfoStep) {
	if cpu.registers.carryFlag() == 0 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BCS  Branch on Carry Set

	branch on C = 1                 N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BCS oper      B0    2     2**
*/
func (cpu *Cpu6502) bcs(info defs.InfoStep) {
	if cpu.registers.carryFlag() == 1 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BEQ  Branch on Result Zero

	branch on Z = 1             N Z C I D V
								- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BEQ oper      F0    2     2**
*/
func (cpu *Cpu6502) beq(info defs.InfoStep) {
	if cpu.registers.zeroFlag() == 1 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BIT  Test Bits in Memory with Accumulator

	bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
	the zeroflag is set to the result of operand AND Accumulator.
	The result is not kept.

	A AND M, M7 -> N, M6 -> V        N Z C I D V
									M7 + - - - M6

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	zeropage      BIT oper      24    2     3
	Absolute      BIT oper      2C    3     4
*/
func (cpu *Cpu6502) bit(info defs.InfoStep) {
	value := cpu.Read(info.OperandAddress)
	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateFlag(overflowFlag, (value>>6)&0x01)
	cpu.registers.updateZeroFlag(value & cpu.registers.A)
}

/*
	BMI  Branch on Result Minus

	branch on N = 1                 N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BMI oper      30    2     2**
*/
func (cpu *Cpu6502) bmi(info defs.InfoStep) {
	if cpu.registers.negativeFlag() == 1 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BNE  Branch on Result not Zero

	branch on Z = 0                  N Z C I D V
									- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BNE oper      D0    2     2**
*/
func (cpu *Cpu6502) bne(info defs.InfoStep) {
	// CHeck how to negate a bit and apply it here
	//if !cpu.Registers.zeroFlag() == 1 {
	if cpu.registers.zeroFlag() == 0 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BPL  Branch on Result Plus

	branch on N = 0             N Z C I D V
								- - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BPL oper      10    2     2**
*/
func (cpu *Cpu6502) bpl(info defs.InfoStep) {
	//if !cpu.Registers.NegativeFlag {
	if cpu.registers.negativeFlag() == 0 {
		cpu.registers.Pc = info.OperandAddress
	}
}

/*
	BRK Force Break
	The BRK Instruction forces the generation of an interrupt request.
    The program counter and processor status are pushed on the stack then
    the IRQ interrupt vector at $FFFE/F is loaded into the PC and the break
    flag in the status set to one.

	interrupt,                       N Z C I D V
	push PC+2, push SR               - - - 1 - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       BRK           00    1     7
*/
func (cpu *Cpu6502) brk(info defs.InfoStep) {
	// Store PC in stack
	pc := cpu.registers.Pc
	cpu.pushStack(defs.HighNibble(pc))
	cpu.pushStack(defs.LowNibble(pc))

	// Push status with Break flag set
	cpu.pushStack(cpu.registers.Status | 0b00010000)

	cpu.registers.updateFlag(interruptFlag, 1)

	cpu.registers.Pc = defs.Address(cpu.read16(0xFFFE))
}

/*
	BVC  Branch on Overflow Clear
	branch on V = 0               N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BVC oper      50    2     2**
*/
func (cpu *Cpu6502) bvc(info defs.InfoStep) {
	if cpu.registers.overflowFlag() == byte(1) {
		return
	}

	cpu.registers.Pc = info.OperandAddress
}

/*
	BVS  Branch on Overflow Set
	branch on V = 1               N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Relative      BVC oper      70    2     2**
*/
func (cpu *Cpu6502) bvs(info defs.InfoStep) {
	if cpu.registers.overflowFlag() == 0 {
		return
	}

	cpu.registers.Pc = info.OperandAddress
}

/*
	CLC  Clear Carry Flag
	0 -> C                        N Z C I D V
								  - - 0 - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       CLC           18    1     2
*/
func (cpu *Cpu6502) clc(info defs.InfoStep) {
	cpu.registers.updateFlag(carryFlag, 0)
}

/*
	CLD  Clear Decimal Mode
	0 -> D                        N Z C I D V
								  - - - - 0 -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       CLD           D8    1     2
*/
func (cpu *Cpu6502) cld(info defs.InfoStep) {
	cpu.registers.updateFlag(decimalFlag, 0)
}

/*
	CLI  Clear Interrupt Disable Bit
	0 -> I                        N Z C I D V
								  - - - 0 - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       CLI           58    1     2
*/
func (cpu *Cpu6502) cli(info defs.InfoStep) {
	cpu.registers.updateFlag(interruptFlag, 0)
}

/*
	CLV  Clear Overflow Flag
	0 -> V                        N Z C I D V
								  - - - - - 0

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       CLV           B8    1     2
*/
func (cpu *Cpu6502) clv(info defs.InfoStep) {
	cpu.registers.updateFlag(overflowFlag, 0)
}

/*
	CMP (CoMPare Accumulator)

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

	Compare sets flags as if a subtraction had been carried out. If the value in the Accumulator is equal or greater than the compared value, the Carry will be set. The equal (Z) and sign (S) flags will be set based on equality or lack thereof and the sign (i.e. A>=$80) of the Accumulator.
*/
func (cpu *Cpu6502) cmp(info defs.InfoStep) {
	operand := cpu.Read(info.OperandAddress)
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
func (cpu *Cpu6502) cpx(info defs.InfoStep) {
	operand := cpu.Read(info.OperandAddress)
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
func (cpu *Cpu6502) cpy(info defs.InfoStep) {
	operand := cpu.Read(info.OperandAddress)
	cpu.compare(cpu.registers.Y, operand)
}

func (cpu *Cpu6502) compare(register byte, operand byte) {
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

func (cpu *Cpu6502) dec(info defs.InfoStep) {
	address := info.OperandAddress
	operand := cpu.Read(address)

	operand--
	cpu.write(address, operand)

	cpu.registers.updateZeroFlag(operand)
	//if operand == 0 {
	//	cpu.Registers.ZeroFlag = true
	//} else {
	//	cpu.Registers.ZeroFlag = false
	//}
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0xFF {
	//	cpu.Registers.updateFlag(negativeFlag, 1)
	//} else {
	//	cpu.Registers.updateFlag(negativeFlag, 0)
	//}
}

func (cpu *Cpu6502) dex(info defs.InfoStep) {
	cpu.registers.X--
	operand := cpu.registers.X

	cpu.registers.updateZeroFlag(operand)
	//if operand == 0 {
	//	cpu.Registers.ZeroFlag = true
	//} else {
	//	cpu.Registers.ZeroFlag = false
	//}
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0xFF {
	//	cpu.Registers.updateFlag(negativeFlag, 1)
	//} else {
	//	cpu.Registers.updateFlag(negativeFlag, )NegativeFlag = false
	//}
}

func (cpu *Cpu6502) dey(info defs.InfoStep) {
	operand := cpu.registers.Y

	operand--
	cpu.registers.Y = operand

	cpu.registers.updateZeroFlag(operand)
	cpu.registers.updateNegativeFlag(operand)
	//if operand == 0 {
	//	cpu.Registers.ZeroFlag = true
	//} else {
	//	cpu.Registers.ZeroFlag = false
	//}
	//
	//if operand == 0xFF {
	//	cpu.Registers.NegativeFlag = true
	//} else {
	//	cpu.Registers.NegativeFlag = false
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
func (cpu *Cpu6502) eor(info defs.InfoStep) {
	value := cpu.Read(info.OperandAddress)

	cpu.registers.A = cpu.registers.A ^ value
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	INC  Increment Memory by One
	M + 1 -> M                    N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	zeropage      INC oper      E6    2     5
	zeropage,X    INC oper,X    F6    2     6
	Absolute      INC oper      EE    3     6
	Absolute,X    INC oper,X    FE    3     7
*/
func (cpu *Cpu6502) inc(info defs.InfoStep) {
	value := cpu.Read(info.OperandAddress)
	value += 1

	cpu.write(info.OperandAddress, value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateNegativeFlag(value)
}

/*
	INX  Increment Index X by One
	X + 1 -> X                N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       INX           E8    1     2
*/
func (cpu *Cpu6502) inx(info defs.InfoStep) {
	cpu.registers.X += 1
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)

}

/*
	INY  Increment Index Y by One
	Y + 1 -> Y                    N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       INY           C8    1     2
*/
func (cpu *Cpu6502) iny(info defs.InfoStep) {
	cpu.registers.Y += 1
	cpu.registers.updateZeroFlag(cpu.registers.Y)
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
}

/*
	JMP  Jump to New Location
	(PC+1) -> PCL                    N Z C I D V
	(PC+2) -> PCH                    - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Absolute      JMP oper      4C    3     3
	Indirect      JMP (oper)    6C    3     5
*/
func (cpu *Cpu6502) jmp(info defs.InfoStep) {
	cpu.registers.Pc = info.OperandAddress
}

/*
	JSR  Jump to New Location Saving Return Address
	push (PC+2),                     N Z C I D V
	(PC+1) -> PCL                    - - - - - -
	(PC+2) -> PCH

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Absolute      JSR oper      20    3     6
*/
func (cpu *Cpu6502) jsr(info defs.InfoStep) {
	pc := cpu.registers.Pc - 1
	cpu.pushStack(byte(pc >> 8))
	cpu.pushStack(byte(pc & 0xFF))

	cpu.registers.Pc = info.OperandAddress
}

/*
	LDA  Load Accumulator with Memory
	M -> A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Immediate     LDA #oper     A9    2     2
	zeropage      LDA oper      A5    2     3
	zeropage,X    LDA oper,X    B5    2     4
	Absolute      LDA oper      AD    3     4
	Absolute,X    LDA oper,X    BD    3     4*
	Absolute,Y    LDA oper,Y    B9    3     4*
	(Indirect,X)  LDA (oper,X)  A1    2     6
	(Indirect),Y  LDA (oper),Y  B1    2     5*
*/
func (cpu *Cpu6502) lda(info defs.InfoStep) {
	cpu.registers.A = cpu.Read(info.OperandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	LDX  Load Index X with Memory
	M -> X                    N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Immediate     LDX #oper     A2    2     2
	zeropage      LDX oper      A6    2     3
	zeropage,Y    LDX oper,Y    B6    2     4
	Absolute      LDX oper      AE    3     4
	Absolute,Y    LDX oper,Y    BE    3     4*
*/
func (cpu *Cpu6502) ldx(info defs.InfoStep) {
	cpu.registers.X = cpu.Read(info.OperandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)
}

/*
	LDY  Load Index Y with Memory
	M -> Y                N Z C I D V
						  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	immidiate     LDY #oper     A0    2     2
	zeropage      LDY oper      A4    2     3
	zeropage,X    LDY oper,X    B4    2     4
	Absolute      LDY oper      AC    3     4
	Absolute,X    LDY oper,X    BC    3     4*
*/
func (cpu *Cpu6502) ldy(info defs.InfoStep) {
	cpu.registers.Y = cpu.Read(info.OperandAddress)
	cpu.registers.updateZeroFlag(cpu.registers.Y)
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
}

/*
	LSR  Shift One Bit Right (Memory or Accumulator)
	0 -> [76543210] -> C      N Z C I D V
							  0 + + - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Accumulator   LSR A         4A    1     2
	zeropage      LSR oper      46    2     5
	zeropage,X    LSR oper,X    56    2     6
	Absolute      LSR oper      4E    3     6
	Absolute,X    LSR oper,X    5E    3     7
*/
func (cpu *Cpu6502) lsr(info defs.InfoStep) {
	var value byte
	if info.AddressMode == defs.Implicit {
		value = cpu.registers.A
	} else {
		value = cpu.Read(info.OperandAddress)
	}

	//cpu.Registers.CarryFlag = value & 0x01
	cpu.registers.updateFlag(carryFlag, value&0x01)

	value >>= 1
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateNegativeFlag(0)

	if info.AddressMode == defs.Implicit {
		cpu.registers.A = value
	} else {
		cpu.write(info.OperandAddress, value)
	}
}

/*
	NOP  No Operation
	---                           N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       NOP           EA    1     2
*/
func (cpu *Cpu6502) nop(info defs.InfoStep) {

}

/*
	ORA  OR Memory with Accumulator
	A OR M -> A               N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	immidiate     ORA #oper     09    2     2
	zeropage      ORA oper      05    2     3
	zeropage,X    ORA oper,X    15    2     4
	Absolute      ORA oper      0D    3     4
	Absolute,X    ORA oper,X    1D    3     4*
	Absolute,Y    ORA oper,Y    19    3     4*
	(Indirect,X)  ORA (oper,X)  01    2     6
	(Indirect),Y  ORA (oper),Y  11    2     5*
*/
func (cpu *Cpu6502) ora(info defs.InfoStep) {
	value := cpu.Read(info.OperandAddress)
	cpu.registers.A |= value
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	PHA  Push Accumulator on Stack
	push A                        N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       PHA           48    1     3
*/
func (cpu *Cpu6502) pha(info defs.InfoStep) {
	cpu.pushStack(cpu.registers.A)
}

/*
	PHP  Push Processor Status on Stack
	push SR                       N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       PHP           08    1     3
*/
func (cpu *Cpu6502) php(info defs.InfoStep) {
	value := cpu.registers.statusRegister()
	value |= 0b00110000
	cpu.pushStack(value)
}

/*
	PLA  Pull Accumulator from Stack
	pull A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       PLA           68    1     4
*/
func (cpu *Cpu6502) pla(info defs.InfoStep) {
	cpu.registers.A = cpu.popStack()
	cpu.registers.updateNegativeFlag(cpu.registers.A)
	cpu.registers.updateZeroFlag(cpu.registers.A)
}

/*
	PLP  Pull Processor Status from Stack
	pull SR                       N Z C I D V
								  from stack

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       PLP           28    1     4
*/
func (cpu *Cpu6502) plp(info defs.InfoStep) {
	value := cpu.popStack()

	// From http://nesdev.com/the%20%27B%27%20flag%20&%20BRK%20instruction.txt
	// ...when the flags are restored (via PLP or RTI), the B bit is discarded.
	// From https://wiki.nesdev.com/w/index.php/Status_flags
	// ...two instructions (PLP and RTI) pull a byte from the stack and set all the flags.
	// They ignore bits 5 and 4.
	cpu.registers.loadStatusRegister(value)
}

/*
	ROL  Rotate One Bit Left (Memory or Accumulator)
	C <- [76543210] <- C          N Z C I D V
								  + + + - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Accumulator   ROL A         2A    1     2
	zeropage      ROL oper      26    2     5
	zeropage,X    ROL oper,X    36    2     6
	Absolute      ROL oper      2E    3     6
	Absolute,X    ROL oper,X    3E    3     7
*/
func (cpu *Cpu6502) rol(info defs.InfoStep) {
	var newCarry byte
	var value byte
	if info.AddressMode == defs.Implicit {
		newCarry = cpu.registers.A & 0x80 >> 7
		cpu.registers.A <<= 1
		cpu.registers.A |= cpu.registers.carryFlag()
		value = cpu.registers.A
	} else {
		value = cpu.Read(info.OperandAddress)
		newCarry = value & 0x80 >> 7
		value <<= 1
		value |= cpu.registers.carryFlag()
		cpu.write(info.OperandAddress, value)
	}

	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateFlag(carryFlag, newCarry)
}

/*
	ROR  Rotate One Bit Right (Memory or Accumulator)
	C -> [76543210] -> C          N Z C I D V
								  + + + - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	Accumulator   ROR A         6A    1     2
	zeropage      ROR oper      66    2     5
	zeropage,X    ROR oper,X    76    2     6
	Absolute      ROR oper      6E    3     6
	Absolute,X    ROR oper,X    7E    3     7
*/
func (cpu *Cpu6502) ror(info defs.InfoStep) {
	var newCarry byte
	var value byte
	if info.AddressMode == defs.Implicit {
		newCarry = cpu.registers.A & 0x01
		cpu.registers.A >>= 1
		cpu.registers.A |= cpu.registers.carryFlag() << 7
		value = cpu.registers.A
	} else {
		value = cpu.Read(info.OperandAddress)
		newCarry = value & 0x01
		value >>= 1
		value |= cpu.registers.carryFlag() << 7
		cpu.write(info.OperandAddress, value)
	}

	cpu.registers.updateNegativeFlag(value)
	cpu.registers.updateZeroFlag(value)
	cpu.registers.updateFlag(carryFlag, newCarry)
}

/*
	RTI  Return from Interrupt

	pull SR, pull PC              N Z C I D V
								  from stack

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       RTI           40    1     6
*/
func (cpu *Cpu6502) rti(info defs.InfoStep) {
	statusRegister := cpu.popStack()
	cpu.registers.loadStatusRegister(statusRegister)

	msb := cpu.popStack()
	lsb := cpu.popStack()
	cpu.registers.Pc = defs.CreateAddress(lsb, msb)
}

/*
	RTS  Return from Subroutine
	pull PC, PC+1 -> PC           N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       RTS           60    1     6
*/
func (cpu *Cpu6502) rts(info defs.InfoStep) {
	msb := cpu.popStack()
	lsb := cpu.popStack()
	cpu.registers.Pc = defs.CreateAddress(lsb, msb)
}

/*
	SBC  Subtract Memory from Accumulator with Borrow
	A - M - C -> A                N Z C I D V
								  + + + - - +

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	immidiate     SBC #oper     E9    2     2
	zeropage      SBC oper      E5    2     3
	zeropage,X    SBC oper,X    F5    2     4
	Absolute      SBC oper      ED    3     4
	Absolute,X    SBC oper,X    FD    3     4*
	Absolute,Y    SBC oper,Y    F9    3     4*
	(Indirect,X)  SBC (oper,X)  E1    2     6
	(Indirect),Y  SBC (oper),Y  F1    2     5*
*/
func (cpu *Cpu6502) sbc(info defs.InfoStep) {
	value := cpu.Read(info.OperandAddress)
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

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       SEC           38    1     2
*/
func (cpu *Cpu6502) sec(info defs.InfoStep) {
	cpu.registers.updateFlag(carryFlag, 1)
}

/*
	SED  Set Decimal Flag
	1 -> D                    N Z C I D V
							  - - - - 1 -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       SED           F8    1     2
*/
func (cpu *Cpu6502) sed(info defs.InfoStep) {
	cpu.registers.updateFlag(decimalFlag, 1)
}

func (cpu *Cpu6502) sei(info defs.InfoStep) {
	cpu.registers.updateFlag(interruptFlag, 1)
}

func (cpu *Cpu6502) sta(info defs.InfoStep) {
	cpu.write(info.OperandAddress, cpu.registers.A)
}

func (cpu *Cpu6502) stx(info defs.InfoStep) {
	cpu.write(info.OperandAddress, cpu.registers.X)
}

func (cpu *Cpu6502) sty(info defs.InfoStep) {
	cpu.write(info.OperandAddress, cpu.registers.Y)
}

/*
	TAX  Transfer Accumulator to Index X
	A -> X                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       TAX           AA    1     2
*/
func (cpu *Cpu6502) tax(info defs.InfoStep) {
	cpu.registers.X = cpu.registers.A
	cpu.registers.updateNegativeFlag(cpu.registers.X)
	cpu.registers.updateZeroFlag(cpu.registers.X)
}

/*
	TAY  Transfer Accumulator to Index Y
	A -> Y                    N Z C I D V
							  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       TAY           A8    1     2
*/
func (cpu *Cpu6502) tay(info defs.InfoStep) {
	cpu.registers.Y = cpu.registers.A
	cpu.registers.updateNegativeFlag(cpu.registers.Y)
	cpu.registers.updateZeroFlag(cpu.registers.Y)
}

/*
	TSX  Transfer Stack Pointer to Index X
	SP -> X                       N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       TSX           BA    1     2
*/
func (cpu *Cpu6502) tsx(info defs.InfoStep) {
	cpu.registers.X = cpu.Registers().Sp
	cpu.registers.updateZeroFlag(cpu.registers.X)
	cpu.registers.updateNegativeFlag(cpu.registers.X)
}

/*
	TXA  Transfer Index X to Accumulator
	X -> A                        N Z C I D V
								  + + - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       TXA           8A    1     2
*/
func (cpu *Cpu6502) txa(info defs.InfoStep) {
	cpu.registers.A = cpu.registers.X
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}

/*
	TXS  Transfer Index X to Stack Pointer
	X -> SP                       N Z C I D V
								  - - - - - -

	addressing    assembler    opc  bytes  cycles
	--------------------------------------------
	implied       TXS           9A    1     2
*/
func (cpu *Cpu6502) txs(info defs.InfoStep) {
	cpu.registers.Sp = cpu.registers.X
}

/*
	TYA  Transfer Index Y to Accumulator
	 Y -> A                           N Z C I D V
									  + + - - - -
	 addressing    assembler    opc  bytes  cycles
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
func (cpu *Cpu6502) tya(info defs.InfoStep) {
	cpu.registers.A = cpu.registers.Y
	cpu.registers.updateZeroFlag(cpu.registers.A)
	cpu.registers.updateNegativeFlag(cpu.registers.A)
}
