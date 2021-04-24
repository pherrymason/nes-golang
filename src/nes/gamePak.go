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
	chrROM []byte
}

func CreateGamePak(header Header, prgROM []byte, chrROM []byte) GamePak {
	return GamePak{header, prgROM, chrROM}
}

func CreateGamePakFromROMFile(romFilePath string) GamePak {
	data, err := ioutil.ReadFile(romFilePath)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read Header
	inesHeader := CreateINes1Header(data[0:16])

	prgLength := int(inesHeader.prgROMSize)*0x4000 + 16
	prgROM := data[16:prgLength]

	chrLength := int(inesHeader.chrROMSize) * 0x2000
	chrROM := data[prgLength : chrLength+prgLength]
	return CreateGamePak(
		inesHeader,
		prgROM,
		chrROM,
	)
}

func (gamePak GamePak) Header() Header {
	return gamePak.header
}
