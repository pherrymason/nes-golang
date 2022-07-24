package gamePak

import (
	"fmt"
	"io/ioutil"
)

func CreateGamePak(header Header, prgROM []byte, chrROM []byte) GamePak {
	return GamePak{header, prgROM, chrROM}
}

func NewGamePakWithINes(flag6 byte, flag7 byte, flag8 byte, flag9 byte, flag10 byte, prgROM []byte, chrROM []byte) GamePak {

	return GamePak{
		header: CreateINes1Header(byte(len(prgROM)/16), byte(len(chrROM)/8), flag6, flag7, flag8, flag9, flag10),
		prgROM: prgROM,
		chrROM: chrROM,
	}
}

func CreateGamePakFromROMFile(romFilePath string) GamePak {
	data, err := ioutil.ReadFile(romFilePath)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	// Read INesHeader
	inesHeader := CreateINes1Header(
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

func NewDummyGamePak(chrROM []byte) *GamePak {
	pak := CreateGamePak(
		CreateINes1Header(1, 1, 0, 0, 0, 0, 0),
		make([]byte, 100),
		chrROM,
	)

	return &pak
}

func NewEmptyCHRROM() []byte {
	return make([]byte, 0x01FFF)
}
