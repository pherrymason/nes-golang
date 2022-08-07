package nes

import (
	"github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"strings"
)

func myHex(n types.Word, d int) string {
	s := strings.Repeat("0", d)
	i := d - 1
	for i >= 0 {
		c := "0123456789ABCDEF"[n&0xF]
		s = s[:i] + string(c) + s[i+1:]
		i--
		n >>= 4
	}

	return s
}

func (cpu6502 *Cpu6502) Disassemble(start types.Address, end types.Address) map[types.Address]string {
	disassembledCode := make(map[types.Address]string)
	addr := start
	value := byte(0x00)
	lo := byte(0x00)
	hi := byte(0x00)
	if end == 0xFFFF {
		end = 0x0000
	}

	for addr >= start {
		lineAddr := addr

		// Prefix line with instruction address
		sInst := "$" + myHex(addr, 4) + ": "

		// Read instruction, and get its readable name
		opcode := cpu6502.memory.Peek(addr)
		addr++
		instruction := cpu6502.instructions[opcode]

		if len(instruction.Name()) == 0 {
			sInst += "0x" + myHex(types.Word(opcode), 2) + "? "
		} else {
			sInst += instruction.Name() + " "
		}

		if instruction.AddressMode() == cpu.Implicit {
			sInst += " {IMP}"
		} else if instruction.AddressMode() == cpu.Immediate {
			value = cpu6502.memory.Peek(addr)
			addr++
			sInst += "#$" + myHex(types.Word(value), 2) + " {IMM}"
		} else if instruction.AddressMode() == cpu.ZeroPage {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = 0x00
			sInst += "$" + myHex(types.Word(lo), 2) + " {ZP0}"
		} else if instruction.AddressMode() == cpu.ZeroPageX {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = 0x00
			sInst += "$" + myHex(types.Word(lo), 2) + ", X {ZPX}"
		} else if instruction.AddressMode() == cpu.ZeroPageY {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = 0x00
			sInst += "$" + myHex(types.Word(lo), 2) + ", Y {ZPY}"
		} else if instruction.AddressMode() == cpu.IndirectX {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = 0x00
			sInst += "($" + myHex(types.Word(lo), 2) + ", X) {IZX}"
		} else if instruction.AddressMode() == cpu.IndirectY {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = 0x00
			sInst += "($" + myHex(types.Word(lo), 2) + "), Y {IZY}"
		} else if instruction.AddressMode() == cpu.Absolute {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = cpu6502.memory.Peek(addr)
			addr++
			sInst += "$" + myHex(types.CreateWord(lo, hi), 4) + " {ABS}"
		} else if instruction.AddressMode() == cpu.AbsoluteXIndexed {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = cpu6502.memory.Peek(addr)
			addr++
			sInst += "$" + myHex(types.CreateWord(lo, hi), 4) + ", X {ABX}"
		} else if instruction.AddressMode() == cpu.AbsoluteYIndexed {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = cpu6502.memory.Peek(addr)
			addr++
			sInst += "$" + myHex(types.CreateWord(lo, hi), 4) + ", Y {ABY}"
		} else if instruction.AddressMode() == cpu.Indirect {
			lo = cpu6502.memory.Peek(addr)
			addr++
			hi = cpu6502.memory.Peek(addr)
			addr++
			sInst += "($" + myHex(types.CreateWord(lo, hi), 4) + ") {IND}"
		} else if instruction.AddressMode() == cpu.Relative {
			value = cpu6502.memory.Peek(addr)
			addr++
			sInst += "$" + myHex(types.Word(value), 2) + " [$" + myHex(addr+types.Word(value), 4) + "] {REL}"
		}

		disassembledCode[lineAddr] = sInst
	}

	return disassembledCode
}

func (cpu6502 *Cpu6502) GetOperation(operation byte) cpu.Instruction {
	return cpu6502.instructions[operation]
}
