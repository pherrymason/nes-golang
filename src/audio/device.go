package audio

import (
	"encoding/binary"
	"fmt"
	r "github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"os"
)

const SamplesCount = 1024

type Oscillator struct {
	phase       float32
	phaseStride float32
}

type Audio struct {
	sampleRate  float32
	audioStream r.AudioStream
	AudioSample *Sample
	frequency   float32
	osc         Oscillator
}

func NewAudio() *Audio {
	return &Audio{
		sampleRate:  44100,
		AudioSample: NewAudioSample(),
	}
}

func (a *Audio) Init() {
	log.Println("Init audio")
	r.InitAudioDevice()
	audioStream := r.LoadAudioStream(
		uint32(a.sampleRate),
		32,
		2,
	)

	r.SetAudioStreamVolume(audioStream, 0.10)
	r.SetAudioStreamBufferSizeDefault(SamplesCount)
	r.PlayAudioStream(audioStream)

	a.audioStream = audioStream
	a.frequency = 400
	sampleDuration := 1 / a.sampleRate
	a.osc.phase = 0
	a.osc.phaseStride = a.frequency * sampleDuration
}

func (a *Audio) Stop() {
	r.UnloadAudioStream(a.audioStream)
}

func (a *Audio) Update() {
	if !r.IsAudioStreamPlaying(a.audioStream) {

	}
	if r.IsAudioStreamProcessed(a.audioStream) {
		a.AudioSample.genDummy(&a.osc)
		r.UpdateAudioStream(
			a.audioStream,
			a.AudioSample.Sample,
			a.AudioSample.SamplesCount,
		)
		a.frequency += 1
		a.osc.phaseStride = a.frequency * (1 / a.sampleRate)
	}
}

const (
	Duration   = 2
	SampleRate = 44100
	Frequency  = 440 // Pitch Standard
)

func Generate() {
	nsamps := Duration * SampleRate
	tau := math.Pi * 2
	var angle float64 = tau / float64(nsamps)
	file := "out.bin"
	f, _ := os.Create(file)
	for i := 0; i < nsamps; i++ {
		sample := math.Sin(angle * Frequency * float64(i))
		var buf [8]byte
		binary.LittleEndian.PutUint32(buf[:],
			math.Float32bits(float32(sample)))
		bw, _ := f.Write(buf[:])
		fmt.Printf("\rWrote: %v bytes to %s", bw, file)
	}
}
