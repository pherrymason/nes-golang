package ppu

import (
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
)

// Address      Size
// $0000-$0FFF 	$1000 	Pattern table 0 \ CHR ROM 4KB
// $1000-$1FFF 	$1000 	Pattern table 1 / CHR ROM 4KB
// $2000-$23FF 	$0400 	Nametable 0		\
// $2400-$27FF 	$0400 	Nametable 1		| NameTable Memory
// $2800-$2BFF 	$0400 	Nametable 2		|
// $2C00-$2FFF 	$0400 	Nametable 3		/
// $3000-$3EFF 	$0F00 	Mirrors of $2000-$2EFF
// $3F00-$3F1F 	$0020 	Palette RAM indexes		} Palette Memory
const PaletteLowAddress = types.Address(0x3F00)
const PaletteHighAddress = types.Address(0x3FFF)
const PPU_HIGH_ADDRESS = types.Address(0x3FFF)

type Memory struct {
	gamePak *gamePak2.GamePak
	vram    [2048]byte

	paletteTable [32]byte
	//$3F00 	    Universal background color
	//$3F01-$3F03 	Background palette 0
	//$3F05-$3F07 	Background palette 1
	//$3F09-$3F0B 	Background palette 2
	//$3F0D-$3F0F 	Background palette 3
	//$3F11-$3F13 	Sprite palette 0
	//$3F15-$3F17 	Sprite palette 1
	//$3F19-$3F1B 	Sprite palette 2
	//$3F1D-$3F1F 	Sprite palette 3
}

func CreateMemory(gamePak *gamePak2.GamePak) *Memory {
	return &Memory{
		gamePak: gamePak,
	}
}

func (memory *Memory) Peek(address types.Address) byte {
	return memory.read(address, false)
}

func (memory *Memory) Read(address types.Address) byte {
	return memory.read(address, true)
}

func (memory *Memory) read(address types.Address, readOnly bool) byte {
	result := byte(0x00)

	// CHR ROM address
	if address < 0x01FFF {
		result = memory.gamePak.ReadCHRROM(address)
	} else if address >= 0x2000 && address <= 0x2FFF {
		// Nametable 0, 1, 2, 3
		mirroring := memory.gamePak.Header().Mirroring()
		realAddress := nameTableMirrorAddress(mirroring, address)
		result = memory.vram[realAddress]
	} else if isPaletteAddress(address) {
		result = readPalette(memory, address)
	}

	return result
}

func (memory *Memory) Write(address types.Address, value byte) {
	if address >= 0x2000 && address <= 0x2FFF {
		realAddress := nameTableMirrorAddress(memory.gamePak.Header().Mirroring(), address)
		memory.vram[realAddress] = value
	} else if address == 0x4010 {
		// OAM DMA: Transfers 256 bytes of data from CPU page $XX00-$XXFF to internal PPU OAM
		// DMA will begin at current OAM write address.
	} else if isPaletteAddress(address) {
		writePalette(memory, address, value)
	} else {
		panic("Unhandled ppu address")
	}
}

func isPaletteAddress(address types.Address) bool {
	return address >= PaletteLowAddress && address <= PaletteHighAddress
}

// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
func readPalette(memory *Memory, address types.Address) byte {
	address &= 0x1F
	// Mirrors
	// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
	if address == 0x10 {
		address = 0x00
	} else if address == 0x14 {
		address = 0x04
	} else if address == 0x18 {
		address = 0x08
	} else if address == 0x1C {
		address = 0x0C
	}

	return memory.paletteTable[address]
}

func writePalette(memory *Memory, address types.Address, colorIndex byte) {
	address &= 0x1F
	// Mirrors
	if address == 0x10 {
		address = 0x00
	} else if address == 0x14 {
		address = 0x04
	} else if address == 0x18 {
		address = 0x08
	} else if address == 0x1C {
		address = 0x0C
	}
	memory.paletteTable[address] = colorIndex
}

func nameTableMirrorAddress(mirrorMode byte, address types.Address) types.Address {
	realAddress := address
	if mirrorMode == gamePak2.VerticalMirroring {
		realAddress = (address - 0x2000) & 0x27FF
		/*
			if address >= 0x2000 && address <= 0x23FF {
				// Nametable 0
				realAddress = address - 0x2000
			} else if address >= 0x2400 && address < 0x27FF {
				// Nametable 2
				realAddress = address - 0x2000
			} else if address >= 0x2800 && address <= 0x2BFF {
				// Nametable 1
				realAddress = address - 0x2800
			} else {
				// Nametable 3
				realAddress = address - 0x2800
			}*/
	} else if mirrorMode == gamePak2.HorizontalMirroring {
		if address >= 0x2000 && address < 0x2400 {
			realAddress = address - 0x2000
		} else if address >= 0x2400 && address <= 0x27FF {
			realAddress = address - 0x2400
		} else if address >= 0x2800 && address <= 0x2BFF {
			realAddress = address - 0x2400
		} else if address >= 0x2C00 && address <= 0x2FFF {
			realAddress = address - 0x2800
		}
	} else if mirrorMode == gamePak2.OneScreenMirroring {
		realAddress = (address - 0x2000) & 0x3FF
	}

	return realAddress
}
