package ppu

type SimplePPUState struct {
	Frame       uint
	RenderCycle uint16
	Scanline    int16
}

func NewSimplePPUState(frame uint, renderCycle uint16, scanline int16) SimplePPUState {
	return SimplePPUState{
		Frame:       frame,
		RenderCycle: renderCycle,
		Scanline:    scanline,
	}
}
