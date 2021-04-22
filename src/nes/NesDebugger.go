package nes

import (
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

type NesDebugger struct {
	debug        bool
	cpu          *cpu2.Cpu6502
	logger       *cpu2.Logger
	cyclesLimit  uint16
	disassembled map[defs.Address]string
}

func (debugger NesDebugger) Disassembled() map[defs.Address]string {
	return debugger.disassembled
}

func (debugger NesDebugger) ProgramCounter() defs.Address {
	return debugger.cpu.Registers().Pc
}
