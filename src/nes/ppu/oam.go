package ppu

type objectAttributeEntry struct {
	y      byte
	tileId byte

	// 76543210
	// ||||||||
	// ||||||++- Palette (4 to 7) of sprite
	// |||+++--- Unimplemented (read 0)
	// ||+------ Priority (0: in front of background; 1: behind background)
	// |+------- Flip sprite horizontally
	// +-------- Flip sprite vertically
	attributes byte
	x          byte
}

func (oae *objectAttributeEntry) isFlippedVertically() bool {
	return oae.attributes&0x80 == 0x80
}

func (oae *objectAttributeEntry) isFlippedHorizontally() bool {
	return oae.attributes&0x40 == 0x40
}

func (oae *objectAttributeEntry) palette() byte {
	// Sprite palette always goes from index 4 up to 7.
	// but to save bytes, only three bits are used.
	return (oae.attributes & 0b111) + 4
}
