package nes

import (
	"github.com/raulferras/nes-golang/src/nes/component"
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
)

type Nes struct {
	cpu   *cpu2.Cpu6502
	bus   *component.Bus
	debug DebuggableNes
}

func CreateNes() Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPU(&bus)

	nes := Nes{
		cpu:   &cpu,
		bus:   &bus,
		debug: DebuggableNes{false, nil, 0},
	}

	return nes
}

type DebuggableNes struct {
	debug       bool
	logger      *cpu2.Logger
	cyclesLimit uint16
}

func CreateDebuggableNes(options DebuggableNes) Nes {
	ram := component.RAM{}
	bus := component.CreateBus(&ram)
	cpu := cpu2.CreateCPUDebuggable(&bus, options.logger)

	nes := Nes{
		&cpu,
		&bus,
		DebuggableNes{
			true,
			options.logger,
			options.cyclesLimit,
		},
	}

	return nes
}

func (nes *Nes) Start() {
	nes.cpu.Init()
	var i uint16 = 1

	//nes.cpu.reset()
	for {
		opCyclesLeft := nes.cpu.Tick()
		if opCyclesLeft == 0 {
			i++
		}
		if nes.debug.cyclesLimit > 0 && i >= nes.debug.cyclesLimit {
			break
		}
	}
}

func (nes *Nes) InsertCartridge(cartridge *component.GamePak) {
	nes.bus.AttachCartridge(cartridge)
}
