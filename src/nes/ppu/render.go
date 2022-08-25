package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
	"image/png"
	"math/bits"
	"os"
)

func (ppu *P2c02) Render() {
	if ppu.renderByPixel == false {
		ppu.renderBackground()
		//ppu.renderSprites()

		ppu.nameTableChanged = false
	}
}

func (ppu *P2c02) scanlineIsVisibleOrIsPreRender() bool {
	scanlineVisible := ppu.currentScanline < 240
	preRenderScanline := ppu.currentScanline == 261

	return scanlineVisible || preRenderScanline
}

func (ppu *P2c02) renderLogic() {
	//renderingEnabled := ppu.PpuMask.ShowBackground || ppu.PpuMask.ShowSprites
	preRenderScanline := ppu.currentScanline == 261
	scanlineVisible := ppu.currentScanline < 240

	// We are in a cycle which falls inside the visible horizontal region
	cycleIsVisible := ppu.renderCycle >= 1 && ppu.renderCycle <= 256

	// On these cycles, we fetch data that will be used in next Scanline
	preFetchCycle := ppu.renderCycle >= 321 && ppu.renderCycle <= 336

	if scanlineVisible || preRenderScanline {
		if ppu.renderCycle == 0 {
			// Idle cycle
			ppu.renderCycle = 0
		}

		// Horizontally...
		if cycleIsVisible || preFetchCycle {
			ppu.updateShifters()

			switch ppu.renderCycle % 8 {
			case 0:
				ppu.loadShifters()
				ppu.incrementX()
			case 1:
				// fetch NameTable byte
				address := 0x2000 | ppu.vRam.nameTableAddress()
				ppu.bgNextTileId = ppu.Read(address)
			case 3:
				// fetch attribute table byte
				// "(vram_addr.coarse_x >> 2)"        : integer divide coarse x by 4,
				//                                      from 5 bits to 3 bits
				// "((vram_addr.coarse_y >> 2) << 3)" : integer divide coarse y by 4,
				//                                      from 5 bits to 3 bits,
				//                                      shift to make room for coarse x

				address := types.Address(0x23C0)
				address |= types.Address(ppu.vRam.NameTableY()) << 11
				address |= types.Address(ppu.vRam.NameTableX()) << 10
				address |= types.Address(ppu.vRam.CoarseX() >> 2)    // Divide coarse by 4
				address |= types.Address(ppu.vRam.CoarseY()>>2) << 3 // Divide coarse by 4, shift to make space to CoarseX
				ppu.bgNextAttribute = ppu.Read(address)

				// We got the right attribute byte, but we need to find the right 2bit section corresponding
				// to the tile we are processing
				if ppu.vRam.CoarseY()&0x02 == 0x02 {
					ppu.bgNextAttribute >>= 4
				}
				if ppu.vRam.CoarseX()&0x02 == 0x02 {
					ppu.bgNextAttribute >>= 2
				}
				ppu.bgNextAttribute &= 0x03
			case 5:
				// fetch low tile byte
				// "(ppu.PpuControl.BackgroundPatternTableAddress << 12)"  : the pattern memory selector
				//                                         from control register, either 0K
				//                                         or 4K offset
				// "(ppu.bgNextTileId << 4)"    : the tile id multiplied by 16, as
				//                                         2 lots of 8 rows of 8 bit pixels
				// "(ppu.vRam.FineY)"                  : Offset into which row based on
				//                                         vertical scroll offset
				// "+ 0"                                 : Mental clarity for plane offset
				address := types.Address(ppu.PpuControl.BackgroundPatternTableAddress) << 12
				address |= types.Address(ppu.bgNextTileId) << 4
				address |= types.Address(ppu.vRam.FineY())
				address |= types.Address(0)

				ppu.bgNextLowTile = ppu.Read(address)
			case 7:
				// fetch high tile byte
				address := types.Address(ppu.PpuControl.BackgroundPatternTableAddress) << 12
				address |= types.Address(ppu.bgNextTileId) << 4
				address |= types.Address(ppu.vRam.FineY())
				address |= types.Address(8)

				ppu.bgNextHighTile = ppu.Read(address)
			}
		} // horizontal cycle visible|prefecth check

		if ppu.renderCycle == 256 {
			ppu.incrementY()
		}

		// When every pixel of a Scanline has been rendered,
		// we need to reset the X coordinate
		if ppu.renderCycle == 257 {
			ppu.transferX()
		}

		if preRenderScanline && ppu.renderCycle >= 280 && ppu.renderCycle < 305 {
			ppu.transferY()
		}

		// Sprites stuff -------------------
		if ppu.renderCycle == 257 && scanlineVisible {
			for i := 0; i < 8; i++ {
				ppu.oamDataScanline[i] = objectAttributeEntry{0, 0, 0, 0}
				ppu.spShifterPatternLow[i] = 0x00
				ppu.spShifterPatternHigh[i] = 0x00
			}
		}

		if ppu.renderCycle == 257 {
			// Pre-fetching sprites for next Scanline
			// This is not fully accurate, sprite loading occurs along different cycles.
			ppu.spriteScanlineCount = 0
			var spriteHeight int16
			if ppu.PpuControl.SpriteSize == 0 {
				spriteHeight = 8
			} else {
				spriteHeight = 16
			}

			for oamIndex := 0; oamIndex < OAMDATA_SIZE && ppu.spriteScanlineCount < 8; oamIndex += 4 {
				spriteY := int16(ppu.oamData[oamIndex])
				diff := ppu.currentScanline - Scanline(spriteY)
				if diff >= 0 && diff <= Scanline(spriteHeight) {
					ppu.oamDataScanline[ppu.spriteScanlineCount] = objectAttributeEntry{
						y:          byte(spriteY),
						tileId:     ppu.oamData[oamIndex+1],
						attributes: ppu.oamData[oamIndex+2],
						x:          ppu.oamData[oamIndex+3],
					}
					ppu.spriteScanlineCount++
				}
			}
		}

		if ppu.renderCycle == 340 {
			// Get sprite pattern information
			// Doing this at cycle 340 is a simplification of what the real NES actually does.
			// More info: https://www.nesdev.org/wiki/PPU_sprite_evaluation
			ppu.fetchSpriteShifters()
		}
		// ---------------------------------
	}

	if ppu.currentScanline == 240 {
		// idle PPU does nothing here
	}

	ppu.finalPixelComposition()
}

