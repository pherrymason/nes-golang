package nes

import (
	"github.com/FMNSSun/hexit"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type Nes struct {
	cpu *Cpu6502
	ppu *ppu.Ppu2c02

	systemClockCounter byte // Controls how many times to call each processor
	debug              *NesDebugger
	vBlankCount        byte
}

func CreateNes(gamePak *gamePak.GamePak, debugger *NesDebugger) Nes {
	hexit.BuildTable()
	thePPU := ppu.CreatePPU(gamePak)

	cpuBus := newNESCPUMemory(thePPU, gamePak)
	cpu := CreateCPU(
		cpuBus,
		Cpu6502DebugOptions{debugger.debug, debugger.outputLogPath},
	)
	debugger.cpu = cpu
	debugger.ppu = thePPU

	nes := Nes{
		cpu:   cpu,
		ppu:   thePPU,
		debug: debugger,
	}

	return nes
}

func (nes *Nes) StartAt(address types.Address) {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.ResetToAddress(address)
}

// Start todo Rename to PowerOn
func (nes *Nes) Start() {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.Reset()
}

func (nes *Nes) Tick() byte {
	nes.ppu.Tick()

	cpuCycles := byte(0)
	if nes.systemClockCounter%3 == 0 {
		cpuCycles = nes.cpu.Tick()
	}

	if nes.ppu.Nmi() {
		nes.cpu.nmi()
		nes.ppu.ResetNmi()
	}

	if nes.ppu.VBlank() {
		if nes.vBlankCount == 60 {
			//nes.ppu.Render()
		}
		nes.vBlankCount++
	}

	nes.systemClockCounter++

	return cpuCycles
}

func (nes *Nes) TickForTime(seconds float64) {
	cycles := int(1789773 * seconds)
	for cycles > 0 {
		cycles -= int(nes.Tick())
	}
}

func (nes *Nes) Stop() {
	nes.cpu.Stop()
}

func (nes *Nes) Debugger() *NesDebugger {
	return nes.debug
}

func (nes Nes) SystemClockCounter() byte {
	return nes.systemClockCounter
}

func (nes Nes) Frame() *types.Frame {
	return nes.ppu.Frame()
}
func (nes Nes) FramePattern() *[1024]byte {
	return nes.ppu.FramePattern()
}
