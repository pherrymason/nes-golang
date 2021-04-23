package mappers

import "github.com/raulferras/nes-golang/src/nes/defs"

type Mapper interface {
	readCpu(defs.Address) byte
	writeCpu(defs.Address, byte)

	readPpu(defs.Address) byte
	writePpu(defs.Address, byte)
}
