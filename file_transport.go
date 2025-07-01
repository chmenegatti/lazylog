package lazylog

import (
	"io"
	"os"
)

// FileTransport escreve logs em um arquivo específico.
type FileTransport struct {
	File      *os.File
	Level     Level
	Formatter Formatter
}

func NewFileTransport(path string, level Level, formatter Formatter) (*FileTransport, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FileTransport{
		File:      file,
		Level:     level,
		Formatter: formatter,
	}, nil
}

func (f *FileTransport) WriteLog(entry *Entry) error {
	bytes, err := f.Formatter.Format(entry)
	if err != nil {
		_, err2 := io.WriteString(f.File, entry.Timestamp.Format("2006-01-02T15:04:05Z07:00")+" ["+entry.Level.String()+"] "+entry.Message+"\n")
		return err2
	}
	_, err = f.File.Write(bytes)
	return err
}

func (f *FileTransport) MinLevel() Level {
	return f.Level
}

func (f *FileTransport) Close() error {
	return f.File.Close()
}
