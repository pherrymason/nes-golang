package nes

import (
	"github.com/raulferras/nes-golang/src/nes/component"
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

type Nes struct {
	cpu   *cpu2.Cpu6502
	bus   *component.Bus
	debug NesDebugger
}

func CreateNes(debugger NesDebugger) Nes {
	ram := component.RAM{}
	cpu := cpu2.CreateCPU()

	bus := component.CreateBus(&ram)
	bus.ConnectCPU(cpu)

	nes := Nes{
		cpu:   cpu,
		bus:   bus,
		debug: debugger,
	}

	return nes
}

func (nes *Nes) StartAt(address defs.Address) {
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.ResetToAddress(address)
}

func (nes *Nes) Start() {
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.Reset()
}

func (nes *Nes) Tick() byte {
	return nes.bus.Tick()
}

func (nes *Nes) InsertGamePak(cartridge *component.GamePak) {
	nes.bus.InsertGamePak(cartridge)
}

func (nes Nes) Debugger() *NesDebugger {
	return &nes.debug
}
