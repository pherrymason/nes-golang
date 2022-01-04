package mappers

import (
	gamePak2 "github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/types"
)

// if PRGROM is 16KB
//     CPU Address Bus          PRG ROM
//     0x8000 -> 0xBFFF: Map    0x0000 -> 0x3FFF
//     0xC000 -> 0xFFFF: Mirror 0x0000 -> 0x3FFF
// if PRGROM is 32KB
//     CPU Address Bus          PRG ROM
//     0x8000 -> 0xFFFF: Map    0x0000 -> 0x7FFF

type Mapper000 struct {
	gamePak     *gamePak2.GamePak
	prgROMBanks byte
	chrROMBanks byte
}

func CreateMapper000(gamePak *gamePak2.GamePak) Mapper000 {
	header := gamePak.Header()
	return Mapper000{
		gamePak:     gamePak,
		prgROMBanks: header.ProgramSize(),
		chrROMBanks: header.CHRSize(),
	}
}

func (mapper Mapper000) PrgBanks() byte {
	return mapper.prgROMBanks
}

func (mapper Mapper000) ChrBanks() byte {
	return mapper.chrROMBanks
}

func (mapper Mapper000) Read(address types.Address) byte {
	if !satisfiableAddress(address) {
		return 0
	}

	if mapper.PrgBanks() == 1 {
		address = address & 0x3FFF
	} else {
		address = address & 0x7FFF
	}

	return mapper.gamePak.ReadPrgROM(address)
}

func (mapper Mapper000) Write(address types.Address, value byte) {
	if !satisfiableAddress(address) {
		return
	}

	mapper.gamePak.WritePrgROM(address, value)
}

func satisfiableAddress(address types.Address) bool {
	if address >= 0x8000 && address <= 0xFFFF {
		return true
	}

	return false
}
