package gamePak

type INesHeader struct {
	prgROMSize byte
	chrROMSize byte
	flags6     byte
	flags7     byte
	flags8     byte
	flags9     byte
	flags10    byte
}

/*
Flags6
76543210
||||||||
|||||||+- Mirroring: 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
|||||||              1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
||||||+-- 1: GamePak contains battery-backed PRG RAM ($6000-7FFF) or other persistent memory
|||||+--- 1: 512-byte trainer at $7000-$71FF (stored before PRG prgROM)
||||+---- 1: Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
++++----- Lower nybble of mapper number
*/
func (ines INesHeader) ProgramSize() byte {
	return ines.prgROMSize
}

func (ines INesHeader) CHRSize() byte {
	return ines.chrROMSize
}

func (ines INesHeader) Mirroring() byte {
	mirroring := byte(ines.flags6 & 0b11)

	return mirroring
}

func (ines INesHeader) HasTrainer() bool {
	if ines.flags6&0x06 == 0x06 {
		return true
	}

	return false
}

func (ines INesHeader) HasPersistentMemory() bool {
	if ines.flags6&0x04 == 0x04 {
		return true
	}

	return false
}

func (ines INesHeader) IgnoreMirroringControl() bool {
	panic("implement me")
}

func (ines INesHeader) PRGRAM() byte {
	panic("implement me")
}

func (ines INesHeader) MapperNumber() byte {
	return (ines.flags6 >> 4) | (ines.flags7 & 0xF0)
}

func (ines INesHeader) TvSystem() byte {
	return ines.flags9 & 0x01
}

func CreateINes1Header(prgRomSize byte, chrRomSize byte, flag6 byte, flag7 byte, flag8 byte, flag9 byte, flag10 byte) INesHeader {
	return INesHeader{
		prgROMSize: prgRomSize,
		chrROMSize: chrRomSize,
		flags6:     flag6,
		flags7:     flag7,
		flags8:     flag8,
		flags9:     flag9,
		flags10:    flag10,
	}
}
