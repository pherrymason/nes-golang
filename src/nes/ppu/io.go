package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
)

// Read made by CPU
func (ppu *P2c02) ReadRegister(register types.Address) byte {
	value := byte(0x00)

	switch register {
	case PPUCTRL:
		panic("trying to read PPUCTRL")

	case PPUMASK:
		//panic("trying to read PPMASK")
		value = 0
		break

	case PPUSTATUS:
		// Source: javid9x reading from status only get top 3 bits. The rest tends to be filled with noise, or more likely what was last in data buffer.
		value = ppu.PpuStatus.Value()

		// Reading from status register alters it
		ppu.PpuStatus.VerticalBlankStarted = false // Reading from status, clears VBlank flag.
		//ppu.registers.status &= 0x7F
		ppu.tRam.resetLatch()
		break

	case OAMADDR:
		break

	case OAMDATA:
		value = ppu.oamData[ppu.oamAddr]
		break

	case PPUSCROLL:
		break

	case PPUADDR:
		break

	case PPUDATA:
		// TODO test delay and not delay from palette
		value = ppu.readBuffer
		ppu.readBuffer = ppu.Read(ppu.vRam.address())

		// If reading from Palette, there is no delay
		if isPaletteAddress(ppu.vRam.address()) {
			value = ppu.readBuffer
		}

		ppu.vRam.increment(ppu.PpuControl.IncrementMode)
		break

	case OAMDMA:
		break
	}

	return value
}

// Write made by CPU
func (ppu *P2c02) WriteRegister(register types.Address, value byte) {
	if !ppu.warmup {
		//log.Printf("Ignoring write register: %40X: %0X\n", register, value)
		return
	}

	switch register {
	case PPUCTRL:
		ppu.ppuCtrlWrite(value)
		ppu.tRam.setNameTableX(ppu.PpuControl.NameTableX)
		ppu.tRam.setNameTableY(ppu.PpuControl.NameTableY)
		// todo trigger nmi if in vblank and generateNMI transitions from 0 to 1
		break

	case PPUMASK:
		ppu.PpuMask.write(value)
		break

	case PPUSTATUS:
		// READONLY!
		panic("tried to write @PPUSTATUS")

	case OAMADDR:
		ppu.oamAddr = value
		break

	case OAMDATA:
		ppu.oamData[ppu.oamAddr] = value
		ppu.oamAddr = (ppu.oamAddr + 1) & 0xFF
		break

	case PPUSCROLL:
		if ppu.tRam.latch == 0 {
			ppu.tRam._coarseX = value >> 3
			ppu.fineX = value & 0x07
			ppu.tRam.latch = 1
		} else {
			ppu.tRam._coarseY = value >> 3
			ppu.tRam._fineY = value & 0b111
			ppu.tRam.latch = 0
		}
		break

	case PPUADDR:
		ppu.tRam.push(value)
		if ppu.tRam.latch == 0 {
			ppu.vRam = ppu.tRam
		}
		break
	case PPUDATA:
		address := ppu.vRam.address()
		ppu.Write(address, value)
		ppu.vRam.increment(ppu.PpuControl.IncrementMode)
		break
	case OAMDMA:
		//fmt.Println("OAMDMA!")
		break
	}
}

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

func (ppu *P2c02) Peek(address types.Address) byte {
	return ppu.read(address, true)
}

func (ppu *P2c02) Read(address types.Address) byte {
	return ppu.read(address, false)
}

func (ppu *P2c02) read(address types.Address, readOnly bool) byte {
	result := byte(0x00)

	// CHR ROM address
	if isCHRAddress(address) {
		result = ppu.cartridge.ReadCHRROM(address)
	} else if isNameTableAddress(address) {
		// Nametable 0, 1, 2, 3
		mirroring := ppu.cartridge.Header().Mirroring()
		nameTableAddress := getNameTableAddress(mirroring, address)
		result = ppu.nameTables[nameTableAddress]
	} else if isPaletteAddress(address) {
		result = ppu.readPalette(address)
	}

	return result
}

