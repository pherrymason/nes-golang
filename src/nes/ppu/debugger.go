package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

func (ppu *Ppu2c02) PatternTable(patternTable byte, palette byte) []types.Color {
	const CanvasWIDTH = 128
	chr := make([]types.Color, CanvasWIDTH*128)

	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			tileID := byte(tileY*16 + tileX)
			tile := ppu.findTile(tileID, patternTable)

			for i := 0; i < 8*8; i++ {
				index := types.CoordinatesToArrayIndex(
					tileX*8+(i%8),
					tileY*8+(i/8),
					CanvasWIDTH)
				chr[index] = tile.Pixels[i]
			}
		}
	}

	return chr
}
