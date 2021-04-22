package cpu

import (
	"github.com/raulferras/nes-golang/src/nes/defs"
	"strings"
)

func hex(n uint32, d int) string {
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

func (cpu *Cpu6502) Disassemble(start defs.Address, end defs.Address) map[defs.Address]string {
	disassembledCode := make(map[defs.Address]string)
	addr := uint32(start)
	value := byte(0x00)
	lo := byte(0x00)
	hi := byte(0x00)

	for addr <= uint32(end) {
		lineAddr := addr

		// Prefix line with instruction address
		sInst := "$" + hex(uint32(addr), 4) + ": "

		// Read instruction, and get its readable name
		opcode := cpu.bus.Read(defs.Address(addr))
		addr++
		sInst += cpu.instructions[opcode].Name() + " "

		if cpu.instructions[opcode].AddressMode() == defs.Implicit {
			sInst += " {IMP}"
		} else if cpu.instructions[opcode].AddressMode() == defs.Immediate {
			value = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "#$" + hex(uint32(value), 2) + " {IMM}"
		} else if cpu.instructions[opcode].AddressMode() == defs.ZeroPage {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = 0x00
			sInst += "$" + hex(uint32(lo), 2) + " {ZP0}"
		} else if cpu.instructions[opcode].AddressMode() == defs.ZeroPageX {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = 0x00
			sInst += "$" + hex(uint32(lo), 2) + ", X {ZPX}"
		} else if cpu.instructions[opcode].AddressMode() == defs.ZeroPageY {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = 0x00
			sInst += "$" + hex(uint32(lo), 2) + ", Y {ZPY}"
		} else if cpu.instructions[opcode].AddressMode() == defs.IndirectX {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = 0x00
			sInst += "($" + hex(uint32(lo), 2) + ", X) {IZX}"
		} else if cpu.instructions[opcode].AddressMode() == defs.IndirectY {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = 0x00
			sInst += "($" + hex(uint32(lo), 2) + "), Y {IZY}"
		} else if cpu.instructions[opcode].AddressMode() == defs.Absolute {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "$" + hex(uint32(defs.CreateWord(hi, lo)), 4) + " {ABS}"
		} else if cpu.instructions[opcode].AddressMode() == defs.AbsoluteXIndexed {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "$" + hex(uint32(defs.CreateWord(hi, lo)), 4) + ", X {ABX}"
		} else if cpu.instructions[opcode].AddressMode() == defs.AbsoluteYIndexed {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "$" + hex(uint32(defs.CreateWord(hi, lo)), 4) + ", Y {ABY}"
		} else if cpu.instructions[opcode].AddressMode() == defs.Indirect {
			lo = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			hi = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "($" + hex(uint32(defs.CreateWord(hi, lo)), 4) + ") {IND}"
		} else if cpu.instructions[opcode].AddressMode() == defs.Relative {
			value = cpu.bus.ReadOnly(defs.Address(addr))
			addr++
			sInst += "$" + hex(uint32(value), 2) + " [$" + hex(addr+uint32(value), 4) + "] {REL}"
		}

		disassembledCode[defs.Address(lineAddr)] = sInst
	}

	return disassembledCode
}
