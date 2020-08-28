package nes

// Address is
type Address uint16

func (address *Address) low() byte {
	return byte(*address & 0x00FF)
}

func (address *Address) high() byte {
	return byte((*address & 0xFF00) >> 8)
}

func CreateAddress(low byte, high byte) Address {
	return Address(uint16(low) + uint16(high)<<8)
}

// RAM is
type RAM struct {
	memory [0xFFFF]byte
}

func (ram *RAM) read(address Address) byte {
	return ram.memory[address]
}

func (ram *RAM) write(address Address, data byte) {
	ram.memory[address] = data
}
