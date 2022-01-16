package ppu

import "github.com/raulferras/nes-golang/src/nes/types"

func (ppu *Ppu2c02) Render() *types.Frame {
	// Render first name table
	//bankAddress := types.Address(1 * 0x1000)
	nameTableStart := 0
	nameTablesEnd := int(PPU_NAMETABLES_0_END - PPU_NAMETABLES_0_START)
	tilesWidth := 32
	backgroudPatternTable := ppu.ppuctrlReadFlag(backgroundPatternTableAddress)
	bankAddress := 0x1000 * int(backgroudPatternTable)
	//tilesHeight := 30
	for addr := nameTableStart; addr < nameTablesEnd; addr++ {
		tileID := ppu.memory.vram[addr]

		tileX := addr % tilesWidth
		tileY := addr / tilesWidth

		//tileAddressInPatternTable := int(tileID) * 16
		//tiles := ppu.Read([tileID : tileID+16]
		ppu.framePattern[addr] = tileID
		offsetAddress := types.Address(bankAddress + int(tileID)*16)
		for y := 0; y <= 7; y++ {
			upper := ppu.Read(offsetAddress + types.Address(y))
			lower := ppu.Read(offsetAddress + types.Address(y+8))

			for x := 0; x <= 7; x++ {
				value := (1&upper)<<1 | (1 & lower)
				upper = upper >> 1
				lower = lower >> 1
				palette := backgroundPalette(tileX, tileY, ppu.memory.vram)
				rgb := ppu.GetColorFromPaletteRam(palette, value)

				ppu.frame.SetPixel(tileX*8+(7-x), tileY*8+y, rgb)
			}
		}
	}

	return &ppu.frame
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
