package nes

import (
	"github.com/FMNSSun/hexit"
	cpu2 "github.com/raulferras/nes-golang/src/nes/cpu"
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"github.com/raulferras/nes-golang/src/nes/types"
	"image"
	"log"
	"runtime/debug"
)

type Nes struct {
	Cpu *Cpu6502
	ppu *ppu.P2c02
	bus *CPUMemory

	systemClockCounter uint64 // Controls how many times to call each processor
	debug              *Debugger
	vBlankCount        byte
	finished           bool
	paused             bool
}

func CreateNes(gamePak *gamePak.GamePak, debugger *Debugger) *Nes {
	hexit.BuildTable()
	thePPU := ppu.CreatePPU(gamePak, debugger.DebugPPU, debugger.logPath+"/ppu.log")

	cpuBus := newNESCPUMemory(thePPU, gamePak)
	cpu := CreateCPU(
		cpuBus,
		cpu2.NewDebugger(debugger.debugCPU, debugger.logPath+"/Cpu.log"),
	)
	debugger.cpu = cpu
	debugger.ppu = thePPU

	nes := &Nes{
		Cpu:   cpu,
		ppu:   thePPU,
		bus:   cpuBus,
		debug: debugger,
	}

	nes.debug.pauseEmulation = nes.Pause

	return nes
}

func (nes *Nes) StartAt(address types.Address) {
	nes.systemClockCounter = 0
	disassembledMap, sortedDisassembled := nes.Cpu.Disassemble(0x8000, 0xFFFF)
	nes.debug.disassembled = disassembledMap
	nes.debug.sortedDisassembled = sortedDisassembled

	nes.Cpu.ResetToAddress(address)
	// Run PPU for 7 cpu cycles. We subtract 1 because next Nes.Tick will first call PPU.
	for i := 0; i < (7 * 3); i++ {
		nes.ppu.Tick()
		nes.systemClockCounter++
	}
}

// Start todo Rename to PowerOn
func (nes *Nes) Start() {
	nes.systemClockCounter = 0
	disassembledMap, sortedDisassembled := nes.Cpu.Disassemble(0x8000, 0xFFFF)
	nes.debug.disassembled = disassembledMap
	nes.debug.sortedDisassembled = sortedDisassembled
	nes.Cpu.Reset()

	// Run PPU for 7 cpu cycles
	for i := 0; i < (7*3)-1; i++ {
		nes.ppu.Tick()
	}
}

func (nes *Nes) Pause() {
	nes.paused = true
}

func (nes *Nes) PausedTick() {
	if nes.Debugger().shouldPauseEmulation() {
		return
	}
	waitForOperationCompletion := false
	initialCompleteStatus := nes.Cpu.Complete()
	if initialCompleteStatus == false {
		waitForOperationCompletion = true
	}
	for {
		nes.Tick()
		if !waitForOperationCompletion {
			if nes.Cpu.Complete() != initialCompleteStatus {
				waitForOperationCompletion = true
			}
		} else {
			if nes.Cpu.Complete() {
				break
			}
		}
	}

	// Run until just before next cpu operation schedules to be called
	if (nes.systemClockCounter)%3 != 0 {
		for {
			nes.Tick()
			if (nes.systemClockCounter)%3 == 0 {
				break
			}
		}
	}

	nes.Debugger().oneCpuOperationRan()
}

func (nes *Nes) TickForTime(seconds float64) {
	cycles := int(1789773 * seconds)
	//waitingForCpuOperation := false
	for cycles > 0 {
		if nes.Cpu.Complete() && nes.Debugger().shouldPauseEmulation() {
			break
		}

		cpuCycles, _ := nes.Tick()
		cycles -= int(cpuCycles)
		if nes.Finished() {
			break
		}
	}
}

func (nes *Nes) TickTillFrameComplete() {
	for !nes.PPU().FrameComplete() {
		nes.Tick()
	}
}

func (nes *Nes) Tick() (byte, bool) {
	defer nes.handlePanic()

	var ppuState ppu.SimplePPUState
	if nes.Cpu.debugger.Enabled {
		ppuState = ppu.NewSimplePPUState(nes.ppu.FrameNumber(), nes.ppu.RenderCycle(), nes.ppu.Scanline())
	}
	//start := time.Now()
	nes.ppu.Tick()
	//elapsed := time.Since(start)
	//log.Printf("ppu took %s", elapsed)

	//start = time.Now()
	cpuCycles := byte(0)
	cpuExecuted := false
	if nes.systemClockCounter%3 == 0 {
		cpuExecuted = true
		if nes.Cpu.memory.IsDMATransfer() {
			// DMA starts on an even Cpu cycle
			if nes.Cpu.memory.IsDMAWaiting() {
				if nes.systemClockCounter%2 == 1 {
					nes.Cpu.memory.DisableDMWaiting()
				}
			} else {
				// On even cycles, read from RAM
				if nes.systemClockCounter%2 == 0 {
					address := uint16(nes.Cpu.memory.GetDMAPage())<<8 | uint16(nes.Cpu.memory.GetDMAAddress())
					nes.Cpu.memory.SetDMAReadBuffer(nes.Cpu.memory.Read(types.Address(address)))
				} else { // On odd cycles, write to OAMDATA
					nes.ppu.WriteRegister(ppu.OAMADDR, nes.Cpu.memory.GetDMAAddress())
					nes.ppu.WriteRegister(ppu.OAMDATA, nes.Cpu.memory.GetDMAReadBuffer())
					nes.Cpu.memory.IncrementDMAAddress()
					if nes.Cpu.memory.GetDMAAddress() == 0 {
						nes.Cpu.memory.ResetDMA()
					}
				}
			}
			cpuCycles = 1
		} else {
			spend, cpuState := nes.Cpu.Tick()
			cpuCycles = spend
			if nes.Cpu.debugger.Enabled {
				nes.Cpu.debugger.LogState(
					cpuState,
					ppuState,
				)
			}
		}
	}
	//elapsed = time.Since(start)
	//log.Printf("cpu took %s", elapsed)

	if nes.ppu.Nmi() {
		nes.Cpu.nmi()
		nes.ppu.ResetNmi()
	}

	nes.systemClockCounter++

	return cpuCycles, cpuExecuted
}

func (nes *Nes) Stop() {
	nes.Cpu.Stop()
	nes.ppu.Stop()
	nes.finished = true
}

func (nes *Nes) Finished() bool {
	return nes.finished
}

func (nes *Nes) Paused() bool {
	return nes.paused
}

func (nes *Nes) Debugger() *Debugger {
	return nes.debug
}

func (nes *Nes) SystemClockCounter() uint64 {
	return nes.systemClockCounter
}

func (nes *Nes) Frame() *image.RGBA {
	return nes.ppu.Frame()
}
func (nes *Nes) FramePattern() []byte {
	return nes.ppu.FramePattern()
}

func (nes *Nes) PPU() *ppu.P2c02 {
	return nes.ppu
}

func (nes *Nes) handlePanic() {
	a := recover()
	if a != nil {
		nes.Stop()
		log.Fatalf("%s\nTrace: %s", a, string(debug.Stack()))
		//os.Exit(3)
	}
}

func (nes *Nes) UpdateController(controllerNumber int, state ControllerState) {
	nes.bus.controllers[controllerNumber-1] = state.value()
}
