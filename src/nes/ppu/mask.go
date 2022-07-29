package ppu

type Mask struct {
	greyScale              byte
	showBackgroundLeftMost byte
	showSpritesLeftMost    byte
	showBackground         byte
	showSprites            byte
	emphasizeRed           byte
	emphasizeGreen         byte
	emphasizeBlue          byte
}

func (register *Mask) write(value byte) {
	register.greyScale = 0
	register.showBackgroundLeftMost = 0
	register.showSpritesLeftMost = 0
	register.showBackground = 0
	register.showSprites = 0
	register.emphasizeRed = 0
	register.emphasizeGreen = 0
	register.emphasizeBlue = 0

	if value&0x01 == 1 {
		register.greyScale = 1
	}
	if (value>>1)&0x01 == 1 {
		register.showBackgroundLeftMost = 1
	}
	if (value>>2)&0x01 == 1 {
		register.showSpritesLeftMost = 1
	}
	if (value>>3)&0x01 == 1 {
		register.showBackground = 1
	}
	if (value>>4)&0x01 == 1 {
		register.showSprites = 1
	}
	if (value>>5)&0x01 == 1 {
		register.emphasizeRed = 1
	}
	if (value>>6)&0x01 == 1 {
		register.emphasizeGreen = 1
	}
	if (value>>7)&0x01 == 1 {
		register.emphasizeGreen = 1
	}
}

func (register *Mask) value() byte {
	value := byte(0)
	value |= register.greyScale
	value |= register.showBackgroundLeftMost << 1
	value |= register.showSpritesLeftMost << 2
	value |= register.showBackground << 3
	value |= register.showSprites << 4
	value |= register.emphasizeRed << 5
	value |= register.emphasizeGreen << 6
	value |= register.emphasizeBlue << 7

	return value
}

func (register *Mask) showBackgroundEnabled() bool {
	return register.showBackground == 1
}

func (register *Mask) showSpritesEnabled() bool {
	return register.showSprites == 1
}

func (register *Mask) renderingEnabled() bool {
	return register.showBackgroundEnabled() || register.showSpritesEnabled()
}
