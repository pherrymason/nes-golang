package cpu

func CreateCPU() *Cpu6502 {
	registers := CreateRegisters()
	cpu := Cpu6502{
		registers: registers,
		bus:       nil,
		debug:     false,
	}

	cpu.Init()

	return &cpu
}

func CreateCPUDebuggable(logger *Logger) *Cpu6502 {
	registers := CreateRegisters()
	cpu := Cpu6502{
		registers: registers,
		debug:     true,
		Logger:    *logger,
	}

	cpu.Init()

	return &cpu
}

// CreateCPUWithBus creates a Cpu6502 with a Bus, Useful for tests
func CreateCPUWithBus() *Cpu6502 {
	registers := CreateRegisters()

	//ram := component.RAM{}
	//gamePak := component.CreateDummyGamePak()

	cpu := Cpu6502{
		registers: registers,
		debug:     false,
	}

	//bus := component.CreateBus(&ram)
	//bus.InsertGamePak(&gamePak)
	//bus.ConnectCPU(&cpu)
	cpu.Init()

	return &cpu
}
