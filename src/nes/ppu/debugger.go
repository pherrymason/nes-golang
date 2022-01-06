package ppu

import (
	"github.com/raulferras/nes-golang/src/graphics"
	"github.com/raulferras/nes-golang/src/nes/types"
)

func (ppu *Ppu2c02) PatternTable(patternTable int, palette uint8) []graphics.Pixel {
	chr := make([]graphics.Pixel, 128*128)

	bankAddress := types.Address(patternTable * 0x1000)
	i := 0
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
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

					pixel := graphics.Pixel{
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

func (ppu *Ppu2c02) getTile(patternTable int, palette uint8, tileX int, tileY int) []graphics.Pixel {
	tile := make([]graphics.Pixel, 8*8)
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

			pixel := graphics.Pixel{
				X:     coordX,
				Y:     coordY,
				Color: ppu.GetColorFromPaletteRam(palette, value),
			}
			tile[row+col*8] = pixel
		}
	}

	return tile
}
