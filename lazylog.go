package lazylog

import (
	"io"
	"time"
)

type Logger struct {
	Out       io.Writer
	MinLevel  Level
	Formatter Formatter // Optional formatter for log entries
}

func New() *Logger {
	return &Logger{
		Out:      io.Discard, // Default to discard output
		MinLevel: INFO,       // Default minimum level to INFO
		Formatter: &TextFormatter{
			TimestampFormat: time.RFC3339, // Use default timestamp format (time.RFC3339)
		}, // No formatter by default
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
	bytes, err := l.Formatter.Format(&entry)
	if err != nil {
		// If formatting fails, we log to the default output without formatting
		io.WriteString(l.Out, entry.Timestamp.Format(time.RFC3339)+" ["+level.String()+"] "+message+"\n")
		return
	}

	l.Out.Write(bytes) // Write the formatted log entry to the output
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
