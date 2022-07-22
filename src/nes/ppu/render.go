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
}

func (ppu *Ppu2c02) renderBackground() {
	// Render first name table
	//bankAddress := types.Address(1 * 0x1000)
	nameTableStart := 0
	nameTablesEnd := int(PPU_NAMETABLES_0_END - PPU_NAMETABLES_0_START)
	tilesWidth := 32
	backgroundPatternTable := ppu.ppuctrlReadFlag(backgroundPatternTableAddress)
	//bankAddress := 0x1000 * int(backgroundPatternTable)
	//tilesHeight := 30
	for addr := nameTableStart; addr < nameTablesEnd; addr++ {
		tileID := ppu.memory.vram[addr]
		tile := ppu.findTile(tileID, backgroundPatternTable)

		tileX := addr % tilesWidth
		tileY := addr / tilesWidth

		//ppu.deprecatedFrame.PushTile(tile, tileX*8, tileY*8)
		ppu.renderTile(tile, tileX, tileY)
		ppu.framePattern[addr] = tileID
	}
}

func (ppu *Ppu2c02) renderTile(tile image.RGBA, coordX int, coordY int) {
	//ppu.screen.Set()
	//baseY := coordY * 256
	//baseX := coordX
	for i := 0; i < types.TILE_PIXELS; i++ {
		//calculatedY := baseY + (i/8)*types.SCREEN_WIDTH
		//calculatedX := baseX + i%8
		//arrayIndex := calculatedX + calculatedY
		//frame.Pixels[arrayIndex] = tile.Pixels[i]
		ppu.screen.Set(coordX, coordY, tile.At(i/types.TILE_WIDTH, i%types.TILE_WIDTH))
	}
}

func (ppu *Ppu2c02) renderSprites() {
	spritePatternTable := ppu.ppuctrlReadFlag(spritePatternTableAddress)
	for i := 0; i < OAMDATA_SIZE; i++ {
		yCoordinate := ppu.oamData[i]
		xCoordinate := ppu.oamData[i+1]
		tileID := ppu.oamData[i+2]
		//renderOption := ppu.oamData[i+3]

		tile := ppu.findTile(tileID, spritePatternTable)

		// Copy tile into frameSprites
		//ppu.deprecatedFrame.PushTile(tile, int(xCoordinate), int(yCoordinate))
		ppu.renderTile(tile, int(xCoordinate), int(yCoordinate))
		i += 3
	}
}

// Finds the palette id to be used given a background Tile coordinate
func backgroundPalette(x int, y int, vram [2048]byte) byte {
	metaTilesByRow := 8
	attributeTableAddress := 0x03C0
	attributeTableIndex := (y/4)*metaTilesByRow + x/4

	attrValue := vram[attributeTableAddress+attributeTableIndex]

	a := x % 4 / 2
	b := y % 4 / 2

	if a == 0 && b == 0 {
		return attrValue & 0b11
	} else if a == 1 && b == 0 {
		return (attrValue >> 2) & 0b11
	} else if a == 0 && b == 1 {
		return (attrValue >> 4) & 0b11
	} else if a == 1 && b == 1 {
		return (attrValue >> 6) & 0b11
	}

	panic("Invalid!")
}

func (ppu *Ppu2c02) findTile(tileID byte, patternTable byte) image.RGBA {
	bankAddress := 0x1000 * int(patternTable)
	offsetAddress := types.Address(bankAddress + int(tileID)*16)
	tile := image.NewRGBA(image.Rect(0, 0, types.TILE_WIDTH, types.TILE_HEIGHT))
	for y := 0; y <= 7; y++ {
		upper := ppu.Read(offsetAddress + types.Address(y))
		lower := ppu.Read(offsetAddress + types.Address(y+8))

		for x := 0; x <= 7; x++ {
			value := (1&upper)<<1 | (1 & lower)
			upper >>= 1
			lower >>= 1
			palette := byte(0) //backgroundPalette(tileX, tileY, ppu.memory.vram)
			rgb := ppu.GetColorFromPaletteRam(palette, value)
			tile.Set(7-x, y, rgb)
		}
	}
	//saveTile(int(tileID), tile)
	return *tile
}

func SaveTile(tileID int, tile *image.RGBA) {
	//myImage := image.NewRGBA(image.Rect(0, 0, 100, 200))

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
