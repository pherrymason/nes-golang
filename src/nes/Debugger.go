package nes

import (
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/raulferras/nes-golang/src/utils"
	"image"
	"image/color"
	"log"
)

// Debugger offers an api to interact externally with
// NES components
type Debugger struct {
	cpu                *Cpu6502
	ppu                *ppu.P2c02
	logPath            string
	disassembled       map[types.Address]string
	sortedDisassembled []utils.ASM

	pauseEmulation func()

	DebugPPU bool
	debugCPU bool
	// debugging related
	cpuBreakPoints                  map[types.Address]bool
	cpuBreakPointTriggeredAt        types.Address // flag breakpoint as used
	cpuStepByStepMode               bool
	waitingNextCPUOperationFinishes bool
}

func CreateNesDebugger(logPath string, debugCPU bool, debugPPU bool) *Debugger {
	return &Debugger{
		cpu:          nil,
		ppu:          nil,
		logPath:      logPath,
		disassembled: nil,
		debugCPU:     debugCPU,
		DebugPPU:     debugPPU,

		pauseEmulation: nil,
		cpuBreakPoints: make(map[types.Address]bool),
	}
}

// Control flow related ---------------------------

func (debugger *Debugger) AddBreakPoint(address types.Address) {
	debugger.cpuBreakPoints[address] = true
}

func (debugger *Debugger) RemoveBreakPoint(address types.Address) {
	debugger.cpuBreakPoints[address] = false
}

func (debugger *Debugger) shouldPauseEmulation() bool {
	if debugger.waitingNextCPUOperationFinishes {
		log.Printf("Allow one Cpu instruction")
		return false
	}

	if debugger.isBreakpointTriggered() {
		debugger.pauseEmulation()
		return true
	}

	if debugger.isManualStepMode() {
		return true
	} else {
	}

	return false
}

func (debugger *Debugger) isBreakpointTriggered() bool {
	pc := debugger.cpu.ProgramCounter()
	if pc == debugger.cpuBreakPointTriggeredAt {
		return false
	}

	enabled, exist := debugger.cpuBreakPoints[pc]
	if enabled && exist {
		log.Printf("Breakpoint reached")
		debugger.cpuStepByStepMode = true
		debugger.cpuBreakPointTriggeredAt = pc
		return true
	}

	return false
}

func (debugger *Debugger) isManualStepMode() bool {
	return debugger.cpuStepByStepMode
}

func (debugger *Debugger) resumeFromBreakpoint() {
	debugger.cpuStepByStepMode = false
}

// RunOneCPUOperationAndPause executed when user wants to just run one single cycle after having
// emulation paused.
func (debugger *Debugger) RunOneCPUOperationAndPause() {
	debugger.waitingNextCPUOperationFinishes = true
}

// oneCpuOperationRan This method is expected to be called after the CPU runs a cycle.
// This is intended to reenable the emulation pause automatically after next cpu instruction has been called.
func (debugger *Debugger) oneCpuOperationRan() {
	debugger.waitingNextCPUOperationFinishes = false
}

// debugging control flow related ^---------------------------

func (debugger *Debugger) Disassembled() map[types.Address]string {
	return debugger.disassembled
}

func (debugger *Debugger) SortedDisassembled() []utils.ASM {
	return debugger.sortedDisassembled
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
