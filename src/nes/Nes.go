package nes

type Nes struct {
	cpu *Cpu6502
	ppu *Ppu2c02

	systemClockCounter byte // Controls how many times to call each processor
	debug              NesDebugger
}

func CreateNes(gamePak *GamePak, debugger NesDebugger) Nes {
	ppuBus := CreatePPUMemory(gamePak)
	ppu := CreatePPU(
		ppuBus,
	)

	cpuBus := CreateCPUMemory(ppu, gamePak)
	cpu := CreateCPU(
		cpuBus,
		Cpu6502DebugOptions{debugger.debug, debugger.outputLogPath},
	)
	debugger.cpu = cpu
	debugger.ppu = ppu

	nes := Nes{
		cpu:   cpu,
		ppu:   ppu,
		debug: debugger,
	}

	return nes
}

func (nes *Nes) StartAt(address Address) {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.ResetToAddress(address)
}

func (nes *Nes) Start() {
	nes.systemClockCounter = 0
	nes.debug.disassembled = nes.cpu.Disassemble(0x8000, 0xFFFF)
	nes.cpu.Reset()
}

func (nes *Nes) Tick() {
	//if nes.systemClockCounter%3 == 0 {
	nes.cpu.Tick()
	//}

	nes.ppu.Tick()
	nes.systemClockCounter++
}

func (nes Nes) Debugger() *NesDebugger {
	return &nes.debug
}

func (nes Nes) SystemClockCounter() byte {
	return nes.systemClockCounter
}
