package mocks

import "github.com/raulferras/nes-golang/src/nes/types"

type SimpleMemory struct {
	ram [0xFFFF + 1]byte
}

func (s *SimpleMemory) Peek(address types.Address) byte {
	return s.ram[address]
}

func (s *SimpleMemory) Read(address types.Address) byte {
	return s.ram[address]
}

func (s *SimpleMemory) Write(address types.Address, b byte) {
	s.ram[address] = b
}

func NewSimpleMemory() *SimpleMemory {
	return &SimpleMemory{}
}
