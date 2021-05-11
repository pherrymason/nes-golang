package nes

import "github.com/raulferras/nes-golang/src/graphics"

type PPU interface {
	WriteRegister(register Address, value byte)
	ReadRegister(register Address) byte
}

// Ppu2c02 Processor
//  Registers mapped to memory locations: $2000 through $2007
//  mirrored in every 8 bytes from $2008 through $3FFF
type Ppu2c02 struct {
	Memory
	registers PPURegisters

	// OAM (Object Attribute Memory) is internal memory inside the PPU.
	// Contains a display list of up to 64 sprites, where each sprite occupies 4 bytes
	oamData [256]byte

	cycle uint32

	patternTable []byte // Decoded pattern table
}

type PPURegisters struct {
	ctrl   byte // Controls PPU operation
	mask   byte // Controls the rendering of sprites and backgrounds
	status byte // Reflects state of various functions inside PPU

	scrollX     byte // Changes scroll position
	scrollY     byte // Changes scroll position
	scrollLatch byte // Controls which scroll needs to be written

	oamAddr      byte
	ppuAddr      Address
	addressLatch byte
	readBuffer   byte
}

type PPUCtrlFlag int

const (
	baseNameTableAddress0 PPUCtrlFlag = iota // Most significant bit of scrolling coordinates (X)
	baseNameTableAddress1                    // Most significant bit of scrolling coordinates (Y)
	incrementMode
	spritePatternTableAddress
	backgroundPatternTableAddress
	spriteSize
	generateNMIAtVBlank
)

type PPUMASKFlag int

const (
	greyScale PPUMASKFlag = iota
	showBackgroundLeftEdge
	showSpritesLeftEdge
	showBackground
	showSprites
	emphasizeRed
	emphasizeGreen
	emphasizeBlue
)

type PPUSTATUSFlag int

const (
	spriteOverflow       = 5
	sprite0Hit           = 6
	verticalBlankStarted = 7
)

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

const NES_PALETTE_COLORS = 64

func CreatePPU(memory Memory) *Ppu2c02 {
	ppu := &Ppu2c02{
		Memory:       memory,
		patternTable: make([]byte, 8*8*512),
	}

	return ppu
}

func (ppu *Ppu2c02) Tick() {

	//bit := ppu.registers.scrollX

	// Load new data into registers
	if ppu.cycle%8 == 0 {

	}

	ppu.cycle++
}

func (ppu *Ppu2c02) ReadRegister(register Address) byte {
	value := byte(0x00)

	switch register {
	case PPUCTRL:
		panic("trying to read PPUCTRL")

	case PPUMASK:
		break

	case PPUSTATUS:
		value = ppu.registers.status
		// set vblank for test
		value |= 1 << verticalBlankStarted
		ppu.registers.status &= 0x7F // Clear VBlank flag
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
		ppu.registers.readBuffer = ppu.Read(ppu.registers.ppuAddr)

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

func (ppu *Ppu2c02) WriteRegister(register Address, value byte) {
	switch register {
	case PPUCTRL:
		if ppu.cycle > 30000 {
			ppu.registers.ctrl = value
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
			ppu.registers.ppuAddr = Address(value) << 8
		} else {
			ppu.registers.ppuAddr |= Address(value)
		}

		ppu.registers.addressLatch++
		ppu.registers.addressLatch &= 0b1

		if ppu.registers.addressLatch == 0 {
			ppu.registers.ppuAddr &= 0x3FFF
		}
		break
	case PPUDATA:
		break
	case OAMDMA:
		break
	}
}

func (ppu *Ppu2c02) PatternTable(patternTable int, palette uint8) []graphics.Pixel {
	chr := make([]graphics.Pixel, 128*128)

	bankAddress := Address(patternTable * 0x1000)
	i := 0
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			offset := bankAddress + Address(tileY*256+tileX*16)

			for row := 0; row < 8; row++ {
				pixelLSB := ppu.Read(offset + Address(row))
				pixelMSB := ppu.Read(offset + Address(row+8))

				for col := 0; col < 8; col++ {
					value := (pixelLSB & 0x01) | ((pixelMSB & 0x01) << 1)

					pixelLSB >>= 1
					pixelMSB >>= 1

					coordY := tileY*8 + row
					coordX := (7 - col) + tileX*8

					pixel := graphics.Pixel{
						X:     coordX,
						Y:     coordY,
						Color: ppu.getColorFromPaletteRam(palette, value),
					}
					chr[i] = pixel
					i++
				}
			}
		}
	}

	return chr
}

