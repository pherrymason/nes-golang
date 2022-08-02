package ppu

const PPU_CONTROL_SPRITE_SIZE_8 = 0
const PPU_CONTROL_SPRITE_SIZE_16 = 1

type Control struct {
	nameTableX                    byte // Most significant bit of scrolling coordinates (X)
	nameTableY                    byte // Most significant bit of scrolling coordinates (Y)
	incrementMode                 byte // Address increment per CPU IO of PPUDATA. (0: add 1, going across; 1: add 32, going down)
	spritePatternTableAddress     byte // Sprite pattern table address for 8x8 sprites. (0: $0000; 1: $1000; ignored in 8x16 mode)
	backgroundPatternTableAddress byte // Background pattern table address (0: $0000; 1: $1000)
	spriteSize                    byte // 0: 8x8   1:8x16
	masterSlaveSelect             byte
	generateNMIAtVBlank           bool // NMI enabled/disabled
}

func (control *Control) value() byte {
	ctrl := byte(0)
	ctrl |= control.nameTableX
	ctrl |= control.nameTableY << 1
	ctrl |= control.incrementMode << 2
	ctrl |= control.spritePatternTableAddress << 3
	ctrl |= control.backgroundPatternTableAddress << 4
	ctrl |= control.spriteSize << 5
	ctrl |= control.masterSlaveSelect << 6
	if control.generateNMIAtVBlank {
		ctrl |= 1 << 7
	}

	return ctrl
}

type MASKFlag int

const (
	greyScale MASKFlag = iota
	showBackgroundLeftEdge
	showSpritesLeftEdge
	showBackground
	showSprites
	emphasizeRed
	emphasizeGreen
	emphasizeBlue
)

type Status struct {
	spriteOverflow       byte
	sprite0Hit           byte
	verticalBlankStarted bool
}

func (status *Status) value() byte {
	value := byte(0)
	value |= status.spriteOverflow << 5
	value |= status.sprite0Hit << 6
	if status.verticalBlankStarted {
		value |= 1 << 7
	}

	return value
}

func (ppu *Ppu2c02) ppuCtrlWrite(value byte) {
	ppu.ppuControl.nameTableX = value & 0x01
	ppu.ppuControl.nameTableY = (value >> 1) & 1
	ppu.ppuControl.incrementMode = (value >> 2) & 1
	ppu.ppuControl.spritePatternTableAddress = (value >> 3) & 1
	ppu.ppuControl.backgroundPatternTableAddress = (value >> 4) & 1
	ppu.ppuControl.spriteSize = (value >> 5) & 1
	ppu.ppuControl.masterSlaveSelect = (value >> 6) & 1
	if (value>>7)&1 == 1 {
		ppu.ppuControl.generateNMIAtVBlank = true
	} else {
		ppu.ppuControl.generateNMIAtVBlank = false
	}
}

func (ppu *Ppu2c02) VBlank() bool {
	return ppu.ppuStatus.verticalBlankStarted
}
