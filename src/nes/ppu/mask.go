package ppu

type Mask struct {
	GreyScale              byte
	ShowBackgroundLeftMost byte
	ShowSpritesLeftMost    byte
	ShowBackground         byte // x8
	ShowSprites            byte
	EmphasizeRed           byte
	EmphasizeGreen         byte
	EmphasizeBlue          byte
}

func (register *Mask) write(value byte) {
	register.GreyScale = 0
	register.ShowBackgroundLeftMost = 0
	register.ShowSpritesLeftMost = 0
	register.ShowBackground = 0
	register.ShowSprites = 0
	register.EmphasizeRed = 0
	register.EmphasizeGreen = 0
	register.EmphasizeBlue = 0

	if value&0x01 == 1 {
		register.GreyScale = 1
	}
	if (value>>1)&0x01 == 1 {
		register.ShowBackgroundLeftMost = 1
	}
	if (value>>2)&0x01 == 1 {
		register.ShowSpritesLeftMost = 1
	}
	if (value>>3)&0x01 == 1 {
		register.ShowBackground = 1
	}
	if (value>>4)&0x01 == 1 {
		register.ShowSprites = 1
	}
	if (value>>5)&0x01 == 1 {
		register.EmphasizeRed = 1
	}
	if (value>>6)&0x01 == 1 {
		register.EmphasizeGreen = 1
	}
	if (value>>7)&0x01 == 1 {
		register.EmphasizeGreen = 1
	}
}

func (register *Mask) Value() byte {
	value := byte(0)
	value |= register.GreyScale
	value |= register.ShowBackgroundLeftMost << 1
	value |= register.ShowSpritesLeftMost << 2
	value |= register.ShowBackground << 3
	value |= register.ShowSprites << 4
	value |= register.EmphasizeRed << 5
	value |= register.EmphasizeGreen << 6
	value |= register.EmphasizeBlue << 7

	return value
}

func (register *Mask) showBackgroundEnabled() bool {
	return register.ShowBackground == 1
}

func (register *Mask) showSpritesEnabled() bool {
	return register.ShowSprites == 1
}

func (register *Mask) renderingEnabled() bool {
	return register.showBackgroundEnabled() || register.showSpritesEnabled()
}
