package component

import (
	"github.com/raulferras/nes-golang/src/nes/defs"
)

type Bus struct {
	//
	// 	0x0000 - 0x00FF  ZeroPage
	//	0x0100 - 0x01FF  Stack
	//	0x0200 - 0x0800  General Purpose RAM
	//  0x0801 - 0x2000  Mirrors previous chunk of memory (0x0 - 0x7FF)
	//  0x2000 - 0x2007  PPU registers
	//  0x2008 - 0x4000  Mirrors PPU registers
	//  0x4000 - 0x4020  I/O Registers
	//  0x4020 - 0x5FFF  Expansion ROM
	//  0x6000 - 0x7FFF  SRAM
	//  0x8000 - 0xBFFF  GamePak prgROM lower bank
	//  0xC000 - 0x10000 GamePak prgROM higher bank

	Ram *RAM // is $0000 -> $07FF
	// RAM mirrors $0800 -> $1FFF
	Cartridge *GamePak // 0x8000 -> $FFFF

	// APU $4000 -> $4017
	// PPU -> $2000 -> $2007
}

func (bus *Bus) Read(address defs.Address) byte {
	if address <= 0x7FF {
		return bus.Ram.read(address)
	} else if address >= 0x8000 {
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
	msb := (lsb & 0xFF00) | (lsb & 0xFF) + 1

	low := bus.Read(lsb)
	high := bus.Read(msb)

	return defs.CreateWord(low, high)
}

func (bus *Bus) Write(address defs.Address, value byte) {
	if address <= 0x7FF {
		bus.Ram.write(address, value)
	} else if address >= 0x8000 {
		bus.Cartridge.write(address, value)
	}
}

func (bus *Bus) AttachCartridge(cartridge *GamePak) {
	bus.Cartridge = cartridge
}

func CreateBus(ram *RAM) Bus {
	return Bus{Ram: ram}
}
