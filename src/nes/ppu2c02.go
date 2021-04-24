package nes

// Ppu2c02 Processor
//  Registers mapped to memory locations: $2000 through $2007
//  mirrored in every 8 bytes from $2008 through $3FFF
type Ppu2c02 struct {
	Memory
	patternTable []byte // Decoded pattern table
}

func CreatePPU(memory Memory) *Ppu2c02 {
	ppu := &Ppu2c02{
		Memory:       memory,
		patternTable: make([]byte, 8*8*512),
	}

	return ppu
}

func (ppu *Ppu2c02) Tick() {

}

func (ppu *Ppu2c02) PatternTable() []byte {
	return ppu.patternTable
}