func (ppu *P2c02) finalPixelComposition() {
	var bgPixel byte = 0x00
	var bgPalette byte = 0x00
	if ppu.PpuMask.showBackgroundEnabled() {
		bitSelector := uint16(0x8000) >> ppu.fineX
		pixel0 := byte(0)
		pixel1 := byte(0)
		if ppu.bgShifterTileLow&bitSelector > 0 {
			pixel0 = 1
		} else {
			pixel0 = 0
		}
		if ppu.bgShifterTileHigh&bitSelector > 0 {
			pixel1 = 1
		} else {
			pixel1 = 0
		}
		bgPixel = pixel1<<1 | pixel0

		palette0 := byte(0)
		palette1 := byte(0)
		if ppu.bgShifterAttributeLow&bitSelector > 0 {
			palette0 = 1
		} else {
			palette0 = 0
		}
		if ppu.bgShifterAttributeHigh&bitSelector > 0 {
			palette1 = 1
		} else {
			palette1 = 0
		}
		bgPalette = palette1<<1 | palette0
	}

	var fgPixel byte
	var fgPalette byte
	//var fgPriority byte //todo implement priority
	if ppu.PpuMask.showSpritesEnabled() {
		for i := byte(0); i < ppu.spriteScanlineCount; i++ {
			if ppu.oamDataScanline[i].x > 0 {
				continue
			}

			fgPixelLow := (ppu.spShifterPatternLow[i] & 0x80) >> 7
			fgPixelHigh := (ppu.spShifterPatternHigh[i] & 0x80) >> 7
			fgPixel = fgPixelHigh<<1 | fgPixelLow

			fgPalette = ppu.oamDataScanline[i].palette()
			if fgPixel != 0 {
				break
			}
		}
	}

	var finalPixel byte
	var finalPalette byte
	if fgPixel != 0 {
		finalPixel = fgPixel
		finalPalette = fgPalette
	} else {
		finalPixel = bgPixel
		finalPalette = bgPalette
	}

	if ppu.renderByPixel {
		ppu.screen.Set(
			int(ppu.renderCycle-1),
			int(ppu.currentScanline),
			ppu.GetRGBColor(finalPalette, finalPixel))
	}
}

