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
	debug       bool
	cpu         *Cpu6502
	ppu         *ppu.Ppu2c02
	logPath     string
	maxCPUCycle int64

	disassembled map[types.Address]string
	DebugPPU     bool
}

func CreateNesDebugger(logPath string, debug bool, debugPPU bool, maxCPUCycle int64) *NesDebugger {
	return &NesDebugger{
		debug:        debug,
		cpu:          nil,
		ppu:          nil,
		logPath:      logPath,
		maxCPUCycle:  maxCPUCycle,
		disassembled: nil,
		DebugPPU:     debugPPU,
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

func (debugger *NesDebugger) PatternTable(patternTable byte, palette uint8) image.RGBA {
	return debugger.ppu.PatternTable(patternTable, palette)
}

func (debugger *NesDebugger) GetPaletteFromRam(paletteIndex uint8) [4]color.Color {
	var colors [4]color.Color

	colors[0] = debugger.ppu.GetRGBColor(paletteIndex, 0)
	colors[1] = debugger.ppu.GetRGBColor(paletteIndex, 1)
	colors[2] = debugger.ppu.GetRGBColor(paletteIndex, 2)
	colors[3] = debugger.ppu.GetRGBColor(paletteIndex, 3)

	return colors
}

func (debugger *NesDebugger) GetPaletteColorFromPaletteRam(paletteIndex byte, colorIndex byte) byte {
	return debugger.ppu.GetPaletteColor(paletteIndex, colorIndex)
}

func (debugger *NesDebugger) OAM(index byte) []byte {
	return debugger.ppu.Oam(index)
}
