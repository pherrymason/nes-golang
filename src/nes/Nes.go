package nes

import (
	"github.com/FMNSSun/hexit"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type Nes struct {
	cpu *Cpu6502
	ppu *Ppu2c02

	systemClockCounter byte // Controls how many times to call each processor
	debug              *NesDebugger
}

func CreateNes(gamePak *gamePak.GamePak, debugger *NesDebugger) Nes {
	hexit.BuildTable()
	ppuBus := CreatePPUMemory(gamePak)
	ppu := CreatePPU(
		ppuBus,
	)

	cpuBus := newNESCPUMemory(ppu, gamePak)
	cpu := CreateCPU(
		cpuBus,
		Cpu6502DebugOptions{debugger.debug, debugger.outputLogPath},
	)
	debugger.cpu = cpu
	debugger.ppu = ppu

	nes := Nes{
		cpu:   cpu,
		ppu:   ppu,
		debug: debugger,
	}

	return nes
}

func (nes *Nes) StartAt(address types.Address) {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.ResetToAddress(address)
}

// Rename to PowerOn
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

	if nes.ppu.nmi {
		nes.cpu.nmi()
		nes.ppu.nmi = false
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
