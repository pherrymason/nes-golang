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

const HorizontalMirroring = byte(0)
const VerticalMirroring = byte(1)
const OneScreenMirroring = byte(2)
const FourScreenMirroring = byte(3)
