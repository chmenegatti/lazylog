package lazylog

import (
	"io"
	"time"
)

type Logger struct {
	Out      io.Writer
	MinLevel Level
}

func New() *Logger {
	return &Logger{
		Out:      io.Discard, // Default to discard output
		MinLevel: INFO,       // Default minimum level to INFO
	}
}

func (l *Logger) log(level Level, message string) {
	if level < l.MinLevel {
		return // Skip logging if the level is below the minimum
	}
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
	}
	io.WriteString(l.Out, entry.Timestamp.Format(time.RFC3339)+" ["+entry.Level.String()+"] "+entry.Message+"\n")
}

// Debug registra uma mensagem no nível DEBUG.
func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

// Info registra uma mensagem no nível INFO.
func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

// Warn registra uma mensagem no nível WARN.
func (l *Logger) Warn(message string) {
	l.log(WARN, message)
}

// Error registra uma mensagem no nível ERROR.
func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}
