package nes

import (
	"fmt"
	"io/ioutil"
)

const GAMEPAK_MEMORY_SIZE = 0xBFE0
const GAMEPAK_LOW_RANGE = 0x4020
const GAMEPAK_HIGH_RANGE = 0xFFFF

const GAMEPAK_ROM_LOWER_BANK_START = 0x8000

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

func CreateGamePakFromROMFile(romFilePath string) GamePak {
	data, err := ioutil.ReadFile(romFilePath)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read Header
	inesHeader := CreateINes1Header(data[0:16])

	return CreateGamePak(inesHeader, data[16:])
}

func (gamePak GamePak) Header() Header {
	return gamePak.header
}

func (gamePak *GamePak) read(address Address) byte {
	romAddress := toRomAddress(address)
	return gamePak.prgROM[romAddress]
}

func (gamePak *GamePak) write(address Address, value byte) {
	romAddress := toRomAddress(address)
	gamePak.prgROM[romAddress] = value
}

func toRomAddress(address Address) Address {
	offset := Address(GAMEPAK_ROM_LOWER_BANK_START)
	//if gamepak.header.HasTrainer() {
	//	offset += 512
	//}

	// NROM has mirroring from 0xC000
	if address >= 0xC000 {
		address -= 0x4000
	}

	romAddress := address - offset

	//if int(romAddress) > len(gamepak.prgROM) {
	//	return 0
	//}

	return romAddress
}
