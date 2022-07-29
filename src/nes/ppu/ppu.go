package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
	"github.com/raulferras/nes-golang/src/utils"
	"image"
	"image/color"
)

type PPU interface {
	WriteRegister(register types.Address, value byte)
	ReadRegister(register types.Address) byte
}

type Ppu2c02 struct {
	ppuControl Control
	ppuStatus  Status
	ppuScroll  Scroll
	ppuMask    Mask // Controls the rendering of sprites and backgrounds
	vRam       loopyRegister
	tRam       loopyRegister
	fineX      uint8
	readBuffer byte

	nextTileId    byte
	nextAttribute byte
	nextLowTile   byte
	nextHighTile  byte

	shifterTileLow       uint16
	shifterTileHigh      uint16
	shifterAttributeLow  uint16
	shifterAttributeHigh uint16

	oamAddr byte

	// OAM (Object Attribute Memory) is internal memory inside the PPU.
	// Contains a display list of up to 64 sprites, where each sprite occupies 4 bytes
	cartridge    *gamePak.GamePak
	nameTables   [2 * NAMETABLE_SIZE]byte
	paletteTable [PALETTE_SIZE]byte
	oamData      [OAMDATA_SIZE]byte

	cycle  uint32 // Current lifetime PPU Cycle. After warmup, ignored.
	warmup bool   // Indicates ppu is already warmed up (cycles went above 30000)

	renderCycle     uint16 // Current cycle inside a scanline. From 0 to PPU_CYCLES_BY_SCANLINE
	currentScanline int16  // Current vertical scanline being rendered
	evenFrame       bool   // Is current frame even?

	nmi              bool // NMI Interrupt thrown
	nameTableChanged bool

	// Render related
	screen          *image.RGBA
	framePatternIDs [1024]byte // Screen representation with pattern ids and its position in screen. For debugging purposes.
}

func CreatePPU(cartridge *gamePak.GamePak) *Ppu2c02 {
	ppu := &Ppu2c02{
		cartridge:       cartridge,
		renderCycle:     0,
		currentScanline: 0,
		cycle:           0,
		vRam:            loopyRegister{0, 0},
		tRam:            loopyRegister{0, 0},
		fineX:           0,

		shifterTileLow:       0,
		shifterTileHigh:      0,
		shifterAttributeLow:  0,
		shifterAttributeHigh: 0,

		warmup: false,
		screen: image.NewRGBA(image.Rect(0, 0, types.SCREEN_WIDTH, types.SCREEN_HEIGHT)),
	}

	return ppu
}

func (ppu *Ppu2c02) Frame() *image.RGBA {
	return ppu.screen
}

func (ppu *Ppu2c02) FramePattern() *[1024]byte {
	return &ppu.framePatternIDs
}

func (ppu *Ppu2c02) Tick() {
	// VBlank logic
	if ppu.currentScanline == VBLANK_START_SCANLINE {
		if ppu.renderCycle == 1 {
			// TODO refactor to a method to set Vblank
			// TODO enabling VBlank only on ==1 and not on >=1 makes it difficult to start emulation inside a VBlank cycle. If changed, nmi triggering should be worked though.
			ppu.ppuStatus.verticalBlankStarted = true

			if ppu.ppuControl.generateNMIAtVBlank {
				ppu.nmi = true
			}
		}
	} else if ppu.currentScanline == VBLANK_END_SCNALINE && ppu.renderCycle == 1 {
		ppu.ppuStatus.verticalBlankStarted = false
	}

	// ------------------------------
	// Render logic
	ppu.renderLogic()

	//bit := ppu.registers.scrollX
	// Load new data into registers
	//if ppu.cycle%8 == 0 {
	//
	//}

	// Render logic end
	// ------------------------------

	// 341 PPU clock cycles have passed
	if ppu.renderCycle == PPU_CYCLES_BY_SCANLINE-1 {
		if ppu.currentScanline == 261 {
			ppu.currentScanline = 0
		} else {
			ppu.currentScanline++
		}
		ppu.renderCycle = 0
	} else {
		ppu.renderCycle++
	}

	if ppu.cycle >= PPU_CYCLES_TO_WARMUP {
		ppu.warmup = true
	} else {
		ppu.cycle++
	}
}

