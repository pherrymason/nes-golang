package cpu

import (
	"fmt"
	"github.com/FMNSSun/hexit"
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"strconv"
	"strings"
)

type CpuState struct {
	Registers          Registers
	CurrentInstruction Instruction
	RawOpcode          [3]byte
	EvaluatedAddress   types.Address
	CyclesSinceReset   uint32
	Waiting            bool
}

func CreateWaitingState() CpuState {
	return CpuState{Waiting: true}
}

func CreateState(registers Registers, opcode [3]byte, instruction Instruction, step OperationMethodArgument, cpuCycle uint32) CpuState {
	state := CpuState{
		registers,
		instruction,
		opcode,
		step.OperandAddress,
		cpuCycle,
		false,
	}

	return state
}

func (state *CpuState) String(ppuState ppu.SimplePPUState) string {
	var msg strings.Builder
	msg.Grow(150)

	// Pointer
	msg.Write(hexit.HexUint16(uint16(state.Registers.Pc)))
	msg.WriteString(" ")

	// Raw OPCode + Operand
	msg.Write(hexit.HexUint8(state.RawOpcode[0]))
	msg.WriteString(" ")
	msg.Write(hexit.HexUint8(state.RawOpcode[1]))
	msg.WriteString(" ")
	msg.Write(hexit.HexUint8(state.RawOpcode[2]))
	msg.WriteString(" ")

	//clampSpace(&msg, 16)
	msg.WriteString(state.CurrentInstruction.Name())
	msg.WriteString(" ")

	if state.CurrentInstruction.AddressMode() == Immediate {
		msg.WriteString("#$")
	} else {
		msg.WriteString("$")
	}
	msg.Write(hexit.HexUint16(uint16(state.EvaluatedAddress)))

	//msg = clampSpace(msg, 48)
	if state.CurrentInstruction.AddressMode() == Immediate {
		msg.WriteString("       ")
	} else {
		msg.WriteString("        ")
	}
	msg.WriteString("A:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.A))
	msg.WriteString(" X:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.X))
	msg.WriteString(" Y:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Y))
	msg.WriteString(" P:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Status))
	msg.WriteString(" SP:")
	msg.WriteString(hexit.HexUint8Str(state.Registers.Sp))
	msg.WriteString(" PPU: ")
	msg.WriteString(strconv.Itoa(int(ppuState.Scanline)))
	msg.WriteString(",")
	msg.WriteString(strconv.Itoa(int(ppuState.RenderCycle)))

	spaces := 7
	if ppuState.Scanline < 10 {
		spaces -= 1
	} else if ppuState.Scanline < 100 {
		spaces -= 2
	} else {
		spaces -= 3
	}
	if ppuState.RenderCycle < 10 {
		spaces -= 1
	} else if ppuState.RenderCycle < 100 {
		spaces -= 2
	} else {
		spaces -= 3
	}
	for s := 0; s < spaces; s++ {
		msg.WriteString(" ")
	}
	msg.WriteString(" CYC:")
	msg.WriteString(strconv.Itoa(int(state.CyclesSinceReset)))

	// Frame number
	msg.WriteString(" f")
	msg.WriteString(strconv.Itoa(int(ppuState.Frame + 1)))

	msg.WriteString("\n")

	return msg.String()
}

func (state CpuState) RegistersEquals(b CpuState) bool {
	if state.Registers.Pc != b.Registers.Pc ||
		state.Registers.A != b.Registers.A ||
		state.Registers.X != b.Registers.X ||
		state.Registers.Y != b.Registers.Y ||
		state.Registers.Sp != b.Registers.Sp ||
		state.Registers.StatusRegister() != b.Registers.StatusRegister() ||
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
		state.Registers.StatusRegister(),
		state.Registers.Sp,
	)
	// Cycles
	msg += fmt.Sprintf("PPU:%d,%d CYC:%d", 0, 0, state.CyclesSinceReset)

	return msg
}
