package nes

import (
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
	"image/color"
)

// NesDebugger offers an api to interact externally with
// stuff from Nes
type NesDebugger struct {
	debug         bool
	cpu           *Cpu6502
	ppu           *ppu.Ppu2c02
	outputLogPath string

	disassembled map[types.Address]string
}

func CreateNesDebugger(logPath string, debug bool) *NesDebugger {
	return &NesDebugger{
		debug,
		nil,
		nil,
		logPath,
		nil,
	}
}

func (debugger *NesDebugger) Disassembled() map[types.Address]string {
	return debugger.disassembled
}

func (debugger *NesDebugger) ProgramCounter() types.Address {
	return debugger.cpu.ProgramCounter()
}

func (debugger *NesDebugger) N() bool {
	bit := debugger.cpu.Registers().NegativeFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) O() bool {
	bit := debugger.cpu.Registers().OverflowFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) B() bool {
	bit := debugger.cpu.Registers().BreakFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) D() bool {
	bit := debugger.cpu.Registers().DecimalFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) I() bool {
	bit := debugger.cpu.Registers().InterruptFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) Z() bool {
	bit := debugger.cpu.Registers().ZeroFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *NesDebugger) C() bool {
	bit := debugger.cpu.Registers().CarryFlag()
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

func (debugger *NesDebugger) PatternTable(patternTable byte) image.RGBA {
	return debugger.ppu.PatternTable(patternTable, 0)
}

func (debugger *NesDebugger) GetPaletteFromRam(paletteIndex uint8) [3]color.Color {
	var colors [3]color.Color

	colors[0] = debugger.ppu.GetColorFromPaletteRam(paletteIndex, 0)
	colors[1] = debugger.ppu.GetColorFromPaletteRam(paletteIndex, 1)
	colors[2] = debugger.ppu.GetColorFromPaletteRam(paletteIndex, 2)

	return colors
}
