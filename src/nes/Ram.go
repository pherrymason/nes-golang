package nes

// RAM is $0000 -> $07FF
type RAM struct {
	memory [0xFFFF + 1]byte
}

func (ram *RAM) read(address Address) byte {
	return ram.memory[address]
}

func (ram *RAM) read16(address Address) Word {
	low := ram.read(address)
	high := ram.read(address + 1)

	return CreateWord(low, high)
}

/*
// Emulates page boundary hardware bug
func (ram *RAM) read16Bugged(address Address) Word {
	lsb := address
	msb := (lsb & 0xFF00) | (lsb & 0xFF) + 1

	flsb := ram.read(lsb)
	fmsb := ram.read(msb)

	return CreateWord(flsb, fmsb)
}*/

func (ram *RAM) write(address Address, data byte) {
	ram.memory[address] = data
}
