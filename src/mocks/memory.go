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

func (s *SimpleMemory) IsDMAWaiting() bool {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) IsDMATransfer() bool {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) DisableDMWaiting() {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) GetDMAPage() byte {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) GetDMAAddress() byte {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) GetDMAReadBuffer() byte {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) SetDMAReadBuffer(value byte) {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) IncrementDMAAddress() {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleMemory) ResetDMA() {
	//TODO implement me
	panic("implement me")
}

func NewSimpleMemory() *SimpleMemory {
	return &SimpleMemory{}
}
