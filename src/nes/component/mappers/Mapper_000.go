package mappers

import (
	"github.com/raulferras/nes-golang/src/nes"
)

/**
	NROM Mapper

	Banks
	CPU $6000-$7FFF:
		Family Basic only: PRG RAM, mirrored as necessary to fill entire 8 KiB window,
		write protectable with an external switch

	if PRGROM is 16KB
	CPU $8000-$BFFF: First 16 KB of ROM.
	CPU $C000-$FFFF: Mirrors first 16 KB of ROM.

	if PRGROM is 32KB
    CPU $8000-$BFFF: First 16 KB of ROM.
	CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256)
*/

// Mapper000 NROM Implementation
type Mapper000 struct {
	prgBanks int
	chrBanks int
}

func (mapper *Mapper000) readCpu(address nes.Address) byte {
	panic("implement me")
}

func (mapper *Mapper000) writeCpu(address nes.Address, b byte) {
	panic("implement me")
}

func (mapper *Mapper000) readPpu(address nes.Address) byte {
	panic("implement me")
}

func (mapper *Mapper000) writePpu(address nes.Address, b byte) {
	panic("implement me")
}
