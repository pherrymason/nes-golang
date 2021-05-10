package nes

import (
	"bufio"
	"fmt"
	"os"
)

type cpu6502Logger struct {
	file       *os.File
	fileBuffer *bufio.Writer
	outputPath string
	snapshots  []CpuState
}

const CPU_LOG_MAX_SIZE = 120000

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
		snapshots:  make([]CpuState, 0, CPU_LOG_MAX_SIZE),
	}
}

func (logger *cpu6502Logger) Log(state CpuState) {
	if len(logger.snapshots) == CPU_LOG_MAX_SIZE {
		for _, state := range logger.snapshots {
			fmt.Fprintf(logger.fileBuffer, state.String())
		}
		logger.snapshots = logger.snapshots[:0]
	}

	logger.snapshots = append(logger.snapshots, state)
}

func (logger *cpu6502Logger) Close() {
	defer logger.file.Close()

	for _, state := range logger.snapshots {
		fmt.Fprintf(logger.fileBuffer, state.String())
	}
	logger.file.Sync()
}

func (logger cpu6502Logger) Snapshots() []CpuState {
	return logger.snapshots
}
