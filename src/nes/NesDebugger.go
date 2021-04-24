package nes

type NesDebugger struct {
	debug         bool
	cpu           *Cpu6502
	ppu           *Ppu2c02
	outputLogPath string

	disassembled map[Address]string
}

func (debugger NesDebugger) Disassembled() map[Address]string {
	return debugger.disassembled
}

func (debugger NesDebugger) ProgramCounter() Address {
	return debugger.cpu.ProgramCounter()
}

//func (debugger NesDebugger) PatternTable(patternTable int) [][]byte {
func (debugger NesDebugger) PatternTable(patternTable int) []Pixel {
	return debugger.ppu.PatternTable(patternTable)
}
