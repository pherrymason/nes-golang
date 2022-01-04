package mocks

import "github.com/raulferras/nes-golang/src/nes/types"

type SimpleMapper struct {
	memory [0x10000]byte
}

func (s SimpleMapper) PrgBanks() byte {
	return 1
}

func (s SimpleMapper) ChrBanks() byte {
	return 0
}

func (s SimpleMapper) Read(address types.Address) byte {
	return s.memory[address]
}

func (s SimpleMapper) Write(address types.Address, value byte) {
	s.memory[address] = value
}
