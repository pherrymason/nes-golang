package gamePak

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type Mapper interface {
	PrgBanks() byte
	ChrBanks() byte

	ReadPrgROM(address types.Address) byte
	WritePrgROM(address types.Address, value byte)
	ReadChrROM(address types.Address) byte
	WriteChrROM(address types.Address, value byte)
}

func CreateMapper(header Header, prgROM []byte, chrROM []byte) Mapper {
	switch header.MapperNumber() {
	case 0:
		return CreateMapper000(header, prgROM, chrROM)
	}

	panic(fmt.Sprintf("mapper %d not supported", header.MapperNumber()))
}