func (ppu *Ppu2c02) renderLogic() {
	//renderingEnabled := ppu.ppuMask.showBackground || ppu.ppuMask.showSprites
	preRenderScanline := ppu.currentScanline == 261
	scanlineVisible := ppu.currentScanline >= 0 && ppu.currentScanline < 240

	// We are in a cycle which falls inside the visible horizontal region
	cycleIsVisible := ppu.renderCycle >= 1 && ppu.renderCycle <= 256

	// On these cycles, we fetch data that will be used in next scanline
	preFetchCycle := ppu.renderCycle >= 321 && ppu.renderCycle <= 336

	//if ppu.ppuMask.
	//if !renderingEnabled {
	//	return
	//}

	if scanlineVisible {
		if ppu.renderCycle == 0 {
			// Idle cycle
		}
		if cycleIsVisible || preFetchCycle {
			ppu.updateShifters()

			switch ppu.renderCycle % 8 {
			case 0:
				// TODO feed appropiate shift registers
				ppu.loadShifters()
				ppu.incrementX()
			case 1:
				// fetch NameTable byte
				ppu.nextTileId = ppu.Read(ppu.vRam.address & 0xFFFF)
			case 3:
				// fetch attribute table byte
				address := types.Address(0x23C0)
				address |= types.Address(ppu.vRam.nameTableY()) << 11
				address |= types.Address(ppu.vRam.nameTableX()) << 10
				address |= types.Address(ppu.vRam.coarseX()>>2) << 3
				address |= types.Address(ppu.vRam.coarseY() >> 2)
				ppu.nextAttribute = ppu.Read(address)

				// I need to understand this
				if ppu.vRam.coarseY()&0x02 == 0x02 {
					ppu.nextAttribute >>= 4
				}
				if ppu.vRam.coarseX()&0x02 == 0x02 {
					ppu.nextAttribute >>= 2
				}
			case 5:
				// fetch low tile byte
				address := types.Address(0)
				address |= types.Address(ppu.ppuControl.backgroundPatternTableAddress) << 12
				address |= types.Address(ppu.nextTileId << 4)
				address |= types.Address(ppu.vRam.fineY())
				address |= types.Address(0)

				ppu.nextLowTile = ppu.Read(address)
			case 7:
				// fetch high tile byte
				address := types.Address(0)
				address |= types.Address(ppu.ppuControl.backgroundPatternTableAddress) << 12
				address |= types.Address(ppu.nextTileId << 4)
				address |= types.Address(ppu.vRam.fineY())
				address |= types.Address(8)

				ppu.nextLowTile = ppu.Read(address)
			}

			if ppu.renderCycle == 257 {
				ppu.incrementY()
			}
		}

		// When every pixel of a scanline has been rendered,
		// we need to reset the X coordinate
		if ppu.renderCycle == 258 {
			ppu.transferX()
		}

		if preRenderScanline && ppu.renderCycle >= 280 && ppu.renderCycle < 305 {
			ppu.transferY()
		}
	}

	if ppu.currentScanline == 240 {
		// idle PPU does nothing here
	}

	var bgPixel byte = 0x00
	var bgPalette byte = 0x00
	if ppu.ppuMask.showBackgroundEnabled() {
		bitSelector := uint16(0x8000) >> ppu.fineX
		pixel0 := byte(0)
		pixel1 := byte(0)
		if ppu.shifterTileLow&bitSelector > 0 {
			pixel0 = 1
		} else {
			pixel0 = 0
		}
		if ppu.shifterTileHigh&bitSelector > 0 {
			pixel1 = 1
		} else {
			pixel1 = 0
		}
		bgPixel = pixel1<<1 | pixel0

		palette0 := byte(0)
		palette1 := byte(0)
		if ppu.shifterAttributeLow&bitSelector > 0 {
			palette0 = 1
		} else {
			palette0 = 0
		}
		if ppu.shifterAttributeHigh&bitSelector > 0 {
			palette1 = 1
		} else {
			palette1 = 0
		}
		bgPalette = palette1<<1 | palette0
	}
	//fmt.Printf("Render %d,%d: %X %X\n", ppu.renderCycle, ppu.currentScanline, bgPalette, bgPixel)
	ppu.screen.Set(int(ppu.renderCycle), int(ppu.currentScanline), ppu.GetRGBColor(bgPalette, bgPixel))
}

