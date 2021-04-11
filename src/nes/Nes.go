package nes

import (
	"github.com/raulferras/nes-golang/src/log"
	"github.com/raulferras/nes-golang/src/nes/component"
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
)

type Nes struct {
	cpu   *cpu2.Cpu6502
	bus   *component.Bus
	debug bool
}

func CreateNes() Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPU(&bus)

	nes := Nes{
		cpu:   &cpu,
		bus:   &bus,
		debug: false,
	}

	return nes
}

func CreateDebuggableNes(logger log.Logger) Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPUDebuggable(&bus, logger)

	nes := Nes{
		&cpu,
		&bus,
		true,
	}

	return nes
}

func (nes *Nes) Start() {
	nes.cpu.Init()
	//nes.cpu.reset()
	for {
		nes.cpu.Tick()
	}
}

func (nes *Nes) InsertCartridge(cartridge *component.GamePak) {
	nes.bus.AttachCartridge(cartridge)
}
