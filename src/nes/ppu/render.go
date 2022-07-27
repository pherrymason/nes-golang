package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
	"image/png"
	"os"
)

func (ppu *Ppu2c02) Render() {
	ppu.renderBackground()
	//ppu.renderSprites()

	ppu.nameTableChanged = false
}

func (ppu *Ppu2c02) renderBackground() {
	// Render first name table
	//bankAddress := types.Address(1 * 0x1000)
	nameTableStart := 0
	nameTablesEnd := int(PPU_NAMETABLES_0_END - NameTableStartAddress)
	tilesWidth := 32
	backgroundPatternTable := ppu.ppuControl.backgroundPatternTableAddress
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

		//ppu.renderTile(tile, tileX, tileY)
		//ppu.framePatternIDs[addr] = tileID
	}
}

func (ppu *Ppu2c02) renderTile(tile image.RGBA, coordX int, coordY int) {
	//ppu.screen.Set()
	//baseY := coordY * 256
	//baseX := coordX
	for i := 0; i < TILE_PIXELS; i++ {
		//calculatedY := baseY + (i/8)*types.SCREEN_WIDTH
		//calculatedX := baseX + i%8
		//arrayIndex := calculatedX + calculatedY
		//frame.Pixels[arrayIndex] = tile.Pixels[i]
		ppu.screen.Set(coordX, coordY, tile.At(i/TILE_WIDTH, i%TILE_WIDTH))
	}
}

func (ppu *Ppu2c02) renderSprites() {
	spritePatternTable := ppu.ppuControl.spritePatterTableAddress
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

	panic("backgroundPalette: Invalid attribute value!")
}

func (ppu *Ppu2c02) findTile(tileID byte, patternTable byte, tileColumn uint8, tileRow uint8, forcedPalette uint8) image.RGBA {
	bankAddress := 0x1000 * int(patternTable)
	offsetAddress := types.Address(bankAddress + int(tileID)*16)
	tile := image.NewRGBA(image.Rect(0, 0, TILE_WIDTH, TILE_HEIGHT))

	var palette byte
	if forcedPalette != 255 {
		palette = forcedPalette
	} else {
		palette = backgroundPalette(tileColumn, tileRow, &ppu.nameTables)
	}

	if tileID == 0xA3 && tileColumn == 6 && tileRow == 4 {
		tileID += 1
		tileID -= 1
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
	//saveTile(int(tileID), tile)
	return *tile
}

func SaveTile(tileID int, tile *image.RGBA) {
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create(fmt.Sprintf("%d.png", tileID))
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
