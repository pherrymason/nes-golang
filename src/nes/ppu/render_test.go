package ppu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/utils"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

// Shifters test
func TestPpu2c02_updateShifters_updates_shifters_while_rendering_is_enabled(t *testing.T) {
	ppu := aPPU()
	ppu.ppuMask.showBackground = 1
	ppu.ppuMask.showSprites = 1
	ppu.bgShifterTileLow = 1
	ppu.bgShifterTileHigh = 1
	ppu.bgShifterAttributeLow = 1
	ppu.bgShifterAttributeHigh = 1

	ppu.updateShifters()

	assert.Equal(t, uint16(0x02), ppu.bgShifterTileLow)
	assert.Equal(t, uint16(0x02), ppu.bgShifterTileHigh)
	assert.Equal(t, uint16(0x02), ppu.bgShifterAttributeLow)
	assert.Equal(t, uint16(0x02), ppu.bgShifterAttributeHigh)
}

func TestPpu2c02_updateShifters_does_not_updates_shifters_while_rendering_is_disabled(t *testing.T) {
	ppu := aPPU()
	ppu.ppuMask.showBackground = 0
	ppu.ppuMask.showSprites = 0
	ppu.bgShifterTileLow = 1
	ppu.bgShifterTileHigh = 1
	ppu.bgShifterAttributeLow = 1
	ppu.bgShifterAttributeHigh = 1

	ppu.updateShifters()

	assert.Equal(t, uint16(0x01), ppu.bgShifterTileLow)
	assert.Equal(t, uint16(0x01), ppu.bgShifterTileHigh)
	assert.Equal(t, uint16(0x01), ppu.bgShifterAttributeLow)
	assert.Equal(t, uint16(0x01), ppu.bgShifterAttributeHigh)
}

func TestPpu2c02_loadShifters(t *testing.T) {
	ppu := aPPU()
	ppu.bgShifterTileLow = 0xCCAA
	ppu.bgShifterTileHigh = 0xCCAA
	ppu.bgShifterAttributeLow = 0xFFFF
	ppu.bgShifterAttributeHigh = 0x0000

	ppu.bgNextLowTile = 0xFA
	ppu.bgNextHighTile = 0xBB
	ppu.bgNextAttribute = 0xAA

	ppu.loadShifters()

	assert.Equal(t, uint16(0xCCFA), ppu.bgShifterTileLow, "tile low with unexpected Value")
	assert.Equal(t, uint16(0xCCBB), ppu.bgShifterTileHigh, "tile high with unexpected Value")
	assert.Equal(t, uint16(0xFF00), ppu.bgShifterAttributeLow, "attribute low with unexpected Value")
	assert.Equal(t, uint16(0x00FF), ppu.bgShifterAttributeHigh, "attribute low with unexpected Value")
}

func Test_getsTile(t *testing.T) {
	chrROM := make([]byte, 0x01FFF)
	// LSB of tile
	chrROM[0] = 0b00000111
	chrROM[1] = 0b00001111
	chrROM[2] = 0b00011001
	chrROM[3] = 0b00110000
	chrROM[4] = 0b01100011 //**
	chrROM[5] = 0b01110010
	chrROM[6] = 0b01110000
	chrROM[7] = 0b00000001

	// MSB of tile
	chrROM[8] = 0b00000111
	chrROM[9] = 0b00001111
	chrROM[10] = 0b00011111
	chrROM[11] = 0b00111111
	chrROM[12] = 0b11111100 //**
	chrROM[13] = 0b11111100
	chrROM[14] = 0b11111111
	chrROM[15] = 0b11111111

	cartridge := gamePak.NewDummyGamePak(chrROM)
	ppu := CreatePPU(cartridge, false, "")
	ppu.nameTables[0] = 0
	ppu.paletteTable[0] = 0x0F
	ppu.paletteTable[1] = 0x30
	ppu.paletteTable[2] = 0x36
	ppu.paletteTable[3] = 0x06

	renderedTile := ppu.findTile(0, 0, 0, 0, 0)

	expectedRenderedTile := expectedTile()

	if assert.Equal(t, *expectedRenderedTile, renderedTile, "generated tile is wrong. Check src/ppu/0.png & src/ppu/1.png") == false {
		SaveTile("0", expectedRenderedTile)
		SaveTile("1", &renderedTile)
	}
}

