package component

import "github.com/raulferras/nes-golang/src/nes/defs"

const GAMEPAK_MEMORY_SIZE = 0xBFE0
const GAMEPAK_LOW_RANGE = 0x4020
const GAMEPAK_HIGH_RANGE = 0xFFFF

const GAMEPAK_ROM_START = 0x8000

type GamePak struct {
	header Header
	prgROM []byte
}

func CreateDummyGamePak() GamePak {
	return GamePak{
		header: Header{},
		prgROM: make([]byte, 0xFFFF),
	}
}

func CreateGamePak(header Header, prgROM []byte) GamePak {
	return GamePak{header, prgROM}
}

func (gamepak *GamePak) read(address defs.Address) byte {

	offset := defs.Address(GAMEPAK_ROM_START)
	//if gamepak.header.HasTrainer() {
	//	offset += 512
	//}

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

func (gamepak *GamePak) write(address defs.Address, value byte) {
	offset := defs.Address(GAMEPAK_ROM_START)
	//if gamepak.header.HasTrainer() {
	//	offset += 512
	//}

	// NROM has mirroring from 0xC000
	if address >= 0xC000 {
		address -= 0x4000
	}

	gamepak.prgROM[address-offset] = value
}
