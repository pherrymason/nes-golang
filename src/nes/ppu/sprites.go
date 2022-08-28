package ppu

import (
	"github.com/raulferras/nes-golang/src/nes/types"
	"math/bits"
)

func (ppu *P2c02) loadNextScanLineSprites() {
	ppu.spriteScanlineCount = 0
	var spriteHeight int16
	if ppu.PpuControl.SpriteSize == 0 {
		spriteHeight = 8
	} else {
		spriteHeight = 16
	}

	for oamIndex := 0; oamIndex < OAMDATA_SIZE && ppu.spriteScanlineCount < 8; oamIndex += 4 {
		spriteY := int16(ppu.oamData[oamIndex])
		diff := int16(ppu.currentScanline+1) - spriteY
		if diff >= 0 && diff < spriteHeight {
			ppu.oamDataScanline[ppu.spriteScanlineCount] = objectAttributeEntry{
				y:          byte(spriteY),
				tileId:     ppu.oamData[oamIndex+1],
				attributes: ppu.oamData[oamIndex+2],
				x:          ppu.oamData[oamIndex+3],
			}
			ppu.spriteScanlineCount++
		}
	}
}

// Get sprite pattern information
// Doing this at cycle 340 is a simplification of what the real NES actually does.
// More info: https://www.nesdev.org/wiki/PPU_sprite_evaluation
func (ppu *P2c02) loadSpriteShifters() {
	var spritePatternAddressLow types.Address
	var spritePatternAddressHigh types.Address
	var spritePatternLow byte
	var spritePatternHigh byte

	for i := byte(0); i < ppu.spriteScanlineCount; i++ {
		object := ppu.oamDataScanline[i]

		if ppu.PpuControl.SpriteSize == PPU_CONTROL_SPRITE_SIZE_8 {
			if !object.isFlippedVertically() {
				spritePatternAddressLow = types.Address(ppu.PpuControl.SpritePatternTableAddress) << 12
				spritePatternAddressLow |= types.Address(object.tileId) << 4                             // Multiply ID per 16 (16 bytes per tile)
				spritePatternAddressLow |= types.Address((ppu.currentScanline + 1) - Scanline(object.y)) // Which tile line we want
			} else {
				spritePatternAddressLow = types.Address(ppu.PpuControl.SpritePatternTableAddress) << 12
				spritePatternAddressLow |= types.Address(object.tileId) << 4                                 // Multiply ID per 16 (16 bytes per tile)
				spritePatternAddressLow |= types.Address(7 - (ppu.currentScanline - Scanline(object.y) - 1)) // Which tile line we want
			}
		} else {
			// TODO implement 8x16 sprites
		}

		spritePatternAddressHigh = spritePatternAddressLow + 8
		spritePatternLow = ppu.Read(spritePatternAddressLow)
		spritePatternHigh = ppu.Read(spritePatternAddressHigh)

		if object.isFlippedHorizontally() {
			spritePatternLow = bits.Reverse8(spritePatternLow)
			spritePatternHigh = bits.Reverse8(spritePatternHigh)
		}

		ppu.spShifterPatternLow[i] = spritePatternLow
		ppu.spShifterPatternHigh[i] = spritePatternHigh
	}
}

func (ppu *P2c02) checkSprite0Hit() {
	if !ppu.PpuMask.showSpritesEnabled() {
		return
	}

	if Scanline(ppu.oamDataScanline[0].y) == ppu.currentScanline {
		if uint16(ppu.oamDataScanline[0].x) <= ppu.renderCycle {
			ppu.PpuStatus.Sprite0Hit = 1
		}
	}
}
