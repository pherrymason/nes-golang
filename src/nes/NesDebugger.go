package nes

import "github.com/raulferras/nes-golang/src/graphics"

// NesDebugger offers an api to interact externally with
// stuff from Nes
type NesDebugger struct {
	debug         bool
	cpu           *Cpu6502
	ppu           *Ppu2c02
	outputLogPath string

	disassembled map[Address]string
}

func CreateNesDebugger(logPath string) *NesDebugger {
	return &NesDebugger{
		false,
		nil,
		nil,
		logPath,
		nil,
	}
}

func (debugger *NesDebugger) Disassembled() map[Address]string {
	return debugger.disassembled
}

func (debugger *NesDebugger) ProgramCounter() Address {
	return debugger.cpu.ProgramCounter()
}

func (debugger *NesDebugger) N() bool {
	bit := debugger.cpu.Registers().negativeFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) O() bool {
	bit := debugger.cpu.Registers().overflowFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) B() bool {
	bit := debugger.cpu.Registers().breakFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) D() bool {
	bit := debugger.cpu.Registers().decimalFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) I() bool {
	bit := debugger.cpu.Registers().interruptFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) Z() bool {
	bit := debugger.cpu.Registers().zeroFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) C() bool {
	bit := debugger.cpu.Registers().carryFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) ARegister() byte {
	return debugger.cpu.registers.A
}

func (debugger *NesDebugger) XRegister() byte {
	return debugger.cpu.registers.X
}

func (debugger *NesDebugger) YRegister() byte {
	return debugger.cpu.registers.Y
}

//func (debugger NesDebugger) PatternTable(patternTable int) [][]byte {
func (debugger *NesDebugger) PatternTable(patternTable int) []graphics.Pixel {
	return debugger.ppu.PatternTable(patternTable, 0)
}

func (debugger *NesDebugger) GetPaletteFromRam(paletteIndex uint8) [3]graphics.Color {
	var colors [3]graphics.Color

	colors[0] = debugger.ppu.getColorFromPaletteRam(byte(paletteIndex), 0)
	colors[1] = debugger.ppu.getColorFromPaletteRam(byte(paletteIndex), 1)
	colors[2] = debugger.ppu.getColorFromPaletteRam(byte(paletteIndex), 2)

	return colors
}
