package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

func (ppu *Ppu2c02) PatternTable(patternTable int, palette uint8) []types.Pixel {
	chr := make([]types.Pixel, 128*128)

	bankAddress := types.Address(patternTable * 0x1000)
	i := 0
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			// Calculate address where the tile is located in memory
			// each row of CHR contains 256 bytes. So current row * 256 bytes = tileY*256
			// each tile of CHR uses 16 bytes. So current column *16 bytes = (tileX*16)
			offset := bankAddress + types.Address(tileY*256+tileX*16)

			for row := 0; row < 8; row++ {
				pixelLSB := ppu.Read(offset + types.Address(row))
				pixelMSB := ppu.Read(offset + types.Address(row+8))

				for col := 0; col < 8; col++ {
					value := (pixelLSB & 0x01) | ((pixelMSB & 0x01) << 1)

					pixelLSB >>= 1
					pixelMSB >>= 1

					coordY := tileY*8 + row
					coordX := (7 - col) + tileX*8

					pixel := types.Pixel{
						X:     coordX,
						Y:     coordY,
						Color: ppu.GetColorFromPaletteRam(palette, value),
					}
					chr[i] = pixel
					i++
				}
			}
		}
	}

	return chr
}

func (ppu *Ppu2c02) getTile(patternTable int, palette uint8, tileX int, tileY int) [types.TILE_PIXELS]types.Pixel {
	tile := [types.TILE_PIXELS]types.Pixel{}
	bankAddress := types.Address(patternTable * 0x1000)

	offset := bankAddress + types.Address(tileY*256+tileX*16)

	for row := 0; row < 8; row++ {
		pixelLSB := ppu.Read(offset + types.Address(row))
		pixelMSB := ppu.Read(offset + types.Address(row+8))

		for col := 0; col < 8; col++ {
			value := (pixelLSB & 0x01) | ((pixelMSB & 0x01) << 1)

			pixelLSB >>= 1
			pixelMSB >>= 1

			coordY := tileY*8 + row
			coordX := (7 - col) + tileX*8

			pixel := types.Pixel{
				X:     coordX,
				Y:     coordY,
				Color: ppu.GetColorFromPaletteRam(palette, value),
			}
			tile[row+col*8] = pixel
		}
	}

	return tile
}
