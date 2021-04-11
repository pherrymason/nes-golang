package log

import (
	"fmt"
	"os"
)

type Logger interface {
	Log(message string)
}

type MemoryLogger struct {
	document string
}

func (logger *MemoryLogger) Log(message string) {
	logger.document += message + "\n"
}

type FileLogger struct {
	outputPath string
}

func CreateFileLogger(outputPath string) FileLogger {
	f, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", outputPath))
	}
	defer f.Close()

	return FileLogger{outputPath: outputPath}
}

func (logger *FileLogger) Log(message string) {
	f, err := os.OpenFile(logger.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(message + "\n")
	f.Sync()
}
