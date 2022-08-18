package ppu

const PPU_CONTROL_SPRITE_SIZE_8 = 0
const PPU_CONTROL_SPRITE_SIZE_16 = 1

type Control struct {
	NameTableX                    byte // Most significant bit of scrolling coordinates (X)
	NameTableY                    byte // Most significant bit of scrolling coordinates (Y)
	IncrementMode                 byte // Address increment per CPU IO of PPUDATA. (0: add 1, going across; 1: add 32, going down)
	SpritePatternTableAddress     byte // Sprite pattern table address for 8x8 sprites. (0: $0000; 1: $1000; ignored in 8x16 mode)
	BackgroundPatternTableAddress byte // Background pattern table address (0: $0000; 1: $1000)
	SpriteSize                    byte // 0: 8x8   1:8x16
	MasterSlaveSelect             byte
	GenerateNMIAtVBlank           bool // NMI enabled/disabled
}

func (control *Control) Value() byte {
	ctrl := byte(0)
	ctrl |= control.NameTableX
	ctrl |= control.NameTableY << 1
	ctrl |= control.IncrementMode << 2
	ctrl |= control.SpritePatternTableAddress << 3
	ctrl |= control.BackgroundPatternTableAddress << 4
	ctrl |= control.SpriteSize << 5
	ctrl |= control.MasterSlaveSelect << 6
	if control.GenerateNMIAtVBlank {
		ctrl |= 1 << 7
	}

	return ctrl
}

type Status struct {
	SpriteOverflow       byte
	Sprite0Hit           byte
	VerticalBlankStarted bool
}

func (status *Status) Value() byte {
	value := byte(0)
	value |= status.SpriteOverflow << 5
	value |= status.Sprite0Hit << 6
	if status.VerticalBlankStarted {
		value |= 1 << 7
	}

	return value
}

func (ppu *P2c02) ppuCtrlWrite(value byte) {
	ppu.PpuControl.NameTableX = value & 0x01
	ppu.PpuControl.NameTableY = (value >> 1) & 1
	ppu.PpuControl.IncrementMode = (value >> 2) & 1
	ppu.PpuControl.SpritePatternTableAddress = (value >> 3) & 1
	ppu.PpuControl.BackgroundPatternTableAddress = (value >> 4) & 1
	ppu.PpuControl.SpriteSize = (value >> 5) & 1
	ppu.PpuControl.MasterSlaveSelect = (value >> 6) & 1
	if (value>>7)&1 == 1 {
		ppu.PpuControl.GenerateNMIAtVBlank = true
	} else {
		ppu.PpuControl.GenerateNMIAtVBlank = false
	}
}

func (ppu *P2c02) VBlank() bool {
	return ppu.PpuStatus.VerticalBlankStarted
}