func expectedTile() *image.RGBA {
	//Palette Colors indexes f 30 36 06
	// 0x0F => {0,0,0}
	// 0x30 => {236, 238, 236}
	// 0x36 => {236, 180, 176},
	// 0x06 => {84, 4, 0},,
	black := utils.NewColorRGB(0, 0, 0)
	white := utils.NewColorRGB(236, 238, 236)
	skin := utils.NewColorRGB(236, 180, 176)
	brown := utils.NewColorRGB(84, 4, 0)
	expectedRenderedTile := image.NewRGBA(image.Rect(0, 0, TILE_WIDTH, TILE_HEIGHT))
	expectedRenderedTile.Set(0, 0, black)
	expectedRenderedTile.Set(1, 0, black)
	expectedRenderedTile.Set(2, 0, black)
	expectedRenderedTile.Set(3, 0, black)
	expectedRenderedTile.Set(4, 0, black)
	expectedRenderedTile.Set(5, 0, brown)
	expectedRenderedTile.Set(6, 0, brown)
	expectedRenderedTile.Set(7, 0, brown)

	expectedRenderedTile.Set(0, 1, black)
	expectedRenderedTile.Set(1, 1, black)
	expectedRenderedTile.Set(2, 1, black)
	expectedRenderedTile.Set(3, 1, black)
	expectedRenderedTile.Set(4, 1, brown)
	expectedRenderedTile.Set(5, 1, brown)
	expectedRenderedTile.Set(6, 1, brown)
	expectedRenderedTile.Set(7, 1, brown)

	expectedRenderedTile.Set(0, 2, black)
	expectedRenderedTile.Set(1, 2, black)
	expectedRenderedTile.Set(2, 2, black)
	expectedRenderedTile.Set(3, 2, brown)
	expectedRenderedTile.Set(4, 2, brown)
	expectedRenderedTile.Set(5, 2, skin)
	expectedRenderedTile.Set(6, 2, skin)
	expectedRenderedTile.Set(7, 2, brown)

	expectedRenderedTile.Set(0, 3, black)
	expectedRenderedTile.Set(1, 3, black)
	expectedRenderedTile.Set(2, 3, brown)
	expectedRenderedTile.Set(3, 3, brown)
	expectedRenderedTile.Set(4, 3, skin)
	expectedRenderedTile.Set(5, 3, skin)
	expectedRenderedTile.Set(6, 3, skin)
	expectedRenderedTile.Set(7, 3, skin)

	expectedRenderedTile.Set(0, 4, skin)
	expectedRenderedTile.Set(1, 4, brown)
	expectedRenderedTile.Set(2, 4, brown)
	expectedRenderedTile.Set(3, 4, skin)
	expectedRenderedTile.Set(4, 4, skin)
	expectedRenderedTile.Set(5, 4, skin)
	expectedRenderedTile.Set(6, 4, white)
	expectedRenderedTile.Set(7, 4, white)

	expectedRenderedTile.Set(0, 5, skin)
	expectedRenderedTile.Set(1, 5, brown)
	expectedRenderedTile.Set(2, 5, brown)
	expectedRenderedTile.Set(3, 5, brown)
	expectedRenderedTile.Set(4, 5, skin)
	expectedRenderedTile.Set(5, 5, skin)
	expectedRenderedTile.Set(6, 5, white)
	expectedRenderedTile.Set(7, 5, black)

	expectedRenderedTile.Set(0, 6, skin)
	expectedRenderedTile.Set(1, 6, brown)
	expectedRenderedTile.Set(2, 6, brown)
	expectedRenderedTile.Set(3, 6, brown)
	expectedRenderedTile.Set(4, 6, skin)
	expectedRenderedTile.Set(5, 6, skin)
	expectedRenderedTile.Set(6, 6, skin)
	expectedRenderedTile.Set(7, 6, skin)

	expectedRenderedTile.Set(0, 7, skin)
	expectedRenderedTile.Set(1, 7, skin)
	expectedRenderedTile.Set(2, 7, skin)
	expectedRenderedTile.Set(3, 7, skin)
	expectedRenderedTile.Set(4, 7, skin)
	expectedRenderedTile.Set(5, 7, skin)
	expectedRenderedTile.Set(6, 7, skin)
	expectedRenderedTile.Set(7, 7, brown)
	return expectedRenderedTile
}

