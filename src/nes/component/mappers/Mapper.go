package mappers

import (
	"github.com/raulferras/nes-golang/src/nes"
)

type Mapper interface {
	readCpu(nes.Address) byte
	writeCpu(nes.Address, byte)

	readPpu(nes.Address) byte
	writePpu(nes.Address, byte)
}
