package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

// Registers
const PPUCTRL = 0x2000 // NMI enable (V), PPU master/slave (P), sprite height (H),
// background tile select (B), sprite tile select (S), increment mode (I),
//nametable select (NN)
const PPUMASK = 0x2001 // color emphasis (BGR), sprite enable (s), background enable (b),
// sprite left column enable (M), background left column enable (m), greyscale (G)
const PPUSTATUS = 0x2002 // vblank (V), sprite 0 hit (S), sprite overflow (O); read resets write pair for $2005/$2006
const OAMADDR = 0x2003
const OAMDATA = 0x2004
const PPUSCROLL = 0x2005
const PPUADDR = 0x2006
const PPUDATA = 0x2007
const OAMDMA = 0x4014

// Registers end

const NES_PALETTE_COLORS = 64
const OAMDATA_SIZE = 256

const PPU_SCREEN_SPACE_CYCLES_BY_SCANLINE = 256
const PPU_CYCLES_BY_SCANLINE = 341
const PPU_SCREEN_SPACE_SCANLINES = 240
const VBLANK_START_SCANLINE = 241
const VBLANK_END_SCNALINE = 261
const PPU_SCANLINES = 261
const PPU_VBLANK_START_CYCLE = (PPU_SCREEN_SPACE_SCANLINES + 1) * PPU_CYCLES_BY_SCANLINE
const PPU_VBLANK_END_CYCLE = PPU_SCANLINES * PPU_CYCLES_BY_SCANLINE

type PPU interface {
	WriteRegister(register types.Address, value byte)
	ReadRegister(register types.Address) byte
}

type Ppu2c02 struct {
	registers Registers
	memory    Memory

	// OAM (Object Attribute Memory) is internal memory inside the PPU.
	// Contains a display list of up to 64 sprites, where each sprite occupies 4 bytes
	oamData [OAMDATA_SIZE]byte

	cycle           uint32
	renderCycle     uint16
	currentScanline uint16
	nmi             bool // NMI Interrupt thrown
	frame           types.Frame
	frameSprites    types.Frame
	framePattern    [1024]byte
}

func CreatePPU(memory Memory) *Ppu2c02 {
	ppu := &Ppu2c02{
		memory:          memory,
		currentScanline: 0,
	}

	return ppu
}

func (ppu *Ppu2c02) Frame() *types.Frame {
	return &ppu.frame
}

func (ppu *Ppu2c02) FramePattern() *[1024]byte {
	return &ppu.framePattern
}

func (ppu *Ppu2c02) Tick() {
	//bit := ppu.registers.scrollX

	// Load new data into registers
	if ppu.cycle%8 == 0 {

	}

	// 341 PPU clock cycles have passed
	if ppu.renderCycle%PPU_CYCLES_BY_SCANLINE == 0 {
		ppu.currentScanline++
		ppu.renderCycle = 0
	}

	// VBlank logic
	if ppu.currentScanline == VBLANK_START_SCANLINE {
		ppu.registers.status |= 1 << verticalBlankStarted // Todo refactor to a method to set Vblank
		if (ppu.registers.ctrl & (1 << generateNMIAtVBlank)) > 0 {
			if ppu.renderCycle == 0 {
				ppu.nmi = true
			}
		}
	} else if ppu.currentScanline == VBLANK_END_SCNALINE && ppu.renderCycle == PPU_CYCLES_BY_SCANLINE-1 {
		ppu.registers.status &= ^byte(1 << verticalBlankStarted)
		ppu.currentScanline = 0
	}
	ppu.cycle++
	ppu.renderCycle++
}

func (ppu *Ppu2c02) VBlank() bool {
	return ppu.currentScanline >= 241
}

func (ppu *Ppu2c02) Peek(address types.Address) byte {
	return ppu.memory.Peek(address)
}

func (ppu *Ppu2c02) Read(address types.Address) byte {
	return ppu.memory.Read(address)
}

func (ppu *Ppu2c02) Write(address types.Address, value byte) {
	ppu.memory.Write(address, value)
}

func (ppu *Ppu2c02) Nmi() bool {
	occurred := ppu.nmi
	ppu.nmi = false

	return occurred
}

func (ppu *Ppu2c02) ResetNmi() {
	ppu.nmi = false
}

