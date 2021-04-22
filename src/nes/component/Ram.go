package component

import (
	"github.com/raulferras/nes-golang/src/nes/defs"
)

const RAM_LOWER_ADDRESS = defs.Address(0x0000)
const RAM_HIGHER_ADDRESS = defs.Address(0x1FFF)
const RAM_LAST_REAL_ADDRESS = defs.Address(0x07FF)

// RAM is $0000 -> $07FF
type RAM struct {
	memory [0xFFFF + 1]byte
}

func (ram *RAM) read(address defs.Address) byte {
	effectiveAddress := address & RAM_LAST_REAL_ADDRESS

	return ram.memory[effectiveAddress]
}

func (ram *RAM) write(address defs.Address, data byte) {
	effectiveAddress := address & RAM_LAST_REAL_ADDRESS

	ram.memory[effectiveAddress] = data
}
