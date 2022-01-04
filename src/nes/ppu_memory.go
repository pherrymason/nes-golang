package nes

import (
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
)

// PPUMemory map
// -------------
// PPU Memory addresses a 16kB space: From $0000-$3FFFF.
// As the CPU, different ranges point in reality to different actual devices:
//
// $0000-$0FFF: Pattern table 0.
// $1000-$1FFF: Pattern Table 1.
// These two banks are usually mapped by the cartridge to a CHR-ROM or CHR-RAM.

// $2000-$23FF

//
// Address      Size
// $0000-$0FFF 	$1000 	Pattern table 0 \ CHR ROM 4KB
// $1000-$1FFF 	$1000 	Pattern table 1 / CHR ROM 4KB
// $2000-$23FF 	$0400 	Nametable 0		\
// $2400-$27FF 	$0400 	Nametable 1		| NameTable Memory
// $2800-$2BFF 	$0400 	Nametable 2		|
// $2C00-$2FFF 	$0400 	Nametable 3		/
// $3000-$3EFF 	$0F00 	Mirrors of $2000-$2EFF
// $3F00-$3F1F 	$0020 	Palette RAM indexes		} Palette Memory
// $3F20-$3FFF 	$00E0 	Mirrors of $3F00-$3F1F
// ---------------------------------------------
// The actual device that the PPU fetches data from, however, may be configured by the cartridge.
//
// - $0000-1FFF is normally mapped by the cartridge to a CHR-ROM or CHR-RAM, often with a bank switching mechanism.
// - $2000-2FFF is normally mapped to the 2kB NES internal VRAM, providing 2 nametables with a mirroring configuration controlled by the cartridge, but it can be partly or fully remapped to RAM on the cartridge, allowing up to 4 simultaneous nametables.
// - $3000-3EFF is usually a mirror of the 2kB region from $2000-2EFF. The PPU does not render from this address range, so this space has negligible utility.
// - $3F00-3FFF is not configurable, always mapped to the internal palette control.
type PPUMemory struct {
	gamePak *gamePak2.GamePak
	vram    [2048]byte

	paletteTable [32]byte
}

func CreatePPUMemory(gamePak *gamePak2.GamePak) *PPUMemory {
	return &PPUMemory{
		gamePak: gamePak,
	}
}

func (ppu *PPUMemory) Peek(address types.Address) byte {
	return ppu.read(address, false)
}

func (ppu *PPUMemory) Read(address types.Address) byte {
	return ppu.read(address, true)
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
	}

	return realAddress
}

func (ppu *PPUMemory) read(address types.Address, readOnly bool) byte {
	result := byte(0x00)

	// CHR ROM address
	if address < 0x01FFF {
		result = ppu.gamePak.ReadCHRROM(address)
	} else if address >= 0x2000 && address <= 0x2FFF {
		// Nametable 0, 1, 2, 3
		mirroring := ppu.gamePak.Header().Mirroring()
		realAddress := nameTableMirrorAddress(mirroring, address)
		result = ppu.vram[realAddress]
	} else if address >= 0x3F00 && address <= 0x3FFF {
		// palette ram indexes
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
		result = ppu.paletteTable[address]
	}

	return result
}

func (ppu *PPUMemory) Write(address types.Address, value byte) {
	if address >= 0x2000 && address <= 0x2FFF {
		realAddress := nameTableMirrorAddress(ppu.gamePak.Header().Mirroring(), address)
		ppu.vram[realAddress] = value
	} else if address == 0x4010 {
		// OAM DMA: Transfers 256 bytes of data from CPU page $XX00-$XXFF to internal PPU OAM
		// DMA will begin at current OAM write address.
	} else if address >= 0x3F00 && address <= 0x3FFF {
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
		ppu.paletteTable[address] = value
	}
}
