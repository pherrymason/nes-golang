package nes

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/types"
)

type Mapper interface {
	prgBanks() byte
	chrBanks() byte

	Read(address types.Address) byte
	Write(address types.Address, value byte)
}

func CreateMapper(gamePak *GamePak) Mapper {
	header := gamePak.Header()
	switch header.MapperNumber() {
	case 0:
		return CreateMapper000(gamePak)
	}

	panic(fmt.Sprintf("mapper %d not implemented", header.MapperNumber()))
}
