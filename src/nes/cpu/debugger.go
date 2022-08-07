package cpu

// Debugger
// Handles logging and debugging features like breakpoints
type Debugger struct {
	Enabled       bool
	OutputLogPath string
	Logger        *cpu6502Logger
}

func NewDebugger(enabled bool, logPath string) *Debugger {
	var logger *cpu6502Logger
	if enabled {
		logger = createCPULogger(logPath)
	} else {
		logger = nil
	}

	return &Debugger{
		Enabled:       enabled,
		OutputLogPath: logPath,
		Logger:        logger,
	}
}

func (debugger *Debugger) Stop() {
	if debugger.Enabled {
		debugger.Logger.Close()
	}
}

func (debugger *Debugger) LogStep(registers Registers, opcode byte, operand [3]byte, instruction Instruction, step OperationMethodArgument, cpuCycle uint32) {
	//state := CreateStateFromCPU(*debugger)
	state := CreateState(
		registers,
		[3]byte{opcode, operand[0], operand[1]},
		instruction,
		step,
		cpuCycle,
	)

	debugger.Logger.Log(state)
}
