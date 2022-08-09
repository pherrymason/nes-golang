package cpu

import (
	"bufio"
	"fmt"
	"github.com/raulferras/nes-golang/src/nes/ppu"
	"os"
)

type cpu6502Logger struct {
	file       *os.File
	fileBuffer *bufio.Writer
	outputPath string
	snapshots  []Snapshot
}

type Snapshot struct {
	CpuState CpuState
	PpuState ppu.SimplePPUState
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
		snapshots:  make([]Snapshot, 0, CPU_LOG_MAX_SIZE),
	}
}

func (logger *cpu6502Logger) Log(state CpuState, ppuState ppu.SimplePPUState) {
	if len(logger.snapshots) == CPU_LOG_MAX_SIZE {
		logger.logToFile()
		logger.snapshots = logger.snapshots[:0]
	}

	logger.snapshots = append(logger.snapshots, Snapshot{state, ppuState})
}

func (logger *cpu6502Logger) Close() {
	defer logger.file.Close()
	logger.logToFile()
	logger.file.Sync()
}

func (logger *cpu6502Logger) logToFile() {
	for _, snapshot := range logger.snapshots {
		logger.fileBuffer.WriteString(snapshot.CpuState.String(snapshot.PpuState))
	}
	logger.fileBuffer.Flush()
	logger.file.Sync()
}

func (logger cpu6502Logger) Snapshots() []Snapshot {
	return logger.snapshots
}
