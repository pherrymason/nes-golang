package cpu

import (
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/defs"
	"os"
)

type Logger struct {
	outputPath string
	snapshots  []State
}

func CreateCPULogger(outputPath string) Logger {
	f, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", outputPath))
	}
	defer f.Close()

	return Logger{outputPath: outputPath}
}

func (logger *Logger) Log(state State) {
	logger.snapshots = append(logger.snapshots, state)

	f, err := os.OpenFile(logger.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	message := stateToString(state)
	f.WriteString(message + "\n")
	f.Sync()
}

func (logger Logger) Snapshots() []State {
	return logger.snapshots
}

func stateToString(state State) string {
	var msg string
	// Pointer
	msg += fmt.Sprintf("%02X", state.Registers.Pc) + "  "

	// Raw OPCode + Operand
	for _, value := range state.RawOpcode {
		msg += fmt.Sprintf("%02X ", value)
	}

	msg = clampSpace(msg, 16)

	msg += state.CurrentInstruction.Name() + " "

	if state.CurrentInstruction.AddressMode() == defs.Immediate {
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
