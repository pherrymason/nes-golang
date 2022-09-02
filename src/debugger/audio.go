package debugger

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/raulferras/nes-golang/src/audio"
)

type audioDebugger struct {
	audio *audio.Audio
}

func NewAudioDebugger(a *audio.Audio) *audioDebugger {
	return &audioDebugger{
		audio: a,
	}
}

func (dbg *audioDebugger) Draw() {
	positionX := int32(10)
	positionY := int32(520)
	rl.DrawRectangle(positionX, positionY, 1024, 100, rl.White)
	for i := int32(0); i < dbg.audio.AudioSample.SamplesCount; i++ {
		x := positionX + i
		y := int32(dbg.audio.AudioSample.Sample[i] * 20)
		rl.DrawPixel(x, positionY+y+50, rl.Red)
	}
}
