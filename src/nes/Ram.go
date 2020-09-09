package nes

// RAM is $0000 -> $07FF
type RAM struct {
	memory [0xFFFF + 1]byte
}

func (ram *RAM) read(address Address) byte {
	return ram.memory[address]
}

func (ram *RAM) write(address Address, data byte) {
	ram.memory[address] = data
}
