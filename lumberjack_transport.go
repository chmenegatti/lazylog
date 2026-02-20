package lazylog

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

// LumberjackTransport faz rotação automática de arquivos de log.
type LumberjackTransport struct {
	Logger    *lumberjack.Logger
	Level     Level
	Formatter Formatter
}

func NewLumberjackTransport(filename string, level Level, formatter Formatter, maxSize, maxBackups, maxAge int, compress bool) *LumberjackTransport {
	return &LumberjackTransport{
		Logger: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		},
		Level:     level,
		Formatter: formatter,
	}
}

func (l *LumberjackTransport) WriteLog(entry *Entry) error {
	formatter := l.Formatter
	if formatter == nil {
		formatter = &TextFormatter{}
	}
	bytes, err := formatter.Format(entry)
	if err != nil {
		_, err2 := l.Logger.Write([]byte(entry.Timestamp.Format("2006-01-02T15:04:05Z07:00") + " [" + entry.Level.String() + "] " + entry.Message + "\n"))
		return err2
	}
	_, err = l.Logger.Write(bytes)
	return err
}

func (l *LumberjackTransport) MinLevel() Level {
	return l.Level
}

// Close fecha o logger lumberjack.
func (l *LumberjackTransport) Close() error {
	return l.Logger.Close()
}
