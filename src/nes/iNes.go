package nes

type GamePak struct {
	header Header
	prgROM []byte
}

func (gamepak *GamePak) read(address Address) byte {

	offset := Address(0x8000)
	if gamepak.header.HasTrainer() {
		offset += 512
	}

	// NROM has mirroring from 0xC000
	if address >= 0xC000 {
		address -= 0x4000
	}

	finalAddress := address - offset

	if int(finalAddress) > len(gamepak.prgROM) {
		return 0
	}

	return gamepak.prgROM[finalAddress]
}

func (gamepak *GamePak) write(address Address, value byte) {
	offset := Address(0x8000)
	if gamepak.header.HasTrainer() {
		offset += 512
	}
	gamepak.prgROM[address-offset] = value
}

type iNes interface {
	ProgramSize() byte
	CHRSize() byte
	Mirroring() byte
	HasPersistentMemory() bool
	HasTrainer() bool
	IgnoreMirroringControl() bool

	MapperNumber() byte
	PRGRAM() byte
	TvSystem() byte
}

type Header struct {
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
func (ines *Header) ProgramSize() byte {
	return ines.prgROMSize
}

func (ines *Header) CHRSize() byte {
	return ines.chrROMSize
}
func (ines *Header) HasTrainer() bool {
	if ines.flags6&0x06 == 0x06 {
		return true
	}

	return false
}

func (ines *Header) MapperNumber() byte {
	return (ines.flags6 >> 4) | (ines.flags7 & 0xF0)

}

func (ines *Header) TvSystem() byte {
	return (ines.flags9 & 0x01)
}

func CreateINes1Header(header []byte) Header {
	return Header{
		prgROMSize: header[4],
		chrROMSize: header[5],
		flags6:     header[6],
		flags7:     header[7],
		flags8:     header[8],
		flags9:     header[9],
		flags10:    header[10],
	}
}
