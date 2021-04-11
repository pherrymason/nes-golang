package component

import (
	"github.com/raulferras/nes-golang/src/nes/defs"
)

// RAM is $0000 -> $07FF
type RAM struct {
	memory [0xFFFF + 1]byte
}

func (ram *RAM) read(address defs.Address) byte {
	return ram.memory[address]
}

func (ram *RAM) write(address defs.Address, data byte) {
	ram.memory[address] = data
}
