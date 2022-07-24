package ppu

import (
	"image"
)

func (ppu *Ppu2c02) PatternTable(patternTable byte, palette byte) image.RGBA {
	const CanvasWIDTH = 128
	chr := image.NewRGBA(image.Rect(0, 0, CanvasWIDTH, 8*8))
	//chr := make([]color.Color, CanvasWIDTH*128)

	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			tileID := byte(tileY*16 + tileX)
			tile := ppu.findTile(tileID, patternTable, 0, 0, palette)

			insertImageAt(
				chr,
				&tile,
				tileX*8,
				tileY*8,
			)
			//if tileID == 255 {
			//	saveTile(int(tileID), &tile)
			//	saveTile(300, chr)
			//}

			//for i := 0; i < 8*8; i++ {
			//index := types.CoordinatesToArrayIndex(
			//	tileX*8+(i%8),
			//	tileY*8+(i/8),
			//	CanvasWIDTH)
			//chr[index] = tile.At(tileX, tileY)

			//chr.Set(
			//	tileX*8+(i%8),
			//	tileY*8+(i/8),
			//	tile.At(tileX, tileY),
			//)
			//}
		}
	}

	//saveTile(300, chr)
	return *chr
}
