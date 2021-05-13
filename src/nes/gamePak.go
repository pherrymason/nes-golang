package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"io/ioutil"
)

const GAMEPAK_MEMORY_SIZE = 0xBFE0
const GAMEPAK_LOW_RANGE = 0x4020
const GAMEPAK_HIGH_RANGE = 0xFFFF

const GAMEPAK_ROM_LOWER_BANK_START = 0x8000

type GamePak struct {
	header gamePak.Header
	prgROM []byte
	chrROM []byte
}

func CreateGamePak(header gamePak.Header, prgROM []byte, chrROM []byte) GamePak {
	return GamePak{header, prgROM, chrROM}
}

func CreateGamePakFromROMFile(romFilePath string) GamePak {
	data, err := ioutil.ReadFile(romFilePath)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read INesHeader
	inesHeader := gamePak.CreateINes1Header(
		data[4],
		data[5],
		data[6],
		data[7],
		data[8],
		data[9],
		data[10],
	)

	prgLength := int(inesHeader.ProgramSize())*0x4000 + 16
	prgROM := data[16:prgLength]

	chrLength := int(inesHeader.CHRSize()) * 0x2000
	chrROM := data[prgLength : chrLength+prgLength]
	return CreateGamePak(
		inesHeader,
		prgROM,
		chrROM,
	)
}

func (gamePak GamePak) Header() gamePak.Header {
	return gamePak.header
}
