package nes

import (
	"fmt"
	"os"
)

type cpu6502Logger struct {
	outputPath string
	snapshots  []CpuState
}

func createCPULogger(outputPath string) cpu6502Logger {
	return cpu6502Logger{outputPath: outputPath}
}

func (logger *cpu6502Logger) Log(state CpuState) {
	logger.snapshots = append(logger.snapshots, state)
}

func (logger *cpu6502Logger) Close() {
	f, err := os.Create(logger.outputPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", logger.outputPath))
	}
	defer f.Close()

	//f, err := os.OpenFile(logger.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()

	for _, state := range logger.snapshots {
		f.WriteString(stateToString(state) + "\n")
	}
	f.Sync()
}

func (logger cpu6502Logger) Snapshots() []CpuState {
	return logger.snapshots
}

func stateToString(state CpuState) string {
	var msg string
	// Pointer
	msg += fmt.Sprintf("%02X", state.Registers.Pc) + "  "

	// Raw OPCode + Operand
	for _, value := range state.RawOpcode {
		msg += fmt.Sprintf("%02X ", value)
	}

	msg = clampSpace(msg, 16)

	msg += state.CurrentInstruction.Name() + " "

	if state.CurrentInstruction.AddressMode() == Immediate {
		msg += fmt.Sprintf("#$%02X", state.EvaluatedAddress)
	} else {
		msg += fmt.Sprintf("$%02X", state.EvaluatedAddress)
	}

	msg = clampSpace(msg, 48)

	msg += fmt.Sprintf(
		"A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:___,___ CYC:%d",
		state.Registers.A,
		state.Registers.X,
		state.Registers.Y,
		state.Registers.Status,
		state.Registers.Sp,
		state.CyclesSinceReset,
	)

	return msg
}

func clampSpace(msg string, clamp int) string {
	for i := len(msg); i <= clamp; i++ {
		msg += " "
	}
	return msg
}
