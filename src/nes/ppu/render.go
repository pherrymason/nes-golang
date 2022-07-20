package ppu

import "github.com/raulferras/nes-golang/src/nes/types"

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

		ppu.frame.PushTile(tile, tileX*8, tileY*8)
		ppu.framePattern[addr] = tileID
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
		ppu.frame.PushTile(tile, int(xCoordinate), int(yCoordinate))
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

func (ppu *Ppu2c02) findTile(tileID byte, patternTable byte) types.Tile {
	//patternTable := ppu.ppuctrlReadFlag(spritePatternTableAddress)
	bankAddress := 0x1000 * int(patternTable)
	offsetAddress := types.Address(bankAddress + int(tileID)*16)
	tile := types.Tile{}
	for y := 0; y <= 7; y++ {
		upper := ppu.Read(offsetAddress + types.Address(y))
		lower := ppu.Read(offsetAddress + types.Address(y+8))

		for x := 0; x <= 7; x++ {
			value := (1&upper)<<1 | (1 & lower)
			upper >>= 1
			lower >>= 1
			palette := byte(0) //backgroundPalette(tileX, tileY, ppu.memory.vram)
			rgb := ppu.GetColorFromPaletteRam(palette, value)
			index := types.CoordinatesToArrayIndex(7-x, y, types.TILE_WIDTH)
			tile.Pixels[index] = rgb
		}
	}

	return tile
}