func (ppu *Ppu2c02) incrementX() {
	if ppu.ppuMask.renderingEnabled() {
		if ppu.vRam.address&0x001F == 31 { // if coarseX == 31
			ppu.vRam.address &= 0b111111111100000 // coarseX = 0
			ppu.vRam.address ^= 0x0400            // switch horizontal nametable
		} else {
			ppu.vRam.address += 1 // coarseX++
		}
	}
}

func (ppu *Ppu2c02) incrementY() {
	if ppu.ppuMask.renderingEnabled() {
		if ppu.vRam.fineY() < 7 {
			ppu.vRam.address += 0x1000 // fineY++
		} else {
			ppu.vRam.address &= 0b000111111111111 // fineY=0
			y := ppu.vRam.coarseY()
			if y == 29 {
				y = 0
				ppu.vRam.address ^= 0x0800 // Switch vertical NameTable
			} else if y == 31 {
				y = 0 // coarseY = 0, NameTable not switched
			} else {
				y += 1
			}
			resetCoarseY := ^0x3e0
			ppu.vRam.address = (ppu.vRam.address & types.Address(resetCoarseY)) | types.Address(y)<<5
		}
	}
}

func (ppu *Ppu2c02) transferX() {
	if ppu.ppuMask.renderingEnabled() {
		ppu.vRam.setCoarseX(ppu.tRam.coarseX())
		ppu.vRam.setNameTableX(ppu.tRam.nameTableX())
	}
}

func (ppu *Ppu2c02) transferY() {
	if ppu.ppuMask.renderingEnabled() {
		ppu.vRam.setCoarseY(ppu.tRam.coarseY())
		ppu.vRam.setNameTableY(ppu.tRam.nameTableY())
	}
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
		panic("trying to read PPMASK")

	case PPUSTATUS:
		// Source: javid9x reading from status only get top 3 bits. The rest tends to be filled with noise, or more likely what was last in data buffer.
		value = ppu.ppuStatus.value()

		// Reading from status register alters it
		ppu.ppuStatus.verticalBlankStarted = false // Reading from status, clears VBlank flag.
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
		// Todo test delay and not delay from palette
		value = ppu.readBuffer
		ppu.readBuffer = ppu.Read(ppu.vRam.address)

		// If reading from Palette, there is no delay
		if isPaletteAddress(ppu.vRam.address) {
			value = ppu.readBuffer
		}

		ppu.vRam.increment(ppu.ppuControl.incrementMode)
		break

	case OAMDMA:
		break
	}

	return value
}

// Write made by CPU
func (ppu *Ppu2c02) WriteRegister(register types.Address, value byte) {
	if !ppu.warmup {
		return
	}

	switch register {
	case PPUCTRL:
		ppu.ppuCtrlWrite(value)
		ppu.tRam.setNameTableX(ppu.ppuControl.nameTableX)
		ppu.tRam.setNameTableY(ppu.ppuControl.nameTableY)
		// todo trigger nmi if in vblank and generateNMI transitions from 0 to 1
		break

	case PPUMASK:
		ppu.ppuMask.write(value)
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
		ppu.ppuScroll.write(value)
		break

	case PPUADDR:
		ppu.tRam.push(value)
		if ppu.tRam.latch == 0 {
			ppu.vRam = ppu.tRam
		}
		break
	case PPUDATA:
		ppu.Write(ppu.vRam.address, value)
		ppu.vRam.increment(ppu.ppuControl.incrementMode)
		break
	case OAMDMA:
		break
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
func (ppu *Ppu2c02) GetRGBColor(palette byte, colorIndex byte) color.Color {
	paletteColor := ppu.GetPaletteColor(palette, colorIndex)

	return utils.NewColorRGB(
		SystemPalette[paletteColor][0],
		SystemPalette[paletteColor][1],
		SystemPalette[paletteColor][2],
	)
}

func (ppu *Ppu2c02) GetPaletteColor(palette byte, colorIndex byte) byte {
	if palette > 0 && colorIndex == 0 {
		palette = 0
	}

	paletteAddress := types.Address((palette * 4) + colorIndex)
	return ppu.Read(PaletteLowAddress + paletteAddress)
}
