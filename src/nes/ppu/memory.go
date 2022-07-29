package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
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
//       $3F00 	        Universal background color
//	     $3F01-$3F03 	Background palette 0
//	     $3F05-$3F07 	Background palette 1
//	     $3F09-$3F0B 	Background palette 2
//	     $3F0D-$3F0F 	Background palette 3
//	     $3F11-$3F13 	Sprite palette 0
//	     $3F15-$3F17 	Sprite palette 1
//	     $3F19-$3F1B 	Sprite palette 2
//	     $3F1D-$3F1F 	Sprite palette 3
const PaletteLowAddress = types.Address(0x3F00)
const PaletteHighAddress = types.Address(0x3FFF)
const NameTableStartAddress = types.Address(0x2000)
const PPU_NAMETABLES_0_END = types.Address(0x23C0)
const NameTableEndAddress = types.Address(0x2FFF)
const PPU_HIGH_ADDRESS = types.Address(0x3FFF)

const PatternTable0Address = types.Address(0x0000)
const PatternTable1Address = types.Address(0x1000)

func (ppu *Ppu2c02) Peek(address types.Address) byte {
	return ppu.read(address, true)
}

func (ppu *Ppu2c02) Read(address types.Address) byte {
	return ppu.read(address, false)
}

func (ppu *Ppu2c02) read(address types.Address, readOnly bool) byte {
	result := byte(0x00)

	// CHR ROM address
	if address < 0x01FFF {
		result = ppu.cartridge.ReadCHRROM(address)
	} else if isNameTableAddress(address) {
		// Nametable 0, 1, 2, 3
		mirroring := ppu.cartridge.Header().Mirroring()
		realAddress := nameTableMirrorAddress(mirroring, address)
		result = ppu.nameTables[realAddress]
	} else if isPaletteAddress(address) {
		result = ppu.readPalette(address)
	}

	return result
}

func (ppu *Ppu2c02) Write(address types.Address, value byte) {
	if isNameTableAddress(address) {
		realAddress := nameTableMirrorAddress(ppu.cartridge.Header().Mirroring(), address)
		if ppu.nameTables[realAddress] != value {
			ppu.nameTableChanged = true
		}

		ppu.nameTables[realAddress] = value

	} else if address == 0x4010 {
		// OAM DMA: Transfers 256 bytes of data from CPU page $XX00-$XXFF to internal PPU OAM
		// DMA will begin at current OAM write address.
	} else if isPaletteAddress(address) {
		ppu.writePalette(address, value)
	} else {
		var vblank string
		if ppu.ppuStatus.verticalBlankStarted {
			vblank = "yes"
		} else {
			vblank = "no"
		}
		err := fmt.Sprintf("Unhandled ppu address: 0x%X, ppu cycle: %d, scanline: %d, vBlank: %s", address, ppu.renderCycle, ppu.currentScanline, vblank)
		panic(err)
	}
}

func isNameTableAddress(address types.Address) bool {
	return address >= NameTableStartAddress && address <= NameTableEndAddress
}

func isPaletteAddress(address types.Address) bool {
	return address >= PaletteLowAddress && address <= PaletteHighAddress
}

// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
func (ppu *Ppu2c02) readPalette(address types.Address) byte {
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

	return ppu.paletteTable[address]
}

func (ppu *Ppu2c02) writePalette(address types.Address, colorIndex byte) {
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
	ppu.paletteTable[address] = colorIndex
}

func nameTableMirrorAddress(mirrorMode byte, address types.Address) types.Address {
	realAddress := address
	if mirrorMode == gamePak.VerticalMirroring {
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
	} else if mirrorMode == gamePak.HorizontalMirroring {
		if address >= 0x2000 && address < 0x2400 {
			realAddress = address - 0x2000
		} else if address >= 0x2400 && address <= 0x27FF {
			realAddress = address - 0x2400
		} else if address >= 0x2800 && address <= 0x2BFF {
			realAddress = address - 0x2400
		} else if address >= 0x2C00 && address <= 0x2FFF {
			realAddress = address - 0x2800
		}
	} else if mirrorMode == gamePak.OneScreenMirroring {
		realAddress = (address - 0x2000) & 0x3FF
	}

	return realAddress
}
