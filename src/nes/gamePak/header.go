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

const HorizontalMirroring = byte(1)
const VerticalMirroring = byte(2)
