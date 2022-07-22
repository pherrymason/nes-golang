package gamePak

import (
	"github.com/raulferras/nes-golang/src/nes/types"
)

const GAMEPAK_MEMORY_SIZE = 0xBFE0
const GAMEPAK_LOW_RANGE = 0x4020
const GAMEPAK_HIGH_RANGE = 0xFFFF

const GAMEPAK_ROM_LOWER_BANK_START = 0x8000

type GamePakInterface interface {
	Header() Header
	ReadPrgROM(address types.Address) byte
	WritePrgROM(address types.Address, value byte)
	ReadCHRROM(address types.Address) byte
}

type GamePak struct {
	header Header
	prgROM []byte
	chrROM []byte
}

func (gamePak GamePak) Header() Header {
	return gamePak.header
}

func (gamePak GamePak) ReadPrgROM(address types.Address) byte {
	return gamePak.prgROM[address]
}

func (gamePak *GamePak) WritePrgROM(address types.Address, value byte) {
	gamePak.prgROM[address] = value
}

func (gamePak *GamePak) ReadCHRROM(address types.Address) byte {
	return gamePak.chrROM[address]
}
