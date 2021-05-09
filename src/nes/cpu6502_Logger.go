package nes

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/raulferras/nes-golang/src/utils"
	"os"
	"strconv"
	"strings"
)

type cpu6502Logger struct {
	file       *os.File
	fileBuffer *bufio.Writer
	outputPath string
	snapshots  []CpuState
}

func createCPULogger(outputPath string) cpu6502Logger {
	f, err := os.Create(outputPath)
	writer := bufio.NewWriter(f)
	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", outputPath))
	}

	return cpu6502Logger{
		file:       f,
		fileBuffer: writer,
		outputPath: outputPath,
		snapshots:  make([]CpuState, 0, 30024),
	}
}

func (logger *cpu6502Logger) Log(state CpuState) {
	if len(logger.snapshots) == 30000 {
		for _, state := range logger.snapshots {
			fmt.Fprintf(logger.fileBuffer, stateToString(state)+"\n")
		}
		//logger.file.Sync()
		logger.snapshots = make([]CpuState, 0, 300024)
	}

	logger.snapshots = append(logger.snapshots, state)
}

func (logger *cpu6502Logger) Close() {
	defer logger.file.Close()

	for _, state := range logger.snapshots {
		fmt.Fprintf(logger.fileBuffer, stateToString(state)+"\n")
	}
	logger.file.Sync()
}

func (logger cpu6502Logger) Snapshots() []CpuState {
	return logger.snapshots
}

func stateToString(state CpuState) string {
	var msg strings.Builder
	msg.Grow(75)

	// Pointer
	pcBytes := []byte{state.Registers.Pc.HighNibble(), state.Registers.Pc.LowNibble()}
	msg.WriteString(hex.EncodeToString(pcBytes) + " ")

	// Raw OPCode + Operand
	msg.WriteString(hex.EncodeToString(state.RawOpcode))
	/*for _, value := range state.RawOpcode {
		msg.WriteString("0x" + hex.EncodeToString()fmt.Sprintf("%02X ", value))
	}*/

	clampSpace(&msg, 16)
	msg.WriteString(state.CurrentInstruction.Name() + " ")

	if state.CurrentInstruction.AddressMode() == Immediate {
		msg.WriteString("#$%02X" + hex.EncodeToString(state.EvaluatedAddress.ToBytes()))
	} else {
		msg.WriteString(
			"$%02X" + hex.EncodeToString(state.EvaluatedAddress.ToBytes()))
	}

	//msg = clampSpace(msg, 48)
	msg.WriteByte(' ')
	msg.WriteString("A:" + utils.ByteToHex(state.Registers.A))
	msg.WriteString("X:" + utils.ByteToHex(state.Registers.X))
	msg.WriteString("Y:" + utils.ByteToHex(state.Registers.Y))
	msg.WriteString("P:" + utils.ByteToHex(state.Registers.Status))
	msg.WriteString("SP:" + utils.ByteToHex(state.Registers.Sp))
	msg.WriteString("PPU:___,___")
	msg.WriteString("CYC:" + strconv.Itoa(int(state.Registers.Sp)))

	/*
		msg.WriteString(
			fmt.Sprintf(
				"A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:___,___ CYC:%d",
				state.Registers.A,
				state.Registers.X,
				state.Registers.Y,
				state.Registers.Status,
				state.Registers.Sp,
				state.CyclesSinceReset,
			))
	*/
	return msg.String()
}

func clampSpace(msg *strings.Builder, clamp int) {
	for i := msg.Len(); i <= clamp; i++ {
		msg.WriteString(" ")
	}
}