// updateShifters
// This method shifts one bit to the left the contents of the shifter registers.
// This, together with the fineX register allows to get the pixel information
// to be rendered together with a smooth per pixel scrolling
func (ppu *P2c02) updateShifters() {
	if ppu.PpuMask.showBackgroundEnabled() {
		ppu.bgShifterTileLow <<= 1
		ppu.bgShifterTileHigh <<= 1
		ppu.bgShifterAttributeLow <<= 1
		ppu.bgShifterAttributeHigh <<= 1
	}

	if ppu.PpuMask.showSpritesEnabled() && ppu.renderCycle >= 1 && ppu.renderCycle < 258 {
		for i := byte(0); i < ppu.spriteScanlineCount; i++ {
			if ppu.oamDataScanline[i].x > 0 {
				ppu.oamDataScanline[i].x--
			} else {
				ppu.spShifterPatternLow[i] <<= 1
				ppu.spShifterPatternHigh[i] <<= 1
			}
		}
	}
}

// loadShifters
// This prepares shifters with the next tile to be rendered
func (ppu *P2c02) loadShifters() {
	ppu.bgShifterTileLow = (ppu.bgShifterTileLow & 0xFF00) | uint16(ppu.bgNextLowTile)
	ppu.bgShifterTileHigh = (ppu.bgShifterTileHigh & 0xFF00) | uint16(ppu.bgNextHighTile)

	// As only two bits are required for the color index,
	// we will use the strategy of shifting bits from a 16bit Value.
	// In this case, we will repeat low bit through 16 bits
	// and the same with high bit
	if ppu.bgNextAttribute&0b01 == 1 {
		ppu.bgShifterAttributeLow = (ppu.bgShifterAttributeLow & 0xFF00) | 0xFF
	} else {
		ppu.bgShifterAttributeLow = (ppu.bgShifterAttributeLow & 0xFF00) | 0x00
	}

	if ppu.bgNextAttribute&0b10 == 2 {
		ppu.bgShifterAttributeHigh = (ppu.bgShifterAttributeHigh & 0xFF00) | 0xFF
	} else {
		ppu.bgShifterAttributeHigh = (ppu.bgShifterAttributeHigh & 0xFF00) | 0x00
	}
}

func (ppu *P2c02) fetchSpriteShifters() {
	var spritePatternAddressLow types.Address
	var spritePatternAddressHigh types.Address
	var spritePatternLow byte
	var spritePatternHigh byte

	for i := byte(0); i < ppu.spriteScanlineCount; i++ {
		object := ppu.oamDataScanline[i]

		if ppu.PpuControl.SpriteSize == PPU_CONTROL_SPRITE_SIZE_8 {
			if !object.isFlippedVertically() {
				spritePatternAddressLow = types.Address(ppu.PpuControl.SpritePatternTableAddress) << 12
				spritePatternAddressLow |= types.Address(object.tileId) << 4                           // Multiply ID per 16 (16 bytes per tile)
				spritePatternAddressLow |= types.Address(ppu.currentScanline - Scanline(object.y) - 1) // Which tile line we want
			} else {
				// TODO implement vertical flip
			}
		} else {
			// TODO implement 8x16 sprites
		}

		spritePatternAddressHigh = spritePatternAddressLow + 8
		spritePatternLow = ppu.Read(spritePatternAddressLow)
		spritePatternHigh = ppu.Read(spritePatternAddressHigh)

		// TODO check for horizontal flip
		if object.isFlippedHorizontally() {
			// todo implement
			spritePatternLow = bits.Reverse8(spritePatternLow)
			spritePatternHigh = bits.Reverse8(spritePatternHigh)
		}

		ppu.spShifterPatternLow[i] = spritePatternLow
		ppu.spShifterPatternHigh[i] = spritePatternHigh
	}
}

func (ppu *P2c02) renderBackground() {
	// Render first name table
	//bankAddress := types.Address(1 * 0x1000)
	nameTableStart := 0
	nameTablesEnd := int(PPU_NAMETABLES_0_END - NameTableStartAddress)
	tilesWidth := 32
	backgroundPatternTable := ppu.PpuControl.BackgroundPatternTableAddress
	//bankAddress := 0x1000 * int(backgroundPatternTable)
	//tilesHeight := 30
	if !ppu.nameTableChanged {
		return
	}
	for addr := nameTableStart; addr < nameTablesEnd; addr++ {
		tileID := ppu.nameTables[addr]
		tileX := addr % tilesWidth
		tileY := addr / tilesWidth
		tile := ppu.findTile(tileID, backgroundPatternTable, uint8(tileX), uint8(tileY), 255)

		insertImageAt(ppu.screen, &tile, tileX*8, tileY*8)
		//SaveTile(999+addr, ppu.screen)
		//ppu.renderTile(tile, tileX, tileY)
		//ppu.framePatternIDs[addr] = tileID
	}
}