func Test_gets_attribute_table_address_from_tile(t *testing.T) {
	var vram [2 * NAMETABLE_SIZE]byte

	// Setup attribute table
	vram[0x03C0] = 0b11011000
	// 0b00 background palette 0
	// 0b01 background palette 1
	// 0b10 background palette 2
	// 0b11 background palette 3
	// Tiles and its palettes --
	//	[00][00] [10][10]
	//  [00][00] [10][10]
	//  [01][01] [11][11]
	//  [01][01] [11][11]

	// First meta tile should
	for tileColumn := uint8(0); tileColumn < 2; tileColumn++ {
		for tileRow := uint8(0); tileRow < 2; tileRow++ {
			palette := backgroundPalette(tileColumn, tileRow, &vram)
			assert.Equal(t, byte(0b00), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b00", tileColumn, tileRow))
		}
	}

	// Second meta tile should
	for x := uint8(2); x < 4; x++ {
		for y := uint8(0); y < 2; y++ {
			palette := backgroundPalette(x, y, &vram)
			assert.Equal(t, byte(0b10), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}

	// Third meta tile should
	for x := uint8(0); x < 2; x++ {
		for y := uint8(2); y < 4; y++ {
			palette := backgroundPalette(x, y, &vram)
			assert.Equal(t, byte(0b01), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}

	// Fourth meta tile should
	for x := uint8(2); x < 4; x++ {
		for y := uint8(2); y < 4; y++ {
			palette := backgroundPalette(x, y, &vram)
			assert.Equal(t, byte(0b11), palette, fmt.Sprintf("Second meta tile (%d,%d) should use palette 0b10", x, y))
		}
	}
}

func TestPPU_should_get_propper_color_for_a_given_pixel_color_and_palette(t *testing.T) {
	ppu := aPPU()
	backgroundColor := utils.NewColorRGB(236, 88, 180)
	cases := []struct {
		name          string
		palette       byte
		colorIndex    byte
		expectedColor color.Color
	}{
		{"", 0, 0, backgroundColor},
		{"", 0, 1, utils.NewColorRGB(0, 30, 116)},
		{"", 0, 2, utils.NewColorRGB(8, 16, 144)},
		{"", 0, 3, utils.NewColorRGB(48, 0, 136)},
		//{"mirroring $0x3F10", 4, 0, backgroundColor},
		//{"mirroring $0x3F14", 5, 0, graphics.Color{68, 0, 100}},
		//{"mirroring $0x3F18", 6, 0, graphics.Color{32, 42, 0}},
		//{"mirroring $0x3F1C", 7, 0, graphics.Color{0, 50, 60}},
	}

	// 0x3F00
	ppu.Write(PaletteLowAddress+0, 0x25) // light blue
	ppu.Write(PaletteLowAddress+1, 0x01) // Dark Blue
	ppu.Write(PaletteLowAddress+2, 0x02) // Blue-Purple
	ppu.Write(PaletteLowAddress+3, 0x03) // Dark Purple

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			paletteColor := ppu.GetRGBColor(tt.palette, tt.colorIndex)
			assert.Equal(t, tt.expectedColor, paletteColor)
		})
	}
}

func TestPpu2c02_GetNesColorFromPaletteRam_should_gets_background_color(t *testing.T) {
	ppu := aPPU()
	backgroundColor := byte(0x10)
	ppu.Write(PaletteLowAddress, backgroundColor)

	color := ppu.GetPaletteColor(0, 0)
	assert.Equal(t, backgroundColor, color)

	color = ppu.GetPaletteColor(1, 0)
	assert.Equal(t, backgroundColor, color)

	color = ppu.GetPaletteColor(2, 0)
	assert.Equal(t, backgroundColor, color)

	color = ppu.GetPaletteColor(3, 0)
	assert.Equal(t, backgroundColor, color)
}
