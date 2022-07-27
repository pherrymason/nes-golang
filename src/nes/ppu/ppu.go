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
	ppuControl     Control
	ppuStatus      Status
	ppuScroll      Scroll
	ppuMask        byte // Controls the rendering of sprites and backgrounds
	ppuDataAddress DataAddress

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
		warmup:          false,
		screen:          image.NewRGBA(image.Rect(0, 0, types.SCREEN_WIDTH, types.SCREEN_HEIGHT)),
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
			ppu.ppuStatus.verticalBlankStarted = true // Todo refactor to a method to set Vblank

			if ppu.ppuControl.generateNMIAtVBlank {
				//if ppu.renderCycle == 0 {
				ppu.nmi = true
				//}
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
	if ppu.renderCycle == 0 {
		// Idle
	}
	if ppu.renderCycle%1 == 0 {
		// Read 1 byte nametable
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
		break

	case PPUSTATUS:
		// Source: javid9x reading from status only get top 3 bits. The rest tends to be filled with noise, or more likely what was last in data buffer.
		value = ppu.ppuStatus.value()

		// Reading from status register alters it
		ppu.ppuStatus.verticalBlankStarted = false // Reading from status, clears VBlank flag.
		//ppu.registers.status &= 0x7F
		ppu.ppuDataAddress.resetLatch()
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
		value = ppu.ppuDataAddress.readBuffer
		ppu.ppuDataAddress.readBuffer = ppu.Read(ppu.ppuDataAddress.at())

		// If reading from Palette, there is no delay
		if isPaletteAddress(ppu.ppuDataAddress.at()) {
			value = ppu.ppuDataAddress.readBuffer
		}

		ppu.ppuDataAddress.increment(ppu.ppuControl.incrementMode)
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
		// todo trigger nmi if in vblank and generateNMI transitions from 0 to 1
		break

	case PPUMASK:
		ppu.ppuMask = value
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
		ppu.ppuDataAddress.set(value)
		break
	case PPUDATA:
		ppu.Write(ppu.ppuDataAddress.at(), value)
		ppu.ppuDataAddress.increment(ppu.ppuControl.incrementMode)
		//ppu.Write(ppu.registers.ppuDataAddr, value)
		/*
			if ppu.ppuControl.incrementMode == 0 {
				ppu.registers.ppuDataAddr++
			} else {
				ppu.registers.ppuDataAddr += 32
			}
			ppu.registers.ppuDataAddr &= 0x3FFF*/
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
func (ppu *Ppu2c02) GetColorFromPaletteRam(palette byte, colorIndex byte) color.Color {
	paletteColor := ppu.GetNesColorFromPaletteRam(palette, colorIndex)

	return utils.NewColorRGB(
		SystemPalette[paletteColor][0],
		SystemPalette[paletteColor][1],
		SystemPalette[paletteColor][2],
	)
}

func (ppu *Ppu2c02) GetNesColorFromPaletteRam(palette byte, colorIndex byte) byte {
	if palette > 0 && colorIndex == 0 {
		palette = 0
	}

	paletteAddress := types.Address((palette * 4) + colorIndex)
	return ppu.Read(PaletteLowAddress + paletteAddress)
}
