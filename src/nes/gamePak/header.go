package gamePak

type Header interface {
	ProgramSize() byte
	CHRSize() byte
	Mirroring() byte
	HasPersistentMemory() bool
	HasTrainer() bool
	IgnoreMirroringControl() bool

	MapperNumber() byte
	PRGRAM() byte
	TvSystem() byte
}

const HorizontalMirroring = byte(0b00)
const VerticalMirroring = byte(0b01)
const OneScreenMirroring = byte(0b10)
const FourScreenMirroring = byte(0b11)
