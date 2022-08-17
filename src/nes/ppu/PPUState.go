package ppu

type SimplePPUState struct {
	Frame       uint16
	RenderCycle uint16
	Scanline    Scanline
}

func NewSimplePPUState(frame uint16, renderCycle uint16, scanline Scanline) SimplePPUState {
	return SimplePPUState{
		Frame:       frame,
		RenderCycle: renderCycle,
		Scanline:    scanline,
	}
}