func (ppu *Ppu2c02) getTile(patternTable int, palette uint8, tileX int, tileY int) []graphics.Pixel {
	tile := make([]graphics.Pixel, 8*8)
	bankAddress := Address(patternTable * 0x1000)

	offset := bankAddress + Address(tileY*256+tileX*16)

	for row := 0; row < 8; row++ {
		pixelLSB := ppu.Read(offset + Address(row))
		pixelMSB := ppu.Read(offset + Address(row+8))

		for col := 0; col < 8; col++ {
			value := (pixelLSB & 0x01) | ((pixelMSB & 0x01) << 1)

			pixelLSB >>= 1
			pixelMSB >>= 1

			coordY := tileY*8 + row
			coordX := (7 - col) + tileX*8

			pixel := graphics.Pixel{
				X:     coordX,
				Y:     coordY,
				Color: ppu.getColorFromPaletteRam(palette, value),
			}
			tile[row+col*8] = pixel
		}
	}

	return tile
}

// blargg's palette
var SystemPalette = [...][3]byte{
	{84, 84, 84},
	{0, 30, 116},
	{8, 16, 144},
	{48, 0, 136},
	{68, 0, 100},
	{92, 0, 48},
	{84, 4, 0},
	{60, 24, 0},
	{32, 42, 0},
	{8, 58, 0},
	{0, 64, 0},
	{0, 60, 0},
	{0, 50, 60},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},

	// 0x10
	{152, 150, 152},
	{8, 76, 196},
	{48, 50, 236},
	{92, 30, 228},
	{136, 20, 176},
	{160, 20, 100},
	{152, 34, 32},
	{120, 60, 0},
	{84, 90, 0},
	{40, 114, 0},
	{8, 124, 0},
	{0, 118, 40},
	{0, 102, 120},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},

	// 0x20
	{236, 238, 236},
	{76, 154, 236},
	{120, 124, 236},
	{176, 98, 236},
	{228, 84, 236},
	{236, 88, 180},
	{236, 106, 100},
	{212, 136, 32},
	{160, 170, 0},
	{116, 196, 0},
	{76, 208, 32},
	{56, 204, 108},
	{56, 180, 204},
	{60, 60, 60},
	{0, 0, 0},
	{0, 0, 0},

	// 0x3F
	{236, 238, 236},
	{168, 204, 236},
	{188, 188, 236},
	{212, 178, 236},
	{236, 174, 236},
	{236, 174, 212},
	{236, 180, 176},
	{228, 196, 144},
	{204, 210, 120},
	{180, 222, 120},
	{168, 226, 144},
	{152, 226, 180},
	{160, 214, 228},
	{160, 162, 160},
	{0, 0, 0},
	{0, 0, 0},
}

func (ppu Ppu2c02) getColorFromPaletteRam(palette byte, pixelColor byte) graphics.Color {
	paletteAddress := Address((palette * 4) + pixelColor)
	/*
		if int(colorIndex) > len(SystemPalette) {
			panic(fmt.Sprintf("pixel color out of palette: %X", pixelColor))
		}*/

	colorIndex := ppu.Read(Address(0x3F00) + paletteAddress)

	return graphics.Color{R: SystemPalette[colorIndex][0], G: SystemPalette[colorIndex][1], B: SystemPalette[colorIndex][2]}
}

func (ppu *Ppu2c02) ppuctrlReadFlag(flag PPUCtrlFlag) byte {
	ctrl := ppu.registers.ctrl
	mask := byte(1) << flag

	return (ctrl & mask) >> flag
}

func (ppu *Ppu2c02) ppuctrlWriteFlag(flag PPUCtrlFlag, value byte) {
	mask := byte(1)
	if flag == incrementMode {
		mask = 1 << flag
	}
	ppu.registers.ctrl |= (value << flag) & mask
}
