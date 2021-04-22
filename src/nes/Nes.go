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

func CreateNes() Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPU(&bus)

	nes := Nes{
		cpu:   cpu,
		bus:   &bus,
		debug: NesDebugger{false, cpu, nil, 0, nil},
	}

	return nes
}

func CreateDebuggableNes(debugger NesDebugger) Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPUDebuggable(&bus, debugger.logger)

	nes := Nes{
		cpu,
		&bus,
		NesDebugger{
			true,
			cpu,
			debugger.logger,
			debugger.cyclesLimit,
			nil,
		},
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
	return nes.cpu.Tick()
}

func (nes *Nes) InsertGamePak(cartridge *component.GamePak) {
	nes.bus.InsertGamePak(cartridge)
}

func (nes Nes) Debugger() *NesDebugger {
	return &nes.debug
}
