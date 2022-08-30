package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createCHRRom() []byte {
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

	return chrROM
}

func Test_loadSpritesShifters_should_render_sprite_line_no_flipping(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")

	// setup oamData Scanline
	ppu.currentScanline = 0
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 0, tileId: 0, attributes: 0x00, x: 0}
	ppu.spriteScanlineCount = 1

	ppu.loadSpriteShifters()

	assert.Equal(t, byte(0b00000111), ppu.spShifterPatternLow[0])
	assert.Equal(t, byte(0b00000111), ppu.spShifterPatternHigh[0])
}

func Test_loadSpritesShifters_should_render_sprite_line_vertical_flipping(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")

	// setup oamData Scanline
	ppu.currentScanline = 1
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 0, tileId: 0, attributes: 0x80, x: 0}
	ppu.spriteScanlineCount = 1

	ppu.loadSpriteShifters()

	assert.Equal(t, byte(0b00000001), ppu.spShifterPatternLow[0])
	assert.Equal(t, byte(0b11111111), ppu.spShifterPatternHigh[0])
}

func Test_loadSpritesShifters_should_render_sprite_line_horizontal_flipping(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")

	// setup oamData Scanline
	ppu.currentScanline = 0
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 0, tileId: 0, attributes: 0x40, x: 0}
	ppu.spriteScanlineCount = 1

	ppu.loadSpriteShifters()

	assert.Equal(t, byte(0b11100000), ppu.spShifterPatternLow[0])
	assert.Equal(t, byte(0b11100000), ppu.spShifterPatternHigh[0])
}

func Test_loadSpritesShifters_should_render_sprite_line_vertical_and_horizontal_flipping(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")

	// setup oamData Scanline
	ppu.currentScanline = 1
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 0, tileId: 0, attributes: 0x80 | 0x40, x: 0}
	ppu.spriteScanlineCount = 1

	ppu.loadSpriteShifters()

	assert.Equal(t, byte(0b10000000), ppu.spShifterPatternLow[0])
	assert.Equal(t, byte(0b11111111), ppu.spShifterPatternHigh[0])
}

func Test_checkSprite0Hit_should_not_enable_flag_if_sprite_rendering_is_disabled(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")
	ppu.currentScanline = 10
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 10, tileId: 0, attributes: 0x80 | 0x40, x: 20}

	ppu.checkSprite0Hit()

	assert.Equal(t, byte(0), ppu.PpuStatus.Sprite0Hit)
}

func Test_checkSprite0Hit_should_not_enable_flag_if_no_sprite_0_is_being_rendered(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")
	ppu.currentScanline = 10
	ppu.renderCycle = 20
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 11, tileId: 0, attributes: 0x80 | 0x40, x: 20}

	ppu.checkSprite0Hit()

	assert.Equal(t, byte(0), ppu.PpuStatus.Sprite0Hit)
}

func Test_checkSprite0Hit_should_enable_flag_if_sprite_0_hits_opaque_background_pixel(t *testing.T) {
	chrRom := createCHRRom()
	cartridge := gamePak.NewDummyGamePak(chrRom)
	ppu := CreatePPU(cartridge, false, "")
	ppu.currentScanline = 10
	ppu.renderCycle = 20
	ppu.oamDataScanline[0] = objectAttributeEntry{y: 10, tileId: 0, attributes: 0x80 | 0x40, x: 20}

	ppu.checkSprite0Hit()

	assert.Equal(t, byte(0), ppu.PpuStatus.Sprite0Hit)
}
