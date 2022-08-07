package cpu

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

func createCPULogger(outputPath string) *cpu6502Logger {
	f, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", outputPath))
	}

	return &cpu6502Logger{
		file:       f,
		fileBuffer: bufio.NewWriterSize(f, CPU_LOG_MAX_SIZE*10),
		outputPath: outputPath,
		snapshots:  make([]CpuState, 0, CPU_LOG_MAX_SIZE),
	}
}

func (logger *cpu6502Logger) Log(state CpuState) {
	if len(logger.snapshots) == CPU_LOG_MAX_SIZE {
		logger.logToFile()
		logger.snapshots = logger.snapshots[:0]
	}

	logger.snapshots = append(logger.snapshots, state)
}

func (logger *cpu6502Logger) Close() {
	defer logger.file.Close()
	logger.logToFile()
	logger.file.Sync()
}

func (logger *cpu6502Logger) logToFile() {
	for _, state := range logger.snapshots {
		logger.fileBuffer.WriteString(state.String())
	}
	logger.fileBuffer.Flush()
	logger.file.Sync()
}

func (logger cpu6502Logger) Snapshots() []CpuState {
	return logger.snapshots
}
