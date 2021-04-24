package nes

import "fmt"

type Mapper interface {
	prgBanks() byte
	chrBanks() byte

	Read(address Address) byte
	Write(address Address, value byte)
}

func CreateMapper(gamePak *GamePak) Mapper {
	header := gamePak.Header()
	switch header.MapperNumber() {
	case 0:
		return CreateMapper000(gamePak)
	}

	panic(fmt.Sprintf("mapper %d not implemented", header.MapperNumber()))
}
