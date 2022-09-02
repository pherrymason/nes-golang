package audio

import "math"

type Sample struct {
	Sample       []float32
	SamplesCount int32
}

func NewAudioSample() *Sample {
	as := Sample{
		Sample:       make([]float32, SamplesCount),
		SamplesCount: SamplesCount,
	}

	return &as
}

func (a *Sample) genDummy(osc *Oscillator) {
	for t := 0; t < SamplesCount; t++ {
		osc.phase = osc.phaseStride
		if osc.phase >= 1.0 {
			osc.phase -= 1.0
		}

		x := 2 * math.Pi * osc.phase * float32(t)
		a.Sample[t] = float32(math.Sin(float64(x)))
	}
}
