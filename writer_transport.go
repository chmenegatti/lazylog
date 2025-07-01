package lazylog

import (
	"io"
)

type WriterTransport struct {
	Writer    io.Writer
	Level     Level
	Formatter Formatter
}

func (w *WriterTransport) WriteLog(entry *Entry) error {
	bytes, err := w.Formatter.Format(entry)
	if err != nil {
		// fallback simples
		_, err2 := w.Writer.Write([]byte(entry.Timestamp.Format("2006-01-02T15:04:05Z07:00") + " [" + entry.Level.String() + "] " + entry.Message + "\n"))
		return err2
	}
	_, err = w.Writer.Write(bytes)
	return err
}

func (w *WriterTransport) MinLevel() Level {
	return w.Level
}
