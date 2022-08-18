package gamePak

import "github.com/raulferras/nes-golang/src/nes/types"

// if PRGROM is 16KB
//     CPU Address Bus          PRG ROM
//     0x8000 -> 0xBFFF: Map    0x0000 -> 0x3FFF
//     0xC000 -> 0xFFFF: Mirror 0x0000 -> 0x3FFF
// if PRGROM is 32KB
//     CPU Address Bus          PRG ROM
//     0x8000 -> 0xFFFF: Map    0x0000 -> 0x7FFF

type Mapper000 struct {
	prgROMBanks byte
	chrROMBanks byte
	prgROM      []byte
	chrROM      []byte
	hasCHRRAM   bool
}

func CreateMapper000(header Header, prgROM []byte, chrROM []byte) *Mapper000 {
	mapper0 := Mapper000{
		//gamePak:     gamePak,
		prgROMBanks: header.ProgramSize(),
		chrROMBanks: header.CHRSize(),
		prgROM:      prgROM,
		chrROM:      chrROM,
		hasCHRRAM:   header.CHRSize() == 0,
	}

	if header.CHRSize() == 0 {
		chrLength := 8 * 0x2000
		mapper0.chrROM = make([]byte, chrLength)
	}

	return &mapper0
}

func (mapper *Mapper000) PrgBanks() byte {
	return mapper.prgROMBanks
}

func (mapper *Mapper000) ChrBanks() byte {
	return mapper.chrROMBanks
}

func (mapper *Mapper000) ReadPrgROM(address types.Address) byte {
	if !satisfiableAddress(address) {
		return 0
	}

	if mapper.PrgBanks() == 1 {
		address = address & 0x3FFF
	} else {
		address = address & 0x7FFF
	}

	return mapper.prgROM[address]
}

func (mapper *Mapper000) WritePrgROM(address types.Address, value byte) {
	if !satisfiableAddress(address) {
		return
	}

	mapper.prgROM[address] = value
}

func (mapper *Mapper000) ReadChrROM(address types.Address) byte {
	return mapper.chrROM[address]
}

func (mapper *Mapper000) WriteChrROM(address types.Address, value byte) {
	mapper.chrROM[address] = value
}

func satisfiableAddress(address types.Address) bool {
	if address >= 0x8000 && address <= 0xFFFF {
		return true
	}

	return false
}
