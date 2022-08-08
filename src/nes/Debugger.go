package nes

import (
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
	"image/color"
)

// Debugger offers an api to interact externally with
// NES components
type Debugger struct {
	cpu          *Cpu6502
	ppu          *ppu.Ppu2c02
	logPath      string
	maxCPUCycle  int64
	disassembled map[types.Address]string
	DebugPPU     bool
	debugCPU     bool
	// debugging related
	cpuBreakPoints    map[types.Address]bool
	cpuStepByStepMode bool
}

func CreateNesDebugger(logPath string, debugCPU bool, debugPPU bool, maxCPUCycle int64) *Debugger {
	return &Debugger{
		cpu:          nil,
		ppu:          nil,
		logPath:      logPath,
		maxCPUCycle:  maxCPUCycle,
		disassembled: nil,
		debugCPU:     debugCPU,
		DebugPPU:     debugPPU,

		cpuBreakPoints: make(map[types.Address]bool),
	}
}

func (debugger *Debugger) AddBreakPoint(address types.Address) {
	debugger.cpuBreakPoints[address] = true
}

func (debugger *Debugger) RemoveBreakPoint(address types.Address) {
	debugger.cpuBreakPoints[address] = false
}

func (debugger *Debugger) shouldPauseBecauseBreakpoint() bool {
	pc := debugger.cpu.ProgramCounter()
	enabled, exist := debugger.cpuBreakPoints[pc]
	if enabled && exist {
		debugger.cpuStepByStepMode = true
		return true
	}

	return false
}

func (debugger *Debugger) Disassembled() map[types.Address]string {
	return debugger.disassembled
}

func (debugger *Debugger) ProgramCounter() types.Address {
	return debugger.cpu.ProgramCounter()
}

func (debugger *Debugger) N() bool {
	bit := debugger.cpu.Registers().NegativeFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) O() bool {
	bit := debugger.cpu.Registers().OverflowFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) B() bool {
	bit := debugger.cpu.Registers().BreakFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) D() bool {
	bit := debugger.cpu.Registers().DecimalFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) I() bool {
	bit := debugger.cpu.Registers().InterruptFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) Z() bool {
	bit := debugger.cpu.Registers().ZeroFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) C() bool {
	bit := debugger.cpu.Registers().CarryFlag()
	if bit == 1 {
		return true
	}

	return false
}

func (debugger *Debugger) ARegister() byte {
	return debugger.cpu.registers.A
}

func (debugger *Debugger) XRegister() byte {
	return debugger.cpu.registers.X
}

func (debugger *Debugger) YRegister() byte {
	return debugger.cpu.registers.Y
}

// PPU Related
func (debugger *Debugger) PatternTable(patternTable byte, palette uint8) image.RGBA {
	return debugger.ppu.PatternTable(patternTable, palette)
}

func (debugger *Debugger) GetPaletteFromRam(paletteIndex uint8) [4]color.Color {
	var colors [4]color.Color

	colors[0] = debugger.ppu.GetRGBColor(paletteIndex, 0)
	colors[1] = debugger.ppu.GetRGBColor(paletteIndex, 1)
	colors[2] = debugger.ppu.GetRGBColor(paletteIndex, 2)
	colors[3] = debugger.ppu.GetRGBColor(paletteIndex, 3)

	return colors
}

func (debugger *Debugger) GetPaletteColorFromPaletteRam(paletteIndex byte, colorIndex byte) byte {
	return debugger.ppu.GetPaletteColor(paletteIndex, colorIndex)
}

func (debugger *Debugger) OAM(index byte) []byte {
	return debugger.ppu.Oam(index)
}
