package ppu

import "github.com/raulferras/nes-golang/src/nes/types"

func (ppu *Ppu2c02) Render() *types.Frame {
	// Render first name table
	//bankAddress := types.Address(1 * 0x1000)
	nameTableStart := 0
	nameTablesEnd := int(PPU_NAMETABLES_0_END - PPU_NAMETABLES_0_START)
	tilesWidth := 32
	//tilesHeight := 30
	for addr := nameTableStart; addr < nameTablesEnd; addr++ {
		tileID := ppu.memory.vram[addr]

		tileX := addr % tilesWidth
		tileY := addr / tilesWidth

		//tileAddressInPatternTable := int(tileID) * 16
		//tiles := ppu.Read([tileID : tileID+16]
		ppu.framePattern[addr] = tileID
		bankAddress := types.Address(0x0 + int(tileID)*16)
		for y := 0; y <= 7; y++ {
			upper := ppu.Read(bankAddress + types.Address(y))
			lower := ppu.Read(bankAddress + types.Address(y+8))

			for x := 0; x <= 7; x++ {
				value := (1&upper)<<1 | (1 & lower)
				upper = upper >> 1
				lower = lower >> 1
				rgb := ppu.GetColorFromPaletteRam(0, value)

				ppu.frame.SetPixel(tileX*8+(7-x), tileY*8+y, rgb)
			}
		}
	}

	return &ppu.frame
}
