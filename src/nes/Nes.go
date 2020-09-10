package nes

import "github.com/raulferras/nes-golang/src/log"

type Nes struct {
	cpu   *CPU
	bus   *Bus
	debug bool
}

func CreateNes() Nes {
	ram := RAM{}
	bus := CreateBus(&ram)

	cpu := CreateCPU(&bus)

	nes := Nes{
		cpu:   &cpu,
		bus:   &bus,
		debug: false,
	}

	return nes
}

func CreateDebugableNes(logger log.Logger) Nes {
	ram := RAM{}
	bus := CreateBus(&ram)

	cpu := CreateCPUDebuggable(&bus, logger)

	nes := Nes{
		&cpu,
		&bus,
		true,
	}

	return nes
}

func (nes *Nes) Start() {
	nes.cpu.initInstructionsTable()
	//nes.cpu.reset()
	for {
		nes.cpu.tick()
	}
}

func (nes *Nes) InsertCartridge(cartridge *GamePak) {
	nes.bus.attachCartridge(cartridge)
}
