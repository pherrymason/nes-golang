package nes

import (
	"fmt"
	"github.com/FMNSSun/hexit"
	"github.com/raulferras/nes-golang/src/utils"
	"regexp"
	"strconv"
	"strings"
)

type CpuState struct {
	Registers          Cpu6502Registers
	CurrentInstruction Instruction
	RawOpcode          byte
	EvaluatedAddress   Address
	CyclesSinceReset   uint32
}

func CreateState(registers Cpu6502Registers, opcode byte, instruction Instruction, step OperationMethodArgument, cpuCycle uint32) CpuState {
	state := CpuState{
		registers,
		instruction,
		opcode,
		step.OperandAddress,
		cpuCycle,
	}

	return state
}

func CreateStateFromNesTestLine(nesTestLine string) CpuState {
	fields := strings.Fields(nesTestLine)
	_ = fields

	blocks := utils.StringSplitByRegex(nesTestLine)

	result := utils.HexStringToByteArray(blocks[0])
	pc := CreateAddress(result[1], result[0])

	fields = strings.Fields(blocks[1])
	opcode := utils.HexStringToByteArray(fields[0])

	flagFields := strings.Fields(blocks[3])

	r, _ := regexp.Compile("CYC:([0-9]+)$")
	cpuCyclesString := r.FindStringSubmatch(nesTestLine)

	cpuCycles, _ := strconv.ParseUint(cpuCyclesString[1], 10, 16)

	state := CpuState{
		Cpu6502Registers{
			utils.NestestDecodeRegisterFlag(flagFields[0]),
			utils.NestestDecodeRegisterFlag(flagFields[1]),
			utils.NestestDecodeRegisterFlag(flagFields[2]),
			pc,
			utils.NestestDecodeRegisterFlag(flagFields[4]),
			utils.NestestDecodeRegisterFlag(flagFields[3]),
		},
		CreateInstruction(
			strings.Fields(blocks[2])[0],
			Implicit,
			nil,
			0,
			0,
		),
		opcode[0],
		CreateAddress(0x00, 0x00),
		uint32(cpuCycles),
	}

	return state
}

func (state *CpuState) String() string {
	var msg strings.Builder
	msg.Grow(150)

	// Pointer
	msg.Write(hexit.HexUint16(uint16(state.Registers.Pc)))
	msg.WriteString(" ")

	// Raw OPCode + Operand
	msg.Write(hexit.HexUint8(state.RawOpcode))

	//clampSpace(&msg, 16)
	msg.WriteString(state.CurrentInstruction.Name())
	msg.WriteString(" ")

	if state.CurrentInstruction.AddressMode() == Immediate {
		msg.WriteString("#$%02X")
	} else {
		msg.WriteString("$%02X")
	}
	msg.Write(hexit.HexUint16(uint16(state.EvaluatedAddress)))

	//msg = clampSpace(msg, 48)
	msg.WriteByte(' ')
	msg.WriteString("A:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.A))
	msg.WriteString("X:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.X))
	msg.WriteString("Y:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Y))
	msg.WriteString("P:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Status))
	msg.WriteString("SP:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Sp))
	msg.WriteString("PPU:___,___")
	msg.WriteString("CYC:")
	msg.WriteString(strconv.Itoa(int(state.Registers.Sp)))
	msg.WriteString("\n")

	/*
		msg.WriteString(
			fmt.Sprintf(
				"A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:___,___ CYC:%d",
				state.Registers.A,
				state.Registers.X,
				state.Registers.Y,
				state.Registers.Status,
				state.Registers.Sp,
				state.CyclesSinceReset,
			))
	*/
	return msg.String()
}

func (state CpuState) Equals(b CpuState) bool {
	if state.Registers.Pc != b.Registers.Pc ||
		state.Registers.A != b.Registers.A ||
		state.Registers.X != b.Registers.X ||
		state.Registers.Y != b.Registers.Y ||
		state.Registers.Sp != b.Registers.Sp ||
		state.Registers.statusRegister() != b.Registers.statusRegister() ||
		state.CyclesSinceReset != b.CyclesSinceReset {
		return false
	}

	return true
}

func (state CpuState) ToString() string {
	msg := fmt.Sprintf("%X ", state.Registers.Pc)
	msg += fmt.Sprintf("%X ", state.RawOpcode)
	msg += fmt.Sprintf("%s - ", state.CurrentInstruction.Name())
	msg += fmt.Sprintf(
		"A:%X X:%X Y:%X P:%X SP:%X ",
		state.Registers.A,
		state.Registers.X,
		state.Registers.Y,
		state.Registers.statusRegister(),
		state.Registers.Sp,
	)
	// Cycles
	msg += fmt.Sprintf("PPU:%d,%d CYC:%d", 0, 0, state.CyclesSinceReset)

	return msg
}

func clampSpace(msg *strings.Builder, clamp int) {
	for i := msg.Len(); i <= clamp; i++ {
		msg.WriteString(" ")
	}
}