func (ppu *P2c02) Write(address types.Address, value byte) {
	if isNameTableAddress(address) {
		nameTableAddress := getNameTableAddress(ppu.cartridge.Header().Mirroring(), address)
		if ppu.nameTables[nameTableAddress] != value {
			ppu.nameTableChanged = true
		}

		ppu.nameTables[nameTableAddress] = value

	} else if address == 0x4010 {
		// OAM DMA: Transfers 256 bytes of data from CPU page $XX00-$XXFF to internal PPU OAM
		// DMA will begin at current OAM write address.
	} else if isPaletteAddress(address) {
		ppu.writePalette(address, value)
	} else if isCHRAddress(address) {
		ppu.cartridge.WriteCHRRAM(address, value)
	} else {
		var vblank string
		if ppu.PpuStatus.VerticalBlankStarted {
			vblank = "yes"
		} else {
			vblank = "no"
		}
		err := fmt.Sprintf("Unhandled ppu write address: 0x%X, ppu cycle: %d, Scanline: %d, vBlank: %s", address, ppu.renderCycle, ppu.currentScanline, vblank)
		panic(err)
	}
}

func isCHRAddress(address types.Address) bool {
	return address >= 0x0000 && address <= 0x1FFF
}

func isNameTableAddress(address types.Address) bool {
	// $3000-$3EFF nametable mirrors!
	return address >= NameTableStartAddress && /*address <= NameTableEndAddress*/
		address <= 0x3EFF
}

func isPaletteAddress(address types.Address) bool {
	return address >= PaletteLowAddress && address <= PaletteHighAddress
}

// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
func (ppu *P2c02) readPalette(address types.Address) byte {
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

func (ppu *P2c02) writePalette(address types.Address, colorIndex byte) {
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

func getNameTableAddress(mirrorMode byte, address types.Address) types.Address {
	realAddress := address
	// $2000-$23FF 	$0400 	Nametable 0
	// $2400-$27FF 	$0400 	Nametable 1
	// $2800-$2BFF 	$0400 	Nametable 2
	// $2C00-$2FFF 	$0400 	Nametable 3
	// $3000-$3EFF 	$0F00
	if mirrorMode == gamePak.VerticalMirroring {
		// -----------------------------
		// |    $2000    |    $2400    |
		// |      A      |      B      |
		// |-------------+-------------|
		// |    $2800    |    $2C00    |
		// |      A      |      B      |
		// |-------------+-------------|
		// |    $3000    |    $3400    |
		// |      A      |      B      |
		// |-------------+-------------|
		// |    $3800    |    $3C00    |
		// |      A      |    $3EFF    |
		// |-------------|-------------|

		// Inspired from fceux https://github.com/TASEmulators/fceux/blob/d1467182046e7ca00d65cd35f20ee011b2a665e6/src/ppu.cpp#L524

		nameTable := (address >> 10) & 0x3
		mask := types.Address(0x3FF)

		if nameTable > 1 {
			nameTable -= 2 // Substracting 2 gives us the real nametable
		}

		realAddress = 0x2000 + (0x400 * nameTable)
		if float32(address)/0x400 > 0xF {
			realAddress += address & 0x2FF
		} else {
			realAddress += address & mask
		}
	} else if mirrorMode == gamePak.HorizontalMirroring {
		// -----------------------------
		// |    $2000    |    $2400    |
		// |      A      |      A      |
		// |-------------+-------------|
		// |    $2800    |    $2C00    |
		// |      B      |      B      |
		// |-------------+-------------|
		// |    $3000    |    $3400    |
		// |      A      |      A      |
		// |-------------+-------------|
		// |    $3800    |    $3C00    |
		// |      B      |    $3EFF    |
		// |-------------|-------------|
		address = (address - 0x2000) % 0x1000 // keep the 0xFFF part
		table := address / 0x0400             // Nametable index
		y := table / 2
		x := table + 1
		realAddress = address - ((y - (x - 1)) * 0x400)
		// Formula explanation
		// Having these tables
		// +----+-----+
		// | 0  |  1  |   Nametable A: 0x000 -> 0x3FF
		// +----+-----+
		// | 2  |  3  |   Nametable B: 0x400 -> 0x7FF
		// +----+-----+
		// We can infer that even rows are nametable a.
		// We need to calculate how many times we need to subtract 0x400 to the address
		// to reach 0x00 (name table A) or 0x400 (nametable B)
		// Know the row, we subtract x -1 to the row and multiply by 0x400
	} else if mirrorMode == gamePak.OneScreenMirroring {
		realAddress = (address) & 0x3FF
	}

	return realAddress % 2048
}