func (ppu *P2c02) renderTile(tile image.RGBA, coordX int, coordY int) {
	//ppu.screen.Set()
	//baseY := coordY * 256
	//baseX := coordX
	for i := 0; i < TILE_PIXELS; i++ {
		//calculatedY := baseY + (i/8)*types.SCREEN_WIDTH
		//calculatedX := baseX + i%8
		//arrayIndex := calculatedX + calculatedY
		//Frame.Pixels[arrayIndex] = tile.Pixels[i]
		ppu.screen.Set(coordX, coordY, tile.At(i/TILE_WIDTH, i%TILE_WIDTH))
	}
}

func (ppu *P2c02) renderSprites() {
	spritePatternTable := ppu.PpuControl.SpritePatternTableAddress
	for i := 0; i < OAMDATA_SIZE; i++ {
		yCoordinate := ppu.oamData[i]
		xCoordinate := ppu.oamData[i+1]
		tileID := ppu.oamData[i+2]
		//renderOption := ppu.oamData[i+3]

		tile := ppu.findTile(tileID, spritePatternTable, 0, 0, 255)

		// Copy tile into frameSprites
		//ppu.deprecatedFrame.PushTile(tile, int(xCoordinate), int(yCoordinate))
		ppu.renderTile(tile, int(xCoordinate), int(yCoordinate))
		i += 3
	}
}

// Finds the palette id to be used given a background Tile coordinate
func backgroundPalette(tileColumn uint8, tileRow uint8, nameTable *[2 * NAMETABLE_SIZE]byte) byte {
	metaTilesByRow := uint8(8)
	attributeTableAddress := types.Address(0x23C0 - 0x2000)

	attributeTableIndex := (tileRow/4)*metaTilesByRow + tileColumn/4

	attrValue := nameTable[attributeTableAddress+types.Address(attributeTableIndex)]

	// Each byte controls the palette of a 32×32 pixel or 4×4 tile part of the nametable and is divided into four 2-bit areas
	// 7654 3210
	// |||| ||++- Color bits 1-0 for top left quadrant of this byte
	// |||| ++--- Color bits 3-2 for top right quadrant of this byte
	// ||++------ Color bits 5-4 for bottom left quadrant of this byte
	// ++-------- Color bits 7-6 for bottom right quadrant of this byte
	//
	//	+-------+
	//  | 0 | 1 |
	//	+---+---+
	//  | 2 | 3 |
	//  +---+---+

	a := tileColumn % 4 / 2
	b := tileRow % 4 / 2

	if a == 0 && b == 0 {
		return attrValue & 0b11
	} else if a == 1 && b == 0 {
		return (attrValue >> 2) & 0b11
	} else if a == 0 && b == 1 {
		return (attrValue >> 4) & 0b11
	} else if a == 1 && b == 1 {
		return (attrValue >> 6) & 0b11
	}

	panic("backgroundPalette: Invalid attribute Value!")
}

func (ppu *P2c02) findTile(tileID byte, patternTable byte, tileColumn uint8, tileRow uint8, forcedPalette uint8) image.RGBA {
	bankAddress := 0x1000 * int(patternTable)
	offsetAddress := types.Address(bankAddress + int(tileID)*16)
	tile := image.NewRGBA(image.Rect(0, 0, TILE_WIDTH, TILE_HEIGHT))

	var palette byte
	if forcedPalette != 255 {
		palette = forcedPalette
	} else {
		palette = backgroundPalette(tileColumn, tileRow, &ppu.nameTables)
	}

	for y := 0; y <= 7; y++ {
		lower := ppu.Read(offsetAddress + types.Address(y))
		upper := ppu.Read(offsetAddress + types.Address(y+8))

		for x := 0; x <= 7; x++ {
			value := (1&upper)<<1 | (1 & lower)
			// todo Should take transparency into account
			rgb := ppu.GetRGBColor(palette, value)
			tile.Set(7-x, y, rgb)
			upper >>= 1
			lower >>= 1
		}
	}
	//SaveTile(int(tileID), tile)
	return *tile
}

func SaveTile(fileName string, tile *image.RGBA) {
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create(fileName)
	if err != nil {
		// Handle error
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, tile)

	// Don't forget to close files
	outputFile.Close()
}

func insertImageAt(canvas *image.RGBA, sprite *image.RGBA, x int, y int) {
	bounds := sprite.Bounds()
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			canvas.Set(x+i, y+j, sprite.At(i, j))
		}
	}
}
