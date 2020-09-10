package log

type Logger interface {
	Log(message string)
}

type MemoryLogger struct {
	document string
}

func (logger *MemoryLogger) Log(message string) {
	logger.document += message + "\n"
}