// Read made by CPU
func (ppu *Ppu2c02) ReadRegister(register types.Address) byte {
	value := byte(0x00)

	switch register {
	case PPUCTRL:
		panic("trying to read PPUCTRL")

	case PPUMASK:
		break

	case PPUSTATUS:
		value = ppu.registers.status
		ppu.registers.status &= 0x7F // Clear VBlank flag. Why are we clearing this?
		break

	case OAMADDR:
		break

	case OAMDATA:
		value = ppu.oamData[ppu.registers.oamAddr]
		break

	case PPUSCROLL:
		break

	case PPUADDR:
		break

	case PPUDATA:
		value = ppu.registers.readBuffer
		ppu.registers.readBuffer = ppu.memory.Read(ppu.registers.ppuAddr)
		if ppu.registers.ppuAddr >= 0x3F00 && ppu.registers.ppuAddr <= 0x3FFF {
			value = ppu.registers.readBuffer
		}

		if ppu.ppuctrlReadFlag(incrementMode) == 0 {
			ppu.registers.ppuAddr++
		} else {
			ppu.registers.ppuAddr += 32
		}
		ppu.registers.ppuAddr &= 0x3FFF
		break

	case OAMDMA:
		break
	}

	return value
}

// Write made by CPU
func (ppu *Ppu2c02) WriteRegister(register types.Address, value byte) {
	switch register {
	case PPUCTRL:
		if ppu.cycle > 30000 {
			ppu.registers.ctrl = value

			// todo trigger nmi if in vblank and generateNMI transitions from 0 to 1
		}
		break

	case PPUMASK:
		ppu.registers.mask = value
		break

	case PPUSTATUS:
		// READONLY!
		panic("tried to write @PPUSTATUS")

	case OAMADDR:
		ppu.registers.oamAddr = value
		break

	case OAMDATA:
		ppu.oamData[ppu.registers.oamAddr] = value
		ppu.registers.oamAddr = (ppu.registers.oamAddr + 1) & 0xFF
		break

	case PPUSCROLL:
		if ppu.registers.scrollLatch == 0 {
			ppu.registers.scrollX = value
		} else {
			ppu.registers.scrollY = value
		}

		ppu.registers.scrollLatch = (ppu.registers.scrollLatch + 1) & 0x01
		break

	case PPUADDR:
		if ppu.registers.addressLatch == 0 {
			ppu.registers.ppuAddr = types.Address(value) << 8
		} else {
			ppu.registers.ppuAddr |= types.Address(value)
		}

		ppu.registers.addressLatch++
		ppu.registers.addressLatch &= 0b1

		if ppu.registers.addressLatch == 0 {
			ppu.registers.ppuAddr &= 0x3FFF
		}
		break
	case PPUDATA:
		ppu.Write(ppu.registers.ppuAddr, value)

		if ppu.ppuctrlReadFlag(incrementMode) == 0 {
			ppu.registers.ppuAddr++
		} else {
			ppu.registers.ppuAddr += 32
		}
		ppu.registers.ppuAddr &= 0x3FFF
		break
	case OAMDMA:
		break
	}
}

func (ppu *Ppu2c02) ppuctrlReadFlag(flag CtrlFlag) byte {
	ctrl := ppu.registers.ctrl
	mask := byte(1) << flag

	return (ctrl & mask) >> flag
}

// Helper method, only used in tests
func (ppu *Ppu2c02) ppuctrlWriteFlag(flag CtrlFlag, value byte) {
	if value == 1 {
		ppu.registers.ctrl |= 1 << flag
	} else {
		ppu.registers.ctrl &= ^(1 << flag)
	}
}

/*
	//$3F00 	    Universal background color
	//$3F01-$3F03 	Background palette 0
	//$3F05-$3F07 	Background palette 1
	//$3F09-$3F0B 	Background palette 2
	//$3F0D-$3F0F 	Background palette 3
	//$3F11-$3F13 	Sprite palette 0
	//$3F15-$3F17 	Sprite palette 1
	//$3F19-$3F1B 	Sprite palette 2
	//$3F1D-$3F1F 	Sprite palette 3
*/
func (ppu *Ppu2c02) GetColorFromPaletteRam(palette byte, colorIndex byte) types.Color {
	paletteAddress := types.Address((palette * 4) + colorIndex)
	/*
		if int(colorIndex) > len(SystemPalette) {
			panic(fmt.Sprintf("pixel color out of palette: %X", colorIndex))
		}*/

	color := ppu.Read(PaletteLowAddress + paletteAddress)

	return types.Color{
		R: SystemPalette[color][0],
		G: SystemPalette[color][1],
		B: SystemPalette[color][2],
	}
}