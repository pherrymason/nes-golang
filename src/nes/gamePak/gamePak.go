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
	WriteCHRAM(address types.Address, value byte)
}

type GamePak struct {
	header Header
	mapper Mapper
	prgROM []byte
	chrROM []byte
}

func (gamePak *GamePak) Header() Header {
	return gamePak.header
}

func (gamePak *GamePak) ReadPrgROM(address types.Address) byte {
	return gamePak.mapper.ReadPrgROM(address)
}

func (gamePak *GamePak) WritePrgROM(address types.Address, value byte) {
	gamePak.mapper.WritePrgROM(address, value)
}

func (gamePak *GamePak) ReadCHRROM(address types.Address) byte {
	return gamePak.mapper.ReadChrROM(address)
}

func (gamePak *GamePak) WriteCHRRAM(address types.Address, value byte) {
	gamePak.mapper.WriteChrROM(address, value)
}
