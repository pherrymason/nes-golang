package component

import (
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/defs"
)

type Bus struct {
	//
	// 	0x0000 - 0x00FF  ZeroPage
	//	0x0100 - 0x01FF  Stack
	//	0x0200 - 0x07FF  General Purpose RAM
	//  0x0800 - 0x1FFF  Mirrors previous chunk of memory (0x0 - 0x7FF)
	//  0x2000 - 0x2007  PPU registers
	//  0x2008 - 0x4000  Mirrors PPU registers
	//  0x4000 - 0x4020  I/O Registers
	//  0x4020 - 0x5FFF  Expansion ROM
	//  0x6000 - 0x7FFF  SRAM
	//  0x8000 - 0xBFFF  GamePak prgROM lower bank
	//  0xC000 - 0x10000 GamePak prgROM higher bank

	Ram *RAM // is $0000 -> $1FFF
	// RAM mirrors $0800 -> $1FFF
	Cartridge *GamePak // 0x8000 -> $FFFF

	cpu *cpu2.Cpu6502
	// APU $4000 -> $4017
	// PPU -> $2000 -> $2007

	ticks byte
}

func CreateBus(ram *RAM) *Bus {
	return &Bus{Ram: ram}
}

func (bus *Bus) ConnectCPU(cpu *cpu2.Cpu6502) {
	bus.cpu = cpu
	bus.cpu.ConnectBus(bus)
}

func (bus *Bus) ConnectPPU() {

}

func (bus *Bus) InsertGamePak(cartridge *GamePak) {
	bus.Cartridge = cartridge
}

func (bus *Bus) Tick() byte {
	if bus.ticks%3 == 0 {
		bus.cpu.Tick()
		bus.ticks = 0xFF
	}

	//bus.ppu.Tick()
	bus.ticks++

	return bus.ticks
}

func (bus *Bus) Read(address defs.Address) byte {
	return bus.ReadOnly(address)
}

func (bus *Bus) ReadOnly(address defs.Address) byte {
	if address <= RAM_HIGHER_ADDRESS {
		return bus.Ram.read(address)
	} else if address >= GAMEPAK_LOW_RANGE {
		return bus.Cartridge.read(address)
	}

	return 0
}

func (bus *Bus) Read16(address defs.Address) defs.Word {
	low := bus.Read(address)
	high := bus.Read(address + 1)

	return defs.CreateWord(low, high)
}

func (bus *Bus) Read16Bugged(address defs.Address) defs.Word {
	lsb := address
	msb := (lsb & 0xFF00) | defs.Address(byte(lsb)+1)

	low := bus.Read(lsb)
	high := bus.Read(msb)

	return defs.CreateWord(low, high)
}

func (bus *Bus) Write(address defs.Address, value byte) {
	if address <= 0x7FF {
		bus.Ram.write(address, value)
	} else if address >= GAMEPAK_ROM_LOWER_BANK_START {
		bus.Cartridge.write(address, value)
	}
}
