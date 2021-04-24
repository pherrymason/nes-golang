package nes

type NesDebugger struct {
	debug         bool
	cpu           *Cpu6502
	outputLogPath string

	disassembled map[Address]string
}

func (debugger NesDebugger) Disassembled() map[Address]string {
	return debugger.disassembled
}

func (debugger NesDebugger) ProgramCounter() Address {
	return debugger.cpu.ProgramCounter()
}
