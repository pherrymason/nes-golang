package component

import "github.com/raulferras/nes-golang/src/nes/defs"

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

	offset := defs.Address(0x8000)
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

func (gamepak *GamePak) write(address defs.Address, value byte) {
	offset := defs.Address(0x8000)
	if gamepak.header.HasTrainer() {
		offset += 512
	}
	gamepak.prgROM[address-offset] = value
}
