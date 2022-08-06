package nes

import (
	"github.com/FMNSSun/hexit"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
)

type Nes struct {
	cpu *Cpu6502
	ppu *ppu.Ppu2c02

	systemClockCounter byte // Controls how many times to call each processor
	debug              *NesDebugger
	vBlankCount        byte
	stopped            bool
}

func CreateNes(gamePak *gamePak.GamePak, debugger *NesDebugger) *Nes {
	hexit.BuildTable()
	thePPU := ppu.CreatePPU(gamePak, debugger.DebugPPU, debugger.logPath+"/ppu.log")

	cpuBus := newNESCPUMemory(thePPU, gamePak)
	cpu := CreateCPU(
		cpuBus,
		Cpu6502DebugOptions{debugger.debug, debugger.logPath + "/cpu.log"},
	)
	debugger.cpu = cpu
	debugger.ppu = thePPU

	nes := &Nes{
		cpu:   cpu,
		ppu:   thePPU,
		debug: debugger,
	}

	return nes
}

func (nes *Nes) StartAt(address types.Address) {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.ResetToAddress(address)
}

// Start todo Rename to PowerOn
func (nes *Nes) Start() {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.Reset()
}

func (nes *Nes) Tick() byte {
	nes.ppu.Tick()

	cpuCycles := byte(0)
	if nes.systemClockCounter%3 == 0 {
		if nes.cpu.memory.IsDMATransfer() {
			// DMA starts on an even cpu cycle
			if nes.cpu.memory.IsDMAWaiting() {
				if nes.systemClockCounter%2 == 1 {
					nes.cpu.memory.DisableDMWaiting()
				}
			} else {
				// On even cycles, read from RAM
				if nes.systemClockCounter%2 == 0 {
					address := uint16(nes.cpu.memory.GetDMAPage())<<8 | uint16(nes.cpu.memory.GetDMAAddress())
					nes.cpu.memory.SetDMAReadBuffer(nes.cpu.memory.Read(types.Address(address)))
				} else {
					nes.ppu.WriteRegister(ppu.OAMADDR, nes.cpu.memory.GetDMAAddress())
					nes.ppu.WriteRegister(ppu.OAMDATA, nes.cpu.memory.GetDMAReadBuffer())
					nes.cpu.memory.IncrementDMAAddress()
					if nes.cpu.memory.GetDMAAddress() == 0 {
						nes.cpu.memory.ResetDMA()
					}
				}
			}
		} else {
			cpuCycles = nes.cpu.Tick()
		}
	}

	if nes.ppu.Nmi() {
		nes.cpu.nmi()
		nes.ppu.ResetNmi()
	}

	if nes.ppu.VBlank() {
		if nes.vBlankCount == 20 {
			nes.ppu.Render()
		}
		nes.vBlankCount++
	}

	nes.systemClockCounter++
	if nes.debug.maxCPUCycle != -1 && nes.debug.maxCPUCycle <= int64(nes.cpu.cycle) {
		nes.Stop()
	}

	return cpuCycles
}

func (nes *Nes) TickForTime(seconds float64) {
	cycles := int(1789773 * seconds)
	for cycles > 0 {
		cycles -= int(nes.Tick())
		if nes.Stopped() {
			break
		}
	}
}

func (nes *Nes) Stop() {
	nes.cpu.Stop()
	nes.ppu.Stop()
	nes.stopped = true
}

func (nes *Nes) Stopped() bool {
	return nes.stopped
}

func (nes *Nes) Debugger() *NesDebugger {
	return nes.debug
}

func (nes *Nes) SystemClockCounter() byte {
	return nes.systemClockCounter
}

func (nes *Nes) Frame() *image.RGBA {
	return nes.ppu.Frame()
}
func (nes *Nes) FramePattern() []byte {
	return nes.ppu.FramePattern()
}

func (nes *Nes) PPU() *ppu.Ppu2c02 {
	return nes.ppu
}
