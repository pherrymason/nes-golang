package nes

import "fmt"

type Pixel struct {
	X     int
	Y     int
	Color []byte
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
	ppu.cycle++
}

func (ppu *Ppu2c02) ReadRegister(register Address) byte {
	value := byte(0x00)

	switch register {
	case PPUCTRL:
		panic("trying to read PPUCTRL")
		break

	case PPUMASK:
		break

	case PPUSTATUS:
		value = ppu.registers.status
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

//func (ppu *Ppu2c02) PatternTable(patternTable int) [][]byte {
func (ppu *Ppu2c02) PatternTable(patternTable int) []Pixel {
	//chr := make([][]byte, 8*8*16*16)
	chr := make([]Pixel, 128*128)

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

					pixel := Pixel{X: coordX, Y: coordY, Color: palette(value)}
					chr[i] = pixel
					i++
				}
			}
		}
	}

	return chr
}

func palette(nesColor byte) []byte {

	SystemPalette := [][]byte{
		{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96}, {0xA1, 0x00, 0x5E},
		{0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00}, {0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00},
		{0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E}, {0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05},
		{0x05, 0x05, 0x05}, {0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
		{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00}, {0xC4, 0x62, 0x00},
		{0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55}, {0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21},
		{0x09, 0x09, 0x09}, {0x09, 0x09, 0x09}, {0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF},
		{0xD4, 0x80, 0xFF}, {0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
		{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4}, {0x05, 0xFB, 0xFF},
		{0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D}, {0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF},
		{0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB}, {0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0},
		{0xFF, 0xEF, 0xA6}, {0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
		{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},

		//{84  84  84},  {0  30 116},    {8  16 144},    {48   0 136 }, { 68   0 100   92   0  48   84   4   0   60  24   0   32  42   0    8  58   0    0  64   0    0  60   0    0  50  60    0   0   0
		//{152 150 152},  {8  76 196},   {48  50 236},   {92  30 228 }, {136  20 176  160  20 100  152  34  32  120  60   0   84  90   0   40 114   0    8 124   0    0 118  40    0 102 120    0   0   0
		//{236 238 236},  {76 154 236},  {120 124 236},  {176  98 236}, {228  84 236  236  88 180  236 106 100  212 136  32  160 170   0  116 196   0   76 208  32   56 204 108   56 180 204   60  60  60
		//{236 238 236},  {168 204 236}  {188 188 236},  {212 178 236}, {236 174 236  236 174 212  236 180 176  228 196 144  204 210 120  180 222 120  168 226 144  152 226 180  160 214 228  160 162 160
	}

	switch nesColor {
	case 0:
		return SystemPalette[0x0D]
	case 1:
		return SystemPalette[0x06]
	case 2:
		return SystemPalette[0x18]
	case 3:
		return SystemPalette[0x08]
	}

	if int(nesColor) > len(SystemPalette) {
		panic(fmt.Sprintf("pixel color out of palette: %X", nesColor))
	}

	return SystemPalette[nesColor]
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
